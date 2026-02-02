// 权限常量
export const PERMISSIONS = {
  // 入库权限
  INBOUND_VIEW: 'inbound_view',
  INBOUND_CREATE: 'inbound_create',
  INBOUND_APPROVE: 'inbound_approve',
  INBOUND_DELETE: 'inbound_delete',

  // 出库权限
  REQUISITION_VIEW: 'requisition_view',
  REQUISITION_CREATE: 'requisition_create',
  REQUISITION_APPROVE: 'requisition_approve',
  REQUISITION_ISSUE: 'requisition_issue',
  REQUISITION_DELETE: 'requisition_delete',

  // 库存权限
  STOCK_VIEW: 'stock_view',
  STOCK_ADJUST: 'stock_adjust',

  // 施工日志
  CONSTRUCTION_LOG_VIEW: 'construction_log_view',
  CONSTRUCTION_LOG_CREATE: 'construction_log_create',

  // AI 分析
  AI_ANALYZE: 'ai_analyze',
}

// 状态常量
export const INBOUND_STATUS = {
  PENDING: 'pending',
  APPROVED: 'approved',
  COMPLETED: 'completed',
  REJECTED: 'rejected',
}

export const INBOUND_STATUS_TEXT = {
  [INBOUND_STATUS.PENDING]: '待审批',
  [INBOUND_STATUS.APPROVED]: '已审批',
  [INBOUND_STATUS.COMPLETED]: '已完成',
  [INBOUND_STATUS.REJECTED]: '已拒绝',
}

export const REQUISITION_STATUS = {
  PENDING: 'pending',
  APPROVED: 'approved',
  ISSUED: 'issued',
  REJECTED: 'rejected',
}

export const REQUISITION_STATUS_TEXT = {
  [REQUISITION_STATUS.PENDING]: '待审批',
  [REQUISITION_STATUS.APPROVED]: '已审批',
  [REQUISITION_STATUS.ISSUED]: '已发料',
  [REQUISITION_STATUS.REJECTED]: '已拒绝',
}

// 本地存储键名
export const STORAGE_KEYS = {
  TOKEN: 'token',
  USER_INFO: 'user_info',
  PERMISSIONS: 'permissions',
}
