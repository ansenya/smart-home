#include "capability_ws2811_hsv.h"
#include "fastled_manager.h"
#include "state.h"
#include <ArduinoJson.h>

WS2811HSVCapability::WS2811HSVCapability() {
}

void WS2811HSVCapability::describe(JsonObject &o) {
  o["type"] = type();
  o["retrievable"] = true;
  o["reportable"] = true;

  JsonObject params = o.createNestedObject("parameters");
  params["color_model"] = "hsv";

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
  StaticJsonDocument<384> d;
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

  // 1) Yandex-style: {"instance":"hsv","value":{...}} or {"color":{...}}
  if (d.containsKey("color") && d["color"].is<JsonObject>()) {
    changed = setFromObj(d["color"]);
  }
  else if (d.containsKey("instance") && d["instance"].is<const char*>()
           && String(d["instance"].as<const char>()) == "hsv"
           && d.containsKey("value") && d["value"].is<JsonObject>()) {
    changed = setFromObj(d["value"]);
  }
  // 2) Direct HSV: {"h":..,"s":..,"v":..} or {"hue":..,"saturation":..,"value":..}
  else if (d.containsKey("h") && d.containsKey("s") && d.containsKey("v")) {
    changed = setFromObj(d);
  }
  else if (d.containsKey("hue") && d.containsKey("saturation") && d.containsKey("value")) {
    changed = setFromObj(d);
  }
  // 3) RGB components: {"r":..,"g":..,"b":..} or {"red":..,"green":..,"blue":..}
  else if (d.containsKey("r") || d.containsKey("red")) {
    uint8_t r = d.containsKey("r") ? d["r"].as<uint8_t>() : d["red"].as<uint8_t>();
    uint8_t g = d.containsKey("g") ? d["g"].as<uint8_t>() : d["green"].as<uint8_t>();
    uint8_t b = d.containsKey("b") ? d["b"].as<uint8_t>() : d["blue"].as<uint8_t>();

    float rf = r / 255.0f, gf = g / 255.0f, bf = b / 255.0f;
    float mx = fmaxf(rf, fmaxf(gf, bf));
    float mn = fminf(rf, fminf(gf, bf));
    float delta = mx - mn;

    float hh = 0;
    if (delta != 0.0f) {
      if (mx == rf) hh = fmodf((gf - bf) / delta, 6.0f);
      else if (mx == gf) hh = ((bf - rf) / delta) + 2.0f;
      else hh = ((rf - gf) / delta) + 4.0f;
      hh *= 60.0f;
      if (hh < 0) hh += 360.0f;
    }
    float ss = (mx == 0.0f) ? 0.0f : (delta / mx) * 100.0f;
    float vv = mx * 100.0f;

    newH = constrain(hh, 0.0f, 360.0f);
    newS = constrain(ss, 0.0f, 100.0f);
    newV = constrain(vv, 0.0f, 100.0f);
    changed = (round(newH) != round(hue)) || (round(newS) != round(saturation)) || (round(newV) != round(value));
  }
  // 4) Integer RGB: {"value": 0xRRGGBB} (value is number)
  else if (d.containsKey("value") && d["value"].is<uint32_t>()) {
    uint32_t rgb = d["value"].as<uint32_t>();
    uint8_t r = (rgb >> 16) & 0xFF;
    uint8_t g = (rgb >> 8) & 0xFF;
    uint8_t b = rgb & 0xFF;

    float rf = r / 255.0f, gf = g / 255.0f, bf = b / 255.0f;
    float mx = fmaxf(rf, fmaxf(gf, bf));
    float mn = fminf(rf, fminf(gf, bf));
    float delta = mx - mn;

    float hh = 0;
    if (delta != 0.0f) {
      if (mx == rf) hh = fmodf((gf - bf) / delta, 6.0f);
      else if (mx == gf) hh = ((bf - rf) / delta) + 2.0f;
      else hh = ((rf - gf) / delta) + 4.0f;
      hh *= 60.0f;
      if (hh < 0) hh += 360.0f;
    }
    float ss = (mx == 0.0f) ? 0.0f : (delta / mx) * 100.0f;
    float vv = mx * 100.0f;

    newH = constrain(hh, 0.0f, 360.0f);
    newS = constrain(ss, 0.0f, 100.0f);
    newV = constrain(vv, 0.0f, 100.0f);
    changed = (round(newH) != round(hue)) || (round(newS) != round(saturation)) || (round(newV) != round(value));
  }
  else {
    // Неизвестный формат
    LOG("[WS_HSV] Unknown payload format");
    return false;
  }

  if (changed) {
    hue = newH;
    saturation = newS;
    value = newV;

    apply();
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
