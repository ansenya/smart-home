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
        <h2 class="text-2xl font-bold text-slate-900 mb-1">Аккаунт создан</h2>
        <p class="text-slate-500 text-sm">{{ user.email }}</p>
      </div>
      <button @click="handleAuthorize" class="btn-primary w-full mb-3">
        Продолжить
      </button>
      <button @click="handleLogout" class="btn-ghost w-full">
        Выйти из аккаунта
      </button>
    </div>

    <!-- Register form -->
    <div v-else>
      <div class="mb-8">
        <h2 class="text-2xl font-bold text-slate-900 mb-1.5">Создать аккаунт</h2>
        <p class="text-slate-500 text-sm">Зарегистрируйтесь, чтобы управлять умным домом</p>
      </div>

      <form @submit.prevent="handleRegister" class="space-y-5">
        <div>
          <label class="block text-sm font-medium text-slate-700 mb-1.5">Email</label>
          <input v-model="email" type="email" required placeholder="you@example.com" class="input"/>
        </div>

        <div>
          <label class="block text-sm font-medium text-slate-700 mb-1.5">Пароль</label>
          <div class="relative">
            <input v-model="password" :type="showPassword ? 'text' : 'password'"
                   required placeholder="Минимум 8 символов" class="input pr-11"/>
            <button type="button" @click="showPassword = !showPassword"
                    class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 transition-colors">
              <svg v-if="showPassword" class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"/>
              </svg>
              <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                <path stroke-linecap="round" stroke-linejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
              </svg>
            </button>
          </div>
          <div v-if="isPasswordShort" class="flex items-center gap-1.5 mt-1.5">
            <div class="flex gap-0.5">
              <div v-for="i in 4" :key="i" class="h-1 w-6 rounded-full"
                   :class="password.length >= i * 2 ? 'bg-indigo-500' : 'bg-slate-200'"></div>
            </div>
            <span class="text-xs text-slate-400">Минимум 8 символов</span>
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-slate-700 mb-1.5">Подтвердите пароль</label>
          <div class="relative">
            <input v-model="password2" :type="showPassword2 ? 'text' : 'password'"
                   required placeholder="Повторите пароль" class="input pr-11"
                   :class="password2 && password2 !== password ? 'border-red-300 focus:ring-red-500' : ''"/>
            <button type="button" @click="showPassword2 = !showPassword2"
                    class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 transition-colors">
              <svg v-if="showPassword2" class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"/>
              </svg>
              <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                <path stroke-linecap="round" stroke-linejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
              </svg>
            </button>
          </div>
          <p v-if="password2 && password2 !== password" class="text-xs text-red-500 mt-1.5">Пароли не совпадают</p>
        </div>

        <div v-if="error" class="flex items-center gap-2 px-3 py-2.5 rounded-lg bg-red-50 border border-red-100">
          <svg class="w-4 h-4 text-red-500 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
          </svg>
          <span class="text-sm text-red-600">{{ error }}</span>
        </div>

        <button type="submit" class="btn-primary w-full">Зарегистрироваться</button>
      </form>

      <p class="text-center text-sm text-slate-500 mt-6">
        Уже есть аккаунт?
        <router-link :to="`/login${queryString}`" class="text-indigo-600 hover:text-indigo-700 font-medium transition-colors">
          Войти
        </router-link>
      </p>
    </div>
  </AuthLayout>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { me, logout, authorize, register } from '../api/auth'
import AuthLayout from '../components/AuthLayout.vue'

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
const isPasswordShort = computed(() => password.value.length > 0 && password.value.length < MIN_PASSWORD_LEN)

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
    .then(res => { user.value = res.data })
    .catch(err => { if (err.response?.status !== 401) console.error(err) })
})

function handleRegister() {
  if (password.value.length < MIN_PASSWORD_LEN) { error.value = 'Пароль должен быть не менее 8 символов'; return }
  if (password.value !== password2.value) { error.value = 'Пароли не совпадают'; return }
  if (!oauthQueries) return
  error.value = ''
  register({ email: email.value, password: password.value })
    .then(() => me())
    .then(res => { user.value = res.data; password.value = ''; password2.value = '' })
    .catch(err => {
      const data = err.response?.data
      if (err.response?.status === 409) {
        error.value = 'Этот email уже зарегистрирован. Попробуйте войти.'
      } else {
        error.value = data?.error || data?.message || 'Ошибка регистрации'
      }
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
    .then(() => { user.value = null })
    .catch(err => console.error(err))
}
</script>
