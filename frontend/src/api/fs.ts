import client from './client'
import type { MountPoint, FileListResponse } from '@/types'

export function getMounts() {
  return client.get<MountPoint[]>('/fs/mounts')
}

export function listFiles(mount: string, path: string) {
  return client.get<FileListResponse>('/fs/list', { params: { mount, path } })
}

export function downloadFile(mount: string, path: string) {
  return client.get('/fs/file', {
    params: { mount, path },
    responseType: 'blob',
  })
}

export function uploadFile(
  mount: string,
  path: string,
  file: File,
  occ?: { mtime: string; size: number }
) {
  const form = new FormData()
  form.append('file', file)
  const headers: Record<string, string> = { 'Content-Type': 'multipart/form-data' }
  if (occ) {
    headers['X-Expected-MTime'] = occ.mtime
    headers['X-Expected-Size'] = String(occ.size)
  }
  return client.post('/fs/upload', form, {
    params: { mount, path },
    headers,
  })
}

export function deleteFile(mount: string, path: string) {
  return client.delete('/fs/delete', { params: { mount, path } })
}

export function mkdir(mount: string, path: string) {
  return client.post('/fs/mkdir', null, { params: { mount, path } })
}
