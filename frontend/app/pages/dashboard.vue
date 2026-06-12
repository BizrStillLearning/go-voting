<script setup>
import { ref, onMounted, computed } from 'vue'
import { PlusCircle, Users, Key, PieChart, Trash2, Loader2, BarChart3, Activity, ShieldCheck } from 'lucide-vue-next'
import Swal from 'sweetalert2' // Import SweetAlert2
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
      customClass: { popup: 'rounded-3xl' }
    })

  } catch (error) {
    Swal.fire({
      icon: 'error',
      title: 'Gagal Membuat',
      text: error.data?.error || "Terjadi kesalahan saat menyimpan data.",
      customClass: { popup: 'rounded-3xl' }
    })
  }
}

const handleDeletePoll = async (pollID) => {
  const result = await Swal.fire({
    title: 'Hapus Permanen?',
    text: "Semua data suara, jabatan, dan token akan lenyap. Anda tidak dapat mengembalikannya!",
    icon: 'warning',
    showCancelButton: true,
    confirmButtonColor: '#ef4444',
    cancelButtonColor: '#94a3b8',
    confirmButtonText: 'Ya, Bumihanguskan!',
    cancelButtonText: 'Batal',
    customClass: { popup: 'rounded-3xl' }
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
      text: 'Data pemilihan berhasil dibersihkan dari server.',
      timer: 1500,
      showConfirmButton: false,
      customClass: { popup: 'rounded-3xl' }
    })
  } catch (error) {
    Swal.fire({
      icon: 'error',
      title: 'Gagal Menghapus',
      text: error.data?.error || "Terjadi kesalahan koneksi database.",
      customClass: { popup: 'rounded-3xl' }
    })
  }
}
</script>

