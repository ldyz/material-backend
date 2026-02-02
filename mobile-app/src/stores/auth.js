import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login as loginApi, logout as logoutApi, getCurrentUser } from '@/api/auth'
import { storage } from '@/utils/storage'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(storage.getToken())
  const isAuthenticated = ref(!!storage.getToken())

  /**
   * 登录
   * @param {string} username - 用户名
   * @param {string} password - 密码
   * @returns {Promise}
   */
  async function login(username, password) {
    try {
      const response = await loginApi(username, password)
      // 登录接口特殊格式：{ success: true, token: "...", data: {...} }
      token.value = response.token
      isAuthenticated.value = true
      storage.setToken(response.token)
      return response
    } catch (error) {
      throw error
    }
  }

  /**
   * 登出
   * @returns {Promise}
   */
  async function logout() {
    try {
      await logoutApi()
    } finally {
      token.value = ''
      isAuthenticated.value = false
      storage.clear()
    }
  }

  /**
   * 获取当前用户信息
   * @returns {Promise}
   */
  async function fetchCurrentUser() {
    try {
      const response = await getCurrentUser()
      // 后端返回标准格式：{ success: true, data: {...} }
      return response.data
    } catch (error) {
      throw error
    }
  }

  return {
    token,
    isAuthenticated,
    login,
    logout,
    fetchCurrentUser,
  }
})
