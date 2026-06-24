<script setup>
import { ref } from 'vue'
import { useRouter } from 'nuxt/app'
import Swal from 'sweetalert2'
import { Lock, User, Loader2, ShieldCheck } from 'lucide-vue-next'

definePageMeta({
  layout: false
})

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
      confirmButtonColor: '#DE1B1E',
      customClass: { popup: 'rounded-2xl' }
    })
    return
  }

  if (!form.value.password) {
    Swal.fire({
      icon: 'warning',
      title: 'Password Kosong',
      text: 'Harap masukkan kata sandi keamanan.',
      confirmButtonColor: '#DE1B1E',
      customClass: { popup: 'rounded-2xl' }
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
      customClass: { popup: 'rounded-2xl' }
    }).then(() => {
      router.push('/dashboard')
    })

  } catch (error) {
    Swal.fire({
      icon: 'error',
      title: 'Akses Ditolak',
      text: error.data?.error || 'Kredensial tidak valid. Silakan coba lagi.',
      confirmButtonColor: '#A51013',
      customClass: { popup: 'rounded-2xl' }
    })
    form.value.password = ''
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-poster-dark via-poster-base to-poster-glow flex items-center justify-center p-4 relative overflow-hidden font-sans">

    <div class="w-full max-w-md relative z-10 animate-in slide-in-from-bottom-8 fade-in duration-500">

      <div class="bg-white/95 backdrop-blur-xl border border-white/20 p-8 md:p-10 rounded-2xl shadow-2xl">

        <div class="flex flex-col items-center mb-10 text-center">
          <div class="w-16 h-16 bg-red-50 border border-red-100 rounded-full flex items-center justify-center mb-6 shadow-sm rotate-3 transition-transform hover:rotate-6">
<!--            <ShieldCheck class="text-poster-base" :size="32" stroke-width="2" />-->
            <img src="../assets/logo.jpg" alt="" class="rounded-full w-15 h-15">
          </div>
          <h1 class="text-2xl font-bold text-slate-900 tracking-tight mb-2">Musykom XVIII</h1>
          <p class="text-slate-500 font-medium text-sm">Masuk untuk mengelola sistem pemilihan</p>
        </div>

        <form @submit.prevent="handleLogin" class="space-y-5">

          <div class="space-y-4">
            <div class="relative group">
              <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none text-slate-400 group-focus-within:text-poster-base transition-colors">
                <User :size="20" />
              </div>
              <input
                  type="text"
                  v-model="form.username"
                  placeholder="Username"
                  class="w-full pl-11 pr-4 py-3.5 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-poster-base focus:border-poster-base outline-none transition-all text-slate-900 font-medium placeholder:text-slate-400"
              >
            </div>

            <div class="relative group">
              <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none text-slate-400 group-focus-within:text-poster-base transition-colors">
                <Lock :size="20" />
              </div>
              <input
                  type="password"
                  v-model="form.password"
                  placeholder="Kata Sandi"
                  class="w-full pl-11 pr-4 py-3.5 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-poster-base focus:border-poster-base outline-none transition-all text-slate-900 font-medium placeholder:text-slate-400"
              >
            </div>
          </div>

          <button
              type="submit"
              :disabled="isLoading"
              class="w-full bg-poster-dark hover:bg-poster-dark/90 disabled:bg-slate-300 disabled:text-slate-500 disabled:cursor-not-allowed text-white py-3.5 rounded-xl font-bold text-base transition-all shadow-sm active:scale-[0.98] flex items-center justify-center gap-2 mt-2"
          >
            <Loader2 v-if="isLoading" class="animate-spin" :size="20" />
            <span>{{ isLoading ? 'Mengesahkan...' : 'Login' }}</span>
          </button>

        </form>

      </div>

      <p class="text-center text-white/90 text-xs font-medium mt-8 drop-shadow-md">
        &copy; 2026 eVoting System. Secure Access.
      </p>

    </div>
  </div>
</template>