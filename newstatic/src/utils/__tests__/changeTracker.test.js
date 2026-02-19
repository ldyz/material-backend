import { describe, it, expect, beforeEach, vi } from 'vitest'
import {
  trackChanges,
  generateDiff,
  exportChanges,
  importChanges,
  compressChanges,
  calculateChangeImpact,
  ChangeType,
} from '../changeTracker.js'

describe('Change Tracker', () => {
  let originalTasks, modifiedTasks

  beforeEach(() => {
    originalTasks = [
      {
        id: 'task-1',
        name: 'Original Task 1',
        start: '2024-01-01T08:00:00Z',
        end: '2024-01-05T17:00:00Z',
        progress: 50,
        assignee: 'user-1',
        priority: 'high',
      },
      {
        id: 'task-2',
        name: 'Original Task 2',
        start: '2024-01-06T08:00:00Z',
        end: '2024-01-10T17:00:00Z',
        progress: 0,
        assignee: 'user-2',
        priority: 'medium',
      },
    ]

    modifiedTasks = [
      {
        id: 'task-1',
        name: 'Modified Task 1',
        start: '2024-01-02T08:00:00Z',
        end: '2024-01-06T17:00:00Z',
        progress: 75,
        assignee: 'user-1',
        priority: 'high',
      },
      {
        id: 'task-2',
        name: 'Original Task 2',
        start: '2024-01-06T08:00:00Z',
        end: '2024-01-10T17:00:00Z',
        progress: 0,
        assignee: 'user-2',
        priority: 'low',
      },
      {
        id: 'task-3',
        name: 'New Task 3',
        start: '2024-01-11T08:00:00Z',
        end: '2024-01-15T17:00:00Z',
        progress: 0,
        assignee: 'user-3',
        priority: 'medium',
      },
    ]
  })

  describe('trackChanges', () => {
    it('should detect name changes', () => {
      const changes = trackChanges(originalTasks, modifiedTasks)

      const task1Change = changes.find(c => c.entityId === 'task-1')
      expect(task1Change).toBeDefined()
      expect(task1Change.changes.name).toBeDefined()
      expect(task1Change.changes.name.from).toBe('Original Task 1')
      expect(task1Change.changes.name.to).toBe('Modified Task 1')
    })

    it('should detect date changes', () => {
      const changes = trackChanges(originalTasks, modifiedTasks)

      const task1Change = changes.find(c => c.entityId === 'task-1')
      expect(task1Change.changes.start).toBeDefined()
      expect(task1Change.changes.end).toBeDefined()
    })

    it('should detect progress changes', () => {
      const changes = trackChanges(originalTasks, modifiedTasks)

      const task1Change = changes.find(c => c.entityId === 'task-1')
      expect(task1Change.changes.progress.from).toBe(50)
      expect(task1Change.changes.progress.to).toBe(75)
    })

    it('should detect priority changes', () => {
      const changes = trackChanges(originalTasks, modifiedTasks)

      const task2Change = changes.find(c => c.entityId === 'task-2')
      expect(task2Change.changes.priority.from).toBe('medium')
      expect(task2Change.changes.priority.to).toBe('low')
    })

    it('should detect new tasks', () => {
      const changes = trackChanges(originalTasks, modifiedTasks)

      const newTask = changes.find(c => c.entityId === 'task-3')
      expect(newTask).toBeDefined()
      expect(newTask.type).toBe(ChangeType.CREATE)
    })

    it('should detect deleted tasks', () => {
      const changes = trackChanges(modifiedTasks, originalTasks)

      const deletedTask = changes.find(c => c.entityId === 'task-3')
      expect(deletedTask).toBeDefined()
      expect(deletedTask.type).toBe(ChangeType.DELETE)
    })

    it('should set correct change type for modifications', () => {
      const changes = trackChanges(originalTasks, modifiedTasks)

      const task1Change = changes.find(c => c.entityId === 'task-1')
      expect(task1Change.type).toBe(ChangeType.UPDATE)
    })

    it('should include timestamp', () => {
      const beforeTime = Date.now()
      const changes = trackChanges(originalTasks, modifiedTasks)
      const afterTime = Date.now()

      changes.forEach(change => {
        expect(change.timestamp).toBeDefined()
        expect(change.timestamp).toBeGreaterThanOrEqual(beforeTime)
        expect(change.timestamp).toBeLessThanOrEqual(afterTime)
      })
    })

    it('should handle empty arrays', () => {
      const changes = trackChanges([], [])

      expect(changes).toEqual([])
    })

    it('should handle completely new list', () => {
      const changes = trackChanges([], modifiedTasks)

      expect(changes.length).toBe(modifiedTasks.length)
      changes.forEach(change => {
        expect(change.type).toBe(ChangeType.CREATE)
      })
    })

    it('should handle completely deleted list', () => {
      const changes = trackChanges(originalTasks, [])

      expect(changes.length).toBe(originalTasks.length)
      changes.forEach(change => {
        expect(change.type).toBe(ChangeType.DELETE)
      })
    })
  })

  describe('generateDiff', () => {
    it('should generate readable diff', () => {
      const changes = trackChanges(originalTasks, modifiedTasks)
      const diff = generateDiff(changes)

      expect(diff).toBeDefined()
      expect(typeof diff).toBe('string')
      expect(diff.length).toBeGreaterThan(0)
    })

    it('should include all changes in diff', () => {
      const changes = trackChanges(originalTasks, modifiedTasks)
      const diff = generateDiff(changes)

      changes.forEach(change => {
        expect(diff).toContain(change.entityId)
      })
    })

    it('should format changes human-readable', () => {
      const changes = trackChanges(originalTasks, modifiedTasks)
      const diff = generateDiff(changes)

      expect(diff).toContain('changed')
      expect(diff).toContain('from')
      expect(diff).toContain('to')
    })

    it('should support HTML format', () => {
      const changes = trackChanges(originalTasks, modifiedTasks)
      const diff = generateDiff(changes, 'html')

      expect(diff).toContain('<')
      expect(diff).toContain('>')
    })

    it('should support markdown format', () => {
      const changes = trackChanges(originalTasks, modifiedTasks)
      const diff = generateDiff(changes, 'markdown')

      expect(diff).toContain('*') // Markdown formatting
    })

    it('should support JSON format', () => {
      const changes = trackChanges(originalTasks, modifiedTasks)
      const diff = generateDiff(changes, 'json')

      const parsed = JSON.parse(diff)
      expect(Array.isArray(parsed)).toBe(true)
    })
  })

  describe('exportChanges', () => {
    it('should export changes to JSON', () => {
      const changes = trackChanges(originalTasks, modifiedTasks)
      const exported = exportChanges(changes, 'json')

      expect(exported).toBeDefined()
      expect(typeof exported).toBe('string')

      const parsed = JSON.parse(exported)
      expect(Array.isArray(parsed)).toBe(true)
    })

    it('should export changes to CSV', () => {
      const changes = trackChanges(originalTasks, modifiedTasks)
      const exported = exportChanges(changes, 'csv')

      expect(exported).toBeDefined()
      expect(exported).toContain(',') // CSV format
    })

    it('should include metadata in export', () => {
      const changes = trackChanges(originalTasks, modifiedTasks)
      const exported = exportChanges(changes, 'json', {
        includeMetadata: true,
        exportedBy: 'test-user',
        exportDate: new Date().toISOString(),
      })

      const parsed = JSON.parse(exported)
      expect(parsed).toHaveProperty('metadata')
      expect(parsed.metadata).toHaveProperty('exportedBy')
    })

    it('should filter changes by type', () => {
      const changes = trackChanges(originalTasks, modifiedTasks)
      const exported = exportChanges(changes, 'json', {
        filter: { type: ChangeType.UPDATE },
      })

      const parsed = JSON.parse(exported)
      parsed.forEach(change => {
        expect(change.type).toBe(ChangeType.UPDATE)
      })
    })

    it('should filter changes by date range', () => {
      const changes = trackChanges(originalTasks, modifiedTasks)
      const startDate = new Date('2024-01-01T00:00:00Z')
      const endDate = new Date('2024-12-31T23:59:59Z')

      const exported = exportChanges(changes, 'json', {
        filter: { startDate, endDate },
      })

      const parsed = JSON.parse(exported)
      expect(parsed.length).toBeGreaterThan(0)
    })
  })

  describe('importChanges', () => {
    it('should import changes from JSON', () => {
      const changes = trackChanges(originalTasks, modifiedTasks)
      const exported = exportChanges(changes, 'json')
      const imported = importChanges(exported, 'json')

      expect(imported).toBeDefined()
      expect(Array.isArray(imported)).toBe(true)
      expect(imported.length).toBe(changes.length)
    })

    it('should import changes from CSV', () => {
      const changes = trackChanges(originalTasks, modifiedTasks)
      const exported = exportChanges(changes, 'csv')
      const imported = importChanges(exported, 'csv')

      expect(imported).toBeDefined()
      expect(Array.isArray(imported)).toBe(true)
    })

    it('should validate imported changes', () => {
      const changes = trackChanges(originalTasks, modifiedTasks)
      const exported = exportChanges(changes, 'json')
      const imported = importChanges(exported, 'json', { validate: true })

      expect(imported.valid).toBe(true)
      expect(imported.changes).toBeDefined()
    })

    it('should reject invalid import data', () => {
      const invalid = importChanges('invalid json', 'json')

      expect(invalid.valid).toBe(false)
      expect(invalid.errors).toBeDefined()
    })

    it('should merge with existing changes', () => {
      const changes1 = trackChanges(originalTasks, modifiedTasks)
      const changes2 = trackChanges(modifiedTasks, originalTasks)

      const exported = exportChanges(changes2, 'json')
      const imported = importChanges(exported, 'json', { merge: changes1 })

      expect(imported.length).toBeGreaterThan(changes1.length)
    })
  })

  describe('compressChanges', () => {
    it('should compress changes by grouping', () => {
      const changes = trackChanges(originalTasks, modifiedTasks)
      const compressed = compressChanges(changes)

      expect(compressed).toBeDefined()
      expect(compressed.length).toBeLessThanOrEqual(changes.length)
    })

    it('should group changes by entity', () => {
      const changes = trackChanges(originalTasks, modifiedTasks)
      const compressed = compressChanges(changes, { groupBy: 'entity' })

      // Check that changes for same entity are grouped
      const entityGroups = {}
      compressed.forEach(change => {
        if (!entityGroups[change.entityId]) {
          entityGroups[change.entityId] = []
        }
        entityGroups[change.entityId].push(change)
      })

      Object.values(entityGroups).forEach(group => {
        expect(group.length).toBe(1) // Compressed to single entry
      })
    })

    it('should remove redundant changes', () => {
      // Create changes where a field is changed multiple times
      const changes = [
        {
          entityId: 'task-1',
          type: ChangeType.UPDATE,
          changes: { progress: { from: 0, to: 25 } },
          timestamp: Date.now() - 2000,
        },
        {
          entityId: 'task-1',
          type: ChangeType.UPDATE,
          changes: { progress: { from: 25, to: 50 } },
          timestamp: Date.now() - 1000,
        },
        {
          entityId: 'task-1',
          type: ChangeType.UPDATE,
          changes: { progress: { from: 50, to: 75 } },
          timestamp: Date.now(),
        },
      ]

      const compressed = compressChanges(changes)

      expect(compressed.length).toBe(1)
      expect(compressed[0].changes.progress.from).toBe(0)
      expect(compressed[0].changes.progress.to).toBe(75)
    })

    it('should preserve change order', () => {
      const changes = trackChanges(originalTasks, modifiedTasks)
      const compressed = compressChanges(changes)

      compressed.forEach((change, index) => {
        if (index > 0) {
          expect(change.timestamp).toBeGreaterThanOrEqual(compressed[index - 1].timestamp)
        }
      })
    })
  })

  describe('calculateChangeImpact', () => {
    it('should calculate impact for date changes', () => {
      const changes = trackChanges(originalTasks, modifiedTasks)
      const task1Change = changes.find(c => c.entityId === 'task-1')

      const impact = calculateChangeImpact(task1Change, originalTasks)

      expect(impact).toHaveProperty('affectedTasks')
      expect(impact).toHaveProperty('timelineShift')
    })

    it('should calculate impact for dependency changes', () => {
      const changes = [
        {
          entityId: 'dep-1',
          entityType: 'dependency',
          type: ChangeType.UPDATE,
          changes: { lag: { from: 0, to: 2 } },
          timestamp: Date.now(),
        },
      ]

      const impact = calculateChangeImpact(changes[0], originalTasks, dependencies)

      expect(impact).toBeDefined()
    })

    it('should identify cascading impacts', () => {
      const changes = trackChanges(originalTasks, modifiedTasks)
      const task1Change = changes.find(c => c.entityId === 'task-1')

      const dependencies = [
        { from: 'task-1', to: 'task-2', type: 'finish-to-start' },
      ]

      const impact = calculateChangeImpact(task1Change, originalTasks, dependencies)

      expect(impact.affectedTasks).toContain('task-2')
    })

    it('should calculate severity', () => {
      const changes = trackChanges(originalTasks, modifiedTasks)
      const task1Change = changes.find(c => c.entityId === 'task-1')

      const impact = calculateChangeImpact(task1Change, originalTasks)

      expect(impact).toHaveProperty('severity')
      expect(['low', 'medium', 'high', 'critical']).toContain(impact.severity)
    })

    it('should estimate resource impact', () => {
      const changes = trackChanges(originalTasks, modifiedTasks)

      changes.forEach(change => {
        const impact = calculateChangeImpact(change, originalTasks)
        expect(impact).toHaveProperty('resourceImpact')
      })
    })

    it('should estimate cost impact', () => {
      const changes = trackChanges(originalTasks, modifiedTasks)

      changes.forEach(change => {
        const impact = calculateChangeImpact(change, originalTasks, [], {
          calculateCost: true,
          hourlyRates: { 'user-1': 50, 'user-2': 60 },
        })
        expect(impact).toHaveProperty('costImpact')
      })
    })
  })

  describe('Edge Cases', () => {
    it('should handle null values', () => {
      const tasksWithNull = [
        { id: 'task-1', name: null, start: '2024-01-01T08:00:00Z' },
      ]

      const modifiedWithNull = [
        { id: 'task-1', name: 'Task 1', start: '2024-01-01T08:00:00Z' },
      ]

      const changes = trackChanges(tasksWithNull, modifiedWithNull)

      expect(changes).toBeDefined()
      expect(changes.length).toBeGreaterThan(0)
    })

    it('should handle undefined values', () => {
      const tasksWithUndefined = [
        { id: 'task-1', name: undefined, start: '2024-01-01T08:00:00Z' },
      ]

      const modifiedWithUndefined = [
        { id: 'task-1', name: 'Task 1', start: '2024-01-01T08:00:00Z' },
      ]

      const changes = trackChanges(tasksWithUndefined, modifiedWithUndefined)

      expect(changes).toBeDefined()
    })

    it('should handle missing fields', () => {
      const tasksMissing = [
        { id: 'task-1', start: '2024-01-01T08:00:00Z' },
      ]

      const modifiedMissing = [
        { id: 'task-1', name: 'Task 1', start: '2024-01-01T08:00:00Z' },
      ]

      const changes = trackChanges(tasksMissing, modifiedMissing)

      expect(changes).toBeDefined()
    })

    it('should handle large datasets', () => {
      const largeOriginal = Array.from({ length: 1000 }, (_, i) => ({
        id: `task-${i}`,
        name: `Task ${i}`,
        start: new Date(Date.now() + i * 86400000).toISOString(),
      }))

      const largeModified = largeOriginal.map((task, i) =>
        i % 2 === 0 ? { ...task, name: `Modified ${i}` } : task
      )

      const start = performance.now()
      const changes = trackChanges(largeOriginal, largeModified)
      const end = performance.now()

      expect(changes.length).toBe(500)
      expect(end - start).toBeLessThan(1000) // Should be fast
    })
  })

  describe('Change History', () => {
    it('should maintain change history', () => {
      const changes1 = trackChanges(originalTasks, modifiedTasks)
      const changes2 = trackChanges(modifiedTasks, originalTasks)

      const history = [...changes1, ...changes2]

      expect(history.length).toBe(5) // 2 updates + 1 create + 2 deletes
    })

    it('should sort history by timestamp', () => {
      const changes1 = trackChanges(originalTasks, modifiedTasks)
      const changes2 = trackChanges(modifiedTasks, originalTasks)

      const history = [...changes1, ...changes2].sort((a, b) =>
        a.timestamp - b.timestamp
      )

      for (let i = 1; i < history.length; i++) {
        expect(history[i].timestamp).toBeGreaterThanOrEqual(history[i - 1].timestamp)
      }
    })

    it('should group history by entity', () => {
      const changes1 = trackChanges(originalTasks, modifiedTasks)
      const changes2 = trackChanges(modifiedTasks, originalTasks)

      const history = [...changes1, ...changes2]

      const grouped = history.reduce((acc, change) => {
        if (!acc[change.entityId]) {
          acc[change.entityId] = []
        }
        acc[change.entityId].push(change)
        return acc
      }, {})

      expect(Object.keys(grouped)).toContain('task-1')
      expect(Object.keys(grouped)).toContain('task-2')
      expect(Object.keys(grouped)).toContain('task-3')
    })
  })
})
