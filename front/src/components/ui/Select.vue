<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { cn } from '@/lib/utils'

export interface SelectOption {
  value: number | null
  label: string
}

interface Props {
  modelValue?: number | null
  placeholder?: string
  options?: SelectOption[]
  class?: string
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: 'Выберите...',
  options: () => [],
})

const emit = defineEmits<{
  'update:modelValue': [value: number | null]
}>()

const isOpen = ref(false)

const selectedLabel = computed(() => {
  const found = props.options.find(o => o.value === props.modelValue)
  return found ? found.label : props.placeholder
})

function select(value: number | null) {
  emit('update:modelValue', value)
  isOpen.value = false
}

function toggle() {
  if (!props.disabled) {
    isOpen.value = !isOpen.value
  }
}

function close(e: MouseEvent) {
  const target = e.target as HTMLElement
  if (!target.closest('[data-select]')) {
    isOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', close)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', close)
})
</script>

<template>
  <div class="relative" data-select @click.stop>
    <!-- Trigger -->
    <button
      type="button"
      :disabled="disabled"
      :class="cn(
        'flex h-9 w-full items-center justify-between gap-2 rounded-md border border-input bg-input-background px-3 py-2 text-sm whitespace-nowrap transition-[color,box-shadow] outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50',
        isOpen && 'ring-2 ring-ring',
        props.class
      )"
      @click="toggle"
      @keydown.esc="isOpen = false"
    >
      <span :class="{ 'text-muted-foreground': !modelValue }">{{ selectedLabel }}</span>
      <svg class="pointer-events-none size-4 shrink-0 opacity-50" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="m6 9 6 6 6-6"/>
      </svg>
    </button>

    <!-- Dropdown -->
    <div
      v-if="isOpen"
      class="absolute z-50 mt-1 w-full min-w-[8rem] overflow-hidden rounded-md border bg-popover p-1 text-popover-foreground shadow-md"
    >
      <div
        v-for="option in options"
        :key="String(option.value)"
        class="relative flex w-full cursor-pointer items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-none hover:bg-accent hover:text-accent-foreground"
        :class="{ 'bg-accent text-accent-foreground': modelValue === option.value }"
        @click="select(option.value)"
      >
        {{ option.label }}
      </div>
    </div>
  </div>
</template>
