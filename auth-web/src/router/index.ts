import {createWebHistory, createRouter} from 'vue-router'

import LoginPage from '../pages/LoginPage.vue'
import RegisterPage from "../pages/RegisterPage.vue";

const routes = [
    {path: '/', redirect: '/login'},
    {path: '/login', component: LoginPage},
    {path: '/register', component: RegisterPage},
]

export const router = createRouter({
    history: createWebHistory(),
    routes,
})