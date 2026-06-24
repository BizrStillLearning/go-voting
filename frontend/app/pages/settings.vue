<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Save, KeyRound, AlertOctagon, Archive, ShieldAlert } from 'lucide-vue-next'
import { useAuthStore } from '~/stores/auth'
import Swal from 'sweetalert2'

const authStore = useAuthStore()
const router = useRouter()

const formPassword = ref({
  oldPass: '',
  newPass: '',
  confirmPass: ''
})

definePageMeta({ middleware: ['auth'] })

const getApiUrl = (endpoint) => {
  const host = typeof window !== 'undefined' ? window.location.hostname : 'localhost'
  return `http://${host}:8080/api${endpoint}`
}

const handleChangePassword = async () => {
  if (formPassword.value.newPass !== formPassword.value.confirmPass) {
    Swal.fire({ icon: 'error', title: 'Kesalahan Validasi', text: 'Sandi baru tidak cocok dengan konfirmasi!', customClass: { popup: 'rounded-3xl' } })
    return
  }
  try {
    await $fetch(getApiUrl('/admin/password'), {
      method: 'PUT',
      headers: { Authorization: `Bearer ${authStore.token}` },
      body: { old_password: formPassword.value.oldPass, new_password: formPassword.value.newPass }
    })

    await Swal.fire({
      icon: 'success',
      title: 'Berhasil!',
      text: 'Sandi berhasil diperbarui! Silakan login kembali.',
      confirmButtonColor: '#053123',
      customClass: { popup: 'rounded-3xl' }
    })

    formPassword.value = { oldPass: '', newPass: '', confirmPass: '' }

    authStore.token = null
    if (typeof localStorage !== 'undefined') {
      localStorage.removeItem('auth_token')
    }
    router.push('/login')

  } catch (error) {
    const errorMsg = error.response?._data?.error || 'Sandi saat ini salah atau terjadi masalah di server.'
    Swal.fire({ icon: 'error', title: 'Gagal', text: errorMsg, customClass: { popup: 'rounded-3xl' } })
  }
}

const handleArchive = async () => {
  const result = await Swal.fire({
    title: 'Arsip Otomatis?',
    text: "Tutup semua pemilihan yang masih aktif dan jadikan arsip?",
    icon: 'question',
    showCancelButton: true,
    confirmButtonColor: '#053123', // poster-dark
    cancelButtonColor: '#94a3b8',
    confirmButtonText: 'Ya, Arsipkan',
    cancelButtonText: 'Batal',
    customClass: { popup: 'rounded-3xl' }
  })

  if (result.isConfirmed) {
    try {
      await $fetch(getApiUrl('/admin/archive'), {
        method: 'POST',
        headers: { Authorization: `Bearer ${authStore.token}` }
      })
      Swal.fire({ icon: 'success', title: 'Selesai!', text: 'Proses arsip berhasil!', timer: 1500, showConfirmButton: false, customClass: { popup: 'rounded-3xl' } })
    } catch (error) {
      Swal.fire({ icon: 'error', title: 'Gagal', text: 'Gagal melakukan arsip.', customClass: { popup: 'rounded-3xl' } })
    }
  }
}

const handleResetDatabase = async () => {
  const result = await Swal.fire({
    title: 'PERINGATAN KERAS!',
    text: "Aksi ini akan MENGHAPUS SEMUA DATA polling, token, dan hasil suara secara permanen. Lanjutkan?",
    icon: 'warning',
    showCancelButton: true,
    confirmButtonColor: '#ef4444',
    cancelButtonColor: '#94a3b8',
    confirmButtonText: 'Ya, Lanjutkan',
    cancelButtonText: 'Batal',
    customClass: { popup: 'rounded-3xl' }
  })

  if (result.isConfirmed) {
    const { value: text } = await Swal.fire({
      title: 'Konfirmasi Reset',
      input: 'text',
      inputLabel: "Ketik 'RESET SAYA' untuk mengonfirmasi:",
      inputPlaceholder: 'RESET SAYA',
      showCancelButton: true,
      customClass: { popup: 'rounded-3xl' }
    })

    if (text === 'RESET SAYA') {
      try {
        await $fetch(getApiUrl('/admin/reset'), {
          method: 'DELETE',
          headers: { Authorization: `Bearer ${authStore.token}` }
        })
        Swal.fire({ icon: 'success', title: 'Berhasil!', text: 'Sistem berhasil di-reset ke kondisi pabrik!', showConfirmButton: false, timer: 2000, customClass: { popup: 'rounded-3xl' } }).then(() => {
          window.location.reload()
        })
      } catch (error) {
        Swal.fire({ icon: 'error', title: 'Gagal', text: 'Gagal melakukan reset.', customClass: { popup: 'rounded-3xl' } })
      }
    } else if (text) {
      Swal.fire({ icon: 'info', title: 'Dibatalkan', text: 'Konfirmasi tidak sesuai. Reset dibatalkan demi keamanan.', customClass: { popup: 'rounded-3xl' } })
    }
  }
}
</script>

