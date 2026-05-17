import client from './client'
import type { User, CreateUserRequest, UpdateUserRequest } from '@/types'

export function listUsers() {
  return client.get<User[]>('/admin/users')
}

export function getUser(id: number) {
  return client.get<User>(`/admin/users/${id}`)
}

export function createUser(data: CreateUserRequest) {
  return client.post<User>('/admin/users', data)
}

export function updateUser(id: number, data: UpdateUserRequest) {
  return client.put<User>(`/admin/users/${id}`, data)
}

export function deleteUser(id: number) {
  return client.delete(`/admin/users/${id}`)
}
