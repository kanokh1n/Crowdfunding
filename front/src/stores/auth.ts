import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import * as authApi from '@/api/auth'
import type { User } from '@/types'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const isLoading = ref(false)

  const isAuthenticated = computed(() => !!user.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  async function login(email: string, password: string) {
    isLoading.value = true
    try {
      const tokens = await authApi.login({ email, password })
      localStorage.setItem('access_token', tokens.access_token)
      localStorage.setItem('refresh_token', tokens.refresh_token)
      user.value = tokens.user
      return tokens
    } finally {
      isLoading.value = false
    }
  }

  async function register(email: string, password: string, fio: string) {
    isLoading.value = true
    try {
      await authApi.register({ email, password, fio })
      // Auto-login after registration
      return await login(email, password)
    } finally {
      isLoading.value = false
    }
  }

  async function fetchMe() {
    const token = localStorage.getItem('access_token')
    if (!token) return

    try {
      user.value = await authApi.getMe()
    } catch {
      logout()
    }
  }

  function logout() {
    user.value = null
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
  }

  return {
    user,
    isLoading,
    isAuthenticated,
    isAdmin,
    login,
    register,
    fetchMe,
    logout,
  }
})
