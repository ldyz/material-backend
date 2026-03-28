/**
 * Vue 权限组合式函数
 *
 * 用于在组件中检查权限，避免响应式循环
 *
 * @module composables/usePermission
 * @author Material Management System
 * @date 2026-02-01
 */

import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { hasPermission, hasAnyPermission, hasAllPermissions } from '@/utils/permissions'

/**
 * 权限检查组合式函数
 *
 * @returns {Object} 权限检查方法
 */
export function usePermission() {
  const authStore = useAuthStore()

  // 创建用户对象（缓存，避免重复创建）
  const user = computed(() => ({
    isAdmin: authStore.isAdmin,
    permissions: authStore.permissions || []
  }))

  /**
   * 检查是否有指定权限
   *
   * @param {string} permission - 权限标识
   * @returns {boolean} 是否拥有该权限
   */
  const checkPermission = (permission) => {
    return hasPermission(user.value, permission)
  }

  /**
   * 检查是否有任一权限
   *
   * @param {string[]} permissions - 权限数组
   * @returns {boolean} 是否拥有至少一个权限
   */
  const checkAnyPermission = (permissions) => {
    return hasAnyPermission(user.value, permissions)
  }

  /**
   * 检查是否有所有权限
   *
   * @param {string[]} permissions - 权限数组
   * @returns {boolean} 是否拥有所有权限
   */
  const checkAllPermissions = (permissions) => {
    return hasAllPermissions(user.value, permissions)
  }

  return {
    user,
    isAdmin: computed(() => user.value.isAdmin),
    permissions: computed(() => user.value.permissions),
    checkPermission,
    checkAnyPermission,
    checkAllPermissions,
    hasPermission: checkPermission,
    hasAnyPermission: checkAnyPermission,
    hasAllPermissions: checkAllPermissions
  }
}
