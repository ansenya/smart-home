import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { authService } from '@/services/auth.service'
import { push } from 'notivue'

interface User {
  ID: string
  name: string
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const isLoading = ref(true)

  const isAuthenticated = computed(() => !!user.value)

  async function fetchUser() {
    isLoading.value = true
    try {
      user.value = await authService.getMe()
    } catch (e: any) {
      user.value = null
      if (e.response?.status !== 401) {
        push.error({
          title: 'Failed to authenticate',
          message: e.response?.data?.message || e.message || 'Unknown error',
        })
      }
    } finally {
      isLoading.value = false
    }
  }

  async function login() {
    await authService.initiateLogin()
  }

  async function finishLogin(code: string, state: string) {
    isLoading.value = true
    try {
      await authService.handleCallback(code, state)
      await fetchUser()
      push.success('Authenticated successfully')
      return true
    } catch (e: any) {
      push.error('Authentication error: ' + (e.message || 'Unknown error'))
      return false
    } finally {
      isLoading.value = false
    }
  }

  async function logout() {
    await authService.logout()
    user.value = null
    push.warning('You are logged out')
  }

  return { user, isAuthenticated, isLoading, fetchUser, login, finishLogin, logout }
})
