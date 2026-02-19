/**
 * Workflow Automation Engine
 * Defines automation rules and executes them based on triggers
 */

import { ElMessage } from 'element-plus'

/**
 * Automation rule types
 */
export const RuleTypes = {
  AUTO_DEPENDENCY: 'auto_dependency',
  AUTO_ASSIGN: 'auto_assign',
  AUTO_MILESTONE: 'auto_milestone',
  NOTIFICATION: 'notification',
  AUTO_DURATION: 'auto_duration',
  AUTO_PROGRESS: 'auto_progress',
  DATE_VALIDATION: 'date_validation'
}

/**
 * Trigger types
 */
export const TriggerTypes = {
  ON_TASK_CREATE: 'on_task_create',
  ON_TASK_UPDATE: 'on_task_update',
  ON_TASK_DELETE: 'on_task_delete',
  ON_DEPENDENCY_ADD: 'on_dependency_add',
  ON_STATUS_CHANGE: 'on_status_change',
  ON_DATE_CHANGE: 'on_date_change',
  ON_ASSIGNMENT_CHANGE: 'on_assignment_change',
  MANUAL: 'manual'
}

/**
 * Predefined automation rules
 */
export const defaultRules = [
  {
    id: 'auto-dependency-naming',
    name: 'Auto-set dependencies based on task naming',
    description: 'Automatically create dependencies when tasks are named with prefixes like "Phase 1", "Phase 2"',
    type: RuleTypes.AUTO_DEPENDENCY,
    trigger: TriggerTypes.ON_TASK_CREATE,
    enabled: true,
    config: {
      patterns: [
        { regex: /^(Phase|阶段)\s*(\d+)/i, group: 2, order: 'numeric' },
        { regex: /^(Step|步骤)\s*(\d+)/i, group: 2, order: 'numeric' },
        { regex: /^(\d+)\./, group: 1, order: 'numeric' }
      ]
    }
  },
  {
    id: 'auto-assign-by-type',
    name: 'Auto-assign based on task type',
    description: 'Automatically assign tasks to default resources based on task type or keywords',
    type: RuleTypes.AUTO_ASSIGN,
    trigger: TriggerTypes.ON_TASK_CREATE,
    enabled: false,
    config: {
      assignments: [
        { keywords: ['design', 'ui', 'ux', '设计'], resource: 'Designer' },
        { keywords: ['development', 'backend', 'frontend', '开发'], resource: 'Developer' },
        { keywords: ['test', 'qa', 'testing', '测试'], resource: 'QA Engineer' },
        { keywords: ['review', 'approve', 'review', '审批'], resource: 'Project Manager' }
      ]
    }
  },
  {
    id: 'auto-milestone-dates',
    name: 'Auto-update milestone dates',
    description: 'Automatically adjust milestone dates based on dependent task completion',
    type: RuleTypes.AUTO_MILESTONE,
    trigger: TriggerTypes.ON_TASK_UPDATE,
    enabled: true,
    config: {
      lookAhead: 7 // days
    }
  },
  {
    id: 'notify-delay-risk',
    name: 'Notify delay risks',
    description: 'Send notifications when tasks are at risk of delay',
    type: RuleTypes.NOTIFICATION,
    trigger: TriggerTypes.ON_STATUS_CHANGE,
    enabled: true,
    config: {
      threshold: 0.8 // 80% of duration passed
    }
  },
  {
    id: 'auto-duration-estimation',
    name: 'Auto-estimate duration',
    description: 'Suggest duration based on task name and type',
    type: RuleTypes.AUTO_DURATION,
    trigger: TriggerTypes.ON_TASK_CREATE,
    enabled: false,
    config: {
      defaults: {
        meeting: 1,
        review: 2,
        development: 5,
        testing: 3,
        deployment: 1
      }
    }
  },
  {
    id: 'date-validation',
    name: 'Validate date logic',
    description: 'Check for invalid date sequences and circular dependencies',
    type: RuleTypes.DATE_VALIDATION,
    trigger: TriggerTypes.ON_DEPENDENCY_ADD,
    enabled: true,
    config: {
      checkCircular: true,
      checkPastDates: true
    }
  }
]

/**
 * Workflow Automation Engine class
 */
export class WorkflowAutomationEngine {
  constructor() {
    this.rules = new Map()
    this.history = []
    this.enabled = true
    this.loadRules()
  }

