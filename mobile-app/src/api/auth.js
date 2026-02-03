import request from '@/utils/request'

/**
 * 认证相关 API (移动端)
 *
 * 提供用户登录、登出、获取用户信息、修改密码等功能
 *
 * @module Auth API
 */

/**
 * 用户登录
 * @param {Object} data - 登录信息
 * @param {string} data.username - 用户名
 * @param {string} data.password - 密码
 * @returns {Promise} 返回包含 token 和用户信息的响应
 *
 * @example
 * const result = await login({ username: 'admin', password: '123456' })
 */
export function login(data) {
  return request.post('/auth/login', data)
}

/**
 * 用户登出
 * @returns {Promise} 返回登出结果
 */
export function logout() {
  return request.post('/auth/logout')
}

/**
 * 获取当前登录用户信息
 * @returns {Promise} 返回用户详细信息，包括角色和权限
 */
export function getCurrentUser() {
  return request.get('/auth/me')
}

/**
 * 修改当前用户密码
 * @param {Object} data - 密码修改信息
 * @param {string} data.oldPassword - 原密码
 * @param {string} data.newPassword - 新密码
 * @returns {Promise} 返回修改结果
 */
export function changePassword(data) {
  return request.post('/auth/change-password', data)
}
