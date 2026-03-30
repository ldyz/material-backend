<template>
  <div class="profile-page">
    <van-nav-bar title="我的" />

    <!-- 用户信息卡片 -->
    <van-cell-group inset class="user-info">
      <van-cell center>
        <template #icon>
          <div class="avatar-wrapper">
            <van-image
              v-if="avatarUrl && !imageError"
              round
              width="60"
              height="60"
              :src="avatarUrl"
              @click="showAvatarPreview"
              @error="handleImageError"
            />
            <van-icon
              v-else
              name="user-circle-o"
              size="60"
              :color="getAvatarColor()"
              @click="showAvatarPreview"
            />
          </div>
        </template>
        <template #title>
          <span class="username">{{ authStore.user?.full_name || authStore.username }}</span>
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
      <van-cell title="上传头像" icon="photo-o" is-link @click="handleAvatarClick" />
      <van-cell
        title="检查更新"
        icon="upgrade"
        is-link
        @click="handleCheckUpdate"
      >
        <template #right-icon>
          <van-tag v-if="updateStatus" :type="updateStatus.type">
            {{ updateStatus.text }}
          </van-tag>
          <van-icon v-else name="arrow" />
        </template>
      </van-cell>
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
        <p class="about-title">化建仪表</p>
        <p class="about-version">v{{ currentVersion || '1.0.9' }}</p>
        <p class="about-date">发布日期: 2026-02-13</p>
        <div class="about-divider"></div>
        <p class="about-description">移动端应用</p>
        <p class="about-features">提供材料计划、入库、出库、预约管理等功能</p>
      </div>
    </van-dialog>

    <!-- 退出登录 -->
    <div class="logout-section">
      <van-button round block type="danger" @click="handleLogout">
        退出登录
      </van-button>
    </div>

    <!-- 隐藏的文件上传input -->
    <input
      ref="avatarInputRef"
      type="file"
      accept="image/*"
      style="display: none"
      @change="handleAvatarChange"
    />

    <!-- 头像预览对话框 -->
    <UserAvatarPreview ref="avatarPreviewRef" :avatar="authStore.user?.avatar" />

    <!-- 头像裁剪对话框 -->
    <AvatarCropperDialog
      v-model="showCropperDialog"
      :image-file="selectedAvatarFile"
      @success="handleAvatarSuccess"
    />
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showConfirmDialog, showSuccessToast, showFailToast, showDialog } from 'vant'
import { useAuthStore } from '@/stores/auth'
import * as authApi from '@/api/auth'
import AvatarCropperDialog from '@/components/AvatarCropperDialog.vue'
import UserAvatarPreview from '@/components/common/UserAvatarPreview.vue'
import { useAppUpdate } from '@/composables/useAppUpdate'
import { getAssetUrl } from '@/utils/request'
import { logger } from '@/utils/logger'

const router = useRouter()
const authStore = useAuthStore()
const showAboutDialog = ref(false)
const avatarInputRef = ref(null)
const avatarPreviewRef = ref(null)
const showCropperDialog = ref(false)
const selectedAvatarFile = ref(null)
const updateStatus = ref(null)
const imageError = ref(false)

// 使用应用更新 composable
const { isChecking, checkUpdate, downloadAndInstall, latestVersion, currentVersion } = useAppUpdate()

// 获取完整的头像URL
const avatarUrl = computed(() => {
  return getAssetUrl(authStore.user?.avatar)
})

// 监听头像URL变化，重置错误状态
watch(avatarUrl, () => {
  imageError.value = false
})

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

function getAvatarColor() {
  const colors = ['#1989fa', '#07c160', '#ff976a', '#ee0a24', '#909399']
  const userId = authStore.user?.id || authStore.userId || 0
  return colors[userId % colors.length]
}

function showAvatarPreview() {
  avatarPreviewRef.value?.show()
}

function handleImageError() {
  logger.warn('[Profile] 头像加载失败，使用默认头像')
  imageError.value = true
}

function handleAvatarClick() {
  avatarInputRef.value?.click()
  // 重置错误状态，以便下次尝试加载
  imageError.value = false
}

async function handleAvatarChange(event) {
  const file = event.target.files?.[0]
  if (!file) return

  // 验证文件大小（5MB - 裁剪前可以大一些）
  if (file.size > 5 * 1024 * 1024) {
    showToast({ type: 'fail', message: '图片文件大小不能超过5MB' })
    return
  }

  // 验证文件类型
  if (!file.type.startsWith('image/')) {
    showToast({ type: 'fail', message: '只支持图片格式的文件' })
    return
  }

  // 打开裁剪对话框
  selectedAvatarFile.value = file
  showCropperDialog.value = true

  // 清空 input
  if (avatarInputRef.value) {
    avatarInputRef.value.value = ''
  }
}

async function handleAvatarSuccess() {
  // 刷新用户信息
  await authStore.fetchCurrentUser()
  selectedAvatarFile.value = null
  // 重置错误状态
  imageError.value = false
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

// 检查更新
async function handleCheckUpdate() {
  if (isChecking.value) {
    showToast({ type: 'loading', message: '正在检查更新...', forbidClick: true })
    return
  }

  try {
    const result = await checkUpdate()

    if (result.hasUpdate) {
      // 有更新可用
      updateStatus.value = { type: 'danger', text: '有新版本' }

      showDialog({
        title: '发现新版本',
        message: `当前版本: ${currentVersion.value}\n最新版本: ${result.version}\n\n${result.message || '建议更新到最新版本以获得更好的体验'}`,
        confirmButtonText: '立即更新',
        cancelButtonText: '稍后再说',
        confirmButtonColor: '#1989fa',
        showCancelButton: true,
        teleport: 'body'
      }).then(() => {
        // 用户点击立即更新
        downloadAndInstall()
        updateStatus.value = { type: 'success', text: '已下载' }
      }).catch(() => {
        // 用户取消
        updateStatus.value = { type: 'primary', text: '待更新' }
      })
    } else {
      // 已是最新版本
      updateStatus.value = { type: 'success', text: '已是最新' }
      showSuccessToast('当前已是最新版本')
    }
  } catch (error) {
    logger.error('检查更新失败:', error)
    updateStatus.value = { type: 'danger', text: '检查失败' }
    showFailToast('检查更新失败，请稍后重试')
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

.avatar-wrapper {
  position: relative;
  margin-right: 12px;
}

.avatar-edit-icon {
  position: absolute;
  bottom: 0;
  right: 0;
  background-color: #fff;
  border-radius: 50%;
  padding: 2px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
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
  margin: 0 0 4px 0;
}

.about-date {
  font-size: 12px;
  color: #c8c9cc;
  margin: 0 0 12px 0;
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
