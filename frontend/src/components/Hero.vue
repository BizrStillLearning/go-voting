<script setup>
import { Rocket, Zap, Users2, Vote, ChevronRight } from 'lucide-vue-next'

defineEmits(['start'])
</script>

<template>
  <section class="relative overflow-hidden bg-slate-900 rounded-[3rem] mx-4 md:mx-0 p-8 md:p-20 text-white shadow-2xl mb-12">
    <div class="absolute top-0 right-0 -translate-y-24 translate-x-24 w-96 h-96 bg-indigo-500/20 blur-[120px] rounded-full"></div>
    <div class="absolute bottom-0 left-0 translate-y-24 -translate-x-24 w-72 h-72 bg-blue-500/10 blur-[100px] rounded-full"></div>

    <div class="relative z-10 grid grid-cols-1 lg:grid-cols-2 gap-16 items-center">
      <div class="space-y-8">
        <div class="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-indigo-500/10 border border-indigo-400/20 text-indigo-300 text-xs font-black uppercase tracking-[0.2em]">
          <Zap :size="14" class="fill-indigo-400" /> Trusted Voting System
        </div>

        <h2 class="text-5xl md:text-7xl font-black leading-[1.1] tracking-tight">
          Create a poll <br />
          <span class="text-transparent bg-clip-text bg-gradient-to-r from-indigo-400 to-blue-400">
            in seconds
          </span>
        </h2>

        <p class="text-slate-400 text-lg md:text-xl leading-relaxed max-w-lg">
          Want to ask your friends where to go or manage a formal election?
          Create a poll - and get answers in no time. **No signup required.**
        </p>

        <div class="flex flex-wrap gap-4 pt-4">
          <button @click="$emit('start')" class="px-10 py-5 bg-white text-slate-900 font-black rounded-2xl hover:bg-indigo-50 hover:scale-105 transition-all flex items-center gap-3 shadow-2xl shadow-white/10">
            Create a poll <Rocket :size="20" />
          </button>
          <button class="px-10 py-5 bg-slate-800/50 text-white font-black rounded-2xl border border-slate-700 hover:bg-slate-800 transition-all flex items-center gap-2">
            Live Demo <ChevronRight :size="20" />
          </button>
        </div>
      </div>

      <div class="grid grid-cols-1 gap-6">
        <div class="flex items-center justify-around p-8 bg-slate-800/30 backdrop-blur-md rounded-[2.5rem] border border-slate-700/50">
          <div class="text-center">
            <div class="text-3xl md:text-4xl font-black mb-1">2.3M+</div>
            <div class="text-[10px] text-slate-500 font-black uppercase tracking-[0.2em]">Users</div>
          </div>
          <div class="w-px h-12 bg-slate-700"></div>
          <div class="text-center">
            <div class="text-3xl md:text-4xl font-black mb-1">13M+</div>
            <div class="text-[10px] text-slate-500 font-black uppercase tracking-[0.2em]">Polls</div>
          </div>
          <div class="w-px h-12 bg-slate-700"></div>
          <div class="text-center">
            <div class="text-3xl md:text-4xl font-black mb-1">290M+</div>
            <div class="text-[10px] text-slate-500 font-black uppercase tracking-[0.2em]">Votes</div>
          </div>
        </div>

        <div class="p-8 bg-indigo-600 rounded-[2.5rem] flex items-center justify-between group cursor-pointer hover:bg-indigo-500 transition-colors">
          <div>
            <h4 class="font-bold text-xl">Ready to start?</h4>
            <p class="text-indigo-200 text-sm">Join millions of users worldwide.</p>
          </div>
          <div class="w-12 h-12 bg-white/20 rounded-full flex items-center justify-center group-hover:rotate-45 transition-transform">
            <Rocket :size="24" />
          </div>
        </div>
      </div>
    </div>
  </section>
</template>