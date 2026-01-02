<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import {exchangeCode} from "../api/auth.ts";

const router = useRouter()

const error = ref<string | null>(null)
const loading = ref(true)

onMounted(async () => {

  const params = new URLSearchParams(window.location.search)
  const code = params.get('code')
  const clientId = params.get('client_id')

  if (!code || !clientId) {
    error.value = 'Invalid OAuth callback parameters'
    loading.value = false
    return
  }

  try {
    await exchangeCode({code: code, client_id: clientId})

    // backend должен установить cookie / сессию
    await router.replace('/')
  } catch (e: any) {
    error.value =
        e.response?.data?.message ||
        'Authorization failed. Please try again.'
    loading.value = false
  }
})
</script>

<template>
  <div class="flex items-center justify-center min-h-screen bg-gray-100">
    <div class="bg-white p-6 rounded shadow text-center w-full max-w-sm">
      <p v-if="loading">Authorizing…</p>
      <p v-else-if="error" class="text-red-500">{{ error }}</p>
    </div>
  </div>
</template>
