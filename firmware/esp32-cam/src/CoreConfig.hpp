#ifndef CoreConfig_h
#define CoreConfig_h

#include <Arduino.h>
#include <ArduinoJson.h>

struct WifiConfig {
  String ssid = "";
  String pass = "";
};

struct SsdModel {
  String name = "";
  String number = "";
  String url = "";
};

struct SsdManufacturer {
  String name = "";
  String url = "";
};

struct SsdConfig {
  String name = "";
  String url = "/";
  String serial = "";
  SsdModel model;
  SsdManufacturer manufacturer;
  String device_type = "upnp:rootdevice";
  String schema_url = "description.xml";
  int http_port = 80;
};

struct MqttConfig {
  String host = "";
  int port = 1883;
  String user = "";
  String pass = "";
  String root = "esp32";
};

struct CameraConfig {
  int mode = 1;
};

struct Config {
  WifiConfig wifi;
  SsdConfig ssdp;
  String mdns = "local.esp32";
  MqttConfig mqtt;
  CameraConfig camera;
};


extern Config config;
extern time_t started;

void setupConfig();
void loopConfig();

bool loadConfig();
bool saveConfig();
bool deleteConfigFile();
void applyConfigDefaults();
void ensureSsdpName();
void mergeConfigObject(const JsonObject &source);
String buildSsdpName();
String mqttBaseTopic();
String mqttActionTopic();
String mqttStateTopic();
String mqttPhotoTopic();

void acessPointConfigurationSetup();
void usualConfigurationSetup();

#endif
