/**
 * Vue Router 配置文件
 *
 * 本文件定义了应用的所有路由及其访问控制规则：
 * - 路由路径与组件的映射关系
 * - 路由元信息（标题、权限要求）
 * - 路由守卫（认证、权限检查）
 *
 * @module Router
 * @author Material Management System
 * @date 2025-01-27
 */

import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'
import { hasAnyPermission } from '@/utils/permissions'

/**
 * 路由配置数组
 *
 * 路由结构说明：
 * - path: URL 路径
 * - name: 路由名称（用于编程式导航）
 * - component: 路由对应的 Vue 组件（懒加载）
 * - meta: 路由元信息
 *   - requiresAuth: 是否需要登录认证
 *   - title: 页面标题
 *   - permissions: 访问该页面需要的权限列表（满足其一即可）
 */
const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { requiresAuth: false, title: '用户登录' }
  },
  {
    path: '/',
    component: () => import('@/components/layout/MainLayout.vue'),
    meta: { requiresAuth: true },
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: {
          requiresAuth: true,
          title: '仪表板'
          // 移除权限要求，所有登录用户都可以访问
        }
      },
      {
        path: 'projects',
        name: 'Projects',
        component: () => import('@/views/Projects.vue'),
        meta: {
          requiresAuth: true,
          title: '项目管理',
          permissions: ['project_view', 'project_create', 'project_edit', 'project_delete', 'project_member_manage']
        }
      },
      {
        path: 'construction-log',
        name: 'ConstructionLog',
        component: () => import('@/views/ConstructionLog.vue'),
        meta: {
          requiresAuth: true,
          title: '施工日志',
          permissions: ['constructionlog_view', 'constructionlog_create', 'constructionlog_edit', 'constructionlog_delete']
        }
      },
      {
        path: 'progress',
        name: 'Progress',
        component: () => import('@/views/Progress.vue'),
        meta: {
          requiresAuth: true,
          title: '进度管理',
          permissions: ['progress_view', 'progress_create', 'progress_edit', 'progress_delete']
        }
      },
      {
        path: 'materials',
        name: 'Materials',
        component: () => import('@/views/Materials.vue'),
        meta: {
          requiresAuth: true,
          title: '物资浏览',
          permissions: ['material_view', 'material_create', 'material_edit', 'material_delete', 'material_import', 'material_export', 'material_in']
        }
      },
      {
        path: 'material-categories',
        name: 'MaterialCategories',
        component: () => import('@/views/MaterialCategories.vue'),
        meta: {
          requiresAuth: true,
          title: '物资分类',
          permissions: ['material_view']
        }
      },
      {
        path: 'workflows',
        name: 'Workflows',
        component: () => import('@/views/WorkflowManagement.vue'),
        meta: {
          requiresAuth: true,
          title: '工作流管理',
          permissions: ['system_config']
        }
      },
      {
        path: 'material-plans',
        name: 'MaterialPlans',
        component: () => import('@/views/MaterialPlans.vue'),
        meta: {
          requiresAuth: true,
          title: '物资计划',
          permissions: ['material_plan_view']
        }
      },
      {
        path: 'plan-statistics',
        name: 'PlanStatistics',
        component: () => import('@/views/PlanStatistics.vue'),
        meta: {
          requiresAuth: true,
          title: '计划统计',
          permissions: ['material_plan_view']
        }
      },
      {
        path: 'stock',
        name: 'Stock',
        component: () => import('@/views/Stock.vue'),
        meta: {
          requiresAuth: true,
          title: '库存浏览',
          permissions: ['stock_view', 'stock_in', 'stock_out', 'stock_edit', 'stock_delete', 'stocklog_view', 'stocklog_delete', 'stock_export']
        }
      },
      {
        path: 'requisitions',
        name: 'Requisitions',
        component: () => import('@/views/Requisitions.vue'),
        meta: {
          requiresAuth: true,
          title: '出库管理',
          permissions: ['requisition_view']
        }
      },
      {
        path: 'inbound',
        name: 'Inbound',
        component: () => import('@/views/Inbound.vue'),
        meta: {
          requiresAuth: true,
          title: '入库管理',
          permissions: ['inbound_view']
        }
      },
      {
        path: 'warehouse',
        name: 'Warehouse',
        component: () => import('@/views/Warehouse.vue'),
        meta: {
          requiresAuth: true,
          title: '出入库管理',
          permissions: [] // 这可能是一个综合页面，暂时不需要特定权限
        }
      },
      {
        path: 'ai-analysis',
        name: 'AIAnalysis',
        component: () => import('@/views/AIAnalysis.vue'),
        meta: {
          requiresAuth: true,
          title: 'AI数据分析',
          permissions: ['system_statistics', 'system_report']
        }
      },
      {
        path: 'operation-logs',
        name: 'OperationLogs',
        component: () => import('@/views/OperationLogs.vue'),
        meta: {
          requiresAuth: true,
          title: '操作日志',
          permissions: ['audit_view']
        }
      },
      {
        path: 'system/users',
        name: 'UserManagement',
        component: () => import('@/views/UserManagement.vue'),
        meta: {
          requiresAuth: true,
          title: '用户管理',
          permissions: ['user_view']
        }
      },
      {
        path: 'system/roles',
        name: 'RoleManagement',
        component: () => import('@/views/RoleManagement.vue'),
        meta: {
          requiresAuth: true,
          title: '角色管理',
          permissions: ['role_view']
        }
      },
      {
        path: 'system',
        name: 'System',
        component: () => import('@/views/System.vue'),
        meta: {
          requiresAuth: true,
          title: '系统管理',
          permissions: ['system_log', 'system_backup', 'system_config', 'system_report']
        }
      },
      {
        path: 'reset-password',
        name: 'ResetPassword',
        component: () => import('@/views/ResetPassword.vue'),
        meta: { requiresAuth: true, title: '修改密码' }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue'),
    meta: { title: '页面未找到' }
  }
]

