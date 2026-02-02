import request from '@/utils/request'

/**
 * 获取物资列表
 * @param {Object} params - 查询参数
 * @returns {Promise}
 */
export function getMaterials(params) {
  return request.get('/material/materials', { params })
}

/**
 * 获取物资详情
 * @param {number} id - 物资ID
 * @returns {Promise}
 */
export function getMaterialDetail(id) {
  return request.get(`/material/materials/${id}`)
}

/**
 * 创建物资
 * @param {Object} data - 物资数据
 * @returns {Promise}
 */
export function createMaterial(data) {
  return request.post('/material/materials', data)
}

/**
 * 更新物资
 * @param {number} id - 物资ID
 * @param {Object} data - 更新数据
 * @returns {Promise}
 */
export function updateMaterial(id, data) {
  return request.put(`/material/materials/${id}`, data)
}

/**
 * 删除物资
 * @param {number} id - 物资ID
 * @returns {Promise}
 */
export function deleteMaterial(id) {
  return request.delete(`/material/materials/${id}`)
}
