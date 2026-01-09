<script setup lang="ts">
import {ref, onMounted, onUnmounted} from "vue";
import {useRoute, useRouter} from "vue-router";
import {pairingStatus} from "../api/pairing.ts";

const route = useRoute();
const router = useRouter();

const code = route.query.code as string;
const status = ref("waiting");

let timer:any;

async function check() {
  try {
    const r = await pairingStatus({code});
    if (r.data.status === "done") {
      status.value = "done";
      clearInterval(timer);
    }
  } catch {
    status.value = "error";
    clearInterval(timer);
  }
}

onMounted(() => {
  timer = setInterval(check, 2000);
});

onUnmounted(() => clearInterval(timer));
</script>

<template>
  <div class="page">
    <h2>Подключение устройства</h2>

    <div v-if="status==='waiting'">
      <p>Ожидаем подключения устройства…</p>
    </div>

    <div v-if="status==='done'">
      <p>Устройство успешно добавлено 🎉</p>
      <button @click="router.push('/devices')">Готово</button>
    </div>

    <div v-if="status==='error'">
      <p>Ошибка подключения</p>
      <button @click="router.back()">Назад</button>
    </div>
  </div>
</template>
