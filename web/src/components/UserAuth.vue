<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import UserIcon from '@/assets/icons/UserIcon.vue'
import { computed, onMounted, onUnmounted, ref } from 'vue'
import BaseButton from '@/components/BaseButton.vue'
import { push } from 'notivue'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()
const router = useRouter()
const isMenuOpen = ref(false)
const menuContainer = ref<HTMLElement | null>(null)

const shouldHighlight = computed(() => authStore.highlightLogin)

const handleSignUp = async () => {
  push.error('Sorry, not implemented yet')
}

const handleLogin = async () => {
  await authStore.login()
}

const handleLogout = async () => {
  await authStore.logout()
  await router.replace('/')
}

const toggleMenu = () => {
  isMenuOpen.value = !isMenuOpen.value
}

const handleClickOutside = (event: MouseEvent) => {
  if (menuContainer.value && !menuContainer.value.contains(event.target as Node)) {
    isMenuOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<template>
  <div class="user-auth">
    <template v-if="authStore.isLoading">
      <BaseButton disabled loading>Loading...</BaseButton>
    </template>
    <template v-else-if="authStore.isAuthenticated">
      <div class="user-icon">
        <button
          type="button"
          @click="toggleMenu"
          class="user-icon-btn"
          aria-label="Меню пользователя"
          :aria-expanded="isMenuOpen"
        >
          <UserIcon />
        </button>
        <div v-if="isMenuOpen" class="dropdown-content">
          <ul>
            <li>
              <BaseButton @click="handleLogout">Logout</BaseButton>
            </li>
          </ul>
        </div>
      </div>
    </template>
    <template v-else>
      <div class="buttons-container">
        <BaseButton variant="secondary" size="md" @click="handleSignUp">SignUp</BaseButton>
        <BaseButton
          variant="primary"
          size="md"
          @click="handleLogin"
          :highlight="authStore.highlightLogin"
        >
          Login
        </BaseButton>
      </div>
    </template>
  </div>
</template>

<style scoped>
.user-auth {
  display: flex;
  flex-shrink: 0;
  justify-content: end;
  align-self: end;
  margin-top: auto;
  margin-bottom: auto;
}

.buttons-container {
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  justify-content: center;
  align-items: center;
  gap: 0.7rem;
}

.user-icon {
  position: relative;
  display: inline-block;
}

.dropdown-content {
  position: absolute;
  top: 100%;
  right: 0;
  min-width: 100px;
  background-color: white;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  border-radius: 4px;
  padding: 8px 0;
  z-index: 1000;
  margin-top: 4px;
}
.dropdown-content ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.dropdown-content li {
  padding: 4px 12px;
}

.dropdown-content :deep(button) {
  width: 100%;
  text-align: left;
  background: none;
  border: none;
  padding: 8px 12px;
  cursor: pointer;
}
</style>
