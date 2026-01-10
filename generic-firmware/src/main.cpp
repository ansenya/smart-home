#include <WiFi.h>
#include <WebServer.h>
#include <Preferences.h>
#include <HTTPClient.h>
#include <ArduinoJson.h>
#include <DNSServer.h>
#include <PubSubClient.h>

#define LOG(x, ...) Serial.printf(x "\n", ##__VA_ARGS__)

// ================= CONFIG =================
#define AP_IP_ADDR IPAddress(192,168,4,1)
#define DNS_PORT 53
#define PAIRING_URL "https://api.smarthome.hipahopa.ru/devices/pairing/confirm"

#define MQTT_HOST "hipahopa.ru"
#define MQTT_PORT 1883
#define MQTT_USER "user"
#define MQTT_PASSWORD "superdifficultpassword"

#define RESET_PIN 0
#define LED_PIN 2

#define WIFI_TIMEOUT 20000
#define STATE_INTERVAL 5000
// ========================================

enum DeviceMode { MODE_PAIRING, MODE_WIFI, MODE_MQTT };
DeviceMode deviceMode = MODE_PAIRING;

Preferences prefs;
WebServer server(80);
DNSServer dnsServer;
WiFiClient net;
PubSubClient mqtt(net);

String deviceUID, macAddress, userID, deviceID;
bool powerState = false;
unsigned long lastStatePublish = 0;

// ================= UTIL =================
String getDeviceUID() {
  uint64_t id = ESP.getEfuseMac();
  char buf[32];
  sprintf(buf, "esp32-%04X%08X", (uint16_t)(id>>32), (uint32_t)id);
  return buf;
}

String getMac() {
  uint8_t mac[6];
  esp_read_mac(mac, ESP_MAC_WIFI_STA);
  char buf[32];
  sprintf(buf,"%02X:%02X:%02X:%02X:%02X:%02X",mac[0],mac[1],mac[2],mac[3],mac[4],mac[5]);
  return buf;
}

String topicBase() {
  return userID + "/" + deviceID;
}

void led(bool s){ digitalWrite(LED_PIN, s?HIGH:LOW); }

