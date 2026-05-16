<script setup lang="ts">
import Header from '@/components/Header.vue'
import { Notivue, Notification } from 'notivue'
import { onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth.ts'

const authStore = useAuthStore()

onMounted(async () => {
  await authStore.fetchUser()
})
</script>

<template>
  <div class="app-layout">
    <Header />
    <Notivue v-slot="item">
      <Notification :item="item" />
    </Notivue>
    <main class="main-content">
      <router-view />
    </main>
  </div>
</template>

<style>
*, *::before, *::after {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body {
  height: 100vh;
  overflow: hidden;
  background: #0f0f0f;
  color: #e5e5e5;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

::-webkit-scrollbar { width: 4px; height: 4px; }
::-webkit-scrollbar-track { background: transparent; }
::-webkit-scrollbar-thumb { background: #2a2a2a; border-radius: 2px; }
::-webkit-scrollbar-thumb:hover { background: #333; }

.app-layout {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: #0f0f0f;
}

.main-content {
  flex: 1;
  overflow: hidden;
}
</style>
