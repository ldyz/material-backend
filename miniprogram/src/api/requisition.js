import request from '@/utils/request'

/**
 * 获取领料单列表
 */
export function getRequisitions(params) {
  return request({
    url: '/requisition/requisitions',
    method: 'GET',
    params
  })
}

/**
 * 获取领料单详情
 */
export function getRequisitionDetail(id) {
  return request({
    url: `/requisition/requisitions/${id}`,
    method: 'GET'
  })
}

/**
 * 创建领料单
 */
export function createRequisition(data) {
  return request({
    url: '/requisition/requisitions',
    method: 'POST',
    data
  })
}

/**
 * 批准领料单
 */
export function approveRequisition(id, data) {
  return request({
    url: `/requisition/requisitions/${id}/approve`,
    method: 'POST',
    data
  })
}

/**
 * 拒绝领料单
 */
export function rejectRequisition(id, data) {
  return request({
    url: `/requisition/requisitions/${id}/reject`,
    method: 'POST',
    data
  })
}

/**
 * 重新提交领料单
 */
export function resubmitRequisition(id, data) {
  return request({
    url: `/requisition/requisitions/${id}/resubmit`,
    method: 'POST',
    data
  })
}

/**
 * 获取领料单审批历史
 */
export function getRequisitionWorkflowHistory(id) {
  return request({
    url: `/requisition/requisitions/${id}/workflow-history`,
    method: 'GET'
  })
}

/**
 * 出库
 */
export function issueRequisition(id, data) {
  return request({
    url: `/requisition/requisitions/${id}/issue`,
    method: 'POST',
    data
  })
}
