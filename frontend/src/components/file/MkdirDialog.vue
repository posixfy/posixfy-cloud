<template>
  <el-dialog v-model="visible" title="Create New Folder" width="400px" @close="name = ''">
    <el-input v-model="name" placeholder="Folder name" @keyup.enter="handleCreate" />
    <template #footer>
      <el-button @click="visible = false">Cancel</el-button>
      <el-button type="primary" :loading="creating" @click="handleCreate">Create</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { mkdir } from '@/api/fs'
import { useFsStore } from '@/stores/fs'

const props = defineProps<{ modelValue: boolean }>()
const emit = defineEmits<{ 'update:modelValue': [v: boolean]; done: [] }>()

const fs = useFsStore()
const name = ref('')
const creating = ref(false)

const visible = computed({
  get: () => props.modelValue,
  set: (v) => emit('update:modelValue', v),
})

async function handleCreate() {
  if (!name.value.trim()) {
    ElMessage.warning('Please enter a folder name')
    return
  }

  creating.value = true
  const dirPath = fs.currentPath.endsWith('/')
    ? fs.currentPath + name.value
    : fs.currentPath + '/' + name.value
  try {
    await mkdir(fs.currentMount, dirPath)
    ElMessage.success('Folder created')
    name.value = ''
    visible.value = false
    emit('done')
  } catch {
    ElMessage.error('Failed to create folder')
  } finally {
    creating.value = false
  }
}
</script>
