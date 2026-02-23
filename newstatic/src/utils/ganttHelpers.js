/**
 * 甘特图工具函数
 * 用于依赖关系计算、路径绘制等功能
 */

import { formatDate } from './dateFormat'

/**
 * 计算依赖关系箭头的SVG路径 - 使用A*算法进行路径规划
 * @param {Object} fromTask - 源任务 {id, start, end, ...}
 * @param {Object} toTask - 目标任务 {id, start, end, ...}
 * @param {Object} timelineDays - 时间轴天数数组
 * @param {number} dayWidth - 每天的像素宽度
 * @param {number} fromRow - 源任务行索引
 * @param {number} toRow - 目标任务行索引
 * @param {number} rowHeight - 行高
 * @param {boolean} isCritical - 是否为关键路径
 * @param {Array} allTasks - 所有任务数组（用于障碍物检测）
 * @returns {Object} { path, startX, startY, endX, endY }
 */
export function calculateDependencyPath(
  fromTask,
  toTask,
  timelineDays,
  dayWidth,
  fromRow,
  toRow,
  rowHeight = 60,
  isCritical = false,
  allTasks = []
) {
  // 获取时间轴起始日期
  const timelineStart = timelineDays[0]?.date
  if (!timelineStart) return null

  const timelineStartDate = new Date(timelineStart)

  // 计算源任务位置
  const fromTaskStart = new Date(fromTask.start)
  const fromTaskEnd = new Date(fromTask.end)
  const fromDaysDiff = Math.ceil((fromTaskStart - timelineStartDate) / (1000 * 60 * 60 * 24))
  const fromDuration = Math.ceil((fromTaskEnd - fromTaskStart) / (1000 * 60 * 60 * 24))

  const startX = fromDaysDiff * dayWidth + fromDuration * dayWidth
  const startY = fromRow * rowHeight + rowHeight / 2

  // 计算目标任务位置
  const toTaskStart = new Date(toTask.start)
  const toDaysDiff = Math.ceil((toTaskStart - timelineStartDate) / (1000 * 60 * 60 * 24))

  const endX = toDaysDiff * dayWidth
  const endY = toRow * rowHeight + rowHeight / 2

  // 同一行：直线连接
  if (Math.abs(startY - endY) < 1) {
    return {
      path: `M ${startX} ${startY} L ${endX} ${endY}`,
      startX,
      startY,
      endX,
      endY,
      isCritical
    }
  }

  // 构建障碍物地图（将任务条位置转换为障碍物）
  const obstacles = buildObstacleMap(allTasks, dayWidth, rowHeight, timelineStartDate, fromTask.id, toTask.id)

  // 使用A*算法计算路径
  const path = findOrthogonalPath(startX, startY, endX, endY, obstacles, dayWidth)

  return {
    path,
    startX,
    startY,
    endX,
    endY,
    isCritical
  }
}

/**
 * 构建障碍物地图
 * @param {Array} tasks - 所有任务
 * @param {number} dayWidth - 每天宽度
 * @param {number} rowHeight - 行高
 * @param {Date} timelineStart - 时间轴起始日期
 * @param {number} excludeFromId - 排除的源任务ID
 * @param {number} excludeToId - 排除的目标任务ID
 * @returns {Array} 障碍物数组
 */
function buildObstacleMap(tasks, dayWidth, rowHeight, timelineStart, excludeFromId, excludeToId) {
  const obstacles = []
  const padding = 10 // 障碍物周围的填充

  for (const task of tasks) {
    // 跳过源任务和目标任务
    if (task.id === excludeFromId || task.id === excludeToId) continue

    const taskStart = new Date(task.start)
    const taskEnd = new Date(task.end)
    const daysDiff = Math.ceil((taskStart - timelineStart) / (1000 * 60 * 60 * 24))
    const duration = Math.ceil((taskEnd - taskStart) / (1000 * 60 * 60 * 24))

    const row = tasks.findIndex(t => t.id === task.id)
    if (row === -1) continue

    const x = daysDiff * dayWidth
    const y = row * rowHeight
    const width = duration * dayWidth
    const height = rowHeight

    obstacles.push({
      minX: x - padding,
      maxX: x + width + padding,
      minY: y - padding,
      maxY: y + height + padding
    })
  }

  return obstacles
}

