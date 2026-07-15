#include "CoreCommon.hpp"
#include "CoreBlink.hpp"
#include "CoreButton.hpp"
#include "CoreConfig.hpp"

// ESP8266 button pin (GPIO 4 = D2 on Wemos D1 mini)
const int buttonPin = 4;
int buttonState = 0;
int lastButtonState = LOW;
unsigned long buttonPressedMillis = 0;
unsigned long previousPollMillis = 0;
const long pollInterval = 100;
const unsigned long longPressTime = 3000;

void setupButton() {
  pinMode(buttonPin, INPUT);
}

void loopButton() {
  if (!espState.access_exists) return;
  if (!checkPoolTimerInterval()) return;

  buttonState = digitalRead(buttonPin);

  if (buttonState == HIGH && lastButtonState == LOW) {
    buttonPressedMillis = millis();
  }

  if (buttonState == LOW && lastButtonState == HIGH) {
    unsigned long pressedDuration = millis() - buttonPressedMillis;
    if (pressedDuration >= longPressTime) {
      blink(".-.-.-");
      blink(".-.-.-");
      blink(".-.-.-");
      deleteConfigFile();
      delay(100);
      ESP.restart();
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
