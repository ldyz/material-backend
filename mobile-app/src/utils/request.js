/**
 * Axios HTTP 请求封装 (移动端)
 *
 * 本文件提供移动端统一的 HTTP 请求封装，包括：
 * - 请求/响应拦截器
 * - 自动添加认证 Token
 * - 统一的错误处理
 * - 后端统一响应格式适配
 * - 自动登出机制
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

// 创建 axios 实例
const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    // 添加 token
    const token = storage.getToken()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response) => {
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

    // 成功响应，返回完整的响应对象
    // 调用方可以访问 res.data, res.pagination, res.message, res.meta 等
    return res
  },
  (error) => {
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
    } else if (error.request) {
      // 网络错误
      showToast({
        type: 'fail',
        message: '网络错误，请检查网络连接',
      })
    } else {
      // 请求配置错误
      showToast({
        type: 'fail',
        message: error.message || '请求失败',
      })
    }

    return Promise.reject(error)
  }
)

export default request
