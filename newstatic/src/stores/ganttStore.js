/**
 * 甘特图集中式状态管理 Store
 * 统一管理甘特图的所有状态和操作
 */

import { reactive, computed, watch } from 'vue'
import { progressApi } from '@/api'
import eventBus, { GanttEvents } from '@/utils/eventBus'
import { formatDate, addDays, diffDays } from '@/utils/dateFormat'
import {
  filterDependencies,
  calculateAllDependencyPaths,
  isMilestone,
  groupTasksByStatus,
  groupTasksByPriority,
  getWeekNumber,
  getQuarter,
  getWeekStart,
  getWeekEnd
} from '@/utils/ganttHelpers'

// ==================== 状态定义 ====================
const state = reactive({
  // 项目信息
  projectId: null,
  projectName: '',

  // 原始数据
  scheduleData: { activities: {}, nodes: {} },

  // 格式化后的任务列表
  tasks: [],
  filteredTasks: [],

  // 本地 parent_id 映射（用于跟踪拖拽更改）
  localParentIdMap: new Map(),
  // 数据版本号，用于检测是否是完全重新加载
  dataVersion: 0,

  // 视图状态
  viewMode: 'day', // day, week, month, quarter
  groupMode: '', // '', status, priority
  collapsedGroups: new Set(),
  collapsedTasks: new Set(), // 折叠的任务ID集合（用于树形结构）

  // 缩放配置
  dayWidth: 40,
  rowHeight: 60,
  taskHeight: 32,
  VIEW_CONFIG: {
    day: { minWidth: 40, maxWidth: 80, default: 40, label: '日' },
    week: { minWidth: 20, maxWidth: 40, default: 28, label: '周' },
    month: { minWidth: 8, maxWidth: 20, default: 12, label: '月' },
    quarter: { minWidth: 3, maxWidth: 8, default: 5, label: '季' }
  },

  // 显示选项
  showDependencies: true,
  showCriticalPath: true,
  showBaseline: false,

  // 搜索和筛选
  searchKeyword: '',

  // 选择状态
  selectedTaskId: null,

  // 拖拽状态（时间轴上的任务拖拽）
  isDragging: false,
  dragMode: 'none', // none, move, resize_left, resize_right
  draggedTask: null,
  originalTask: null,
  previewTask: null,
  tooltipPosition: { x: 0, y: 0 },
  tooltipVisible: false,

  // 任务表格拖拽状态（改变层级和顺序）
  tableDraggedTask: null,

  // 依赖关系创建状态
  isCreatingDependency: false,
  dependencySourceTask: null,
  tempLineEnd: null,

  // 编辑状态
  editingTask: null,
  taskDetailVisible: false,

  // 对话框状态
  editDialogVisible: false,
  resourceDialogVisible: false,
  resourceManagementDialogVisible: false,
  currentTaskForResource: null,

  // 右键菜单
  contextMenuVisible: false,
  contextMenuTask: null,
  contextMenuPosition: { x: 0, y: 0 },

  // 全屏状态
  isFullscreen: false,

  // 容器调整大小
  containerSize: { width: null, height: null },
  isResizing: false,

  // 资源库
  resourceLibrary: [],

  // 未保存更改
  hasUnsavedChanges: false,
  pendingTaskUpdates: new Map(), // Map<taskId, updateData>
  pendingDependencyCreations: [],
  pendingDependencyDeletions: [],

  // 加载状态
  loading: false,
  saving: false,

  // 自动保存计时器
  autoSaveTimer: null,
  AUTO_SAVE_INTERVAL: 30000 // 30秒
})

