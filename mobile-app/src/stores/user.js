import { defineStore } from 'pinia'
import { ref } from 'vue'
import { storage } from '@/utils/storage'

export const useUserStore = defineStore('user', () => {
  const userInfo = ref(storage.getUserInfo())
  const roles = ref([])
  const permissions = ref(storage.getPermissions())

  // 初始化时处理 roles（去重）
  if (userInfo.value?.roles && Array.isArray(userInfo.value.roles)) {
    const roleNamesSet = new Set()
    userInfo.value.roles.forEach(r => {
      const name = typeof r === 'string' ? r : (r.name || '')
      if (name) {
        roleNamesSet.add(name)
      }
    })
    roles.value = Array.from(roleNamesSet)
  }

  /**
   * 设置用户信息
   * @param {Object} info - 用户信息
   */
  function setUserInfo(info) {
    userInfo.value = info

    // 处理 roles：后端返回的是对象数组 [{id, name, permissions}]，需要转换为角色名数组
    if (info.roles && Array.isArray(info.roles) && info.roles.length > 0) {
      // 提取角色名并去重（使用 Set 避免重复角色）
      const roleNamesSet = new Set()
      info.roles.forEach(r => {
        if (r.name) {
          roleNamesSet.add(r.name)
        }
      })
      roles.value = Array.from(roleNamesSet)

      // 合并所有角色的 permissions
      const allPermissions = new Set()
      info.roles.forEach(role => {
        if (role.permissions && Array.isArray(role.permissions)) {
          role.permissions.forEach(perm => allPermissions.add(perm))
        }
      })
      permissions.value = Array.from(allPermissions)
    } else {
      roles.value = []
      permissions.value = []
    }

    storage.setUserInfo(info)
    storage.setPermissions(permissions.value)
  }

  /**
   * 清除用户信息
   */
  function clearUserInfo() {
    userInfo.value = null
    roles.value = []
    permissions.value = []
    storage.removeUserInfo()
    storage.removePermissions()
  }

  /**
   * 检查是否有指定权限
   * @param {string} permission - 权限标识
   * @returns {boolean}
   */
  function hasPermission(permission) {
    return permissions.value?.includes(permission) || false
  }

  /**
   * 检查是否有任一权限
   * @param {Array} perms - 权限列表
   * @returns {boolean}
   */
  function hasAnyPermission(perms) {
    return perms.some(p => hasPermission(p))
  }

  /**
   * 检查是否有所有权限
   * @param {Array} perms - 权限列表
   * @returns {boolean}
   */
  function hasAllPermissions(perms) {
    return perms.every(p => hasPermission(p))
  }

  /**
   * 检查是否是管理员
   * @returns {boolean}
   */
  function isAdmin() {
    return roles.value?.includes('admin') || false
  }

  return {
    userInfo,
    roles,
    permissions,
    setUserInfo,
    clearUserInfo,
    hasPermission,
    hasAnyPermission,
    hasAllPermissions,
    isAdmin,
  }
})
