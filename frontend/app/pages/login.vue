<script setup>
import { ref } from 'vue'
import { useRouter } from 'nuxt/app'
import Swal from 'sweetalert2'
import { Lock, User, Loader2, ShieldCheck } from 'lucide-vue-next'

const router = useRouter()

const form = ref({
  username: '',
  password: ''
})

const isLoading = ref(false)

const handleLogin = async () => {
  if (!form.value.username.trim()) {
    Swal.fire({
      icon: 'warning',
      title: 'Username Kosong',
      text: 'Harap isi username Anda terlebih dahulu.',
      confirmButtonColor: '#4f46e5',
      customClass: { popup: 'rounded-3xl' }
    })
    return
  }

  if (!form.value.password) {
    Swal.fire({
      icon: 'warning',
      title: 'Password Kosong',
      text: 'Harap masukkan kata sandi keamanan.',
      confirmButtonColor: '#4f46e5',
      customClass: { popup: 'rounded-3xl' }
    })
    return
  }

  isLoading.value = true

  try {
    const res = await $fetch('http://localhost:8080/api/admin/login', {
      method: 'POST',
      body: form.value
    })

    const authCookie = useCookie('admin_token', { maxAge: 86400 })

    authCookie.value = res.token || 'AUTH-SESSION-' + new Date().getTime()

    Swal.fire({
      icon: 'success',
      title: 'Akses Diberikan',
      text: 'Otentikasi berhasil. Mengalihkan ke Dashboard...',
      showConfirmButton: false,
      timer: 1500,
      customClass: { popup: 'rounded-3xl' }
    }).then(() => {
      router.push('/dashboard')
    })

  } catch (error) {
    Swal.fire({
      icon: 'error',
      title: 'Akses Ditolak',
      text: error.data?.error || 'Kredensial tidak valid. Silakan coba lagi.',
      confirmButtonColor: '#ef4444',
      customClass: { popup: 'rounded-3xl' }
    })
    form.value.password = ''
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-slate-900 flex items-center justify-center p-4 relative overflow-hidden font-sans">

    <div class="absolute top-0 left-1/2 -translate-x-1/2 w-[800px] h-[500px] bg-indigo-600/20 rounded-full blur-[120px] pointer-events-none"></div>
    <div class="absolute bottom-0 right-0 w-[600px] h-[400px] bg-blue-600/10 rounded-full blur-[100px] pointer-events-none"></div>

    <div class="w-full max-w-md relative z-10 animate-in slide-in-from-bottom-8 fade-in duration-500">

      <div class="bg-white/10 backdrop-blur-xl border border-white/20 p-8 md:p-10 rounded-[2.5rem] shadow-2xl">

        <div class="flex flex-col items-center mb-10 text-center">
          <div class="w-20 h-20 bg-gradient-to-br from-indigo-500 to-blue-600 rounded-3xl flex items-center justify-center mb-6 shadow-lg shadow-indigo-500/30 rotate-3 transition-transform hover:rotate-6">
            <ShieldCheck class="text-white" :size="40" stroke-width="1.5" />
          </div>
          <h1 class="text-3xl font-black text-white tracking-tight mb-2">Admin Portal</h1>
          <p class="text-slate-400 font-medium">Masuk untuk mengelola sistem e-Voting</p>
        </div>

        <form @submit.prevent="handleLogin" class="space-y-6">

          <div class="space-y-4">
            <div class="relative group">
              <div class="absolute inset-y-0 left-0 pl-5 flex items-center pointer-events-none text-slate-400 group-focus-within:text-indigo-400 transition-colors">
                <User :size="20" />
              </div>
              <input
                  type="text"
                  v-model="form.username"
                  placeholder="Username"
                  class="w-full pl-12 pr-5 py-4 bg-slate-800/50 border border-slate-700/50 rounded-2xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 focus:bg-slate-800 outline-none transition-all text-white font-medium placeholder:text-slate-500 shadow-inner"
              >
            </div>

            <div class="relative group">
              <div class="absolute inset-y-0 left-0 pl-5 flex items-center pointer-events-none text-slate-400 group-focus-within:text-indigo-400 transition-colors">
                <Lock :size="20" />
              </div>
              <input
                  type="password"
                  v-model="form.password"
                  placeholder="Kata Sandi"
                  class="w-full pl-12 pr-5 py-4 bg-slate-800/50 border border-slate-700/50 rounded-2xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 focus:bg-slate-800 outline-none transition-all text-white font-medium placeholder:text-slate-500 shadow-inner"
              >
            </div>
          </div>

          <button
              type="submit"
              :disabled="isLoading"
              class="w-full bg-indigo-600 hover:bg-indigo-500 disabled:bg-slate-700 disabled:text-slate-400 disabled:cursor-not-allowed text-white py-4 rounded-2xl font-black text-lg transition-all shadow-lg shadow-indigo-600/25 active:scale-[0.98] flex items-center justify-center gap-3 mt-4"
          >
            <Loader2 v-if="isLoading" class="animate-spin" :size="24" />
            <span>{{ isLoading ? 'Mengesahkan...' : 'Masuk ke Sistem' }}</span>
          </button>

        </form>

      </div>

      <p class="text-center text-slate-500 text-sm font-medium mt-8">
        &copy; 2026 eVoting System. Secure Access.
      </p>

    </div>
  </div>
</template>