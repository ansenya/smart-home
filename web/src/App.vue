<script setup lang="ts">
import Header from '@/components/Header.vue'
import { Notivue, Notification } from 'notivue'
import { onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth.ts'

const authStore = useAuthStore()

onMounted(async () => {
  await authStore.fetchUser()
})
</script>

<template>
  <div>
    <Header></Header>
    <Notivue v-slot="item">
      <Notification :item="item" />
    </Notivue>
  </div>
  <router-view></router-view>
</template>

<style scoped>
</style>
