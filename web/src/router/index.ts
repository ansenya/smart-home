import {createWebHistory, createRouter} from 'vue-router'
import CallbackPage from "../pages/CallbackPage.vue";
import HomePage from "../pages/HomePage.vue";
import AddDevicePage from "../pages/AddDevicePage.vue";
import PairingCallbackPage from "../pages/PairingCallbackPage.vue";


const routes = [
    {path: '/', component: HomePage},
    {path: '/auth/callback', component: CallbackPage},
    {path: '/devices/add', component: AddDevicePage},
    {path: '/pairing/callback', component: PairingCallbackPage},
]

export const router = createRouter({
    history: createWebHistory(),
    routes,
})