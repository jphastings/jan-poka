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

AccelStepper stepHour(AccelStepper::FULL4WIRE, 16, 17, 18, 15);
AccelStepper stepMin(AccelStepper::FULL4WIRE, 19, 21, 23, 22);

NeoPixelBus<NeoGrbwFeature, NeoSk6812Method> leds(LED_COUNT, 27);

HTTPClient httpClient;
WiFiClient wifiClient;
PubSubClient mqttClient(wifiClient);
time_t epoch;
struct tm *now;
int lastUpdated = -1;

int sunriseMins = -1;
int sunsetMins = -1;

void setupWifi() {
  WiFi.begin(WIFI_SSID, WIFI_PASS);
  while(WiFi.status()!=WL_CONNECTED) {
    delay(1);
  }
}

void setupMQTT() {
  mqttClient.setServer(MQTT_HOST, MQTT_PORT);
  mqttConnectionLoop();
}

void setupTime() {
  sntp_setoperatingmode(SNTP_OPMODE_POLL);
  sntp_setservername(0, NTP_SERVER);
  sntp_init();

  // Home:
  setTimezone(51.53647542276452, -0.08639983104800102);
  // Chris:
  //setTimezone(49.27103836424926, -123.15245006680813);
  // Venezuela
  //setTimezone(10.436620855606654, -66.84207546938133);
}

void setupMotors() {
  stepMin.setMaxSpeed(200.0);
  stepMin.setAcceleration(200.0);
  stepMin.setCurrentPosition(0);

  stepHour.setMaxSpeed(200.0);
  stepHour.setAcceleration(200.0);
  stepHour.setCurrentPosition(0);
}

void setupLEDs() {
  leds.Begin();
}

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
  time(&epoch);
  now = localtime(&epoch);
  
  stepMin.run();
  stepHour.run();

  updateClock();
}

void mqttConnectionLoop() {
  if (!mqttClient.connected()) {
    if (!mqttClient.connect(APP_NAME, MQTT_USER, MQTT_PASS)) {
      Serial.println("Failed to connect to MQTT broker");
      // TODO: Backoff
      return;
    }
    Serial.println("MQTT connected");
    mqttClient.setCallback(handleGeoTarget);
    mqttClient.subscribe(MQTT_TOPIC);
  }
  mqttClient.loop();
}

void handleGeoTarget(char* topic, byte* payload, unsigned int length) {
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
  
  // Second, we add the string termination code at the end of the payload and we convert it to a String object
  payload[strTerminationPos] = '\0';
  String payloadStr((char*)payload);
  /* end */

  StaticJsonDocument<512> jsonDoc;
  DeserializationError err = deserializeJson(jsonDoc, payload);
  if (err != DeserializationError::Ok) {
    Serial.print("Unsuccessful at parsing MQTT JSON: ");
    Serial.println(err.f_str());
    return;
  }
  setTimezone(jsonDoc["lat"], jsonDoc["lng"]);
}

bool steppersMoving() {
  return stepMin.distanceToGo() != 0 || stepHour.distanceToGo() != 0;
}

void setTimezone(double latitude, double longitude) {
  char url[255];
  sprintf(url, TIMEZONE_QUERY, latitude, longitude);
  
  httpClient.useHTTP10(true);
  httpClient.begin(wifiClient, url);
  httpClient.GET();
  
  StaticJsonDocument<384> doc;
  DeserializationError error = deserializeJson(doc, httpClient.getStream());

  if (error) {
    Serial.print("Grabbing timezone details failed: ");
    Serial.println(error.c_str());
    return;
  }

  // TODO: Deal with DST
  configTime(doc["rawOffset"].as<float>() * 3600, 0, NTP_SERVER);
  setSunriseAndSunset(doc["sunrise"], doc["sunset"]);

  httpClient.end();
  updateClock();
}

void setSunriseAndSunset(String sunriseStr, String sunsetStr) {
  sunriseMins = timeToDayMins(sunriseStr.substring(11, 13).toInt(), sunriseStr.substring(14, 16).toInt());
  sunsetMins = timeToDayMins(sunsetStr.substring(11, 13).toInt(), sunsetStr.substring(14, 16).toInt());
}

int timeToDayMins(int hours, int mins) {
  return (hours * 60 + mins);
}

void updateClock() {
  if (now == 0 || steppersMoving() || noNeedToUpdate())
    return;

  Serial.print("Setting time to: "); Serial.print(now->tm_hour); Serial.print(":"); Serial.println(now->tm_min); 

  updateMotors();
  updateLEDs();
    
  lastUpdated = updatedInt();
}

int updatedInt() {
  return (now->tm_min * 60 + now->tm_sec) / UPDATE_FREQ_S;
}

bool noNeedToUpdate() {
  return lastUpdated == updatedInt();
}

void updateMotors() {
  moveCircular(&stepMin, (now->tm_min + now->tm_sec/60.0) * MINS_STEPS);
  moveCircular(&stepHour, (now->tm_hour + now->tm_min/60.0) * HOURS_STEPS);
}

// TODO: Pick good colours; include Warm White?
RgbwColor dayCol(128, 64, 0, 128);
RgbwColor nightCol(0, 0, 128, 0);
RgbwColor offCol(0, 0, 0, 0);

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
    normalizeStepper(stepper);

    // Move Backwards if it's faster
    long curPos = -stepper->currentPosition(); // -ve for clockwise correction
    long newPos = steps % MOTOR_STEPS;
    if (newPos - curPos > MOTOR_STEPS/2) {
      newPos = MOTOR_STEPS - newPos;
    } else if (newPos - curPos + MOTOR_STEPS < MOTOR_STEPS/2) {
      newPos += MOTOR_STEPS;
    }
    
    stepper->moveTo(-newPos); // -ve for clockwise correction
}

void normalizeStepper(AccelStepper* stepper) {
  long pos = stepper->currentPosition();
  if (pos < 0 || pos >= MOTOR_STEPS) {
    int mult = pos / MOTOR_STEPS;
    stepper->setCurrentPosition((pos - mult * MOTOR_STEPS) % MOTOR_STEPS);
  }
}
