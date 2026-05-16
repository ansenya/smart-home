<script setup lang="ts">
import WrenchIcon from '@/assets/icons/WrenchIcon.vue'
import { onMounted, onUnmounted, ref } from 'vue'
import PairingModal from '@/components/modals/PairingModal.vue'
import { push } from 'notivue'
import { useAuthStore } from '@/stores/auth.ts'
import { useFloatingButtonDevicesButtonStore } from '@/stores/floatingButtonDevicesButtonStore.ts'

const authStore = useAuthStore()
const floatingButtonStore = useFloatingButtonDevicesButtonStore()

const isOpen = ref(false)
const isPairingOpen = ref(false)
const menuRef = ref<HTMLElement | null>(null)
const isShaking = ref(false)

const toggleMenu = () => {
  floatingButtonStore.incr()
  if (!authStore.isAuthenticated) {
    isShaking.value = true
    setTimeout(() => { isShaking.value = false }, 500)
    push.error({ title: 'Action prohibited', message: 'You are not authenticated' })
  } else {
    isOpen.value = !isOpen.value
  }
}

const closeMenu = () => { isOpen.value = false }
const openPairing = () => { isPairingOpen.value = true; closeMenu() }

const handleClickOutside = (e: MouseEvent) => {
  if (menuRef.value && !menuRef.value.contains(e.target as Node)) closeMenu()
}
const handleEscape = (e: KeyboardEvent) => {
  if (e.key === 'Escape' && isOpen.value) closeMenu()
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  document.addEventListener('keydown', handleEscape)
})
onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  document.removeEventListener('keydown', handleEscape)
})
</script>

<template>
  <div class="floating-container" ref="menuRef">
    <transition name="menu-fade">
      <div v-if="isOpen" class="floating-menu">
        <button class="menu-item" @click="openPairing">
          <svg width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
          </svg>
          Add Device
        </button>
      </div>
    </transition>

    <button
      class="fab"
      :class="{ 'fab--shake': isShaking }"
      @click="toggleMenu"
      :aria-expanded="isOpen"
      aria-label="Open Menu"
    >
      <WrenchIcon
        class="fab-icon"
        :class="{ 'fab-icon--open': isOpen, 'fan-spin': floatingButtonStore.isFan }"
      />
    </button>

    <PairingModal v-model="isPairingOpen" @success="() => push.info('Pairing success')" />
  </div>
</template>

<style scoped>
.floating-container {
  position: fixed;
  bottom: 24px;
  right: 24px;
  z-index: 1000;
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 10px;
}

.fab {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  border: 1px solid #2a2a2a;
  background: #1a1a1a;
  color: #737373;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 16px rgba(0,0,0,0.4);
  transition: background 0.15s, color 0.15s, transform 0.15s, box-shadow 0.15s;
}

.fab:hover:not(.fab--shake) {
  background: #222;
  color: #e5e5e5;
  transform: scale(1.05);
  box-shadow: 0 6px 20px rgba(0,0,0,0.5);
}

.fab--shake {
  background: #7f1d1d;
  color: #fca5a5;
  border-color: #991b1b;
  animation: shake 0.5s ease-in-out;
}

.fab-icon {
  transition: transform 0.25s ease;
  width: 20px;
  height: 20px;
}

.fab-icon--open {
  transform: rotate(-90deg);
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  20%, 60% { transform: translateX(-4px); }
  40%, 80% { transform: translateX(4px); }
}

@keyframes fan-spin {
  to { transform: rotate(360deg); }
}
.fan-spin { animation: fan-spin 0.1s linear infinite; }

/* Floating menu popup */
.floating-menu {
  background: #1a1a1a;
  border: 1px solid #2a2a2a;
  border-radius: 12px;
  padding: 4px;
  box-shadow: 0 8px 32px rgba(0,0,0,0.5);
  min-width: 160px;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  padding: 9px 12px;
  border-radius: 8px;
  border: none;
  background: transparent;
  color: #a3a3a3;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  font-family: inherit;
  transition: background 0.12s, color 0.12s;
  text-align: left;
}

.menu-item:hover {
  background: #222;
  color: #e5e5e5;
}

.menu-fade-enter-active,
.menu-fade-leave-active {
  transition: opacity 0.15s, transform 0.15s;
}
.menu-fade-enter-from,
.menu-fade-leave-to {
  opacity: 0;
  transform: translateY(8px) scale(0.96);
}
</style>