/**
 * 使用简化的A*算法查找直角路径
 * @param {number} startX - 起点X
 * @param {number} startY - 起点Y
 * @param {number} endX - 终点X
 * @param {number} endY - 终点Y
 * @param {Array} obstacles - 障碍物数组
 * @returns {string} SVG路径字符串
 */
export function findOrthogonalPath(startX, startY, endX, endY, obstacles = []) {
  // 如果没有障碍物或起点终点之间没有障碍，使用简单折线
  const directPath = tryDirectPath(startX, startY, endX, endY, obstacles)
  if (directPath) {
    return directPath
  }

  // 尝试几种常见的绕行策略
  const strategies = [
    // 策略1: 从上方绕过
    () => tryRouteAbove(startX, startY, endX, endY, obstacles),
    // 策略2: 从下方绕过
    () => tryRouteBelow(startX, startY, endX, endY, obstacles),
    // 策略3: 从右侧绕过
    () => tryRouteRight(startX, startY, endX, endY, obstacles),
    // 策略4: 避开中间障碍物的Z型路径
    () => tryZigzagPath(startX, startY, endX, endY, obstacles)
  ]

  for (const strategy of strategies) {
    const path = strategy()
    if (path) return path
  }

  // 如果所有策略都失败，返回最简单的直接路径（可能会穿过障碍物）
  return `M ${startX} ${startY} L ${(startX + endX) / 2} ${startY} L ${(startX + endX) / 2} ${endY} L ${endX} ${endY}`
}

/**
 * 尝试直接路径（中间点）
 */
function tryDirectPath(startX, startY, endX, endY, obstacles) {
  const midX = (startX + endX) / 2

  // 检查水平路径是否安全
  if (!isPathBlocked(startX, startY, midX, startY, obstacles) &&
      !isPathBlocked(midX, startY, midX, endY, obstacles) &&
      !isPathBlocked(midX, endY, endX, endY, obstacles)) {
    return `M ${startX} ${startY} L ${midX} ${startY} L ${midX} ${endY} L ${endX} ${endY}`
  }

  return null
}

/**
 * 尝试从上方绕过
 */
function tryRouteAbove(startX, startY, endX, endY, obstacles) {
  const minY = Math.min(startY, endY) - 40

  if (!isPathBlocked(startX, startY, startX, minY, obstacles) &&
      !isPathBlocked(startX, minY, endX, minY, obstacles) &&
      !isPathBlocked(endX, minY, endX, endY, obstacles)) {
    return `M ${startX} ${startY} L ${startX} ${minY} L ${endX} ${minY} L ${endX} ${endY}`
  }

  return null
}

/**
 * 尝试从下方绕过
 */
function tryRouteBelow(startX, startY, endX, endY, obstacles) {
  const maxY = Math.max(startY, endY) + 40

  if (!isPathBlocked(startX, startY, startX, maxY, obstacles) &&
      !isPathBlocked(startX, maxY, endX, maxY, obstacles) &&
      !isPathBlocked(endX, maxY, endX, endY, obstacles)) {
    return `M ${startX} ${startY} L ${startX} ${maxY} L ${endX} ${maxY} L ${endX} ${endY}`
  }

  return null
}

/**
 * 尝试从右侧绕过
 */
function tryRouteRight(startX, startY, endX, endY, obstacles) {
  const maxX = Math.max(startX, endX) + 60

  if (!isPathBlocked(startX, startY, maxX, startY, obstacles) &&
      !isPathBlocked(maxX, startY, maxX, endY, obstacles) &&
      !isPathBlocked(maxX, endY, endX, endY, obstacles)) {
    return `M ${startX} ${startY} L ${maxX} ${startY} L ${maxX} ${endY} L ${endX} ${endY}`
  }

  return null
}

/**
 * 尝试Z型路径避开中间障碍物
 */
