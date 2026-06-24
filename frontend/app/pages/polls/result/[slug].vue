<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'nuxt/app'
import { Trophy, AlertTriangle, ArrowLeft, Activity } from 'lucide-vue-next'

definePageMeta({ layout: false })

const route = useRoute()
const slug = route.params.slug
const poll = ref(null)
const error = ref(false)
let pollingInterval = null

const getImageUrl = (path) => {
  if (!path) return ''
  if (path.startsWith('http')) return path

  const host = typeof window !== 'undefined' ? window.location.hostname : 'localhost'
  const prefix = path.startsWith('/') ? '' : '/'
  return `http://${host}:8080${prefix}${path}`
}

const fetchLiveResults = async () => {
  try {
    const host = typeof window !== 'undefined' ? window.location.hostname : 'localhost'
    const response = await $fetch(`http://${host}:8080/api/poll/${slug}/results`)
    poll.value = response.poll || response.data || response
  } catch (err) {
    error.value = true
  }
}

onMounted(() => {
  fetchLiveResults()
  pollingInterval = setInterval(fetchLiveResults, 3000)
})

onUnmounted(() => {
  if (pollingInterval) clearInterval(pollingInterval)
})

const calculatePercentage = (voteCount, totalVotes) => {
  if (!totalVotes || totalVotes === 0) return 0
  return Math.round((voteCount / totalVotes) * 100)
}

const getPositionTotal = (options) => {
  if (!options) return 0
  return options.reduce((sum, opt) => sum + opt.vote_count, 0)
}
</script>

