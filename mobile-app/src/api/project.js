import request from '@/utils/request'

/**
 * 获取项目列表（仅返回当前用户关联的项目）
 * @param {Object} params - 查询参数
 * @param {string} params.search - 搜索关键词
 * @param {string} params.status - 状态筛选
 * @param {number} params.page - 页码
 * @param {number} params.page_size - 每页数量
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
 * @param {number} id - 项目ID
 */
export function getProjectDetail(id) {
  return request({
    url: `/project/projects/${id}`,
    method: 'GET'
  })
}

/**
 * 获取项目成员列表
 * @param {number} id - 项目ID
 */
export function getProjectMembers(id) {
  return request({
    url: `/project/projects/${id}/members`,
    method: 'GET'
  })
}

/**
 * 获取项目子项目列表
 * @param {number} id - 项目ID
 */
export function getProjectChildren(id) {
  return request({
    url: `/project/projects/${id}/children`,
    method: 'GET'
  })
}
