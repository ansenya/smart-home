<script setup lang="ts">
import { nextTick, ref } from 'vue'
import ChatMessage from '@/components/ChatMessage.vue'
import { getHistory } from '@/api/chat'
import type { Message } from '@/types/chat'

const props = defineProps<{ chatId: string | null }>()

const messages = ref<Message[]>([])
const messagesContainer = ref<HTMLElement | null>(null)
const isLoading = ref(false)

const scrollToBottom = async () => {
  await nextTick()
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

const loadHistory = async (chatId: string) => {
  isLoading.value = true
  try {
    const res = await getHistory(chatId)
    messages.value = res.data.messages
    await scrollToBottom()
  } catch (error) {
    console.error('Failed to load history:', error)
  } finally {
    isLoading.value = false
  }
}

const clearMessages = () => {
  messages.value = []
}

defineExpose({ loadHistory, clearMessages, messages, scrollToBottom })
</script>

<template>
  <div class="messages-wrap">
    <!-- Empty state -->
    <div v-if="!chatId && messages.length === 0" class="empty-state">
      <div class="empty-icon">
        <svg width="28" height="28" fill="currentColor" viewBox="0 0 20 20">
          <path d="M10.707 2.293a1 1 0 00-1.414 0l-7 7a1 1 0 001.414 1.414L4 10.414V17a1 1 0 001 1h2a1 1 0 001-1v-2a1 1 0 011-1h2a1 1 0 011 1v2a1 1 0 001 1h2a1 1 0 001-1v-6.586l.293.293a1 1 0 001.414-1.414l-7-7z"/>
        </svg>
      </div>
      <h2 class="empty-title">Smart Home Assistant</h2>
      <p class="empty-sub">Ask me about your devices, automations, or anything about your home.</p>
      <div class="suggestions">
        <div class="suggestion">What devices are online?</div>
        <div class="suggestion">Turn off all lights</div>
        <div class="suggestion">Set bedroom to 22°C</div>
      </div>
    </div>

    <!-- Loading skeleton -->
    <div v-else-if="isLoading" class="messages-container">
      <div v-for="i in 4" :key="i" class="skeleton-msg" :class="i % 2 === 0 ? 'skeleton-right' : 'skeleton-left'">
        <div class="skeleton-bubble" :style="{ width: (50 + (i * 17) % 30) + '%' }" />
      </div>
    </div>

    <!-- Messages -->
    <div v-else ref="messagesContainer" class="messages-container">
      <ChatMessage v-for="message in messages" :key="message.id" :message="message" />
    </div>
  </div>
</template>

<style scoped>
.messages-wrap {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 24px 16px;
  display: flex;
  flex-direction: column;
  gap: 2px;
  scroll-behavior: smooth;
}

.messages-container::-webkit-scrollbar {
  width: 4px;
}
.messages-container::-webkit-scrollbar-track {
  background: transparent;
}
.messages-container::-webkit-scrollbar-thumb {
  background: #2a2a2a;
  border-radius: 2px;
}

/* Empty state */
.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 24px;
  text-align: center;
  gap: 12px;
}

.empty-icon {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  background: #1e1e1e;
  border: 1px solid #2a2a2a;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #6366f1;
  margin-bottom: 4px;
}

.empty-title {
  font-size: 20px;
  font-weight: 600;
  color: #e5e5e5;
  letter-spacing: -0.02em;
}

.empty-sub {
  font-size: 14px;
  color: #555;
  max-width: 340px;
  line-height: 1.6;
}

.suggestions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  justify-content: center;
  margin-top: 8px;
}

.suggestion {
  padding: 8px 14px;
  border-radius: 20px;
  border: 1px solid #262626;
  background: #1a1a1a;
  color: #737373;
  font-size: 13px;
  cursor: pointer;
  transition: background 0.15s, color 0.15s, border-color 0.15s;
}

.suggestion:hover {
  background: #212121;
  color: #a3a3a3;
  border-color: #333;
}

/* Skeletons */
.skeleton-msg {
  display: flex;
  margin-bottom: 16px;
}

.skeleton-right {
  justify-content: flex-end;
}

.skeleton-left {
  justify-content: flex-start;
}

.skeleton-bubble {
  height: 40px;
  border-radius: 12px;
  background: linear-gradient(90deg, #1e1e1e 25%, #252525 50%, #1e1e1e 75%);
  background-size: 200% 100%;
  animation: shimmer 1.4s infinite;
}

@keyframes shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}
</style>
