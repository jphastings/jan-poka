#define ARDUINOJSON_USE_DOUBLE 1

#include <Arduino.h>
#include <AccelStepper.h>
#include <NeoPixelBus.h>
#include <PubSubClient.h>
#include <ArduinoJson.h>
#include <HTTPClient.h>
#include <WiFi.h>
#include <time.h>
#include <sys/time.h>
#include <esp_sntp.h>
#include "vars.h"

#define APP_NAME "jan-poka:clock"
#define TIMEZONE_QUERY "http://api.geonames.org/timezoneJSON?lat=%f&lng=%f&username=jphastings"
#define NTP_SERVER "time.google.com"
#define MOTOR_STEPS 720
#define MINS_STEPS MOTOR_STEPS / 60.0
#define HOURS_STEPS MOTOR_STEPS / 12.0
#define LED_COUNT 96
#define LED_FADE_AT LED_COUNT/3
#define LED_OFF_AT LED_COUNT*2/3
#define LED_MINS 12 * 60 / (float) LED_COUNT
#define UPDATE_FREQ_S 10

// Pins
#define LED_DATA_PIN 27
#define HALL_QANTIZATION 64

#define HOURS_STEPPER_A1 16
#define HOURS_STEPPER_A2 17
#define HOURS_STEPPER_B1 18
#define HOURS_STEPPER_B2 15
#define HOURS_HALL_SENSOR_PIN 34
#define HOURS_HALL_SENSOR_POS 360
#define HOURS_STEP_SPEED 200
#define HOURS_CALIB_STEP_SPEED 85

#define MINS_STEPPER_A1 19
#define MINS_STEPPER_A2 21
#define MINS_STEPPER_B1 23
#define MINS_STEPPER_B2 22
#define MINS_HALL_SENSOR_PIN 35
#define MINS_HALL_SENSOR_POS 356
#define MINS_STEP_SPEED 200
#define MINS_CALIB_STEP_SPEED 110

AccelStepper stepHour(AccelStepper::FULL4WIRE, HOURS_STEPPER_A1, HOURS_STEPPER_A2, HOURS_STEPPER_B1, HOURS_STEPPER_B2);
AccelStepper stepMin(AccelStepper::FULL4WIRE, MINS_STEPPER_A1, MINS_STEPPER_A2, MINS_STEPPER_B1, MINS_STEPPER_B2);

NeoPixelBus<NeoGrbwFeature, NeoSk6812Method> leds(LED_COUNT, LED_DATA_PIN);

WiFiClient wifiClient;
PubSubClient mqttClient(wifiClient);
time_t epoch;
struct tm *now;
int lastUpdated = -1;

int sunriseMins = -1;
int sunsetMins = -1;

// TODO: Pick good colours; include Warm White?
RgbwColor dayCol(128, 64, 0, 128);
RgbwColor nightCol(0, 0, 128, 0);
RgbwColor offCol(0, 0, 0, 0);

struct StepperCalibration {
  long cumulativeReadings[MOTOR_STEPS];
  int countReadings[MOTOR_STEPS];
  bool calibrated;
  bool isMins;
};
struct StepperCalibration minsCalibrate;
struct StepperCalibration hoursCalibrate;
bool motorsCalibrating;

void setupWifi();
void setupMQTT();
void setupTime();
void setupMotors();
void setupLEDs();

void updateTime();
void updateClock();
void updateMotors();
void updateLEDs();
void updateMQTT();

void handleGeoTarget(char*, byte*, unsigned int);
bool steppersMoving();
void calibrateMotors();
void setSunriseAndSunset(String, String);
int timeToDayMins(int, int);
void startMotorCalibration();

int updatedInt();
bool noNeedToUpdate();
void setTimezone(double, double);
void normalizeStepper(AccelStepper*);
void moveCircular(AccelStepper*, long);

int normalizeStepCount(int pos) {
  if (pos < 0) {
      pos += -(pos/MOTOR_STEPS)*MOTOR_STEPS + MOTOR_STEPS;
  }
  return pos % MOTOR_STEPS;
}

void setupWifi() {
  WiFi.begin(WIFI_SSID, WIFI_PASS);
  while(WiFi.status()!=WL_CONNECTED) {
    delay(1);
  }
}

void updateMQTT() {
  if (!mqttClient.connected()) {
    if (!mqttClient.connect(APP_NAME, MQTT_USER, MQTT_PASS)) {
      Serial.println("Failed to connect to MQTT broker");
      // TODO: Backoff
      return;
    }
    mqttClient.setCallback(handleGeoTarget);
    mqttClient.subscribe(MQTT_TOPIC);
  }
  mqttClient.loop();
}

