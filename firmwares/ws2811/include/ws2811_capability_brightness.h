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

private:
  uint8_t brightness = 255; // 0-255, default max
  void apply();
};

Capability* createBrightness();

#endif // CAPABILITY_WS2811_BRIGHTNESS_H