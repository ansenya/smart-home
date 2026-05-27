<script setup lang="ts">
import { ref, nextTick, computed } from 'vue'
import { createChat } from '@/api/chat'

const AVAILABLE_MODELS = [
  { value: 'gpt-4o', label: 'GPT-4o' },
  { value: 'gpt-4o-mini', label: 'GPT-4o mini' },
  { value: 'claude-sonnet-4-5', label: 'Claude Sonnet 4.5' },
  { value: 'claude-opus-4-5', label: 'Claude Opus 4.5' },
  { value: 'claude-haiku-4-5', label: 'Claude Haiku 4.5' },
]

const props = defineProps<{ chatId: string | null }>()

const emit = defineEmits<{
  messageSent: [chatId: string, content: string]
  chatCreated: [chatId: string]
  streaming: [isStreaming: boolean]
  error: [error: string]
}>()

const textarea = ref<HTMLTextAreaElement | null>(null)
const content = ref('')
const isLoading = ref(false)
const selectedModel = ref('gpt-4o')
const showModelPicker = ref(false)

const canSend = computed(() => content.value.trim().length > 0 && !isLoading.value)

const currentModel = computed(() => AVAILABLE_MODELS.find(m => m.value === selectedModel.value) ?? AVAILABLE_MODELS[0]!)

const autoResize = () => {
  if (textarea.value) {
    textarea.value.style.height = 'auto'
    textarea.value.style.height = Math.min(textarea.value.scrollHeight, 200) + 'px'
  }
}

const focusTextarea = async () => {
  await nextTick()
  textarea.value?.focus()
}

const handleKeydown = (event: KeyboardEvent) => {
  if (event.key === 'Enter' && !event.shiftKey && !event.ctrlKey && !event.metaKey) {
    event.preventDefault()
    if (canSend.value) send()
  }
  if ((event.ctrlKey || event.metaKey) && event.key === 'Enter') {
    event.preventDefault()
    if (canSend.value) send()
  }
}

const send = async () => {
  const message = content.value.trim()
  if (!message || isLoading.value) return

  isLoading.value = true
  emit('streaming', true)

  try {
    let chatId = props.chatId
    if (!chatId) {
      const chatRes = await createChat({ model: selectedModel.value, title: 'New Chat' })
      chatId = chatRes.data.id
      emit('chatCreated', chatId)
    }

    content.value = ''
    if (textarea.value) textarea.value.style.height = 'auto'
    await focusTextarea()

    emit('messageSent', chatId, message)
  } catch (error: any) {
    emit('error', error.response?.data?.error || 'Failed to send message')
  } finally {
    isLoading.value = false
    emit('streaming', false)
  }
}

const sendNow = (text: string) => {
  content.value = text
  send()
}

defineExpose({ clear: () => { content.value = ''; if (textarea.value) textarea.value.style.height = 'auto' }, focusTextarea, selectedModel, sendNow })

focusTextarea()
</script>

<template>
  <div class="input-area">
    <div class="input-box" :class="{ 'input-box--active': content.length > 0 }">
      <textarea
        ref="textarea"
        v-model="content"
        placeholder="Message Smart Home..."
        rows="1"
        :disabled="isLoading"
        @input="autoResize"
        @keydown="handleKeydown"
      />

      <div class="input-footer">
        <!-- Model picker -->
        <div class="model-picker-wrap">
          <button class="model-btn" :disabled="isLoading" @click="showModelPicker = !showModelPicker">
            <svg width="12" height="12" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17H3a2 2 0 01-2-2V5a2 2 0 012-2h14a2 2 0 012 2v10a2 2 0 01-2 2h-2"/>
            </svg>
            {{ currentModel.label }}
            <svg width="10" height="10" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7"/>
            </svg>
          </button>

          <div v-if="showModelPicker" class="model-dropdown">
            <button
              v-for="m in AVAILABLE_MODELS"
              :key="m.value"
              class="model-option"
              :class="{ 'model-option--active': m.value === selectedModel }"
              @click="selectedModel = m.value; showModelPicker = false"
            >
              <span class="model-option-check">
                <svg v-if="m.value === selectedModel" width="12" height="12" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"/>
                </svg>
              </span>
              {{ m.label }}
            </button>
          </div>
        </div>

        <!-- Send button -->
        <button class="send-btn" :class="{ 'send-btn--active': canSend }" :disabled="!canSend" @click="send">
          <svg v-if="!isLoading" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M5 12h14M12 5l7 7-7 7"/>
          </svg>
          <div v-else class="loader" />
        </button>
      </div>
    </div>

    <p class="hint">Enter to send · Shift+Enter for new line</p>
  </div>
