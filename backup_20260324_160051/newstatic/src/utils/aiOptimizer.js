/**
 * AI-Powered Schedule Optimizer
 * Provides intelligent suggestions for schedule optimization
 */

/**
 * Analyze schedule and provide suggestions
 * @param {Array} tasks - All tasks in the project
 * @param {Array} dependencies - Task dependencies
 * @param {Array} resources - Available resources
 * @returns {Object} Analysis results with risks, suggestions, and optimizations
 */
export function analyzeSchedule(tasks, dependencies, resources = []) {
  const analysis = {
    overallScore: 100,
    risks: [],
    suggestions: [],
    optimizations: [],
    criticalPath: [],
    overAllocations: [],
    analyzedAt: new Date().toISOString()
  }

  // 1. Critical Path Analysis
  const criticalPathResult = analyzeCriticalPath(tasks, dependencies)
  analysis.criticalPath = criticalPathResult.path
  analysis.criticalPathDuration = criticalPathResult.duration

  // 2. Risk Analysis
  const riskAnalysis = analyzeRisks(tasks, dependencies)
  analysis.risks = riskAnalysis.risks
  analysis.overallScore -= riskAnalysis.scorePenalty

  // 3. Resource Analysis
  const resourceAnalysis = analyzeResources(tasks, resources)
  analysis.overAllocations = resourceAnalysis.overAllocations
  analysis.suggestions.push(...resourceAnalysis.suggestions)
  analysis.overallScore -= resourceAnalysis.scorePenalty

  // 4. Dependency Analysis
  const dependencyAnalysis = analyzeDependencies(tasks, dependencies)
  analysis.suggestions.push(...dependencyAnalysis.suggestions)
  analysis.overallScore -= dependencyAnalysis.scorePenalty

  // 5. Schedule Optimization
  const optimizationResult = generateOptimizations(tasks, dependencies, resources)
  analysis.optimizations = optimizationResult.optimizations
  analysis.suggestions.push(...optimizationResult.suggestions)

  // Ensure score is within 0-100
  analysis.overallScore = Math.max(0, Math.min(100, analysis.overallScore))

  return analysis
}

/**
 * Analyze critical path
 */
function analyzeCriticalPath(tasks, dependencies) {
  if (!tasks || tasks.length === 0) {
    return { path: [], duration: 0 }
  }

  // Build dependency graph
  const graph = new Map()
  const inDegree = new Map()
  const dist = new Map()

  tasks.forEach(task => {
    graph.set(task.id, [])
    inDegree.set(task.id, 0)
    dist.set(task.id, task.start + task.duration)
  })

  dependencies.forEach(dep => {
    const pred = dep.predecessorId
    const succ = dep.successorId
    if (graph.has(pred)) {
      graph.get(pred).push(succ)
      inDegree.set(succ, (inDegree.get(succ) || 0) + 1)
    }
  })

  // Find tasks with no predecessors (start tasks)
  const startTasks = tasks.filter(t => inDegree.get(t.id) === 0)

  // Calculate longest path using topological sort
  const path = []
  let maxDuration = 0

  // Simple approach: find the path with maximum duration
  // In production, use proper longest path algorithm
  const visited = new Set()

  function dfs(taskId, currentPath, currentDuration) {
    if (visited.has(taskId)) return

    visited.add(taskId)
    const task = tasks.find(t => t.id === taskId)
    if (!task) return

    const newPath = [...currentPath, taskId]
    const newDuration = currentDuration + task.duration

    const successors = graph.get(taskId) || []
    if (successors.length === 0) {
      // End of path
      if (newDuration > maxDuration) {
        maxDuration = newDuration
        path.length = 0
        path.push(...newPath)
      }
    } else {
      successors.forEach(succId => {
        dfs(succId, newPath, newDuration)
      })
    }
  }

  startTasks.forEach(task => {
    visited.clear()
    dfs(task.id, [], 0)
  })

  return { path, duration: maxDuration }
}

