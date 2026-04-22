<script setup lang="ts">
import { useRouter } from 'vue-router'
import { Facebook, Twitter, Instagram, Linkedin, Github, Mail, MapPin, Phone } from 'lucide-vue-next'

const router = useRouter()

function navigateTo(name: string) {
  router.push({ name })
}
</script>

<template>
  <footer class="bg-neutral-900 text-white mt-auto">
    <div class="max-w-7xl mx-auto px-6 py-12">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-8">
        <!-- О компании -->
        <div>
          <h3 class="mb-4">IT Crowdfunding</h3>
          <p class="text-neutral-400 mb-4">
            Платформа для финансирования инновационных IT-проектов
          </p>
          <div class="flex gap-3">
            <a href="#" class="w-10 h-10 bg-neutral-800 rounded-full flex items-center justify-center hover:bg-neutral-700 transition-colors">
              <Facebook class="w-5 h-5" />
            </a>
            <a href="#" class="w-10 h-10 bg-neutral-800 rounded-full flex items-center justify-center hover:bg-neutral-700 transition-colors">
              <Twitter class="w-5 h-5" />
            </a>
            <a href="#" class="w-10 h-10 bg-neutral-800 rounded-full flex items-center justify-center hover:bg-neutral-700 transition-colors">
              <Instagram class="w-5 h-5" />
            </a>
            <a href="#" class="w-10 h-10 bg-neutral-800 rounded-full flex items-center justify-center hover:bg-neutral-700 transition-colors">
              <Linkedin class="w-5 h-5" />
            </a>
            <a href="#" class="w-10 h-10 bg-neutral-800 rounded-full flex items-center justify-center hover:bg-neutral-700 transition-colors">
              <Github class="w-5 h-5" />
            </a>
          </div>
        </div>

        <!-- Навигация -->
        <div>
          <h4 class="mb-4">Навигация</h4>
          <ul class="space-y-2 text-neutral-400">
            <li>
              <button @click="navigateTo('home')" class="hover:text-white transition-colors">
                Главная
              </button>
            </li>
            <li>
              <button @click="navigateTo('about')" class="hover:text-white transition-colors">
                О платформе
              </button>
            </li>
            <li><a href="#" class="hover:text-white transition-colors">Проекты</a></li>
            <li>
              <button @click="navigateTo('how-it-works')" class="hover:text-white transition-colors">
                Как это работает
              </button>
            </li>
            <li><a href="#" class="hover:text-white transition-colors">FAQ</a></li>
          </ul>
        </div>

        <!-- Для создателей -->
        <div>
          <h4 class="mb-4">Для создателей</h4>
          <ul class="space-y-2 text-neutral-400">
            <li>
              <button @click="navigateTo('create-project')" class="hover:text-white transition-colors">
                Создать проект
              </button>
            </li>
            <li><a href="#" class="hover:text-white transition-colors">Руководство</a></li>
            <li><a href="#" class="hover:text-white transition-colors">Правила</a></li>
            <li><a href="#" class="hover:text-white transition-colors">Поддержка</a></li>
          </ul>
        </div>

        <!-- Контакты -->
        <div>
          <h4 class="mb-4">Контакты</h4>
          <ul class="space-y-3 text-neutral-400">
            <li class="flex items-start gap-2">
              <MapPin class="w-5 h-5 mt-0.5 flex-shrink-0" />
              <span>г. Москва, ул. Тверская, д. 1</span>
            </li>
            <li class="flex items-center gap-2">
              <Phone class="w-5 h-5 flex-shrink-0" />
              <a href="tel:+74951234567" class="hover:text-white transition-colors">
                +7 (495) 123-45-67
              </a>
            </li>
            <li class="flex items-center gap-2">
              <Mail class="w-5 h-5 flex-shrink-0" />
              <a href="mailto:info@itcrowdfunding.ru" class="hover:text-white transition-colors">
                info@itcrowdfunding.ru
              </a>
            </li>
          </ul>
        </div>
      </div>

      <div class="border-t border-neutral-800 mt-8 pt-8 text-center text-neutral-400">
        <p>&copy; 2026 IT Crowdfunding. Все права защищены.</p>
      </div>
    </div>
  </footer>
</template>
