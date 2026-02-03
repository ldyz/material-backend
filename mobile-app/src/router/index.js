/**
 * Vue Router 配置 (移动端 - 优化版)
 *
 * 包含：
 * - 路由懒加载（按页面分割代码）
 * - 路由分组（按功能模块）
 * - 权限控制
 * - 页面缓存控制
 * - 滚动行为
 * - 过渡动画支持
 */

import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useUserStore } from '@/stores/user'

// ==================== 路由分组（Webpack魔术注释） ====================

/**
 * 认证相关路由
 */
const AuthRoutes = () => import(/* webpackChunkName: "auth" */ '@/views/Auth/Login.vue')

/**
 * 首页相关路由
 */
const HomeRoutes = () => import(/* webpackChunkName: "home" */ '@/views/Home/index.vue')

/**
 * 物资计划相关路由
 */
const PlanList = () => import(/* webpackChunkName: "plans" */ '@/views/Plans/index.vue')
const PlanDetail = () => import(/* webpackChunkName: "plans" */ '@/views/Plans/Detail.vue')
const PlanApprove = () => import(/* webpackChunkName: "plans" */ '@/views/Plans/Approve.vue')

/**
 * 入库相关路由
 */
const InboundList = () => import(/* webpackChunkName: "inbound" */ '@/views/Inbound/index.vue')
const InboundDetail = () => import(/* webpackChunkName: "inbound" */ '@/views/Inbound/Detail.vue')
const InboundApprove = () => import(/* webpackChunkName: "inbound" */ '@/views/Inbound/Approve.vue')
const InboundCreate = () => import(/* webpackChunkName: "inbound" */ '@/views/Inbound/Create.vue')

/**
 * 出库相关路由
 */
const OutboundList = () => import(/* webpackChunkName: "outbound" */ '@/views/Outbound/index.vue')
const OutboundDetail = () => import(/* webpackChunkName: "outbound" */ '@/views/Outbound/Detail.vue')
const OutboundApprove = () => import(/* webpackChunkName: "outbound" */ '@/views/Outbound/Approve.vue')
const OutboundIssue = () => import(/* webpackChunkName: "outbound" */ '@/views/Outbound/Issue.vue')
const OutboundCreate = () => import(/* webpackChunkName: "outbound" */ '@/views/Outbound/Create.vue')

/**
 * 物资相关路由
 */
const MaterialList = () => import(/* webpackChunkName: "materials" */ '@/views/Materials/index.vue')
const MaterialDetail = () => import(/* webpackChunkName: "materials" */ '@/views/Materials/Detail.vue')

/**
 * 库存相关路由
 */
const StockList = () => import(/* webpackChunkName: "stock" */ '@/views/Stock/index.vue')

/**
 * 工作流相关路由
 */
const TaskList = () => import(/* webpackChunkName: "workflow" */ '@/views/Tasks/index.vue')

/**
 * 用户相关路由
 */
const Profile = () => import(/* webpackChunkName: "user" */ '@/views/Profile/index.vue')
const NotificationList = () => import(/* webpackChunkName: "user" */ '@/views/Notification/List.vue')

/**
 * 主布局
 */
const TabbarLayout = () => import(/* webpackChunkName: "layout" */ '@/layouts/TabbarLayout.vue')

// ==================== 路由配置 ====================

