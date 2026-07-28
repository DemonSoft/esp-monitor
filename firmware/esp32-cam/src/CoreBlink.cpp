#include "CoreBlink.hpp"
#include "CoreCamera.h"
/*
  Внимание! Эти функции работают синхронно!!!
  Т.е. микроконтроллер будет замораживаться на заданный интервал!
  ИСПОЛЬЗОВАТЬ ТОЛЬКО ДЛЯ ОТЛАДКИ!
*/


void blink(String pattern) {
  int minus   = 500;
  int point   = 50;
  int space   = 200;


  for (uint i = 0; i < pattern.length(); i++)
  {
    char str = pattern[i];
    if (str == '-')  {
      durationBlink(minus);
    } else if (str == '.') {
      durationBlink(point);
    } else {
      delay(space);
    }
  }
}

void durationBlink(int duration)  {
  int waiting = 100;
  digitalWrite(FLASH_GPIO_NUM, HIGH);
  delay(duration);
  digitalWrite(FLASH_GPIO_NUM, LOW);
  delay(waiting);
}
