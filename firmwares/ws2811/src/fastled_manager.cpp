#include "fastled_manager.h"

void FastLEDManager::begin(uint8_t pin, uint16_t count) {
  if (_inited) return;
  
  _count = count;
  _leds = (CRGB*)malloc(sizeof(CRGB) * _count);
  
  // Важно: FastLED требует константный пин на этапе компиляции
  FastLED.addLeds<WS2811, STRIP_PIN, BRG>(_leds, _count);
  
  FastLED.setBrightness(0);
  fill_solid(_leds, _count, CRGB::Black);
  FastLED.show();
  
  _inited = true;
  LOG("[FASTLED] Initialized %d LEDs on pin %d", _count, STRIP_PIN);
}

void FastLEDManager::setAll(CRGB color) {
   if (!_inited) return;
  fill_solid(_leds, _count, color);
}