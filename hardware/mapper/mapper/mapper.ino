#include <WiFiManager.h>
#include <PubSubClient.h>
#include <AccelStepper.h>
#include <ArduinoJson.h>
#include <math.h>
#include "vars.h";

#define APP_NAME "jan-poka:mapper"
#define ONE_ROTATION_STEPS 8192
#define POWER_OFF_DELAY 2000
#define MAX_SPEED 24576
#define ACCELERATION 4096

int INNER_DIR     = 14; // D5
int INNER_STP     = 13; // D7
int INNER_DISABLE = 12; // D6
int OUTER_DIR     = 4;  // D2
int OUTER_STP     = 5;  // D1
int OUTER_DISABLE = 16; // D0
int BOOTING       = 2;  // D4

WiFiClient wifiClient;
PubSubClient mqttClient(wifiClient);

AccelStepper inner = AccelStepper(AccelStepper::DRIVER, INNER_STP, INNER_DIR);
AccelStepper outer = AccelStepper(AccelStepper::DRIVER, OUTER_STP, OUTER_DIR);

void setup() {
  Serial.begin(115200);
  pinMode(BOOTING, OUTPUT);
  digitalWrite(BOOTING, HIGH);
  pinMode(INNER_DISABLE, OUTPUT);
  digitalWrite(INNER_DISABLE, LOW);
  pinMode(OUTER_DISABLE, OUTPUT);
  digitalWrite(OUTER_DISABLE, LOW);

  mqttClient.setServer(MQTT_HOST, MQTT_PORT);

  inner.setMaxSpeed(MAX_SPEED);
  outer.setMaxSpeed(MAX_SPEED);
  inner.setAcceleration(ACCELERATION);
  outer.setAcceleration(ACCELERATION);
  inner.setPinsInverted(true, false, false);
  outer.setPinsInverted(true, false, false);

  mqttConnectionLoop();
  Serial.println("\nBooted");
  digitalWrite(BOOTING, LOW);
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

  if (jsonDoc["reset"]) {
    inner.setCurrentPosition(0);
    outer.setCurrentPosition(0);
  }
  goTo(jsonDoc["r1"], jsonDoc["r2"]);
}

void goTo(double r1, double r2) {
  Serial.print("Going to, r1=");
  Serial.print(r1);
  Serial.print(", r2=");
  Serial.println(r2);

  inner.moveTo(r1 * ONE_ROTATION_STEPS);
  outer.moveTo(r2 * ONE_ROTATION_STEPS);
}

unsigned long powerOffAt = 0;

void loop() {
  bool noMove = inner.distanceToGo() == 0 && outer.distanceToGo() == 0;
  if (noMove) {
    if (powerOffAt == 0) {
      powerOffAt = millis() + POWER_OFF_DELAY;
    } else if (millis() >= powerOffAt) {
      digitalWrite(INNER_DISABLE, HIGH);
      digitalWrite(OUTER_DISABLE, HIGH);
      powerOffAt = 0;
    }
  } else {
    digitalWrite(INNER_DISABLE, LOW);
    digitalWrite(OUTER_DISABLE, LOW);
  }

  inner.run();
  outer.run();
  mqttConnectionLoop();
}
