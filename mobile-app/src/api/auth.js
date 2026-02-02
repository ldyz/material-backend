import request from '@/utils/request'

/**
 * 登录
 * @param {string} username - 用户名
 * @param {string} password - 密码
 * @returns {Promise}
 */
export function login(username, password) {
  return request.post('/auth/login', {
    username,
    password,
  })
}

/**
 * 登出
 * @returns {Promise}
 */
export function logout() {
  return request.post('/auth/logout')
}

/**
 * 获取当前用户信息
 * @returns {Promise}
 */
export function getCurrentUser() {
  return request.get('/auth/me')
}
