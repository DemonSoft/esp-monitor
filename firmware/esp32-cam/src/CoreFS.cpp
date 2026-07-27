#include "CoreFS.hpp"

void setupFS() {
    bool ok = LittleFS.begin(true);
    if (!ok) {
        Serial.println("LittleFS mount failed");
        return;
    }
    Serial.println("LittleFS mounted successfully");
}

void loopFS() {
} 