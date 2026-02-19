/**
 * Dependency Validation Utilities
 *
 * Provides validation for task dependencies including:
 * - Circular dependency detection using DFS algorithm
 * - Dependency type validation
 * - Lag/Lead time validation
 * - Dependency path analysis
 */

/**
 * Dependency types
 */
export const DependencyTypes = {
  FINISH_TO_START: 'FS', // Default: Task A must finish before Task B can start
  FINISH_TO_FINISH: 'FF', // Task A must finish before Task B can finish
  START_TO_START: 'SS', // Task A must start before Task B can start
  START_TO_FINISH: 'SF' // Task A must start before Task B can finish
}

/**
 * Validates a dependency between two tasks
 *
 * @param {Object} predecessor - Predecessor task object
 * @param {Object} successor - Successor task object
 * @param {Object} dependency - Dependency object with type and lag
 * @returns {Object} Validation result with valid flag and errors
 */
export function validateDependency(predecessor, successor, dependency = {}) {
  const errors = []
  const warnings = []

  if (!predecessor) {
    errors.push('Predecessor task is required')
  }

  if (!successor) {
    errors.push('Successor task is required')
  }

  if (predecessor && successor && predecessor.id === successor.id) {
    errors.push('Task cannot depend on itself')
  }

  // Validate dependency type
  const type = dependency.type || DependencyTypes.FINISH_TO_START
  const validTypes = Object.values(DependencyTypes)

  if (!validTypes.includes(type)) {
    errors.push(`Invalid dependency type: ${type}. Must be one of: ${validTypes.join(', ')}`)
  }

  // Validate lag/lead time
  if (dependency.lag !== undefined) {
    if (typeof dependency.lag !== 'number') {
      errors.push('Lag must be a number')
    } else if (dependency.lag < -10000 || dependency.lag > 10000) {
      warnings.push('Lag value is unusually large')
    }
  }

  // Validate dates based on dependency type
  if (predecessor && successor) {
    const predStart = new Date(predecessor.start_date)
    const predEnd = new Date(predecessor.end_date)
    const succStart = new Date(successor.start_date)
    const succEnd = new Date(successor.end_date)

    const lag = dependency.lag || 0
    const lagMs = lag * 24 * 60 * 60 * 1000

    switch (type) {
      case DependencyTypes.FINISH_TO_START:
        // Successor must start after predecessor finishes + lag
        if (succStart < new Date(predEnd.getTime() + lagMs)) {
          errors.push(
            `FS dependency violated: Successor must start after predecessor finishes. ` +
            `Predecessor ends: ${predecessor.end_date}, ` +
            `Successor starts: ${successor.start_date}`
          )
        }
        break

      case DependencyTypes.FINISH_TO_FINISH:
        // Successor must finish after predecessor finishes + lag
        if (succEnd < new Date(predEnd.getTime() + lagMs)) {
          errors.push(
            `FF dependency violated: Successor must finish after predecessor finishes. ` +
            `Predecessor ends: ${predecessor.end_date}, ` +
            `Successor ends: ${successor.end_date}`
          )
        }
        break

      case DependencyTypes.START_TO_START:
        // Successor must start after predecessor starts + lag
        if (succStart < new Date(predStart.getTime() + lagMs)) {
          errors.push(
            `SS dependency violated: Successor must start after predecessor starts. ` +
            `Predecessor starts: ${predecessor.start_date}, ` +
            `Successor starts: ${successor.start_date}`
          )
        }
        break

      case DependencyTypes.START_TO_FINISH:
        // Successor must finish after predecessor starts + lag
        if (succEnd < new Date(predStart.getTime() + lagMs)) {
          errors.push(
            `SF dependency violated: Successor must finish after predecessor starts. ` +
            `Predecessor starts: ${predecessor.start_date}, ` +
            `Successor ends: ${successor.end_date}`
          )
        }
        break
    }

    // Check for negative lag (lead time)
    if (lag < 0) {
      warnings.push(
        `Negative lag (lead time) detected: ${lag} days. ` +
        `This may cause scheduling conflicts.`
      )
    }
  }

  return {
    valid: errors.length === 0,
    errors,
    warnings
  }
}

