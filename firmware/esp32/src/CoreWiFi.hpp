#ifndef CoreWiFi_h
#define CoreWiFi_h
#include <WiFi.h>
#include <WiFiClient.h>

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