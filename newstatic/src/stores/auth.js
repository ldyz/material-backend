/**
 * 认证管理 Store
 *
 * 使用 Pinia 管理用户认证状态，包括：
 * - Token 管理
 * - 用户信息
 * - 权限列表
 * - 登录/登出
 * - 权限检查
 *
 * @module AuthStore
 * @author Material Management System
 * @date 2025-01-27
 */

import { defineStore } from 'pinia'
import { authApi } from '@/api'
import { ElMessage } from 'element-plus'

/**
 * 认证状态管理 Store
 *
 * 状态设计说明：
 * - token: JWT 认证令牌
 * - user: 当前登录用户信息
 * - permissions: 用户权限列表
 * - isAdmin: 是否为管理员
 * - loginTime: 登录时间戳（用于判断 token 是否过期）
 */
export const useAuthStore = defineStore('auth', {
  /**
   * 状态定义
   *
   * 使用 localStorage 持久化存储，刷新页面后自动恢复
   */
  state: () => ({
    // JWT 认证令牌
    token: localStorage.getItem('token') || '',

    // 当前登录用户信息
    user: JSON.parse(localStorage.getItem('user') || 'null'),

    // 用户权限列表（从角色的 permissions 中提取）
    permissions: JSON.parse(localStorage.getItem('permissions') || '[]'),

    // 是否为管理员（拥有所有权限）
    isAdmin: JSON.parse(localStorage.getItem('isAdmin') || 'false'),

    // 登录时间戳（毫秒），用于判断 72 小时后 token 过期
    loginTime: parseInt(localStorage.getItem('loginTime') || '0')
  }),

  /**
   * Getters - 计算属性
   *
   * 类似 Vue 的 computed，从 state 派生出有用的值
   */
  getters: {
    /**
     * 检查用户是否已登录
     *
     * 判断条件：
     * 1. token 存在
     * 2. 登录时间未超过 72 小时
     *
     * @param {Object} state - store state
     * @returns {boolean} 是否已登录
     */
    isAuthenticated: (state) => {
      // 没有 token 说明未登录
      if (!state.token) return false

      // 检查登录是否过期（72小时 = 72 * 3600 * 1000 毫秒）
      const maxAge = 72 * 3600 * 1000
      if (state.loginTime && Date.now() - state.loginTime > maxAge) {
        return false
      }

      return true
    },

    /**
     * 获取用户显示名称
     *
     * 优先级：username > full_name > '用户'
     *
     * @param {Object} state - store state
     * @returns {string} 用户显示名称
     */
    displayName: (state) => {
      return state.user?.username || state.user?.full_name || '用户'
    }
  },

  /**
   * Actions - 异步方法和业务逻辑
   *
   * Actions 可以是异步的，用于处理业务逻辑和 API 调用
   */
  actions: {
    /**
     * 用户登录
     *
     * 流程：
     * 1. 调用登录 API
     * 2. 保存 token 和用户信息
     * 3. 提取用户权限（支持多角色）
     * 4. 保存到 localStorage
     *
     * @param {Object} credentials - 登录凭证
     * @param {string} credentials.username - 用户名
     * @param {string} credentials.password - 密码
     * @returns {Promise<boolean>} 登录是否成功
     */
    async login(credentials) {
      try {
        // 调用登录 API
        const response = await authApi.login(credentials)

        // 从响应中获取 token（支持 meta.token 和直接返回 token 两种格式）
        const token = response.meta?.token || response.token

        if (token) {
          // 保存 token
          this.token = token

          // 保存用户信息（从 data 中获取）
          this.user = response.data

          // 提取用户权限
          const userData = response.data || response.user

          if (userData.permissions && Array.isArray(userData.permissions)) {
            // 格式1: 直接返回 permissions 数组
            this.permissions = userData.permissions
          } else if (userData.roles && Array.isArray(userData.roles)) {
            // 格式2: 通过角色关联权限，需要合并所有角色的权限
            this.permissions = [
              ...new Set(
                userData.roles.flatMap(role => {
                  // 处理权限字符串（逗号分隔）
                  if (typeof role.permissions === 'string') {
                    return role.permissions.split(',').filter(p => p.trim())
                  } else if (Array.isArray(role.permissions)) {
                    return role.permissions
                  }
                  return []
                })
              )
            ]
          } else {
            // 没有权限信息
            this.permissions = []
          }

          // 判断是否为管理员
          this.isAdmin = userData.is_admin || userData.role === 'admin'

          // 记录登录时间
          this.loginTime = Date.now()

          // 持久化到 localStorage
          this.saveToStorage()

          ElMessage.success('登录成功')
          return true
        }

        return false
      } catch (error) {
        console.error('登录失败:', error)
        throw error
      }
    },

    /**
     * 用户登出
     *
     * 流程：
     * 1. 调用登出 API（可选失败）
     * 2. 清除本地认证信息
     * 3. 跳转到登录页
     */
    async logout() {
      try {
        // 尝试调用登出 API（即使失败也继续执行）
        await authApi.logout()
      } catch (error) {
        console.error('登出API调用失败:', error)
      } finally {
        // 无论 API 调用成功与否，都清除本地认证信息
        this.clearAuth()
        ElMessage.info('已退出登录')
      }
    },

    /**
     * 获取当前用户信息
     *
     * 用于：
     * - 应用启动时验证用户身份
     * - 刷新用户权限信息
     *
     * @returns {Promise<Object>} 用户信息
     */
    async getCurrentUser() {
      try {
        // 调用获取当前用户 API
        const response = await authApi.getCurrentUser()

        // 兼容两种数据格式
        const userData = response.data || response

        // 更新用户状态
        this.user = userData

        // 提取权限（与 login 方法相同的逻辑）
        if (userData.permissions && Array.isArray(userData.permissions)) {
          // 格式1: 直接返回 permissions 数组
          this.permissions = userData.permissions
        } else if (userData.roles && Array.isArray(userData.roles)) {
          // 格式2: 通过角色关联权限，需要合并所有角色的权限
          this.permissions = [
            ...new Set(
              userData.roles.flatMap(role => {
                // 解析角色的 permissions 字符串
                if (typeof role.permissions === 'string') {
                  return role.permissions.split(',').filter(p => p.trim())
                } else if (Array.isArray(role.permissions)) {
                  return role.permissions
                }
                return []
              })
            )
          ]
        } else {
          this.permissions = []
        }

        this.isAdmin = userData.is_admin || userData.role === 'admin'

        // 持久化到 localStorage
        this.saveToStorage()

        return userData
      } catch (error) {
        console.error('获取用户信息失败:', error)
        // API 调用失败，清除认证信息（可能 token 过期）
        this.clearAuth()
        throw error
      }
    },

    /**
     * 修改当前用户密码
     *
     * 修改成功后会自动登出，需要用新密码重新登录
     *
     * @param {Object} passwords - 密码信息
     * @param {string} passwords.oldPassword - 原密码
     * @param {string} passwords.newPassword - 新密码
     * @returns {Promise<boolean>} 修改是否成功
     */
    async changePassword(passwords) {
      try {
        const response = await authApi.changePassword(passwords)

        if (response.success) {
          ElMessage.success('密码修改成功，请重新登录')

          // 清除认证信息，强制重新登录
          this.clearAuth()

          return true
        }

        return false
      } catch (error) {
        console.error('修改密码失败:', error)
        throw error
      }
    },

    /**
     * 保存认证信息到 localStorage
     *
     * 用于持久化存储，刷新页面后自动恢复
     */
    saveToStorage() {
      localStorage.setItem('token', this.token)
      localStorage.setItem('user', JSON.stringify(this.user))
      localStorage.setItem('permissions', JSON.stringify(this.permissions))
      localStorage.setItem('isAdmin', JSON.stringify(this.isAdmin))
      localStorage.setItem('loginTime', this.loginTime.toString())
    },

    /**
     * 清除认证信息
     *
     * 清除所有用户相关的本地存储
     */
    clearAuth() {
      this.token = ''
      this.user = null
      this.permissions = []
      this.isAdmin = false
      this.loginTime = 0

      localStorage.removeItem('token')
      localStorage.removeItem('user')
      localStorage.removeItem('permissions')
      localStorage.removeItem('isAdmin')
      localStorage.removeItem('loginTime')
    },

    /**
     * 检查是否有指定权限
     *
     * 管理员默认拥有所有权限
     *
     * @param {string} permission - 权限标识（如 'material_create'）
     * @returns {boolean} 是否拥有该权限
     */
    hasPermission(permission) {
      // 管理员拥有所有权限
      if (this.isAdmin) return true

      // 普通用户检查权限列表
      return this.permissions.includes(permission)
    },

    /**
     * 检查是否有权限组中的权限
     *
     * 支持两种模式：
     * - requireAll=false（默认）: 只需要拥有权限组中的任意一个权限
     * - requireAll=true: 需要拥有权限组中的所有权限
     *
     * @param {string[]} permissions - 权限标识数组
     * @param {boolean} requireAll - 是否需要所有权限
     * @returns {boolean} 是否拥有权限
     *
     * @example
     * // 只需要以下任一权限即可
     * hasPermissionGroup(['material_create', 'material_edit'], false)
     *
     * // 需要同时拥有以下两个权限
     * hasPermissionGroup(['material_create', 'material_delete'], true)
     */
    hasPermissionGroup(permissions, requireAll = false) {
      // 管理员拥有所有权限
      if (this.isAdmin) return true

      // 权限组为空，默认允许访问
      if (!permissions || permissions.length === 0) return true

      if (requireAll) {
        // 需要拥有所有权限
        return permissions.every(permission => this.permissions.includes(permission))
      } else {
        // 只需要拥有任意一个权限
        return permissions.some(permission => this.permissions.includes(permission))
      }
    },

    /**
     * 刷新用户信息
     *
     * 从服务器重新获取用户信息和权限
     * 用于权限更新后刷新本地状态
     */
    async refreshUserInfo() {
      await this.getCurrentUser()
    }
  }
})
