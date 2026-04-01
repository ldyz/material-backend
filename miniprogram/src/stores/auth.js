import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { storage } from '@/utils/storage'
import { login as loginApi, logout as logoutApi, getCurrentUser as getCurrentUserApi, wechatLogin as wechatLoginApi, bindWechat as bindWechatApi } from '@/api/auth'

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
    if (user.value.roles && Array.isArray(user.value.roles)) {
      const perms = new Set()
      user.value.roles.forEach(role => {
        if (role.permissions && Array.isArray(role.permissions)) {
          role.permissions.forEach(p => perms.add(p))
        }
      })
      if (perms.size > 0) {
        return Array.from(perms)
      }
    }
    if (user.value.permissions && Array.isArray(user.value.permissions)) {
      return user.value.permissions
    }
    return getPermissionsFromRoles(user.value.roles || user.value.role)
  })

  const isAdmin = computed(() => {
    if (!user.value) return false
    if (user.value.roles && Array.isArray(user.value.roles)) {
      return user.value.roles.some(r => {
        const name = r.name || r
        return name === 'admin' || name === '管理员'
      })
    }
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
  }

  /**
   * 用户名密码登录
   */
  async function login(username, password) {
    try {
      const response = await loginApi({ username, password })
      const newToken = response.meta?.token || response.token
      setToken(newToken)

      if (response.data) {
        setUser(response.data)
      } else {
        await fetchCurrentUser()
      }

      return response
    } catch (error) {
      throw error
    }
  }

  /**
   * 微信登录
   */
  async function wechatLogin() {
    return new Promise((resolve, reject) => {
      // #ifdef MP-WEIXIN
      uni.login({
        provider: 'weixin',
        success: async (loginRes) => {
          try {
            const response = await wechatLoginApi(loginRes.code)

            // 检查是否需要绑定账号
            if (response.needBind) {
              reject({ needBind: true, openid: response.openid })
              return
            }

            // 登录成功
            const newToken = response.meta?.token || response.token
            setToken(newToken)

            if (response.data) {
              setUser(response.data)
            }

            resolve(response)
          } catch (error) {
            reject(error)
          }
        },
        fail: (error) => {
          reject(new Error(error.errMsg || '微信登录失败'))
        }
      })
      // #endif

      // #ifndef MP-WEIXIN
      reject(new Error('当前平台不支持微信登录'))
      // #endif
    })
  }

  /**
   * 绑定微信账号
   */
  async function bindWechat(username, password) {
    return new Promise((resolve, reject) => {
      // #ifdef MP-WEIXIN
      uni.login({
        provider: 'weixin',
        success: async (loginRes) => {
          try {
            const response = await bindWechatApi({
              username,
              password,
              code: loginRes.code
            })

            const newToken = response.meta?.token || response.token
            setToken(newToken)

            if (response.data) {
              setUser(response.data)
            }

            resolve(response)
          } catch (error) {
            reject(error)
          }
        },
        fail: (error) => {
          reject(new Error(error.errMsg || '微信登录失败'))
        }
      })
      // #endif

      // #ifndef MP-WEIXIN
      reject(new Error('当前平台不支持微信登录'))
      // #endif
    })
  }

  async function fetchCurrentUser() {
    try {
      const response = await getCurrentUserApi()
      if (response.data) {
        setUser(response.data)
      }
      return response.data
    } catch (error) {
      throw error
    }
  }

  /**
   * 初始化认证状态
   */
  async function initAuth() {
    if (token.value) {
      try {
        await fetchCurrentUser()
      } catch (error) {
        console.error('刷新用户信息失败:', error)
      }
    }
  }

  async function logout() {
    try {
      await logoutApi()
    } finally {
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
    setToken,
    setUser,
    clearAuth,
    login,
    wechatLogin,
    bindWechat,
    fetchCurrentUser,
    initAuth,
    logout
  }
})
