import {createWebHistory, createRouter} from 'vue-router'
import CallbackPage from "../pages/CallbackPage.vue";
import HomePage from "../pages/HomePage.vue";
import PairingPage from "../pages/PairingPage.vue";

const routes = [
    {path: '/', component: HomePage},
    {path: '/auth/callback', component: CallbackPage},
    {path: '/devices/start-pairing', component: PairingPage},
]

export const router = createRouter({
    history: createWebHistory(),
    routes,
})