  /**
   * Load rules from storage or use defaults
   */
  loadRules() {
    try {
      const stored = localStorage.getItem('workflow-rules')
      if (stored) {
        const rules = JSON.parse(stored)
        rules.forEach(rule => this.rules.set(rule.id, rule))
      } else {
        // Load default rules
        defaultRules.forEach(rule => this.rules.set(rule.id, rule))
        this.saveRules()
      }
    } catch (error) {
      console.error('Error loading workflow rules:', error)
      defaultRules.forEach(rule => this.rules.set(rule.id, rule))
    }
  }

  /**
   * Save rules to storage
   */
  saveRules() {
    try {
      const rules = Array.from(this.rules.values())
      localStorage.setItem('workflow-rules', JSON.stringify(rules))
    } catch (error) {
      console.error('Error saving workflow rules:', error)
    }
  }

  /**
   * Add or update a rule
   */
  addRule(rule) {
    this.rules.set(rule.id, rule)
    this.saveRules()
  }

  /**
   * Remove a rule
   */
  removeRule(ruleId) {
    this.rules.delete(ruleId)
    this.saveRules()
  }

  /**
   * Enable/disable a rule
   */
  toggleRule(ruleId, enabled) {
    const rule = this.rules.get(ruleId)
    if (rule) {
      rule.enabled = enabled
      this.saveRules()
    }
  }

  /**
   * Execute rules based on trigger
   */
  async execute(trigger, context) {
    if (!this.enabled) {
      return []
    }

    const results = []

    for (const rule of this.rules.values()) {
      if (!rule.enabled || rule.trigger !== trigger) {
        continue
      }

      try {
        const result = await this.executeRule(rule, context)
        if (result) {
          results.push({
            rule: rule.name,
            ...result
          })
          this.history.push({
            timestamp: new Date().toISOString(),
            rule: rule.name,
            trigger,
            result
          })
        }
      } catch (error) {
        console.error(`Error executing rule ${rule.name}:`, error)
      }
    }

    return results
  }

  /**
   * Execute a single rule
   */
  async executeRule(rule, context) {
    switch (rule.type) {
      case RuleTypes.AUTO_DEPENDENCY:
        return this.executeAutoDependency(rule, context)

      case RuleTypes.AUTO_ASSIGN:
        return this.executeAutoAssign(rule, context)

      case RuleTypes.AUTO_MILESTONE:
        return this.executeAutoMilestone(rule, context)

      case RuleTypes.NOTIFICATION:
        return this.executeNotification(rule, context)

      case RuleTypes.AUTO_DURATION:
        return this.executeAutoDuration(rule, context)

      case RuleTypes.DATE_VALIDATION:
        return this.executeDateValidation(rule, context)

      default:
        console.warn(`Unknown rule type: ${rule.type}`)
        return null
    }
  }

  /**
   * Execute auto-dependency rule
   */
  executeAutoDependency(rule, context) {
    const { task, allTasks } = context
    if (!task || !allTasks) return null

    const results = []

    for (const pattern of rule.config.patterns) {
      const match = task.name.match(pattern.regex)
      if (match) {
        const currentValue = parseInt(match[pattern.group])

        // Find previous task
        const previousValue = currentValue - 1
        const previousTask = allTasks.find(t => {
          const prevMatch = t.name.match(pattern.regex)
          return prevMatch && parseInt(prevMatch[pattern.group]) === previousValue
        })

        if (previousTask) {
          results.push({
            type: 'dependency',
            from: previousTask.id,
            to: task.id,
            message: `Auto-created dependency: "${previousTask.name}" → "${task.name}"`
          })
        }
      }
    }

    return results.length > 0 ? { actions: results } : null
  }

  /**
   * Execute auto-assign rule
   */
  executeAutoAssign(rule, context) {
    const { task } = context
    if (!task || !task.name) return null

    const taskNameLower = task.name.toLowerCase()

    for (const assignment of rule.config.assignments) {
      if (assignment.keywords.some(keyword => taskNameLower.includes(keyword))) {
        return {
          type: 'assignment',
          resource: assignment.resource,
          message: `Auto-assigned "${task.name}" to ${assignment.resource}`
        }
      }
    }

    return null
  }

