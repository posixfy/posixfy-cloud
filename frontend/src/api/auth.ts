import client from './client'
import type { LoginRequest, LoginResponse, User } from '@/types'

export function login(data: LoginRequest) {
  return client.post<LoginResponse>('/auth/login', data)
}

export function getMe() {
  return client.get<User>('/auth/me')
}
