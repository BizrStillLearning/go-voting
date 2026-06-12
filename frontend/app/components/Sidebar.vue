<script setup>
import { LayoutDashboard, CheckSquare, Users, Settings, LogOut } from 'lucide-vue-next'
import { useRouter } from 'nuxt/app'
import { useAuthStore } from '~/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const handleLogout = () => {
  if (confirm("Apakah Anda yakin ingin keluar dari sistem?")) {
    authStore.logout()
    router.push('/login')
  }
}
</script>

<template>
  <aside class="w-72 bg-slate-900 text-slate-300 flex flex-col h-screen transition-all hidden md:flex">

    <div class="h-20 flex items-center px-8 border-b border-slate-800">
      <div class="bg-indigo-500 p-2 rounded-xl text-white mr-3 shadow-lg shadow-indigo-500/30">
        <CheckSquare :size="24" />
      </div>
      <h1 class="text-2xl font-black text-white tracking-tight">GoVoting</h1>
    </div>

    <nav class="flex-1 px-4 py-8 space-y-2 overflow-y-auto">
      <p class="px-4 text-xs font-black text-slate-500 uppercase tracking-widest mb-4">Menu Utama</p>

      <NuxtLink to="/dashboard" active-class="bg-indigo-600 text-white shadow-md" class="flex items-center gap-3 px-4 py-3 hover:bg-slate-800 hover:text-white rounded-xl font-bold transition-all">
        <LayoutDashboard :size="20" /> Dashboard
      </NuxtLink>

      <NuxtLink to="/polls" active-class="bg-indigo-600 text-white shadow-md" class="flex items-center gap-3 px-4 py-3 hover:bg-slate-800 hover:text-white rounded-xl font-bold transition-all">
        <CheckSquare :size="20" /> Data Polling
      </NuxtLink>

      <NuxtLink to="/settings" active-class="bg-indigo-600 text-white shadow-md" class="flex items-center gap-3 px-4 py-3 hover:bg-slate-800 hover:text-white rounded-xl font-bold transition-all">
        <Settings :size="20" /> Pengaturan
      </NuxtLink>
    </nav>

    <div class="p-4 border-t border-slate-800">
      <button @click="handleLogout" class="w-full flex items-center justify-center gap-2 px-4 py-3 bg-slate-800 hover:bg-red-500 hover:text-white text-slate-400 rounded-xl font-bold transition-all cursor-pointer">
        <LogOut :size="18" /> Keluar Sistem
      </button>
    </div>
  </aside>
</template>

