import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue')
  },
  {
    path: '/',
    component: () => import('@/layouts/Layout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue')
      },
      {
        path: 'apps',
        name: 'Apps',
        component: () => import('@/views/app/index.vue')
      },
      {
        path: 'apps/:id/config',
        name: 'AppConfig',
        component: () => import('@/views/app/config/index.vue'),
        children: [
          { path: '', redirect: 'dashboard' },
          { path: 'dashboard', name: 'AppConfigDashboard', component: () => import('@/views/app/config/Dashboard.vue') },
          { path: 'user', name: 'UserConfig', component: () => import('@/views/app/config/UserConfig.vue') },
          { path: 'message', name: 'MessageConfig', component: () => import('@/views/app/config/MessageConfig.vue') },
          { path: 'payment', name: 'PaymentConfig', component: () => import('@/views/app/config/PaymentConfig.vue') },
          { path: 'analytics', name: 'AnalyticsConfig', component: () => import('@/views/app/config/AnalyticsConfig.vue') },
          { path: 'security', name: 'SecurityConfig', component: () => import('@/views/app/config/SecurityConfig.vue') },
          { path: 'version', name: 'VersionConfig', component: () => import('@/views/app/config/VersionConfig.vue') },
          { path: 'modules', name: 'ModulesConfig', component: () => import('@/views/app/config/ModulesConfig.vue') }
        ]
      },
      {
        path: 'modules',
        name: 'Modules',
        component: () => import('@/views/module/index.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.path !== '/login' && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router