function tryZigzagPath(startX, startY, endX, endY, obstacles) {
  // 多个转折点的Z型路径
  const midY = (startY + endY) / 2
  const offset = Math.abs(endX - startX) * 0.3

  // 向右偏移的Z型
  if (endX > startX) {
    const x1 = startX + offset
    const x2 = endX - offset

    if (!isPathBlocked(startX, startY, x1, startY, obstacles) &&
        !isPathBlocked(x1, startY, x1, midY, obstacles) &&
        !isPathBlocked(x1, midY, x2, midY, obstacles) &&
        !isPathBlocked(x2, midY, x2, endY, obstacles) &&
        !isPathBlocked(x2, endY, endX, endY, obstacles)) {
      return `M ${startX} ${startY} L ${x1} ${startY} L ${x1} ${midY} L ${x2} ${midY} L ${x2} ${endY} L ${endX} ${endY}`
    }
  }

  return null
}

/**
 * 检查路径是否被障碍物阻挡
 * @param {number} x1 - 起点X
 * @param {number} y1 - 起点Y
 * @param {number} x2 - 终点X
 * @param {number} y2 - 终点Y
 * @param {Array} obstacles - 障碍物数组
 * @returns {boolean} 是否被阻挡
 */
function isPathBlocked(x1, y1, x2, y2, obstacles) {
  // 检查线段是否与任何障碍物相交
  for (const obs of obstacles) {
    if (lineIntersectsRect(x1, y1, x2, y2, obs)) {
      return true
    }
  }
  return false
}

/**
 * 检查线段是否与矩形相交
 */
function lineIntersectsRect(x1, y1, x2, y2, rect) {
  // 检查线段端点是否在矩形内
  if ((x1 >= rect.minX && x1 <= rect.maxX && y1 >= rect.minY && y1 <= rect.maxY) ||
      (x2 >= rect.minX && x2 <= rect.maxX && y2 >= rect.minY && y2 <= rect.maxY)) {
    return true
  }

  // 检查水平线段
  if (y1 === y2) {
    return y1 >= rect.minY && y1 <= rect.maxY &&
           ((x1 <= rect.minX && x2 >= rect.minX) ||
            (x1 <= rect.maxX && x2 >= rect.maxX) ||
            (x1 >= rect.minX && x2 <= rect.maxX))
  }

  // 检查垂直线段
  if (x1 === x2) {
    return x1 >= rect.minX && x1 <= rect.maxX &&
           ((y1 <= rect.minY && y2 >= rect.minY) ||
            (y1 <= rect.maxY && y2 >= rect.maxY) ||
            (y1 >= rect.minY && y2 <= rect.maxY))
  }

  return false
}

/**
 * 获取依赖关系箭头的颜色
 * @param {boolean} isCritical - 是否为关键路径
 * @returns {string} 颜色值
 */
export function getDependencyArrowColor(isCritical) {
  return isCritical ? '#f56c6c' : '#909399'
}

/**
 * 计算箭头标记的路径
 * @param {number} x - 终点X坐标
 * @param {number} y - 终点Y坐标
 * @param {number} angle - 角度（弧度）
 * @param {number} size - 箭头大小
 * @returns {string} SVG路径
 */
export function createArrowMarker(x, y, angle, size = 8) {
  const cos = Math.cos(angle)
  const sin = Math.sin(angle)

  // 箭头三角形顶点
  const x1 = x - size * cos + size * 0.5 * sin
  const y1 = y - size * sin - size * 0.5 * cos
  const x2 = x - size * cos - size * 0.5 * sin
  const y2 = y - size * sin + size * 0.5 * cos

  return `M ${x} ${y} L ${x1} ${y1} L ${x2} ${y2} Z`
}

/**
 * 计算两个任务之间的依赖角度
 * @param {number} startX - 起点X
 * @param {number} startY - 起点Y
 * @param {number} endX - 终点X
 * @param {number} endY - 终点Y
 * @returns {number} 角度（弧度）
 */
export function calculateArrowAngle(startX, startY, endX, endY) {
  return Math.atan2(endY - startY, endX - startX)
}

/**
 * 检查依赖关系是否需要交叉线（避免重叠）
 * @param {Array} dependencies - 所有依赖关系
 * @param {Object} currentDep - 当前依赖关系
 * @returns {boolean} 是否需要交叉线
 */
export function needsCrossoverLine(dependencies, currentDep) {
  // 简单实现：检查是否有其他依赖线在同一水平位置
  const sameLevelDeps = dependencies.filter(dep => {
    return dep.fromRow === currentDep.fromRow && dep.id !== currentDep.id
  })
  return sameLevelDeps.length > 0
}

