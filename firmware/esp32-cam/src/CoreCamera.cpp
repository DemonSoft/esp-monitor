#include "CoreCamera.h"
#include "CaptureManager.h"

CameraManager camera;
CaptureManager capturer(camera);

// Настройка интервала
const unsigned long INTERVAL_N_SEC = 5; // Снимок каждые 5 секунд
unsigned long lastCaptureTime = 0;

// Флаги и таймеры неблокирующего старта
bool systemReady = false;
const unsigned long BOOT_DELAY_MS = 5000; // Пауза 5 секунд после включения
unsigned long bootTime = 0;

void setupCamera() {    
    delay(1000);

    Serial.println("Инициализация камеры...");
    if (!camera.begin()) {
        Serial.println("Сбой! Проверьте подключения/питание.");
        while (true) { delay(1000); } // Зависаем в случае ошибки
    }
    Serial.println("Камера готова к работе!");
    bootTime = millis();
}

void loopCamera() {
    // 1. Ожидаем завершения стартовой паузы для стабилизации питания
    if (!systemReady) {
        if (millis() - bootTime >= BOOT_DELAY_MS) {
            systemReady = true;
            Serial.println("\n[SYSTEM] Питание стабильно. Запускаем ТЕСТ 1: BURST_COUNT");
            
            // Запуск ТЕСТА 1: 3 снимка с интервалом в 2 секунды
            capturer.startBurstCount(/*N sec=*/2, /*M photos=*/3);
        }
        return; // Пока плата "греется", ничего дальше не делаем
    }

    // 2. Обязательный постоянный вызов обновления состояния менеджера
    capturer.update();

    // 3. Логика переключения фаз теста (аналогично твоему варианту)
    static bool testPhase2Started = false;
    static bool testPhase3Started = false;
    static unsigned long phaseDelayTimer = 0;
    static bool waitingForPause = false;

    // Если текущий режим завершился и мы не в процессе паузы
    if (!capturer.isActive()) {
        
        // Переход к ФАЗЕ 2: Съемка по времени (BURST_TIME)
        if (!testPhase2Started) {
            if (!waitingForPause) {
                waitingForPause = true;
                phaseDelayTimer = millis();
                Serial.println("[TEST] Тест 1 завершен. Ожидание 3 секунды...");
            } 
            else if (millis() - phaseDelayTimer >= 3000) {
                waitingForPause = false;
                testPhase2Started = true;
                Serial.println("\n[TEST] Запуск ТЕСТА 2: BURST_TIME (съемка 7 секунд каждые 2 сек)");
                capturer.startBurstTime(/*N sec=*/2, /*K sec=*/7);
            }
        }
        // Переход к ФАЗЕ 3: Непрерывная съемка (CONTINUOUS)
        else if (testPhase2Started && !testPhase3Started) {
            if (!waitingForPause) {
                waitingForPause = true;
                phaseDelayTimer = millis();
                Serial.println("[TEST] Тест 2 завершен. Ожидание 3 секунды...");
            } 
            else if (millis() - phaseDelayTimer >= 3000) {
                waitingForPause = false;
                testPhase3Started = true;
                Serial.println("\n[TEST] Запуск ТЕСТА 3: CONTINUOUS (каждые 3 сек)");
                capturer.startContinuous(/*N sec=*/3);
            }
        }
    }
}

void processPhoto(camera_fb_t* fb) {
    // Тут будет твоя логика обработчика!
    // Например: отправка по Wi-Fi, запись на SD-карту и т.д.
    Serial.printf("Снимок успешно сделан! Размер файла: %u байт\n", fb->len);
}



/*
    CameraManager class implementation for managing the ESP32-CAM camera module.
*/

CameraManager::CameraManager() {
    config.ledc_channel = LEDC_CHANNEL_0;
    config.ledc_timer = LEDC_TIMER_0;
    config.pin_d0 = Y2_GPIO_NUM;
    config.pin_d1 = Y3_GPIO_NUM;
    config.pin_d2 = Y4_GPIO_NUM;
    config.pin_d3 = Y5_GPIO_NUM;
    config.pin_d4 = Y6_GPIO_NUM;
    config.pin_d5 = Y7_GPIO_NUM;
    config.pin_d6 = Y8_GPIO_NUM;
    config.pin_d7 = Y9_GPIO_NUM;
    config.pin_xclk = XCLK_GPIO_NUM;
    config.pin_pclk = PCLK_GPIO_NUM;
    config.pin_vsync = VSYNC_GPIO_NUM;
    config.pin_href = HREF_GPIO_NUM;
    config.pin_sccb_sda = SIOD_GPIO_NUM;
    config.pin_sccb_scl = SIOC_GPIO_NUM;
    config.pin_pwdn = PWDN_GPIO_NUM;
    config.pin_reset = RESET_GPIO_NUM;
    config.xclk_freq_hz = 20000000;
    config.pixel_format = PIXFORMAT_JPEG; // Формат готового снимка

    // Настройка качества и размера в зависимости от наличия PSRAM
    if (psramFound()) {
        config.frame_size = FRAMESIZE_UXGA; // 1600x1200 (можно уменьшить, например FRAMESIZE_VGA)
        config.jpeg_quality = 10;            // 0-63 (чем меньше число, тем выше качество)
        config.fb_count = 2;
    } else {
        config.frame_size = FRAMESIZE_SVGA; // 800x600
        config.jpeg_quality = 12;
        config.fb_count = 1;
    }
}

bool CameraManager::begin() {
    pinMode(FLASH_GPIO_NUM, OUTPUT);
    digitalWrite(FLASH_GPIO_NUM, LOW); // Отключаем вспышку по умолчанию

    pinMode(4, OUTPUT);
    digitalWrite(4, LOW);


    esp_err_t err = esp_camera_init(&config);
    if (err != ESP_OK) {
        Serial.printf("Ошибка инициализации камеры: 0x%x\n", err);
        return false;
    }
    return true;
}

camera_fb_t* CameraManager::capture() {
    // Получение фреймбуфера (кадра) с камеры
    camera_fb_t* fb = esp_camera_fb_get();
    if (!fb) {
        Serial.println("Ошибка захвата кадра!");
        return nullptr;
    }
    return fb;
}

void CameraManager::release(camera_fb_t* fb) {
    if (fb) {
        // ОБЯЗАТЕЛЬНО возвращаем буфер обратно драйверу, иначе память утечет!
        esp_camera_fb_return(fb);
    }
}

void CameraManager::setFlash(bool state) {
    digitalWrite(FLASH_GPIO_NUM, state ? HIGH : LOW);
}

