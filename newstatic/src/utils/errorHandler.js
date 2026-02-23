/**
 * 统一的 API 错误处理工具
 * 提供友好的错误提示和日志记录
 */

import { ElMessage } from 'element-plus'

/**
 * 处理 API 错误
 * @param {Error} error - 错误对象
 * @param {string} context - 错误上下文（操作描述）
 * @param {Object} options - 配置选项
 * @param {boolean} options.showMessage - 是否显示错误消息（默认 true）
 * @param {boolean} options.logError - 是否记录错误日志（默认 true）
 * @returns {string} 错误消息
 */
export function handleApiError(error, context = '', options = {}) {
  const {
    showMessage = true,
    logError = true
  } = options

  if (logError) {
    console.error(`[API Error] ${context ? `${context}: ` : ''}`, error)
  }

  let message = '操作失败，请重试'

  if (error.response) {
    // 服务器返回了错误响应
    const status = error.response.status
    const data = error.response.data

    switch (status) {
      case 400:
        message = data?.message || '请求参数错误'
        break
      case 401:
        message = '未授权，请重新登录'
        // 可以在这里触发重新登录逻辑
        break
      case 403:
        message = '没有权限执行此操作'
        break
      case 404:
        message = '请求的资源不存在'
        break
      case 409:
        message = data?.message || '数据冲突，请刷新后重试'
        break
      case 422:
        message = data?.message || '数据验证失败'
        break
      case 429:
        message = '请求过于频繁，请稍后再试'
        break
      case 500:
        message = '服务器错误，请稍后重试'
        break
      case 502:
        message = '网关错误，请稍后重试'
        break
      case 503:
        message = '服务暂时不可用，请稍后重试'
        break
      default:
        message = data?.message || `操作失败 (${status})`
    }
  } else if (error.request) {
    // 请求已发出但没有收到响应
    message = '网络连接失败，请检查网络设置'
  } else if (error.message) {
    // 其他错误
    message = error.message
  }

  if (showMessage) {
    ElMessage.error(message)
  }

  return message
}

/**
 * 包装异步函数，自动处理错误
 * @param {Function} fn - 异步函数
 * @param {string} context - 错误上下文
 * @param {Object} options - 配置选项
 * @returns {Function} 包装后的函数
 */
export function withErrorHandling(fn, context = '', options = {}) {
  return async (...args) => {
    try {
      return await fn(...args)
    } catch (error) {
      handleApiError(error, context, options)
      throw error // 重新抛出错误以便调用者处理
    }
  }
}

/**
 * 创建带有错误处理的 API 调用包装器
 * @param {Object} apiObject - API 对象
 * @param {Object} options - 配置选项
 * @returns {Object} 包装后的 API 对象
 */
export function createApiWrapper(apiObject, options = {}) {
  const wrapped = {}

  for (const [key, value] of Object.entries(apiObject)) {
    if (typeof value === 'function') {
      wrapped[key] = withErrorHandling(value, key, options)
    } else if (typeof value === 'object' && value !== null) {
      wrapped[key] = createApiWrapper(value, options)
    } else {
      wrapped[key] = value
    }
  }

  return wrapped
}

/**
 * 验证错误类型
 * @param {Error} error - 错误对象
 * @returns {string} 错误类型
 */
export function getErrorType(error) {
  if (!error) return 'unknown'

  if (error.response) {
    const status = error.response.status
    if (status >= 400 && status < 500) {
      return 'client_error'
    } else if (status >= 500) {
      return 'server_error'
    }
  } else if (error.request) {
    return 'network_error'
  } else if (error.code === 'ECONNABORTED') {
    return 'timeout'
  }

  return 'unknown'
}

/**
 * 是否为网络错误
 * @param {Error} error - 错误对象
 * @returns {boolean}
 */
export function isNetworkError(error) {
  return getErrorType(error) === 'network_error'
}

/**
 * 是否为超时错误
 * @param {Error} error - 错误对象
 * @returns {boolean}
 */
export function isTimeoutError(error) {
  return getErrorType(error) === 'timeout'
}

/**
 * 是否为认证错误
 * @param {Error} error - 错误对象
 * @returns {boolean}
 */
export function isAuthError(error) {
  return error.response?.status === 401
}

/**
 * 是否为权限错误
 * @param {Error} error - 错误对象
 * @returns {boolean}
 */
export function isPermissionError(error) {
  return error.response?.status === 403
}

/**
 * 是否为验证错误
 * @param {Error} error - 错误对象
 * @returns {boolean}
 */
export function isValidationError(error) {
  return error.response?.status === 422
}

/**
 * 重试装饰器
 * @param {Function} fn - 异步函数
 * @param {number} maxRetries - 最大重试次数
 * @param {number} delay - 重试延迟（毫秒）
 * @returns {Function} 包装后的函数
 */
export function withRetry(fn, maxRetries = 3, delay = 1000) {
  return async (...args) => {
    let lastError

    for (let i = 0; i < maxRetries; i++) {
      try {
        return await fn(...args)
      } catch (error) {
        lastError = error

        // 如果是客户端错误（如 404），不重试
        if (error.response?.status >= 400 && error.response?.status < 500) {
          throw error
        }

        // 如果是最后一次重试，不等待
        if (i < maxRetries - 1) {
          await new Promise(resolve => setTimeout(resolve, delay * (i + 1)))
        }
      }
    }

    throw lastError
  }
}

export default {
  handleApiError,
  withErrorHandling,
  createApiWrapper,
  getErrorType,
  isNetworkError,
  isTimeoutError,
  isAuthError,
  isPermissionError,
  isValidationError,
  withRetry
}
