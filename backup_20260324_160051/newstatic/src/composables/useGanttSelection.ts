/**
 * 甘特图选择管理 Composable
 * 处理任务选择、多选、范围选择等功能
 */

import { ref, computed, Ref, watch } from 'vue'
import type { GanttTask } from '@/types/gantt'

/**
 * 选择模式
 */
export type SelectionMode = 'single' | 'multiple' | 'range'

/**
 * 选择状态
 */
export interface SelectionState {
  selectedTaskIds: Set<string | number>
  lastSelectedTaskId: string | number | null
  anchorTaskId: string | number | null  // 范围选择的起始点
}

/**
 * 使用甘特图选择功能
 * @param options 配置选项
 * @returns 选择状态和方法
 */
export function useGanttSelection(options: {
  tasks: Ref<GanttTask[]>
  mode?: Ref<SelectionMode>
  onSelectionChange?: (selectedTasks: GanttTask[]) => void
}) {
  const { tasks, mode = ref('single'), onSelectionChange } = options

  // 选择状态
  const selectedTaskIds = ref<Set<string | number>>(new Set())
  const lastSelectedTaskId = ref<string | number | null>(null)
  const anchorTaskId = ref<string | number | null>(null)

  /**
   * 当前选中的任务列表
   */
  const selectedTasks = computed(() => {
    return tasks.value.filter(task => selectedTaskIds.value.has(task.id))
  })

  /**
   * 第一个选中的任务
   */
  const firstSelectedTask = computed(() => {
    if (selectedTaskIds.value.size === 0) return null
    const firstId = Array.from(selectedTaskIds.value)[0]
    return tasks.value.find(t => t.id === firstId) || null
  })

  /**
   * 是否已选中指定任务
   */
  function isSelected(taskId: string | number): boolean {
    return selectedTaskIds.value.has(taskId)
  }

  /**
   * 选择单个任务
   */
  function selectTask(taskId: string | number, addToSelection = false) {
    if (mode.value === 'single' || !addToSelection) {
      // 单选模式：清除其他选择
      selectedTaskIds.value.clear()
      selectedTaskIds.value.add(taskId)
      anchorTaskId.value = taskId
    } else {
      // 多选模式：切换选择状态
      if (selectedTaskIds.value.has(taskId)) {
        selectedTaskIds.value.delete(taskId)
      } else {
        selectedTaskIds.value.add(taskId)
      }
    }

    lastSelectedTaskId.value = taskId
    notifyChange()
  }

  /**
   * 选择任务范围
   */
  function selectRange(startTaskId: string | number, endTaskId: string | number) {
    const startIndex = tasks.value.findIndex(t => t.id === startTaskId)
    const endIndex = tasks.value.findIndex(t => t.id === endTaskId)

    if (startIndex === -1 || endIndex === -1) return

    const [from, to] = startIndex < endIndex ? [startIndex, endIndex] : [endIndex, startIndex]

    selectedTaskIds.value.clear()

    for (let i = from; i <= to; i++) {
      selectedTaskIds.value.add(tasks.value[i].id)
    }

    lastSelectedTaskId.value = endTaskId
    notifyChange()
  }

  /**
   * 选择所有任务
   */
  function selectAll() {
    selectedTaskIds.value.clear()
    tasks.value.forEach(task => {
      selectedTaskIds.value.add(task.id)
    })
    lastSelectedTaskId.value = tasks.value[tasks.value.length - 1]?.id || null
    notifyChange()
  }

  /**
   * 清除所有选择
   */
  function clearSelection() {
    selectedTaskIds.value.clear()
    lastSelectedTaskId.value = null
    anchorTaskId.value = null
    notifyChange()
  }

  /**
   * 反选
   */
  function invertSelection() {
    const newSelection = new Set<string | number>()

    tasks.value.forEach(task => {
      if (!selectedTaskIds.value.has(task.id)) {
        newSelection.add(task.id)
      }
    })

    selectedTaskIds.value = newSelection
    notifyChange()
  }

  /**
   * 处理点击事件（支持 Shift 范围选择、Ctrl 多选）
   */
  function handleClick(taskId: string | number, event: MouseEvent) {
    const isShiftKey = event.shiftKey
    const isCtrlKey = event.ctrlKey || event.metaKey
    const isAltKey = event.altKey

    if (mode.value === 'single') {
      selectTask(taskId, false)
    } else if (isShiftKey && anchorTaskId.value !== null) {
      // Shift+点击：范围选择
      selectRange(anchorTaskId.value, taskId)
    } else if (isCtrlKey || isMetaKey) {
      // Ctrl+点击：多选切换
      if (selectedTaskIds.value.size === 0 || !selectedTaskIds.value.has(taskId)) {
        anchorTaskId.value = taskId
      }
      selectTask(taskId, true)
    } else if (isAltKey) {
      // Alt+点击：仅设置锚点，不改变选择
      anchorTaskId.value = taskId
    } else {
      // 普通点击：单选
      selectTask(taskId, false)
    }
  }

  /**
   * 获取选中任务的统计信息
   */
  function getSelectionStats() {
    const selected = selectedTasks.value

    return {
      total: selected.length,
      completed: selected.filter(t => t.status === 'completed').length,
      inProgress: selected.filter(t => t.status === 'in_progress').length,
      notStarted: selected.filter(t => t.status === 'not_started').length,
      delayed: selected.filter(t => t.status === 'delayed').length,
      critical: selected.filter(t => t.is_critical).length,
      totalProgress: selected.reduce((sum, t) => sum + t.progress, 0),
      avgProgress: selected.length > 0
        ? Math.round(selected.reduce((sum, t) => sum + t.progress, 0) / selected.length)
        : 0
    }
  }

  /**
   * 获取选中任务的时间范围
   */
  function getSelectionTimeRange() {
    if (selectedTasks.value.length === 0) {
      return { start: null, end: null }
    }

    let minDate = new Date(selectedTasks.value[0].start)
    let maxDate = new Date(selectedTasks.value[0].end)

    selectedTasks.value.forEach(task => {
      const start = new Date(task.start)
      const end = new Date(task.end)

      if (start < minDate) minDate = start
      if (end > maxDate) maxDate = end
    })

    return {
      start: minDate.toISOString().split('T')[0],
      end: maxDate.toISOString().split('T')[0]
    }
  }

  /**
   * 通知选择变化
   */
  function notifyChange() {
    if (onSelectionChange) {
      onSelectionChange(selectedTasks.value)
    }
  }

  // 监听任务列表变化，清理无效的选择
  watch(
    () => tasks.value.map(t => t.id),
    () => {
      const validIds = new Set(tasks.value.map(t => t.id))
      const invalidIds: Array<string | number> = []

      selectedTaskIds.value.forEach(id => {
        if (!validIds.has(id)) {
          invalidIds.push(id)
        }
      })

      invalidIds.forEach(id => {
        selectedTaskIds.value.delete(id)
      })

      if (invalidIds.length > 0) {
        notifyChange()
      }
    },
    { deep: true }
  )

  return {
    // 状态
    selectedTaskIds,
    selectedTasks,
    firstSelectedTask,
    lastSelectedTaskId,
    anchorTaskId,

    // 方法
    isSelected,
    selectTask,
    selectRange,
    selectAll,
    clearSelection,
    invertSelection,
    handleClick,
    getSelectionStats,
    getSelectionTimeRange
  }
}