// ================= CAPTIVE =================
String captiveHTML(){
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

  String ap = "SmartHome-" + deviceUID.substring(deviceUID.length()-6);
  WiFi.softAP(ap.c_str(), NULL);

  LOG("[AP] AP started: %s ip=%s",
      ap.c_str(), WiFi.softAPIP().toString().c_str());

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
    String ssid=server.arg("ssid");
    String pass=server.arg("pass");
    String code=server.arg("code");

    LOG("[AP] Received config ssid='%s' code='%s'", ssid.c_str(), code.c_str());

    prefs.putString("ssid",ssid);
    prefs.putString("pass",pass);
    prefs.putString("pair_code",code);
    prefs.putBool("paired",false);

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

// ================= WIFI =================
bool connectWiFi(){
  String ssid=prefs.getString("ssid","");
  String pass=prefs.getString("pass","");

  if(ssid==""){
    LOG("[WIFI] No SSID saved");
    return false;
  }

  LOG("[WIFI] Connecting to '%s'", ssid.c_str());

  WiFi.mode(WIFI_STA);
  WiFi.begin(ssid.c_str(),pass.c_str());

  unsigned long t=millis();
  while(WiFi.status()!=WL_CONNECTED && millis()-t<WIFI_TIMEOUT){
    delay(250);
    Serial.print(".");
  }
  Serial.println();

  if(WiFi.status()==WL_CONNECTED){
    LOG("[WIFI] Connected IP=%s", WiFi.localIP().toString().c_str());
    return true;
  }

  LOG("[WIFI] Connection failed");
  return false;
}

// ================= PAIR =================
bool confirmPair(){
  LOG("[PAIR] Sending confirm request");

  HTTPClient http;
  http.begin(PAIRING_URL);
  http.addHeader("Content-Type","application/json");

  StaticJsonDocument<256> d;
  d["code"]=prefs.getString("pair_code","");
  d["device_uid"]=deviceUID;
  d["mac_address"]=macAddress;
  d["type"]="devices.types.light";

  String body; serializeJson(d,body);
  LOG("[PAIR] Body: %s", body.c_str());

  int code=http.POST(body);
  String resp=http.getString();
  http.end();

  LOG("[PAIR] HTTP %d resp=%s", code, resp.c_str());

  if(code!=200 && code!=201) return false;

  StaticJsonDocument<256> r;
  if(deserializeJson(r,resp)){
    LOG("[PAIR] JSON parse error");
    return false;
  }

  userID=r["user_id"].as<String>();
  deviceID=r["device_id"].as<String>();

  LOG("[PAIR] user_id=%s device_id=%s", userID.c_str(), deviceID.c_str());

  prefs.putString("user_id",userID);
  prefs.putString("device_id",deviceID);
  prefs.putBool("paired",true);
  prefs.remove("pair_code");

  return true;
}

// ================= MQTT =================
void publishDescribe(){
  StaticJsonDocument<256> d;
  auto caps=d.createNestedArray("capabilities");
  auto c=caps.createNestedObject();
  c["type"]="devices.capabilities.on_off";
  c["retrievable"]=true;
  c["reportable"]=true;

  String out; serializeJson(d,out);

  String t=topicBase()+"/describe";
  mqtt.publish(t.c_str(),out.c_str(),true);
  LOG("[MQTT] describe -> %s", t.c_str());
}

void publishState(){
  String t=topicBase()+"/state";
  mqtt.publish(t.c_str(),"{\"status\":\"online\"}",true);
  LOG("[MQTT] state online -> %s", t.c_str());
}

void publishOnOff(){
  StaticJsonDocument<64> d;
  d["value"]=powerState;
  String out; serializeJson(d,out);

  String t=topicBase()+"/capabilities/devices.capabilities.on_off/state";
  mqtt.publish(t.c_str(),out.c_str(),true);
  LOG("[MQTT] on_off=%s -> %s", powerState?"true":"false", t.c_str());
}

void mqttCb(char* t, byte* p, unsigned int l){
  String msg;
  for(uint i=0;i<l;i++) msg+=(char)p[i];

  LOG("[MQTT] RX %s payload=%s", t, msg.c_str());

  StaticJsonDocument<64> d;
  if(!deserializeJson(d,msg)){
    bool v=d["value"];
    powerState=v;
    led(v);
    LOG("[ACTION] powerState=%s", v?"true":"false");
    publishOnOff();
  }
}

bool connectMQTT(){
  mqtt.setServer(MQTT_HOST,MQTT_PORT);
  mqtt.setCallback(mqttCb);

  LOG("[MQTT] Connecting as %s", deviceUID.c_str());

  if(!mqtt.connect(deviceUID.c_str(),MQTT_USER,MQTT_PASSWORD)){
    LOG("[MQTT] connect failed state=%d", mqtt.state());
    return false;
  }

  String sub=topicBase()+"/capabilities/devices.capabilities.on_off/set";
  mqtt.subscribe(sub.c_str());

  LOG("[MQTT] Subscribed %s", sub.c_str());

  publishDescribe();
  publishState();
  publishOnOff();

  deviceMode=MODE_MQTT;
  return true;
}

// ================= RESET =================
void checkReset(){
  pinMode(RESET_PIN,INPUT_PULLUP);
  if(digitalRead(RESET_PIN)==LOW){
    LOG("[RESET] Button pressed");
    delay(3000);
    if(digitalRead(RESET_PIN)==LOW){
      LOG("[RESET] Factory reset!");
      prefs.clear();
      ESP.restart();
    }
  }
}

// ================= SETUP =================
void setup(){
  Serial.begin(115200);
  pinMode(LED_PIN,OUTPUT);

  prefs.begin("device",false);

  deviceUID=getDeviceUID();
  macAddress=getMac();

  LOG("==================================");
  LOG("UID: %s", deviceUID.c_str());
  LOG("MAC: %s", macAddress.c_str());
  LOG("==================================");

  checkReset();

  if(prefs.getBool("paired",false)){
    userID=prefs.getString("user_id","");
    deviceID=prefs.getString("device_id","");

    LOG("[INIT] paired user=%s device=%s", userID.c_str(), deviceID.c_str());

    if(connectWiFi() && connectMQTT()) return;
  }

  if(prefs.getString("pair_code","")!=""){
    LOG("[INIT] Found pairing code");
    if(connectWiFi() && confirmPair() && connectMQTT()) return;
  }

  LOG("[INIT] Entering AP pairing mode");
  startCaptive();
}

// ================= LOOP =================
void loop(){
  if(deviceMode==MODE_PAIRING){
    dnsServer.processNextRequest();
    server.handleClient();
    led(millis()%1000<500);
    return;
  }

  if(deviceMode==MODE_MQTT){
    mqtt.loop();
    if(millis()-lastStatePublish>STATE_INTERVAL){
      publishState();
      lastStatePublish=millis();
    }
  }

  delay(10);
}
