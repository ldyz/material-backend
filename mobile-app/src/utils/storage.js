const TOKEN_KEY = 'token'
const USER_KEY = 'user_info'
const PUSH_TOKEN_KEY = 'push_token'

export const storage = {
  getToken() {
    return localStorage.getItem(TOKEN_KEY)
  },

  setToken(token) {
    localStorage.setItem(TOKEN_KEY, token)
  },

  removeToken() {
    localStorage.removeItem(TOKEN_KEY)
  },

  getUser() {
    const userStr = localStorage.getItem(USER_KEY)
    return userStr ? JSON.parse(userStr) : null
  },

  setUser(user) {
    localStorage.setItem(USER_KEY, JSON.stringify(user))
  },

  removeUser() {
    localStorage.removeItem(USER_KEY)
  },

  getPushToken() {
    return localStorage.getItem(PUSH_TOKEN_KEY)
  },

  setPushToken(token) {
    localStorage.setItem(PUSH_TOKEN_KEY, token)
  },

  removePushToken() {
    localStorage.removeItem(PUSH_TOKEN_KEY)
  },

  clear() {
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(USER_KEY)
    localStorage.removeItem(PUSH_TOKEN_KEY)
  }
}
