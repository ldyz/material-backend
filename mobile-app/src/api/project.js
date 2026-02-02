import request from '@/utils/request'

/**
 * 获取项目列表
 * @param {Object} params - 查询参数
 * @returns {Promise}
 */
export function getProjects(params) {
  return request.get('/project/projects', { params })
}

/**
 * 获取项目详情
 * @param {number} id - 项目ID
 * @returns {Promise}
 */
export function getProjectDetail(id) {
  return request.get(`/project/projects/${id}`)
}

/**
 * 创建项目
 * @param {Object} data - 项目数据
 * @returns {Promise}
 */
export function createProject(data) {
  return request.post('/project/projects', data)
}

/**
 * 更新项目
 * @param {number} id - 项目ID
 * @param {Object} data - 更新数据
 * @returns {Promise}
 */
export function updateProject(id, data) {
  return request.put(`/project/projects/${id}`, data)
}

/**
 * 删除项目
 * @param {number} id - 项目ID
 * @returns {Promise}
 */
export function deleteProject(id) {
  return request.delete(`/project/projects/${id}`)
}
