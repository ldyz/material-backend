import request from '@/utils/request'

/**
 * 施工日志管理 API (移动端)
 *
 * 提供施工日志的增删改查、图片上传等功能
 *
 * @module Construction Log API
 */

/**
 * 获取施工日志列表
 * @param {Object} params - 查询参数
 * @param {number} params.page - 页码
 * @param {number} params.pageSize - 每页数量
 * @param {number} params.project_id - 项目ID筛选
 * @returns {Promise} 返回施工日志列表
 */
export function getConstructionLogs(params) {
  return request.get('/construction_log/logs', { params })
}

/**
 * 获取施工日志详情
 * @param {number} id - 日志ID
 * @returns {Promise} 返回施工日志详情
 */
export function getConstructionLogDetail(id) {
  return request.get(`/construction_log/logs/${id}`)
}

/**
 * 创建施工日志
 * @param {Object} data - 日志数据
 * @returns {Promise} 返回创建的日志
 */
export function createConstructionLog(data) {
  return request.post('/construction_log/logs', data)
}

/**
 * 更新施工日志
 * @param {number} id - 日志ID
 * @param {Object} data - 更新数据
 * @returns {Promise} 返回更新后的日志
 */
export function updateConstructionLog(id, data) {
  return request.put(`/construction_log/logs/${id}`, data)
}

/**
 * 删除施工日志
 * @param {number} id - 日志ID
 * @returns {Promise} 返回删除结果
 */
export function deleteConstructionLog(id) {
  return request.delete(`/construction_log/logs/${id}`)
}

/**
 * 上传图片
 * @param {FormData} formData - 图片数据
 * @returns {Promise} 返回上传后的图片URL
 */
export function uploadImage(formData) {
  return request.post('/construction_log/upload_image', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })
}

/**
 * 获取天气信息
 * @param {string} date - 日期 (YYYY-MM-DD)
 * @returns {Promise} 返回天气信息
 */
export function getWeather(date) {
  return request.get('/construction_log/logs/weather', { params: { date } })
}
