/**
 * 甘特图拖拽逻辑 Composable
 * 处理任务条的拖拽和边缘调整功能
 * 优化版本：包含节流和性能优化
 */

import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { formatDate, addDays, diffDays } from '@/utils/dateFormat'

/**
 * 节流函数 - 限制函数执行频率
 * @param {Function} fn - 需要节流的函数
 * @param {number} delay - 延迟时间（毫秒）
 * @returns {Function} 节流后的函数
 */
function throttle(fn, delay) {
  let lastTime = 0
  let timer = null

  return function (...args) {
    const now = Date.now()
    const remaining = delay - (now - lastTime)

    if (remaining <= 0) {
      if (timer) {
        clearTimeout(timer)
        timer = null
      }
      lastTime = now
      fn.apply(this, args)
    } else if (!timer) {
      timer = setTimeout(() => {
        lastTime = Date.now()
        timer = null
        fn.apply(this, args)
      }, remaining)
    }
  }
}

/**
 * 防抖函数 - 延迟执行函数
 * @param {Function} fn - 需要防抖的函数
 * @param {number} delay - 延迟时间（毫秒）
 * @returns {Function} 防抖后的函数
 */
function debounce(fn, delay) {
  let timer = null

  return function (...args) {
    if (timer) {
      clearTimeout(timer)
    }
    timer = setTimeout(() => {
      fn.apply(this, args)
      timer = null
    }, delay)
  }
}

/**
 * 请求动画帧节流 - 专门用于动画的高性能节流
 * 目标 60fps (~16.67ms)
 */
function rafThrottle(fn) {
  let rafId = null

  return function (...args) {
    if (rafId === null) {
      rafId = requestAnimationFrame(() => {
        fn.apply(this, args)
        rafId = null
      })
    }
  }
}

/**
 * 拖拽模式枚举
 */
export const DragMode = {
  NONE: 'none',
  MOVE: 'move',           // 移动整个任务条
  RESIZE_LEFT: 'resize_left',  // 调整左边缘（开始时间）
  RESIZE_RIGHT: 'resize_right' // 调整右边缘（结束时间）
}

/**
 * 使用甘特图拖拽功能
 * @param {Object} options - 配置选项
 * @returns {Object} 拖拽相关的状态和方法
 */
