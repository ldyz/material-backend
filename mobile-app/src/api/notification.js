import request from '@/utils/request'

/**
 * 获取通知列表
 * @param {Object} params - 查询参数
 * @param {number} params.page - 页码
 * @param {number} params.page_size - 每页数量
 * @param {boolean} params.unread_only - 是否只获取未读
 */
export function getNotifications(params) {
  return request({
    url: '/notification/notifications',
    method: 'GET',
    params
  })
}

/**
 * 获取未读通知数量
 */
export function getUnreadCount() {
  return request({
    url: '/notification/notifications/count',
    method: 'GET'
  })
}

/**
 * 标记通知为已读
 * @param {number} id - 通知ID
 */
export function markAsRead(id) {
  return request({
    url: `/notification/notifications/${id}/read`,
    method: 'PUT'
  })
}

/**
 * 标记所有通知为已读
 */
export function markAllAsRead() {
  return request({
    url: '/notification/notifications/read-all',
    method: 'PUT'
  })
}

/**
 * 删除通知
 * @param {number} id - 通知ID
 */
export function deleteNotification(id) {
  return request({
    url: `/notification/notifications/${id}`,
    method: 'DELETE'
  })
}

/**
 * 清空所有通知
 */
export function clearAllNotifications() {
  return request({
    url: '/notification/notifications',
    method: 'DELETE'
  })
}

/**
 * 注册推送令牌
 * @param {Object} data - 令牌信息
 * @param {string} data.token - 推送令牌
 * @param {string} data.platform - 平台 (ios/android/web)
 * @param {string} data.device_id - 设备ID (可选)
 */
export function registerPushToken(data) {
  return request({
    url: '/notification/register-token',
    method: 'POST',
    data
  })
}

/**
 * 注销推送令牌
 * @param {Object} data - 令牌信息
 * @param {string} data.token - 推送令牌
 */
export function unregisterPushToken(data) {
  return request({
    url: '/notification/unregister-token',
    method: 'DELETE',
    data
  })
}
