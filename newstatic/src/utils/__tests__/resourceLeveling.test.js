import { describe, it, expect, beforeEach } from 'vitest'
import {
  detectConflicts,
  calculateAllocation,
  levelResources,
  findOverAllocations,
  calculateResourceUtilization,
  suggestLevelingActions,
  AllocationType,
} from '../resourceLeveling.js'

describe('Resource Leveling', () => {
  let tasks, resources, calendar

  beforeEach(() => {
    tasks = [
      {
        id: 'task-1',
        name: 'Task 1',
        start: '2024-01-01T08:00:00Z',
        end: '2024-01-05T17:00:00Z',
        duration: 5,
        resources: ['res-1'],
        assignee: 'user-1',
        priority: 'high',
        progress: 0,
      },
      {
        id: 'task-2',
        name: 'Task 2',
        start: '2024-01-01T08:00:00Z',
        end: '2024-01-05T17:00:00Z',
        duration: 5,
        resources: ['res-1', 'res-2'],
        assignee: 'user-1',
        priority: 'medium',
        progress: 0,
      },
      {
        id: 'task-3',
        name: 'Task 3',
        start: '2024-01-03T08:00:00Z',
        end: '2024-01-08T17:00:00Z',
        duration: 5,
        resources: ['res-1'],
        assignee: 'user-1',
        priority: 'low',
        progress: 0,
      },
    ]

    resources = [
      {
        id: 'res-1',
        name: 'John Developer',
        type: 'work',
        capacity: 8,
        unit: 'hours',
      },
      {
        id: 'res-2',
        name: 'Jane Designer',
        type: 'work',
        capacity: 8,
        unit: 'hours',
      },
      {
        id: 'res-3',
        name: 'Meeting Room',
        type: 'material',
        capacity: 1,
        unit: 'rooms',
      },
    ]

    calendar = {
      id: 'cal-1',
      workingDays: [1, 2, 3, 4, 5],
      workingHours: { start: '08:00', end: '17:00' },
      holidays: [],
    }
  })

  describe('detectConflicts', () => {
    it('should detect resource overallocation', () => {
      const conflicts = detectConflicts(tasks, resources, calendar)

      expect(conflicts.length).toBeGreaterThan(0)
      expect(conflicts[0].resourceId).toBe('res-1')
      expect(conflicts[0].type).toBe(AllocationType.OVERALLOCATED)
    })

    it('should detect multiple resource conflicts', () => {
      const conflictingTasks = [...tasks]
      conflictingTasks[1].resources = ['res-2', 'res-3']

      const conflicts = detectConflicts(conflictingTasks, resources, calendar)

      expect(conflicts.length).toBeGreaterThan(1)
    })

    it('should return empty array when no conflicts exist', () => {
      const nonConflictingTasks = [
        {
          ...tasks[0],
          start: '2024-01-01T08:00:00Z',
          end: '2024-01-03T17:00:00Z',
        },
        {
          ...tasks[1],
          start: '2024-01-04T08:00:00Z',
          end: '2024-01-08T17:00:00Z',
        },
      ]

      const conflicts = detectConflicts(nonConflictingTasks, resources, calendar)

      expect(conflicts).toHaveLength(0)
    })

    it('should calculate conflict severity', () => {
      const conflicts = detectConflicts(tasks, resources, calendar)

      expect(conflicts[0]).toHaveProperty('severity')
      expect(conflicts[0].severity).toBeGreaterThan(0)
    })

    it('should include affected tasks in conflict', () => {
      const conflicts = detectConflicts(tasks, resources, calendar)

      expect(conflicts[0]).toHaveProperty('affectedTasks')
      expect(conflicts[0].affectedTasks.length).toBeGreaterThan(1)
    })

    it('should handle material resource conflicts', () => {
      const materialTasks = [
        {
          id: 'task-1',
          start: '2024-01-01T08:00:00Z',
          end: '2024-01-05T17:00:00Z',
          resources: ['res-3'],
        },
        {
          id: 'task-2',
          start: '2024-01-01T08:00:00Z',
          end: '2024-01-05T17:00:00Z',
          resources: ['res-3'],
        },
      ]

      const conflicts = detectConflicts(materialTasks, resources, calendar)

      expect(conflicts.length).toBeGreaterThan(0)
      expect(conflicts[0].resourceId).toBe('res-3')
    })
  })

  describe('calculateAllocation', () => {
    it('should calculate correct allocation for single task', () => {
      const allocation = calculateAllocation(tasks[0], resources[0], calendar)

      expect(allocation).toHaveProperty('hours')
      expect(allocation).toHaveProperty('percentage')
      expect(allocation.percentage).toBeLessThanOrEqual(100)
    })

    it('should calculate allocation for overlapping tasks', () => {
      const allocation = calculateAllocation(
        [tasks[0], tasks[1]],
        resources[0],
        calendar
      )

      expect(allocation.hours).toBeGreaterThan(resources[0].capacity)
      expect(allocation.percentage).toBeGreaterThan(100)
    })

    it('should handle non-overlapping tasks', () => {
      const allocation = calculateAllocation(
        [
          { ...tasks[0], end: '2024-01-02T17:00:00Z' },
          { ...tasks[1], start: '2024-01-03T08:00:00Z' },
        ],
        resources[0],
        calendar
      )

      expect(allocation.percentage).toBeLessThanOrEqual(100)
    })

    it('should consider resource capacity', () => {
      const allocation = calculateAllocation(tasks, resources[0], calendar)

      expect(allocation.capacityUsed).toBeLessThanOrEqual(resources[0].capacity)
    })

    it('should handle partial overlap', () => {
      const allocation = calculateAllocation(
        [tasks[0], tasks[2]],
        resources[0],
        calendar
      )

      expect(allocation).toHaveProperty('overlapDays')
      expect(allocation.overlapDays).toBeGreaterThan(0)
    })

    it('should calculate allocation by date range', () => {
      const startDate = '2024-01-01T00:00:00Z'
      const endDate = '2024-01-10T23:59:59Z'

      const allocation = calculateAllocation(
        tasks,
        resources[0],
        calendar,
        startDate,
        endDate
      )

      expect(allocation).toHaveProperty('dailyBreakdown')
      expect(allocation.dailyBreakdown).toBeInstanceOf(Array)
    })
  })

  describe('levelResources', () => {
    it('should level overallocated resources', () => {
      const conflicts = detectConflicts(tasks, resources, calendar)
      const leveled = levelResources(tasks, resources, conflicts, calendar)

      expect(leveled).toHaveProperty('adjustedTasks')
      expect(leveled.adjustedTasks.length).toBeGreaterThan(0)
    })

    it('should respect task priority during leveling', () => {
      const conflicts = detectConflicts(tasks, resources, calendar)
      const leveled = levelResources(tasks, resources, conflicts, calendar, {
        prioritizeBy: 'priority',
      })

      // High priority tasks should not be delayed
      const highPriorityTask = leveled.adjustedTasks.find(t => t.priority === 'high')
      expect(highPriorityTask.start).toBe(tasks[0].start)
    })

    it('should minimize total project duration', () => {
      const conflicts = detectConflicts(tasks, resources, calendar)
      const leveled = levelResources(tasks, resources, conflicts, calendar, {
        minimizeProjectDuration: true,
      })

      const originalEnd = new Date(tasks[tasks.length - 1].end)
      const leveledEnd = new Date(leveled.adjustedTasks[leveled.adjustedTasks.length - 1].end)

      // Leveled should not significantly extend project
      expect(leveledEnd - originalEnd).toBeLessThan(7 * 24 * 60 * 60 * 1000) // 1 week
    })

    it('should preserve task dependencies', () => {
      const tasksWithDeps = [
        ...tasks,
        {
          id: 'task-4',
          name: 'Task 4',
          start: '2024-01-06T08:00:00Z',
          end: '2024-01-10T17:00:00Z',
          duration: 5,
          resources: ['res-1'],
          dependsOn: ['task-1'],
        },
      ]

      const conflicts = detectConflicts(tasksWithDeps, resources, calendar)
      const leveled = levelResources(tasksWithDeps, resources, conflicts, calendar)

      const task4 = leveled.adjustedTasks.find(t => t.id === 'task-4')
      const task1 = leveled.adjustedTasks.find(t => t.id === 'task-1')

      expect(new Date(task4.start)).toBeGreaterThanOrEqual(new Date(task1.end))
    })

    it('should return leveling summary', () => {
      const conflicts = detectConflicts(tasks, resources, calendar)
      const leveled = levelResources(tasks, resources, conflicts, calendar)

      expect(leveled).toHaveProperty('summary')
      expect(leveled.summary).toHaveProperty('tasksAdjusted')
      expect(leveled.summary).toHaveProperty('conflictsResolved')
    })

    it('should respect working calendar', () => {
      calendar.holidays = [{ date: '2024-01-02', name: 'Holiday' }]

      const conflicts = detectConflicts(tasks, resources, calendar)
      const leveled = levelResources(tasks, resources, conflicts, calendar)

      // Tasks should be scheduled around holidays
      expect(leveled.adjustedTasks).toBeDefined()
    })
  })

  describe('findOverAllocations', () => {
    it('should find all overallocated resources', () => {
      const overAllocations = findOverAllocations(tasks, resources, calendar)

      expect(overAllocations.length).toBeGreaterThan(0)
      expect(overAllocations[0].allocation.percentage).toBeGreaterThan(100)
    })

    it('should include allocation details', () => {
      const overAllocations = findOverAllocations(tasks, resources, calendar)

      expect(overAllocations[0]).toHaveProperty('resourceId')
      expect(overAllocations[0]).toHaveProperty('allocation')
      expect(overAllocations[0]).toHaveProperty('tasks')
    })

    it('should return empty array when no overallocations', () => {
      const nonConflictingTasks = [
        {
          ...tasks[0],
          resources: ['res-1'],
          end: '2024-01-03T17:00:00Z',
        },
        {
          ...tasks[1],
          resources: ['res-2'],
          start: '2024-01-04T08:00:00Z',
        },
      ]

      const overAllocations = findOverAllocations(nonConflictingTasks, resources, calendar)

      expect(overAllocations).toHaveLength(0)
    })

    it('should calculate overallocation percentage', () => {
      const overAllocations = findOverAllocations(tasks, resources, calendar)

      expect(overAllocations[0].allocation.percentage).toBeGreaterThan(100)
      expect(overAllocations[0].allocation.overAllocation).toBeGreaterThan(0)
    })
  })

  describe('calculateResourceUtilization', () => {
    it('should calculate overall utilization', () => {
      const utilization = calculateResourceUtilization(tasks, resources, calendar)

      expect(utilization).toHaveProperty('overall')
      expect(utilization.overall).toBeGreaterThanOrEqual(0)
      expect(utilization.overall).toBeLessThanOrEqual(100)
    })

    it('should calculate per-resource utilization', () => {
      const utilization = calculateResourceUtilization(tasks, resources, calendar)

      expect(utilization).toHaveProperty('byResource')
      expect(Object.keys(utilization.byResource)).toContain('res-1')
      expect(Object.keys(utilization.byResource)).toContain('res-2')
    })

    it('should calculate utilization by time period', () => {
      const utilization = calculateResourceUtilization(tasks, resources, calendar, {
        groupBy: 'week',
      })

      expect(utilization).toHaveProperty('byTimePeriod')
      expect(utilization.byTimePeriod).toBeInstanceOf(Array)
    })

    it('should identify underutilized resources', () => {
      const utilization = calculateResourceUtilization(tasks, resources, calendar)

      expect(utilization).toHaveProperty('underutilized')
      expect(utilization.underutilized).toBeInstanceOf(Array)
    })

    it('should identify overutilized resources', () => {
      const utilization = calculateResourceUtilization(tasks, resources, calendar)

      expect(utilization).toHaveProperty('overutilized')
      expect(utilization.overutilized).toBeInstanceOf(Array)
    })
  })

  describe('suggestLevelingActions', () => {
    it('should suggest task splitting for long tasks', () => {
      const longTask = {
        ...tasks[0],
        duration: 20,
        end: '2024-01-25T17:00:00Z',
      }

      const actions = suggestLevelingActions([longTask], resources, calendar)

      expect(actions.some(a => a.type === 'split')).toBe(true)
    })

    it('should suggest resource substitution', () => {
      const actions = suggestLevelingActions(tasks, resources, calendar)

      expect(actions.some(a => a.type === 'substitute')).toBe(true)
    })

    it('should suggest task delay for low priority tasks', () => {
      const actions = suggestLevelingActions(tasks, resources, calendar)

      expect(actions.some(a => a.type === 'delay')).toBe(true)
    })

    it('should suggest adding resources for critical tasks', () => {
      const criticalTasks = tasks.map(t => ({ ...t, priority: 'critical' }))

      const actions = suggestLevelingActions(criticalTasks, resources, calendar)

      expect(actions.some(a => a.type === 'addResource')).toBe(true)
    })

    it('should provide action priority', () => {
      const actions = suggestLevelingActions(tasks, resources, calendar)

      actions.forEach(action => {
        expect(action).toHaveProperty('priority')
      })
    })

    it('should estimate impact of each action', () => {
      const actions = suggestLevelingActions(tasks, resources, calendar)

      actions.forEach(action => {
        expect(action).toHaveProperty('impact')
      })
    })

    it('should not suggest actions for balanced schedule', () => {
      const balancedTasks = [
        {
          ...tasks[0],
          resources: ['res-1'],
          end: '2024-01-03T17:00:00Z',
        },
        {
          ...tasks[1],
          resources: ['res-2'],
          start: '2024-01-04T08:00:00Z',
          end: '2024-01-08T17:00:00Z',
        },
      ]

      const actions = suggestLevelingActions(balancedTasks, resources, calendar)

      expect(actions.length).toBe(0)
    })
  })

  describe('Edge Cases', () => {
    it('should handle empty task list', () => {
      const conflicts = detectConflicts([], resources, calendar)

      expect(conflicts).toHaveLength(0)
    })

    it('should handle tasks with no resources', () => {
      const tasksWithoutResources = tasks.map(t => ({ ...t, resources: [] }))

      const conflicts = detectConflicts(tasksWithoutResources, resources, calendar)

      expect(conflicts).toHaveLength(0)
    })

    it('should handle null calendar', () => {
      const conflicts = detectConflicts(tasks, resources, null)

      expect(conflicts).toBeDefined()
    })

    it('should handle invalid resource capacity', () => {
      const invalidResources = [
        { ...resources[0], capacity: -1 },
      ]

      const allocation = calculateAllocation(tasks, invalidResources[0], calendar)

      expect(allocation).toBeDefined()
    })

    it('should handle zero duration tasks', () => {
      const milestone = {
        ...tasks[0],
        duration: 0,
        end: tasks[0].start,
      }

      const allocation = calculateAllocation([milestone], resources[0], calendar)

      expect(allocation.hours).toBe(0)
    })
  })
})