/**
 * 创建 Vue Router 实例
 *
 * 使用 HTML5 History 模式，URL 更美观（不需要 # 前缀）
 * 需要服务器配置支持，所有路由都指向 index.html
 */
const router = createRouter({
  history: createWebHistory(),
  routes
})

/**
 * 全局前置路由守卫
 *
 * 在每次路由跳转前执行，用于：
 * 1. 设置页面标题
 * 2. 检查用户登录状态
 * 3. 验证用户访问权限
 * 4. 处理登录页重定向
 *
 * @param {Object} to - 目标路由对象
 * @param {Object} from - 来源路由对象
 * @param {Function} next - 路由跳转控制函数
 */
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()

  // ========== 1. 设置页面标题 ==========
  document.title = to.meta.title ? `${to.meta.title} - 材料管理系统` : '材料管理系统'

  // ========== 2. 检查是否需要认证 ==========
  if (to.meta.requiresAuth) {
    if (!authStore.isAuthenticated) {
      next({ name: 'Login', query: { redirect: to.fullPath } })
      return
    }

    // ========== 3. 检查访问权限（使用静态工具函数，避免循环） ==========
    if (to.meta.permissions && to.meta.permissions.length > 0) {
      // 创建用户对象
      const user = {
        isAdmin: authStore.isAdmin,
        permissions: Array.isArray(authStore.permissions) ? authStore.permissions : []
      }

      // 管理员直接通过
      if (user.isAdmin) {
        next()
        return
      }

      // 检查是否有任一权限
      if (!hasAnyPermission(user, to.meta.permissions)) {
        // 用户没有权限，重定向到仪表板
        // 使用控制台警告，避免 ElMessage 触发组件更新
        console.warn('[路由守卫] 权限不足:', {
          path: to.path,
          required: to.meta.permissions,
          user: user.permissions
        })
        next({ name: 'Dashboard' })
        return
      }
    }
  }

  // ========== 4. 已登录用户访问登录页的处理 ==========
  if (to.name === 'Login' && authStore.isAuthenticated) {
    next({ name: 'Dashboard' })
    return
  }

  // ========== 5. 允许路由跳转 ==========
  next()
})

/**
 * 导出 Router 实例
 *
 * 在 main.js 中导入并注册到 Vue 应用
 *
 * @example
 * import router from './router'
 * app.use(router)
 */
export default router
