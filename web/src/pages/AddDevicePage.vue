<script setup lang="ts">
import {ref, onUnmounted} from "vue";
import {useRouter} from "vue-router";
import {pairingStatus, startPairing} from "../api/pairing.ts";

const router = useRouter();

const step = ref<
    "intro" |
    "requesting" |
    "connecting" |
    "sending" |
    "waiting" |
    "done" |
    "error"
>("intro");

const code = ref("");
const ssid = ref("");
const password = ref("");

let timer: any = null;

async function begin() {
  step.value = "requesting";

  try {
    const r = await startPairing();
    code.value = r.data.code;

    step.value = "connecting";
    await connectToDeviceAP();

    step.value = "sending";
    await sendConfigToDevice();

    step.value = "waiting";
    startPolling();
  } catch (e) {
    step.value = "error";
  }
}

async function connectToDeviceAP() {
  // Пользователь подключается вручную к WiFi AP устройства
  alert("Подключитесь к Wi-Fi сети устройства (SmartHome-XXXX), затем нажмите OK");
}

async function sendConfigToDevice() {
  await fetch("http://192.168.4.1/config", {
    method: "POST",
    headers: {"Content-Type": "application/json"},
    body: JSON.stringify({
      wifi_ssid: ssid.value,
      wifi_password: password.value,
      pairing_code: code.value
    })
  });
}

function startPolling() {
  timer = setInterval(checkStatus, 2000);
}

async function checkStatus() {
  try {
    const r = await pairingStatus({code: code.value});
    if (r.data.status === "done") {
      clearInterval(timer);
      step.value = "done";
    }
  } catch {
    clearInterval(timer);
    step.value = "error";
  }
}

function finish() {
  router.push("/devices");
}

onUnmounted(() => clearInterval(timer));
</script>

<template>
  <div class="page">

    <h2>Добавление устройства</h2>

    <div v-if="step==='intro'">
      <p>
        Сейчас мы автоматически подключим устройство к вашей Wi-Fi сети.
        Для этого устройство временно создаст собственную сеть.
      </p>
      <p>
        Выберите вашу Wi-Fi сеть и введите пароль — дальше всё произойдёт автоматически.
      </p>

      <input v-model="ssid" placeholder="Wi-Fi сеть"/>
      <input v-model="password" placeholder="Пароль Wi-Fi" type="password"/>

      <button @click="begin">Начать подключение</button>
    </div>

    <div v-if="step==='requesting'">
      <p>Получаем код подключения…</p>
    </div>

    <div v-if="step==='connecting'">
      <p>Подключитесь к сети устройства…</p>
    </div>

    <div v-if="step==='sending'">
      <p>Передаём данные устройству…</p>
    </div>

    <div v-if="step==='waiting'">
      <p>Ожидание подтверждения…</p>
      <div class="code">{{ code }}</div>
    </div>

    <div v-if="step==='done'">
      <p>Устройство успешно добавлено 🎉</p>
      <button @click="finish">Готово</button>
    </div>

    <div v-if="step==='error'">
      <p>Ошибка подключения</p>
      <button @click="router.back()">Назад</button>
    </div>

  </div>
</template>

<style scoped>
.page {
  max-width: 420px;
  margin: auto;
  padding: 24px;
}

input {
  width: 100%;
  margin: 8px 0;
  padding: 8px;
}

button {
  margin-top: 16px;
  width: 100%;
  padding: 10px;
  background: #2563eb;
  color: white;
  border: none;
  border-radius: 6px;
}

.code {
  font-size: 32px;
  letter-spacing: 4px;
  margin-top: 16px;
}
</style>
