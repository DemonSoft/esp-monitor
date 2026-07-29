#include "CaptureManager.h"

CaptureManager::CaptureManager(CameraManager& cameraRef) : camera(cameraRef) {}

void CaptureManager::startContinuous(uint32_t intervalSec) {
    currentMode = CaptureMode::CONTINUOUS;
    intervalMs = intervalSec * 1000;
    lastCaptureTime = millis() - intervalMs; // Чтобы первый снимок сделался сразу
    Serial.printf("[CaptureManager] Старт режимов: CONTINUOUS (Интервал: %u сек)\n", intervalSec);
}

void CaptureManager::startBurstCount(uint32_t intervalSec, uint32_t count) {
    currentMode = CaptureMode::BURST_COUNT;
    intervalMs = intervalSec * 1000;
    maxPhotos = count;
    photosTaken = 0;
    lastCaptureTime = millis() - intervalMs;
    Serial.printf("[CaptureManager] Старт режимов: BURST_COUNT (Всего: %u, Интервал: %u сек)\n", count, intervalSec);
}

void CaptureManager::startBurstTime(uint32_t intervalSec, uint32_t durationSec) {
    currentMode = CaptureMode::BURST_TIME;
    intervalMs = intervalSec * 1000;
    durationMs = durationSec * 1000;
    modeStartTime = millis();
    lastCaptureTime = millis() - intervalMs;
    Serial.printf("[CaptureManager] Старт режимов: BURST_TIME (Длительность: %u сек, Интервал: %u сек)\n", durationSec, intervalSec);
}

void CaptureManager::stop() {
    if (currentMode != CaptureMode::IDLE) {
        currentMode = CaptureMode::IDLE;
        Serial.println("[CaptureManager] Съемка остановлена.");
    }
}

void CaptureManager::triggerCapture() {
    camera_fb_t* fb = camera.capture();
    if (fb) {
        photosTaken++;
        Serial.printf("-> Снимок #%u выполнен | Размер: %u байт\n", photosTaken, fb->len);
        
        // На следующем шаге здесь будет отправка по MQTT!
        
        camera.release(fb);
    } else {
        // Мягкое уведомление вместо спама об ошибке
        Serial.println("[CaptureManager] Камера еще не готова или пропуск кадра...");
    }
}

void CaptureManager::update() {
    if (currentMode == CaptureMode::IDLE) return;

    uint32_t now = millis();

    // 1. Проверка условия завершения по времени (K секунд)
    if (currentMode == CaptureMode::BURST_TIME) {
        if (now - modeStartTime >= durationMs) {
            Serial.println("[CaptureManager] Лимит времени K истек.");
            stop();
            return;
        }
    }

    // 2. Проверка таймера интервала (N секунд)
    if (now - lastCaptureTime >= intervalMs) {
        lastCaptureTime = now;
        triggerCapture();

        // 3. Проверка условия завершения по количеству (M кадров)
        if (currentMode == CaptureMode::BURST_COUNT && photosTaken >= maxPhotos) {
            Serial.println("[CaptureManager] Лимит кадров M достигнут.");
            stop();
        }
    }
}