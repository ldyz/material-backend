/**
 * 通知管理 Store
 *
 * 使用 Pinia 管理系统通知，包括：
 * - 通知列表
 * - 未读数量
 * - WebSocket 实时更新
 * - 标记已读、删除操作
 *
 * @module NotificationStore
 * @date 2025-02-07
 */

import { defineStore } from 'pinia'
import { notificationApi } from '@/api'
import { ElMessage } from 'element-plus'
import { useWebSocket } from '@/utils/websocket'
import eventBus from '@/utils/eventBus'

/**
 * 通知状态管理 Store
 *
 * 状态设计说明：
 * - notifications: 通知列表
 * - unreadCount: 未读通知数量
 * - wsConnected: WebSocket 连接状态
 * - loading: 加载状态
 */
export const useNotificationStore = defineStore('notification', {
  state: () => ({
    // 通知列表
    notifications: [],

    // 未读通知数量
    unreadCount: 0,

    // WebSocket 连接状态
    wsConnected: false,

    // 加载状态
    loading: false,

    // 分页信息
    pagination: {
      page: 1,
      pageSize: 20,
      total: 0
    }
  }),

  /**
   * Getters - 计算属性
   */
  getters: {
    /**
     * 获取未读通知列表
     */
    unreadNotifications: (state) => {
      return state.notifications.filter(n => !n.is_read)
    },

    /**
     * 获取已读通知列表
     */
    readNotifications: (state) => {
      return state.notifications.filter(n => n.is_read)
    },

    /**
     * 检查是否有未读通知
     */
    hasUnread: (state) => {
      return state.unreadCount > 0
    }
  },

  /**
   * Actions - 异步方法和业务逻辑
   */
  actions: {
    /**
     * 获取通知列表
     *
     * @param {Object} params - 查询参数
     * @param {number} params.page - 页码
     * @param {number} params.page_size - 每页数量
     * @param {boolean} params.unreadOnly - 是否只获取未读
     */
    async fetchNotifications(params = {}) {
      this.loading = true
      try {
        const response = await notificationApi.getList({
          page: this.pagination.page,
          page_size: this.pagination.pageSize,
          ...params
        })

        if (response.success) {
          this.notifications = response.data || []
          this.pagination.total = response.pagination?.total || 0

          // 从响应中获取未读数量
          if (response.meta?.unread_count !== undefined) {
            this.unreadCount = response.meta.unread_count
          }
        }
      } catch (error) {
        console.error('获取通知列表失败:', error)
        ElMessage.error('获取通知列表失败')
      } finally {
        this.loading = false
      }
    },

    /**
     * 获取未读通知数量
     */
    async fetchUnreadCount() {
      try {
        const response = await notificationApi.getUnreadCount()
        if (response.success) {
          this.unreadCount = response.data?.unread_count || 0
        }
      } catch (error) {
        console.error('获取未读数量失败:', error)
      }
    },

    /**
     * 标记通知为已读
     *
     * @param {number} id - 通知ID
     */
    async markAsRead(id) {
      try {
        const response = await notificationApi.markAsRead(id)
        if (response.success) {
          // 更新本地状态
          const notification = this.notifications.find(n => n.id === id)
          if (notification && !notification.is_read) {
            notification.is_read = true
            notification.read_at = new Date().toISOString()
            this.unreadCount = Math.max(0, this.unreadCount - 1)
          }
        }
      } catch (error) {
        console.error('标记已读失败:', error)
        ElMessage.error('标记已读失败')
      }
    },

    /**
     * 标记所有通知为已读
     */
    async markAllAsRead() {
      try {
        const response = await notificationApi.markAllAsRead()
        if (response.success) {
          // 更新本地状态
          this.notifications.forEach(n => {
            if (!n.is_read) {
              n.is_read = true
              n.read_at = new Date().toISOString()
            }
          })
          this.unreadCount = 0
          ElMessage.success('已全部标记为已读')
        }
      } catch (error) {
        console.error('标记全部已读失败:', error)
        ElMessage.error('标记全部已读失败')
      }
    },

    /**
     * 删除通知
     *
     * @param {number} id - 通知ID
     */
    async deleteNotification(id) {
      try {
        const response = await notificationApi.delete(id)
        if (response.success) {
          // 更新本地状态
          const index = this.notifications.findIndex(n => n.id === id)
          if (index !== -1) {
            const notification = this.notifications[index]
            if (!notification.is_read) {
              this.unreadCount = Math.max(0, this.unreadCount - 1)
            }
            this.notifications.splice(index, 1)
          }
          ElMessage.success('通知已删除')
        }
      } catch (error) {
        console.error('删除通知失败:', error)
        ElMessage.error('删除通知失败')
      }
    },

    /**
     * 清空所有通知
     */
    async clearAll() {
      try {
        const response = await notificationApi.clearAll()
        if (response.success) {
          this.notifications = []
          this.unreadCount = 0
          ElMessage.success('通知已清空')
        }
      } catch (error) {
        console.error('清空通知失败:', error)
        ElMessage.error('清空通知失败')
      }
    },

    /**
     * 添加新通知（来自 WebSocket）
     *
     * @param {Object} notification - 通知对象
     */
    addNotification(notification) {
      // 添加到列表开头
      this.notifications.unshift(notification)

      // 如果是未读通知，增加未读数量
      if (!notification.is_read) {
        this.unreadCount++
      }

      // 显示通知提示
      ElMessage({
        message: notification.content || notification.title,
        type: 'info',
        duration: 3000,
        showClose: true
      })
    },

    /**
     * 更新未读数量（来自 WebSocket）
     *
     * @param {number} count - 未读数量
     */
    updateUnreadCount(count) {
      this.unreadCount = count
    },

    /**
     * 初始化 WebSocket 连接
     */
    initWebSocket() {
      const token = localStorage.getItem('token')
      if (!token) {
        return
      }

      // 构建 WebSocket URL
      const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
      const wsHost = window.location.host
      const wsUrl = `${wsProtocol}//${wsHost}/api/notification/ws`

      const { connect, disconnect, isConnected } = useWebSocket(wsUrl, {
        token,
        debug: false,  // 关闭调试日志
        maxReconnectAttempts: 5,  // 最多重连5次
        onMessage: (data) => {
          this.handleWebSocketMessage(data)
        },
        onOpen: () => {
          this.wsConnected = true
        },
        onClose: () => {
          this.wsConnected = false
        },
        onError: () => {
          this.wsConnected = false
        }
      })

      // 连接 WebSocket
      connect()

      // 保存 disconnect 方法以便后续使用
      this.disconnectWebSocket = disconnect
      this.isWebSocketConnected = isConnected
    },

    /**
     * 处理 WebSocket 消息
     *
     * @param {Object} data - 消息数据
     */
    handleWebSocketMessage(data) {
      switch (data.type) {
        case 'notification':
          // 新通知
          this.addNotification(data.data)
          break
        case 'unread_count':
          // 未读数量更新
          this.updateUnreadCount(data.count)
          break
        case 'appointment_approval_update':
          // 预约审批更新 - 通过 EventBus 广播
          console.log('[NotificationStore] 收到审批更新消息:', data)
          eventBus.emit('appointment:approval-updated', data.data)
          break
        default:
          console.warn('未知的 WebSocket 消息类型:', data.type)
      }
    },

    /**
     * 断开 WebSocket 连接
     */
    closeWebSocket() {
      if (this.disconnectWebSocket) {
        this.disconnectWebSocket()
        this.wsConnected = false
      }
    },

    /**
     * 重置状态
     */
    reset() {
      this.notifications = []
      this.unreadCount = 0
      this.pagination = {
        page: 1,
        pageSize: 20,
        total: 0
      }
      this.closeWebSocket()
    }
  }
})
