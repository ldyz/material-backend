import request from '@/utils/request'

/**
 * 获取预约单列表
 */
export function getAppointments(params) {
  return request({
    url: '/appointments',
    method: 'GET',
    params
  })
}

/**
 * 获取我的预约
 */
export function getMyAppointments(params) {
  return request({
    url: '/appointments/my',
    method: 'GET',
    params
  })
}

/**
 * 获取我的历史联系人
 */
export function getMyContacts() {
  return request({
    url: '/appointments/my/contacts',
    method: 'GET'
  })
}

/**
 * 获取待审批列表
 */
export function getPendingApprovals(params) {
  return request({
    url: '/appointments/pending',
    method: 'GET',
    params
  })
}

/**
 * 获取预约单详情
 */
export function getAppointmentDetail(id) {
  return request({
    url: `/appointments/${id}`,
    method: 'GET'
  })
}

/**
 * 创建预约单
 */
export function createAppointment(data) {
  return request({
    url: '/appointments',
    method: 'POST',
    data
  })
}

/**
 * 批量创建预约单
 */
export function batchCreateAppointments(data) {
  return request({
    url: '/appointments/batch',
    method: 'POST',
    data
  })
}

/**
 * 更新预约单
 */
export function updateAppointment(id, data) {
  return request({
    url: `/appointments/${id}`,
    method: 'PUT',
    data
  })
}

/**
 * 删除预约单
 */
export function deleteAppointment(id) {
  return request({
    url: `/appointments/${id}`,
    method: 'DELETE'
  })
}

/**
 * 审批预约单
 */
export function approveAppointment(id, data) {
  return request({
    url: `/appointments/${id}/approve`,
    method: 'POST',
    data
  })
}

/**
 * 撤回预约单
 */
export function recallAppointment(id) {
  return request({
    url: `/appointments/${id}/recall`,
    method: 'POST'
  })
}

/**
 * 分配作业人员
 */
export function assignWorker(id, data) {
  return request({
    url: `/appointments/${id}/assign`,
    method: 'POST',
    data
  })
}

/**
 * 完成作业
 */
export function completeAppointment(id, data) {
  return request({
    url: `/appointments/${id}/complete`,
    method: 'POST',
    data
  })
}

/**
 * 取消预约
 */
export function cancelAppointment(id, data) {
  return request({
    url: `/appointments/${id}/cancel`,
    method: 'POST',
    data
  })
}

/**
 * 获取审批历史
 */
export function getApprovalHistory(id) {
  return request({
    url: `/appointments/${id}/approval-history`,
    method: 'GET'
  })
}

/**
 * 获取日历视图数据
 */
export function getCalendarView(params) {
  return request({
    url: '/appointments/calendar/view',
    method: 'GET',
    params
  })
}

/**
 * 获取可用作业人员
 */
export function getAvailableWorkers(params) {
  return request({
    url: '/appointments/calendar/available-workers',
    method: 'GET',
    params
  })
}

/**
 * 获取作业人员列表
 */
export function getWorkersList() {
  return request({
    url: '/appointments/workers',
    method: 'GET'
  })
}

/**
 * 获取待审批数量
 */
export function getPendingApprovalCount() {
  return request({
    url: '/appointments/ending-count',
    method: 'GET'
  })
}

/**
 * 获取统计数据
 */
export function getAppointmentStats(params) {
  return request({
    url: '/appointments/stats',
    method: 'GET',
    params
  })
}

// ========== 工具函数 ==========

/**
 * 获取时间段标签
 */
export function getTimeSlotLabel(timeSlot) {
  const labels = {
    morning: '上午',
    noon: '中午',
    afternoon: '下午',
    full_day: '全天'
  }
  return labels[timeSlot] || timeSlot
}

/**
 * 获取状态标签
 */
export function getStatusLabel(status) {
  const labels = {
    draft: '草稿',
    pending: '待审批',
    scheduled: '已排期',
    in_progress: '进行中',
    completed: '已完成',
    cancelled: '已取消',
    rejected: '已拒绝'
  }
  return labels[status] || status
}

/**
 * 获取状态颜色
 */
export function getStatusColor(status) {
  const colors = {
    draft: 'gray',
    pending: 'orange',
    scheduled: 'blue',
    in_progress: 'cyan',
    completed: 'green',
    cancelled: 'gray',
    rejected: 'red'
  }
  return colors[status] || 'gray'
}

/**
 * 获取时间段选项
 */
export function getTimeSlotOptions() {
  return [
    { value: 'morning', label: '上午 (8:00-11:30)' },
    { value: 'noon', label: '中午 (12:00-13:30)' },
    { value: 'afternoon', label: '下午 (13:30-16:30)' },
    { value: 'full_day', label: '全天' }
  ]
}

/**
 * 获取状态选项
 */
export function getStatusOptions() {
  return [
    { value: '', label: '全部' },
    { value: 'draft', label: '草稿' },
    { value: 'pending', label: '待审批' },
    { value: 'scheduled', label: '已排期' },
    { value: 'in_progress', label: '进行中' },
    { value: 'completed', label: '已完成' },
    { value: 'cancelled', label: '已取消' },
    { value: 'rejected', label: '已拒绝' }
  ]
}
