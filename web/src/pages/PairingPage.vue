<script setup lang="ts">
import {ref, onUnmounted} from "vue"
import {useRouter} from "vue-router"
import {startPairing, pairingStatus} from "../api/pairing"

const router = useRouter()

const code = ref("")
const step = ref<"loading" | "waiting" | "done" | "expired" | "error">("loading")
const connectionError = ref(false)

let timer: any = null

async function init() {
  try {
    const r = await startPairing()
    code.value = r.data.code
    step.value = "waiting"
    startPolling()
  } catch {
    step.value = "error"
  }
}

function startPolling() {
  timer = setInterval(checkStatus, 2000)
}

async function checkStatus() {
  try {
    const r = await pairingStatus({code: code.value})

    connectionError.value = false

    if (r.data.status === "done") {
      clearInterval(timer)
      step.value = "done"
      return
    }

    if (r.data.status === "expired") {
      clearInterval(timer)
      step.value = "expired"
      return
    }

    step.value = "waiting"
  } catch {
    connectionError.value = true
  }
}

function finish() {
  router.push("/devices")
}

onUnmounted(() => clearInterval(timer))

init()
</script>

<template>
  <div class="page">
    <h2>Добавление устройства</h2>

    <div v-if="step === 'loading'">
      Получаем код…
    </div>

    <div v-if="step === 'waiting'">
      <p>Включите устройство и подключитесь к его Wi-Fi сети.</p>
      <p>В открывшейся странице введите этот код:</p>

      <div class="code">{{ code }}</div>

      <p class="hint">
        После ввода устройство автоматически появится в системе.
      </p>

      <p v-if="connectionError" class="error-message">
        Не удаётся связаться с сервером. Проверьте подключение.
      </p>
    </div>

    <div v-if="step === 'done'">
      <p>Устройство успешно добавлено</p>
      <button @click="finish">Готово</button>
    </div>

    <div v-if="step === 'expired'">
      <p>Код истёк</p>
      <button @click="router.go(0)">Получить новый</button>
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
  text-align: center;
}

.code {
  font-size: 36px;
  letter-spacing: 6px;
  font-weight: bold;
  margin: 16px 0;
}

.hint {
  opacity: 0.7;
  font-size: 14px;
}

.error-message {
  color: #dc2626;
  font-size: 14px;
  margin-top: 12px;
}

button {
  margin-top: 16px;
  padding: 10px;
  width: 100%;
  background: #2563eb;
  color: white;
  border: none;
  border-radius: 6px;
}
</style>