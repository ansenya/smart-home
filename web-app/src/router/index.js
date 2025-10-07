import {createRouter, createWebHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'
import ResetPassword from '../views/ResetPassword.vue'

const routes = [
    {
        path: '/',
        name: 'home',
        component: HomeView,
        meta: {title: 'Hopahome | Home Page'}
    },
    {
        path: '/login',
        name: 'login',
        component: LoginView,
        meta: {title: 'Hopahome | Login'}
    },
    {
        path: '/register',
        name: 'register',
        component: RegisterView,
        meta: {title: 'Hopahome | Register'}
    },
    {
        path: '/reset-password',
        name: 'reset-password',
        component: ResetPassword,
        meta: {title: 'Hopahome | Reset Password'}

    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

router.beforeEach((to, from, next) => {
    document.title = to.meta.title || 'Hopahome | Default Title';
    next();
});

export default router