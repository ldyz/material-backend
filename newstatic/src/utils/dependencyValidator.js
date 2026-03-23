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

// Alias for backward compatibility with tests
export const DependencyType = DependencyTypes

/**
 * Validates a dependency between two tasks
 *
 * @param {Object} dependency - Dependency object with from, to, type, lag
 * @param {Array} tasks - Array of all tasks
 * @param {Array} existingDependencies - Array of existing dependencies
 * @returns {Object} Validation result with valid flag and errors
 */
export function validateDependency(dependency, tasks, existingDependencies = []) {
  const errors = []
  const warnings = []

  // Check if dependency exists
  if (!dependency) {
    errors.push('Dependency is required')
    return { valid: false, errors, warnings }
  }

  // Build task map
  const taskMap = new Map(tasks.map(t => [t.id, t]))

  // Get predecessor and successor
  const predecessor = taskMap.get(dependency.from)
  const successor = taskMap.get(dependency.to)

  // Check if tasks exist
  if (!predecessor) {
    errors.push(`Predecessor task "${dependency.from}" not found`)
  }

  if (!successor) {
    errors.push(`Successor task "${dependency.to}" not found`)
  }

  // Check for self-reference
  if (dependency.from === dependency.to) {
    errors.push('Task cannot depend on the same task')
  }

  // Check for duplicate dependency
  const duplicate = existingDependencies.some(d =>
    d.from === dependency.from && d.to === dependency.to
  )
  if (duplicate) {
    errors.push('Dependency between these tasks already exists')
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

  // Validate dates based on dependency type (if both tasks exist)
  // Only check for clear violations that would make the dependency impossible
  if (predecessor && successor) {
    const predStart = new Date(predecessor.start_date || predecessor.start)
    const predEnd = new Date(predecessor.end_date || predecessor.end)
    const succStart = new Date(successor.start_date || successor.start)
    const succEnd = new Date(successor.end_date || successor.end)

    const lag = dependency.lag || 0
    const lagMs = lag * 24 * 60 * 60 * 1000

    switch (type) {
      case DependencyTypes.FINISH_TO_START:
        // Only flag as error if successor starts WAY before predecessor finishes
        // Allow some flexibility for editing
        const minAllowedStartFS = new Date(predEnd.getTime() + lagMs - (365 * 24 * 60 * 60 * 1000)) // Allow up to 1 year early
        if (succStart < minAllowedStartFS) {
          errors.push(
            `FS dependency violated: Successor starts too far before predecessor finishes. ` +
            `Predecessor ends: ${predecessor.end_date || predecessor.end}, ` +
            `Successor starts: ${successor.start_date || successor.start}`
          )
        } else if (succStart < new Date(predEnd.getTime() + lagMs)) {
          // Just a warning for minor violations
          warnings.push(
            `FS dependency: Successor starts before predecessor finishes + lag. ` +
            `Predecessor ends: ${predecessor.end_date || predecessor.end}, ` +
            `Successor starts: ${successor.start_date || successor.start}, ` +
            `Required: ${lag} days after predecessor end`
          )
        }
        break

      case DependencyTypes.FINISH_TO_FINISH:
        const minAllowedFinishFF = new Date(predEnd.getTime() + lagMs - (365 * 24 * 60 * 60 * 1000))
        if (succEnd < minAllowedFinishFF) {
          errors.push(
            `FF dependency violated: Successor finishes too far before predecessor finishes. ` +
            `Predecessor ends: ${predecessor.end_date || predecessor.end}, ` +
            `Successor ends: ${successor.end_date || successor.end}`
          )
        } else if (succEnd < new Date(predEnd.getTime() + lagMs)) {
          warnings.push(
            `FF dependency: Successor finishes before predecessor finishes + lag. ` +
            `Predecessor ends: ${predecessor.end_date || predecessor.end}, ` +
            `Successor ends: ${successor.end_date || successor.end}`
          )
        }
        break

      case DependencyTypes.START_TO_START:
        const minAllowedStartSS = new Date(predStart.getTime() + lagMs - (365 * 24 * 60 * 60 * 1000))
        if (succStart < minAllowedStartSS) {
          errors.push(
            `SS dependency violated: Successor starts too far before predecessor starts. ` +
            `Predecessor starts: ${predecessor.start_date || predecessor.start}, ` +
            `Successor starts: ${successor.start_date || successor.start}`
          )
        } else if (succStart < new Date(predStart.getTime() + lagMs)) {
          warnings.push(
            `SS dependency: Successor starts before predecessor starts + lag. ` +
            `Predecessor starts: ${predecessor.start_date || predecessor.start}, ` +
            `Successor starts: ${successor.start_date || successor.start}`
          )
        }
        break

      case DependencyTypes.START_TO_FINISH:
        const minAllowedFinishSF = new Date(predStart.getTime() + lagMs - (365 * 24 * 60 * 60 * 1000))
        if (succEnd < minAllowedFinishSF) {
          errors.push(
            `SF dependency violated: Successor finishes too far before predecessor starts. ` +
            `Predecessor starts: ${predecessor.start_date || predecessor.start}, ` +
            `Successor ends: ${successor.end_date || successor.end}`
          )
        } else if (succEnd < new Date(predStart.getTime() + lagMs)) {
          warnings.push(
            `SF dependency: Successor finishes before predecessor starts + lag. ` +
            `Predecessor starts: ${predecessor.start_date || predecessor.start}, ` +
            `Successor ends: ${successor.end_date || successor.end}`
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
 * @param {Array} dependencies - Array of dependency objects with from/to
 * @returns {Object} Result with hasCycles boolean and cycles array
 */
export function detectCircularDependencies(tasks, dependencies = []) {
  if (!tasks || !Array.isArray(tasks) || tasks.length === 0) {
    return { hasCycles: false, cycles: [] }
  }

  // Build adjacency list from dependencies
  const graph = new Map()
  const taskMap = new Map()

  tasks.forEach(task => {
    graph.set(task.id, [])
    taskMap.set(task.id, task)
  })

  // Build edges from dependencies array
  dependencies.forEach(dep => {
    if (dep.from && dep.to && graph.has(dep.from)) {
      graph.get(dep.from).push(dep.to)
    }
  })

  // Also check task.dependencies property for backward compatibility
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
        const cyclePathIds = path.slice(cycleStartIndex)
        const cyclePath = cyclePathIds.map(id => taskMap.get(id)).filter(t => t)

        circularPaths.push({
          tasks: cyclePath,
          path: cyclePathIds,
          description: cyclePath.length > 0
            ? `Circular dependency: ${cyclePath.map(t => t.name).join(' → ')} → ${cyclePath[0].name}`
            : 'Circular dependency detected'
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

  return {
    hasCycles: circularPaths.length > 0,
    cycles: circularPaths.map(cp => cp.path) // Return just the path arrays
  }
}

/**
 * Check if there are any circular dependencies
 *
 * @param {Array} tasks - Array of task objects
 * @returns {boolean} True if circular dependencies exist
 */
export function hasCircularDependencies(tasks) {
  const cycles = detectCircularDependencies(tasks)
  return cycles && cycles.length > 0
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

/**
 * Validates lag/lead time for a dependency
 *
 * @param {number} lag - Lag time in days
 * @param {Object} predecessor - Predecessor task
 * @param {Object} successor - Successor task
 * @param {string} type - Dependency type
 * @returns {Object} Validation result with valid flag and suggestions
 */
export function validateLag(lag, predecessor, successor, type = DependencyTypes.FINISH_TO_START) {
  const errors = []
  const warnings = []
  const suggestions = []

  if (!predecessor || !successor) {
    return { valid: false, errors: ['Both tasks are required'], warnings, suggestions }
  }

  // Check if lag is excessive (lead time more than half predecessor duration)
  const predDuration = predecessor.duration || 5
  if (lag < 0 && Math.abs(lag) > predDuration / 2) {
    errors.push('Lead time exceeds 50% of predecessor duration')
    suggestions.push('Consider reducing lead time or using SS dependency instead')
  }

  // Check temporal consistency
  if (type === DependencyTypes.FINISH_TO_START && lag >= 0) {
    // With positive lag, successor starts after predecessor finishes + lag
    const predEnd = new Date(predecessor.end_date || predecessor.end)
    const expectedStart = new Date(predEnd)
    expectedStart.setDate(expectedStart.getDate() + lag)

    if (successor.end_date && new Date(successor.end_date || successor.end) < expectedStart) {
      errors.push('Lag causes successor to finish before it starts')
    }
  }

  // Calculate limits
  const { minLag, maxLag } = calculateLagLimits(predecessor, successor, type)

  if (lag < minLag) {
    errors.push(`Lag ${lag} is below minimum ${minLag}`)
    suggestions.push(`Increase lag to at least ${minLag} days`)
  }

  if (lag > maxLag) {
    errors.push(`Lag ${lag} exceeds maximum ${maxLag}`)
    suggestions.push(`Reduce lag to at most ${maxLag} days`)
  }

  return {
    valid: errors.length === 0,
    errors,
    warnings,
    suggestions,
    limits: { min: minLag, max: maxLag }
  }
}

// Import and re-export critical path function
export { calculateCriticalPath } from './criticalPath.js'

/**
 * Analyze a single dependency path (alias for analyzeDependencyPaths with singular form)
 *
 * @param {string} taskId - Starting task ID
 * @param {Array} tasks - All tasks
 * @param {Array} dependencies - Array of dependency objects
 * @param {string} direction - 'successors' or 'predecessors'
 * @returns {Object} Path analysis result
 */
export function analyzeDependencyPath(taskId, tasks, dependencies = [], direction = 'successors') {
  // Build dependency maps
  const taskMap = new Map(tasks.map(t => [t.id, t]))
  const predecessors = new Map() // taskId -> array of predecessor task IDs
  const successors = new Map() // taskId -> array of successor task IDs

  // Initialize maps
  tasks.forEach(task => {
    predecessors.set(task.id, [])
    successors.set(task.id, [])
  })

  // Build from dependencies array
  dependencies.forEach(dep => {
    if (dep.from && dep.to) {
      successors.get(dep.from)?.push(dep.to)
      predecessors.get(dep.to)?.push(dep.from)
    }
  })

  // Also check task.dependencies property
  tasks.forEach(task => {
    if (task.dependencies && Array.isArray(task.dependencies)) {
      task.dependencies.forEach(depId => {
        if (taskMap.has(depId)) {
          successors.get(depId)?.push(task.id)
          predecessors.get(task.id)?.push(depId)
        }
      })
    }
  })

  // Get successors and predecessors
  const taskSuccessors = successors.get(taskId) || []
  const taskPredecessors = predecessors.get(taskId) || []

  // Build path
  const path = [taskId]
  const visited = new Set([taskId])

  const trace = (currentId) => {
    const succs = successors.get(currentId) || []
    for (const succId of succs) {
      if (!visited.has(succId)) {
        visited.add(succId)
        path.push(succId)
        trace(succId)
      }
    }
  }

  trace(taskId)

  // Calculate duration
  const pathTasks = path.map(id => taskMap.get(id)).filter(t => t)
  const duration = pathTasks.reduce((sum, task) => {
    const start = new Date(task.start_date || task.start)
    const end = new Date(task.end_date || task.end)
    return sum + Math.ceil((end - start) / (1000 * 60 * 60 * 24))
  }, 0)

  // Calculate flexibility (simplified - use path length as indicator)
  const flexibility = path.length > 1 ? path.length - 1 : 0

  // Identify merge points (tasks with multiple predecessors)
  const mergePoints = path.filter(id => {
    const preds = predecessors.get(id) || []
    return preds.length > 1
  })

  // Identify burst points (tasks with multiple successors)
  const burstPoints = path.filter(id => {
    const succs = successors.get(id) || []
    return succs.length > 1
  })

  return {
    path,
    length: path.length,
    duration,
    tasks: pathTasks,
    successors: taskSuccessors,
    predecessors: taskPredecessors,
    flexibility,
    mergePoints,
    burstPoints
  }
}
