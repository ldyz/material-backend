import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useUserStore } from '@/stores/user'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Auth/Login.vue'),
    meta: { title: '登录', noAuth: true },
  },
  {
    path: '/',
    component: () => import('@/layouts/TabbarLayout.vue'),
    children: [
      {
        path: '',
        name: 'Home',
        component: () => import('@/views/Home/index.vue'),
        meta: { title: '首页', keepAlive: true },
      },
      {
        path: 'inbound',
        name: 'InboundList',
        component: () => import('@/views/Inbound/index.vue'),
        meta: { title: '入库', keepAlive: true, permission: 'inbound_view' },
      },
      {
        path: 'inbound/:id',
        name: 'InboundDetail',
        component: () => import('@/views/Inbound/Detail.vue'),
        meta: { title: '入库详情', permission: 'inbound_view' },
      },
      {
        path: 'inbound/:id/approve',
        name: 'InboundApprove',
        component: () => import('@/views/Inbound/Approve.vue'),
        meta: { title: '入库审批', permission: 'inbound_approve' },
      },
      {
        path: 'inbound/create',
        name: 'InboundCreate',
        component: () => import('@/views/Inbound/Create.vue'),
        meta: { title: '新建入库单', permission: 'inbound_create' },
      },
      {
        path: 'outbound',
        name: 'OutboundList',
        component: () => import('@/views/Outbound/index.vue'),
        meta: { title: '出库', keepAlive: true, permission: 'requisition_view' },
      },
      {
        path: 'outbound/:id',
        name: 'OutboundDetail',
        component: () => import('@/views/Outbound/Detail.vue'),
        meta: { title: '出库详情', permission: 'requisition_view' },
      },
      {
        path: 'outbound/:id/approve',
        name: 'OutboundApprove',
        component: () => import('@/views/Outbound/Approve.vue'),
        meta: { title: '出库审批', permission: 'requisition_approve' },
      },
      {
        path: 'outbound/create',
        name: 'OutboundCreate',
        component: () => import('@/views/Outbound/Create.vue'),
        meta: { title: '新建出库单', permission: 'requisition_create' },
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/Profile/index.vue'),
        meta: { title: '我的', keepAlive: true },
      },
      {
        path: 'notifications',
        name: 'NotificationList',
        component: () => import('@/views/Notification/List.vue'),
        meta: { title: '通知中心', keepAlive: true },
      },
      // 施工日志
      {
        path: 'construction',
        name: 'ConstructionList',
        component: () => import('@/views/Construction/index.vue'),
        meta: { title: '施工日志', keepAlive: true, permission: 'construction_log_view' },
      },
      {
        path: 'construction/create',
        name: 'ConstructionCreate',
        component: () => import('@/views/Construction/Form.vue'),
        meta: { title: '新建日志', permission: 'construction_log_create' },
      },
      {
        path: 'construction/:id',
        name: 'ConstructionDetail',
        component: () => import('@/views/Construction/Detail.vue'),
        meta: { title: '日志详情', permission: 'construction_log_view' },
      },
      {
        path: 'construction/:id/edit',
        name: 'ConstructionEdit',
        component: () => import('@/views/Construction/Form.vue'),
        meta: { title: '编辑日志', permission: 'construction_log_create' },
      },
      // 库存
      {
        path: 'stock',
        name: 'StockList',
        component: () => import('@/views/Stock/index.vue'),
        meta: { title: '库存', keepAlive: true, permission: 'stock_view' },
      },
      {
        path: 'stock/:id',
        name: 'StockDetail',
        component: () => import('@/views/Stock/Detail.vue'),
        meta: { title: '库存详情', permission: 'stock_view' },
      },
      // AI分析
      {
        path: 'ai',
        name: 'AI',
        component: () => import('@/views/AI/index.vue'),
        meta: { title: 'AI分析', keepAlive: true, permission: 'ai_analyze' },
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    redirect: '/',
  },
]

const router = createRouter({
  history: createWebHistory('/mobile/'),  // 移动端部署在 /mobile/ 子目录
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    }
    return { top: 0 }
  },
})

// 全局前置守卫
router.beforeEach((to, from, next) => {
  document.title = to.meta.title
    ? `${to.meta.title} - ${import.meta.env.VITE_APP_TITLE || '材料管理'}`
    : import.meta.env.VITE_APP_TITLE || '材料管理'

  const authStore = useAuthStore()

  // 不需要认证的页面
  if (to.meta.noAuth) {
    if (authStore.isAuthenticated) {
      next('/')
    } else {
      next()
    }
    return
  }

  // 需要认证
  if (!authStore.isAuthenticated) {
    next({
      path: '/login',
      query: { redirect: to.fullPath },
    })
    return
  }

  // 权限检查
  if (to.meta.permission) {
    const userStore = useUserStore()
    if (!userStore.hasPermission(to.meta.permission)) {
      next('/profile')
      return
    }
  }

  next()
})

export default router
