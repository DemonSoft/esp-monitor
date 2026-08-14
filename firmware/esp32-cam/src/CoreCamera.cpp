#include "CoreCamera.hpp"
#include "CaptureManager.hpp"

CameraManager camera;
CaptureManager capturer(camera);

// Настройка интервала
const unsigned long INTERVAL_N_SEC = 5; // Снимок каждые 5 секунд
unsigned long lastCaptureTime = 0;

// Флаги и таймеры неблокирующего старта
bool systemReady = false;
const unsigned long BOOT_DELAY_MS = 5000; // Пауза 5 секунд после включения
unsigned long bootTime = 0;

CameraManager::CameraManager() {}

bool CameraManager::begin() {
    // 1. Вспышка OFF
    pinMode(FLASH_GPIO_NUM, OUTPUT);
    digitalWrite(FLASH_GPIO_NUM, LOW);

// 1. Принудительный сброс шины SCCB (I2C)
    pinMode(SIOC_GPIO_NUM, OUTPUT);
    pinMode(SIOD_GPIO_NUM, OUTPUT);
    digitalWrite(SIOC_GPIO_NUM, HIGH);
    digitalWrite(SIOD_GPIO_NUM, HIGH);
    
    // Генерируем 9 тактовых импульсов для освобождения SDA
    for (int i = 0; i < 9; i++) {
        digitalWrite(SIOC_GPIO_NUM, LOW);
        delayMicroseconds(5);
        digitalWrite(SIOC_GPIO_NUM, HIGH);
        delayMicroseconds(5);
    }

    // 2. Аппаратный сброс матрицы через RESET/PWDN пины (если разведены)
    if (RESET_GPIO_NUM != -1) {
        pinMode(RESET_GPIO_NUM, OUTPUT);
        digitalWrite(RESET_GPIO_NUM, LOW);      // Удерживаем в сбросе
        delay(20);
        digitalWrite(RESET_GPIO_NUM, HIGH);     // Выводим из сброса
        delay(20);
    }
    
    if (PWDN_GPIO_NUM != -1) {
        pinMode(PWDN_GPIO_NUM, OUTPUT);
        digitalWrite(PWDN_GPIO_NUM, HIGH); // Sensor turn off
        delay(100);
        digitalWrite(PWDN_GPIO_NUM, LOW);  // Sensor turn on
        delay(100);
    }

    // 2. Пины
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
    
    config.pixel_format = PIXFORMAT_JPEG;

    if (psramFound()) {
        Serial.println("[Camera] PSRAM обнаружена! Используем PSRAM буферы.");
        config.frame_size = FRAMESIZE_SVGA;        // 800x600
        config.jpeg_quality = 12;                  // Качество
        config.fb_count = 1;                       // 1 буфер для стабильного одиночного захвата
        config.xclk_freq_hz = 20000000;            // 20 МГц
        config.grab_mode = CAMERA_GRAB_WHEN_EMPTY; // Захват строго по запросу
        config.fb_location = CAMERA_FB_IN_PSRAM;
    } else {
        Serial.println("[Camera] ВНИМАНИЕ: PSRAM не найдена! Используем DRAM.");
        config.frame_size = FRAMESIZE_QVGA;
        config.jpeg_quality = 15;
        config.fb_count = 1;
        config.xclk_freq_hz = 10000000;
        config.grab_mode = CAMERA_GRAB_WHEN_EMPTY;
        config.fb_location = CAMERA_FB_IN_DRAM;
    }

    esp_err_t err = esp_camera_init(&config);
    if (err != ESP_OK) {
        Serial.printf("[Camera] Ошибка инициализации камеры: 0x%x\n", err);
        return false;
    }

    // Прогрев сенсора после успешной инициализации
    for (int i = 0; i < 4; i++) {
        camera_fb_t* fb = esp_camera_fb_get();
        if (fb) esp_camera_fb_return(fb);
        delay(100);
    }    
    
    return true;
}

camera_fb_t* CameraManager::capture() {
    camera_fb_t* fb = esp_camera_fb_get();
    if (!fb) {
        Serial.println("Ошибка захвата кадра!");
        return nullptr;
    }
    return fb;
}

void CameraManager::release(camera_fb_t* fb) {
    if (fb) {
        esp_camera_fb_return(fb);
    }
}

void CameraManager::setFlash(bool state) {
    digitalWrite(FLASH_GPIO_NUM, state ? HIGH : LOW);
}

void setupCamera() {
    Serial.println("[Camera] Инициализация модуля камеры...");
    
    if (!camera.begin()) {
        Serial.println("[Camera] ОШИБКА: Не удалось запустить камеру!");
        return;
    }
    
    Serial.println("[Camera] Камера успешно инициализирована.");
    capturer.startContinuous(5); // UNCOMMENT: Снимок каждые 5 секунд
    // capturer.startBurstCount(5, 3); // UNCOMMENT: Сделать 3 снимка каждые 5 секунд
    // capturer.startBurstTime(5, 20); // UNCOMMENT: Снимать каждые 5 секунд в течение 20 секунд
}

void loopCamera() {
    capturer.update();
}