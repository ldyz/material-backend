<template>
  <div class="notification-center">
    <!-- 通知图标按钮 -->
    <van-badge
      :content="unreadCount > 99 ? '99+' : unreadCount"
      :show-zero="false"
      @click="showNotifications = true"
    >
      <van-icon
        name="bell"
        size="20"
        color="#323233"
      />
    </van-badge>

    <!-- 通知列表弹窗 -->
    <van-popup
      v-model:show="showNotifications"
      position="right"
      :style="{ width: '100%', height: '100%' }"
    >
      <div class="notifications-page">
        <van-nav-bar
          title="通知中心"
          :left-text="unreadCount > 0 ? '全部已读' : ''"
          left-arrow
          @click-left="handleClose"
        >
          <template #right>
            <van-icon
              name="delete-o"
              size="18"
              @click="handleClearAll"
            />
          </template>
        </van-nav-bar>

        <!-- 通知筛选 -->
        <van-tabs v-model:active="activeTab" sticky>
          <van-tab title="全部" title-style="font-size: 14px">
            <NotificationList
              :notifications="filteredNotifications"
              :unread-only="false"
              @read="handleRead"
              @delete="handleDelete"
            />
          </van-tab>
          <van-tab title="未读" title-style="font-size: 14px">
            <NotificationList
              :notifications="filteredNotifications"
              :unread-only="true"
              @read="handleRead"
              @delete="handleDelete"
            />
          </van-tab>
        </van-tabs>
      </div>
    </van-popup>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { showConfirmDialog, showToast } from 'vant'
import { useRouter } from 'vue-router'
import { useNotification } from '@/composables/useNotification'

const router = useRouter()

const {
  notifications,
  unreadCount,
  markAsRead,
  markAllAsRead,
  removeNotification,
  clearNotifications,
  NOTIFICATION_TYPES,
} = useNotification()

const showNotifications = ref(false)
const activeTab = ref(0)
const filterType = ref('')

// 监听弹窗状态变化，处理返回键
watch(showNotifications, (newValue) => {
  if (newValue) {
    // 弹窗打开时，添加一个历史记录
    history.pushState({ notificationOpen: true }, '', '#notification')
    // 监听返回键
    window.addEventListener('popstate', handlePopState)
  } else {
    // 弹窗关闭时，移除监听
    window.removeEventListener('popstate', handlePopState)
  }
})

// 处理返回键
function handlePopState(event) {
  if (showNotifications.value) {
    event.stopPropagation()
    showNotifications.value = false
    // 移除添加的历史记录
    history.back()
  }
}

// 组件卸载时清理
onUnmounted(() => {
  window.removeEventListener('popstate', handlePopState)
})

// 筛选后的通知列表
const filteredNotifications = computed(() => {
  let filtered = notifications.value

  // 按类型筛选
  if (filterType.value) {
    filtered = filtered.filter(n => n.type === filterType.value)
  }

  return filtered
})

// 标记已读并跳转
async function handleRead(id) {
  const notification = notifications.value.find(n => n.id === id)
  if (!notification) return

  // 标记为已读
  await markAsRead(id)

  // 关闭通知中心
  showNotifications.value = false

  // 根据通知类型跳转
  const data = notification.data || {}
  switch (notification.type) {
    case NOTIFICATION_TYPES.APPROVAL_PENDING:
      if (data.entity_type === 'inbound') {
        router.push(`/inbound/${data.entity_id}/approve`)
      } else if (data.entity_type === 'requisition') {
        router.push(`/outbound/${data.entity_id}/approve`)
      }
      break

    case NOTIFICATION_TYPES.APPROVAL_APPROVED:
    case NOTIFICATION_TYPES.APPROVAL_REJECTED:
    case NOTIFICATION_TYPES.STATUS_CHANGED:
      if (data.entity_type === 'inbound') {
        router.push(`/inbound/${data.entity_id}`)
      } else if (data.entity_type === 'requisition') {
        router.push(`/outbound/${data.entity_id}`)
      }
      break

    case NOTIFICATION_TYPES.STOCK_ALERT:
      router.push('/stock')
      break

    default:
      break
  }
}

// 关闭通知中心
function handleClose() {
  showNotifications.value = false
}

// 标记全部已读
function handleMarkAllRead() {
  if (unreadCount.value === 0) return

  showConfirmDialog({
    title: '确认',
    message: '确定要将所有通知标记为已读吗？',
  })
    .then(() => {
      markAllAsRead()
      showToast({
        type: 'success',
        message: '已全部标记为已读',
      })
    })
    .catch(() => {
      // 取消
    })
}

// 删除单个通知
function handleDelete(id) {
  showConfirmDialog({
    title: '删除通知',
    message: '确定要删除这条通知吗？',
  })
    .then(() => {
      removeNotification(id)
      showToast({
        type: 'success',
        message: '删除成功',
      })
    })
    .catch(() => {
      // 取消
    })
}

// 清空所有通知
function handleClearAll() {
  if (notifications.value.length === 0) {
    showToast('暂无通知')
    return
  }

  showConfirmDialog({
    title: '清空通知',
    message: '确定要清空所有通知吗？',
  })
    .then(() => {
      clearNotifications()
      showToast({
        type: 'success',
        message: '已清空',
      })
    })
    .catch(() => {
      // 取消
    })
}
</script>

