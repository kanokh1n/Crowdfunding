<script setup lang="ts">
import { inject, computed, type Ref } from 'vue'
import { cn } from '@/lib/utils'

interface Props {
  value: string
  class?: string
}

const props = defineProps<Props>()

const activeTab = inject<Ref<string>>('activeTab')
const activateTab = inject<(value: string) => void>('activateTab')

const isActive = computed(() => activeTab?.value === props.value)

function onClick() {
  activateTab?.(props.value)
}
</script>

<template>
  <button
    :class="cn(
      'inline-flex h-[calc(100%-1px)] flex-1 items-center justify-center rounded-xl border border-transparent px-2 py-1.5 text-sm font-medium whitespace-nowrap transition-colors',
      isActive ? 'bg-card text-foreground' : 'text-foreground/60 hover:text-foreground',
      props.class
    )"
    @click="onClick"
  >
    <slot />
  </button>
</template>
