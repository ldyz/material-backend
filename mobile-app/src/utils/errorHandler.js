/**
 * 统一错误处理器
 * 提供全局错误处理和统一错误提示
 */

import { showToast } from 'vant'
import { logger } from './logger'

// 错误码映射
const ERROR_MESSAGES = {
  400: '请求参数错误',
  401: '登录已过期',
  403: '没有权限访问',
  404: '资源不存在',
  405: '请求方法不允许',
  408: '请求超时',
  429: '请求过于频繁',
  500: '服务器内部错误',
  502: '网关错误',
  503: '服务暂时不可用',
  504: '网关超时',
  NETWORK_ERROR: '网络连接失败',
  TIMEOUT: '请求超时',
  UNKNOWN: '未知错误'
}

/**
 * 处理错误并显示提示
 * @param {Error} error - 错误对象
 * @param {boolean} showToastMsg - 是否显示错误提示
 * @returns {object} 标准化的错误信息
 */
export function handleError(error, showToastMsg = true) {
  // 获取错误码
  const code = error.code || error.response?.status || 'UNKNOWN'

  // 获取错误消息
  let message = ERROR_MESSAGES[code]

  // 如果没有预定义消息，尝试从错误对象获取
  if (!message) {
    message = error.message || error.error || ERROR_MESSAGES.UNKNOWN
  }

  // 特殊处理：401 错误已在 request.js 中处理，不重复显示
  if (showToastMsg && code !== 401) {
    showToast({
      type: 'fail',
      message,
      duration: 3000
    })
  }

  // 记录错误日志
  logger.error('[ErrorHandler]', {
    code,
    message,
    originalError: error
  })

  return { code, message }
}

/**
 * 包装异步函数，自动处理错误
 * @param {Function} fn - 异步函数
 * @param {boolean} showToastMsg - 是否显示错误提示
 * @returns {Function} 包装后的函数
 */
export function wrapAsync(fn, showToastMsg = true) {
  return async (...args) => {
    try {
      return await fn(...args)
    } catch (error) {
      handleError(error, showToastMsg)
      throw error
    }
  }
}

/**
 * 创建带错误处理的异步操作
 * @param {Function} fn - 异步函数
 * @param {object} options - 选项
 * @param {boolean} options.showToast - 是否显示错误提示
 * @param {Function} options.onError - 错误回调
 * @param {any} options.defaultValue - 错误时的默认返回值
 * @returns {Function} 包装后的函数
 */
export function createSafeAsync(fn, options = {}) {
  const {
    showToast: showToastMsg = true,
    onError = null,
    defaultValue = null
  } = options

  return async (...args) => {
    try {
      return await fn(...args)
    } catch (error) {
      if (showToastMsg) {
        handleError(error, true)
      } else {
        logger.error('[SafeAsync]', error)
      }

      if (onError && typeof onError === 'function') {
        onError(error)
      }

      return defaultValue
    }
  }
}

/**
 * 判断是否为网络错误
 * @param {Error} error - 错误对象
 * @returns {boolean}
 */
export function isNetworkError(error) {
  return (
    !error.response &&
    (error.message === 'Network Error' || error.code === 'NETWORK_ERROR')
  )
}

/**
 * 判断是否为超时错误
 * @param {Error} error - 错误对象
 * @returns {boolean}
 */
export function isTimeoutError(error) {
  return (
    error.code === 'ECONNABORTED' ||
    error.message?.includes('timeout') ||
    error.code === 'TIMEOUT'
  )
}

/**
 * 判断是否为取消请求错误
 * @param {Error} error - 错误对象
 * @returns {boolean}
 */
export function isCancelError(error) {
  return error.__CANCEL__ || error.message === 'canceled'
}

/**
 * 格式化错误信息用于显示
 * @param {Error} error - 错误对象
 * @returns {string} 格式化后的错误信息
 */
export function formatErrorMessage(error) {
  if (typeof error === 'string') return error

  const code = error.code || error.response?.status
  const message = error.message || error.error

  if (code && ERROR_MESSAGES[code]) {
    return ERROR_MESSAGES[code]
  }

  return message || ERROR_MESSAGES.UNKNOWN
}

export default {
  handleError,
  wrapAsync,
  createSafeAsync,
  isNetworkError,
  isTimeoutError,
  isCancelError,
  formatErrorMessage,
  ERROR_MESSAGES
}
