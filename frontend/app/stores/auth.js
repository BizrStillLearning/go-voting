import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAuthStore = defineStore('auth', () => {
    const token = useCookie('auth_token', { maxAge: 60 * 60 * 24 })

    const isAuthenticated = ref(!!token.value)

    const login = async (username, password) => {
        try {
            const response = await $fetch('http://localhost:8080/api/admin/login', {
                method: 'POST',
                body: { username, password }
            })

            token.value = response.token
            isAuthenticated.value = true

            return { success: true }
        } catch (error) {
            return {
                success: false,
                message: error.data?.error || 'Gagal terhubung ke server backend'
            }
        }
    }

    const logout = () => {
        token.value = null
        isAuthenticated.value = false
    }

    return { token, isAuthenticated, login, logout }
})