export function useGanttDrag(options = {}) {
  const {
    dayWidth = ref(40),
    timelineDays = ref([]),
    onDragEnd = null,
    onDragChange = null,
    // 新增：性能配置
    throttleMs = 16,  // ~60fps
    enableRAF = true  // 是否启用 requestAnimationFrame 优化
  } = options

  // 拖拽状态
  const isDragging = ref(false)
  const dragMode = ref(DragMode.NONE)
  const draggedTask = ref(null)
  const originalTask = ref(null)

  // 鼠标位置
  const mouseX = ref(0)
  const startMouseX = ref(0)

  // 拖拽预览
  const previewTask = ref(null)

  // 拖拽提示框位置
  const tooltipPosition = ref({ x: 0, y: 0 })
  const tooltipVisible = ref(false)

  // 边缘检测区域宽度（像素）
  const EDGE_RESIZE_THRESHOLD = 10

  // 性能优化：缓存计算结果
  let cachedDayOffset = 0
  let lastClientX = 0

  /**
   * 检测鼠标是否在任务条边缘
   * @param {MouseEvent} event - 鼠标事件
   * @param {HTMLElement} taskBar - 任务条DOM元素
   * @returns {string} 拖拽模式
   */
  function detectDragMode(event, taskBar) {
    if (!taskBar) return DragMode.NONE

    const rect = taskBar.getBoundingClientRect()
    const mouseX = event.clientX

    // 检查是否在左边缘
    if (mouseX <= rect.left + EDGE_RESIZE_THRESHOLD) {
      return DragMode.RESIZE_LEFT
    }

    // 检查是否在右边缘
    if (mouseX >= rect.right - EDGE_RESIZE_THRESHOLD) {
      return DragMode.RESIZE_RIGHT
    }

    // 在中间区域，可以移动
    return DragMode.MOVE
  }

  /**
   * 开始拖拽
   * @param {MouseEvent} event - 鼠标事件
   * @param {Object} task - 任务对象
   * @param {HTMLElement} taskBar - 任务条DOM元素
   */
  function startDrag(event, task, taskBar) {
    event.preventDefault()

    dragMode.value = detectDragMode(event, taskBar)

    if (dragMode.value === DragMode.NONE) {
      return
    }

    isDragging.value = true
    draggedTask.value = task
    originalTask.value = { ...task }

    startMouseX.value = event.clientX
    lastClientX = event.clientX
    mouseX.value = event.clientX

    // 初始化预览任务
    previewTask.value = {
      ...task,
      start: task.start,
      end: task.end
    }

    // 显示提示框
    updateTooltip(event)

    // 创建节流版本的鼠标移动处理
    const throttledMouseMove = enableRAF
      ? rafThrottle(handleMouseMoveCore)
      : throttle(handleMouseMoveCore, throttleMs)

    // 添加全局事件监听
    document.addEventListener('mousemove', throttledMouseMove)
    // 保存引用以便后续清理
    event.target._throttledMouseMove = throttledMouseMove
    document.addEventListener('mouseup', handleMouseUp)

    // 改变光标样式
    document.body.style.cursor = getCursorForMode(dragMode.value)
  }

  /**
   * 处理鼠标移动核心逻辑（不含节流）
   * @param {MouseEvent} event - 鼠标事件
   */
  function handleMouseMoveCore(event) {
    if (!isDragging.value) return

    mouseX.value = event.clientX

    // 计算拖拽的像素距离
    const deltaX = event.clientX - startMouseX.value

    // 转换为天数偏移
    const dayOffset = Math.round(deltaX / dayWidth.value)

    // 性能优化：只在偏移量变化时更新
    if (dayOffset !== cachedDayOffset) {
      cachedDayOffset = dayOffset

      // 根据拖拽模式更新预览
      updatePreviewTask(dayOffset)

      // 触发变更回调
      if (onDragChange && previewTask.value) {
        onDragChange(previewTask.value)
      }
    }

    // 更新提示框位置（使用 RAF 节流）
    updateTooltipRAF(event)
  }

  /**
   * RAF 节流的提示框更新
   */
  let rafId = null
  function updateTooltipRAF(event) {
    if (rafId !== null) {
      return
    }

    rafId = requestAnimationFrame(() => {
      if (previewTask.value) {
        tooltipPosition.value = {
          x: event.clientX + 15,
          y: event.clientY + 15
        }
        tooltipVisible.value = true
      }
      rafId = null
    })
  }

  /**
   * 处理鼠标松开
   * @param {MouseEvent} event - 鼠标事件
   */
  function handleMouseUp(event) {
    if (!isDragging.value) return

    // 获取节流函数引用并清理
    const throttledMouseMove = event.target?._throttledMouseMove

    // 移除全局事件监听
    if (throttledMouseMove) {
      document.removeEventListener('mousemove', throttledMouseMove)
      delete event.target._throttledMouseMove
    }
    document.removeEventListener('mouseup', handleMouseUp)

    // 恢复光标样式
    document.body.style.cursor = ''

    // 隐藏提示框
    tooltipVisible.value = false

    // 清理 RAF
    if (rafId !== null) {
      cancelAnimationFrame(rafId)
      rafId = null
    }

    // 检查是否有实际变化
    if (hasTaskChanged(previewTask.value, originalTask.value)) {
      // 触发拖拽结束回调
      if (onDragEnd) {
        onDragEnd(previewTask.value, originalTask.value)
      }
    }

    // 重置状态
    isDragging.value = false
    dragMode.value = DragMode.NONE
    draggedTask.value = null
    originalTask.value = null
    previewTask.value = null
    cachedDayOffset = 0
    lastClientX = 0
  }

  /**
   * 根据拖拽模式更新预览任务
   * @param {number} dayOffset - 天数偏移
   */
  function updatePreviewTask(dayOffset) {
    if (!originalTask.value) return

    const original = originalTask.value

    switch (dragMode.value) {
      case DragMode.MOVE:
        // 移动整个任务条
        previewTask.value = {
          ...original,
          start: formatDate(addDays(original.start, dayOffset)),
          end: formatDate(addDays(original.end, dayOffset))
        }
        break

      case DragMode.RESIZE_LEFT:
        // 调整左边缘（开始时间）
        const newStart = addDays(original.start, dayOffset)
        // 确保开始日期不晚于结束日期
        if (newStart <= new Date(original.end)) {
          previewTask.value = {
            ...original,
            start: formatDate(newStart)
          }
        }
        break

      case DragMode.RESIZE_RIGHT:
        // 调整右边缘（结束时间）
        const newEnd = addDays(original.end, dayOffset)
        // 确保结束日期不早于开始日期
        if (newEnd >= new Date(original.start)) {
          previewTask.value = {
            ...original,
            end: formatDate(newEnd)
          }
        }
        break
    }

    // 更新工期
    if (previewTask.value) {
      const duration = diffDays(previewTask.value.start, previewTask.value.end)
      previewTask.value.duration = Math.max(duration, 0)
    }
  }

  /**
   * 更新提示框位置和内容
   * @param {MouseEvent} event - 鼠标事件
   */
  function updateTooltip(event) {
    if (!previewTask.value) return

    tooltipPosition.value = {
      x: event.clientX + 15,
      y: event.clientY + 15
    }

    tooltipVisible.value = true
  }

  /**
   * 检查任务是否有变化
   * @param {Object} newTask - 新任务
   * @param {Object} oldTask - 旧任务
   * @returns {boolean} 是否有变化
   */
  function hasTaskChanged(newTask, oldTask) {
    if (!newTask || !oldTask) return false
    return newTask.start !== oldTask.start || newTask.end !== oldTask.end
  }

  /**
   * 根据拖拽模式获取光标样式
   * @param {string} mode - 拖拽模式
   * @returns {string} 光标样式
   */
  function getCursorForMode(mode) {
    switch (mode) {
      case DragMode.MOVE:
        return 'grab'
      case DragMode.RESIZE_LEFT:
        return 'w-resize'
      case DragMode.RESIZE_RIGHT:
        return 'e-resize'
      default:
        return 'default'
    }
  }

  /**
   * 获取任务条的样式类（包含拖拽状态）
   * @param {Object} task - 任务对象
   * @returns {string} 样式类
   */
  function getTaskBarClass(task) {
    const classes = []

    if (isDragging.value && draggedTask.value?.id === task.id) {
      classes.push('is-dragging')
    }

    if (dragMode.value === DragMode.RESIZE_LEFT) {
      classes.push('resize-left')
    }

    if (dragMode.value === DragMode.RESIZE_RIGHT) {
      classes.push('resize-right')
    }

    return classes.join(' ')
  }

  /**
   * 获取任务条的样式（包含拖拽预览）
   * @param {Object} task - 任务对象
   * @returns {Object} 样式对象
   */
  function getTaskBarStyle(task) {
    if (isDragging.value && draggedTask.value?.id === task.id && previewTask.value) {
      // 返回预览任务的位置
      const timelineStart = timelineDays.value[0]?.date
      if (!timelineStart) return {}

      const taskStart = new Date(previewTask.value.start)
      const taskEnd = new Date(previewTask.value.end)

      const daysDiff = Math.ceil((taskStart - new Date(timelineStart)) / (1000 * 60 * 60 * 24))
      const duration = Math.ceil((taskEnd - taskStart) / (1000 * 60 * 60 * 24))

      const left = daysDiff * dayWidth.value
      const width = duration * dayWidth.value

      return {
        left: left + 'px',
        width: Math.max(width - 4, 10) + 'px',
        opacity: '0.8'
      }
    }

    return {}
  }

  /**
   * 取消拖拽
   */
  function cancelDrag() {
    if (isDragging.value) {
      // 获取节流函数引用并清理
      const throttledMouseMove = document.body?._throttledMouseMove

      // 移除全局事件监听
      if (throttledMouseMove) {
        document.removeEventListener('mousemove', throttledMouseMove)
        delete document.body._throttledMouseMove
      }
      document.removeEventListener('mouseup', handleMouseUp)

      // 恢复光标样式
      document.body.style.cursor = ''

      // 清理 RAF
      if (rafId !== null) {
        cancelAnimationFrame(rafId)
        rafId = null
      }

      // 重置状态
      isDragging.value = false
      dragMode.value = DragMode.NONE
      draggedTask.value = null
      originalTask.value = null
      previewTask.value = null
      tooltipVisible.value = false
      cachedDayOffset = 0
      lastClientX = 0
    }
  }

  /**
   * 计算提示文本
   * @returns {string} 提示文本
   */
  const tooltipText = computed(() => {
    if (!previewTask.value) return ''

    const modeText = {
      [DragMode.MOVE]: '移动',
      [DragMode.RESIZE_LEFT]: '调整开始',
      [DragMode.RESIZE_RIGHT]: '调整结束'
    }

    const mode = modeText[dragMode.value] || ''
    const duration = diffDays(previewTask.value.start, previewTask.value.end)

    return `${mode}: ${previewTask.value.start} ~ ${previewTask.value.end} (${duration}天)`
  })

  /**
   * 是否可以调整左边缘
   */
  const canResizeLeft = computed(() => {
    return dragMode.value === DragMode.RESIZE_LEFT || !isDragging.value
  })

  /**
   * 是否可以调整右边缘
   */
  const canResizeRight = computed(() => {
    return dragMode.value === DragMode.RESIZE_RIGHT || !isDragging.value
  })

  return {
    // 状态
    isDragging,
    dragMode,
    draggedTask,
    previewTask,
    tooltipPosition,
    tooltipVisible,
    tooltipText,
    canResizeLeft,
    canResizeRight,

    // 方法
    startDrag,
    cancelDrag,
    getTaskBarClass,
    getTaskBarStyle,
    detectDragMode
  }
}
