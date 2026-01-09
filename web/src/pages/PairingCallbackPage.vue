<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { pairingStatus } from "../api/pairing";

const route = useRoute();
const router = useRouter();

const status = ref<"waiting" | "done" | "error">("waiting");
const code = ref<string>((route.query.code as string) || "");
let timer: any = null;

function startPolling() {
  timer = setInterval(checkStatus, 2000);
}

async function checkStatus() {
  try {
    const r = await pairingStatus({ code: code.value });
    if (r.data.status === "done") {
      clearInterval(timer);
      status.value = "done";
      setTimeout(() => router.push("/devices"), 800);
    } else {
      status.value = "waiting";
    }
  } catch (e) {
    clearInterval(timer);
    status.value = "error";
  }
}

onMounted(() => {
  if (!code.value) {
    status.value = "error";
    return;
  }
  startPolling();
});

onUnmounted(() => clearInterval(timer));
</script>

<template>
  <div class="page">
    <h2>Завершение привязки</h2>

    <div v-if="status === 'waiting'">
      <p>Идёт привязка устройства…</p>
      <div class="code">{{ code }}</div>
      <p class="muted">После завершения вы будете перенаправлены в список устройств.</p>
    </div>

    <div v-if="status === 'done'">
      <p>Устройство успешно добавлено. Сейчас перенаправим вас.</p>
    </div>

    <div v-if="status === 'error'">
      <p>Не удалось завершить привязку. Попробуйте снова.</p>
      <button @click="router.push('/devices/add')">Попробовать снова</button>
    </div>
  </div>
</template>

<style scoped>
.page { max-width: 420px; margin: auto; padding: 24px; text-align:center; }
.code { font-size: 28px; letter-spacing: 4px; margin-top: 16px; }
.muted { opacity:.7; font-size:13px; }
button { margin-top: 16px; padding: 10px 16px; background:#2563eb; color:white; border:none; border-radius:6px; }
</style>