void setupMQTT() {
  mqttClient.setBufferSize(1024);
  mqttClient.setServer(MQTT_HOST, MQTT_PORT);
  updateMQTT();
}

void setupTime() {
  sntp_setoperatingmode(SNTP_OPMODE_POLL);
  sntp_setservername(0, NTP_SERVER);
  sntp_init();
}

void setupMotors() {
  stepMin.setMaxSpeed(MINS_CALIB_STEP_SPEED);
  stepMin.setAcceleration(200.0);
  stepMin.setCurrentPosition(0);

  stepHour.setMaxSpeed(HOURS_CALIB_STEP_SPEED);
  stepHour.setAcceleration(200.0);
  stepHour.setCurrentPosition(0);

  startMotorCalibration();
}

void setupLEDs() {
  leds.Begin();
  leds.ClearTo(offCol);
  leds.Show();
}

void handleGeoTarget(char* topic, byte* payload, unsigned int length) {
  Serial.println("MQTT message received");
  /* Copied from EspMQTTClient: https://github.com/plapointe6/EspMQTTClient/blob/master/src/EspMQTTClient.cpp#L649 */
  // Convert the payload into a String
  // First, We ensure that we dont bypass the maximum size of the PubSubClient library buffer that originated the payload
  // This buffer has a maximum length of _mqttClient.getBufferSize() and the payload begin at "headerSize + topicLength + 1"
  unsigned int strTerminationPos;
  if (strlen(topic) + length + 9 >= mqttClient.getBufferSize()) {
    strTerminationPos = length - 1;
  } else {
    strTerminationPos = length;
  }
  
  // Second, we add the string termination code at the end of the payload
  payload[strTerminationPos] = '\0';
  /* end */

  StaticJsonDocument<1024> doc;
  DeserializationError err = deserializeJson(doc, payload, strTerminationPos + 1);
  if (err != DeserializationError::Ok) {
    Serial.print("Unsuccessful at parsing MQTT JSON: ");
    Serial.println(err.c_str());
    return;
  }
  setTimezone(doc["lat"], doc["lng"]);
}

bool steppersMoving() {
  return stepMin.distanceToGo() != 0 || stepHour.distanceToGo() != 0;
}

void setSunriseAndSunset(String sunriseStr, String sunsetStr) {
  sunriseMins = timeToDayMins(sunriseStr.substring(11, 13).toInt(), sunriseStr.substring(14, 16).toInt());
  sunsetMins = timeToDayMins(sunsetStr.substring(11, 13).toInt(), sunsetStr.substring(14, 16).toInt());
}

int timeToDayMins(int hours, int mins) {
  return (hours * 60 + mins);
}

void updateClock() {
  if (now == 0 || sunsetMins == -1 || sunriseMins == -1 || steppersMoving() || noNeedToUpdate())
    return;

  Serial.print("Setting time to: "); Serial.print(now->tm_hour); Serial.print(":"); Serial.print(now->tm_min); Serial.print(":"); Serial.println(now->tm_sec);

  updateMotors();
  updateLEDs();
    
  lastUpdated = updatedInt();
}

int updatedInt() {
  return (now->tm_min * 60 + now->tm_sec) / UPDATE_FREQ_S;
}

bool noNeedToUpdate() {
  return motorsCalibrating || lastUpdated == updatedInt();
}

void updateMotors() {
  moveCircular(&stepMin, (now->tm_min + now->tm_sec/60.0) * MINS_STEPS);
  moveCircular(&stepHour, (now->tm_hour + now->tm_min/60.0) * HOURS_STEPS);
}

void updateLEDs() {
  int nowMins = timeToDayMins(now->tm_hour, now->tm_min);
  // The time in minutes now, scaled to fit into the number of LEDs, offset by half as 0th pixel is at 6 o'clock
  int nowPos = (nowMins * LED_COUNT / (12 * 60) + (LED_COUNT / 2)) % LED_COUNT;

  leds.ClearTo(offCol);
  for (int i = 0; i < LED_OFF_AT; i++) {
    int iMins = (int)(i*LED_MINS + nowMins) % (24 * 60);
    bool nightTime = iMins < sunriseMins || iMins >= sunsetMins;
    RgbwColor col = nightTime ? nightCol : dayCol;

    float fadeAmount = (i < LED_FADE_AT) ? 0 : ((i - LED_FADE_AT) / (float)(LED_OFF_AT - LED_FADE_AT));
    col = RgbwColor::LinearBlend(col, offCol, fadeAmount);
    
    int pos = (i + nowPos) % LED_COUNT;
    leds.SetPixelColor(pos, col);
  }
  leds.Show();
}

