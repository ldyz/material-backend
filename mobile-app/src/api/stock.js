import request from '@/utils/request'

// 便利导出
export const getStockList = getStocks

/**
 * 获取库存列表
 * @param {Object} params - 查询参数
 * @returns {Promise} 返回库存列表和分页信息
 */
export function getStocks(params) {
  return request({
    url: '/stock/stocks',
    method: 'GET',
    params
  })
}

/**
 * 获取库存预警列表
 * @param {Object} params - 查询参数
 * @returns {Promise} 返回库存预警列表
 */
export function getStockAlerts(params) {
  return request({
    url: '/stock/stocks/alerts',
    method: 'GET',
    params
  })
}

/**
 * 获取库存详情
 * @param {number} id - 库存ID
 * @returns {Promise} 返回库存详情
 */
export function getStockDetail(id) {
  return request({
    url: `/stock/stocks/${id}`,
    method: 'GET'
  })
}

/**
 * 创建库存记录
 * @param {Object} data - 库存信息
 * @returns {Promise} 返回创建的库存
 */
export function createStock(data) {
  return request({
    url: '/stock/stocks',
    method: 'POST',
    data
  })
}

/**
 * 更新库存
 * @param {number} id - 库存ID
 * @param {Object} data - 更新的库存信息
 * @returns {Promise} 返回更新后的库存
 */
export function updateStock(id, data) {
  return request({
    url: `/stock/stocks/${id}`,
    method: 'PUT',
    data
  })
}

/**
 * 删除库存
 * @param {number} id - 库存ID
 * @returns {Promise} 返回删除结果
 */
export function deleteStock(id) {
  return request({
    url: `/stock/stocks/${id}`,
    method: 'DELETE'
  })
}

/**
 * 获取库存变动日志
 * @param {Object} params - 查询参数
 * @returns {Promise} 返回库存日志列表
 */
export function getStockLogs(params) {
  return request({
    url: '/stock/stock-logs',
    method: 'GET',
    params
  })
}

/**
 * 获取库存操作日志（按库存ID）
 * @param {number} stockId - 库存ID
 * @param {Object} params - 查询参数
 * @returns {Promise} 返回库存日志列表
 */
export function getStockLogsById(stockId, params) {
  return request({
    url: `/stock/stocks/${stockId}/logs`,
    method: 'GET',
    params
  })
}

/**
 * 删除库存日志
 * @param {number} id - 日志ID
 * @returns {Promise} 返回删除结果
 */
export function deleteStockLog(id) {
  return request({
    url: `/stock/stock-logs/${id}`,
    method: 'DELETE'
  })
}

/**
 * 物资入库
 * @param {Object} data - 入库信息
 * @returns {Promise} 返回入库后的库存信息
 */
export function stockIn(data) {
  return request({
    url: `/stock/stocks/${data.id}/in`,
    method: 'POST',
    data
  })
}

/**
 * 物资出库
 * @param {Object} data - 出库信息
 * @returns {Promise} 返回出库后的库存信息
 */
export function stockOut(data) {
  return request({
    url: `/stock/stocks/${data.id}/out`,
    method: 'POST',
    data
  })
}

/**
 * 库存调整
 * @param {number} id - 库存ID
 * @param {Object} data - 调整信息
 * @returns {Promise} 返回调整结果
 */
export function adjustStock(id, data) {
  return request({
    url: `/stock/stocks/${id}/adjust`,
    method: 'POST',
    data
  })
}

/**
 * 导出库存数据
 * @param {Object} params - 查询参数
 * @returns {Promise} 返回 Excel 文件 Blob
 */
export function exportStocks(params) {
  return request({
    url: '/stock/stocks/export',
    method: 'GET',
    params,
    responseType: 'blob'
  })
}
