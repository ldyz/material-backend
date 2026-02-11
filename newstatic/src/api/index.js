/**
 * API 接口统一管理
 *
 * 本文件集中管理所有后端 API 接口的调用
 * 使用 axios 实例进行 HTTP 请求，baseURL 已配置为 '/api'
 *
 * @module API
 * @author Material Management System
 * @date 2025-01-27
 */

import request from './request'

// ==================== 认证相关 API ====================

/**
 * 认证相关 API 接口
 *
 * 提供用户登录、登出、获取用户信息、修改密码等功能
 *
 * @namespace authApi
 */
export const authApi = {
  /**
   * 用户登录
   *
   * @param {Object} data - 登录信息
   * @param {string} data.username - 用户名
   * @param {string} data.password - 密码
   * @returns {Promise} 返回包含 token 和用户信息的响应
   *
   * @example
   * const result = await authApi.login({ username: 'admin', password: '123456' })
   */
  login(data) {
    return request({
      url: '/auth/login',
      method: 'POST',
      data
    })
  },

  /**
   * 用户登出
   *
   * 清除服务端的 session 或 token，使当前用户登出
   *
   * @returns {Promise} 返回登出结果
   */
  logout() {
    return request({
      url: '/auth/logout',
      method: 'POST'
    })
  },

  /**
   * 获取当前登录用户信息
   *
   * 用于验证用户身份和权限，通常在应用启动时调用
   *
   * @returns {Promise} 返回用户详细信息，包括角色和权限
   */
  getCurrentUser() {
    return request({
      url: '/auth/me',
      method: 'GET'
    })
  },

  /**
   * 修改当前用户密码
   *
   * @param {Object} data - 密码修改信息
   * @param {string} data.oldPassword - 原密码
   * @param {string} data.newPassword - 新密码
   * @returns {Promise} 返回修改结果
   */
  changePassword(data) {
    return request({
      url: '/auth/change-password',
      method: 'POST',
      data
    })
  },

  /**
   * 上传头像
   *
   * @param {FormData} formData - 包含头像文件的表单数据
   * @returns {Promise} 返回上传结果
   */
  uploadAvatar(formData) {
    return request({
      url: '/auth/avatar',
      method: 'POST',
      data: formData,
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  }
}

// ==================== 项目管理 API ====================

/**
 * 项目管理 API 接口
 *
 * 提供项目的增删改查功能，支持分页、搜索和筛选
 *
 * @namespace projectApi
 */
export const projectApi = {
  /**
   * 获取项目列表
   *
   * 支持分页、搜索和状态筛选
   *
   * @param {Object} params - 查询参数
   * @param {number} params.page - 页码（从1开始）
   * @param {number} params.per_page - 每页数量（最大100）
   * @param {string} params.search - 搜索关键词（项目名称、编号）
   * @param {string} params.status - 状态筛选（planning/active/suspended/completed/cancelled）
   * @returns {Promise} 返回项目列表和分页信息
   */
  getList(params) {
    return request({
      url: '/project/projects',
      method: 'GET',
      params
    })
  },

  /**
   * 获取项目详情
   *
   * @param {number} id - 项目ID
   * @returns {Promise} 返回项目详细信息
   */
  getDetail(id) {
    return request({
      url: `/project/projects/${id}`,
      method: 'GET'
    })
  },

  /**
   * 创建新项目
   *
   * @param {Object} data - 项目信息
   * @param {string} data.name - 项目名称（必填）
   * @param {string} data.code - 项目编号（可选，为空则自动生成）
   * @param {string} data.manager - 项目负责人
   * @param {string} data.contact - 联系电话
   * @param {string} data.start_date - 开始日期（YYYY-MM-DD）
   * @param {string} data.end_date - 结束日期（YYYY-MM-DD）
   * @param {number} data.budget - 预算金额
   * @param {string} data.status - 状态（planning/active/suspended/completed/cancelled）
   * @param {string} data.location - 项目地址
   * @param {string} data.description - 项目描述
   * @returns {Promise} 返回创建的项目信息
   */
  create(data) {
    return request({
      url: '/project/projects',
      method: 'POST',
      data
    })
  },

  /**
   * 更新项目信息
   *
   * @param {number} id - 项目ID
   * @param {Object} data - 要更新的项目信息
   * @returns {Promise} 返回更新后的项目信息
   */
  update(id, data) {
    return request({
      url: `/project/projects/${id}`,
      method: 'PUT',
      data
    })
  },

  /**
   * 删除项目
   *
   * 注意：删除项目前请确保没有关联的物资、入库单等数据
   *
   * @param {number} id - 项目ID
   * @returns {Promise} 返回删除结果
   */
  delete(id) {
    return request({
      url: `/project/projects/${id}`,
      method: 'DELETE'
    })
  },

  /**
   * 获取项目成员列表
   *
   * @param {number} id - 项目ID
   * @returns {Promise} 返回项目成员列表
   */
  getMembers(id) {
    return request({
      url: `/project/projects/${id}/members`,
      method: 'GET'
    })
  },

  /**
   * 添加项目成员（批量替换）
   *
   * 注意：这会清除现有成员并添加新成员
   *
   * @param {number} id - 项目ID
   * @param {Object} data - 成员信息
   * @param {number[]} data.user_ids - 用户ID数组
   * @returns {Promise} 返回添加结果
   */
  addMember(id, data) {
    return request({
      url: `/project/projects/${id}/members`,
      method: 'POST',
      data
    })
  },

  /**
   * 删除项目成员
   *
   * @param {number} id - 项目ID
   * @param {number} userId - 用户ID
   * @returns {Promise} 返回删除结果
   */
  removeMember(id, userId) {
    return request({
      url: `/project/projects/${id}/members/${userId}`,
      method: 'DELETE'
    })
  },

  /**
   * 获取项目树
   *
   * 获取项目的完整层级结构（包含所有子项目）
   *
   * @param {number} id - 项目ID
   * @returns {Promise} 返回项目树结构
   */
  getProjectTree(id) {
    return request({
      url: `/project/projects/${id}/tree`,
      method: 'GET'
    })
  },

  /**
   * 获取子项目列表
   *
   * 获取项目的直接子项目列表
   *
   * @param {number} id - 项目ID
   * @returns {Promise} 返回子项目列表
   */
  getChildren(id) {
    return request({
      url: `/project/projects/${id}/children`,
      method: 'GET'
    })
  },

  /**
   * 聚合子项目进度
   *
   * 从所有子项目聚合进度到主项目
   *
   * @param {number} id - 项目ID
   * @returns {Promise} 返回聚合后的进度
   */
  aggregateProgress(id) {
    return request({
      url: `/project/projects/${id}/aggregate-progress`,
      method: 'POST'
    })
  }
}

// ==================== 物资管理 API ====================

/**
 * 物资管理 API 接口
 *
 * 提供物资的增删改查、导入导出功能
 *
 * @namespace materialApi
 */
export const materialApi = {
  // ========== 物资主数据API ==========

  /**
   * 获取物资主数据列表
   *
   * @param {Object} params - 查询参数
   * @returns {Promise} 返回物资主数据列表
   */
  getMasterList(params) {
    return request({
      url: '/materials/master',
      method: 'GET',
      params
    })
  },

  /**
   * 获取物资主数据详情
   *
   * @param {number} id - 物资主数据ID
   * @returns {Promise} 返回物资主数据详情
   */
  getMasterDetail(id) {
    return request({
      url: `/materials/master/${id}`,
      method: 'GET'
    })
  },

  /**
   * 创建物资主数据
   *
   * @param {Object} data - 物资主数据信息
   * @returns {Promise} 返回创建的物资主数据
   */
  createMaster(data) {
    return request({
      url: '/materials/master',
      method: 'POST',
      data
    })
  },

  /**
   * 更新物资主数据
   *
   * @param {number} id - 物资主数据ID
   * @param {Object} data - 更新的物资主数据信息
   * @returns {Promise} 返回更新后的物资主数据
   */
  updateMaster(id, data) {
    return request({
      url: `/materials/master/${id}`,
      method: 'PUT',
      data
    })
  },

  /**
   * 删除物资主数据
   *
   * @param {number} id - 物资主数据ID
   * @returns {Promise} 返回删除结果
   */
  deleteMaster(id) {
    return request({
      url: `/materials/master/${id}`,
      method: 'DELETE'
    })
  },

  /**
   * 获取项目物资列表（带库存）
   *
   * @param {Object} params - 查询参数
   * @returns {Promise} 返回项目物资列表
   */
  getProjectMaterials(params) {
    return request({
      url: '/materials/master/project',
      method: 'GET',
      params
    })
  },

  // ========== 物资管理API ==========

  /**
   * 获取物资列表
   *
   * 支持分页、搜索和项目筛选
   *
   * @param {Object} params - 查询参数
   * @param {number} params.page - 页码
   * @param {number} params.per_page - 每页数量
   * @param {string} params.search - 搜索关键词（物资名称、编码、规格）
   * @param {number} params.project_id - 项目ID筛选
   * @returns {Promise} 返回物资列表和分页信息
   */
  getList(params) {
    return request({
      url: '/material/materials',
      method: 'GET',
      params
    })
  },

  /**
   * 获取物资详情
   *
   * @param {number} id - 物资ID
   * @returns {Promise} 返回物资详细信息
   */
  getDetail(id) {
    return request({
      url: `/material/materials/${id}`,
      method: 'GET'
    })
  },

  /**
   * 创建新物资
   *
   * @param {Object} data - 物资信息
   * @param {string} data.name - 物资名称（必填）
   * @param {string} data.code - 物资编码（可选，为空则自动生成）
   * @param {string} data.specification - 规格型号
   * @param {string} data.unit - 单位（个、台、米、kg等）
   * @param {number} data.price - 单价
   * @param {number} data.quantity - 数量
   * @param {string} data.category - 分类（建材、五金、电器等）
   * @param {number} data.project_id - 关联项目ID
   * @returns {Promise} 返回创建的物资信息
   */
  create(data) {
    return request({
      url: '/material/materials',
      method: 'POST',
      data
    })
  },

  /**
   * 更新物资信息
   *
   * @param {number} id - 物资ID
   * @param {Object} data - 要更新的物资信息
   * @returns {Promise} 返回更新后的物资信息
   */
  update(id, data) {
    return request({
      url: `/material/materials/${id}`,
      method: 'PUT',
      data
    })
  },

  /**
   * 删除物资
   *
   * @param {number} id - 物资ID
   * @returns {Promise} 返回删除结果
   */
  delete(id) {
    return request({
      url: `/material/materials/${id}`,
      method: 'DELETE'
    })
  },

  /**
   * 导入物资数据
   *
   * 支持 Excel 文件导入，批量创建物资记录
   *
   * @param {FormData} data - 表单数据，包含 file 字段
   * @param {File} data.file - Excel 文件（.xlsx 或 .xls）
   * @param {number} data.project_id - 关联的项目ID
   * @returns {Promise} 返回导入结果（包含成功和失败记录数）
   */
  import(data) {
    return request({
      url: '/material/material/materials/import',
      method: 'POST',
      data,
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  /**
   * 批量创建物资
   *
   * @param {Object} data - 导入数据
   * @returns {Promise} 返回导入结果
   */
  batchCreate(data) {
    return request({
      url: '/material/materials/batch-create',
      method: 'POST',
      data
    })
  },

  /**
   * 导出物资数据
   *
   * 导出为 Excel 文件
   *
   * @param {Object} params - 查询参数（与 getList 相同）
   * @returns {Promise} 返回 Excel 文件 Blob
   */
  export(params) {
    return request({
      url: '/material/materials/export',
      method: 'GET',
      params,
      responseType: 'blob'
    })
  },

  /**
   * 获取未入库物资列表
   *
   * @param {Object} params - 查询参数
   * @returns {Promise} 返回未入库物资列表
   */
  getUnstored(params) {
    return request({
      url: '/material/materials/unstored',
      method: 'GET',
      params
    })
  },

  /**
   * 导出未入库物资
   *
   * @param {Object} params - 查询参数
   * @returns {Promise} 返回 Excel 文件 Blob
   */
  exportUnstored(params) {
    return request({
      url: '/material/materials/unstored/export',
      method: 'GET',
      params,
      responseType: 'blob'
    })
  },

  /**
   * 获取物资日志
   *
   * @param {number} id - 物资ID
   * @returns {Promise} 返回物资日志列表
   */
  getLogs(id) {
    return request({
      url: `/materials/${id}/logs`,
      method: 'GET'
    })
  },

  // ==================== 物资分类 API ====================

  /**
   * 获取所有物资分类
   *
   * @returns {Promise} 返回分类列表
   */
  getCategories() {
    return request({
      url: '/material/categories',
      method: 'GET'
    })
  },

  /**
   * 获取单个分类详情
   *
   * @param {number} id - 分类ID
   * @returns {Promise} 返回分类详细信息
   */
  getCategory(id) {
    return request({
      url: `/material/categories/${id}`,
      method: 'GET'
    })
  },

  /**
   * 创建物资分类
   *
   * @param {Object} data - 分类信息
   * @param {string} data.name - 分类名称（必填）
   * @param {string} data.code - 分类编码
   * @param {number} data.sort - 排序
   * @param {string} data.remark - 备注
   * @returns {Promise} 返回创建的分类信息
   */
  createCategory(data) {
    return request({
      url: '/material/categories',
      method: 'POST',
      data
    })
  },

  /**
   * 更新物资分类
   *
   * @param {number} id - 分类ID
   * @param {Object} data - 要更新的分类信息
   * @returns {Promise} 返回更新后的分类信息
   */
  updateCategory(id, data) {
    return request({
      url: `/material/categories/${id}`,
      method: 'PUT',
      data
    })
  },

  /**
   * 删除物资分类
   *
   * @param {number} id - 分类ID
   * @returns {Promise} 返回删除结果
   */
  deleteCategory(id) {
    return request({
      url: `/material/categories/${id}`,
      method: 'DELETE'
    })
  },

  /**
   * 批量更新分类排序
   *
   * @param {Array} sorts - 排序数组 [{id: 1, sort: 1}, {id: 2, sort: 2}]
   * @returns {Promise} 返回更新结果
   */
  updateCategorySort(sorts) {
    return request({
      url: '/material/categories/sort',
      method: 'POST',
      data: { sorts }
    })
  },

  /**
   * 批量创建物资（用于计划导入自动创建物资）
   *
   * 当计划项没有关联物资库时，自动创建物资记录
   *
   * @param {Array} materials - 物资数据数组
   * @returns {Promise} 返回创建的物资列表
   */
  batchCreateMaterials(materials) {
    return request({
      url: '/material/materials/batch-create',
      method: 'POST',
      data: { materials }
    })
  }
}

// ==================== 库存管理 API ====================

/**
 * 库存管理 API 接口
 *
 * 提供库存查询、出入库操作、库存日志等功能
 *
 * @namespace stockApi
 */
export const stockApi = {
  /**
   * 获取库存列表
   *
   * 支持分页、搜索、项目筛选和状态筛选
   *
   * @param {Object} params - 查询参数
   * @param {number} params.page - 页码
   * @param {number} params.per_page - 每页数量
   * @param {string} params.search - 搜索关键词
   * @param {number} params.project_id - 项目ID筛选
   * @param {string} params.status - 状态筛选（normal/low/shortage）
   * @returns {Promise} 返回库存列表和分页信息
   */
  getList(params) {
    return request({
      url: '/stock/stocks',
      method: 'GET',
      params
    })
  },

  /**
   * 获取库存预警列表
   *
   * @returns {Promise} 返回库存预警列表
   */
  getAlerts() {
    return request({
      url: '/stock/stocks/alerts',
      method: 'GET'
    })
  },

  /**
   * 获取库存详情
   *
   * @param {number} id - 库存ID
   * @returns {Promise} 返回库存详情
   */
  getDetail(id) {
    return request({
      url: `/stock/stocks/${id}`,
      method: 'GET'
    })
  },

  /**
   * 创建库存记录
   *
   * @param {Object} data - 库存信息
   * @returns {Promise} 返回创建的库存
   */
  create(data) {
    return request({
      url: '/stock/stocks',
      method: 'POST',
      data
    })
  },

  /**
   * 更新库存
   *
   * @param {number} id - 库存ID
   * @param {Object} data - 更新的库存信息
   * @returns {Promise} 返回更新后的库存
   */
  update(id, data) {
    return request({
      url: `/stock/stocks/${id}`,
      method: 'PUT',
      data
    })
  },

  /**
   * 删除库存
   *
   * @param {number} id - 库存ID
   * @returns {Promise} 返回删除结果
   */
  delete(id) {
    return request({
      url: `/stock/stocks/${id}`,
      method: 'DELETE'
    })
  },

  /**
   * 获取库存变动日志
   *
   * 用于追溯库存的出入库历史记录
   *
   * @param {Object} params - 查询参数
   * @param {number} params.stock_id - 库存ID筛选
   * @param {number} params.page - 页码
   * @param {number} params.per_page - 每页数量
   * @returns {Promise} 返回库存日志列表
   */
  getLogs(params) {
    return request({
      url: '/stock/stock-logs',
      method: 'GET',
      params
    })
  },

  /**
   * 获取库存操作日志
   *
   * @param {number} id - 库存ID
   * @returns {Promise} 返回库存日志列表
   */
  getStockLogs(id) {
    return request({
      url: `/stocks/${id}/logs`,
      method: 'GET'
    })
  },

  /**
   * 删除库存日志
   *
   * @param {number} id - 日志ID
   * @returns {Promise} 返回删除结果
   */
  deleteLog(id) {
    return request({
      url: `/stock/stock-logs/${id}`,
      method: 'DELETE'
    })
  },

  /**
   * 物资入库
   *
   * 增加指定物资的库存数量
   *
   * @param {Object} data - 入库信息
   * @param {number} data.id - 库存ID
   * @param {number} data.quantity - 入库数量
   * @param {string} data.remark - 备注信息
   * @returns {Promise} 返回入库后的库存信息
   */
  in(data) {
    return request({
      url: `/stock/stocks/${data.id}/in`,
    method: 'POST',
      data
    })
  },

  /**
   * 物资入库
   *
   * 增加指定物资的库存数量
   *
   * @param {number} id - 库存ID
   * @param {Object} data - 入库信息
   * @param {number} data.quantity - 入库数量（必填）
   * @param {string} data.remark - 备注信息（可选）
   * @returns {Promise} 返回入库后的库存信息
   */
  in(id, data) {
    return request({
      url: `/stock/stocks/${id}/in`,
      method: 'POST',
      data
    })
  },

  /**
   * 物资出库
   *
   * 减少指定物资的库存数量
   *
   * @param {number} id - 库存ID
   * @param {Object} data - 出库信息
   * @param {number} data.quantity - 出库数量（必填）
   * @param {string} data.remark - 备注信息（可选）
   * @returns {Promise} 返回出库后的库存信息
   */
  out(id, data) {
    return request({
      url: `/stock/stocks/${id}/out`,
    method: 'POST',
      data
    })
  },

  /**
   * 库存调整
   *
   * @param {number} id - 库存ID
   * @param {Object} data - 调整信息
   * @returns {Promise} 返回调整结果
   */
  adjust(id, data) {
    return request({
      url: `/stocks/${id}/adjust`,
      method: 'POST',
      data
    })
  },

  /**
   * 导出库存数据
   *
   * 导出为 Excel 文件
   *
   * @param {Object} params - 查询参数（与 getList 相同）
   * @returns {Promise} 返回 Excel 文件 Blob
   */
  export(params) {
    return request({
      url: '/stock/stocks/export',
      method: 'GET',
      params,
      responseType: 'blob'
    })
  }
}

/**
 * 出库单管理 API
 *
 * 出库单（领料单）用于记录和审批材料出库流程
 * 支持完整的工作流审批机制
 *
 * @namespace requisitionApi
 */
export const requisitionApi = {
  /**
   * 获取待处理出库单数量
   *
   * @returns {Promise} 返回待处理数量
   */
  getPendingCount() {
    return request({
      url: '/requisition/requisitions/pending/count',
      method: 'GET'
    })
  },

  /**
   * 获取出库单列表
   *
   * 支持分页、搜索、排序和状态过滤
   *
   * @param {Object} params - 查询参数
   * @param {number} params.page - 页码（从 1 开始）
   * @param {number} params.page_size - 每页数量
   * @param {string} params.status - 状态过滤（pending/approved/rejected/issued）
   * @param {string} params.search - 搜索关键字（单号、项目名称）
   * @returns {Promise} 返回出库单列表和分页信息
   */
  getList(params) {
    return request({
      url: '/requisition/requisitions',
      method: 'GET',
      params
    })
  },

  /**
   * 获取出库单详情
   *
   * 包含完整的出库单信息和出库物资明细
   *
   * @param {number} id - 出库单ID
   * @returns {Promise} 返回出库单详情
   */
  getDetail(id) {
    return request({
      url: `/requisition/requisitions/${id}`,
      method: 'GET'
    })
  },

  /**
   * 创建出库单
   *
   * 创建后状态为 pending，需要经过审批流程
   *
   * @param {Object} data - 出库单信息
   * @param {number} data.project_id - 项目ID
   * @param {Array} data.items - 出库物资明细
   * @param {string} data.remark - 备注
   * @returns {Promise} 返回创建的出库单信息
   */
  create(data) {
    return request({
      url: '/requisition/requisitions',
      method: 'POST',
      data
    })
  },

  /**
   * 更新出库单
   *
   * 只能更新 pending 状态的出库单
   *
   * @param {number} id - 出库单ID
   * @param {Object} data - 更新的出库单信息
   * @returns {Promise} 返回更新后的出库单信息
   */
  update(id, data) {
    return request({
      url: `/requisition/requisitions/${id}`,
      method: 'PUT',
      data
    })
  },

  /**
   * 删除出库单
   *
   * 只能删除 pending 状态的出库单
   *
   * @param {number} id - 出库单ID
   * @returns {Promise} 返回删除结果
   */
  delete(id) {
    return request({
      url: `/requisition/requisitions/${id}`,
      method: 'DELETE'
    })
  },

  /**
   * 审核通过出库单
   *
   * 审批通过后，单据状态变为 approved
   * 可以继续执行发货操作
   *
   * @param {number} id - 出库单ID
   * @param {Object} data - 审批信息
   * @param {string} data.approval_comment - 审批意见
   * @returns {Promise} 返回审批后的出库单信息
   */
  approve(id, data) {
    return request({
      url: `/requisition/requisitions/${id}/approve`,
      method: 'POST',
      data
    })
  },

  /**
   * 审批拒绝出库单
   *
   * 拒绝后，单据状态变为 rejected
   * 可以修改后重新提交
   *
   * @param {number} id - 出库单ID
   * @param {Object} data - 拒绝信息
   * @param {string} data.rejection_reason - 拒绝原因
   * @returns {Promise} 返回审批后的出库单信息
   */
  reject(id, data) {
    return request({
      url: `/requisition/requisitions/${id}/reject`,
    method: 'POST',
      data
    })
  },

  /**
   * 发货（实际出库）
   *
   * 执行实际的库存扣减操作
   * 发货后，单据状态变为 issued
   *
   * @param {number} id - 出库单ID
   * @param {Object} data - 发货信息
   * @param {string} data.issuer - 发货人
   * @param {string} data.receiver - 收货人
   * @param {string} data.issue_remark - 发货备注
   * @returns {Promise} 返回发货后的出库单信息
   */
  issue(id, data) {
    return request({
      url: `/requisition/requisitions/${id}/issue`,
      method: 'POST',
      data
    })
  },

  /**
   * 获取领料单明细列表
   *
   * @param {Object} params - 查询参数
   * @returns {Promise} 返回领料单明细列表
   */
  getItems(params) {
    return request({
      url: '/requisition/requisition-items',
      method: 'GET',
      params
    })
  },

  /**
   * 获取工作流审批历史
   *
   * 查看出库单的完整审批记录
   *
   * @param {number} id - 出库单ID（关联工作流实例）
   * @returns {Promise} 返回审批历史记录列表
   */
  getWorkflowHistory(id) {
    return request({
      url: `/workflow-instances/${id}/approvals`,
      method: 'GET'
    })
  }
}

/**
 * 入库单管理 API
 *
 * 入库单用于记录和审批材料入库流程
 * 审批通过后会自动增加库存数量
 *
 * @namespace inboundApi
 */
export const inboundApi = {
  /**
   * 获取待处理入库单数量
   *
   * @returns {Promise} 返回待处理数量
   */
  getPendingCount() {
    return request({
      url: '/inbound/inbound-orders/pending/count',
      method: 'GET'
    })
  },

  /**
   * 获取入库单列表
   *
   * 支持分页、搜索、排序和状态过滤
   *
   * @param {Object} params - 查询参数
   * @param {number} params.page - 页码（从 1 开始）
   * @param {number} params.page_size - 每页数量
   * @param {string} params.status - 状态过滤（pending/approved/rejected）
   * @param {string} params.search - 搜索关键字（单号、供应商名称）
   * @returns {Promise} 返回入库单列表和分页信息
   */
  getList(params) {
    return request({
      url: '/inbound/inbound-orders',
      method: 'GET',
      params
    })
  },

  /**
   * 获取入库单模板
   *
   * @returns {Promise} 返回入库单模板
   */
  getTemplate() {
    return request({
      url: '/inbound/inbound/template',
      method: 'GET'
    })
  },

  /**
   * 提交入库单
   *
   * @param {Object} data - 入库单信息
   * @returns {Promise} 返回提交结果
   */
  submit(data) {
    return request({
      url: '/inbound/inbound/submit',
      method: 'POST',
      data
    })
  },

  /**
   * 导入入库单
   *
   * @param {FormData} data - 包含文件的FormData
   * @returns {Promise} 返回导入结果
   */
  import(data) {
    return request({
      url: '/inbound/inbound/import',
      method: 'POST',
      data,
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  /**
   * 获取入库单详情
   *
   * 包含完整的入库单信息和入库物资明细
   *
   * @param {number} id - 入库单ID
   * @returns {Promise} 返回入库单详情
   */
  getDetail(id) {
    return request({
      url: `/inbound/inbound-orders/${id}`,
      method: 'GET'
    })
  },

  /**
   * 创建入库单
   *
   * 创建后状态为 pending，需要经过审批流程
   * 审批通过后会自动增加库存
   *
   * @param {Object} data - 入库单信息
   * @param {string} data.supplier - 供应商名称
   * @param {string} data.inbound_date - 入库日期
   * @param {Array} data.items - 入库物资明细
   * @param {string} data.remark - 备注
   * @returns {Promise} 返回创建的入库单信息
   */
  create(data) {
    return request({
      url: '/inbound/inbound-orders',
      method: 'POST',
      data
    })
  },

  /**
   * 更新入库单
   *
   * 只能更新 pending 状态的入库单
   *
   * @param {number} id - 入库单ID
   * @param {Object} data - 更新的入库单信息
   * @returns {Promise} 返回更新后的入库单信息
   */
  update(id, data) {
    return request({
      url: `/inbound/inbound-orders/${id}`,
      method: 'PUT',
      data
    })
  },

  /**
   * 删除入库单
   *
   * 只能删除 pending 状态的入库单
   *
   * @param {number} id - 入库单ID
   * @returns {Promise} 返回删除结果
   */
  delete(id) {
    return request({
      url: `/inbound/inbound-orders/${id}`,
      method: 'DELETE'
    })
  },

  /**
   * 审批通过入库单
   *
   * 审批通过后会：
   * 1. 更新入库单状态为 approved
   * 2. 自动增加对应物资的库存数量
   * 3. 记录库存变动日志
   *
   * @param {number} id - 入库单ID
   * @param {Object} data - 审批信息
   * @param {string} data.approval_comment - 审批意见
   * @returns {Promise} 返回审批后的入库单信息
   */
  approve(id, data) {
    return request({
      url: `/inbound/inbound-orders/${id}/approve`,
    method: 'POST',
      data
    })
  },

  /**
   * 审批拒绝入库单
   *
   * 拒绝后，单据状态变为 rejected
   * 不会增加库存
   *
   * @param {number} id - 入库单ID
   * @param {Object} data - 拒绝信息
   * @param {string} data.rejection_reason - 拒绝原因
   * @returns {Promise} 返回审批后的入库单信息
   */
  reject(id, data) {
    return request({
      url: `/inbound/inbound-orders/${id}/reject`,
      method: 'POST',
      data
    })
  },

  /**
   * 获取工作流审批历史
   *
   * 查看入库单的完整审批记录
   *
   * @param {number} id - 入库单ID（关联工作流实例）
   * @returns {Promise} 返回审批历史记录列表
   */
  getWorkflowHistory(id) {
    return request({
      url: `/inbound/inbound-orders/${id}/workflow-history`,
      method: 'GET'
    })
  },

  /**
   * 获取可执行的操作
   *
   * 根据当前用户权限和单据状态
   * 返回可以执行的操作列表
   *
   * @param {number} id - 入库单ID
   * @returns {Promise} 返回可执行的操作列表
   */
  getAvailableActions(id) {
    return request({
      url: `/inbound/inbound-orders/${id}/available-actions`,
      method: 'GET'
    })
  },

  /**
   * 批量审批入库单
   *
   * 一次审批多个入库单
   * 适用于快速处理多个待审批单据
   *
   * @param {number[]} ids - 入库单ID数组
   * @param {Object} data - 审批信息
   * @param {string} data.approval_comment - 审批意见
   * @returns {Promise} 返回批量审批结果
   */
  batchApprove(ids, data) {
    return request({
      url: '/inbound/inbound-orders/batch-approve',
      method: 'POST',
      data: { ids, ...data }
    })
  }
}

/**
 * 系统管理 API
 *
 * 提供系统级别的管理功能，包括：
 * - 系统配置管理
 * - 系统日志查询
 * - 数据备份与恢复
 * - 系统统计信息
 *
 * @namespace systemApi
 */
export const systemApi = {
  /**
   * 获取公开系统设置（不需要认证）
   *
   * 获取系统名称等公开配置信息
   *
   * @returns {Promise} 返回系统配置信息
   */
  getPublicSettings() {
    return request({
      url: '/system/public/settings',
      method: 'GET'
    })
  },

  /**
   * 获取系统设置
   *
   * 获取当前系统的配置参数
   *
   * @returns {Promise} 返回系统配置信息
   */
  getSettings() {
    return request({
      url: '/system/settings',
      method: 'GET'
    })
  },

  /**
   * 保存系统设置
   *
   * 更新系统配置参数
   *
   * @param {Object} data - 系统配置信息
   * @returns {Promise} 返回更新后的配置
   */
  saveSettings(data) {
    return request({
      url: '/system/settings',
      method: 'PUT',
      data
    })
  },

  /**
   * 创建数据备份
   *
   * 创建数据库的完整备份
   *
   * @param {Object} data - 备份配置
   * @param {string} data.name - 备份名称
   * @param {string} data.description - 备份描述
   * @returns {Promise} 返回备份信息
   */
  createBackup(data) {
    return request({
      url: '/system/backup',
      method: 'POST',
      data
    })
  },

  /**
   * 获取备份列表
   *
   * 查询所有可用的数据备份
   *
   * @returns {Promise} 返回备份列表
   */
  getBackups() {
    return request({
      url: '/system/backup',
      method: 'GET'
    })
  },

  /**
   * 删除备份
   *
   * 删除指定的数据备份文件
   *
   * @param {number} id - 备份ID
   * @returns {Promise} 返回删除结果
   */
  deleteBackup(id) {
    return request({
      url: `/system/backup/${id}`,
      method: 'DELETE'
    })
  },

  /**
   * 恢复备份
   *
   * 从备份文件恢复数据库
   * 注意：此操作会覆盖当前数据
   *
   * @param {string} filename - 备份文件名
   * @returns {Promise} 返回恢复结果
   */
  restoreBackup(filename) {
    return request({
      url: `/system/backup/restore`,
      method: 'POST',
      data: { backup_name: filename }
    })
  },

  /**
   * 获取统计数据
   *
   * 获取系统运行统计数据
   * 包括用户数量、项目数量、物资数量等
   *
   * @returns {Promise} 返回统计信息
   */
  getStats() {
    return request({
      url: '/system/stats',
      method: 'GET'
    })
  },

  /**
   * 获取系统信息
   *
   * 获取系统版本和运行环境信息
   *
   * @returns {Promise} 返回系统信息
   */
  getSystemInfo() {
    return request({
      url: '/system/info',
      method: 'GET'
    })
  },

  /**
   * 获取最近活动
   *
   * @returns {Promise} 返回最近活动列表
   */
  getRecentActivities() {
    return request({
      url: '/system/recent-activities',
      method: 'GET'
    })
  },

  /**
   * 获取物资分类统计
   *
   * @returns {Promise} 返回物资分类统计
   */
  getMaterialCategoryStats() {
    return request({
      url: '/system/material-category-stats',
      method: 'GET'
    })
  },

  /**
   * 获取项目物资统计
   *
   * @returns {Promise} 返回项目物资统计
   */
  getProjectMaterialStats() {
    return request({
      url: '/system/project-material-stats',
      method: 'GET'
    })
  },

  // ========== 数据备份扩展 ==========

  /**
   * 获取备份历史
   *
   * @returns {Promise} 返回备份历史
   */
  getBackupHistory() {
    return request({
      url: '/system/backup/history',
      method: 'GET'
    })
  },

  /**
   * 下载备份
   *
   * @param {string} backupName - 备份名称
   * @returns {Promise} 返回备份文件
   */
  downloadBackup(backupName) {
    return request({
      url: `/system/backup/${backupName}/download`,
      method: 'GET',
      responseType: 'blob'
    })
  },

  // ========== 报表管理 ==========

  /**
   * 获取仪表板数据
   *
   * @returns {Promise} 返回仪表板数据
   */
  getDashboard() {
    return request({
      url: '/system/reports/dashboard',
      method: 'GET'
    })
  },

  /**
   * 获取报表列表
   *
   * @returns {Promise} 返回报表列表
   */
  getReports() {
    return request({
      url: '/system/reports',
      method: 'GET'
    })
  },

  /**
   * 生成报表
   *
   * @param {Object} data - 报表参数
   * @returns {Promise} 返回生成结果
   */
  generateReport(data) {
    return request({
      url: '/system/reports/generate',
      method: 'POST',
      data
    })
  },

  /**
   * 下载报表
   *
   * @param {Object} params - 查询参数
   * @returns {Promise} 返回报表文件
   */
  downloadReport(params) {
    return request({
      url: '/system/reports/download',
      method: 'GET',
      params,
      responseType: 'blob'
    })
  },

  /**
   * 下载指定报表
   *
   * @param {string} reportName - 报表名称
   * @returns {Promise} 返回报表文件
   */
  downloadReportByName(reportName) {
    return request({
      url: `/system/reports/${reportName}/download`,
      method: 'GET',
      responseType: 'blob'
    })
  },

  /**
   * 删除报表
   *
   * @param {Object} data - 删除参数
   * @returns {Promise} 返回删除结果
   */
  deleteReport(data) {
    return request({
      url: '/system/reports/delete',
      method: 'POST',
      data
    })
  }
}

/**
 * 施工日志管理 API
 *
 * 记录项目施工过程中的日常日志
 * 包括施工内容、进度、质量、安全等信息
 *
 * @namespace constructionLogApi
 */
export const constructionLogApi = {
  /**
   * 获取施工日志列表
   *
   * 支持按项目、日期范围等条件过滤
   *
   * @param {Object} params - 查询参数
   * @param {number} params.page - 页码
   * @param {number} params.page_size - 每页数量
   * @param {number} params.project_id - 项目ID
   * @param {string} params.start_date - 开始日期
   * @param {string} params.end_date - 结束日期
   * @returns {Promise} 返回施工日志列表
   */
  getList(params) {
    return request({
      url: '/construction_log/logs',
      method: 'GET',
      params
    })
  },

  /**
   * 获取施工日志详情
   *
   * 获取单条施工日志的完整信息
   *
   * @param {number} id - 日志ID
   * @returns {Promise} 返回施工日志详情
   */
  getDetail(id) {
    return request({
      url: `/construction_log/${id}`,
      method: 'GET'
    })
  },

  /**
   * 创建施工日志
   *
   * 记录当天的施工情况
   *
   * @param {Object} data - 施工日志信息
   * @param {number} data.project_id - 项目ID
   * @param {string} data.log_date - 日志日期
   * @param {string} data.weather - 天气情况
   * @param {string} data.content - 施工内容
   * @param {string} data.progress - 进度说明
   * @param {string} data.quality_issues - 质量问题
   * @param {string} data.safety_issues - 安全问题
   * @param {Array} data.images - 图片附件
   * @returns {Promise} 返回创建的日志信息
   */
  create(data) {
    return request({
      url: '/construction_log/logs',
      method: 'POST',
      data
    })
  },

  /**
   * 更新施工日志
   *
   * 修改已提交的施工日志
   *
   * @param {number} id - 日志ID
   * @param {Object} data - 更新的日志信息
   * @returns {Promise} 返回更新后的日志信息
   */
  update(id, data) {
    return request({
      url: `/construction_log/${id}`,
      method: 'PUT',
      data
    })
  },

  /**
   * 删除施工日志
   *
   * 删除指定的施工日志记录
   *
   * @param {number} id - 日志ID
   * @returns {Promise} 返回删除结果
   */
  delete(id) {
    return request({
      url: `/construction_log/${id}`,
      method: 'DELETE'
    })
  },

  /**
   * 上传日志图片
   *
   * @param {FormData} formData - 包含图片的FormData
   * @returns {Promise} 返回上传结果
   */
  uploadImage(formData) {
    return request({
      url: '/construction_log/upload_image',
      method: 'POST',
      data: formData,
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  }
}

/**
 * 进度管理 API
 *
 * 提供项目进度管理、任务管理、依赖关系管理、AI生成等功能
 *
 * @namespace progressApi
 */
export const progressApi = {
  // ========== 项目进度API ==========

  /**
   * 获取项目进度计划
   *
   * 获取指定项目的完整进度计划和里程碑
   *
   * @param {number} projectId - 项目ID
   * @returns {Promise} 返回项目进度计划
   */
  getProjectSchedule(projectId) {
    return request({
      url: `/progress/project/${projectId}`,
      method: 'GET'
    })
  },

  /**
   * 更新项目进度计划
   *
   * 更新项目的进度信息和里程碑
   *
   * @param {number} projectId - 项目ID
   * @param {Object} data - 进度计划数据（ScheduleData格式）
   * @returns {Promise} 返回更新后的进度计划
   */
  updateProjectSchedule(projectId, data) {
    return request({
      url: `/progress/project/${projectId}`,
      method: 'PUT',
      data
    })
  },

  /**
   * 获取所有项目进度计划
   *
   * 获取所有项目的进度计划状态
   *
   * @returns {Promise} 返回所有项目的进度计划
   */
  getAllProjectSchedules() {
    return request({
      url: '/progress/project-schedules',
      method: 'GET'
    })
  },

  /**
   * 删除项目进度计划
   *
   * 删除项目的进度计划数据
   *
   * @param {number} projectId - 项目ID
   * @returns {Promise} 返回删除结果
   */
  deleteProjectSchedule(projectId) {
    return request({
      url: `/progress/project/${projectId}/schedule`,
      method: 'DELETE'
    })
  },

  /**
   * 检查项目进度计划是否存在
   *
   * 用于判断项目是否已创建进度计划
   *
   * @param {number} projectId - 项目ID
   * @returns {Promise} 返回存在状态
   */
  checkScheduleExists(projectId) {
    return request({
      url: `/progress/project/${projectId}/exists`,
      method: 'GET'
    })
  },

  // ========== 任务管理API ==========

  /**
   * 获取项目任务列表
   *
   * @param {number} projectId - 项目ID
   * @param {Object} params - 查询参数
   * @returns {Promise} 返回任务列表
   */
  getTasks(projectId, params) {
    return request({
      url: `/progress/project/${projectId}/tasks`,
      method: 'GET',
      params
    })
  },

  /**
   * 创建任务
   *
   * @param {number} projectId - 项目ID
   * @param {Object} data - 任务信息
   * @returns {Promise} 返回创建的任务
   */
  createTask(projectId, data) {
    return request({
      url: `/progress/project/${projectId}/tasks`,
      method: 'POST',
      data
    })
  },

  /**
   * 更新任务
   *
   * @param {number} taskId - 任务ID
   * @param {Object} data - 更新的任务信息
   * @returns {Promise} 返回更新后的任务
   */
  updateTask(taskId, data) {
    return request({
      url: `/progress/tasks/${taskId}`,
      method: 'PUT',
      data
    })
  },

  /**
   * 删除任务
   *
   * @param {number} taskId - 任务ID
   * @returns {Promise} 返回删除结果
   */
  deleteTask(taskId) {
    return request({
      url: `/progress/tasks/${taskId}`,
      method: 'DELETE'
    })
  },

  // ========== 任务依赖关系API ==========

  /**
   * 获取任务依赖
   *
   * @param {number} taskId - 任务ID
   * @returns {Promise} 返回依赖关系列表
   */
  getDependencies(taskId) {
    return request({
      url: `/progress/tasks/${taskId}/dependencies`,
      method: 'GET'
    })
  },

  /**
   * 添加任务依赖
   *
   * @param {number} taskId - 任务ID
   * @param {Object} data - 依赖信息
   * @param {number} data.depends_on - 依赖的任务ID
   * @param {string} data.type - 依赖类型（FS/FF/SS/SF）
   * @param {number} data.lag - 延迟天数
   * @returns {Promise} 返回创建的依赖关系
   */
  addDependency(taskId, data) {
    return request({
      url: `/progress/tasks/${taskId}/dependencies`,
    method: 'POST',
      data
    })
  },

  /**
   * 删除任务依赖
   *
   * @param {number} depId - 依赖关系ID
   * @returns {Promise} 返回删除结果
   */
  removeDependency(depId) {
    return request({
      url: `/progress/dependencies/${depId}`,
      method: 'DELETE'
    })
  },

  // ========== 位置持久化API ==========

  /**
   * 更新任务位置
   *
   * 用于网络图节点位置保存
   *
   * @param {number} taskId - 任务ID
   * @param {Object} position - 位置信息
   * @param {number} position.position_x - X坐标
   * @param {number} position.position_y - Y坐标
   * @returns {Promise} 返回更新结果
   */
  updateTaskPosition(taskId, position) {
    return request({
      url: `/progress/tasks/${taskId}/position`,
      method: 'PUT',
      data: position
    })
  },

  // ========== AI生成API ==========

  /**
   * AI生成进度计划
   *
   * 使用AI自动生成项目进度计划
   *
   * @param {number} projectId - 项目ID
   * @param {Object} options - 生成选项
   * @param {string} options.mode - 生成模式（auto/manual）
   * @param {string} options.project_type - 项目类型
   * @param {number} options.task_count - 任务数量
   * @param {string} options.requirements - 需求描述
   * @param {string} options.start_date - 开始日期
   * @param {string} options.end_date - 结束日期
   * @returns {Promise} 返回生成的进度计划
   */
  generatePlan(projectId, options) {
    return request({
      url: `/progress/project/${projectId}/generate-plan`,
      method: 'POST',
      data: options
    })
  },

  /**
   * 聚合子项目计划
   *
   * 将所有子项目的进度计划聚合到主项目
   *
   * @param {number} projectId - 主项目ID
   * @returns {Promise} 返回聚合后的进度计划
   */
  aggregateChildren(projectId) {
    return request({
      url: `/progress/project/${projectId}/aggregate-plan`,
      method: 'POST'
    })
  },

  // ========== 兼容旧API ==========

  /**
   * 获取进度列表（兼容旧版本）
   *
   * @deprecated 建议使用 getTasks
   * @param {Object} params - 查询参数
   * @returns {Promise} 返回进度列表
   */
  getList(params) {
    return request({
      url: '/progress',
      method: 'GET',
      params
    })
  },

  /**
   * 获取进度详情（兼容旧版本）
   *
   * @deprecated 建议使用 getTasks
   * @param {number} id - 进度ID
   * @returns {Promise} 返回进度详情
   */
  getDetail(id) {
    return request({
      url: `/progress/${id}`,
      method: 'GET'
    })
  },

  /**
   * 创建进度（兼容旧版本）
   *
   * @deprecated 建议使用 createTask
   * @param {Object} data - 进度信息
   * @returns {Promise} 返回创建的进度
   */
  create(data) {
    return request({
      url: `/progress/project/${data.project_id}/tasks`,
    method: 'POST',
      data
    })
  },

  /**
   * 更新进度（兼容旧版本）
   *
   * @deprecated 建议使用 updateTask
   * @param {number} id - 进度ID
   * @param {Object} data - 更新的进度信息
   * @returns {Promise} 返回更新后的进度
   */
  update(id, data) {
    return request({
      url: `/progress/tasks/${id}`,
      method: 'PUT',
      data
    })
  },

  /**
   * 删除进度（兼容旧版本）
   *
   * @deprecated 建议使用 deleteTask
   * @param {number} id - 进度ID
   * @returns {Promise} 返回删除结果
   */
  delete(id) {
    return request({
      url: `/progress/tasks/${id}`,
      method: 'DELETE'
    })
  },

  /**
   * 导出进度数据
   *
   * @param {Object} params - 查询参数
   * @returns {Promise} 返回Excel文件Blob
   */
  export(params) {
    return request({
      url: '/progress/export',
      method: 'GET',
      params,
      responseType: 'blob'
    })
  },

  // ========== 子任务进度管理API ==========

  /**
   * 计算父任务进度
   *
   * 根据子任务按工期加权平均计算父任务进度
   *
   * @param {number} taskId - 任务ID
   * @returns {Promise} 返回更新结果
   */
  calculateParentProgress(taskId) {
    return request({
      url: `/progress/tasks/${taskId}/calculate-parent-progress`,
      method: 'POST'
    })
  },

  /**
   * 更新任务及其所有父任务的进度
   *
   * @param {number} taskId - 任务ID
   * @returns {Promise} 返回更新结果
   */
  updateParentProgress(taskId) {
    return request({
      url: `/progress/tasks/${taskId}/update-parent-progress`,
      method: 'POST'
    })
  },

  // ========== 资源管理API ==========

  /**
   * 获取项目资源列表
   *
   * @param {number} projectId - 项目ID
   * @returns {Promise} 返回资源列表
   */
  getProjectResources(projectId) {
    return request({
      url: `/progress/project/${projectId}/resources`,
      method: 'GET'
    })
  },

  /**
   * 创建资源
   *
   * @param {number} projectId - 项目ID
   * @param {Object} data - 资源信息
   * @param {string} data.name - 资源名称
   * @param {string} data.type - 资源类型（labor/equipment/material）
   * @param {string} data.unit - 单位
   * @param {number} data.quantity - 可用数量
   * @param {number} data.cost_per_unit - 单位成本
   * @param {string} data.color - 显示颜色
   * @returns {Promise} 返回创建的资源
   */
  createResource(projectId, data) {
    return request({
      url: `/progress/project/${projectId}/resources`,
      method: 'POST',
      data
    })
  },

  /**
   * 更新资源
   *
   * @param {number} projectId - 项目ID
   * @param {number} resourceId - 资源ID
   * @param {Object} data - 更新的资源信息
   * @returns {Promise} 返回更新后的资源
   */
  updateResource(projectId, resourceId, data) {
    return request({
      url: `/progress/project/${projectId}/resources/${resourceId}`,
      method: 'PUT',
      data
    })
  },

  /**
   * 删除资源
   *
   * @param {number} projectId - 项目ID
   * @param {number} resourceId - 资源ID
   * @returns {Promise} 返回删除结果
   */
  deleteResource(projectId, resourceId) {
    return request({
      url: `/progress/project/${projectId}/resources/${resourceId}`,
      method: 'DELETE'
    })
  },

  // ========== 任务资源分配API ==========

  /**
   * 获取任务资源分配
   *
   * @param {number} taskId - 任务ID
   * @returns {Promise} 返回任务资源列表
   */
  getTaskResources(taskId) {
    return request({
      url: `/progress/tasks/${taskId}/resources`,
      method: 'GET'
    })
  },

  /**
   * 分配资源给任务
   *
   * @param {number} taskId - 任务ID
   * @param {Object} data - 分配信息
   * @param {number} data.resource_id - 资源ID
   * @param {number} data.quantity - 分配数量
   * @returns {Promise} 返回分配结果
   */
  allocateTaskResource(taskId, data) {
    return request({
      url: `/progress/tasks/${taskId}/resources`,
      method: 'POST',
      data
    })
  },

  /**
   * 移除任务资源
   *
   * @param {number} taskId - 任务ID
   * @param {number} resourceId - 资源ID
   * @returns {Promise} 返回移除结果
   */
  removeTaskResource(taskId, resourceId) {
    return request({
      url: `/progress/tasks/${taskId}/resources/${resourceId}`,
      method: 'DELETE'
    })
  },

  // ========== 可视化创建依赖关系 ==========

  /**
   * 可视化创建依赖关系
   *
   * @param {number} fromTaskId - 源任务ID（前置任务）
   * @param {number} toTaskId - 目标任务ID（后置任务）
   * @param {Object} data - 依赖信息
   * @param {string} data.type - 依赖类型（FS/FF/SS/SF）
   * @param {number} data.lag - 延迟天数
   * @returns {Promise} 返回创建的依赖关系
   */
  createDependencyVisual(fromTaskId, toTaskId, data = {}) {
    return request({
      url: `/progress/dependencies/visual/${fromTaskId}/${toTaskId}`,
      method: 'POST',
      data
    })
  },

  /**
   * 更新资源
   *
   * 更新项目资源的基本信息
   *
   * @param {number} resourceId - 资源ID
   * @param {Object} data - 资源信息
   * @param {string} data.name - 资源名称
   * @param {string} data.type - 资源类型
   * @param {string} data.unit - 计量单位
   * @param {number} data.quantity - 可用数量
   * @param {number} data.cost_per_unit - 单位成本
   * @param {string} data.color - 标识颜色
   * @param {boolean} data.is_active - 是否启用
   * @returns {Promise} 返回更新后的资源信息
   */
  updateResource(resourceId, data) {
    return request({
      url: `/progress/project/${data.project_id}/resources/${resourceId}`,
      method: 'PUT',
      data
    })
  },

  /**
   * 删除资源
   *
   * 删除指定的资源
   *
   * @param {number} projectId - 项目ID
   * @param {number} resourceId - 资源ID
   * @returns {Promise} 返回删除结果
   */
  deleteResource(projectId, resourceId) {
    return request({
      url: `/progress/project/${projectId}/resources/${resourceId}`,
      method: 'DELETE'
    })
  }
}

/**
 * 用户管理 API
 *
 * 管理系统用户账号，包括用户的增删改查和密码重置
 *
 * @namespace userApi
 */
export const userApi = {
  /**
   * 获取用户列表
   *
   * 支持分页、搜索和角色过滤
   *
   * @param {Object} params - 查询参数
   * @param {number} params.page - 页码
   * @param {number} params.page_size - 每页数量
   * @param {string} params.search - 搜索关键字（用户名、姓名）
   * @param {number} params.role_id - 角色ID过滤
   * @returns {Promise} 返回用户列表
   */
  getList(params) {
    return request({
      url: '/auth/users',
      method: 'GET',
      params
    })
  },

  /**
   * 获取用户详情
   *
   * 获取用户的完整信息，包括角色和权限
   *
   * @param {number} id - 用户ID
   * @returns {Promise} 返回用户详情
   */
  getDetail(id) {
    return request({
      url: `/auth/users/${id}`,
      method: 'GET'
    })
  },

  /**
   * 创建用户
   *
   * 创建新的系统用户账号
   *
   * @param {Object} data - 用户信息
   * @param {string} data.username - 用户名（唯一）
   * @param {string} data.password - 密码
   * @param {string} data.full_name - 姓名
   * @param {string} data.email - 邮箱
   * @param {string} data.phone - 电话
   * @param {number[]} data.role_ids - 角色ID列表
   * @returns {Promise} 返回创建的用户信息
   */
  create(data) {
    return request({
      url: '/auth/users',
      method: 'POST',
      data
    })
  },

  /**
   * 更新用户
   *
   * 更新用户基本信息和角色
   *
   * @param {number} id - 用户ID
   * @param {Object} data - 更新的用户信息
   * @param {string} data.full_name - 姓名
   * @param {string} data.email - 邮箱
   * @param {string} data.phone - 电话
   * @param {number[]} data.role_ids - 角色ID列表
   * @returns {Promise} 返回更新后的用户信息
   */
  update(id, data) {
    return request({
      url: `/auth/users/${id}`,
      method: 'PUT',
      data
    })
  },

  /**
   * 删除用户
   *
   * 删除指定的用户账号
   * 注意：不能删除自己的账号
   *
   * @param {number} id - 用户ID
   * @returns {Promise} 返回删除结果
   */
  delete(id) {
    return request({
      url: `/auth/users/${id}`,
      method: 'DELETE'
    })
  },

  /**
   * 重置用户密码
   *
   * 管理员重置指定用户的密码
   *
   * @param {number} id - 用户ID
   * @param {Object} data - 密码信息
   * @param {string} data.new_password - 新密码
   * @returns {Promise} 返回重置结果
   */
  resetPassword(id, data) {
    return request({
      url: `/auth/users/${id}/reset-password`,
    method: 'POST',
      data
    })
  }
}

/**
 * 角色管理 API
 *
 * 管理系统角色和权限分配
 * 支持基于角色的访问控制（RBAC）
 *
 * @namespace roleApi
 */
export const roleApi = {
  /**
   * 获取角色列表
   *
   * 获取系统中所有的角色
   *
   * @returns {Promise} 返回角色列表
   */
  getList() {
    return request({
      url: '/auth/roles',
      method: 'GET'
    })
  },

  /**
   * 获取角色详情
   *
   * 获取角色的完整信息和权限列表
   *
   * @param {number} id - 角色ID
   * @returns {Promise} 返回角色详情
   */
  getDetail(id) {
    return request({
      url: `/auth/roles/${id}`,
      method: 'GET'
    })
  },

  /**
   * 创建角色
   *
   * 创建新的角色并分配权限
   *
   * @param {Object} data - 角色信息
   * @param {string} data.name - 角色名称
   * @param {string} data.description - 角色描述
   * @param {number[]} data.permission_ids - 权限ID列表
   * @returns {Promise} 返回创建的角色信息
   */
  create(data) {
    return request({
      url: '/auth/roles',
      method: 'POST',
      data
    })
  },

  /**
   * 更新角色
   *
   * 更新角色信息和权限
   *
   * @param {number} id - 角色ID
   * @param {Object} data - 更新的角色信息
   * @param {string} data.name - 角色名称
   * @param {string} data.description - 角色描述
   * @param {number[]} data.permission_ids - 权限ID列表
   * @returns {Promise} 返回更新后的角色信息
   */
  update(id, data) {
    return request({
      url: `/auth/roles/${id}`,
      method: 'PUT',
      data
    })
  },

  /**
   * 删除角色
   *
   * 删除指定的角色
   * 注意：如果有用户正在使用该角色，无法删除
   *
   * @param {number} id - 角色ID
   * @returns {Promise} 返回删除结果
   */
  delete(id) {
    return request({
      url: `/auth/roles/${id}`,
      method: 'DELETE'
    })
  },

  /**
   * 为角色分配权限
   *
   * 更新角色的权限列表
   *
   * @param {number} id - 角色ID
   * @param {Object} data - 权限信息
   * @param {number[]} data.permission_ids - 权限ID列表
   * @returns {Promise} 返回更新后的角色权限
   */
  assignPermissions(id, data) {
    return request({
      url: `/auth/roles/${id}/permissions`,
    method: 'POST',
      data
    })
  }
}

/**
 * AI 智能分析 API
 *
 * 提供自然语言查询、数据洞察、智能推荐等功能
 *
 * @namespace aiApi
 */
export const aiApi = {
  /**
   * 自然语言分析
   *
   * 使用AI分析自然语言问题并返回结果
   *
   * @param {Object} data - 分析请求
   * @param {string} data.question - 自然语言问题
   * @param {boolean} data.conversation_mode - 是否为对话模式
   * @param {string} data.conversation_id - 对话ID
   * @param {Array} data.conversation_history - 对话历史
   * @param {number} data.max_iterations - 最大迭代次数
   * @returns {Promise} 返回分析结果
   */
  analyze(data) {
    return request({
      url: '/system/ai/analyze',
      method: 'POST',
      data
    })
  },

  /**
   * 获取数据洞察
   *
   * 获取AI生成的数据洞察信息
   *
   * @param {Object} params - 查询参数
   * @param {string} params.type - 类型（dashboard/inventory/requisitions/users/projects/inbound）
   * @returns {Promise} 返回洞察数据
   */
  getInsights(params) {
    return request({
      url: '/system/ai/insights',
      method: 'GET',
      params
    })
  },

  /**
   * 获取建议/推荐
   *
   * 获取快捷问题列表或智能推荐建议
   *
   * @param {Object} params - 查询参数
   * @param {string} params.type - 类型（questions/recommendations）
   * @returns {Promise} 返回建议或推荐列表
   */
  getSuggestions(params) {
    return request({
      url: '/system/ai/suggestions',
      method: 'GET',
      params
    })
  },

  /**
   * 获取分析历史
   *
   * 获取AI分析的历史记录
   *
   * @param {Object} params - 查询参数
   * @param {number} params.page - 页码
   * @param {number} params.per_page - 每页数量
   * @param {string} params.status - 状态筛选
   * @param {string} params.my_only - 是否只看我的记录
   * @returns {Promise} 返回历史记录列表
   */
  getHistory(params) {
    return request({
      url: '/system/ai/history',
      method: 'GET',
      params
    })
  },

  /**
   * 获取历史详情
   *
   * 获取单条AI分析历史的详细信息
   *
   * @param {number} id - 历史记录ID
   * @returns {Promise} 返回历史详情
   */
  getHistoryDetail(id) {
    return request({
      url: `/system/ai/history/${id}`,
      method: 'GET'
    })
  },

  /**
   * 删除历史记录
   *
   * 删除指定的AI分析历史记录
   *
   * @param {number} id - 历史记录ID
   * @returns {Promise} 返回删除结果
   */
  deleteHistory(id) {
    return request({
      url: `/system/ai/history/${id}`,
      method: 'DELETE'
    })
  },

  /**
   * 获取统计信息
   *
   * 获取AI分析的统计信息
   *
   * @returns {Promise} 返回统计信息
   */
  getStats() {
    return request({
      url: '/system/ai/stats',
      method: 'GET'
    })
  },

  /**
   * 获取AI配置
   *
   * 获取当前的AI配置参数
   *
   * @returns {Promise} 返回配置信息
   */
  getConfig() {
    return request({
      url: '/system/ai/config',
      method: 'GET'
    })
  },

  /**
   * 更新AI配置
   *
   * 更新AI配置参数
   *
   * @param {Object} data - 配置信息
   * @returns {Promise} 返回更新结果
   */
  updateConfig(data) {
    return request({
      url: '/system/ai/config',
      method: 'POST',
      data
    })
  },

  /**
   * 检查AI状态
   *
   * 检查AI服务是否正常运行
   *
   * @returns {Promise} 返回状态信息
   */
  getStatus() {
    return request({
      url: '/system/ai/status',
      method: 'GET'
    })
  }
}

/**
 * 工作流管理 API
 *
 * 提供完整的工作流审批系统功能，包括：
 * - 工作流定义管理（设计、激活、停用）
 * - 工作流实例管理（查看、跟踪）
 * - 工作流任务管理（审批、拒绝、退回、评论）
 *
 * @namespace workflowApi
 */
export const workflowApi = {
  // ========== 工作流定义管理 ==========

  /**
   * 获取工作流列表
   *
   * 获取所有工作流定义
   *
   * @param {Object} params - 查询参数
   * @param {number} params.page - 页码
   * @param {number} params.page_size - 每页数量
   * @param {string} params.status - 状态过滤（active/inactive）
   * @returns {Promise} 返回工作流列表
   */
  getList(params) {
    return request({
      url: '/workflows',
      method: 'GET',
      params
    })
  },

  /**
   * 获取工作流详情
   *
   * 获取工作流的完整定义，包括节点和连接线
   *
   * @param {number} id - 工作流ID
   * @returns {Promise} 返回工作流详情
   */
  getDetail(id) {
    return request({
      url: `/workflows/${id}`,
      method: 'GET'
    })
  },

  /**
   * 创建工作流
   *
   * 创建新的工作流定义
   *
   * @param {Object} data - 工作流定义
   * @param {string} data.name - 工作流名称
   * @param {string} data.description - 工作流描述
   * @param {string} data.entity_type - 关联实体类型（inbound_order/requisition）
   * @param {Array} data.nodes - 工作流节点定义
   * @param {Array} data.edges - 节点连接线定义
   * @returns {Promise} 返回创建的工作流信息
   */
  create(data) {
    return request({
      url: '/workflows',
      method: 'POST',
      data
    })
  },

  /**
   * 更新工作流
   *
   * 更新工作流定义
   *
   * @param {number} id - 工作流ID
   * @param {Object} data - 更新的工作流定义
   * @returns {Promise} 返回更新后的工作流信息
   */
  update(id, data) {
    return request({
      url: `/workflows/${id}`,
      method: 'PUT',
      data
    })
  },

  /**
   * 删除工作流
   *
   * 删除工作流定义
   * 注意：只能删除没有活跃实例的工作流
   *
   * @param {number} id - 工作流ID
   * @returns {Promise} 返回删除结果
   */
  delete(id) {
    return request({
      url: `/workflows/${id}`,
      method: 'DELETE'
    })
  },

  /**
   * 激活工作流
   *
   * 激活工作流定义，使其可以被使用
   *
   * @param {number} id - 工作流ID
   * @returns {Promise} 返回激活结果
   */
  activate(id) {
    return request({
      url: `/workflows/${id}/activate`,
      method: 'PUT'
    })
  },

  /**
   * 停用工作流
   *
   * 停用工作流定义，不再创建新实例
   *
   * @param {number} id - 工作流ID
   * @returns {Promise} 返回停用结果
   */
  deactivate(id) {
    return request({
      url: `/workflows/${id}/deactivate`,
      method: 'PUT'
    })
  },

  // ========== 工作流实例管理 ==========

  /**
   * 获取工作流实例列表
   *
   * 查询所有的工作流实例（运行中的流程）
   *
   * @param {Object} params - 查询参数
   * @param {number} params.page - 页码
   * @param {number} params.page_size - 每页数量
   * @param {string} params.status - 状态过滤（running/completed/cancelled）
   * @returns {Promise} 返回实例列表
   */
  getInstances(params) {
    return request({
      url: '/workflow-instances',
      method: 'GET',
      params
    })
  },

  /**
   * 获取实例详情
   *
   * 获取工作流实例的完整信息
   *
   * @param {number} id - 实例ID
   * @returns {Promise} 返回实例详情
   */
  getInstance(id) {
    return request({
      url: `/workflow-instances/${id}`,
      method: 'GET'
    })
  },

  /**
   * 获取实例审批记录
   *
   * 查看实例的所有审批历史记录
   *
   * @param {number} id - 实例ID
   * @returns {Promise} 返回审批记录列表
   */
  getInstanceApprovals(id) {
    return request({
      url: `/workflow-instances/${id}/approvals`,
      method: 'GET'
    })
  },

  /**
   * 获取实例日志
   *
   * 查看实例的完整操作日志
   *
   * @param {number} id - 实例ID
   * @returns {Promise} 返回日志列表
   */
  getInstanceLogs(id) {
    return request({
      url: `/workflow-instances/${id}/logs`,
      method: 'GET'
    })
  },

  // ========== 工作流任务管理 ==========

  /**
   * 获取待办任务列表
   *
   * 获取当前用户需要处理的任务
   *
   * @param {Object} params - 查询参数
   * @param {number} params.page - 页码
   * @param {number} params.page_size - 每页数量
   * @returns {Promise} 返回待办任务列表
   */
  getPendingTasks(params) {
    return request({
      url: '/workflow-tasks/pending',
      method: 'GET',
      params
    })
  },

  /**
   * 审批通过任务
   *
   * 同意当前节点的审批，流转到下一节点
   *
   * @param {number} id - 任务ID
   * @param {Object} data - 审批信息
   * @param {string} data.comment - 审批意见
   * @returns {Promise} 返回审批结果
   */
  approveTask(id, data) {
    return request({
      url: `/workflow-tasks/${id}/approve`,
    method: 'POST',
      data
    })
  },

  /**
   * 审批拒绝任务
   *
   * 拒绝当前节点的审批，终止流程
   *
   * @param {number} id - 任务ID
   * @param {Object} data - 拒绝信息
   * @param {string} data.reason - 拒绝原因
   * @returns {Promise} 返回拒绝结果
   */
  rejectTask(id, data) {
    return request({
      url: `/workflow-tasks/${id}/reject`,
    method: 'POST',
      data
    })
  },

  /**
   * 退回任务
   *
   * 将任务退回到上一节点或指定节点
   *
   * @param {number} id - 任务ID
   * @param {Object} data - 退回信息
   * @param {number} data.to_node_id - 目标节点ID
   * @param {string} data.comment - 退回原因
   * @returns {Promise} 返回退回结果
   */
  returnTask(id, data) {
    return request({
      url: `/workflow-tasks/${id}/return`,
      method: 'POST',
      data
    })
  },

  /**
   * 任务评论
   *
   * 对任务添加评论或说明
   *
   * @param {number} id - 任务ID
   * @param {Object} data - 评论信息
   * @param {string} data.comment - 评论内容
   * @returns {Promise} 返回评论结果
   */
  commentTask(id, data) {
    return request({
      url: `/workflow-tasks/${id}/comment`,
      method: 'POST',
      data
    })
  }
}

// ==================== 物资计划管理 API ====================

/**
 * 物资计划管理 API 接口
 *
 * 提供物资计划的增删改查功能，支持分页、搜索和筛选
 *
 * @namespace materialPlanApi
 */
export const materialPlanApi = {
  /**
   * 获取物资计划列表
   *
   * @param {Object} params - 查询参数
   * @param {number} params.page - 页码
   * @param {number} params.page_size - 每页数量
   * @param {string} params.status - 状态过滤
   * @param {string} params.plan_type - 计划类型
   * @param {string} params.priority - 优先级
   * @param {string} params.search - 搜索关键词
   * @param {string} params.project_id - 项目ID
   * @param {string} params.project_ids - 项目ID列表（逗号分隔）
   * @param {string} params.start_date - 开始日期
   * @param {string} params.end_date - 结束日期
   * @returns {Promise} 返回计划列表
   */
  getPlans(params) {
    return request({
      url: '/material-plan/plans',
      method: 'GET',
      params
    })
  },

  /**
   * 获取计划详情
   *
   * @param {number} id - 计划ID
   * @returns {Promise} 返回计划详情
   */
  getPlanDetail(id) {
    return request({
      url: `/material-plan/plans/${id}`,
      method: 'GET'
    })
  },

  /**
   * 创建物资计划
   *
   * @param {Object} data - 计划数据
   * @param {string} data.plan_name - 计划名称
   * @param {number} data.project_id - 项目ID
   * @param {string} data.plan_type - 计划类型
   * @param {string} data.priority - 优先级
   * @param {string} data.planned_start_date - 计划开始日期
   * @param {string} data.planned_end_date - 计划结束日期
   * @param {number} data.total_budget - 总预算
   * @param {string} data.description - 描述
   * @param {string} data.remark - 备注
   * @param {Array} data.items - 计划项目列表
   * @returns {Promise} 返回创建的计划信息
   */
  createPlan(data) {
    return request({
      url: '/material-plan/plans',
      method: 'POST',
      data
    })
  },

  /**
   * 更新物资计划
   *
   * @param {number} id - 计划ID
   * @param {Object} data - 更新的计划数据
   * @returns {Promise} 返回更新后的计划信息
   */
  updatePlan(id, data) {
    return request({
      url: `/material-plan/plans/${id}`,
      method: 'PUT',
      data
    })
  },

  /**
   * 删除物资计划
   *
   * @param {number} id - 计划ID
   * @returns {Promise} 返回删除结果
   */
  deletePlan(id) {
    return request({
      url: `/material-plan/plans/${id}`,
      method: 'DELETE'
    })
  },

  /**
   * 提交计划审批
   *
   * @param {number} id - 计划ID
   * @returns {Promise} 返回提交结果
   */
  submitPlan(id) {
    return request({
      url: `/material-plan/plans/${id}/submit`,
      method: 'POST'
    })
  },

  /**
   * 批准计划
   *
   * @param {number} id - 计划ID
   * @param {Object} data - 审批信息
   * @param {string} data.remark - 审批备注
   * @returns {Promise} 返回审批结果
   */
  approvePlan(id, data) {
    return request({
      url: `/material-plan/plans/${id}/approve`,
    method: 'POST',
      data
    })
  },

  /**
   * 拒绝计划
   *
   * @param {number} id - 计划ID
   * @param {Object} data - 拒绝信息
   * @param {string} data.remark - 拒绝原因
   * @returns {Promise} 返回拒绝结果
   */
  rejectPlan(id, data) {
    return request({
      url: `/material-plan/plans/${id}/reject`,
    method: 'POST',
      data
    })
  },

  /**
   * 激活计划
   *
   * @param {number} id - 计划ID
   * @returns {Promise} 返回激活结果
   */
  activatePlan(id) {
    return request({
      url: `/material-plan/plans/${id}/activate`,
      method: 'POST'
    })
  },

  /**
   * 重新提交计划
   *
   * @param {number} id - 计划ID
   * @returns {Promise} 返回提交结果
   */
  resubmitPlan(id) {
    return request({
      url: `/material-plan/plans/${id}/resubmit`,
      method: 'POST'
    })
  },

  /**
   * 取消计划
   *
   * @param {number} id - 计划ID
   * @param {Object} data - 取消信息
   * @param {string} data.reason - 取消原因
   * @returns {Promise} 返回取消结果
   */
  cancelPlan(id, data) {
    return request({
      url: `/material-plan/plans/${id}/cancel`,
    method: 'POST',
      data
    })
  },

  /**
   * 获取计划项目列表
   *
   * @param {number} id - 计划ID
   * @returns {Promise} 返回项目列表
   */
  getPlanItems(id) {
    return request({
      url: `/material-plan/plans/${id}/items`,
      method: 'GET'
    })
  },

  /**
   * 添加计划项目
   *
   * @param {number} id - 计划ID
   * @param {Object} data - 项目数据
   * @returns {Promise} 返回添加的项目信息
   */
  addPlanItem(id, data) {
    return request({
      url: `/material-plan/plans/${id}/items`,
    method: 'POST',
      data
    })
  },

  /**
   * 更新计划项目
   *
   * @param {number} id - 项目ID
   * @param {Object} data - 更新的项目数据
   * @returns {Promise} 返回更新后的项目信息
   */
  updatePlanItem(id, data) {
    return request({
      url: `/material-plan/items/${id}`,
      method: 'PUT',
      data
    })
  },

  /**
   * 删除计划项目
   *
   * @param {number} id - 项目ID
   * @returns {Promise} 返回删除结果
   */
  deletePlanItem(id) {
    return request({
      url: `/material-plan/items/${id}`,
      method: 'DELETE'
    })
  },

  /**
   * 获取计划统计概览
   *
   * @returns {Promise} 返回统计数据
   */
  getStatistics() {
    return request({
      url: '/material-plan/statistics/overview',
      method: 'GET'
    })
  },

  /**
   * 获取计划详细统计
   *
   * @param {number} id - 计划ID
   * @returns {Promise} 返回统计数据
   */
  getPlanStatistics(id) {
    return request({
      url: `/material-plan/statistics/plan/${id}`,
      method: 'GET'
    })
  },

  /**
   * 获取计划工作流状态
   *
   * @param {number} id - 计划ID
   * @returns {Promise} 返回工作流状态信息
   */
  getPlanWorkflow(id) {
    return request({
      url: `/material-plan/plans/${id}/workflow`,
      method: 'GET'
    })
  },

  /**
   * 获取计划审批记录
   *
   * @param {number} id - 计划ID
   * @returns {Promise} 返回审批记录列表
   */
  getPlanApprovals(id) {
    return request({
      url: `/material-plan/plans/${id}/approvals`,
      method: 'GET'
    })
  },

  /**
   * 获取待办任务列表
   *
   * @returns {Promise} 返回待办任务列表
   */
  getPendingTasks() {
    return request({
      url: '/material-plan/workflow/pending',
      method: 'GET'
    })
  },

  /**
   * 同步计划项的物资ID
   *
   * 更新物资计划项的material_id
   *
   * @param {number} planId - 计划ID
   * @param {Object} data - 更新数据
   * @param {Array} data.items - 计划项列表
   * @returns {Promise} 返回更新结果
   */
  syncPlanMaterialIds(planId, data) {
    return request({
      url: `/material-plan/plans/${planId}/sync-materials`,
      method: 'POST',
      data
    })
  }
}

// ==================== 文件上传 API ====================

/**
 * 文件上传 API 接口
 *
 * 提供图片、文件等资源上传功能
 *
 * @namespace uploadApi
 */
export const uploadApi = {
  /**
   * 上传图片
   *
   * 用于富文本编辑器等场景的图片上传
   *
   * @param {FormData} formData - 包含文件的FormData对象
   * @param {File} formData.file - 图片文件
   * @returns {Promise} 返回上传结果，包含图片URL
   *
   * @example
   * const formData = new FormData()
   * formData.append('file', file)
   * const result = await uploadApi.uploadImage(formData)
   * console.log(result.data.url) // 图片URL
   */
  uploadImage(formData) {
    return request({
      url: '/upload/image',
      method: 'POST',
      data: formData,
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  /**
   * 上传文件
   *
   * 用于通用文件上传
   *
   * @param {FormData} formData - 包含文件的FormData对象
   * @param {File} formData.file - 文件
   * @returns {Promise} 返回上传结果，包含文件URL
   */
  uploadFile(formData) {
    return request({
      url: '/upload/file',
      method: 'POST',
      data: formData,
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  /**
   * 批量上传图片
   *
   * @param {FormData} formData - 包含多个文件的FormData对象
   * @returns {Promise} 返回上传结果，包含所有文件URL
   */
  uploadImages(formData) {
    return request({
      url: '/upload/images',
      method: 'POST',
      data: formData,
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  }
}

// ==================== 通知管理 API ====================

/**
 * 通知管理 API 接口
 *
 * 提供系统通知的查询、标记已读、删除等功能
 *
 * @namespace notificationApi
 */
export const notificationApi = {
  /**
   * 获取通知列表
   *
   * @param {Object} params - 查询参数
   * @param {number} params.page - 页码
   * @param {number} params.per_page - 每页数量
   * @returns {Promise} 返回通知列表
   */
  getList(params) {
    return request({
      url: '/notification/notifications',
      method: 'GET',
      params
    })
  },

  /**
   * 获取未读通知数量
   *
   * @returns {Promise} 返回未读数量
   */
  getUnreadCount() {
    return request({
      url: '/notification/notifications/count',
      method: 'GET'
    })
  },

  /**
   * 标记通知为已读
   *
   * @param {number} id - 通知ID
   * @returns {Promise} 返回标记结果
   */
  markAsRead(id) {
    return request({
      url: `/notification/notifications/${id}/read`,
      method: 'PUT'
    })
  },

  /**
   * 全部标记为已读
   *
   * @returns {Promise} 返回标记结果
   */
  markAllAsRead() {
    return request({
      url: '/notification/notifications/read-all',
      method: 'PUT'
    })
  },

  /**
   * 删除通知
   *
   * @param {number} id - 通知ID
   * @returns {Promise} 返回删除结果
   */
  delete(id) {
    return request({
      url: `/notification/notifications/${id}`,
      method: 'DELETE'
    })
  },

  /**
   * 清空所有通知
   *
   * @returns {Promise} 返回删除结果
   */
  clearAll() {
    return request({
      url: '/notification/notifications',
      method: 'DELETE'
    })
  }
}

/**
 * AI Agent API
 *
 * 提供AI Agent操作接口，支持未来接入Claude Desktop等AI Agent
 *
 * @namespace agentApi
 */
export const agentApi = {
  /**
   * 获取可用能力列表
   *
   * 获取当前AI Agent支持的所有操作和资源
   *
   * @returns {Promise} 返回能力列表
   */
  getCapabilities() {
    return request({
      url: '/agent/capabilities',
      method: 'GET'
    })
  },

  /**
   * 验证操作合法性
   *
   * 预验证操作是否被允许，不执行实际操作
   *
   * @param {Object} data - 验证请求
   * @param {string} data.operation - 操作类型
   * @param {string} data.resource - 资源类型
   * @param {Object} data.parameters - 操作参数
   * @returns {Promise} 返回验证结果
   */
  validateOperation(operation, resource, parameters) {
    return request({
      url: '/agent/validate',
      method: 'POST',
      data: {
        operation,
        resource,
        parameters
      }
    })
  },

  /**
   * 执行 AI Agent 操作
   *
   * 执行指定的AI Agent操作
   *
   * @param {string} operation - 操作类型 (query/analyze/create_material_plan/update_stock/approve_workflow/generate_report)
   * @param {string} resource - 资源类型 (material/stock/workflow/material_plan)
   * @param {Object} parameters - 操作参数
   * @param {string} reasoning - AI推理过程说明
   * @returns {Promise} 返回操作结果
   */
  operate(operation, resource, parameters, reasoning) {
    return request({
      url: '/agent/operate',
      method: 'POST',
      data: {
        operation,
        resource,
        parameters,
        reasoning
      }
    })
  },

  /**
   * AI 查询
   *
   * 使用自然语言查询数据
   *
   * @param {string} question - 查询问题
   * @param {number} limit - 结果数量限制 (默认10)
   * @param {Array} fields - 返回字段列表
   * @param {Object} filters - 过滤条件
   * @returns {Promise} 返回查询结果
   */
  query(question, limit = 10, fields = null, filters = null) {
    return request({
      url: '/agent/query',
      method: 'POST',
      data: {
        question,
        limit,
        fields,
        filters
      }
    })
  },

  /**
   * 工作流操作
   *
   * 执行工作流审批操作
   *
   * @param {number} taskId - 任务ID
   * @param {string} action - 操作类型 (approve/reject/return)
   * @param {string} remark - 备注/意见
   * @param {number} toNodeId - 目标节点ID (用于退回操作)
   * @returns {Promise} 返回操作结果
   */
  workflow(taskId, action, remark = '', toNodeId = null) {
    return request({
      url: '/agent/workflow',
      method: 'POST',
      data: {
        task_id: taskId,
        action,
        remark,
        to_node_id: toNodeId
      }
    })
  },

  /**
   * 获取操作日志
   *
   * 查询AI Agent操作日志
   *
   * @param {Object} params - 查询参数
   * @param {number} params.page - 页码
   * @param {number} params.page_size - 每页数量
   * @param {string} params.operation - 操作类型过滤
   * @param {string} params.resource - 资源类型过滤
   * @param {string} params.status - 状态过滤 (pending/completed/failed)
   * @param {string} params.agent_id - Agent ID过滤
   * @param {string} params.start_date - 开始日期
   * @param {string} params.end_date - 结束日期
   * @returns {Promise} 返回日志列表
   */
  getLogs(params) {
    return request({
      url: '/agent/logs',
      method: 'GET',
      params
    })
  },

  // ========== 便捷方法 ==========

  /**
   * 查询物资
   *
   * @param {string} search - 搜索关键词
   * @param {number} limit - 结果数量
   */
  queryMaterials(search = '', limit = 10) {
    return this.operate(
      'query',
      'material',
      { search, limit },
      'Query materials from user request'
    )
  },

  /**
   * 分析库存
   *
   * @param {string} question - 分析问题
   */
  analyzeInventory(question = '库存分析') {
    return this.operate(
      'analyze',
      'inventory',
      { question },
      `Inventory analysis: ${question}`
    )
  },

  /**
   * 获取库存预警
   */
  getStockAlerts() {
    return this.operate(
      'query',
      'stock',
      { low_stock_alert: true },
      'Get low stock alerts'
    )
  },

  /**
   * 获取待办任务
   */
  getPendingTasks() {
    return this.operate(
      'query',
      'workflow',
      { status: 'pending' },
      'Get pending workflow tasks'
    )
  },

  /**
   * 创建物资计划
   *
   * @param {number} projectId - 项目ID
   * @param {Array} items - 计划明细
   * @param {string} remark - 备注
   */
  createMaterialPlan(projectId, items, remark = '') {
    return this.operate(
      'create_material_plan',
      'material_plan',
      {
        project_id: projectId,
        items,
        remark
      },
      `AI generated material plan for project ${projectId}`
    )
  },

  /**
   * 更新库存
   *
   * @param {number} stockId - 库存ID
   * @param {number} quantity - 新数量
   * @param {string} remark - 备注
   */
  updateStock(stockId, quantity, remark = '') {
    return this.operate(
      'update_stock',
      'stock',
      {
        stock_id: stockId,
        quantity,
        remark
      },
      `Stock update for stock_id ${stockId}`
    )
  },

  /**
   * 生成报表
   *
   * @param {string} reportType - 报表类型 (inventory_summary/material_plan_summary)
   */
  generateReport(reportType = 'inventory_summary') {
    return this.operate(
      'generate_report',
      'report',
      { report_type: reportType },
      `Generate ${reportType} report`
    )
  }
}

// ==================== 施工预约管理 API ====================

/**
 * 施工预约管理 API 接口
 *
 * 提供施工预约单的增删改查功能，支持日历排期、作业人员分配、工作流审批等
 *
 * @namespace appointmentApi
 */
export const appointmentApi = {
  /**
   * 获取预约单列表
   *
   * 支持分页、搜索、排序和状态过滤
   *
   * @param {Object} params - 查询参数
   * @param {number} params.page - 页码（从 1 开始）
   * @param {number} params.page_size - 每页数量
   * @param {string} params.status - 状态过滤
   * @param {boolean} params.is_urgent - 是否加急
   * @param {string} params.start_date - 开始日期
   * @param {string} params.end_date - 结束日期
   * @param {number} params.applicant_id - 申请人ID
   * @param {number} params.worker_id - 作业人员ID
   * @param {string} params.work_type - 作业类型
   * @returns {Promise} 返回预约单列表和分页信息
   */
  getList(params) {
    return request({
      url: '/appointments',
      method: 'GET',
      params
    })
  },

  /**
   * 获取我的预约列表
   *
   * @param {Object} params - 查询参数
   * @returns {Promise} 返回我的预约列表
   */
  getMyList(params) {
    return request({
      url: '/appointments/my',
      method: 'GET',
      params
    })
  },

  /**
   * 获取待审批列表
   *
   * @param {Object} params - 查询参数
   * @returns {Promise} 返回待审批列表
   */
  getPendingApprovals(params) {
    return request({
      url: '/appointments/pending',
      method: 'GET',
      params
    })
  },

  /**
   * 获取预约单详情
   *
   * @param {number} id - 预约单ID
   * @returns {Promise} 返回预约单详情
   */
  getDetail(id) {
    return request({
      url: `/appointments/${id}`,
      method: 'GET'
    })
  },

  /**
   * 创建预约单
   *
   * @param {Object} data - 预约单信息
   * @param {number} data.project_id - 项目ID
   * @param {string} data.contact_phone - 联系电话
   * @param {string} data.contact_person - 联系人
   * @param {string} data.work_date - 作业日期
   * @param {string} data.time_slot - 时间段 (morning/afternoon/evening/full_day)
   * @param {string} data.work_location - 作业地点
   * @param {string} data.work_content - 作业内容
   * @param {string} data.work_type - 作业类型
   * @param {boolean} data.is_urgent - 是否加急
   * @param {number} data.priority - 优先级 (0-10)
   * @param {string} data.urgent_reason - 加急原因
   * @param {number} data.assigned_worker_id - 指派的作业人员ID
   * @returns {Promise} 返回创建的预约单信息
   */
  create(data) {
    return request({
      url: '/appointments',
      method: 'POST',
      data
    })
  },

  /**
   * 批量创建预约单
   *
   * @param {Object} data - 批量创建数据
   * @param {Array} data.appointments - 预约单数组
   * @returns {Promise} 返回批量创建结果
   */
  batchCreate(data) {
    return request({
      url: '/appointments/batch',
      method: 'POST',
      data
    })
  },

  /**
   * 更新预约单
   *
   * 只能更新草稿状态的预约单
   *
   * @param {number} id - 预约单ID
   * @param {Object} data - 更新的预约单信息
   * @returns {Promise} 返回更新后的预约单信息
   */
  update(id, data) {
    return request({
      url: `/appointments/${id}`,
      method: 'PUT',
      data
    })
  },

  /**
   * 删除预约单
   *
   * 只能删除草稿状态的预约单
   *
   * @param {number} id - 预约单ID
   * @returns {Promise} 返回删除结果
   */
  delete(id) {
    return request({
      url: `/appointments/${id}`,
      method: 'DELETE'
    })
  },

  /**
   * 提交审批
   *
   * @param {number} id - 预约单ID
   * @returns {Promise} 返回提交结果
   */
  submit(id) {
    return request({
      url: `/appointments/${id}/submit`,
      method: 'POST'
    })
  },

  /**
   * 启动工作流
   *
   * @param {number} id - 预约单ID
   * @param {Object} data - 工作流信息
   * @param {number} data.workflow_id - 工作流ID
   * @returns {Promise} 返回启动结果
   */
  startWorkflow(id, data) {
    return request({
      url: `/appointments/${id}/workflow/start`,
      method: 'POST',
      data
    })
  },

  /**
   * 审批预约单
   *
   * @param {number} id - 预约单ID
   * @param {Object} data - 审批信息
   * @param {string} data.action - 操作 (approve/reject)
   * @param {string} data.comment - 审批意见
   * @param {boolean} data.assign_now - 是否立即分配作业人员
   * @param {number} data.worker_id - 指定的作业人员ID
   * @returns {Promise} 返回审批结果
   */
  approve(id, data) {
    return request({
      url: `/appointments/${id}/approve`,
      method: 'POST',
      data
    })
  },

  /**
   * 撤回预约单
   *
   * @param {number} id - 预约单ID
   * @returns {Promise} 返回撤回结果
   */
  recall(id) {
    return request({
      url: `/appointments/${id}/recall`,
      method: 'POST'
    })
  },

  /**
   * 分配作业人员
   *
   * @param {number} id - 预约单ID
   * @param {Object} data - 分配信息
   * @param {number} data.worker_id - 作业人员ID
   * @returns {Promise} 返回分配结果
   */
  assignWorker(id, data) {
    return request({
      url: `/appointments/${id}/assign`,
      method: 'POST',
      data
    })
  },

  /**
   * 开始作业
   *
   * @param {number} id - 预约单ID
   * @returns {Promise} 返回操作结果
   */
  startWork(id) {
    return request({
      url: `/appointments/${id}/start`,
      method: 'POST'
    })
  },

  /**
   * 完成作业
   *
   * @param {number} id - 预约单ID
   * @param {Object} data - 完成信息
   * @param {string} data.completion_note - 完成备注
   * @param {Array} data.photos - 完成照片URL列表
   * @returns {Promise} 返回完成结果
   */
  complete(id, data) {
    return request({
      url: `/appointments/${id}/complete`,
      method: 'POST',
      data
    })
  },

  /**
   * 取消预约
   *
   * @param {number} id - 预约单ID
   * @param {Object} data - 取消信息
   * @param {string} data.reason - 取消原因
   * @returns {Promise} 返回取消结果
   */
  cancel(id, data) {
    return request({
      url: `/appointments/${id}/cancel`,
      method: 'POST',
      data
    })
  },

  /**
   * 获取审批历史
   *
   * @param {number} id - 预约单ID
   * @returns {Promise} 返回审批历史记录
   */
  getApprovalHistory(id) {
    return request({
      url: `/appointments/${id}/approval-history`,
      method: 'GET'
    })
  },

  /**
   * 获取工作流进度
   *
   * @param {number} id - 预约单ID
   * @returns {Promise} 返回工作流进度信息
   */
  getWorkflowProgress(id) {
    return request({
      url: `/appointments/${id}/workflow-progress`,
      method: 'GET'
    })
  },

  /**
   * 获取当前审批节点
   *
   * @param {number} id - 预约单ID
   * @returns {Promise} 返回当前审批节点信息
   */
  getCurrentApproval(id) {
    return request({
      url: `/appointments/${id}/current-approval`,
      method: 'GET'
    })
  },

  /**
   * 批量审批
   *
   * @param {Object} data - 批量审批数据
   * @param {Array} data.instance_ids - 工作流实例ID数组
   * @param {string} data.action - 操作 (approve/reject)
   * @param {string} data.comment - 审批意见
   * @returns {Promise} 返回批量审批结果
   */
  batchApprove(data) {
    return request({
      url: '/appointments/batch-approve',
      method: 'POST',
      data
    })
  },

  /**
   * 获取统计数据
   *
   * @param {Object} params - 查询参数
   * @param {string} params.date - 统计日期
   * @param {number} params.applicant_id - 申请人ID
   * @returns {Promise} 返回统计数据
   */
  getStats(params) {
    return request({
      url: '/appointments/stats',
      method: 'GET',
      params
    })
  },

  /**
   * 获取作业人员列表
   *
   * @returns {Promise} 返回作业人员列表
   */
  getWorkersList() {
    return request({
      url: '/appointments/workers',
      method: 'GET'
    })
  },

  /**
   * 获取每日预约统计数据
   *
   * @param {Object} params - 查询参数
   * @param {string} params.start_date - 开始日期
   * @param {string} params.end_date - 结束日期
   * @returns {Promise} 返回每日统计数据
   */
  getDailyStatistics(params) {
    return request({
      url: '/appointments/daily-statistics',
      method: 'GET',
      params
    })
  },

  // ========== 日历相关API ==========

  /**
   * 获取作业人员日历
   *
   * @param {number} workerId - 作业人员ID
   * @param {Object} params - 查询参数
   * @param {string} params.start_date - 开始日期
   * @param {string} params.end_date - 结束日期
   * @returns {Promise} 返回日历数据
   */
  getWorkerCalendar(workerId, params) {
    return request({
      url: `/appointments/calendar/worker/${workerId}`,
      method: 'GET',
      params
    })
  },

  /**
   * 检查可用性
   *
   * @param {Object} data - 检查数据
   * @param {number} data.worker_id - 作业人员ID
   * @param {string} data.work_date - 作业日期
   * @param {string} data.time_slot - 时间段
   * @returns {Promise} 返回可用性检查结果
   */
  checkAvailability(data) {
    return request({
      url: '/appointments/calendar/check-availability',
      method: 'POST',
      data
    })
  },

  /**
   * 批量锁定日历
   *
   * @param {Object} data - 锁定数据
   * @param {number} data.worker_id - 作业人员ID
   * @param {string} data.start_date - 开始日期
   * @param {string} data.end_date - 结束日期
   * @param {Array} data.time_slots - 时间段数组
   * @param {string} data.blocked_reason - 锁定原因
   * @returns {Promise} 返回锁定结果
   */
  batchBlockCalendar(data) {
    return request({
      url: '/appointments/calendar/batch-block',
      method: 'POST',
      data
    })
  },

  /**
   * 获取可用作业人员
   *
   * @param {Object} params - 查询参数
   * @param {string} params.work_date - 作业日期
   * @param {string} params.time_slot - 时间段
   * @returns {Promise} 返回可用作业人员列表
   */
  getAvailableWorkers(params) {
    return request({
      url: '/appointments/calendar/available-workers',
      method: 'GET',
      params
    })
  },

  /**
   * 获取日历视图数据
   *
   * @param {Object} params - 查询参数
   * @param {string} params.start_date - 开始日期
   * @param {string} params.end_date - 结束日期
   * @param {number} params.worker_id - 作业人员ID（可选）
   * @returns {Promise} 返回日历视图数据
   */
  getCalendarView(params) {
    return request({
      url: '/appointments/calendar/view',
      method: 'GET',
      params
    })
  },

  /**
   * 搜索预约单
   *
   * @param {Object} params - 查询参数
   * @param {string} params.keyword - 搜索关键词
   * @param {number} params.page - 页码
   * @param {number} params.page_size - 每页数量
   * @returns {Promise} 返回搜索结果
   */
  search(params) {
    return request({
      url: '/appointments/search',
      method: 'GET',
      params
    })
  },

  /**
   * 获取作业人员的预约列表
   *
   * @param {number} workerId - 作业人员ID
   * @param {Object} params - 查询参数
   * @returns {Promise} 返回预约列表
   */
  getWorkerAppointments(workerId, params) {
    return request({
      url: `/appointments/worker/${workerId}`,
      method: 'GET',
      params
    })
  },

  /**
   * 导出预约单
   *
   * @param {Object} params - 查询参数
   * @param {string} params.ids - 预约单ID列表（逗号分隔）
   * @returns {Promise} 返回Excel文件Blob
   */
  export(params) {
    return request({
      url: '/appointments/export',
      method: 'GET',
      params,
      responseType: 'blob'
    })
  },

  // ========== 工具函数 ==========

  /**
   * 获取时间段标签
   *
   * @param {string} timeSlot - 时间段代码
   * @returns {string} 时间段标签
   */
  getTimeSlotLabel(timeSlot) {
    const labels = {
      morning: '上午',
      afternoon: '下午',
      evening: '晚上',
      full_day: '全天'
    }
    return labels[timeSlot] || timeSlot
  },

  /**
   * 获取状态标签
   *
   * @param {string} status - 状态代码
   * @returns {string} 状态标签
   */
  getStatusLabel(status) {
    const labels = {
      draft: '草稿',
      pending: '待审批',
      scheduled: '已排期',
      in_progress: '进行中',
      completed: '已完成',
      cancelled: '已取消',
      rejected: '已拒绝'
    }
    return labels[status] || status
  },

  /**
   * 获取状态颜色类型
   *
   * @param {string} status - 状态代码
   * @returns {string} Element Plus Tag类型
   */
  getStatusType(status) {
    const types = {
      draft: '',
      pending: 'warning',
      scheduled: 'primary',
      in_progress: 'info',
      completed: 'success',
      cancelled: 'info',
      rejected: 'danger'
    }
    return types[status] || ''
  },

  /**
   * 判断是否可编辑
   *
   * @param {string} status - 状态代码
   * @returns {boolean} 是否可编辑
   */
  isEditable(status) {
    return status === 'draft'
  },

  /**
   * 判断是否可取消
   *
   * @param {string} status - 状态代码
   * @returns {boolean} 是否可取消
   */
  isCancellable(status) {
    return ['draft', 'pending', 'scheduled'].includes(status)
  },

  /**
   * 判断是否可完成
   *
   * @param {string} status - 状态代码
   * @returns {boolean} 是否可完成
   */
  canComplete(status) {
    return ['in_progress', 'scheduled'].includes(status)
  },

  /**
   * 判断是否可开始
   *
   * @param {string} status - 状态代码
   * @returns {boolean} 是否可开始
   */
  canStart(status) {
    return status === 'scheduled'
  }
}

