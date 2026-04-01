import request from '@/utils/request'

/**
 * 用户名密码登录
 */
export function login(data) {
  return request({
    url: '/auth/login',
    method: 'POST',
    data
  })
}

/**
 * 微信登录
 * @param {string} code - 微信登录凭证
 */
export function wechatLogin(code) {
  return request({
    url: '/auth/wechat-login',
    method: 'POST',
    data: { code }
  })
}

/**
 * 绑定微信
 * @param {Object} data - 包含 username, password, code
 */
export function bindWechat(data) {
  return request({
    url: '/auth/bind-wechat',
    method: 'POST',
    data
  })
}

/**
 * 退出登录
 */
export function logout() {
  return request({
    url: '/auth/logout',
    method: 'POST'
  })
}

/**
 * 获取当前用户信息
 */
export function getCurrentUser() {
  return request({
    url: '/auth/me',
    method: 'GET'
  })
}

/**
 * 更新用户头像
 */
export function uploadAvatar(filePath) {
  const token = uni.getStorageSync('token')
  return new Promise((resolve, reject) => {
    uni.uploadFile({
      url: 'https://home.mbed.org.cn/auth/avatar',
      filePath: filePath,
      name: 'file',
      header: {
        'Authorization': `Bearer ${token}`
      },
      success: (response) => {
        if (response.statusCode === 200) {
          try {
            const data = JSON.parse(response.data)
            resolve(data)
          } catch (e) {
            reject(new Error('解析响应失败'))
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