const routes = [
  // 登录页
  {
    path: '/login',
    name: 'Login',
    component: AuthRoutes,
    meta: {
      title: '登录',
      noAuth: true,
      transition: 'fade',
    },
  },

  // 主应用
  {
    path: '/',
    component: TabbarLayout,
    children: [
      // 首页
      {
        path: '',
        name: 'Home',
        component: HomeRoutes,
        meta: {
          title: '首页',
          keepAlive: true,
          transition: 'slide',
        },
      },

      // ========== 物资计划模块 ==========
      {
        path: 'plans',
        name: 'PlanList',
        component: PlanList,
        meta: {
          title: '物资计划',
          keepAlive: true,
          permission: 'material_plan_view',
          transition: 'slide',
        },
      },
      {
        path: 'plans/:id',
        name: 'PlanDetail',
        component: PlanDetail,
        meta: {
          title: '计划详情',
          permission: 'material_plan_view',
          transition: 'slide-left',
        },
      },
      {
        path: 'plans/:id/approve',
        name: 'PlanApprove',
        component: PlanApprove,
        meta: {
          title: '计划审批',
          permission: 'material_plan_approve',
          transition: 'slide-up',
        },
      },

      // ========== 入库模块 ==========
      {
        path: 'inbound',
        name: 'InboundList',
        component: InboundList,
        meta: {
          title: '入库单',
          keepAlive: true,
          permission: 'inbound_view',
          transition: 'slide',
        },
      },
      {
        path: 'inbound/:id',
        name: 'InboundDetail',
        component: InboundDetail,
        meta: {
          title: '入库详情',
          permission: 'inbound_view',
          transition: 'slide-left',
        },
      },
      {
        path: 'inbound/:id/approve',
        name: 'InboundApprove',
        component: InboundApprove,
        meta: {
          title: '入库审批',
          permission: 'inbound_approve',
          transition: 'slide-up',
        },
      },
      {
        path: 'inbound/create',
        name: 'InboundCreate',
        component: InboundCreate,
        meta: {
          title: '新建入库单',
          permission: 'inbound_create',
          transition: 'slide-up',
        },
      },

      // ========== 出库模块 ==========
      {
        path: 'outbound',
        name: 'OutboundList',
        component: OutboundList,
        meta: {
          title: '出库单',
          keepAlive: true,
          permission: 'requisition_view',
          transition: 'slide',
        },
      },
      {
        path: 'outbound/:id',
        name: 'OutboundDetail',
        component: OutboundDetail,
        meta: {
          title: '出库详情',
          permission: 'requisition_view',
          transition: 'slide-left',
        },
      },
      {
        path: 'outbound/:id/approve',
        name: 'OutboundApprove',
        component: OutboundApprove,
        meta: {
          title: '出库审批',
          permission: 'requisition_approve',
          transition: 'slide-up',
        },
      },
      {
        path: 'outbound/:id/issue',
        name: 'OutboundIssue',
        component: OutboundIssue,
        meta: {
          title: '出库发放',
          permission: 'requisition_issue',
          transition: 'slide-up',
        },
      },
      {
        path: 'outbound/create',
        name: 'OutboundCreate',
        component: OutboundCreate,
        meta: {
          title: '新建出库单',
          permission: 'requisition_create',
          transition: 'slide-up',
        },
      },

      // ========== 物资浏览模块 ==========
      {
        path: 'materials',
        name: 'MaterialList',
        component: MaterialList,
        meta: {
          title: '物资浏览',
          keepAlive: true,
          permission: 'material_view',
          transition: 'slide',
        },
      },
      {
        path: 'materials/:id',
        name: 'MaterialDetail',
        component: MaterialDetail,
        meta: {
          title: '物资详情',
          permission: 'material_view',
          transition: 'slide-left',
        },
      },

      // ========== 库存模块 ==========
      {
        path: 'stock',
        name: 'StockList',
        component: StockList,
        meta: {
          title: '库存',
          keepAlive: true,
          permission: 'stock_view',
          transition: 'slide',
        },
      },

      // ========== 工作流模块 ==========
      {
        path: 'tasks',
        name: 'TaskList',
        component: TaskList,
        meta: {
          title: '待办任务',
          keepAlive: true,
          transition: 'slide',
        },
      },

      // ========== 用户模块 ==========
      {
        path: 'profile',
        name: 'Profile',
        component: Profile,
        meta: {
          title: '我的',
          keepAlive: true,
          transition: 'slide',
        },
      },
      {
        path: 'notifications',
        name: 'NotificationList',
        component: NotificationList,
        meta: {
          title: '通知中心',
          keepAlive: true,
          transition: 'slide',
        },
      },
    ],
  },

  // 404页面
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    redirect: '/',
  },
]

// ==================== 创建路由实例 ====================

const router = createRouter({
  history: createWebHistory('/mobile/'),
  routes,
  scrollBehavior(to, from, savedPosition) {
    // 如果有保存的位置（浏览器前进/后退）
    if (savedPosition) {
      return savedPosition
    }

    // 如果是详情页，保持滚动位置
    if (to.meta.keepAlive && from.meta.keepAlive) {
      return false
    }

    // 默认滚动到顶部
    return { top: 0, behavior: 'smooth' }
  },
})

// ==================== 路由守卫 ====================

/**
 * 全局前置守卫
 */
router.beforeEach((to, from, next) => {
  // 设置页面标题
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
      // 显示无权限提示
      next({ name: 'Profile' })
      return
    }
  }

  next()
})

/**
 * 全局后置钩子（用于页面访问统计等）
 */
router.afterEach((to, from) => {
  // 可以在这里添加页面访问统计
  if (import.meta.env.NODE_ENV === 'development') {
    console.log(`[Router] ${from.path} -> ${to.path}`)
  }
})

/**
 * 全局错误处理
 */
router.onError((error) => {
  console.error('[Router Error]', error)
  // 可以在这里添加错误上报
})

// ==================== 导出 ====================

export default router

/**
 * 路由过渡类型
 * @type {Object}
 */
export const TransitionTypes = {
  FADE: 'fade',
  SLIDE: 'slide',
  SLIDE_LEFT: 'slide-left',
  SLIDE_RIGHT: 'slide-right',
  SLIDE_UP: 'slide-up',
  SLIDE_DOWN: 'slide-down',
}

/**
 * 路由过渡名称
 * @param {Object} to - 目标路由
 * @param {Object} from - 来源路由
 * @returns {string} 过渡名称
 */
export function getTransitionName(to, from) {
  const toDepth = to.path.split('/').length
  const fromDepth = from.path.split('/').length

  if (toDepth > fromDepth) {
    return 'slide-left'
  } else if (toDepth < fromDepth) {
    return 'slide-right'
  } else {
    return to.meta.transition || 'fade'
  }
}
