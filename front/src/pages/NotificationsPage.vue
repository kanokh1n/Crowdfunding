<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Bell, CheckCheck, ChevronRight } from 'lucide-vue-next'
import * as notifApi from '@/api/notifications'
import type { Notification } from '@/types'

const router = useRouter()
const notifications = ref<Notification[]>([])
const isLoading = ref(true)

async function load() {
  try {
    notifications.value = await notifApi.listNotifications()
  } finally {
    isLoading.value = false
  }
}

async function handleRead(notif: Notification) {
  if (!notif.is_read) {
    await notifApi.markRead(notif.id)
    notif.is_read = true
  }
  if (notif.project_id) {
    router.push({ name: 'project-detail', params: { id: notif.project_id } })
  }
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString('ru-RU', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function typeColor(type: Notification['type']) {
  if (type === 'ai_passed') return 'bg-green-100 text-green-700'
  if (type === 'ai_failed') return 'bg-red-100 text-red-700'
  return 'bg-blue-100 text-blue-700'
}

function typeLabel(type: Notification['type']) {
  if (type === 'ai_passed') return 'AI: одобрено'
  if (type === 'ai_failed') return 'AI: отклонено'
  return 'Приглашение'
}

onMounted(load)
</script>

<template>
  <div class="min-h-screen py-8 sm:py-12 px-4 sm:px-6">
    <div class="max-w-3xl mx-auto">
      <div class="flex items-center gap-3 mb-8">
        <Bell class="w-7 h-7 text-blue-600" />
        <h1>Уведомления</h1>
      </div>

      <div v-if="isLoading" class="text-center py-20 text-neutral-400">
        Загрузка...
      </div>

      <div v-else-if="notifications.length === 0" class="bg-white rounded-xl border border-neutral-200 p-12 text-center">
        <Bell class="w-12 h-12 mx-auto mb-4 text-neutral-300" />
        <p class="text-neutral-500">Нет уведомлений</p>
      </div>

      <div v-else class="space-y-3">
        <div
          v-for="notif in notifications"
          :key="notif.id"
          class="bg-white rounded-xl border border-neutral-200 p-4 sm:p-5 flex items-start gap-4 cursor-pointer hover:border-blue-300 transition-colors"
          :class="{ 'opacity-60': notif.is_read }"
          @click="handleRead(notif)"
        >
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 mb-1">
              <span class="text-xs px-2 py-0.5 rounded-full font-medium" :class="typeColor(notif.type)">
                {{ typeLabel(notif.type) }}
              </span>
              <span v-if="!notif.is_read" class="w-2 h-2 rounded-full bg-blue-500 shrink-0" />
            </div>
            <p class="font-semibold text-neutral-900">{{ notif.title }}</p>
            <p class="text-neutral-600 text-sm mt-1">{{ notif.body }}</p>
            <p class="text-neutral-400 text-xs mt-2">{{ formatDate(notif.created_at) }}</p>
          </div>
          <ChevronRight v-if="notif.project_id" class="w-5 h-5 text-neutral-400 shrink-0 mt-1" />
        </div>
      </div>
    </div>
  </div>
</template>
