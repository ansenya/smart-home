#ifndef MQTT_MANAGER_H
#define MQTT_MANAGER_H

#include "capability_manager.h"

extern CapabilityManager capman;

void mqttInit();
void mqttLoop();

#endif // MQTT_MANAGER_H
