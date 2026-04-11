<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import ImageWithFallback from '@/components/ui/ImageWithFallback.vue'
import { Check, X, ShieldCheck, Settings } from 'lucide-vue-next'
import * as adminApi from '@/api/admin'
import type { Project } from '@/types'

const router = useRouter()
const projects = ref<Project[]>([])
const isLoading = ref(true)
const actionLoading = ref<Record<string, boolean>>({})

async function loadProjects() {
  isLoading.value = true
  try {
    projects.value = await adminApi.getModerationQueue()
  } catch (err) {
    console.error('Failed to load moderation queue:', err)
  } finally {
    isLoading.value = false
  }
}

async function handleModerate(projectId: number, decision: 'approve' | 'reject') {
  actionLoading.value[projectId] = true
  try {
    await adminApi.moderationDecision(projectId, { decision })
    // Remove from list
    projects.value = projects.value.filter(p => p.id !== projectId)
  } catch (err: any) {
    console.error('Moderation error:', err)
  } finally {
    actionLoading.value[projectId] = false
  }
}

function handleChooseMethod(projectId: number) {
  router.push({ name: 'moderation-choice', params: { id: projectId } })
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
  loadProjects()
})
</script>

<template>
  <div class="min-h-screen py-8 sm:py-12 px-4 sm:px-6">
    <div class="max-w-7xl mx-auto">
      <div class="mb-6 sm:mb-8">
        <div class="flex items-center gap-2 sm:gap-3 mb-3 sm:mb-4">
          <ShieldCheck class="w-6 h-6 sm:w-8 sm:h-8 text-blue-600" />
          <h1>Модерация проектов</h1>
        </div>
        <p class="text-neutral-600">
          Проверяйте и одобряйте проекты перед публикацией на платформе
        </p>
      </div>

      <div v-if="isLoading" class="text-center py-20">
        <div class="text-neutral-400">Загрузка проектов на модерации...</div>
      </div>

      <template v-else>
        <div v-if="projects.length > 0" class="space-y-4 sm:space-y-6">
          <div
            v-for="project in projects"
            :key="project.id"
            class="bg-white rounded-xl shadow-lg border border-neutral-200 overflow-hidden"
          >
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4 sm:gap-6">
              <!-- Изображение -->
              <div class="relative h-48 sm:h-64 md:h-auto">
                <ImageWithFallback
                  :src="project.project_img"
                  :alt="project.title"
                  class="w-full h-full object-cover"
                />
                <div class="absolute top-3 sm:top-4 left-3 sm:left-4 px-2.5 sm:px-3 py-1 bg-yellow-500 text-white rounded-full text-xs sm:text-sm">
                  На модерации
                </div>
              </div>

              <!-- Контент -->
              <div class="md:col-span-2 p-4 sm:p-6">
                <div class="mb-3 sm:mb-4">
                  <div class="inline-block px-2.5 sm:px-3 py-1 bg-neutral-100 rounded-full mb-2 sm:mb-3 text-xs sm:text-sm">
                    {{ project.categories?.map(c => c.title).join(', ') || 'Без категории' }}
                  </div>
                  <h2 class="mb-1 sm:mb-2">{{ project.title }}</h2>
                  <p class="text-neutral-600 mb-3 sm:mb-4 text-sm">{{ project.description }}</p>
                </div>

                <div class="grid grid-cols-2 gap-3 sm:gap-4 mb-4 sm:mb-6 pb-4 sm:pb-6 border-b">
                  <div>
                    <div class="text-neutral-500 mb-1 text-xs sm:text-sm">Автор</div>
                    <div class="text-sm sm:text-base">{{ project.user?.fio || 'Аноним' }}</div>
                  </div>
                  <div>
                    <div class="text-neutral-500 mb-1 text-xs sm:text-sm">Целевая сумма</div>
                    <div class="text-sm sm:text-base">{{ formatAmount(project.goal_amount) }} ₽</div>
                  </div>
                  <div>
                    <div class="text-neutral-500 mb-1 text-xs sm:text-sm">Дата создания</div>
                    <div class="text-sm sm:text-base">{{ formatDate(project.created_at) }}</div>
                  </div>
                  <div>
                    <div class="text-neutral-500 mb-1 text-xs sm:text-sm">Статус</div>
                    <div class="text-yellow-600 text-sm sm:text-base">
                      {{ project.status === 'pending_ai' ? 'AI проверка' : 'Ожидает модерации' }}
                    </div>
                  </div>
                </div>

                <div class="mb-4 sm:mb-6">
                  <h3 class="mb-1 sm:mb-2 font-semibold">Полное описание</h3>
                  <p class="text-neutral-600 leading-relaxed text-sm line-clamp-3">
                    {{ project.description }}
                  </p>
                </div>

                <!-- Кнопки модерации -->
                <div class="flex flex-col sm:flex-row gap-2 sm:gap-3">
                  <Button
                    @click="handleChooseMethod(project.id)"
                    class="bg-blue-600 hover:bg-blue-700 text-sm"
                    :disabled="actionLoading[project.id]"
                  >
                    <Settings class="w-4 h-4 mr-1 sm:mr-2" />
                    Выбрать метод проверки
                  </Button>
                  <Button
                    @click="handleModerate(project.id, 'approve')"
                    class="bg-green-600 hover:bg-green-700 text-sm"
                    :disabled="actionLoading[project.id]"
                  >
                    <Check class="w-4 h-4 mr-1 sm:mr-2" />
                    Одобрить
                  </Button>
                  <Button
                    @click="handleModerate(project.id, 'reject')"
                    variant="destructive"
                    class="text-sm"
                    :disabled="actionLoading[project.id]"
                  >
                    <X class="w-4 h-4 mr-1 sm:mr-2" />
                    Отклонить
                  </Button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div v-else class="bg-white rounded-xl shadow-lg p-8 sm:p-12 text-center border border-neutral-200">
          <ShieldCheck class="w-12 h-12 sm:w-16 sm:h-16 mx-auto mb-3 sm:mb-4 text-neutral-400" />
          <h2 class="mb-1 sm:mb-2">Нет проектов на модерации</h2>
          <p class="text-neutral-600">
            Все проекты проверены. Новые проекты появятся здесь автоматически.
          </p>
        </div>
      </template>
    </div>
  </div>
</template>
