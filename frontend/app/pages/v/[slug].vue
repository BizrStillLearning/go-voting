<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'nuxt/app'
import Swal from 'sweetalert2'
import { Check, Send, AlertTriangle, Loader2, ArrowRight, ArrowLeft, UserCircle2 } from 'lucide-vue-next'

definePageMeta({
  layout: false
})

const route = useRoute()
const slug = route.params.slug

const poll = ref(null)
const isLoading = ref(true)
const fetchError = ref('')

const currentStep = ref(0)
const tokenInput = ref('')
const selections = ref({})

const isSubmitting = ref(false)
const isVerifying = ref(false)
const submitError = ref('')

const activeVisionId = ref(null)

const getApiUrl = () => {
  if (typeof window !== 'undefined') {
    return `http://${window.location.hostname}:8080/api`
  }
  return 'http://localhost:8080/api'
}

onMounted(async () => {
  try {
    const res = await $fetch(`${getApiUrl()}/poll/${slug}`)
    if (res.is_closed) {
      fetchError.value = "Mohon maaf, sesi pemilihan ini telah ditutup."
    } else {
      poll.value = res.poll || res.data || res
    }
  } catch (error) {
    fetchError.value = "Pemilihan tidak ditemukan atau URL tidak valid."
  } finally {
    isLoading.value = false
  }
})

const totalPositions = computed(() => poll.value?.positions?.length || 0)
const currentPosition = computed(() => {
  if (currentStep.value > 0 && currentStep.value <= totalPositions.value) {
    return poll.value.positions[currentStep.value - 1]
  }
  return null
})

const verifyAndNext = async () => {
  if (!tokenInput.value) return
  isVerifying.value = true

  try {
    await $fetch(`${getApiUrl()}/v/${poll.value.ID}/verify`, {
      method: 'POST',
      body: { token: String(tokenInput.value).toUpperCase() }
    })

    currentStep.value = 1
  } catch (error) {
    const isDarkMode = document.documentElement.classList.contains('dark')
    Swal.fire({
      icon: 'error',
      title: 'Akses Ditolak!',
      text: error.data?.error || 'Token tidak valid atau sudah digunakan.',
      confirmButtonColor: '#A51013', // poster-dark
      background: isDarkMode ? '#1e293b' : '#ffffff',
      color: isDarkMode ? '#f8fafc' : '#0f172a',
      customClass: { popup: 'rounded-2xl' }
    })
    tokenInput.value = ''
  } finally {
    isVerifying.value = false
  }
}

const nextStep = () => {
  if (currentStep.value > 0 && !selections.value[currentPosition.value.ID]) return
  submitError.value = ''
  currentStep.value++
}

const prevStep = () => {
  if (currentStep.value > 0) currentStep.value--
}

const toggleVision = (id, event) => {
  event.stopPropagation()
  activeVisionId.value = activeVisionId.value === id ? null : id
}

const submitVote = async () => {
  isSubmitting.value = true
  submitError.value = ''

  const payloadVotes = Object.keys(selections.value).map(posId => ({
    position_id: Number(posId),
    option_id: selections.value[posId]
  }))

  try {
    await $fetch(`${getApiUrl()}/v/${poll.value.ID}`, {
      method: 'POST',
      body: {
        token: String(tokenInput.value).toUpperCase(),
        votes: payloadVotes
      }
    })

    const isDarkMode = document.documentElement.classList.contains('dark')

    Swal.fire({
      icon: 'success',
      title: 'Suara Masuk!',
      text: 'Terima kasih atas partisipasi Anda.',
      showConfirmButton: false,
      timer: 2500,
      timerProgressBar: true,
      background: isDarkMode ? '#1e293b' : '#ffffff',
      color: isDarkMode ? '#f8fafc' : '#0f172a',
      customClass: { popup: 'rounded-2xl' }
    }).then(() => {
      localStorage.clear()
      window.location.reload()
    })

  } catch (error) {
    const isDarkMode = document.documentElement.classList.contains('dark')
    Swal.fire({
      icon: 'error',
      title: 'Gagal Menyimpan',
      text: error.data?.error || "Terjadi kesalahan koneksi ke server.",
      confirmButtonColor: '#A51013',
      background: isDarkMode ? '#1e293b' : '#ffffff',
      color: isDarkMode ? '#f8fafc' : '#0f172a',
      customClass: { popup: 'rounded-2xl' }
    })
    if (error.response?.status === 401 || error.response?.status === 403) {
      currentStep.value = 0
      tokenInput.value = ''
    }
  } finally {
    isSubmitting.value = false
  }
}

const getImageUrl = (path) => {
  if (!path) return ''

  if (path.startsWith('http')) return path

  const host = typeof window !== 'undefined' ? window.location.hostname : 'localhost'
  const prefix = path.startsWith('/') ? '' : '/'

  return `http://${host}:8080${prefix}${path}`
}
</script>

