<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Search } from 'lucide-vue-next'
import Input from '@/components/ui/Input.vue'
import ProjectCard from '@/components/ProjectCard.vue'
import * as projectApi from '@/api/projects'
import * as adminApi from '@/api/admin'
import type { Project, Category } from '@/types'

const router = useRouter()
const projects = ref<Project[]>([])
const categories = ref<Category[]>([])
const searchQuery = ref('')
const activeCategoryId = ref<number | null>(null)
const isLoading = ref(true)

const allCategories = computed(() => [
  { id: null, title: 'Все проекты' },
  ...categories.value,
])

const filteredProjects = computed(() => {
  let result = projects.value

  if (activeCategoryId.value) {
    result = result.filter(p =>
      p.categories?.some(c => c.id === activeCategoryId.value)
    )
  }

  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    result = result.filter(p =>
      p.title.toLowerCase().includes(q) ||
      p.description.toLowerCase().includes(q)
    )
  }

  return result
})

async function loadData() {
  isLoading.value = true
  try {
    const [projectsData, categoriesData] = await Promise.all([
      projectApi.listProjects({ status: 'active' }),
      adminApi.listCategories(),
    ])
    projects.value = projectsData
    categories.value = categoriesData
  } catch (err) {
    console.error('Failed to load projects:', err)
  } finally {
    isLoading.value = false
  }
}

function onProjectClick(id: number) {
  router.push({ name: 'project-detail', params: { id: String(id) } })
}

function selectCategory(id: number | null) {
  activeCategoryId.value = id
}

function formatAmount(amount: number) {
  return amount.toLocaleString('ru-RU')
}

const totalCurrent = computed(() =>
  projects.value.reduce((sum, p) => sum + p.current_amount, 0)
)

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="min-h-screen">
    <!-- Hero Section -->
    <section class="bg-gradient-to-br from-blue-600 via-purple-600 to-pink-600 text-white py-16 sm:py-20 px-4 sm:px-6">
      <div class="max-w-7xl mx-auto text-center">
        <h1 class="mb-4 sm:mb-6">Финансируйте будущее технологий</h1>
        <p class="max-w-2xl mx-auto mb-6 sm:mb-8 text-base sm:text-lg opacity-90">
          Поддержите инновационные IT-проекты или запустите свой собственный.
          Вместе мы создаем технологии завтрашнего дня.
        </p>
        <div class="max-w-xl mx-auto relative">
          <Search class="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-neutral-400" />
          <Input
            type="text"
            placeholder="Поиск проектов..."
            v-model="searchQuery"
            class="pl-12 h-12 sm:h-14 bg-white/10 backdrop-blur-sm border-white/20 text-white placeholder:text-white/60"
          />
        </div>
      </div>
    </section>

    <!-- Filters -->
    <section class="sticky top-[57px] sm:top-16 z-40 bg-white border-b border-neutral-200 shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 py-3 sm:py-4">
        <div class="flex gap-2 overflow-x-auto scrollbar-hide">
          <button
            v-for="cat in allCategories"
            :key="String(cat.id)"
            @click="selectCategory(cat.id)"
            class="px-4 sm:px-6 py-2 rounded-full whitespace-nowrap transition-all text-sm"
            :class="activeCategoryId === cat.id
              ? 'bg-gradient-to-r from-blue-600 to-purple-600 text-white shadow-lg'
              : 'bg-neutral-100 text-neutral-700 hover:bg-neutral-200'"
          >
            {{ cat.title }}
          </button>
        </div>
      </div>
    </section>

    <!-- Stats -->
    <section class="bg-neutral-100 py-8 sm:py-12 px-4 sm:px-6">
      <div class="max-w-7xl mx-auto grid grid-cols-2 lg:grid-cols-4 gap-4 sm:gap-6">
        <div class="bg-white rounded-lg p-4 sm:p-6 text-center">
          <div class="text-blue-600 mb-2">Активных проектов</div>
          <div class="text-neutral-900">{{ projects.length }}</div>
        </div>
        <div class="bg-white rounded-lg p-4 sm:p-6 text-center">
          <div class="text-purple-600 mb-2">Собрано средств</div>
          <div class="text-neutral-900">{{ formatAmount(totalCurrent) }} ₽</div>
        </div>
        <div class="bg-white rounded-lg p-4 sm:p-6 text-center">
          <div class="text-pink-600 mb-2">Успешных запусков</div>
          <div class="text-neutral-900">142</div>
        </div>
        <div class="bg-white rounded-lg p-4 sm:p-6 text-center">
          <div class="text-orange-600 mb-2">Спонсоров</div>
          <div class="text-neutral-900">3,847</div>
        </div>
      </div>
    </section>

    <!-- Projects Grid -->
    <section class="py-8 sm:py-12 px-4 sm:px-6">
      <div class="max-w-7xl mx-auto">
        <div v-if="isLoading" class="text-center py-20">
          <div class="text-neutral-400">Загрузка проектов...</div>
        </div>
        <div v-else-if="filteredProjects.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 sm:gap-8">
          <ProjectCard
            v-for="project in filteredProjects"
            :key="project.id"
            :project="project"
            @click="onProjectClick(project.id)"
          />
        </div>
        <div v-else class="text-center py-20">
          <div class="text-neutral-400 mb-2">Проекты не найдены</div>
          <p class="text-neutral-500">Попробуйте изменить параметры поиска</p>
        </div>
      </div>
    </section>
  </div>
</template>
