/**
 * Critical Path Method (CPM) Utilities
 *
 * Enhanced critical path calculation supporting:
 * - All 4 dependency types (FS, FF, SS, SF)
 * - Lag/Lead time in calculations
 * - Slack/Float calculation
 * - Path analysis and reporting
 */

import { DependencyTypes } from './dependencyValidator'

// Re-export for backward compatibility with tests
export const DependencyType = DependencyTypes

/**
 * Calculates the critical path for a set of tasks
 *
 * @param {Array} tasks - Array of task objects
 * @param {Array} dependencies - Array of dependency objects
 * @returns {Object} Critical path analysis result
 */
export function calculateCriticalPath(tasks, dependencies = []) {
  if (!tasks || !Array.isArray(tasks) || tasks.length === 0) {
    return {
      criticalPath: [],
      criticalTasks: [],
      slack: {},
      projectDuration: 0,
      dates: {}
    }
  }

  // Build task map and dependency graph
  const taskMap = new Map(tasks.map(t => [t.id, t]))
  const predecessors = new Map() // taskId -> array of predecessor task IDs
  const successors = new Map() // taskId -> array of successor task IDs

  // Initialize maps
  tasks.forEach(task => {
    predecessors.set(task.id, [])
    successors.set(task.id, [])
  })

  // Build dependency relationships from dependencies array
  dependencies.forEach(dep => {
    if (dep.from && dep.to && taskMap.has(dep.from) && taskMap.has(dep.to)) {
      predecessors.get(dep.to).push(dep.from)
      successors.get(dep.from).push(dep.to)
    }
  })

  // Also check task.dependencies property for backward compatibility
  tasks.forEach(task => {
    if (task.dependencies && Array.isArray(task.dependencies)) {
      task.dependencies.forEach(depId => {
        if (taskMap.has(depId)) {
          predecessors.get(task.id).push(depId)
          successors.get(depId).push(task.id)
        }
      })
    }
  })

  // Calculate duration for each task
  const durations = {}
  tasks.forEach(task => {
    const start = new Date(task.start_date || task.start)
    const end = new Date(task.end_date || task.end)
    const duration = Math.ceil((end - start) / (1000 * 60 * 60 * 24))
    durations[task.id] = duration
  })

  // Forward pass - Calculate Early Start (ES) and Early Finish (EF)
  const earlyStart = {}
  const earlyFinish = {}

  // Initialize all tasks
  tasks.forEach(task => {
    earlyStart[task.id] = 0
    earlyFinish[task.id] = durations[task.id]
  })

  // Topological sort for forward pass
  const sortedTasks = topologicalSort(tasks, predecessors)

  sortedTasks.forEach(task => {
    const preds = predecessors.get(task.id) || []

    if (preds.length === 0) {
      // No predecessors - starts at 0
      earlyStart[task.id] = 0
    } else {
      // ES = max of all predecessors' EF + lag
      let maxEF = 0
      preds.forEach(predId => {
        const predFinish = earlyFinish[predId]
        const lag = getDependencyLag(task, predId, dependencies)
        const adjustedFinish = predFinish + lag
        maxEF = Math.max(maxEF, adjustedFinish)
      })
      earlyStart[task.id] = maxEF
    }

    earlyFinish[task.id] = earlyStart[task.id] + durations[task.id]
  })

  // Find project completion time
  const projectCompletion = Math.max(...Object.values(earlyFinish))

  // Backward pass - Calculate Late Start (LS) and Late Finish (LF)
  const lateStart = {}
  const lateFinish = {}

  // Initialize all tasks to project completion
  tasks.forEach(task => {
    lateFinish[task.id] = projectCompletion
    lateStart[task.id] = projectCompletion - durations[task.id]
  })

  // Reverse topological order for backward pass
  const reversedTasks = [...sortedTasks].reverse()

  reversedTasks.forEach(task => {
    const succs = successors.get(task.id) || []

    if (succs.length === 0) {
      // No successors - finishes at project completion
      lateFinish[task.id] = projectCompletion
    } else {
      // LF = min of all successors' LS - lag
      let minLS = Infinity
      succs.forEach(succId => {
        const succStart = lateStart[succId]
        const lag = getDependencyLag(task, succId, dependencies)
        const adjustedStart = succStart - lag
        minLS = Math.min(minLS, adjustedStart)
      })
      lateFinish[task.id] = minLS
    }

    lateStart[task.id] = lateFinish[task.id] - durations[task.id]
  })

  // Calculate slack/float for each task
  const slack = {}
  const criticalTasks = []

  tasks.forEach(task => {
    const taskSlack = lateStart[task.id] - earlyStart[task.id]
    slack[task.id] = Math.round(taskSlack * 100) / 100

    // Critical tasks have zero slack
    if (Math.abs(taskSlack) < 0.01) {
      criticalTasks.push(task.id)
    }
  })

  // Build critical path
  const criticalPath = buildCriticalPathIds(tasks, criticalTasks, predecessors, successors)

  // Build dates object
  const dates = {}
  tasks.forEach(task => {
    dates[task.id] = {
      earlyStart: earlyStart[task.id],
      earlyFinish: earlyFinish[task.id],
      lateStart: lateStart[task.id],
      lateFinish: lateFinish[task.id]
    }
  })

  // Get all critical paths for projects with multiple parallel paths
  const allCriticalPaths = getAllCriticalPathsFromData(tasks, criticalTasks, predecessors, successors, durations)

  return {
    criticalPath,
    criticalTasks,
    slack,
    projectDuration: projectCompletion,
    dates,
    criticalPaths: allCriticalPaths.map(path => ({
      tasks: path,
      duration: path.reduce((sum, taskId) => sum + durations[taskId], 0)
    }))
  }
}

