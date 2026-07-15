
#ifndef CoreCommon_h
#define CoreCommon_h

#include <Arduino.h>

typedef struct {
  bool mqtt_connected = false;
  bool access_exists  = false;
  String mac          = "unknown";
  String name         = "unknown";
} EspState;
extern EspState espState;


#endif
