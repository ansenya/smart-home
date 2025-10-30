import {createWebHistory, createRouter} from 'vue-router'

import LoginPage from '../pages/LoginPage.vue'

const routes = [
    {path: '/', redirect: '/login'},
    {path: '/login', component: LoginPage},
]

export const router = createRouter({
    history: createWebHistory(),
    routes,
})