/**
 * 使用键盘导航选择
 * 配合 useGanttKeyboard 使用
 */
export function useKeyboardNavigation(options: {
  tasks: Ref<GanttTask[]>
  onNavigate: (taskId: string | number) => void
}) {
  const { tasks, onNavigate } = options

  /**
   * 根据方向查找下一个任务
   */
  function findNextTask(
    currentTaskId: string | number,
    direction: 'up' | 'down' | 'left' | 'right'
  ): GanttTask | null {
    const currentIndex = tasks.value.findIndex(t => t.id === currentTaskId)
    if (currentIndex === -1) return null

    const currentTask = tasks.value[currentIndex]

    switch (direction) {
      case 'up':
        return tasks.value[currentIndex - 1] || null
      case 'down':
        return tasks.value[currentIndex + 1] || null
      case 'left':
        // 查找前驱任务
        if (currentTask.predecessors && currentTask.predecessors.length > 0) {
          const predecessorId = currentTask.predecessors[0].predecessor_id
          return tasks.value.find(t => t.id === predecessorId) || null
        }
        return null
      case 'right':
        // 查找后继任务
        if (currentTask.successors && currentTask.successors.length > 0) {
          const successorId = currentTask.successors[0].successor_id
          return tasks.value.find(t => t.id === successorId) || null
        }
        return null
    }
  }

  /**
   * 导航到指定方向的任务
   */
  function navigate(currentTaskId: string | number, direction: 'up' | 'down' | 'left' | 'right') {
    const nextTask = findNextTask(currentTaskId, direction)
    if (nextTask) {
      onNavigate(nextTask.id)
    }
  }

  return {
    navigate,
    findNextTask
  }
}
