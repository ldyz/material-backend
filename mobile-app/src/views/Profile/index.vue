<template>
  <div class="profile-page">
    <van-nav-bar title="我的" />

    <van-cell-group inset class="user-info">
      <van-cell center>
        <template #icon>
          <van-icon name="user-circle-o" size="60" color="#1989fa" />
        </template>
        <template #title>
          <span class="username">用户</span>
        </template>
      </van-cell>
    </van-cell-group>

    <van-cell-group inset title="设置">
      <van-cell title="关于" icon="info-o" is-link />
    </van-cell-group>

    <div class="logout-section">
      <van-button round block type="danger" @click="handleLogout">
        退出登录
      </van-button>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { showConfirmDialog } from 'vant'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

async function handleLogout() {
  try {
    await showConfirmDialog({
      title: '提示',
      message: '确定要退出登录吗？'
    })
    authStore.logout()
    router.push('/login')
  } catch {
    // 用户取消
  }
}
</script>

<style scoped>
.profile-page {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding: 16px 0;
}

.username {
  font-size: 18px;
  font-weight: bold;
  margin-left: 12px;
}

.logout-section {
  padding: 0 16px;
  margin-top: 24px;
}
</style>
