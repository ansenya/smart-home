#include "captive.h"
#include "state.h"

static String captiveHTML(){
  return "<html><body><h3>SmartHome setup</h3>"
         "<form method='POST' action='/config'>"
         "SSID:<br><input name='ssid'><br>"
         "PASS:<br><input name='pass' type='password'><br>"
         "CODE:<br><input name='code'><br>"
         "<button>Save</button></form></body></html>";
}

void startCaptive(){
  LOG("[AP] Starting captive portal");

  deviceMode = MODE_PAIRING;

  WiFi.disconnect(true,true);
  WiFi.mode(WIFI_AP);
  WiFi.softAPConfig(AP_IP_ADDR,AP_IP_ADDR,IPAddress(255,255,255,0));

  // безопасный индекс начала подстроки
  int start = 0;
  if(deviceUID.length() > 6) start = (int)deviceUID.length() - 6;
  String ap = "SmartHome-" + deviceUID.substring(start);

  WiFi.softAP(ap.c_str(), NULL);

  LOG("[AP] AP started: %s ip=%s", ap.c_str(), WiFi.softAPIP().toString().c_str());

  dnsServer.start(DNS_PORT,"*",AP_IP_ADDR);
  LOG("[AP] DNS captive started");

  server.on("/",HTTP_GET,[](){
    LOG("[AP] GET /");
    server.send(200,"text/html",captiveHTML());
  });

  server.on("/hotspot-detect.html",HTTP_GET,[](){
    LOG("[AP] iOS captive check");
    server.send(200,"text/html",captiveHTML());
  });

  server.on("/generate_204",HTTP_GET,[](){
    LOG("[AP] Android captive check");
    server.sendHeader("Location","/",true);
    server.send(302);
  });

  server.on("/config",HTTP_POST,[](){
    String ssid = server.arg("ssid");
    String pass = server.arg("pass");
    String code = server.arg("code");

    LOG("[AP] Received config ssid='%s' code='%s'", ssid.c_str(), code.c_str());

    prefs.putString("ssid", ssid);
    prefs.putString("pass", pass);
    prefs.putString("pair_code", code);
    prefs.putBool("paired", false);

    server.send(200,"text/html","Saved. Rebooting...");
    delay(500);
    ESP.restart();
  });

  server.onNotFound([](){
    server.sendHeader("Location","/",true);
    server.send(302);
  });

  server.begin();
  LOG("[AP] Web server started");
}
