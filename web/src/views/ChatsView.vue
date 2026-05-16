<script setup lang="ts">
import { ref } from 'vue'
import ChatSidebar from '@/components/ChatSidebar.vue'
import ChatMessages from '@/components/ChatMessages.vue'
import ChatInput from '@/components/ChatInput.vue'
import { streamMessage, generateTitle } from '@/api/chat'
import type { Message } from '@/types/chat'
import { push } from 'notivue'

const currentChatId = ref<string | null>(null)
const sidebarRef = ref<InstanceType<typeof ChatSidebar> | null>(null)
const messagesRef = ref<InstanceType<typeof ChatMessages> | null>(null)
const inputRef = ref<InstanceType<typeof ChatInput> | null>(null)
const sidebarOpen = ref(false)

// Tracks chat IDs that already have a generated title
const titledChats = new Set<string>()

const handleChatSelected = (chatId: string) => {
  currentChatId.value = chatId
  messagesRef.value?.loadHistory(chatId)
  sidebarOpen.value = false
}

const handleNewChat = (chatId: string) => {
  currentChatId.value = chatId
  messagesRef.value?.loadHistory(chatId)
  sidebarRef.value?.refreshChats()
  sidebarOpen.value = false
}

const handleChatDeleted = (chatId: string) => {
  if (currentChatId.value === chatId) {
    currentChatId.value = null
    messagesRef.value?.clearMessages()
  }
}

const handleMessageSent = (chatId: string, content: string) => {
  const messages = messagesRef.value?.messages
  if (!messages) return

  const isFirstMessage = messages.length === 0

  const userMsg: Message = {
    id: crypto.randomUUID(),
    role: 'user',
    content,
    status: 'completed',
    created_at: new Date().toISOString(),
  }
  messages.push(userMsg)

  const assistantMsg: Message = {
    id: 'streaming',
    role: 'assistant',
    content: '',
    status: 'pending',
    created_at: new Date().toISOString(),
  }
  messages.push(assistantMsg)

  messagesRef.value?.scrollToBottom()

  streamMessage(
    chatId,
    content,
    (token) => {
      const msg = messages.find((m) => m.id === 'streaming')
      if (msg) {
        msg.content += token
        messagesRef.value?.scrollToBottom()
      }
    },
    async () => {
      messagesRef.value?.loadHistory(chatId)

      // Generate title after the first message in a new chat
      if (isFirstMessage && !titledChats.has(chatId)) {
        titledChats.add(chatId)
        try {
          const res = await generateTitle(chatId)
          sidebarRef.value?.updateChatTitle(chatId, res.data.title)
        } catch {
          // Non-critical — title generation failure doesn't break the chat
        }
      }
    },
    (error) => {
      const idx = messages.findIndex((m) => m.id === 'streaming')
      if (idx !== -1) messages.splice(idx, 1)
      push.error(error)
    },
  )
}

const handleError = (error: string) => push.error(error)
</script>

<template>
  <div class="chats-view">
    <!-- Mobile overlay -->
    <transition name="fade">
      <div v-if="sidebarOpen" class="sidebar-overlay" @click="sidebarOpen = false" />
    </transition>

    <ChatSidebar
      ref="sidebarRef"
      :open="sidebarOpen"
      @chat-selected="handleChatSelected"
      @new-chat="handleNewChat"
      @chat-deleted="handleChatDeleted"
      @close="sidebarOpen = false"
    />

    <div class="chat-main">
      <!-- Mobile header bar -->
      <div class="mobile-bar">
        <button class="mobile-menu-btn" @click="sidebarOpen = true">
          <svg width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h16M4 18h16"/>
          </svg>
        </button>
        <span class="mobile-title">{{ currentChatId ? 'Чат' : 'Smart Home' }}</span>
        <div style="width: 36px;" />
      </div>

      <ChatMessages ref="messagesRef" :chat-id="currentChatId" />
      <ChatInput
        ref="inputRef"
        :chat-id="currentChatId"
        @chat-created="handleNewChat"
        @message-sent="handleMessageSent"
        @error="handleError"
      />
    </div>
  </div>
</template>

<style scoped>
.chats-view {
  display: flex;
  height: 100%;
  background: #171717;
  position: relative;
  overflow: hidden;
}

.chat-main {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-width: 0;
  overflow: hidden;
}

.sidebar-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  z-index: 40;
  backdrop-filter: blur(2px);
}

.mobile-bar {
  display: none;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid #262626;
  background: #171717;
}

.mobile-menu-btn {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  border: 1px solid #333;
  background: transparent;
  color: #a3a3a3;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.15s, color 0.15s;
}
.mobile-menu-btn:hover { background: #262626; color: #fff; }

.mobile-title {
  font-size: 15px;
  font-weight: 500;
  color: #e5e5e5;
}

@media (max-width: 768px) {
  .mobile-bar { display: flex; }
}

.fade-enter-active,
.fade-leave-active { transition: opacity 0.2s; }
.fade-enter-from,
.fade-leave-to { opacity: 0; }
</style>
