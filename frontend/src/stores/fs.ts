import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getMounts, listFiles } from '@/api/fs'
import type { MountPoint, FileEntry } from '@/types'

interface FsEvent {
  event_type: string
  name: string
  is_dir: boolean
  size: number
  modified: number
}

export const useFsStore = defineStore('fs', () => {
  const mounts = ref<MountPoint[]>([])
  const currentMount = ref('')
  const currentPath = ref('/')
  const files = ref<FileEntry[]>([])
  const loading = ref(false)

  async function fetchMounts() {
    const res = await getMounts()
    mounts.value = res.data
    if (res.data.length > 0 && !currentMount.value) {
      currentMount.value = res.data[0].name
    }
  }

  async function fetchFiles() {
    if (!currentMount.value) return
    loading.value = true
    try {
      const res = await listFiles(currentMount.value, currentPath.value)
      files.value = res.data.entries
    } finally {
      loading.value = false
    }
  }

  function navigateTo(mount: string, path: string) {
    currentMount.value = mount
    currentPath.value = path
  }

  function enterDir(name: string) {
    const p = currentPath.value.endsWith('/')
      ? currentPath.value + name
      : currentPath.value + '/' + name
    currentPath.value = p
  }

  function goUp() {
    if (currentPath.value === '/' || currentPath.value === '') return
    const parts = currentPath.value.split('/').filter(Boolean)
    parts.pop()
    currentPath.value = '/' + parts.join('/')
  }

  function applyEvent(event: FsEvent) {
    switch (event.event_type) {
      case 'created':
        // Avoid duplicates (e.g. from own upload)
        if (!files.value.find((f) => f.name === event.name)) {
          files.value.push({
            name: event.name,
            is_dir: event.is_dir,
            size: event.size,
            modified: String(event.modified),
            permissions: '',
          })
        }
        break
      case 'modified': {
        const existing = files.value.find((f) => f.name === event.name)
        if (existing) {
          existing.size = event.size
          existing.modified = String(event.modified)
        }
        break
      }
      case 'deleted':
        files.value = files.value.filter((f) => f.name !== event.name)
        break
    }
  }

  return {
    mounts,
    currentMount,
    currentPath,
    files,
    loading,
    fetchMounts,
    fetchFiles,
    navigateTo,
    enterDir,
    goUp,
    applyEvent,
  }
})
