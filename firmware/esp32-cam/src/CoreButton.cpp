#include "CoreCommon.hpp"
#include "CoreButton.hpp"
#include "CoreConfig.hpp"

// ESP32 button pin example (GPIO 13 for ESP32-CAM)
const int buttonPin = 13;
int buttonState = HIGH;
int lastButtonState = HIGH;
bool buttonWasPressed = false;
unsigned long buttonPressedMillis = 0;
unsigned long previousPollMillis = 0;
const long pollInterval = 100;
const unsigned long longPressTime = 3000;

void setupButton() {
  pinMode(buttonPin, INPUT_PULLUP);
}

void loopButton() {
  if (!espState.access_exists) return;
  if (!checkPoolTimerInterval()) return;

  buttonState = digitalRead(buttonPin);

  // Нажали кнопку (переход в LOW)
  if (buttonState == LOW && lastButtonState == HIGH) {
    buttonPressedMillis = millis();
    buttonWasPressed = true;
  }

  // Отпустили кнопку (переход в HIGH)
  if (buttonState == HIGH && lastButtonState == LOW) {
    if (buttonWasPressed) {
      unsigned long pressedDuration = millis() - buttonPressedMillis;
      if (pressedDuration >= longPressTime) {
        Serial.println("[Button] Сброс настроек по удержанию кнопки!");
        deleteConfigFile();
        delay(100);
        ESP.restart();
      }
      buttonWasPressed = false;
    }
  }

  lastButtonState = buttonState;
}

bool checkPoolTimerInterval() {
  unsigned long currentMillis = millis();
  if (currentMillis - previousPollMillis < pollInterval) return false;
  previousPollMillis = currentMillis;
  return true;
}
