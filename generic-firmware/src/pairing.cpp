#include "pairing.h"
#include "state.h"
#include <HTTPClient.h>
#include <ArduinoJson.h>

bool confirmPair(){
  LOG("[PAIR] Sending confirm request");

  HTTPClient http;
  http.begin(PAIRING_URL);
  http.addHeader("Content-Type","application/json");

  StaticJsonDocument<256> d;
  d["code"] = prefs.getString("pair_code", "");
  d["device_uid"] = deviceUID;
  d["mac_address"] = macAddress;
  d["type"] = "devices.types.light";

  String body; serializeJson(d, body);
  LOG("[PAIR] Body: %s", body.c_str());

  int code = http.POST(body);
  String resp = http.getString();
  http.end();

  LOG("[PAIR] HTTP %d resp=%s", code, resp.c_str());

  if(code != 200 && code != 201) return false;

  StaticJsonDocument<256> r;
  DeserializationError err = deserializeJson(r, resp);
  if(err){
    LOG("[PAIR] JSON parse error");
    return false;
  }

  userID = r["user_id"].as<String>();
  deviceID = r["device_id"].as<String>();

  LOG("[PAIR] user_id=%s device_id=%s", userID.c_str(), deviceID.c_str());

  prefs.putString("user_id", userID);
  prefs.putString("device_id", deviceID);
  prefs.putBool("paired", true);
  prefs.remove("pair_code");
  return true;
}
