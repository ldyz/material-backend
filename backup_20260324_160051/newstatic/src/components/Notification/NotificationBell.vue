<template>
  <div class="notification-bell">
    <el-dropdown
      trigger="click"
      placement="bottom-end"
      @visible-change="handleDropdownVisible"
    >
      <div class="bell-trigger" @click="handleClick">
        <el-badge :value="unreadCount" :hidden="unreadCount === 0" :max="99">
          <el-icon class="bell-icon" :size="20">
            <Bell />
          </el-icon>
        </el-badge>
      </div>
      <template #dropdown>
        <el-dropdown-menu class="notification-dropdown">
          <div class="notification-header">
            <div class="notification-title">
              <span>通知</span>
              <el-badge
                v-if="unreadCount > 0"
                :value="unreadCount"
                :max="99"
                class="header-badge"
              />
            </div>
            <div class="notification-actions">
              <el-button
                link
                type="primary"
                size="small"
                :disabled="unreadCount === 0"
                @click.stop="markAllAsRead"
              >
                全部已读
              </el-button>
              <el-button
                link
                type="primary"
                size="small"
                @click.stop="goToNotifications"
              >
                查看全部
              </el-button>
            </div>
          </div>

          <el-scrollbar max-height="400px">
            <div v-if="loading" class="notification-loading">
              <el-icon class="is-loading"><Loading /></el-icon>
              <span>加载中...</span>
            </div>
            <div v-else-if="displayNotifications.length === 0" class="notification-empty">
              <el-empty description="暂无通知" :image-size="80" />
            </div>
            <div v-else class="notification-list">
              <NotificationItem
                v-for="notification in displayNotifications"
                :key="notification.id"
                :notification="notification"
                @mark-read="handleMarkRead"
                @delete="handleDelete"
                @click="handleNotificationClick"
              />
            </div>
          </el-scrollbar>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useNotificationStore } from '@/stores/notificationStore'
import { Bell, Loading } from '@element-plus/icons-vue'
import NotificationItem from './NotificationItem.vue'

const props = defineProps({
  // 显示的通知数量
  displayCount: {
    type: Number,
    default: 5
  }
})

const emit = defineEmits(['click'])

const router = useRouter()
const notificationStore = useNotificationStore()
const { notifications, unreadCount, loading } = storeToRefs(notificationStore)

// 显示的通知列表
const displayNotifications = computed(() => {
  return notifications.value.slice(0, props.displayCount)
})

// 点击铃铛
const handleClick = () => {
  emit('click')
}

// 下拉框显示/隐藏事件
const handleDropdownVisible = (visible) => {
  if (visible) {
    // 打开时刷新通知列表
    notificationStore.fetchNotifications({ page_size: props.displayCount })
  }
}

// 标记单个通知为已读
const handleMarkRead = (id) => {
  notificationStore.markAsRead(id)
}

// 标记所有通知为已读
const markAllAsRead = () => {
  notificationStore.markAllAsRead()
}

// 删除通知
const handleDelete = (id) => {
  notificationStore.deleteNotification(id)
}

// 通知点击事件
const handleNotificationClick = (notification) => {
  // 关闭下拉框（点击外部自动关闭）
}

// 跳转到通知页面
const goToNotifications = () => {
  router.push('/notifications')
}
</script>

<style scoped>
.notification-bell {
  display: inline-block;
}

.bell-trigger {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  cursor: pointer;
  border-radius: 50%;
  transition: background-color 0.2s;
}

.bell-trigger:hover {
  background-color: #f5f7fa;
}

.bell-icon {
  color: #606266;
  transition: color 0.2s;
}

.bell-trigger:hover .bell-icon {
  color: #409eff;
}

.notification-dropdown {
  width: 380px;
  padding: 0;
}

.notification-header {
  padding: 12px 16px;
  border-bottom: 1px solid #ebeef5;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.notification-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 500;
  color: #303133;
}

.notification-actions {
  display: flex;
  gap: 8px;
}

.notification-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  color: #909399;
  gap: 8px;
}

.notification-empty {
  padding: 20px;
}

.notification-list {
  padding: 0;
}

:deep(.el-dropdown-menu__item) {
  padding: 0;
}

:deep(.el-scrollbar__wrap) {
  max-height: 400px;
}
</style>
