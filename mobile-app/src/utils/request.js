/**
 * Axios HTTP 请求封装 (移动端 - 增强版)
 *
 * 本文件提供移动端统一的 HTTP 请求封装，包括：
 * - 请求/响应拦截器
 * - 自动添加认证 Token
 * - 统一的错误处理
 * - 后端统一响应格式适配
 * - 自动登出机制
 * - 请求重试机制
 * - 响应缓存
 * - 请求取消（防重复请求）
 * - 请求日志
 *
 * 后端统一响应格式：
 * - 成功：{ success: true, data: ..., pagination?: {...}, message?: "...", meta?: {...} }
 * - 失败：{ success: false, error: "错误信息", code?: "ERROR_CODE" }
 *
 * @module Request
 * @author Material Management System (Mobile)
 * @date 2025-01-28
 */

import axios from 'axios'
import { showToast } from 'vant'
import { storage } from './storage'
import router from '@/router'

// ==================== 配置 ====================

const CONFIG = {
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 15000,
  retryTimes: 3, // 重试次数
  retryDelay: 1000, // 重试延迟（毫秒）
  enableCache: true, // 是否启用缓存
  cacheTimeout: 5 * 60 * 1000, // 缓存时间（5分钟）
  enableLog: import.meta.env.NODE_ENV === 'development', // 是否启用日志
}

// ==================== 缓存管理 ====================

const cache = new Map()

/**
 * 生成缓存键
 */
function generateCacheKey(config) {
  const { method, url, params, data } = config
  return `${method}:${url}:${JSON.stringify(params)}:${JSON.stringify(data)}`
}

/**
 * 获取缓存
 */
function getCache(config) {
  if (!CONFIG.enableCache || config.method !== 'get') {
    return null
  }

  const key = generateCacheKey(config)
  const cached = cache.get(key)

  if (cached && Date.now() - cached.timestamp < CONFIG.cacheTimeout) {
    log('Cache hit:', key)
    return cached.data
  }

  // 清除过期缓存
  if (cached) {
    cache.delete(key)
  }

  return null
}

/**
 * 设置缓存
 */
function setCache(config, data) {
  if (!CONFIG.enableCache || config.method !== 'get') {
    return
  }

  const key = generateCacheKey(config)
  cache.set(key, {
    data,
    timestamp: Date.now(),
  })

  log('Cache set:', key)
}

/**
 * 清除缓存
 */
export function clearCache(url) {
  if (url) {
    // 清除指定URL的缓存
    for (const key of cache.keys()) {
      if (key.includes(url)) {
        cache.delete(key)
      }
    }
  } else {
    // 清除所有缓存
    cache.clear()
  }
}

// ==================== 请求取消管理 ====================

const pendingRequests = new Map()

/**
 * 生成请求键
 */
function generateRequestKey(config) {
  const { method, url, params, data } = config
  return `${method}:${url}:${JSON.stringify(params)}:${JSON.stringify(data)}`
}

/**
 * 添加待处理请求
 */
function addPendingRequest(config) {
  const key = generateRequestKey(config)

  // 如果存在相同请求，取消之前的请求
  if (pendingRequests.has(key)) {
    const controller = pendingRequests.get(key)
    controller.abort()
  }

  // 创建新的AbortController
  const controller = new AbortController()
  config.signal = controller.signal
  pendingRequests.set(key, controller)

  log('Pending request added:', key)
}

/**
 * 移除待处理请求
 */
function removePendingRequest(config) {
  const key = generateRequestKey(config)
  pendingRequests.delete(key)
  log('Pending request removed:', key)
}

/**
 * 取消所有待处理请求
 */
export function cancelAllRequests() {
  for (const [key, controller] of pendingRequests.entries()) {
    controller.abort()
    log('Request cancelled:', key)
  }
  pendingRequests.clear()
}

// ==================== 日志工具 ====================

function log(...args) {
  if (CONFIG.enableLog) {
    console.log('[Request]', ...args)
  }
}

function logError(...args) {
  if (CONFIG.enableLog) {
    console.error('[Request Error]', ...args)
  }
}

// ==================== 重试机制 ====================

/**
 * 请求重试
 */
async function retryRequest(originalRequest, retryTimes = CONFIG.retryTimes) {
  return new Promise((resolve, reject) => {
    const retry = (attempt) => {
      axios(originalRequest)
        .then(resolve)
        .catch((error) => {
          if (attempt < retryTimes && shouldRetry(error)) {
            log(`Retrying request (${attempt + 1}/${retryTimes})...`)
            setTimeout(() => retry(attempt + 1), CONFIG.retryDelay * attempt)
          } else {
            reject(error)
          }
        })
    }

    retry(1)
  })
}

/**
 * 判断是否应该重试
 */
function shouldRetry(error) {
  // 网络错误或超时错误才重试
  return (
    !error.response &&
    (error.code === 'ECONNABORTED' || error.code === 'ETIMEDOUT' || error.message === 'Network Error')
  )
}

