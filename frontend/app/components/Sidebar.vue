<script setup>
import { LayoutDashboard, Database, Settings, LogOut, Vote } from 'lucide-vue-next'
import { useRouter } from 'nuxt/app'
import Swal from 'sweetalert2'

const router = useRouter()
const isCollapsed = useState('sidebar-collapsed')

const handleLogout = () => {
  const isDarkMode = document.documentElement.classList.contains('dark')

  Swal.fire({
    title: 'Keluar Sistem?',
    text: "Sesi Anda akan diakhiri dan butuh otentikasi ulang.",
    icon: 'warning',
    showCancelButton: true,
    background: isDarkMode ? '#0f172a' : '#ffffff',
    color: isDarkMode ? '#f8fafc' : '#1e293b',
    confirmButtonColor: '#A51013', // poster-dark
    cancelButtonColor: isDarkMode ? '#334155' : '#f1f5f9',
    cancelButtonText: `<span style="color: ${isDarkMode ? '#cbd5e1' : '#64748b'};">BATAL</span>`,
    confirmButtonText: 'YA, KELUAR',
    reverseButtons: true,
    customClass: { popup: 'rounded-2xl' }
  }).then((result) => {
    if (result.isConfirmed) {
      const authCookie = useCookie('admin_token')
      authCookie.value = null
      window.location.href = '/login'
    }
  })
}
</script>

<template>
  <aside
      class="bg-white dark:bg-slate-950 text-slate-800 dark:text-white flex flex-col h-screen border-r border-slate-200 dark:border-slate-800 transition-all duration-300 ease-in-out z-20 relative shadow-sm"
      :class="isCollapsed ? 'w-20' : 'w-72'"
  >
    <div class="h-20 flex items-center border-b border-slate-200 dark:border-slate-800 px-6 overflow-hidden whitespace-nowrap transition-colors duration-300">
      <Vote class="text-poster-base dark:text-poster-light flex-shrink-0" :size="28" />
      <span class="ml-4 text-xl font-bold tracking-wider transition-opacity duration-300" :class="isCollapsed ? 'opacity-0' : 'opacity-100'">
        E-VOTING
      </span>
    </div>

    <nav class="flex-1 py-6 flex flex-col gap-2 px-4 overflow-x-hidden">
      <NuxtLink to="/dashboard" active-class="bg-red-50 dark:bg-poster-base/10 text-poster-base dark:text-poster-light" class="flex items-center px-3 py-3 rounded-xl text-slate-500 dark:text-slate-400 font-medium hover:bg-slate-50 dark:hover:bg-slate-900 hover:text-slate-900 dark:hover:text-slate-200 transition-colors group">
        <LayoutDashboard :size="22" class="flex-shrink-0" />
        <span class="ml-4 transition-opacity duration-300 whitespace-nowrap" :class="isCollapsed ? 'opacity-0 w-0' : 'opacity-100'">Dashboard</span>
      </NuxtLink>

      <NuxtLink to="/polls" active-class="bg-red-50 dark:bg-poster-base/10 text-poster-base dark:text-poster-light" class="flex items-center px-3 py-3 rounded-xl text-slate-500 dark:text-slate-400 font-medium hover:bg-slate-50 dark:hover:bg-slate-900 hover:text-slate-900 dark:hover:text-slate-200 transition-colors group">
        <Database :size="22" class="flex-shrink-0" />
        <span class="ml-4 transition-opacity duration-300 whitespace-nowrap" :class="isCollapsed ? 'opacity-0 w-0' : 'opacity-100'">Data Polling</span>
      </NuxtLink>

      <NuxtLink to="/settings" active-class="bg-red-50 dark:bg-poster-base/10 text-poster-base dark:text-poster-light" class="flex items-center px-3 py-3 rounded-xl text-slate-500 dark:text-slate-400 font-medium hover:bg-slate-50 dark:hover:bg-slate-900 hover:text-slate-900 dark:hover:text-slate-200 transition-colors group">
        <Settings :size="22" class="flex-shrink-0" />
        <span class="ml-4 transition-opacity duration-300 whitespace-nowrap" :class="isCollapsed ? 'opacity-0 w-0' : 'opacity-100'">Pengaturan</span>
      </NuxtLink>
    </nav>

    <div class="p-4 border-t border-slate-200 dark:border-slate-800 transition-colors duration-300">
      <button @click="handleLogout" class="w-full flex items-center px-3 py-3 rounded-xl text-red-500 dark:text-red-400 font-medium hover:bg-red-50 dark:hover:bg-red-500/10 transition-colors group">
        <LogOut :size="22" class="flex-shrink-0 group-hover:-translate-x-1 transition-transform" />
        <span class="ml-4 transition-opacity duration-300 whitespace-nowrap" :class="isCollapsed ? 'opacity-0 w-0' : 'opacity-100'">Keluar</span>
      </button>
    </div>
  </aside>
</template>