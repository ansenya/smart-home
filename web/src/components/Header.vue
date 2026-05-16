<script setup lang="ts">
import UserAuth from '@/components/UserAuth.vue'
import { useAuthStore } from '@/stores/auth.ts'

const authStore = useAuthStore()
</script>

<template>
  <header class="header">
    <router-link to="/" class="logo">
      <div class="logo-icon">
        <svg width="14" height="14" fill="currentColor" viewBox="0 0 20 20">
          <path d="M10.707 2.293a1 1 0 00-1.414 0l-7 7a1 1 0 001.414 1.414L4 10.414V17a1 1 0 001 1h2a1 1 0 001-1v-2a1 1 0 011-1h2a1 1 0 011 1v2a1 1 0 001 1h2a1 1 0 001-1v-6.586l.293.293a1 1 0 001.414-1.414l-7-7z"/>
        </svg>
      </div>
      <span>Hiphome</span>
    </router-link>

    <nav class="nav">
      <router-link to="/" class="nav-link" exact-active-class="nav-link--active">Home</router-link>
      <router-link
        to="/devices"
        class="nav-link"
        exact-active-class="nav-link--active"
        :class="{ 'nav-link--disabled': !authStore.isAuthenticated }"
      >Devices</router-link>
      <router-link
        to="/chats"
        class="nav-link"
        exact-active-class="nav-link--active"
        :class="{ 'nav-link--disabled': !authStore.isAuthenticated }"
      >Chats</router-link>
    </nav>

    <UserAuth />
  </header>
</template>

<style scoped>
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  height: 52px;
  flex-shrink: 0;
  border-bottom: 1px solid #1a1a1a;
  background: #0f0f0f;
  position: relative;
}

.logo {
  display: flex;
  align-items: center;
  gap: 8px;
  text-decoration: none;
  color: #e5e5e5;
  font-size: 15px;
  font-weight: 600;
  letter-spacing: -0.01em;
  flex-shrink: 0;
  z-index: 1;
  transition: color 0.15s;
}

.logo:hover {
  color: #fff;
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

.nav {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 4px;
}

.nav-link {
  color: #737373;
  text-decoration: none;
  padding: 6px 12px;
  border-radius: 7px;
  font-size: 14px;
  font-weight: 500;
  transition: background 0.15s, color 0.15s;
  white-space: nowrap;
}

.nav-link:hover:not(.nav-link--disabled) {
  background: #1a1a1a;
  color: #e5e5e5;
}

.nav-link--active {
  color: #e5e5e5;
  background: #1a1a1a;
}

.nav-link--disabled {
  color: #333;
  cursor: not-allowed;
  pointer-events: none;
}

@media (max-width: 600px) {
  .nav { display: none; }
}
</style>
