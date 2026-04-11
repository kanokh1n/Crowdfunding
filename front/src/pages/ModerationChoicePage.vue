<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import { ShieldCheck, Map, Brain, ArrowLeft } from 'lucide-vue-next'
import * as adminApi from '@/api/admin'

const props = defineProps<{
  id: string
}>()

const router = useRouter()
const isLoading = ref(false)
const submitError = ref('')
const showSuccess = ref(false)
const successMessage = ref('')

async function handleChoice(choice: 'roadmap' | 'ai') {
  isLoading.value = true
  submitError.value = ''
  try {
    if (choice === 'ai') {
      await adminApi.recheckProject(parseInt(props.id))
    } else {
      // Roadmap — ручное одобрение модератором
      await adminApi.moderationDecision(parseInt(props.id), { decision: 'approve' })
    }
    // Показываем уведомление
    const methodName = choice === 'roadmap' ? 'Roadmap проверка' : 'AI проверка'
    successMessage.value = `Выбран метод модерации: ${methodName}`
    showSuccess.value = true

    // Переходим к модерации после показа уведомления
    setTimeout(() => {
      router.push({ name: 'moderation' })
    }, 1500)
  } catch (err: any) {
    submitError.value = err.message || 'Ошибка при выполнении действия'
    isLoading.value = false
  }
}

function goBack() {
  router.push({ name: 'moderation' })
}
</script>

