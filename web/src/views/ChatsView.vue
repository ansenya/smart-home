<script setup lang="ts">
import ChatsSidebar from '@/components/ChatsSidebar.vue'
import ChatsMessages from '@/components/ChatsMessages.vue'
</script>

<template>
  <div class="chats-layout">
    <ChatsSidebar />
    <ChatsMessages />
  </div>
</template>

<style scoped>
.chats-layout {
  display: flex;
  height: 100%;
  overflow: hidden;
  padding-bottom: 4px;
  padding-left: 4px;
}
</style>
