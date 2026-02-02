import dayjs from 'dayjs'
import 'dayjs/locale/zh-cn'
import relativeTime from 'dayjs/plugin/relativeTime'
import isSameOrBefore from 'dayjs/plugin/isSameOrBefore'
import isSameOrAfter from 'dayjs/plugin/isSameOrAfter'

dayjs.locale('zh-cn')
dayjs.extend(relativeTime)
dayjs.extend(isSameOrBefore)
dayjs.extend(isSameOrAfter)

/**
 * 格式化日期
 * @param {string|Date|number} date - 日期
 * @param {string} format - 格式
 * @returns {string}
 */
export function formatDate(date, format = 'YYYY-MM-DD HH:mm:ss') {
  if (!date) return '-'
  return dayjs(date).format(format)
}

/**
 * 格式化相对时间
 * @param {string|Date|number} date - 日期
 * @returns {string}
 */
export function formatRelativeTime(date) {
  if (!date) return '-'
  return dayjs(date).fromNow()
}

/**
 * 格式化日期为简短格式
 * @param {string|Date|number} date - 日期
 * @returns {string}
 */
export function formatShortDate(date) {
  if (!date) return '-'
  return dayjs(date).format('MM-DD')
}

/**
 * 格式化日期为月份格式
 * @param {string|Date|number} date - 日期
 * @returns {string}
 */
export function formatMonth(date) {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM')
}

/**
 * 判断日期是否在今天之前
 * @param {string|Date|number} date - 日期
 * @returns {boolean}
 */
export function isBeforeToday(date) {
  if (!date) return false
  return dayjs(date).isBefore(dayjs(), 'day')
}

/**
 * 判断日期是否在今天之后
 * @param {string|Date|number} date - 日期
 * @returns {boolean}
 */
export function isAfterToday(date) {
  if (!date) return false
  return dayjs(date).isAfter(dayjs(), 'day')
}

/**
 * 判断日期是否是今天
 * @param {string|Date|number} date - 日期
 * @returns {boolean}
 */
export function isToday(date) {
  if (!date) return false
  return dayjs(date).isSame(dayjs(), 'day')
}

export default dayjs
