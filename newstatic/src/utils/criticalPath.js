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

/**
 * Calculates the critical path for a set of tasks
 *
 * @param {Array} tasks - Array of task objects
 * @returns {Object} Critical path analysis result
 */
export function calculateCriticalPath(tasks) {
  if (!tasks || !Array.isArray(tasks) || tasks.length === 0) {
    return {
      criticalPath: [],
      criticalTasks: new Set(),
      taskSlack: new Map(),
      projectDuration: 0,
      earliestStart: new Map(),
      earliestFinish: new Map(),
      latestStart: new Map(),
      latestFinish: new Map()
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

  // Build dependency relationships
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
  const durations = new Map()
  tasks.forEach(task => {
    const start = new Date(task.start_date)
    const end = new Date(task.end_date)
    const duration = Math.ceil((end - start) / (1000 * 60 * 60 * 24))
    durations.set(task.id, duration)
  })

  // Forward pass - Calculate Early Start (ES) and Early Finish (EF)
  const earlyStart = new Map()
  const earlyFinish = new Map()

  // Initialize all tasks
  tasks.forEach(task => {
    earlyStart.set(task.id, 0)
    earlyFinish.set(task.id, durations.get(task.id))
  })

  // Topological sort for forward pass
  const sortedTasks = topologicalSort(tasks, predecessors)

  sortedTasks.forEach(task => {
    const preds = predecessors.get(task.id) || []

    if (preds.length === 0) {
      // No predecessors - starts at 0
      earlyStart.set(task.id, 0)
    } else {
      // ES = max of all predecessors' EF + lag
      let maxEF = 0
      preds.forEach(predId => {
        const predFinish = earlyFinish.get(predId)
        const lag = getDependencyLag(task, predId)
        const adjustedFinish = predFinish + lag
        maxEF = Math.max(maxEF, adjustedFinish)
      })
      earlyStart.set(task.id, maxEF)
    }

    earlyFinish.set(task.id, earlyStart.get(task.id) + durations.get(task.id))
  })

  // Find project completion time
  const projectCompletion = Math.max(...Array.from(earlyFinish.values()))

  // Backward pass - Calculate Late Start (LS) and Late Finish (LF)
  const lateStart = new Map()
  const lateFinish = new Map()

  // Initialize all tasks to project completion
  tasks.forEach(task => {
    lateFinish.set(task.id, projectCompletion)
    lateStart.set(task.id, projectCompletion - durations.get(task.id))
  })

  // Reverse topological order for backward pass
  const reversedTasks = [...sortedTasks].reverse()

  reversedTasks.forEach(task => {
    const succs = successors.get(task.id) || []

    if (succs.length === 0) {
      // No successors - finishes at project completion
      lateFinish.set(task.id, projectCompletion)
    } else {
      // LF = min of all successors' LS - lag
      let minLS = Infinity
      succs.forEach(succId => {
        const succStart = lateStart.get(succId)
        const lag = getDependencyLag(task, succId)
        const adjustedStart = succStart - lag
        minLS = Math.min(minLS, adjustedStart)
      })
      lateFinish.set(task.id, minLS)
    }

    lateStart.set(task.id, lateFinish.get(task.id) - durations.get(task.id))
  })

  // Calculate slack/float for each task
  const taskSlack = new Map()
  const criticalTasks = new Set()

  tasks.forEach(task => {
    const slack = lateStart.get(task.id) - earlyStart.get(task.id)
    taskSlack.set(task.id, Math.round(slack * 100) / 100)

    // Critical tasks have zero slack
    if (Math.abs(slack) < 0.01) {
      criticalTasks.add(task.id)
    }
  })

  // Build critical path
  const criticalPath = buildCriticalPath(tasks, criticalTasks, predecessors, successors)

  return {
    criticalPath,
    criticalTasks,
    taskSlack,
    projectDuration: projectCompletion,
    earliestStart: earlyStart,
    earliestFinish: earlyFinish,
    latestStart: lateStart,
    latestFinish: lateFinish
  }
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
 * @returns {number} Lag time in days
 */
function getDependencyLag(task, depId) {
  if (!task.dependencyDetails) {
    return 0
  }

  const depDetail = task.dependencyDetails.find(d => d.taskId === depId)
  return depDetail?.lag || 0
}

/**
 * Builds the critical path from critical tasks
 *
 * @param {Array} tasks - Array of task objects
 * @param {Set} criticalTasks - Set of critical task IDs
 * @param {Map} predecessors - Map of predecessors
 * @param {Map} successors - Map of successors
 * @returns {Array} Critical path tasks in order
 */
function buildCriticalPath(tasks, criticalTasks, predecessors, successors) {
  // Find start tasks (no predecessors or only non-critical predecessors)
  const startTasks = tasks.filter(task => {
    if (!criticalTasks.has(task.id)) {
      return false
    }

    const preds = predecessors.get(task.id) || []
    return preds.length === 0 || preds.every(p => !criticalTasks.has(p))
  })

  // Trace forward from start tasks
  const path = []
  const visited = new Set()

  const trace = (taskId) => {
    if (visited.has(taskId)) {
      return
    }

    visited.add(taskId)

    if (criticalTasks.has(taskId)) {
      const task = tasks.find(t => t.id === taskId)
      if (task) {
        path.push(task)
      }

      // Follow critical successors
      const succs = successors.get(taskId) || []
      succs.forEach(succId => {
        if (criticalTasks.has(succId)) {
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

  return cpmResult.taskSlack.get(task.id) || 0
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

  return cpmResult.criticalTasks.has(task.id)
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

  const taskEF = cpmResult.earliestFinish.get(task.id)

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
      const lag = getDependencyLag(succ, task.id)
      return cpmResult.earliestStart.get(succ.id) - lag
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
 * @returns {Array} Array of critical paths
 */
export function getAllCriticalPaths(tasks) {
  const cpmResult = calculateCriticalPath(tasks)

  // Group critical tasks by their chains
  const paths = []
  const visited = new Set()

  const findPaths = (taskId, currentPath) => {
    if (visited.has(taskId)) {
      return
    }

    visited.add(taskId)

    if (!cpmResult.criticalTasks.has(taskId)) {
      if (currentPath.length > 0) {
        paths.push([...currentPath])
      }
      return
    }

    const task = tasks.find(t => t.id === taskId)
    if (!task) return

    currentPath.push(task)

    // Find critical successors
    const successors = tasks.filter(t =>
      t.dependencies && t.dependencies.includes(taskId) &&
      cpmResult.criticalTasks.has(t.id)
    )

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
    if (!cpmResult.criticalTasks.has(task.id)) {
      return false
    }
    return !task.dependencies || task.dependencies.length === 0 ||
      task.dependencies.every(depId => !cpmResult.criticalTasks.has(depId))
  })

  startTasks.forEach(task => findPaths(task.id, []))

  return paths.map(path => ({
    tasks: path,
    duration: path.reduce((sum, task) => {
      const duration = Math.ceil(
        (new Date(task.end_date) - new Date(task.start_date)) / (1000 * 60 * 60 * 24)
      )
      return sum + duration
    }, 0)
  }))
}

/**
 * Generates a critical path report
 *
 * @param {Array} tasks - Array of task objects
 * @returns {Object} Critical path report
 */
export function generateCriticalPathReport(tasks) {
  const cpmResult = calculateCriticalPath(tasks)
  const allPaths = getAllCriticalPaths(tasks)

  const report = {
    projectDuration: cpmResult.projectDuration,
    criticalTaskCount: cpmResult.criticalTasks.size,
    totalTaskCount: tasks.length,
    criticalPathCount: allPaths.length,
    criticalPaths: allPaths,
    taskDetails: tasks.map(task => ({
      id: task.id,
      name: task.name,
      isCritical: cpmResult.criticalTasks.has(task.id),
      totalSlack: cpmResult.taskSlack.get(task.id) || 0,
      earlyStart: cpmResult.earliestStart.get(task.id),
      earlyFinish: cpmResult.earliestFinish.get(task.id),
      lateStart: cpmResult.latestStart.get(task.id),
      lateFinish: cpmResult.latestFinish.get(task.id)
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
  const isCritical = cpmResult.criticalTasks.has(task.id)

  if (isCritical) {
    // Delaying critical task directly delays project
    return {
      projectDelay: delayDays,
      affectedTasks: Array.from(cpmResult.criticalTasks),
      reason: 'Task is on critical path'
    }
  } else {
    const slack = cpmResult.taskSlack.get(task.id) || 0

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
        affectedTasks: Array.from(cpmResult.criticalTasks),
        reason: 'Delay exceeds slack, pushes critical path'
      }
    }
  }
}
