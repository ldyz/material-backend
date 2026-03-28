/**
 * 状态映射和显示工具
 * 提供统一的状态类型和文本映射
 */

// 入库单状态映射
const inboundStatusMap = {
  draft: { type: 'default', text: '草稿' },
  pending: { type: 'warning', text: '待审批' },
  approved: { type: 'success', text: '已批准' },
  rejected: { type: 'danger', text: '已拒绝' },
  completed: { type: 'primary', text: '已完成' }
}

// 计划单状态映射
const planStatusMap = {
  draft: { type: 'default', text: '草稿' },
  pending: { type: 'warning', text: '待审批' },
  approved: { type: 'success', text: '已批准' },
  rejected: { type: 'danger', text: '已拒绝' },
  completed: { type: 'primary', text: '已完成' }
}

// 请购单状态映射
const requisitionStatusMap = {
  draft: { type: 'default', text: '草稿' },
  pending: { type: 'warning', text: '待审批' },
  approved: { type: 'success', text: '已批准' },
  rejected: { type: 'danger', text: '已拒绝' },
  completed: { type: 'primary', text: '已完成' }
}

// 预约单状态映射
const appointmentStatusMap = {
  draft: { type: 'default', text: '草稿' },
  pending: { type: 'warning', text: '待审批' },
  approved: { type: 'success', text: '已批准' },
  rejected: { type: 'danger', text: '已拒绝' },
  in_progress: { type: 'primary', text: '进行中' },
  completed: { type: 'success', text: '已完成' },
  cancelled: { type: 'danger', text: '已取消' }
}

// 获取状态映射表
function getStatusMap(type) {
  const maps = {
    inbound: inboundStatusMap,
    plan: planStatusMap,
    requisition: requisitionStatusMap,
    appointment: appointmentStatusMap
  }
  return maps[type] || inboundStatusMap
}

/**
 * 获取状态标签类型
 * @param {string} status - 状态值
 * @param {string} type - 业务类型 ('inbound', 'plan', 'requisition', 'appointment')
 * @returns {string} Vant Tag type
 */
export function getStatusType(status, type = 'inbound') {
  if (!status) return 'default'

  const statusMap = getStatusMap(type)
  const statusInfo = statusMap[status]
  return statusInfo?.type || 'default'
}

/**
 * 获取状态文本
 * @param {string} status - 状态值
 * @param {string} type - 业务类型 ('inbound', 'plan', 'requisition', 'appointment')
 * @returns {string} 状态显示文本
 */
export function getStatusText(status, type = 'inbound') {
  if (!status) return '-'

  const statusMap = getStatusMap(type)
  const statusInfo = statusMap[status]
  return statusInfo?.text || status
}

/**
 * 获取状态完整信息
 * @param {string} status - 状态值
 * @param {string} type - 业务类型
 * @returns {object} { type, text }
 */
export function getStatusInfo(status, type = 'inbound') {
  return {
    type: getStatusType(status, type),
    text: getStatusText(status, type)
  }
}

/**
 * 判断状态是否为指定值
 * @param {string} status - 当前状态
 * @param {string|string[]} targetStatus - 目标状态或状态数组
 * @returns {boolean}
 */
export function isStatus(status, targetStatus) {
  if (Array.isArray(targetStatus)) {
    return targetStatus.includes(status)
  }
  return status === targetStatus
}

/**
 * 判断是否可编辑
 * @param {string} status - 当前状态
 * @returns {boolean}
 */
export function canEdit(status) {
  return isStatus(status, ['draft', 'rejected'])
}

/**
 * 判断是否可审批
 * @param {string} status - 当前状态
 * @returns {boolean}
 */
export function canApprove(status) {
  return status === 'pending'
}

/**
 * 判断是否可删除
 * @param {string} status - 当前状态
 * @returns {boolean}
 */
export function canDelete(status) {
  return isStatus(status, ['draft', 'rejected'])
}

/**
 * 获取状态颜色（用于自定义样式）
 * @param {string} status - 状态值
 * @param {string} type - 业务类型
 * @returns {string} 颜色值
 */
export function getStatusColor(status, type = 'inbound') {
  const typeColorMap = {
    default: '#969799',
    primary: '#1989fa',
    success: '#07c160',
    warning: '#ff976a',
    danger: '#ee0a24'
  }
  const statusType = getStatusType(status, type)
  return typeColorMap[statusType] || typeColorMap.default
}

export function useStatus() {
  return {
    getStatusType,
    getStatusText,
    getStatusInfo,
    isStatus,
    canEdit,
    canApprove,
    canDelete,
    getStatusColor
  }
}
