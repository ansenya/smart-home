<script setup>
import {ref, watch} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import TextInput from "@/components/TextInput.vue";
import PasswordInput from "@/components/PasswordInput.vue";
import Cookies from 'js-cookie';

import {useAuth, useAuthValidation} from "@/composables/useAuth";
import {authService} from "@/services/api";

const {validateEmail, validatePassword} = useAuthValidation()
const {errors, resetErrors} = useAuth()

const router = useRouter()
const route = useRoute()

const queryParams = route.query
const decodedParams = {
  state: decodeURIComponent(queryParams.state || ''),
  scope: decodeURIComponent(queryParams.scope || ''),
  redirect_uri: decodeURIComponent(queryParams.redirect_uri || ''),
  response_type: queryParams.response_type || '',
  client_id: queryParams.client_id || ''
}

const form = ref({
  email: '',
  password: ''
})

const handleLogin = async () => {
  resetErrors()

  errors.value.email = validateEmail(form.value.email)
  errors.value.password = validatePassword(form.value.password)

  if (errors.value.email || errors.value.password) return

  const params = decodedParams.state && decodedParams.redirect_uri &&
  decodedParams.response_type && decodedParams.client_id ? {
    'smart-home': 'true',
    state: decodedParams.state,
    scope: decodedParams.scope,
    client_id: decodedParams.client_id
  } : null
  try {
    const response = await authService.login({
      email: form.value.email,
      password: form.value.password
    }, params)

    if (!response.ok) {
      handleLoginError(response)
      return
    }

    const data = await response.json()
    handleLoginSuccess(data)
  } catch (error) {
    errors.value.form = 'Ошибка соединения с сервером'
    console.error('Error:', error)
  }
}

const handleLoginError = (response) => {
  if (response.status === 401) {
    errors.value.form = 'Неверная почта или пароль'
  } else {
    errors.value.form = 'Неизвестная ошибка'
    console.error(response)
  }
}

const handleLoginSuccess = (data) => {
  Cookies.set('access_token', data.access_token, {
    expires: new Date(Date.now() + data.expires_in * 1000), // Convert seconds to milliseconds
    secure: true, // Use HTTPS
    sameSite: 'Strict',
    path: '/'
  });

  // Set refresh token cookie (expires in 7 days)
  Cookies.set('refresh_token', data.refresh_token, {
    expires: 7, // Days
    secure: true,
    sameSite: 'Strict',
    path: '/'
  });

  if (data.redirect_url) {
    window.location.href = data.redirect_url
  } else {
    router.push('/')
  }
}

watch(() => form.value.email, (newVal) => {
  errors.value.form = ''
  errors.value.email = validateEmail(newVal)
})

watch(() => form.value.password, (newVal) => {
  errors.value.form = ''
  errors.value.password = validatePassword(newVal)
})
</script>

<template>
  <div class="auth-container">
    <h2>Авторизация</h2>
    <form @submit.prevent="handleLogin" class="auth-form">
      <TextInput
          v-model:model-value="form.email"
          placeholder="mail@mail.ru"
          :error="errors.email"
          :formError="errors.form"
          :src="require('@/assets/mail.svg')"
          required
      />

      <PasswordInput
          v-model:model-value="form.password"
          :error="errors.password"
          :formError="errors.form"
          required
      />

      <router-link to="/reset-password" class="auth-link">Забыли пароль?</router-link>

      <p class="auth-text">Нет аккаунта?
        <router-link to="/register" class="auth-link">Зарегистрироваться</router-link>
      </p>

      <button type="submit" class="auth-button">Авторизоваться</button>

      <span v-if="errors.form" class="error-text">{{ errors.form }}</span>
    </form>
  </div>
</template>

<style src="@/assets/styles/auth.css"></style>
