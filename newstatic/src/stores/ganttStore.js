/**
 * 甘特图集中式状态管理 Store
 * 统一管理甘特图的所有状态和操作
 */

import { reactive, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
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
  rowHeight: 50,
  taskHeight: 24,
  VIEW_CONFIG: {
    day: { minWidth: 20, maxWidth: 100, default: 40, label: '日' },
    'month-day': { minWidth: 20, maxWidth: 100, default: 40, label: '月/日' },
    'year-month-day': { minWidth: 20, maxWidth: 100, default: 40, label: '年/月/日' },
    'year-month': { minWidth: 30, maxWidth: 150, default: 60, label: '年/月' },
    week: { minWidth: 20, maxWidth: 80, default: 40, label: '周' },
    month: { minWidth: 15, maxWidth: 60, default: 30, label: '月' },
    quarter: { minWidth: 10, maxWidth: 50, default: 25, label: '季' }
  },

  // 时间轴格式预设配置
  TIMELINE_FORMATS: {
    'day': { label: '日期', layers: 1, format: ['day'] },
    'month-day': { label: '月/日', layers: 2, format: ['month', 'day'] },
    'year-month': { label: '年/月', layers: 2, format: ['year', 'month'] },
    'year-month-day': { label: '年/月/日', layers: 3, format: ['year', 'month', 'day'] },
    'week': { label: '周', layers: 1, format: ['week'] },
    'month': { label: '月', layers: 1, format: ['month'] },
    'quarter': { label: '季度', layers: 1, format: ['quarter'] }
  },

  // 日期显示格式预设
  DATE_FORMATS: {
    'all': { label: '全部', interval: 1 },
    'odd': { label: '奇数 (1,3,5)', interval: 2 },
    'interval3': { label: '间隔3天 (1,4,7)', interval: 3 },
    'interval5': { label: '间隔5天 (1,6,11)', interval: 5 },
    'first': { label: '每月1号', interval: 'first' }
  },

  // 显示选项
  showDependencies: true,
  showCriticalPath: true,
  showBaseline: false,
  showTaskList: true, // 显示任务列表

  // 时间轴格式配置
  timelineFormat: 'month-day', // 时间轴头部显示格式
  // 时间轴格式预设：
  // - 'day': 单层 - 只显示日期
  // - 'month-day': 双层 - 上月、下日（默认）
  // - 'year-month': 双层 - 上年、下月
  // - 'year-month-day': 三层 - 年月日
  // - 'week': 单层 - 只显示周
  // - 'month': 单层 - 只显示月
  // - 'quarter': 单层 - 只显示季度
  // 日期显示格式：
  // - 'all': 显示所有日期 1 2 3 4 5...
  // - 'odd': 只显示奇数 1 3 5 7...
  // - 'interval2': 间隔2天 1 3 5 7... (等同于odd)
  // - 'interval3': 间隔3天 1 4 7 10...
  // - 'interval5': 间隔5天 1 6 11 16...
  // - 'first': 每月1号
  dateDisplayFormat: 'all', // 日期显示格式

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

  // 平移模式（手形工具）
  panMode: false,
  scrollLeft: 0, // 时间轴横向滚动位置

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
    console.log('timelineDays - filteredTasks:', tasks.length)

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

    // 计算可用宽度（假设最小可视宽度为800px）
    const availableWidth = Math.max(800, state.containerSize.width || window.innerWidth - 400)
    const minDayWidth = 30

    // 计算最大可显示的天数
    const maxVisibleDays = Math.floor(availableWidth / minDayWidth)

    // 如果总天数超过可显示天数，使用稀疏显示
    const useSparseDisplay = totalDays > maxVisibleDays * 0.8

    const days = []
    const currentDate = new Date(startDate)
    const today = new Date()

    if (useSparseDisplay) {
      // 稀疏显示模式：使用1-3-5-7-9或1-3-5-9模式
      let displayPattern = '13579'
      if (totalDays > maxVisibleDays * 1.5) {
        displayPattern = '159'
      }

      let displayIndex = 0
      const patternDigits = displayPattern.split('').map(Number)

      while (currentDate <= endDate) {
        const dateStr = formatDate(currentDate)
        const shouldDisplay = patternDigits[displayIndex] || 1

        if (shouldDisplay) {
          days.push({
            date: dateStr,
            day: currentDate.getDate(),
            weekday: ['日', '一', '二', '三', '四', '五', '六'][currentDate.getDay()],
            isToday: dateStr === formatDate(today),
            isWeekend: currentDate.getDay() === 0 || currentDate.getDay() === 6,
            position: 0
          })
        }

        displayIndex = (displayIndex + 1) % patternDigits.length
        currentDate.setDate(currentDate.getDate() + 1)
      }
    } else {
      // 正常显示模式：显示所有日期
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

  // 时间轴月份行（用于日视图双层显示的上层）
  timelineHeaderMonths: computed(() => {
    const days = getters.timelineDays.value
    if (days.length === 0) return []

    const months = []
    let currentMonth = null

    days.forEach((day, index) => {
      const date = new Date(day.date)
      const year = date.getFullYear()
      const month = date.getMonth() + 1
      const monthKey = `${year}-${month}`

      if (monthKey !== currentMonth) {
        // 开始新月份
        currentMonth = monthKey
        months.push({
          key: monthKey,
          label: `${year}年${month}月`,
          position: day.position,
          width: 0  // 稍后计算
        })
      }

      // 累加当前月份的宽度
      if (months.length > 0) {
        const lastMonth = months[months.length - 1]
        lastMonth.width += state.dayWidth
      }
    })

    return months
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
        const checkTask = (tid, visited = new Set()) => {
          if (visited.has(tid)) {
            return false // 循环引用，停止检查
          }
          visited.add(tid)
          const t = tasks.find(task => task.id === tid)
          if (!t) return false
          if (!t.parent_id) return false
          if (state.collapsedTasks.has(t.parent_id)) return true
          return checkTask(t.parent_id, visited)
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

    // 不过滤依赖关系，渲染所有依赖关系
    // 路径计算逻辑会自动处理异常情况（避免穿越任务条）
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
  }),

  // 暴露时间轴月份行（用于双层显示）
  timelineHeaderMonths: computed(() => state.timelineHeaderMonths)
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

      // 加载任务表以获取最新的 parent_id 信息
      try {
        const tasksResponse = await progressApi.getTasks(state.projectId)
        const tasks = tasksResponse.data || []

        // 更新本地 parent_id 映射
        state.localParentIdMap.clear()
        for (const task of tasks) {
          if (task.parent_id !== null && task.parent_id !== undefined) {
            state.localParentIdMap.set(task.id, task.parent_id)
          }
        }

        console.log('已从数据库加载 parent_id 映射:', Object.fromEntries(state.localParentIdMap))
      } catch (error) {
        console.warn('加载任务 parent_id 失败，将使用 schedule 数据:', error)
      }

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

    console.log('formatTasks - activities count:', Object.keys(activities).length)
    console.log('formatTasks - activities sample:', Object.entries(activities).slice(0, 3))

    for (const [key, activity] of Object.entries(activities)) {
      console.log('formatTasks - processing activity:', key, 'is_dummy:', activity.is_dummy, 'name:', activity.name)
      if (!activity.is_dummy) {
        const startDate = new Date(activity.earliest_start * 1000)
        const endDate = new Date(activity.earliest_finish * 1000)

        // 验证日期是否有效
        if (isNaN(startDate.getTime()) || isNaN(endDate.getTime())) {
          console.warn('无效的日期:', activity)
          continue
        }

        // 使用后端提供的 task_id（数字类型），如果没有则从 id 字符串中提取
        let taskId = activity.task_id
        if (!taskId) {
          taskId = activity.id
          if (typeof taskId === 'string' && taskId.includes('_')) {
            taskId = parseInt(taskId.split('_').pop()) || 0
          }
        }

        // 使用本地映射中的 parent_id（如果存在）
        const parentId = state.localParentIdMap.has(taskId)
          ? state.localParentIdMap.get(taskId)
          : (activity.parent_id || null)

        const taskObj = {
          id: taskId,
          name: activity.name || '未命名任务',
          start: formatDate(startDate),
          end: formatDate(endDate),
          startDate: startDate,
          endDate: endDate,
          baseline_start: activity.baseline_start ? formatDate(new Date(activity.baseline_start * 1000)) : null,
          baseline_end: activity.baseline_end ? formatDate(new Date(activity.baseline_end * 1000)) : null,
          duration: activity.duration || 1,
          // 处理进度值：如果是小数(0-1范围)，转换为百分比(0-100)；如果是百分比，直接使用
          progress: (activity.progress || 0) < 1
            ? Math.round((activity.progress || 0) * 100)
            : Math.round(activity.progress || 0),
          status: this.getActivityStatus(activity),
          priority: activity.priority || 'medium',
          is_critical: activity.is_critical || false,
          predecessors: activity.predecessors || [],
          successors: activity.successors || [],
          resources: activity.resources || [],
          parent_id: parentId,
          sort_order: activity.sort_order || 0
        }
        console.log('创建任务对象:', taskObj)
        tasks.push(taskObj)
      }
    }

    // 按 sort_order 和 parent_id 构建树形结构
    const sortedTasks = tasks.sort((a, b) => {
      if (a.sort_order !== b.sort_order) {
        return (a.sort_order || 0) - (b.sort_order || 0)
      }
      return a.id - b.id
    })

    console.log('formatTasks - sortedTasks 的 parent_id 情况:', sortedTasks.map(t => ({ id: t.id, name: t.name, parent_id: t.parent_id })))

    // 构建树形结构的显示顺序
    const buildTreeOrder = (tasks, parentId = null, visited = new Set()) => {
      // 检测循环
      const key = `${parentId}`
      if (visited.has(key)) {
        console.warn(`检测到循环: parent_id=${parentId}`)
        return []
      }
      visited.add(key)

      const result = []
      const children = tasks.filter(t => (t.parent_id || null) === parentId)
      console.log(`buildTreeOrder - parentId: ${parentId}, 找到 ${children.length} 个子任务`)

      for (const child of children) {
        result.push(child)
        result.push(...buildTreeOrder(tasks, child.id, visited))
      }

      return result
    }

    let finalTasks = buildTreeOrder(sortedTasks, null)

    // 如果没有找到根任务（所有任务都有parent_id），则显示所有任务
    if (finalTasks.length === 0 && sortedTasks.length > 0) {
      console.warn('没有找到根任务，显示所有任务（按sort_order排序）')
      finalTasks = sortedTasks
    }

    // 确保 finalTasks 包含所有任务（修复树形结构可能遗漏任务的问题）
    if (finalTasks.length < sortedTasks.length) {
      const includedIds = new Set(finalTasks.map(t => t.id))
      const missingTasks = sortedTasks.filter(t => !includedIds.has(t.id))
      if (missingTasks.length > 0) {
        console.warn(`树形结构遗漏了 ${missingTasks.length} 个任务，将它们添加到末尾:`, missingTasks.map(t => ({ id: t.id, name: t.name, parent_id: t.parent_id })))
        finalTasks = [...finalTasks, ...missingTasks]
      }
    }

    state.tasks = finalTasks

    // 自动处理汇总任务（父任务）
    // 如果一个任务有子任务，自动将其标记为汇总任务，并更新日期
    this.updateSummaryTasks()

    console.log('formatTasks - 最终 state.tasks 数量:', state.tasks.length)
    console.log('formatTasks - 最终 state.tasks:', state.tasks)
    this.filterTasks(state.searchKeyword)
  },

  /**
   * 更新汇总任务（父任务）
   * 有子任务的父任务自动变成汇总任务，日期由子任务决定
   */
  updateSummaryTasks() {
    // 找出所有有子任务的任务
    const parentTaskIds = new Set()
    for (const task of state.tasks) {
      if (task.parent_id) {
        parentTaskIds.add(task.parent_id)
      }
    }

    // 对每个父任务，更新其日期为子任务的日期范围
    for (const parentId of parentTaskIds) {
      const parentTask = state.tasks.find(t => t.id === parentId)
      if (!parentTask) continue

      // 找出所有子任务（递归查找所有后代）
      const getAllDescendants = (taskId, visited = new Set()) => {
        // 防止循环引用
        if (visited.has(taskId)) {
          console.warn(`检测到循环引用: 任务 ${taskId}`)
          return []
        }
        visited.add(taskId)

        const descendants = []
        const children = state.tasks.filter(t => t.parent_id === taskId)
        for (const child of children) {
          descendants.push(child)
          descendants.push(...getAllDescendants(child.id, visited))
        }
        return descendants
      }

      const descendants = getAllDescendants(parentId)
      if (descendants.length === 0) continue

      // 计算日期范围：最早开始到最晚结束
      let minStartDate = null
      let maxEndDate = null

      for (const child of descendants) {
        const childStart = new Date(child.start)
        const childEnd = new Date(child.end)

        if (!minStartDate || childStart < minStartDate) {
          minStartDate = childStart
        }
        if (!maxEndDate || childEnd > maxEndDate) {
          maxEndDate = childEnd
        }
      }

      // 更新父任务的日期和标记
      if (minStartDate && maxEndDate) {
        parentTask.start = formatDate(minStartDate)
        parentTask.end = formatDate(maxEndDate)
        parentTask.is_summary = true // 标记为汇总任务
        parentTask.startDate = minStartDate
        parentTask.endDate = maxEndDate
      }
    }
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
   * 设置时间轴格式
   */
  setTimelineFormat(format) {
    if (state.TIMELINE_FORMATS[format]) {
      state.timelineFormat = format
      eventBus.emit(GanttEvents.TIMELINE_FORMAT_CHANGED, { format })
    }
  },

  /**
   * 设置日期显示格式
   */
  setDateDisplayFormat(format) {
    if (state.DATE_FORMATS[format]) {
      state.dateDisplayFormat = format
      eventBus.emit(GanttEvents.DATE_FORMAT_CHANGED, { format })
    }
  },

  /**
   * 缩放控制
   * zoomIn: 放大（让日期单元格变小，显示更多日期）
   * zoomOut: 缩小（让日期单元格变大，显示更少日期但更详细）
   */
  zoomIn() {
    const format = state.timelineFormat
    const config = state.VIEW_CONFIG[format] || state.VIEW_CONFIG['month-day']
    // 放大 = 减少 dayWidth，让更多日期显示在屏幕上
    state.dayWidth = Math.max(
      state.dayWidth - 5,
      config.minWidth
    )
    eventBus.emit(GanttEvents.ZOOM_CHANGED, { width: state.dayWidth })
  },

  zoomOut() {
    const format = state.timelineFormat
    const config = state.VIEW_CONFIG[format] || state.VIEW_CONFIG['month-day']
    // 缩小 = 增加 dayWidth，让每个日期单元格更大更详细
    state.dayWidth = Math.min(
      state.dayWidth + 5,
      config.maxWidth
    )
    eventBus.emit(GanttEvents.ZOOM_CHANGED, { width: state.dayWidth })
  },

  zoomReset() {
    const format = state.timelineFormat
    const config = state.VIEW_CONFIG[format] || state.VIEW_CONFIG['month-day']
    state.dayWidth = config.default
    eventBus.emit(GanttEvents.ZOOM_CHANGED, { width: state.dayWidth })
  },

  /**
   * 自动适应容器尺寸
   * 根据实际任务条的时间范围调整时间轴，使所有任务都显示在可见屏幕中
   */
  autoFit() {
    if (state.filteredTasks.length === 0) return

    // 获取时间轴上实际渲染的日期范围
    const timelineDays = getters.timelineDays.value
    if (!timelineDays || timelineDays.length === 0) return

    // 计算时间轴的总天数
    const totalDays = timelineDays.length

    // 获取时间轴容器的可用宽度
    const containerWidth = state.containerSize.width || 1200
    // 减去任务列表宽度（550px）和其他边距（20px）
    const availableWidth = Math.max(600, containerWidth - 570)

    // 计算每天的最佳像素宽度，确保所有天都能显示
    const optimalDayWidth = Math.floor(availableWidth / totalDays)

    console.log('autoFit - totalDays:', totalDays, 'availableWidth:', availableWidth, 'optimalDayWidth:', optimalDayWidth)

    // 根据时间轴格式获取配置
    const format = state.timelineFormat
    const config = state.VIEW_CONFIG[format] || state.VIEW_CONFIG['month-day']

    // 设置新的 dayWidth，确保在合理范围内
    state.dayWidth = Math.max(config.minWidth, Math.min(config.maxWidth, optimalDayWidth))

    console.log('autoFit - timelineFormat:', format, 'dayWidth:', state.dayWidth)

    // 发送事件通知其他组件
    eventBus.emit(GanttEvents.ZOOM_CHANGED, { width: state.dayWidth })
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
    let preview = null

    switch (state.dragMode) {
      case 'move':
        preview = {
          ...original,
          start: formatDate(addDays(original.start, dayOffset)),
          end: formatDate(addDays(original.end, dayOffset))
        }
        break

      case 'resize_left':
        const newStart = addDays(original.start, dayOffset)
        if (newStart <= new Date(original.end)) {
          preview = {
            ...original,
            start: formatDate(newStart)
          }
        }
        break

      case 'resize_right':
        const newEnd = addDays(original.end, dayOffset)
        if (newEnd >= new Date(original.start)) {
          preview = {
            ...original,
            end: formatDate(newEnd)
          }
        }
        break
    }

    if (preview) {
      const duration = diffDays(preview.start, preview.end)
      preview.duration = Math.max(duration, 0)

      // 应用依赖关系约束到预览任务
      state.previewTask = this.applyDependencyConstraints(preview, state.draggedTask?.id)
    } else {
      state.previewTask = null
    }

    eventBus.emit(GanttEvents.TASK_DRAG_MOVE, { preview: state.previewTask })
  },

  /**
   * 更新父任务的日期（基于子任务计算）
   * 当子任务的日期发生变化时，需要更新所有父任务的日期范围
   */
  updateParentTaskDates(taskId) {
    const task = state.tasks.find(t => t.id === taskId)
    if (!task || !task.parent_id) return

    const updateParentRecursive = (parentId) => {
      const parent = state.tasks.find(t => t.id === parentId)
      if (!parent) return

      // 获取所有子任务
      const children = state.tasks.filter(t => t.parent_id === parentId)
      if (children.length === 0) return

      // 计算子任务的最早开始和最晚结束日期
      let earliestStart = null
      let latestEnd = null

      children.forEach(child => {
        const childStart = new Date(child.start)
        const childEnd = new Date(child.end)

        if (!earliestStart || childStart < earliestStart) {
          earliestStart = childStart
        }
        if (!latestEnd || childEnd > latestEnd) {
          latestEnd = childEnd
        }
      })

      // 更新父任务的日期
      const startDate = formatDate(earliestStart)
      const endDate = formatDate(latestEnd)

      parent.start = startDate
      parent.end = endDate
      parent.startDate = earliestStart
      parent.endDate = latestEnd

      // 计算工期（天数）
      const duration = Math.ceil((latestEnd - earliestStart) / (1000 * 60 * 60 * 24))
      parent.duration = duration

      // 递归更新上层的父任务
      if (parent.parent_id) {
        updateParentRecursive(parent.parent_id)
      }
    }

    // 开始更新父任务
    updateParentRecursive(task.parent_id)
  },

  /**
   * 结束拖拽
   */
  endDrag(newTask, originalTask) {
    // 根据拖动模式决定是否应用依赖约束
    let adjustedTask
    if (state.dragMode === 'resize_right') {
      // 调整结束日期：不应用前置约束，只更新结束日期，保持开始日期不变
      adjustedTask = {
        ...newTask,
        start: originalTask.start // 保持原始开始日期
      }
    } else {
      // 移动或调整开始日期：应用依赖约束
      adjustedTask = this.applyDependencyConstraints(newTask, originalTask.id)
    }

    const updateData = {
      name: adjustedTask.name || newTask.name,
      start_date: adjustedTask.start,
      end_date: adjustedTask.end,
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
      Object.assign(task, adjustedTask)
      task.startDate = new Date(adjustedTask.start)
      task.endDate = new Date(adjustedTask.end)
    }

    // 只在非 resize_right 模式下调整后置任务的开始时间
    // resize_right 模式只改变结束时间，不应该影响后置任务
    if (state.dragMode !== 'resize_right') {
      this.adjustSuccessorTasks(originalTask.id)
    }

    // 更新父任务的日期（基于子任务计算）
    this.updateParentTaskDates(originalTask.id)

    this.markUnsaved()

    eventBus.emit(GanttEvents.TASK_DRAG_END, {
      fromTask: originalTask,
      toTask: adjustedTask
    })

    // 重置拖拽状态
    state.isDragging = false
    state.dragMode = 'none'
    state.draggedTask = null
    state.originalTask = null
    state.previewTask = null
  },

  /**
   * 应用依赖关系约束：确保任务开始时间不早于所有前置任务的结束时间
   */
  applyDependencyConstraints(task, taskId) {
    const activities = state.scheduleData.activities
    const activityKey = `task_${taskId}`
    const activity = activities[activityKey]

    if (!activity || !activity.predecessors || activity.predecessors.length === 0) {
      return task
    }

    let latestPredecessorEnd = null

    // 找到所有前置任务的最新结束时间
    for (const predId of activity.predecessors) {
      const predTask = state.tasks.find(t => t.id === predId)
      if (predTask) {
        const predEnd = new Date(predTask.end)
        if (!latestPredecessorEnd || predEnd > latestPredecessorEnd) {
          latestPredecessorEnd = predEnd
        }
      }
    }

    // 如果有前置任务且任务开始时间早于前置任务结束时间，调整开始时间
    if (latestPredecessorEnd) {
      const taskStart = new Date(task.start)
      const taskEnd = new Date(task.end)
      const duration = Math.ceil((taskEnd - taskStart) / (1000 * 60 * 60 * 24))

      if (taskStart < latestPredecessorEnd) {
        // 计算新的开始和结束时间
        const newStart = new Date(latestPredecessorEnd)
        const newEnd = addDays(newStart, duration)

        return {
          ...task,
          start: formatDate(newStart),
          end: formatDate(newEnd)
        }
      }
    }

    return task
  },

  /**
   * 调整所有后置任务的开始时间
   */
  adjustSuccessorTasks(taskId) {
    const activities = state.scheduleData.activities
    const task = state.tasks.find(t => t.id === taskId)

    if (!task) return

    // 递归调整后置任务
    const adjustTask = (currentTaskId) => {
      const currentTask = state.tasks.find(t => t.id === currentTaskId)
      if (!currentTask) return

      const currentActivityKey = `task_${currentTaskId}`
      const currentActivity = activities[currentActivityKey]

      if (!currentActivity || !currentActivity.predecessors) return

      // 检查所有前置任务的结束时间
      let latestPredecessorEnd = null
      for (const predId of currentActivity.predecessors) {
        const predTask = state.tasks.find(t => t.id === predId)
        if (predTask) {
          const predEnd = new Date(predTask.end)
          if (!latestPredecessorEnd || predEnd > latestPredecessorEnd) {
            latestPredecessorEnd = predEnd
          }
        }
      }

      // 如果当前任务开始时间早于最新前置任务结束时间，调整
      if (latestPredecessorEnd) {
        const taskStart = new Date(currentTask.start)
        const taskEnd = new Date(currentTask.end)
        const duration = Math.ceil((taskEnd - taskStart) / (1000 * 60 * 60 * 24))

        if (taskStart < latestPredecessorEnd) {
          const newStart = new Date(latestPredecessorEnd)
          const newEnd = addDays(newStart, duration)

          // 更新任务数据
          currentTask.start = formatDate(newStart)
          currentTask.end = formatDate(newEnd)
          currentTask.startDate = newStart
          currentTask.endDate = newEnd

          // 记录待保存的更新（包含所有必要字段）
          state.pendingTaskUpdates.set(currentTask.id, {
            id: currentTask.id,
            name: currentTask.name,
            start_date: currentTask.start,
            end_date: currentTask.end,
            progress: currentTask.progress || 0,
            priority: currentTask.priority || 1
          })
        }
      }
    }

    // 找到所有后置任务
    const successorIds = []
    for (const [key, activity] of Object.entries(activities)) {
      if (activity.predecessors && activity.predecessors.includes(taskId)) {
        const id = parseInt(key.replace('task_', ''))
        successorIds.push(id)
      }
    }

    // 按依赖关系顺序调整后置任务
    for (const successorId of successorIds) {
      adjustTask(successorId)
      this.adjustSuccessorTasks(successorId)
    }
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
    // 初始化 tempLineEnd 为源任务位置（需要从组件传递）
    // 先设为 null，等待组件第一次鼠标移动时更新
    state.tempLineEnd = null
    eventBus.emit(GanttEvents.DEPENDENCY_CREATING, { task })
  },

  /**
   * 完成依赖关系创建
   */
  async completeDependencyCreation(targetTask) {
    if (!state.dependencySourceTask || !targetTask) return

    // 检查是否是父子关系，如果是则不允许创建
    if (this.checkParentChildRelation(state.dependencySourceTask.id, targetTask.id)) {
      console.warn('不能在父子任务之间创建依赖关系')
      // 显示提示消息
      eventBus.emit(GanttEvents.DEPENDENCY_ERROR, {
        message: '不能在父子任务之间创建依赖关系'
      })
      // 取消依赖关系创建
      this.cancelDependencyCreation()
      return
    }

    // 检查是否是从右到左的连线（源任务的结束日期晚于目标任务的开始日期）
    const sourceEndDate = new Date(state.dependencySourceTask.end)
    const targetStartDate = new Date(targetTask.start)

    if (sourceEndDate > targetStartDate) {
      console.warn('不能创建从右到左的依赖关系')
      // 显示提示消息
      eventBus.emit(GanttEvents.DEPENDENCY_ERROR, {
        message: '不能创建从右到左的依赖关系（前置任务必须先于后置任务）'
      })
      // 取消依赖关系创建
      this.cancelDependencyCreation()
      return
    }

    // 不再强制调整日期 - 只允许从左到右的依赖关系

    // 立即调用 API 创建依赖关系
    try {
      console.log('[创建依赖关系] 源任务:', state.dependencySourceTask.id, state.dependencySourceTask.name)
      console.log('[创建依赖关系] 目标任务:', targetTask.id, targetTask.name)

      const response = await progressApi.createDependencyVisual(
        state.dependencySourceTask.id,
        targetTask.id,
        { type: 'FS', lag: 0 }
      )

      console.log('[创建依赖关系] API响应:', response.data)
      console.log('[创建依赖关系] 正在重新加载数据...')

      // 创建成功后重新加载数据
      await this.loadData()
      ElMessage.success('依赖关系创建成功')

      console.log('[创建依赖关系] 数据重新加载完成')
    } catch (error) {
      console.error('[创建依赖关系] 失败:', error)
      console.error('[创建依赖关系] 错误响应:', error.response?.data)
      const errorMsg = error.response?.data?.error || error.response?.data?.message || error.message || '创建依赖关系失败'
      ElMessage.error(errorMsg)
    }

    state.isCreatingDependency = false
    state.dependencySourceTask = null
    state.tempLineEnd = null
  },

  /**
   * 检查两个任务之间是否存在父子关系
   */
  checkParentChildRelation(taskId1, taskId2) {
    const task1 = state.tasks.find(t => t.id === taskId1)
    const task2 = state.tasks.find(t => t.id === taskId2)

    if (!task1 || !task2) return false

    // 检查task1是否是task2的父任务
    if (task2.parent_id === taskId1) return true

    // 检查task2是否是task1的父任务
    if (task1.parent_id === taskId2) return true

    // 检查task1的祖先任务是否包含task2
    let current = task1
    while (current.parent_id) {
      if (current.parent_id === taskId2) return true
      current = state.tasks.find(t => t.id === current.parent_id)
      if (!current) break
    }

    // 检查task2的祖先任务是否包含task1
    current = task2
    while (current.parent_id) {
      if (current.parent_id === taskId1) return true
      current = state.tasks.find(t => t.id === current.parent_id)
      if (!current) break
    }

    return false
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
   * 获取任务持续天数
   */
  getTaskDuration(task) {
    const startDate = new Date(task.start_date)
    const endDate = new Date(task.end_date)
    return Math.ceil((endDate - startDate) / (1000 * 60 * 60 * 24))
  },

  /**
   * 格式化日期为 YYYY-MM-DD 格式
   */
  formatDate(date) {
    return formatDate(date)
  },

  /**
   * 上移任务（交换任务顺序）
   */
  async moveTaskUp(taskId) {
    const tasks = state.tasks
    const currentIndex = tasks.findIndex(t => t.id === taskId)
    if (currentIndex <= 0) {
      throw new Error('任务已经在最前面')
    }

    const targetIndex = currentIndex - 1
    const currentTask = tasks[currentIndex]
    const targetTask = tasks[targetIndex]

    // 交换 parent_id、wbs_code 等层级相关字段
    const currentParentId = currentTask.parent_id
    const targetParentId = targetTask.parent_id

    // 通过 API 更新
    await progressApi.update(currentTask.id, {
      parent_id: targetParentId
    })
    await progressApi.update(targetTask.id, {
      parent_id: currentParentId
    })

    // 重新加载和格式化任务
    await this.loadData()
  },

  /**
   * 下移任务（交换任务顺序）
   */
  async moveTaskDown(taskId) {
    const tasks = state.tasks
    const currentIndex = tasks.findIndex(t => t.id === taskId)
    if (currentIndex >= tasks.length - 1) {
      throw new Error('任务已经在最后面')
    }

    const targetIndex = currentIndex + 1
    const currentTask = tasks[currentIndex]
    const targetTask = tasks[targetIndex]

    // 交换 parent_id、wbs_code 等层级相关字段
    const currentParentId = currentTask.parent_id
    const targetParentId = targetTask.parent_id

    // 通过 API 更新
    await progressApi.update(currentTask.id, {
      parent_id: targetParentId
    })
    await progressApi.update(targetTask.id, {
      parent_id: currentParentId
    })

    // 重新加载和格式化任务
    await this.loadData()
  },

  /**
   * 将子任务转为独立任务（解除父子关系）
   */
  async convertToIndependentTask(taskId) {
    const task = state.tasks.find(t => t.id === taskId)
    if (!task) {
      throw new Error('任务不存在')
    }

    if (!task.parent_id) {
      throw new Error('任务已经是独立任务')
    }

    // 将 parent_id 设置为 null
    await progressApi.update(taskId, {
      parent_id: null
    })

    // 重新加载和格式化任务
    await this.loadData()
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
      console.log('待保存的任务更新:', Array.from(state.pendingTaskUpdates.entries()))

      // 1. 保存任务更新
      const taskPromises = []
      for (const [taskId, data] of state.pendingTaskUpdates) {
        console.log(`保存任务 ${taskId}:`, data)
        taskPromises.push(progressApi.updateTask(taskId, data))
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
  isTaskHidden(taskId, tasks = state.filteredTasks, visited = new Set()) {
    if (visited.has(taskId)) {
      return false // 循环引用，停止检查
    }
    visited.add(taskId)

    const task = tasks.find(t => t.id === taskId)
    if (!task) return false
    if (!task.parent_id) return false

    const parent = tasks.find(t => t.id === task.parent_id)
    if (!parent) return false

    if (state.collapsedTasks.has(parent.id)) return true
    return this.isTaskHidden(parent.id, tasks, visited)
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

  toggleTaskList() {
    state.showTaskList = !state.showTaskList
  },

  /**
   * 切换平移模式（手形工具）
   */
  togglePanMode() {
    state.panMode = !state.panMode
    eventBus.emit(GanttEvents.PAN_MODE_CHANGED, { enabled: state.panMode })
  },

  /**
   * 设置时间轴滚动位置
   */
  setScrollLeft(scrollLeft) {
    state.scrollLeft = scrollLeft
    eventBus.emit(GanttEvents.TIMELINE_SCROLLED, { scrollLeft })
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
