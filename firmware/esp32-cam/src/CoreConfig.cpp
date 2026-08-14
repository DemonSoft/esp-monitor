#include <LittleFS.h>
#include <ArduinoJson.h>
#include "CoreCommon.hpp"
#include "CoreConfig.hpp"
#include "CoreWiFi.hpp"
#include "CoreWebServer.hpp"
#include "CoreMQTT.hpp"

Config config;
time_t started = 0;

void setupConfig() {
  if (!loadConfig()) {
    acessPointConfigurationSetup();
    return;
  }
  
  usualConfigurationSetup();
}

void loopConfig() {
}

void acessPointConfigurationSetup() {
  // Start the access point for configuration
  Serial.println("Starting setup AP.");
  espState.access_exists = false;

  startAccessPoint();
}

void usualConfigurationSetup() {
    espState.access_exists = true;

    if (connectToWifi()) {
      webServerSetup(config.ssdp.http_port);
      startSsdp();
      mqttClientSetup();
    } else {
      Serial.println("Failed to connect to WiFi. Starting setup AP.");
      startAccessPoint();
    }  
}



static String getStringValue(const JsonVariant &value) {
  return value.isNull() ? String("") : value.as<String>();
}

static int getIntValue(const JsonVariant &value, int fallback) {
  return value.isNull() ? fallback : value.as<int>();
}

String buildSsdpName() {
  time_t now = time(nullptr);
  if (now <= 0) {
    now = millis() / 1000;
  }
  return String("esp32-") + String(now);
}

void ensureSsdpName() {
  if (config.ssdp.name.length() == 0) {
    createMacAddress(); // Ensure MAC address and name are set
    config.ssdp.name = espState.name; // buildSsdpName(); // Creation name based on time, but we can use the MAC address name instead
  }
}

void applyConfigDefaults() {
  if (config.mdns.length() == 0) {
    config.mdns = "local.esp32";
  }
  if (config.ssdp.http_port == 0) {
    config.ssdp.http_port = 80;
  }
  if (config.mqtt.port == 0) {
    config.mqtt.port = 1883;
  }
  if (config.mqtt.root.length() == 0) {
    config.mqtt.root = "esp32";
  }
  if (config.ssdp.url.length() == 0) {
    config.ssdp.url = "/";
  }
  ensureSsdpName();
}

bool loadConfig() {
  if (!LittleFS.exists("/config.json")) {
    Serial.println("No config.json found.");
    return false;
  }

  File configFile = LittleFS.open("/config.json", "r");
  if (!configFile) {
    return false;
  }

  JsonDocument doc;
  DeserializationError error = deserializeJson(doc, configFile);
  configFile.close();
  if (error) {
    return false;
  }

  if (doc["wifi"].is<JsonObject>()) {
    JsonObject wifi = doc["wifi"].as<JsonObject>();
    config.wifi.ssid = getStringValue(wifi["ssid"]);
    config.wifi.pass = getStringValue(wifi["pass"]);
  }
  if (doc["ssdp"].is<JsonObject>()) {
    JsonObject ssdp = doc["ssdp"].as<JsonObject>();
    config.ssdp.name = getStringValue(ssdp["name"]);
    config.ssdp.url = getStringValue(ssdp["url"]);
    config.ssdp.serial = getStringValue(ssdp["serial"]);
    config.ssdp.device_type = getStringValue(ssdp["device_type"]);
    config.ssdp.schema_url = getStringValue(ssdp["schema_url"]);
    config.ssdp.http_port = getIntValue(ssdp["http_port"], config.ssdp.http_port);
    if (ssdp["model"].is<JsonObject>()) {
      JsonObject model = ssdp["model"].as<JsonObject>();
      config.ssdp.model.name = getStringValue(model["name"]);
      config.ssdp.model.number = getStringValue(model["number"]);
      config.ssdp.model.url = getStringValue(model["url"]);
    }
    if (ssdp["manufacturer"].is<JsonObject>()) {
      JsonObject manufacturer = ssdp["manufacturer"].as<JsonObject>();
      config.ssdp.manufacturer.name = getStringValue(manufacturer["name"]);
      config.ssdp.manufacturer.url = getStringValue(manufacturer["url"]);
    }
  }
  config.mdns = getStringValue(doc["mdns"]);
  if (doc["mqtt"].is<JsonObject>()) {
    JsonObject mqtt = doc["mqtt"].as<JsonObject>();
    config.mqtt.host = getStringValue(mqtt["host"]);
    config.mqtt.port = getIntValue(mqtt["port"], config.mqtt.port);
    config.mqtt.user = getStringValue(mqtt["user"]);
    config.mqtt.pass = getStringValue(mqtt["pass"]);
    config.mqtt.root = getStringValue(mqtt["root"]);
  }
  if (doc["camera"].is<JsonObject>()) {
    JsonObject camera = doc["camera"].as<JsonObject>();
    config.camera.mode = getIntValue(camera["mode"], config.camera.mode);
  }

  applyConfigDefaults();
  Serial.println("Configuration loaded from config.json");
  Serial.println(config.wifi.ssid);
  Serial.println(config.wifi.pass);
  Serial.println(config.ssdp.name);

  return config.wifi.ssid.length() && config.wifi.pass.length() && config.mqtt.host.length();
}