/**
 * Detects circular dependencies using DFS algorithm
 *
 * @param {Array} tasks - Array of task objects
 * @returns {Array} Array of circular dependency paths found
 */
export function detectCircularDependencies(tasks) {
  if (!tasks || !Array.isArray(tasks) || tasks.length === 0) {
    return []
  }

  // Build adjacency list
  const graph = new Map()
  const taskMap = new Map()

  tasks.forEach(task => {
    graph.set(task.id, [])
    taskMap.set(task.id, task)
  })

  // Build edges from dependencies
  tasks.forEach(task => {
    if (task.dependencies && Array.isArray(task.dependencies)) {
      task.dependencies.forEach(depId => {
        if (graph.has(depId)) {
          graph.get(depId).push(task.id)
        }
      })
    }
  })

  const circularPaths = []
  const visited = new Set()
  const recStack = new Set()
  const path = []

  // DFS helper function
  const dfs = (nodeId) => {
    visited.add(nodeId)
    recStack.add(nodeId)
    path.push(nodeId)

    const neighbors = graph.get(nodeId) || []

    for (const neighborId of neighbors) {
      if (!visited.has(neighborId)) {
        if (dfs(neighborId)) {
          return true
        }
      } else if (recStack.has(neighborId)) {
        // Found a cycle - extract the cycle path
        const cycleStartIndex = path.indexOf(neighborId)
        const cyclePath = path.slice(cycleStartIndex).map(id => taskMap.get(id))

        circularPaths.push({
          tasks: cyclePath,
          path: path.slice(cycleStartIndex),
          description: `Circular dependency: ${cyclePath.map(t => t.name).join(' → ')} → ${cyclePath[0].name}`
        })

        return true
      }
    }

    path.pop()
    recStack.delete(nodeId)
    return false
  }

  // Run DFS on all nodes
  for (const taskId of graph.keys()) {
    if (!visited.has(taskId)) {
      dfs(taskId)
    }
  }

  return circularPaths
}

/**
 * Validates all dependencies in a task list
 *
 * @param {Array} tasks - Array of task objects
 * @returns {Object} Validation result with valid flag, errors, and warnings
 */
export function validateAllDependencies(tasks) {
  if (!tasks || !Array.isArray(tasks)) {
    return {
      valid: false,
      errors: ['Tasks must be an array'],
      warnings: []
    }
  }

  const errors = []
  const warnings = []
  const taskMap = new Map(tasks.map(t => [t.id, t]))

  // Check for circular dependencies
  const circularPaths = detectCircularDependencies(tasks)
  if (circularPaths.length > 0) {
    circularPaths.forEach(path => {
      errors.push(path.description)
    })
  }

  // Validate each dependency
  tasks.forEach(task => {
    if (task.dependencies && Array.isArray(task.dependencies)) {
      task.dependencies.forEach(depId => {
        const predecessor = taskMap.get(depId)
        const successor = task

        if (!predecessor) {
          errors.push(
            `Task "${task.name}" has dependency on non-existent task ID: ${depId}`
          )
          return
        }

        // Get dependency details if available
        const dependencyDetails = task.dependencyDetails?.find(d => d.taskId === depId) || {}
        const validation = validateDependency(predecessor, successor, dependencyDetails)

        errors.push(...validation.errors.map(e =>
          `Dependency error on task "${task.name}": ${e}`
        ))
        warnings.push(...validation.warnings.map(w =>
          `Dependency warning on task "${task.name}": ${w}`
        ))
      })
    }
  })

  return {
    valid: errors.length === 0,
    errors,
    warnings,
    circularPaths
  }
}

/**
 * Analyzes dependency paths from a task
 *
 * @param {number} taskId - Starting task ID
 * @param {Array} tasks - Array of task objects
 * @param {string} direction - 'predecessors' or 'successors'
 * @returns {Array} Array of dependency paths
 */
