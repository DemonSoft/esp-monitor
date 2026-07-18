#include "CoreCommon.hpp"
#include "CoreWiFi.hpp"
#include "CoreWebServer.hpp"
#include "main.hpp"
#include "CoreConfig.hpp"
#include "CoreTime.hpp"
#include "CoreMQTT.hpp"
#include <ESPmDNS.h>
#include <ESP32SSDP.h>

// Access point credentials.
const char *apSsid = "esp-setup";
const char *apPassword = "1qazxsw2";
const int wifiPOWER = 10; // Wi-Fi output power in dBm (0-20). Default is 20 dBm. Lower values reduce power consumption and range.
// Data for Wi-Fi reconnection
const int wifiTimeOutMS = 5000; // Maximum time to wait for WiFi connection before retrying
unsigned long wifiReconnectTime = 0;
bool wifiTimerActive = false;
const int maxDuration = 20000; // Maximum duration to wait for WiFi connection in milliseconds

void setupWiFi() {
  WiFi.disconnect(true);
  wifiReconnectTime = 0;
  wifiTimerActive = false;
}

void loopWiFi() { 
  waitingForWiFiDisconnection();
  checkWifiConnection();
}

void waitingForWiFiDisconnection() {
  // If the WiFi reconnection timer is active and the timeout has elapsed, attempt to reconnect to WiFi. 
  // About 5 seconds after the WiFi connection is lost, the device will try to reconnect to WiFi.
  if (wifiTimerActive && (millis() - wifiReconnectTime >= wifiTimeOutMS)) {
    connectToWifi();
    wifiTimerActive = false;
  }
}  

bool connectToWifi() {
  Serial.print("Wi-Fi Connecting...");

  if (config.wifi.ssid.length() == 0 || config.wifi.pass.length() == 0) {
    Serial.println("WiFi credentials missing");
    return false;
  }

  WiFi.mode(WIFI_STA);
  WiFi.setTxPower(WIFI_POWER_11dBm);
  WiFi.begin(config.wifi.ssid.c_str(), config.wifi.pass.c_str());

  unsigned long start = millis();
  while (WiFi.status() != WL_CONNECTED && millis() - start < maxDuration) {
    delay(500);
    Serial.print('.');
  }

  if (WiFi.status() != WL_CONNECTED) {
    Serial.println("\nWiFi connection timed out");
    return false;
  }

  createMacAddress();
  Serial.println();
  Serial.print("Connected to ");
  Serial.println(config.wifi.ssid);
  Serial.print("IP address: ");
  Serial.println(WiFi.localIP());

  if (MDNS.begin(config.mdns.c_str())) {
    Serial.println("MDNS started with name: " + config.mdns);
  }

  setupTime();
  syncTime();
  started = getUnixTime();

  Serial.print("UNIX time: ");
  Serial.println(started);

  return true;
}

void createMacAddress() {
  byte mac[6];
  WiFi.macAddress(mac);

  String address = 
                    String(mac[5], HEX) +
                    String(mac[4], HEX) +
                    String(mac[3], HEX) +
                    String(mac[2], HEX) +
                    String(mac[1], HEX) +
                    String(mac[0], HEX);

  String name = "esp-" + address;
  espState.mac = address;
  espState.name = name;

  Serial.println("");
  Serial.println("MAC: " + address);
  Serial.println("NAME: " + name);
}

void startAccessPoint() {
  Serial.println();
  Serial.println("Configuring access point...");
  WiFi.mode(WIFI_AP);
  WiFi.setSleep(false);
  WiFi.setTxPower(WIFI_POWER_11dBm);
  WiFi.softAPConfig(IPAddress(192, 168, 4, 1), IPAddress(192, 168, 4, 1), IPAddress(255, 255, 255, 0));
  bool apStarted = WiFi.softAP(apSsid, apPassword, 1, false, 4);
  if (!apStarted) {
    Serial.println("Failed to start access point");
  }
  IPAddress myIP = WiFi.softAPIP();
  Serial.print("AP IP address: ");
  Serial.println(myIP);

  createMacAddress();
  webServerSetup(80);
}

void startSsdp() {
  ensureSsdpName();

  SSDP.setDeviceType(config.ssdp.device_type.c_str());
  SSDP.setSchemaURL(config.ssdp.schema_url.c_str());
  SSDP.setHTTPPort(config.ssdp.http_port);
  SSDP.setName(config.ssdp.name.c_str());
  SSDP.setURL(config.ssdp.url.c_str());
  SSDP.setSerialNumber(config.ssdp.serial.c_str());
  SSDP.setModelName(config.ssdp.model.name.c_str());
  SSDP.setModelNumber(config.ssdp.model.number.c_str());
  SSDP.setModelURL(config.ssdp.model.url.c_str());
  SSDP.setManufacturer(config.ssdp.manufacturer.name.c_str());
  SSDP.setManufacturerURL(config.ssdp.manufacturer.url.c_str());
  SSDP.begin();
  Serial.println("SSDP started with name: " + config.ssdp.name);
}

void resetWiFiTimers() {
  // Stop MQTT reconnection and start WiFi reconnection timer
  wifiTimerActive = true;
  wifiReconnectTime = millis();
}

void checkWifiConnection() {
  static bool mqttConnected = false;

  if (WiFi.isConnected() && !mqttConnected) {
    Serial.println("WiFi connected!");
    connectToMqtt();
    mqttConnected = true;
  } else if (!WiFi.isConnected() && mqttConnected) {
    Serial.println("WiFi lost connection");
    resetMQTTTimers();
    resetWiFiTimers();
    mqttConnected = false;
  }
}
