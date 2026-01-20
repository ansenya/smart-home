#ifndef CAPABILITY_WS2811_HSV_H
#define CAPABILITY_WS2811_HSV_H
#include "capability.h"
#include <Arduino.h>

class WS2811HSVCapability : public Capability {
public:
  WS2811HSVCapability();

  String type() override { return "devices.capabilities.color_setting"; } // canonical
  void describe(JsonObject &o) override;
  void state(JsonObject &o) override;
  bool handleSet(const String &payload) override;

  void apply();
  void update();
private:
  // target values (what user requested)
  float hue = 0.0f;          // 0..360
  float saturation = 100.0f; // 0..100
  float value = 100.0f;      // 0..100

  // current (actual) values used to render LEDs (float for smooth steps)
  float curHue = 0.0f;
  float curSaturation = 100.0f;
  float curValue = 100.0f;

  // transition parameters
  const float easing = 0.15f;             // k in cur += (target-cur)*k
  const unsigned long updateInterval = 30; // ms
  unsigned long lastUpdate = 0;
  bool transitioning = false;

  void hsvToRgb(float h, float s, float v, uint8_t &r, uint8_t &g, uint8_t &b);
};

#endif // CAPABILITY_WS2811_HSV_H