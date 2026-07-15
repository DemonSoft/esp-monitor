#ifndef CoreWiFi_h
#define CoreWiFi_h
// ESP8266 specific WiFi includes
#include <ESP8266WiFi.h>
#include <WiFiClient.h>
#include <ESP8266WebServer.h>
#include <ESP8266mDNS.h>
#include <ESP8266SSDP.h>

bool connectToWifi();
void createMacAddress();
void startAccessPoint();
void startSsdp();
void resetWiFiTimers();
void checkWifiConnection();
void waitingForWiFiDisconnection();

void setupWiFi();
void loopWiFi();

#endif