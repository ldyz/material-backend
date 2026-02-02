/**
 * 物资计划管理 Store
 *
 * 使用 Pinia 管理物资计划状态，包括：
 * - 计划列表
 * - 当前计划
 * - 统计数据
 * - 筛选条件
 *
 * @module PlanStore
 * @author Material Management System
 * @date 2025-02-01
 */

import { defineStore } from 'pinia'
import { materialPlanApi } from '@/api'
import { ElMessage } from 'element-plus'

/**
 * 物资计划状态管理 Store
 *
 * 状态设计说明：
 * - plans: 计划列表
 * - currentPlan: 当前查看/编辑的计划
 * - statistics: 统计数据
 * - loading: 加载状态
 * - filters: 筛选条件
 * - pagination: 分页信息
 */
export const usePlanStore = defineStore('plan', {
  /**
   * 状态定义
   */
  state: () => ({
    // 计划列表
    plans: [],

    // 当前计划
    currentPlan: null,

    // 统计数据
    statistics: null,

    // 当前计划统计
    currentPlanStatistics: null,

    // 加载状态
    loading: false,

    // 筛选条件
    filters: {
      status: '',
      plan_type: '',
      priority: '',
      search: '',
      project_id: null,
      start_date: '',
      end_date: ''
    },

    // 分页信息
    pagination: {
      page: 1,
      page_size: 20,
      total: 0,
      pages: 0
    }
  }),

  /**
   * Getters - 计算属性
   */
  getters: {
    /**
     * 获取草稿状态的计划
     */
    draftPlans: (state) => {
      return state.plans.filter(p => p.status === 'draft')
    },

    /**
     * 获取待审批的计划
     */
    pendingPlans: (state) => {
      return state.plans.filter(p => p.status === 'pending')
    },

    /**
     * 获取活跃的计划
     */
    activePlans: (state) => {
      return state.plans.filter(p => p.status === 'active')
    },

    /**
     * 获取已完成的计划
     */
    completedPlans: (state) => {
      return state.plans.filter(p => p.status === 'completed')
    },

    /**
     * 根据状态筛选计划
     */
    plansByStatus: (state) => (status) => {
      if (!status) return state.plans
      return state.plans.filter(p => p.status === status)
    },

    /**
     * 检查是否有数据
     */
    hasPlans: (state) => {
      return state.plans.length > 0
    },

    /**
     * 获取当前计划的项目数量
     */
    currentPlanItemsCount: (state) => {
      return state.currentPlan?.items_count || 0
    }
  },

  /**
   * Actions - 异步方法和业务逻辑
   */
  actions: {
    /**
     * 获取计划列表
     *
     * @param {Object} params - 查询参数
     * @returns {Promise<Object>} 响应数据
     */
    async fetchPlans(params = {}) {
      this.loading = true
      try {
        // 合并筛选条件
        const queryParams = {
          page: this.pagination.page,
          page_size: this.pagination.page_size,
          ...this.filters,
          ...params
        }

        const response = await materialPlanApi.getPlans(queryParams)

        if (response.success) {
          this.plans = response.data || []

          // 更新分页信息
          if (response.pagination) {
            this.pagination = {
              page: response.pagination.page || 1,
              page_size: response.pagination.per_page || 20,
              total: response.pagination.total || 0,
              pages: response.pagination.pages || 0
            }
          }

          return response
        }

        throw new Error(response.error || '获取计划列表失败')
      } catch (error) {
        console.error('获取计划列表失败:', error)
        ElMessage.error(error.message || '获取计划列表失败')
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * 获取计划详情
     *
     * @param {number} id - 计划ID
     * @returns {Promise<Object>} 计划详情
     */
    async fetchPlanDetail(id) {
      this.loading = true
      try {
        const response = await materialPlanApi.getPlanDetail(id)

        if (response.success) {
          this.currentPlan = response.data
          return response.data
        }

        throw new Error(response.error || '获取计划详情失败')
      } catch (error) {
        console.error('获取计划详情失败:', error)
        ElMessage.error(error.message || '获取计划详情失败')
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * 创建计划
     *
     * @param {Object} data - 计划数据
     * @returns {Promise<Object>} 创建的计划
     */
    async createPlan(data) {
      this.loading = true
      try {
        const response = await materialPlanApi.createPlan(data)

        if (response.success) {
          ElMessage.success('计划创建成功')

          // 刷新列表
          await this.fetchPlans()

          return response.data
        }

        throw new Error(response.error || '创建计划失败')
      } catch (error) {
        console.error('创建计划失败:', error)
        ElMessage.error(error.message || '创建计划失败')
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * 更新计划
     *
     * @param {number} id - 计划ID
     * @param {Object} data - 更新数据
     * @returns {Promise<Object>} 更新后的计划
     */
    async updatePlan(id, data) {
      this.loading = true
      try {
        const response = await materialPlanApi.updatePlan(id, data)

        if (response.success) {
          ElMessage.success('计划更新成功')

          // 更新当前计划
          if (this.currentPlan?.id === id) {
            this.currentPlan = response.data
          }

          // 刷新列表
          await this.fetchPlans()

          return response.data
        }

        throw new Error(response.error || '更新计划失败')
      } catch (error) {
        console.error('更新计划失败:', error)
        ElMessage.error(error.message || '更新计划失败')
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * 删除计划
     *
     * @param {number} id - 计划ID
     * @returns {Promise<boolean>} 是否成功
     */
    async deletePlan(id) {
      this.loading = true
      try {
        const response = await materialPlanApi.deletePlan(id)

        if (response.success) {
          ElMessage.success('计划删除成功')

          // 清除当前计划
          if (this.currentPlan?.id === id) {
            this.currentPlan = null
          }

          // 刷新列表
          await this.fetchPlans()

          return true
        }

        throw new Error(response.error || '删除计划失败')
      } catch (error) {
        console.error('删除计划失败:', error)
        ElMessage.error(error.message || '删除计划失败')
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * 提交计划审批
     *
     * @param {number} id - 计划ID
     * @returns {Promise<boolean>} 是否成功
     */
    async submitPlan(id) {
      this.loading = true
      try {
        const response = await materialPlanApi.submitPlan(id)

        console.log('提交审批响应:', response)

        if (response.success) {
          ElMessage.success('计划已提交审批')

          // 更新当前计划状态
          if (this.currentPlan?.id === id) {
            this.currentPlan.status = 'pending'
          }

          // 刷新列表
          await this.fetchPlans()

          console.log('刷新后的计划列表:', this.plans)

          return true
        }

        throw new Error(response.error || '提交审批失败')
      } catch (error) {
        // 如果用户取消了操作（通过 MessageBox），不显示错误
        if (error !== 'cancel' && error.message !== 'cancel') {
          console.error('提交审批失败:', error)
          ElMessage.error(error.message || '提交审批失败')
        }
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * 批准计划
     *
     * @param {number} id - 计划ID
     * @param {Object} data - 审批信息
     * @returns {Promise<boolean>} 是否成功
     */
    async approvePlan(id, data = {}) {
      this.loading = true
      try {
        const response = await materialPlanApi.approvePlan(id, data)

        if (response.success) {
          ElMessage.success('计划已批准')

          // 更新当前计划状态
          if (this.currentPlan?.id === id) {
            this.currentPlan.status = 'approved'
          }

          // 刷新列表
          await this.fetchPlans()

          return true
        }

        throw new Error(response.error || '批准计划失败')
      } catch (error) {
        console.error('批准计划失败:', error)
        ElMessage.error(error.message || '批准计划失败')
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * 拒绝计划
     *
     * @param {number} id - 计划ID
     * @param {Object} data - 拒绝信息
     * @returns {Promise<boolean>} 是否成功
     */
    async rejectPlan(id, data = {}) {
      this.loading = true
      try {
        const response = await materialPlanApi.rejectPlan(id, data)

        if (response.success) {
          ElMessage.success('计划已拒绝')

          // 更新当前计划状态
          if (this.currentPlan?.id === id) {
            this.currentPlan.status = 'rejected'
          }

          // 刷新列表
          await this.fetchPlans()

          return true
        }

        throw new Error(response.error || '拒绝计划失败')
      } catch (error) {
        console.error('拒绝计划失败:', error)
        ElMessage.error(error.message || '拒绝计划失败')
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * 激活计划
     *
     * @param {number} id - 计划ID
     * @returns {Promise<boolean>} 是否成功
     */
    async activatePlan(id) {
      this.loading = true
      try {
        const response = await materialPlanApi.activatePlan(id)

        if (response.success) {
          ElMessage.success('计划已激活')

          // 更新当前计划状态
          if (this.currentPlan?.id === id) {
            this.currentPlan.status = 'active'
          }

          // 刷新列表
          await this.fetchPlans()

          return true
        }

        throw new Error(response.error || '激活计划失败')
      } catch (error) {
        console.error('激活计划失败:', error)
        ElMessage.error(error.message || '激活计划失败')
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * 重新提交计划
     *
     * @param {number} id - 计划ID
     * @returns {Promise<boolean>} 是否成功
     */
    async resubmitPlan(id) {
      this.loading = true
      try {
        const response = await materialPlanApi.resubmitPlan(id)

        if (response.success) {
          ElMessage.success('计划已重新提交')

          // 更新当前计划状态
          if (this.currentPlan?.id === id) {
            this.currentPlan.status = 'pending'
          }

          // 刷新列表
          await this.fetchPlans()

          return true
        }

        throw new Error(response.error || '重新提交失败')
      } catch (error) {
        console.error('重新提交失败:', error)
        ElMessage.error(error.message || '重新提交失败')
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * 取消计划
     *
     * @param {number} id - 计划ID
     * @param {Object} data - 取消信息
     * @returns {Promise<boolean>} 是否成功
     */
    async cancelPlan(id, data = {}) {
      this.loading = true
      try {
        const response = await materialPlanApi.cancelPlan(id, data)

        if (response.success) {
          ElMessage.success('计划已取消')

          // 更新当前计划状态
          if (this.currentPlan?.id === id) {
            this.currentPlan.status = 'cancelled'
          }

          // 刷新列表
          await this.fetchPlans()

          return true
        }

        throw new Error(response.error || '取消计划失败')
      } catch (error) {
        console.error('取消计划失败:', error)
        ElMessage.error(error.message || '取消计划失败')
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * 获取计划项目列表
     *
     * @param {number} id - 计划ID
     * @returns {Promise<Array>} 项目列表
     */
    async fetchPlanItems(id) {
      try {
        const response = await materialPlanApi.getPlanItems(id)

        if (response.success) {
          return response.data || []
        }

        throw new Error(response.error || '获取项目列表失败')
      } catch (error) {
        console.error('获取项目列表失败:', error)
        ElMessage.error(error.message || '获取项目列表失败')
        throw error
      }
    },

    /**
     * 添加计划项目
     *
     * @param {number} planId - 计划ID
     * @param {Object} data - 项目数据
     * @returns {Promise<Object>} 添加的项目
     */
    async addPlanItem(planId, data) {
      try {
        const response = await materialPlanApi.addPlanItem(planId, data)

        if (response.success) {
          ElMessage.success('项目添加成功')

          // 刷新当前计划
          if (this.currentPlan?.id === planId) {
            await this.fetchPlanDetail(planId)
          }

          return response.data
        }

        throw new Error(response.error || '添加项目失败')
      } catch (error) {
        console.error('添加项目失败:', error)
        ElMessage.error(error.message || '添加项目失败')
        throw error
      }
    },

    /**
     * 更新计划项目
     *
     * @param {number} itemId - 项目ID
     * @param {Object} data - 更新数据
     * @returns {Promise<Object>} 更新后的项目
     */
    async updatePlanItem(itemId, data) {
      try {
        const response = await materialPlanApi.updatePlanItem(itemId, data)

        if (response.success) {
          ElMessage.success('项目更新成功')

          // 刷新当前计划
          if (this.currentPlan) {
            await this.fetchPlanDetail(this.currentPlan.id)
          }

          return response.data
        }

        throw new Error(response.error || '更新项目失败')
      } catch (error) {
        console.error('更新项目失败:', error)
        ElMessage.error(error.message || '更新项目失败')
        throw error
      }
    },

    /**
     * 删除计划项目
     *
     * @param {number} itemId - 项目ID
     * @returns {Promise<boolean>} 是否成功
     */
    async deletePlanItem(itemId) {
      try {
        const response = await materialPlanApi.deletePlanItem(itemId)

        if (response.success) {
          ElMessage.success('项目删除成功')

          // 刷新当前计划
          if (this.currentPlan) {
            await this.fetchPlanDetail(this.currentPlan.id)
          }

          return true
        }

        throw new Error(response.error || '删除项目失败')
      } catch (error) {
        console.error('删除项目失败:', error)
        ElMessage.error(error.message || '删除项目失败')
        throw error
      }
    },

    /**
     * 获取统计概览
     *
     * @returns {Promise<Object>} 统计数据
     */
    async fetchStatistics() {
      try {
        const response = await materialPlanApi.getStatistics()

        if (response.success) {
          this.statistics = response.data
          return response.data
        }

        throw new Error(response.error || '获取统计数据失败')
      } catch (error) {
        console.error('获取统计数据失败:', error)
        throw error
      }
    },

    /**
     * 获取计划详细统计
     *
     * @param {number} id - 计划ID
     * @returns {Promise<Object>} 统计数据
     */
    async fetchPlanStatistics(id) {
      try {
        const response = await materialPlanApi.getPlanStatistics(id)

        if (response.success) {
          this.currentPlanStatistics = response.data
          return response.data
        }

        throw new Error(response.error || '获取计划统计失败')
      } catch (error) {
        console.error('获取计划统计失败:', error)
        throw error
      }
    },

    /**
     * 设置筛选条件
     *
     * @param {Object} filters - 筛选条件
     */
    setFilters(filters) {
      this.filters = { ...this.filters, ...filters }
      // 重置到第一页
      this.pagination.page = 1
    },

    /**
     * 重置筛选条件
     */
    resetFilters() {
      this.filters = {
        status: '',
        plan_type: '',
        priority: '',
        search: '',
        project_id: null,
        start_date: '',
        end_date: ''
      }
      this.pagination.page = 1
    },

    /**
     * 设置分页
     *
     * @param {Object} pagination - 分页信息
     */
    setPagination(pagination) {
      this.pagination = { ...this.pagination, ...pagination }
    },

    /**
     * 清除当前计划
     */
    clearCurrentPlan() {
      this.currentPlan = null
      this.currentPlanStatistics = null
    }
  }
})
