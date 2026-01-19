#include "mqtt_manager.h"
#include "state.h"
#include <ArduinoJson.h>

CapabilityManager capman;

static void mqttCallback(char* topic, byte* payload, unsigned int length){
  String t = String(topic);
  String p;
  for(unsigned int i=0;i<length;i++) p += (char)payload[i];

  LOG("[MQTT] RX %s payload=%s", t.c_str(), p.c_str());

  if(capman.routeSet(t, p, topicBase())){
    // опубликовать обновлённые состояния только для простоты: публикуем все
    capman.publishAllStates(mqtt, topicBase());
  } else {
    LOG("[MQTT] message not routed");
  }
}

void mqttInit(){
  mqtt.setServer(MQTT_HOST, MQTT_PORT);
  mqtt.setCallback(mqttCallback);
}

bool ensureMqttConnected(){
  if(mqtt.connected()) return true;
  LOG("[MQTT] Connecting as %s", deviceUID.c_str());
  if(!mqtt.connect(deviceUID.c_str(), MQTT_USER, MQTT_PASSWORD)){
    LOG("[MQTT] connect failed state=%d", mqtt.state());
    return false;
  }

  capman.subscribeAll(mqtt, topicBase());

  StaticJsonDocument<512> d;
  JsonArray arr = d.createNestedArray("capabilities");
  capman.describeAll(arr);
  String out; serializeJson(d,out);
  String descTopic = topicBase() + "/describe";
  mqtt.publish(descTopic.c_str(), out.c_str(), true);
  LOG("[MQTT] describe -> %s", descTopic.c_str());

  capman.publishAllStates(mqtt, topicBase());
  return true;
}

void mqttLoop(){
  if(!ensureMqttConnected()) return;
  mqtt.loop();
  if(millis() - lastStatePublish > STATE_INTERVAL){
    capman.publishAllStates(mqtt, topicBase());
    lastStatePublish = millis();
  }
}
