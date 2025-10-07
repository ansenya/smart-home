<script setup>

import TextInput from "@/components/TextInput.vue";
import {onMounted, ref, watch} from "vue";
import {useAuth, useAuthValidation} from "@/composables/useAuth";
import {authService} from "@/services/api";
import {useRoute} from "vue-router";
import PasswordInput from "@/components/PasswordInput.vue";

const {errors, resetErrors} = useAuth()
const {validatePassword, validateEmail} = useAuthValidation()

const route = useRoute();
const email = ref('')

// 0 - show form to send link
// 1 - link sent, show "success" text or smth
// 2 - form for inserting code and new password
// 3 - password changed
const pageStatus = ref(0)
const isLoading = ref(false);

const form = ref({
  email: '',
  code: '',
  password: '',
  confirmedPassword: ''
});

const handleResetPasswordRequest = async () => {
  resetErrors()
  isLoading.value = true;

  errors.value.email = validateEmail(form.value.email)

  if (errors.value.email) {
    isLoading.value = false;
    return
  }

  try {
    const response = await authService.resetPasswordRequest({
      email: form.value.email,
    })

    if (!response.ok) {
      handleResetPasswordError(response)
      isLoading.value = false;
      return
    }

    pageStatus.value = 1;
  } catch (error) {
    isLoading.value = false;
    errors.value.form = 'Ошибка соединения с сервером'
    console.error(error)
  }
}

const handleResetPasswordError = (response) => {
  if (response.status === 403) {
    errors.value.form = 'Пользователя с таким email не существует'
  } else {
    errors.value.form = 'Неизвестная ошибка'
    console.error(response)
  }
}

const handleSetNewPassword = async () => {
  resetErrors();
  isLoading.value = true;

  errors.value.password = validatePassword(form.value.password);
  if (form.value.password !== form.value.confirmedPassword) {
    errors.value.confirmedPassword = 'Пароли не совпадают';
    isLoading.value = false;
    return;
  }

  if (errors.value.password) {
    isLoading.value = false;
    return
  }

  try {
    const response = await authService.resetPasswordConfirm({
      code: form.value.code,
      password: form.value.password
    });

    if (!response.ok) {
      errors.value.form = 'Не удалось сбросить пароль';
      isLoading.value = false;
      return;
    }

    pageStatus.value = 3;
  } catch (error) {
    isLoading.value = false;
    errors.value.form = 'Ошибка соединения с сервером';
    console.error(error);
  }
};

onMounted(() => {
  const code = route.query.code;
  email.value = route.query.email;
  if (code && email) {
    form.value.code = code;
    pageStatus.value = 2;
  }
});

watch(() => form.value.email, (newVal) => {
  errors.value.form = ''
  errors.value.email = validateEmail(newVal)
})

watch(() => form.value.password, (newVal) => {
  errors.value.form = ''
  errors.value.password = validatePassword(newVal)
  if (form.value.confirmedPassword !== newVal) {
    errors.value.confirmedPassword = 'Пароли не совпадают'
  } else {
    errors.value.confirmedPassword = ''
  }
})

watch(() => form.value.confirmedPassword, (newVal) => {
  errors.value.form = ''
  if (form.value.password !== newVal) {
    errors.value.confirmedPassword = 'Пароли не совпадают'
  } else {
    errors.value.confirmedPassword = ''
  }
})
</script>

<template>
  <div v-if="pageStatus === 0" class="auth-container">
    <h2>Забыли пароль?</h2>
    <span class="reset-password-hint">Укажите email, чтобы сбросить пароль</span>
    <form @submit.prevent="handleResetPasswordRequest" class="reset-password-form">
      <TextInput
          v-model:model-value="form.email"
          placeholder="mail@mail.ru"
          :error="errors.email"
          :formError="errors.form"
          :src="require('@/assets/mail.svg')"
          required
      />
      <p class="auth-text">Уже есть аккаунт?
        <router-link to="/login" class="auth-link">Войти</router-link>
      </p>
      <button :disabled="isLoading" type="submit" class="reset-password-button">
        {{ isLoading ? 'Отправляем...' : 'Отправить ссылку' }}
      </button>

      <span v-if="errors.form" class="error-text">{{ errors.form }}</span>
    </form>
  </div>

  <div v-if="pageStatus === 1" class="password-reset-message-container">
    <h2>Ссылка отправлена</h2>
    <p>Мы выслали инструкции для восстановления пароля на {{ form.email }}.</p>
    <p>Проверьте папку "Входящие", возможно письмо попало в спам.</p>
    <button @click="pageStatus = 0" class="reset-password-button">Назад</button>
  </div>

  <div v-if="pageStatus === 2" class="auth-container">
    <h2>Введите новый пароль</h2>
    <span class="reset-password-hint">Вы меняете пароль для {{ email }}</span>

    <form @submit.prevent="handleSetNewPassword" class="reset-password-form">
      <PasswordInput
          v-model:model-value="form.password"
          type="password"
          placeholder="Новый пароль"
          :error="errors.password"
          label="Новый пароль"
          required
      />
      <PasswordInput
          v-model:model-value="form.confirmedPassword"
          type="password"
          placeholder="Подтвердите пароль"
          :error="errors.confirmedPassword"
          label="Подтвердите пароль"
          required
      />
      <p class="auth-text">Уже есть аккаунт?
        <router-link to="/login" class="auth-link">Войти</router-link>
      </p>
      <button :disabled="isLoading" type="submit" class="reset-password-button">
        {{ isLoading ? 'Сохраняем...' : 'Сохранить пароль' }}
      </button>
      <span v-if="errors.form" class="error-text">{{ errors.form }}</span>
    </form>
  </div>

  <div v-if="pageStatus === 3" class="password-reset-message-container">
    <h2>Пароль изменён</h2>
    <p>Ваш пароль успешно обновлён. Теперь вы можете войти с новым паролем.</p>
    <router-link to="/login" class="reset-password-button">Войти</router-link>
  </div>
</template>

<style src="@/assets/styles/reset-password.css"></style>
