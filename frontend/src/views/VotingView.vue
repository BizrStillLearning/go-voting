<script setup>
import { ref, onMounted, computed } from 'vue'
import {
  CheckCircle2, BarChart3, Info,
  RotateCcw, Send, Trophy
} from 'lucide-vue-next'

const paslonList = ref([])
const selectedId = ref(null)
const hasVoted = ref(false)

// Warna untuk Progress Bar
const warnaProgress = {
  blue: "bg-blue-500",
  emerald: "bg-emerald-500",
  rose: "bg-rose-500",
  amber: "bg-amber-500"
}

const loadData = () => {
  const saved = localStorage.getItem('go_voting_paslon')
  if (saved) paslonList.value = JSON.parse(saved)

  const voted = localStorage.getItem('user_has_voted')
  if (voted) hasVoted.value = true
}

const totalSuara = computed(() => {
  return paslonList.value.reduce((acc, curr) => acc + curr.suara, 0)
})

const handleVote = () => {
  if (!selectedId.value) return alert('Pilih salah satu paslon!')

  const index = paslonList.value.findIndex(p => p.id === selectedId.value)
  if (index !== -1) {
    paslonList.value[index].suara += 1
    hasVoted.value = true

    // Simpan ke local storage
    localStorage.setItem('go_voting_paslon', JSON.stringify(paslonList.value))
    localStorage.setItem('user_has_voted', 'true')
  }
}

const resetSession = () => {
  if (confirm('Reset sesi voting ini? (Hanya untuk demo)')) {
    hasVoted.value = false
    localStorage.removeItem('user_has_voted')
    selectedId.value = null
  }
}

onMounted(loadData)
</script>

<template>
  <div class="max-w-5xl mx-auto animate-in fade-in slide-in-from-bottom-4 duration-500">
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-8">

      <div class="lg:col-span-7 space-y-6">
        <div class="bg-white p-8 rounded-[2.5rem] border border-slate-200 shadow-sm">
          <div class="flex justify-between items-center mb-8">
            <h2 class="text-2xl font-black text-slate-800 flex items-center gap-3">
              <CheckCircle2 :size="28" class="text-emerald-500" />
              {{ hasVoted ? 'Hasil Pemungutan Suara' : 'Berikan Suara Anda' }}
            </h2>
            <button @click="resetSession" class="p-2 text-slate-400 hover:text-slate-600 transition-colors" title="Reset Sesi">
              <RotateCcw :size="18" />
            </button>
          </div>

          <div v-if="!hasVoted" class="grid grid-cols-1 gap-4">
            <label
                v-for="p in paslonList" :key="p.id"
                :class="selectedId === p.id ? 'border-indigo-500 ring-2 ring-indigo-500/10 bg-indigo-50/30' : 'border-slate-100 bg-slate-50 hover:bg-white hover:border-slate-200'"
                class="relative flex items-center p-5 rounded-3xl border-2 cursor-pointer transition-all group"
            >
              <input type="radio" :value="p.id" v-model="selectedId" class="hidden" />

              <div class="flex-1 flex items-center gap-5">
                <div class="w-14 h-14 rounded-2xl bg-white flex items-center justify-center text-xl font-black text-slate-800 shadow-sm border border-slate-100">
                  {{ p.nomor }}
                </div>
                <div>
                  <h3 class="font-bold text-slate-900">{{ p.ketua }}</h3>
                  <p class="text-sm text-slate-500">Wakil: {{ p.wakil }}</p>
                </div>
              </div>

              <div class="w-6 h-6 rounded-full border-2 border-slate-300 flex items-center justify-center" :class="{'border-indigo-500 bg-indigo-500': selectedId === p.id}">
                <div v-if="selectedId === p.id" class="w-2 h-2 bg-white rounded-full"></div>
              </div>
            </label>

            <button
                @click="handleVote"
                :disabled="!selectedId"
                class="w-full mt-6 bg-slate-900 hover:bg-indigo-600 disabled:bg-slate-200 disabled:text-slate-400 text-white font-black py-4 rounded-2xl shadow-xl transition-all flex items-center justify-center gap-2"
            >
              <Send :size="18" /> Kirim Suara Sekarang
            </button>
          </div>

          <div v-else class="text-center py-10 space-y-4">
            <div class="w-20 h-20 bg-emerald-100 text-emerald-600 rounded-full flex items-center justify-center mx-auto mb-4">
              <Trophy :size="40" />
            </div>
            <h3 class="text-2xl font-black text-slate-900">Terima Kasih!</h3>
            <p class="text-slate-500 max-w-xs mx-auto text-sm">Suara Anda telah berhasil direkam dalam sistem Go-Voting secara anonim.</p>
          </div>
        </div>
      </div>

      <div class="lg:col-span-5">
        <div class="bg-white p-8 rounded-[2.5rem] border border-slate-200 shadow-sm sticky top-28">
          <div class="flex items-center gap-3 mb-8">
            <div class="p-2 bg-indigo-50 rounded-lg text-indigo-600">
              <BarChart3 :size="20" />
            </div>
            <h2 class="text-xl font-bold text-slate-800">Live Statistics</h2>
          </div>

          <div class="space-y-6">
            <div v-for="p in paslonList" :key="p.id" class="space-y-2">
              <div class="flex justify-between items-end text-sm">
                <div>
                  <span class="font-black text-slate-900">#{{ p.nomor }}</span>
                  <span class="ml-2 font-medium text-slate-600 truncate inline-block max-w-[150px] align-bottom">{{ p.ketua }}</span>
                </div>
                <div class="text-right">
                  <span class="font-black text-slate-900">{{ p.suara }}</span>
                  <span class="text-slate-400 ml-1">Votes</span>
                </div>
              </div>

              <div class="h-3 w-full bg-slate-100 rounded-full overflow-hidden">
                <div
                    class="h-full rounded-full transition-all duration-1000 ease-out"
                    :class="warnaProgress[p.warna]"
                    :style="{ width: totalSuara === 0 ? '0%' : (p.suara / totalSuara * 100) + '%' }"
                ></div>
              </div>

              <div class="flex justify-end">
                <span class="text-[10px] font-bold text-slate-400 tracking-widest uppercase">
                  {{ totalSuara === 0 ? 0 : (p.suara / totalSuara * 100).toFixed(1) }}%
                </span>
              </div>
            </div>

            <div class="pt-6 border-t border-slate-50 mt-4">
              <div class="flex justify-between items-center text-slate-400">
                <span class="text-xs font-bold uppercase tracking-widest flex items-center gap-1">
                   <Info :size="12" /> Total Suara Masuk
                </span>
                <span class="text-lg font-black text-slate-900">{{ totalSuara }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>