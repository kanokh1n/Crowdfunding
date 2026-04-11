import { apiFetch, fetchList } from './client'
import type { AuthTokens, RegisterInput, LoginInput, User, Project } from '@/types'

export function register(data: RegisterInput) {
  return apiFetch<void>('/auth/register', {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

export function login(data: LoginInput) {
  return apiFetch<AuthTokens>('/auth/login', {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

export function verifyEmail(token: string) {
  return apiFetch<void>(`/auth/verify-email?token=${encodeURIComponent(token)}`, {
    method: 'POST',
  })
}

export function refreshToken() {
  return apiFetch<AuthTokens>('/auth/refresh', {
    method: 'POST',
  })
}

export function getMe() {
  return apiFetch<User>('/users/me', { requireAuth: true })
}

export function updateMe(data: Partial<Pick<User, 'fio' | 'description' | 'phone' | 'profile_img' | 'email'>>) {
  return apiFetch<User>('/users/me', {
    method: 'PATCH',
    requireAuth: true,
    body: JSON.stringify(data),
  })
}

export function getUserById(id: number) {
  return apiFetch<User>(`/users/${id}`, { requireAuth: true })
}

export function getMyProjects(params: { page?: number; limit?: number } = {}) {
  const qs = new URLSearchParams()
  if (params.page) qs.set('page', String(params.page))
  if (params.limit) qs.set('limit', String(params.limit))
  const query = qs.toString()
  return fetchList<Project>(`/users/me/projects${query ? `?${query}` : ''}`, { requireAuth: true })
}
