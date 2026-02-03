import request from '@/utils/request'

// 便利导出，保持与命名规范一致
export const getPlans = getMaterialPlans
export const getPlanDetail = getMaterialPlanDetail
export const approvePlan = approveMaterialPlan
export const rejectPlan = rejectMaterialPlan

/**
 * 获取物资计划列表
 * @param {Object} params - 查询参数
 * @returns {Promise} 返回物资计划列表
 */
export function getMaterialPlans(params) {
  return request({
    url: '/material-plan/plans',
    method: 'GET',
    params
  })
}

/**
 * 获取物资计划详情
 * @param {number} id - 计划ID
 * @returns {Promise} 返回物资计划详情
 */
export function getMaterialPlanDetail(id) {
  return request({
    url: `/material-plan/plans/${id}`,
    method: 'GET'
  })
}

/**
 * 创建物资计划
 * @param {Object} data - 计划数据
 * @returns {Promise} 返回创建的计划
 */
export function createMaterialPlan(data) {
  return request({
    url: '/material-plan/plans',
    method: 'POST',
    data
  })
}

/**
 * 更新物资计划
 * @param {number} id - 计划ID
 * @param {Object} data - 更新数据
 * @returns {Promise} 返回更新后的计划
 */
export function updateMaterialPlan(id, data) {
  return request({
    url: `/material-plan/plans/${id}`,
    method: 'PUT',
    data
  })
}

/**
 * 删除物资计划
 * @param {number} id - 计划ID
 * @returns {Promise} 返回删除结果
 */
export function deleteMaterialPlan(id) {
  return request({
    url: `/material-plan/plans/${id}`,
    method: 'DELETE'
  })
}

/**
 * 提交物资计划
 * @param {number} id - 计划ID
 * @returns {Promise} 返回提交结果
 */
export function submitMaterialPlan(id) {
  return request({
    url: `/material-plan/plans/${id}/submit`,
    method: 'POST'
  })
}

/**
 * 审批通过物资计划
 * @param {number} id - 计划ID
 * @param {Object} data - 审批数据
 * @returns {Promise} 返回审批结果
 */
export function approveMaterialPlan(id, data) {
  return request({
    url: `/material-plan/plans/${id}/approve`,
    method: 'POST',
    data
  })
}

/**
 * 审批拒绝物资计划
 * @param {number} id - 计划ID
 * @param {Object} data - 拒绝原因
 * @returns {Promise} 返回拒绝结果
 */
export function rejectMaterialPlan(id, data) {
  return request({
    url: `/material-plan/plans/${id}/reject`,
    method: 'POST',
    data
  })
}

/**
 * 激活物资计划
 * @param {number} id - 计划ID
 * @returns {Promise} 返回激活结果
 */
export function activateMaterialPlan(id) {
  return request({
    url: `/material-plan/plans/${id}/activate`,
    method: 'POST'
  })
}

/**
 * 获取计划物资列表
 * @param {number} id - 计划ID
 * @returns {Promise} 返回物资列表
 */
export function getPlanItems(id) {
  return request({
    url: `/material-plan/plans/${id}/items`,
    method: 'GET'
  })
}

/**
 * 添加计划物资
 * @param {number} id - 计划ID
 * @param {Object} data - 物资数据
 * @returns {Promise} 返回添加结果
 */
export function addPlanItem(id, data) {
  return request({
    url: `/material-plan/plans/${id}/items`,
    method: 'POST',
    data
  })
}

/**
 * 同步物资ID（批量更新）
 * @param {number} id - 计划ID
 * @param {Object} data - 包含items数组的数据
 * @returns {Promise} 返回同步结果
 */
export function syncPlanMaterialIds(id, data) {
  return request({
    url: `/material-plan/plans/${id}/sync-materials`,
    method: 'POST',
    data
  })
}

/**
 * 获取待审批任务列表
 * @returns {Promise} 返回待审批任务
 */
export function getPendingTasks() {
  return request({
    url: '/material-plan/workflow/pending',
    method: 'GET'
  })
}
