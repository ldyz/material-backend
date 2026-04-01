/**
 * uni-app 请求封装
 * 替代 axios，使用 uni.request
 */

// 基础 URL
const baseURL = 'https://home.mbed.org.cn'

/**
 * 请求拦截器
 */
function requestInterceptor(config) {
  const token = uni.getStorageSync('token')
  if (token) {
    config.header = config.header || {}
    config.header['Authorization'] = `Bearer ${token}`
  }
  return config
}

/**
 * 响应拦截器
 */
function responseInterceptor(response) {
  const { statusCode, data } = response

  // HTTP 状态码判断
  if (statusCode === 200) {
    // 业务状态码判断
    if (data.code === 0 || data.code === 200 || !data.code) {
      return data
    }
    // 业务错误
    return Promise.reject(data)
  }

  // HTTP 错误
  if (statusCode === 401) {
    // 未授权，清除登录状态，跳转登录页
    uni.removeStorageSync('token')
    uni.removeStorageSync('user')
    uni.reLaunch({ url: '/pages/login/index' })
    return Promise.reject(new Error('登录已过期，请重新登录'))
  }

  if (statusCode === 403) {
    return Promise.reject(new Error('没有权限访问'))
  }

  if (statusCode === 404) {
    return Promise.reject(new Error('请求的资源不存在'))
  }

  if (statusCode >= 500) {
    return Promise.reject(new Error('服务器错误'))
  }

  return Promise.reject(new Error(`请求失败: ${statusCode}`))
}

/**
 * 通用请求方法
 */
export function request(options) {
  return new Promise((resolve, reject) => {
    // 构建请求配置
    const config = {
      url: (options.baseURL || baseURL) + options.url,
      method: options.method || 'GET',
      data: options.data || options.params,
      header: {
        'Content-Type': 'application/json',
        ...options.header,
        ...options.headers
      },
      timeout: options.timeout || 30000
    }

    // 请求拦截
    const processedConfig = requestInterceptor(config)

    // 发起请求
    uni.request({
      ...processedConfig,
      success: (response) => {
        responseInterceptor(response)
          .then(resolve)
          .catch(reject)
      },
      fail: (error) => {
        reject(new Error(error.errMsg || '网络请求失败'))
      }
    })
  })
}

/**
 * GET 请求
 */
export function get(url, params, options = {}) {
  return request({
    url,
    method: 'GET',
    params,
    ...options
  })
}

/**
 * POST 请求
 */
export function post(url, data, options = {}) {
  return request({
    url,
    method: 'POST',
    data,
    ...options
  })
}

/**
 * PUT 请求
 */
export function put(url, data, options = {}) {
  return request({
    url,
    method: 'PUT',
    data,
    ...options
  })
}

/**
 * DELETE 请求
 */
export function del(url, data, options = {}) {
  return request({
    url,
    method: 'DELETE',
    data,
    ...options
  })
}

/**
 * 文件上传
 */
export function uploadFile(filePath, options = {}) {
  return new Promise((resolve, reject) => {
    const token = uni.getStorageSync('token')

    uni.uploadFile({
      url: baseURL + (options.url || '/upload/image'),
      filePath: filePath,
      name: options.name || 'file',
      formData: options.formData,
      header: {
        'Authorization': `Bearer ${token}`
      },
      success: (response) => {
        if (response.statusCode === 200) {
          try {
            const data = JSON.parse(response.data)
            resolve(data)
          } catch (e) {
            resolve(response.data)
          }
        } else {
          reject(new Error('上传失败'))
        }
      },
      fail: (error) => {
        reject(new Error(error.errMsg || '上传失败'))
      }
    })
  })
}

export default request
