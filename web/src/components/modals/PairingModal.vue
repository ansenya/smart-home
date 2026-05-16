<script setup lang="ts">
import { ref, onUnmounted, watch } from 'vue'
import { startPairing, pairingStatus } from '@/api/pairing'

const props = defineProps<{ modelValue: boolean }>()
const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'success'): void
}>()

const code = ref('')
const status = ref('loading')
let timer: any = null

watch(() => props.modelValue, (isOpen) => {
  if (isOpen) startPairingProcess()
  else cleanup()
})

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
    if (r.data.status === 'done') { cleanup(); status.value = 'done'; emit('success') }
    else if (r.data.status === 'expired') { cleanup(); status.value = 'expired' }
  } catch {
    status.value = 'error'
  }
}

function cleanup() { if (timer) clearInterval(timer) }
function close() { cleanup(); emit('update:modelValue', false) }

onUnmounted(cleanup)
</script>

<template>
  <Teleport to="body">
    <transition name="modal-fade">
      <div v-if="modelValue" class="backdrop" @click.self="close">
        <div class="modal">
          <div class="modal-header">
            <h2 class="modal-title">Pair Device</h2>
            <button class="close-btn" @click="close">
              <svg width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>

          <div class="modal-body">
            <div v-if="status === 'loading'" class="state-row">
              <div class="spinner" />
              <span>Generating code...</span>
            </div>

            <div v-else-if="status === 'waiting'">
              <p class="state-label">Enter this code on your device:</p>
              <div class="code-display">{{ code }}</div>
              <p class="state-hint">Waiting for device to connect...</p>
              <div class="actions">
                <button class="btn-secondary" @click="close">Cancel</button>
              </div>
            </div>

            <div v-else-if="status === 'done'" class="state-row state-row--success">
              <div class="state-icon state-icon--success">
                <svg width="20" height="20" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"/>
                </svg>
              </div>
              <div>
                <div class="state-title">Device paired!</div>
                <div class="state-sub">Successfully connected.</div>
              </div>
              <div class="actions">
                <button class="btn-primary" @click="close">Done</button>
              </div>
            </div>

            <div v-else-if="status === 'expired'" class="state-row state-row--warn">
              <div class="state-icon state-icon--warn">
                <svg width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/>
                </svg>
              </div>
              <div>
                <div class="state-title">Code expired</div>
                <div class="state-sub">Generate a new one.</div>
              </div>
              <div class="actions">
                <button class="btn-primary" @click="startPairingProcess">Try again</button>
                <button class="btn-secondary" @click="close">Close</button>
              </div>
            </div>

            <div v-else-if="status === 'error'" class="state-row state-row--error">
              <div class="state-icon state-icon--error">
                <svg width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                </svg>
              </div>
              <div>
                <div class="state-title">Connection error</div>
                <div class="state-sub">Check your network.</div>
              </div>
              <div class="actions">
                <button class="btn-secondary" @click="close">Close</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </transition>
  </Teleport>
</template>

<style scoped>
.backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
  backdrop-filter: blur(4px);
}

.modal {
  background: #1a1a1a;
  border: 1px solid #2a2a2a;
  border-radius: 16px;
  width: 100%;
  max-width: 360px;
  box-shadow: 0 24px 64px rgba(0,0,0,0.6);
  overflow: hidden;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 18px 20px 14px;
  border-bottom: 1px solid #222;
}

.modal-title {
  font-size: 15px;
  font-weight: 600;
  color: #e5e5e5;
}

.close-btn {
  width: 28px;
  height: 28px;
  border-radius: 7px;
  border: none;
  background: transparent;
  color: #555;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.15s, color 0.15s;
}
.close-btn:hover { background: #222; color: #a3a3a3; }

.modal-body {
  padding: 24px 20px 20px;
}

.state-row {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  text-align: center;
  color: #737373;
  font-size: 14px;
}

.state-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.state-icon--success { background: #052e16; color: #4ade80; }
.state-icon--warn    { background: #2d1b00; color: #fbbf24; }
.state-icon--error   { background: #1f0f0f; color: #f87171; }

.state-title { font-size: 14px; font-weight: 600; color: #e5e5e5; margin-bottom: 2px; }
.state-sub   { font-size: 13px; color: #555; }
.state-label { font-size: 13px; color: #737373; margin-bottom: 14px; text-align: center; }

.spinner {
  width: 28px;
  height: 28px;
  border: 2px solid #2a2a2a;
  border-top-color: #6366f1;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }

.code-display {
  font-size: 28px;
  font-weight: 700;
  letter-spacing: 6px;
  color: #e5e5e5;
  text-align: center;
  padding: 14px;
  background: #111;
  border: 1px solid #2a2a2a;
  border-radius: 10px;
  font-family: ui-monospace, monospace;
  margin-bottom: 10px;
}

.state-hint {
  font-size: 12px;
  color: #404040;
  text-align: center;
}

.actions {
  display: flex;
  gap: 8px;
  justify-content: center;
  margin-top: 16px;
  width: 100%;
}

.btn-primary, .btn-secondary {
  padding: 8px 18px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  font-family: inherit;
  transition: background 0.15s, color 0.15s;
}

.btn-primary {
  background: #6366f1;
  color: #fff;
  border: 1px solid #6366f1;
}
.btn-primary:hover { background: #4f46e5; }

.btn-secondary {
  background: transparent;
  color: #737373;
  border: 1px solid #2a2a2a;
}
.btn-secondary:hover { background: #222; color: #a3a3a3; }

.modal-fade-enter-active, .modal-fade-leave-active {
  transition: opacity 0.2s;
}
.modal-fade-enter-active .modal, .modal-fade-leave-active .modal {
  transition: transform 0.2s;
}
.modal-fade-enter-from, .modal-fade-leave-to { opacity: 0; }
.modal-fade-enter-from .modal, .modal-fade-leave-to .modal {
  transform: scale(0.96) translateY(8px);
}
</style>
