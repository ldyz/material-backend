import request from '@/utils/request'

/**
 * 获取今日待打卡任务
 */
export function getTodayAppointments() {
  return request({
    url: '/attendance/today-appointments',
    method: 'GET'
  })
}

/**
 * 打卡
 */
export function clockIn(data) {
  return request({
    url: '/attendance/clock-in',
    method: 'POST',
    data
  })
}

/**
 * 获取我的打卡记录
 */
export function getMyRecords(params) {
  return request({
    url: '/attendance/my-records',
    method: 'GET',
    params
  })
}

/**
 * 获取我的月度汇总
 */
export function getMySummary(params) {
  return request({
    url: '/attendance/my-summary',
    method: 'GET',
    params
  })
}

/**
 * 获取打卡日历统计
 */
export function getCalendarStatistics(params) {
  return request({
    url: '/attendance/calendar-statistics',
    method: 'GET',
    params
  })
}

// ========== 工具函数 ==========

/**
 * 获取打卡类型标签
 */
export function getAttendanceTypeLabel(type) {
  const labels = {
    morning: '上午打卡',
    afternoon: '下午打卡',
    noon_overtime: '中午加班',
    night_overtime: '晚上加班'
  }
  return labels[type] || type
}

/**
 * 获取打卡类型颜色
 */
export function getAttendanceTypeColor(type) {
  const colors = {
    morning: 'blue',
    afternoon: 'green',
    noon_overtime: 'orange',
    night_overtime: 'purple'
  }
  return colors[type] || 'gray'
}

/**
 * 获取状态标签
 */
export function getStatusLabel(status) {
  const labels = {
    pending: '待确认',
    confirmed: '已确认',
    rejected: '已驳回'
  }
  return labels[status] || status
}

/**
 * 获取状态颜色
 */
export function getStatusColor(status) {
  const colors = {
    pending: 'orange',
    confirmed: 'green',
    rejected: 'red'
  }
  return colors[status] || 'gray'
}

/**
 * 获取打卡类型选项
 */
export function getAttendanceTypeOptions() {
  return [
    { value: 'morning', label: '上午打卡', color: '#1989fa' },
    { value: 'afternoon', label: '下午打卡', color: '#07c160' },
    { value: 'noon_overtime', label: '中午加班', color: '#ff976a' },
    { value: 'night_overtime', label: '晚上加班', color: '#7232dd' }
  ]
}

/**
 * 判断是否是加班类型
 */
export function isOvertimeType(type) {
  return type === 'noon_overtime' || type === 'night_overtime'
}
