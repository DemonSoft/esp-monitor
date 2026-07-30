#include "CoreCommon.hpp"
#include "CoreMQTT.hpp"
#include "CoreWiFi.hpp"
#include "CoreCamera.hpp"
#include "CoreConfig.hpp"
#include "CoreTime.hpp"
#include <ArduinoJson.h>

// ESP32 uses simple timer management with millis()
AsyncMqttClient mqttClient;
int mqttPort;

bool checkPinStateChanged();
void publishStateMessage();
void publishState();

// Simple timer management for ESP32
unsigned long mqttReconnectTime = 0;
bool mqttTimerActive = false;
const int mqttTimeOutMS = 5000; // Maximum time to wait for MQTT connection before retrying

unsigned long previousMillis = 0;  // хранится время последней публикации
const long interval = 1000;        // максимальная частота публикации состояния
                                   // данных от датчика


bool hasPreviousPinState = false;

// ESP32-DevKit-V1 exposes many GPIOs, but not all are safe or practical to sample as general-purpose inputs.
// We read a broad set of usable GPIOs here. Excluded pins are reserved for boot/flash/UART or are not exposed as regular GPIOs.
constexpr int kPinCount = 22;
const int gpioPins[kPinCount] = {
  4, 5, 12, 13, 14, 15, 16, 17, 18, 19,
  21, 22, 23, 25, 26, 27, 32, 33, 34, 35,
  36, 39
};
int previousPinStates[kPinCount] = {0};
 
void mqttClientSetup() {

  mqttPort = config.mqtt.port;
  mqttClient.setServer(config.mqtt.host.c_str(), mqttPort);

  // User credentials for MQTT authentication
  if (config.mqtt.user.length() && config.mqtt.pass.length()) {
    mqttClient.setCredentials(config.mqtt.user.c_str(), config.mqtt.pass.c_str());
  }

  // Set the client ID for MQTT connection
  mqttClient.setClientId(espState.name.c_str());

  // Initialize timers for ESP32
  mqttReconnectTime = 0;
  mqttTimerActive = false;

  // Setup MQTT callbacks
  mqttClient.onConnect(onMqttConnect);
  mqttClient.onDisconnect(onMqttDisconnect);
  mqttClient.onSubscribe(onMqttSubscribe);
  mqttClient.onUnsubscribe(onMqttUnsubscribe);
  mqttClient.onMessage(onMqttMessage);
  mqttClient.onPublish(onMqttPublish);

  Serial.println("MQTT HOST: " + config.mqtt.host);
  Serial.println("MQTT PORT: " + String(config.mqtt.port));
  Serial.println("MQTT USER: " + config.mqtt.user);
  Serial.println("MQTT PASS: " + config.mqtt.pass);
  Serial.println("MQTT CLIENT ID: " + espState.name);
  Serial.println("MQTT CLIENT ID: " + String(mqttClient.getClientId()));

  Serial.println("MQTT ACTION TOPIC: " + mqttActionTopic());
  Serial.println("MQTT STATE TOPIC: " + mqttStateTopic());
  Serial.println("MQTT SUBSCRIBE TOPIC: " + config.mqtt.root + "/" + config.ssdp.name + "/#");
  Serial.println("MQTT BASE TOPIC: " + mqttBaseTopic());
  Serial.println("MQTT ROOT: " + config.mqtt.root);

}

void loopMqtt() {
  waitingMqttDisconnect();
  publishStateMessage();
}

void waitingMqttDisconnect() {
     // Check MQTT reconnection timer
   if (mqttTimerActive && (millis() - mqttReconnectTime >= mqttTimeOutMS)) {
     connectToMqtt();
     mqttTimerActive = false;
   }
}

void connectToMqtt() { 
  Serial.println("Connecting to MQTT...");
  mqttClient.connect();
}

void resetMQTTTimers() {
  // Stop MQTT reconnection and start WiFi reconnection timer
  mqttTimerActive = false;
}

