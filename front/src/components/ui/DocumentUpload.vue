<script setup lang="ts">
import { ref } from 'vue'
import { Upload, FileText, X } from 'lucide-vue-next'

const props = defineProps<{
  modelValue?: string
  fileName?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
  'update:fileName': [value: string]
}>()

const isUploading = ref(false)
const error = ref('')
const localFileName = ref(props.fileName || '')
const fileInput = ref<HTMLInputElement>()

async function handleFile(file: File) {
  error.value = ''
  isUploading.value = true
  try {
    const form = new FormData()
    form.append('file', file)
    const res = await fetch('/api/upload/document', {
      method: 'POST',
      headers: { Authorization: `Bearer ${localStorage.getItem('access_token')}` },
      body: form,
    })
    if (!res.ok) {
      const data = await res.json().catch(() => ({}))
      throw new Error(data.error || 'Ошибка загрузки')
    }
    const data = await res.json()
    emit('update:modelValue', data.url)
    localFileName.value = data.name || file.name
    emit('update:fileName', localFileName.value)
  } catch (err: any) {
    error.value = err.message || 'Ошибка загрузки'
  } finally {
    isUploading.value = false
  }
}

function onFileChange(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (file) handleFile(file)
}

function onDrop(e: DragEvent) {
  const file = e.dataTransfer?.files[0]
  if (file) handleFile(file)
}

function clear() {
  emit('update:modelValue', '')
  emit('update:fileName', '')
  localFileName.value = ''
  if (fileInput.value) fileInput.value.value = ''
}
</script>

<template>
  <div class="space-y-2">
    <div
      v-if="!modelValue"
      class="relative border-2 border-dashed border-neutral-300 rounded-xl p-8 text-center cursor-pointer hover:border-blue-400 hover:bg-blue-50 transition-colors"
      @dragover.prevent
      @drop.prevent="onDrop"
      @click="fileInput?.click()"
    >
      <input
        ref="fileInput"
        type="file"
        accept=".pdf,.pptx"
        class="hidden"
        @change="onFileChange"
      />
      <Upload class="w-8 h-8 mx-auto mb-2 text-neutral-400" />
      <p class="text-neutral-600 text-sm">Нажмите или перетащите файл</p>
      <p class="text-neutral-400 text-xs mt-1">PDF или PPTX — до 20 МБ</p>
      <div
        v-if="isUploading"
        class="absolute inset-0 flex items-center justify-center bg-white/80 rounded-xl"
      >
        <span class="text-sm text-blue-600">Загрузка...</span>
      </div>
    </div>

    <div v-else class="flex items-center gap-3 p-4 bg-neutral-50 border border-neutral-200 rounded-xl">
      <FileText class="w-8 h-8 text-blue-600 shrink-0" />
      <span class="flex-1 text-sm text-neutral-700 truncate">{{ localFileName || 'Файл загружен' }}</span>
      <button
        type="button"
        class="p-1 rounded hover:bg-neutral-200 transition-colors"
        @click="clear"
      >
        <X class="w-4 h-4 text-neutral-500" />
      </button>
    </div>

    <p v-if="error" class="text-red-500 text-sm">{{ error }}</p>
  </div>
</template>
