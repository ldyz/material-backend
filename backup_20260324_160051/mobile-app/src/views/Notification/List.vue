<template>
  <div class="notification-list-page">
    <van-nav-bar title="通知中心" left-arrow @click-left="onClickLeft" />

    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <van-list
        v-model:loading="loading"
        :finished="finished"
        finished-text="没有更多了"
        @load="onLoad"
      >
        <van-cell-group inset class="notification-group">
          <van-cell
            v-for="notification in notifications"
            :key="notification.id"
            class="notification-item"
            :class="{ 'is-unread': !notification.is_read }"
            is-link
            @click="handleClick(notification)"
          >
            <template #title>
              <div class="notification-header">
                <van-tag :type="getNotificationType(notification.type)" size="small">
                  {{ getNotificationTypeText(notification.type) }}
                </van-tag>
                <span class="notification-time">{{ formatTime(notification.created_at) }}</span>
              </div>
            </template>
            <template #label>
              <div class="notification-title">{{ notification.title }}</div>
              <div class="notification-content">{{ notification.content }}</div>
            </template>
            <template #icon>
              <van-icon :name="getNotificationIcon(notification.type)" :color="getNotificationIconColor(notification.type)" size="24" />
            </template>
          </van-cell>
        </van-cell-group>

        <van-empty v-if="!loading && notifications.length === 0" description="暂无通知" />
      </van-list>
    </van-pull-refresh>

    <div class="action-bar">
      <van-button
        v-if="unreadCount > 0"
        type="primary"
        size="small"
        @click="handleMarkAllRead"
      >
        全部已读 ({{ unreadCount }})
      </van-button>
      <van-button
        v-if="notifications.length > 0"
        type="danger"
        size="small"
        @click="handleClearAll"
      >
        清空全部
      </van-button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useNotificationStore } from '@/stores/notification'
import { showToast, showConfirmDialog } from 'vant'
import { storeToRefs } from 'pinia'

const router = useRouter()
const notificationStore = useNotificationStore()
const { notifications, unreadCount } = storeToRefs(notificationStore)

const loading = ref(false)
const finished = ref(false)
const refreshing = ref(false)
const page = ref(1)
const pageSize = 20

// 返回
const onClickLeft = () => {
  router.back()
}

// 加载通知列表
const onLoad = async () => {
  if (refreshing.value) {
    notifications.value = []
    page.value = 1
    finished.value = false
  }

  loading.value = true
  try {
    await notificationStore.fetchNotifications({
      page: page.value,
      page_size: pageSize
    })

    if (notifications.value.length < pageSize) {
      finished.value = true
    } else {
      page.value++
    }
  } finally {
    loading.value = false
    refreshing.value = false
  }
}

// 刷新
const onRefresh = () => {
  finished.value = false
  onLoad()
}

// 获取通知图标
const getNotificationIcon = (type) => {
  if (type.includes('approve')) return 'checked'
  if (type.includes('reject')) return 'close'
  if (type.includes('stock_alert')) return 'warning-o'
  return 'bell'
}

// 获取通知图标颜色
const getNotificationIconColor = (type) => {
  if (type.includes('approve')) return '#07c160'
  if (type.includes('reject')) return '#ee0a24'
  if (type.includes('stock_alert')) return '#ff976a'
  return '#1989fa'
}

// 获取通知类型
const getNotificationType = (type) => {
  if (type.includes('approve')) return 'success'
  if (type.includes('reject')) return 'danger'
  if (type.includes('stock_alert')) return 'warning'
  return 'primary'
}

// 获取通知类型文本
const getNotificationTypeText = (type) => {
  if (type.includes('approve')) return '审批'
  if (type.includes('reject')) return '拒绝'
  if (type.includes('stock_alert')) return '预警'
  return '通知'
}

// 格式化时间
const formatTime = (timeStr) => {
  if (!timeStr) return ''
  const time = new Date(timeStr)
  const now = new Date()
  const diffMs = now - time
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMins / 60)
  const diffDays = Math.floor(diffHours / 24)

  if (diffMins < 1) return '刚刚'
  if (diffMins < 60) return `${diffMins}分钟前`
  if (diffHours < 24) return `${diffHours}小时前`
  if (diffDays < 7) return `${diffDays}天前`

  return time.toLocaleDateString()
}

// 点击通知
const handleClick = async (notification) => {
  // 标记为已读
  if (!notification.is_read) {
    await notificationStore.markAsRead(notification.id)
  }

  // 解析数据并跳转
  try {
    const data = JSON.parse(notification.data || '{}')
    if (data.business_type && data.business_id) {
      switch (data.business_type) {
        case 'inbound_order':
          router.push(`/inbound/${data.business_id}`)
          break
        case 'requisition':
          router.push(`/requisition/${data.business_id}`)
          break
        case 'material_plan':
          router.push(`/plan/${data.business_id}`)
          break
      }
    }
  } catch (e) {
    console.error('解析通知数据失败:', e)
  }
}

// 全部已读
const handleMarkAllRead = async () => {
  try {
    await showConfirmDialog({
      title: '提示',
      message: '确定要将所有通知标记为已读吗？'
    })
    await notificationStore.markAllAsRead()
  } catch (e) {
    // 用户取消
  }
}

// 清空全部
const handleClearAll = async () => {
  try {
    await showConfirmDialog({
      title: '警告',
      message: '确定要清空所有通知吗？'
    })
    await notificationStore.clearAll()
    notifications.value = []
    finished.value = true
  } catch (e) {
    // 用户取消
  }
}

onMounted(() => {
  onLoad()
  notificationStore.fetchUnreadCount()
})
</script>

<style scoped>
.notification-list-page {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.notification-group {
  margin: 12px;
}

.notification-item {
  position: relative;
}

.notification-item.is-unread {
  background-color: #f0f9ff;
}

.notification-item.is-unread::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 4px;
  height: 40px;
  background-color: #1989fa;
  border-radius: 0 4px 4px 0;
}

.notification-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.notification-time {
  font-size: 12px;
  color: #969799;
}

.notification-title {
  font-size: 14px;
  font-weight: 500;
  color: #323233;
  margin-bottom: 4px;
}

.notification-content {
  font-size: 13px;
  color: #646566;
  line-height: 1.5;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.action-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  justify-content: center;
  gap: 12px;
  padding: 12px 16px;
  background-color: #fff;
  border-top: 1px solid #ebedf0;
  z-index: 100;
}

.action-bar .van-button {
  flex: 1;
  max-width: 200px;
}
</style>
