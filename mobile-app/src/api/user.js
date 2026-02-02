import request from '@/utils/request'

/**
 * 获取用户信息
 * @param {number} id - 用户ID
 * @returns {Promise}
 */
export function getUserInfo(id) {
  return request.get(`/users/${id}`)
}

/**
 * 获取当前登录用户的完整信息
 * @returns {Promise}
 */
export function getCurrentUserFullInfo() {
  return request.get('/auth/me')
}
