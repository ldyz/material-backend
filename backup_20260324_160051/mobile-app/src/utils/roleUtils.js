// 角色工具函数 - 统一使用 roles 数组进行角色判断

/**
 * 角色名称到代码的映射
 */
const ROLE_MAPPING = {
  // 中文角色名称
  '施工员': 'foreman',      // 施工员负责审批
  '作业人员': 'worker',     // 作业人员负责执行施工
  '项目经理': 'project_manager',
  '管理员': 'admin',
  '保管员': 'keeper',
  '材料员': 'material_staff',
  '分包材料员': 'subcontractor_material_staff',
  '预约管理员': 'appointment_admin',
  // 英文角色代码
  'foreman': 'foreman',
  'worker': 'worker',
  'admin': 'admin',
  'project_manager': 'project_manager',
  'keeper': 'keeper',
  'material_staff': 'material_staff',
  'subcontractor_material_staff': 'subcontractor_material_staff',
  'appointment_admin': 'appointment_admin'
}

/**
 * 角色代码对应的角色名称列表（用于反向查找）
 */
const ROLE_CODE_TO_NAMES = {
  'foreman': ['施工员', 'foreman'],
  'worker': ['作业人员', 'worker'],
  'admin': ['管理员', 'admin'],
  'project_manager': ['项目经理', 'project_manager'],
  'keeper': ['保管员', 'keeper'],
  'material_staff': ['材料员', 'material_staff'],
  'subcontractor_material_staff': ['分包材料员', 'subcontractor_material_staff'],
  'appointment_admin': ['预约管理员', 'appointment_admin']
}

/**
 * 从用户信息中获取角色代码列表
 * @param {Object} userInfo - 用户信息对象
 * @returns {string[]} 角色代码数组
 */
export function getUserRoleCodes(userInfo) {
  if (!userInfo) return []

  // 优先使用 roles 数组（新系统）
  if (userInfo.roles && Array.isArray(userInfo.roles) && userInfo.roles.length > 0) {
    return userInfo.roles.map(role => {
      const name = role.name || role
      return ROLE_MAPPING[name] || name
    })
  }

  // 兼容旧系统：使用 role 字符串字段
  if (userInfo.role) {
    return [userInfo.role]
  }

  return []
}

/**
 * 检查用户是否拥有指定角色
 * @param {Object} userInfo - 用户信息对象
 * @param {string|string[]} roleCodes - 要检查的角色代码或代码数组
 * @returns {boolean} 是否拥有该角色
 */
export function hasRole(userInfo, roleCodes) {
  const userRoles = getUserRoleCodes(userInfo)
  const checkRoles = Array.isArray(roleCodes) ? roleCodes : [roleCodes]
  return userRoles.some(role => checkRoles.includes(role))
}

/**
 * 检查用户是否是管理员
 * @param {Object} userInfo - 用户信息对象
 * @returns {boolean}
 */
export function isAdmin(userInfo) {
  return hasRole(userInfo, 'admin')
}

/**
 * 检查用户是否是施工员（负责审批）
 * @param {Object} userInfo - 用户信息对象
 * @returns {boolean}
 */
export function isForeman(userInfo) {
  return hasRole(userInfo, 'foreman')
}

/**
 * 检查用户是否是作业人员（执行施工）
 * @param {Object} userInfo - 用户信息对象
 * @returns {boolean}
 */
export function isWorker(userInfo) {
  return hasRole(userInfo, 'worker')
}

/**
 * 检查用户是否是项目经理
 * @param {Object} userInfo - 用户信息对象
 * @returns {boolean}
 */
export function isProjectManager(userInfo) {
  return hasRole(userInfo, 'project_manager')
}

/**
 * 检查用户是否可以管理预约单（管理员或项目经理或施工员）
 * @param {Object} userInfo - 用户信息对象
 * @returns {boolean}
 */
export function canManageAppointments(userInfo) {
  return isAdmin(userInfo) || isProjectManager(userInfo) || isForeman(userInfo)
}

/**
 * 检查用户是否可以执行施工（作业人员）
 * @param {Object} userInfo - 用户信息对象
 * @returns {boolean}
 */
export function canExecuteWork(userInfo) {
  return isWorker(userInfo)
}

/**
 * 获取用户的主要角色代码（用于显示和简单判断）
 * @param {Object} userInfo - 用户信息对象
 * @returns {string} 主要角色代码
 */
export function getPrimaryRoleCode(userInfo) {
  const roleCodes = getUserRoleCodes(userInfo)

  // 优先级顺序：admin > project_manager > foreman > worker > 其他
  if (roleCodes.includes('admin')) return 'admin'
  if (roleCodes.includes('project_manager')) return 'project_manager'
  if (roleCodes.includes('foreman')) return 'foreman'
  if (roleCodes.includes('worker')) return 'worker'

  return roleCodes[0] || ''
}

/**
 * 获取角色显示名称
 * @param {string} roleCode - 角色代码
 * @returns {string} 角色显示名称
 */
export function getRoleDisplayName(roleCode) {
  const displayNames = {
    'admin': '管理员',
    'foreman': '施工员',
    'worker': '作业人员',
    'project_manager': '项目经理',
    'keeper': '保管员',
    'material_staff': '材料员',
    'subcontractor_material_staff': '分包材料员',
    'appointment_admin': '预约管理员'
  }
  return displayNames[roleCode] || roleCode
}

/**
 * 将角色名称映射为角色代码
 * @param {string} roleName - 角色名称（中文或英文）
 * @returns {string} 角色代码
 */
export function mapRoleNameToCode(roleName) {
  return ROLE_MAPPING[roleName] || roleName
}
