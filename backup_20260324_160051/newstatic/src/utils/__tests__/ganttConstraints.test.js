import { describe, it, expect, beforeEach } from 'vitest'
import {
  validateConstraint,
  applyConstraint,
  calculateConstraintImpact,
  checkConstraintConflicts,
  getConstraintSuggestion,
  CONSTRAINT_TYPES,
} from '../ganttConstraints.js'

describe('Gantt Constraints', () => {
  let task, calendar, dependencies

  beforeEach(() => {
    task = {
      id: 'task-1',
      name: 'Test Task',
      start: '2024-01-10T08:00:00Z',
      end: '2024-01-15T17:00:00Z',
      duration: 5,
      progress: 50,
    }

    calendar = {
      id: 'cal-1',
      workingDays: [1, 2, 3, 4, 5],
      workingHours: { start: '08:00', end: '17:00' },
      holidays: [],
    }

    dependencies = []
  })

  describe('validateConstraint', () => {
    it('should validate a valid must-start-on constraint', () => {
      const constraint = {
        type: CONSTRAINT_TYPES.MUST_START_ON,
        date: '2024-01-10T08:00:00Z',
      }

      const result = validateConstraint(task, constraint, calendar, dependencies)

      expect(result.valid).toBe(true)
      expect(result.errors).toHaveLength(0)
    })

    it('should reject must-start-on constraint with conflicting dependencies', () => {
      const constraint = {
        type: CONSTRAINT_TYPES.MUST_START_ON,
        date: '2024-01-05T08:00:00Z',
      }

      dependencies = [
        {
          id: 'dep-1',
          from: 'task-2',
          to: 'task-1',
          type: 'finish-to-start',
          lag: 0,
        },
      ]

      const result = validateConstraint(task, constraint, calendar, dependencies)

      expect(result.valid).toBe(false)
      expect(result.errors.length).toBeGreaterThan(0)
    })

    it('should validate a valid must-finish-by constraint', () => {
      const constraint = {
        type: CONSTRAINT_TYPES.MUST_FINISH_BY,
        date: '2024-01-20T17:00:00Z',
      }

      const result = validateConstraint(task, constraint, calendar, dependencies)

      expect(result.valid).toBe(true)
    })

    it('should reject must-finish-by constraint that violates task duration', () => {
      const constraint = {
        type: CONSTRAINT_TYPES.MUST_FINISH_BY,
        date: '2024-01-12T17:00:00Z',
      }

      const result = validateConstraint(task, constraint, calendar, dependencies)

      expect(result.valid).toBe(false)
    })

    it('should validate start-no-earlier-than constraint', () => {
      const constraint = {
        type: CONSTRAINT_TYPES.START_NO_EARLIER_THAN,
        date: '2024-01-08T08:00:00Z',
      }

      const result = validateConstraint(task, constraint, calendar, dependencies)

      expect(result.valid).toBe(true)
    })

    it('should reject start-no-earlier-than constraint that is violated', () => {
      const constraint = {
        type: CONSTRAINT_TYPES.START_NO_EARLIER_THAN,
        date: '2024-01-12T08:00:00Z',
      }

      const result = validateConstraint(task, constraint, calendar, dependencies)

      expect(result.valid).toBe(false)
    })

    it('should validate finish-no-later-than constraint', () => {
      const constraint = {
        type: CONSTRAINT_TYPES.FINISH_NO_LATER_THAN,
        date: '2024-01-20T17:00:00Z',
      }

      const result = validateConstraint(task, constraint, calendar, dependencies)

      expect(result.valid).toBe(true)
    })

    it('should reject invalid constraint type', () => {
      const constraint = {
        type: 'invalid-type',
        date: '2024-01-10T08:00:00Z',
      }

      const result = validateConstraint(task, constraint, calendar, dependencies)

      expect(result.valid).toBe(false)
      expect(result.errors.some(e => e.includes('Invalid constraint type'))).toBe(true)
    })

    it('should handle holidays in validation', () => {
      const constraint = {
        type: CONSTRAINT_TYPES.MUST_START_ON,
        date: '2024-01-15T08:00:00Z',
      }

      calendar.holidays = [{ date: '2024-01-15', name: 'Holiday' }]

      const result = validateConstraint(task, constraint, calendar, dependencies)

      expect(result.valid).toBe(false)
      expect(result.errors.some(e => e.includes('holiday'))).toBe(true)
    })
  })

  describe('applyConstraint', () => {
    it('should apply must-start-on constraint correctly', () => {
      const constraint = {
        type: CONSTRAINT_TYPES.MUST_START_ON,
        date: '2024-01-12T08:00:00Z',
      }

      const updated = applyConstraint(task, constraint, calendar)

      expect(updated.start).toBe('2024-01-12T08:00:00Z')
      expect(updated.end).toBe('2024-01-17T17:00:00Z') // Duration preserved
    })

    it('should apply must-finish-by constraint correctly', () => {
      const constraint = {
        type: CONSTRAINT_TYPES.MUST_FINISH_BY,
        date: '2024-01-20T17:00:00Z',
      }

      const updated = applyConstraint(task, constraint, calendar)

      expect(updated.end).toBe('2024-01-20T17:00:00Z')
      expect(updated.start).toBe('2024-01-15T08:00:00Z') // Duration preserved
    })

    it('should apply start-no-earlier-than constraint', () => {
      const constraint = {
        type: CONSTRAINT_TYPES.START_NO_EARLIER_THAN,
        date: '2024-01-12T08:00:00Z',
      }

      const updated = applyConstraint(task, constraint, calendar)

      expect(updated.start).toBe('2024-01-12T08:00:00Z')
    })

    it('should not adjust task if it already satisfies start-no-earlier-than', () => {
      const constraint = {
        type: CONSTRAINT_TYPES.START_NO_EARLIER_THAN,
        date: '2024-01-08T08:00:00Z',
      }

      const updated = applyConstraint(task, constraint, calendar)

      expect(updated.start).toBe(task.start)
    })

    it('should apply finish-no-later-than constraint', () => {
      const constraint = {
        type: CONSTRAINT_TYPES.FINISH_NO_LATER_THAN,
        date: '2024-01-14T17:00:00Z',
      }

      const updated = applyConstraint(task, constraint, calendar)

      expect(updated.end).toBe('2024-01-14T17:00:00Z')
    })

    it('should handle flexible constraints without changing dates', () => {
      const constraint = {
        type: CONSTRAINT_TYPES.AS_SOON_AS_POSSIBLE,
      }

      const updated = applyConstraint(task, constraint, calendar)

      expect(updated.start).toBe(task.start)
      expect(updated.end).toBe(task.end)
    })
  })

  describe('calculateConstraintImpact', () => {
    it('should calculate impact for must-start-on constraint', () => {
      const constraint = {
        type: CONSTRAINT_TYPES.MUST_START_ON,
        date: '2024-01-12T08:00:00Z',
      }

      const impact = calculateConstraintImpact(task, constraint, calendar, dependencies)

      expect(impact).toHaveProperty('startDateChange')
      expect(impact).toHaveProperty('endDateChange')
      expect(impact).toHaveProperty('affectedTasks')
      expect(impact.affectedTasks).toContain(task.id)
    })

    it('should calculate impact with dependent tasks', () => {
      const constraint = {
        type: CONSTRAINT_TYPES.MUST_START_ON,
        date: '2024-01-12T08:00:00Z',
      }

      dependencies = [
        {
          id: 'dep-1',
          from: task.id,
          to: 'task-2',
          type: 'finish-to-start',
          lag: 0,
        },
      ]

      const impact = calculateConstraintImpact(task, constraint, calendar, dependencies)

      expect(impact.affectedTasks).toContain('task-2')
    })

    it('should return zero impact for flexible constraints', () => {
      const constraint = {
        type: CONSTRAINT_TYPES.AS_SOON_AS_POSSIBLE,
      }

      const impact = calculateConstraintImpact(task, constraint, calendar, dependencies)

      expect(impact.startDateChange).toBe(0)
      expect(impact.endDateChange).toBe(0)
    })

    it('should calculate impact in working days', () => {
      const constraint = {
        type: CONSTRAINT_TYPES.MUST_START_ON,
        date: '2024-01-15T08:00:00Z',
      }

      const impact = calculateConstraintImpact(task, constraint, calendar, dependencies)

      expect(impact.startDateChange).toBeGreaterThan(0)
      expect(impact.workingDaysChanged).toBe(true)
    })
  })

  describe('checkConstraintConflicts', () => {
    it('should detect no conflicts when constraints are compatible', () => {
      const constraints = [
        {
          id: 'con-1',
          taskId: task.id,
          type: CONSTRAINT_TYPES.START_NO_EARLIER_THAN,
          date: '2024-01-08T08:00:00Z',
        },
        {
          id: 'con-2',
          taskId: task.id,
          type: CONSTRAINT_TYPES.FINISH_NO_LATER_THAN,
          date: '2024-01-20T17:00:00Z',
        },
      ]

      const conflicts = checkConstraintConflicts(constraints, calendar)

      expect(conflicts).toHaveLength(0)
    })

    it('should detect conflicts between incompatible constraints', () => {
      const constraints = [
        {
          id: 'con-1',
          taskId: task.id,
          type: CONSTRAINT_TYPES.MUST_START_ON,
          date: '2024-01-10T08:00:00Z',
        },
        {
          id: 'con-2',
          taskId: task.id,
          type: CONSTRAINT_TYPES.MUST_START_ON,
          date: '2024-01-15T08:00:00Z',
        },
      ]

      const conflicts = checkConstraintConflicts(constraints, calendar)

      expect(conflicts.length).toBeGreaterThan(0)
    })

    it('should detect conflicts with duration constraints', () => {
      const constraints = [
        {
          id: 'con-1',
          taskId: task.id,
          type: CONSTRAINT_TYPES.START_NO_EARLIER_THAN,
          date: '2024-01-10T08:00:00Z',
        },
        {
          id: 'con-2',
          taskId: task.id,
          type: CONSTRAINT_TYPES.FINISH_NO_LATER_THAN,
          date: '2024-01-12T17:00:00Z',
        },
      ]

      const conflicts = checkConstraintConflicts(constraints, calendar)

      expect(conflicts.length).toBeGreaterThan(0)
    })

    it('should check conflicts across multiple tasks', () => {
      const constraints = [
        {
          id: 'con-1',
          taskId: 'task-1',
          type: CONSTRAINT_TYPES.MUST_START_ON,
          date: '2024-01-10T08:00:00Z',
        },
        {
          id: 'con-2',
          taskId: 'task-2',
          type: CONSTRAINT_TYPES.MUST_FINISH_BY,
          date: '2024-01-12T17:00:00Z',
        },
      ]

      dependencies = [
        {
          id: 'dep-1',
          from: 'task-1',
          to: 'task-2',
          type: 'finish-to-start',
          lag: 0,
        },
      ]

      const conflicts = checkConstraintConflicts(constraints, calendar, dependencies)

      expect(conflicts.length).toBeGreaterThan(0)
    })
  })

  describe('getConstraintSuggestion', () => {
    it('should suggest as-soon-as-possible for typical forward scheduling', () => {
      const suggestion = getConstraintSuggestion(task, [], calendar)

      expect(suggestion.type).toBe(CONSTRAINT_TYPES.AS_SOON_AS_POSSIBLE)
    })

    it('should suggest must-start-on when task has fixed start', () => {
      const suggestion = getConstraintSuggestion(
        task,
        [{ type: CONSTRAINT_TYPES.MUST_START_ON, date: task.start }],
        calendar
      )

      expect(suggestion.type).toBeDefined()
    })

    it('should suggest finish-no-later-than for deadline-driven tasks', () => {
      const deadlineTask = { ...task, priority: 'high', hasDeadline: true }

      const suggestion = getConstraintSuggestion(deadlineTask, [], calendar)

      expect(suggestion.type).toBe(CONSTRAINT_TYPES.FINISH_NO_LATER_THAN)
    })

    it('should provide reasoning for suggestion', () => {
      const suggestion = getConstraintSuggestion(task, [], calendar)

      expect(suggestion).toHaveProperty('reasoning')
      expect(typeof suggestion.reasoning).toBe('string')
    })

    it('should suggest appropriate constraint for milestone', () => {
      const milestone = { ...task, type: 'milestone' }

      const suggestion = getConstraintSuggestion(milestone, [], calendar)

      expect(suggestion.type).toBeDefined()
    })
  })

  describe('Edge Cases', () => {
    it('should handle null calendar gracefully', () => {
      const constraint = {
        type: CONSTRAINT_TYPES.MUST_START_ON,
        date: '2024-01-10T08:00:00Z',
      }

      const result = validateConstraint(task, constraint, null, dependencies)

      expect(result).toBeDefined()
    })

    it('should handle empty dependencies', () => {
      const constraint = {
        type: CONSTRAINT_TYPES.MUST_START_ON,
        date: '2024-01-10T08:00:00Z',
      }

      const result = validateConstraint(task, constraint, calendar, [])

      expect(result.valid).toBe(true)
    })

    it('should handle task with zero duration', () => {
      const milestone = { ...task, duration: 0, start: task.end }

      const constraint = {
        type: CONSTRAINT_TYPES.MUST_START_ON,
        date: '2024-01-10T08:00:00Z',
      }

      const result = validateConstraint(milestone, constraint, calendar, dependencies)

      expect(result).toBeDefined()
    })

    it('should handle invalid date formats', () => {
      const constraint = {
        type: CONSTRAINT_TYPES.MUST_START_ON,
        date: 'invalid-date',
      }

      const result = validateConstraint(task, constraint, calendar, dependencies)

      expect(result.valid).toBe(false)
    })

    it('should handle weekend dates for work calendar', () => {
      const constraint = {
        type: CONSTRAINT_TYPES.MUST_START_ON,
        date: '2024-01-13T08:00:00Z', // Saturday
      }

      const result = validateConstraint(task, constraint, calendar, dependencies)

      expect(result.valid).toBe(false)
    })
  })
})