</template>

<style scoped>
.input-area {
  padding: 12px 16px 16px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.input-box {
  width: 100%;
  max-width: 720px;
  background: #1e1e1e;
  border: 1px solid #2a2a2a;
  border-radius: 14px;
  transition: border-color 0.15s;
  overflow: hidden;
}

.input-box:focus-within {
  border-color: #3a3a3a;
}

.input-box--active {
  border-color: #333;
}

textarea {
  width: 100%;
  background: transparent;
  border: none;
  color: #e5e5e5;
  padding: 14px 16px 10px;
  font-size: 14px;
  line-height: 1.6;
  resize: none;
  max-height: 200px;
  font-family: inherit;
  overflow-y: auto;
}

textarea::placeholder {
  color: #404040;
}

textarea:focus {
  outline: none;
}

textarea:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

textarea::-webkit-scrollbar {
  width: 4px;
}
textarea::-webkit-scrollbar-track { background: transparent; }
textarea::-webkit-scrollbar-thumb { background: #333; border-radius: 2px; }

.input-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 10px 10px;
}

/* Model picker */
.model-picker-wrap {
  position: relative;
}

.model-btn {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 5px 10px;
  border-radius: 8px;
  border: 1px solid #2a2a2a;
  background: transparent;
  color: #555;
  font-size: 12px;
  cursor: pointer;
  transition: background 0.15s, color 0.15s;
  font-family: inherit;
  white-space: nowrap;
}

.model-btn:hover:not(:disabled) {
  background: #252525;
  color: #888;
}

.model-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.model-dropdown {
  position: absolute;
  bottom: calc(100% + 6px);
  left: 0;
  background: #1a1a1a;
  border: 1px solid #2a2a2a;
  border-radius: 10px;
  padding: 4px;
  min-width: 180px;
  z-index: 10;
  box-shadow: 0 8px 32px rgba(0,0,0,0.6);
}

.model-option {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  padding: 8px 10px;
  border-radius: 7px;
  border: none;
  background: transparent;
  color: #a3a3a3;
  font-size: 13px;
  cursor: pointer;
  text-align: left;
  transition: background 0.12s, color 0.12s;
  font-family: inherit;
}

.model-option:hover {
  background: #222;
  color: #e5e5e5;
}

.model-option--active {
  color: #e5e5e5;
}

.model-option-check {
  width: 16px;
  height: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #6366f1;
  flex-shrink: 0;
}

/* Send button */
.send-btn {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  border: none;
  background: #2a2a2a;
  color: #555;
  cursor: not-allowed;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.15s, color 0.15s, transform 0.1s;
  flex-shrink: 0;
}

.send-btn--active {
  background: #e5e5e5;
  color: #000;
  cursor: pointer;
}

.send-btn--active:hover {
  background: #fff;
  transform: scale(1.05);
}

.send-btn--active:active {
  transform: scale(0.95);
}

.loader {
  width: 14px;
  height: 14px;
  border: 2px solid #555;
  border-top-color: transparent;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.hint {
  font-size: 11px;
  color: #333;
}

@media (max-width: 768px) {
  .hint { display: none; }
  .input-area { padding: 10px 12px 14px; }
}
</style>
