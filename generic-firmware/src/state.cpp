#include "state.h"
#include <esp_system.h>

DeviceMode deviceMode = MODE_PAIRING;
Preferences prefs;
WebServer server(80);
DNSServer dnsServer;
WiFiClient net;
PubSubClient mqtt(net);

String deviceUID, macAddress, userID, deviceID;
bool powerState = false;
unsigned long lastStatePublish = 0;

String getDeviceUID() {
  uint64_t id = ESP.getEfuseMac();
  char buf[32];
  sprintf(buf, "esp32-%04X%08X", (uint16_t)(id>>32), (uint32_t)id);
  return String(buf);
}

String getMac() {
  uint8_t mac[6];
  esp_read_mac(mac, ESP_MAC_WIFI_STA);
  char buf[32];
  sprintf(buf, "%02X:%02X:%02X:%02X:%02X:%02X", mac[0],mac[1],mac[2],mac[3],mac[4],mac[5]);
  return String(buf);
}

String topicBase() {
  return userID + "/" + deviceID;
}

void led(bool s){ digitalWrite(LED_PIN, s?HIGH:LOW); }
