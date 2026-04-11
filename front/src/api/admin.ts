import { apiFetch, fetchList } from './client'
import type { Category, Project, ProjectModeration, User } from '@/types'

// Categories
export function listCategories() {
  return apiFetch<Category[]>('/categories')
}

export function createCategory(title: string) {
  return apiFetch<Category>('/admin/categories', {
    method: 'POST',
    requireAuth: true,
    body: JSON.stringify({ title }),
  })
}

export function deleteCategory(id: number) {
  return apiFetch<void>(`/admin/categories/${id}`, {
    method: 'DELETE',
    requireAuth: true,
  })
}

// Moderation
export function getModerationQueue() {
  return fetchList<Project>('/admin/moderation', { requireAuth: true })
}

export function getModerationProject(projectId: number) {
  return apiFetch<Project & { moderation: ProjectModeration }>(
    `/admin/moderation/${projectId}`,
    { requireAuth: true }
  )
}

export function moderationDecision(
  projectId: number,
  data: { decision: 'approve' | 'reject'; moderator_note?: string }
) {
  return apiFetch<Project>(`/admin/moderation/${projectId}`, {
    method: 'PATCH',
    requireAuth: true,
    body: JSON.stringify(data),
  })
}

export function recheckProject(projectId: number) {
  return apiFetch<Project>(`/admin/moderation/${projectId}/recheck`, {
    method: 'POST',
    requireAuth: true,
  })
}

// Admin projects
export function adminListProjects(params: { page?: number; limit?: number; status?: string } = {}) {
  const qs = new URLSearchParams()
  if (params.page) qs.set('page', String(params.page))
  if (params.limit) qs.set('limit', String(params.limit))
  if (params.status) qs.set('status', params.status)
  const query = qs.toString()
  return fetchList<Project>(`/admin/projects${query ? `?${query}` : ''}`, { requireAuth: true })
}

export function adminUpdateProject(id: number, data: { status?: string }) {
  return apiFetch<Project>(`/admin/projects/${id}`, {
    method: 'PATCH',
    requireAuth: true,
    body: JSON.stringify(data),
  })
}

export function adminDeleteProject(id: number) {
  return apiFetch<void>(`/admin/projects/${id}`, {
    method: 'DELETE',
    requireAuth: true,
  })
}

// Admin users
export function adminListUsers(params: { page?: number; limit?: number } = {}) {
  const qs = new URLSearchParams()
  if (params.page) qs.set('page', String(params.page))
  if (params.limit) qs.set('limit', String(params.limit))
  const query = qs.toString()
  return fetchList<any>(`/admin/users${query ? `?${query}` : ''}`, { requireAuth: true })
}

export function adminUpdateUser(id: number, data: { role?: string; is_verified?: boolean }) {
  return apiFetch<User>(`/admin/users/${id}`, {
    method: 'PATCH',
    requireAuth: true,
    body: JSON.stringify(data),
  })
}
