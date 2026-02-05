import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login/index.vue'),
    meta: { noAuth: true }
  },
  {
    path: '/',
    component: () => import('@/layouts/TabbarLayout.vue'),
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard/index.vue'),
        meta: { keepAlive: true }
      },
      {
        path: 'plans',
        name: 'PlanList',
        component: () => import('@/views/Plan/List.vue'),
        meta: { keepAlive: true }
      },
      {
        path: 'plans/:id',
        name: 'PlanDetail',
        component: () => import('@/views/Plan/Detail.vue')
      },
      {
        path: 'plans/:id/approve',
        name: 'PlanApprove',
        component: () => import('@/views/Plan/Approve.vue')
      },
      {
        path: 'inbound',
        name: 'InboundList',
        component: () => import('@/views/Inbound/List.vue'),
        meta: { keepAlive: true }
      },
      {
        path: 'inbound/:id',
        name: 'InboundDetail',
        component: () => import('@/views/Inbound/Detail.vue')
      },
      {
        path: 'inbound/:id/approve',
        name: 'InboundApprove',
        component: () => import('@/views/Inbound/Approve.vue')
      },
      {
        path: 'inbound/create',
        name: 'InboundCreate',
        component: () => import('@/views/Inbound/Create.vue')
      },
      {
        path: 'requisition',
        name: 'RequisitionList',
        component: () => import('@/views/Requisition/List.vue'),
        meta: { keepAlive: true }
      },
      {
        path: 'requisition/:id',
        name: 'RequisitionDetail',
        component: () => import('@/views/Requisition/Detail.vue')
      },
      {
        path: 'requisition/:id/approve',
        name: 'RequisitionApprove',
        component: () => import('@/views/Requisition/Approve.vue')
      },
      {
        path: 'requisition/:id/issue',
        name: 'RequisitionIssue',
        component: () => import('@/views/Requisition/Issue.vue')
      },
      {
        path: 'requisition/create',
        name: 'RequisitionCreate',
        component: () => import('@/views/Requisition/Create.vue')
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/Profile/index.vue'),
        meta: { keepAlive: true }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/'
  }
]

const router = createRouter({
  history: createWebHistory('/mobile/'),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})

router.beforeEach((to, from, next) => {
  document.title = to.meta.title ? `${to.meta.title} - 材料管理` : '材料管理'

  const authStore = useAuthStore()

  if (to.meta.noAuth) {
    if (authStore.isAuthenticated) {
      next('/')
    } else {
      next()
    }
    return
  }

  if (!authStore.isAuthenticated) {
    next({ path: '/login', query: { redirect: to.fullPath } })
    return
  }

  next()
})

export default router
