import request from '@/utils/request'

// 便利导出，用于移动端物资浏览
export const getMasterList = getMaterialMasterList
export const getMasterDetail = getMaterialMasterDetail

// ==================== 物资主数据 API ====================

/**
 * 获取物资主数据列表
 * @param {Object} params - 查询参数
 * @returns {Promise} 返回物资主数据列表
 */
export function getMaterialMasterList(params) {
  return request({
    url: '/materials/master',
    method: 'GET',
    params
  })
}

/**
 * 获取物资主数据详情
 * @param {number} id - 物资主数据ID
 * @returns {Promise} 返回物资主数据详情
 */
export function getMaterialMasterDetail(id) {
  return request({
    url: `/materials/master/${id}`,
    method: 'GET'
  })
}

/**
 * 创建物资主数据
 * @param {Object} data - 物资主数据信息
 * @returns {Promise} 返回创建的物资主数据
 */
export function createMaterialMaster(data) {
  return request({
    url: '/materials/master',
    method: 'POST',
    data
  })
}

/**
 * 更新物资主数据
 * @param {number} id - 物资主数据ID
 * @param {Object} data - 更新的物资主数据信息
 * @returns {Promise} 返回更新后的物资主数据
 */
export function updateMaterialMaster(id, data) {
  return request({
    url: `/materials/master/${id}`,
    method: 'PUT',
    data
  })
}

/**
 * 删除物资主数据
 * @param {number} id - 物资主数据ID
 * @returns {Promise} 返回删除结果
 */
export function deleteMaterialMaster(id) {
  return request({
    url: `/materials/master/${id}`,
    method: 'DELETE'
  })
}

/**
 * 获取项目物资列表（带库存）
 * @param {Object} params - 查询参数
 * @returns {Promise} 返回项目物资列表
 */
export function getProjectMaterials(params) {
  return request({
    url: '/materials/master/project',
    method: 'GET',
    params
  })
}

// ==================== 物资管理 API ====================

/**
 * 获取物资列表
 * @param {Object} params - 查询参数
 * @returns {Promise} 返回物资列表和分页信息
 */
export function getMaterials(params) {
  return request({
    url: '/material/materials',
    method: 'GET',
    params
  })
}

/**
 * 获取物资详情
 * @param {number} id - 物资ID
 * @returns {Promise} 返回物资详细信息
 */
export function getMaterialDetail(id) {
  return request({
    url: `/material/materials/${id}`,
    method: 'GET'
  })
}

/**
 * 创建新物资
 * @param {Object} data - 物资信息
 * @returns {Promise} 返回创建的物资信息
 */
export function createMaterial(data) {
  return request({
    url: '/material/materials',
    method: 'POST',
    data
  })
}

/**
 * 更新物资信息
 * @param {number} id - 物资ID
 * @param {Object} data - 要更新的物资信息
 * @returns {Promise} 返回更新后的物资信息
 */
export function updateMaterial(id, data) {
  return request({
    url: `/material/materials/${id}`,
    method: 'PUT',
    data
  })
}

/**
 * 删除物资
 * @param {number} id - 物资ID
 * @returns {Promise} 返回删除结果
 */
export function deleteMaterial(id) {
  return request({
    url: `/material/materials/${id}`,
    method: 'DELETE'
  })
}

/**
 * 导出物资数据
 * @param {Object} params - 查询参数
 * @returns {Promise} 返回 Excel 文件 Blob
 */
export function exportMaterials(params) {
  return request({
    url: '/material/materials/export',
    method: 'GET',
    params,
    responseType: 'blob'
  })
}

/**
 * 获取物资日志
 * @param {number} id - 物资ID
 * @returns {Promise} 返回物资日志列表
 */
export function getMaterialLogs(id) {
  return request({
    url: `/materials/${id}/logs`,
    method: 'GET'
  })
}

// ==================== 物资分类 API ====================

/**
 * 获取所有物资分类
 * @returns {Promise} 返回分类列表
 */
export function getMaterialCategories() {
  return request({
    url: '/material/categories',
    method: 'GET'
  })
}

/**
 * 获取单个分类详情
 * @param {number} id - 分类ID
 * @returns {Promise} 返回分类详细信息
 */
export function getMaterialCategory(id) {
  return request({
    url: `/material/categories/${id}`,
    method: 'GET'
  })
}

/**
 * 创建物资分类
 * @param {Object} data - 分类信息
 * @returns {Promise} 返回创建的分类信息
 */
export function createMaterialCategory(data) {
  return request({
    url: '/material/categories',
    method: 'POST',
    data
  })
}

/**
 * 更新物资分类
 * @param {number} id - 分类ID
 * @param {Object} data - 要更新的分类信息
 * @returns {Promise} 返回更新后的分类信息
 */
export function updateMaterialCategory(id, data) {
  return request({
    url: `/material/categories/${id}`,
    method: 'PUT',
    data
  })
}

/**
 * 删除物资分类
 * @param {number} id - 分类ID
 * @returns {Promise} 返回删除结果
 */
export function deleteMaterialCategory(id) {
  return request({
    url: `/material/categories/${id}`,
    method: 'DELETE'
  })
}

/**
 * 批量更新分类排序
 * @param {Array} sorts - 排序数组
 * @returns {Promise} 返回更新结果
 */
export function updateMaterialCategorySort(sorts) {
  return request({
    url: '/material/categories/sort',
    method: 'POST',
    data: { sorts }
  })
}
