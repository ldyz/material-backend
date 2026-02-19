/**
 * Resource Leveling Utilities
 *
 * Provides resource conflict detection and leveling algorithms
 * for optimizing task schedules based on resource availability.
 */

/**
 * Detects resource conflicts across all tasks
 *
 * @param {Array} tasks - Array of task objects
 * @param {Array} resources - Array of resource objects
 * @returns {Promise<Array>} Array of conflict objects
 */
export async function detectResourceConflicts(tasks, resources) {
  if (!tasks || !Array.isArray(tasks) || tasks.length === 0) {
    return []
  }

  if (!resources || !Array.isArray(resources) || resources.length === 0) {
    return []
  }

  const conflicts = []
  const resourceMap = new Map()

  // Build resource map
  resources.forEach(resource => {
    resourceMap.set(resource.id, {
      ...resource,
      allocations: []
    })
  })

  // Track allocations by resource and date
  const allocationsByResourceAndDate = new Map()

  // Process each task's resource assignments
  tasks.forEach(task => {
    if (!task.resources || !Array.isArray(task.resources)) {
      return
    }

    const startDate = new Date(task.start_date)
    const endDate = new Date(task.end_date)
    const days = Math.ceil((endDate - startDate) / (1000 * 60 * 60 * 24))

    task.resources.forEach(assignment => {
      const resourceId = assignment.resource_id
      const allocation = assignment.units || 100 // Percentage allocation

      // Check each day of the task
      for (let i = 0; i <= days; i++) {
        const currentDate = new Date(startDate)
        currentDate.setDate(currentDate.getDate() + i)
        const dateKey = currentDate.toISOString().split('T')[0]

        const key = `${resourceId}_${dateKey}`
        if (!allocationsByResourceAndDate.has(key)) {
          allocationsByResourceAndDate.set(key, {
            resourceId,
            date: dateKey,
            tasks: [],
            totalAllocation: 0
          })
        }

        const allocationData = allocationsByResourceAndDate.get(key)
        allocationData.tasks.push({
          id: task.id,
          name: task.name,
          allocation
        })
        allocationData.totalAllocation += allocation
      }
    })
  })

  // Find conflicts where total allocation exceeds capacity
  allocationsByResourceAndDate.forEach(allocationData => {
    const resource = resourceMap.get(allocationData.resourceId)
    if (!resource) return

    const capacity = resource.capacity || 100 // Default 100%

    if (allocationData.totalAllocation > capacity) {
      conflicts.push({
        resourceId: allocationData.resourceId,
        resourceName: resource.name,
        date: allocationData.date,
        assigned: Math.round(allocationData.totalAllocation),
        capacity,
        overallocation: allocationData.totalAllocation - capacity,
        tasks: allocationData.tasks
      })
    }
  })

  return conflicts
}

/**
 * Applies heuristic resource leveling algorithm
 *
 * @param {Array} tasks - Array of task objects
 * @param {Array} resources - Array of resource objects
 * @param {Object} options - Leveling options
 * @returns {Promise<Array>} Leveled tasks array
 */
