#ifndef CAPABILITY_WS2811_ONOFF_H
#define CAPABILITY_WS2811_ONOFF_H
#include "capability.h"
#include "capability_ws2811_brightness.h"
#include <Arduino.h>

class WS2811OnOffCapability : public Capability {
public:
  WS2811OnOffCapability(WS2811BrightnessCapability* brightnessCap);
  String type() override { return "devices.capabilities.on_off"; }
  
  void describe(JsonObject &o) override;
  void state(JsonObject &o) override;
  bool handleSet(const String &payload) override;

private:
  bool getCurrentPower();
  WS2811BrightnessCapability* brightnessCap;
};

#endif // CAPABILITY_WS2811_ONOFF_H