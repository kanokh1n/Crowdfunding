<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/lib/utils'

interface Props {
  modelValue?: number
  class?: string
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: 0,
})

const style = computed(() => ({
  transform: `translateX(-${100 - props.modelValue}%)`,
}))
</script>

<template>
  <div :class="cn('relative h-2 w-full overflow-hidden rounded-full bg-primary/20', props.class)">
    <div
      class="h-full w-full flex-1 bg-primary transition-all"
      :style="style"
    />
  </div>
</template>
