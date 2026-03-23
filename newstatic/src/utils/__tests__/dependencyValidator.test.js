import { describe, it, expect, beforeEach } from 'vitest'
import {
  validateDependency,
  detectCircularDependencies,
  validateLag,
  calculateCriticalPath,
  analyzeDependencyPath,
  DependencyType,
} from '../dependencyValidator.js'

describe('Dependency Validator', () => {
  let tasks, dependencies

  beforeEach(() => {
    tasks = [
      {
        id: 'task-1',
        name: 'Task 1',
        start: '2024-01-01T08:00:00Z',
        end: '2024-01-05T17:00:00Z',
        duration: 5,
      },
      {
        id: 'task-2',
        name: 'Task 2',
        start: '2024-01-06T08:00:00Z',
        end: '2024-01-10T17:00:00Z',
        duration: 5,
      },
      {
        id: 'task-3',
        name: 'Task 3',
        start: '2024-01-11T08:00:00Z',
        end: '2024-01-15T17:00:00Z',
        duration: 5,
      },
      {
        id: 'task-4',
        name: 'Task 4',
        start: '2024-01-16T08:00:00Z',
        end: '2024-01-20T17:00:00Z',
        duration: 5,
      },
    ]

    dependencies = [
      {
        id: 'dep-1',
        from: 'task-1',
        to: 'task-2',
        type: DependencyType.FINISH_TO_START,
        lag: 0,
      },
      {
        id: 'dep-2',
        from: 'task-2',
        to: 'task-3',
        type: DependencyType.FINISH_TO_START,
        lag: 0,
      },
    ]
  })

  describe('validateDependency', () => {
    it('should validate valid finish-to-start dependency', () => {
      const dependency = {
        id: 'dep-new',
        from: 'task-3',
        to: 'task-4',
        type: DependencyType.FINISH_TO_START,
        lag: 0,
      }

      const result = validateDependency(dependency, tasks, dependencies)

      expect(result.valid).toBe(true)
      expect(result.errors).toHaveLength(0)
    })

    it('should reject dependency to non-existent task', () => {
      const dependency = {
        id: 'dep-new',
        from: 'task-1',
        to: 'task-non-existent',
        type: DependencyType.FINISH_TO_START,
        lag: 0,
      }

      const result = validateDependency(dependency, tasks, dependencies)

      expect(result.valid).toBe(false)
      expect(result.errors.some(e => e.includes('not found'))).toBe(true)
    })

    it('should reject self-referencing dependency', () => {
      const dependency = {
        id: 'dep-new',
        from: 'task-1',
        to: 'task-1',
        type: DependencyType.FINISH_TO_START,
        lag: 0,
      }

      const result = validateDependency(dependency, tasks, dependencies)

      expect(result.valid).toBe(false)
      expect(result.errors.some(e => e.includes('same task'))).toBe(true)
    })

    it('should reject duplicate dependency', () => {
      const dependency = {
        id: 'dep-new',
        from: 'task-1',
        to: 'task-2',
        type: DependencyType.FINISH_TO_START,
        lag: 0,
      }

      const result = validateDependency(dependency, tasks, dependencies)

      expect(result.valid).toBe(false)
      expect(result.errors.some(e => e.includes('already exists'))).toBe(true)
    })

    it('should validate start-to-start dependency', () => {
      const dependency = {
        id: 'dep-new',
        from: 'task-1',
        to: 'task-2',
        type: DependencyType.START_TO_START,
        lag: 0,
      }

      const existingDeps = dependencies.filter(d => d.from !== 'task-1')
      const result = validateDependency(dependency, tasks, existingDeps)

      expect(result.valid).toBe(true)
    })

    it('should validate finish-to-finish dependency', () => {
      const dependency = {
        id: 'dep-new',
        from: 'task-2',
        to: 'task-3',
        type: DependencyType.FINISH_TO_FINISH,
        lag: 0,
      }

      const existingDeps = dependencies.filter(d => d.from !== 'task-2')
      const result = validateDependency(dependency, tasks, existingDeps)

      expect(result.valid).toBe(true)
    })

    it('should validate start-to-finish dependency', () => {
      const dependency = {
        id: 'dep-new',
        from: 'task-1',
        to: 'task-2',
        type: DependencyType.START_TO_FINISH,
        lag: 0,
      }

      const existingDeps = dependencies.filter(d => d.from !== 'task-1')
      const result = validateDependency(dependency, tasks, existingDeps)

      expect(result.valid).toBe(true)
    })

    it('should reject invalid dependency type', () => {
      const dependency = {
        id: 'dep-new',
        from: 'task-1',
        to: 'task-2',
        type: 'invalid-type',
        lag: 0,
      }

      const result = validateDependency(dependency, tasks, [])

      expect(result.valid).toBe(false)
    })

    it('should validate dependency with positive lag', () => {
      const dependency = {
        id: 'dep-new',
        from: 'task-3',
        to: 'task-4',
        type: DependencyType.FINISH_TO_START,
        lag: 2,
      }

      const result = validateDependency(dependency, tasks, [])

      expect(result.valid).toBe(true)
    })

    it('should validate dependency with negative lag (lead)', () => {
      const dependency = {
        id: 'dep-new',
        from: 'task-1',
        to: 'task-2',
        type: DependencyType.FINISH_TO_START,
        lag: -1,
      }

      const existingDeps = dependencies.filter(d => d.from !== 'task-1')
      const result = validateDependency(dependency, tasks, existingDeps)

      expect(result.valid).toBe(true)
    })

    it('should check temporal consistency', () => {
      const dependency = {
        id: 'dep-new',
        from: 'task-2',
        to: 'task-1',
        type: DependencyType.FINISH_TO_START,
        lag: 0,
      }

      const result = validateDependency(dependency, tasks, [])

      expect(result.valid).toBe(false)
      expect(result.errors.some(e => e.includes('temporal'))).toBe(true)
    })
  })

  describe('detectCircularDependencies', () => {
    it('should detect simple circular dependency', () => {
      const circularDeps = [
        ...dependencies,
        {
          id: 'dep-3',
          from: 'task-3',
          to: 'task-1',
          type: DependencyType.FINISH_TO_START,
          lag: 0,
        },
      ]

      const result = detectCircularDependencies(tasks, circularDeps)

      expect(result.hasCycles).toBe(true)
      expect(result.cycles.length).toBeGreaterThan(0)
    })

    it('should detect complex circular dependency chain', () => {
      const circularDeps = [
        ...dependencies,
        {
          id: 'dep-3',
          from: 'task-3',
          to: 'task-4',
          type: DependencyType.FINISH_TO_START,
          lag: 0,
        },
        {
          id: 'dep-4',
          from: 'task-4',
          to: 'task-1',
          type: DependencyType.FINISH_TO_START,
          lag: 0,
        },
      ]

      const result = detectCircularDependencies(tasks, circularDeps)

      expect(result.hasCycles).toBe(true)
    })

    it('should return no cycles for valid dependencies', () => {
      const result = detectCircularDependencies(tasks, dependencies)

      expect(result.hasCycles).toBe(false)
      expect(result.cycles).toHaveLength(0)
    })

    it('should identify all tasks in cycle', () => {
      const circularDeps = [
        ...dependencies,
        {
          id: 'dep-3',
          from: 'task-3',
          to: 'task-1',
          type: DependencyType.FINISH_TO_START,
          lag: 0,
        },
      ]

      const result = detectCircularDependencies(tasks, circularDeps)

      expect(result.cycles[0]).toContain('task-1')
      expect(result.cycles[0]).toContain('task-2')
      expect(result.cycles[0]).toContain('task-3')
    })

    it('should handle multiple independent cycles', () => {
      const circularDeps = [
        {
          id: 'dep-1',
          from: 'task-1',
          to: 'task-2',
          type: DependencyType.FINISH_TO_START,
          lag: 0,
        },
        {
          id: 'dep-2',
          from: 'task-2',
          to: 'task-1',
          type: DependencyType.FINISH_TO_START,
          lag: 0,
        },
        {
          id: 'dep-3',
          from: 'task-3',
          to: 'task-4',
          type: DependencyType.FINISH_TO_START,
          lag: 0,
        },
        {
          id: 'dep-4',
          from: 'task-4',
          to: 'task-3',
          type: DependencyType.FINISH_TO_START,
          lag: 0,
        },
      ]

      const result = detectCircularDependencies(tasks, circularDeps)

      expect(result.hasCycles).toBe(true)
      expect(result.cycles.length).toBeGreaterThanOrEqual(2)
    })

    it('should provide cycle path information', () => {
      const circularDeps = [
        ...dependencies,
        {
          id: 'dep-3',
          from: 'task-3',
          to: 'task-1',
          type: DependencyType.FINISH_TO_START,
          lag: 0,
        },
      ]

      const result = detectCircularDependencies(tasks, circularDeps)

      expect(result.cycles[0]).toHaveProperty('path')
      expect(Array.isArray(result.cycles[0].path)).toBe(true)
    })
  })

  describe('validateLag', () => {
    it('should validate zero lag', () => {
      const result = validateLag(0, tasks[0], tasks[1], DependencyType.FINISH_TO_START)

      expect(result.valid).toBe(true)
    })

    it('should validate positive lag', () => {
      const result = validateLag(3, tasks[0], tasks[1], DependencyType.FINISH_TO_START)

      expect(result.valid).toBe(true)
    })

    it('should validate negative lag (lead)', () => {
      const result = validateLag(-2, tasks[0], tasks[1], DependencyType.FINISH_TO_START)

      expect(result.valid).toBe(true)
    })

    it('should reject excessive lead time', () => {
      const result = validateLag(-10, tasks[0], tasks[1], DependencyType.FINISH_TO_START)

      expect(result.valid).toBe(false)
    })

    it('should check lag against task duration', () => {
      const result = validateLag(-6, tasks[0], tasks[1], DependencyType.FINISH_TO_START)

      expect(result.valid).toBe(false)
    })

    it('should validate lag for start-to-start dependencies', () => {
      const result = validateLag(2, tasks[0], tasks[1], DependencyType.START_TO_START)

      expect(result.valid).toBe(true)
    })

    it('should validate lag for finish-to-finish dependencies', () => {
      const result = validateLag(1, tasks[0], tasks[1], DependencyType.FINISH_TO_FINISH)

      expect(result.valid).toBe(true)
    })

    it('should provide suggestions for invalid lag', () => {
      const result = validateLag(-10, tasks[0], tasks[1], DependencyType.FINISH_TO_START)

      expect(result.valid).toBe(false)
      expect(result.suggestions).toBeDefined()
      expect(result.suggestions.length).toBeGreaterThan(0)
    })
  })

  describe('calculateCriticalPath', () => {
    it('should calculate critical path for linear dependencies', () => {
      const result = calculateCriticalPath(tasks, dependencies)

      expect(result).toHaveProperty('criticalPath')
      expect(result.criticalPath.length).toBeGreaterThan(0)
    })

    it('should include all tasks when all are critical', () => {
      const result = calculateCriticalPath(tasks, dependencies)

      // All tasks should be on critical path for linear chain
      expect(result.criticalPath).toHaveLength(3) // task-1, task-2, task-3
    })

    it('should calculate slack for non-critical tasks', () => {
      const nonCriticalTask = {
        id: 'task-5',
        name: 'Non-Critical Task',
        start: '2024-01-06T08:00:00Z',
        end: '2024-01-12T17:00:00Z',
        duration: 5,
      }

      const tasksWithNonCritical = [...tasks, nonCriticalTask]
      const result = calculateCriticalPath(tasksWithNonCritical, dependencies)

      expect(result.slack).toHaveProperty('task-5')
      expect(result.slack['task-5']).toBeGreaterThan(0)
    })

    it('should handle multiple dependency types', () => {
      const mixedDeps = [
        {
          id: 'dep-1',
          from: 'task-1',
          to: 'task-2',
          type: DependencyType.FINISH_TO_START,
          lag: 0,
        },
        {
          id: 'dep-2',
          from: 'task-1',
          to: 'task-3',
          type: DependencyType.START_TO_START,
          lag: 1,
        },
      ]

      const result = calculateCriticalPath(tasks, mixedDeps)

      expect(result.criticalPath).toBeDefined()
    })

    it('should calculate project duration', () => {
      const result = calculateCriticalPath(tasks, dependencies)

      expect(result).toHaveProperty('projectDuration')
      expect(result.projectDuration).toBeGreaterThan(0)
    })

    it('should identify multiple critical paths', () => {
      const parallelDeps = [
        {
          id: 'dep-1',
          from: 'task-1',
          to: 'task-2',
          type: DependencyType.FINISH_TO_START,
          lag: 0,
        },
        {
          id: 'dep-2',
          from: 'task-1',
          to: 'task-3',
          type: DependencyType.FINISH_TO_START,
          lag: 0,
        },
        {
          id: 'dep-3',
          from: 'task-2',
          to: 'task-4',
          type: DependencyType.FINISH_TO_START,
          lag: 0,
        },
        {
          id: 'dep-4',
          from: 'task-3',
          to: 'task-4',
          type: DependencyType.FINISH_TO_START,
          lag: 0,
        },
      ]

      const result = calculateCriticalPath(tasks, parallelDeps)

      expect(result).toHaveProperty('criticalPaths')
      expect(Array.isArray(result.criticalPaths)).toBe(true)
    })

    it('should calculate early and late dates', () => {
      const result = calculateCriticalPath(tasks, dependencies)

      result.criticalPath.forEach(taskId => {
        expect(result.dates[taskId]).toHaveProperty('earlyStart')
        expect(result.dates[taskId]).toHaveProperty('earlyFinish')
        expect(result.dates[taskId]).toHaveProperty('lateStart')
        expect(result.dates[taskId]).toHaveProperty('lateFinish')
      })
    })
  })

  describe('analyzeDependencyPath', () => {
    it('should analyze path from start to finish', () => {
      const result = analyzeDependencyPath('task-1', tasks, dependencies)

      expect(result).toHaveProperty('path')
      expect(result).toHaveProperty('length')
      expect(result.path[0]).toBe('task-1')
    })

    it('should calculate path duration', () => {
      const result = analyzeDependencyPath('task-1', tasks, dependencies)

      expect(result).toHaveProperty('duration')
      expect(result.duration).toBeGreaterThan(0)
    })

    it('should identify all successors', () => {
      const result = analyzeDependencyPath('task-1', tasks, dependencies)

      expect(result).toHaveProperty('successors')
      expect(result.successors).toContain('task-2')
    })

    it('should identify all predecessors', () => {
      const result = analyzeDependencyPath('task-3', tasks, dependencies)

      expect(result).toHaveProperty('predecessors')
      expect(result.predecessors).toContain('task-2')
    })

    it('should handle tasks with no dependencies', () => {
      const result = analyzeDependencyPath('task-4', tasks, dependencies)

      expect(result).toHaveProperty('path')
      expect(result.predecessors).toBeDefined()
    })

    it('should calculate path flexibility', () => {
      const result = analyzeDependencyPath('task-1', tasks, dependencies)

      expect(result).toHaveProperty('flexibility')
    })

    it('should identify merge points', () => {
      const mergeDeps = [
        {
          id: 'dep-1',
          from: 'task-1',
          to: 'task-3',
          type: DependencyType.FINISH_TO_START,
          lag: 0,
        },
        {
          id: 'dep-2',
          from: 'task-2',
          to: 'task-3',
          type: DependencyType.FINISH_TO_START,
          lag: 0,
        },
      ]

      const result = analyzeDependencyPath('task-1', tasks, mergeDeps)

      expect(result).toHaveProperty('mergePoints')
    })

    it('should identify burst points', () => {
      const burstDeps = [
        {
          id: 'dep-1',
          from: 'task-1',
          to: 'task-2',
          type: DependencyType.FINISH_TO_START,
          lag: 0,
        },
        {
          id: 'dep-2',
          from: 'task-1',
          to: 'task-3',
          type: DependencyType.FINISH_TO_START,
          lag: 0,
        },
      ]

      const result = analyzeDependencyPath('task-1', tasks, burstDeps)

      expect(result).toHaveProperty('burstPoints')
      expect(result.burstPoints.length).toBeGreaterThan(0)
    })
  })

  describe('Edge Cases', () => {
    it('should handle empty task list', () => {
      const result = detectCircularDependencies([], dependencies)

      expect(result.hasCycles).toBe(false)
    })

    it('should handle empty dependency list', () => {
      const result = detectCircularDependencies(tasks, [])

      expect(result.hasCycles).toBe(false)
    })

    it('should handle null lag values', () => {
      const result = validateLag(null, tasks[0], tasks[1], DependencyType.FINISH_TO_START)

      expect(result).toBeDefined()
    })

    it('should handle tasks with same dates', () => {
      const sameDateTasks = [
        { ...tasks[0] },
        { ...tasks[1], start: tasks[0].start, end: tasks[0].end },
      ]

      const result = validateDependency(
        {
          id: 'dep-1',
          from: 'task-1',
          to: 'task-2',
          type: DependencyType.FINISH_TO_START,
          lag: 0,
        },
        sameDateTasks,
        []
      )

      expect(result).toBeDefined()
    })

    it('should handle zero duration tasks (milestones)', () => {
      const milestone = {
        id: 'milestone-1',
        name: 'Milestone',
        start: '2024-01-05T17:00:00Z',
        end: '2024-01-05T17:00:00Z',
        duration: 0,
      }

      const result = validateDependency(
        {
          id: 'dep-1',
          from: 'task-1',
          to: 'milestone-1',
          type: DependencyType.FINISH_TO_START,
          lag: 0,
        },
        [...tasks, milestone],
        []
      )

      expect(result).toBeDefined()
    })
  })
})
