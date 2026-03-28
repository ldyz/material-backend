import request from '@/utils/request'

export function getPlans(params) {
  return request({
    url: '/material-plan/plans',
    method: 'GET',
    params
  })
}

export function getPlanDetail(id) {
  return request({
    url: `/material-plan/plans/${id}`,
    method: 'GET'
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
