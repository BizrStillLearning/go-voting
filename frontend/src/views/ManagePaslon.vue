<script setup>
import { ref, onMounted } from 'vue'
import {
  UserPlus, PenTool, Users,
  Save, Trash2, LayoutGrid, Info
} from 'lucide-vue-next'

const paslonList = ref([])
const form = ref({
  id: null,
  ketua: '',
  wakil: '',
  nomor: null,
  visi: '',
  warna: 'blue'
})

const warnaMap = {
  blue: "from-blue-500 to-indigo-500",
  emerald: "from-emerald-500 to-teal-500",
  rose: "from-rose-400 to-pink-500",
  amber: "from-amber-400 to-orange-500"
}

const loadData = () => {
  const saved = localStorage.getItem('go_voting_paslon')
  if (saved) paslonList.value = JSON.parse(saved)
}

const saveData = () => {
  localStorage.setItem('go_voting_paslon', JSON.stringify(paslonList.value))
}

const handleSubmit = () => {
  if (form.value.id) {
    const index = paslonList.value.findIndex(p => p.id === form.value.id)
    paslonList.value[index] = { ...form.value }
  } else {
    paslonList.value.push({
      ...form.value,
      id: Date.now(),
      suara: 0
    })
  }
  resetForm()
  saveData()
}

const resetForm = () => {
  form.value = { id: null, ketua: '', wakil: '', nomor: null, visi: '', warna: 'blue' }
}

const editPaslon = (item) => {
  form.value = { ...item }
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

const deletePaslon = (id) => {
  if(confirm('Hapus paslon ini? Semua data terkait akan hilang.')) {
    paslonList.value = paslonList.value.filter(p => p.id !== id)
    saveData()
  }
}

onMounted(loadData)
</script>

<template>
  <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 animate-in fade-in slide-in-from-bottom-4 duration-500">

    <section class="lg:col-span-4">
      <div class="bg-white p-6 rounded-3xl shadow-sm border border-slate-200 sticky top-8">
        <div class="flex items-center gap-3 mb-6">
          <div class="p-2 bg-indigo-50 rounded-lg text-indigo-600">
            <PenTool :size="20" />
          </div>
          <h2 class="text-xl font-bold text-slate-800">
            {{ form.id ? 'Edit Paslon' : 'Tambah Paslon' }}
          </h2>
        </div>

        <form @submit.prevent="handleSubmit" class="space-y-5">
          <div class="space-y-1.5">
            <label class="text-xs font-black uppercase text-slate-400 tracking-wider">Calon Ketua</label>
            <input v-model="form.ketua" type="text" required class="w-full px-4 py-3 rounded-2xl bg-slate-50 border border-slate-200 focus:ring-2 focus:ring-indigo-500 outline-none transition-all" placeholder="Masukkan nama ketua">
          </div>

          <div class="space-y-1.5">
            <label class="text-xs font-black uppercase text-slate-400 tracking-wider">Calon Wakil</label>
            <input v-model="form.wakil" type="text" required class="w-full px-4 py-3 rounded-2xl bg-slate-50 border border-slate-200 focus:ring-2 focus:ring-indigo-500 outline-none transition-all" placeholder="Masukkan nama wakil">
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-1.5">
              <label class="text-xs font-black uppercase text-slate-400 tracking-wider">No. Urut</label>
              <input v-model.number="form.nomor" type="number" required class="w-full px-4 py-3 rounded-2xl bg-slate-50 border border-slate-200 focus:ring-2 focus:ring-indigo-500 outline-none transition-all" placeholder="0">
            </div>
            <div class="space-y-1.5">
              <label class="text-xs font-black uppercase text-slate-400 tracking-wider">Tema</label>
              <select v-model="form.warna" class="w-full px-4 py-3 rounded-2xl bg-slate-50 border border-slate-200 outline-none appearance-none cursor-pointer">
                <option value="blue">Biru</option>
                <option value="emerald">Hijau</option>
                <option value="rose">Merah</option>
                <option value="amber">Kuning</option>
              </select>
            </div>
          </div>

          <div class="space-y-1.5">
            <label class="text-xs font-black uppercase text-slate-400 tracking-wider">Visi & Misi</label>
            <textarea v-model="form.visi" rows="3" class="w-full px-4 py-3 rounded-2xl bg-slate-50 border border-slate-200 focus:ring-2 focus:ring-indigo-500 outline-none transition-all" placeholder="Apa visi misi paslon ini?"></textarea>
          </div>

          <div class="pt-2 space-y-3">
            <button type="submit" class="w-full bg-slate-900 hover:bg-indigo-600 text-white font-bold py-4 rounded-2xl shadow-lg transition-all flex items-center justify-center gap-2 group">
              <Save :size="18" class="group-hover:scale-110 transition-transform" />
              {{ form.id ? 'Perbarui Data' : 'Simpan Paslon' }}
            </button>
            <button v-if="form.id" @click="resetForm" type="button" class="w-full py-2 text-slate-400 text-sm font-bold hover:text-slate-600">
              Batalkan Perubahan
            </button>
          </div>
        </form>
      </div>
    </section>

    <section class="lg:col-span-8">
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-xl font-bold text-slate-800 flex items-center gap-2">
          <LayoutGrid :size="20" class="text-indigo-500" />
          Preview Paslon
        </h2>
        <span class="px-4 py-1 bg-white border border-slate-200 rounded-full text-xs font-bold text-slate-500">
          Total: {{ paslonList.length }} Pasangan
        </span>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div v-for="p in paslonList" :key="p.id" class="bg-white rounded-3xl border border-slate-200 overflow-hidden group hover:shadow-xl hover:shadow-slate-200/50 transition-all duration-300">
          <div :class="warnaMap[p.warna]" class="h-3 bg-gradient-to-r"></div>

          <div class="p-6">
            <div class="flex justify-between items-start mb-4">
              <div class="w-12 h-12 rounded-2xl bg-slate-50 flex items-center justify-center text-xl font-black text-slate-800 border border-slate-100">
                {{ p.nomor }}
              </div>
              <div class="flex gap-2">
                <button @click="editPaslon(p)" class="p-2.5 bg-slate-50 hover:bg-indigo-50 text-slate-400 hover:text-indigo-600 rounded-xl transition-colors">
                  <PenTool :size="18" />
                </button>
                <button @click="deletePaslon(p)" class="p-2.5 bg-slate-50 hover:bg-rose-50 text-slate-400 hover:text-rose-600 rounded-xl transition-colors">
                  <Trash2 :size="18" />
                </button>
              </div>
            </div>

            <h3 class="text-xl font-bold text-slate-900 leading-tight">
              {{ p.ketua }}
              <span class="block text-slate-400 text-sm font-medium mt-1">Wakil: {{ p.wakil }}</span>
            </h3>

            <div class="mt-6 pt-5 border-t border-slate-50">
              <p class="text-sm text-slate-500 italic leading-relaxed">
                <Info :size="14" class="inline mr-1 mb-1" />
                {{ p.visi || 'Visi misi belum ditentukan untuk pasangan ini.' }}
              </p>
            </div>
          </div>
        </div>

        <div v-if="paslonList.length === 0" class="col-span-full py-32 flex flex-col items-center justify-center border-2 border-dashed border-slate-200 rounded-[2.5rem] bg-white/50">
          <Users :size="48" class="text-slate-200 mb-4" />
          <h3 class="text-slate-400 font-bold text-lg">Daftar Paslon Kosong</h3>
          <p class="text-slate-300 text-sm">Gunakan form di samping untuk mulai menambahkan.</p>
        </div>
      </div>
    </section>
  </div>
</template>