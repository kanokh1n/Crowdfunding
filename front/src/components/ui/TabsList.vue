<script setup lang="ts">
import { cn } from '@/lib/utils'

interface Props {
  class?: string
}

const props = defineProps<Props>()
</script>

<template>
  <div :class="cn('flex h-10 w-full items-center justify-center rounded-xl bg-muted p-[3px] text-muted-foreground', props.class)">
    <slot />
  </div>
</template>
