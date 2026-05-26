import {createWebHistory, createRouter} from 'vue-router'

import LoginPage from '../pages/LoginPage.vue'
import RegisterPage from "../pages/RegisterPage.vue";
import ResetPasswordPage from "../pages/ResetPasswordPage.vue";
import ResetPasswordConfirmPage from "../pages/ResetPasswordConfirmPage.vue";

const routes = [
    {path: '/', redirect: (to: { fullPath: string }) => ({path: '/login', query: parseQuery(to.fullPath)})},
    {path: '/login', component: LoginPage},
    {path: '/register', component: RegisterPage},
    {path: '/reset-password', component: ResetPasswordPage},
    {path: '/reset-password/confirm', component: ResetPasswordConfirmPage},
    {path: '/:pathMatch(.*)*', redirect: '/login'},
]

function parseQuery(fullPath: string): Record<string, string> {
    const i = fullPath.indexOf('?')
    if (i < 0) return {}
    const params = new URLSearchParams(fullPath.slice(i + 1))
    const out: Record<string, string> = {}
    params.forEach((v, k) => { out[k] = v })
    return out
}

export const router = createRouter({
    history: createWebHistory(),
    routes,
})