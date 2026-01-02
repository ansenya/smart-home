import {createWebHistory, createRouter} from 'vue-router'
import CallbackPage from "../pages/CallbackPage.vue";


const routes = [
    {path: '/', redirect: '/'},
    {path: '/auth/callback', component: CallbackPage},
]

export const router = createRouter({
    history: createWebHistory(),
    routes,
})