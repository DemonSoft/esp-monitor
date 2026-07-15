#include "CoreTime.hpp"
#include <time.h>

static const char *ntpServer = "pool.ntp.org";
static const long utcOffsetSeconds = 0;

void setupTime() {
  configTime(utcOffsetSeconds, 0, ntpServer);
}

void syncTime() {
  for (int attempt = 0; attempt < 100; attempt++) {
    time_t now = time(nullptr);
    if (now > 1000) {
      return;
    }
    delay(1000);
  }
}

unsigned long getUnixTime() {
  time_t t = time(nullptr);
  return t > 0 ? (unsigned long)t : 0;
}