/**
 * 计算交叉线路径（避免箭头重叠）
 * @param {Object} dep - 依赖关系
 * @param {number} index - 索引
 * @returns {string} SVG路径
 */
export function calculateCrossoverPath(dep, index) {
  const offset = index * 15 // 每条线偏移15px
  const midY = (dep.startY + dep.endY) / 2

  return `M ${dep.startX} ${dep.startY} L ${dep.startX + 20} ${dep.startY} L ${dep.startX + 20} ${midY + offset} L ${dep.endX - 20} ${midY + offset} L ${dep.endX - 20} ${dep.endY} L ${dep.endX} ${dep.endY}`
}

/**
 * 获取任务在列表中的索引
 * @param {Array} tasks - 任务数组
 * @param {string} taskId - 任务ID
 * @returns {number} 索引，未找到返回-1
 */
export function getTaskIndex(tasks, taskId) {
  return tasks.findIndex(t => t.id === taskId)
}

/**
 * 过滤显示的依赖关系
 * @param {Object} scheduleData - 调度数据
 * @param {Array} tasks - 格式化后的任务数组
 * @param {Object} options - 选项 { showCriticalOnly, showNonCritical }
 * @returns {Array} 过滤后的依赖关系数组
 */
export function filterDependencies(scheduleData, tasks, options = {}) {
  const { showCriticalOnly = false, showNonCritical = true } = options
  const dependencies = []
  const activities = scheduleData.activities || {}

  console.log('filterDependencies - activities:', activities)
  console.log('filterDependencies - tasks:', tasks)

  for (const task of tasks) {
    // activities 的 key 是字符串格式 "task_3"，需要拼接
    const activityKey = `task_${task.id}`
    const activity = activities[activityKey]
    console.log(`任务 ${task.id} (${task.name}) 的 activityKey: ${activityKey}, activity:`, activity)

    if (!activity) {
      console.log(`任务 ${task.id} 没有 activity (查找 key: ${activityKey})`)
      continue
    }

    if (!activity.predecessors || activity.predecessors.length === 0) {
      console.log(`任务 ${task.id} 没有前置任务`)
      continue
    }

    console.log(`任务 ${task.id} 的前置任务:`, activity.predecessors)

    for (const predId of activity.predecessors) {
      const predTask = tasks.find(t => t.id === predId)
      if (!predTask) {
        console.log(`前置任务 ${predId} 不存在于任务列表中`)
        continue
      }

      // 检查是否是父子关系，如果是则跳过
      const isParentChild = checkParentChildRelation(predTask.id, task.id, tasks)
      if (isParentChild) {
        console.log(`跳过父子任务依赖关系: ${predTask.id} -> ${task.id}`)
        continue
      }

      const predActivityKey = `task_${predId}`
      const predActivity = activities[predActivityKey]
      const isCritical = activity.is_critical && predActivity?.is_critical

      // 根据选项过滤
      if (showCriticalOnly && !isCritical) continue
      if (!showNonCritical && !isCritical) continue

      dependencies.push({
        id: `${predId}-${task.id}`,
        fromTask: predTask,
        toTask: task,
        fromId: predId,
        toId: task.id,
        isCritical
      })
    }
  }

  console.log('filterDependencies - 最终依赖关系:', dependencies)
  return dependencies
}

/**
 * 检查两个任务之间是否存在父子关系
 * @param {number} taskId1 - 任务1的ID
 * @param {number} taskId2 - 任务2的ID
 * @param {Array} tasks - 任务数组
 * @returns {boolean} 是否存在父子关系
 */
function checkParentChildRelation(taskId1, taskId2, tasks) {
  const task1 = tasks.find(t => t.id === taskId1)
  const task2 = tasks.find(t => t.id === taskId2)

  if (!task1 || !task2) return false

  // 检查task1是否是task2的父任务
  if (task2.parent_id === taskId1) return true

  // 检查task2是否是task1的父任务
  if (task1.parent_id === taskId2) return true

  // 检查task1的子孙任务是否包含task2
  let current = task1
  while (current.parent_id) {
    if (current.parent_id === taskId2) return true
    current = tasks.find(t => t.id === current.parent_id)
    if (!current) break
  }

  // 检查task2的子孙任务是否包含task1
  current = task2
  while (current.parent_id) {
    if (current.parent_id === taskId1) return true
    current = tasks.find(t => t.id === current.parent_id)
    if (!current) break
  }

  return false
}

