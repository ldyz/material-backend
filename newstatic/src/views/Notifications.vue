<template>
  <div class="notifications-page">
    <el-page-header @back="goBack" class="page-header">
      <template #content>
        <div class="header-content">
          <el-icon class="header-icon"><Bell /></el-icon>
          <span class="header-title">通知中心</span>
        </div>
      </template>
    </el-page-header>

    <div class="page-content">
      <NotificationList @notification-click="handleNotificationClick" />
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { Bell } from '@element-plus/icons-vue'
import NotificationList from '@/components/Notification/NotificationList.vue'

const router = useRouter()

// 返回
const goBack = () => {
  router.back()
}

// 处理通知点击
const handleNotificationClick = (notification) => {
  // 导航逻辑已经在 NotificationItem 中处理
  console.log('Notification clicked:', notification)
}
</script>

<style scoped>
.notifications-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: #f5f7fa;
}

.page-header {
  background-color: #fff;
  padding: 16px 24px;
  border-radius: 4px;
  margin-bottom: 16px;
  flex-shrink: 0;
}

.header-content {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-icon {
  font-size: 20px;
  color: #409eff;
}

.header-title {
  font-size: 18px;
  font-weight: 500;
  color: #303133;
}

.page-content {
  flex: 1;
  min-height: 0;
}
</style>
