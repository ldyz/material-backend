import { computed } from 'vue'
import { useUserStore } from '@/stores/user'
import { PERMISSIONS } from '@/utils/constants'

export function usePermission() {
  const userStore = useUserStore()

  /**
   * 检查是否有指定权限
   * @param {string} permission - 权限标识
   * @returns {boolean}
   */
  function hasPermission(permission) {
    return userStore.hasPermission(permission)
  }

  /**
   * 检查是否有任一权限
   * @param {Array<string>} perms - 权限列表
   * @returns {boolean}
   */
  function hasAnyPermission(perms) {
    return userStore.hasAnyPermission(perms)
  }

  /**
   * 检查是否有所有权限
   * @param {Array<string>} perms - 权限列表
   * @returns {boolean}
   */
  function hasAllPermissions(perms) {
    return userStore.hasAllPermissions(perms)
  }

  /**
   * 检查是否是管理员
   * @returns {boolean}
   */
  const isAdmin = computed(() => userStore.isAdmin())

  // 入库权限
  const canViewInbound = computed(() => hasPermission(PERMISSIONS.INBOUND_VIEW))
  const canCreateInbound = computed(() => hasPermission(PERMISSIONS.INBOUND_CREATE))
  const canApproveInbound = computed(() => hasPermission(PERMISSIONS.INBOUND_APPROVE))
  const canDeleteInbound = computed(() => hasPermission(PERMISSIONS.INBOUND_DELETE))

  // 出库权限
  const canViewRequisition = computed(() => hasPermission(PERMISSIONS.REQUISITION_VIEW))
  const canCreateRequisition = computed(() => hasPermission(PERMISSIONS.REQUISITION_CREATE))
  const canApproveRequisition = computed(() => hasPermission(PERMISSIONS.REQUISITION_APPROVE))
  const canIssueRequisition = computed(() => hasPermission(PERMISSIONS.REQUISITION_ISSUE))
  const canDeleteRequisition = computed(() => hasPermission(PERMISSIONS.REQUISITION_DELETE))

  // 库存权限
  const canViewStock = computed(() => hasPermission(PERMISSIONS.STOCK_VIEW))
  const canAdjustStock = computed(() => hasPermission(PERMISSIONS.STOCK_ADJUST))

  // 施工日志权限
  const canViewConstructionLog = computed(() => hasPermission(PERMISSIONS.CONSTRUCTION_LOG_VIEW))
  const canCreateConstructionLog = computed(() => hasPermission(PERMISSIONS.CONSTRUCTION_LOG_CREATE))

  // AI 分析权限
  const canUseAI = computed(() => hasPermission(PERMISSIONS.AI_ANALYZE))

  return {
    hasPermission,
    hasAnyPermission,
    hasAllPermissions,
    isAdmin,
    canViewInbound,
    canCreateInbound,
    canApproveInbound,
    canDeleteInbound,
    canViewRequisition,
    canCreateRequisition,
    canApproveRequisition,
    canIssueRequisition,
    canDeleteRequisition,
    canViewStock,
    canAdjustStock,
    canViewConstructionLog,
    canCreateConstructionLog,
    canUseAI,
  }
}
