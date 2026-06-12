<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'nuxt/app'
import Swal from 'sweetalert2'
import { Check, Send, CheckCircle, AlertTriangle, Loader2, ArrowRight, ArrowLeft, UserCircle2 } from 'lucide-vue-next'

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
const isSuccess = ref(false)

const activeVisionId = ref(null)

onMounted(async () => {
  try {
    const res = await $fetch(`http://localhost:8080/api/poll/${slug}`)
    if (res.is_closed) {
      fetchError.value = "Mohon maaf, sesi pemilihan ini telah ditutup."
    } else {
      poll.value = res.poll
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
    await $fetch(`http://localhost:8080/api/vote/${poll.value.ID}/verify`, {
      method: 'POST',
      body: { token: String(tokenInput.value).toUpperCase() }
    })

    currentStep.value = 1
  } catch (error) {
    Swal.fire({
      icon: 'error',
      title: 'Akses Ditolak!',
      text: error.data?.error || 'Token tidak valid',
      confirmButtonColor: '#4f46e5',
      customClass: {
        popup: 'rounded-3xl'
      }
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
    await $fetch(`http://localhost:8080/api/vote/${poll.value.ID}`, {
      method: 'POST',
      body: {
        token: String(tokenInput.value).toUpperCase(),
        votes: payloadVotes
      }
    })

    Swal.fire({
      icon: 'success',
      title: 'Suara Masuk!',
      text: 'Terima kasih atas partisipasi Anda.',
      showConfirmButton: false,
      timer: 3000,
      timerProgressBar: true,
      customClass: {
        popup: 'rounded-3xl'
      }
    }).then(() => {
      tokenInput.value = ''
      selections.value = {}
      currentStep.value = 0
    })

  } catch (error) {
    Swal.fire({
      icon: 'error',
      title: 'Gagal Menyimpan',
      text: error.data?.error || "Terjadi kesalahan sistem.",
      confirmButtonColor: '#4f46e5',
      customClass: { popup: 'rounded-3xl' }
    })
    if (error.response?.status === 401 || error.response?.status === 403) {
      currentStep.value = 0
    }
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="bg-slate-100 min-h-screen flex items-center justify-center p-4 py-12 selection:bg-indigo-100 font-sans">
    <div class="bg-white p-8 md:p-12 rounded-[3rem] shadow-xl w-full max-w-3xl border border-slate-100 relative overflow-hidden transition-all duration-500" :class="{'max-w-xl': currentStep === 0}">

      <div v-if="isLoading" class="flex flex-col items-center justify-center py-20">
        <Loader2 class="animate-spin mb-4 text-indigo-500" :size="48" />
        <p class="text-slate-500 font-bold">Mempersiapkan bilik suara...</p>
      </div>

      <div v-else-if="fetchError" class="text-center py-10">
        <AlertTriangle class="mx-auto text-amber-500 mb-4" :size="64" />
        <h2 class="text-2xl font-black text-slate-800 mb-2">Akses Ditolak</h2>
        <p class="text-slate-500 font-medium">{{ fetchError }}</p>
      </div>

      <div v-else class="transition-opacity duration-300 relative">

        <div class="text-center mb-10">
          <h1 class="text-3xl font-black text-slate-800">{{ poll?.title }}</h1>
          <div class="mt-4 flex justify-center gap-2" v-if="currentStep > 0">
            <div v-for="n in totalPositions" :key="n" class="h-2 rounded-full transition-all duration-300" :class="n === currentStep ? 'w-8 bg-indigo-600' : (n < currentStep ? 'w-4 bg-indigo-300' : 'w-4 bg-slate-200')"></div>
          </div>
        </div>

        <form @submit.prevent class="space-y-8">

          <div v-if="currentStep === 0" class="animate-in slide-in-from-right-4 fade-in duration-300">
            <div class="bg-slate-50 p-8 rounded-[2rem] border border-slate-200 focus-within:border-indigo-500 focus-within:ring-4 ring-indigo-50 transition-all text-center">
              <label class="block text-sm font-black text-slate-500 uppercase tracking-widest mb-4">Masukkan Kode Token Anda</label>
              <input type="text" v-model="tokenInput" placeholder="VOTE-XXXX" required
                     class="w-full p-5 bg-white border-2 border-slate-200 focus:border-indigo-500 rounded-2xl text-3xl text-center font-mono font-black text-indigo-700 outline-none transition-all uppercase tracking-widest shadow-sm mb-6">

              <button type="button" @click="verifyAndNext" :disabled="!tokenInput || isVerifying" class="w-full bg-slate-900 hover:bg-indigo-600 disabled:bg-slate-300 disabled:cursor-not-allowed text-white p-5 rounded-2xl font-black text-xl shadow-xl active:scale-95 transition-all flex justify-center items-center gap-3">
                <Loader2 v-if="isVerifying" class="animate-spin" :size="20" />
                {{ isVerifying ? 'Mengecek...' : 'Mulai Memilih' }} <ArrowRight v-if="!isVerifying" :size="20" />
              </button>
            </div>
          </div>

          <div v-else-if="currentPosition" class="animate-in slide-in-from-right-4 fade-in duration-300">
            <div class="flex items-center justify-between mb-4 px-2">
              <label class="block text-sm font-black text-indigo-500 uppercase tracking-widest">Pilih {{ currentPosition.title }}</label>
              <span class="text-xs font-bold text-slate-400 bg-slate-100 px-3 py-1 rounded-full">Langkah {{ currentStep }} dari {{ totalPositions }}</span>
            </div>

            <div class="grid grid-cols-1 gap-4">
              <div v-for="opt in currentPosition.options" :key="opt.ID"
                   class="border-2 border-slate-200 rounded-[2rem] cursor-pointer hover:border-indigo-300 transition-all overflow-hidden bg-white group"
                   :class="{ 'border-indigo-600 ring-4 ring-indigo-50': selections[currentPosition.ID] === opt.ID }"
                   @click="selections[currentPosition.ID] = opt.ID">

                <div class="p-5 flex flex-col md:flex-row items-center gap-5 relative">

                  <div class="w-24 h-32 rounded-xl overflow-hidden bg-slate-100 flex-shrink-0 border border-slate-200 flex items-center justify-center">
                    <img v-if="opt.photo_url" :src="opt.photo_url" class="w-full h-full object-cover" />
                    <UserCircle2 v-else :size="48" class="text-slate-300" />
                  </div>

                  <div class="flex-1 text-center md:text-left w-full">
                    <span class="text-2xl font-black text-slate-800 block mb-2">{{ opt.value }}</span>
                    <button type="button" @click="toggleVision(opt.ID, $event)" class="text-xs font-bold text-indigo-600 bg-indigo-50 hover:bg-indigo-100 px-3 py-1.5 rounded-xl transition-colors">
                      {{ activeVisionId === opt.ID ? 'Tutup Info' : 'Lihat Visi & Misi' }}
                    </button>
                  </div>

                  <div class="w-8 h-8 rounded-full border-2 border-slate-300 flex items-center justify-center text-white transition-all absolute top-5 right-5 md:relative md:top-0 md:right-0"
                       :class="{ 'bg-indigo-600 border-indigo-600 scale-110 shadow-lg shadow-indigo-200': selections[currentPosition.ID] === opt.ID }">
                    <Check :size="16" v-if="selections[currentPosition.ID] === opt.ID" />
                  </div>
                </div>

                <div v-if="activeVisionId === opt.ID" class="px-6 pb-6 pt-4 bg-slate-50 border-t border-slate-100 text-sm space-y-3 animate-in slide-in-from-top-2 duration-200">
                  <div>
                    <h4 class="font-black text-indigo-600 uppercase text-[11px] tracking-wider mb-0.5">Visi:</h4>
                    <p class="text-slate-700 font-medium leading-relaxed">{{ opt.vision || 'Tidak ada visi tertulis.' }}</p>
                  </div>
                  <div>
                    <h4 class="font-black text-indigo-600 uppercase text-[11px] tracking-wider mb-0.5">Misi:</h4>
                    <p class="text-slate-700 font-medium leading-relaxed whitespace-pre-line">{{ opt.mission || 'Tidak ada misi tertulis.' }}</p>
                  </div>
                </div>
              </div>
            </div>

            <div class="flex gap-4 mt-8">
              <button type="button" @click="prevStep" class="px-6 py-5 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-2xl font-black transition-colors flex items-center gap-2">
                <ArrowLeft :size="20" /> Kembali
              </button>

              <button v-if="currentStep < totalPositions" type="button" @click="nextStep" :disabled="!selections[currentPosition.ID]" class="flex-1 bg-slate-900 hover:bg-indigo-600 disabled:bg-slate-300 disabled:cursor-not-allowed text-white p-5 rounded-2xl font-black text-xl shadow-xl transition-all flex justify-center items-center gap-3">
                Lanjut <ArrowRight :size="20" />
              </button>

              <button v-else type="button" @click="submitVote" :disabled="isSubmitting || !selections[currentPosition.ID]" class="flex-1 bg-indigo-600 hover:bg-indigo-500 disabled:bg-slate-300 disabled:cursor-not-allowed text-white p-5 rounded-2xl font-black text-xl shadow-xl shadow-indigo-600/30 active:scale-95 transition-all flex justify-center items-center gap-3">
                <Loader2 v-if="isSubmitting" class="animate-spin" :size="20" />
                <Send v-else :size="20" />
                {{ isSubmitting ? 'Memproses...' : 'Selesai & Kirim Suara' }}
              </button>
            </div>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>