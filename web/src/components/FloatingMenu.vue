<script setup lang="ts">
import BaseButton from '@/components/BaseButton.vue'
import WrenchIcon from '@/assets/icons/WrenchIcon.vue'
import { onMounted, onUnmounted, ref } from 'vue'
import PairingModal from '@/components/modals/PairingModal.vue'
import { push } from 'notivue'
import { useAuthStore } from '@/stores/auth.ts'

const authStore = useAuthStore()

const isOpen = ref(false)
const isPairingOpen = ref(false)
const menuRef = ref<HTMLElement | null>(null)
const isShaking = ref(false)

const toggleMenu = () => {
  if (!authStore.isAuthenticated) {
    isShaking.value = true
    setTimeout(() => {
      isShaking.value = false
    }, 500)

    push.error({
      title: 'Action prohibited',
      message: 'You are not authenticated',
    })
  } else {
    isOpen.value = !isOpen.value
  }
}

const closeMenu = () => {
  isOpen.value = false
}

const openPairing = () => {
  isPairingOpen.value = true
  closeMenu()
}

const handleClickOutside = (event: MouseEvent) => {
  if (menuRef.value && !menuRef.value.contains(event.target as Node)) {
    closeMenu()
  }
}

const handleEscape = (event: KeyboardEvent) => {
  if (event.key === 'Escape' && isOpen.value) {
    closeMenu()
  }
}

const handlePairingSuccess = () => {
  push.info('Pairing success')
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
        <ul>
          <li>
            <BaseButton variant="ghost" size="sm" @click="openPairing">
              <span class="floating-menu-text">Add Device</span>
            </BaseButton>
          </li>
        </ul>
      </div>
    </transition>

    <BaseButton
      variant="primary"
      size="icon"
      @click="toggleMenu"
      :aria-expanded="isOpen"
      aria-label="Open Menu"
      class="floating-btn"
      :class="{ shake: isShaking }"
    >
      <WrenchIcon class="icon" :class="{ 'rotate-icon': isOpen }" />
    </BaseButton>

    <PairingModal v-model="isPairingOpen" @success="handlePairingSuccess" />
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
  gap: 12px;
}
.floating-btn {
  width: 56px;
  height: 56px;
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);
  transition:
    transform 0.2s,
    box-shadow 0.2s,
    background-color 0.1s;
}

.floating-btn:not(.shake):hover {
  transform: scale(1.05);
  box-shadow: 0 6px 16px rgba(59, 130, 246, 0.5);
}

.floating-btn.shake:hover {
  background-color: darkred;
}

@keyframes shake {
  0%,
  100% {
    transform: translateX(0);
  }
  10%,
  30%,
  50%,
  70%,
  90% {
    transform: translateX(-4px);
  }
  20%,
  40%,
  60%,
  80% {
    transform: translateX(4px);
  }
}

.shake {
  animation: shake 0.5s ease-in-out;
  background-color: darkred;
}
.icon {
  transition: transform 0.3s ease;
}
.rotate-icon {
  transform: rotate(-90deg);
  transition: transform 0.3s ease;
  font-size: 24px;
  line-height: 1;
}
.floating-menu {
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  border: 1px solid #e5e7eb;
  overflow: hidden;
  min-width: 160px;
}
.floating-menu ul {
  list-style: none;
  padding: 8px 0;
  margin: 0;
}
.floating-menu li {
  padding: 0 8px;
}
.floating-menu-text {
  width: 100%;
  font-size: 14px;
  text-align: center;
}
.floating-menu li button {
  width: 100%;
  justify-content: flex-start;
  border-radius: 8px;
}
.menu-fade-enter-active,
.menu-fade-leave-active {
  transition:
    opacity 0.2s,
    transform 0.2s;
}
.menu-fade-enter-from,
.menu-fade-leave-to {
  opacity: 0;
  transform: translateY(10px) scale(0.95);
}
</style>
