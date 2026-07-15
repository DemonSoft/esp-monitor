#ifndef CoreTime_h
#define CoreTime_h

#include <Arduino.h>

void setupTime();
void syncTime();
unsigned long getUnixTime();

#endif
