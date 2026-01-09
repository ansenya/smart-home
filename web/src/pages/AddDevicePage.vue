<script setup lang="ts">
import {ref} from "vue";
import {useRouter} from "vue-router";
import {startPairing} from "../api/pairing";

const router = useRouter();

const step = ref<
    "intro" |
    "requesting" |
    "connecting" |
    "redirected" |
    "waiting" |
    "done" |
    "error"
>("intro");

const code = ref("");
const ssid = ref("");
const password = ref("");

async function begin() {
  step.value = "requesting";

  try {
    const r = await startPairing();
    code.value = r.data.code;

    step.value = "connecting";
    // Инструкция пользователю: подключиться к AP устройства вручную
    // (можно добавить UI с подсказками)
    // После подключения браузер будет редиректнут на ESP (см. redirectToDevice)
    redirectToDeviceAP();
  } catch (e) {
    console.error(e);
    step.value = "error";
  }
}

function redirectToDeviceAP() {
  // redirect_url — куда ESP вернёт браузер после получения параметров
  const callbackUrl = `${window.location.origin}/pairing/callback?code=${encodeURIComponent(code.value)}`;

  // формируем URL устройства; IP стандартный для ESP softAP
  const deviceUrl = `http://192.168.4.1/` +
      `?wifi_ssid=${encodeURIComponent(ssid.value)}` +
      `&wifi_password=${encodeURIComponent(password.value)}` +
      `&pairing_code=${encodeURIComponent(code.value)}` +
      `&redirect_url=${encodeURIComponent(callbackUrl)}`;

  // делаем переход — браузер переходит в HTTP на устройство (allowed)
  window.location.href = deviceUrl;

  step.value = "redirected";
}
</script>

<template>
  <div class="page">
    <h2>Добавление устройства</h2>

    <div v-if="step === 'intro'">
      <p>
        Мы автоматически подключим устройство к вашей Wi-Fi сети.
        Устройство создаст временную сеть — подключитесь к ней (или используйте тот же телефон).
      </p>
      <p>Введите сеть и пароль, затем нажмите «Начать подключение» — браузер откроет страницу устройства.</p>

      <input v-model="ssid" placeholder="Wi-Fi сеть"/>
      <input v-model="password" placeholder="Пароль Wi-Fi" type="password"/>

      <button @click="begin">Начать подключение</button>
    </div>

    <div v-if="step === 'requesting'">
      <p>Получаем код подключения…</p>
    </div>

    <div v-if="step === 'connecting'">
      <p>Через секунду вы будете перенаправлены на страницу устройства. Подключитесь к Wi-Fi сети устройства
        (SmartHome-XXXX) и дождитесь автоматического возврата.</p>
      <p class="muted">Если браузер не перешёл автоматически — подключитесь к сети вручную и откройте <code>http://192.168.4.1</code>.
      </p>
    </div>

    <div v-if="step === 'redirected'">
      <p>Вы были перенаправлены на устройство. После отправки данных устройство вернёт вас обратно на страницу
        ожидания.</p>
      <p class="muted">Если ничего не произошло — проверьте подключение к сети устройства.</p>
    </div>

    <div v-if="step === 'waiting'">
      <p>Ожидаем подтверждения от backend…</p>
      <div class="code">{{ code }}</div>
    </div>

    <div v-if="step === 'done'">
      <p>Устройство успешно добавлено 🎉</p>
      <button @click="router.push('/devices')">Перейти к устройствам</button>
    </div>

    <div v-if="step === 'error'">
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
  font-size: 28px;
  letter-spacing: 4px;
  margin-top: 16px;
  text-align: center;
}

.muted {
  opacity: .7;
  font-size: 13px;
}
</style>
