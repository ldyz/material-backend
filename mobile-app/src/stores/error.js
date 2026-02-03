import { defineStore } from 'pinia'
import { ref } from 'vue'

/**
 * 错误处理Store
 * 统一管理应用中的错误状态
 */
export const useErrorStore = defineStore('error', () => {
  // 当前错误
  const currentError = ref(null)
  // 错误历史
  const errorHistory = ref([])
  // 最大历史记录数
  const maxHistorySize = 50

  /**
   * 设置错误
   */
  function setError(error) {
    const errorObj = {
      id: Date.now(),
      message: error?.message || '未知错误',
      code: error?.code || null,
      stack: error?.stack || null,
      timestamp: new Date().toISOString(),
      context: error?.context || null,
    }

    currentError.value = errorObj

    // 添加到历史记录
    errorHistory.value.unshift(errorObj)

    // 限制历史记录大小
    if (errorHistory.value.length > maxHistorySize) {
      errorHistory.value = errorHistory.value.slice(0, maxHistorySize)
    }

    // 记录到控制台
    console.error('[Error Store]', errorObj)

    return errorObj
  }

  /**
   * 清除当前错误
   */
  function clearError() {
    currentError.value = null
  }

  /**
   * 清除所有错误历史
   */
  function clearHistory() {
    errorHistory.value = []
  }

  /**
   * 获取最近的错误
   */
  function getRecentErrors(count = 10) {
    return errorHistory.value.slice(0, count)
  }

  /**
   * 按类型筛选错误
   */
  function getErrorsByCode(code) {
    return errorHistory.value.filter(error => error.code === code)
  }

  /**
   * 处理API错误
   */
  function handleApiError(error, context = {}) {
    let errorCode = 'UNKNOWN_ERROR'
    let errorMessage = '操作失败，请稍后重试'

    // HTTP错误
    if (error.response) {
      const { status, data } = error.response
      errorCode = `HTTP_${status}`

      switch (status) {
        case 400:
          errorMessage = data?.error || '请求参数错误'
          break
        case 401:
          errorMessage = '登录已失效，请重新登录'
          errorCode = 'UNAUTHORIZED'
          break
        case 403:
          errorMessage = '没有权限执行此操作'
          errorCode = 'FORBIDDEN'
          break
        case 404:
          errorMessage = '请求的资源不存在'
          errorCode = 'NOT_FOUND'
          break
        case 500:
          errorMessage = '服务器错误'
          errorCode = 'SERVER_ERROR'
          break
        default:
          errorMessage = data?.error || `请求失败 (${status})`
      }
    }
    // 网络错误
    else if (error.request) {
      errorMessage = '网络错误，请检查网络连接'
      errorCode = 'NETWORK_ERROR'
    }
    // 其他错误
    else {
      errorMessage = error.message || errorMessage
    }

    return setError({
      message: errorMessage,
      code: errorCode,
      original: error,
      context,
    })
  }

  /**
   * 处理表单验证错误
   */
  function handleValidationError(errors) {
    const message = Object.values(errors).join(', ')
    return setError({
      message,
      code: 'VALIDATION_ERROR',
      details: errors,
    })
  }

  return {
    currentError,
    errorHistory,
    setError,
    clearError,
    clearHistory,
    getRecentErrors,
    getErrorsByCode,
    handleApiError,
    handleValidationError,
  }
})
