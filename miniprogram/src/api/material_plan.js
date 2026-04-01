import request from '@/utils/request'

/**
 * 获取物资计划列表
 */
export function getPlans(params) {
  return request({
    url: '/material-plan/plans',
    method: 'GET',
    params
  })
}

/**
 * 获取已批准的计划列表
 */
export function getApprovedPlans(params) {
  return request({
    url: '/material-plan/plans',
    method: 'GET',
    params: { ...params, status: 'approved' }
  })
}

/**
 * 获取计划详情
 */
export function getPlanDetail(id) {
  return request({
    url: `/material-plan/plans/${id}`,
    method: 'GET'
  })
}

/**
 * 创建物资计划
 */
export function createPlan(data) {
  return request({
    url: '/material-plan/plans',
    method: 'POST',
    data
  })
}

/**
 * 批准计划
 */
export function approvePlan(id, data) {
  return request({
    url: `/material-plan/plans/${id}/approve`,
    method: 'POST',
    data
  })
}

/**
 * 拒绝计划
 */
export function rejectPlan(id, data) {
  return request({
    url: `/material-plan/plans/${id}/reject`,
    method: 'POST',
    data
  })
}

/**
 * 重新提交计划
 */
export function resubmitPlan(id, data) {
  return request({
    url: `/material-plan/plans/${id}/resubmit`,
    method: 'POST',
    data
  })
}