// ==================== Getters (计算属性) ====================
const getters = {
  // 获取当前选中的任务
  selectedTask: computed(() => {
    return state.tasks.find(t => t.id === state.selectedTaskId) || null
  }),

  // 获取任务统计
  taskStats: computed(() => {
    const tasks = state.filteredTasks
    return {
      total: tasks.length,
      completed: tasks.filter(t => t.status === 'completed').length,
      inProgress: tasks.filter(t => t.status === 'in_progress').length,
      notStarted: tasks.filter(t => t.status === 'not_started').length,
      delayed: tasks.filter(t => t.status === 'delayed').length,
      critical: tasks.filter(t => t.is_critical).length,
      progressRate: tasks.length > 0
        ? Math.round(tasks.reduce((sum, t) => sum + (t.progress || 0), 0) / tasks.length)
        : 0
    }
  }),

  // 时间轴数据
  timelineDays: computed(() => {
    if (state.filteredTasks.length === 0) return []

    const tasks = state.filteredTasks
    const minDate = new Date(Math.min(...tasks.map(t => t.startDate.getTime())))
    const maxDate = new Date(Math.max(...tasks.map(t => t.endDate.getTime())))

    const totalDays = Math.ceil((maxDate - minDate) / (1000 * 60 * 60 * 24))

    // 根据任务时间跨度动态调整缓冲区
    let bufferDays = 7
    if (totalDays <= 14) {
      bufferDays = 14
    } else if (totalDays <= 30) {
      bufferDays = 7
    } else if (totalDays <= 90) {
      bufferDays = 7
    } else {
      bufferDays = 3
    }

    const startDate = new Date(minDate)
    startDate.setDate(startDate.getDate() - bufferDays)

    const endDate = new Date(maxDate)
    endDate.setDate(endDate.getDate() + bufferDays)

    const days = []
    const currentDate = new Date(startDate)
    const today = new Date()

    while (currentDate <= endDate) {
      const dateStr = formatDate(currentDate)
      days.push({
        date: dateStr,
        day: currentDate.getDate(),
        weekday: ['日', '一', '二', '三', '四', '五', '六'][currentDate.getDay()],
        isToday: dateStr === formatDate(today),
        isWeekend: currentDate.getDay() === 0 || currentDate.getDay() === 6,
        position: 0
      })
      currentDate.setDate(currentDate.getDate() + 1)
    }

    return days.map((day, index) => ({
      ...day,
      position: index * state.dayWidth
    }))
  }),

  timelineWeeks: computed(() => {
    if (state.filteredTasks.length === 0) return []

    const tasks = state.filteredTasks
    const minDate = new Date(Math.min(...tasks.map(t => t.startDate.getTime())))
    const maxDate = new Date(Math.max(...tasks.map(t => t.endDate.getTime())))

    const weeks = []
    const currentWeekStart = getWeekStart(minDate)
    const today = new Date()
    let position = 0

    while (currentWeekStart <= maxDate) {
      const weekEnd = getWeekEnd(currentWeekStart)
      const weekNumber = getWeekNumber(currentWeekStart)
      const startStr = formatDate(currentWeekStart)
      const endStr = formatDate(weekEnd)

      weeks.push({
        key: `week-${startStr}`,
        weekNumber,
        start: startStr,
        end: endStr,
        position,
        width: 7 * state.dayWidth,
        isCurrent: today >= currentWeekStart && today <= weekEnd
      })

      position += 7 * state.dayWidth
      currentWeekStart.setDate(currentWeekStart.getDate() + 7)
    }

    return weeks
  }),

  timelineMonths: computed(() => {
    if (state.filteredTasks.length === 0) return []

    const tasks = state.filteredTasks
    const minDate = new Date(Math.min(...tasks.map(t => t.startDate.getTime())))
    const maxDate = new Date(Math.max(...tasks.map(t => t.endDate.getTime())))

    const months = []
    const currentDate = new Date(minDate.getFullYear(), minDate.getMonth(), 1)
    let position = 0

    while (currentDate <= maxDate) {
      const year = currentDate.getFullYear()
      const month = currentDate.getMonth() + 1
      const daysInMonth = new Date(year, month, 0).getDate()
      const width = daysInMonth * state.dayWidth

      months.push({
        key: `${year}-${month}`,
        year,
        month,
        dayCount: daysInMonth,
        position,
        width
      })

      position += width
      currentDate.setMonth(currentDate.getMonth() + 1)
    }

    return months
  }),

  timelineQuarters: computed(() => {
    if (state.filteredTasks.length === 0) return []

    const tasks = state.filteredTasks
    const minDate = new Date(Math.min(...tasks.map(t => t.startDate.getTime())))
    const maxDate = new Date(Math.max(...tasks.map(t => t.endDate.getTime())))

    const quarters = []
    const currentYear = minDate.getFullYear()
    const currentQuarter = getQuarter(minDate)
    let position = 0
    let year = currentYear
    let quarter = currentQuarter

    while (true) {
      const quarterStart = new Date(year, (quarter - 1) * 3, 1)
      const quarterEnd = new Date(year, quarter * 3, 0)
      const days = (quarterEnd - quarterStart) / (1000 * 60 * 60 * 24) + 1
      const width = days * state.dayWidth

      quarters.push({
        key: `Q${quarter}-${year}`,
        quarter,
        year,
        position,
        width,
        isCurrent: new Date().getFullYear() === year && getQuarter(new Date()) === quarter
      })

      position += width

      if (quarterEnd >= maxDate) break

      quarter++
      if (quarter > 4) {
        quarter = 1
        year++
      }
    }

    return quarters
  }),

  // 时间轴宽度
  timelineWidth: computed(() => {
    if (state.filteredTasks.length === 0) return 800

    switch (state.viewMode) {
      case 'day':
        return (getters.timelineDays.value.length * state.dayWidth) || 800
      case 'week':
        return getters.timelineWeeks.value.reduce((sum, w) => sum + w.width, 0) || 800
      case 'month':
        return getters.timelineMonths.value.reduce((sum, m) => sum + m.width, 0) || 800
      case 'quarter':
        return getters.timelineQuarters.value.reduce((sum, q) => sum + q.width, 0) || 800
      default:
        return 800
    }
  }),

  // 今天的位置
  todayPosition: computed(() => {
    const todayStr = formatDate(new Date())

    switch (state.viewMode) {
      case 'day':
        const todayDay = getters.timelineDays.value.find(d => d.date === todayStr)
        return todayDay ? todayDay.position + state.dayWidth / 2 : null

      case 'week':
        const todayWeek = getters.timelineWeeks.value.find(w => w.isCurrent)
        return todayWeek ? todayWeek.position + todayWeek.width / 2 : null

      case 'month':
        const today = new Date()
        const todayMonth = getters.timelineMonths.value.find(m =>
          m.year === today.getFullYear() && m.month === today.getMonth() + 1
        )
        return todayMonth ? todayMonth.position + todayMonth.width / 2 : null

      case 'quarter':
        const currentQuarter = getQuarter(new Date())
        const currentQ = getters.timelineQuarters.value.find(q =>
          q.quarter === currentQuarter && q.year === new Date().getFullYear()
        )
        return currentQ ? currentQ.position + currentQ.width / 2 : null

      default:
        return null
    }
  }),

  // 分组后的任务
  groupedTasks: computed(() => {
    if (!state.groupMode) return []

    const tasks = state.filteredTasks

    if (state.groupMode === 'status') {
      const groups = groupTasksByStatus(tasks)
      return Object.entries(groups)
        .filter(([_, group]) => group.tasks.length > 0)
        .map(([key, group]) => ({
          name: key,
          label: group.name,
          tasks: group.tasks
        }))
    }

    if (state.groupMode === 'priority') {
      const groups = groupTasksByPriority(tasks)
      return Object.entries(groups)
        .filter(([_, group]) => group.tasks.length > 0)
        .map(([key, group]) => ({
          name: key,
          label: group.name,
          tasks: group.tasks
        }))
    }

    return []
  }),

  // 可见的依赖关系
  visibleDependencies: computed(() => {
    if (!state.showDependencies) return []

    // 内联获取可见任务的逻辑
    const getVisibleTasks = (tasks) => {
      return tasks.filter(task => {
        const taskId = task.id
        const checkTask = (tid) => {
          const t = tasks.find(task => task.id === tid)
          if (!t) return false
          if (!t.parent_id) return false
          if (state.collapsedTasks.has(t.parent_id)) return true
          return checkTask(t.parent_id)
        }
        return !checkTask(taskId)
      })
    }

    const visibleTasks = getVisibleTasks(state.filteredTasks)

    const dependencies = filterDependencies(
      state.scheduleData,
      visibleTasks,
      {
        showCriticalOnly: false,
        showNonCritical: true
      }
    )

    const paths = calculateAllDependencyPaths(
      dependencies,
      visibleTasks,
      getters.timelineDays.value,
      state.dayWidth,
      state.rowHeight
    ).filter(Boolean)

    return paths
  }),

  // 当前周期文本
  currentPeriodText: computed(() => {
    if (state.filteredTasks.length === 0) return ''

    const tasks = state.filteredTasks
    const minDate = tasks[0].start
    const maxDate = tasks[tasks.length - 1].end

    return `${minDate} ~ ${maxDate}`
  }),

  // 当前缩放标签
  currentZoomLabel: computed(() => {
    const config = state.VIEW_CONFIG[state.viewMode]
    return `${Math.round(state.dayWidth)}${config.label}`
  }),

  // 容器样式
  containerStyle: computed(() => {
    const style = {}
    if (state.containerSize.width) {
      style.width = state.containerSize.width + 'px'
    }
    if (state.containerSize.height) {
      style.height = state.containerSize.height + 'px'
      style.maxHeight = 'none'
    }
    return style
  }),

  // 箭头相关
  arrowColor: computed(() => '#909399'),
  arrowMarkerId: computed(() => `arrow-marker-${state.showCriticalPath ? 'critical' : 'normal'}`),

  // 空状态描述
  emptyDescription: computed(() => {
    if (state.searchKeyword) return '未找到匹配的任务'
    return '暂无进度计划数据'
  }),

  // 拖拽提示文本
  tooltipText: computed(() => {
    if (!state.previewTask) return ''

    const modeText = {
      move: '移动',
      resize_left: '调整开始',
      resize_right: '调整结束'
    }

    const mode = modeText[state.dragMode] || ''
    const duration = diffDays(state.previewTask.start, state.previewTask.end)

    return `${mode}: ${state.previewTask.start} ~ ${state.previewTask.end} (${duration}天)`
  })
}

