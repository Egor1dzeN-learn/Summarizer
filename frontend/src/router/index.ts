import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import { useAuthStore } from '../stores/auth'
import BlankView from '../views/home/BlankView.vue'
import ChatView from '../views/home/ChatView.vue'
import LoginView from '../views/LoginView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: LoginView,
    },
    {
      path: '/',
      name: 'home',
      component: HomeView,
      children: [
        {
          path: '',
          name: 'blank',
          component: BlankView,
        },
        {
          path: ':id',
          name: 'chat',
          component: ChatView,
        },
      ],
    },
  ],
})

router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  if (!authStore.initialized) {
    await authStore.fetchUser()
  }

  if (to.name !== 'login' && !authStore.isAuthenticated) {
    next({ name: 'login' })
    return
  }

  next()
})

export default router