<template>
  <div class="flex h-screen bg-[#f8fafc] font-sans overflow-hidden selection:bg-indigo-200">
    <Sidebar />

    <div class="flex-1 flex flex-col h-screen overflow-hidden relative">
      <Header />

      <main class="flex-1 overflow-x-hidden overflow-y-auto p-6 md:p-8 lg:px-12">

        <div class="mb-10">
          <div class="flex flex-col md:flex-row md:items-end justify-between gap-6 mb-8">
            <div>
              <h1 class="text-4xl font-black text-slate-800 tracking-tight">Tinjauan Sistem</h1>
              <p class="text-slate-500 font-medium mt-2 text-lg">Pusat kendali seluruh aktivitas pemilihan digital.</p>
            </div>
            <button @click="showModal = true" class="flex items-center justify-center gap-2 bg-slate-900 hover:bg-indigo-600 text-white px-8 py-4 rounded-2xl font-black shadow-xl shadow-slate-900/10 hover:shadow-indigo-600/30 hover:-translate-y-1 transition-all duration-300 cursor-pointer">
              <PlusCircle :size="22" /> Buat Pemilihan Baru
            </button>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div class="bg-white p-6 rounded-[2rem] border border-slate-100 shadow-sm flex items-center gap-5 relative overflow-hidden group">
              <div class="absolute -right-6 -top-6 w-24 h-24 bg-blue-50 rounded-full group-hover:scale-150 transition-transform duration-500"></div>
              <div class="p-4 bg-blue-100 text-blue-600 rounded-2xl relative z-10"><BarChart3 :size="28" /></div>
              <div class="relative z-10">
                <p class="text-sm font-black text-slate-400 uppercase tracking-widest mb-1">Total Pemilihan</p>
                <p class="text-3xl font-black text-slate-800">{{ totalPolls }} <span class="text-sm text-slate-400 font-medium">Acara</span></p>
              </div>
            </div>

            <div class="bg-white p-6 rounded-[2rem] border border-slate-100 shadow-sm flex items-center gap-5 relative overflow-hidden group">
              <div class="absolute -right-6 -top-6 w-24 h-24 bg-emerald-50 rounded-full group-hover:scale-150 transition-transform duration-500"></div>
              <div class="p-4 bg-emerald-100 text-emerald-600 rounded-2xl relative z-10"><Activity :size="28" /></div>
              <div class="relative z-10">
                <p class="text-sm font-black text-slate-400 uppercase tracking-widest mb-1">Pemilihan Aktif</p>
                <p class="text-3xl font-black text-slate-800">{{ activePolls }} <span class="text-sm text-slate-400 font-medium">Berjalan</span></p>
              </div>
            </div>

            <div class="bg-indigo-600 p-6 rounded-[2rem] shadow-xl shadow-indigo-600/20 flex items-center gap-5 relative overflow-hidden group text-white">
              <div class="absolute -right-6 -top-6 w-32 h-32 bg-white/10 rounded-full group-hover:scale-150 transition-transform duration-700"></div>
              <div class="p-4 bg-white/20 backdrop-blur-sm text-white rounded-2xl relative z-10"><ShieldCheck :size="28" /></div>
              <div class="relative z-10">
                <p class="text-sm font-black text-indigo-200 uppercase tracking-widest mb-1">Total Suara Masuk</p>
                <p class="text-3xl font-black">{{ totalVotesCast }} <span class="text-sm text-indigo-200 font-medium">Pemilih</span></p>
              </div>
            </div>
          </div>
        </div>

        <div class="h-px w-full bg-slate-200 mb-10"></div>

        <div v-if="isLoading" class="flex flex-col items-center justify-center py-20 text-slate-400">
          <Loader2 class="animate-spin mb-4 text-indigo-500" :size="48" />
          <p class="font-bold text-lg">Membaca arsip satelit...</p>
        </div>

        <div v-else-if="polls.length === 0" class="max-w-2xl mx-auto bg-white border-2 border-slate-200 border-dashed rounded-[3rem] p-16 flex flex-col items-center justify-center text-center shadow-sm">
          <div class="bg-slate-50 p-6 rounded-full mb-6">
            <PlusCircle :size="56" class="text-slate-300" />
          </div>
          <h3 class="text-3xl font-black text-slate-800 mb-3">Server Bersih</h3>
          <p class="text-slate-500 text-lg">Anda belum memulai kegiatan pemilihan apa pun. Jadilah inisiator demokrasi hari ini.</p>
        </div>

        <div v-else class="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-8">
          <div v-for="poll in polls" :key="poll.ID" class="bg-white border border-slate-100 p-8 rounded-[2.5rem] shadow-sm hover:shadow-2xl hover:shadow-indigo-500/10 transition-all duration-300 group flex flex-col h-full relative overflow-hidden">

            <div class="absolute top-0 left-0 w-full h-1.5" :class="poll.is_active ? 'bg-emerald-400' : 'bg-slate-300'"></div>

            <div class="flex justify-between items-start mb-6">
              <span v-if="poll.is_active" class="px-4 py-1.5 bg-emerald-50 text-emerald-600 border border-emerald-100 rounded-full text-xs font-black uppercase tracking-widest flex items-center gap-1.5">
                <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span> LIVE
              </span>
              <span v-else class="px-4 py-1.5 bg-slate-50 text-slate-500 border border-slate-200 rounded-full text-xs font-black uppercase tracking-widest">
                Ditutup
              </span>

              <span class="text-sm font-black text-slate-400 flex items-center gap-1 bg-slate-50 px-3 py-1 rounded-lg">
                <Users :size="14" /> {{ poll.positions?.reduce((total, pos) => total + (pos.options?.length || 0), 0) || 0 }} Kandidat
              </span>
            </div>

            <div class="flex-1">
              <h3 class="text-2xl font-black mb-3 text-slate-800 leading-tight group-hover:text-indigo-600 transition-colors">{{ poll.title }}</h3>
              <p class="text-slate-500 text-sm font-medium mb-8 line-clamp-2 leading-relaxed">{{ poll.description || 'Tidak ada catatan khusus untuk pemilihan ini.' }}</p>

              <div class="mb-8 p-5 bg-slate-50 rounded-3xl border border-slate-100 relative overflow-hidden">
                <div class="flex justify-between items-end mb-3 relative z-10">
                  <span class="text-xs font-black text-slate-400 uppercase tracking-widest flex items-center gap-1.5"><Key :size="14" /> Partisipasi (Token)</span>
                  <div class="text-right">
                    <span class="text-xl font-black text-slate-800">{{ poll.tokens?.filter(t => t.is_used).length || 0 }}</span>
                    <span class="text-sm font-bold text-slate-400"> / {{ poll.tokens?.length || 0 }}</span>
                  </div>
                </div>
                <div class="w-full h-3 bg-white rounded-full overflow-hidden shadow-inner border border-slate-200 relative z-10">
                  <div class="h-full bg-gradient-to-r from-indigo-500 to-blue-500 transition-all duration-1000 ease-out"
                       :style="{ width: `${ poll.tokens?.length ? ((poll.tokens.filter(t => t.is_used).length / poll.tokens.length) * 100) : 0 }%` }">
                  </div>
                </div>
              </div>
            </div>

            <div class="grid grid-cols-2 gap-3 mt-auto pt-4">
              <NuxtLink :to="`/results/${poll.slug}`" target="_blank" class="flex items-center justify-center gap-2 bg-slate-900 text-white py-3.5 rounded-2xl text-sm font-black hover:bg-indigo-600 shadow-lg shadow-slate-900/10 active:scale-95 transition-all">
                <PieChart :size="18" /> Monitor Live
              </NuxtLink>

              <button @click="handleDeletePoll(poll.ID)" class="flex items-center justify-center gap-2 bg-white border-2 border-slate-200 text-slate-600 hover:text-red-600 hover:border-red-200 hover:bg-red-50 py-3.5 rounded-2xl text-sm font-black active:scale-95 transition-all">
                <Trash2 :size="18" /> Hapus
              </button>

              <NuxtLink :to="`/print/${poll.ID}`" target="_blank" class="flex items-center justify-center gap-2 bg-indigo-50 text-indigo-600 py-3.5 rounded-2xl text-sm font-black hover:bg-indigo-100 active:scale-95 transition-all col-span-2 border border-indigo-100 mt-1">
                🖨️ Cetak Token Fisik
              </NuxtLink>
            </div>

          </div>
        </div>

      </main>
    </div>
  </div>
  <CreatePollModal
      v-if="showModal"
      @close="showModal = false"
      @submit="handleCreatePoll"
  />
</template>