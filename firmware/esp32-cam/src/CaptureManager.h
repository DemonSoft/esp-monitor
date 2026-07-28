#ifndef CAPTURE_MANAGER_H
#define CAPTURE_MANAGER_H

#include <Arduino.h>
#include "CoreCamera.h"

// Режимы работы
enum class CaptureMode {
    IDLE,             // Простой (ничего не делаем)
    CONTINUOUS,       // Съемка каждые N сек до остановки
    BURST_COUNT,      // Сделать M снимков каждые N сек
    BURST_TIME        // Снимать каждые N сек в течение K секунд
};

class CaptureManager {
public:
    CaptureManager(CameraManager& cameraRef);

    // Запуск режимов
    void startContinuous(uint32_t intervalSec);
    void startBurstCount(uint32_t intervalSec, uint32_t maxPhotos);
    void startBurstTime(uint32_t intervalSec, uint32_t durationSec);

    // Остановка съёмки вручную (по внешнему сигналу/команде)
    void stop();

    // Главный метод обработки (вызывается в loop())
    void update();

    // Проверка статуса
    bool isActive() const { return currentMode != CaptureMode::IDLE; }
    CaptureMode getMode() const { return currentMode; }

private:
    CameraManager& camera;
    CaptureMode currentMode = CaptureMode::IDLE;

    uint32_t intervalMs = 0;      // Интервал N (в мс)
    uint32_t maxPhotos = 0;       // Лимит снимков M
    uint32_t durationMs = 0;      // Лимит времени K (в мс)

    uint32_t lastCaptureTime = 0; // Время последнего снимка
    uint32_t modeStartTime = 0;   // Время старта режима
    uint32_t photosTaken = 0;     // Счетчик сделанных кадров

    void triggerCapture();
};

#endif // CAPTURE_MANAGER_H