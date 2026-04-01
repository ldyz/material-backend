<template>
  <view class="index-page">
    <!-- 顶部用户信息 -->
    <view class="header">
      <view class="user-info">
        <image class="avatar" :src="userAvatar" mode="aspectFill" />
        <view class="user-text">
          <text class="username">{{ username }}</text>
          <text class="welcome">欢迎使用项目管理系统</text>
        </view>
      </view>
    </view>

    <!-- 功能入口 -->
    <view class="menu-grid">
      <view class="menu-item" v-for="item in menuList" :key="item.id" @click="navigateTo(item)">
        <view class="menu-icon" :style="{ backgroundColor: item.color }">
          <text class="icon-text">{{ item.icon }}</text>
        </view>
        <text class="menu-text">{{ item.title }}</text>
      </view>
    </view>

    <!-- 待办事项 -->
    <view class="section">
      <view class="section-header">
        <text class="section-title">待办事项</text>
      </view>
      <view class="todo-list" v-if="todoList.length > 0">
        <view class="todo-item" v-for="(item, index) in todoList" :key="index" @click="handleTodoClick(item)">
          <view class="todo-info">
            <text class="todo-title">{{ item.title }}</text>
            <text class="todo-desc">{{ item.desc }}</text>
          </view>
          <StatusTag :text="item.statusText" :color="item.statusColor" plain />
        </view>
      </view>
      <view class="empty-todo" v-else>
        <text>暂无待办事项</text>
      </view>
    </view>
  </view>
</template>

<script>
import { useAuthStore } from '@/stores/auth'
import StatusTag from '@/components/StatusTag.vue'

export default {
  components: {
    StatusTag
  },
  data() {
    return {
      menuList: [
        { id: 'plan', title: '物资计划', icon: '📋', color: '#1989fa', url: '/pages/plan/list' },
        { id: 'inbound', title: '入库管理', icon: '📥', color: '#07c160', url: '/pages/inbound/list' },
        { id: 'requisition', title: '领料管理', icon: '📤', color: '#ff976a', url: '/pages/requisition/list' },
        { id: 'appointment', title: '预约管理', icon: '📅', color: '#7232dd', url: '/pages/appointment/list' },
        { id: 'attendance', title: '考勤打卡', icon: '⏰', color: '#ee0a24', url: '/pages/attendance/clockIn' },
        { id: 'stock', title: '库存查询', icon: '📦', color: '#909399', url: '/pages/stock/list' }
      ],
      todoList: []
    }
  },
  computed: {
    username() {
      const authStore = useAuthStore()
      return authStore.username || '用户'
    },
    userAvatar() {
      const authStore = useAuthStore()
      return authStore.user?.avatar || '/static/images/default-avatar.png'
    }
  },
  onShow() {
    this.initData()
  },
  methods: {
    async initData() {
      const authStore = useAuthStore()
      if (!authStore.isAuthenticated) {
        uni.redirectTo({ url: '/pages/login/index' })
        return
      }

      // 刷新用户信息
      try {
        await authStore.initAuth()
      } catch (e) {
        console.error('刷新用户信息失败', e)
      }
    },
    navigateTo(item) {
      uni.navigateTo({ url: item.url })
    },
    handleTodoClick(item) {
      if (item.url) {
        uni.navigateTo({ url: item.url })
      }
    }
  }
}
</script>

<style scoped>
.index-page {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.header {
  background: linear-gradient(135deg, #1989fa 0%, #07c160 100%);
  padding: 48rpx 32rpx;
}

.user-info {
  display: flex;
  align-items: center;
}

.avatar {
  width: 96rpx;
  height: 96rpx;
  border-radius: 50%;
  border: 4rpx solid rgba(255, 255, 255, 0.3);
}

.user-text {
  margin-left: 24rpx;
}

.username {
  display: block;
  font-size: 36rpx;
  font-weight: 600;
  color: #ffffff;
  margin-bottom: 8rpx;
}

.welcome {
  display: block;
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.8);
}

.menu-grid {
  display: flex;
  flex-wrap: wrap;
  background-color: #ffffff;
  margin: 24rpx;
  border-radius: 16rpx;
  padding: 24rpx 0;
}

.menu-item {
  width: 33.33%;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 24rpx 0;
}

.menu-icon {
  width: 96rpx;
  height: 96rpx;
  border-radius: 24rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16rpx;
}

.icon-text {
  font-size: 48rpx;
}

.menu-text {
  font-size: 28rpx;
  color: #323233;
}

.section {
  background-color: #ffffff;
  margin: 0 24rpx 24rpx;
  border-radius: 16rpx;
  padding: 24rpx;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24rpx;
}

.section-title {
  font-size: 32rpx;
  font-weight: 500;
  color: #323233;
}

.todo-list {
  padding: 0;
}

.todo-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20rpx 0;
  border-bottom: 1rpx solid #ebedf0;
}

.todo-item:last-child {
  border-bottom: none;
}

.todo-info {
  flex: 1;
  min-width: 0;
}

.todo-title {
  display: block;
  font-size: 28rpx;
  color: #323233;
  margin-bottom: 8rpx;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.todo-desc {
  display: block;
  font-size: 24rpx;
  color: #969799;
}

.empty-todo {
  text-align: center;
  padding: 48rpx 0;
  color: #969799;
  font-size: 28rpx;
}
</style>
