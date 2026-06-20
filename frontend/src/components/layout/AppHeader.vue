<template>
  <header class="bg-white border-b border-gray-200 h-14 flex items-center justify-between px-6">
    <div class="flex items-center gap-8">
      <span
        class="font-semibold text-lg text-gray-900 cursor-pointer"
        @click="$router.push(auth.isAdmin ? '/admin' : '/')"
      >
        Posixfy
      </span>
      <nav class="flex items-center gap-1">
        <a
          v-if="!auth.isAdmin"
          class="px-3 py-1.5 text-sm rounded-md cursor-pointer transition-colors"
          :class="route.path === '/' ? 'text-gray-900 bg-gray-100 font-medium' : 'text-gray-500 hover:text-gray-900 hover:bg-gray-50'"
          @click="router.push('/')"
        >
          Files
        </a>
        <a
          v-if="auth.isAdmin"
          class="px-3 py-1.5 text-sm rounded-md cursor-pointer transition-colors"
          :class="route.path === '/admin' ? 'text-gray-900 bg-gray-100 font-medium' : 'text-gray-500 hover:text-gray-900 hover:bg-gray-50'"
          @click="router.push('/admin')"
        >
          Admin
        </a>
      </nav>
    </div>
    <div class="flex items-center gap-4">
      <span class="text-sm text-gray-600">{{ auth.user?.username }}</span>
      <button
        class="text-sm text-gray-400 hover:text-gray-900 cursor-pointer transition-colors"
        @click="handleLogout"
      >
        Logout
      </button>
    </div>
  </header>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

function handleLogout() {
  auth.logout()
  router.push('/login')
}
</script>