/**
 * Analyze schedule risks
 */
function analyzeRisks(tasks, dependencies) {
  const risks = []
  let scorePenalty = 0

  tasks.forEach(task => {
    const taskRisks = []

    // Check for long duration (higher risk)
    if (task.duration > 20) {
      taskRisks.push({
        type: 'long_duration',
        level: 'medium',
        message: `Task "${task.name}" has long duration (${task.duration} days), consider breaking it down`
      })
    }

    // Check for many dependencies (complexity risk)
    const dependencyCount = dependencies.filter(
      d => d.predecessorId === task.id || d.successorId === task.id
    ).length

    if (dependencyCount > 5) {
      taskRisks.push({
        type: 'high_complexity',
        level: 'high',
        message: `Task "${task.name}" has ${dependencyCount} dependencies, high complexity risk`
      })
    }

    // Check for slack/float (tasks with no slack are risky)
    const hasSlack = calculateTaskSlack(task, tasks, dependencies)
    if (hasSlack === 0 && !task.isMilestone) {
      taskRisks.push({
        type: 'no_slack',
        level: 'high',
        message: `Task "${task.name}" has no slack, any delay will impact project end date`
      })
    }

    // Add risks to list
    risks.push(...taskRisks.map(r => ({ ...r, taskId: task.id })))

    // Calculate score penalty
    taskRisks.forEach(risk => {
      if (risk.level === 'high') scorePenalty += 5
      else if (risk.level === 'medium') scorePenalty += 2
    })
  })

  return { risks, scorePenalty }
}

/**
 * Calculate task slack (float)
 */
function calculateTaskSlack(task, tasks, dependencies) {
  // Early start = task.start
  // Late start = early start + slack
  // Simplified calculation
  const successors = dependencies
    .filter(d => d.predecessorId === task.id)
    .map(d => tasks.find(t => t.id === d.successorId))
    .filter(Boolean)

  if (successors.length === 0) {
    return 0 // No successors, on critical path
  }

  const minSuccessorStart = Math.min(...successors.map(s => s.start))
  const taskEnd = task.start + task.duration

  return Math.max(0, minSuccessorStart - taskEnd)
}

/**
 * Analyze resource allocation
 */
function analyzeResources(tasks, resources) {
  const overAllocations = []
  const suggestions = []
  let scorePenalty = 0

  if (!resources || resources.length === 0) {
    return { overAllocations, suggestions, scorePenalty }
  }

  // Calculate resource utilization per day
  const utilization = new Map()

  tasks.forEach(task => {
    if (!task.assignedTo) return

    for (let day = task.start; day < task.start + task.duration; day++) {
      const key = `${task.assignedTo}-${day}`
      utilization.set(key, (utilization.get(key) || 0) + 1)
    }
  })

  // Find over-allocations
  resources.forEach(resource => {
    const overloadedDays = []

    for (let day = 0; day < 365; day++) { // Assuming 1 year project
      const key = `${resource.id}-${day}`
      const load = utilization.get(key) || 0

      if (load > 1) {
        overloadedDays.push(day)
      }
    }

    if (overloadedDays.length > 0) {
      overAllocations.push({
        resourceId: resource.id,
        resourceName: resource.name,
        overloadedDays,
        severity: overloadedDays.length > 10 ? 'high' : 'medium'
      })

      scorePenalty += 10

      suggestions.push({
        id: `resource-${resource.id}`,
        type: 'resource',
        title: 'Resource Overallocation Detected',
        description: `${resource.name} is assigned to multiple tasks on ${overloadedDays.length} days`,
        impact: 'high',
        actions: [
          {
            type: 'reassign',
            description: 'Reassign some tasks to other resources'
          },
          {
            type: 'reschedule',
            description: 'Adjust task dates to reduce overlap'
          }
        ]
      })
    }
  })

  return { overAllocations, suggestions, scorePenalty }
}

