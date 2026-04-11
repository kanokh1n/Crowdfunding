<script setup lang="ts">
import { cn } from '@/lib/utils'

interface Props {
  variant?: 'default' | 'secondary' | 'destructive' | 'outline'
  class?: string
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'default',
})
</script>

<template>
  <span :class="cn(
    'inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium transition-colors',
    {
      'bg-primary text-primary-foreground': variant === 'default',
      'bg-secondary text-secondary-foreground': variant === 'secondary',
      'bg-destructive text-destructive-foreground': variant === 'destructive',
      'border border-border bg-background text-foreground': variant === 'outline',
    },
    props.class
  )">
    <slot />
  </span>
</template>
