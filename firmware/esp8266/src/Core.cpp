#include "Core.hpp"
#include "CoreCommon.hpp"
#include "CoreButton.hpp"
#include "CoreWiFi.hpp"
#include "CoreFS.hpp"
#include "CoreConfig.hpp"
#include "CoreMQTT.hpp"
#include "CoreWebServer.hpp"

EspState espState;

void coreSetup() {
    setupStart();
    setupButton();
    setupFS();
    setupWiFi();
    setupConfig();
    setupFinish();
} 

void coreLoop() {
    loopWiFi();
    loopWebServer();
    loopMqtt();
    loopButton();
} 

void setupStart() {
    Serial.begin(115200);
    Serial.flush();
    // Serial.begin(74880);
    Serial.println("\n\n\n");
    Serial.println("================================");
    Serial.println("Initializing...");
    Serial.println("================================");
    delay(10);

    pinMode(LED_BUILTIN, OUTPUT);
    Serial.print("Memory: ");
    Serial.print(ESP.getFreeHeap());
    Serial.println("");
    Serial.print("Flash: ");
    Serial.print(ESP.getFlashChipSize());
    Serial.println("");
}

void setupFinish() {
    Serial.println("================================");
    Serial.println("Setup finshed!");
    Serial.println("================================");
    Serial.println("\n");
}


