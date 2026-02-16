/**
 * 审批工作流统一的映射工具
 * 确保移动端和web端显示一致
 */

// ========== 节点Key到中文名称的映射 ==========
const NODE_NAME_MAP = {
  // 通用节点
  'start': '提交申请',
  'submit': '提交审批',

  // 预约单审批节点
  'foreman_approve': '施工员审批',
  'project_manager_approve': '项目经理审批',
  'supervisor_approve': '主管审批',
  'director_approve': '总监审批',

  // 入库单审批节点
  'warehouse_manager_approve': '库管员审批',
  'finance_approve': '财务审批',

  // 计划单审批节点
  'department_head_approve': '部门主管审批',

  // 请购单审批节点
  'section_chief_approve': '科长审批',
  'procurement_approve': '采购审批'
}

// ========== Action到中文操作名称的映射 ==========
const ACTION_NAME_MAP = {
  'start': '提交',
  'submit': '提交审批',
  'approve': '审批通过',
  'reject': '审批拒绝',
  'return': '退回修改',
  'comment': '添加意见',
  'cancel': '取消审批'
}

// ========== Action到状态的映射 ==========
const ACTION_STATUS_MAP = {
  'start': 'approved',
  'submit': 'approved',
  'approve': 'approved',
  'reject': 'rejected',
  'return': 'returned',
  'comment': 'commented',
  'cancel': 'cancelled'
}

// ========== 状态文本映射 ==========
const STATUS_TEXT_MAP = {
  'approved': '已通过',
  'rejected': '已拒绝',
  'pending': '待审批',
  'returned': '已退回',
  'cancelled': '已取消',
  'commented': '已评论',
  'draft': '草稿',
  'in_progress': '进行中',
  'completed': '已完成'
}

// ========== 状态类型（用于标签颜色）映射 ==========
const STATUS_TYPE_MAP = {
  'approved': 'success',    // 绿色
  'rejected': 'danger',     // 红色
  'pending': 'warning',     // 橙色
  'returned': 'warning',    // 橙色
  'cancelled': 'default',  // 灰色
  'commented': 'primary',   // 蓝色
  'draft': 'info',          // 青色
  'in_progress': 'primary', // 蓝色
  'completed': 'success'   // 绿色
}

// ========== 颜色代码（用于自定义样式）==========
const STATUS_COLOR_MAP = {
  'approved': '#07c160',   // 绿色
  'rejected': '#ee0a24',   // 红色
  'pending': '#ff976a',    // 橙色
  'returned': '#ff976a',   // 橙色
  'cancelled': '#c8c9cc', // 灰色
  'commented': '#1989fa',  // 蓝色
  'draft': '#969799',      // 灰色
  'in_progress': '#1989fa',// 蓝色
  'completed': '#07c160'  // 绿色
}

/**
 * 获取节点中文名称
 * @param {string} nodeKey - 节点key
 * @returns {string} 中文名称
 */
export function getNodeName(nodeKey) {
  return NODE_NAME_MAP[nodeKey] || nodeKey
}

/**
 * 获取操作中文名称
 * @param {string} action - 操作类型
 * @returns {string} 中文操作名称
 */
export function getActionName(action) {
  return ACTION_NAME_MAP[action] || action
}

/**
 * 根据action获取状态
 * @param {string} action - 操作类型
 * @returns {string} 状态
 */
export function getActionStatus(action) {
  return ACTION_STATUS_MAP[action] || 'pending'
}

/**
 * 获取状态文本
 * @param {string} status - 状态值
 * @returns {string} 中文文本
 */
export function getStatusText(status) {
  return STATUS_TEXT_MAP[status] || status
}

/**
 * 获取状态类型（用于标签）
 * @param {string} status - 状态值
 * @returns {string} 类型
 */
export function getStatusType(status) {
  return STATUS_TYPE_MAP[status] || 'default'
}

/**
 * 获取状态颜色代码
 * @param {string} status - 状态值
 * @returns {string} 颜色代码
 */
export function getStatusColor(status) {
  return STATUS_COLOR_MAP[status] || '#969799'
}

/**
 * 格式化审批记录为统一格式
 * @param {object} log - 原始审批记录
 * @returns {object} 格式化后的记录
 */
export function formatApprovalLog(log) {
  const nodeName = getNodeName(log.node_name || log.node_key)
  const actionName = getActionName(log.action)
  const status = getActionStatus(log.action)
  const statusText = getStatusText(status)

  return {
    id: log.id,
    title: nodeName,
    action: log.action,
    action_name: actionName,
    status: status,
    status_text: statusText,
    approver: log.approver_name || log.actor_name,
    approver_role: log.approver_role,
    approved_at: log.created_at || log.approved_at,
    remark: log.remark || log.action_data || log.comment
  }
}

/**
 * 批量格式化审批记录
 * @param {Array} logs - 审批记录数组
 * @returns {Array} 格式化后的记录数组
 */
export function formatApprovalLogs(logs) {
  if (!Array.isArray(logs)) return []
  return logs.map(log => formatApprovalLog(log))
}

export function useApprovalWorkflow() {
  return {
    getNodeName,
    getActionName,
    getActionStatus,
    getStatusText,
    getStatusType,
    getStatusColor,
    formatApprovalLog,
    formatApprovalLogs
  }
}
