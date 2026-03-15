import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/LoginView.vue'),
    meta: { public: true }
  },
  {
    path: '/',
    component: () => import('@/views/AppShell.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('@/views/DashboardView.vue')
      },
      {
        path: 'tournaments/new',
        meta: {requiresAdmin: true},
        name: 'NewTournament',
        component: () => import('@/views/TournamentFormView.vue')
      },
      {
        path: 'tournaments/join',
        name: 'JoinTournament',
        component: () => import('@/views/JoinTournamentView.vue')
      },
      {
        path: 'tournaments/:id',
        name: 'Tournament',
        component: () => import('@/views/TournamentView.vue')
      },
      {
        path: 'tournaments/:id/manage',
        meta: {requiresAdmin: true},
        name: 'ManageTournament',
        component: () => import('@/views/ManageTournamentView.vue')
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/ProfileView.vue')
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior: () => ({ top: 0 })
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return { name: 'Login', query: { redirect: to.fullPath } }
  }
  if (to.meta.requiresAdmin && !auth.isAdmin) {
    return { name: 'Dashboard' }
  }
  if (to.name === 'Login' && auth.isAuthenticated) {
    return { name: 'Dashboard' }
  }
})

export default router
