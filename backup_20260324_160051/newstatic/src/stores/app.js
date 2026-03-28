/**
 * 应用全局状态管理 Store
 *
 * 本 Store 管理应用级别的全局状态，包括：
 * - 侧边栏折叠状态
 * - 设备类型检测（桌面/移动）
 * - 全局加载状态
 * - 系统名称配置
 *
 * @module AppStore
 * @author Material Management System
 * @date 2025-01-27
 */

import { defineStore } from 'pinia'
import { systemApi } from '@/api'

/**
 * 应用全局状态管理 Store
 *
 * 状态设计说明：
 * - sidebarCollapsed: 控制侧边栏的展开/收起状态
 * - device: 当前设备类型，影响布局响应式行为
 * - loading: 全局加载状态，用于显示加载动画
 * - systemName: 系统名称，可动态配置
 */
export const useAppStore = defineStore('app', {
  /**
   * 状态定义
   *
   * 这些状态不需要持久化到 localStorage，
   * 因为它们会在每次页面加载时重置为初始值
   */
  state: () => ({
    /**
     * 侧边栏折叠状态
     *
     * - false: 侧边栏展开（默认状态，桌面端）
     * - true: 侧边栏收起（移动端默认，或用户手动折叠）
     *
     * @type {boolean}
     */
    sidebarCollapsed: false,

    /**
     * 设备类型
     *
     * 可能的值：
     * - 'desktop': 桌面设备（屏幕宽度 ≥ 992px）
     * - 'tablet': 平板设备（屏幕宽度 768px - 991px）
     * - 'mobile': 移动设备（屏幕宽度 < 768px）
     *
     * @type {string}
     * @default 'desktop'
     */
    device: 'desktop',

    /**
     * 全局加载状态
     *
     * 用于控制全屏加载动画的显示/隐藏
     * 适用于需要长时间等待的操作（如大数据量导出）
     *
     * @type {boolean}
     * @default false
     */
    loading: false,

    /**
     * 系统名称
     *
     * 显示在浏览器标签页、登录页等位置
     * 可以通过后端配置动态修改
     *
     * @type {string}
     * @default '材料管理系统'
     */
    systemName: '材料管理系统'
  }),

  /**
   * Getters - 计算属性
   *
   * 从 state 派生出有用的值
   */
  getters: {
    /**
     * 检查是否为移动设备
     *
     * 用于响应式布局判断，移动设备需要特殊处理：
     * - 侧边栏默认收起
     * - 表格显示简化版
     * - 触摸操作优化
     *
     * @param {Object} state - store state
     * @returns {boolean} 是否为移动设备
     */
    isMobile: (state) => state.device === 'mobile'
  },

  /**
   * Actions - 业务逻辑方法
   *
   * 用于修改状态和执行业务逻辑
   */
  actions: {
    /**
     * 切换侧边栏折叠状态
     *
     * 用途：
     * - 响应用户点击侧边栏切换按钮
     * - 桌面端：展开/收起切换
     * - 移动端：显示/隐藏侧边栏遮罩
     *
     * @example
     * // 在布局组件的侧边栏切换按钮上调用
     * <el-button @click="appStore.toggleSidebar()">
     *   <i :class="appStore.sidebarCollapsed ? 'el-icon-s-unfold' : 'el-icon-s-fold'" />
     * </el-button>
     */
    toggleSidebar() {
      this.sidebarCollapsed = !this.sidebarCollapsed
    },

    /**
     * 设置设备类型
     *
     * 用途：
     * - 响应窗口大小变化事件
     * - 根据屏幕宽度自动调整布局
     * - 移动端自动折叠侧边栏
     *
     * @param {string} device - 设备类型 ('desktop' | 'tablet' | 'mobile')
     *
     * @example
     * // 在 App.vue 中监听窗口大小变化
     * window.addEventListener('resize', () => {
     *   const width = window.innerWidth
     *   let device = 'desktop'
     *   if (width < 768) device = 'mobile'
     *   else if (width < 992) device = 'tablet'
     *   appStore.setDevice(device)
     * })
     */
    setDevice(device) {
      this.device = device

      // 移动设备自动折叠侧边栏，节省屏幕空间
      if (device === 'mobile') {
        this.sidebarCollapsed = true
      }
    },

    /**
     * 设置全局加载状态
     *
     * 用途：
     * - 在耗时操作前显示加载动画
     * - 操作完成后隐藏加载动画
     *
     * @param {boolean} loading - 加载状态
     *
     * @example
     * // 在异步操作中使用
     * async function exportData() {
     *   appStore.setLoading(true)
     *   try {
     *     await api.exportLargeData()
     *   } finally {
     *     appStore.setLoading(false)
     *   }
     * }
     */
    setLoading(loading) {
      this.loading = loading
    },

    /**
     * 设置系统名称
     *
     * 用途：
     * - 从后端配置中获取系统名称
     * - 动态更新系统名称（如多租户场景）
     *
     * @param {string} name - 新的系统名称
     *
     * @example
     * // 在应用初始化时从后端获取配置
     * const config = await api.getSystemConfig()
     * appStore.setSystemName(config.system_name || '材料管理系统')
     */
    setSystemName(name) {
      this.systemName = name
    },

    /**
     * 从后端加载系统配置
     *
     * 用途：
     * - 应用初始化时获取系统名称等配置
     * - 动态更新系统显示信息
     *
     * @example
     * // 在 App.vue 的 onMounted 中调用
     * onMounted(async () => {
     *   await appStore.loadSystemConfig()
     * })
     */
    async loadSystemConfig() {
      try {
        const { data } = await systemApi.getPublicSettings()
        if (data && data.system_name) {
          this.systemName = data.system_name
        }
      } catch (error) {
        console.error('加载系统配置失败:', error)
        // 保持默认的系统名称
      }
    }
  }
})
