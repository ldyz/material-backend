/**
 * 本地存储封装
 * 适配 uni-app 的存储 API
 */

export const storage = {
  getToken() {
    return uni.getStorageSync('token') || null
  },

  setToken(token) {
    uni.setStorageSync('token', token)
  },

  getUser() {
    const user = uni.getStorageSync('user')
    return user ? JSON.parse(user) : null
  },

  setUser(user) {
    uni.setStorageSync('user', JSON.stringify(user))
  },

  get(key) {
    return uni.getStorageSync(key)
  },

  set(key, value) {
    uni.setStorageSync(key, value)
  },

  remove(key) {
    uni.removeStorageSync(key)
  },

  clear() {
    uni.removeStorageSync('token')
    uni.removeStorageSync('user')
  }
}

export default storage
