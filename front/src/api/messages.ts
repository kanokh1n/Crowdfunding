import { apiFetch, fetchList } from './client'
import type { Message } from '@/types'

export function listMessages(params: { page?: number; limit?: number } = {}) {
  const qs = new URLSearchParams()
  if (params.page) qs.set('page', String(params.page))
  if (params.limit) qs.set('limit', String(params.limit))
  const query = qs.toString()
  return fetchList<Message>(`/messages${query ? `?${query}` : ''}`, { requireAuth: true })
}

export function sendMessage(data: { project_id: number; title: string; content: string }) {
  return apiFetch<Message>('/messages', {
    method: 'POST',
    requireAuth: true,
    body: JSON.stringify(data),
  })
}

export function markMessageRead(messageId: number) {
  return apiFetch<Message>(`/messages/${messageId}/read`, {
    method: 'PATCH',
    requireAuth: true,
  })
}