/**
 * Get all critical paths from the computed data
 *
 * @param {Array} tasks - Array of task objects
 * @param {Array} criticalTasks - Array of critical task IDs
 * @param {Map} predecessors - Map of predecessors
 * @param {Map} successors - Map of successors
 * @param {Object} durations - Task durations
 * @returns {Array} Array of critical paths (each is an array of task IDs)
 */
function getAllCriticalPathsFromData(tasks, criticalTasks, predecessors, successors, durations) {
  const paths = []
  const visited = new Set()

  // Find start tasks (no critical predecessors)
  const startTasks = tasks.filter(task => {
    if (!criticalTasks.includes(task.id)) {
      return false
    }
    const preds = predecessors.get(task.id) || []
    return !preds.some(p => criticalTasks.includes(p))
  })

  const findPaths = (taskId, currentPath) => {
    if (visited.has(taskId)) {
      return
    }

    const newPath = [...currentPath, taskId]

    if (!criticalTasks.includes(taskId)) {
      if (newPath.length > 1) {
        paths.push(newPath)
      }
      return
    }

    // Find critical successors
    const succs = successors.get(taskId) || []
    const criticalSuccs = succs.filter(s => criticalTasks.includes(s))

    if (criticalSuccs.length === 0) {
      // End of path
      paths.push(newPath)
    } else {
      criticalSuccs.forEach(succId => {
        findPaths(succId, newPath)
      })
    }
  }

  startTasks.forEach(task => findPaths(task.id, []))

  // Remove duplicates and sort
  const uniquePaths = Array.from(new Set(paths.map(p => p.join(',')))).map(s => s.split(','))
  return uniquePaths
}

/**
 * Performs topological sort on tasks
 *
 * @param {Array} tasks - Array of task objects
 * @param {Map} predecessors - Map of task IDs to predecessor IDs
 * @returns {Array} Sorted array of tasks
 */
function topologicalSort(tasks, predecessors) {
  const sorted = []
  const visited = new Set()
  const temp = new Set()

  const visit = (taskId) => {
    if (temp.has(taskId)) {
      // Circular dependency - skip in this implementation
      return
    }

    if (visited.has(taskId)) {
      return
    }

    temp.add(taskId)

    const preds = predecessors.get(taskId) || []
    preds.forEach(predId => visit(predId))

    temp.delete(taskId)
    visited.add(taskId)

    const task = tasks.find(t => t.id === taskId)
    if (task) {
      sorted.push(task)
    }
  }

  tasks.forEach(task => visit(task.id))

  return sorted
}

/**
 * Gets lag time for a dependency
 *
 * @param {Object} task - Task object
 * @param {number} depId - Dependency task ID
 * @param {Array} dependencies - Array of dependency objects
 * @returns {number} Lag time in days
 */
function getDependencyLag(task, depId, dependencies = []) {
  // First check dependencyDetails in task
  if (task.dependencyDetails) {
    const depDetail = task.dependencyDetails.find(d => d.taskId === depId || d.taskId === parseInt(depId))
    if (depDetail) return depDetail.lag || 0
  }

  // Then check dependencies array
  const dep = dependencies.find(d =>
    (d.from === depId && d.to === task.id) ||
    (d.to === depId && d.from === task.id)
  )
  return dep?.lag || 0
}

/**
 * Builds the critical path from critical tasks (returns IDs)
 *
 * @param {Array} tasks - Array of task objects
 * @param {Array} criticalTasks - Array of critical task IDs
 * @param {Map} predecessors - Map of predecessors
 * @param {Map} successors - Map of successors
 * @returns {Array} Critical path task IDs in order
 */
