#include "CoreFS.hpp"

void setupFS() {
    if (!LittleFS.begin()) {
    Serial.println("LittleFS mount failed");
    return;
    }
    Serial.println("LittleFS mounted successfully");
} 

void loopFS() {
} 