<script setup lang="ts">
import { useRouter } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import {
  Lightbulb, FileText, Bot, ShieldCheck, Users,
  Banknote, CheckCircle, ArrowRight, Rocket
} from 'lucide-vue-next'

const router = useRouter()
</script>

<template>
  <div class="min-h-screen">

    <!-- Hero -->
    <section class="bg-gradient-to-br from-blue-600 via-blue-700 to-purple-700 text-white py-20 px-4 sm:px-6">
      <div class="max-w-4xl mx-auto text-center">
        <div class="inline-flex items-center gap-2 bg-white/20 rounded-full px-4 py-2 mb-6 text-sm">
          <Rocket class="w-4 h-4" />
          Путь от идеи до финансирования
        </div>
        <h1 class="text-3xl sm:text-5xl font-bold mb-6 leading-tight">
          Как работает<br />IT Crowdfunding
        </h1>
        <p class="text-lg sm:text-xl text-blue-100 max-w-2xl mx-auto leading-relaxed">
          Мы соединяем перспективные IT-проекты с инвесторами через прозрачный процесс проверки и поэтапное финансирование
        </p>
      </div>
    </section>

    <!-- Steps -->
    <section class="py-16 sm:py-24 px-4 sm:px-6 bg-white">
      <div class="max-w-5xl mx-auto">
        <h2 class="text-2xl sm:text-3xl font-bold text-center mb-4">Процесс для стартапера</h2>
        <p class="text-neutral-500 text-center mb-14">Шесть шагов от идеи до финансирования</p>

        <div class="relative">
          <!-- Vertical line -->
          <div class="hidden md:block absolute left-1/2 -translate-x-1/2 top-0 bottom-0 w-0.5 bg-neutral-200" />

          <div class="space-y-12">

            <!-- Step 1 -->
            <div class="relative flex flex-col md:flex-row items-center gap-6 md:gap-0">
              <div class="md:w-1/2 md:pr-12 md:text-right order-2 md:order-1">
                <h3 class="text-xl font-bold mb-2">У вас есть идея</h3>
                <p class="text-neutral-600 leading-relaxed">
                  Всё начинается с идеи. Неважно, на каком этапе вы находитесь — это может быть концепция, MVP или работающий прототип. Главное — чёткое понимание проблемы, которую вы решаете.
                </p>
              </div>
              <div class="relative z-10 w-14 h-14 rounded-full bg-blue-600 flex items-center justify-center shadow-lg order-1 md:order-2 shrink-0">
                <Lightbulb class="w-6 h-6 text-white" />
              </div>
              <div class="md:w-1/2 md:pl-12 order-3" />
            </div>

            <!-- Step 2 -->
            <div class="relative flex flex-col md:flex-row items-center gap-6 md:gap-0">
              <div class="md:w-1/2 md:pr-12 order-2 md:order-1" />
              <div class="relative z-10 w-14 h-14 rounded-full bg-purple-600 flex items-center justify-center shadow-lg order-1 shrink-0">
                <FileText class="w-6 h-6 text-white" />
              </div>
              <div class="md:w-1/2 md:pl-12 order-3">
                <h3 class="text-xl font-bold mb-2">Готовите roadmap</h3>
                <p class="text-neutral-600 leading-relaxed">
                  Подготовьте roadmap в формате PDF или PPTX: опишите этапы разработки, сроки, команду и бюджет. Roadmap — ключевой документ для прохождения проверки и будущего финансирования по этапам.
                </p>
              </div>
            </div>

            <!-- Step 3 -->
            <div class="relative flex flex-col md:flex-row items-center gap-6 md:gap-0">
              <div class="md:w-1/2 md:pr-12 md:text-right order-2 md:order-1">
                <h3 class="text-xl font-bold mb-2">Создаёте проект на платформе</h3>
                <p class="text-neutral-600 leading-relaxed">
                  Зарегистрируйтесь, заполните карточку проекта: название, описание, цель сбора, ссылки на соцсети и приложите roadmap. После отправки проект уходит на автоматическую проверку.
                </p>
              </div>
              <div class="relative z-10 w-14 h-14 rounded-full bg-blue-600 flex items-center justify-center shadow-lg order-1 md:order-2 shrink-0">
                <Rocket class="w-6 h-6 text-white" />
              </div>
              <div class="md:w-1/2 md:pl-12 order-3" />
            </div>

            <!-- Step 4 -->
            <div class="relative flex flex-col md:flex-row items-center gap-6 md:gap-0">
              <div class="md:w-1/2 md:pr-12 order-2 md:order-1" />
              <div class="relative z-10 w-14 h-14 rounded-full bg-purple-600 flex items-center justify-center shadow-lg order-1 shrink-0">
                <Bot class="w-6 h-6 text-white" />
              </div>
              <div class="md:w-1/2 md:pl-12 order-3">
                <h3 class="text-xl font-bold mb-2">Автоматическая AI-проверка</h3>
                <p class="text-neutral-600 leading-relaxed">
                  В течение нескольких минут система проверяет наличие roadmap и содержание проекта. Если всё в порядке — вы получите уведомление о прохождении первого этапа и ожидании модератора.
                </p>
                <div class="mt-3 flex flex-wrap gap-2">
                  <span class="text-xs px-3 py-1 bg-green-100 text-green-700 rounded-full">✓ Roadmap приложен</span>
                  <span class="text-xs px-3 py-1 bg-green-100 text-green-700 rounded-full">✓ Контент проверен</span>
                </div>
              </div>
            </div>

            <!-- Step 5 -->
            <div class="relative flex flex-col md:flex-row items-center gap-6 md:gap-0">
              <div class="md:w-1/2 md:pr-12 md:text-right order-2 md:order-1">
                <h3 class="text-xl font-bold mb-2">Проверка модератором</h3>
                <p class="text-neutral-600 leading-relaxed">
                  Модератор платформы изучает ваш проект и roadmap, затем назначает встречу в удобное время. На созвоне он задаёт вопросы о продукте, команде и планах монетизации. После встречи выносится решение об одобрении или отклонении.
                </p>
              </div>
              <div class="relative z-10 w-14 h-14 rounded-full bg-blue-600 flex items-center justify-center shadow-lg order-1 md:order-2 shrink-0">
                <ShieldCheck class="w-6 h-6 text-white" />
              </div>
              <div class="md:w-1/2 md:pl-12 order-3" />
            </div>

            <!-- Step 6 -->
            <div class="relative flex flex-col md:flex-row items-center gap-6 md:gap-0">
              <div class="md:w-1/2 md:pr-12 order-2 md:order-1" />
              <div class="relative z-10 w-14 h-14 rounded-full bg-green-600 flex items-center justify-center shadow-lg order-1 shrink-0">
                <Banknote class="w-6 h-6 text-white" />
              </div>
              <div class="md:w-1/2 md:pl-12 order-3">
                <h3 class="text-xl font-bold mb-2">Поэтапное финансирование</h3>
                <p class="text-neutral-600 leading-relaxed">
                  После одобрения проект публикуется и открывается для инвесторов. Собранные средства выплачиваются поэтапно — по мере выполнения этапов из roadmap. Завершили этап — получили часть финансирования. Это мотивирует команду и защищает инвесторов.
                </p>
                <div class="mt-3 flex flex-wrap gap-2">
                  <span class="text-xs px-3 py-1 bg-blue-100 text-blue-700 rounded-full">Этап 1 → Выплата</span>
                  <span class="text-xs px-3 py-1 bg-blue-100 text-blue-700 rounded-full">Этап 2 → Выплата</span>
                  <span class="text-xs px-3 py-1 bg-blue-100 text-blue-700 rounded-full">Этап N → Выплата</span>
                </div>
              </div>
            </div>

          </div>
        </div>
      </div>
    </section>

    <!-- For investors -->
    <section class="py-16 sm:py-20 px-4 sm:px-6 bg-neutral-50">
      <div class="max-w-5xl mx-auto">
        <h2 class="text-2xl sm:text-3xl font-bold text-center mb-4">Для инвесторов</h2>
        <p class="text-neutral-500 text-center mb-12">Каждый проект на платформе прошёл двойную проверку</p>
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-6">
          <div class="bg-white rounded-2xl p-6 border border-neutral-200 shadow-sm">
            <div class="w-12 h-12 bg-blue-100 rounded-xl flex items-center justify-center mb-4">
              <CheckCircle class="w-6 h-6 text-blue-600" />
            </div>
            <h3 class="font-bold mb-2">Проверенные проекты</h3>
            <p class="text-neutral-600 text-sm leading-relaxed">Каждый проект проходит AI-фильтрацию и личную встречу с модератором — случайных или мошеннических проектов не бывает</p>
          </div>
          <div class="bg-white rounded-2xl p-6 border border-neutral-200 shadow-sm">
            <div class="w-12 h-12 bg-purple-100 rounded-xl flex items-center justify-center mb-4">
              <FileText class="w-6 h-6 text-purple-600" />
            </div>
            <h3 class="font-bold mb-2">Прозрачный roadmap</h3>
            <p class="text-neutral-600 text-sm leading-relaxed">Вы видите план развития проекта до того, как вкладываете деньги. Roadmap — обязательный документ для каждого проекта</p>
          </div>
          <div class="bg-white rounded-2xl p-6 border border-neutral-200 shadow-sm">
            <div class="w-12 h-12 bg-green-100 rounded-xl flex items-center justify-center mb-4">
              <Banknote class="w-6 h-6 text-green-600" />
            </div>
            <h3 class="font-bold mb-2">Поэтапные выплаты</h3>
            <p class="text-neutral-600 text-sm leading-relaxed">Деньги выплачиваются команде только по мере достижения этапов roadmap — ваши вложения защищены от нецелевого использования</p>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA -->
    <section class="py-16 sm:py-20 px-4 sm:px-6 bg-gradient-to-br from-blue-600 to-purple-700 text-white">
      <div class="max-w-3xl mx-auto text-center">
        <Users class="w-12 h-12 mx-auto mb-6 opacity-80" />
        <h2 class="text-2xl sm:text-3xl font-bold mb-4">Готовы запустить проект?</h2>
        <p class="text-blue-100 mb-8 text-lg">Присоединяйтесь к платформе и начните путь к финансированию своей идеи</p>
        <Button
          class="bg-white text-blue-700 hover:bg-blue-50 font-semibold px-8 py-3 text-base"
          @click="router.push({ name: 'create-project' })"
        >
          Создать проект
          <ArrowRight class="w-4 h-4 ml-2" />
        </Button>
      </div>
    </section>

  </div>
</template>
