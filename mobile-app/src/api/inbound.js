import request from '@/utils/request'

/**
 * 获取入库单列表
 * @param {Object} params - 查询参数
 * @returns {Promise}
 */
export function getInboundOrders(params) {
  return request.get('/inbound/inbound-orders', { params })
}

/**
 * 获取入库单详情
 * @param {number} id - 入库单ID
 * @returns {Promise}
 */
export function getInboundDetail(id) {
  return request.get(`/inbound/inbound-orders/${id}`)
}

/**
 * 创建入库单
 * @param {Object} data - 入库单数据
 * @returns {Promise}
 */
export function createInbound(data) {
  return request.post('/inbound/inbound-orders', data)
}

/**
 * 更新入库单
 * @param {number} id - 入库单ID
 * @param {Object} data - 更新数据
 * @returns {Promise}
 */
export function updateInbound(id, data) {
  return request.put(`/inbound/inbound-orders/${id}`, data)
}

/**
 * 删除入库单
 * @param {number} id - 入库单ID
 * @returns {Promise}
 */
export function deleteInbound(id) {
  return request.delete(`/inbound/inbound-orders/${id}`)
}

/**
 * 审批通过入库单
 * @param {number} id - 入库单ID
 * @param {Object} data - 审批数据 { items: [{id, approved_quantity}], remark: string }
 * @returns {Promise}
 */
export function approveInbound(id, data) {
  // 兼容 notes 和 remark 参数
  const requestData = {
    remark: data?.remark || data?.notes || ''
  }
  // 如果原来有 items，保留
  if (data?.items) {
    requestData.items = data.items
  }
  return request.post(`/inbound/inbound-orders/${id}/approve`, requestData)
}

/**
 * 拒绝入库单
 * @param {number} id - 入库单ID
 * @param {Object} data - 拒绝数据 { remark: string }
 * @returns {Promise}
 */
export function rejectInbound(id, data) {
  return request.post(`/inbound/inbound-orders/${id}/reject`, {
    remark: data?.remark || data?.notes || ''
  })
}

/**
 * 完成入库单
 * @param {number} id - 入库单ID
 * @returns {Promise}
 */
export function completeInbound(id) {
  return request.post(`/inbound/inbound-orders/${id}/complete`)
}

/**
 * 获取待审批入库单数量
 * @returns {Promise}
 */
export function getPendingInboundCount() {
  return request.get('/inbound/inbound-orders/pending/count')
}
