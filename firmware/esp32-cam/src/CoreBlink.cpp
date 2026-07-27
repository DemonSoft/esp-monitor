#include "CoreBlink.hpp"
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
  digitalWrite(LED_BUILTIN, HIGH);
  delay(duration);
  digitalWrite(LED_BUILTIN, LOW);
  delay(waiting);
}
