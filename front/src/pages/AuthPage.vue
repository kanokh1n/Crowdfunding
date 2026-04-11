<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import Label from '@/components/ui/Label.vue'
import Tabs from '@/components/ui/Tabs.vue'
import TabsList from '@/components/ui/TabsList.vue'
import TabsTrigger from '@/components/ui/TabsTrigger.vue'
import TabsContent from '@/components/ui/TabsContent.vue'
import { Lock, Mail, User } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

// Login form
const loginEmail = ref('')
const loginPassword = ref('')
const loginError = ref('')

// Register form
const registerName = ref('')
const registerEmail = ref('')
const registerPassword = ref('')
const registerError = ref('')
const agreeToTerms = ref(false)

const isLoading = ref(false)

async function handleLogin() {
  loginError.value = ''
  isLoading.value = true
  try {
    await auth.login(loginEmail.value, loginPassword.value)
    const redirect = route.query.redirect as string
    router.push(redirect || { name: 'home' })
  } catch (err: any) {
    loginError.value = err.message || 'Ошибка входа. Проверьте данные.'
  } finally {
    isLoading.value = false
  }
}

async function handleRegister() {
  registerError.value = ''
  if (!agreeToTerms.value) {
    registerError.value = 'Необходимо согласиться с условиями использования'
    return
  }
  if (registerPassword.value.length < 6) {
    registerError.value = 'Пароль должен содержать минимум 6 символов'
    return
  }
  isLoading.value = true
  try {
    await auth.register(registerEmail.value, registerPassword.value, registerName.value)
    const redirect = route.query.redirect as string
    router.push(redirect || { name: 'home' })
  } catch (err: any) {
    registerError.value = err.message || 'Ошибка регистрации. Проверьте данные.'
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="min-h-[calc(100vh-200px)] flex items-center justify-center py-12 px-4 sm:px-6">
    <div class="w-full max-w-md">
      <div class="text-center mb-6 sm:mb-8">
        <h1 class="mb-2">Добро пожаловать</h1>
        <p class="text-neutral-600">Войдите или создайте новый аккаунт</p>
      </div>

      <div class="bg-white rounded-xl shadow-lg p-6 sm:p-8 border border-neutral-200">
        <Tabs default-value="login">
          <TabsList class="grid w-full grid-cols-2 mb-6">
            <TabsTrigger value="login">Вход</TabsTrigger>
            <TabsTrigger value="register">Регистрация</TabsTrigger>
          </TabsList>

          <TabsContent value="login">
            <form @submit.prevent="handleLogin" class="space-y-4">
              <div class="space-y-2">
                <Label for="login-email">Email</Label>
                <div class="relative">
                  <Mail class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-neutral-400" />
                  <Input
                    id="login-email"
                    type="email"
                    placeholder="your@email.com"
                    v-model="loginEmail"
                    class="pl-10"
                    required
                  />
                </div>
              </div>

              <div class="space-y-2">
                <Label for="login-password">Пароль</Label>
                <div class="relative">
                  <Lock class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-neutral-400" />
                  <Input
                    id="login-password"
                    type="password"
                    placeholder="••••••••"
                    v-model="loginPassword"
                    class="pl-10"
                    required
                  />
                </div>
              </div>

              <div v-if="loginError" class="text-red-500 text-sm text-center">
                {{ loginError }}
              </div>

              <div class="flex items-center justify-between">
                <label class="flex items-center gap-2 text-neutral-600 cursor-pointer text-sm">
                  <input type="checkbox" class="rounded" />
                  <span>Запомнить меня</span>
                </label>
                <a href="#" class="text-blue-600 hover:text-blue-700 text-sm">
                  Забыли пароль?
                </a>
              </div>

              <Button type="submit" class="w-full" :disabled="isLoading">
                {{ isLoading ? 'Вход...' : 'Войти' }}
              </Button>
            </form>
          </TabsContent>

          <TabsContent value="register">
            <form @submit.prevent="handleRegister" class="space-y-4">
              <div class="space-y-2">
                <Label for="register-name">Имя</Label>
                <div class="relative">
                  <User class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-neutral-400" />
                  <Input
                    id="register-name"
                    type="text"
                    placeholder="Иван Иванов"
                    v-model="registerName"
                    class="pl-10"
                    required
                  />
                </div>
              </div>

              <div class="space-y-2">
                <Label for="register-email">Email</Label>
                <div class="relative">
                  <Mail class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-neutral-400" />
                  <Input
                    id="register-email"
                    type="email"
                    placeholder="your@email.com"
                    v-model="registerEmail"
                    class="pl-10"
                    required
                  />
                </div>
              </div>

              <div class="space-y-2">
                <Label for="register-password">Пароль</Label>
                <div class="relative">
                  <Lock class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-neutral-400" />
                  <Input
                    id="register-password"
                    type="password"
                    placeholder="••••••••"
                    v-model="registerPassword"
                    class="pl-10"
                    required
                    minlength="6"
                  />
                </div>
              </div>

              <div v-if="registerError" class="text-red-500 text-sm text-center">
                {{ registerError }}
              </div>

              <div class="text-neutral-600 text-sm">
                <label class="flex items-start gap-2 cursor-pointer">
                  <input v-model="agreeToTerms" type="checkbox" class="mt-1 rounded" required />
                  <span>
                    Я согласен с
                    <a href="#" class="text-blue-600 hover:text-blue-700">условиями использования</a>
                    и
                    <a href="#" class="text-blue-600 hover:text-blue-700">политикой конфиденциальности</a>
                  </span>
                </label>
              </div>

              <Button type="submit" class="w-full" :disabled="isLoading">
                {{ isLoading ? 'Регистрация...' : 'Зарегистрироваться' }}
              </Button>
            </form>
          </TabsContent>
        </Tabs>
      </div>
    </div>
  </div>
</template>
