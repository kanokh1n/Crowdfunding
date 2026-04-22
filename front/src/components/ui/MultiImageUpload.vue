<script setup lang="ts">
import { ref } from 'vue'
import { Upload, X } from 'lucide-vue-next'

const MAX = 5

const props = defineProps<{
  modelValue: string[]
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string[]]
}>()

const isUploading = ref(false)
const error = ref('')
const fileInput = ref<HTMLInputElement>()

async function handleFiles(files: FileList | File[]) {
  const arr = Array.from(files)
  const remaining = MAX - props.modelValue.length
  if (remaining <= 0) return

  error.value = ''
  isUploading.value = true
  try {
    const toUpload = arr.slice(0, remaining)
    const urls = await Promise.all(
      toUpload.map(async (file) => {
        if (file.size > 5 * 1024 * 1024) throw new Error(`Файл ${file.name} больше 5 МБ`)
        const form = new FormData()
        form.append('file', file)
        const res = await fetch('/api/upload', {
          method: 'POST',
          headers: { Authorization: `Bearer ${localStorage.getItem('access_token')}` },
          body: form,
        })
        if (!res.ok) {
          const data = await res.json().catch(() => ({}))
          throw new Error(data.error || 'Ошибка загрузки')
        }
        const data = await res.json()
        return data.url as string
      })
    )
    emit('update:modelValue', [...props.modelValue, ...urls])
  } catch (err: any) {
    error.value = err.message || 'Ошибка загрузки'
  } finally {
    isUploading.value = false
    if (fileInput.value) fileInput.value.value = ''
  }
}

function onFileChange(e: Event) {
  const files = (e.target as HTMLInputElement).files
  if (files) handleFiles(files)
}

function onDrop(e: DragEvent) {
  if (e.dataTransfer?.files) handleFiles(e.dataTransfer.files)
}

function remove(index: number) {
  const updated = [...props.modelValue]
  updated.splice(index, 1)
  emit('update:modelValue', updated)
}
</script>

<template>
  <div class="space-y-3">
    <!-- Превью загруженных -->
    <div v-if="modelValue.length > 0" class="grid grid-cols-3 sm:grid-cols-5 gap-2">
      <div
        v-for="(url, i) in modelValue"
        :key="url"
        class="relative aspect-square rounded-lg overflow-hidden border border-neutral-200 bg-neutral-100"
      >
        <img :src="url" :alt="`Фото ${i + 1}`" class="w-full h-full object-contain" />
        <button
          type="button"
          class="absolute top-1 right-1 w-6 h-6 bg-white/90 rounded-full flex items-center justify-center hover:bg-red-50 shadow transition-colors"
          @click="remove(i)"
        >
          <X class="w-3.5 h-3.5 text-neutral-600" />
        </button>
        <div v-if="i === 0" class="absolute bottom-1 left-1 px-1.5 py-0.5 bg-blue-600 text-white rounded text-xs">
          Главная
        </div>
      </div>
    </div>

    <!-- Зона загрузки -->
    <div
      v-if="modelValue.length < MAX"
      class="relative border-2 border-dashed border-neutral-300 rounded-xl p-6 text-center cursor-pointer hover:border-blue-400 hover:bg-blue-50 transition-colors"
      @dragover.prevent
      @drop.prevent="onDrop"
      @click="fileInput?.click()"
    >
      <input
        ref="fileInput"
        type="file"
        accept="image/jpeg,image/png,image/webp,image/gif"
        multiple
        class="hidden"
        @change="onFileChange"
      />
      <Upload class="w-7 h-7 mx-auto mb-2 text-neutral-400" />
      <p class="text-neutral-600 text-sm">Нажмите или перетащите изображения</p>
      <p class="text-neutral-400 text-xs mt-1">JPG, PNG, WebP — до 5 МБ каждое</p>
      <p class="text-neutral-400 text-xs">Рекомендуемый размер: 1280×720 (16:9) · Загружено: {{ modelValue.length }}/{{ MAX }}</p>
      <div
        v-if="isUploading"
        class="absolute inset-0 flex items-center justify-center bg-white/80 rounded-xl"
      >
        <span class="text-sm text-blue-600">Загрузка...</span>
      </div>
    </div>

    <p v-if="error" class="text-red-500 text-sm">{{ error }}</p>
  </div>
</template>
