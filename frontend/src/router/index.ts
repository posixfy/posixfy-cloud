import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/LoginView.vue'),
      meta: { guest: true },
    },
    {
      path: '/',
      name: 'Files',
      component: () => import('@/views/FileBrowserView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/admin',
      name: 'Admin',
      component: () => import('@/views/AdminView.vue'),
      meta: { requiresAuth: true, requiresAdmin: true },
    },
  ],
})

router.beforeEach((to) => {
  const auth = useAuthStore()

  if (to.meta.requiresAuth && !auth.isLoggedIn) {
    return '/login'
  }
  if (to.meta.guest && auth.isLoggedIn) {
    return auth.isAdmin ? '/admin' : '/'
  }
  if (to.meta.requiresAdmin && !auth.isAdmin) {
    return '/'
  }
  // admin cannot access file browser
  if (to.path === '/' && auth.isLoggedIn && auth.isAdmin) {
    return '/admin'
  }
})

export default router
