<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import Label from '@/components/ui/Label.vue'
import Textarea from '@/components/ui/Textarea.vue'
import Select from '@/components/ui/Select.vue'
import type { SelectOption } from '@/components/ui/Select.vue'
import { ArrowLeft, Image as ImageIcon } from 'lucide-vue-next'
import * as projectApi from '@/api/projects'
import * as adminApi from '@/api/admin'
import type { Category } from '@/types'

const router = useRouter()
const categories = ref<Category[]>([])
const isLoading = ref(false)
const submitError = ref('')

const title = ref('')
const categoryId = ref<number | null>(null)
const description = ref('')
const fullDescription = ref('')
const goalAmount = ref('')
const endDate = ref('')
const projectImg = ref('')

const socialLinks = ref({
  website: '',
  github: '',
  twitter: '',
  linkedin: '',
})

async function loadCategories() {
  try {
    categories.value = await adminApi.listCategories()
  } catch {
    console.error('Failed to load categories')
  }
}

async function handleSubmit() {
  submitError.value = ''
  isLoading.value = true

  try {
    await projectApi.createProject({
      title: title.value,
      description: description.value,
      goal_amount: parseFloat(goalAmount.value),
      end_date: endDate.value || undefined,
      project_img: projectImg.value || undefined,
      category_ids: categoryId.value ? [categoryId.value] : undefined,
    })
    router.push({ name: 'home' })
  } catch (err: any) {
    submitError.value = err.message || 'Ошибка создания проекта'
  } finally {
    isLoading.value = false
  }
}

function handleCancel() {
  router.back()
}

const categoryOptions = computed<SelectOption[]>(() =>
  categories.value.map(c => ({ value: c.id, label: c.title }))
)

onMounted(() => {
  loadCategories()
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
          Создать новый проект
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

          <!-- Категория -->
          <div class="space-y-2">
            <Label for="category">Категория</Label>
            <Select
              id="category"
              v-model="categoryId"
              placeholder="Выберите категорию"
              :options="categoryOptions"
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
              v-model="fullDescription"
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
            <p class="text-neutral-500 text-sm">
              Оставьте пустым для использования изображения по умолчанию
            </p>
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
              {{ isLoading ? 'Создание...' : 'Создать проект' }}
            </Button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
