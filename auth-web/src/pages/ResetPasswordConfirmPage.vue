<template>
  <AuthLayout>
    <!-- Missing token -->
    <div v-if="!token" class="text-center">
      <div class="w-12 h-12 rounded-2xl bg-red-50 flex items-center justify-center mx-auto mb-4">
        <svg class="w-6 h-6 text-red-500" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
        </svg>
      </div>
      <h2 class="text-xl font-bold text-slate-900 mb-1.5">Некорректная ссылка</h2>
      <p class="text-slate-500 text-sm">В ссылке отсутствует токен сброса.</p>
      <router-link to="/login" class="btn-ghost w-full mt-6 text-center block">
        Назад ко входу
      </router-link>
    </div>

    <!-- Success -->
    <div v-else-if="done">
      <div class="mb-8">
        <div class="w-12 h-12 rounded-2xl bg-emerald-50 flex items-center justify-center mb-4">
          <svg class="w-6 h-6 text-emerald-600" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"/>
          </svg>
        </div>
        <h2 class="text-2xl font-bold text-slate-900 mb-1.5">Пароль обновлён</h2>
        <p class="text-slate-500 text-sm">Войдите в аккаунт с новым паролем.</p>
      </div>
      <router-link to="/login" class="btn-primary w-full text-center block">
        Перейти ко входу
      </router-link>
    </div>

    <!-- Form -->
    <div v-else>
      <div class="mb-8">
        <h2 class="text-2xl font-bold text-slate-900 mb-1.5">Новый пароль</h2>
        <p class="text-slate-500 text-sm">Придумайте новый пароль для входа.</p>
      </div>

      <form @submit.prevent="handleSubmit" class="space-y-5">
        <div>
          <label class="block text-sm font-medium text-slate-700 mb-1.5">Пароль</label>
          <input v-model="password" type="password" required placeholder="Минимум 8 символов" class="input" :disabled="pending"/>
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
          <input v-model="password2" type="password" required placeholder="Повторите пароль" class="input" :disabled="pending"
                 :class="password2 && password2 !== password ? 'border-red-300 focus:ring-red-500' : ''"/>
          <p v-if="password2 && password2 !== password" class="text-xs text-red-500 mt-1.5">Пароли не совпадают</p>
        </div>

        <div v-if="error" class="flex items-center gap-2 px-3 py-2.5 rounded-lg bg-red-50 border border-red-100">
          <svg class="w-4 h-4 text-red-500 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
          </svg>
          <span class="text-sm text-red-600">{{ error }}</span>
        </div>

        <button type="submit" class="btn-primary w-full" :disabled="pending">
          {{ pending ? 'Сохранение…' : 'Сохранить пароль' }}
        </button>
      </form>
    </div>
  </AuthLayout>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { confirmPasswordReset } from '../api/auth'
import AuthLayout from '../components/AuthLayout.vue'

const params = new URLSearchParams(window.location.search)
const token = params.get('token') || ''

const password = ref('')
const password2 = ref('')
const error = ref('')
const pending = ref(false)
const done = ref(false)

const MIN_PASSWORD_LEN = 8
const isPasswordShort = computed(() => password.value.length > 0 && password.value.length < MIN_PASSWORD_LEN)

function handleSubmit() {
  if (password.value.length < MIN_PASSWORD_LEN) { error.value = 'Пароль должен быть не менее 8 символов'; return }
  if (password.value !== password2.value) { error.value = 'Пароли не совпадают'; return }
  error.value = ''
  pending.value = true
  confirmPasswordReset(token, password.value)
    .then(() => { done.value = true })
    .catch(err => {
      const data = err.response?.data
      if (data?.error === 'invalid or expired reset link') {
        error.value = 'Ссылка недействительна или истекла. Запросите новую.'
      } else if (data?.error === 'password too weak') {
        error.value = 'Слишком простой пароль.'
      } else {
        error.value = data?.error || 'Не удалось сохранить пароль'
      }
    })
    .finally(() => { pending.value = false })
}
</script>
