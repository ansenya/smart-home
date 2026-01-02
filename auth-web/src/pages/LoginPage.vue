<template>
  <div class="flex items-center justify-center min-h-screen bg-gray-100">
    <div class="bg-white p-8 rounded-lg shadow-md w-full max-w-md text-center">

      <h1 class="text-2xl font-bold mb-6">Login</h1>

      <!-- Ошибка в query -->
      <p v-if="queryError" class="text-red-500">{{ queryError }}</p>

      <!-- Пользователь авторизован -->
      <div v-else-if="user">
        <p class="mb-4">Вы авторизованы по почте <strong>{{ user.email }}</strong></p>
        <button @click="handleAuthorize"
                class="w-full bg-green-500 text-white py-2 rounded hover:bg-green-600 transition mb-2">
          Далее
        </button>
        <button @click="handleLogout"
                class="w-full bg-red-500 text-white py-2 rounded hover:bg-red-600 transition">
          Выйти
        </button>
      </div>

      <!-- Пользователь не авторизован -->
      <form v-else @submit.prevent="handleLogin" class="space-y-4">
        <div>
          <label class="block mb-1 font-medium">Email</label>
          <input v-model="email" type="text" required
                 class="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"/>
        </div>
        <div>
          <label class="block mb-1 font-medium">Password</label>
          <input v-model="password" type="password" required
                 class="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"/>
        </div>
        <button type="submit"
                class="w-full bg-blue-500 text-white py-2 rounded hover:bg-blue-600 transition">
          Log In
        </button>
        <div>
          <router-link class="text-blue-500 underline" :to="`/reset-password${queryString}`">Forgot password?</router-link>
          <br>
          <router-link class="text-blue-500 underline" :to="`/register${queryString}`">Register</router-link>
        </div>
        <p v-if="error" class="text-red-500 mt-2">{{ error }}</p>
      </form>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {ref, onMounted} from 'vue'
import {me, login, logout, authorize} from '../api/auth'

interface User {
  email: string
}

const queryString = window.location.search

const email = ref('')
const password = ref('')
const error = ref('')
const user = ref<User | null>(null)
const queryError = ref<string | null>(null)

// ---------------------
// Helper: проверка query на корректность
// ---------------------
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

// ---------------------
// Проверка при монтировании
// ---------------------
let oauthQueries: Record<string, string> | null = null
onMounted(() => {
  const qs = window.location.search.startsWith('?') ? window.location.search : ''
  oauthQueries = validateOAuthQuery(qs)
  if (!oauthQueries) return // query неправильные → показываем только ошибку

  // Проверка текущей сессии
  me()
      .then(res => {
        user.value = res.data
      })
      .catch(err => {
        if (err.response?.status !== 401) console.error('Error checking auth:', err)
      })
})

// ---------------------
// Actions
// ---------------------
function handleLogin() {
  if (!oauthQueries) return // защита от неправильного query
  error.value = ''

  login({email: email.value, password: password.value})
      .then(() => me())
      .then(res => {
        user.value = res.data
        email.value = ''
        password.value = ''
      })
      .catch(err => {
        error.value = err.response?.data?.message || 'Login failed'
      })
}

function handleAuthorize() {
  if (!oauthQueries) return
  authorize(oauthQueries)
      .then(res => {
        if (res.data.redirect_url) window.location.replace(res.data.redirect_url)
        else console.error('No redirect_url in response')
      })
      .catch(err => console.error('Authorization error:', err))
}

function handleLogout() {
  if (!oauthQueries) return
  logout()
      .then(() => {
        user.value = null
        email.value = ''
        password.value = ''
      })
      .catch(err => console.error('Logout error:', err))
}
</script>
