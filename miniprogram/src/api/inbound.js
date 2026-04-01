import request from '@/utils/request'

/**
 * 获取入库单列表
 */
export function getInboundOrders(params) {
  return request({
    url: '/inbound/inbound-orders',
    method: 'GET',
    params
  })
}

/**
 * 获取入库单详情
 */
export function getInboundDetail(id) {
  return request({
    url: `/inbound/inbound-orders/${id}`,
    method: 'GET'
  })
}

/**
 * 创建入库单
 */
export function createInbound(data) {
  return request({
    url: '/inbound/inbound-orders',
    method: 'POST',
    data
  })
}

/**
 * 批准入库单
 */
export function approveInbound(id, data) {
  return request({
    url: `/inbound/inbound-orders/${id}/approve`,
    method: 'POST',
    data
  })
}

/**
 * 拒绝入库单
 */
export function rejectInbound(id, data) {
  return request({
    url: `/inbound/inbound-orders/${id}/reject`,
    method: 'POST',
    data
  })
}

/**
 * 重新提交入库单
 */
export function resubmitInbound(id, data) {
  return request({
    url: `/inbound/inbound-orders/${id}/resubmit`,
    method: 'POST',
    data
  })
}

/**
 * 获取入库单审批历史
 */
export function getInboundWorkflowHistory(id) {
  return request({
    url: `/inbound/inbound-orders/${id}/workflow-history`,
    method: 'GET'
  })
}

/**
 * 获取物资主数据列表
 */
export function getMaterialMasters(params) {
  return request({
    url: '/materials/master',
    method: 'GET',
    params
  })
}
