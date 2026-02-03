import request from '@/utils/request'

/**
 * 获取待审批任务列表
 * @param {Object} params - 查询参数
 * @returns {Promise} 返回待审批任务列表
 */
export function getPendingTasks(params) {
  return request({
    url: '/workflow-tasks/pending',
    method: 'GET',
    params
  })
}

/**
 * 审批通过任务
 * @param {number} id - 任务ID
 * @param {Object} data - 审批数据
 * @returns {Promise} 返回审批结果
 */
export function approveTask(id, data) {
  return request({
    url: `/workflow-tasks/${id}/approve`,
    method: 'POST',
    data
  })
}

/**
 * 审批拒绝任务
 * @param {number} id - 任务ID
 * @param {Object} data - 拒绝原因
 * @returns {Promise} 返回拒绝结果
 */
export function rejectTask(id, data) {
  return request({
    url: `/workflow-tasks/${id}/reject`,
    method: 'POST',
    data
  })
}

/**
 * 任务退回
 * @param {number} id - 任务ID
 * @param {Object} data - 退回原因
 * @returns {Promise} 返回退回结果
 */
export function returnTask(id, data) {
  return request({
    url: `/workflow-tasks/${id}/return`,
    method: 'POST',
    data
  })
}

/**
 * 任务评论
 * @param {number} id - 任务ID
 * @param {Object} data - 评论内容
 * @returns {Promise} 返回评论结果
 */
export function commentTask(id, data) {
  return request({
    url: `/workflow-tasks/${id}/comment`,
    method: 'POST',
    data
  })
}

/**
 * 获取工作流定义列表
 * @returns {Promise} 返回工作流定义列表
 */
export function getWorkflows() {
  return request({
    url: '/workflows',
    method: 'GET'
  })
}

/**
 * 获取工作流实例列表
 * @param {Object} params - 查询参数
 * @returns {Promise} 返回工作流实例列表
 */
export function getWorkflowInstances(params) {
  return request({
    url: '/workflow-instances',
    method: 'GET',
    params
  })
}

/**
 * 获取工作流实例详情
 * @param {number} id - 实例ID
 * @returns {Promise} 返回实例详情
 */
export function getWorkflowInstance(id) {
  return request({
    url: `/workflow-instances/${id}`,
    method: 'GET'
  })
}

/**
 * 获取实例审批历史
 * @param {number} id - 实例ID
 * @returns {Promise} 返回审批历史
 */
export function getInstanceApprovals(id) {
  return request({
    url: `/workflow-instances/${id}/approvals`,
    method: 'GET'
  })
}

/**
 * 重新提交实例
 * @param {number} id - 实例ID
 * @returns {Promise} 返回提交结果
 */
export function resubmitInstance(id) {
  return request({
    url: `/workflow-instances/${id}/resubmit`,
    method: 'POST'
  })
}