void moveCircular(AccelStepper* stepper, long steps) {
    // normalizeStepper(stepper);
    // Move Backwards if it's faster
    // long curPos = -stepper->currentPosition(); // -ve for clockwise correction
    // long newPos = steps % MOTOR_STEPS;
    // if (newPos - curPos > MOTOR_STEPS/2) {
    //   newPos = MOTOR_STEPS - newPos;
    // }
    // else if (newPos - curPos + MOTOR_STEPS < MOTOR_STEPS/2) {
    //   newPos += MOTOR_STEPS;
    // }
    
    stepper->moveTo(normalizeStepCount(-steps)); // -ve for clockwise correction
}

void setTimezone(double latitude, double longitude) {
  Serial.print("Setting timezone for: ");
    Serial.print(latitude);
    Serial.print(", ");
    Serial.println(longitude);

  char url[255];
  sprintf(url, TIMEZONE_QUERY, latitude, longitude);
  
  HTTPClient httpClient;
  httpClient.useHTTP10(true);
  httpClient.begin(wifiClient, url);
  httpClient.GET();
  
  StaticJsonDocument<512> doc;
  DeserializationError err = deserializeJson(doc, httpClient.getStream());
  if (err != DeserializationError::Ok) {
    Serial.print("Grabbing timezone details failed: ");
    Serial.println(err.c_str());
    return;
  }

  // TODO: Deal with DST
  configTime(doc["rawOffset"].as<float>() * 3600, 0, NTP_SERVER);
  setSunriseAndSunset(doc["sunrise"], doc["sunset"]);

  httpClient.end();
  updateClock();
}

void updateTime() {
  time(&epoch);
  now = localtime(&epoch);
}

int guessHomePosition(long cumulativeReadings[MOTOR_STEPS], int countReadings[MOTOR_STEPS]) {
  int bestGuessPosFirst = 0;
  int bestGuessPosLast = 0;
  int lowestReading = 10000;

  for (int pos = 0; pos < MOTOR_STEPS; pos++) {
    if (countReadings[pos] == 0)
      continue;

    int val = (cumulativeReadings[pos] / countReadings[pos]) / HALL_QANTIZATION;
    // Serial.print(calibration->isMins); Serial.print(","); Serial.print(pos); Serial.print(","); Serial.println(val);
    if (val < lowestReading) {
      lowestReading = val;
      bestGuessPosFirst = pos;
      bestGuessPosLast = pos;
    } else if (val == lowestReading) {
      bestGuessPosLast = pos;
    }
  }

  return (bestGuessPosFirst + bestGuessPosLast) / 2;
}

void startMotorCalibration() {
  minsCalibrate = {.cumulativeReadings = {}, .countReadings = {}, .calibrated = false, .isMins = true};
  hoursCalibrate = {.cumulativeReadings = {}, .countReadings = {}, .calibrated = false, .isMins = false};
  motorsCalibrating = true;

  stepMin.moveTo(-MOTOR_STEPS*2);
  stepHour.moveTo(MOTOR_STEPS*2);
}

bool calibrateMotorStep(AccelStepper* stepper, struct StepperCalibration* calibration, int pin, int sensorPos) {
  if (calibration->calibrated)
    return true;

  int pos = normalizeStepCount(stepper->currentPosition());
  int reading = analogRead(pin);
  calibration->cumulativeReadings[pos] += reading;
  calibration->countReadings[pos]++;

  if (stepper->isRunning())
    return false;

  int bestGuessPos = guessHomePosition(calibration->cumulativeReadings, calibration->countReadings);

  int actualPos = normalizeStepCount(pos - bestGuessPos + sensorPos);

  stepper->setCurrentPosition(actualPos);
  stepper->moveTo(0);
  calibration->calibrated = true;
  return true;
}

void calibrateMotors() {
  bool minsGood = calibrateMotorStep(&stepMin, &minsCalibrate, MINS_HALL_SENSOR_PIN, MINS_HALL_SENSOR_POS);
  bool hoursGood = calibrateMotorStep(&stepHour, &hoursCalibrate, HOURS_HALL_SENSOR_PIN, HOURS_HALL_SENSOR_POS);
  if (minsGood && hoursGood) {
    Serial.println("Finished calibrating motors.");
    stepMin.setMaxSpeed(MINS_STEP_SPEED);
    stepHour.setMaxSpeed(HOURS_STEP_SPEED);
    motorsCalibrating = false;
  }
};

void setup() {
  Serial.begin(115200);
  
  setupWifi();
  setupMotors();
  setupLEDs();
  setupTime();
  setupMQTT();
  Serial.println("Booted");
}

void loop() {
  updateTime();
  
  stepMin.run();
  stepHour.run();

  if (motorsCalibrating)
    calibrateMotors();
  
  updateClock();
  updateMQTT();
}