/**
 * Analyze dependencies for optimization
 */
function analyzeDependencies(tasks, dependencies) {
  const suggestions = []
  let scorePenalty = 0

  // Check for redundant dependencies
  const redundantDeps = findRedundantDependencies(tasks, dependencies)
  if (redundantDeps.length > 0) {
    scorePenalty += redundantDeps.length * 2

    redundantDeps.forEach(dep => {
      const pred = tasks.find(t => t.id === dep.predecessorId)
      const succ = tasks.find(t => t.id === dep.successorId)

      suggestions.push({
        id: `dep-${dep.id}`,
        type: 'dependency',
        title: 'Redundant Dependency',
        description: `Dependency "${pred?.name} → ${succ?.name}" may be redundant`,
        impact: 'low',
        actions: [
          {
            type: 'remove',
            description: 'Remove this dependency',
            dependencyId: dep.id
          }
        ]
      })
    })
  }

  return { suggestions, scorePenalty }
}

/**
 * Find redundant dependencies (transitive dependencies)
 */
function findRedundantDependencies(tasks, dependencies) {
  const redundant = []

  // Build adjacency list
  const graph = new Map()
  tasks.forEach(task => {
    graph.set(task.id, new Set())
  })

  dependencies.forEach(dep => {
    if (graph.has(dep.predecessorId)) {
      graph.get(dep.predecessorId).add(dep.successorId)
    }
  })

  // Check for transitive relationships
  dependencies.forEach(dep => {
    const visited = new Set()
    const queue = [dep.predecessorId]

    while (queue.length > 0) {
      const current = queue.shift()
      if (current === dep.successorId && current !== dep.predecessorId) {
        // Found indirect path
        redundant.push(dep)
        break
      }

      if (visited.has(current)) continue
      visited.add(current)

      const neighbors = graph.get(current) || new Set()
      neighbors.forEach(neighbor => {
        if (neighbor !== dep.successorId) {
          queue.push(neighbor)
        }
      })
    }
  })

  return redundant
}

/**
 * Generate schedule optimizations
 */
function generateOptimizations(tasks, dependencies, resources) {
  const optimizations = []
  const suggestions = []

  // 1. Parallel opportunities
  const parallelOpportunities = findParallelOpportunities(tasks, dependencies)
  if (parallelOpportunities.length > 0) {
    optimizations.push({
      type: 'parallelization',
      title: 'Parallel Execution Opportunities',
      description: `${parallelOpportunities.length} tasks can be executed in parallel`,
      potentialSavings: parallelOpportunities.reduce((sum, t) => sum + t.savings, 0),
      tasks: parallelOpportunities
    })

    suggestions.push({
      id: 'parallel-1',
      type: 'optimization',
      title: 'Parallel Task Execution',
      description: `Execute ${parallelOpportunities.length} tasks in parallel to reduce project duration`,
      impact: 'high',
      actions: parallelOpportunities.map(t => ({
        type: 'reschedule',
        taskId: t.taskId,
        newStart: t.newStart,
        description: `Move "${t.taskName}" to start at day ${t.newStart}`
      }))
    })
  }

  // 2. Task splitting opportunities
  const splittableTasks = findSplittableTasks(tasks)
  if (splittableTasks.length > 0) {
    optimizations.push({
      type: 'task_splitting',
      title: 'Task Splitting Opportunities',
      description: `${splittableTasks.length} long tasks can be split for better resource utilization`,
      tasks: splittableTasks
    })
  }

  return { optimizations, suggestions }
}

/**
 * Find tasks that can be executed in parallel
 */
