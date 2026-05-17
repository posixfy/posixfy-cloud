<template>
  <nav class="flex items-center gap-1.5 mb-5 text-sm">
    <span
      class="cursor-pointer font-medium text-gray-900 hover:text-gray-600 transition-colors"
      @click="fs.navigateTo(fs.currentMount, '/')"
    >
      {{ fs.currentMount }}
    </span>
    <template v-for="(part, index) in pathParts" :key="index">
      <span class="text-gray-300">/</span>
      <span
        class="cursor-pointer transition-colors"
        :class="index === pathParts.length - 1
          ? 'text-gray-500'
          : 'text-gray-900 font-medium hover:text-gray-600'"
        @click="navigateToIndex(index)"
      >
        {{ part }}
      </span>
    </template>
  </nav>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useFsStore } from '@/stores/fs'

const fs = useFsStore()

const pathParts = computed(() =>
  fs.currentPath.split('/').filter(Boolean)
)

function navigateToIndex(index: number) {
  const parts = pathParts.value.slice(0, index + 1)
  fs.navigateTo(fs.currentMount, '/' + parts.join('/'))
}
</script>
