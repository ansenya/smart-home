<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth.ts'
import router from '@/router'

const route = useRoute()
const authStore = useAuthStore()

onMounted(async () => {
  const code = route.query.code as string
  const state = route.query.state as string

  if (code && state) {
    const success = await authStore.finishLogin(code, state)
    if (success) {
      await router.replace('/')
    } else {
    }
  } else {
    await router.replace('/')
  }
})
</script>

<template></template>

<style scoped></style>
