#ifndef STATE_H
#define STATE_H

#include <Preferences.h>
#include <WebServer.h>
#include <DNSServer.h>
#include <PubSubClient.h>
#include <WiFi.h>
#include <Arduino.h>
#include "capability_ws2811_onoff.h"
#include "capability_ws2811_brightness.h"
#include "capability_ws2811_hsv.h"

#include "config.h"

enum DeviceMode { MODE_PAIRING, MODE_WIFI, MODE_MQTT };

extern DeviceMode deviceMode;
extern Preferences prefs;
extern WebServer server;
extern DNSServer dnsServer;
extern WiFiClient net;
extern PubSubClient mqtt;

extern String deviceUID;
extern String macAddress;
extern String userID;
extern String deviceID;
extern bool powerState;
extern unsigned long lastStatePublish;

String getDeviceUID();
String getMac();
String topicBase();
void led(bool s);

extern WS2811OnOffCapability* onOffCap;
extern WS2811BrightnessCapability* brightnessCap;
extern WS2811HSVCapability* hsvCap;

#endif // STATE_H
