#include "CoreCamera.h"


CameraManager camera;

// Настройка интервала
const unsigned long INTERVAL_N_SEC = 5; // Снимок каждые 5 секунд
unsigned long lastCaptureTime = 0;


void setupCamera() {
    // 1. Принудительно отключаем вспышку на GPIO4 ДО инициализации чего-либо еще
    pinMode(4, OUTPUT);
    digitalWrite(4, LOW);
    
    delay(1000);

    Serial.println("Инициализация камеры...");
    if (!camera.begin()) {
        Serial.println("Сбой! Проверьте подключения/питание.");
        while (true) { delay(1000); } // Зависаем в случае ошибки
    }
    Serial.println("Камера готова к работе!");
}

void processPhoto(camera_fb_t* fb);

void loopCamera() {
    unsigned long currentMillis = millis();

    // Проверяем, прошло ли N секунд (без использования блокирующего delay)
    if (currentMillis - lastCaptureTime >= (INTERVAL_N_SEC * 1000)) {
        lastCaptureTime = currentMillis;

        Serial.println("Делаем снимок...");
        
        // 1. Делаем снимок
        camera_fb_t* fb = camera.capture();

        // 2. Если фото захвачено успешно, обрабатываем
        if (fb) {
            processPhoto(fb);

            // 3. ОБЯЗАТЕЛЬНО освобождаем память!
            camera.release(fb);
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
    // pinMode(FLASH_GPIO_NUM, OUTPUT);
    // digitalWrite(FLASH_GPIO_NUM, LOW); // Отключаем вспышку по умолчанию

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

