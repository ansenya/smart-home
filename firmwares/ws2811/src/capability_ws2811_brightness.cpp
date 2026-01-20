#include "capability_ws2811_brightness.h"
#include "fastled_manager.h"
#include <ArduinoJson.h>

WS2811BrightnessCapability::WS2811BrightnessCapability() {
  // В конструкторе ничего не инициализируем — FastLEDManager управляет лентой
}

void WS2811BrightnessCapability::describe(JsonObject &o) {
  o["type"] = type();
  o["retrievable"] = true;
  o["reportable"] = true;

  JsonObject p = o.createNestedObject("parameters");
  p["instance"] = "brightness";
  p["random_access"] = true;
  p["unit"] = "unit.percent";

  JsonObject r = p.createNestedObject("range");
  r["max"] = 100;
  r["min"] = 1;
  r["precision"] = 1;

  JsonObject st = o.createNestedObject("state");
  st["instance"] = "brightness";
  st["value"] = map(brightness, 0, 255, 0, 100);;
}

void WS2811BrightnessCapability::state(JsonObject &o) {
  o["instance"] = "brightness";
  o["value"] = map(brightness, 0, 255, 0, 100); // 0-100%
}

bool WS2811BrightnessCapability::handleSet(const String &payload) {
  StaticJsonDocument<192> d;
  DeserializationError err = deserializeJson(d, payload);
  if (err) {
    LOG("[WS_BRIGHT] JSON parse error: %s", err.c_str());
    return false;
  }

  uint8_t target = brightness;
  bool changed = false;

  // Поддерживаем разные форматы входных данных
  if (d.containsKey("level")) { // Yandex format
    int percent = d["level"].as<int>();
    target = map(constrain(percent, 0, 100), 0, 100, 0, 255);
    changed = (target != brightness);
  } 
  else if (d.containsKey("value")) { // Our format
    int value = d["value"].as<int>();
    // Автоматически определяем диапазон: 0-100 или 0-255
    if (value <= 100) {
      target = map(value, 0, 100, 0, 255);
    } else {
      target = constrain(value, 0, 255);
    }
    changed = (target != brightness);
  }

  if (changed) {
    brightness = target;
    apply();
    
    // Сохраняем в Preferences (опционально)
    // prefs.putUChar("brightness", brightness);
    
    LOG("[WS_BRIGHT] Set brightness: %d (raw=%d)", 
        map(brightness, 0, 255, 0, 100), brightness);
  }

  return changed;
}

void WS2811BrightnessCapability::setTargetBrightness(uint8_t target){
  targetBrightness = target;
  LOG("[WS_BRIGHTNESS] target brightness: %d", targetBrightness);
}

void WS2811BrightnessCapability::update(){
  if(millis() - lastUpdate < updateInterval) return;
  
  if(brightness != targetBrightness){
    if(brightness < targetBrightness){
      brightness += fadeStep;
      if(brightness > targetBrightness) brightness = targetBrightness;
    } else {
      if(brightness > fadeStep){
          brightness -= fadeStep;
      } else {
          brightness = 0;
      }
      if(brightness < targetBrightness) brightness = targetBrightness;
    }
    
    FastLED.setBrightness(brightness);
    FastLED.show();
    LOG("[WS_BRIGHTNESS] current: %d, target: %d", brightness, targetBrightness);
  }
  
  lastUpdate = millis();
}

void WS2811BrightnessCapability::apply() {
  auto& ledmgr = FastLEDManager::instance();
  if (!ledmgr.initialized()) return;
  
  ledmgr.setBrightness(brightness);
  ledmgr.show();
}

// фабрика
Capability* createBrightness(){
  return new WS2811BrightnessCapability();
}
