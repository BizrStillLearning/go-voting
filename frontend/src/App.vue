<script setup>
import { useRouter, useRoute } from 'vue-router'
import Navbar from './components/Navbar.vue'

const router = useRouter()
const route = useRoute()

const handleSetTab = (tabName) => {
  router.push({ name: tabName })
}
</script>

<template>
  <div class="min-h-screen bg-[#F8FAFC] dark:bg-slate-950 transition-colors duration-500">
    <Navbar :activeTab="route.name" @setTab="handleSetTab" />

    <main class="max-w-7xl mx-auto p-4 md:p-8 pt-32 md:pt-40">
      <router-view v-slot="{ Component }">
        <transition
            mode="out-in"
            enter-active-class="transition duration-300 ease-out"
            enter-from-class="opacity-0 translate-y-4"
            enter-to-class="opacity-100 translate-y-0"
        >
          <component :is="Component" />
        </transition>
      </router-view>
    </main>
  </div>
</template>