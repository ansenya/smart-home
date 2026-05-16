<template>
  <AuthLayout>
    <!-- OAuth error -->
    <div v-if="queryError" class="text-center">
      <div class="w-12 h-12 rounded-2xl bg-red-50 flex items-center justify-center mx-auto mb-4">
        <svg class="w-6 h-6 text-red-500" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
        </svg>
      </div>
      <p class="text-sm text-slate-500">{{ queryError }}</p>
    </div>

    <!-- Already logged in -->
    <div v-else-if="user">
      <div class="mb-8">
        <div class="w-12 h-12 rounded-2xl bg-indigo-50 flex items-center justify-center mb-4">
          <svg class="w-6 h-6 text-indigo-600" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
          </svg>
        </div>
        <h2 class="text-2xl font-bold text-slate-900 mb-1">Вы вошли</h2>
        <p class="text-slate-500 text-sm">{{ user.email }}</p>
      </div>
      <button @click="handleAuthorize" class="btn-primary w-full mb-3">
        Продолжить
      </button>
      <button @click="handleLogout" class="btn-ghost w-full">
        Выйти из аккаунта
      </button>
    </div>

    <!-- Login form -->
    <div v-else>
      <div class="mb-8">
        <h2 class="text-2xl font-bold text-slate-900 mb-1.5">Добро пожаловать</h2>
        <p class="text-slate-500 text-sm">Войдите в аккаунт, чтобы продолжить</p>
      </div>

      <form @submit.prevent="handleLogin" class="space-y-5">
        <div>
          <label class="block text-sm font-medium text-slate-700 mb-1.5">Email</label>
          <input v-model="email" type="email" required placeholder="you@example.com" class="input"/>
        </div>

        <div>
          <div class="flex items-center justify-between mb-1.5">
            <label class="block text-sm font-medium text-slate-700">Пароль</label>
            <router-link :to="`/reset-password${queryString}`"
                         class="text-xs text-indigo-600 hover:text-indigo-700 transition-colors">
              Забыли пароль?
            </router-link>
          </div>
          <input v-model="password" type="password" required placeholder="••••••••" class="input"/>
        </div>

        <div v-if="error" class="flex items-center gap-2 px-3 py-2.5 rounded-lg bg-red-50 border border-red-100">
          <svg class="w-4 h-4 text-red-500 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
          </svg>
          <span class="text-sm text-red-600">{{ error }}</span>
        </div>

        <button type="submit" class="btn-primary w-full">Войти</button>
      </form>

      <p class="text-center text-sm text-slate-500 mt-6">
        Нет аккаунта?
        <router-link :to="`/register${queryString}`" class="text-indigo-600 hover:text-indigo-700 font-medium transition-colors">
          Зарегистрироваться
        </router-link>
      </p>
    </div>
  </AuthLayout>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { me, login, logout, authorize } from '../api/auth'
import AuthLayout from '../components/AuthLayout.vue'

interface User {
  email: string
}

const queryString = window.location.search
const email = ref('')
const password = ref('')
const error = ref('')
const user = ref<User | null>(null)
const queryError = ref<string | null>(null)

let oauthQueries: Record<string, string> | null = null

function validateOAuthQuery(qs: string): Record<string, string> | null {
  const params = new URLSearchParams(qs)
  const required = ['client_id', 'redirect_uri', 'response_type', 'state']
  const missing = required.filter(key => !params.has(key))
  if (missing.length) {
    queryError.value = `Некорректные query параметры: отсутствуют ${missing.join(', ')}`
    return null
  }
  const result: Record<string, string> = {}
  params.forEach((value, key) => result[key] = value)
  return result
}

onMounted(() => {
  const qs = window.location.search.startsWith('?') ? window.location.search : ''
  oauthQueries = validateOAuthQuery(qs)
  if (!oauthQueries) return
  me()
    .then(res => { user.value = res.data })
    .catch(err => { if (err.response?.status !== 401) console.error('Error checking auth:', err) })
})

function handleLogin() {
  if (!oauthQueries) return
  error.value = ''
  login({ email: email.value, password: password.value })
    .then(() => me())
    .then(res => {
      user.value = res.data
      email.value = ''
      password.value = ''
    })
    .catch(err => { error.value = err.response?.data?.message || 'Ошибка входа' })
}

function handleAuthorize() {
  if (!oauthQueries) return
  authorize(oauthQueries)
    .then(res => { if (res.data.redirect_url) window.location.replace(res.data.redirect_url) })
    .catch(err => console.error('Authorization error:', err))
}

function handleLogout() {
  if (!oauthQueries) return
  logout()
    .then(() => { user.value = null; email.value = ''; password.value = '' })
    .catch(err => console.error('Logout error:', err))
}
</script>
