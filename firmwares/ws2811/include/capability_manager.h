#ifndef CAPABILITY_MANAGER_H
#define CAPABILITY_MANAGER_H

#include <vector>
#include "capability.h"
#include <PubSubClient.h>

class CapabilityManager {
public:
  void registerCapability(Capability* c);
  void describeAll(JsonArray &arr);
  void publishAllStates(PubSubClient &mqtt, const String &base);
  bool routeSet(const String &topic, const String &payload, const String &base);
  void subscribeAll(PubSubClient &mqtt, const String &base);
  void publishHeartbeatState(PubSubClient &mqtt, const String &base);
  std::vector<Capability*> caps;
};

#endif // CAPABILITY_MANAGER_H
