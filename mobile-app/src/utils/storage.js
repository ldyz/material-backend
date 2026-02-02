import { STORAGE_KEYS } from './constants'

// 本地存储工具
export const storage = {
  // 设置 token
  setToken(token) {
    localStorage.setItem(STORAGE_KEYS.TOKEN, token)
  },

  // 获取 token
  getToken() {
    return localStorage.getItem(STORAGE_KEYS.TOKEN)
  },

  // 移除 token
  removeToken() {
    localStorage.removeItem(STORAGE_KEYS.TOKEN)
  },

  // 设置用户信息
  setUserInfo(userInfo) {
    localStorage.setItem(STORAGE_KEYS.USER_INFO, JSON.stringify(userInfo))
  },

  // 获取用户信息
  getUserInfo() {
    const info = localStorage.getItem(STORAGE_KEYS.USER_INFO)
    return info ? JSON.parse(info) : null
  },

  // 移除用户信息
  removeUserInfo() {
    localStorage.removeItem(STORAGE_KEYS.USER_INFO)
  },

  // 设置权限
  setPermissions(permissions) {
    localStorage.setItem(STORAGE_KEYS.PERMISSIONS, JSON.stringify(permissions))
  },

  // 获取权限
  getPermissions() {
    const perms = localStorage.getItem(STORAGE_KEYS.PERMISSIONS)
    return perms ? JSON.parse(perms) : []
  },

  // 移除权限
  removePermissions() {
    localStorage.removeItem(STORAGE_KEYS.PERMISSIONS)
  },

  // 清除所有数据
  clear() {
    localStorage.clear()
  },
}
