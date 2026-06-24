<script setup>
import { onMounted } from 'vue'
import { Search, CheckCircle2, XCircle, Loader2 } from 'lucide-vue-next'
import { usePollStore } from '~/stores/poll'
import { storeToRefs } from 'pinia'

definePageMeta({ middleware: ['auth'] })

const pollStore = usePollStore()
const { polls, isLoading } = storeToRefs(pollStore)

onMounted(() => {
  pollStore.fetchPolls()
})

const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return new Intl.DateTimeFormat('id-ID', { dateStyle: 'medium', timeStyle: 'short' }).format(date)
}
</script>

<template>
  <div>
    <div class="mb-8">
      <h2 class="text-3xl font-black text-slate-800 dark:text-white tracking-tight">Data & Arsip Polling</h2>
      <p class="text-slate-500 dark:text-slate-400 font-medium mt-1">Riwayat seluruh kegiatan pemilihan (Aktif maupun Arsip).</p>
    </div>
    <div class="bg-white dark:bg-slate-900 rounded-[2rem] border border-slate-200 dark:border-slate-800 shadow-sm overflow-hidden transition-colors duration-500">
      <div class="p-6 border-b border-slate-100 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-800/50 transition-colors duration-500">
        <div class="relative w-full md:w-96">
          <Search :size="18" class="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400" />
          <input type="text" placeholder="Cari judul pemilihan..." class="w-full pl-11 pr-4 py-3 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-xl focus:ring-2 focus:ring-indigo-500 outline-none transition-all font-medium text-slate-800 dark:text-white placeholder:text-slate-400">
        </div>
      </div>

      <div v-if="isLoading" class="p-20 text-center flex flex-col items-center">
        <Loader2 class="animate-spin text-indigo-500 mb-4" :size="40" />
        <p class="text-slate-500 dark:text-slate-400 font-bold">Memuat arsip data...</p>
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
          <tr class="bg-slate-50 dark:bg-slate-800/50 text-slate-500 dark:text-slate-400 text-xs uppercase tracking-widest border-b border-slate-100 dark:border-slate-800 transition-colors duration-500">
            <th class="p-6 font-black">Judul Pemilihan</th>
            <th class="p-6 font-black text-center">Status (Arsip)</th>
            <th class="p-6 font-black text-center">Kandidat</th>
            <th class="p-6 font-black text-center">Total Token</th>
            <th class="p-6 font-black">Dibuat Pada</th>
          </tr>
          </thead>
          <tbody>
          <tr v-for="poll in polls" :key="poll.ID" class="border-b border-slate-50 dark:border-slate-800 hover:bg-slate-50/80 dark:hover:bg-slate-800/80 transition-colors">
            <td class="p-6">
              <p class="font-black text-slate-800 dark:text-white text-lg mb-1">{{ poll.title }}</p>
              <p class="text-sm text-slate-500 dark:text-slate-400">{{ poll.description || 'Tidak ada deskripsi' }}</p>
            </td>
            <td class="p-6 text-center">
                <span v-if="poll.is_active" class="inline-flex items-center gap-1.5 px-3 py-1 bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border border-emerald-200 dark:border-emerald-500/20 rounded-full text-xs font-black uppercase tracking-wider transition-colors duration-500">
                  <CheckCircle2 :size="14" /> Aktif
                </span>
              <span v-else class="inline-flex items-center gap-1.5 px-3 py-1 bg-slate-100 dark:bg-slate-800 text-slate-500 dark:text-slate-400 border border-slate-200 dark:border-slate-700 rounded-full text-xs font-black uppercase tracking-wider transition-colors duration-500">
                  <XCircle :size="14" /> Diarsipkan
                </span>
            </td>
            <td class="p-6 text-center font-bold text-slate-700 dark:text-slate-300">
              {{ poll.positions?.reduce((total, pos) => total + (pos.options?.length || 0), 0) || 0 }} Opsi
            </td>
            <td class="p-6 text-center font-mono font-bold text-indigo-600 dark:text-indigo-400">
              {{ poll.tokens?.length || 0 }}
            </td>
            <td class="p-6 font-medium text-slate-500 dark:text-slate-400">{{ formatDate(poll.CreatedAt) }}</td>
          </tr>

          <tr v-if="polls?.length === 0">
            <td colspan="5" class="p-10 text-center font-bold text-slate-400 dark:text-slate-500">Belum ada data atau arsip pemilihan.</td>
          </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

