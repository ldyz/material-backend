import request from '@/utils/request'

/**
 * 获取用户通知列表
 * @param {Object} params - 查询参数
 * @returns {Promise}
 */
export function getNotifications(params) {
  return request.get('/notification/notifications', { params })
}

/**
 * 获取未读通知数量
 * @returns {Promise}
 */
export function getUnreadCount() {
  return request.get('/notification/notifications/count')
}

/**
 * 标记通知为已读
 * @param {number} id - 通知ID
 * @returns {Promise}
 */
export function markAsRead(id) {
  return request.put(`/notification/notifications/${id}/read`)
}

/**
 * 标记所有通知为已读
 * @returns {Promise}
 */
export function markAllAsRead() {
  return request.put('/notification/notifications/read-all')
}

/**
 * 删除通知
 * @param {number} id - 通知ID
 * @returns {Promise}
 */
export function deleteNotification(id) {
  return request.delete(`/notification/notifications/${id}`)
}

/**
 * 清空所有通知
 * @returns {Promise}
 */
export function clearAllNotifications() {
  return request.delete('/notification/notifications')
}
