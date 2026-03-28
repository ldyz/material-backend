/**
 * Vue 权限指令
 *
 * 用于在模板中根据权限控制元素的显示
 * 避免在模板中使用复杂的响应式逻辑
 *
 * @module directives/permission
 * @author Material Management System
 * @date 2026-02-01
 */

import { hasPermission, hasAnyPermission, hasAllPermissions } from '@/utils/permissions'

/**
 * v-permission 指令
 * 用法：v-permission="'material_create'"
 */
export const permissionDirective = {
  mounted(el, binding) {
    const authStore = window.__authStore__
    if (!authStore) {
      el.style.display = 'none'
      return
    }

    const permission = binding.value
    const user = {
      isAdmin: authStore.isAdmin,
      permissions: authStore.permissions || []
    }

    if (!hasPermission(user, permission)) {
      el.style.display = 'none'
    }
  },
  updated(el, binding) {
    const authStore = window.__authStore__
    if (!authStore) {
      el.style.display = 'none'
      return
    }

    const permission = binding.value
    const user = {
      isAdmin: authStore.isAdmin,
      permissions: authStore.permissions || []
    }

    if (!hasPermission(user, permission)) {
      el.style.display = 'none'
    } else {
      el.style.display = ''
    }
  }
}

/**
 * v-permission-any 指令
 * 用法：v-permission-any="['material_create', 'material_edit']"
 */
export const permissionAnyDirective = {
  mounted(el, binding) {
    const authStore = window.__authStore__
    if (!authStore) {
      el.style.display = 'none'
      return
    }

    const permissions = binding.value
    const user = {
      isAdmin: authStore.isAdmin,
      permissions: authStore.permissions || []
    }

    if (!hasAnyPermission(user, permissions)) {
      el.style.display = 'none'
    }
  },
  updated(el, binding) {
    const authStore = window.__authStore__
    if (!authStore) {
      el.style.display = 'none'
      return
    }

    const permissions = binding.value
    const user = {
      isAdmin: authStore.isAdmin,
      permissions: authStore.permissions || []
    }

    if (!hasAnyPermission(user, permissions)) {
      el.style.display = 'none'
    } else {
      el.style.display = ''
    }
  }
}
