import request from '@/utils/request'

/**
 * 获取库存列表
 * @param {Object} params - 查询参数
 * @returns {Promise}
 */
export function getStocks(params) {
  return request.get('/stock/stocks', { params })
}

/**
 * 获取库存详情
 * @param {number} id - 库存ID
 * @returns {Promise}
 */
export function getStockDetail(id) {
  return request.get(`/stock/stocks/${id}`)
}

/**
 * 获取库存操作日志
 * @param {number} stockId - 库存ID
 * @param {Object} params - 查询参数
 * @returns {Promise}
 */
export function getStockLogs(stockId, params) {
  return request.get(`/stock/stocks/${stockId}/logs`, { params })
}

/**
 * 获取库存操作记录（所有）
 * @param {Object} params - 查询参数
 * @returns {Promise}
 */
export function getStockOpLogs(params) {
  return request.get('/stock/stock-logs', { params })
}

/**
 * 调整库存
 * @param {number} id - 库存ID
 * @param {Object} data - 调整数据
 * @returns {Promise}
 */
export function adjustStock(id, data) {
  return request.post(`/stock/stocks/${id}/adjust`, data)
}

/**
 * 获取库存预警列表
 * @param {Object} params - 查询参数
 * @returns {Promise}
 */
export function getStockAlerts(params) {
  return request.get('/stock/stocks/alerts', { params })
}