export function analyzeDependencyPaths(taskId, tasks, direction = 'successors') {
  if (!tasks || !Array.isArray(tasks)) {
    return []
  }

  const taskMap = new Map(tasks.map(t => [t.id, t]))
  const paths = []
  const visited = new Set()

  const traverse = (currentId, currentPath) => {
    if (visited.has(currentId)) {
      return
    }

    visited.add(currentId)
    const currentTask = taskMap.get(currentId)

    if (!currentTask) {
      return
    }

    const newPath = [...currentPath, currentTask]

    if (direction === 'successors') {
      // Find tasks that depend on current task
      const successors = tasks.filter(t =>
        t.dependencies && t.dependencies.includes(currentId)
      )

      if (successors.length === 0) {
        // End of path
        paths.push(newPath)
      } else {
        successors.forEach(successor => {
          traverse(successor.id, newPath)
        })
      }
    } else {
      // Find predecessors
      const predecessors = currentTask.dependencies || []

      if (predecessors.length === 0) {
        // End of path
        paths.push(newPath)
      } else {
        predecessors.forEach(predId => {
          traverse(predId, newPath)
        })
      }
    }

    visited.delete(currentId)
  }

  traverse(taskId, [])

  return paths.map(path => ({
    tasks: path,
    length: path.length,
    duration: calculatePathDuration(path)
  }))
}

/**
 * Calculates the total duration of a dependency path
 *
 * @param {Array} path - Array of tasks in the path
 * @returns {number} Duration in days
 */
export function calculatePathDuration(path) {
  if (!path || path.length === 0) {
    return 0
  }

  const startDates = path.map(t => new Date(t.start_date).getTime())
  const endDates = path.map(t => new Date(t.end_date).getTime())

  const minStart = Math.min(...startDates)
  const maxEnd = Math.max(...endDates)

  return Math.ceil((maxEnd - minStart) / (1000 * 60 * 60 * 24))
}

/**
 * Gets all tasks that would be affected by a task date change
 *
 * @param {number} taskId - Task ID to check
 * @param {Array} tasks - Array of task objects
 * @returns {Set} Set of affected task IDs
 */
export function getAffectedTasks(taskId, tasks) {
  const affected = new Set()

  const propagate = (currentId) => {
    const successors = tasks.filter(t =>
      t.dependencies && t.dependencies.includes(currentId)
    )

    successors.forEach(successor => {
      if (!affected.has(successor.id)) {
        affected.add(successor.id)
        propagate(successor.id)
      }
    })
  }

  affected.add(taskId)
  propagate(taskId)

  return affected
}

/**
 * Calculates the maximum lag/lead time that can be applied
 *
 * @param {Object} predecessor - Predecessor task
 * @param {Object} successor - Successor task
 * @param {string} type - Dependency type
 * @returns {Object} Min and max lag values
 */
export function calculateLagLimits(predecessor, successor, type) {
  if (!predecessor || !successor) {
    return { minLag: 0, maxLag: 0 }
  }

  const predStart = new Date(predecessor.start_date)
  const predEnd = new Date(predecessor.end_date)
  const succStart = new Date(successor.start_date)
  const succEnd = new Date(successor.end_date)

  let minLag = -Infinity
  let maxLag = Infinity

  switch (type) {
    case DependencyTypes.FINISH_TO_START:
      // Successor starts after predecessor finishes
      const currentFSLag = Math.ceil((succStart - predEnd) / (1000 * 60 * 60 * 24))
      minLag = currentFSLag - 365 // Arbitrary limit
      maxLag = currentFSLag + 365
      break

    case DependencyTypes.FINISH_TO_FINISH:
      // Successor finishes after predecessor finishes
      const currentFFLag = Math.ceil((succEnd - predEnd) / (1000 * 60 * 60 * 24))
      minLag = currentFFLag - 365
      maxLag = currentFFLag + 365
      break

    case DependencyTypes.START_TO_START:
      // Successor starts after predecessor starts
      const currentSSLag = Math.ceil((succStart - predStart) / (1000 * 60 * 60 * 24))
      minLag = currentSSLag - 365
      maxLag = currentSSLag + 365
      break

    case DependencyTypes.START_TO_FINISH:
      // Successor finishes after predecessor starts
      const currentSFLag = Math.ceil((succEnd - predStart) / (1000 * 60 * 60 * 24))
      minLag = currentSFLag - 365
      maxLag = currentSFLag + 365
      break
  }

  return { minLag, maxLag }
}
