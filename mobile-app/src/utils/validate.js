/**
 * 验证工具函数
 */

/**
 * 验证手机号
 */
export function isValidPhone(phone) {
  return /^1[3-9]\d{9}$/.test(phone)
}

/**
 * 验证邮箱
 */
export function isValidEmail(email) {
  return /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$/.test(email)
}

/**
 * 验证身份证号
 */
export function isValidIdCard(idCard) {
  return /^[1-9]\d{5}(18|19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]$/.test(idCard)
}

/**
 * 验证URL
 */
export function isValidUrl(url) {
  try {
    new URL(url)
    return true
  } catch {
    return false
  }
}

/**
 * 验证数字（包括小数）
 */
export function isNumber(value) {
  return !isNaN(parseFloat(value)) && isFinite(value)
}

/**
 * 验证整数
 */
export function isInteger(value) {
  return Number.isInteger(Number(value))
}

/**
 * 验证正数
 */
export function isPositive(value) {
  return isNumber(value) && Number(value) > 0
}

/**
 * 验证非负数
 */
export function isNonNegative(value) {
  return isNumber(value) && Number(value) >= 0
}

/**
 * 验证是否在范围内
 */
export function isInRange(value, min, max) {
  const num = Number(value)
  return !isNaN(num) && num >= min && num <= max
}

/**
 * 验证字符串长度
 */
export function isValidLength(str, min, max) {
  const len = String(str).length
  return len >= min && len <= max
}

/**
 * 验证是否为空
 */
export function isEmpty(value) {
  if (value === null || value === undefined) return true
  if (typeof value === 'string') return value.trim() === ''
  if (Array.isArray(value)) return value.length === 0
  if (typeof value === 'object') return Object.keys(value).length === 0
  return false
}

/**
 * 表单验证器生成器
 */
export const Validators = {
  required: (message = '此项为必填项') => (value) => {
    return !isEmpty(value) || message
  },

  phone: (message = '请输入正确的手机号') => (value) => {
    return !value || isValidPhone(value) || message
  },

  email: (message = '请输入正确的邮箱地址') => (value) => {
    return !value || isValidEmail(value) || message
  },

  idCard: (message = '请输入正确的身份证号') => (value) => {
    return !value || isValidIdCard(value) || message
  },

  minLength: (min, message) => (value) => {
    return !value || String(value).length >= min || message
  },

  maxLength: (max, message) => (value) => {
    return !value || String(value).length <= max || message
  },

  range: (min, max, message) => (value) => {
    return !value || isInRange(value, min, max) || message
  },

  positive: (message = '请输入正数') => (value) => {
    return !value || isPositive(value) || message
  },

  pattern: (pattern, message) => (value) => {
    return !value || pattern.test(value) || message
  },
}