bool checkTimerInterval() {
  unsigned long currentMillis = millis();

  if (previousMillis > currentMillis)  // Correction after counter overflow through 50 days.
    previousMillis = interval - currentMillis;
  
  
  if (currentMillis - previousMillis < interval) return false;
   
    previousMillis = currentMillis;

   return true;
}

String createTopic(const String &name) {
  return mqttBaseTopic() + "/" + name;
}

bool checkPinStateChanged() {
  int currentPins[kPinCount];
  for (int i = 0; i < kPinCount; i++) {
    currentPins[i] = digitalRead(gpioPins[i]);
  }

  bool changed = false;
  for (int i = 0; i < kPinCount; i++) {
    if (!hasPreviousPinState || currentPins[i] != previousPinStates[i]) {
      changed = true;
      break;
    }
  }

  if (changed) {
    for (int i = 0; i < kPinCount; i++) {
      previousPinStates[i] = currentPins[i];
    }
    hasPreviousPinState = true;
  }

  return changed;
}

void publishStateMessage() {
    if (!espState.mqtt_connected) return;
    if (!checkTimerInterval()) return;
    if (!checkPinStateChanged()) return;
    publishState();
}

void publishState() {
    JsonDocument doc;
    doc["SSDP"] = config.ssdp.name;
    doc["MDNS"] = config.mdns;
    doc["Started"] = started;
    doc["Updated"] = getUnixTime();
    JsonObject pins = doc["Pins"].to<JsonObject>();
    for (int i = 0; i < kPinCount; i++) {
      String label = "GPIO" + String(gpioPins[i]);
      pins[label] = digitalRead(gpioPins[i]);
    }

    String payload;
    serializeJson(doc, payload);
    String topic = mqttStateTopic();
    mqttClient.publish(topic.c_str(), 1, false, payload.c_str());
}

void handleMqttAction(const String &payload) {
    JsonDocument doc;
    DeserializationError err = deserializeJson(doc, payload);
    if (err) {
      Serial.print("Invalid MQTT action JSON: ");
      Serial.println(err.c_str());
      return;
    }

    if (!doc["config"].is<JsonObject>()) {
      return;
    }

    String oldActionTopic = mqttActionTopic(); // Store the old action topic before merging the new config
    JsonObject configObject = doc["config"].as<JsonObject>();
    mergeConfigObject(configObject);
    if (!saveConfig()) {
      Serial.println("Failed to save config from MQTT action");
      return;
    }
    
    removeTopic(oldActionTopic); // Clean up the action topic after processing the action

    Serial.println("Config updated from MQTT action. Rebooting...");
    delay(1000); // We need to wait a bit before restarting to ensure the message is sent
    ESP.restart();
}

void publishPhoto(const char* photoData, size_t length) {

    if (!mqttClient.connected()) return;
    
    String topic = mqttPhotoTopic();
    uint16_t packetId = mqttClient.publish(
    topic.c_str(), 
    0,                  // QoS 0 (for photo stream)
    false,              // Retain
    photoData, 
    length);
}

bool isMqttConnected() {
    return mqttClient.connected();
}

// MQTT callback functions
// Thre are subscribed to MQTT events and handle them accordingly. 
// These functions are called when the corresponding MQTT event occurs.
void onMqttConnect(bool sessionPresent) {
  Serial.println("Connected to MQTT.");
  Serial.print("Session present: ");
  Serial.println(sessionPresent);
  String actionTopic = mqttActionTopic();
  String subscribeTopic = config.mqtt.root + "/" + config.ssdp.name + "/#";
  uint16_t packetIdSub = mqttClient.subscribe(subscribeTopic.c_str(), 0);  // QoS 0 for subscription
  Serial.print("Subscribing at QoS 0, packetId: ");
  Serial.println(packetIdSub);
  espState.mqtt_connected = true;
  publishStateMessage();
}


