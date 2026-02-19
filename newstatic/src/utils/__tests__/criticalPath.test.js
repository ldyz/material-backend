import { describe, it, expect, beforeEach } from 'vitest'
import {
  calculateCriticalPath,
  calculateSlack,
  identifyCriticalTasks,
  calculateEarlyStart,
  calculateLateFinish,
  getCriticalPathLength,
  DependencyType,
} from '../criticalPath.js'

describe('Critical Path Calculator', () => {
  let tasks, dependencies, calendar

  beforeEach(() => {
    tasks = [
      {
        id: 'task-1',
        name: 'Start',
        start: '2024-01-01T08:00:00Z',
        end: '2024-01-03T17:00:00Z',
        duration: 3,
      },
      {
        id: 'task-2',
        name: 'Task A',
        start: '2024-01-04T08:00:00Z',
        end: '2024-01-08T17:00:00Z',
        duration: 5,
      },
      {
        id: 'task-3',
        name: 'Task B',
        start: '2024-01-04T08:00:00Z',
        end: '2024-01-10T17:00:00Z',
        duration: 7,
      },
      {
        id: 'task-4',
        name: 'Task C',
        start: '2024-01-09T08:00:00Z',
        end: '2024-01-12T17:00:00Z',
        duration: 4,
      },
      {
        id: 'task-5',
        name: 'End',
        start: '2024-01-13T08:00:00Z',
        end: '2024-01-15T17:00:00Z',
        duration: 3,
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
        to: 'task-5',
        type: DependencyType.FINISH_TO_START,
        lag: 0,
      },
      {
        id: 'dep-5',
        from: 'task-4',
        to: 'task-5',
        type: DependencyType.FINISH_TO_START,
        lag: 0,
      },
    ]

    calendar = {
      workingDays: [1, 2, 3, 4, 5],
      workingHours: { start: '08:00', end: '17:00' },
      holidays: [],
    }
  })

  describe('calculateCriticalPath', () => {
    it('should calculate the critical path correctly', () => {
      const result = calculateCriticalPath(tasks, dependencies, calendar)

      expect(result.criticalPath).toBeDefined()
      expect(result.criticalPath.length).toBeGreaterThan(0)
    })

    it('should include start and end tasks in critical path', () => {
      const result = calculateCriticalPath(tasks, dependencies, calendar)

      expect(result.criticalPath).toContain('task-1')
      expect(result.criticalPath).toContain('task-5')
    })

    it('should identify the longest path through the network', () => {
      const result = calculateCriticalPath(tasks, dependencies, calendar)

      // task-1 -> task-3 -> task-5 is the longest path (3 + 7 + 3 = 13 days)
      expect(result.criticalPath).toContain('task-3')
    })

    it('should calculate project duration', () => {
      const result = calculateCriticalPath(tasks, dependencies, calendar)

      expect(result.projectDuration).toBeDefined()
      expect(result.projectDuration).toBeGreaterThan(0)
    })

    it('should calculate early start/finish times', () => {
      const result = calculateCriticalPath(tasks, dependencies, calendar)

      tasks.forEach(task => {
        expect(result.earlyStart[task.id]).toBeDefined()
        expect(result.earlyFinish[task.id]).toBeDefined()
      })
    })

    it('should calculate late start/finish times', () => {
      const result = calculateCriticalPath(tasks, dependencies, calendar)

      tasks.forEach(task => {
        expect(result.lateStart[task.id]).toBeDefined()
        expect(result.lateFinish[task.id]).toBeDefined()
      })
    })

    it('should handle finish-to-start dependencies', () => {
      const result = calculateCriticalPath(tasks, dependencies, calendar)

      expect(result.criticalPath.length).toBeGreaterThan(0)
    })

    it('should handle start-to-start dependencies', () => {
      const ssDeps = [
        {
          id: 'dep-1',
          from: 'task-1',
          to: 'task-2',
          type: DependencyType.START_TO_START,
          lag: 0,
        },
      ]

      const result = calculateCriticalPath(tasks, ssDeps, calendar)

      expect(result).toBeDefined()
    })

    it('should handle finish-to-finish dependencies', () => {
      const ffDeps = [
        {
          id: 'dep-1',
          from: 'task-1',
          to: 'task-2',
          type: DependencyType.FINISH_TO_FINISH,
          lag: 0,
        },
      ]

      const result = calculateCriticalPath(tasks, ffDeps, calendar)

      expect(result).toBeDefined()
    })

    it('should handle start-to-finish dependencies', () => {
      const sfDeps = [
        {
          id: 'dep-1',
          from: 'task-1',
          to: 'task-2',
          type: DependencyType.START_TO_FINISH,
          lag: 0,
        },
      ]

      const result = calculateCriticalPath(tasks, sfDeps, calendar)

      expect(result).toBeDefined()
    })

    it('should handle lag time', () => {
      const lagDeps = [
        {
          id: 'dep-1',
          from: 'task-1',
          to: 'task-2',
          type: DependencyType.FINISH_TO_START,
          lag: 2,
        },
      ]

      const result = calculateCriticalPath(tasks, lagDeps, calendar)

      expect(result).toBeDefined()
    })

    it('should handle negative lag (lead time)', () => {
      const leadDeps = [
        {
          id: 'dep-1',
          from: 'task-1',
          to: 'task-2',
          type: DependencyType.FINISH_TO_START,
          lag: -1,
        },
      ]

      const result = calculateCriticalPath(tasks, leadDeps, calendar)

      expect(result).toBeDefined()
    })
  })

  describe('calculateSlack', () => {
    it('should calculate zero slack for critical tasks', () => {
      const result = calculateCriticalPath(tasks, dependencies, calendar)

      result.criticalPath.forEach(taskId => {
        expect(result.slack[taskId]).toBe(0)
      })
    })

    it('should calculate positive slack for non-critical tasks', () => {
      const result = calculateCriticalPath(tasks, dependencies, calendar)

      // task-2 and task-4 should have slack
      if (!result.criticalPath.includes('task-2')) {
        expect(result.slack['task-2']).toBeGreaterThan(0)
      }
    })

    it('should calculate total slack', () => {
      const result = calculateCriticalPath(tasks, dependencies, calendar)

      tasks.forEach(task => {
        expect(result.slack[task.id]).toBeGreaterThanOrEqual(0)
      })
    })

    it('should calculate free slack', () => {
      const result = calculateCriticalPath(tasks, dependencies, calendar)

      expect(result.freeSlack).toBeDefined()
      tasks.forEach(task => {
        expect(result.freeSlack[task.id]).toBeGreaterThanOrEqual(0)
      })
    })

    it('should calculate independent slack', () => {
      const result = calculateCriticalPath(tasks, dependencies, calendar)

      expect(result.independentSlack).toBeDefined()
    })

    it('should handle tasks with no dependencies', () => {
      const noDepsResult = calculateCriticalPath(tasks, [], calendar)

      Object.keys(noDepsResult.slack).forEach(taskId => {
        expect(noDepsResult.slack[taskId]).toBeDefined()
      })
    })
  })

  describe('identifyCriticalTasks', () => {
    it('should identify all critical tasks', () => {
      const result = calculateCriticalPath(tasks, dependencies, calendar)
      const criticalTasks = identifyCriticalTasks(result)

      expect(criticalTasks).toEqual(expect.arrayContaining(result.criticalPath))
    })

    it('should return array of task IDs', () => {
      const result = calculateCriticalPath(tasks, dependencies, calendar)
      const criticalTasks = identifyCriticalTasks(result)

      criticalTasks.forEach(taskId => {
        expect(typeof taskId).toBe('string')
      })
    })

    it('should handle empty critical path', () => {
      const criticalTasks = identifyCriticalTasks({ criticalPath: [] })

      expect(criticalTasks).toEqual([])
    })
  })

  describe('calculateEarlyStart', () => {
    it('should calculate early start for first task', () => {
      const es = calculateEarlyStart('task-1', tasks, dependencies, calendar)

      expect(es).toBeDefined()
    })

    it('should consider predecessor tasks', () => {
      const es = calculateEarlyStart('task-2', tasks, dependencies, calendar)

      expect(es).toBeDefined()
      const task1End = new Date(tasks[0].end)
      const task2Start = new Date(es)
      expect(task2Start.getTime()).toBeGreaterThanOrEqual(task1End.getTime())
    })

    it('should handle multiple predecessors', () => {
      const es = calculateEarlyStart('task-5', tasks, dependencies, calendar)

      expect(es).toBeDefined()
    })

    it('should respect lag time', () => {
      const lagDeps = [
        {
          id: 'dep-1',
          from: 'task-1',
          to: 'task-2',
          type: DependencyType.FINISH_TO_START,
          lag: 3,
        },
      ]

      const es = calculateEarlyStart('task-2', tasks, lagDeps, calendar)

      expect(es).toBeDefined()
    })

    it('should handle different dependency types', () => {
      const ssDep = [
        {
          id: 'dep-1',
          from: 'task-1',
          to: 'task-2',
          type: DependencyType.START_TO_START,
          lag: 0,
        },
      ]

      const es = calculateEarlyStart('task-2', tasks, ssDep, calendar)

      expect(es).toBeDefined()
    })
  })

  describe('calculateLateFinish', () => {
    it('should calculate late finish for last task', () => {
      const lf = calculateLateFinish('task-5', tasks, dependencies, calendar)

      expect(lf).toBeDefined()
    })

    it('should consider successor tasks', () => {
      const lf = calculateLateFinish('task-3', tasks, dependencies, calendar)

      expect(lf).toBeDefined()
    })

    it('should handle multiple successors', () => {
      const lf = calculateLateFinish('task-1', tasks, dependencies, calendar)

      expect(lf).toBeDefined()
    })

    it('should equal early finish for critical tasks', () => {
      const result = calculateCriticalPath(tasks, dependencies, calendar)

      result.criticalPath.forEach(taskId => {
        const earlyFinish = result.earlyFinish[taskId]
        const lateFinish = result.lateFinish[taskId]
        expect(new Date(earlyFinish).getTime()).toBe(new Date(lateFinish).getTime())
      })
    })
  })

  describe('getCriticalPathLength', () => {
    it('should calculate path length in days', () => {
      const result = calculateCriticalPath(tasks, dependencies, calendar)
      const length = getCriticalPathLength(result)

      expect(length).toBeGreaterThan(0)
      expect(typeof length).toBe('number')
    })

    it('should calculate path length in hours', () => {
      const result = calculateCriticalPath(tasks, dependencies, calendar)
      const length = getCriticalPathLength(result, 'hours')

      expect(length).toBeGreaterThan(0)
    })

    it('should return 0 for empty critical path', () => {
      const length = getCriticalPathLength({ criticalPath: [] })

      expect(length).toBe(0)
    })

    it('should account for working days only', () => {
      const result = calculateCriticalPath(tasks, dependencies, calendar)
      const length = getCriticalPathLength(result)

      // Length should be in working days
      expect(length).toBeLessThanOrEqual(tasks.reduce((sum, t) => sum + t.duration, 0))
    })
  })

  describe('Edge Cases', () => {
    it('should handle circular dependencies gracefully', () => {
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
      ]

      const result = calculateCriticalPath(tasks, circularDeps, calendar)

      expect(result).toBeDefined()
      expect(result.error).toBeDefined()
    })

    it('should handle tasks with zero duration', () => {
      const milestone = {
        id: 'milestone-1',
        name: 'Milestone',
        start: '2024-01-05T17:00:00Z',
        end: '2024-01-05T17:00:00Z',
        duration: 0,
      }

      const milestoneDeps = [
        {
          id: 'dep-1',
          from: 'task-1',
          to: 'milestone-1',
          type: DependencyType.FINISH_TO_START,
          lag: 0,
        },
      ]

      const result = calculateCriticalPath([...tasks, milestone], milestoneDeps, calendar)

      expect(result).toBeDefined()
    })

    it('should handle empty task list', () => {
      const result = calculateCriticalPath([], dependencies, calendar)

      expect(result.criticalPath).toEqual([])
    })

    it('should handle empty dependency list', () => {
      const result = calculateCriticalPath(tasks, [], calendar)

      expect(result).toBeDefined()
    })

    it('should handle holidays in calendar', () => {
      const holidayCalendar = {
        ...calendar,
        holidays: [
          { date: '2024-01-02', name: 'Holiday' },
          { date: '2024-01-03', name: 'Holiday' },
        ],
      }

      const result = calculateCriticalPath(tasks, dependencies, holidayCalendar)

      expect(result).toBeDefined()
      expect(result.projectDuration).toBeGreaterThan(0)
    })

    it('should handle weekend calendar', () => {
      const weekendCalendar = {
        ...calendar,
        workingDays: [1, 2, 3, 4], // Mon-Thu only
      }

      const result = calculateCriticalPath(tasks, dependencies, weekendCalendar)

      expect(result).toBeDefined()
    })
  })

  describe('Multiple Critical Paths', () => {
    it('should detect multiple critical paths', () => {
      // Create two equal-length paths
      const equalDeps = [
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
          to: 'task-5',
          type: DependencyType.FINISH_TO_START,
          lag: 0,
        },
        {
          id: 'dep-4',
          from: 'task-3',
          to: 'task-5',
          type: DependencyType.FINISH_TO_START,
          lag: 0,
        },
      ]

      const result = calculateCriticalPath(tasks, equalDeps, calendar)

      expect(result).toHaveProperty('multiplePaths')
      expect(result.multiplePaths).toBe(true)
    })

    it('should identify all critical paths', () => {
      const result = calculateCriticalPath(tasks, dependencies, calendar)

      if (result.multiplePaths) {
        expect(result.allCriticalPaths).toBeDefined()
        expect(Array.isArray(result.allCriticalPaths)).toBe(true)
      }
    })
  })

  describe('Performance', () => {
    it('should handle large task lists efficiently', () => {
      const largeTasks = Array.from({ length: 100 }, (_, i) => ({
        id: `task-${i}`,
        name: `Task ${i}`,
        start: new Date(Date.now() + i * 86400000).toISOString(),
        end: new Date(Date.now() + (i + 1) * 86400000).toISOString(),
        duration: 1,
      }))

      const largeDeps = Array.from({ length: 99 }, (_, i) => ({
        id: `dep-${i}`,
        from: `task-${i}`,
        to: `task-${i + 1}`,
        type: DependencyType.FINISH_TO_START,
        lag: 0,
      }))

      const start = performance.now()
      const result = calculateCriticalPath(largeTasks, largeDeps, calendar)
      const end = performance.now()

      expect(end - start).toBeLessThan(1000) // Should complete in less than 1 second
      expect(result).toBeDefined()
    })
  })
})
