import request from '@/utils/request'

/**
 * 项目管理 API (移动端)
 *
 * 提供项目的增删改查功能，支持分页、搜索和筛选
 *
 * @module Project API
 */

/**
 * 获取项目列表
 * @param {Object} params - 查询参数
 * @param {number} params.page - 页码（从1开始）
 * @param {number} params.pageSize - 每页数量
 * @param {string} params.search - 搜索关键词（项目名称、编号）
 * @param {string} params.status - 状态筛选
 * @returns {Promise} 返回项目列表和分页信息
 */
export function getProjects(params) {
  return request.get('/project/projects', { params })
}

/**
 * 获取项目详情
 * @param {number} id - 项目ID
 * @returns {Promise} 返回项目详细信息
 */
export function getProjectDetail(id) {
  return request.get(`/project/projects/${id}`)
}

/**
 * 创建新项目
 * @param {Object} data - 项目信息
 * @returns {Promise} 返回创建的项目信息
 */
export function createProject(data) {
  return request.post('/project/projects', data)
}

/**
 * 更新项目信息
 * @param {number} id - 项目ID
 * @param {Object} data - 要更新的项目信息
 * @returns {Promise} 返回更新后的项目信息
 */
export function updateProject(id, data) {
  return request.put(`/project/projects/${id}`, data)
}

/**
 * 删除项目
 * @param {number} id - 项目ID
 * @returns {Promise} 返回删除结果
 */
export function deleteProject(id) {
  return request.delete(`/project/projects/${id}`)
}

/**
 * 获取项目成员列表
 * @param {number} id - 项目ID
 * @returns {Promise} 返回项目成员列表
 */
export function getProjectMembers(id) {
  return request.get(`/project/projects/${id}/members`)
}

/**
 * 添加项目成员（批量替换）
 * @param {number} id - 项目ID
 * @param {Object} data - 成员信息
 * @param {number[]} data.user_ids - 用户ID数组
 * @returns {Promise} 返回添加结果
 */
export function addProjectMembers(id, data) {
  return request.post(`/project/projects/${id}/members`, data)
}

/**
 * 删除项目成员
 * @param {number} id - 项目ID
 * @param {number} userId - 用户ID
 * @returns {Promise} 返回删除结果
 */
export function removeProjectMember(id, userId) {
  return request.delete(`/project/projects/${id}/members/${userId}`)
}

/**
 * 获取项目树
 * @param {number} id - 项目ID
 * @returns {Promise} 返回项目树结构
 */
export function getProjectTree(id) {
  return request.get(`/project/projects/${id}/tree`)
}

/**
 * 获取子项目列表
 * @param {number} id - 项目ID
 * @returns {Promise} 返回子项目列表
 */
export function getProjectChildren(id) {
  return request.get(`/project/projects/${id}/children`)
}

/**
 * 聚合子项目进度
 * @param {number} id - 项目ID
 * @returns {Promise} 返回聚合后的进度
 */
export function aggregateProjectProgress(id) {
  return request.post(`/project/projects/${id}/aggregate-progress`)
}