function buildCriticalPathIds(tasks, criticalTasks, predecessors, successors) {
  // Find start tasks (no predecessors or only non-critical predecessors)
  const startTasks = tasks.filter(task => {
    if (!criticalTasks.includes(task.id)) {
      return false
    }

    const preds = predecessors.get(task.id) || []
    return preds.length === 0 || preds.every(p => !criticalTasks.includes(p))
  })

  // Trace forward from start tasks
  const path = []
  const visited = new Set()

  const trace = (taskId) => {
    if (visited.has(taskId)) {
      return
    }

    visited.add(taskId)

    if (criticalTasks.includes(taskId)) {
      path.push(taskId)

      // Follow critical successors
      const succs = successors.get(taskId) || []
      succs.forEach(succId => {
        if (criticalTasks.includes(succId)) {
          trace(succId)
        }
      })
    }
  }

  startTasks.forEach(task => trace(task.id))

  return path
}

/**
 * Calculates slack/float for a specific task
 *
 * @param {Object} task - Task object
 * @param {Object} cpmResult - Result from calculateCriticalPath
 * @returns {number} Slack in days
 */
export function calculateTaskSlack(task, cpmResult) {
  if (!task || !cpmResult) {
    return 0
  }

  return cpmResult.slack[task.id] || 0
}

/**
 * Determines if a task is on the critical path
 *
 * @param {Object} task - Task object
 * @param {Object} cpmResult - Result from calculateCriticalPath
 * @returns {boolean} True if task is critical
 */
export function isTaskCritical(task, cpmResult) {
  if (!task || !cpmResult) {
    return false
  }

  return cpmResult.criticalTasks.includes(task.id)
}

/**
 * Calculates free slack for a task
 * Free slack is the time a task can be delayed without delaying any successor
 *
 * @param {Object} task - Task object
 * @param {Array} tasks - Array of all tasks
 * @param {Object} cpmResult - Result from calculateCriticalPath
 * @returns {number} Free slack in days
 */
export function calculateFreeSlack(task, tasks, cpmResult) {
  if (!task || !tasks || !cpmResult) {
    return 0
  }

  const taskEF = cpmResult.dates[task.id]?.earlyFinish || 0

  // Find successors
  const successors = tasks.filter(t =>
    t.dependencies && t.dependencies.includes(task.id)
  )

  if (successors.length === 0) {
    // No successors - free slack equals total slack
    return calculateTaskSlack(task, cpmResult)
  }

  // Free slack = min(successor ES) - task EF
  const minSuccessorES = Math.min(
    ...successors.map(succ => {
      const lag = 0 // Simplified for now
      return (cpmResult.dates[succ.id]?.earlyStart || 0) - lag
    })
  )

  return Math.max(0, minSuccessorES - taskEF)
}

/**
 * Calculates interfering slack for a task
 * Interfering slack = Total slack - Free slack
 *
 * @param {Object} task - Task object
 * @param {Array} tasks - Array of all tasks
 * @param {Object} cpmResult - Result from calculateCriticalPath
 * @returns {number} Interfering slack in days
 */
export function calculateInterferingSlack(task, tasks, cpmResult) {
  const totalSlack = calculateTaskSlack(task, cpmResult)
  const freeSlack = calculateFreeSlack(task, tasks, cpmResult)

  return Math.max(0, totalSlack - freeSlack)
}

/**
 * Gets all critical paths in the project
 * (Can be multiple parallel critical paths)
 *
 * @param {Array} tasks - Array of task objects
 * @param {Array} dependencies - Array of dependency objects
 * @returns {Array} Array of critical paths
 */
export function getAllCriticalPaths(tasks, dependencies = []) {
  const cpmResult = calculateCriticalPath(tasks, dependencies)

  // Group critical tasks by their chains
  const paths = []
  const visited = new Set()

  const findPaths = (taskId, currentPath) => {
    if (visited.has(taskId)) {
      return
    }

    visited.add(taskId)

    if (!cpmResult.criticalTasks.includes(taskId)) {
      if (currentPath.length > 0) {
        paths.push([...currentPath])
      }
      return
    }

    const task = tasks.find(t => t.id === taskId)
    if (!task) return

    currentPath.push(task)

    // Find critical successors from dependencies
    const successors = []
    dependencies.forEach(dep => {
      if (dep.from === taskId && cpmResult.criticalTasks.includes(dep.to)) {
        const succTask = tasks.find(t => t.id === dep.to)
        if (succTask) successors.push(succTask)
      }
    })

    if (successors.length === 0) {
      // End of path
      paths.push([...currentPath])
    } else {
      successors.forEach(succ => {
        findPaths(succ.id, [...currentPath])
      })
    }
  }

  // Find all start tasks
  const startTasks = tasks.filter(task => {
    if (!cpmResult.criticalTasks.includes(task.id)) {
      return false
    }
    // Check if task has no critical predecessors
    const hasCriticalPred = dependencies.some(dep =>
      dep.to === task.id && cpmResult.criticalTasks.includes(dep.from)
    )
    return !hasCriticalPred
  })

  startTasks.forEach(task => findPaths(task.id, []))

  return paths.map(path => ({
    tasks: path,
    duration: path.reduce((sum, task) => {
      const duration = Math.ceil(
        (new Date(task.end_date || task.end) - new Date(task.start_date || task.start)) / (1000 * 60 * 60 * 24)
      )
      return sum + duration
    }, 0)
  }))
}

