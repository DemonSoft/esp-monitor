#ifndef CORE_CAMERA_H
#define CORE_CAMERA_H

#include "esp_camera.h"
#include <Arduino.h>

// Определение пинов для популярной платы AI-THINKER ESP32-CAM
#define PWDN_GPIO_NUM  32
#define RESET_GPIO_NUM -1
#define XCLK_GPIO_NUM  0
#define SIOD_GPIO_NUM  26
#define SIOC_GPIO_NUM  27

#define Y9_GPIO_NUM    35
#define Y8_GPIO_NUM    34
#define Y7_GPIO_NUM    39
#define Y6_GPIO_NUM    36
#define Y5_GPIO_NUM    21
#define Y4_GPIO_NUM    19
#define Y3_GPIO_NUM    18
#define Y2_GPIO_NUM    5
#define VSYNC_GPIO_NUM 25
#define HREF_GPIO_NUM  23
#define PCLK_GPIO_NUM  22

#define FLASH_GPIO_NUM 4

class CameraManager {
public:
    CameraManager();
    bool begin();
    camera_fb_t* capture();
    void release(camera_fb_t* fb);
    
    // Вспомогательные функции для вспышки
    void setFlash(bool state);

private:
    camera_config_t config;
};

void setupCamera();
void loopCamera();

// void processPhoto(camera_fb_t* fb);


#endif // CORE_CAMERA_H