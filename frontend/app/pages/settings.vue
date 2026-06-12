<script setup>
import { ref } from 'vue'
import { Save, KeyRound, AlertOctagon, Archive, ShieldAlert } from 'lucide-vue-next'
import { useAuthStore } from '~/stores/auth'

const authStore = useAuthStore()

const formPassword = ref({
  oldPass: '',
  newPass: '',
  confirmPass: ''
})

definePageMeta({ middleware: ['auth'] })

const handleChangePassword = async () => {
  if (formPassword.value.newPass !== formPassword.value.confirmPass) {
    alert("Sandi baru tidak cocok dengan konfirmasi!")
    return
  }
  try {
    await $fetch('http://localhost:8080/api/admin/password', {
      method: 'PUT',
      headers: { Authorization: `Bearer ${authStore.token}` },
      body: { old_password: formPassword.value.oldPass, new_password: formPassword.value.newPass }
    })
    alert("Sandi berhasil diperbarui!")
    formPassword.value = { oldPass: '', newPass: '', confirmPass: '' }
  } catch (error) {
    alert("Gagal memperbarui sandi.")
  }
}

const handleArchive = async () => {
  if(confirm("Tutup semua pemilihan yang masih aktif dan jadikan arsip?")) {
    try {
      await $fetch('http://localhost:8080/api/admin/archive', {
        method: 'POST',
        headers: { Authorization: `Bearer ${authStore.token}` }
      })
      alert("Proses arsip berhasil!")
    } catch (error) {
      alert("Gagal melakukan arsip.")
    }
  }
}

const handleResetDatabase = async () => {
  if (confirm("PERINGATAN KERAS! Aksi ini akan MENGHAPUS SEMUA DATA polling, token, dan hasil suara secara permanen. Lanjutkan?")) {
    if (prompt("Ketik 'RESET SAYA' untuk mengonfirmasi:") === "RESET SAYA") {
      try {
        await $fetch('http://localhost:8080/api/admin/reset', {
          method: 'DELETE',
          headers: { Authorization: `Bearer ${authStore.token}` }
        })
        alert("Sistem berhasil di-reset ke kondisi pabrik!")
        window.location.reload()
      } catch (error) {
        alert("Gagal melakukan reset.")
      }
    } else {
      alert("Reset dibatalkan demi keamanan.")
    }
  }
}
</script>

<template>
  <div class="flex h-screen bg-slate-50 font-sans overflow-hidden">
    <Sidebar />
    <div class="flex-1 flex flex-col h-screen overflow-hidden relative">
      <Header />
      <main class="flex-1 overflow-x-hidden overflow-y-auto p-6 md:p-8">

        <div class="mb-8">
          <h2 class="text-3xl font-black text-slate-800">Manajemen Sistem Utama</h2>
          <p class="text-slate-500 font-medium mt-1">Pusat kontrol kredensial Superadmin dan pembersihan operasional server.</p>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">

          <div class="bg-white p-8 rounded-[2rem] border border-slate-200 shadow-sm h-fit">
            <div class="flex items-center gap-3 mb-6 border-b border-slate-100 pb-4">
              <KeyRound class="text-indigo-600" />
              <h3 class="text-xl font-black text-slate-800">Kredensial Superadmin</h3>
            </div>

            <form @submit.prevent="handleChangePassword" class="space-y-5">
              <div class="space-y-2">
                <label class="text-xs font-black text-slate-500 uppercase tracking-widest">Sandi Saat Ini</label>
                <input type="password" v-model="formPassword.oldPass" required class="w-full p-4 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-indigo-500 outline-none">
              </div>
              <div class="space-y-2">
                <label class="text-xs font-black text-slate-500 uppercase tracking-widest">Sandi Baru</label>
                <input type="password" v-model="formPassword.newPass" required class="w-full p-4 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-indigo-500 outline-none">
              </div>
              <div class="space-y-2">
                <label class="text-xs font-black text-slate-500 uppercase tracking-widest">Konfirmasi Sandi Baru</label>
                <input type="password" v-model="formPassword.confirmPass" required class="w-full p-4 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-indigo-500 outline-none">
              </div>
              <div class="pt-4">
                <button type="submit" class="w-full flex items-center justify-center gap-2 bg-slate-900 hover:bg-indigo-600 text-white px-8 py-3.5 rounded-xl font-black shadow-lg transition-all">
                  <Save :size="18" /> Perbarui Sandi Akses
                </button>
              </div>
            </form>
          </div>

          <div class="space-y-6">

            <div class="bg-white p-8 rounded-[2rem] border border-slate-200 shadow-sm">
              <div class="flex items-center gap-3 mb-4">
                <Archive class="text-indigo-600" />
                <h3 class="font-black text-slate-800 text-xl">Arsip Otomatis</h3>
              </div>
              <p class="text-sm text-slate-500 mb-6 leading-relaxed">Pindahkan semua pemilihan yang kedaluwarsa ke tabel arsip agar Dashboard tidak penuh.</p>
              <button @click="handleArchive" class="px-6 py-3 bg-slate-100 hover:bg-slate-200 text-slate-700 font-black rounded-xl transition-all">
                Jalankan Pengarsipan
              </button>
            </div>

            <div class="bg-red-50 p-8 rounded-[2rem] border border-red-100 shadow-sm">
              <div class="flex items-center gap-3 mb-4">
                <AlertOctagon class="text-red-600" />
                <h3 class="font-black text-red-800 text-xl">Pembersihan Total</h3>
              </div>
              <p class="text-sm text-red-600/80 mb-6 leading-relaxed">Menghapus seluruh tabel dan me-reset sistem ke kondisi pabrik. Gunakan hanya jika periode tahunan berakhir.</p>
              <button @click="handleResetDatabase" class="flex items-center justify-center gap-2 w-full py-3 bg-red-600 hover:bg-red-700 text-white font-black rounded-xl shadow-lg shadow-red-600/30 transition-all cursor-pointer">
                <ShieldAlert :size="18" /> Format / Reset Database
              </button>
            </div>

          </div>
        </div>
      </main>
    </div>
  </div>
</template>