export async function applyResourceLeveling(tasks, resources, options = {}) {
  if (!tasks || !Array.isArray(tasks)) {
    return tasks
  }

  // Make deep copy to avoid mutating original
  const leveledTasks = JSON.parse(JSON.stringify(tasks))

  const {
    priority = 'priority', // 'priority', 'duration', 'slack'
    range = 'all', // 'all', 'selected'
    allowSplitting = false,
    adjustDependencies = false
  } = options

  // Sort tasks by priority (descending)
  const sortedTasks = [...leveledTasks].sort((a, b) => {
    switch (priority) {
      case 'priority':
        // Higher priority number = higher priority
        return (b.priority || 0) - (a.priority || 0)
      case 'duration':
        const aDuration = new Date(a.end_date) - new Date(a.start_date)
        const bDuration = new Date(b.end_date) - new Date(b.start_date)
        return bDuration - aDuration // Longer tasks first
      case 'slack':
        // More slack first (easier to move)
        return (b.slack || 0) - (a.slack || 0)
      default:
        return 0
    }
  })

  // Track resource usage by date
  const resourceUsage = new Map()

  // Initialize resource usage tracking
  resources.forEach(resource => {
    resourceUsage.set(resource.id, {
      capacity: resource.capacity || 100,
      dailyUsage: new Map()
    })
  })

  // Process tasks in priority order
  const processedTasks = new Set()

  for (const task of sortedTasks) {
    if (processedTasks.has(task.id)) {
      continue
    }

    // Find valid start date that doesn't cause conflicts
    let validStartDate = new Date(task.start_date)
    let validEndDate = new Date(task.end_date)
    let shifted = false

    // Check each day for conflicts
    let maxDelay = 0
    const duration = Math.ceil((validEndDate - validStartDate) / (1000 * 60 * 60 * 24))

    for (let day = 0; day <= duration; day++) {
      const checkDate = new Date(validStartDate)
      checkDate.setDate(checkDate.getDate() + day)
      const dateKey = checkDate.toISOString().split('T')[0]

      // Check all resource assignments for this task
      if (task.resources && Array.isArray(task.resources)) {
        for (const assignment of task.resources) {
          const resourceId = assignment.resource_id
          const usage = resourceUsage.get(resourceId)

          if (!usage) continue

          const dailyUsage = usage.dailyUsage.get(dateKey) || 0
          const allocation = assignment.units || 100

          if (dailyUsage + allocation > usage.capacity) {
            // Conflict found - need to delay task
            const daysUntilAvailable = calculateDaysUntilAvailable(
              resourceUsage,
              resourceId,
              checkDate,
              allocation
            )

            if (daysUntilAvailable > maxDelay) {
              maxDelay = daysUntilAvailable
              shifted = true
            }
          }
        }
      }
    }

    // Apply delay if needed
    if (shifted && maxDelay > 0) {
      validStartDate.setDate(validStartDate.getDate() + maxDelay)
      validEndDate.setDate(validEndDate.getDate() + maxDelay)

      const taskIndex = leveledTasks.findIndex(t => t.id === task.id)
      if (taskIndex !== -1) {
        leveledTasks[taskIndex].start_date = validStartDate.toISOString().split('T')[0]
        leveledTasks[taskIndex].end_date = validEndDate.toISOString().split('T')[0]
        leveledTasks[taskIndex].leveled = true
      }
    }

    // Update resource usage with this task's allocations
    if (task.resources && Array.isArray(task.resources)) {
      for (let assignment of task.resources) {
        const resourceId = assignment.resource_id
        const usage = resourceUsage.get(resourceId)

        if (!usage) continue

        const allocation = assignment.units || 100

        for (let day = 0; day <= duration; day++) {
          const allocDate = new Date(validStartDate)
          allocDate.setDate(allocDate.getDate() + day)
          const dateKey = allocDate.toISOString().split('T')[0]

          const currentUsage = usage.dailyUsage.get(dateKey) || 0
          usage.dailyUsage.set(dateKey, currentUsage + allocation)
        }
      }
    }

    processedTasks.add(task.id)
  }

  // Adjust dependencies if needed
  if (adjustDependencies) {
    await adjustTaskDependencies(leveledTasks)
  }

  return leveledTasks
}

/**
 * Calculates days until a resource has available capacity
 *
 * @param {Map} resourceUsage - Resource usage tracking map
 * @param {number} resourceId - Resource ID
 * @param {Date} startDate - Start date to check from
 * @param {number} requiredAllocation - Required allocation percentage
 * @returns {number} Days until available
 */
function calculateDaysUntilAvailable(resourceUsage, resourceId, startDate, requiredAllocation) {
  const usage = resourceUsage.get(resourceId)
  if (!usage) return 0

  let daysUntilAvailable = 0
  const checkDate = new Date(startDate)

  while (daysUntilAvailable < 365) { // Prevent infinite loop
    const dateKey = checkDate.toISOString().split('T')[0]
    const dailyUsage = usage.dailyUsage.get(dateKey) || 0

    if (dailyUsage + requiredAllocation <= usage.capacity) {
      break
    }

    checkDate.setDate(checkDate.getDate() + 1)
    daysUntilAvailable++
  }

  return daysUntilAvailable
}

/**
 * Adjusts task dependencies after leveling
 *
 * @param {Array} tasks - Leveled tasks array
 * @returns {Promise<void>}
 */
async function adjustTaskDependencies(tasks) {
  // Build dependency graph
  const taskMap = new Map(tasks.map(t => [t.id, t]))

  // Update successor tasks if their predecessor was delayed
  for (const task of tasks) {
    if (!task.dependencies || !Array.isArray(task.dependencies)) {
      continue
    }

    for (const depId of task.dependencies) {
      const predecessor = taskMap.get(depId)
      if (!predecessor) continue

      const predEndDate = new Date(predecessor.end_date)
      const taskStartDate = new Date(task.start_date)

      if (predEndDate > taskStartDate) {
        // Adjust successor to start after predecessor
        const taskIndex = tasks.findIndex(t => t.id === task.id)
        if (taskIndex !== -1) {
          const duration = Math.ceil(
            (new Date(task.end_date) - new Date(task.start_date)) / (1000 * 60 * 60 * 24)
          )

          const newStartDate = new Date(predEndDate)
          newStartDate.setDate(newStartDate.getDate() + 1) // Start next day

          const newEndDate = new Date(newStartDate)
          newEndDate.setDate(newEndDate.getDate() + duration)

          tasks[taskIndex].start_date = newStartDate.toISOString().split('T')[0]
          tasks[taskIndex].end_date = newEndDate.toISOString().split('T')[0]
          tasks[taskIndex].dependencyAdjusted = true
        }
      }
    }
  }
}

