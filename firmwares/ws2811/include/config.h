#ifndef CONFIG_H
#define CONFIG_H

#include <Arduino.h>

#define AP_IP_ADDR IPAddress(192,168,4,1)
#define DNS_PORT 53
#define PAIRING_URL "https://api.smarthome.hipahopa.ru/devices/pairing/confirm"

#define MQTT_HOST "hipahopa.ru"
#define MQTT_PORT 1883
#define MQTT_USER "user"
#define MQTT_PASSWORD "superdifficultpassword"

#define RESET_PIN 0
#define LED_PIN 2

#ifndef STRIP_PIN
#define STRIP_PIN 5
#endif
#ifndef STRIP_COUNT
#define STRIP_COUNT 300
#endif

#define WIFI_TIMEOUT 10000
#define STATE_INTERVAL 5000

// Надёжная лог функция: печатает формат и перевод строки
#define LOG(fmt, ...) do { Serial.printf((fmt), ##__VA_ARGS__); Serial.println(); } while(0)

#endif // CONFIG_H
