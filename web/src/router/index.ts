import {createWebHistory, createRouter} from 'vue-router'
import CallbackPage from "../pages/CallbackPage.vue";
import HomePage from "../pages/HomePage.vue";


const routes = [
    {path: '/', component: HomePage},
    {path: '/auth/callback', component: CallbackPage},
]

export const router = createRouter({
    history: createWebHistory(),
    routes,
})