#include "capability_ws2811_hsv.h"
#include "fastled_manager.h"
#include "state.h"
#include <ArduinoJson.h>

WS2811HSVCapability::WS2811HSVCapability() {
  curHue = hue;
  curSaturation = saturation;
  curValue = value;
  transitioning = false;
  lastUpdate = 0;
}

void WS2811HSVCapability::describe(JsonObject &o) {
  o["type"] = type();
  o["retrievable"] = true;
  o["reportable"] = true;

  JsonObject params = o.createNestedObject("parameters");
  params["color_model"] = "hsv";

  JsonObject cs = params.createNestedObject("color_scene");
  JsonArray scenes = cs.createNestedArray("scenes");

  const char* ids[] = {"party","fantasy","reading"};
  for(auto id : ids){
    JsonObject s = scenes.createNestedObject();
    s["id"] = id;
  }

  JsonObject st = o.createNestedObject("state");
  st["instance"] = "hsv";
  JsonObject hsv = st.createNestedObject("value");
  hsv["h"] = (int)round(hue);
  hsv["s"] = (int)round(saturation);
  hsv["v"] = (int)round(value);
}

void WS2811HSVCapability::state(JsonObject &o) {
  o["instance"] = "hsv";
  JsonObject hsv = o.createNestedObject("value");
  hsv["h"] = (int)round(hue);
  hsv["s"] = (int)round(saturation);
  hsv["v"] = (int)round(value);
}

bool WS2811HSVCapability::handleSet(const String &payload) {
  LOG("%s", payload);

  StaticJsonDocument<512> d;
  DeserializationError err = deserializeJson(d, payload);
  if (err) {
    LOG("[WS_HSV] JSON parse error");
    return false;
  }

  float newH = hue;
  float newS = saturation;
  float newV = value;
  bool changed = false;

  // helper to set from HSV object variant
  auto setFromObj = [&](JsonVariant obj) {
    bool localChanged = false;
    if (!obj.isNull()) {
      if (obj.containsKey("h")) {
        newH = constrain(obj["h"].as<float>(), 0.0f, 360.0f);
      } else if (obj.containsKey("hue")) {
        newH = constrain(obj["hue"].as<float>(), 0.0f, 360.0f);
      }
      if (obj.containsKey("s")) {
        newS = constrain(obj["s"].as<float>(), 0.0f, 100.0f);
      } else if (obj.containsKey("saturation")) {
        newS = constrain(obj["saturation"].as<float>(), 0.0f, 100.0f);
      }
      if (obj.containsKey("v")) {
        newV = constrain(obj["v"].as<float>(), 0.0f, 100.0f);
      } else if (obj.containsKey("value")) {
        newV = constrain(obj["value"].as<float>(), 0.0f, 100.0f);
      }
      localChanged = (round(newH) != round(hue)) || (round(newS) != round(saturation)) || (round(newV) != round(value));
    }
    return localChanged;
  };


  changed = setFromObj(d["value"]);

  // // 1) Yandex-style: {"instance":"hsv","value":{...}} or {"color":{...}}
  // if (d.containsKey("color") && d["color"].is<JsonObject>()) {
  //   changed = setFromObj(d["color"]);
  // }
  // else if (d.containsKey("instance") && d["instance"].is<const char*>()
  //          && String(d["instance"].as<const char>()) == "hsv"
  //          && d.containsKey("value") && d["value"].is<JsonObject>()) {
  //   changed = setFromObj(d["value"]);
  // }
  // // 2) Direct HSV: {"h":..,"s":..,"v":..} or {"hue":..,"saturation":..,"value":..}
  // else if (d.containsKey("h") && d.containsKey("s") && d.containsKey("v")) {
  //   changed = setFromObj(d);
  // }
  // else if (d.containsKey("hue") && d.containsKey("saturation") && d.containsKey("value")) {
  //   changed = setFromObj(d);
  // }
  // // 3) RGB components: {"r":..,"g":..,"b":..} or {"red":..,"green":..,"blue":..}
  // else if (d.containsKey("r") || d.containsKey("red")) {
  //   uint8_t r = d.containsKey("r") ? d["r"].as<uint8_t>() : d["red"].as<uint8_t>();
  //   uint8_t g = d.containsKey("g") ? d["g"].as<uint8_t>() : d["green"].as<uint8_t>();
  //   uint8_t b = d.containsKey("b") ? d["b"].as<uint8_t>() : d["blue"].as<uint8_t>();

  //   float rf = r / 255.0f, gf = g / 255.0f, bf = b / 255.0f;
  //   float mx = fmaxf(rf, fmaxf(gf, bf));
  //   float mn = fminf(rf, fminf(gf, bf));
  //   float delta = mx - mn;

  //   float hh = 0;
  //   if (delta != 0.0f) {
  //     if (mx == rf) hh = fmodf((gf - bf) / delta, 6.0f);
  //     else if (mx == gf) hh = ((bf - rf) / delta) + 2.0f;
  //     else hh = ((rf - gf) / delta) + 4.0f;
  //     hh *= 60.0f;
  //     if (hh < 0) hh += 360.0f;
  //   }
  //   float ss = (mx == 0.0f) ? 0.0f : (delta / mx) * 100.0f;
  //   float vv = mx * 100.0f;

  //   newH = constrain(hh, 0.0f, 360.0f);
  //   newS = constrain(ss, 0.0f, 100.0f);
  //   newV = constrain(vv, 0.0f, 100.0f);
  //   changed = (round(newH) != round(hue)) || (round(newS) != round(saturation)) || (round(newV) != round(value));
  // }
  // // 4) Integer RGB: {"value": 0xRRGGBB} (value is number)
  // else if (d.containsKey("value") && d["value"].is<uint32_t>()) {
  //   uint32_t rgb = d["value"].as<uint32_t>();
  //   uint8_t r = (rgb >> 16) & 0xFF;
  //   uint8_t g = (rgb >> 8) & 0xFF;
  //   uint8_t b = rgb & 0xFF;

  //   float rf = r / 255.0f, gf = g / 255.0f, bf = b / 255.0f;
  //   float mx = fmaxf(rf, fmaxf(gf, bf));
  //   float mn = fminf(rf, fminf(gf, bf));
  //   float delta = mx - mn;

  //   float hh = 0;
  //   if (delta != 0.0f) {
  //     if (mx == rf) hh = fmodf((gf - bf) / delta, 6.0f);
  //     else if (mx == gf) hh = ((bf - rf) / delta) + 2.0f;
  //     else hh = ((rf - gf) / delta) + 4.0f;
  //     hh *= 60.0f;
  //     if (hh < 0) hh += 360.0f;
  //   }
  //   float ss = (mx == 0.0f) ? 0.0f : (delta / mx) * 100.0f;
  //   float vv = mx * 100.0f;

  //   newH = constrain(hh, 0.0f, 360.0f);
  //   newS = constrain(ss, 0.0f, 100.0f);
  //   newV = constrain(vv, 0.0f, 100.0f);
  //   changed = (round(newH) != round(hue)) || (round(newS) != round(saturation)) || (round(newV) != round(value));
  // }
  // else {
  //   // Неизвестный формат
  //   LOG("[WS_HSV] Unknown payload format");
  //   return false;
  // }

  if (changed) {
    hue = newH;
    saturation = newS;
    value = newV;

    transitioning = true;
    // apply();
    LOG("[WS_HSV] Set color: H=%.0f S=%.0f V=%.0f", round(hue), round(saturation), round(value));
  } else {
    LOG("[WS_HSV] No change");
  }

  return changed;
}

