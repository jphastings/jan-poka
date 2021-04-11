#include <WiFiManager.h>
#include <PubSubClient.h>
#include <AccelStepper.h>
#include <ArduinoJson.h>
#include <math.h>

#define APP_NAME "concentric-gear-pointer"
#define ONE_ROTATION_STEPS 8192

int INNER_DIR = 14; // D5
int INNER_STP = 12; // D6
int OUTER_DIR = 5;  // D1
int OUTER_STP = 4;  // D2
int DISABLE   = 16; // D0
int BOOTING   = 2;  // D4

int MAX_SPEED = 9000;

WiFiClient wifiClient;
PubSubClient mqttClient(wifiClient);
String mqttUsername;
String mqttPassword;
String mqttTopic;

AccelStepper inner = AccelStepper(AccelStepper::DRIVER, INNER_STP, INNER_DIR);
AccelStepper outer = AccelStepper(AccelStepper::DRIVER, OUTER_STP, OUTER_DIR);

void setup() {
  pinMode(BOOTING, OUTPUT);
  digitalWrite(BOOTING, HIGH);
  pinMode(DISABLE, OUTPUT);
  digitalWrite(DISABLE, HIGH);

  WiFiManager wifiManager;
  WiFiManagerParameter custom_mqtt_host("host", "MQTT Host", "mqtt.local", 64);
  wifiManager.addParameter(&custom_mqtt_host);
  WiFiManagerParameter custom_mqtt_port("port", "MQTT Port", "1883", 5);
  wifiManager.addParameter(&custom_mqtt_port);
  WiFiManagerParameter custom_mqtt_username("username", "MQTT Username", "jan-poka", 32);
  wifiManager.addParameter(&custom_mqtt_username);
  WiFiManagerParameter custom_mqtt_password("password", "MQTT Password", "", 32);
  wifiManager.addParameter(&custom_mqtt_password);
  WiFiManagerParameter custom_mqtt_topic("topic", "MQTT Topic", "home/geo/target", 32);
  wifiManager.addParameter(&custom_mqtt_topic);

  wifiManager.autoConnect(APP_NAME);

  mqttClient.setServer(custom_mqtt_host.getValue(), custom_mqtt_port.getValue().toInt());
  strcpy(mqttUsername, custom_mqtt_username.getValue());
  strcpy(mqttPassword, custom_mqtt_password.getValue());
  strcpy(mqttTopic, custom_mqtt_topic.getValue());

  Serial.begin(115200);

  inner.setMaxSpeed(MAX_SPEED);
  outer.setMaxSpeed(MAX_SPEED);

  mqttConnectionLoop();
  Serial.println("\nBooted");
  digitalWrite(BOOTING, LOW);
}

void mqttConnectionLoop() {
  if (!mqttClient.connected()) {
    if (!mqttClient.connect(APP_NAME, mqttUsername, mqttPassword)) {
      Serial.println("Failed to connect to MQTT broker");
      // TODO: Backoff
      return
    }
    Serial.println("MQTT connected");
    mqttClient.setCallback(handleGeoTarget);
    mqttClient.subscribe(mqttTopic);
  }
  mqttClient.loop();
}

void handleGeoTarget(const char[] topic, byte* payloadPtr, unsigned int length) {
  String payload;
  // Ugh. There's definitely a better way than this, but I can't test right now, so TODO.
  for (int i = 0; i < length; i++) {
    payload += (char)payloadPtr[i];
  }

  StaticJsonDocument<512> jsonDoc;
  DeserializationError err = deserializeJson(jsonDoc, payload);
  if (err != DeserializationError::Ok) {
    Serial.print("Unsuccessful at parsing MQTT JSON: ");
    Serial.println(err.f_str());
    return;
  }
  goTo(jsonDoc["azi"], jsonDoc["ele"]);
}

void goTo(double azimuth, double elevation) {
  double zRotate = azimuth / 180.0;
  double xyRotate = elevation / 270.0;
  long innerRotation = floor(zRotate * ONE_ROTATION_STEPS);
  long outerRotation = floor((xyRotate - 2 * zRotate) * ONE_ROTATION_STEPS);

  long stepsForInner = innerRotation - inner.currentPosition();
  long stepsForOuter = outerRotation - outer.currentPosition();

  inner.moveTo(innerRotation);
  outer.moveTo(outerRotation);

  double innerSpeed = MAX_SPEED;
  double outerSpeed = MAX_SPEED;
  if (stepsForInner > stepsForOuter) {
    outerSpeed = innerSpeed * stepsForInner / stepsForOuter;
  } else {
    innerSpeed = outerSpeed * stepsForOuter / stepsForInner;
  }

  inner.setSpeed(innerSpeed);
  outer.setSpeed(outerSpeed);
}

void loop() {
  bool noMove = inner.distanceToGo() == 0 && outer.distanceToGo() == 0;
  digitalWrite(DISABLE, noMove ? HIGH : LOW);
  
  inner.runSpeedToPosition();
  outer.runSpeedToPosition();
  mqttConnectionLoop();
}
