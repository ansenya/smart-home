<script setup lang="ts">
import ChatListItem from '@/components/ChatListItem.vue'
import BaseButton from '@/components/BaseButton.vue'
import SquarePlusIcon from '@/assets/icons/SquarePlusIcon.vue'
</script>

<template>
  <aside class="sidebar">
    <BaseButton variant="secondary" class="new-chat-btn"
      >New Chat<SquarePlusIcon class="icon"
    /></BaseButton>

    <div class="chat-list">
      <ChatListItem v-for="chat in chats" :key="chat.id" :chat="chat" />
    </div>
  </aside>
</template>

<style scoped>
.sidebar {
  display: flex;
  flex-direction: column;
  padding: 12px;
  gap: 12px;
  width: 280px;
  background: #f4f4f4;
  border-radius: 4px;
}
.icon {
  width: 24px;
  height: 24px;
}
</style>
