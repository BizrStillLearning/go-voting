<script setup>
import { ref, onMounted, computed } from 'vue'
import { PlusCircle, Users, Key, PieChart, Trash2, Loader2, BarChart3, Activity, ShieldCheck, Lock, Download, Printer, Link, ExternalLink } from 'lucide-vue-next'
import Swal from 'sweetalert2'
import { usePollStore } from '~/stores/poll'
import { useAuthStore } from '~/stores/auth'
import { storeToRefs } from 'pinia'

definePageMeta({ middleware: ['auth'] })

const pollStore = usePollStore()
const authStore = useAuthStore()
const { polls, isLoading } = storeToRefs(pollStore)

const showModal = ref(false)

onMounted(() => {
  pollStore.fetchPolls()
})

const totalPolls = computed(() => polls.value.length)
const activePolls = computed(() => polls.value.filter(p => p.is_active).length)
const totalVotesCast = computed(() => {
  return polls.value.reduce((total, poll) => {
    return total + (poll.tokens?.filter(t => t.is_used).length || 0)
  }, 0)
})

const handleCreatePoll = async (formData) => {
  try {
    await $fetch('http://localhost:8080/api/admin/poll/create', {
      method: 'POST',
      headers: { Authorization: `Bearer ${authStore.token}` },
      body: formData
    })

    showModal.value = false
    pollStore.fetchPolls()

    Swal.fire({
      icon: 'success',
      title: 'Berhasil!',
      text: 'Pemilihan baru siap digunakan.',
      timer: 2000,
      showConfirmButton: false,
      background: document.documentElement.classList.contains('dark') ? '#1e293b' : '#ffffff',
      color: document.documentElement.classList.contains('dark') ? '#f8fafc' : '#0f172a',
      customClass: { popup: 'rounded-2xl' }
    })

  } catch (error) {
    Swal.fire({
      icon: 'error',
      title: 'Gagal Membuat',
      text: error.data?.error || "Terjadi kesalahan saat menyimpan data.",
      background: document.documentElement.classList.contains('dark') ? '#1e293b' : '#ffffff',
      color: document.documentElement.classList.contains('dark') ? '#f8fafc' : '#0f172a',
      customClass: { popup: 'rounded-2xl' }
    })
  }
}

const handleClosePoll = async (pollID) => {
  const isDarkMode = document.documentElement.classList.contains('dark')
  const result = await Swal.fire({
    title: 'Tutup Pemilihan?',
    text: "Pemilih tidak akan bisa memasukkan token lagi setelah ini ditutup.",
    icon: 'warning',
    showCancelButton: true,
    confirmButtonColor: '#f59e0b',
    cancelButtonColor: isDarkMode ? '#334155' : '#94a3b8',
    confirmButtonText: 'Ya, Tutup!',
    cancelButtonText: `<span style="color: ${isDarkMode ? '#cbd5e1' : '#f8fafc'};">Batal</span>`,
    background: isDarkMode ? '#1e293b' : '#ffffff',
    color: isDarkMode ? '#f8fafc' : '#0f172a',
    customClass: { popup: 'rounded-2xl' }
  })

  if (!result.isConfirmed) return

  try {
    await $fetch(`http://localhost:8080/api/admin/poll/${pollID}/close`, {
      method: 'PUT',
      headers: { Authorization: `Bearer ${authStore.token}` }
    })
    pollStore.fetchPolls()
    Swal.fire({
      icon: 'success',
      title: 'Ditutup!',
      text: 'Pemilihan telah resmi dihentikan.',
      timer: 1500,
      showConfirmButton: false,
      background: isDarkMode ? '#1e293b' : '#ffffff',
      color: isDarkMode ? '#f8fafc' : '#0f172a',
      customClass: { popup: 'rounded-2xl' }
    })
  } catch (error) {
    Swal.fire({ icon: 'error', title: 'Gagal', text: 'Terjadi kesalahan sistem.', background: isDarkMode ? '#1e293b' : '#ffffff', color: isDarkMode ? '#f8fafc' : '#0f172a', customClass: { popup: 'rounded-2xl' } })
  }
}

