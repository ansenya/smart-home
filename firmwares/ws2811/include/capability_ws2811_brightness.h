#ifndef CAPABILITY_WS2811_BRIGHTNESS_H
#define CAPABILITY_WS2811_BRIGHTNESS_H

#include "capability.h"
#include <Arduino.h>

class WS2811BrightnessCapability : public Capability {
public:
  WS2811BrightnessCapability();
  String type() override { return "devices.capabilities.range"; }
  
  void describe(JsonObject &o) override;
  void state(JsonObject &o) override;
  bool handleSet(const String &payload) override;
  void update(); // Для плавных переходов
  void setTargetBrightness(uint8_t target);

private:
  uint8_t brightness = 255; // 0-255, default max
  uint8_t targetBrightness;
  unsigned long lastUpdate;
  const uint8_t fadeStep = 2; // Шаг изменения яркости
  const uint8_t updateInterval = 20; // Интервал в мс
  void apply();
};

Capability* createBrightness();

#endif // CAPABILITY_WS2811_BRIGHTNESS_H