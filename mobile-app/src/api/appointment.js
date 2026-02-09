import request from '@/utils/request'

// 获取预约单列表
export function getAppointments(params) {
  return request({
    url: '/appointments',
    method: 'GET',
    params
  })
}

// 获取我的预约
export function getMyAppointments(params) {
  return request({
    url: '/appointments/my',
    method: 'GET',
    params
  })
}

// 获取待审批列表
export function getPendingApprovals(params) {
  return request({
    url: '/appointments/pending',
    method: 'GET',
    params
  })
}

// 获取作业人员的预约列表
export function getWorkerAppointments(workerId, params) {
  return request({
    url: `/appointments/worker/${workerId}`,
    method: 'GET',
    params
  })
}

// 搜索预约单
export function searchAppointments(params) {
  return request({
    url: '/appointments/search',
    method: 'GET',
    params
  })
}

// 获取预约单详情
export function getAppointmentDetail(id) {
  return request({
    url: `/appointments/${id}`,
    method: 'GET'
  })
}

// 创建预约单
export function createAppointment(data) {
  return request({
    url: '/appointments',
    method: 'POST',
    data
  })
}

// 批量创建预约单
export function batchCreateAppointments(data) {
  return request({
    url: '/appointments/batch',
    method: 'POST',
    data
  })
}

// 更新预约单
export function updateAppointment(id, data) {
  return request({
    url: `/appointments/${id}`,
    method: 'PUT',
    data
  })
}

// 删除预约单
export function deleteAppointment(id) {
  return request({
    url: `/appointments/${id}`,
    method: 'DELETE'
  })
}

// 提交审批
export function submitAppointment(id) {
  return request({
    url: `/appointments/${id}/submit`,
    method: 'POST'
  })
}

// 启动工作流
export function startWorkflow(id, data) {
  return request({
    url: `/appointments/${id}/workflow/start`,
    method: 'POST',
    data
  })
}

// 审批预约单
export function approveAppointment(id, data) {
  return request({
    url: `/appointments/${id}/approve`,
    method: 'POST',
    data
  })
}

// 撤回预约单
export function recallAppointment(id) {
  return request({
    url: `/appointments/${id}/recall`,
    method: 'POST'
  })
}

// 分配作业人员
export function assignWorker(id, data) {
  return request({
    url: `/appointments/${id}/assign`,
    method: 'POST',
    data
  })
}

// 开始作业
export function startWork(id) {
  return request({
    url: `/appointments/${id}/start`,
    method: 'POST'
  })
}

// 完成作业
export function completeAppointment(id, data) {
  return request({
    url: `/appointments/${id}/complete`,
    method: 'POST',
    data
  })
}

// 取消预约
export function cancelAppointment(id, data) {
  return request({
    url: `/appointments/${id}/cancel`,
    method: 'POST',
    data
  })
}

// 获取审批历史
export function getApprovalHistory(id) {
  return request({
    url: `/appointments/${id}/approval-history`,
    method: 'GET'
  })
}

// 获取工作流进度
export function getWorkflowProgress(id) {
  return request({
    url: `/appointments/${id}/workflow-progress`,
    method: 'GET'
  })
}

// 获取当前审批节点
export function getCurrentApproval(id) {
  return request({
    url: `/appointments/${id}/current-approval`,
    method: 'GET'
  })
}

// 批量审批
export function batchApprove(data) {
  return request({
    url: '/appointments/batch-approve',
    method: 'POST',
    data
  })
}

// 获取统计数据
export function getAppointmentStats(params) {
  return request({
    url: '/appointments/stats',
    method: 'GET',
    params
  })
}

// 日历相关API

// 获取作业人员日历
export function getWorkerCalendar(workerId, params) {
  return request({
    url: `/appointments/calendar/worker/${workerId}`,
    method: 'GET',
    params
  })
}

// 检查可用性
export function checkAvailability(data) {
  return request({
    url: '/appointments/calendar/check-availability',
    method: 'POST',
    data
  })
}

// 批量锁定日历
export function batchBlockCalendar(data) {
  return request({
    url: '/appointments/calendar/batch-block',
    method: 'POST',
    data
  })
}

// 获取可用作业人员
export function getAvailableWorkers(params) {
  return request({
    url: '/appointments/calendar/available-workers',
    method: 'GET',
    params
  })
}

// 获取日历视图数据
export function getCalendarView(params) {
  return request({
    url: '/appointments/calendar/view',
    method: 'GET',
    params
  })
}

// 导出预约单
export function exportAppointments(params) {
  return request({
    url: '/appointments/export',
    method: 'GET',
    params,
    responseType: 'blob'
  })
}

// 工具函数

// 获取时间段标签
export function getTimeSlotLabel(timeSlot) {
  const labels = {
    morning: '上午',
    afternoon: '下午',
    evening: '晚上',
    full_day: '全天'
  }
  return labels[timeSlot] || timeSlot
}

// 获取状态标签
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

// 获取状态颜色
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

// 获取时间段选项
export function getTimeSlotOptions() {
  return [
    { value: 'morning', label: '上午 (08:00-12:00)' },
    { value: 'afternoon', label: '下午 (14:00-18:00)' },
    { value: 'evening', label: '晚上 (19:00-22:00)' },
    { value: 'full_day', label: '全天' }
  ]
}

// 获取状态选项（用于筛选）
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

// 判断是否可编辑
export function isEditable(status) {
  return status === 'draft'
}

// 判断是否可取消
export function isCancellable(status) {
  return status === 'pending' || status === 'scheduled' || status === 'draft'
}

// 判断是否可完成
export function canComplete(status) {
  return status === 'in_progress' || status === 'scheduled'
}

// 判断是否可开始
export function canStart(status) {
  return status === 'scheduled'
}

// 格式化日期显示
export function formatDateRange(appointment) {
  const date = new Date(appointment.work_date)
  const dateStr = date.toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })
  const timeSlot = getTimeSlotLabel(appointment.time_slot)
  return `${dateStr} ${timeSlot}`
}

// 计算预约距离现在的天数
export function getDaysUntilWork(workDate) {
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  const work = new Date(workDate)
  work.setHours(0, 0, 0, 0)
  const diffTime = work - today
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
  return diffDays
}

// 获取紧急程度标签
export function getUrgentLabel(isUrgent, priority) {
  if (isUrgent && priority >= 7) {
    return '紧急'
  } else if (isUrgent) {
    return '加急'
  } else if (priority >= 5) {
    return '重要'
  }
  return '普通'
}

// 获取紧急程度颜色
export function getUrgentColor(isUrgent, priority) {
  if (isUrgent && priority >= 7) {
    return 'red'
  } else if (isUrgent) {
    return 'orange'
  } else if (priority >= 5) {
    return 'blue'
  }
  return 'gray'
}
