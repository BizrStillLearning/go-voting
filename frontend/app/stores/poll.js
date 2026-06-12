import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useAuthStore } from './auth'

export const usePollStore = defineStore('poll', () => {
    const polls = ref([])
    const isLoading = ref(false)

    const fetchPolls = async () => {
        isLoading.value = true

        const authStore = useAuthStore()

        try {
            const response = await $fetch('http://localhost:8080/api/admin/polls', {
                method: 'GET',
                headers: {
                    Authorization: `Bearer ${authStore.token}`
                }
            })

            polls.value = response.data || []
        } catch (error) {
            console.error("Gagal mengambil data polling:", error)
        } finally {
            isLoading.value = false
        }
    }

    return { polls, isLoading, fetchPolls }
})