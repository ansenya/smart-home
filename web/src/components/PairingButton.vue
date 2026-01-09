<script setup lang="ts">
import {ref, onUnmounted} from "vue";
import {pairingStatus, startPairing} from "../api/pairing.ts";

const open = ref(false);
const code = ref("");
const status = ref<"waiting" | "done" | "error">("waiting");

let timer: any = null;

async function start() {
  const res = await startPairing();
  code.value = res.data.code;
  status.value = "waiting";
  open.value = true;

  timer = setInterval(checkStatus, 2000);
}

async function checkStatus() {
  try {
    const r = await pairingStatus({code: code.value});

    if (r.data.status === "done") {
      status.value = "done";
      clearInterval(timer);
    }
  } catch {
    status.value = "error";
    clearInterval(timer);
  }
}

function close() {
  open.value = false;
  clearInterval(timer);
}

onUnmounted(() => clearInterval(timer));
</script>

<template>
  <button class="pair-btn" @click="start">+</button>

  <div v-if="open" class="modal">
    <div class="modal-box">
      <h3>Добавление устройства</h3>

      <div v-if="status==='waiting'">
        <p>Введите этот код на устройстве:</p>
        <div class="code">{{ code }}</div>
        <p class="wait">Ожидание подключения…</p>
      </div>

      <div v-if="status==='done'">
        <p>Устройство успешно добавлено</p>
        <button @click="close">Готово</button>
      </div>

      <div v-if="status==='error'">
        <p>Ошибка подключения</p>
        <button @click="close">Закрыть</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.pair-btn {
  position: fixed;
  right: 24px;
  bottom: 24px;
  width: 56px;
  height: 56px;
  border-radius: 50%;
  font-size: 28px;
  background: #2563eb;
  color: white;
  border: none;
}

.modal {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, .4);
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal-box {
  background: white;
  padding: 24px;
  border-radius: 12px;
  width: 300px;
  text-align: center;
}

.code {
  font-size: 32px;
  letter-spacing: 4px;
  margin: 16px 0;
  font-weight: bold;
}

.wait {
  opacity: .7;
}
</style>
