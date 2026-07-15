#ifndef CoreWebServer_h
#define CoreWebServer_h
#include <ESP8266WebServer.h>
#include <ESP8266mDNS.h>

void loopWebServer();
void webServerSetup(int port);
void shutdownWebServer(bool stopAp);
void handleConfigurationRequest();
void handleDescriptionRequest();

#endif