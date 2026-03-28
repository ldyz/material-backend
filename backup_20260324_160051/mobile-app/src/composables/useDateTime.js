/**
 * 日期时间格式化工具
 * 提供统一的日期时间格式化函数
 */

/**
 * 格式化日期
 * @param {string|Date} date - 日期字符串或Date对象
 * @param {object} options - 格式化选项
 * @returns {string} 格式化后的日期字符串
 */
export function formatDate(date, options = {}) {
  if (!date) return '-'

  const dateObj = typeof date === 'string' ? new Date(date) : date

  // 检查日期是否有效
  if (isNaN(dateObj.getTime())) return '-'

  const defaultOptions = {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    ...options
  }

  try {
    return dateObj.toLocaleDateString('zh-CN', defaultOptions)
  } catch (error) {
    console.error('日期格式化失败:', error)
    return '-'
  }
}

/**
 * 格式化日期时间
 * @param {string|Date} date - 日期字符串或Date对象
 * @param {object} options - 格式化选项
 * @returns {string} 格式化后的日期时间字符串
 */
export function formatDateTime(date, options = {}) {
  if (!date) return '-'

  const dateObj = typeof date === 'string' ? new Date(date) : date

  if (isNaN(dateObj.getTime())) return '-'

  const defaultOptions = {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    ...options
  }

  try {
    return dateObj.toLocaleString('zh-CN', defaultOptions)
  } catch (error) {
    console.error('日期时间格式化失败:', error)
    return '-'
  }
}

/**
 * 格式化时间
 * @param {string|Date} date - 日期字符串或Date对象
 * @returns {string} 格式化后的时间字符串
 */
export function formatTime(date) {
  if (!date) return '-'

  const dateObj = typeof date === 'string' ? new Date(date) : date

  if (isNaN(dateObj.getTime())) return '-'

  try {
    return dateObj.toLocaleTimeString('zh-CN', {
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch (error) {
    console.error('时间格式化失败:', error)
    return '-'
  }
}

/**
 * 获取时间段标签
 * @param {string} timeSlot - 时间段 ('morning', 'afternoon', 'full_day')
 * @returns {string} 时间段显示文本
 */
export function getTimeSlotLabel(timeSlot) {
  const labels = {
    morning: '上午',
    afternoon: '下午',
    full_day: '全天',
    AM: '上午',
    PM: '下午',
    FULL: '全天'
  }
  return labels[timeSlot] || timeSlot || '-'
}

/**
 * 格式化预约日期（包含时间段）
 * @param {string} date - 日期字符串
 * @param {string} timeSlot - 时间段
 * @returns {string} 格式化后的日期和时间段
 */
export function formatAppointmentDate(date, timeSlot) {
  if (!date) return '-'

  const dateStr = formatDate(date)
  const slot = getTimeSlotLabel(timeSlot)

  return timeSlot ? `${dateStr} ${slot}` : dateStr
}

/**
 * 获取相对时间描述
 * @param {string|Date} date - 日期
 * @returns {string} 相对时间描述（如：今天、昨天、3天前）
 */
export function getRelativeTime(date) {
  if (!date) return '-'

  const dateObj = typeof date === 'string' ? new Date(date) : date

  if (isNaN(dateObj.getTime())) return '-'

  const now = new Date()
  const diffMs = now - dateObj
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))

  if (diffDays === 0) return '今天'
  if (diffDays === 1) return '昨天'
  if (diffDays === -1) return '明天'
  if (diffDays < -1 && diffDays > -7) return `${Math.abs(diffDays)}天后`
  if (diffDays > 1 && diffDays < 7) return `${diffDays}天前`

  return formatDate(dateObj)
}

/**
 * 判断日期是否为今天
 * @param {string|Date} date - 日期
 * @returns {boolean}
 */
export function isToday(date) {
  if (!date) return false

  const dateObj = typeof date === 'string' ? new Date(date) : date
  const today = new Date()

  return (
    dateObj.getFullYear() === today.getFullYear() &&
    dateObj.getMonth() === today.getMonth() &&
    dateObj.getDate() === today.getDate()
  )
}

/**
 * 判断日期是否过期
 * @param {string|Date} date - 日期
 * @returns {boolean}
 */
export function isExpired(date) {
  if (!date) return false

  const dateObj = typeof date === 'string' ? new Date(date) : date
  return dateObj < new Date()
}

/**
 * 计算日期差（天数）
 * @param {string|Date} date1 - 日期1
 * @param {string|Date} date2 - 日期2
 * @returns {number} 天数差
 */
export function daysDiff(date1, date2) {
  const d1 = typeof date1 === 'string' ? new Date(date1) : date1
  const d2 = typeof date2 === 'string' ? new Date(date2) : date2

  const diffMs = d1 - d2
  return Math.floor(diffMs / (1000 * 60 * 60 * 24))
}

/**
 * 格式化日期范围
 * @param {string|Date} startDate - 开始日期
 * @param {string|Date} endDate - 结束日期
 * @returns {string} 格式化后的日期范围
 */
export function formatDateRange(startDate, endDate) {
  const start = startDate ? formatDate(startDate) : '-'
  const end = endDate ? formatDate(endDate) : '-'
  return `${start} ~ ${end}`
}

/**
 * 解析日期字符串为 Date 对象
 * @param {string} dateStr - 日期字符串
 * @returns {Date|null} Date 对象或 null
 */
export function parseDate(dateStr) {
  if (!dateStr) return null

  const date = new Date(dateStr)
  return isNaN(date.getTime()) ? null : date
}

export function useDateTime() {
  return {
    formatDate,
    formatDateTime,
    formatTime,
    getTimeSlotLabel,
    formatAppointmentDate,
    getRelativeTime,
    isToday,
    isExpired,
    daysDiff,
    formatDateRange,
    parseDate
  }
}
