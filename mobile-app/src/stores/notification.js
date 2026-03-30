import { defineStore } from 'pinia'
import { showToast } from 'vant'
import {
  getNotifications,
  getUnreadCount,
  markAsRead,
  markAllAsRead,
  deleteNotification,
  clearAllNotifications
} from '@/api/notification'
import { logger } from '@/utils/logger'

export const useNotificationStore = defineStore('notification', {
  state: () => ({
    notifications: [],
    unreadCount: 0,
    loading: false,
    wsConnected: false,
    pagination: {
      page: 1,
      pageSize: 20,
      total: 0
    }
  }),

  getters: {
    hasUnread: (state) => state.unreadCount > 0,
    unreadNotifications: (state) => state.notifications.filter(n => !n.is_read)
  },

  actions: {
    /**
     * 获取通知列表
     * 添加防重复请求逻辑，避免并发调用导致数据覆盖
     */
    async fetchNotifications(params = {}) {
      // 防止重复请求
      if (this.loading) return

      this.loading = true
      try {
        const response = await getNotifications({
          page: this.pagination.page,
          page_size: this.pagination.pageSize,
          ...params
        })

        if (response.success) {
          this.notifications = response.data || []
          this.pagination.total = response.pagination?.total || 0

          if (response.meta?.unread_count !== undefined) {
            this.unreadCount = response.meta.unread_count
          }
        }
      } catch (error) {
        logger.error('获取通知列表失败:', error)
        showToast({ type: 'fail', message: '获取通知列表失败' })
      } finally {
        this.loading = false
      }
    },

    /**
     * 获取未读数量
     */
    async fetchUnreadCount() {
      try {
        const response = await getUnreadCount()
        if (response.success) {
          this.unreadCount = response.data?.unread_count || 0
        }
      } catch (error) {
        logger.error('获取未读数量失败:', error)
      }
    },

    /**
     * 标记为已读
     */
    async markAsRead(id) {
      // 先更新本地状态（乐观更新）
      const notification = this.notifications.find(n => n.id === id)
      const wasUnread = notification && !notification.is_read

      if (notification && !notification.is_read) {
        notification.is_read = true
        notification.read_at = new Date().toISOString()
        this.unreadCount = Math.max(0, this.unreadCount - 1)
      }

      try {
        await markAsRead(id)
        // API 调用成功，本地状态已经更新
      } catch (error) {
        logger.error('标记已读失败:', error)
        // 回滚本地状态
        if (notification && wasUnread) {
          notification.is_read = false
          notification.read_at = null
          this.unreadCount++
        }
        showToast({ type: 'fail', message: '标记已读失败' })
      }
    },

    /**
     * 全部标记为已读
     */
    async markAllAsRead() {
      try {
        const response = await markAllAsRead()
        if (response.success) {
          this.notifications.forEach(n => {
            if (!n.is_read) {
              n.is_read = true
              n.read_at = new Date().toISOString()
            }
          })
          this.unreadCount = 0
          showToast({ type: 'success', message: '已全部标记为已读' })
        }
      } catch (error) {
        logger.error('标记全部已读失败:', error)
        showToast({ type: 'fail', message: '标记全部已读失败' })
      }
    },

    /**
     * 删除通知
     */
    async deleteNotification(id) {
      try {
        const response = await deleteNotification(id)
        if (response.success) {
          const index = this.notifications.findIndex(n => n.id === id)
          if (index !== -1) {
            const notification = this.notifications[index]
            if (!notification.is_read) {
              this.unreadCount = Math.max(0, this.unreadCount - 1)
            }
            this.notifications.splice(index, 1)
          }
          showToast({ type: 'success', message: '通知已删除' })
        }
      } catch (error) {
        logger.error('删除通知失败:', error)
        showToast({ type: 'fail', message: '删除通知失败' })
      }
    },

    /**
     * 清空所有通知
     */
    async clearAll() {
      try {
        const response = await clearAllNotifications()
        if (response.success) {
          this.notifications = []
          this.unreadCount = 0
          showToast({ type: 'success', message: '通知已清空' })
        }
      } catch (error) {
        logger.error('清空通知失败:', error)
        showToast({ type: 'fail', message: '清空通知失败' })
      }
    },

    /**
     * 添加新通知（来自 WebSocket）
     */
    addNotification(notification) {
      this.notifications.unshift(notification)
      if (!notification.is_read) {
        this.unreadCount++
      }
      showToast({ type: 'success', message: notification.content || notification.title })
    },

    /**
     * 更新未读数量
     */
    updateUnreadCount(count) {
      this.unreadCount = count
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
    }
  }
})
