#pragma once
#include "capability.h"
#include <Arduino.h>

class WS2811HSVCapability : public Capability {
public:
  WS2811HSVCapability();

  String type() override { return "devices.capabilities.color_setting"; } // canonical
  void describe(JsonObject &o) override;
  void state(JsonObject &o) override;
  bool handleSet(const String &payload) override;

  void update();
private:
  // state
  String curInstance;

  // target values
  float hue = 0.0f;
  float saturation = 100.0f;
  float value = 100.0f;
  int temperatureK = 4500;
  String scene = "";
  
  // current (actual) values
  float curHue = 0.0f;
  float curSaturation = 100.0f;
  float curValue = 100.0f;

  int curTemperatureK = 4500;

  // transition parameters
  const float easing = 0.15f;             // k in cur += (target-cur)*k
  const unsigned long updateInterval = 30; // ms
  unsigned long lastUpdate = 0;
  bool transitioning = false;

  void hsvToRgb(float h, float s, float v, uint8_t &r, uint8_t &g, uint8_t &b);
  void temperatureToHsv(int kelvin);  
};
