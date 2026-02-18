<script setup lang="ts">
import UserAuth from '@/components/UserAuth.vue'
import { useAuthStore } from '@/stores/auth.ts'
import { ref } from 'vue'

const authStore = useAuthStore()
</script>

<template>
  <header class="header">
    <router-link to="/" class="logo">Hiphome</router-link>

    <nav class="nav">
      <router-link to="/" class="nav-link" exact-active-class="exact-active-nav-link">
        Home
      </router-link>
      <router-link
        to="/devices"
        class="nav-link"
        exact-active-class="exact-active-nav-link"
        :class="{ 'nav-link--disabled': !authStore.isAuthenticated }"
      >
        My Devices
      </router-link>
      <router-link
        to="/chats"
        class="nav-link"
        exact-active-class="exact-active-nav-link"
        :class="{ 'nav-link--disabled': !authStore.isAuthenticated }"
      >
        Chats
      </router-link>
    </nav>
    <UserAuth />
  </header>
</template>

<style scoped>
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin: 15px 20px;
  position: relative;
  min-height: 44px;
}

.logo {
  text-decoration: none;
  color: black;
  font-size: 24px;
  font-weight: bold;
  flex-shrink: 0;
  z-index: 1;
  transition: font-size 0.2s;
  text-align: center;
}

.logo:hover {
  font-size: 25px;
  color: darkblue;
}

.nav {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 40px;
}

.nav-link {
  color: #333;
  text-decoration: none;
  padding: 8px 16px;
  border-radius: 4px;
  transition: all 0.3s ease;
  outline: 2px solid transparent;
  white-space: nowrap;
}

.nav-link--disabled {
  color: #999;
  cursor: not-allowed;
}

.nav-link--disabled:hover {
  background-color: transparent;
  color: #999;
}

.exact-active-nav-link {
  color: darkblue;
  outline-color: darkblue;
}

.nav-link:not(.nav-link--disabled):hover {
  background-color: #f0f0f0;
  color: #000;
}
</style>
