import axios from 'axios'
import { showToast } from 'vant'
import { storage } from './storage'

// 检测是否在 Capacitor 原生环境中
const isCapacitor = typeof window !== 'undefined' && window.Capacitor
const baseURL = isCapacitor ? 'https://home.mbed.org.cn:9090/api' : '/api'

let router = null

const request = axios.create({
  baseURL,
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 设置 router 实例（延迟设置，避免循环依赖）
export function setRouter(routerInstance) {
  router = routerInstance
}

request.interceptors.request.use(
  config => {
    const token = storage.getToken()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

request.interceptors.response.use(
  response => {
    const res = response.data

    if (res.success === false) {
      // 业务逻辑错误，不在这里显示 toast，让调用方处理
      const errorMsg = res.error || '操作失败'
      const errorObj = new Error(errorMsg)
      errorObj.error = errorMsg
      errorObj.code = res.code
      errorObj.response = response
      return Promise.reject(errorObj)
    }

    return res
  },
  error => {
    // 网络错误或 HTTP 错误状态码 - 构建错误对象，由调用方显示 toast
    if (error.response?.status === 401) {
      // 401 错误特殊处理 - 拦截器直接处理并跳转
      showToast({ type: 'fail', message: '登录已过期', duration: 3000 })
      storage.clear()
      if (router) {
        router.push('/login').catch(err => {
          console.error('Router navigation error:', err)
        })
      } else {
        // Fallback: reload page to go to login
        window.location.href = '/login'
      }
      return Promise.reject(error)
    } else if (error.response) {
      // HTTP 错误状态码 (400, 500 等) - 构建错误对象，让调用方显示
      let errorMsg = '请求失败'

      // 尝试从响应中提取错误信息
      const data = error.response.data
      if (typeof data === 'string') {
        errorMsg = data
      } else if (data) {
        // 尝试多个可能的错误字段
        errorMsg = data.error || data.message || data.msg || data.detail || '请求失败'
      }

      const errorObj = new Error(errorMsg)
      errorObj.error = errorMsg
      errorObj.response = error.response
      errorObj.code = error.response.status
      // 不在这里显示 toast，让组件层处理
      return Promise.reject(errorObj)
    } else if (error.request) {
      // 请求已发出但没有收到响应 - 网络错误
      const errorMsg = '网络连接失败，请检查网络'
      const errorObj = new Error(errorMsg)
      errorObj.error = errorMsg
      // 网络错误在拦截器显示 toast
      showToast({ type: 'fail', message: errorMsg })
      return Promise.reject(errorObj)
    } else {
      // 其他错误
      const errorMsg = '网络错误'
      const errorObj = new Error(errorMsg)
      errorObj.error = errorMsg
      showToast({ type: 'fail', message: errorMsg })
      return Promise.reject(errorObj)
    }
  }
)

export default request
