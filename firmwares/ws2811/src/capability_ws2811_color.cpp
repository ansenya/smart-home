#include "capability_ws2811_color.h"
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

  JsonObject colorTemp = params.createNestedObject("temperature_k");
  colorTemp["min"] = 2700;
  colorTemp["max"] = 6500;

  JsonObject cs = params.createNestedObject("color_scene");
  JsonArray scenes = cs.createNestedArray("scenes");

  const char* ids[] = {"party","fantasy","reading"};
  for(auto id : ids){
    JsonObject s = scenes.createNestedObject();
    s["id"] = id;
  }

  JsonObject st = o.createNestedObject("state");
  this->state(st);
}

void WS2811HSVCapability::state(JsonObject &o) {
  o["instance"] = this->curInstance;
  if (this->curInstance == "hsv") {
    JsonObject hsv = o.createNestedObject("value");
    hsv["h"] = (int)round(hue);
    hsv["s"] = (int)round(saturation);
    hsv["v"] = (int)round(value);
  } else if (this->curInstance == "temperature_k") {
    o["value"] = this->temperatureK;
  } else if (this->curInstance == "scene") {
    o["value"] = this->scene;
  }
}

bool WS2811HSVCapability::handleSet(const String &payload) {
  LOG("%s", payload);

  StaticJsonDocument<512> doc;
  if (deserializeJson(doc, payload)) {
    LOG("[WS_HSV] JSON parse error");
    return false;
  }

  String instance = doc["instance"] | "hsv";
  JsonObject val = doc["value"];

  this->curInstance = instance;
  if (instance == "hsv") {
    this->hue = val["h"] | hue;
    this->saturation = val["s"] | this->saturation;
    this->value = val["v"] | this->value;

    this->scene = "";
  } else if (instance == "temperature_k") {
    this->temperatureK = doc["value"] | temperatureK;
    temperatureToHsv(temperatureK);

    this->scene = "";
  } else if (instance == "scene") {
    const char* sceneId = val["id"] | "";
    this->scene = String(sceneId);
  }

  return true;
}

void WS2811HSVCapability::temperatureToHsv(int kelvin) {
    // Convert color temperature to RGB
    float temp = kelvin / 100.0f;
    float r, g, b;
    
    // Red
    if (temp <= 66) {
        r = 255;
    } else {
        r = temp - 60;
        r = 329.698727446 * pow(r, -0.1332047592);
        r = constrain(r, 0, 255);
    }
    
    // Green
    if (temp <= 66) {
        g = temp;
        g = 99.4708025861 * log(g) - 161.1195681661;
    } else {
        g = temp - 60;
        g = 288.1221695283 * pow(g, -0.0755148492);
    }
    g = constrain(g, 0, 255);
    
    // Blue
    if (temp >= 66) {
        b = 255;
    } else if (temp <= 19) {
        b = 0;
    } else {
        b = temp - 10;
        b = 138.5177312231 * log(b) - 305.0447927307;
        b = constrain(b, 0, 255);
    }
    
    // Convert RGB to HSV
    float rf = r / 255.0f;
    float gf = g / 255.0f;
    float bf = b / 255.0f;
    
    float maxVal = max(max(rf, gf), bf);
    float minVal = min(min(rf, gf), bf);
    float delta = maxVal - minVal;
    
    this->value= maxVal * 100.0f;
    
    if (delta < 0.00001f) {
        this->saturation= 0;
        this->hue = 0;
        return;
    }
    
    this->saturation = (delta / maxVal) * 100.0f;
    
    if (rf >= maxVal) {
        this->hue = (gf - bf) / delta;
    } else if (gf >= maxVal) {
        this->hue = 2.0f + (bf - rf) / delta;
    } else {
        this->hue = 4.0f + (rf - gf) / delta;
    }
    
    this->hue *= 60.0f;
    if (this->hue < 0) this->hue += 360.0f;
}


void WS2811HSVCapability::update() {
  auto &ledmgr = FastLEDManager::instance();
  if (!ledmgr.initialized()) return;
  unsigned long now = millis();
  if (this->curInstance == "hsv" || this->curInstance == "temperature_k") {
    if (now - lastUpdate < updateInterval) return;

    if (this->curHue == this->hue && this->curSaturation == this->saturation && this->curValue == this->value) {
      return;
    }

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

  this->lastUpdate = now;
}
