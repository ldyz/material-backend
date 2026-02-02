<template>
  <router-view />
</template>

<script setup>
import { onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useUserStore } from '@/stores/user'

const authStore = useAuthStore()
const userStore = useUserStore()

onMounted(async () => {
  // 如果已登录，获取用户信息
  if (authStore.isAuthenticated && !userStore.userInfo) {
    try {
      const user = await authStore.fetchCurrentUser()
      userStore.setUserInfo(user)
    } catch (error) {
      // Token 可能已过期，清除数据
      await authStore.logout()
      userStore.clearUserInfo()
    }
  }
})
</script>

<style>
/* 全局样式 */
* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

html, body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  font-size: 14px;
  line-height: 1.5;
  color: #323233;
  background-color: #f7f8fa;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

#app {
  min-height: 100vh;
  /* 状态栏安全区域适配 */
  padding-top: env(safe-area-inset-top);
}

/* 安全区域适配 - 顶部 */
.safe-area-inset-top {
  padding-top: env(safe-area-inset-top);
}

/* 清除浮动 */
.clearfix::after {
  content: '';
  display: block;
  clear: both;
}

/* 文本溢出省略 */
.ellipsis {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 多行文本溢出省略 */
.multi-ellipsis-2 {
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 安全区域适配 */
.safe-area-inset-bottom {
  padding-bottom: env(safe-area-inset-bottom);
}

/* 页面容器 */
.page-container {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding-bottom: 50px; /* tabbar 高度 */
}

.loading-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 200px;
}

.empty-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 300px;
  color: #969799;
}

.empty-container .empty-text {
  margin-top: 16px;
  font-size: 14px;
}

/* 修复 Toast 样式 - Toast 挂载在 body 上，不需要 :deep() */
.van-toast {
  background-color: rgba(0, 0, 0, 0.8) !important;
  color: #fff !important;
}

.van-toast--text {
  min-width: 120px;
  padding: 8px 16px;
}

.van-toast--fail {
  background-color: rgba(238, 10, 36, 0.9) !important;
  color: #fff !important;
}

.van-toast--success {
  background-color: rgba(7, 193, 96, 0.9) !important;
  color: #fff !important;
}

.van-toast--loading {
  background-color: rgba(0, 0, 0, 0.8) !important;
  color: #fff !important;
}

/* 确保 Toast 内的文字也是白色 */
.van-toast__text {
  color: #fff !important;
}
</style>