<template>
  <div class="min-h-screen bg-slate-50 dark:bg-slate-950 flex items-center justify-center p-4 py-12 selection:bg-red-100 dark:selection:bg-red-900/30 font-sans transition-colors duration-300">
    <div class="bg-white dark:bg-slate-900 p-8 md:p-12 rounded-3xl shadow-xl w-full max-w-3xl border border-slate-200 dark:border-slate-800 relative overflow-hidden transition-all duration-500" :class="{'max-w-xl': currentStep === 0}">

      <div v-if="isLoading" class="flex flex-col items-center justify-center py-20">
        <Loader2 class="animate-spin mb-4 text-poster-base dark:text-poster-light" :size="48" />
        <p class="text-slate-500 dark:text-slate-400 font-semibold">Mempersiapkan bilik suara...</p>
      </div>

      <div v-else-if="fetchError" class="text-center py-10">
        <AlertTriangle class="mx-auto text-poster-base dark:text-poster-light mb-4" :size="64" />
        <h2 class="text-2xl font-bold text-slate-800 dark:text-white mb-2">Akses Ditolak</h2>
        <p class="text-slate-500 dark:text-slate-400 font-medium">{{ fetchError }}</p>
      </div>

      <div v-else class="transition-opacity duration-300 relative">

        <div class="text-center mb-10">
          <h1 class="text-2xl md:text-3xl font-bold text-slate-800 dark:text-white">{{ poll?.title }}</h1>
          <div class="mt-4 flex justify-center gap-2" v-if="currentStep > 0">
            <div v-for="n in totalPositions" :key="n" class="h-2 rounded-full transition-all duration-300"
                 :class="n === currentStep ? 'w-8 bg-poster-base dark:bg-poster-light' : (n < currentStep ? 'w-4 bg-poster-base/30 dark:bg-poster-base/50' : 'w-4 bg-slate-200 dark:bg-slate-700')"></div>
          </div>
        </div>

        <form @submit.prevent class="space-y-8">

          <div v-if="currentStep === 0" class="animate-in slide-in-from-right-4 fade-in duration-300">
            <div class="bg-slate-50 dark:bg-slate-800/50 p-8 rounded-2xl border border-slate-200 dark:border-slate-700 focus-within:border-poster-base dark:focus-within:border-poster-light focus-within:ring-4 ring-poster-base/10 dark:ring-poster-base/20 transition-all text-center">
              <label class="block text-sm font-bold text-slate-500 dark:text-slate-400 uppercase tracking-widest mb-4">Masukkan Kode Token Anda</label>

              <input type="text" v-model="tokenInput" placeholder="VOTE-XXXX" required
                     class="w-full p-5 bg-white dark:bg-slate-900 border-2 border-slate-200 dark:border-slate-700 focus:border-poster-base dark:focus:border-poster-light rounded-xl text-2xl md:text-3xl text-center font-mono font-bold text-poster-dark dark:text-poster-light outline-none transition-all uppercase tracking-widest shadow-sm mb-6">

              <button type="button" @click="verifyAndNext" :disabled="!tokenInput || isVerifying"
                      class="w-full bg-poster-base dark:bg-poster-dark hover:bg-poster-dark dark:hover:bg-poster-base disabled:bg-slate-300 dark:disabled:bg-slate-700 disabled:text-slate-500 disabled:cursor-not-allowed text-white p-4 rounded-xl font-bold text-lg shadow-lg active:scale-95 transition-all flex justify-center items-center gap-3">
                <Loader2 v-if="isVerifying" class="animate-spin" :size="20" />
                {{ isVerifying ? 'Mengecek...' : 'Mulai Memilih' }} <ArrowRight v-if="!isVerifying" :size="20" />
              </button>
            </div>
          </div>

          <div v-else-if="currentPosition" class="animate-in slide-in-from-right-4 fade-in duration-300">
            <div class="flex items-center justify-between mb-4 px-2">
              <label class="block text-sm font-bold text-poster-base dark:text-poster-light uppercase tracking-widest">Pilih {{ currentPosition.title }}</label>
              <span class="text-xs font-semibold text-slate-500 dark:text-slate-400 bg-slate-100 dark:bg-slate-800 px-3 py-1 rounded-full">Langkah {{ currentStep }} dari {{ totalPositions }}</span>
            </div>

            <div class="grid grid-cols-1 gap-4">
              <div v-for="opt in currentPosition.options" :key="opt.ID"
                   class="border-2 border-slate-200 dark:border-slate-700 rounded-2xl cursor-pointer hover:border-poster-base/50 dark:hover:border-poster-light/50 transition-all overflow-hidden bg-white dark:bg-slate-800 group"
                   :class="{ 'border-poster-base dark:border-poster-light ring-4 ring-poster-base/10 dark:ring-poster-base/20': selections[currentPosition.ID] === opt.ID }"
                   @click="selections[currentPosition.ID] = opt.ID">

                <div class="p-5 flex flex-col md:flex-row items-center gap-5 relative">
                  <div class="w-24 h-32 rounded-xl overflow-hidden bg-slate-100 dark:bg-slate-900 flex-shrink-0 border border-slate-200 dark:border-slate-700 flex items-center justify-center">
                    <img v-if="opt.photo_url" :src="getImageUrl(opt.photo_url)" class="w-full h-full object-cover" />
                    <UserCircle2 v-else :size="48" class="text-slate-300 dark:text-slate-600" />
                  </div>

                  <div class="flex-1 text-center md:text-left w-full">
                    <span class="text-xl md:text-2xl font-bold text-slate-800 dark:text-white block mb-2">{{ opt.value }}</span>
                    <button type="button" @click="toggleVision(opt.ID, $event)"
                            class="text-xs font-semibold text-poster-base dark:text-poster-light bg-red-50 dark:bg-poster-base/10 hover:bg-red-100 dark:hover:bg-poster-base/20 px-3 py-1.5 rounded-lg transition-colors">
                      {{ activeVisionId === opt.ID ? 'Tutup Info' : 'Lihat Visi & Misi' }}
                    </button>
                  </div>

                  <div class="w-8 h-8 rounded-full border-2 border-slate-300 dark:border-slate-600 flex items-center justify-center text-white transition-all absolute top-5 right-5 md:relative md:top-0 md:right-0"
                       :class="{ 'bg-poster-base border-poster-base dark:bg-poster-light dark:border-poster-light scale-110 shadow-lg shadow-poster-base/20 dark:shadow-none': selections[currentPosition.ID] === opt.ID }">
                    <Check :size="16" v-if="selections[currentPosition.ID] === opt.ID" />
                  </div>
                </div>

                <div v-if="activeVisionId === opt.ID" class="px-6 pb-6 pt-4 bg-slate-50 dark:bg-slate-900/50 border-t border-slate-100 dark:border-slate-700 text-sm space-y-3 animate-in slide-in-from-top-2 duration-200">
                  <div>
                    <h4 class="font-bold text-poster-base dark:text-poster-light uppercase text-[11px] tracking-wider mb-0.5">Visi:</h4>
                    <p class="text-slate-700 dark:text-slate-300 font-medium leading-relaxed">{{ opt.vision || 'Tidak ada visi tertulis.' }}</p>
                  </div>
                  <div>
                    <h4 class="font-bold text-poster-base dark:text-poster-light uppercase text-[11px] tracking-wider mb-0.5">Misi:</h4>
                    <p class="text-slate-700 dark:text-slate-300 font-medium leading-relaxed whitespace-pre-line">{{ opt.mission || 'Tidak ada misi tertulis.' }}</p>
                  </div>
                </div>
              </div>
            </div>

            <div class="flex gap-4 mt-8">
              <button type="button" @click="prevStep"
                      class="px-5 md:px-6 py-4 bg-slate-100 dark:bg-slate-800 hover:bg-slate-200 dark:hover:bg-slate-700 text-slate-700 dark:text-slate-300 rounded-xl font-bold transition-colors flex items-center gap-2">
                <ArrowLeft :size="20" /> <span class="hidden md:inline">Kembali</span>
              </button>

              <button v-if="currentStep < totalPositions" type="button" @click="nextStep" :disabled="!selections[currentPosition.ID]"
                      class="flex-1 bg-slate-900 dark:bg-slate-700 hover:bg-poster-dark dark:hover:bg-poster-base disabled:bg-slate-300 dark:disabled:bg-slate-800 disabled:text-slate-500 disabled:cursor-not-allowed text-white p-4 rounded-xl font-bold text-lg shadow-lg transition-all flex justify-center items-center gap-3">
                Lanjut <ArrowRight :size="20" />
              </button>

              <button v-else type="button" @click="submitVote" :disabled="isSubmitting || !selections[currentPosition.ID]"
                      class="flex-1 bg-poster-base hover:bg-poster-dark disabled:bg-slate-300 dark:disabled:bg-slate-800 disabled:text-slate-500 disabled:cursor-not-allowed text-white p-4 rounded-xl font-bold text-lg shadow-lg shadow-poster-base/30 dark:shadow-none active:scale-95 transition-all flex justify-center items-center gap-3">
                <Loader2 v-if="isSubmitting" class="animate-spin" :size="20" />
                <Send v-else :size="20" />
                <span class="hidden md:inline">{{ isSubmitting ? 'Memproses...' : 'Selesai & Kirim Suara' }}</span>
                <span class="md:hidden">{{ isSubmitting ? 'Proses...' : 'Kirim' }}</span>
              </button>
            </div>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>