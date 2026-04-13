<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useThemeStore } from '../store/theme'
import {
  LayoutDashboard, UserPlus, Vote,
  Sun, Moon
} from 'lucide-vue-next'

const themeStore = useThemeStore()
const props = defineProps(['activeTab'])
const emit = defineEmits(['setTab'])

const isScrolled = ref(false)
const handleScroll = () => { isScrolled.value = window.scrollY > 20 }

onMounted(() => window.addEventListener('scroll', handleScroll))
onUnmounted(() => window.removeEventListener('scroll', handleScroll))
</script>

<template>
  <nav
      :class="[
      isScrolled
        ? 'py-3 bg-white/70 dark:bg-slate-950/70 shadow-lg border-slate-200/50 dark:border-slate-800/50'
        : 'py-6 bg-transparent border-transparent',
      'fixed top-0 inset-x-0 z-50 transition-all duration-500 backdrop-blur-md border-b px-4 md:px-8'
    ]"
  >
    <div class="max-w-7xl mx-auto flex items-center justify-between gap-4">

      <div @click="$emit('setTab', 'home')" class="flex items-center gap-3 cursor-pointer group shrink-0 transition-transform" :class="isScrolled ? 'scale-90' : 'scale-100'">
        <div class="bg-indigo-600 p-2.5 rounded-2xl shadow-lg shadow-indigo-500/20 group-hover:rotate-6 transition-transform">
          <LayoutDashboard :size="isScrolled ? 20 : 24" class="text-white" />
        </div>
        <h1 class="text-xl font-black tracking-tighter text-slate-900 dark:text-white hidden sm:block">Go-Voting</h1>
      </div>

      <div class="flex bg-slate-200/50 dark:bg-slate-800/40 p-1 rounded-2xl border border-white/50 dark:border-white/5 shadow-inner transition-transform" :class="isScrolled ? 'scale-95' : 'scale-100'">
        <button @click="$emit('setTab', 'manage')" :class="activeTab === 'manage' ? 'bg-white dark:bg-slate-700 shadow-md text-indigo-600 dark:text-indigo-400' : 'text-slate-500 dark:text-slate-400'" class="px-5 py-2 rounded-xl text-xs font-black transition-all flex items-center gap-2">
          <UserPlus :size="16" /> <span class="hidden xs:block">Manage</span>
        </button>
        <button @click="$emit('setTab', 'vote')" :class="activeTab === 'vote' ? 'bg-white dark:bg-slate-700 shadow-md text-indigo-600 dark:text-indigo-400' : 'text-slate-500 dark:text-slate-400'" class="px-5 py-2 rounded-xl text-xs font-black transition-all flex items-center gap-2">
          <Vote :size="16" /> <span class="hidden xs:block">Voting</span>
        </button>
      </div>

      <button
          @click="themeStore.toggleTheme"
          class="p-3 rounded-2xl bg-white/50 dark:bg-slate-800/50 border border-slate-200 dark:border-slate-700 text-slate-600 dark:text-slate-300 hover:scale-110 active:scale-95 transition-all"
      >
        <Transition mode="out-in" name="rotate-fade">
          <Sun v-if="themeStore.isDark" :size="20" key="sun" />
          <Moon v-else :size="20" key="moon" />
        </Transition>
      </button>

    </div>
  </nav>
</template>

<style scoped>
.rotate-fade-enter-active, .rotate-fade-leave-active {
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}
.rotate-fade-enter-from { opacity: 0; transform: rotate(-90deg) scale(0.5); }
.rotate-fade-leave-to { opacity: 0; transform: rotate(90deg) scale(0.5); }
</style>