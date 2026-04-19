<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import Button from '@/components/ui/Button.vue'
import Textarea from '@/components/ui/Textarea.vue'
import Progress from '@/components/ui/Progress.vue'
import ImageWithFallback from '@/components/ui/ImageWithFallback.vue'
import {
  ArrowLeft, Calendar, User, Target, TrendingUp, Edit,
  MessageCircle, Clock, DollarSign, Github, Linkedin, Send
} from 'lucide-vue-next'
import * as projectApi from '@/api/projects'
import type { Project, Comment as CommentType } from '@/types'

const props = defineProps<{
  id: string
}>()

const router = useRouter()
const auth = useAuthStore()

const project = ref<Project | null>(null)
const comments = ref<CommentType[]>([])
const commentText = ref('')
const donationAmount = ref('')
const isLoading = ref(true)
const submitError = ref('')
const isLiked = ref(false)

const projectId = computed(() => parseInt(props.id))

const progress = computed(() => {
  if (!project.value || project.value.goal_amount === 0) return 0
  return Math.min((project.value.current_amount / project.value.goal_amount) * 100, 100)
})

const daysLeft = computed(() => {
  if (!project.value?.end_date) return 0
  const end = new Date(project.value.end_date)
  const now = new Date()
  return Math.max(Math.ceil((end.getTime() - now.getTime()) / (1000 * 60 * 60 * 24)), 0)
})

const canEdit = computed(() => {
  if (!auth.user || !project.value) return false
  return auth.user.role === 'admin' || auth.user.id === project.value.user_id
})

async function loadData() {
  isLoading.value = true
  try {
    const [projectData, commentsData] = await Promise.all([
      projectApi.getProject(projectId.value),
      projectApi.getComments(projectId.value),
    ])
    project.value = projectData
    comments.value = commentsData
    isLiked.value = projectData.is_liked || false
  } catch (err) {
    console.error('Failed to load project:', err)
  } finally {
    isLoading.value = false
  }
}

async function handleAddComment() {
  if (!commentText.value.trim() || !auth.isAuthenticated) return
  submitError.value = ''
  try {
    const newComment = await projectApi.createComment(projectId.value, commentText.value)
    comments.value.unshift(newComment)
    commentText.value = ''
  } catch (err: any) {
    submitError.value = err.message || 'Ошибка отправки комментария'
  }
}

const showSuccess = ref(false)
const successMessage = ref('')

async function handleDonate() {
  if (!donationAmount.value) return
  submitError.value = ''
  try {
    const amount = parseFloat(donationAmount.value)
    await projectApi.createPledge(projectId.value, amount)
    // Reload project to get updated amount
    project.value = await projectApi.getProject(projectId.value)
    donationAmount.value = ''

    // In-app toast notification
    successMessage.value = `Спасибо за поддержку! Вы поддержали проект на ${amount.toLocaleString('ru-RU')} ₽`
    showSuccess.value = true
    setTimeout(() => { showSuccess.value = false }, 4000)
  } catch (err: any) {
    submitError.value = err.message || 'Ошибка поддержки проекта'
  }
}

async function handleLike() {
  if (!auth.isAuthenticated) {
    router.push({ name: 'auth' })
    return
  }
  try {
    if (isLiked.value) {
      await projectApi.unlikeProject(projectId.value)
      isLiked.value = false
    } else {
      await projectApi.likeProject(projectId.value)
      isLiked.value = true
    }
    project.value = await projectApi.getProject(projectId.value)
  } catch (err) {
    console.error('Like error:', err)
  }
}

function handleEdit() {
  router.push({ name: 'edit-project', params: { id: String(projectId.value) } })
}