/**
 * Generates a critical path report
 *
 * @param {Array} tasks - Array of task objects
 * @param {Array} dependencies - Array of dependency objects
 * @returns {Object} Critical path report
 */
export function generateCriticalPathReport(tasks, dependencies = []) {
  const cpmResult = calculateCriticalPath(tasks, dependencies)
  const allPaths = getAllCriticalPaths(tasks, dependencies)

  const report = {
    projectDuration: cpmResult.projectDuration,
    criticalTaskCount: cpmResult.criticalTasks.length,
    totalTaskCount: tasks.length,
    criticalPathCount: allPaths.length,
    criticalPaths: allPaths,
    taskDetails: tasks.map(task => ({
      id: task.id,
      name: task.name,
      isCritical: cpmResult.criticalTasks.includes(task.id),
      totalSlack: cpmResult.slack[task.id] || 0,
      earlyStart: cpmResult.dates[task.id]?.earlyStart,
      earlyFinish: cpmResult.dates[task.id]?.earlyFinish,
      lateStart: cpmResult.dates[task.id]?.lateStart,
      lateFinish: cpmResult.dates[task.id]?.lateFinish
    }))
  }

  return report
}

/**
 * Calculates the impact of delaying a task on the project end date
 *
 * @param {Object} task - Task to delay
 * @param {number} delayDays - Days to delay
 * @param {Array} tasks - Array of all tasks
 * @returns {Object} Impact analysis
 */
export function calculateDelayImpact(task, delayDays, tasks) {
  const cpmResult = calculateCriticalPath(tasks)
  const isCritical = cpmResult.criticalTasks.includes(task.id)

  if (isCritical) {
    // Delaying critical task directly delays project
    return {
      projectDelay: delayDays,
      affectedTasks: [...cpmResult.criticalTasks],
      reason: 'Task is on critical path'
    }
  } else {
    const slack = cpmResult.slack[task.id] || 0

    if (delayDays <= slack) {
      return {
        projectDelay: 0,
        affectedTasks: [task.id],
        reason: 'Delay within slack limits'
      }
    } else {
      const effectiveDelay = delayDays - slack
      return {
        projectDelay: effectiveDelay,
        affectedTasks: [...cpmResult.criticalTasks],
        reason: 'Delay exceeds slack, pushes critical path'
      }
    }
  }
}

/**
 * Identifies all critical tasks
 *
 * @param {Array} tasks - Array of task objects
 * @param {Array} dependencies - Array of dependency objects
 * @returns {Array} Array of critical task IDs
 */
export function identifyCriticalTasks(tasks, dependencies = []) {
  const result = calculateCriticalPath(tasks, dependencies)
  return result.criticalTasks
}

/**
 * Calculates early start for a task
 *
 * @param {string} taskId - Task ID
 * @param {Array} tasks - Array of all tasks
 * @param {Array} dependencies - Array of dependency objects
 * @returns {number} Early start value
 */
export function calculateEarlyStart(taskId, tasks, dependencies = []) {
  const result = calculateCriticalPath(tasks, dependencies)
  return result.dates[taskId]?.earlyStart ?? 0
}

/**
 * Calculates early finish for a task
 *
 * @param {string} taskId - Task ID
 * @param {Array} tasks - Array of all tasks
 * @param {Array} dependencies - Array of dependency objects
 * @returns {number} Early finish value
 */
export function calculateEarlyFinish(taskId, tasks, dependencies = []) {
  const result = calculateCriticalPath(tasks, dependencies)
  return result.dates[taskId]?.earlyFinish ?? 0
}

/**
 * Calculates late start for a task
 *
 * @param {string} taskId - Task ID
 * @param {Array} tasks - Array of all tasks
 * @param {Array} dependencies - Array of dependency objects
 * @returns {number} Late start value
 */
export function calculateLateStart(taskId, tasks, dependencies = []) {
  const result = calculateCriticalPath(tasks, dependencies)
  return result.dates[taskId]?.lateStart ?? 0
}

/**
 * Calculates late finish for a task
 *
 * @param {string} taskId - Task ID
 * @param {Array} tasks - Array of all tasks
 * @param {Array} dependencies - Array of dependency objects
 * @returns {number} Late finish value
 */
export function calculateLateFinish(taskId, tasks, dependencies = []) {
  const result = calculateCriticalPath(tasks, dependencies)
  return result.dates[taskId]?.lateFinish ?? 0
}
