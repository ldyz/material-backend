import { ref, onMounted, onUnmounted } from 'vue'
import { showToast, showNotify } from 'vant'
import { getNotifications, markAsRead as apiMarkAsRead, markAllAsRead as apiMarkAllAsRead, deleteNotification as apiDeleteNotification, clearAllNotifications as apiClearAllNotifications } from '@/api/notification'

// 通知类型
const NOTIFICATION_TYPES = {
  APPROVAL_PENDING: 'approval_pending', // 待审批提醒
  APPROVAL_APPROVED: 'approval_approved', // 审批通过
  APPROVAL_REJECTED: 'approval_rejected', // 审批拒绝
  STATUS_CHANGED: 'status_changed', // 状态变更
  STOCK_ALERT: 'stock_alert', // 库存预警
  SYSTEM: 'system', // 系统通知
}

// 通知存储键
const STORAGE_KEY = 'notifications'

let pollingTimer = null
let eventSource = null

export function useNotification() {
  const notifications = ref([])
  const unreadCount = ref(0)

  // 从后端加载通知
  async function loadNotifications() {
    try {
      const response = await getNotifications({ per_page: 50 })
      // 后端返回标准格式：{ success: true, data: [...], meta: { unread_count: 5 } }
      const serverNotifications = response.data || []

      // 合并本地和服务器通知，去重
      const localIds = new Set(notifications.value.map(n => n.id))
      serverNotifications.forEach(n => {
        if (!localIds.has(n.id)) {
          notifications.value.push({
            id: n.id,
            type: n.type,
            title: n.title,
            content: n.content,
            data: n.data ? JSON.parse(n.data) : {},
            read: n.is_read,
            created_at: n.created_at,
          })
        }
      })

      // 更新未读数量（从meta中获取）
      unreadCount.value = response.meta?.unread_count || 0
    } catch (error) {
      console.error('加载通知失败:', error)
    }
  }

  // 保存通知到本地
  function saveNotifications() {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify({
        notifications: notifications.value,
        unreadCount: unreadCount.value,
      }))
    } catch (error) {
      console.error('保存通知失败:', error)
    }
  }

  // 添加通知
  function addNotification(notification) {
    const newNotification = {
      id: Date.now(),
      type: notification.type || NOTIFICATION_TYPES.SYSTEM,
      title: notification.title,
      content: notification.content,
      data: notification.data || {},
      read: false,
      created_at: new Date().toISOString(),
    }

    notifications.value.unshift(newNotification)
    unreadCount.value++

    // 保存到本地
    saveNotifications()

    // 显示通知
    showNotify({
      type: getNotifyType(newNotification.type),
      message: newNotification.title,
      duration: 3000,
      onClick: () => {
        handleNotificationClick(newNotification)
      },
    })

    return newNotification
  }

  // 获取 Vant Notify 类型
  function getNotifyType(notificationType) {
    const typeMap = {
      [NOTIFICATION_TYPES.APPROVAL_PENDING]: 'warning',
      [NOTIFICATION_TYPES.APPROVAL_APPROVED]: 'success',
      [NOTIFICATION_TYPES.APPROVAL_REJECTED]: 'danger',
      [NOTIFICATION_TYPES.STATUS_CHANGED]: 'primary',
      [NOTIFICATION_TYPES.STOCK_ALERT]: 'warning',
      [NOTIFICATION_TYPES.SYSTEM]: 'primary',
    }
    return typeMap[notificationType] || 'primary'
  }

  // 处理通知点击
  function handleNotificationClick(notification) {
    // 标记为已读
    markAsRead(notification.id)

    // 根据通知类型跳转
    const router = window.router // 需要在 main.js 中设置
    if (!router) return

    switch (notification.type) {
      case NOTIFICATION_TYPES.APPROVAL_PENDING:
        if (notification.data.entity_type === 'inbound') {
          router.push(`/inbound/${notification.data.entity_id}/approve`)
        } else if (notification.data.entity_type === 'requisition') {
          router.push(`/outbound/${notification.data.entity_id}/approve`)
        }
        break

      case NOTIFICATION_TYPES.APPROVAL_APPROVED:
      case NOTIFICATION_TYPES.APPROVAL_REJECTED:
      case NOTIFICATION_TYPES.STATUS_CHANGED:
        if (notification.data.entity_type === 'inbound') {
          router.push(`/inbound/${notification.data.entity_id}`)
        } else if (notification.data.entity_type === 'requisition') {
          router.push(`/outbound/${notification.data.entity_id}`)
        }
        break

      case NOTIFICATION_TYPES.STOCK_ALERT:
        router.push(`/stock/${notification.data.stock_id}`)
        break

      default:
        break
    }
  }

  // 标记为已读
  async function markAsRead(id) {
    const notification = notifications.value.find(n => n.id === id)
    if (notification && !notification.read) {
      try {
        // 调用后端 API
        await apiMarkAsRead(id)
        notification.read = true
        unreadCount.value--
      } catch (error) {
        console.error('标记已读失败:', error)
      }
    }
  }

  // 标记全部已读
  async function markAllAsRead() {
    try {
      await apiMarkAllAsRead()
      notifications.value.forEach(n => {
        n.read = true
      })
      unreadCount.value = 0
    } catch (error) {
      console.error('标记全部已读失败:', error)
    }
  }

  // 删除通知
  async function removeNotification(id) {
    const index = notifications.value.findIndex(n => n.id === id)
    if (index > -1) {
      try {
        await apiDeleteNotification(id)
        if (!notifications.value[index].read) {
          unreadCount.value--
        }
        notifications.value.splice(index, 1)
      } catch (error) {
        console.error('删除通知失败:', error)
      }
    }
  }

  // 清空通知
  async function clearNotifications() {
    try {
      await apiClearAllNotifications()
      notifications.value = []
      unreadCount.value = 0
      showToast({
        type: 'success',
        message: '通知已清空'
      })
    } catch (error) {
      console.error('清空通知失败:', error)
      showToast({
        type: 'fail',
        message: '清空失败'
      })
    }
  }

  // 轮询服务器通知
  function startPolling(interval = 30000) {
    if (pollingTimer) return

    pollingTimer = setInterval(async () => {
      try {
        await loadNotifications()
      } catch (error) {
        console.error('轮询通知失败:', error)
      }
    }, interval)
  }

  // 停止轮询
  function stopPolling() {
    if (pollingTimer) {
      clearInterval(pollingTimer)
      pollingTimer = null
    }
  }

  // 使用 Server-Sent Events (SSE)
  function startSSE() {
    if (eventSource) return

    const token = localStorage.getItem('token')
    const baseURL = import.meta.env.VITE_API_BASE_URL || '/api'

    eventSource = new EventSource(`${baseURL}/notifications/stream?token=${token}`)

    eventSource.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        addNotification(data)
      } catch (error) {
        console.error('解析通知失败:', error)
      }
    }

    eventSource.onerror = (error) => {
      console.error('SSE 错误:', error)
      eventSource?.close()
      eventSource = null
    }
  }

  // 停止 SSE
  function stopSSE() {
    if (eventSource) {
      eventSource.close()
      eventSource = null
    }
  }

  // 快捷方法：添加待审批通知
  function notifyApprovalPending(data) {
    return addNotification({
      type: NOTIFICATION_TYPES.APPROVAL_PENDING,
      title: data.title || '有待审批单据',
      content: data.content || '请及时处理',
      data: {
        entity_type: data.entity_type,
        entity_id: data.entity_id,
      },
    })
  }

  // 快捷方法：添加审批结果通知
  function notifyApprovalResult(data) {
    const type = data.approved
      ? NOTIFICATION_TYPES.APPROVAL_APPROVED
      : NOTIFICATION_TYPES.APPROVAL_REJECTED

    return addNotification({
      type,
      title: data.title || (data.approved ? '审批已通过' : '审批已拒绝'),
      content: data.content || '',
      data: {
        entity_type: data.entity_type,
        entity_id: data.entity_id,
      },
    })
  }

  // 快捷方法：添加状态变更通知
  function notifyStatusChanged(data) {
    return addNotification({
      type: NOTIFICATION_TYPES.STATUS_CHANGED,
      title: data.title || '状态已变更',
      content: data.content || '',
      data: {
        entity_type: data.entity_type,
        entity_id: data.entity_id,
      },
    })
  }

  // 快捷方法：添加库存预警通知
  function notifyStockAlert(data) {
    return addNotification({
      type: NOTIFICATION_TYPES.STOCK_ALERT,
      title: '库存预警',
      content: data.content || '有材料库存不足',
      data: {
        stock_id: data.stock_id,
      },
    })
  }

  onMounted(() => {
    loadNotifications()
  })

  onUnmounted(() => {
    stopPolling()
    stopSSE()
  })

  return {
    notifications,
    unreadCount,
    addNotification,
    markAsRead,
    markAllAsRead,
    removeNotification,
    clearNotifications,
    startPolling,
    stopPolling,
    startSSE,
    stopSSE,
    notifyApprovalPending,
    notifyApprovalResult,
    notifyStatusChanged,
    notifyStockAlert,
    NOTIFICATION_TYPES,
  }
}