<template>
  <div class="animate-in fade-in duration-500 min-h-screen">
    <div class="mb-8">
      <h2 class="text-3xl font-black text-slate-800 dark:text-white tracking-tight">Manajemen Sistem Utama</h2>
      <p class="text-slate-500 dark:text-slate-400 font-medium mt-1">Pusat kontrol kredensial Superadmin dan pembersihan operasional server.</p>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <div class="bg-white dark:bg-slate-900 p-8 rounded-[2rem] border border-slate-200 dark:border-slate-800 shadow-sm h-fit transition-colors duration-500">
        <div class="flex items-center gap-3 mb-6 border-b border-slate-100 dark:border-slate-800 pb-4">
          <KeyRound class="text-poster-base dark:text-poster-light" />
          <h3 class="text-xl font-black text-slate-800 dark:text-white">Kredensial Superadmin</h3>
        </div>

        <form @submit.prevent="handleChangePassword" class="space-y-5">
          <div class="space-y-2">
            <label class="text-xs font-black text-slate-500 dark:text-slate-400 uppercase tracking-widest">Sandi Saat Ini</label>
            <input type="password" v-model="formPassword.oldPass" required class="w-full p-4 bg-slate-50 dark:bg-slate-950 border border-slate-200 dark:border-slate-700 rounded-xl focus:ring-2 focus:ring-poster-base outline-none text-slate-800 dark:text-white transition-colors duration-500">
          </div>
          <div class="space-y-2">
            <label class="text-xs font-black text-slate-500 dark:text-slate-400 uppercase tracking-widest">Sandi Baru</label>
            <input type="password" v-model="formPassword.newPass" required class="w-full p-4 bg-slate-50 dark:bg-slate-950 border border-slate-200 dark:border-slate-700 rounded-xl focus:ring-2 focus:ring-poster-base outline-none text-slate-800 dark:text-white transition-colors duration-500">
          </div>
          <div class="space-y-2">
            <label class="text-xs font-black text-slate-500 dark:text-slate-400 uppercase tracking-widest">Konfirmasi Sandi Baru</label>
            <input type="password" v-model="formPassword.confirmPass" required class="w-full p-4 bg-slate-50 dark:bg-slate-950 border border-slate-200 dark:border-slate-700 rounded-xl focus:ring-2 focus:ring-poster-base outline-none text-slate-800 dark:text-white transition-colors duration-500">
          </div>
          <div class="pt-4">
            <button type="submit" class="w-full flex items-center justify-center gap-2 bg-poster-base dark:bg-poster-light hover:bg-poster-dark dark:hover:bg-poster-base text-white px-8 py-3.5 rounded-xl font-black shadow-lg transition-all">
              <Save :size="18" /> Perbarui Sandi Akses
            </button>
          </div>
        </form>
      </div>

      <div class="space-y-6">
        <div class="bg-white dark:bg-slate-900 p-8 rounded-[2rem] border border-slate-200 dark:border-slate-800 shadow-sm transition-colors duration-500">
          <div class="flex items-center gap-3 mb-4">
            <Archive class="text-poster-base dark:text-poster-light" />
            <h3 class="font-black text-slate-800 dark:text-white text-xl">Arsip Otomatis</h3>
          </div>
          <p class="text-sm text-slate-500 dark:text-slate-400 mb-6 leading-relaxed">Pindahkan semua pemilihan yang kedaluwarsa ke tabel arsip agar Dashboard tidak penuh.</p>
          <button @click="handleArchive" class="px-6 py-3 bg-slate-100 dark:bg-slate-800 hover:bg-poster-base hover:text-white dark:hover:bg-poster-light text-slate-700 dark:text-slate-200 font-black rounded-xl transition-all">
            Jalankan Pengarsipan
          </button>
        </div>

        <div class="bg-red-50 dark:bg-red-950/30 p-8 rounded-[2rem] border border-red-100 dark:border-red-900/50 shadow-sm transition-colors duration-500">
          <div class="flex items-center gap-3 mb-4">
            <AlertOctagon class="text-red-600 dark:text-red-500" />
            <h3 class="font-black text-red-800 dark:text-red-400 text-xl">Pembersihan Total</h3>
          </div>
          <p class="text-sm text-red-600/80 dark:text-red-400/80 mb-6 leading-relaxed">Menghapus seluruh tabel dan me-reset sistem ke kondisi pabrik. Gunakan hanya jika periode tahunan berakhir.</p>
          <button @click="handleResetDatabase" class="flex items-center justify-center gap-2 w-full py-3 bg-red-600 hover:bg-red-700 text-white font-black rounded-xl shadow-lg shadow-red-600/30 transition-all cursor-pointer">
            <ShieldAlert :size="18" /> Format / Reset Database
          </button>
        </div>
      </div>
    </div>
  </div>
</template>


n