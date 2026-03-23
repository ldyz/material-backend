import request from '@/utils/request'

export function getRequisitions(params) {
  return request({
    url: '/requisition/requisitions',
    method: 'GET',
    params
  })
}

export function getRequisitionDetail(id) {
  return request({
    url: `/requisition/requisitions/${id}`,
    method: 'GET'
  })
}

export function createRequisition(data) {
  return request({
    url: '/requisition/requisitions',
    method: 'POST',
    data
  })
}

export function approveRequisition(id, data) {
  return request({
    url: `/requisition/requisitions/${id}/approve`,
    method: 'POST',
    data
  })
}

export function rejectRequisition(id, data) {
  return request({
    url: `/requisition/requisitions/${id}/reject`,
    method: 'POST',
    data
  })
}

export function resubmitRequisition(id, data) {
  return request({
    url: `/requisition/requisitions/${id}/resubmit`,
    method: 'POST',
    data
  })
}

export function issueRequisition(id, data) {
  return request({
    url: `/requisition/requisitions/${id}/issue`,
    method: 'POST',
    data
  })
}

// 获取项目列表
export function getProjects(params) {
  return request({
    url: '/projects',
    method: 'GET',
    params
  })
}

// 获取库存列表（按项目筛选）
export function getStock(params) {
  return request({
    url: '/stock',
    method: 'GET',
    params
  })
}
