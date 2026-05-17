<template>
  <div class="h-screen flex flex-col">
    <AppHeader />
    <main class="flex-1 overflow-y-auto p-6">
      <div class="mb-6">
        <h1 class="text-xl font-semibold text-gray-900">User Management</h1>
      </div>
      <div class="mb-5">
        <el-button :icon="Plus" @click="openCreate">Create User</el-button>
      </div>
      <UserTable :users="users" @edit="openEdit" @delete="handleDelete" @refresh="fetchUsers" />
      <UserForm
        v-model="formVisible"
        :user="editingUser"
        @done="onFormDone"
      />
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listUsers, deleteUser } from '@/api/admin'
import type { User } from '@/types'
import AppHeader from '@/components/layout/AppHeader.vue'
import UserTable from '@/components/admin/UserTable.vue'
import UserForm from '@/components/admin/UserForm.vue'

const users = ref<User[]>([])
const formVisible = ref(false)
const editingUser = ref<User | null>(null)

onMounted(fetchUsers)

async function fetchUsers() {
  const res = await listUsers()
  users.value = res.data
}

function openCreate() {
  editingUser.value = null
  formVisible.value = true
}

function openEdit(user: User) {
  editingUser.value = user
  formVisible.value = true
}

async function handleDelete(user: User) {
  await ElMessageBox.confirm(`Delete user "${user.username}"?`, 'Confirm', { type: 'warning' })
  try {
    await deleteUser(user.id)
    ElMessage.success('User deleted')
    fetchUsers()
  } catch {
    ElMessage.error('Failed to delete user')
  }
}

function onFormDone() {
  fetchUsers()
}
</script>
