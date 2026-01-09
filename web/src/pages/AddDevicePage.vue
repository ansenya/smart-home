<script setup lang="ts">
import {ref} from "vue";
import {useRouter} from "vue-router";
import {startPairing} from "../api/pairing.ts";

const router = useRouter();

const ssid = ref("");
const password = ref("");

async function begin() {
  const r = await startPairing();

  const code = r.data.code;

  const redirect = encodeURIComponent(
      window.location.origin + "/devices/add/wait?code=" + code
  );

  const url =
      `http://192.168.4.1/setup` +
      `?wifi_ssid=${encodeURIComponent(ssid.value)}` +
      `&wifi_password=${encodeURIComponent(password.value)}` +
      `&pairing_code=${encodeURIComponent(code)}` +
      `&redirect_url=${redirect}`;

  window.location.href = url;
}
</script>

<template>
  <div class="page">
    <h2>Добавление устройства</h2>

    <p>
      Подключитесь к Wi-Fi сети устройства, затем введите данные вашей сети.
    </p>

    <input v-model="ssid" placeholder="Wi-Fi сеть"/>
    <input v-model="password" placeholder="Пароль Wi-Fi" type="password"/>

    <button @click="begin">Продолжить</button>
  </div>
</template>
