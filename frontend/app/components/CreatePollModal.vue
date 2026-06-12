<script setup>
import { ref } from 'vue'
import { X, Plus, Trash2, ImagePlus, CheckCircle2 } from 'lucide-vue-next'

const emit = defineEmits(['close', 'submit'])
const isSubmitting = ref(false)

const form = ref({
  title: '',
  description: '',
  voter_count: 50,
  expires_at: '',
  positions: [
    {
      title: 'Ketua Umum',
      options: [
        { value: '', vision: '', mission: '', photo: '' },
        { value: '', vision: '', mission: '', photo: '' }
      ]
    }
  ]
})

const addPosition = () => {
  form.value.positions.push({
    title: '',
    options: [
      { value: '', vision: '', mission: '', photo: '' },
      { value: '', vision: '', mission: '', photo: '' }
    ]
  })
}

const removePosition = (posIndex) => {
  if (form.value.positions.length > 1) {
    form.value.positions.splice(posIndex, 1)
  }
}

const addOption = (posIndex) => {
  form.value.positions[posIndex].options.push({ value: '', vision: '', mission: '', photo: '' })
}

const removeOption = (posIndex, optIndex) => {
  if (form.value.positions[posIndex].options.length > 1) {
    form.value.positions[posIndex].options.splice(optIndex, 1)
  }
}

const handlePhotoUpload = (event, posIndex, optIndex) => {
  const file = event.target.files[0]
  if (!file) return

  if (file.size > 2 * 1024 * 1024) {
    alert("Ukuran foto terlalu besar! Maksimal 2MB.")
    event.target.value = ''
    return
  }

  const reader = new FileReader()
  reader.onload = (e) => {
    form.value.positions[posIndex].options[optIndex].photo = e.target.result
  }
  reader.readAsDataURL(file)
}

const submitForm = () => {
  if (!form.value.title.trim()) {
    alert("Judul pemilihan tidak boleh kosong!")
    return
  }

  for (const pos of form.value.positions) {
    if (!pos.title.trim()) {
      alert("Ada nama jabatan yang masih kosong!")
      return
    }
    for (const opt of pos.options) {
      if (!opt.value.trim()) {
        alert(`Nama kandidat di jabatan "${pos.title}" tidak boleh kosong!`)
        return
      }
    }
  }

  isSubmitting.value = true

  const payload = JSON.parse(JSON.stringify(form.value))

  if (payload.expires_at) {
    payload.expires_at = new Date(payload.expires_at).toISOString()
  } else {
    delete payload.expires_at
  }

  payload.voter_count = Number(payload.voter_count)

  emit('submit', payload)

  setTimeout(() => {
    isSubmitting.value = false
  }, 1000)
}
</script>