// ==================== Actions (同步操作) ====================
const actions = {
  /**
   * 设置项目
   */
  setProject(projectId, projectName) {
    state.projectId = projectId
    state.projectName = projectName
  },

  /**
   * 加载数据
   */
  async loadData() {
    state.loading = true
    try {
      const response = await progressApi.getProjectSchedule(state.projectId)
      state.scheduleData = response.data
      await this.formatTasks()

      eventBus.emit(GanttEvents.DATA_LOADED, { source: 'load' })
      eventBus.emit(GanttEvents.DATA_REFRESHED, { source: 'load' })
    } catch (error) {
      console.error('加载甘特图数据失败:', error)
      throw error
    } finally {
      state.loading = false
    }
  },

  /**
   * 格式化任务数据
   */
  async formatTasks() {
    const tasks = []
    const activities = state.scheduleData.activities || {}

    for (const [key, activity] of Object.entries(activities)) {
      if (!activity.is_dummy) {
        const startDate = new Date(activity.earliest_start * 1000)
        const endDate = new Date(activity.latest_finish * 1000)

        // 验证日期是否有效
        if (isNaN(startDate.getTime()) || isNaN(endDate.getTime())) {
          console.warn('无效的日期:', activity)
          continue
        }

        // 提取数字ID
        let taskId = activity.id
        if (typeof taskId === 'string' && taskId.includes('_')) {
          taskId = parseInt(taskId.split('_').pop()) || 0
        }

        // 使用本地映射中的 parent_id（如果存在）
        const parentId = state.localParentIdMap.has(taskId)
          ? state.localParentIdMap.get(taskId)
          : (activity.parent_id || null)

        tasks.push({
          id: taskId,
          name: activity.name || '未命名任务',
          start: formatDate(startDate),
          end: formatDate(endDate),
          startDate: startDate,
          endDate: endDate,
          baseline_start: activity.baseline_start ? formatDate(new Date(activity.baseline_start * 1000)) : null,
          baseline_end: activity.baseline_end ? formatDate(new Date(activity.baseline_end * 1000)) : null,
          duration: activity.duration || 1,
          progress: Math.round(activity.progress || 0),
          status: this.getActivityStatus(activity),
          priority: activity.priority || 'medium',
          is_critical: activity.is_critical || false,
          predecessors: activity.predecessors || [],
          successors: activity.successors || [],
          resources: activity.resources || [],
          parent_id: parentId,
          sort_order: activity.sort_order || 0
        })
      }
    }

    // 按 sort_order 和 parent_id 构建树形结构
    const sortedTasks = tasks.sort((a, b) => {
      if (a.sort_order !== b.sort_order) {
        return (a.sort_order || 0) - (b.sort_order || 0)
      }
      return a.id - b.id
    })

    // 构建树形结构的显示顺序
    const buildTreeOrder = (tasks, parentId = null) => {
      const result = []
      const children = tasks.filter(t => (t.parent_id || null) === parentId)

      for (const child of children) {
        result.push(child)
        result.push(...buildTreeOrder(tasks, child.id))
      }

      return result
    }

    state.tasks = buildTreeOrder(sortedTasks, null)
    this.filterTasks(state.searchKeyword)
  },

  /**
   * 获取活动状态
   */
  getActivityStatus(activity) {
    const now = new Date()
    const endDate = new Date(activity.latest_finish * 1000)

    if (activity.progress >= 100) return 'completed'
    if (endDate < now && activity.progress < 100) return 'delayed'
    if (activity.progress > 0) return 'in_progress'
    return 'not_started'
  },

  /**
   * 过滤任务
   */
  filterTasks(keyword) {
    state.searchKeyword = keyword
    let tasks = [...state.tasks]

    if (keyword) {
      const lowerKeyword = keyword.toLowerCase()
      tasks = tasks.filter(task =>
        task.name.toLowerCase().includes(lowerKeyword)
      )
    }

    state.filteredTasks = tasks
  },

  /**
   * 更新视图模式
   */
  setViewMode(mode) {
    state.viewMode = mode
    state.dayWidth = state.VIEW_CONFIG[mode].default
    eventBus.emit(GanttEvents.VIEW_CHANGED, { mode })
  },

  /**
   * 缩放控制
   */
  zoomIn() {
    const config = state.VIEW_CONFIG[state.viewMode]
    state.dayWidth = Math.min(
      state.dayWidth + (state.viewMode === 'day' ? 10 : 5),
      config.maxWidth
    )
    eventBus.emit(GanttEvents.ZOOM_CHANGED, { width: state.dayWidth })
  },

  zoomOut() {
    const config = state.VIEW_CONFIG[state.viewMode]
    state.dayWidth = Math.max(
      state.dayWidth - (state.viewMode === 'day' ? 10 : 5),
      config.minWidth
    )
    eventBus.emit(GanttEvents.ZOOM_CHANGED, { width: state.dayWidth })
  },

  zoomReset() {
    state.dayWidth = state.VIEW_CONFIG[state.viewMode].default
    eventBus.emit(GanttEvents.ZOOM_CHANGED, { width: state.dayWidth })
  },

  /**
   * 自动适应容器尺寸
   */
  autoFit() {
    if (state.tasks.length === 0) return

    state.containerSize = { width: null, height: null }

    const tasks = state.tasks
    const minDate = new Date(Math.min(...tasks.map(t => t.startDate.getTime())))
    const maxDate = new Date(Math.max(...tasks.map(t => t.endDate.getTime())))
    const totalDays = Math.ceil((maxDate - minDate) / (1000 * 60 * 60 * 24))

    const availableWidth = 1200
    const optimalDayWidth = Math.floor(availableWidth / totalDays)

    if (optimalDayWidth >= 30) {
      state.viewMode = 'day'
      state.dayWidth = Math.max(
        state.VIEW_CONFIG.day.minWidth,
        Math.min(state.VIEW_CONFIG.day.maxWidth, optimalDayWidth)
      )
    } else if (optimalDayWidth >= 10) {
      state.viewMode = 'week'
      state.dayWidth = Math.max(
        state.VIEW_CONFIG.week.minWidth,
        Math.min(state.VIEW_CONFIG.week.maxWidth, optimalDayWidth * 7)
      )
    } else if (optimalDayWidth >= 3) {
      state.viewMode = 'month'
      state.dayWidth = Math.max(
        state.VIEW_CONFIG.month.minWidth,
        Math.min(state.VIEW_CONFIG.month.maxWidth, optimalDayWidth * 30)
      )
    } else {
      state.viewMode = 'quarter'
      state.dayWidth = Math.max(
        state.VIEW_CONFIG.quarter.minWidth,
        Math.min(state.VIEW_CONFIG.quarter.maxWidth, optimalDayWidth * 90)
      )
    }

    eventBus.emit(GanttEvents.VIEW_AUTO_FITTED, {
      mode: state.viewMode,
      dayWidth: state.dayWidth
    })
  },

  /**
   * 选择任务
   */
  selectTask(taskId) {
    state.selectedTaskId = taskId
    const task = state.tasks.find(t => t.id === taskId)
    eventBus.emit(GanttEvents.TASK_SELECTED, { task })
  },

  /**
   * 开始拖拽（时间轴上的任务拖拽）
   */
  startDrag(task, mode = 'move') {
    state.isDragging = true
    state.dragMode = mode
    state.draggedTask = task
    state.originalTask = { ...task }

    // 初始化预览任务
    state.previewTask = {
      ...task,
      start: task.start,
      end: task.end
    }

    eventBus.emit(GanttEvents.TASK_DRAG_START, { task, mode })
  },

  /**
   * 更新拖拽预览
   */
  updateDragPreview(dayOffset) {
    if (!state.originalTask) return

    const original = state.originalTask

    switch (state.dragMode) {
      case 'move':
        state.previewTask = {
          ...original,
          start: formatDate(addDays(original.start, dayOffset)),
          end: formatDate(addDays(original.end, dayOffset))
        }
        break

      case 'resize_left':
        const newStart = addDays(original.start, dayOffset)
        if (newStart <= new Date(original.end)) {
          state.previewTask = {
            ...original,
            start: formatDate(newStart)
          }
        }
        break

      case 'resize_right':
        const newEnd = addDays(original.end, dayOffset)
        if (newEnd >= new Date(original.start)) {
          state.previewTask = {
            ...original,
            end: formatDate(newEnd)
          }
        }
        break
    }

    if (state.previewTask) {
      const duration = diffDays(state.previewTask.start, state.previewTask.end)
      state.previewTask.duration = Math.max(duration, 0)
    }

    eventBus.emit(GanttEvents.TASK_DRAG_MOVE, { preview: state.previewTask })
  },

  /**
   * 结束拖拽
   */
  endDrag(newTask, originalTask) {
    const updateData = {
      task_name: newTask.name,
      start_date: newTask.start,
      end_date: newTask.end,
      progress: newTask.progress,
      priority: newTask.priority
    }

    state.pendingTaskUpdates.set(originalTask.id, {
      id: originalTask.id,
      ...updateData
    })

    // 更新本地状态
    const task = state.tasks.find(t => t.id === originalTask.id)
    if (task) {
      Object.assign(task, newTask)
      task.startDate = new Date(newTask.start)
      task.endDate = new Date(newTask.end)
    }

    this.markUnsaved()

    eventBus.emit(GanttEvents.TASK_DRAG_END, {
      fromTask: originalTask,
      toTask: newTask
    })

    // 重置拖拽状态
    state.isDragging = false
    state.dragMode = 'none'
    state.draggedTask = null
    state.originalTask = null
    state.previewTask = null
  },

  /**
   * 取消拖拽
   */
  cancelDrag() {
    state.isDragging = false
    state.dragMode = 'none'
    state.draggedTask = null
    state.originalTask = null
    state.previewTask = null
    state.tooltipVisible = false
  },

  /**
   * 任务排序（拖拽改变位置/层级）
   */
  reorderTask(fromTaskId, toTaskId, position) {
    const fromTask = state.tasks.find(t => t.id === fromTaskId)
    const toTask = state.tasks.find(t => t.id === toTaskId)

    if (!fromTask) return

    let newParentId = null
    let targetSortOrder = null

    if (position === 'child') {
      newParentId = toTaskId
    } else {
      newParentId = toTask?.parent_id || null

      if (toTask) {
        const siblings = state.tasks.filter(t =>
          (t.parent_id || null) === newParentId &&
          t.id !== fromTask.id
        )

        const targetIndex = siblings.findIndex(t => t.id === toTask.id)

        if (position === 'before') {
          if (targetIndex > 0) {
            const prevTask = siblings[targetIndex - 1]
            targetSortOrder = ((prevTask.sort_order || 0) + (toTask.sort_order || 100)) / 2
          } else {
            targetSortOrder = (toTask.sort_order || 100) - 10
          }
        } else {
          if (targetIndex < siblings.length - 1) {
            const nextTask = siblings[targetIndex + 1]
            targetSortOrder = ((toTask.sort_order || 100) + (nextTask.sort_order || 200)) / 2
          } else {
            targetSortOrder = (toTask.sort_order || 100) + 10
          }
        }
      }
    }

    // 更新本地映射
    state.localParentIdMap.set(fromTask.id, newParentId)
    fromTask.parent_id = newParentId

    // 保存到待更新列表
    let updateData = {
      id: fromTask.id,
      project_id: state.projectId,
      parent_id: newParentId
    }

    if (position === 'before' || position === 'after') {
      updateData.sort_order = targetSortOrder !== null ? Math.round(targetSortOrder * 100) / 100 : null
    }

    state.pendingTaskUpdates.set(fromTask.id, updateData)

    this.markUnsaved()

    eventBus.emit(GanttEvents.TASK_REORDERED, {
      fromTask,
      toTask,
      position,
      newParentId,
      sortOrder: targetSortOrder
    })

    // 重新排序显示
    this.formatTasks()
  },

  /**
   * 添加依赖关系
   */
  addDependency(fromId, toId, type = 'FS', lag = 0) {
    state.pendingDependencyCreations.push({
      fromId,  // 前置任务（源）
      toId,    // 后续任务（目标）
      type,
      lag
    })

    this.markUnsaved()
    eventBus.emit(GanttEvents.DEPENDENCY_CREATED, { fromId, toId, type, lag })
  },

  /**
   * 删除依赖关系
   */
  removeDependency(depId) {
    state.pendingDependencyDeletions.push(depId)
    this.markUnsaved()
    eventBus.emit(GanttEvents.DEPENDENCY_DELETED, { depId })
  },

  /**
   * 开始创建依赖关系
   */
  startDependencyCreation(task) {
    state.isCreatingDependency = true
    state.dependencySourceTask = task
    eventBus.emit(GanttEvents.DEPENDENCY_CREATING, { task })
  },

  /**
   * 完成依赖关系创建
   */
  completeDependencyCreation(targetTask) {
    if (!state.dependencySourceTask || !targetTask) return

    this.addDependency(state.dependencySourceTask.id, targetTask.id, 'FS', 0)

    state.isCreatingDependency = false
    state.dependencySourceTask = null
    state.tempLineEnd = null
  },

  /**
   * 取消依赖关系创建
   */
  cancelDependencyCreation() {
    state.isCreatingDependency = false
    state.dependencySourceTask = null
    state.tempLineEnd = null
  },

  /**
   * 分配资源
   */
  allocateResource(taskId, resourceId, quantity, cost) {
    const task = state.tasks.find(t => t.id === taskId)
    if (task) {
      if (!task.resources) task.resources = []

      const existing = task.resources.find(r => r.resource_id === resourceId)
      if (existing) {
        existing.quantity = quantity
        existing.cost = cost
      } else {
        task.resources.push({ resource_id: resourceId, quantity, cost })
      }
    }

    state.pendingTaskUpdates.set(taskId, {
      resources: task?.resources || []
    })

    this.markUnsaved()
    eventBus.emit(GanttEvents.RESOURCE_ALLOCATED, { taskId, resourceId, quantity })
  },

  /**
   * 标记未保存
   */
  markUnsaved() {
    state.hasUnsavedChanges = true
    this.startAutoSave()

    eventBus.emit(GanttEvents.DATA_CHANGED, {
      hasUnsavedChanges: true,
      pendingUpdates: state.pendingTaskUpdates.size
    })
  },

  /**
   * 启动自动保存
   */
  startAutoSave() {
    if (state.autoSaveTimer) {
      clearTimeout(state.autoSaveTimer)
    }
    state.autoSaveTimer = setTimeout(() => {
      this.saveAll()
    }, state.AUTO_SAVE_INTERVAL)
  },

  /**
   * 停止自动保存
   */
  stopAutoSave() {
    if (state.autoSaveTimer) {
      clearTimeout(state.autoSaveTimer)
      state.autoSaveTimer = null
    }
  },

  /**
   * 保存所有更改
   */
  async saveAll() {
    if (!state.hasUnsavedChanges) return

    state.saving = true
    this.stopAutoSave()

    try {
      console.log('开始保存所有更改...')

      // 1. 保存任务更新
      const taskPromises = []
      for (const [taskId, data] of state.pendingTaskUpdates) {
        taskPromises.push(progressApi.update(taskId, data))
      }

      if (taskPromises.length > 0) {
        await Promise.all(taskPromises)
        console.log(`已保存 ${taskPromises.length} 个任务更新`)
      }

      // 2. 保存依赖关系
      for (const dep of state.pendingDependencyCreations) {
        await progressApi.createDependencyVisual(dep.fromId, dep.toId, {
          type: dep.type || 'FS',
          lag: dep.lag || 0
        })
      }

      // 3. 删除依赖关系
      for (const depId of state.pendingDependencyDeletions) {
        await progressApi.removeDependency(depId)
      }

      // 清空待保存列表
      state.pendingTaskUpdates.clear()
      state.pendingDependencyCreations = []
      state.pendingDependencyDeletions = []
      state.hasUnsavedChanges = false

      // 重新加载数据
      await this.loadData()

      eventBus.emit(GanttEvents.DATA_SAVED, {
        timestamp: new Date()
      })
    } catch (error) {
      console.error('保存失败:', error)
      eventBus.emit(GanttEvents.DATA_SAVE_ERROR, { error })
      throw error
    } finally {
      state.saving = false
    }
  },

  /**
   * 加载资源库
   */
  async loadResources() {
    try {
      const response = await progressApi.getProjectResources(state.projectId)
      state.resourceLibrary = response.data || []
    } catch (error) {
      console.error('加载资源库失败:', error)
      state.resourceLibrary = []
    }
  },

  /**
   * 切换分组
   */
  setGroupMode(mode) {
    state.groupMode = mode
    state.collapsedGroups.clear()
  },

  /**
   * 切换分组折叠状态
   */
  toggleGroup(groupName) {
    if (state.collapsedGroups.has(groupName)) {
      state.collapsedGroups.delete(groupName)
    } else {
      state.collapsedGroups.add(groupName)
    }
  },

  /**
   * 切换任务折叠状态（树形结构）
   */
  toggleTaskCollapse(taskId) {
    if (state.collapsedTasks.has(taskId)) {
      state.collapsedTasks.delete(taskId)
    } else {
      state.collapsedTasks.add(taskId)
    }
  },

  /**
   * 检查任务是否被折叠
   */
  isTaskCollapsed(taskId) {
    return state.collapsedTasks.has(taskId)
  },

  /**
   * 检查任务是否应该被隐藏（因为父任务被折叠）
   */
  isTaskHidden(taskId, tasks = state.filteredTasks) {
    const task = tasks.find(t => t.id === taskId)
    if (!task) return false
    if (!task.parent_id) return false

    const parent = tasks.find(t => t.id === task.parent_id)
    if (!parent) return false

    if (state.collapsedTasks.has(parent.id)) return true
    return this.isTaskHidden(parent.id, tasks)
  },

  /**
   * 获取应该显示的任务列表（考虑折叠状态）
   */
  getVisibleTasks(tasks = state.filteredTasks) {
    return tasks.filter(task => !this.isTaskHidden(task.id, tasks))
  },

  /**
   * 切换显示选项
   */
  toggleDependencies() {
    state.showDependencies = !state.showDependencies
  },

  toggleCriticalPath() {
    state.showCriticalPath = !state.showCriticalPath
  },

  toggleBaseline() {
    state.showBaseline = !state.showBaseline
  },

  /**
   * 对话框控制
   */
  openEditDialog(task = null) {
    state.editingTask = task
    state.editDialogVisible = true
    state.taskDetailVisible = false
    eventBus.emit(GanttEvents.DIALOG_OPEN, { type: 'edit', task })
  },

  closeEditDialog() {
    state.editDialogVisible = false
    state.editingTask = null
    eventBus.emit(GanttEvents.DIALOG_CLOSE, { type: 'edit' })
  },

  openResourceDialog(task) {
    state.currentTaskForResource = task
    state.resourceDialogVisible = true
    eventBus.emit(GanttEvents.DIALOG_OPEN, { type: 'resource', task })
  },

  closeResourceDialog() {
    state.resourceDialogVisible = false
    state.currentTaskForResource = null
    eventBus.emit(GanttEvents.DIALOG_CLOSE, { type: 'resource' })
  },

  openResourceManagementDialog() {
    state.resourceManagementDialogVisible = true
    eventBus.emit(GanttEvents.DIALOG_OPEN, { type: 'resource-management' })
  },

  closeResourceManagementDialog() {
    state.resourceManagementDialogVisible = false
    eventBus.emit(GanttEvents.DIALOG_CLOSE, { type: 'resource-management' })
  },

  openTaskDetail(task) {
    this.selectTask(task.id)
    state.taskDetailVisible = true
  },

  closeTaskDetail() {
    state.taskDetailVisible = false
  },

  /**
   * 右键菜单控制
   */
  showContextMenu(task, position) {
    state.contextMenuTask = task
    state.contextMenuPosition = position
    state.contextMenuVisible = true
    eventBus.emit(GanttEvents.CONTEXT_MENU_SHOW, { task, position })
  },

  hideContextMenu() {
    state.contextMenuVisible = false
    state.contextMenuTask = null
    eventBus.emit(GanttEvents.CONTEXT_MENU_HIDE)
  },

  /**
   * 全屏控制
   */
  setFullscreen(isFullscreen) {
    state.isFullscreen = isFullscreen
  },

  /**
   * 容器大小控制
   */
  setContainerSize(width, height) {
    state.containerSize = { width, height }
  },

  /**
   * 设置调整大小状态
   */
  setResizing(isResizing) {
    state.isResizing = isResizing
  },

  /**
   * 更新临时连线位置（依赖创建时）
   */
  updateTempLineEnd(position) {
    state.tempLineEnd = position
  },

  /**
   * 设置拖拽提示位置
   */
  setTooltipPosition(position) {
    state.tooltipPosition = position
  },

  /**
   * 显示/隐藏拖拽提示
   */
  setTooltipVisible(visible) {
    state.tooltipVisible = visible
  }
}

// ==================== 导出 ====================
export const ganttStore = {
  state,
  getters,
  actions
}

export default ganttStore