function findParallelOpportunities(tasks, dependencies) {
  const opportunities = []

  // Group tasks by dependency depth
  const levels = new Map()
  const visited = new Set()

  function getLevel(taskId) {
    if (levels.has(taskId)) return levels.get(taskId)

    const predecessors = dependencies
      .filter(d => d.successorId === taskId)
      .map(d => d.predecessorId)

    if (predecessors.length === 0) {
      levels.set(taskId, 0)
      return 0
    }

    const maxPredLevel = Math.max(...predecessors.map(p => getLevel(p)))
    levels.set(taskId, maxPredLevel + 1)
    return levels.get(taskId)
  }

  tasks.forEach(task => getLevel(task.id))

  // Find tasks at same level that can be parallelized
  const levelGroups = new Map()
  levels.forEach((level, taskId) => {
    if (!levelGroups.has(level)) {
      levelGroups.set(level, [])
    }
    levelGroups.get(level).push(taskId)
  })

  levelGroups.forEach((taskIds, level) => {
    if (taskIds.length > 1) {
      // Calculate potential savings
      const tasksInLevel = taskIds.map(id => tasks.find(t => t.id === id)).filter(Boolean)
      const avgDuration = tasksInLevel.reduce((sum, t) => sum + t.duration, 0) / tasksInLevel.length
      const maxDuration = Math.max(...tasksInLevel.map(t => t.duration))
      const savings = maxDuration - avgDuration

      if (savings > 2) { // Only if significant savings
        opportunities.push({
          level,
          taskIds,
          savings: Math.round(savings),
          tasks: tasksInLevel.map(t => ({
            id: t.id,
            name: t.name,
            duration: t.duration
          }))
        })
      }
    }
  })

  return opportunities
}

/**
 * Find tasks that can be split
 */
function findSplittableTasks(tasks) {
  const SPLIT_THRESHOLD = 10 // days

  return tasks
    .filter(t => t.duration > SPLIT_THRESHOLD && !t.isMilestone)
    .map(t => ({
      taskId: t.id,
      taskName: t.name,
      duration: t.duration,
      suggestedSplits: Math.ceil(t.duration / SPLIT_THRESHOLD)
    }))
}

/**
 * Predict delays using historical data (simplified)
 * @param {Array} tasks - Current tasks
 * @param {Object} historicalData - Historical completion data
 * @returns {Array} Tasks at risk of delay
 */
export function predictDelays(tasks, historicalData = {}) {
  const atRiskTasks = []

  tasks.forEach(task => {
    let riskScore = 0
    const factors = []

    // Factor 1: Task complexity (based on dependencies)
    if (task.dependencies?.length > 3) {
      riskScore += 20
      factors.push('High complexity (many dependencies)')
    }

    // Factor 2: Duration (longer tasks have higher risk)
    if (task.duration > 14) {
      riskScore += 15
      factors.push('Long duration')
    }

    // Factor 3: Historical data (if available)
    if (historicalData[task.type]) {
      const avgDuration = historicalData[task.type].avgDuration
      const variance = historicalData[task.type].variance

      if (task.duration < avgDuration + variance) {
        riskScore += 30
        factors.push('Historically tends to exceed estimates')
      }
    }

    // Factor 4: Resource availability
    if (task.assignedTo && task.assignedTo.workload > 1.0) {
      riskScore += 25
      factors.push('Resource overallocation')
    }

    // Factor 5: Current progress
    if (task.progress < 50 && task.duration > 7) {
      const elapsed = task.duration - task.remaining
      if (elapsed > task.duration * 0.7) {
        riskScore += 35
        factors.push('Behind schedule')
      }
    }

    if (riskScore >= 50) {
      atRiskTasks.push({
        taskId: task.id,
        taskName: task.name,
        riskScore,
        riskLevel: riskScore >= 80 ? 'high' : riskScore >= 60 ? 'medium' : 'low',
        factors,
        predictedDelay: Math.round(task.duration * (riskScore / 100))
      })
    }
  })

  return atRiskTasks.sort((a, b) => b.riskScore - a.riskScore)
}

