#ifndef FASTLED_MANAGER_H
#define FASTLED_MANAGER_H

#include <FastLED.h>
#include "config.h"

class FastLEDManager {
public:
  static FastLEDManager& instance() {
    static FastLEDManager inst;
    return inst;
  }

  void begin(uint8_t pin = STRIP_PIN, uint16_t count = STRIP_COUNT);
  bool initialized() const { return _inited; }
  
  void setAll(CRGB color);
  void setBrightness(uint8_t b) { FastLED.setBrightness(b); }
  void show() { if (_inited) FastLED.show(); }
  CRGB* leds() { return _leds; }
  uint16_t count() const { return _count; }

private:
  FastLEDManager() = default;
  bool _inited = false;
  CRGB* _leds = nullptr;
  uint16_t _count = 0;
};

#endif // FASTLED_MANAGER_H