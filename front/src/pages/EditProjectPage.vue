<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import Label from '@/components/ui/Label.vue'
import Textarea from '@/components/ui/Textarea.vue'
import { ArrowLeft, Image as ImageIcon } from 'lucide-vue-next'
import * as projectApi from '@/api/projects'
import * as adminApi from '@/api/admin'
import type { Category } from '@/types'

const props = defineProps<{
  id: string
}>()

const router = useRouter()
const isLoading = ref(false)
const submitError = ref('')
const categories = ref<Category[]>([])

const title = ref('')
const description = ref('')
const goalAmount = ref('')
const endDate = ref('')
const projectImg = ref('')
const categoryIds = ref<number[]>([])

async function loadData() {
  try {
    const [project, cats] = await Promise.all([
      projectApi.getProject(parseInt(props.id)),
      adminApi.listCategories(),
    ])
    title.value = project.title
    description.value = project.description
    goalAmount.value = String(project.goal_amount)
    endDate.value = project.end_date ? project.end_date.slice(0, 10) : ''
    projectImg.value = project.project_img ?? ''
    categoryIds.value = project.categories?.map(c => c.id) ?? []
    categories.value = cats
  } catch {
    submitError.value = 'Не удалось загрузить проект'
  }
}

async function handleSubmit() {
  submitError.value = ''
  isLoading.value = true

  try {
    await projectApi.updateProject(parseInt(props.id), {
      title: title.value,
      description: description.value,
      goal_amount: parseFloat(goalAmount.value),
      end_date: endDate.value ? endDate.value + 'T00:00:00Z' : undefined,
      project_img: projectImg.value || undefined,
      category_ids: categoryIds.value,
    })
    router.push({ name: 'project-detail', params: { id: props.id } })
  } catch (err: any) {
    submitError.value = err.message || 'Ошибка сохранения проекта'
  } finally {
    isLoading.value = false
  }
}

function handleCancel() {
  router.push({ name: 'project-detail', params: { id: props.id } })
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="min-h-screen py-8 sm:py-12 px-4 sm:px-6">
    <div class="max-w-4xl mx-auto">
      <Button
        variant="ghost"
        @click="handleCancel"
        class="mb-6"
      >
        <ArrowLeft class="w-4 h-4 mr-2" />
        Назад
      </Button>

      <div class="bg-white rounded-xl shadow-lg p-6 sm:p-8 border border-neutral-200">
        <h1 class="mb-4 sm:mb-6">
          Редактировать проект
        </h1>

        <form @submit.prevent="handleSubmit" class="space-y-6">
          <!-- Название -->
          <div class="space-y-2">
            <Label for="title">Название проекта *</Label>
            <Input
              id="title"
              v-model="title"
              placeholder="Введите название проекта"
              required
            />
          </div>

          <!-- Краткое описание -->
          <div class="space-y-2">
            <Label for="description">Краткое описание *</Label>
            <Textarea
              id="description"
              v-model="description"
              placeholder="Краткое описание проекта (до 200 символов)"
              :rows="3"
              :maxlength="200"
              required
            />
            <div class="text-neutral-500 text-right text-sm">
              {{ description.length }}/200
            </div>
          </div>

          <!-- Полное описание -->
          <div class="space-y-2">
            <Label for="fullDescription">Полное описание *</Label>
            <Textarea
              id="fullDescription"
              v-model="description"
              placeholder="Подробное описание проекта, его целей и особенностей"
              :rows="8"
              required
            />
          </div>

          <!-- Целевая сумма -->
          <div class="space-y-2">
            <Label for="goalAmount">Целевая сумма (₽) *</Label>
            <Input
              id="goalAmount"
              type="number"
              v-model="goalAmount"
              placeholder="100000"
              min="1"
              step="0.01"
              required
            />
          </div>

          <!-- Дата окончания -->
          <div class="space-y-2">
            <Label for="endDate">Дата окончания (необязательно)</Label>
            <Input
              id="endDate"
              type="date"
              v-model="endDate"
            />
          </div>

          <!-- Изображение -->
          <div class="space-y-2">
            <Label for="image">URL изображения</Label>
            <div class="relative">
              <ImageIcon class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-neutral-400" />
              <Input
                id="image"
                v-model="projectImg"
                placeholder="https://example.com/image.jpg"
                class="pl-10"
              />
            </div>
          </div>

          <!-- Ошибка -->
          <div v-if="submitError" class="text-red-500 text-sm text-center">
            {{ submitError }}
          </div>

          <!-- Кнопки -->
          <div class="flex gap-4 pt-6 border-t">
            <Button type="button" variant="outline" @click="handleCancel" class="flex-1" :disabled="isLoading">
              Отмена
            </Button>
            <Button type="submit" class="flex-1" :disabled="isLoading">
              {{ isLoading ? 'Сохранение...' : 'Сохранить изменения' }}
            </Button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
