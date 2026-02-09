<template>
  <div class="profile-page">
    <van-nav-bar title="我的" />

    <!-- 用户信息卡片 -->
    <van-cell-group inset class="user-info">
      <van-cell center>
        <template #icon>
          <van-icon name="user-circle-o" size="60" color="#1989fa" />
        </template>
        <template #title>
          <span class="username">{{ authStore.username }}</span>
        </template>
        <template #label>
          <span class="user-role" v-if="authStore.user?.role">{{ getRoleText(authStore.user.role) }}</span>
        </template>
      </van-cell>
    </van-cell-group>

    <!-- 权限信息 -->
    <van-cell-group inset title="权限信息">
      <van-cell title="用户ID" :value="String(authStore.userId || '-')" />
      <van-cell title="角色" :value="getRoleText(authStore.user?.role || '-')" />
      <van-cell title="用户名" :value="authStore.user?.username || authStore.user?.name || '-'" />
      <van-cell
        v-if="authStore.user?.email"
        title="邮箱"
        :value="authStore.user.email"
      />
      <van-cell
        v-if="authStore.user?.phone"
        title="手机号"
        :value="authStore.user.phone"
      />
      <van-cell
        v-if="authStore.user?.department"
        title="部门"
        :value="authStore.user.department"
      />
    </van-cell-group>

    <!-- 设置 -->
    <van-cell-group inset title="设置">
      <van-cell title="刷新数据" icon="replay" is-link @click="handleRefresh" />
      <van-cell title="关于" icon="info-o" is-link @click="showAbout" />
    </van-cell-group>

    <!-- 关于弹窗 -->
    <van-dialog
      v-model:show="showAboutDialog"
      title="关于"
      :show-confirm-button="true"
      :show-cancel-button="false"
      confirm-button-text="我知道了"
      confirm-button-color="#1989fa"
      :close-on-click-overlay="true"
      teleport="body"
    >
      <div class="about-content">
        <p class="about-title">材料管理系统</p>
        <p class="about-version">v1.0.0</p>
        <div class="about-divider"></div>
        <p class="about-description">移动端应用</p>
        <p class="about-features">提供材料计划、入库、出库等功能</p>
      </div>
    </van-dialog>

    <!-- 退出登录 -->
    <div class="logout-section">
      <van-button round block type="danger" @click="handleLogout">
        退出登录
      </van-button>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showConfirmDialog } from 'vant'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const showAboutDialog = ref(false)

function getRoleText(role) {
  const roleMap = {
    admin: '管理员',
    manager: '经理',
    supervisor: '主管',
    user: '普通用户',
    super_admin: '超级管理员'
  }
  return roleMap[role] || role || '-'
}

async function handleRefresh() {
  try {
    showToast({ type: 'loading', message: '刷新中...', forbidClick: true })
    await authStore.fetchCurrentUser()
    showToast({ type: 'success', message: '刷新成功' })
  } catch (error) {
    showToast({ type: 'fail', message: '刷新失败' })
  }
}

function showAbout() {
  showAboutDialog.value = true
}

async function handleLogout() {
  try {
    await showConfirmDialog({
      title: '提示',
      message: '确定要退出登录吗？',
      teleport: 'body',
      confirmButtonColor: '#ee0a24'
    })
    await authStore.logout()
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

.user-info {
  margin-bottom: 16px;
}

.username {
  font-size: 20px;
  font-weight: bold;
  margin-left: 12px;
}

.user-role {
  font-size: 14px;
  color: #666;
  margin-left: 12px;
  display: inline-block;
  background-color: #f0f0f0;
  padding: 2px 8px;
  border-radius: 4px;
  margin-top: 4px;
}

:deep(.van-cell-group) {
  margin-bottom: 12px;
}

:deep(.van-cell-group__title) {
  padding-left: 16px;
  font-weight: bold;
}

.logout-section {
  padding: 0 16px;
  margin-top: 24px;
}

.about-dialog {
  text-align: center;
  line-height: 1.8;
  white-space: pre-line;
}

.about-content {
  padding: 20px 0;
}

.about-title {
  font-size: 20px;
  font-weight: bold;
  margin: 0 0 8px 0;
  color: #323233;
}

.about-version {
  font-size: 14px;
  color: #969799;
  margin: 0 0 16px 0;
}

.about-divider {
  height: 1px;
  background-color: #ebedf0;
  margin: 16px 0;
}

.about-description {
  font-size: 15px;
  color: #646566;
  margin: 0 0 8px 0;
}

.about-features {
  font-size: 14px;
  color: #969799;
  margin: 0;
}

:deep(.about-dialog .van-dialog__message) {
  text-align: center;
  font-size: 16px;
  line-height: 1.8;
  padding: 20px 0;
}
</style>
