import request from '@/utils/request'

export function getPlans(params) {
  return request({
    url: '/material-plan/plans',
    method: 'GET',
    params
  })
}

// 获取已批准的计划列表
export function getApprovedPlans(params) {
  return request({
    url: '/material-plan/plans',
    method: 'GET',
    params: { ...params, status: 'approved' }
  })
}

export function getPlanDetail(id) {
  return request({
    url: `/material-plan/plans/${id}`,
    method: 'GET'
  })
}

export function createPlan(data) {
  return request({
    url: '/material-plan/plans',
    method: 'POST',
    data
  })
}

export function getProjects() {
  return request({
    url: '/project/projects',
    method: 'GET',
    params: { pageSize: 1000 }
  })
}

export function approvePlan(id, data) {
  return request({
    url: `/material-plan/plans/${id}/approve`,
    method: 'POST',
    data
  })
}

export function rejectPlan(id, data) {
  return request({
    url: `/material-plan/plans/${id}/reject`,
    method: 'POST',
    data
  })
}

export function resubmitPlan(id, data) {
  return request({
    url: `/material-plan/plans/${id}/resubmit`,
    method: 'POST',
    data
  })
}
