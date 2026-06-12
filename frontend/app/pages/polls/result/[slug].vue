<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'nuxt/app'
import { ArrowLeft, Users, Loader2, AlertTriangle, Trophy, UserCircle2 } from 'lucide-vue-next'

const route = useRoute()
const slug = route.params.slug

const pollResult = ref(null)
const isLoading = ref(true)
const fetchError = ref('')
let refreshInterval = null

const getPositionTotal = (options) => {
  if (!options) return 0
  return options.reduce((sum, opt) => sum + opt.vote_count, 0)
}

const getPercentage = (count, total) => {
  if (total === 0) return 0
  return ((count / total) * 100).toFixed(1)
}

const getHighestVote = (options) => {
  if (!options || options.length === 0) return 0
  return Math.max(...options.map(opt => opt.vote_count))
}

const fetchLiveResults = async () => {
  try {
    const res = await $fetch(`http://localhost:8080/api/poll/${slug}/results`)
    pollResult.value = res
  } catch (error) {
    fetchError.value = "Gagal mengambil data Live. Pastikan server berjalan."
    console.error(error)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchLiveResults()
  refreshInterval = setInterval(fetchLiveResults, 3000)
})

onUnmounted(() => {
  clearInterval(refreshInterval)
})
</script>

<template>
  <div class="min-h-screen bg-slate-900 text-white font-sans flex flex-col selection:bg-indigo-500">

    <header class="p-6 md:p-8 flex items-center justify-between border-b border-slate-800 bg-slate-900/50 backdrop-blur-md sticky top-0 z-10">
      <NuxtLink to="/dashboard" class="flex items-center gap-2 text-slate-400 hover:text-white transition-colors font-bold bg-slate-800 hover:bg-slate-700 px-4 py-2 rounded-xl">
        <ArrowLeft :size="18" /> Kembali ke Dashboard
      </NuxtLink>

      <div class="flex items-center gap-3">
        <div class="px-4 py-2 bg-indigo-500/10 border border-indigo-500/30 text-indigo-400 rounded-xl font-mono text-sm font-bold tracking-widest uppercase">
          ID: {{ slug }}
        </div>
        <div class="flex items-center gap-2 px-4 py-2 bg-red-500/10 border border-red-500/30 text-red-400 rounded-xl font-black text-sm uppercase tracking-widest shadow-[0_0_15px_rgba(239,68,68,0.3)]">
          <span class="w-2.5 h-2.5 bg-red-500 rounded-full animate-pulse"></span>
          LIVE COUNT
        </div>
      </div>
    </header>

    <main class="flex-1 max-w-6xl w-full mx-auto p-6 md:p-12 flex flex-col">

      <div v-if="isLoading" class="flex flex-col items-center justify-center py-20 flex-1">
        <Loader2 class="animate-spin mb-4 text-indigo-500" :size="64" />
        <p class="text-slate-400 font-bold text-xl">Menghubungkan ke satelit suara...</p>
      </div>

      <div v-else-if="fetchError" class="text-center py-20 flex-1">
        <AlertTriangle class="mx-auto text-amber-500 mb-4" :size="80" />
        <h1 class="text-3xl font-black text-white mb-2">Koneksi Terputus</h1>
        <p class="text-slate-400 text-lg">{{ fetchError }}</p>
      </div>

      <div v-else-if="pollResult?.poll">

        <div class="text-center mb-16">
          <h1 class="text-4xl md:text-6xl font-black mb-6 leading-tight text-transparent bg-clip-text bg-gradient-to-r from-white to-slate-400">
            {{ pollResult.poll.title }}
          </h1>
          <div class="inline-flex items-center justify-center gap-4 bg-slate-800/80 border border-slate-700 px-8 py-5 rounded-3xl shadow-2xl backdrop-blur-sm">
            <Users class="text-indigo-400" :size="36" />
            <div class="text-left">
              <p class="text-slate-400 text-sm font-black uppercase tracking-widest">Total Partisipasi Pemilih</p>
              <p class="text-5xl font-black font-mono text-white tracking-tight">{{ pollResult.total_votes }} <span class="text-lg text-slate-500 font-sans">Suara Masuk</span></p>
            </div>
          </div>
        </div>

        <div class="space-y-12">
          <div v-for="pos in pollResult.poll.positions" :key="pos.ID" class="bg-slate-800/40 border border-slate-700/50 p-8 rounded-[3rem] shadow-xl relative overflow-hidden">

            <div class="absolute top-0 right-0 w-64 h-64 bg-indigo-500/5 rounded-full blur-3xl -translate-y-1/2 translate-x-1/2"></div>

            <div class="flex items-center justify-between mb-8 border-b border-slate-700/50 pb-4">
              <h2 class="text-3xl font-black text-indigo-300 uppercase tracking-widest">{{ pos.title }}</h2>
              <span class="text-slate-400 font-bold bg-slate-900/50 px-4 py-2 rounded-xl text-sm border border-slate-700">
                Total: {{ getPositionTotal(pos.options) }} Suara
              </span>
            </div>

            <div class="space-y-6">
              <div v-for="(opt, index) in pos.options" :key="opt.ID" class="relative group">

                <div class="flex items-center gap-4 mb-3 relative z-10">
                  <div class="w-12 h-12 rounded-full overflow-hidden border-2 border-slate-600 bg-slate-800 flex-shrink-0">
                    <img v-if="opt.photo_url" :src="opt.photo_url" class="w-full h-full object-cover" />
                    <UserCircle2 v-else :size="44" class="text-slate-500 m-auto mt-0.5" />
                  </div>

                  <div class="flex-1 flex justify-between items-end">
                    <div class="flex items-center gap-3">
                      <h3 class="text-2xl font-bold text-slate-200">{{ opt.value }}</h3>
                      <Trophy v-if="opt.vote_count > 0 && opt.vote_count === getHighestVote(pos.options)" :size="20" class="text-yellow-400 animate-bounce" />
                    </div>
                    <div class="text-right">
                      <span class="text-3xl font-black font-mono text-white">{{ opt.vote_count }}</span>
                      <span class="text-slate-400 font-bold ml-2">suara</span>
                    </div>
                  </div>
                </div>

                <div class="w-full h-12 bg-slate-900/80 rounded-full overflow-hidden shadow-inner relative border border-slate-700">
                  <div class="h-full rounded-full flex items-center justify-end px-4 font-black transition-all duration-1000 ease-out"
                       :class="[
                         index % 4 === 0 ? 'bg-gradient-to-r from-blue-600 to-indigo-500' : '',
                         index % 4 === 1 ? 'bg-gradient-to-r from-teal-500 to-emerald-400' : '',
                         index % 4 === 2 ? 'bg-gradient-to-r from-amber-500 to-orange-400' : '',
                         index % 4 === 3 ? 'bg-gradient-to-r from-purple-600 to-pink-500' : ''
                       ]"
                       :style="{ width: `${Math.max(getPercentage(opt.vote_count, getPositionTotal(pos.options)), 5)}%` }">
                    <span class="text-white drop-shadow-md text-lg">{{ getPercentage(opt.vote_count, getPositionTotal(pos.options)) }}%</span>
                  </div>
                </div>

              </div>
            </div>

          </div>
        </div>
      </div>

    </main>
  </div>
</template>