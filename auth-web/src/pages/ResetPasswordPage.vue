<template>
  <AuthLayout>
    <div v-if="sent">
      <div class="mb-8">
        <div class="w-12 h-12 rounded-2xl bg-emerald-50 flex items-center justify-center mb-4">
          <svg class="w-6 h-6 text-emerald-600" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
          </svg>
        </div>
        <h2 class="text-2xl font-bold text-slate-900 mb-1.5">Письмо отправлено</h2>
        <p class="text-slate-500 text-sm leading-relaxed">
          Если адрес <span class="font-medium text-slate-700">{{ submittedEmail }}</span>
          зарегистрирован, на него отправлено письмо со ссылкой для сброса пароля.
          Ссылка действует 30 минут.
        </p>
      </div>
      <router-link :to="`/login${queryString}`" class="btn-ghost w-full text-center block">
        Назад ко входу
      </router-link>
    </div>

    <div v-else>
      <div class="mb-8">
        <h2 class="text-2xl font-bold text-slate-900 mb-1.5">Сброс пароля</h2>
        <p class="text-slate-500 text-sm">Введите email — мы пришлём ссылку для смены пароля.</p>
      </div>

      <form @submit.prevent="handleSubmit" class="space-y-5">
        <div>
          <label class="block text-sm font-medium text-slate-700 mb-1.5">Email</label>
          <input v-model="email" type="email" required placeholder="you@example.com" class="input" :disabled="pending"/>
        </div>

        <div v-if="error" class="flex items-center gap-2 px-3 py-2.5 rounded-lg bg-red-50 border border-red-100">
          <svg class="w-4 h-4 text-red-500 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
          </svg>
          <span class="text-sm text-red-600">{{ error }}</span>
        </div>

        <button type="submit" class="btn-primary w-full" :disabled="pending">
          {{ pending ? 'Отправка…' : 'Отправить ссылку' }}
        </button>
      </form>

      <p class="text-center text-sm text-slate-500 mt-6">
        Вспомнили пароль?
        <router-link :to="`/login${queryString}`" class="text-indigo-600 hover:text-indigo-700 font-medium transition-colors">
          Войти
        </router-link>
      </p>
    </div>
  </AuthLayout>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { requestPasswordReset } from '../api/auth'
import AuthLayout from '../components/AuthLayout.vue'

const queryString = window.location.search
const email = ref('')
const submittedEmail = ref('')
const error = ref('')
const sent = ref(false)
const pending = ref(false)

function handleSubmit() {
  error.value = ''
  pending.value = true
  const target = email.value
  requestPasswordReset(target)
    .then(() => { submittedEmail.value = target; sent.value = true })
    .catch(err => {
      const data = err.response?.data
      error.value = data?.error || 'Не удалось отправить запрос'
    })
    .finally(() => { pending.value = false })
}
</script>