function goBack() {
  router.push({ name: 'home' })
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString('ru-RU', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

function formatAmount(amount: number) {
  return amount.toLocaleString('ru-RU')
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <!-- Toast notification -->
  <Transition name="toast">
    <div
      v-if="showSuccess"
      class="fixed top-20 right-4 sm:right-6 z-[9999] max-w-sm shadow-lg rounded-xl border border-green-200 bg-white px-4 py-3 flex items-start gap-3"
    >
      <div class="flex-shrink-0 w-8 h-8 rounded-full bg-green-100 flex items-center justify-center">
        <svg class="w-5 h-5 text-green-600" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
        </svg>
      </div>
      <div class="flex-1 min-w-0">
        <div class="text-sm font-medium text-neutral-900">IT Crowdfunding</div>
        <div class="text-sm text-neutral-600">{{ successMessage }}</div>
      </div>
      <button @click="showSuccess = false" class="flex-shrink-0 text-neutral-400 hover:text-neutral-600">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>
  </Transition>

  <div v-if="isLoading" class="min-h-screen flex items-center justify-center">
    <div class="text-neutral-400">Загрузка проекта...</div>
  </div>

  <div v-else-if="project" class="min-h-screen py-6 sm:py-8 px-4 sm:px-6">
    <div class="max-w-7xl mx-auto">
      <Button variant="ghost" @click="goBack" class="mb-4 sm:mb-6">
        <ArrowLeft class="w-4 h-4 mr-2" />
        Назад к проектам
      </Button>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 sm:gap-8">
        <!-- Основной контент -->
        <div class="lg:col-span-2 space-y-6 sm:space-y-8">
          <!-- Изображение -->
          <div class="bg-white rounded-xl shadow-lg overflow-hidden border border-neutral-200">
            <div class="relative h-64 sm:h-80 lg:h-96">
              <ImageWithFallback
                :src="project.project_img"
                :alt="project.title"
                class="w-full h-full object-cover"
              />
              <div class="absolute top-4 left-4 px-3 sm:px-4 py-1.5 sm:py-2 bg-white/90 backdrop-blur-sm rounded-full text-xs sm:text-sm">
                {{ project.categories?.map(c => c.title).join(', ') || 'Без категории' }}
              </div>
              <Button
                v-if="canEdit"
                variant="outline"
                size="sm"
                @click="handleEdit"
                class="absolute top-4 right-4"
              >
                <Edit class="w-4 h-4 mr-1" />
                Редактировать
              </Button>
              <button
                v-else
                @click="handleLike"
                class="absolute top-4 right-4 p-2 sm:p-3 bg-white/90 backdrop-blur-sm rounded-full hover:bg-white transition-colors"
              >
                <TrendingUp
                  class="w-5 h-5 sm:w-6 sm:h-6"
                  :class="isLiked ? 'text-red-500 fill-red-500' : 'text-neutral-600'"
                />
              </button>
            </div>

            <div class="p-6 sm:p-8">
              <h1 class="mb-3 sm:mb-4">{{ project.title }}</h1>
              <div class="flex items-center gap-3 sm:gap-4 text-neutral-600 mb-4 sm:mb-6 text-sm">
                <div class="flex items-center gap-2">
                  <User class="w-4 h-4 sm:w-5 sm:h-5" />
                  <span>{{ project.user?.fio || 'Аноним' }}</span>
                </div>
                <div class="flex items-center gap-2">
                  <Calendar class="w-4 h-4 sm:w-5 sm:h-5" />
                  <span>{{ formatDate(project.created_at) }}</span>
                </div>
              </div>

              <div class="prose max-w-none">
                <h3 class="mb-3 sm:mb-4">О проекте</h3>
                <p class="text-neutral-700 leading-relaxed whitespace-pre-line">
                  {{ project.description }}
                </p>
              </div>

              <!-- Социальные сети -->
              <div
                v-if="project.link_telegram || project.link_github || project.link_linkedin"
                class="mt-6 sm:mt-8 pt-6 sm:pt-8 border-t"
              >
                <h3 class="mb-3 sm:mb-4">Ссылки</h3>
                <div class="flex flex-wrap gap-2 sm:gap-3">
                  <a
                    v-if="project.link_telegram"
                    :href="project.link_telegram"
                    target="_blank"
                    rel="noopener noreferrer"
                    class="inline-flex items-center gap-2 px-4 py-2 bg-neutral-100 hover:bg-neutral-200 text-neutral-800 rounded-lg transition-colors text-sm font-medium"
                  >
                    <Send class="w-4 h-4" />
                    Telegram
                  </a>
                  <a
                    v-if="project.link_github"
                    :href="project.link_github"
                    target="_blank"
                    rel="noopener noreferrer"
                    class="inline-flex items-center gap-2 px-4 py-2 bg-neutral-100 hover:bg-neutral-200 text-neutral-800 rounded-lg transition-colors text-sm font-medium"
                  >
                    <Github class="w-4 h-4" />
                    GitHub
                  </a>
                  <a
                    v-if="project.link_linkedin"
                    :href="project.link_linkedin"
                    target="_blank"
                    rel="noopener noreferrer"
                    class="inline-flex items-center gap-2 px-4 py-2 bg-neutral-100 hover:bg-neutral-200 text-neutral-800 rounded-lg transition-colors text-sm font-medium"
                  >
                    <Linkedin class="w-4 h-4" />
                    LinkedIn
                  </a>
                </div>
              </div>

              <!-- Место для отслеживания -->
              <div class="mt-6 sm:mt-8 p-4 sm:p-6 bg-blue-50 border border-blue-200 rounded-xl">
                <h3 class="mb-2 sm:mb-3 flex items-center gap-2 text-blue-600">
                  <TrendingUp class="w-4 h-4 sm:w-5 sm:h-5" />
                  Отслеживание прогресса
                </h3>
                <p class="text-neutral-600 mb-3 sm:mb-4 text-sm">
                  Здесь будет размещена детальная информация о прогрессе проекта, обновления от создателей и ключевые вехи разработки.
                </p>
                <div class="space-y-1 sm:space-y-2 text-neutral-600 text-sm">
                  <div class="flex items-center gap-2">
                    <Clock class="w-4 h-4" />
                    <span>Последнее обновление: {{ new Date().toLocaleDateString('ru-RU') }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Комментарии -->
          <div class="bg-white rounded-xl shadow-lg p-6 sm:p-8 border border-neutral-200">
            <h2 class="mb-4 sm:mb-6 flex items-center gap-2 text-xl sm:text-2xl font-bold">
              <MessageCircle class="w-5 h-5 sm:w-6 sm:h-6" />
              Комментарии ({{ comments.length }})
            </h2>

            <!-- Форма добавления комментария -->
            <div v-if="auth.isAuthenticated">
              <form @submit.prevent="handleAddComment" class="mb-6 sm:mb-8">
                <Textarea
                  v-model="commentText"
                  placeholder="Оставьте свой комментарий..."
                  :rows="4"
                  class="mb-3"
                />
                <Button type="submit" :disabled="!commentText.trim()">
                  Отправить комментарий
                </Button>
              </form>
            </div>
            <div v-else class="mb-6 sm:mb-8 p-4 bg-neutral-50 rounded-lg text-center text-neutral-600">
              Войдите, чтобы оставить комментарий
            </div>

            <!-- Список комментариев -->
            <div class="space-y-4 sm:space-y-6">
              <div
                v-for="comment in comments"
                :key="comment.id"
                class="border-b border-neutral-200 pb-4 sm:pb-6 last:border-0"
              >
                <div class="flex items-center gap-3 mb-2">
                  <div class="w-8 h-8 sm:w-10 sm:h-10 bg-gradient-to-br from-blue-500 to-purple-500 rounded-full flex items-center justify-center text-white text-sm">
                    {{ comment.user?.fio?.charAt(0)?.toUpperCase() || '?' }}
                  </div>
                  <div>
                    <div class="font-medium text-sm sm:text-base">{{ comment.user?.fio || 'Аноним' }}</div>
                    <div class="text-neutral-500 text-xs sm:text-sm">
                      {{ formatDate(comment.created_at) }}
                    </div>
                  </div>
                </div>
                <p class="text-neutral-700 text-sm sm:text-base ml-10 sm:ml-13">{{ comment.content }}</p>
              </div>
              <div v-if="comments.length === 0" class="text-center py-6 sm:py-8 text-neutral-500">
                Пока нет комментариев. Будьте первым!
              </div>
            </div>
          </div>
        </div>

        <!-- Боковая панель -->
        <div class="space-y-4 sm:space-y-6">
          <!-- Финансирование -->
          <div class="bg-white rounded-xl shadow-lg p-5 sm:p-6 border border-neutral-200 lg:sticky lg:top-24">
            <div class="mb-4 sm:mb-6">
              <div class="text-xl sm:text-2xl font-bold mb-1 sm:mb-2">
                {{ formatAmount(project.current_amount) }} ₽
              </div>
              <div class="text-neutral-600 mb-3 sm:mb-4 text-sm">
                собрано из {{ formatAmount(project.goal_amount) }} ₽
              </div>
              <Progress :model-value="progress" class="h-2 sm:h-3 mb-2" />
              <div class="text-green-600 font-medium text-sm">{{ Math.floor(progress) }}% достигнуто</div>
            </div>

            <div class="grid grid-cols-2 gap-3 sm:gap-4 mb-4 sm:mb-6 pb-4 sm:pb-6 border-b">
              <div>
                <div class="text-neutral-500 mb-1 text-xs sm:text-sm">Спонсоров</div>
                <div class="font-semibold">{{ project.likes_count || 0 }}</div>
              </div>
              <div>
                <div class="text-neutral-500 mb-1 text-xs sm:text-sm">Осталось дней</div>
                <div class="font-semibold">{{ daysLeft || '—' }}</div>
              </div>
            </div>

            <div v-if="submitError" class="text-red-500 text-sm text-center mb-3">
              {{ submitError }}
            </div>

            <template v-if="auth.isAuthenticated">
              <div class="space-y-2 sm:space-y-3">
                <div class="flex gap-2">
                  <input
                    type="number"
                    placeholder="Сумма"
                    v-model="donationAmount"
                    class="flex-1 px-3 py-2 border border-neutral-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                    min="1"
                    step="0.01"
                  />
                  <span class="flex items-center text-neutral-600">₽</span>
                </div>
                <Button @click="handleDonate" class="w-full" :disabled="!donationAmount">
                  <DollarSign class="w-4 h-4 mr-2" />
                  Поддержать проект
                </Button>
              </div>
            </template>
            <div v-else class="text-center p-3 sm:p-4 bg-neutral-50 rounded-lg text-neutral-600 text-sm">
              Войдите, чтобы поддержать проект
            </div>
          </div>

          <!-- Информация о проекте -->
          <div class="bg-white rounded-xl shadow-lg p-5 sm:p-6 border border-neutral-200">
            <h3 class="mb-3 sm:mb-4 font-semibold text-lg">Информация</h3>
            <div class="space-y-2 sm:space-y-3 text-neutral-600 text-sm">
              <div class="flex items-center gap-2 sm:gap-3">
                <Target class="w-4 h-4 sm:w-5 sm:h-5 flex-shrink-0" />
                <div>
                  <div class="text-neutral-500 text-xs">Цель</div>
                  <div class="font-medium">{{ formatAmount(project.goal_amount) }} ₽</div>
                </div>
              </div>
              <div class="flex items-center gap-2 sm:gap-3">
                <User class="w-4 h-4 sm:w-5 sm:h-5 flex-shrink-0" />
                <div>
                  <div class="text-neutral-500 text-xs">Создатель</div>
                  <div class="font-medium">{{ project.user?.fio || 'Аноним' }}</div>
                </div>
              </div>
              <div class="flex items-center gap-2 sm:gap-3">
                <Calendar class="w-4 h-4 sm:w-5 sm:h-5 flex-shrink-0" />
                <div>
                  <div class="text-neutral-500 text-xs">Создан</div>
                  <div class="font-medium">{{ formatDate(project.created_at) }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <div v-else class="min-h-screen flex items-center justify-center">
    <div class="text-neutral-400">Проект не найден</div>
  </div>
</template>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}
.toast-enter-from {
  opacity: 0;
  transform: translateX(40px);
}
.toast-leave-to {
  opacity: 0;
  transform: translateX(40px);
}
</style>
