<template>
  <el-dialog v-model="visible" :title="isEdit ? 'Edit User' : 'Create User'" width="500px" @close="resetForm">
    <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
      <el-form-item label="Username" prop="username">
        <el-input v-model="form.username" :disabled="isEdit" />
      </el-form-item>
      <el-form-item label="Password" :prop="isEdit ? undefined : 'password'">
        <el-input v-model="form.password" type="password" show-password
          :placeholder="isEdit ? 'Leave empty to keep current' : ''" />
      </el-form-item>
      <el-form-item v-if="!isEdit" label="UID" prop="uid">
        <el-input-number v-model="form.uid" :min="0" />
      </el-form-item>
      <el-form-item label="GID" prop="gid">
        <el-input-number v-model="form.gid" :min="0" />
      </el-form-item>
      <el-form-item label="Groups" prop="groups">
        <el-input v-model="form.groups" placeholder='[100,101]' />
      </el-form-item>
      <el-form-item label="Role" prop="role">
        <el-select v-model="form.role">
          <el-option label="User" value="user" />
          <el-option label="Admin" value="admin" />
        </el-select>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="visible = false">Cancel</el-button>
      <el-button type="primary" :loading="saving" @click="handleSave">Save</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { createUser, updateUser } from '@/api/admin'
import type { User } from '@/types'

const props = defineProps<{ modelValue: boolean; user: User | null }>()
const emit = defineEmits<{ 'update:modelValue': [v: boolean]; done: [] }>()

const formRef = ref<FormInstance>()
const saving = ref(false)
const isEdit = computed(() => !!props.user)

const form = ref({
  username: '',
  password: '',
  uid: 1000,
  gid: 1000,
  groups: '[]',
  role: 'user',
})

const rules = {
  username: [{ required: true, message: 'Required', trigger: 'blur' }],
  password: [{ required: true, message: 'Required', trigger: 'blur' }],
  uid: [{ required: true, message: 'Required', trigger: 'blur' }],
  gid: [{ required: true, message: 'Required', trigger: 'blur' }],
}

const visible = computed({
  get: () => props.modelValue,
  set: (v) => emit('update:modelValue', v),
})

watch(() => props.user, (u) => {
  if (u) {
    form.value = {
      username: u.username,
      password: '',
      uid: u.uid,
      gid: u.gid,
      groups: u.groups,
      role: u.role,
    }
  } else {
    resetForm()
  }
})

function resetForm() {
  form.value = { username: '', password: '', uid: 1000, gid: 1000, groups: '[]', role: 'user' }
  formRef.value?.clearValidate()
}

async function handleSave() {
  if (!isEdit.value) {
    const valid = await formRef.value?.validate().catch(() => false)
    if (!valid) return
  }

  saving.value = true
  try {
    if (isEdit.value && props.user) {
      await updateUser(props.user.id, {
        password: form.value.password || undefined,
        gid: form.value.gid,
        groups: form.value.groups,
        role: form.value.role,
      })
      ElMessage.success('User updated')
    } else {
      await createUser({
        username: form.value.username,
        password: form.value.password,
        uid: form.value.uid,
        gid: form.value.gid,
        groups: form.value.groups,
        role: form.value.role,
      })
      ElMessage.success('User created')
    }
    visible.value = false
    emit('done')
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || 'Operation failed')
  } finally {
    saving.value = false
  }
}
</script>
