#pragma once
#include <Arduino.h>
#include <ArduinoJson.h>

class Capability {
public:
  virtual String type() = 0;

  virtual void describe(JsonObject &o) = 0;
  virtual void state(JsonObject &o) = 0;
  virtual bool handleSet(const String &payload) = 0;

  virtual String setTopic(const String &base){
    return base + "/capabilities/" + type() + "/set";
  }

  virtual String stateTopic(const String &base){
    return base + "/capabilities/" + type() + "/state";
  }

  virtual ~Capability() {}
};
