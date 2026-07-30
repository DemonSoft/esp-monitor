#include "Core.hpp"
#include "CoreCommon.hpp"
#include "CoreButton.hpp"
#include "CoreWiFi.hpp"
#include "CoreFS.hpp"
#include "CoreConfig.hpp"
#include "CoreMQTT.hpp"
#include "CoreWebServer.hpp"
#include "CoreCamera.hpp"

#include "soc/soc.h"
#include "soc/rtc_cntl_reg.h"

// Если RTC_CNTL_BROWNOUT_REG не объявлен, используем его прямое смещение
#ifndef RTC_CNTL_BROWNOUT_REG
#define RTC_CNTL_BROWNOUT_REG (DR_REG_RTCCNTL_BASE + 0x00d4)
#endif

EspState espState;

void coreSetup() {
    setupStart();
    setupButton();
    setupFS();
    setupWiFi();
    setupConfig();
    setupCamera();
    setupFinish();
} 

void coreLoop() {
    loopWiFi();
    loopWebServer();
    loopMqtt();
    loopButton();
    loopCamera();
} 

void setupStart() {

    WRITE_PERI_REG(RTC_CNTL_BROWNOUT_REG, 0);

    Serial.begin(115200);
    Serial.flush();
    Serial.println("\n\n\n");
    Serial.println("================================");
    Serial.println("Initializing...");
    Serial.println("================================");
    delay(10);

    // Подробный вывод состояния памяти
    Serial.printf("Total heap: %d\n", ESP.getHeapSize());
    Serial.printf("Free heap: %d\n", ESP.getFreeHeap());
    Serial.printf("Total PSRAM: %d\n", ESP.getPsramSize());
    Serial.printf("Free PSRAM: %d\n", ESP.getFreePsram());

    if (!psramFound()) {
        Serial.println("[RAM] ОШИБКА: PSRAM не обнаружена на плате!");
    } else {
        Serial.println("[RAM] УСПЕХ: PSRAM инициализирована.");
    }    

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


