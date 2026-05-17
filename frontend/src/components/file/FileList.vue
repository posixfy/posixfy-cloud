<template>
  <el-table
    :data="sortedFiles"
    v-loading="fs.loading"
    @row-click="handleRowClick"
    style="width: 100%"
    :header-cell-style="{ background: 'transparent', color: '#9ca3af', fontSize: '12px', textTransform: 'uppercase', fontWeight: '500', letterSpacing: '0.05em' }"
    :row-style="{ cursor: 'pointer' }"
    :row-class-name="() => 'file-row'"
  >
    <el-table-column label="Name" sortable sort-by="name" min-width="300">
      <template #default="{ row }">
        <div class="flex items-center gap-2.5">
          <svg v-if="row.is_dir" class="w-4.5 h-4.5 text-gray-400 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 12.75V12A2.25 2.25 0 014.5 9.75h15A2.25 2.25 0 0121.75 12v.75m-8.69-6.44l-2.12-2.12a1.5 1.5 0 00-1.061-.44H4.5A2.25 2.25 0 002.25 6v12a2.25 2.25 0 002.25 2.25h15A2.25 2.25 0 0021.75 18V9a2.25 2.25 0 00-2.25-2.25h-5.379a1.5 1.5 0 01-1.06-.44z" />
          </svg>
          <svg v-else class="w-4.5 h-4.5 text-gray-300 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m2.25 0H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />
          </svg>
          <span class="text-sm text-gray-700">{{ row.name }}</span>
        </div>
      </template>
    </el-table-column>
    <el-table-column label="Size" width="120" sortable sort-by="size">
      <template #default="{ row }">
        <span class="text-sm text-gray-400">{{ row.is_dir ? '-' : formatSize(row.size) }}</span>
      </template>
    </el-table-column>
    <el-table-column label="Modified" width="180" sortable sort-by="modified">
      <template #default="{ row }">
        <span class="text-sm text-gray-400">{{ formatDate(row.modified) }}</span>
      </template>
    </el-table-column>
    <el-table-column label="Actions" width="200" fixed="right">
      <template #default="{ row }">
        <div class="flex items-center gap-1">
          <el-button
            v-if="!row.is_dir"
            size="small"
            text
            :icon="Download"
            @click.stop="handleDownload(row)"
          >
            Download
          </el-button>
          <el-button
            size="small"
            text
            type="danger"
            :icon="Delete"
            @click.stop="handleDelete(row)"
          >
            Delete
          </el-button>
        </div>
      </template>
    </el-table-column>
  </el-table>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Download, Delete } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useFsStore } from '@/stores/fs'
import { downloadFile, deleteFile } from '@/api/fs'
import type { FileEntry } from '@/types'

const emit = defineEmits<{ refresh: [] }>()
const fs = useFsStore()

const sortedFiles = computed(() => {
  const items = [...fs.files]
  items.sort((a, b) => {
    if (a.is_dir !== b.is_dir) return a.is_dir ? -1 : 1
    return a.name.localeCompare(b.name)
  })
  return items
})

function handleRowClick(row: FileEntry) {
  if (row.is_dir) {
    fs.enterDir(row.name)
  }
}

async function handleDownload(row: FileEntry) {
  const filePath = fs.currentPath.endsWith('/')
    ? fs.currentPath + row.name
    : fs.currentPath + '/' + row.name
  try {
    const res = await downloadFile(fs.currentMount, filePath)
    const url = URL.createObjectURL(res.data)
    const a = document.createElement('a')
    a.href = url
    a.download = row.name
    a.click()
    URL.revokeObjectURL(url)
  } catch {
    ElMessage.error('Download failed')
  }
}

async function handleDelete(row: FileEntry) {
  await ElMessageBox.confirm(`Delete "${row.name}"?`, 'Confirm', { type: 'warning' })
  const filePath = fs.currentPath.endsWith('/')
    ? fs.currentPath + row.name
    : fs.currentPath + '/' + row.name
  try {
    await deleteFile(fs.currentMount, filePath)
    ElMessage.success('Deleted')
    emit('refresh')
  } catch {
    ElMessage.error('Delete failed')
  }
}

function formatSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  return (bytes / Math.pow(1024, i)).toFixed(i > 0 ? 1 : 0) + ' ' + units[i]
}

function formatDate(iso: string): string {
  if (!iso) return '-'
  return new Date(iso).toLocaleString()
}
</script>

<style scoped>
:deep(.file-row:hover > td) {
  background-color: #f9fafb !important;
}
:deep(.el-table th.el-table__cell) {
  border-bottom: 1px solid #f3f4f6 !important;
}
:deep(.el-table td.el-table__cell) {
  border-bottom: 1px solid #f3f4f6 !important;
  padding: 14px 0;
}
</style>
