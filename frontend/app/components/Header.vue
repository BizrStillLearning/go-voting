<script setup>
import { Search, Bell, User } from 'lucide-vue-next'
</script>

<template>
  <header class="h-20 bg-white/80 backdrop-blur-md border-b border-slate-200 flex items-center justify-between px-8 sticky top-0 z-40">

    <div class="hidden md:flex items-center bg-slate-100 px-4 py-2.5 rounded-2xl w-96 border border-slate-200 focus-within:border-indigo-500 focus-within:bg-white transition-all">
      <Search :size="18" class="text-slate-400 mr-3" />
      <input type="text" placeholder="Cari data polling atau token..." class="bg-transparent border-none outline-none w-full text-sm font-medium text-slate-700">
    </div>

    <div class="flex items-center gap-6 ml-auto">
      <button class="relative p-2 text-slate-400 hover:text-indigo-600 transition-colors cursor-pointer">
        <Bell :size="22" />
        <span class="absolute top-1 right-1 w-2.5 h-2.5 bg-red-500 rounded-full border-2 border-white"></span>
      </button>

      <div class="h-8 w-px bg-slate-200"></div>

      <div class="flex items-center gap-3 cursor-pointer group">
        <div class="text-right hidden sm:block">
          <p class="text-sm font-black text-slate-800 group-hover:text-indigo-600 transition-colors">Admin Utama</p>
          <p class="text-xs font-bold text-slate-400">Superuser</p>
        </div>
        <div class="bg-indigo-100 p-2.5 rounded-xl text-indigo-600 group-hover:bg-indigo-600 group-hover:text-white transition-all">
          <User :size="20" />
        </div>
      </div>
    </div>
  </header>
</template>