<script setup lang="ts">
import { ref, onUnmounted, watch } from 'vue'
import { startPairing, pairingStatus } from '@/api/pairing'
import BaseButton from '@/components/BaseButton.vue'

const props = defineProps<{ modelValue: boolean }>()
const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'success'): void
}>()

const code = ref('')
const status = ref('loading')
let timer: any = null

watch(
  () => props.modelValue,
  (isOpen) => {
    if (isOpen) startPairingProcess()
    else cleanup()
  },
)

async function startPairingProcess() {
  status.value = 'loading'
  try {
    const r = await startPairing()
    code.value = r.data.code
    status.value = 'waiting'
    timer = setInterval(checkStatus, 2000)
  } catch {
    status.value = 'error'
  }
}

async function checkStatus() {
  try {
    const r = await pairingStatus({ code: code.value })
    if (r.data.status === 'done') {
      cleanup()
      status.value = 'done'
      emit('success')
    } else if (r.data.status === 'expired') {
      cleanup()
      status.value = 'expired'
    }
  } catch {
    status.value = 'error'
  }
}

function cleanup() {
  if (timer) clearInterval(timer)
}

function close() {
  cleanup()
  emit('update:modelValue', false)
}

onUnmounted(() => cleanup())
</script>
<template>
  <Teleport to="body">
    <div v-if="modelValue" class="modal-backdrop" @click.self="close">
      <div class="modal">
        <h2>Pair device</h2>

        <div v-if="status === 'loading'">
          <p>Loading...</p>
        </div>

        <div v-if="status === 'waiting'">
          <p>Pass the code the the device:</p>
          <div class="code">{{ code }}</div>
          <div class="actions">
            <BaseButton variant="primary" @click="close">Cancel</BaseButton>
          </div>
        </div>

        <div v-if="status === 'done'">
          <p>Device successfully paired!!</p>
          <div class="actions">
            <BaseButton variant="primary" @click="close">Ok</BaseButton>
          </div>
        </div>

        <div v-if="status === 'expired'">
          <p>Code has expired</p>
          <div class="actions">
            <BaseButton variant="primary" @click="startPairingProcess">Try again</BaseButton>
            <BaseButton variant="secondary" @click="close">Close</BaseButton>
          </div>
        </div>

        <div v-if="status === 'error'">
          <p>Ошибка подключения</p>
          <div class="actions">
            <BaseButton variant="primary" @click="close">Close</BaseButton>
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
}

.modal {
  background: white;
  padding: 24px;
  border-radius: 8px;
  min-width: 300px;
  max-width: 90%;
  text-align: center;
}

.code {
  font-size: 24px;
  font-weight: bold;
  letter-spacing: 4px;
  margin: 16px 0;
  padding: 12px;
  background: #f3f4f6;
  border-radius: 4px;
  text-align: center;
}

.actions {
  display: flex;
  justify-content: center;
  gap: 8px;
  margin-top: 16px;
}
</style>