// ==================== 创建 Axios 实例 ====================

const request = axios.create({
  baseURL: CONFIG.baseURL,
  timeout: CONFIG.timeout,
  headers: {
    'Content-Type': 'application/json',
  },
})

// ==================== 请求拦截器 ====================

request.interceptors.request.use(
  (config) => {
    // 检查缓存
    const cached = getCache(config)
    if (cached) {
      // 返回一个特殊的Promise，直接返回缓存数据
      config.adapter = () => Promise.resolve({
        data: cached,
        status: 200,
        statusText: 'OK (cached)',
        headers: {},
        config,
      })
    }

    // 添加 token
    const token = storage.getToken()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }

    // 添加请求取消
    addPendingRequest(config)

    // 记录请求日志
    log('Request:', config.method.toUpperCase(), config.url, config.params || config.data)

    return config
  },
  (error) => {
    logError('Request error:', error)
    return Promise.reject(error)
  }
)

// ==================== 响应拦截器 ====================

request.interceptors.response.use(
  (response) => {
    // 移除待处理请求
    removePendingRequest(response.config)

    // 记录响应日志
    log('Response:', response.config.url, response.data)

    const res = response.data

    /**
     * 后端统一响应格式处理
     *
     * 成功响应格式：
     * { success: true, data: ..., pagination?: {...}, message?: "...", meta?: {...} }
     *
     * 失败响应格式：
     * { success: false, error: "错误信息", code?: "ERROR_CODE" }
     */

    // 处理业务逻辑错误（success: false）
    if (res.success === false) {
      const errorMsg = res.error || '操作失败'
      showToast({
        type: 'fail',
        message: errorMsg,
      })
      return Promise.reject(new Error(errorMsg))
    }

    // 缓存GET请求的响应
    setCache(response.config, res)

    // 成功响应，返回完整的响应对象
    return res
  },
  async (error) => {
    // 移除待处理请求
    if (error.config) {
      removePendingRequest(error.config)
    }

    // 请求被取消
    if (axios.isCancel(error)) {
      log('Request cancelled:', error.message)
      return Promise.reject(error)
    }

    // HTTP 错误处理
    if (error.response) {
      const { status, data } = error.response

      switch (status) {
        case 401:
          // 未授权或token失效
          showToast({
            type: 'fail',
            message: data?.error || '登录已失效，请重新登录',
          })
          storage.clear()
          router.push('/login')
          break
        case 403:
          // 权限不足
          showToast({
            type: 'fail',
            message: data?.error || '没有权限执行此操作',
          })
          break
        case 404:
          // 资源不存在
          showToast({
            type: 'fail',
            message: data?.error || '请求的资源不存在',
          })
          break
        case 500:
          // 服务器错误
          showToast({
            type: 'fail',
            message: data?.error || '服务器错误',
          })
          break
        default:
          // 其他HTTP错误
          showToast({
            type: 'fail',
            message: data?.error || `请求失败 (${status})`,
          })
      }

      logError('HTTP error:', status, data)
    } else if (error.request) {
      // 网络错误或超时，尝试重试
      if (shouldRetry(error)) {
        log('Network error, retrying...', error.message)
        try {
          const response = await retryRequest(error.config)
          return response
        } catch (retryError) {
          showToast({
            type: 'fail',
            message: '网络错误，请检查网络连接',
          })
          logError('Retry failed:', retryError)
        }
      } else {
        showToast({
          type: 'fail',
          message: '网络错误，请检查网络连接',
        })
      }

      logError('Network error:', error.message)
    } else {
      // 请求配置错误
      showToast({
        type: 'fail',
        message: error.message || '请求失败',
      })
      logError('Request config error:', error.message)
    }

    return Promise.reject(error)
  }
)

// ==================== 导出 ====================

export default request

/**
 * 请求助手函数
 */
export const requestHelper = {
  // GET 请求
  get(url, params, config = {}) {
    return request.get(url, { params, ...config })
  },

  // POST 请求
  post(url, data, config = {}) {
    return request.post(url, data, config)
  },

  // PUT 请求
  put(url, data, config = {}) {
    return request.put(url, data, config)
  },

  // DELETE 请求
  delete(url, params, config = {}) {
    return request.delete(url, { params, ...config })
  },

  // 文件上传
  upload(url, file, onProgress, config = {}) {
    const formData = new FormData()
    formData.append('file', file)

    return request.post(url, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      onUploadProgress: (progressEvent) => {
        if (onProgress && progressEvent.total) {
          const percent = Math.round((progressEvent.loaded * 100) / progressEvent.total)
          onProgress(percent)
        }
      },
      ...config,
    })
  },

  // 批量请求（并行）
  async all(requests) {
    return Promise.all(requests)
  },

  // 批量请求（串行）
  async series(requests) {
    const results = []
    for (const req of requests) {
      results.push(await req())
    }
    return results
  },
}