  /**
   * Execute auto-milestone rule
   */
  executeAutoMilestone(rule, context) {
    const { tasks, dependencies } = context
    if (!tasks || !dependencies) return null

    const results = []
    const milestones = tasks.filter(t => t.isMilestone)

    for (const milestone of milestones) {
      // Find predecessors
      const predecessors = dependencies
        .filter(d => d.successorId === milestone.id)
        .map(d => tasks.find(t => t.id === d.predecessorId))
        .filter(Boolean)

      if (predecessors.length === 0) continue

      // Calculate latest finish date
      const latestFinish = Math.max(...predecessors.map(t => {
        return t.start + t.duration
      }))

      // Check if milestone date needs adjustment
      if (latestFinish > milestone.start) {
        results.push({
          type: 'milestone_adjustment',
          milestone: milestone.id,
          currentDate: milestone.start,
          suggestedDate: latestFinish,
          message: `Milestone "${milestone.name}" should move to day ${latestFinish}`
        })
      }
    }

    return results.length > 0 ? { actions: results } : null
  }

  /**
   * Execute notification rule
   */
  executeNotification(rule, context) {
    const { task } = context
    if (!task) return null

    // Check if task is at risk
    const elapsed = task.duration - task.remaining
    const progressRatio = elapsed / task.duration

    if (progressRatio >= rule.config.threshold && task.status !== 'completed') {
      return {
        type: 'notification',
        level: 'warning',
        message: `Task "${task.name}" is at risk of delay (${Math.round(progressRatio * 100)}% complete)`,
        taskId: task.id
      }
    }

    return null
  }

  /**
   * Execute auto-duration rule
   */
  executeAutoDuration(rule, context) {
    const { task } = context
    if (!task || !task.name) return null

    const taskNameLower = task.name.toLowerCase()

    for (const [type, duration] of Object.entries(rule.config.defaults)) {
      if (taskNameLower.includes(type)) {
        return {
          type: 'duration_suggestion',
          duration,
          message: `Suggested duration for "${task.name}": ${duration} days`
        }
      }
    }

    return null
  }

  /**
   * Execute date validation rule
   */
  executeDateValidation(rule, context) {
    const { dependency, tasks } = context
    if (!dependency || !tasks) return null

    const results = []

    // Check for circular dependencies
    if (rule.config.checkCircular) {
      const hasCircular = this.checkCircularDependency(dependency, tasks)
      if (hasCircular) {
        results.push({
          type: 'circular_dependency',
          severity: 'error',
          message: 'Circular dependency detected'
        })
      }
    }

    // Check for past dates
    if (rule.config.checkPastDates) {
      const predecessor = tasks.find(t => t.id === dependency.predecessorId)
      const successor = tasks.find(t => t.id === dependency.successorId)

      if (predecessor && successor) {
        if (predecessor.start + predecessor.duration > successor.start) {
          results.push({
            type: 'invalid_date_sequence',
            severity: 'warning',
            message: `Task sequence: "${predecessor.name}" should finish before "${successor.name}" starts`
          })
        }
      }
    }

    return results.length > 0 ? { actions: results } : null
  }

  /**
   * Check for circular dependencies
   */
  checkCircularDependency(newDependency, tasks) {
    const visited = new Set()
    const recStack = new Set()

    // Build adjacency list
    const graph = new Map()
    tasks.forEach(task => {
      graph.set(task.id, [])
    })

    // Add existing dependencies (mock - would come from context)
    // Add new dependency
    if (graph.has(newDependency.successorId)) {
      graph.get(newDependency.successorId).push(newDependency.predecessorId)
    }

    // DFS cycle detection
    const hasCycle = (node) => {
      if (recStack.has(node)) return true
      if (visited.has(node)) return false

      visited.add(node)
      recStack.add(node)

      const neighbors = graph.get(node) || []
      for (const neighbor of neighbors) {
        if (hasCycle(neighbor)) return true
      }

      recStack.delete(node)
      return false
    }

    for (const task of tasks) {
      if (hasCycle(task.id)) return true
    }

    return false
  }

  /**
   * Get execution history
   */
  getHistory(limit = 50) {
    return this.history.slice(-limit)
  }

  /**
   * Clear history
   */
  clearHistory() {
    this.history = []
  }

  /**
   * Enable/disable all rules
   */
  setEnabled(enabled) {
    this.enabled = enabled
  }

  /**
   * Get all rules
   */
  getAllRules() {
    return Array.from(this.rules.values())
  }

  /**
   * Get rule by ID
   */
  getRule(ruleId) {
    return this.rules.get(ruleId)
  }
}

// Singleton instance
let engineInstance = null

/**
 * Get workflow automation engine instance
 */
export function getWorkflowEngine() {
  if (!engineInstance) {
    engineInstance = new WorkflowAutomationEngine()
  }
  return engineInstance
}

/**
 * Execute automation based on trigger (convenience function)
 */
export async function executeAutomation(trigger, context) {
  const engine = getWorkflowEngine()
  return await engine.execute(trigger, context)
}