<script>
// 通知列表子组件
import { defineComponent, h } from 'vue'
import { formatRelativeTime } from '@/utils/date'

const NotificationList = defineComponent({
  props: {
    notifications: Array,
    unreadOnly: Boolean,
  },
  emits: ['read', 'delete'],
  setup(props, { emit }) {
    return () => {
      const list = props.unreadOnly
        ? props.notifications.filter(n => !n.read)
        : props.notifications

      if (list.length === 0) {
        return h('div', { class: 'empty-notifications' }, [
          h('van-empty', { description: '暂无通知' }),
        ])
      }

      return h('div', { class: 'notifications-list' },
        list.map(notification =>
          h('div', {
            class: [
              'notification-item',
              { unread: !notification.read },
            ],
            onClick: () => emit('read', notification.id),
          }, [
            // 左侧图标
            h('div', { class: 'notification-icon' }, [
              h('div', { 
                class: ['icon-circle', 'icon-' + getTagType(notification.type)]
              }, [
                h('van-icon', {
                  name: getNotificationIcon(notification.type),
                  size: '20',
                }),
              ]),
            ]),

            // 中间内容
            h('div', { class: 'notification-content' }, [
              h('div', { class: 'notification-header' }, [
                h('span', { class: 'notification-type' }, getTypeText(notification.type)),
                h('span', { class: 'notification-time' }, formatRelativeTime(notification.created_at)),
              ]),
              h('div', { class: 'notification-title' }, notification.title),
              h('div', { class: 'notification-text' }, notification.content),
            ]),

            // 未读标识
            !notification.read ? h('div', { class: 'notification-dot' }) : null,
          ])
        )
      )
    }
  },
})

function getTagType(type) {
  const typeMap = {
    approval_pending: 'warning',
    approval_approved: 'success',
    approval_rejected: 'danger',
    status_changed: 'primary',
    stock_alert: 'warning',
    system: 'default',
  }
  return typeMap[type] || 'default'
}

function getTypeText(type) {
  const textMap = {
    approval_pending: '待审批',
    approval_approved: '已通过',
    approval_rejected: '已拒绝',
    status_changed: '状态变更',
    stock_alert: '库存预警',
    system: '系统通知',
  }
  return textMap[type] || '通知'
}

function getNotificationIcon(type) {
  const iconMap = {
    approval_pending: 'clock-o',
    approval_approved: 'checked',
    approval_rejected: 'close',
    status_changed: 'refresh',
    stock_alert: 'warning-o',
    system: 'bell',
  }
  return iconMap[type] || 'bell'
}

function getNotificationColor(type) {
  const colorMap = {
    approval_pending: '#ff976a',
    approval_approved: '#07c160',
    approval_rejected: '#ee0a24',
    status_changed: '#1989fa',
    stock_alert: '#ff976a',
    system: '#969799',
  }
  return colorMap[type] || '#969799'
}
</script>

<style scoped>
.notification-center {
  display: inline-block;
}

.notifications-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: #f7f8fa;
}

.notifications-page :deep(.van-tabs__content) {
  flex: 1;
  overflow-y: auto;
}

.notifications-list {
  padding: 0;
  background: #f7f8fa;
}

/* 仿手机通知样式 */
.notification-item {
  position: relative;
  display: flex;
  align-items: flex-start;
  padding: 12px 16px;
  background: white;
  border-bottom: 1px solid #f0f0f0;
  transition: background-color 0.2s;
  cursor: pointer;
}

.notification-item:active {
  background-color: #f5f5f5;
}

.notification-item.unread {
  background: linear-gradient(90deg, #f0f9ff 0%, #ffffff 100%);
}

.notification-item.unread::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 3px;
  background: #1989fa;
}

/* 左侧图标 */
.notification-icon {
  margin-right: 12px;
  flex-shrink: 0;
}

:deep(.icon-circle) {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

:deep(.icon-warning) {
  background: linear-gradient(135deg, #ff976a, #ff6b35);
}

:deep(.icon-success) {
  background: linear-gradient(135deg, #07c160, #00d976);
}

:deep(.icon-danger) {
  background: linear-gradient(135deg, #ee0a24, #ff4757);
}

:deep(.icon-primary) {
  background: linear-gradient(135deg, #1989fa, #3b9eff);
}

:deep(.icon-default) {
  background: linear-gradient(135deg, #969799, #b0b0b0);
}

/* 中间内容区 */
.notification-content {
  flex: 1;
  min-width: 0;
  padding-right: 8px;
}

.notification-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.notification-type {
  font-size: 13px;
  font-weight: 600;
  color: #323233;
}

.notification-time {
  font-size: 11px;
  color: #969799;
  white-space: nowrap;
  margin-left: 8px;
}

.notification-title {
  font-size: 14px;
  font-weight: 500;
  color: #323233;
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.notification-text {
  font-size: 13px;
  color: #646566;
  line-height: 1.4;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
}

/* 未读标识 */
.notification-dot {
  position: absolute;
  right: 16px;
  top: 50%;
  transform: translateY(-50%);
  width: 8px;
  height: 8px;
  background: #ee0a24;
  border-radius: 50%;
  box-shadow: 0 0 0 2px rgba(238, 10, 36, 0.2);
}

.empty-notifications {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 200px;
}
</style>
