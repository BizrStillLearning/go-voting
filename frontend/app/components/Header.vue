<script setup>
import { ref, onMounted } from 'vue'
import { Menu, ShieldCheck, UserCircle, Sun, Moon } from 'lucide-vue-next'
import { useAuthStore } from '~/stores/auth'

const isCollapsed = useState('sidebar-collapsed', () => false)
const authStore = useAuthStore()

const username = authStore.user?.username || 'Administrator'
const role = authStore.user?.role || 'Super Admin'

const isDark = ref(false)

const toggleTheme = () => {
  isDark.value = !isDark.value
  if (isDark.value) {
    document.documentElement.classList.add('dark')
    localStorage.setItem('theme', 'dark')
  } else {
    document.documentElement.classList.remove('dark')
    localStorage.setItem('theme', 'light')
  }
}

onMounted(() => {
  if (localStorage.getItem('theme') === 'dark' || (!('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
    isDark.value = true
    document.documentElement.classList.add('dark')
  } else {
    isDark.value = false
    document.documentElement.classList.remove('dark')
  }
})
</script>

<template>
  <header class="bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border-b border-slate-200 dark:border-slate-800 h-20 flex items-center justify-between px-6 z-10 sticky top-0 shadow-sm transition-colors duration-300">

    <div class="flex items-center gap-4">
      <button @click="isCollapsed = !isCollapsed" class="p-2.5 text-slate-500 dark:text-slate-400 hover:text-poster-base dark:hover:text-poster-light hover:bg-red-50 dark:hover:bg-slate-800 rounded-xl transition-all cursor-pointer">
        <Menu :size="24" />
      </button>
    </div>

    <div class="flex items-center gap-2 md:gap-5">
      <button @click="toggleTheme" class="p-2.5 text-slate-500 dark:text-slate-400 hover:text-poster-glow dark:hover:text-poster-glow hover:bg-slate-100 dark:hover:bg-slate-800 rounded-xl transition-all cursor-pointer" title="Ganti Tema">
        <Sun v-if="isDark" :size="20" />
        <Moon v-else :size="20" />
      </button>

      <div class="flex items-center gap-3 border-l border-slate-200 dark:border-slate-800 pl-4 md:pl-5 transition-colors duration-300">
        <div class="text-right hidden md:block">
          <p class="text-sm font-bold text-slate-900 dark:text-white leading-tight transition-colors duration-300">{{ username }}</p>
          <p class="text-[11px] font-medium text-slate-500 dark:text-slate-400 uppercase tracking-wider transition-colors duration-300">{{ role }}</p>
        </div>
        <div class="w-10 h-10 bg-red-50 dark:bg-slate-800 text-poster-base dark:text-poster-light rounded-xl flex items-center justify-center border border-red-100 dark:border-slate-700 transition-colors duration-300">
          <UserCircle :size="24" stroke-width="2" />
        </div>
      </div>
    </div>

  </header>
</template>