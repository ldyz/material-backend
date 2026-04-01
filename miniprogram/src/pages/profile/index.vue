<template>
  <view class="profile-page">
    <!-- 用户信息卡片 -->
    <view class="user-card">
      <view class="user-header">
        <image class="avatar" :src="userAvatar" mode="aspectFill" @click="changeAvatar" />
        <view class="user-info">
          <text class="username">{{ username }}</text>
          <text class="role">{{ roleName }}</text>
        </view>
      </view>
    </view>

    <!-- 功能列表 -->
    <view class="menu-list">
      <view class="menu-group">
        <view class="menu-item" @click="navigateTo('/pages/attendance/recordList')">
          <text class="menu-label">打卡记录</text>
          <text class="menu-arrow">›</text>
        </view>
        <view class="menu-item" @click="navigateTo('/pages/notification/list')">
          <text class="menu-label">通知消息</text>
          <text class="menu-arrow">›</text>
        </view>
        <view class="menu-item" @click="navigateTo('/pages/about/index')">
          <text class="menu-label">关于</text>
          <text class="menu-arrow">›</text>
        </view>
      </view>
    </view>

    <!-- 退出登录 -->
    <view class="logout-wrapper">
      <button class="logout-btn" @click="handleLogout">退出登录</button>
    </view>
  </view>
</template>

<script>
import { useAuthStore } from '@/stores/auth'
import { chooseFromAlbum, uploadImage } from '@/utils/device'

export default {
  computed: {
    username() {
      const authStore = useAuthStore()
      return authStore.username || '用户'
    },
    userAvatar() {
      const authStore = useAuthStore()
      return authStore.user?.avatar || '/static/images/default-avatar.png'
    },
    roleName() {
      const authStore = useAuthStore()
      const roles = authStore.user?.roles
      if (roles && roles.length > 0) {
        return roles.map(r => r.name || r).join('、')
      }
      return '普通用户'
    }
  },
  methods: {
    navigateTo(url) {
      uni.navigateTo({ url })
    },
    async changeAvatar() {
      try {
        const result = await chooseFromAlbum(1)
        if (result.tempFilePaths && result.tempFilePaths.length > 0) {
          uni.showLoading({ title: '上传中...' })
          const uploadResult = await uploadImage(result.tempFilePaths[0])
          uni.hideLoading()

          if (uploadResult.data && uploadResult.data.url) {
            const authStore = useAuthStore()
            authStore.user.avatar = uploadResult.data.url
            authStore.setUser(authStore.user)
            uni.showToast({ title: '头像更新成功', icon: 'success' })
          }
        }
      } catch (error) {
        uni.hideLoading()
        uni.showToast({
          title: error.message || '上传失败',
          icon: 'none'
        })
      }
    },
    handleLogout() {
      uni.showModal({
        title: '提示',
        content: '确定要退出登录吗？',
        success: async (res) => {
          if (res.confirm) {
            const authStore = useAuthStore()
            await authStore.logout()
            uni.reLaunch({ url: '/pages/login/index' })
          }
        }
      })
    }
  }
}
</script>

<style scoped>
.profile-page {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.user-card {
  background: linear-gradient(135deg, #1989fa 0%, #07c160 100%);
  padding: 48rpx 32rpx;
  margin-bottom: 24rpx;
}

.user-header {
  display: flex;
  align-items: center;
}

.avatar {
  width: 128rpx;
  height: 128rpx;
  border-radius: 50%;
  border: 4rpx solid rgba(255, 255, 255, 0.3);
}

.user-info {
  margin-left: 24rpx;
}

.username {
  display: block;
  font-size: 40rpx;
  font-weight: 600;
  color: #ffffff;
  margin-bottom: 8rpx;
}

.role {
  display: block;
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.8);
}

.menu-list {
  padding: 0 24rpx;
}

.menu-group {
  background-color: #ffffff;
  border-radius: 16rpx;
  overflow: hidden;
  margin-bottom: 24rpx;
}

.menu-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 32rpx;
  border-bottom: 1rpx solid #ebedf0;
}

.menu-item:last-child {
  border-bottom: none;
}

.menu-item:active {
  background-color: #f7f8fa;
}

.menu-label {
  font-size: 32rpx;
  color: #323233;
}

.menu-arrow {
  font-size: 36rpx;
  color: #c8c9cc;
}

.logout-wrapper {
  padding: 48rpx 24rpx;
}

.logout-btn {
  width: 100%;
  height: 96rpx;
  background-color: #ffffff;
  border-radius: 16rpx;
  color: #ee0a24;
  font-size: 32rpx;
}

.logout-btn:active {
  background-color: #f7f8fa;
}
</style>
