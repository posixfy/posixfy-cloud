export interface User {
  id: number
  username: string
  uid: number
  gid: number
  groups: string
  role: 'admin' | 'user'
  created_at: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
}

export interface CreateUserRequest {
  username: string
  password: string
  uid: number
  gid: number
  groups: string
  role: string
}

export interface UpdateUserRequest {
  password?: string
  gid?: number
  groups?: string
  role?: string
}

export interface MountPoint {
  name: string
  path: string
}

export interface FileEntry {
  name: string
  is_dir: boolean
  size: number
  modified: string
  permissions: string
}

export interface FileListResponse {
  entries: FileEntry[]
  total: number
  page: number
  limit: number
  has_more: boolean
}
