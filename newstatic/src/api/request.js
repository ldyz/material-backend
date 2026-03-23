/**
 * Axios HTTP 请求封装
 *
 * 本文件提供统一的 HTTP 请求封装，包括：
 * - 请求/响应拦截器
 * - 自动添加认证 Token
 * - 统一的错误处理
 * - 自动登出机制
 *
 * @module Request
 * @author Material Management System
 * @date 2025-01-27
 */

import axios from 'axios'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'

/**
 * 创建 Axios 实例
 *
 * 基础配置：
 * - baseURL: '/api' - 所有 API 请求的基础路径
 * - timeout: 30000ms (30秒) - 请求超时时间
 * - headers: 默认 Content-Type 为 application/json
 */
const request = axios.create({
  baseURL: '/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

/**
 * 请求拦截器
 *
 * 在每个请求发送前自动执行，用于：
 * 1. 自动添加认证 Token 到请求头
 * 2. 可以在这里添加其他全局请求处理逻辑
 *
 * @param {Object} config - Axios 请求配置对象
 * @returns {Object} 处理后的请求配置
 */
request.interceptors.request.use(
  config => {
    // 从 Pinia store 获取认证 token
    const authStore = useAuthStore()

    // 如果用户已登录，自动在请求头中添加 Authorization
    // 格式：Bearer <token>
    if (authStore.token) {
      config.headers.Authorization = `Bearer ${authStore.token}`
    }

    // Debug log for material-plan requests to trace material field
    if (config.url?.includes('material-plan') && config.method === 'post') {
      console.log('[API Request] URL:', config.url)
      console.log('[API Request] Method:', config.method?.toUpperCase())
      if (config.data) {
        console.log('[API Request] Payload:', JSON.parse(JSON.stringify(config.data)))
        if (config.data.items) {
          console.log('[API Request] Items material field check:')
          config.data.items.forEach((item, index) => {
            console.log(`  Item ${index}: material="${item.material}", material_name="${item.material_name}"`)
          })
        }
      } else {
        console.log('[API Request] Payload: (no body)')
      }
    }

    return config
  },
  error => {
    // 请求配置错误处理
    console.error('请求错误:', error)
    return Promise.reject(error)
  }
)

/**
 * 响应拦截器
 *
 * 在收到响应后自动执行，用于：
 * 1. 检查统一的响应格式中的 success 字段
 * 2. 处理业务逻辑错误（success: false）
 * 3. 统一处理 HTTP 状态码错误
 * 4. 自动处理 401 未授权错误（自动登出）
 * 5. 统一的错误提示显示
 *
 * 后端统一响应格式：
 * - 成功：{ success: true, data: ..., pagination?: {...}, message?: "..." }
 * - 失败：{ success: false, error: "错误信息", code?: "ERROR_CODE" }
 *
 * @param {Object} response - Axios 响应对象
 * @returns {*} 处理后的响应数据（自动解包 response.data）
 */
request.interceptors.response.use(
  response => {
    const res = response.data

    // 检查后端统一响应格式的 success 字段
    // 后端规范：所有响应都包含 success 字段
    // success: false 表示业务逻辑错误（如参数验证失败、资源不存在等）
    if (res.success === false) {
      // 业务逻辑错误，显示错误提示
      const errorMsg = res.error || '操作失败'
      ElMessage.error(errorMsg)

      // 返回一个被拒绝的 Promise，让调用方可以捕获错误
      return Promise.reject(new Error(errorMsg))
    }

    // 成功响应：直接返回响应数据（解包 response.data）
    // 调用方可以直接使用 res.data、res.pagination 等
    return res
  },
  error => {
    // HTTP 错误响应处理
    const authStore = useAuthStore()

    if (error.response) {
      // 服务器返回了错误状态码
      const { status, data } = error.response

      switch (status) {
        case 401:
          // 未授权 - Token 过期或无效
          // 防止重复处理 401 错误
          if (!window.__isHandling401) {
            window.__isHandling401 = true
            ElMessage.error('登录已过期，请重新登录')
            // 清除本地存储的用户信息（不调用 logout API，避免额外请求）
            authStore.clearAuth()
            // 跳转到登录页
            window.location.href = '/login'
            // 延迟重置标记，防止短时间内重复触发
            setTimeout(() => {
              window.__isHandling401 = false
            }, 1000)
          }
          break

        case 403:
          // 禁止访问 - 用户没有权限执行此操作
          ElMessage.error(data?.error || data?.message || '没有权限执行此操作')
          break

        case 404:
          // 资源不存在
          ElMessage.error(data?.error || data?.message || '请求的资源不存在')
          break

        case 500:
          // 服务器内部错误
          ElMessage.error(data?.error || data?.message || '服务器错误')
          break

        default:
          // 其他 HTTP 错误
          ElMessage.error(data?.error || data?.message || `请求失败 (${status})`)
      }

    } else if (error.request) {
      // 请求已发出但没有收到响应
      // 通常是网络连接问题或服务器无响应
      ElMessage.error('网络错误，请检查您的网络连接')

    } else {
      // 请求配置错误（如取消的请求）
      ElMessage.error(error.message || '请求失败')
    }

    // 将错误继续传递，让调用方可以单独处理
    return Promise.reject(error)
  }
)

/**
 * 导出 Axios 实例
 *
 * 使用示例：
 * ```javascript
 * import request from '@/api/request'
 *
 * // GET 请求
 * const data = await request({ url: '/users', method: 'GET' })
 *
 * // POST 请求
 * const result = await request({
 *   url: '/users',
 *   method: 'POST',
 *   data: { name: 'John' }
 * })
 * ```
 */
export default request
