import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { storage } from '@/utils/storage'
import { login as loginApi, logout as logoutApi, getCurrentUser as getCurrentUserApi } from '@/api/auth'
import { initWebSocket, disconnectWebSocket } from '@/utils/websocket'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(storage.getToken())
  const user = ref(storage.getUser())
  const isAuthenticated = ref(!!token.value)

  const username = computed(() => user.value?.username || user.value?.name || '用户')
  const userId = computed(() => user.value?.id)
  const userRole = computed(() => user.value?.role || '')

  // 权限相关
  const permissions = computed(() => {
    if (!user.value) return []
    // 优先使用 permissions 字段
    if (user.value.permissions && Array.isArray(user.value.permissions)) {
      return user.value.permissions
    }
    // 兼容旧系统：从角色映射权限
    return getPermissionsFromRoles(user.value.roles || user.value.role)
  })

  const isAdmin = computed(() => {
    if (!user.value) return false
    // 检查 roles 数组
    if (user.value.roles && Array.isArray(user.value.roles)) {
      return user.value.roles.some(r => {
        const name = r.name || r
        return name === 'admin' || name === '管理员'
      })
    }
    // 检查 role 字段
    return user.value.role === 'admin'
  })

  /**
   * 从角色列表获取权限
   */
  function getPermissionsFromRoles(roles) {
    if (!roles) return []

    const roleArray = Array.isArray(roles) ? roles : [roles]
    const perms = new Set()

    const rolePermissionMap = {
      'admin': [
        'user_view', 'user_create', 'user_edit', 'user_delete',
        'role_view', 'role_create', 'role_edit', 'role_delete',
        'project_view', 'project_create', 'project_edit', 'project_delete',
        'material_view', 'material_plan_view', 'material_plan_create', 'material_plan_edit', 'material_plan_delete', 'material_plan_approve',
        'stock_view', 'inbound_view', 'inbound_create', 'inbound_edit', 'inbound_delete', 'inbound_approve',
        'requisition_view', 'requisition_create', 'requisition_edit', 'requisition_delete', 'requisition_approve',
        'appointment_view', 'appointment_create', 'appointment_edit', 'appointment_delete', 'appointment_approve', 'appointment_assign',
        'constructionlog_view', 'constructionlog_create', 'constructionlog_edit', 'constructionlog_delete',
        'progress_view', 'progress_edit',
        'system_config', 'audit_view', 'system_log', 'system_backup', 'system_report',
        'notification_view', 'notification_manage'
      ],
      'project_manager': [
        'project_view', 'project_create', 'project_edit',
        'material_view', 'material_plan_view', 'material_plan_create', 'material_plan_approve',
        'appointment_view', 'appointment_create', 'appointment_approve', 'appointment_assign',
        'constructionlog_view', 'constructionlog_create', 'constructionlog_edit',
        'progress_view', 'progress_edit'
      ],
      'foreman': [
        'appointment_view', 'appointment_approve',
        'constructionlog_view', 'constructionlog_create'
      ],
      'worker': [
        'appointment_view', 'appointment_create',
        'constructionlog_view', 'constructionlog_create'
      ],
      'keeper': [
        'stock_view', 'inbound_view', 'inbound_create', 'inbound_approve',
        'requisition_view', 'requisition_approve'
      ],
      'material_staff': [
        'material_view', 'material_plan_view', 'material_plan_create', 'material_plan_edit',
        'stock_view', 'requisition_view'
      ],
      'appointment_admin': [
        'appointment_view', 'appointment_create', 'appointment_edit', 'appointment_delete', 'appointment_approve', 'appointment_assign'
      ]
    }

    roleArray.forEach(role => {
      const roleName = role.name || role
      const rolePerms = rolePermissionMap[roleName] || []
      rolePerms.forEach(p => perms.add(p))
    })

    return Array.from(perms)
  }

  /**
   * 检查是否有指定权限
   */
  function hasPermission(permission) {
    if (isAdmin.value) return true
    return permissions.value.includes(permission)
  }

  /**
   * 检查是否有权限列表中的任一权限
   */
  function hasAnyPermission(permissionList) {
    if (!permissionList || permissionList.length === 0) return true
    if (isAdmin.value) return true
    return permissionList.some(p => permissions.value.includes(p))
  }

  /**
   * 检查是否有权限列表中的所有权限
   */
  function hasAllPermissions(permissionList) {
    if (!permissionList || permissionList.length === 0) return true
    if (isAdmin.value) return true
    return permissionList.every(p => permissions.value.includes(p))
  }

  function setToken(newToken) {
    token.value = newToken
    isAuthenticated.value = true
    storage.setToken(newToken)
  }

  function setUser(newUser) {
    user.value = newUser
    storage.setUser(newUser)
  }

  function clearAuth() {
    token.value = null
    user.value = null
    isAuthenticated.value = false
    storage.clear()
    // 断开 WebSocket 连接
    disconnectWebSocket()
  }

  async function login(username, password) {
    try {
      const response = await loginApi({ username, password })
      const newToken = response.meta?.token || response.token
      setToken(newToken)

      // 如果返回中包含用户信息，保存用户信息
      if (response.data) {
        setUser(response.data)
      } else {
        // 否则获取当前用户信息
        await fetchCurrentUser()
      }

      // 登录成功后连接 WebSocket
      initWebSocket()

      return response
    } catch (error) {
      throw error
    }
  }

  async function fetchCurrentUser() {
    try {
      const response = await getCurrentUserApi()
      if (response.data) {
        setUser(response.data)
      }
      return response.data
    } catch (error) {
      console.error('获取用户信息失败:', error)
      throw error
    }
  }

  async function logout() {
    try {
      await logoutApi()
    } finally {
      // 断开 WebSocket 连接
      disconnectWebSocket()
      clearAuth()
    }
  }

  return {
    token,
    user,
    isAuthenticated,
    username,
    userId,
    userRole,
    permissions,
    isAdmin,
    hasPermission,
    hasAnyPermission,
    hasAllPermissions,
    setToken,
    setUser,
    clearAuth,
    login,
    fetchCurrentUser,
    logout
  }
})
