import request from '@/utils/request'

export function login(data) {
  return request({
    url: '/auth/login',
    method: 'POST',
    data
  })
}

export function logout() {
  return request({
    url: '/auth/logout',
    method: 'POST'
  })
}

export function getCurrentUser() {
  return request({
    url: '/auth/me',
    method: 'GET'
  })
}