<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm animate-in fade-in">
    <div class="bg-white rounded-[2rem] shadow-2xl w-full max-w-4xl max-h-[90vh] flex flex-col overflow-hidden">

      <div class="flex items-center justify-between p-6 md:p-8 border-b border-slate-100">
        <div>
          <h2 class="text-2xl font-black text-slate-800">Buat Pemilihan Baru</h2>
          <p class="text-slate-500 font-medium text-sm mt-1">Sistem Multi-Jabatan & Unggah Profil Kandidat</p>
        </div>
        <button @click="$emit('close')" class="p-2 bg-slate-100 hover:bg-slate-200 text-slate-600 rounded-full transition-colors cursor-pointer">
          <X :size="24" />
        </button>
      </div>

      <div class="flex-1 overflow-y-auto p-6 md:p-8 bg-slate-50/50">
        <form id="pollForm" @submit.prevent="submitForm" class="space-y-8">

          <div class="space-y-6 bg-white p-6 rounded-2xl border border-slate-200">
            <div>
              <label class="block text-xs font-black text-slate-500 uppercase tracking-widest mb-2">Judul Pemilihan</label>
              <input type="text" v-model="form.title" placeholder="Contoh: Pemilu Raya OSIS 2026" class="w-full p-4 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-indigo-500 outline-none font-bold text-slate-800 transition-all">
            </div>
            <div>
              <label class="block text-xs font-black text-slate-500 uppercase tracking-widest mb-2">Deskripsi (Opsional)</label>
              <textarea v-model="form.description" rows="2" placeholder="Tuliskan aturan atau deskripsi singkat..." class="w-full p-4 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-indigo-500 outline-none font-medium text-slate-600 transition-all"></textarea>
            </div>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label class="block text-xs font-black text-slate-500 uppercase tracking-widest mb-2">Jumlah Token (Pemilih)</label>
                <input type="number" v-model="form.voter_count" min="1" class="w-full p-4 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-indigo-500 outline-none font-bold text-slate-800 transition-all">
              </div>
              <div>
                <label class="block text-xs font-black text-slate-500 uppercase tracking-widest mb-2">Tenggat Waktu (Opsional)</label>
                <input type="datetime-local" v-model="form.expires_at" class="w-full p-4 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-indigo-500 outline-none font-bold text-slate-800 transition-all">
              </div>
            </div>
          </div>

          <div class="flex items-center gap-4">
            <div class="h-px bg-slate-300 flex-1"></div>
            <span class="text-xs font-black text-slate-400 uppercase tracking-widest">Struktur Jabatan & Kandidat</span>
            <div class="h-px bg-slate-300 flex-1"></div>
          </div>

          <div v-for="(pos, posIndex) in form.positions" :key="posIndex" class="p-6 bg-slate-100 border-2 border-slate-200 rounded-[2rem] relative space-y-6">

            <button v-if="form.positions.length > 1" type="button" @click="removePosition(posIndex)" class="absolute -top-3 -right-3 p-2 bg-red-100 text-red-600 hover:bg-red-500 hover:text-white rounded-full shadow-sm transition-colors z-10">
              <Trash2 :size="16" />
            </button>

            <div>
              <label class="block text-xs font-black text-indigo-600 uppercase tracking-widest mb-2">Nama Jabatan {{ posIndex + 1 }}</label>
              <input type="text" v-model="pos.title" placeholder="Misal: Ketua, Sekretaris, Bendahara" class="w-full p-4 bg-white border border-slate-300 rounded-xl focus:ring-2 focus:ring-indigo-500 outline-none font-black text-indigo-900 text-lg transition-all shadow-sm">
            </div>

            <div class="space-y-4">
              <div v-for="(opt, optIndex) in pos.options" :key="optIndex" class="p-5 bg-white border border-slate-200 rounded-2xl relative shadow-sm flex flex-col md:flex-row gap-5">

                <div class="w-full md:w-32 flex-shrink-0 flex flex-col items-center justify-center">
                  <label class="w-24 h-32 border-2 border-dashed border-slate-300 hover:border-indigo-500 rounded-xl flex flex-col items-center justify-center cursor-pointer overflow-hidden bg-slate-50 group transition-colors relative">
                    <img v-if="opt.photo" :src="opt.photo" class="w-full h-full object-cover absolute inset-0 z-0" />
                    <div class="z-10 flex flex-col items-center justify-center text-slate-400 group-hover:text-indigo-500 bg-white/60 w-full h-full p-2 text-center backdrop-blur-[2px]" :class="{'opacity-0 hover:opacity-100 transition-opacity': opt.photo}">
                      <ImagePlus :size="24" class="mb-1" />
                      <span class="text-[10px] font-bold uppercase">Foto</span>
                    </div>
                    <input type="file" accept="image/png, image/jpeg, image/jpg" class="hidden" @change="handlePhotoUpload($event, posIndex, optIndex)">
                  </label>
                </div>

                <div class="flex-1 space-y-3">
                  <div class="flex gap-2 items-center">
                    <input type="text" v-model="opt.value" :placeholder="`Nama Kandidat ${optIndex + 1}`" class="flex-1 p-3 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-indigo-500 outline-none font-black text-slate-700">
                    <button type="button" @click="removeOption(posIndex, optIndex)" class="p-3 text-slate-400 hover:text-red-500 transition-colors rounded-xl bg-slate-50 hover:bg-red-50 border border-slate-200">
                      <Trash2 :size="20" />
                    </button>
                  </div>
                  <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
                    <textarea v-model="opt.vision" rows="2" placeholder="Visi kandidat..." class="p-3 bg-slate-50 border border-slate-200 rounded-xl text-sm outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/20"></textarea>
                    <textarea v-model="opt.mission" rows="2" placeholder="Misi kandidat..." class="p-3 bg-slate-50 border border-slate-200 rounded-xl text-sm outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/20"></textarea>
                  </div>
                </div>

              </div>

              <button type="button" @click="addOption(posIndex)" class="text-sm font-bold text-indigo-600 hover:text-indigo-800 flex items-center gap-1.5 px-2">
                <Plus :size="16" /> Tambah Kandidat untuk {{ pos.title || 'jabatan ini' }}
              </button>
            </div>
          </div>

          <button type="button" @click="addPosition" class="w-full py-4 border-2 border-dashed border-indigo-300 text-indigo-600 hover:bg-indigo-50 rounded-2xl font-black transition-colors flex items-center justify-center gap-2">
            <Plus :size="20" /> Tambah Jabatan Baru (Misal: Wakil / Sekretaris)
          </button>

        </form>
      </div>

      <div class="p-6 md:p-8 border-t border-slate-100 bg-white">
        <button type="button" @click="submitForm" :disabled="isSubmitting" class="w-full flex items-center justify-center gap-2 bg-slate-900 hover:bg-indigo-600 text-white p-4 rounded-2xl font-black text-lg shadow-xl shadow-slate-900/20 active:scale-95 transition-all cursor-pointer">
          <CheckCircle2 v-if="!isSubmitting" :size="22" />
          <span v-if="isSubmitting" class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
          {{ isSubmitting ? 'Menyimpan Konfigurasi...' : 'Simpan & Buat Polling' }}
        </button>
      </div>

    </div>
  </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background-color: #cbd5e1; border-radius: 20px; }
</style>


