<script setup lang="ts">
import { cn } from '@/lib/utils'

interface Props {
  variant?: 'default' | 'destructive' | 'outline' | 'secondary' | 'ghost' | 'link'
  size?: 'default' | 'sm' | 'lg' | 'icon'
  asChild?: boolean
  class?: string
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'default',
  size: 'default',
  asChild: false,
})

const emit = defineEmits<{
  click: [e: MouseEvent]
}>()

function handleClick(e: MouseEvent) {
  emit('click', e)
}
</script>

<template>
  <button
    :class="cn(
      'inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-all disabled:pointer-events-none disabled:opacity-50 outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2',
      {
        'bg-primary text-primary-foreground hover:bg-primary/90': variant === 'default',
        'bg-destructive text-white hover:bg-destructive/90': variant === 'destructive',
        'border border-border bg-background text-foreground hover:bg-accent hover:text-accent-foreground': variant === 'outline',
        'bg-secondary text-secondary-foreground hover:bg-secondary/80': variant === 'secondary',
        'hover:bg-accent hover:text-accent-foreground': variant === 'ghost',
        'text-primary underline-offset-4 hover:underline': variant === 'link',
      },
      {
        'h-9 px-4 py-2': size === 'default',
        'h-8 rounded-md px-3': size === 'sm',
        'h-10 rounded-md px-6': size === 'lg',
        'size-9': size === 'icon',
      },
      props.class
    )"
    :disabled="props.disabled"
    @click="handleClick"
  >
    <slot />
  </button>
</template>
