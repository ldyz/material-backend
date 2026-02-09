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
    setToken,
    setUser,
    clearAuth,
    login,
    fetchCurrentUser,
    logout
  }
})