const handleExportCSV = async (pollID) => {
  try {
    const blob = await $fetch(`http://localhost:8080/api/admin/poll/${pollID}/export`, {
      headers: { Authorization: `Bearer ${authStore.token}` },
      responseType: 'blob'
    })

    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `laporan-evoting-${pollID}.csv`)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)

  } catch (error) {
    Swal.fire({ icon: 'error', title: 'Gagal Unduh', text: 'Tidak dapat mengekspor laporan dari server.', customClass: { popup: 'rounded-2xl' } })
  }
}

const handleCopyLink = (slug) => {
  const url = `${window.location.origin}/v/${slug}`

  navigator.clipboard.writeText(url).then(() => {
    const isDarkMode = document.documentElement.classList.contains('dark')
    Swal.fire({
      icon: 'success',
      title: 'Tautan Tersalin!',
      text: 'Tempel (Paste) tautan ini di laptop bilik suara.',
      timer: 2000,
      showConfirmButton: false,
      background: isDarkMode ? '#1e293b' : '#ffffff',
      color: isDarkMode ? '#f8fafc' : '#0f172a',
      customClass: { popup: 'rounded-2xl' }
    })
  }).catch(() => {
    Swal.fire({ icon: 'error', title: 'Gagal', text: 'Tidak dapat menyalin tautan.' })
  })
}

