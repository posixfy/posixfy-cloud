<template>
  <el-table
    :data="users"
    style="width: 100%"
    :header-cell-style="{ background: 'transparent', color: '#9ca3af', fontSize: '12px', textTransform: 'uppercase', fontWeight: '500', letterSpacing: '0.05em' }"
    :row-class-name="() => 'user-row'"
  >
    <el-table-column prop="id" label="ID" width="60" />
    <el-table-column prop="username" label="Username" width="150" />
    <el-table-column prop="uid" label="UID" width="80" />
    <el-table-column prop="gid" label="GID" width="80" />
    <el-table-column prop="groups" label="Groups" width="150" />
    <el-table-column prop="role" label="Role" width="100">
      <template #default="{ row }">
        <span
          class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
          :class="row.role === 'admin'
            ? 'bg-red-50 text-red-700'
            : 'bg-gray-100 text-gray-600'"
        >
          {{ row.role }}
        </span>
      </template>
    </el-table-column>
    <el-table-column prop="created_at" label="Created" width="180">
      <template #default="{ row }">
        <span class="text-sm text-gray-400">{{ new Date(row.created_at).toLocaleString() }}</span>
      </template>
    </el-table-column>
    <el-table-column label="Actions" width="180" fixed="right">
      <template #default="{ row }">
        <div class="flex items-center gap-1">
          <el-button size="small" text :icon="Edit" @click="$emit('edit', row)">Edit</el-button>
          <el-button size="small" text type="danger" :icon="Delete" @click="$emit('delete', row)">
            Delete
          </el-button>
        </div>
      </template>
    </el-table-column>
  </el-table>
</template>

<script setup lang="ts">
import { Edit, Delete } from '@element-plus/icons-vue'
import type { User } from '@/types'

defineProps<{ users: User[] }>()
defineEmits<{ edit: [user: User]; delete: [user: User]; refresh: [] }>()
</script>

<style scoped>
:deep(.user-row:hover > td) {
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