void mqttDisconnect() {
  
  espState.mqtt_connected  = false;
  Serial.println("");
  Serial.println("Disconnected from MQTT.");
  
  previousMillis = 0;
  
  // Start MQTT reconnection timer for ESP32
  if (WiFi.isConnected()) {
    mqttTimerActive = true;
    mqttReconnectTime = millis();
  }
}
 
void onMqttSubscribe(uint16_t packetId, uint8_t qos) {
  Serial.println("Subscribe acknowledged.");
  Serial.print("  packetId: ");
  Serial.println(packetId);
  Serial.print("  qos: ");
  Serial.println(qos);
}
 
void onMqttUnsubscribe(uint16_t packetId) {
  Serial.println("Unsubscribe acknowledged.");
  Serial.print("  packetId: ");
  Serial.println(packetId);
}
 
void onMqttPublish(uint16_t packetId) {
}
 
void mqttMessage(char* topic, char* payload, size_t len, size_t index, size_t total, int qos, int dup, int retain) {
  String messageTemp;
  for (size_t i = 0; i < len; i++) {
    messageTemp += (char)payload[i];
  }

  if (messageTemp == "{}") {
    Serial.printf("-  Topic %s was removed.\n", topic);
    return;
  }

  String topicStr = String(topic);
  String actionTopic = mqttActionTopic();

  if (topicStr == actionTopic) {
    handleMqttAction(messageTemp);
  }

  Serial.println("");
  Serial.println("Publish received.");
  Serial.print("  topic: ");
  Serial.println(topicStr);
  Serial.print("  message: ");
  Serial.println(messageTemp);
  Serial.print("  qos: ");
  Serial.println(qos);
  Serial.print("  dup: ");
  Serial.println(dup);
  Serial.print("  retain: ");
  Serial.println(retain);
  Serial.print("  len: ");
  Serial.println(len);
  Serial.print("  index: ");
  Serial.println(index);
  Serial.print("  total: ");
  Serial.println(total);
}

void onMqttDisconnect(AsyncMqttClientDisconnectReason reason) {

  switch ((int)reason) {
    case 0: // TCP_DISCONNECTED
      Serial.println("Disconnected from MQTT: TCP_DISCONNECTED");
      break;
    case 1: // MQTT_UNACCEPTABLE_PROTOCOL_VERSION
      Serial.println("Disconnected from MQTT: MQTT_UNACCEPTABLE_PROTOCOL_VERSION");
      break;
    case 2: // MQTT_IDENTIFIER_REJECTED          
      Serial.println("Disconnected from MQTT: MQTT_IDENTIFIER_REJECTED");
      break;
    case 3: // MQTT_SERVER_UNAVAILABLE
      Serial.println("Disconnected from MQTT: MQTT_SERVER_UNAVAILABLE");
      break;
    case 4: // MQTT_MALFORMED_CREDENTIALS
      Serial.println("Disconnected from MQTT: MQTT_MALFORMED_CREDENTIALS");
      break;
    case 5: // MQTT_NOT_AUTHORIZED
      Serial.println("Disconnected from MQTT: MQTT_NOT_AUTHORIZED");
      break;
    case 6: // ESP32_NOT_ENOUGH_SPACE
      Serial.println("Disconnected from MQTT: ESP32_NOT_ENOUGH_SPACE");
      break;
    case 7: // TLS_BAD_FINGERPRINT
      Serial.println("Disconnected from MQTT: TLS_BAD_FINGERPRINT");
      break;  

    default:
      Serial.println("Disconnected from MQTT: UNKNOWN");
      break;
  }

  mqttDisconnect();
}

void onMqttMessage(char* topic, char* payload, 
     struct AsyncMqttClientMessageProperties properties, size_t len, size_t index, size_t total) {
  mqttMessage(topic, payload, len, index, total, properties.qos, properties.dup, properties.retain);
}

void removeTopic(String topic) {
    mqttClient.publish(topic.c_str(), 1, false, "{}");
}
