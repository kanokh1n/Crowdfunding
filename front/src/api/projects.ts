import { apiFetch, fetchList } from './client'
import type {
  Project,
  ProjectListParams,
  CreateProjectInput,
  UpdateProjectInput,
  Comment,
  Pledge,
  Like,
} from '@/types'

// Projects
export function listProjects(params: ProjectListParams = {}) {
  const qs = new URLSearchParams()
  if (params.page) qs.set('page', String(params.page))
  if (params.limit) qs.set('limit', String(params.limit))
  if (params.category_id) qs.set('category_id', String(params.category_id))
  if (params.search) qs.set('search', params.search)
  if (params.status) qs.set('status', params.status)
  if (params.sort) qs.set('sort', params.sort)
  const query = qs.toString()
  return fetchList<Project>(`/projects${query ? `?${query}` : ''}`)
}

export function getProject(id: number) {
  return apiFetch<Project>(`/projects/${id}`)
}

export function createProject(data: CreateProjectInput) {
  return apiFetch<Project>('/projects', {
    method: 'POST',
    requireAuth: true,
    body: JSON.stringify(data),
  })
}

export function updateProject(id: number, data: UpdateProjectInput) {
  return apiFetch<Project>(`/projects/${id}`, {
    method: 'PATCH',
    requireAuth: true,
    body: JSON.stringify(data),
  })
}

export function deleteProject(id: number) {
  return apiFetch<void>(`/projects/${id}`, {
    method: 'DELETE',
    requireAuth: true,
  })
}

// Comments
export function getComments(projectId: number, params: { page?: number; limit?: number } = {}) {
  const qs = new URLSearchParams()
  if (params.page) qs.set('page', String(params.page))
  if (params.limit) qs.set('limit', String(params.limit))
  const query = qs.toString()
  return fetchList<Comment>(`/projects/${projectId}/comments${query ? `?${query}` : ''}`)
}

export function createComment(projectId: number, content: string) {
  return apiFetch<Comment>(`/projects/${projectId}/comments`, {
    method: 'POST',
    requireAuth: true,
    body: JSON.stringify({ content }),
  })
}

export function updateComment(commentId: number, content: string) {
  return apiFetch<Comment>(`/comments/${commentId}`, {
    method: 'PATCH',
    requireAuth: true,
    body: JSON.stringify({ content }),
  })
}

export function deleteComment(commentId: number) {
  return apiFetch<void>(`/comments/${commentId}`, {
    method: 'DELETE',
    requireAuth: true,
  })
}

// Pledges
export function createPledge(projectId: number, amount: number) {
  return apiFetch<Pledge>(`/projects/${projectId}/pledges`, {
    method: 'POST',
    requireAuth: true,
    body: JSON.stringify({ amount }),
  })
}

export function getPledges(projectId: number, params: { page?: number; limit?: number } = {}) {
  const qs = new URLSearchParams()
  if (params.page) qs.set('page', String(params.page))
  if (params.limit) qs.set('limit', String(params.limit))
  const query = qs.toString()
  return fetchList<Pledge>(`/projects/${projectId}/pledges${query ? `?${query}` : ''}`, {
    requireAuth: true,
  })
}

// Likes
export function likeProject(projectId: number) {
  return apiFetch<Like>(`/projects/${projectId}/like`, {
    method: 'POST',
    requireAuth: true,
  })
}

export function unlikeProject(projectId: number) {
  return apiFetch<void>(`/projects/${projectId}/like`, {
    method: 'DELETE',
    requireAuth: true,
  })
}
