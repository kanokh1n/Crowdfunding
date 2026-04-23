<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import Label from '@/components/ui/Label.vue'
import Textarea from '@/components/ui/Textarea.vue'
import Select from '@/components/ui/Select.vue'
import MultiImageUpload from '@/components/ui/MultiImageUpload.vue'
import DocumentUpload from '@/components/ui/DocumentUpload.vue'
import type { SelectOption } from '@/components/ui/Select.vue'
import { ArrowLeft } from 'lucide-vue-next'
import * as projectApi from '@/api/projects'
import * as adminApi from '@/api/admin'
import type { Category } from '@/types'
import { useErrorToast } from '@/composables/useErrorToast'

const router = useRouter()
const { showError, showSuccess } = useErrorToast()
const categories = ref<Category[]>([])
const isLoading = ref(false)
const submitError = ref('')

const title = ref('')
const categoryId = ref<number | null>(null)
const shortDescription = ref('')
const description = ref('')
const goalAmount = ref('')
const endDate = ref('')
const projectImages = ref<string[]>([])
const roadmapFile = ref('')
const roadmapFileName = ref('')

const socialLinks = ref({
  github: '',
  telegram: '',
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
      short_description: shortDescription.value || undefined,
      description: description.value,
      goal_amount: parseFloat(goalAmount.value),
      end_date: endDate.value ? new Date(endDate.value).toISOString() : undefined,
      project_img: projectImages.value[0] || undefined,
      images: projectImages.value.length > 0 ? projectImages.value : undefined,
      roadmap_file: roadmapFile.value || undefined,
      category_ids: categoryId.value ? [categoryId.value] : undefined,
      link_telegram: socialLinks.value.telegram || undefined,
      link_github: socialLinks.value.github || undefined,
      link_linkedin: socialLinks.value.linkedin || undefined,
    })
    showSuccess('Проект успешно создан и отправлен на проверку!')
    router.push({ name: 'home' })
  } catch (err: any) {
    showError(err)
    submitError.value = err.userMessage ?? err.message ?? 'Ошибка создания проекта'
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
            <Label for="shortDescription">Краткое описание <span class="text-neutral-400 font-normal text-xs">(выводится на карточке, до 200 символов)</span></Label>
            <Textarea
              id="shortDescription"
              v-model="shortDescription"
              placeholder="Одно-два предложения о том, что делает ваш проект"
              :rows="2"
              :maxlength="200"
            />
            <div class="text-neutral-400 text-right text-xs">{{ shortDescription.length }}/200</div>
          </div>

          <!-- Полное описание -->
          <div class="space-y-2">
            <Label for="description">Полное описание *</Label>
            <Textarea
              id="description"
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

          <!-- Изображения -->
          <div class="space-y-2">
            <Label>Изображения проекта <span class="text-neutral-400 font-normal text-xs">(до 5 штук, первое — главное)</span></Label>
            <MultiImageUpload v-model="projectImages" />
          </div>

          <!-- Roadmap -->
          <div class="space-y-2">
            <Label>Roadmap * <span class="text-neutral-400 font-normal text-xs">(PDF или PPTX — обязателен для прохождения проверки)</span></Label>
            <DocumentUpload v-model="roadmapFile" v-model:fileName="roadmapFileName" />
          </div>

          <!-- Соцсети -->
          <div class="space-y-4">
            <div class="space-y-2">
              <Label for="telegramLink">Telegram</Label>
              <Input id="telegramLink" v-model="socialLinks.telegram" placeholder="https://t.me/yourproject" />
            </div>
            <div class="space-y-2">
              <Label for="githubLink">GitHub</Label>
              <Input id="githubLink" v-model="socialLinks.github" placeholder="https://github.com/yourproject" />
            </div>
            <div class="space-y-2">
              <Label for="linkedinLink">LinkedIn</Label>
              <Input id="linkedinLink" v-model="socialLinks.linkedin" placeholder="https://linkedin.com/in/yourprofile" />
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
              {{ isLoading ? 'Создание...' : 'Создать проект' }}
            </Button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
