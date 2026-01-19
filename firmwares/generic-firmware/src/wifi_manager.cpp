#include "wifi_manager.h"
#include "state.h"

bool connectWiFi(){
  String ssid = prefs.getString("ssid", "");
  String pass = prefs.getString("pass", "");

  if(ssid == ""){
    LOG("[WIFI] No SSID saved");
    return false;
  }

  LOG("[WIFI] Connecting to '%s'", ssid.c_str());

  WiFi.mode(WIFI_STA);
  WiFi.begin(ssid.c_str(), pass.c_str());

  unsigned long t = millis();
  while(WiFi.status() != WL_CONNECTED && millis() - t < WIFI_TIMEOUT){
    delay(250);
    Serial.print(".");
  }
  Serial.println();

  if(WiFi.status() == WL_CONNECTED){
    LOG("[WIFI] Connected IP=%s", WiFi.localIP().toString().c_str());
    deviceMode = MODE_WIFI;
    return true;
  }

  LOG("[WIFI] Connection failed");
  return false;
}
