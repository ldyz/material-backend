import request from '@/utils/request'

/**
 * 获取项目列表
 */
export function getProjects(params) {
  return request({
    url: '/project/projects',
    method: 'GET',
    params
  })
}

/**
 * 获取项目详情
 */
export function getProjectDetail(id) {
  return request({
    url: `/project/projects/${id}`,
    method: 'GET'
  })
}

/**
 * 获取项目成员
 */
export function getProjectMembers(id) {
  return request({
    url: `/project/projects/${id}/members`,
    method: 'GET'
  })
}

/**
 * 获取项目子项目
 */
export function getProjectChildren(id) {
  return request({
    url: `/project/projects/${id}/children`,
    method: 'GET'
  })
}
