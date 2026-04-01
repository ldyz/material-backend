import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore('user', () => {
  const userInfo = ref(null)
  const roles = ref([])
  const permissions = ref([])

  function setUserInfo(info) {
    userInfo.value = info
    if (info.roles && Array.isArray(info.roles)) {
      roles.value = info.roles.map(r => typeof r === 'string' ? r : r.name).filter(Boolean)
      permissions.value = info.roles.flatMap(r => r.permissions || [])
    }
  }

  function clearUserInfo() {
    userInfo.value = null
    roles.value = []
    permissions.value = []
  }

  function hasPermission(permission) {
    return permissions.value?.includes(permission) || false
  }

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
    isAdmin
  }
})
