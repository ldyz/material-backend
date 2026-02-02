<template>
  <div class="notifications-page">
    <van-nav-bar
      title="通知中心"
      left-arrow
      @click-left="onClickLeft"
    >
      <template #right>
        <van-button
          v-if="unreadCount > 0"
          size="small"
          type="primary"
          plain
          round
          @click="markAllRead"
        >
          全部已读
        </van-button>
      </template>
    </van-nav-bar>

    <!-- 未读统计 -->
    <div v-if="unreadCount > 0" class="unread-stats">
      <van-icon name="bell" />
      <span>您有 <strong>{{ unreadCount }}</strong> 条未读通知</span>
    </div>

    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <van-list
        v-model:loading="loading"
        :finished="finished"
        finished-text="没有更多了"
        @load="onLoad"
      >
        <div
          v-for="item in list"
          :key="item.id"
          class="notification-card"
          :class="{ unread: !item.is_read }"
          @click="handleClick(item)"
        >
          <!-- 卡片头部 -->
          <div class="card-header">
            <div class="header-left">
              <div class="icon-wrapper" :class="'icon-' + getTypeColor(item.type)">
                <van-icon :name="getTypeIcon(item.type)" size="18" />
              </div>
              <div class="header-info">
                <van-tag
                  :type="getTypeColor(item.type)"
                  size="medium"
                  plain
                >
                  {{ getTypeText(item.type) }}
                </van-tag>
              </div>
            </div>
            <div class="header-right">
              <span class="time-text">{{ formatTime(item.created_at) }}</span>
              <div v-if="!item.is_read" class="unread-badge"></div>
            </div>
          </div>
          
          <!-- 卡片内容表格 -->
          <div class="card-table">
            <div class="table-row">
              <div class="table-label">标题</div>
              <div class="table-value">{{ item.title }}</div>
            </div>
            <div class="table-row">
              <div class="table-label">内容</div>
              <div class="table-value">{{ item.content }}</div>
            </div>
            <div class="table-row">
              <div class="table-label">时间</div>
              <div class="table-value">{{ new Date(item.created_at).toLocaleString('zh-CN') }}</div>
            </div>
            <div class="table-row" v-if="!item.is_read">
              <div class="table-label">状态</div>
              <div class="table-value status-unread">
                <van-tag type="danger" size="small">未读</van-tag>
              </div>
            </div>
          </div>
        </div>
      </van-list>
    </van-pull-refresh>

    <van-empty
      v-if="list.length === 0 && !loading"
      description="暂无通知"
      image="search"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'
import { getNotifications, markAsRead, markAllAsRead, getUnreadCount } from '@/api/notification'

const router = useRouter()

const list = ref([])
const loading = ref(false)
const finished = ref(false)
const refreshing = ref(false)
const page = ref(1)
const unreadCount = ref(0)

// 获取通知类型图标
function getTypeIcon(type) {
  const iconMap = {
    requisition_approve: 'notes-o',
    requisition_approved: 'passed',
    requisition_rejected: 'close',
    inbound_approve: 'logistics',
    inbound_approved: 'passed',
    inbound_rejected: 'close',
    stock_alert: 'warning-o',
    system: 'info-o'
  }
  return iconMap[type] || 'info-o'
}

// 获取通知类型颜色
function getTypeColor(type) {
  const colorMap = {
    requisition_approve: 'warning',
    requisition_approved: 'success',
    requisition_rejected: 'danger',
    inbound_approve: 'primary',
    inbound_approved: 'success',
    inbound_rejected: 'danger',
    stock_alert: 'danger',
    system: 'default'
  }
  return colorMap[type] || 'default'
}

// 获取通知类型文本
function getTypeText(type) {
  const textMap = {
    requisition_approve: '出库审批',
    requisition_approved: '出库已通过',
    requisition_rejected: '出库已拒绝',
    inbound_approve: '入库审批',
    inbound_approved: '入库已通过',
    inbound_rejected: '入库已拒绝',
    stock_alert: '库存预警',
    system: '系统通知'
  }
  return textMap[type] || '通知'
}

// 格式化时间
function formatTime(timeStr) {
  const time = new Date(timeStr)
  const now = new Date()
  const diff = now - time

  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return Math.floor(diff / 60000) + '分钟前'
  if (diff < 86400000) return Math.floor(diff / 3600000) + '小时前'
  if (diff < 604800000) return Math.floor(diff / 86400000) + '天前'

  return time.toLocaleDateString()
}

// 加载通知列表
async function onLoad() {
  try {
    // 适配统一响应格式
    const { data, meta } = await getNotifications({
      page: page.value,
      per_page: 20
    })
    const notifications = data || []

    if (refreshing.value) {
      list.value = notifications
      refreshing.value = false
    } else {
      list.value.push(...notifications)
    }

    if (notifications.length < 20) {
      finished.value = true
    } else {
      page.value++
    }

    // 更新未读数量（从meta中获取）
    unreadCount.value = meta?.unread_count || 0
  } catch (error) {
    showToast({
      type: 'fail',
      message: '加载失败'
    })
  } finally {
    loading.value = false
  }
}