void WS2811HSVCapability::apply() {
  auto &ledmgr = FastLEDManager::instance();
  if (!ledmgr.initialized()) return;

  uint8_t h = (uint8_t)(hue * 255.0f / 360.0f);
  uint8_t s = (uint8_t)(saturation * 255.0f / 100.0f);
  uint8_t v = (uint8_t)(value * 255.0f / 100.0f);

  CRGB color = CHSV(h, s, v);

  ledmgr.setAll(color);
  ledmgr.show();
}


void WS2811HSVCapability::update() {
  auto &ledmgr = FastLEDManager::instance();
  if (!ledmgr.initialized()) return;
  if (!transitioning) return;

  unsigned long now = millis();
  if (now - lastUpdate < updateInterval) return;
  lastUpdate = now;

  // Нормализуем разницу hue в диапазон (-180, 180]
  float dh = fmodf((hue - curHue) + 540.0f, 360.0f) - 180.0f;
  float ds = saturation - curSaturation;
  float dv = value - curValue;

  // Интерполяция (exponential easing)
  float stepH = dh * easing;
  float stepS = ds * easing;
  float stepV = dv * easing;

  curHue += stepH;
  curSaturation += stepS;
  curValue += stepV;

  // Если близко к цели — сразу ставим ровно
  const float H_EPS = 0.5f;   // градусы
  const float SV_EPS = 0.25f; // процентов

  // recompute residuals to check exact closeness
  float remH = fmodf((hue - curHue) + 540.0f, 360.0f) - 180.0f;
  float remS = saturation - curSaturation;
  float remV = value - curValue;

  if (fabsf(remH) < H_EPS) curHue = hue;
  if (fabsf(remS) < SV_EPS) curSaturation = saturation;
  if (fabsf(remV) < SV_EPS) curValue = value;

  // нормализуем curHue в 0..360
  while (curHue < 0.0f) curHue += 360.0f;
  while (curHue >= 360.0f) curHue -= 360.0f;

  // Преобразование для CHSV: 0..255
  uint8_t h8 = (uint8_t)roundf(curHue * 255.0f / 360.0f);
  uint8_t s8 = (uint8_t)roundf(curSaturation * 255.0f / 100.0f);
  uint8_t v8 = (uint8_t)roundf(curValue * 255.0f / 100.0f);

  CRGB color = CHSV(h8, s8, v8);
  ledmgr.setAll(color);
  ledmgr.show();

  // Завершаем переход, когда всё выставлено точно
  if (curHue == hue && curSaturation == saturation && curValue == value) {
    transitioning = false;
    LOG("[WS_HSV] Transition complete H=%.0f S=%.0f V=%.0f", round(curHue), round(curSaturation), round(curValue));
  }
}
