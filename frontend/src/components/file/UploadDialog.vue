<template>
  <el-dialog v-model="visible" title="Upload File" width="400px" @close="handleClose">
    <el-upload
      ref="uploadRef"
      drag
      :auto-upload="false"
      :limit="1"
      :on-change="handleChange"
      :on-exceed="handleExceed"
    >
      <el-icon style="font-size: 48px; color: #909399"><Upload /></el-icon>
      <div>Drop file here or <em>click to select</em></div>
    </el-upload>
    <template #footer>
      <el-button @click="visible = false">Cancel</el-button>
      <el-button type="primary" :loading="uploading" @click="handleUpload">Upload</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Upload } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { UploadInstance, UploadFile, UploadRawFile } from 'element-plus'
import { uploadFile } from '@/api/fs'
import { useFsStore } from '@/stores/fs'

const props = defineProps<{ modelValue: boolean }>()
const emit = defineEmits<{ 'update:modelValue': [v: boolean]; done: [] }>()

const fs = useFsStore()
const uploadRef = ref<UploadInstance>()
const uploading = ref(false)
const selectedFile = ref<File | null>(null)

const visible = computed({
  get: () => props.modelValue,
  set: (v) => emit('update:modelValue', v),
})

function handleChange(uploadFile: UploadFile) {
  selectedFile.value = uploadFile.raw || null
}

function handleExceed(files: File[]) {
  uploadRef.value?.clearFiles()
  const file = files[0] as UploadRawFile
  file.uid = Date.now()
  uploadRef.value?.handleStart(file)
  selectedFile.value = file
}

async function handleUpload() {
  if (!selectedFile.value) {
    ElMessage.warning('Please select a file')
    return
  }

  uploading.value = true
  try {
    // OCC: find existing file with same name to pass mtime/size
    const existing = fs.files.find((f) => f.name === selectedFile.value!.name && !f.is_dir)
    const occ = existing
      ? { mtime: String(existing.modified), size: existing.size }
      : undefined
    await uploadFile(fs.currentMount, fs.currentPath, selectedFile.value, occ)
    ElMessage.success('Upload successful')
    visible.value = false
    emit('done')
  } catch (err: any) {
    if (err.response?.status === 409) {
      ElMessage.error('File has been modified externally, please refresh and try again')
    } else {
      ElMessage.error('Upload failed')
    }
  } finally {
    uploading.value = false
  }
}

function handleClose() {
  uploadRef.value?.clearFiles()
  selectedFile.value = null
}
</script>
