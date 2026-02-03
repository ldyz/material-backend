import request from '@/utils/request'

/**
 * 获取仪表板数据
 */
export function getDashboardData() {
  return request({
    url: '/system/reports/dashboard',
    method: 'GET'
  })
}

/**
 * 获取统计数据
 */
export function getStats() {
  return request({
    url: '/system/stats',
    method: 'GET'
  })
}

/**
 * 获取系统配置
 */
export function getSettings() {
  return request({
    url: '/system/settings',
    method: 'GET'
  })
}