bool saveConfig() {
  File configFile = LittleFS.open("/config.json", "w");
  if (!configFile) {
    return false;
  }

  JsonDocument doc;
  JsonObject wifi = doc["wifi"].to<JsonObject>();
  wifi["ssid"] = config.wifi.ssid;
  wifi["pass"] = config.wifi.pass;

  JsonObject ssdp = doc["ssdp"].to<JsonObject>();
  ssdp["name"] = config.ssdp.name;
  ssdp["url"] = config.ssdp.url;
  ssdp["serial"] = config.ssdp.serial;
  ssdp["device_type"] = config.ssdp.device_type;
  ssdp["schema_url"] = config.ssdp.schema_url;
  ssdp["http_port"] = config.ssdp.http_port;
  JsonObject model = ssdp["model"].to<JsonObject>();
  model["name"] = config.ssdp.model.name;
  model["number"] = config.ssdp.model.number;
  model["url"] = config.ssdp.model.url;
  JsonObject manufacturer = ssdp["manufacturer"].to<JsonObject>();
  manufacturer["name"] = config.ssdp.manufacturer.name;
  manufacturer["url"] = config.ssdp.manufacturer.url;

  doc["mdns"] = config.mdns;

  JsonObject mqtt = doc["mqtt"].to<JsonObject>();
  mqtt["host"] = config.mqtt.host;
  mqtt["port"] = config.mqtt.port;
  mqtt["user"] = config.mqtt.user;
  mqtt["pass"] = config.mqtt.pass;
  mqtt["root"] = config.mqtt.root;

  JsonObject camera = doc["camera"].to<JsonObject>();
  camera["mode"] = config.camera.mode;

  if (serializeJson(doc, configFile) == 0) {
    configFile.close();
    return false;
  }

  configFile.close();
  return true;
}

bool deleteConfigFile() {
  if (LittleFS.exists("/config.json")) {
    return LittleFS.remove("/config.json");
  }
  return true;
}

