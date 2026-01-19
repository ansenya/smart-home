#include "capability_onoff.h"
#include "state.h"
#include <ArduinoJson.h>

OnOffCapability::OnOffCapability(uint8_t p)
  : pin(p), value(false)
{
  pinMode(pin, OUTPUT);
  digitalWrite(pin, LOW);
}

void OnOffCapability::describe(JsonObject &o) {
  o["type"] = type(); // "devices.capabilities.on_off"
  o["retrievable"] = true;
  o["reportable"] = false;

  JsonObject stateObj = o.createNestedObject("state");
  this->state(stateObj);
}

void OnOffCapability::state(JsonObject &o){
  o["instance"] = "on";
  o["value"] = value;
}

bool OnOffCapability::handleSet(const String &payload){
  StaticJsonDocument<128> d;
  DeserializationError err = deserializeJson(d, payload);
  if(err) {
    LOG("[OnOff] handleSet: json parse error");
    return false;
  }

  if(!d.containsKey("value")){
    LOG("[OnOff] handleSet: no 'value' key");
    return false;
  }

  // value может быть boolean или числом/строкой
  // пробуем безопасно получить boolean
  if(d["value"].is<bool>()){
    value = d["value"].as<bool>();
  } else if(d["value"].is<int>()){
    value = (d["value"].as<int>() != 0);
  } else {
    // как текст: "true"/"false"/"1"/"0"
    String s;
    serializeJson(d["value"], s);
    s.trim();
    if(s == "true" || s == "\"true\"" || s == "1") value = true;
    else value = false;
  }

  digitalWrite(pin, value ? HIGH : LOW);
  LOG("[OnOff] set -> %s", value ? "true" : "false");

  return true;
}

// фабрика
Capability* createOnOff(){
  // используем LED_PIN по умолчанию; можно передавать другой пин при необходимости
  return new OnOffCapability(LED_PIN);
}
