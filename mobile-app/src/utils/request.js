import axios from 'axios'
import { showToast } from 'vant'
import { storage } from './storage'
import router from '@/router'

const request = axios.create({
  baseURL: '/api',
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json'
  }
})

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
      const errorMsg = res.error || '操作失败'
      showToast({ type: 'fail', message: errorMsg })
      return Promise.reject(new Error(errorMsg))
    }

    return res
  },
  error => {
    if (error.response?.status === 401) {
      showToast({ type: 'fail', message: '登录已过期' })
      storage.clear()
      router.push('/login')
    } else if (error.response) {
      showToast({ type: 'fail', message: error.response.data?.error || '请求失败' })
    } else {
      showToast({ type: 'fail', message: '网络错误' })
    }
    return Promise.reject(error)
  }
)

export default request