<template>
  <div class="min-h-screen bg-slate-50 dark:bg-slate-950 text-slate-800 dark:text-slate-200 font-sans selection:bg-poster-base/20 dark:selection:bg-poster-base/40 pb-12 transition-colors duration-300">

    <Header />

    <div v-if="error" class="flex flex-col items-center justify-center h-[calc(100vh-80px)] text-center p-6">
      <div class="bg-red-50 dark:bg-poster-base/10 p-6 rounded-full mb-6">
        <AlertTriangle class="text-poster-base dark:text-poster-light" :size="48" />
      </div>
      <h1 class="text-2xl font-bold mb-2 text-slate-900 dark:text-white">Data Tidak Ditemukan</h1>
      <p class="text-slate-500 dark:text-slate-400 max-w-md">Pemilihan ini mungkin sudah dihapus, ditutup, atau tautan yang Anda masukkan tidak valid.</p>
      <NuxtLink to="/dashboard" class="mt-6 px-6 py-2 bg-poster-dark dark:bg-poster-base text-white rounded-lg font-medium hover:bg-poster-dark/90 dark:hover:bg-poster-light transition-colors">
        Kembali ke Dashboard
      </NuxtLink>
    </div>

    <div v-else-if="!poll" class="flex flex-col items-center justify-center h-[calc(100vh-80px)]">
      <div class="w-12 h-12 border-4 border-poster-base dark:border-poster-light border-t-transparent dark:border-t-transparent rounded-full animate-spin mb-4"></div>
      <p class="text-lg font-medium text-slate-500 dark:text-slate-400 animate-pulse">Memuat data pemilihan...</p>
    </div>

    <div v-else class="max-w-4xl mx-auto p-4 md:p-8 pt-8 md:pt-12">
      <header class="mb-12 relative text-center md:text-left flex flex-col md:flex-row items-center md:items-start justify-between gap-6">
        <div>
          <div class="inline-flex items-center gap-2 px-3 py-1 bg-red-50 dark:bg-poster-base/10 text-poster-base dark:text-poster-light rounded-full text-xs font-bold uppercase tracking-wider mb-4 border border-red-100 dark:border-poster-base/20">
            <span class="w-2 h-2 rounded-full bg-poster-base dark:bg-poster-light animate-pulse"></span> Live Quick Count
          </div>
          <h1 class="text-3xl md:text-4xl font-bold mb-2 text-slate-900 dark:text-white">{{ poll.title }}</h1>
          <p class="text-base text-slate-500 dark:text-slate-400">{{ poll.description || 'Penghitungan suara sedang berlangsung secara real-time.' }}</p>
        </div>

        <NuxtLink to="/dashboard" class="flex items-center gap-2 text-slate-500 dark:text-slate-400 hover:text-poster-base dark:hover:text-poster-light transition-colors bg-white dark:bg-slate-900 hover:bg-red-50 dark:hover:bg-slate-800 border border-slate-200 dark:border-slate-800 px-4 py-2 rounded-xl text-sm font-medium shadow-sm">
          <ArrowLeft :size="18" />
          <span>Kembali</span>
        </NuxtLink>
      </header>

      <div class="space-y-8">
        <div v-for="position in poll.positions" :key="position.ID" class="bg-white dark:bg-slate-900 p-6 md:p-8 rounded-2xl border border-slate-200 dark:border-slate-800 shadow-sm transition-colors duration-300">

          <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8 border-b border-slate-100 dark:border-slate-800 pb-4">
            <div class="flex items-center gap-3">
              <div class="p-2 bg-red-50 dark:bg-poster-base/10 text-poster-base dark:text-poster-light rounded-lg">
                <Trophy :size="24" />
              </div>
              <h2 class="text-2xl font-bold text-slate-800 dark:text-slate-100">{{ position.title }}</h2>
            </div>
            <div class="inline-flex items-center gap-2 text-sm font-medium text-slate-600 dark:text-slate-300 bg-slate-50 dark:bg-slate-800 px-4 py-1.5 rounded-lg border border-slate-200 dark:border-slate-700">
              <Activity :size="16" class="text-slate-400 dark:text-slate-500" />
              Total Suara: <span class="text-slate-900 dark:text-white font-bold">{{ getPositionTotal(position.options) }}</span>
            </div>
          </div>

          <div class="space-y-5">
            <div v-for="(option, index) in position.options" :key="option.ID" class="bg-slate-50/50 dark:bg-slate-800/50 border border-slate-100 dark:border-slate-700 p-5 rounded-xl hover:border-poster-base/30 dark:hover:border-poster-light/30 transition-colors group">

              <div class="flex flex-col md:flex-row md:items-center justify-between mb-4 gap-4">
                <div class="flex items-center gap-4">
                  <div class="w-10 h-10 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-600 text-slate-700 dark:text-slate-200 font-bold text-lg flex items-center justify-center rounded-lg shadow-sm flex-shrink-0">
                    {{ index + 1 }}
                  </div>

                  <div v-if="option.photo_url" class="w-12 h-12 md:w-14 md:h-14 rounded-lg overflow-hidden bg-slate-100 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 flex-shrink-0 shadow-sm">
                    <img :src="getImageUrl(option.photo_url)" class="w-full h-full object-cover" />
                  </div>

                  <h3 class="text-xl font-semibold text-slate-800 dark:text-slate-100">{{ option.value }}</h3>
                </div>

                <div class="text-left md:text-right flex items-baseline gap-2 md:block">
                  <span class="text-2xl font-bold text-slate-900 dark:text-white">{{ option.vote_count }}</span>
                  <span class="text-sm font-medium text-slate-500 dark:text-slate-400">Suara</span>
                </div>
              </div>

              <div class="flex items-center gap-4">
                <div class="flex-grow h-3 bg-slate-200 dark:bg-slate-700 rounded-full overflow-hidden">
                  <div class="h-full rounded-full transition-all duration-1000 ease-out"
                       :class="[
                         index % 4 === 0 ? 'bg-poster-base' : '',
                         index % 4 === 1 ? 'bg-poster-dark' : '',
                         index % 4 === 2 ? 'bg-amber-500 dark:bg-amber-400' : '',
                         index % 4 === 3 ? 'bg-poster-light' : ''
                       ]"
                       :style="{ width: `${calculatePercentage(option.vote_count, getPositionTotal(position.options))}%` }">
                  </div>
                </div>
                <div class="w-12 text-right text-sm font-bold text-slate-700 dark:text-slate-300">
                  {{ calculatePercentage(option.vote_count, getPositionTotal(position.options)) }}%
                </div>
              </div>

            </div>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>