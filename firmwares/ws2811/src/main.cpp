#include "state.h"
#include "wifi_manager.h"
#include "pairing.h"
#include "captive.h"
#include "mqtt_manager.h"
#include "reset.h"
#include "fastled_manager.h"
#include "capability_ws2811_onoff.h"
#include "capability_ws2811_brightness.h"
#include "capability_ws2811_hsv.h"

void setup() {
  Serial.begin(115200);
  pinMode(LED_PIN, OUTPUT);

  prefs.begin("device", false);

  deviceUID = getDeviceUID();
  macAddress = getMac();

  LOG("==================================");
  LOG("UID: %s", deviceUID.c_str());
  LOG("MAC: %s", macAddress.c_str());
  LOG("==================================");

  checkReset();

  // Инициализируем менеджер ленты ОДИН РАЗ
  FastLEDManager::instance().begin(STRIP_PIN, STRIP_COUNT);

  // Register caps
  onOffCap = new WS2811OnOffCapability();
  brightnessCap = new WS2811BrightnessCapability();
  hsvCap = new WS2811HSVCapability();
  capman.registerCapability(onOffCap);
  capman.registerCapability(brightnessCap);
  capman.registerCapability(hsvCap);

  mqttInit();

  if (prefs.getBool("paired", false)) {
    userID = prefs.getString("user_id", "");
    deviceID = prefs.getString("device_id", "");
    LOG("[INIT] paired user=%s device=%s", userID.c_str(), deviceID.c_str());
    if (connectWiFi()) {
      // mqttLoop will handle connection
    }
  }

  if (prefs.getString("pair_code", "") != "") {
    LOG("[INIT] Found pairing code");
    if (connectWiFi() && confirmPair()) {
      // Pairing completed
    }
  }

  if (!WiFi.isConnected()) {
    LOG("[INIT] Entering AP pairing mode");
    startCaptive();
  }
}

void loop() {
  if (deviceMode == MODE_PAIRING) {
    dnsServer.processNextRequest();
    server.handleClient();
    led(millis() % 1000 < 500);
    return;
  }

  if(brightnessCap) brightnessCap->update();
  mqttLoop();
}