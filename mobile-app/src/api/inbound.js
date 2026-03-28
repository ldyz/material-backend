import request from '@/utils/request'

export function getInboundOrders(params) {
  return request({
    url: '/inbound/inbound-orders',
    method: 'GET',
    params
  })
}

export function getInboundDetail(id) {
  return request({
    url: `/inbound/inbound-orders/${id}`,
    method: 'GET'
  })
}

export function createInbound(data) {
  return request({
    url: '/inbound/inbound-orders',
    method: 'POST',
    data
  })
}

export function approveInbound(id, data) {
  return request({
    url: `/inbound/inbound-orders/${id}/approve`,
    method: 'POST',
    data
  })
}

export function rejectInbound(id, data) {
  return request({
    url: `/inbound/inbound-orders/${id}/reject`,
    method: 'POST',
    data
  })
}

export function resubmitInbound(id, data) {
  return request({
    url: `/inbound/inbound-orders/${id}/resubmit`,
    method: 'POST',
    data
  })
}

// 获取已批准的计划列表（用于创建入库单）
export function getApprovedPlans(params) {
  return request({
    url: '/material-plan/plans',
    method: 'GET',
    params: {
      ...params,
      status: 'approved'
    }
  })
}

// 获取项目列表
export function getProjects(params) {
  return request({
    url: '/project/projects',
    method: 'GET',
    params
  })
}

// 获取物资主数据列表（用于物资搜索）
export function getMaterialMasters(params) {
  return request({
    url: '/material-master',
    method: 'GET',
    params
  })
}

// 获取计划详情（包含物料列表）
export function getPlanDetail(id) {
  return request({
    url: `/material-plan/plans/${id}`,
    method: 'GET'
  })
}
