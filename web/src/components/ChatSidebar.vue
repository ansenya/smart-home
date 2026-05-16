<script setup lang="ts">
import { onMounted, ref } from 'vue'
import type { Chat } from '@/types/chat'
import { createChat, getChats, deleteChat, updateChat } from '@/api/chat'
import { push } from 'notivue'
import ChatSidebarItem from '@/components/ChatSidebarItem.vue'

defineProps<{ open: boolean }>()

const emit = defineEmits<{
  chatSelected: [chatId: string]
  newChat: [chatId: string]
  chatDeleted: [chatId: string]
  close: []
}>()

const chats = ref<Chat[]>([])
const isLoading = ref(true)
const currentChatID = ref<string | null>(null)

const loadChats = async () => {
  isLoading.value = true
  try {
    const res = await getChats(50)
    chats.value = res.data.chats
  } catch (error: any) {
    push.error({ title: 'Failed to load chats', message: error.message })
  } finally {
    isLoading.value = false
  }
}

const handleNewChat = async () => {
  try {
    const res = await createChat({ title: 'Новый чат' })
    const newChat = res.data
    chats.value.unshift(newChat)
    currentChatID.value = newChat.id
    emit('newChat', newChat.id)
    emit('chatSelected', newChat.id)
  } catch (error: any) {
    push.error({ title: 'Failed to create chat', message: error.message })
  }
}

const selectChat = (chatId: string) => {
  currentChatID.value = chatId
  emit('chatSelected', chatId)
}

const handleDelete = async (chatId: string) => {
  try {
    await deleteChat(chatId)
    chats.value = chats.value.filter((c) => c.id !== chatId)
    if (currentChatID.value === chatId) {
      currentChatID.value = null
    }
    emit('chatDeleted', chatId)
  } catch (error: any) {
    push.error({ title: 'Failed to delete chat', message: error.message })
  }
}

const handleRename = async (chatId: string, title: string) => {
  try {
    const res = await updateChat(chatId, { title })
    const idx = chats.value.findIndex((c) => c.id === chatId)
    if (idx !== -1) chats.value[idx] = res.data
  } catch (error: any) {
    push.error({ title: 'Failed to rename chat', message: error.message })
  }
}

const updateChatTitle = (chatId: string, title: string) => {
  const idx = chats.value.findIndex((c) => c.id === chatId)
  if (idx !== -1) chats.value[idx] = { ...chats.value[idx]!, title }
}

onMounted(loadChats)

defineExpose({ refreshChats: loadChats, updateChatTitle })
</script>

<template>
  <aside class="sidebar" :class="{ 'sidebar--open': open }">
    <div class="sidebar-header">
      <div class="sidebar-logo">
        <div class="logo-icon">
          <svg width="14" height="14" fill="currentColor" viewBox="0 0 20 20">
            <path d="M10.707 2.293a1 1 0 00-1.414 0l-7 7a1 1 0 001.414 1.414L4 10.414V17a1 1 0 001 1h2a1 1 0 001-1v-2a1 1 0 011-1h2a1 1 0 011 1v2a1 1 0 001 1h2a1 1 0 001-1v-6.586l.293.293a1 1 0 001.414-1.414l-7-7z"/>
          </svg>
        </div>
        <span class="logo-text">Smart Home</span>
      </div>
      <button class="close-btn" @click="emit('close')">
        <svg width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
        </svg>
      </button>
    </div>

    <button class="new-chat-btn" :disabled="isLoading" @click="handleNewChat">
      <svg width="15" height="15" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
      </svg>
      Новый чат
    </button>

    <div class="chat-list">
      <template v-if="isLoading">
        <div v-for="i in 5" :key="i" class="skeleton" />
      </template>
      <template v-else-if="chats.length === 0">
        <div class="empty-list">
          <svg width="32" height="32" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24" style="color: #404040;">
            <path stroke-linecap="round" stroke-linejoin="round" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/>
          </svg>
          <span>Нет чатов</span>
        </div>
      </template>
      <template v-else>
        <ChatSidebarItem
          v-for="chat in chats"
          :key="chat.id"
          :chat="chat"
          :active="chat.id === currentChatID"
          @click="selectChat(chat.id)"
          @delete="handleDelete(chat.id)"
          @rename="(title) => handleRename(chat.id, title)"
        />
      </template>
    </div>
  </aside>
</template>

<style scoped>
.sidebar {
  width: 260px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  background: #111111;
  border-right: 1px solid #1f1f1f;
  height: 100%;
  overflow: hidden;
}

.sidebar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 16px 12px;
  flex-shrink: 0;
}

.sidebar-logo {
  display: flex;
  align-items: center;
  gap: 8px;
}

.logo-icon {
  width: 26px;
  height: 26px;
  border-radius: 7px;
  background: #6366f1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.logo-text {
  font-size: 14px;
  font-weight: 600;
  color: #e5e5e5;
  letter-spacing: -0.01em;
}

.close-btn {
  display: none;
  width: 28px;
  height: 28px;
  border-radius: 6px;
  border: none;
  background: transparent;
  color: #666;
  cursor: pointer;
  align-items: center;
  justify-content: center;
  transition: background 0.15s, color 0.15s;
}
.close-btn:hover { background: #222; color: #aaa; }

.new-chat-btn {
  margin: 0 12px 12px;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 9px 12px;
  border-radius: 8px;
  border: 1px solid #2a2a2a;
  background: transparent;
  color: #a3a3a3;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.15s, color 0.15s, border-color 0.15s;
  text-align: left;
  flex-shrink: 0;
  font-family: inherit;
}
.new-chat-btn:hover:not(:disabled) {
  background: #1e1e1e;
  color: #e5e5e5;
  border-color: #333;
}
.new-chat-btn:disabled { opacity: 0.4; cursor: not-allowed; }

.chat-list {
  flex: 1;
  overflow-y: auto;
  padding: 0 8px 16px;
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.chat-list::-webkit-scrollbar { width: 4px; }
.chat-list::-webkit-scrollbar-track { background: transparent; }
.chat-list::-webkit-scrollbar-thumb { background: #2a2a2a; border-radius: 2px; }

.skeleton {
  height: 40px;
  border-radius: 8px;
  background: linear-gradient(90deg, #1a1a1a 25%, #222 50%, #1a1a1a 75%);
  background-size: 200% 100%;
  animation: shimmer 1.4s infinite;
  margin-bottom: 2px;
}
@keyframes shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

.empty-list {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 40px 16px;
  color: #555;
  font-size: 13px;
}

@media (max-width: 768px) {
  .sidebar {
    position: fixed;
    left: 0; top: 0; bottom: 0;
    z-index: 50;
    transform: translateX(-100%);
    transition: transform 0.25s cubic-bezier(0.4, 0, 0.2, 1);
    width: 280px;
  }
  .sidebar--open { transform: translateX(0); }
  .close-btn { display: flex; }
}
</style>
