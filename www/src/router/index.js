import { createRouter, createWebHistory } from 'vue-router'

export default createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/dashboard',
            component: () => import('@/views/dashboard.vue')
        },
        {
            path: '/login',
            component: () => import('@/views/login.vue')
        },
        {
            path: '/',
            component: () => import('@/views/index.vue')
        }
    ]
})