/**
 * 计算所有依赖关系的路径
 * @param {Array} dependencies - 依赖关系数组
 * @param {Array} tasks - 任务数组
 * @param {Array} timelineDays - 时间轴天数
 * @param {number} dayWidth - 每天宽度
 * @param {number} rowHeight - 行高
 * @returns {Array} 包含路径信息的数组
 */
export function calculateAllDependencyPaths(dependencies, tasks, timelineDays, dayWidth, rowHeight = 60) {
  return dependencies.map(dep => {
    const fromRow = getTaskIndex(tasks, dep.fromId)
    const toRow = getTaskIndex(tasks, dep.toId)

    if (fromRow === -1 || toRow === -1) return null

    return calculateDependencyPath(
      dep.fromTask,
      dep.toTask,
      timelineDays,
      dayWidth,
      fromRow,
      toRow,
      rowHeight,
      dep.isCritical,
      tasks // 传递所有任务用于障碍物检测
    )
  }).filter(Boolean)
}

/**
 * 判断两个日期范围是否重叠
 * @param {Date|string} start1 - 范围1开始
 * @param {Date|string} end1 - 范围1结束
 * @param {Date|string} start2 - 范围2开始
 * @param {Date|string} end2 - 范围2结束
 * @returns {boolean} 是否重叠
 */
export function isDateRangeOverlapping(start1, end1, start2, end2) {
  const s1 = new Date(start1)
  const e1 = new Date(end1)
  const s2 = new Date(start2)
  const e2 = new Date(end2)

  return s1 <= e2 && e1 >= s2
}

/**
 * 计算任务的实际工期（考虑工作日）
 * @param {Date|string} start - 开始日期
 * @param {Date|string} end - 结束日期
 * @param {Array} holidays - 节假日数组（可选）
 * @returns {number} 工作日数
 */
export function calculateWorkingDays(start, end, holidays = []) {
  const startDate = new Date(start)
  const endDate = new Date(end)
  let count = 0
  let current = new Date(startDate)

  while (current <= endDate) {
    const dayOfWeek = current.getDay()
    const dateStr = formatDate(current)

    // 0=周日, 6=周六
    const isWeekend = dayOfWeek === 0 || dayOfWeek === 6
    const isHoliday = holidays.includes(dateStr)

    if (!isWeekend && !isHoliday) {
      count++
    }

    current.setDate(current.getDate() + 1)
  }

  return count
}

/**
 * 检查任务是否为里程碑（工期为0）
 * @param {Object} task - 任务对象
 * @returns {boolean} 是否为里程碑
 */
export function isMilestone(task) {
  return task.duration === 0 || task.start === task.end
}

/**
 * 获取优先级权重（用于排序）
 * @param {string} priority - 优先级
 * @returns {number} 权重值
 */
export function getPriorityWeight(priority) {
  const weights = {
    urgent: 4,
    high: 3,
    medium: 2,
    low: 1
  }
  return weights[priority] || 0
}

/**
 * 按优先级排序任务
 * @param {Array} tasks - 任务数组
 * @returns {Array} 排序后的任务数组
 */
export function sortTasksByPriority(tasks) {
  return [...tasks].sort((a, b) => {
    return getPriorityWeight(b.priority) - getPriorityWeight(a.priority)
  })
}

/**
 * 按状态分组任务
 * @param {Array} tasks - 任务数组
 * @returns {Object} 分组后的任务对象
 */
export function groupTasksByStatus(tasks) {
  const groups = {
    completed: { name: '已完成', tasks: [] },
    in_progress: { name: '进行中', tasks: [] },
    not_started: { name: '未开始', tasks: [] },
    delayed: { name: '已延期', tasks: [] }
  }

  for (const task of tasks) {
    const status = task.status || 'not_started'
    if (groups[status]) {
      groups[status].tasks.push(task)
    }
  }

  return groups
}

/**
 * 按优先级分组任务
 * @param {Array} tasks - 任务数组
 * @returns {Object} 分组后的任务对象
 */
