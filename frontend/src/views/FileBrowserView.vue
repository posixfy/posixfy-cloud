<template>
  <div class="h-screen flex flex-col">
    <AppHeader />
    <div class="flex flex-1 overflow-hidden">
      <AppSidebar />
      <main class="flex-1 overflow-y-auto p-6">
        <Breadcrumb />
        <div class="flex items-center gap-3 mb-5">
          <el-button :icon="Upload" @click="uploadVisible = true">Upload</el-button>
          <el-button :icon="FolderAdd" @click="mkdirVisible = true">New Folder</el-button>
          <el-switch
            v-model="fs.showHidden"
            class="ml-auto"
            inline-prompt
            active-text="Show hidden"
            inactive-text="Hide system files"
          />
        </div>
        <FileList @refresh="refresh" />
        <UploadDialog v-model="uploadVisible" @done="refresh" />
        <MkdirDialog v-model="mkdirVisible" @done="refresh" />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { Upload, FolderAdd } from '@element-plus/icons-vue'
import { useFsStore } from '@/stores/fs'
import { useFileWatcher } from '@/composables/useFileWatcher'
import AppHeader from '@/components/layout/AppHeader.vue'
import AppSidebar from '@/components/layout/AppSidebar.vue'
import Breadcrumb from '@/components/file/Breadcrumb.vue'
import FileList from '@/components/file/FileList.vue'
import UploadDialog from '@/components/file/UploadDialog.vue'
import MkdirDialog from '@/components/file/MkdirDialog.vue'

const fs = useFsStore()
const uploadVisible = ref(false)
const mkdirVisible = ref(false)

useFileWatcher()

onMounted(async () => {
  await fs.fetchMounts()
  await fs.fetchFiles()
})

watch(
  () => [fs.currentMount, fs.currentPath],
  () => fs.fetchFiles()
)

function refresh() {
  fs.fetchFiles()
}
</script>
