import request from '@/utils/request'

/**
 * 获取施工日志列表
 * @param {Object} params - 查询参数
 * @returns {Promise}
 */
export function getConstructionLogs(params) {
  return request.get('/construction_log/logs', { params })
}

/**
 * 获取施工日志详情
 * @param {number} id - 日志ID
 * @returns {Promise}
 */
export function getConstructionLogDetail(id) {
  return request.get(`/construction_log/${id}`)
}

/**
 * 创建施工日志
 * @param {Object} data - 日志数据
 * @returns {Promise}
 */
export function createConstructionLog(data) {
  return request.post('/construction_log/', data)
}

/**
 * 更新施工日志
 * @param {number} id - 日志ID
 * @param {Object} data - 更新数据
 * @returns {Promise}
 */
export function updateConstructionLog(id, data) {
  return request.put(`/construction_log/${id}`, data)
}

/**
 * 删除施工日志
 * @param {number} id - 日志ID
 * @returns {Promise}
 */
export function deleteConstructionLog(id) {
  return request.delete(`/construction_log/${id}`)
}

/**
 * 上传图片
 * @param {FormData} formData - 图片数据
 * @returns {Promise}
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
 * @param {string} date - 日期
 * @returns {Promise}
 */
export function getWeather(date) {
  return request.get('/construction_log/logs/weather', { params: { date } })
}