export function groupTasksByPriority(tasks) {
  const groups = {
    urgent: { name: '紧急', tasks: [] },
    high: { name: '高', tasks: [] },
    medium: { name: '中', tasks: [] },
    low: { name: '低', tasks: [] }
  }

  for (const task of tasks) {
    const priority = task.priority || 'medium'
    if (groups[priority]) {
      groups[priority].tasks.push(task)
    }
  }

  return groups
}

/**
 * 计算关键路径任务百分比
 * @param {Array} tasks - 任务数组
 * @returns {number} 百分比
 */
export function calculateCriticalPathPercentage(tasks) {
  if (tasks.length === 0) return 0
  const criticalCount = tasks.filter(t => t.is_critical).length
  return Math.round((criticalCount / tasks.length) * 100)
}

/**
 * 验证任务日期是否有效
 * @param {string} start - 开始日期
 * @param {string} end - 结束日期
 * @returns {boolean} 是否有效
 */
export function validateTaskDates(start, end) {
  const startDate = new Date(start)
  const endDate = new Date(end)

  if (isNaN(startDate.getTime()) || isNaN(endDate.getTime())) {
    return false
  }

  return startDate <= endDate
}

/**
 * 计算任务拖拽后的新日期
 * @param {string} originalDate - 原始日期
 * @param {number} dayOffset - 天数偏移量
 * @returns {string} 新日期（YYYY-MM-DD格式）
 */
export function calculateDraggedDate(originalDate, dayOffset) {
  const date = new Date(originalDate)
  date.setDate(date.getDate() + dayOffset)
  return formatDate(date)
}

/**
 * 限制日期在时间轴范围内
 * @param {string} date - 日期
 * @param {string} minDate - 最小日期
 * @param {string} maxDate - 最大日期
 * @returns {string} 限制后的日期
 */
export function clampDateToTimeline(date, minDate, maxDate) {
  const d = new Date(date)
  const min = new Date(minDate)
  const max = new Date(maxDate)

  if (d < min) return formatDate(min)
  if (d > max) return formatDate(max)
  return formatDate(d)
}

/**
 * 检测任务依赖循环
 * @param {Object} scheduleData - 调度数据
 * @param {string} taskId - 任务ID
 * @param {Set} visited - 已访问节点
 * @param {Set} recStack - 递归栈
 * @returns {boolean} 是否存在循环
 */
export function detectDependencyCycle(scheduleData, taskId, visited = new Set(), recStack = new Set()) {
  if (recStack.has(taskId)) return true
  if (visited.has(taskId)) return false

  visited.add(taskId)
  recStack.add(taskId)

  const activity = scheduleData.activities?.[taskId]
  if (activity?.successors) {
    for (const successorId of activity.successors) {
      if (detectDependencyCycle(scheduleData, successorId, visited, recStack)) {
        return true
      }
    }
  }

  recStack.delete(taskId)
  return false
}

/**
 * 获取周的开始日期（周一）
 * @param {Date} date - 日期
 * @returns {Date} 周开始日期
 */
export function getWeekStart(date) {
  const d = new Date(date)
  const day = d.getDay()
  const diff = d.getDate() - day + (day === 0 ? -6 : 1) // 周一
  return new Date(d.setDate(diff))
}

/**
 * 获取周的结束日期（周日）
 * @param {Date} date - 日期
 * @returns {Date} 周结束日期
 */
export function getWeekEnd(date) {
  const d = new Date(date)
  const day = d.getDay()
  const diff = d.getDate() - day + (day === 0 ? 0 : 7) // 周日
  return new Date(d.setDate(diff))
}

/**
 * 获取周数（ISO 8601）
 * @param {Date} date - 日期
 * @returns {number} 周数
 */
export function getWeekNumber(date) {
  const d = new Date(date)
  d.setHours(0, 0, 0, 0)
  d.setDate(d.getDate() + 4 - (d.getDay() || 7))
  const yearStart = new Date(d.getFullYear(), 0, 1)
  return Math.ceil((((d - yearStart) / 86400000) + 1) / 7)
}

/**
 * 获取季度
 * @param {Date} date - 日期
 * @returns {number} 季度 (1-4)
 */
export function getQuarter(date) {
  return Math.floor((date.getMonth() + 3) / 3)
}
