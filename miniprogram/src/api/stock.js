import request from '@/utils/request'

/**
 * 获取库存列表
 */
export function getStockList(params) {
  return request({
    url: '/stock/stocks',
    method: 'GET',
    params
  })
}

/**
 * 获取库存详情
 */
export function getStockDetail(id) {
  return request({
    url: `/stock/stocks/${id}`,
    method: 'GET'
  })
}

/**
 * 获取库存出入库记录
 */
export function getStockLogs(id, params) {
  return request({
    url: `/stock/stocks/${id}/logs`,
    method: 'GET',
    params
  })
}

/**
 * 获取库存预警列表
 */
export function getStockAlerts() {
  return request({
    url: '/stock/stocks/alerts',
    method: 'GET'
  })
}
