/**
 * 日期格式化工具
 * 统一使用 YYYY-MM-DD 格式
 */

/**
 * 格式化日期为 YYYY-MM-DD
 * @param {Date|string|number} date - 日期对象、日期字符串或时间戳
 * @returns {string} 格式化后的日期字符串
 */
export function formatDate(date) {
  if (!date) return ''

  const d = new Date(date)
  if (isNaN(d.getTime())) return ''

  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')

  return `${year}-${month}-${day}`
}

/**
 * 格式化日期时间为 YYYY-MM-DD HH:mm:ss
 * @param {Date|string|number} date - 日期对象、日期字符串或时间戳
 * @returns {string} 格式化后的日期时间字符串
 */
export function formatDateTime(date) {
  if (!date) return ''

  const d = new Date(date)
  if (isNaN(d.getTime())) return ''

  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const hours = String(d.getHours()).padStart(2, '0')
  const minutes = String(d.getMinutes()).padStart(2, '0')
  const seconds = String(d.getSeconds()).padStart(2, '0')

  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

/**
 * 格式化时间 HH:mm
 * @param {Date|string|number} date - 日期对象、日期字符串或时间戳
 * @returns {string} 格式化后的时间字符串
 */
export function formatTime(date) {
  if (!date) return ''

  const d = new Date(date)
  if (isNaN(d.getTime())) return ''

  const hours = String(d.getHours()).padStart(2, '0')
  const minutes = String(d.getMinutes()).padStart(2, '0')

  return `${hours}:${minutes}`
}

/**
 * 给日期添加天数
 * @param {Date|string|number} date - 日期对象、日期字符串或时间戳
 * @param {number} days - 要添加的天数（可以为负数）
 * @returns {Date} 新的日期对象
 */
export function addDays(date, days) {
  const d = new Date(date)
  if (isNaN(d.getTime())) return new Date()
  d.setDate(d.getDate() + days)
  return d
}

/**
 * 获取ISO周数
 * @param {Date|string|number} date - 日期对象、日期字符串或时间戳
 * @returns {number} 周数（1-53）
 */
export function getWeekNumber(date) {
  const d = new Date(date)
  if (isNaN(d.getTime())) return 0

  // 复制日期避免修改原对象
  const target = new Date(d.valueOf())
  // 设置到周四（ISO 8601标准）
  target.setDate(target.getDate() - (target.getDay() + 6) % 7 + 3)
  // 获取该年的第一个周四
  const firstThursday = new Date(target.getFullYear(), 0, 4)
  firstThursday.setDate(firstThursday.getDate() - (firstThursday.getDay() + 6) % 7 + 3)
  // 计算周数
  const weekNumber = 1 + Math.round((target.getTime() - firstThursday.getTime()) / 86400000 / 7)
  return weekNumber
}

/**
 * 获取季度
 * @param {Date|string|number} date - 日期对象、日期字符串或时间戳
 * @returns {number} 季度（1-4）
 */
export function getQuarter(date) {
  const d = new Date(date)
  if (isNaN(d.getTime())) return 1
  return Math.floor(d.getMonth() / 3) + 1
}

/**
 * 计算两个日期之间的天数差
 * @param {Date|string|number} date1 - 第一个日期
 * @param {Date|string|number} date2 - 第二个日期
 * @returns {number} 天数差（date2 - date1）
 */
export function diffDays(date1, date2) {
  const d1 = new Date(date1)
  const d2 = new Date(date2)
  if (isNaN(d1.getTime()) || isNaN(d2.getTime())) return 0
  const diff = d2.getTime() - d1.getTime()
  return Math.round(diff / (1000 * 60 * 60 * 24))
}

/**
 * 获取季度的第一天
 * @param {number} quarter - 季度（1-4）
 * @param {number} year - 年份
 * @returns {Date} 季度第一天的日期对象
 */
export function getQuarterStart(quarter, year) {
  const month = (quarter - 1) * 3
  return new Date(year, month, 1)
}

/**
 * 获取季度的最后一天
 * @param {number} quarter - 季度（1-4）
 * @param {number} year - 年份
 * @returns {Date} 季度最后一天的日期对象
 */
export function getQuarterEnd(quarter, year) {
  const month = quarter * 3 - 1
  return new Date(year, month + 1, 0)
}

/**
 * 获取周的第一天（周一）
 * @param {Date|string|number} date - 日期对象、日期字符串或时间戳
 * @returns {Date} 周一的日期对象
 */
export function getWeekStart(date) {
  const d = new Date(date)
  if (isNaN(d.getTime())) return new Date()
  const day = d.getDay()
  const diff = d.getDate() - day + (day === 0 ? -6 : 1)
  return new Date(d.setDate(diff))
}

/**
 * 获取周的最后一天（周日）
 * @param {Date|string|number} date - 日期对象、日期字符串或时间戳
 * @returns {Date} 周日的日期对象
 */
export function getWeekEnd(date) {
  const d = new Date(date)
  if (isNaN(d.getTime())) return new Date()
  const day = d.getDay()
  const diff = d.getDate() - day + (day === 0 ? 0 : 7)
  return new Date(d.setDate(diff))
}

/**
 * 获取月份的第一天
 * @param {Date|string|number} date - 日期对象、日期字符串或时间戳
 * @returns {Date} 月第一天的日期对象
 */
export function getMonthStart(date) {
  const d = new Date(date)
  if (isNaN(d.getTime())) return new Date()
  return new Date(d.getFullYear(), d.getMonth(), 1)
}

/**
 * 获取月份的最后一天
 * @param {Date|string|number} date - 日期对象、日期字符串或时间戳
 * @returns {Date} 月最后一天的日期对象
 */
export function getMonthEnd(date) {
  const d = new Date(date)
  if (isNaN(d.getTime())) return new Date()
  return new Date(d.getFullYear(), d.getMonth() + 1, 0)
}
