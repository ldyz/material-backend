import { defineStore } from 'pinia'
import { ref } from 'vue'
import { storage } from '@/utils/storage'
import { login as loginApi, logout as logoutApi } from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(storage.getToken())
  const isAuthenticated = ref(!!token.value)

  function setToken(newToken) {
    token.value = newToken
    isAuthenticated.value = true
    storage.setToken(newToken)
  }

  function clearAuth() {
    token.value = null
    isAuthenticated.value = false
    storage.clear()
  }

  async function login(username, password) {
    try {
      const response = await loginApi({ username, password })
      const newToken = response.meta?.token || response.token
      setToken(newToken)
      return response
    } catch (error) {
      throw error
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
    isAuthenticated,
    setToken,
    clearAuth,
    login,
    logout
  }
})
