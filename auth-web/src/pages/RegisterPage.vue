<template>
  <div class="flex items-center justify-center min-h-screen bg-gray-100">
    <div class="bg-white p-8 rounded-lg shadow-md w-full max-w-md text-center">

      <h1 class="text-2xl font-bold mb-6">Register</h1>

      <!-- Ошибка query -->
      <p v-if="queryError" class="text-red-500">{{ queryError }}</p>

      <!-- Пользователь уже авторизован -->
      <div v-else-if="user">
        <p class="mb-4">
          Вы зарегистрированы как <strong>{{ user.email }}</strong>
        </p>
        <button @click="handleAuthorize"
                class="w-full bg-green-500 text-white py-2 rounded hover:bg-green-600 transition mb-2">
          Далее
        </button>
        <button @click="handleLogout"
                class="w-full bg-red-500 text-white py-2 rounded hover:bg-red-600 transition">
          Выйти
        </button>
      </div>

      <!-- Форма регистрации -->
      <form v-else @submit.prevent="handleRegister" class="space-y-4">
        <div>
          <label class="block mb-1 font-medium">Email</label>
          <input v-model="email" type="email" required
                 class="w-full border rounded px-3 py-2 focus:ring-2 focus:ring-blue-500"/>
        </div>

        <div>
          <label class="block mb-1 font-medium">Password</label>

          <div class="relative">
            <input
                v-model="password"
                :type="showPassword ? 'text' : 'password'"
                required
                class="w-full border rounded px-3 py-2 pr-10 focus:ring-2 focus:ring-blue-500"
            />

            <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute right-2 top-1/2 -translate-y-1/2 text-sm text-gray-500"
            >
              {{ showPassword ? 'Hide' : 'Show' }}
            </button>
          </div>

          <p v-if="isPasswordShort" class="text-sm text-orange-500 mt-1">
            Password must be at least 8 characters
          </p>
        </div>


        <div>
          <label class="block mb-1 font-medium">Confirm password</label>

          <div class="relative">
            <input
                v-model="password2"
                :type="showPassword2 ? 'text' : 'password'"
                required
                class="w-full border rounded px-3 py-2 pr-10 focus:ring-2 focus:ring-blue-500"
            />

            <button
                type="button"
                @click="showPassword2 = !showPassword2"
                class="absolute right-2 top-1/2 -translate-y-1/2 text-sm text-gray-500"
            >
              {{ showPassword2 ? 'Hide' : 'Show' }}
            </button>
          </div>
        </div>


        <button type="submit"
                class="w-full bg-blue-500 text-white py-2 rounded hover:bg-blue-600 transition">
          Register
        </button>

        <router-link class="text-blue-500 underline block"
                     :to="`/login${queryString}`">
          Уже есть аккаунт?
        </router-link>

        <p v-if="error" class="text-red-500 mt-2">{{ error }}</p>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import {ref, onMounted, computed} from 'vue'
import {me, logout, authorize, register} from '../api/auth'

interface User {
  email: string
}

const queryString = window.location.search

const email = ref('')
const password = ref('')
const password2 = ref('')
const error = ref('')
const user = ref<User | null>(null)
const queryError = ref<string | null>(null)

const showPassword = ref(false)
const showPassword2 = ref(false)


const MIN_PASSWORD_LEN = 8
const isPasswordShort = computed(
    () => password.value.length > 0 && password.value.length < MIN_PASSWORD_LEN
)

let oauthQueries: Record<string, string> | null = null

function validateOAuthQuery(qs: string): Record<string, string> | null {
  const params = new URLSearchParams(qs)
  const required = ['client_id', 'redirect_uri', 'response_type', 'state']
  const missing = required.filter(k => !params.has(k))

  if (missing.length) {
    queryError.value = `Некорректные query параметры: отсутствуют ${missing.join(', ')}`
    return null
  }

  const result: Record<string, string> = {}
  params.forEach((v, k) => result[k] = v)
  return result
}

onMounted(() => {
  oauthQueries = validateOAuthQuery(window.location.search)
  if (!oauthQueries) return

  me()
      .then(res => user.value = res.data)
      .catch(err => {
        if (err.response?.status !== 401) console.error(err)
      })
})

// ---------------------
// Actions
// ---------------------
function handleRegister() {
  if (password.value.length < MIN_PASSWORD_LEN) {
    error.value = 'Password must be at least 8 characters'
    return
  }

  if (!oauthQueries) return
  error.value = ''

  if (password.value !== password2.value) {
    error.value = 'Passwords do not match'
    return
  }

  register({email: email.value, password: password.value})
      .then(() => me())
      .then(res => {
        user.value = res.data
        password.value = ''
        password2.value = ''
      })
      .catch(err => {
        error.value = err.response?.data?.message || 'Registration failed'
      })
}

function handleAuthorize() {
  if (!oauthQueries) return
  authorize(oauthQueries)
      .then(res => window.location.replace(res.data.redirect_url))
      .catch(err => console.error(err))
}

function handleLogout() {
  logout()
      .then(() => user.value = null)
      .catch(err => console.error(err))
}
</script>
