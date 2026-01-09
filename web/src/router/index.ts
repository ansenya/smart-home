import {createWebHistory, createRouter} from 'vue-router'
import CallbackPage from "../pages/CallbackPage.vue";
import HomePage from "../pages/HomePage.vue";
import AddDevicePage from "../pages/AddDevicePage.vue";
import AddDeviceWaitPage from "../pages/AddDeviceWaitPage.vue";


const routes = [
    {path: '/', component: HomePage},
    {path: '/auth/callback', component: CallbackPage},
    {path: '/devices/add', component: AddDevicePage},
    {path: '/devices/add/wait', component: AddDeviceWaitPage},
]

export const router = createRouter({
    history: createWebHistory(),
    routes,
})