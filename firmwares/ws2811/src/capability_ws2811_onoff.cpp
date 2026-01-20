#include "capability_ws2811_onoff.h"
#include "fastled_manager.h"
#include "state.h"
#include <ArduinoJson.h>

WS2811OnOffCapability::WS2811OnOffCapability(WS2811BrightnessCapability* brightnessCap): brightnessCap(brightnessCap) {
  // В конструкторе ничего не инициализируем — FastLEDManager управляет лентой
}

void WS2811OnOffCapability::describe(JsonObject &o) {
  o["type"] = type();
  o["retrievable"] = true;
  o["reportable"] = true;

  JsonObject st = o.createNestedObject("state");
  st["instance"] = "on";
  st["value"] = getCurrentPower();
}

void WS2811OnOffCapability::state(JsonObject &o) {
  o["instance"] = "on";
  o["value"] = getCurrentPower();
}

bool WS2811OnOffCapability::handleSet(const String &payload) {
  StaticJsonDocument<192> d;
  DeserializationError err = deserializeJson(d, payload);
  if (err) {
    LOG("[WS_ONOFF] JSON parse error: %s", err.c_str());
    return false;
  }

  bool targetState = getCurrentPower(); // default to current state
  bool changed = false;

  // Поддерживаем несколько форматов payload
  if (d.containsKey("value")) {
    targetState = d["value"].as<bool>();
    changed = (targetState != getCurrentPower());
  }

  if (changed) {
    if (targetState) {
      uint8_t saved = prefs.getUChar("brightness", 255);
      brightnessCap->setTargetBrightness(saved);
    } else {
      prefs.putUChar("brightness", brightnessCap->getTargetBrightness());
      brightnessCap->setTargetBrightness(0);
    }
    // apply();
    LOG("[WS_ONOFF] Power %s", getCurrentPower() ? "ON" : "OFF");
  }

  return changed;
}

bool WS2811OnOffCapability::getCurrentPower() {
  return brightnessCap->getTargetBrightness() > 1;
}