void mergeConfigObject(const JsonObject &source) {
  if (source["wifi"].is<JsonObject>()) {
    JsonObject wifi = source["wifi"].as<JsonObject>();
    if (!wifi["ssid"].isNull()) {
      String v = getStringValue(wifi["ssid"]);
      if (v.length()) config.wifi.ssid = v;
    }
    if (!wifi["pass"].isNull()) {
      String v = getStringValue(wifi["pass"]);
      if (v.length()) config.wifi.pass = v;
    }
  }
  if (source["ssdp"].is<JsonObject>()) {
    JsonObject ssdp = source["ssdp"].as<JsonObject>();
    if (!ssdp["name"].isNull()) {
      String v = getStringValue(ssdp["name"]);
      if (v.length()) config.ssdp.name = v;
    }
    if (!ssdp["url"].isNull()) {
      String v = getStringValue(ssdp["url"]);
      if (v.length()) config.ssdp.url = v;
    }
    if (!ssdp["serial"].isNull()) {
      String v = getStringValue(ssdp["serial"]);
      if (v.length()) config.ssdp.serial = v;
    }
    if (!ssdp["device_type"].isNull()) {
      String v = getStringValue(ssdp["device_type"]);
      if (v.length()) config.ssdp.device_type = v;
    }
    if (!ssdp["schema_url"].isNull()) {
      String v = getStringValue(ssdp["schema_url"]);
      if (v.length()) config.ssdp.schema_url = v;
    }
    if (!ssdp["http_port"].isNull()) {
      String tmp = getStringValue(ssdp["http_port"]);
      if (tmp.length()) config.ssdp.http_port = getIntValue(ssdp["http_port"], config.ssdp.http_port);
    }
    if (ssdp["model"].is<JsonObject>()) {
      JsonObject model = ssdp["model"].as<JsonObject>();
      if (!model["name"].isNull()) {
        String v = getStringValue(model["name"]);
        if (v.length()) config.ssdp.model.name = v;
      }
      if (!model["number"].isNull()) {
        String v = getStringValue(model["number"]);
        if (v.length()) config.ssdp.model.number = v;
      }
      if (!model["url"].isNull()) {
        String v = getStringValue(model["url"]);
        if (v.length()) config.ssdp.model.url = v;
      }
    }
    if (ssdp["manufacturer"].is<JsonObject>()) {
      JsonObject manufacturer = ssdp["manufacturer"].as<JsonObject>();
      if (!manufacturer["name"].isNull()) {
        String v = getStringValue(manufacturer["name"]);
        if (v.length()) config.ssdp.manufacturer.name = v;
      }
      if (!manufacturer["url"].isNull()) {
        String v = getStringValue(manufacturer["url"]);
        if (v.length()) config.ssdp.manufacturer.url = v;
      }
    }
  }
  if (!source["mdns"].isNull()) {
    String v = getStringValue(source["mdns"]);
    if (v.length()) config.mdns = v;
  }
  if (source["mqtt"].is<JsonObject>()) {
    JsonObject mqtt = source["mqtt"].as<JsonObject>();
    if (!mqtt["host"].isNull()) {
      String v = getStringValue(mqtt["host"]);
      if (v.length()) config.mqtt.host = v;
    }
    if (!mqtt["port"].isNull()) {
      String tmp = getStringValue(mqtt["port"]);
      if (tmp.length()) config.mqtt.port = getIntValue(mqtt["port"], config.mqtt.port);
    }
    if (!mqtt["user"].isNull()) {
      String v = getStringValue(mqtt["user"]);
      if (v.length()) config.mqtt.user = v;
    }
    if (!mqtt["pass"].isNull()) {
      String v = getStringValue(mqtt["pass"]);
      if (v.length()) config.mqtt.pass = v;
    }
    if (!mqtt["root"].isNull()) {
      String v = getStringValue(mqtt["root"]);
      if (v.length()) config.mqtt.root = v;
    }
  }
  if (!source["camera"].isNull()) {
     JsonObject camera = source["camera"].as<JsonObject>();
     if (!camera["mode"].isNull()) {
      String v = getStringValue(camera["mode"]);
      if (v.length()) config.camera.mode = getIntValue(camera["mode"], config.camera.mode);
     }
  }

  applyConfigDefaults();
}

String mqttBaseTopic() {
  return config.mqtt.root + "/" + config.ssdp.name;
}

String mqttActionTopic() {
  return mqttBaseTopic() + "/action";
}

String mqttStateTopic() {
  return mqttBaseTopic() + "/state";
}

String mqttPhotoTopic() {
  return mqttBaseTopic() + "/photo";
}
