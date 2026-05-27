<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { push } from 'notivue'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()
const router = useRouter()
const isMenuOpen = ref(false)
const menuContainer = ref<HTMLElement | null>(null)

const handleLogin = async () => await authStore.login()

const handleLogout = async () => {
  await authStore.logout()
  await router.replace('/')
}

const handleClickOutside = (event: MouseEvent) => {
  if (menuContainer.value && !menuContainer.value.contains(event.target as Node)) {
    isMenuOpen.value = false
  }
}

onMounted(() => document.addEventListener('click', handleClickOutside))
onUnmounted(() => document.removeEventListener('click', handleClickOutside))
</script>

<template>
  <div class="user-auth">
    <!-- Loading -->
    <div v-if="authStore.isLoading" class="skeleton-avatar" />

    <!-- Logged in -->
    <div v-else-if="authStore.isAuthenticated" ref="menuContainer" class="user-menu-wrap">
      <button class="avatar-btn" @click="isMenuOpen = !isMenuOpen" :aria-expanded="isMenuOpen">
        <svg width="16" height="16" fill="none" stroke="currentColor" stroke-width="1.8" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
        </svg>
      </button>
      <div v-if="isMenuOpen" class="dropdown">
        <button class="dropdown-item dropdown-item--danger" @click="handleLogout">
          <svg width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"/>
          </svg>
          Sign out
        </button>
      </div>
    </div>

    <!-- Guest -->
    <div v-else class="guest-btns">
      <button class="btn-login" :class="{ 'btn-login--highlight': authStore.highlightLogin }" @click="handleLogin">
        Sign in
      </button>
    </div>
  </div>
</template>

<style scoped>
.user-auth {
  display: flex;
  align-items: center;
  flex-shrink: 0;
}

.skeleton-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: #1e1e1e;
  animation: pulse 1.4s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.4; }
}

/* Avatar button */
.user-menu-wrap {
  position: relative;
}

.avatar-btn {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: 1px solid #2a2a2a;
  background: #1a1a1a;
  color: #737373;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.15s, color 0.15s, border-color 0.15s;
}

.avatar-btn:hover {
  background: #222;
  color: #e5e5e5;
  border-color: #333;
}

/* Dropdown */
.dropdown {
  position: absolute;
  top: calc(100% + 6px);
  right: 0;
  background: #1a1a1a;
  border: 1px solid #2a2a2a;
  border-radius: 10px;
  padding: 4px;
  min-width: 140px;
  z-index: 1000;
  box-shadow: 0 8px 24px rgba(0,0,0,0.5);
}

.dropdown-item {
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
  font-weight: 500;
  cursor: pointer;
  text-align: left;
  font-family: inherit;
  transition: background 0.12s, color 0.12s;
}

.dropdown-item:hover {
  background: #222;
  color: #e5e5e5;
}

.dropdown-item--danger {
  color: #f87171;
}

.dropdown-item--danger:hover {
  background: #1f0f0f;
  color: #fca5a5;
}

/* Login button */
.guest-btns {
  display: flex;
  gap: 8px;
  align-items: center;
}

.btn-login {
  padding: 6px 14px;
  border-radius: 8px;
  border: 1px solid #2a2a2a;
  background: #1a1a1a;
  color: #a3a3a3;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  font-family: inherit;
  transition: background 0.15s, color 0.15s, border-color 0.15s;
}

.btn-login:hover {
  background: #222;
  color: #e5e5e5;
  border-color: #333;
}

.btn-login--highlight {
  background: #7f1d1d;
  border-color: #991b1b;
  color: #fca5a5;
  animation: shake 0.5s ease-in-out;
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  20%, 60% { transform: translateX(-4px); }
  40%, 80% { transform: translateX(4px); }
}
</style>
