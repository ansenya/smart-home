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
    { path: '/chats/:chatId?', component: ChatsView, meta: { requiresAuth: true } },
    { path: '/auth/callback', component: AuthCallbackView },
  ],
})

router.beforeEach(async (to, _from, next) => {
  const authStore = useAuthStore()

  // Wait for the initial auth check to finish before evaluating guards.
  // Without this, the guard sees isAuthenticated=false on every page reload
  // because fetchUser() hasn't completed yet.
  await authStore.ready

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    push.error({
      title: 'Access denied',
      message: 'Please login to access this page',
    })
    authStore.triggerLoginHighlight()
    next('/')
  } else {
    next()
  }
})

export default router
