import request from '@/utils/request'

/**
 * 获取领料单列表
 * @param {Object} params - 查询参数
 * @returns {Promise}
 */
export function getRequisitions(params) {
  return request.get('/requisition/requisitions', { params })
}

/**
 * 获取领料单详情
 * @param {number} id - 领料单ID
 * @returns {Promise}
 */
export function getRequisitionDetail(id) {
  return request.get(`/requisition/requisitions/${id}`)
}

/**
 * 创建领料单
 * @param {Object} data - 领料单数据
 * @returns {Promise}
 */
export function createRequisition(data) {
  return request.post('/requisition/requisitions', data)
}

/**
 * 更新领料单
 * @param {number} id - 领料单ID
 * @param {Object} data - 更新数据
 * @returns {Promise}
 */
export function updateRequisition(id, data) {
  return request.put(`/requisition/requisitions/${id}`, data)
}

/**
 * 删除领料单
 * @param {number} id - 领料单ID
 * @returns {Promise}
 */
export function deleteRequisition(id) {
  return request.delete(`/requisition/requisitions/${id}`)
}

/**
 * 审批通过领料单
 * @param {number} id - 领料单ID
 * @param {Object} data - 审批数据 { items: [{id, approved_quantity}], remark: string }
 * @returns {Promise}
 */
export function approveRequisition(id, data) {
  // 兼容 notes 和 remark 参数
  const requestData = {
    ...data,
    remark: data?.remark || data?.notes || ''
  }
  // 如果原来有 items，保留
  if (data?.items) {
    requestData.items = data.items
  }
  return request.post(`/requisition/requisitions/${id}/approve`, requestData)
}

/**
 * 拒绝领料单
 * @param {number} id - 领料单ID
 * @param {Object} data - 拒绝数据 { remark: string }
 * @returns {Promise}
 */
export function rejectRequisition(id, data) {
  return request.post(`/requisition/requisitions/${id}/reject`, {
    remark: data?.remark || data?.notes || ''
  })
}

/**
 * 发料（确认出库）
 * @param {number} id - 领料单ID
 * @param {Object} data - 发料数据 { items: [{id, actual_quantity}], remark: string }
 * @returns {Promise}
 */
export function issueRequisition(id, data) {
  // 兼容 notes 和 remark 参数
  const requestData = {
    ...data,
    remark: data?.remark || data?.notes || ''
  }
  // 如果原来有 items，保留
  if (data?.items) {
    requestData.items = data.items
  }
  return request.post(`/requisition/requisitions/${id}/issue`, requestData)
}

/**
 * 获取待审批领料单数量
 * @returns {Promise}
 */
export function getPendingRequisitionCount() {
  return request.get('/requisition/requisitions/pending/count')
}
