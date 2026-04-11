<script setup lang="ts">
import { inject, computed, type Ref } from 'vue'

interface Props {
  value: string
}

const props = defineProps<Props>()

const activeTab = inject<Ref<string>>('activeTab')
const isActive = computed(() => activeTab?.value === props.value)
</script>

<template>
  <div v-if="isActive">
    <slot />
  </div>
</template>
