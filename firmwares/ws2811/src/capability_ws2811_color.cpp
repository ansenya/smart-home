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

  this->curInstance = instance;
  if (instance == "hsv") {
    JsonObject val = doc["value"];

    this->hue = val["h"] | hue;
    this->saturation = val["s"] | this->saturation;
    this->value = val["v"] | this->value;

    this->scene = "";
  } else if (instance == "temperature_k") {
    this->temperatureK = doc["value"] | temperatureK;
    temperatureToHsv(temperatureK);

    this->scene = "";
  } else if (instance == "scene") {
    this->scene = doc["value"] | "";
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
  auto &led = FastLEDManager::instance();
  if (!led.initialized()) return;

  unsigned long now = millis();

  if (curInstance == "scene") {
    updateScene(now);
    return;
  }

  if (curInstance == "hsv" || curInstance == "temperature_k") {
    updateHSV(now);
  }
}

static inline CHSV hsvToCHSV(float h, float s, float v) {
  uint8_t h8 = h * 255 / 360;
  uint8_t s8 = s * 255 / 100;
  uint8_t v8 = v * 255 / 100;
  return CHSV(h8, s8, v8);
}

void WS2811HSVCapability::applyHSV(float h, float s, float v) {
  auto &led = FastLEDManager::instance();
  led.setAll(hsvToCHSV(h, s, v));
  led.show();
}


void WS2811HSVCapability::updateScene(unsigned long now) {
  auto &led = FastLEDManager::instance();

  if (scene == "reading") {
    temperatureToHsv(3000);
    applyHSV(hue, saturation, value);
    return;
  }

  if (scene == "fantasy") {
    if (now - lastUpdate < 25) return;

    sceneHue += 0.7f;
    if (sceneHue >= 360) sceneHue -= 360;

    led.setAll(CHSV(sceneHue * 255 / 360, 255, 255));
    led.show();

    lastUpdate = now;
    return;
  }

  // party — зарезервировано
}

void WS2811HSVCapability::updateHSV(unsigned long now) {
  if (now - lastUpdate < updateInterval) return;

  if (curHue == hue &&
      curSaturation == saturation &&
      curValue == value) return;

  float dh = fmodf((hue-curHue)+540,360)-180;
  float ds = saturation-curSaturation;
  float dv = value-curValue;

  curHue        += dh * easing;
  curSaturation += ds * easing;
  curValue      += dv * easing;

  const float H_EPS = 0.5f;
  const float SV_EPS = 0.25f;

  if (fabsf(dh) < H_EPS) curHue = hue;
  if (fabsf(ds) < SV_EPS) curSaturation = saturation;
  if (fabsf(dv) < SV_EPS) curValue = value;

  while (curHue < 0) curHue += 360;
  while (curHue >= 360) curHue -= 360;

  applyHSV(curHue, curSaturation, curValue);

  lastUpdate = now;

  if (curHue == hue &&
      curSaturation == saturation &&
      curValue == value) {
    transitioning = false;
  }
}
