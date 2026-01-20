#include "capability_manager.h"
#include <ArduinoJson.h>
#include "state.h"

void CapabilityManager::registerCapability(Capability* c){
  caps.push_back(c);
}

void CapabilityManager::describeAll(JsonArray &arr){
  LOG("describing");
  for(auto c: caps){
    JsonObject o = arr.createNestedObject();
    c->describe(o);
  }
}

void CapabilityManager::publishHeartbeatState(PubSubClient &mqtt, const String &base){
  String t = base + "/state";
  mqtt.publish(t.c_str(), "{\"status\":\"online\"}", true);
  LOG("[MQTT] heartbeat -> %s", t.c_str());
}

void CapabilityManager::publishAllStates(PubSubClient &mqtt, const String &base){
  publishHeartbeatState(mqtt, base);

  for(auto c: caps){
    StaticJsonDocument<2048> d;
    JsonObject root = d.to<JsonObject>();
    c->state(root);

    String out;
    serializeJson(d, out);

    String topic = c->stateTopic(base);
    mqtt.publish(topic.c_str(), out.c_str(), true);

    LOG("[MQTT] publish state %s -> %s", topic.c_str(), out.c_str());
  }
}

bool CapabilityManager::routeSet(const String &topic, const String &payload, const String &base){
  for(auto c: caps){
    String expected = c->setTopic(base);
    if(topic == expected){
      bool ok = c->handleSet(payload);
      LOG("[CAPMAN] routed to %s, ok=%d", c->type().c_str(), ok);
      return ok;
    }
  }
  return false;
}

void CapabilityManager::subscribeAll(PubSubClient &mqtt, const String &base){
  for(auto c: caps){
    String sub = c->setTopic(base);
    mqtt.subscribe(sub.c_str());
    LOG("[MQTT] Subscribed %s", sub.c_str());
  }
}