// 下拉刷新
async function onRefresh() {
  finished.value = false
  loading.value = true
  page.value = 1
  await onLoad()
}

// 点击通知
async function handleClick(item) {
  // 标记为已读
  if (!item.is_read) {
    try {
      await markAsRead(item.id)
      item.is_read = true
      unreadCount.value--
    } catch (error) {
      console.error('标记已读失败:', error)
    }
  }

  // 解析通知数据并跳转
  try {
    const data = typeof item.data === 'string' ? JSON.parse(item.data) : (item.data || {})

    // 出库单相关通知
    if (item.type === 'requisition_approve' || item.type === 'requisition_approved' || item.type === 'requisition_rejected') {
      if (data.requisition_id) {
        // 如果是待审批，跳转到审批页面
        if (item.type === 'requisition_approve') {
          router.push(`/outbound/${data.requisition_id}/approve`)
        } else {
          // 已审批或已拒绝，跳转到详情页
          router.push(`/outbound/${data.requisition_id}`)
        }
      }
    }
    // 入库单相关通知
    else if (item.type === 'inbound_approve' || item.type === 'inbound_approved' || item.type === 'inbound_rejected') {
      if (data.inbound_id) {
        // 如果是待审批，跳转到审批页面
        if (item.type === 'inbound_approve') {
          router.push(`/inbound/${data.inbound_id}/approve`)
        } else {
          // 已审批或已拒绝，跳转到详情页
          router.push(`/inbound/${data.inbound_id}`)
        }
      }
    }
    // 库存预警通知
    else if (item.type === 'stock_alert') {
      if (data.stock_id) {
        router.push(`/stock/${data.stock_id}`)
      } else {
        // 跳转到库存列表
        router.push('/stock')
      }
    }
    // 系统通知
    else if (item.type === 'system') {
      // 系统通知暂时不跳转
      console.log('系统通知:', data)
    }
  } catch (error) {
    console.error('解析通知数据失败:', error)
  }
}

// 全部标记已读
async function markAllRead() {
  try {
    await markAllAsRead()
    list.value.forEach(item => {
      item.is_read = true
    })
    unreadCount.value = 0
    showToast({
      type: 'success',
      message: '已全部标记为已读'
    })
  } catch (error) {
    showToast({
      type: 'fail',
      message: '操作失败'
    })
  }
}

// 返回
function onClickLeft() {
  router.back()
}

onMounted(() => {
  onLoad()
})
</script>

<style scoped>
.notifications-page {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.unread-stats {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  font-size: 14px;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.2);
}

.unread-stats strong {
  font-size: 16px;
  font-weight: bold;
}

/* 通知卡片 */
.notification-card {
  margin: 12px 16px;
  background: white;
  border-radius: 12px;
  border: 2px solid #e5e7eb;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  overflow: hidden;
  transition: all 0.3s ease;
}

.notification-card:active {
  transform: scale(0.98);
}

.notification-card.unread {
  border-color: #1989fa;
  background: linear-gradient(to right, #f0f9ff 0%, #ffffff 100%);
  box-shadow: 0 2px 12px rgba(25, 137, 250, 0.2);
}

/* 卡片头部 */
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #f9fafb;
  border-bottom: 1px solid #e5e7eb;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
}

.icon-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 8px;
  flex-shrink: 0;
}

.icon-warning {
  background: linear-gradient(135deg, #ff9a44 0%, #fc6076 100%);
  color: white;
}

.icon-success {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  color: white;
}

.icon-danger {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  color: white;
}

.icon-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.icon-default {
  background: linear-gradient(135deg, #e0e7ff 0%, #f0f4ff 100%);
  color: #646566;
}

.header-info {
  flex: 1;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.time-text {
  font-size: 12px;
  color: #9ca3af;
  white-space: nowrap;
}

.unread-badge {
  width: 8px;
  height: 8px;
  background: #ee0a24;
  border-radius: 50%;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% {
    box-shadow: 0 0 0 0 rgba(238, 10, 36, 0.7);
  }
  70% {
    box-shadow: 0 0 0 6px rgba(238, 10, 36, 0);
  }
  100% {
    box-shadow: 0 0 0 0 rgba(238, 10, 36, 0);
  }
}

/* 卡片表格 */
.card-table {
  padding: 0;
}

.table-row {
  display: flex;
  padding: 12px 16px;
  border-bottom: 1px solid #f3f4f6;
  transition: background-color 0.2s;
}

.table-row:last-child {
  border-bottom: none;
}

.table-row:hover {
  background-color: #f9fafb;
}

.table-label {
  width: 60px;
  flex-shrink: 0;
  font-size: 13px;
  color: #6b7280;
  font-weight: 500;
  line-height: 1.6;
}

.table-value {
  flex: 1;
  font-size: 13px;
  color: #323233;
  line-height: 1.6;
  word-break: break-word;
}

.status-unread {
  display: flex;
  align-items: center;
}

/* 空状态样式 */
:deep(.van-empty) {
  padding: 80px 0;
}
</style>
