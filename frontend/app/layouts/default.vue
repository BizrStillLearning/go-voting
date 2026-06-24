<template>
  <div class="flex h-screen bg-slate-50 dark:bg-slate-900 font-sans overflow-hidden selection:bg-indigo-200 transition-colors duration-500">

    <Sidebar />

    <div class="flex-1 flex flex-col h-screen overflow-hidden relative">

      <Header />

      <main class="flex-1 overflow-x-hidden overflow-y-auto p-6 md:p-10 lg:px-12 pb-6 flex flex-col">

        <div class="flex-1">
          <slot />
        </div>

        <footer class="mt-12 pt-6 border-t border-slate-200 dark:border-slate-800 text-center transition-colors duration-500">
          <p class="text-sm font-bold text-slate-400 dark:text-slate-500">
            &copy; 2026 Administrator E-Voting. Sistem Digital Terpadu.
          </p>
        </footer>

      </main>

    </div>
  </div>
</template>