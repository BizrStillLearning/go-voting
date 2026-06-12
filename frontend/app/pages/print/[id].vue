<script setup>
import { onMounted, computed } from 'vue'
import { useRoute } from 'nuxt/app'
import { usePollStore } from '~/stores/poll'

const route = useRoute()
const pollStore = usePollStore()

const pollId = Number(route.params.id)

const currentPoll = computed(() => {
  return pollStore.polls.find(p => p.ID === pollId)
})

onMounted(() => {
  if (pollStore.polls.length === 0) {
    pollStore.fetchPolls().then(() => {
      setTimeout(() => window.print(), 800)
    })
  } else {
    setTimeout(() => window.print(), 800)
  }
})
</script>

<template>
  <div class="bg-white text-black p-8 min-h-screen font-sans" v-if="currentPoll">

    <div class="mb-8 text-center border-b-2 border-black pb-4">
      <h1 class="text-2xl font-black uppercase tracking-widest">Dokumen Rahasia</h1>
      <p class="font-bold text-lg">Kupon Token Pemilih: {{ currentPoll.title }}</p>
      <p class="text-sm italic">Gunting kupon di bawah ini dan berikan ke peserta yang telah diverifikasi.</p>
    </div>

    <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
      <div v-for="tokenObj in currentPoll.tokens" :key="tokenObj.ID"
           class="border-2 border-dashed border-black p-4 rounded-xl flex flex-col items-center justify-center text-center">
        <p class="text-xs font-bold uppercase mb-2 border-b border-black w-full pb-1">Token Pemilih</p>
        <p class="text-2xl font-black font-mono tracking-widest"
           :class="{'line-through text-gray-300': tokenObj.is_used}">
          {{ tokenObj.TokenString }}
        </p>
        <p class="text-[10px] mt-2 font-bold">
          {{ tokenObj.is_used ? 'SUDAH TERPAKAI' : 'Berlaku untuk 1 kali suara' }}
        </p>
      </div>
    </div>
  </div>

  <div v-else class="p-10 text-center font-bold">Memuat dokumen cetak...</div>
</template>

<style scoped>
@media print {
  body { background-color: white !important; }
  @page { margin: 1cm; }
}
</style>

