<template>
  <aside class="bg-white border-r border-gray-200 w-56 shrink-0 overflow-y-auto">
    <div class="px-4 pt-5 pb-2">
      <span class="text-xs uppercase tracking-wider text-gray-400 font-medium">Mount Points</span>
    </div>
    <nav class="px-2 pb-4 flex flex-col gap-0.5">
      <div
        v-for="m in fs.mounts"
        :key="m.name"
        class="flex items-center gap-2.5 px-3 py-2 rounded-lg text-sm cursor-pointer transition-colors"
        :class="fs.currentMount === m.name
          ? 'bg-gray-100 text-gray-900 font-medium'
          : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'"
        @click="handleSelect(m.name)"
      >
        <svg class="w-4 h-4 shrink-0" :class="fs.currentMount === m.name ? 'text-gray-700' : 'text-gray-400'" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 12.75V12A2.25 2.25 0 014.5 9.75h15A2.25 2.25 0 0121.75 12v.75m-8.69-6.44l-2.12-2.12a1.5 1.5 0 00-1.061-.44H4.5A2.25 2.25 0 002.25 6v12a2.25 2.25 0 002.25 2.25h15A2.25 2.25 0 0021.75 18V9a2.25 2.25 0 00-2.25-2.25h-5.379a1.5 1.5 0 01-1.06-.44z" />
        </svg>
        <span>{{ m.name }}</span>
      </div>
    </nav>
  </aside>
</template>

<script setup lang="ts">
import { useFsStore } from '@/stores/fs'

const fs = useFsStore()

function handleSelect(name: string) {
  fs.navigateTo(name, '/')
}
</script>
