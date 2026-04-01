import request from '@/utils/request'

/**
 * 获取通知列表
 */
export function getNotifications(params) {
  return request({
    url: '/notifications',
    method: 'GET',
    params
  })
}

/**
 * 获取未读通知数量
 */
export function getUnreadCount() {
  return request({
    url: '/notifications/unread-count',
    method: 'GET'
  })
}

/**
 * 标记通知为已读
 */
export function markAsRead(id) {
  return request({
    url: `/notifications/${id}/read`,
    method: 'POST'
  })
}

/**
 * 标记所有通知为已读
 */
export function markAllAsRead() {
  return request({
    url: '/notifications/read-all',
    method: 'POST'
  })
}

/**
 * 删除通知
 */
export function deleteNotification(id) {
  return request({
    url: `/notifications/${id}`,
    method: 'DELETE'
  })
}
