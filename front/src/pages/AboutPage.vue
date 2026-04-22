<script setup lang="ts">
import { useRouter } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import { ShieldCheck, FileText, Banknote, Zap, Users, CheckCircle, ArrowRight, Rocket } from 'lucide-vue-next'

const router = useRouter()
</script>

<template>
  <div class="min-h-screen">

    <!-- Hero / Миссия -->
    <section class="bg-gradient-to-br from-blue-600 via-blue-700 to-purple-700 text-white py-20 px-4 sm:px-6">
      <div class="max-w-4xl mx-auto text-center">
        <div class="inline-flex items-center gap-2 bg-white/20 rounded-full px-4 py-2 mb-6 text-sm">
          <Rocket class="w-4 h-4" />
          О платформе
        </div>
        <h1 class="text-3xl sm:text-5xl font-bold mb-6 leading-tight">
          Финансирование IT-проектов,<br />которым можно доверять
        </h1>
        <p class="text-lg sm:text-xl text-blue-100 max-w-2xl mx-auto leading-relaxed">
          Мы создали платформу, где каждый проект проходит реальную проверку, а инвесторы защищены поэтапной системой выплат. Только работающие команды — только IT.
        </p>
      </div>
    </section>

    <!-- Для инвесторов -->
    <section class="py-16 sm:py-24 px-4 sm:px-6 bg-white">
      <div class="max-w-5xl mx-auto">
        <div class="text-center mb-12">
          <h2 class="text-2xl sm:text-3xl font-bold mb-4">Как мы защищаем инвесторов</h2>
          <p class="text-neutral-500 max-w-xl mx-auto">Прежде чем вы вложите рубль, проект уже прошёл несколько уровней проверки</p>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
          <div class="flex gap-4 p-6 bg-neutral-50 rounded-2xl border border-neutral-100">
            <div class="w-12 h-12 bg-blue-100 rounded-xl flex items-center justify-center shrink-0">
              <Zap class="w-6 h-6 text-blue-600" />
            </div>
            <div>
              <h3 class="font-bold mb-1">AI-фильтрация</h3>
              <p class="text-neutral-600 text-sm leading-relaxed">Автоматическая проверка содержания проекта и наличия roadmap — случайные или некачественные заявки отсеиваются сразу</p>
            </div>
          </div>
          <div class="flex gap-4 p-6 bg-neutral-50 rounded-2xl border border-neutral-100">
            <div class="w-12 h-12 bg-purple-100 rounded-xl flex items-center justify-center shrink-0">
              <ShieldCheck class="w-6 h-6 text-purple-600" />
            </div>
            <div>
              <h3 class="font-bold mb-1">Встреча с модератором</h3>
              <p class="text-neutral-600 text-sm leading-relaxed">Каждый проект проходит личный созвон с модератором платформы — проверяется команда, продукт и план монетизации</p>
            </div>
          </div>
          <div class="flex gap-4 p-6 bg-neutral-50 rounded-2xl border border-neutral-100">
            <div class="w-12 h-12 bg-green-100 rounded-xl flex items-center justify-center shrink-0">
              <FileText class="w-6 h-6 text-green-600" />
            </div>
            <div>
              <h3 class="font-bold mb-1">Обязательный roadmap</h3>
              <p class="text-neutral-600 text-sm leading-relaxed">Каждый проект публикуется с roadmap — вы видите конкретный план развития до того, как принимаете решение об инвестиции</p>
            </div>
          </div>
          <div class="flex gap-4 p-6 bg-neutral-50 rounded-2xl border border-neutral-100">
            <div class="w-12 h-12 bg-orange-100 rounded-xl flex items-center justify-center shrink-0">
              <Banknote class="w-6 h-6 text-orange-600" />
            </div>
            <div>
              <h3 class="font-bold mb-1">Поэтапные выплаты</h3>
              <p class="text-neutral-600 text-sm leading-relaxed">Команда получает деньги только по мере выполнения этапов из roadmap — ваши вложения защищены от нецелевого использования</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Для стартаперов -->
    <section class="py-16 sm:py-24 px-4 sm:px-6 bg-neutral-50">
      <div class="max-w-5xl mx-auto">
        <div class="text-center mb-12">
          <h2 class="text-2xl sm:text-3xl font-bold mb-4">Как мы помогаем стартаперам</h2>
          <p class="text-neutral-500 max-w-xl mx-auto">Мы не просто площадка — мы партнёр, который помогает дойти до финансирования</p>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
          <div class="flex gap-4 p-6 bg-white rounded-2xl border border-neutral-200 shadow-sm">
            <div class="w-12 h-12 bg-blue-100 rounded-xl flex items-center justify-center shrink-0">
              <Zap class="w-6 h-6 text-blue-600" />
            </div>
            <div>
              <h3 class="font-bold mb-1">Быстрая проверка</h3>
              <p class="text-neutral-600 text-sm leading-relaxed">AI-проверка занимает несколько минут. Вы сразу узнаёте, прошёл ли проект первый этап, и получаете уведомление</p>
            </div>
          </div>
          <div class="flex gap-4 p-6 bg-white rounded-2xl border border-neutral-200 shadow-sm">
            <div class="w-12 h-12 bg-purple-100 rounded-xl flex items-center justify-center shrink-0">
              <Users class="w-6 h-6 text-purple-600" />
            </div>
            <div>
              <h3 class="font-bold mb-1">Живой диалог</h3>
              <p class="text-neutral-600 text-sm leading-relaxed">Модератор назначает встречу в удобное время, задаёт вопросы и помогает понять, что улучшить в презентации проекта</p>
            </div>
          </div>
          <div class="flex gap-4 p-6 bg-white rounded-2xl border border-neutral-200 shadow-sm">
            <div class="w-12 h-12 bg-green-100 rounded-xl flex items-center justify-center shrink-0">
              <CheckCircle class="w-6 h-6 text-green-600" />
            </div>
            <div>
              <h3 class="font-bold mb-1">Целевая аудитория</h3>
              <p class="text-neutral-600 text-sm leading-relaxed">На платформе собраны инвесторы, которые целенаправленно ищут IT-проекты — не случайные люди, а заинтересованные вкладчики</p>
            </div>
          </div>
          <div class="flex gap-4 p-6 bg-white rounded-2xl border border-neutral-200 shadow-sm">
            <div class="w-12 h-12 bg-orange-100 rounded-xl flex items-center justify-center shrink-0">
              <Banknote class="w-6 h-6 text-orange-600" />
            </div>
            <div>
              <h3 class="font-bold mb-1">Прозрачное финансирование</h3>
              <p class="text-neutral-600 text-sm leading-relaxed">Выплаты привязаны к этапам roadmap — это дисциплинирует команду и одновременно демонстрирует инвесторам прогресс</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA -->
    <section class="py-16 sm:py-20 px-4 sm:px-6 bg-gradient-to-br from-blue-600 to-purple-700 text-white">
      <div class="max-w-3xl mx-auto text-center">
        <Rocket class="w-12 h-12 mx-auto mb-6 opacity-80" />
        <h2 class="text-2xl sm:text-3xl font-bold mb-4">Готовы запустить проект?</h2>
        <p class="text-blue-100 mb-8 text-lg">Присоединяйтесь к платформе и начните путь к финансированию своей идеи</p>
        <Button
          class="bg-white text-blue-700 hover:bg-blue-50 font-semibold px-8 py-3 text-base"
          @click="router.push({ name: 'create-project' })"
        >
          Запустить проект
          <ArrowRight class="w-4 h-4 ml-2" />
        </Button>
      </div>
    </section>

  </div>
</template>
