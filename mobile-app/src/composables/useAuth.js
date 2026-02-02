import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useUserStore } from '@/stores/user'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'

export function useAuth() {
  const authStore = useAuthStore()
  const userStore = useUserStore()
  const router = useRouter()

  const isAuthenticated = computed(() => authStore.isAuthenticated)
  const userInfo = computed(() => userStore.userInfo)
  const roles = computed(() => userStore.roles)
  const permissions = computed(() => userStore.permissions)

  /**
   * 登录
   * @param {string} username - 用户名
   * @param {string} password - 密码
   */
  async function login(username, password) {
    try {
      const data = await authStore.login(username, password)

      // 获取用户完整信息（如果失败也不影响登录）
      try {
        const user = await authStore.fetchCurrentUser()
        userStore.setUserInfo(user)
      } catch (userError) {
        // 获取用户信息失败，但登录已成功，继续
        console.error('获取用户信息失败，但登录已成功:', userError.message)
        // 不抛出错误，让登录流程继续
      }

      showToast({
        type: 'success',
        message: '登录成功',
      })

      return data
    } catch (error) {
      throw error
    }
  }

  /**
   * 登出
   */
  async function logout() {
    try {
      await authStore.logout()
      userStore.clearUserInfo()
      showToast({
        type: 'success',
        message: '已退出登录',
      })
      router.push('/login')
    } catch (error) {
      // 即使失败也清除本地数据
      userStore.clearUserInfo()
      router.push('/login')
    }
  }

  return {
    isAuthenticated,
    userInfo,
    roles,
    permissions,
    login,
    logout,
  }
}
