import { useAuthStore } from '~/stores/auth'
import { defineNuxtRouteMiddleware, navigateTo } from 'nuxt/app'

export default defineNuxtRouteMiddleware((to, from) => {
    const authCookie = useCookie('admin_token')

    if (!authCookie.value && to.path !== '/login') {
        return navigateTo('/login')
    }

    if (authCookie.value && to.path === '/login') {
        return navigateTo('/dashboard')
    }
})