/**
 * Optimize schedule (generate optimized task dates)
 * @param {Array} tasks - Current tasks
 * @param {Array} dependencies - Task dependencies
 * @param {Array} resources - Available resources
 * @returns {Object} Optimized schedule
 */
export function optimizeSchedule(tasks, dependencies, resources) {
  // Clone tasks to avoid modifying original
  const optimizedTasks = JSON.parse(JSON.stringify(tasks))

  // 1. Apply resource leveling
  const leveledTasks = applyResourceLeveling(optimizedTasks, resources)

  // 2. Optimize task sequence
  const sequencedTasks = optimizeTaskSequence(leveledTasks, dependencies)

  // 3. Minimize project duration
  const compressedTasks = compressSchedule(sequencedTasks, dependencies)

  // Calculate improvements
  const originalDuration = Math.max(...tasks.map(t => t.start + t.duration))
  const optimizedDuration = Math.max(...compressedTasks.map(t => t.start + t.duration))
  const timeSaved = originalDuration - optimizedDuration

  return {
    tasks: compressedTasks,
    originalDuration,
    optimizedDuration,
    timeSaved,
    improvementPercentage: ((timeSaved / originalDuration) * 100).toFixed(1)
  }
}

/**
 * Apply resource leveling
 */
function applyResourceLeveling(tasks, resources) {
  // Simplified resource leveling
  const resourceLoad = new Map()

  return tasks.map(task => {
    if (!task.assignedTo) return task

    const load = resourceLoad.get(task.assignedTo) || 0
    if (load >= 1.0) {
      // Shift task to reduce overallocation
      const delay = Math.ceil(load)
      return {
        ...task,
        start: task.start + delay
      }
    }

    resourceLoad.set(task.assignedTo, load + 1)
    return task
  })
}

/**
 * Optimize task sequence
 */
function optimizeTaskSequence(tasks, dependencies) {
  // Sort tasks by dependency level
  const levels = new Map()

  function getLevel(taskId) {
    if (levels.has(taskId)) return levels.get(taskId)

    const predecessors = dependencies
      .filter(d => d.successorId === taskId)
      .map(d => d.predecessorId)

    if (predecessors.length === 0) {
      levels.set(taskId, 0)
      return 0
    }

    const maxPredLevel = Math.max(...predecessors.map(p => getLevel(p)))
    levels.set(taskId, maxPredLevel + 1)
    return levels.get(taskId)
  }

  tasks.forEach(task => getLevel(task.id))

  // Sort by level and start date
  return [...tasks].sort((a, b) => {
    const levelA = levels.get(a.id)
    const levelB = levels.get(b.id)
    if (levelA !== levelB) return levelA - levelB
    return a.start - b.start
  })
}

/**
 * Compress schedule (reduce gaps)
 */
function compressSchedule(tasks, dependencies) {
  const compressed = [...tasks]

  // Shift tasks left where possible
  compressed.forEach(task => {
    const predecessors = dependencies
      .filter(d => d.successorId === task.id)
      .map(d => compressed.find(t => t.id === d.predecessorId))
      .filter(Boolean)

    if (predecessors.length > 0) {
      const maxPredecessorEnd = Math.max(...predecessors.map(p => p.start + p.duration))
      if (task.start > maxPredecessorEnd) {
        task.start = maxPredecessorEnd
      }
    } else {
      // No predecessors, can start at 0
      task.start = 0
    }
  })

  return compressed
}

/**
 * Generate AI suggestion object
 * @param {string} type - Suggestion type
 * @param {string} title - Suggestion title
 * @param {string} description - Detailed description
 * @param {string} impact - Impact level (high, medium, low)
 * @param {Array} actions - Suggested actions
 * @returns {Object} Formatted suggestion
 */
export function createSuggestion(type, title, description, impact, actions) {
  return {
    id: `suggestion-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
    type,
    title,
    description,
    impact,
    actions,
    status: 'pending',
    createdAt: new Date().toISOString()
  }
}
