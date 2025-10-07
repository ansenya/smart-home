<script setup>

import PasswordInput from "@/components/PasswordInput.vue";
import TextInput from "@/components/TextInput.vue";
import {ref, watch} from "vue";
import {useAuth, useAuthValidation} from "@/composables/useAuth";
import {authService} from "@/services/api";
import router from "@/router";
import Cookies from 'js-cookie';

const {validateEmail, validatePassword} = useAuthValidation()

const {errors, resetErrors} = useAuth()

const form = ref({
  email: '',
  password: '',
  confirmPassword: ''
})

const handleRegistration = async () => {
  resetErrors()

  errors.value.email = validateEmail(form.value.email)
  errors.value.password = validatePassword(form.value.password)

  if (errors.value.email || errors.value.password) return

  try {
    const response = await authService.register({
      email: form.value.email,
      password: form.value.password
    })

    if (!response.ok) {
      handleRegistrationError(response)
      return
    }

    const data = await response.json()
    handleRegistrationSuccess(data)

    errors.value.form = 'Регистрация успешна! Проверьте вашу почту'
  } catch (error) {
    errors.value.form = 'Ошибка соединения с сервером'
    console.error(error)
  }
}

const handleRegistrationError = (response) => {
  if (response.status === 409) {
    errors.value.form = 'Пользователь с таким email уже зарегистрирован'
  } else {
    errors.value.form = 'Неизвестная ошибка'
    console.error(response)
  }
}

const handleRegistrationSuccess = (data) => {
  Cookies.set('access_token', data.access_token, {
    expires: new Date(Date.now() + data.expires_in * 1000),
    secure: true, // Use HTTPS
    sameSite: 'Strict',
    path: '/'
  });

  Cookies.set('refresh_token', data.refresh_token, {
    expires: 7, // Days
    secure: true,
    sameSite: 'Strict',
    path: '/'
  });

  if (data.redirect_url) {
    window.location.href = data.redirect_url
  } else {
    router.push('/');
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

watch(() => form.value.confirmPassword, (newVal) => {
  errors.value.form = ''
  errors.value.confirmedPassword = validatePassword(newVal)
})
</script>

<template>
  <div class="auth-container">
    <h2>Регистрация</h2>
    <form @submit.prevent="handleRegistration" class="auth-form">
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
          placeholder="password"
          required
      />

      <PasswordInput
          v-model:model-value="form.confirmPassword"
          :error="errors.confirmedPassword"
          :formError="errors.form"
          placeholder="password again"
          required
      />

      <p class="auth-text">Уже есть аккаунт?
        <router-link to="/login" class="auth-link">Войти</router-link>
      </p>

      <button type="submit" class="auth-button">Зарегистрироваться</button>

      <span v-if="errors.form" class="error-text">{{ errors.form }}</span>
    </form>
  </div>
</template>

<style src="@/assets/styles/auth.css"></style>