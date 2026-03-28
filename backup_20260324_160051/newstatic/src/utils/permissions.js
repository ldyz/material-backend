/**
 * 权限检查工具模块
 *
 * 提供静态的权限检查方法，避免 Vue 响应式系统的循环依赖
 *
 * @module permissions
 * @author Material Management System
 * @date 2026-02-01
 */

/**
 * 检查用户是否有指定权限
 *
 * @param {Object} user - 用户对象
 * @param {string[]} user.permissions - 用户权限列表
 * @param {boolean} user.isAdmin - 是否为管理员
 * @param {string} permission - 要检查的权限
 * @returns {boolean} 是否拥有该权限
 */
export function hasPermission(user, permission) {
  if (!user) return false
  if (user.isAdmin) return true
  if (!user.permissions) return false
  return user.permissions.includes(permission)
}

/**
 * 检查用户是否有权限组中的任一权限
 *
 * @param {Object} user - 用户对象
 * @param {string[]} user.permissions - 用户权限列表
 * @param {boolean} user.isAdmin - 是否为管理员
 * @param {string[]} permissions - 要检查的权限数组
 * @returns {boolean} 是否拥有至少一个权限
 */
export function hasAnyPermission(user, permissions) {
  if (!user) return false
  if (user.isAdmin) return true
  if (!permissions || permissions.length === 0) return true
  if (!user.permissions) return false
  return permissions.some(perm => user.permissions.includes(perm))
}

/**
 * 检查用户是否有权限组中的所有权限
 *
 * @param {Object} user - 用户对象
 * @param {string[]} user.permissions - 用户权限列表
 * @param {boolean} user.isAdmin - 是否为管理员
 * @param {string[]} permissions - 要检查的权限数组
 * @returns {boolean} 是否拥有所有权限
 */
export function hasAllPermissions(user, permissions) {
  if (!user) return false
  if (user.isAdmin) return true
  if (!permissions || permissions.length === 0) return true
  if (!user.permissions) return false
  return permissions.every(perm => user.permissions.includes(perm))
}

/**
 * 过滤需要权限的菜单项
 *
 * @param {Array} menus - 菜单配置数组
 * @param {Object} user - 用户对象
 * @returns {Array} 过滤后的菜单数组
 */
export function filterMenusByPermission(menus, user) {
  if (!menus) return []
  if (!user) return []
  if (user.isAdmin) return menus

  const result = []
  for (const menu of menus) {
    // 检查父菜单权限
    const hasMenuPermission = !menu.permissions || menu.permissions.length === 0 || hasAnyPermission(user, menu.permissions)

    if (hasMenuPermission) {
      // 如果有子菜单，递归过滤子菜单
      if (menu.children && menu.children.length > 0) {
        const filteredChildren = filterMenusByPermission(menu.children, user)
        // 只有当有可见的子菜单时，才添加父菜单
        if (filteredChildren.length > 0) {
          result.push({
            ...menu,
            children: filteredChildren
          })
        }
      } else {
        // 没有子菜单，直接添加
        result.push(menu)
      }
    }
  }
  return result
}

/**
 * 创建权限检查后的菜单缓存
 * 这个函数用于组件初始化时一次性计算菜单
 *
 * @param {Array} menus - 菜单配置数组
 * @param {Object} authStore - 认证store
 * @returns {Array} 过滤后的菜单数组
 */
export function createVisibleMenus(menus, authStore) {
  const user = {
    isAdmin: authStore.isAdmin,
    permissions: authStore.permissions || []
  }
  return filterMenusByPermission(menus, user)
}
