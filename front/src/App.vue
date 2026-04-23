<script setup lang="ts">
import { RouterView } from 'vue-router'
import Header from '@/components/Header.vue'
import Footer from '@/components/Footer.vue'
import ToastContainer from '@/components/ui/ToastContainer.vue'
</script>

<template>
  <div class="min-h-screen flex flex-col bg-neutral-50">
    <Header />
    <main class="flex-1">
      <RouterView />
    </main>
    <Footer />
    <ToastContainer />
  </div>
</template>