/**
 * Calculates resource overallocation statistics
 *
 * @param {Array} tasks - Array of task objects
 * @param {Array} resources - Array of resource objects
 * @returns {Promise<Object>} Statistics object
 */
export async function calculateResourceOverallocation(tasks, resources) {
  const conflicts = await detectResourceConflicts(tasks, resources)

  const totalOverallocation = conflicts.reduce(
    (sum, conflict) => sum + conflict.overallocation,
    0
  )

  const affectedResources = new Set(conflicts.map(c => c.resourceId))
  const affectedTasks = new Set()
  conflicts.forEach(conflict => {
    conflict.tasks.forEach(task => affectedTasks.add(task.id))
  })

  return {
    conflictCount: conflicts.length,
    totalOverallocation,
    affectedResourceCount: affectedResources.size,
    affectedTaskCount: affectedTasks.size,
    conflicts
  }
}

/**
 * Generates leveling suggestions for conflicts
 *
 * @param {Array} conflicts - Array of conflict objects
 * @returns {Array} Array of suggestion objects
 */
export function generateLevelingSuggestions(conflicts) {
  if (!conflicts || conflicts.length === 0) {
    return []
  }

  return conflicts.map(conflict => {
    const suggestions = []

    // Suggest delaying low-priority tasks
    const lowPriorityTasks = conflict.tasks.filter(t => !t.priority || t.priority < 3)
    if (lowPriorityTasks.length > 0) {
      suggestions.push({
        type: 'delay',
        description: `Delay low-priority tasks: ${lowPriorityTasks.map(t => t.name).join(', ')}`,
        tasks: lowPriorityTasks.map(t => t.id)
      })
    }

    // Suggest adding more resources
    suggestions.push({
      type: 'addResource',
      description: `Add additional capacity for ${conflict.resourceName}`,
      resourceId: conflict.resourceId
    })

    // Suggest reducing allocation
    suggestions.push({
      type: 'reduceAllocation',
      description: `Reduce allocation for tasks on ${conflict.resourceName}`,
      tasks: conflict.tasks.map(t => t.id)
    })

    return {
      conflictId: `${conflict.resourceId}_${conflict.date}`,
      conflict,
      suggestions
    }
  })
}

/**
 * Calculates leveling statistics for before/after comparison
 *
 * @param {Array} originalTasks - Tasks before leveling
 * @param {Array} leveledTasks - Tasks after leveling
 * @returns {Object} Statistics object
 */
export function calculateLevelingStatistics(originalTasks, leveledTasks) {
  const originalMap = new Map(originalTasks.map(t => [t.id, t]))
  const leveledMap = new Map(leveledTasks.map(t => [t.id, t]))

  let tasksDelayed = 0
  let maxDelay = 0

  leveledTasks.forEach(leveledTask => {
    const originalTask = originalMap.get(leveledTask.id)
    if (!originalTask) return

    const originalStart = new Date(originalTask.start_date)
    const leveledStart = new Date(leveledTask.start_date)
    const delay = Math.ceil((leveledStart - originalStart) / (1000 * 60 * 60 * 24))

    if (delay > 0) {
      tasksDelayed++
      if (delay > maxDelay) {
        maxDelay = delay
      }
    }
  })

  // Calculate project extension
  const originalEnd = originalTasks.reduce((max, task) => {
    const endDate = new Date(task.end_date)
    return endDate > max ? endDate : max
  }, new Date(0))

  const leveledEnd = leveledTasks.reduce((max, task) => {
    const endDate = new Date(task.end_date)
    return endDate > max ? endDate : max
  }, new Date(0))

  const projectExtension = Math.ceil((leveledEnd - originalEnd) / (1000 * 60 * 60 * 24))

  return {
    tasksDelayed,
    maxDelay: Math.max(0, maxDelay),
    projectExtension: Math.max(0, projectExtension)
  }
}

/**
 * Optimizes task scheduling based on resource availability
 *
 * @param {Array} tasks - Array of task objects
 * @param {Array} resources - Array of resource objects
 * @returns {Promise<Array>} Optimized tasks array
 */
export async function optimizeTaskSchedule(tasks, resources) {
  // Apply leveling with optimal settings
  return await applyResourceLeveling(tasks, resources, {
    priority: 'priority',
    range: 'all',
    allowSplitting: false,
    adjustDependencies: true
  })
}
