import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import AuthCallbackView from '@/views/AuthCallbackView.vue'
import HomeView from '@/views/HomeView.vue'
import DevicesView from '@/views/DevicesView.vue'
import ChatsView from '@/views/ChatsView.vue'
import { push } from 'notivue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', component: HomeView },
    { path: '/devices', component: DevicesView, meta: { requiresAuth: true } },
    { path: '/chats', component: ChatsView, meta: { requiresAuth: true } },
    { path: '/auth/callback', component: AuthCallbackView },
  ],
})

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    push.error({
      title: 'Access denied',
      message: 'Please login to access this page',
    })
    authStore.triggerLoginHighlight()
  } else {
    next()
  }
})

export default router
