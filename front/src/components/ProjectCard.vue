<script setup lang="ts">
import { computed } from 'vue'
import type { Project } from '@/types'
import { Calendar, User, TrendingUp } from 'lucide-vue-next'
import Progress from '@/components/ui/Progress.vue'
import ImageWithFallback from '@/components/ui/ImageWithFallback.vue'

interface Props {
  project: Project
}

const props = defineProps<Props>()
const emit = defineEmits<{
  click: []
}>()

const progress = computed(() => {
  if (props.project.goal_amount === 0) return 0
  return Math.min((props.project.current_amount / props.project.goal_amount) * 100, 100)
})

const daysLeft = computed(() => {
  if (!props.project.end_date) return null
  const end = new Date(props.project.end_date)
  const now = new Date()
  const diff = Math.ceil((end.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))
  return Math.max(diff, 0)
})

function formatAmount(amount: number) {
  return amount.toLocaleString('ru-RU')
}

function getCategoryName(project: Project): string {
  if (project.categories && project.categories.length > 0) {
    return project.categories.map(c => c.title).join(', ')
  }
  return 'Без категории'
}
</script>

<template>
  <div
    class="group cursor-pointer bg-white rounded-xl overflow-hidden shadow-sm hover:shadow-xl transition-all duration-300 border border-neutral-200"
    @click="emit('click')"
  >
    <div class="relative h-48 overflow-hidden bg-neutral-100">
      <ImageWithFallback
        :src="project.project_img"
        :alt="project.title"
        class="w-full h-full object-contain"
      />
      <div class="absolute top-4 left-4 px-3 py-1 bg-white/90 backdrop-blur-sm rounded-full text-xs font-medium">
        {{ getCategoryName(project) }}
      </div>
      <div class="absolute top-4 right-4 px-3 py-1 bg-green-500 text-white rounded-full flex items-center gap-1 text-xs">
        <TrendingUp class="w-4 h-4" />
        {{ Math.floor(progress) }}%
      </div>
    </div>

    <div class="p-5 sm:p-6">
      <h3 class="mb-2 group-hover:text-blue-600 transition-colors line-clamp-1 font-semibold text-lg">
        {{ project.title }}
      </h3>
      <p class="text-neutral-600 mb-4 line-clamp-2 text-sm">
        {{ project.description }}
      </p>

      <div class="space-y-3 mb-4">
        <div>
          <div class="flex justify-between mb-2 text-sm">
            <span class="font-medium text-neutral-900">
              {{ formatAmount(project.current_amount) }} ₽
            </span>
            <span class="text-neutral-500">
              из {{ formatAmount(project.goal_amount) }} ₽
            </span>
          </div>
          <Progress :model-value="progress" class="h-2" />
        </div>

        <div class="flex items-center justify-between text-neutral-600 text-sm">
          <div class="flex items-center gap-2">
            <Calendar class="w-4 h-4" />
            <span v-if="daysLeft === null">Без срока</span>
            <span v-else-if="daysLeft > 0">{{ daysLeft }} дней</span>
            <span v-else>Завершён</span>
          </div>
          <div class="flex items-center gap-2">
            <User class="w-4 h-4" />
            <span>{{ project.user?.fio || 'Аноним' }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