const handleDeletePoll = async (pollID) => {
  const isDarkMode = document.documentElement.classList.contains('dark')
  const result = await Swal.fire({
    title: 'Hapus Permanen?',
    text: "Semua data akan lenyap. Anda tidak dapat mengembalikannya!",
    icon: 'warning',
    showCancelButton: true,
    confirmButtonColor: '#A51013',
    cancelButtonColor: isDarkMode ? '#334155' : '#94a3b8',
    confirmButtonText: 'Ya, Hapus!',
    cancelButtonText: `<span style="color: ${isDarkMode ? '#cbd5e1' : '#f8fafc'};">Batal</span>`,
    background: isDarkMode ? '#1e293b' : '#ffffff',
    color: isDarkMode ? '#f8fafc' : '#0f172a',
    customClass: { popup: 'rounded-2xl' }
  })

  if (!result.isConfirmed) return

  try {
    await $fetch(`http://localhost:8080/api/admin/poll/${pollID}/delete`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${authStore.token}` }
    })

    pollStore.fetchPolls()
    Swal.fire({
      icon: 'success',
      title: 'Terhapus!',
      text: 'Data pemilihan berhasil dibersihkan.',
      timer: 1500,
      showConfirmButton: false,
      background: isDarkMode ? '#1e293b' : '#ffffff',
      color: isDarkMode ? '#f8fafc' : '#0f172a',
      customClass: { popup: 'rounded-2xl' }
    })
  } catch (error) {
    Swal.fire({
      icon: 'error',
      title: 'Gagal Menghapus',
      text: error.data?.error || "Terjadi kesalahan koneksi database.",
      background: isDarkMode ? '#1e293b' : '#ffffff',
      color: isDarkMode ? '#f8fafc' : '#0f172a',
      customClass: { popup: 'rounded-2xl' }
    })
  }
}
</script>

<template>
  <div class="font-sans text-slate-800 dark:text-slate-200 transition-colors duration-300">

    <div class="mb-10">
      <div class="flex flex-col md:flex-row md:items-end justify-between gap-6 mb-8">
        <div>
          <h1 class="text-3xl font-bold text-slate-900 dark:text-white tracking-tight">Tinjauan Sistem</h1>
          <p class="text-slate-500 dark:text-slate-400 font-medium mt-1 text-base">Pusat kendali seluruh aktivitas pemilihan digital.</p>
        </div>
        <button @click="showModal = true" class="flex items-center justify-center gap-2 bg-poster-base dark:bg-poster-light hover:bg-poster-dark dark:hover:bg-poster-base text-white px-6 py-3 rounded-xl font-semibold shadow-sm active:scale-95 transition-all duration-200 cursor-pointer">
          <PlusCircle :size="18" /> Buat Pemilihan
        </button>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-5">
        <div class="bg-white dark:bg-slate-900 p-6 rounded-2xl border border-slate-200 dark:border-slate-800 flex items-center gap-5 shadow-sm hover:shadow-md transition-shadow">
          <div class="p-3.5 bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300 rounded-xl"><BarChart3 :size="24" /></div>
          <div>
            <p class="text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wider mb-1">Total Acara</p>
            <p class="text-2xl font-bold text-slate-900 dark:text-white">{{ totalPolls }}</p>
          </div>
        </div>

        <div class="bg-white dark:bg-slate-900 p-6 rounded-2xl border border-slate-200 dark:border-slate-800 flex items-center gap-5 shadow-sm hover:shadow-md transition-shadow">
          <div class="p-3.5 bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 rounded-xl"><Activity :size="24" /></div>
          <div>
            <p class="text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wider mb-1">Berjalan</p>
            <p class="text-2xl font-bold text-slate-900 dark:text-white">{{ activePolls }}</p>
          </div>
        </div>

        <div class="bg-white dark:bg-slate-900 p-6 rounded-2xl border border-slate-200 dark:border-slate-800 flex items-center gap-5 shadow-sm hover:shadow-md transition-shadow">
          <div class="p-3.5 bg-red-50 dark:bg-poster-base/10 text-poster-base dark:text-poster-light rounded-xl"><ShieldCheck :size="24" /></div>
          <div>
            <p class="text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wider mb-1">Suara Masuk</p>
            <p class="text-2xl font-bold text-slate-900 dark:text-white">{{ totalVotesCast }}</p>
          </div>
        </div>
      </div>
    </div>

    <div v-if="isLoading" class="flex flex-col items-center justify-center py-20 text-slate-500 dark:text-slate-400">
      <Loader2 class="animate-spin mb-4 text-poster-base dark:text-poster-light" :size="36" />
      <p class="font-medium">Memuat data...</p>
    </div>

    <div v-else-if="polls?.length === 0" class="max-w-2xl mx-auto bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 border-dashed rounded-2xl p-12 flex flex-col items-center justify-center text-center">
      <div class="bg-slate-50 dark:bg-slate-800 p-4 rounded-full mb-5">
        <PlusCircle :size="36" class="text-slate-400 dark:text-slate-500" />
      </div>
      <h3 class="text-xl font-bold text-slate-900 dark:text-white mb-2">Belum Ada Pemilihan</h3>
      <p class="text-slate-500 dark:text-slate-400 mb-6">Anda belum memulai kegiatan pemilihan apa pun.</p>
      <button @click="showModal = true" class="text-poster-base dark:text-poster-light font-bold hover:underline">
        Buat pemilihan pertama Anda
      </button>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
      <div v-for="poll in polls" :key="poll.ID" class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 p-6 rounded-2xl shadow-sm hover:shadow-md transition-all duration-200 flex flex-col h-full relative group">

        <div class="flex justify-between items-start mb-4">
          <span v-if="poll.is_active" class="px-2.5 py-1 bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border border-emerald-200 dark:border-emerald-500/20 rounded-lg text-[11px] font-bold uppercase tracking-wide flex items-center gap-1.5">
            <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></span> LIVE
          </span>
          <span v-else class="px-2.5 py-1 bg-slate-100 dark:bg-slate-800 text-slate-500 dark:text-slate-400 rounded-lg text-[11px] font-bold uppercase tracking-wide border border-slate-200 dark:border-slate-700">
            Selesai
          </span>

          <span class="text-xs font-medium text-slate-500 dark:text-slate-400 bg-slate-50 dark:bg-slate-800 px-2.5 py-1 rounded-lg border border-slate-100 dark:border-slate-800">
            {{ poll.positions?.reduce((total, pos) => total + (pos.options?.length || 0), 0) || 0 }} Kandidat
          </span>
        </div>

        <div class="flex-1 mb-6">
          <h3 class="text-xl font-bold mb-1.5 text-slate-900 dark:text-white leading-snug">{{ poll.title }}</h3>
          <p class="text-slate-500 dark:text-slate-400 text-sm mb-5 line-clamp-2">{{ poll.description || 'Tanpa deskripsi khusus.' }}</p>

          <div class="mb-2">
            <div class="flex justify-between items-end mb-2">
              <span class="text-[11px] font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wide">Partisipasi Token</span>
              <span class="text-sm font-semibold text-slate-800 dark:text-slate-200">{{ poll.tokens?.filter(t => t.is_used).length || 0 }} / {{ poll.tokens?.length || 0 }}</span>
            </div>
            <div class="w-full h-1.5 bg-slate-100 dark:bg-slate-800 rounded-full overflow-hidden">
              <div class="h-full bg-poster-base dark:bg-poster-light rounded-full transition-all duration-500"
                   :style="{ width: `${ poll.tokens?.length ? ((poll.tokens.filter(t => t.is_used).length / poll.tokens.length) * 100) : 0 }%` }">
              </div>
            </div>
          </div>
        </div>

        <div class="grid grid-cols-2 gap-2 mt-auto pt-5 border-t border-slate-100 dark:border-slate-800">
          <div v-if="poll.is_active" class="col-span-2 flex gap-2 mb-2">
            <NuxtLink :to="`/v/${poll.slug}`" target="_blank" class="flex-1 flex items-center justify-center gap-1.5 bg-poster-base hover:bg-poster-dark text-white py-2 rounded-xl text-xs font-bold transition-colors shadow-sm">
              <ExternalLink :size="14" /> Buka Bilik Suara
            </NuxtLink>
            <button @click="handleCopyLink(poll.slug)" class="flex-1 flex items-center justify-center gap-1.5 bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300 hover:bg-slate-200 dark:hover:bg-slate-700 py-2 rounded-xl text-xs font-bold transition-colors">
              <Link :size="14" /> Salin Tautan
            </button>
          </div>

          <NuxtLink :to="`/polls/result/${poll.slug}`" target="_blank" class="flex items-center justify-center gap-1.5 bg-slate-900 dark:bg-slate-800 text-white py-2 rounded-xl text-xs font-medium hover:bg-slate-800 dark:hover:bg-slate-700 transition-colors">
            <PieChart :size="14" /> Monitor
          </NuxtLink>

          <button @click="handleExportCSV(poll.ID)" class="flex items-center justify-center gap-1.5 bg-white dark:bg-slate-900 text-slate-700 dark:text-slate-300 border border-slate-200 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-800 py-2 rounded-xl text-xs font-medium transition-colors">
            <Download :size="14" /> CSV
          </button>

          <NuxtLink :to="`/print/${poll.ID}`" target="_blank" class="flex items-center justify-center gap-1.5 bg-red-50 dark:bg-poster-base/10 text-poster-base dark:text-poster-light hover:bg-red-100 dark:hover:bg-poster-base/20 py-2 rounded-xl text-xs font-medium transition-colors col-span-2">
            <Printer :size="14" /> Cetak Token
          </NuxtLink>

          <button v-if="poll.is_active" @click="handleClosePoll(poll.ID)" class="flex items-center justify-center gap-1.5 bg-amber-50 dark:bg-amber-500/10 text-amber-700 dark:text-amber-400 hover:bg-amber-100 dark:hover:bg-amber-500/20 py-2 rounded-xl text-xs font-medium transition-colors col-span-2">
            <Lock :size="14" /> Hentikan Pemilihan
          </button>

          <button @click="handleDeletePoll(poll.ID)" class="flex items-center justify-center gap-1.5 text-slate-400 hover:text-poster-dark dark:hover:text-red-400 hover:bg-red-50 dark:hover:bg-red-500/10 py-2 rounded-xl text-xs font-medium transition-colors col-span-2">
            <Trash2 :size="14" /> Hapus Permanen
          </button>
        </div>

      </div>
    </div>

    <CreatePollModal
        v-if="showModal"
        @close="showModal = false"
        @submit="handleCreatePoll"
    />

  </div>
</template>