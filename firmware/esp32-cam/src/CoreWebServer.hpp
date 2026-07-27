#ifndef CoreWebServer_h
#define CoreWebServer_h
#include <WebServer.h>

void loopWebServer();
void webServerSetup(int port);
void shutdownWebServer(bool stopAp);
void handleConfigurationRequest();
void handleDescriptionRequest();

#endif