<template>
  <!-- Toast notification -->
  <Transition name="toast">
    <div
      v-if="showSuccess"
      class="fixed top-20 right-4 sm:right-6 z-[9999] max-w-sm shadow-lg rounded-xl border border-blue-200 bg-white px-4 py-3 flex items-start gap-3"
    >
      <div class="flex-shrink-0 w-8 h-8 rounded-full bg-blue-100 flex items-center justify-center">
        <ShieldCheck class="w-5 h-5 text-blue-600" />
      </div>
      <div class="flex-1 min-w-0">
        <div class="text-sm font-medium text-neutral-900">IT Crowdfunding</div>
        <div class="text-sm text-neutral-600">{{ successMessage }}</div>
      </div>
    </div>
  </Transition>

  <div class="min-h-screen py-8 sm:py-12 px-4 sm:px-6">
    <div class="max-w-4xl mx-auto">
      <Button variant="ghost" @click="goBack" class="mb-4 sm:mb-6">
        <ArrowLeft class="w-4 h-4 mr-2" />
        Назад к модерации
      </Button>

      <div class="text-center mb-8 sm:mb-12">
        <div class="flex items-center justify-center gap-2 sm:gap-3 mb-3 sm:mb-4">
          <ShieldCheck class="w-8 h-8 sm:w-10 sm:h-10 text-blue-600" />
          <h1>Выбор метода модерации</h1>
        </div>
        <p class="text-neutral-600 max-w-2xl mx-auto text-sm sm:text-base">
          Выберите способ модерации проекта. Вы можете использовать проверку через Roadmap
          или автоматическую проверку с помощью нейронной сети.
        </p>
      </div>

      <div v-if="submitError" class="text-red-500 text-sm text-center mb-4">
        {{ submitError }}
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-6 sm:gap-8">
        <!-- Roadmap -->
        <div
          class="bg-white rounded-xl shadow-lg border-2 border-neutral-200 hover:border-blue-500 transition-all p-6 sm:p-8 cursor-pointer group"
          @click="handleChoice('roadmap')"
        >
          <div class="flex flex-col items-center text-center">
            <div class="w-16 h-16 sm:w-20 sm:h-20 bg-gradient-to-br from-blue-500 to-blue-600 rounded-2xl flex items-center justify-center mb-4 sm:mb-6 group-hover:scale-110 transition-transform">
              <Map class="w-8 h-8 sm:w-10 sm:h-10 text-white" />
            </div>

            <h2 class="mb-3 sm:mb-4">Roadmap проверка</h2>

            <p class="text-neutral-600 mb-4 sm:mb-6 leading-relaxed text-sm sm:text-base">
              Ручная проверка проекта с анализом дорожной карты, планов развития и ключевых этапов.
              Модератор оценивает реалистичность целей и сроков.
            </p>

            <div class="space-y-2 sm:space-y-3 mb-6 sm:mb-8 w-full text-left text-sm">
              <div class="flex items-center gap-2 text-neutral-700">
                <div class="w-2 h-2 bg-blue-500 rounded-full flex-shrink-0"></div>
                <span>Детальный анализ планов</span>
              </div>
              <div class="flex items-center gap-2 text-neutral-700">
                <div class="w-2 h-2 bg-blue-500 rounded-full flex-shrink-0"></div>
                <span>Оценка этапов развития</span>
              </div>
              <div class="flex items-center gap-2 text-neutral-700">
                <div class="w-2 h-2 bg-blue-500 rounded-full flex-shrink-0"></div>
                <span>Проверка реалистичности</span>
              </div>
              <div class="flex items-center gap-2 text-neutral-700">
                <div class="w-2 h-2 bg-blue-500 rounded-full flex-shrink-0"></div>
                <span>Рекомендации по улучшению</span>
              </div>
            </div>

            <Button class="w-full group-hover:bg-blue-700 text-sm sm:text-base" :disabled="isLoading">
              Выбрать Roadmap
            </Button>
          </div>
        </div>

        <!-- AI проверка -->
        <div
          class="bg-white rounded-xl shadow-lg border-2 border-neutral-200 hover:border-purple-500 transition-all p-6 sm:p-8 cursor-pointer group"
          @click="handleChoice('ai')"
        >
          <div class="flex flex-col items-center text-center">
            <div class="w-16 h-16 sm:w-20 sm:h-20 bg-gradient-to-br from-purple-500 to-purple-600 rounded-2xl flex items-center justify-center mb-4 sm:mb-6 group-hover:scale-110 transition-transform">
              <Brain class="w-8 h-8 sm:w-10 sm:h-10 text-white" />
            </div>

            <h2 class="mb-3 sm:mb-4">AI проверка</h2>

            <p class="text-neutral-600 mb-4 sm:mb-6 leading-relaxed text-sm sm:text-base">
              Автоматическая проверка проекта с использованием нейронной сети.
              AI анализирует описание, категорию, цели и выявляет потенциальные проблемы.
            </p>

            <div class="space-y-2 sm:space-y-3 mb-6 sm:mb-8 w-full text-left text-sm">
              <div class="flex items-center gap-2 text-neutral-700">
                <div class="w-2 h-2 bg-purple-500 rounded-full flex-shrink-0"></div>
                <span>Быстрая автоматическая проверка</span>
              </div>
              <div class="flex items-center gap-2 text-neutral-700">
                <div class="w-2 h-2 bg-purple-500 rounded-full flex-shrink-0"></div>
                <span>Анализ содержимого</span>
              </div>
              <div class="flex items-center gap-2 text-neutral-700">
                <div class="w-2 h-2 bg-purple-500 rounded-full flex-shrink-0"></div>
                <span>Выявление несоответствий</span>
              </div>
              <div class="flex items-center gap-2 text-neutral-700">
                <div class="w-2 h-2 bg-purple-500 rounded-full flex-shrink-0"></div>
                <span>Мгновенный результат</span>
              </div>
            </div>

            <Button class="w-full bg-purple-600 hover:bg-purple-700 group-hover:bg-purple-700 text-sm sm:text-base" :disabled="isLoading">
              Выбрать AI проверку
            </Button>
          </div>
        </div>
      </div>

      <!-- Информационная панель -->
      <div class="mt-8 sm:mt-12 bg-blue-50 border border-blue-200 rounded-xl p-4 sm:p-6">
        <div class="flex items-start gap-2 sm:gap-3">
          <ShieldCheck class="w-5 h-5 sm:w-6 sm:h-6 text-blue-600 flex-shrink-0 mt-1" />
          <div>
            <h3 class="mb-1 sm:mb-2 font-semibold text-sm sm:text-base">Рекомендации по выбору</h3>
            <p class="text-neutral-700 leading-relaxed text-xs sm:text-sm">
              <strong>Roadmap</strong> — используйте для крупных и сложных проектов, требующих детального анализа.<br />
              <strong>AI проверка</strong> — подходит для быстрой предварительной оценки и простых проектов.
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}
.toast-enter-from {
  opacity: 0;
  transform: translateX(40px);
}
.toast-leave-to {
  opacity: 0;
  transform: translateX(40px);
}
</style>
