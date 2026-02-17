import { createRouter, createWebHistory } from 'vue-router'
import AuthCallbackView from '@/views/AuthCallbackView.vue'
import HomeView from '@/views/HomeView.vue'
import DevicesView from '@/views/DevicesView.vue'
import ChatsView from '@/views/ChatsView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', component: HomeView },
    { path: '/devices', component: DevicesView },
    { path: '/chats', component: ChatsView },
    { path: '/auth/callback', component: AuthCallbackView },
  ],
})

export default router
