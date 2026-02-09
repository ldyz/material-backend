<template>
  <div class="notification-list-container">
    <div class="notification-toolbar">
      <el-radio-group v-model="filterType" size="small" @change="handleFilterChange">
        <el-radio-button label="all">全部 ({{ pagination.total }})</el-radio-button>
        <el-radio-button label="unread">未读 ({{ unreadCount }})</el-radio-button>
        <el-radio-button label="read">已读</el-radio-button>
      </el-radio-group>

      <div class="toolbar-actions">
        <el-button
          v-if="unreadCount > 0"
          type="primary"
          size="small"
          :icon="Check"
          @click="markAllAsRead"
        >
          全部已读
        </el-button>
        <el-button
          type="danger"
          size="small"
          :icon="Delete"
          @click="handleClearAll"
        >
          清空全部
        </el-button>
        <el-button
          size="small"
          :icon="Refresh"
          @click="refresh"
          :loading="loading"
        >
          刷新
        </el-button>
      </div>
    </div>

    <div v-if="loading && notifications.length === 0" class="loading-container">
      <el-icon class="is-loading" :size="32"><Loading /></el-icon>
      <p>加载中...</p>
    </div>

    <div v-else-if="notifications.length === 0" class="empty-container">
      <el-empty
        :description="filterType === 'unread' ? '暂无未读通知' : '暂无通知'"
        :image-size="120"
      />
    </div>

    <div v-else class="notification-items">
      <TransitionGroup name="list">
        <NotificationItem
          v-for="notification in notifications"
          :key="notification.id"
          :notification="notification"
          @mark-read="handleMarkRead"
          @delete="handleDelete"
          @click="handleClick"
        />
      </TransitionGroup>
    </div>

    <div v-if="pagination.total > 0" class="pagination-container">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useNotificationStore } from '@/stores/notificationStore'
import NotificationItem from './NotificationItem.vue'
import {
  Check,
  Delete,
  Refresh,
  Loading
} from '@element-plus/icons-vue'
import { ElMessageBox, ElMessage } from 'element-plus'

const emit = defineEmits(['notification-click'])

const notificationStore = useNotificationStore()
const { notifications, unreadCount, loading, pagination } = storeToRefs(notificationStore)

// 筛选类型
const filterType = ref('all')

// 获取通知列表
const fetchNotifications = async () => {
  const params = {
    page: pagination.value.page,
    page_size: pagination.value.pageSize
  }

  if (filterType.value === 'unread') {
    params.unread_only = true
  } else if (filterType.value === 'read') {
    // 已读筛选，可以通过 unread_only=false 并在前端过滤
  }

  await notificationStore.fetchNotifications(params)
}

// 筛选类型变化
const handleFilterChange = () => {
  pagination.value.page = 1
  fetchNotifications()
}

// 标记单个通知为已读
const handleMarkRead = (id) => {
  notificationStore.markAsRead(id)
}

// 标记所有通知为已读
const markAllAsRead = () => {
  ElMessageBox.confirm('确定要将所有通知标记为已读吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'info'
  }).then(() => {
    notificationStore.markAllAsRead()
  }).catch(() => {})
}

// 删除通知
const handleDelete = (id) => {
  ElMessageBox.confirm('确定要删除这条通知吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    notificationStore.deleteNotification(id)
  }).catch(() => {})
}

// 清空所有通知
const handleClearAll = () => {
  const message = filterType.value === 'unread'
    ? '确定要清空所有未读通知吗？'
    : '确定要清空所有通知吗？'

  ElMessageBox.confirm(message, '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    notificationStore.clearAll()
  }).catch(() => {})
}

// 通知点击事件
const handleClick = (notification) => {
  emit('notification-click', notification)
}

// 刷新
const refresh = () => {
  fetchNotifications()
  notificationStore.fetchUnreadCount()
}

// 分页大小变化
const handleSizeChange = (size) => {
  pagination.value.pageSize = size
  pagination.value.page = 1
  fetchNotifications()
}

// 页码变化
const handlePageChange = (page) => {
  pagination.value.page = page
  fetchNotifications()
}

// 初始化
onMounted(() => {
  fetchNotifications()
  notificationStore.fetchUnreadCount()
})
</script>

<style scoped>
.notification-list-container {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.notification-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background-color: #fff;
  border-radius: 4px;
  margin-bottom: 16px;
  flex-shrink: 0;
}

.toolbar-actions {
  display: flex;
  gap: 8px;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: #909399;
}

.loading-container p {
  margin-top: 16px;
  font-size: 14px;
}

.empty-container {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
}

.notification-items {
  flex: 1;
  background-color: #fff;
  border-radius: 4px;
  overflow: hidden;
}

.pagination-container {
  display: flex;
  justify-content: center;
  padding: 16px 0;
  flex-shrink: 0;
}

:deep(.el-radio-button__inner) {
  border-color: #dcdfe6;
  padding: 8px 15px;
}

:deep(.el-radio-button:first-child .el-radio-button__inner) {
  border-left: 1px solid #dcdfe6;
}

/* 列表动画 */
.list-enter-active,
.list-leave-active {
  transition: all 0.3s ease;
}

.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateX(-20px);
}
</style>
