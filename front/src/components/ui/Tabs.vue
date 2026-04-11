<script setup lang="ts">
import { ref, provide, computed } from 'vue'
import { cn } from '@/lib/utils'

interface Props {
  defaultValue?: string
  modelValue?: string
  class?: string
}

const props = withDefaults(defineProps<Props>(), {
  defaultValue: '',
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const activeTab = ref(props.modelValue || props.defaultValue)

function activateTab(value: string) {
  activeTab.value = value
  emit('update:modelValue', value)
}

provide('activeTab', activeTab)
provide('activateTab', activateTab)
</script>

<template>
  <div :class="cn('flex flex-col gap-2', props.class)">
    <slot />
  </div>
</template>
