/**
 * API 接口统一导出 (移动端)
 *
 * 本文件集中导出所有后端 API 接口
 *
 * @module API
 * @author Material Management System (Mobile)
 * @date 2025-02-03
 */

// 认证相关 API
export * as authApi from './auth.js'

// 项目管理 API
export * as projectApi from './project.js'

// 物资管理 API
export * as materialApi from './material.js'

// 库存管理 API
export * as stockApi from './stock.js'

// 领料申请 API
export * as requisitionApi from './requisition.js'

// 入库管理 API
export * as inboundApi from './inbound.js'

// 施工日志 API
export * as constructionApi from './construction.js'

// 物资计划 API
export * as materialPlanApi from './material_plan.js'

// 进度管理 API
export * as progressApi from './progress.js'

// 工作流 API
export * as workflowApi from './workflow.js'

// 通知 API
export * as notificationApi from './notification.js'

// 文件上传 API
export * as uploadApi from './upload.js'

// AI 相关 API
export * as aiApi from './ai.js'

// 用户相关 API
export * as userApi from './user.js'
