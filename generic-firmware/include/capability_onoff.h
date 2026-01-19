#ifndef CAPABILITY_ONOFF_H
#define CAPABILITY_ONOFF_H

#include "capability.h"
#include <Arduino.h>

class OnOffCapability : public Capability {
public:
  explicit OnOffCapability(uint8_t pin);

  String type() override { return "devices.capabilities.on_off"; }

  void describe(JsonObject &o) override;
  void state(JsonObject &o) override;
  bool handleSet(const String &payload) override;

private:
  uint8_t pin;
  bool value;
};

// фабричная функция — реализована в cpp
Capability* createOnOff();

#endif // CAPABILITY_ONOFF_H
