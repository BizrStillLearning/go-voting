<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>

<style>
@import url('https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700;800;900&display=swap');

body {
  font-family: 'Plus Jakarta Sans', sans-serif;
  -webkit-font-smoothing: antialiased;
}

.swal2-popup {
  border-radius: 2.5rem !important;
  padding: 2.5rem 1.5rem !important;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25) !important;
  font-family: 'Plus Jakarta Sans', sans-serif !important;
  border: 1px solid #f1f5f9 !important;
}

.swal2-title {
  font-weight: 900 !important;
  font-size: 2.2rem !important;
  color: #0f172a !important;
  margin-bottom: 0.5rem !important;
}

.swal2-html-container {
  font-weight: 600 !important;
  color: #64748b !important;
  font-size: 1.15rem !important;
  line-height: 1.6 !important;
}

.swal2-confirm, .swal2-cancel {
  font-weight: 900 !important;
  border-radius: 1.2rem !important;
  padding: 1.2rem 2.5rem !important;
  font-size: 1.1rem !important;
  letter-spacing: 0.5px !important;
  text-transform: uppercase !important;
}

.swal2-confirm {
  box-shadow: 0 10px 15px -3px rgba(79, 70, 229, 0.4) !important;
  transition: all 0.2s ease !important;
}

.swal2-confirm:active {
  transform: scale(0.95) !important;
}
</style>