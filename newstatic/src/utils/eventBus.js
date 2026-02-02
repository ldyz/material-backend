/**
 * 轻量级事件总线
 * 用于组件之间的消息同步
 */

class EventBus {
  constructor() {
    this.events = {}
  }

  /**
   * 订阅事件
   * @param {string} event - 事件名称
   * @param {Function} callback - 回调函数
   * @returns {Function} 取消订阅的函数
   */
  on(event, callback) {
    if (!this.events[event]) {
      this.events[event] = []
    }
    this.events[event].push(callback)

    // 返回取消订阅的函数
    return () => {
      this.events[event] = this.events[event].filter(cb => cb !== callback)
    }
  }

  /**
   * 订阅事件（只执行一次）
   * @param {string} event - 事件名称
   * @param {Function} callback - 回调函数
   * @returns {Function} 取消订阅的函数
   */
  once(event, callback) {
    const onceCallback = (...args) => {
      callback(...args)
      this.off(event, onceCallback)
    }
    return this.on(event, onceCallback)
  }

  /**
   * 取消订阅事件
   * @param {string} event - 事件名称
   * @param {Function} callback - 回调函数
   */
  off(event, callback) {
    if (!this.events[event]) return
    this.events[event] = this.events[event].filter(cb => cb !== callback)
  }

  /**
   * 触发事件
   * @param {string} event - 事件名称
   * @param {...any} args - 传递给回调函数的参数
   */
  emit(event, ...args) {
    if (!this.events[event]) return
    this.events[event].forEach(callback => {
      try {
        callback(...args)
      } catch (error) {
        console.error(`[EventBus] Error in event "${event}":`, error)
      }
    })
  }

  /**
   * 清除所有事件监听器
   */
  clear() {
    this.events = {}
  }

  /**
   * 清除指定事件的所有监听器
   * @param {string} event - 事件名称
   */
  clearEvent(event) {
    delete this.events[event]
  }

  /**
   * 获取指定事件的监听器数量
   * @param {string} event - 事件名称
   * @returns {number} 监听器数量
   */
  listenerCount(event) {
    return this.events[event]?.length || 0
  }
}

// 创建全局事件总线实例
const eventBus = new EventBus()

// 甘特图相关的事件名称常量
export const GanttEvents = {
  // ==================== 数据相关 ====================
  DATA_LOADED: 'gantt:data:loaded',
  DATA_REFRESHED: 'gantt:data:refreshed',
  DATA_CHANGED: 'gantt:data:changed',
  DATA_SAVED: 'gantt:data:saved',
  DATA_SAVE_ERROR: 'gantt:data:save-error',

  // ==================== 任务相关 ====================
  TASK_SELECTED: 'gantt:task:selected',
  TASK_HOVERED: 'gantt:task:hovered',
  TASK_DRAG_START: 'gantt:task:drag-start',
  TASK_DRAG_MOVE: 'gantt:task:drag-move',
  TASK_DRAG_END: 'gantt:task:drag-end',
  TASK_REORDERED: 'gantt:task:reordered',
  TASK_CREATED: 'gantt:task:created',
  TASK_UPDATED: 'gantt:task:updated',
  TASK_DELETED: 'gantt:task:deleted',
  TASK_DRAGGED: 'gantt:task:dragged', // 保留用于向后兼容

  // ==================== 依赖关系相关 ====================
  DEPENDENCY_CREATING: 'gantt:dependency:creating',
  DEPENDENCY_CREATED: 'gantt:dependency:created',
  DEPENDENCY_DELETED: 'gantt:dependency:deleted',
  DEPENDENCY_UPDATED: 'gantt:dependency:updated',

  // ==================== 视图相关 ====================
  VIEW_CHANGED: 'gantt:view:changed',
  ZOOM_CHANGED: 'gantt:zoom:changed',
  VIEW_AUTO_FITTED: 'gantt:view:auto-fitted',

  // ==================== 编辑相关 ====================
  EDIT_START: 'gantt:edit:start',
  EDIT_CANCEL: 'gantt:edit:cancel',
  EDIT_SAVE: 'gantt:edit:save',

  // ==================== 资源相关 ====================
  RESOURCE_ALLOCATED: 'gantt:resource:allocated',
  RESOURCE_UPDATED: 'gantt:resource:updated',
  RESOURCE_DELETED: 'gantt:resource:deleted',

  // ==================== UI状态 ====================
  CONTEXT_MENU_SHOW: 'gantt:context-menu:show',
  CONTEXT_MENU_HIDE: 'gantt:context-menu:hide',
  DIALOG_OPEN: 'gantt:dialog:open',
  DIALOG_CLOSE: 'gantt:dialog:close',

  // ==================== 状态相关 (保留用于向后兼容) ====================
  STATUS_CHANGED: 'gantt:status:changed',
  PROGRESS_UPDATED: 'gantt:progress:updated'
}

export default eventBus
