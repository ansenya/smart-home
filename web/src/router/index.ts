import {createWebHistory, createRouter} from 'vue-router'
import CallbackPage from "../pages/CallbackPage.vue";
import HomePage from "../pages/HomePage.vue";
import PairingButton from "../components/PairingButton.vue";

const routes = [
    {path: '/', component: HomePage},
    {path: '/auth/callback', component: CallbackPage},
    {path: '/devices/start-pairing', component: PairingButton},
]

export const router = createRouter({
    history: createWebHistory(),
    routes,
})