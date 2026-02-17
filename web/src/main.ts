import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createNotivue } from 'notivue'
import 'notivue/notification.css'
import 'notivue/animations.css'

import App from './App.vue'
import router from './router'

const app = createApp(App)
const notivue = createNotivue({
  position: 'bottom-left',
  limit: 5,
  avoidDuplicates: false,
  teleportTo: false,
  notifications: {
    global: {
      duration: 5000,
    },
  },
})

app.use(createPinia())
app.use(router)
app.use(notivue)
app.mount('#app')
