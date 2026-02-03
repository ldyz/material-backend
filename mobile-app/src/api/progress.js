import request from '@/utils/request'

/**
 * 进度管理 API (移动端)
 *
 * 提供项目进度计划的查询功能
 *
 * @module Progress API
 */

/**
 * 获取进度计划列表
 * @param {Object} params - 查询参数
 * @param {number} params.project_id - 项目ID
 * @returns {Promise} 返回进度计划列表
 */
export function getProgressList(params) {
  return request.get('/progress', { params })
}

/**
 * 获取进度计划详情
 * @param {number} id - 进度ID
 * @returns {Promise} 返回进度详情
 */
export function getProgressDetail(id) {
  return request.get(`/progress/${id}`)
}

/**
 * 获取项目进度概览
 * @param {number} projectId - 项目ID
 * @returns {Promise} 返回项目进度概览
 */
export function getProjectProgress(projectId) {
  return request.get('/progress/project/' + projectId)
}

/**
 * 获取任务进度列表
 * @param {number} scheduleId - 进度计划ID
 * @returns {Promise} 返回任务进度列表
 */
export function getTaskProgress(scheduleId) {
  return request.get(`/progress/${scheduleId}/tasks`)
}
