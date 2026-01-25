# Hipahome

Hipahome — smart home middleware для пользовательских устройств. Моя курсовая работа.

Цель проекта — дать **plug-and-play интеграцию любых девайсов** (ESP32, ПК, кастомные утилиты) без YAML, скриптов и ручных конфигов, при этом сохранить прозрачность архитектуры и совместимость с существующими экосистемами (Yandex Smart Home, Home Assistant).

Система построена как набор микросервисов + MQTT + единая прошивка устройств.

---

## Stack

Backend:

* Go (microservices)
* PostgreSQL
* Redis
* MQTT (Mosquitto)

Frontend:

* Vue 3 + Vite

Devices:

* ESP32 (PlatformIO)
* generic firmware
* ws2811 firmware

Infra:

* Docker / docker-compose
* nginx reverse proxy

---

## Repo structure

```
auth-service   — OAuth, JWT, sessions
panel-api      — control panel API
devices-api    — pairing + device lifecycle
adapter        — bridges (Yandex / HA) + MQTT sync
web            — main UI
auth-web       — auth UI
firmwares      — ESP32 firmware
mosquitto      — broker config
nginx          — reverse proxy
db/init        — SQL bootstrap
```

---

## First setup (once)

### Generate RSA keys

Auth-service подписывает access/refresh токены RSA ключами.

```bash
mkdir -p keys && \
cd keys && \
openssl genrsa -out access_private.pem 4096 && \
openssl rsa -in access_private.pem -pubout -out access_public.pem && \
openssl genrsa -out refresh_private.pem 4096 && \
openssl rsa -in refresh_private.pem -pubout -out refresh_public.pem
```

Ожидаемая структура:

```
keys/
  access_private.pem
  access_public.pem
  refresh_private.pem
  refresh_public.pem
```

---

## Run

```bash
./pull-and-restart.sh
```

Скрипт пересобирает и перезапускает весь стек.

---

## Notes

* все устройства общаются только через MQTT
* pairing через one-time code
* hardware-bound identity
* нет обязательных конфигов/скриптов
* HA и Yandex работают через adapter, не напрямую

---
