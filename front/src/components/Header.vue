<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import { LogIn, PlusCircle, ShieldCheck, UserPlus, Bell } from 'lucide-vue-next'
import * as notifApi from '@/api/notifications'

const auth = useAuthStore()
const router = useRouter()
const unreadCount = ref(0)
let pollTimer: ReturnType<typeof setInterval> | null = null

async function fetchUnread() {
  if (!auth.isAuthenticated) return
  try {
    const notifs = await notifApi.listNotifications()
    unreadCount.value = notifs.filter(n => !n.is_read).length
  } catch {}
}

onMounted(() => {
  fetchUnread()
  pollTimer = setInterval(fetchUnread, 60_000)
})

onUnmounted(() => {
  if (pollTimer) clearInterval(pollTimer)
})

function navigateTo(name: string) {
  router.push({ name })
}

function handleLogout() {
  auth.logout()
  router.push({ name: 'home' })
}

function getInitials(name: string) {
  return name.charAt(0).toUpperCase()
}
</script>

<template>
  <header class="sticky top-0 z-50 bg-white border-b border-neutral-200 shadow-sm">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 py-3 sm:py-4">
      <div class="flex items-center justify-between">
        <!-- Logo -->
        <div
          class="flex items-center gap-2 cursor-pointer"
          @click="navigateTo('home')"
        >
          <div class="w-10 h-10 bg-gradient-to-br from-blue-600 to-purple-600 rounded-lg flex items-center justify-center">
            <span class="text-white font-bold">IT</span>
          </div>
          <span class="font-semibold hidden sm:inline">IT Crowdfunding</span>
        </div>

        <!-- Nav -->
        <nav class="hidden md:flex items-center gap-6">
          <button
            class="text-neutral-700 hover:text-neutral-900 transition-colors"
            @click="navigateTo('home')"
          >
            Проекты
          </button>
          <span class="text-neutral-500">О платформе</span>
          <button
            class="text-neutral-700 hover:text-neutral-900 transition-colors"
            @click="navigateTo('how-it-works')"
          >
            Как это работает
          </button>
        </nav>

        <!-- Actions -->
        <div class="flex items-center gap-2 sm:gap-3">
          <template v-if="auth.isAuthenticated">
            <Button
              v-if="auth.isAdmin"
              variant="outline"
              size="sm"
              @click="navigateTo('moderation')"
            >
              <ShieldCheck class="w-4 h-4 mr-1 sm:mr-2" />
              <span class="hidden sm:inline">Модерация</span>
            </Button>
            <Button
              variant="outline"
              size="sm"
              @click="navigateTo('create-project')"
            >
              <PlusCircle class="w-4 h-4 mr-1 sm:mr-2" />
              <span class="hidden sm:inline">Создать проект</span>
            </Button>
            <!-- Bell -->
            <button
              class="relative p-1.5 rounded-lg hover:bg-neutral-100 transition-colors"
              @click="navigateTo('notifications')"
            >
              <Bell class="w-5 h-5 text-neutral-600" />
              <span
                v-if="unreadCount > 0"
                class="absolute -top-1 -right-1 w-4 h-4 bg-red-500 text-white text-[10px] font-bold rounded-full flex items-center justify-center"
              >{{ unreadCount > 9 ? '9+' : unreadCount }}</span>
            </button>
            <div class="flex items-center gap-2">
              <div class="w-8 h-8 bg-gradient-to-br from-blue-500 to-purple-500 rounded-full flex items-center justify-center text-white text-sm">
                {{ auth.user ? getInitials(auth.user.fio) : 'U' }}
              </div>
              <span class="hidden lg:inline text-sm">{{ auth.user?.fio }}</span>
            </div>
            <Button variant="ghost" size="sm" @click="handleLogout">
              Выйти
            </Button>
          </template>
          <template v-else>
            <Button variant="outline" size="sm" @click="navigateTo('auth')">
              <LogIn class="w-4 h-4 mr-1 sm:mr-2" />
              <span class="hidden sm:inline">Войти</span>
            </Button>
            <Button size="sm" @click="navigateTo('auth')">
              <UserPlus class="w-4 h-4 mr-1 sm:mr-2" />
              <span class="hidden sm:inline">Регистрация</span>
            </Button>
          </template>
        </div>
      </div>
    </div>
  </header>
</template>
