<script setup lang="ts">
import ChatInput from '@/components/ChatInput.vue'
import ChatMessage from '@/components/ChatMessage.vue'
</script>

<template>
  <main class="chat-area">
    <div class="messages-container">
      <ChatMessage v-for="message in messages" :key="message.id" :message="message" />
    </div>

    <ChatInput />
  </main>
</template>

<style scoped>
.chat-area {
  display: flex;
  flex-direction: column;
  flex: 1;
}
.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}
.input-container {
  padding: 24px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}
.btn:hover {
  opacity: 0.8;
}
</style>
