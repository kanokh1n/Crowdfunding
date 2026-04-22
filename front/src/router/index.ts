import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('@/pages/HomePage.vue'),
    },
    {
      path: '/auth',
      name: 'auth',
      component: () => import('@/pages/AuthPage.vue'),
      meta: { guestOnly: true },
    },
    {
      path: '/create-project',
      name: 'create-project',
      component: () => import('@/pages/CreateProjectPage.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/project/:id/edit',
      name: 'edit-project',
      component: () => import('@/pages/EditProjectPage.vue'),
      props: true,
      meta: { requiresAuth: true },
    },
    {
      path: '/project/:id',
      name: 'project-detail',
      component: () => import('@/pages/ProjectDetailPage.vue'),
      props: true,
    },
    {
      path: '/notifications',
      name: 'notifications',
      component: () => import('@/pages/NotificationsPage.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/about',
      name: 'about',
      component: () => import('@/pages/AboutPage.vue'),
    },
    {
      path: '/how-it-works',
      name: 'how-it-works',
      component: () => import('@/pages/HowItWorksPage.vue'),
    },
    {
      path: '/moderation',
      name: 'moderation',
      component: () => import('@/pages/ModerationPage.vue'),
      meta: { requiresAuth: true, requiresAdmin: true },
    },
  ],
})

router.beforeEach(async (to) => {
  const auth = useAuthStore()

  // Load user on first navigation if not loaded yet
  if (!auth.user && localStorage.getItem('access_token')) {
    await auth.fetchMe()
  }

  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return { name: 'auth', query: { redirect: to.fullPath } }
  }

  if (to.meta.requiresAdmin && !auth.isAdmin) {
    return { name: 'home' }
  }

  if (to.meta.guestOnly && auth.isAuthenticated) {
    return { name: 'home' }
  }
})

export default router
