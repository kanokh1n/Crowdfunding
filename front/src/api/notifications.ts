import { apiFetch } from './client'
import type { Notification } from '@/types'

export function listNotifications() {
  return apiFetch<Notification[]>('/notifications', { requireAuth: true })
}

export function markRead(id: number) {
  return apiFetch<Notification>(`/notifications/${id}/read`, {
    method: 'PATCH',
    requireAuth: true,
  })
}

export function sendInvite(projectId: number, message: string) {
  return apiFetch<{ ok: boolean }>(`/admin/moderation/${projectId}/invite`, {
    method: 'POST',
    requireAuth: true,
    body: JSON.stringify({ message }),
  })
}
