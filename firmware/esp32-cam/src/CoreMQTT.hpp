#ifndef CoreMQTT_h
#define CoreMQTT_h
#include <AsyncMqttClient.h>

void loopMqtt();
void mqttClientSetup();
void mqttDisconnect();
void publishStateMessage();
void publishPhoto(const char* photoData, size_t length);
void connectToMqtt();
void resetMQTTTimers();
void waitingMqttDisconnect();
bool checkTimerInterval();
bool checkPinStateChanged();
void publishState();
void loopPublish();
void removeTopic(String topic);
bool isMqttConnected();

String createTopic(const String &name);

void onMqttConnect(bool sessionPresent);
void onMqttSubscribe(uint16_t packetId, uint8_t qos);
void onMqttUnsubscribe(uint16_t packetId);
void onMqttPublish(uint16_t packetId);
void mqttMessage(char* topic, char* payload, size_t len, size_t index, size_t total, int qos, int dup, int retain);


void onMqttDisconnect(AsyncMqttClientDisconnectReason reason);
void onMqttMessage(char* topic, char* payload, 
     struct AsyncMqttClientMessageProperties properties, size_t len, size_t index, size_t total);


#endif