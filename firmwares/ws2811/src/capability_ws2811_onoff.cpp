#include "capability_ws2811_onoff.h"
#include "fastled_manager.h"
#include <ArduinoJson.h>

WS2811OnOffCapability::WS2811OnOffCapability() {
  // В конструкторе ничего не инициализируем — FastLEDManager управляет лентой
}

void WS2811OnOffCapability::describe(JsonObject &o) {
  o["type"] = type();
  o["retrievable"] = true;
  o["reportable"] = true;

  JsonObject st = o.createNestedObject("state");
  st["instance"] = "on";
  st["value"] = power;
}

void WS2811OnOffCapability::state(JsonObject &o) {
  o["instance"] = "on";
  o["value"] = power;
}

bool WS2811OnOffCapability::handleSet(const String &payload) {
  StaticJsonDocument<192> d;
  DeserializationError err = deserializeJson(d, payload);
  if (err) {
    LOG("[WS_ONOFF] JSON parse error: %s", err.c_str());
    return false;
  }

  bool targetState = power; // default to current state
  bool changed = false;

  // Поддерживаем несколько форматов payload
  if (d.containsKey("value")) {
    targetState = d["value"].as<bool>();
    changed = (targetState != power);
  }

  if (changed) {
    power = targetState;
    apply();
    LOG("[WS_ONOFF] Power %s", power ? "ON" : "OFF");
  }

  return changed;
}

void WS2811OnOffCapability::apply() {
  auto& ledmgr = FastLEDManager::instance();
  if (!ledmgr.initialized()) {
    LOG("[FASTLED] NOT INITIALIZED");
    return;
  };

  if (power) {
    // Базовый цвет при включении (можно заменить на сохранённый цвет/эффект)
    ledmgr.setAll(CRGB::White); // или CRGB::Red для теста
  } else {
    ledmgr.setAll(CRGB::Black);
  }
  ledmgr.show();
}

// фабрика
Capability* createOnOff(){
  // используем LED_PIN по умолчанию; можно передавать другой пин при необходимости
  return new WS2811OnOffCapability();
}
