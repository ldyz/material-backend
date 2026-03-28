/**
 * Change Tracker Utility
 *
 * Tracks all changes to tasks, dependencies, and resources.
 * Generates diffs between versions and exports change logs.
 *
 * @module changeTracker
 * @author Material Management System
 */

/**
 * @typedef {Object} ChangeEntry
 * @property {string} id - Unique change ID
 * @property {number} projectId - Project ID
 * @property {number} userId - User ID who made the change
 * @property {string} entityType - 'task', 'dependency', 'resource'
 * @property {number} entityId - Entity ID
 * @property {string} actionType - 'create', 'update', 'delete'
 * @property {object} changes - Before/After data
 * @property {string} description - Human-readable description
 * @property {Date} timestamp - Change timestamp
 */

/**
 * @typedef {Object} DiffResult
 * @property {Array<string>} added - Added fields
 * @property {Array<string>} removed - Removed fields
 * @property {Array<string>} modified - Modified fields
 * @property {object} before - Before state
 * @property {object} after - After state
 */

class ChangeTracker {
  constructor() {
    this.changes = []
    this.maxChanges = 1000
  }

  /**
   * Track a create operation
   *
   * @param {string} entityType - Entity type
   * @param {number} entityId - Entity ID
   * @param {object} data - Created data
   * @param {number} userId - User ID
   * @param {number} projectId - Project ID
   * @returns {ChangeEntry}
   *
   * @example
   * tracker.trackCreate('task', 123, { name: 'New Task' }, 1, 456)
   */
  trackCreate(entityType, entityId, data, userId, projectId) {
    const change = {
      id: this._generateId(),
      projectId,
      userId,
      entityType,
      entityId,
      actionType: 'create',
      changes: {
        before: {},
        after: this._sanitizeData(data)
      },
      description: this._generateDescription('create', entityType, entityId),
      timestamp: new Date()
    }

    this._addChange(change)
    return change
  }

  /**
   * Track an update operation
   *
   * @param {string} entityType - Entity type
   * @param {number} entityId - Entity ID
   * @param {object} beforeData - Data before change
   * @param {object} afterData - Data after change
   * @param {number} userId - User ID
   * @param {number} projectId - Project ID
   * @returns {ChangeEntry}
   *
   * @example
   * tracker.trackUpdate('task', 123, oldData, newData, 1, 456)
   */
  trackUpdate(entityType, entityId, beforeData, afterData, userId, projectId) {
    const before = this._sanitizeData(beforeData)
    const after = this._sanitizeData(afterData)

    // Only track if there are actual changes
    const diff = this._diffObjects(before, after)
    if (diff.modified.length === 0 && diff.added.length === 0 && diff.removed.length === 0) {
      return null
    }

    const change = {
      id: this._generateId(),
      projectId,
      userId,
      entityType,
      entityId,
      actionType: 'update',
      changes: {
        before,
        after
      },
      description: this._generateDescription('update', entityType, entityId, diff),
      timestamp: new Date()
    }

    this._addChange(change)
    return change
  }

  /**
   * Track a delete operation
   *
   * @param {string} entityType - Entity type
   * @param {number} entityId - Entity ID
   * @param {object} data - Data before deletion
   * @param {number} userId - User ID
   * @param {number} projectId - Project ID
   * @returns {ChangeEntry}
   *
   * @example
   * tracker.trackDelete('task', 123, taskData, 1, 456)
   */
  trackDelete(entityType, entityId, data, userId, projectId) {
    const change = {
      id: this._generateId(),
      projectId,
      userId,
      entityType,
      entityId,
      actionType: 'delete',
      changes: {
        before: this._sanitizeData(data),
        after: {}
      },
      description: this._generateDescription('delete', entityType, entityId),
      timestamp: new Date()
    }

    this._addChange(change)
    return change
  }

  /**
   * Track changes between two task arrays
   *
   * @param {Array} oldTasks - Previous task array
   * @param {Array} newTasks - New task array
   * @param {number} userId - User ID who made the changes
   * @param {number} projectId - Project ID (optional)
   * @returns {Array<ChangeEntry>} Array of change entries
   *
   * @example
   * const changes = tracker.trackChanges(oldTasks, newTasks, 1, 456)
   */
  trackChanges(oldTasks, newTasks, userId = null, projectId = null) {
    const changes = []
    const oldMap = new Map()
    const newMap = new Map()

    // Create maps for O(1) lookup
    oldTasks?.forEach(task => oldMap.set(task.id || task._id, task))
    newTasks?.forEach(task => newMap.set(task.id || task._id, task))

    // Detect new and modified tasks
    newMap.forEach((newTask, id) => {
      const oldTask = oldMap.get(id)
      if (!oldTask) {
        // New task
        const change = this.trackCreate('task', id, newTask, userId, projectId)
        if (change) changes.push(change)
      } else {
        // Check for modifications
        const diff = this._diffObjects(oldTask, newTask)
        if (diff.modified.length > 0 || diff.added.length > 0 || diff.removed.length > 0) {
          const change = this.trackUpdate('task', id, oldTask, newTask, userId, projectId)
          if (change) changes.push(change)
        }
      }
    })

    // Detect deleted tasks
    oldMap.forEach((oldTask, id) => {
      if (!newMap.has(id)) {
        const change = this.trackDelete('task', id, oldTask, userId, projectId)
        if (change) changes.push(change)
      }
    })

    return changes
  }

  /**
   * Generate diff between two objects
   *
   * @param {object} before - Before state
   * @param {object} after - After state
   * @returns {DiffResult}
   *
   * @example
   * const diff = tracker.diff(oldTask, newTask)
   */
  diff(before, after) {
    return this._diffObjects(before, after)
  }

  /**
   * Get all changes for a project
   *
   * @param {number} projectId - Project ID
   * @param {object} options - Filter options
   * @param {string} options.entityType - Filter by entity type
   * @param {string} options.actionType - Filter by action type
   * @param {Date} options.startDate - Filter by start date
   * @param {Date} options.endDate - Filter by end date
   * @returns {Array<ChangeEntry>}
   */
  getChanges(projectId, options = {}) {
    let changes = this.changes.filter(c => c.projectId === projectId)

    if (options.entityType) {
      changes = changes.filter(c => c.entityType === options.entityType)
    }

    if (options.actionType) {
      changes = changes.filter(c => c.actionType === options.actionType)
    }

    if (options.startDate) {
      changes = changes.filter(c => c.timestamp >= options.startDate)
    }

    if (options.endDate) {
      changes = changes.filter(c => c.timestamp <= options.endDate)
    }

    return changes.sort((a, b) => b.timestamp - a.timestamp)
  }

  /**
   * Get changes for an entity
   *
   * @param {string} entityType - Entity type
   * @param {number} entityId - Entity ID
   * @returns {Array<ChangeEntry>}
   */
  getEntityChanges(entityType, entityId) {
    return this.changes
      .filter(c => c.entityType === entityType && c.entityId === entityId)
      .sort((a, b) => b.timestamp - a.timestamp)
  }

  /**
   * Export changes to various formats
   *
   * @param {Array<ChangeEntry>} changes - Changes to export
   * @param {string} format - Export format: 'json', 'csv', 'xlsx'
   * @returns {string|Blob} Exported data
   */
  exportChanges(changes, format = 'json') {
    switch (format) {
      case 'json':
        return this._exportJSON(changes)
      case 'csv':
        return this._exportCSV(changes)
      case 'xlsx':
        return this._exportXLSX(changes)
      default:
        throw new Error(`Unsupported export format: ${format}`)
    }
  }

  /**
   * Clear old changes
   *
   * @param {number} keepDays - Number of days to keep
   */
  clearOldChanges(keepDays = 90) {
    const cutoff = new Date()
    cutoff.setDate(cutoff.getDate() - keepDays)

    this.changes = this.changes.filter(c => c.timestamp > cutoff)
  }

  /**
   * Get statistics about changes
   *
   * @param {number} projectId - Project ID
   * @returns {object} Statistics
   */
  getStatistics(projectId) {
    const changes = this.getChanges(projectId)

    const stats = {
      total: changes.length,
      byActionType: {},
      byEntityType: {},
      byUser: {},
      dateRange: {
        earliest: null,
        latest: null
      }
    }

    changes.forEach(change => {
      // By action type
      stats.byActionType[change.actionType] = (stats.byActionType[change.actionType] || 0) + 1

      // By entity type
      stats.byEntityType[change.entityType] = (stats.byEntityType[change.entityType] || 0) + 1

      // By user
      stats.byUser[change.userId] = (stats.byUser[change.userId] || 0) + 1

      // Date range
      if (!stats.dateRange.earliest || change.timestamp < stats.dateRange.earliest) {
        stats.dateRange.earliest = change.timestamp
      }
      if (!stats.dateRange.latest || change.timestamp > stats.dateRange.latest) {
        stats.dateRange.latest = change.timestamp
      }
    })

    return stats
  }

  // ==================== Private Methods ====================

  /**
   * Add change to history
   * @private
   */
  _addChange(change) {
    this.changes.unshift(change)

    // Limit history size
    if (this.changes.length > this.maxChanges) {
      this.changes = this.changes.slice(0, this.maxChanges)
    }
  }

  /**
   * Generate unique ID
   * @private
   */
  _generateId() {
    return `change-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
  }

  /**
   * Sanitize data for storage
   * @private
   */
  _sanitizeData(data) {
    if (!data) return {}

    const sanitized = { ...data }

    // Remove sensitive fields
    delete sanitized.password
    delete sanitized.token
    delete sanitized.secret

    // Convert dates to ISO strings
    Object.keys(sanitized).forEach(key => {
      if (sanitized[key] instanceof Date) {
        sanitized[key] = sanitized[key].toISOString()
      }
    })

    return sanitized
  }

  /**
   * Diff two objects
   * @private
   */
  _diffObjects(before, after) {
    const beforeKeys = new Set(Object.keys(before))
    const afterKeys = new Set(Object.keys(after))

    const added = [...afterKeys].filter(k => !beforeKeys.has(k))
    const removed = [...beforeKeys].filter(k => !afterKeys.has(k))
    const common = [...beforeKeys].filter(k => afterKeys.has(k))

    const modified = common.filter(key => {
      return JSON.stringify(before[key]) !== JSON.stringify(after[key])
    })

    return { added, removed, modified, before, after }
  }

  /**
   * Generate human-readable description
   * @private
   */
  _generateDescription(action, entityType, entityId, diff = null) {
    const actionText = {
      create: 'created',
      update: 'updated',
      delete: 'deleted'
    }[action] || action

    let description = `${entityType} #${entityId} ${actionText}`

    if (action === 'update' && diff) {
      const changedFields = [...diff.added, ...diff.removed, ...diff.modified]
      if (changedFields.length > 0) {
        description += ` (${changedFields.length} fields changed)`
      }
    }

    return description
  }

  /**
   * Export as JSON
   * @private
   */
  _exportJSON(changes) {
    return JSON.stringify(changes, null, 2)
  }

  /**
   * Export as CSV
   * @private
   */
  _exportCSV(changes) {
    const headers = ['Date', 'User ID', 'Entity Type', 'Entity ID', 'Action', 'Description', 'Changes']
    const rows = changes.map(c => [
      c.timestamp.toISOString(),
      c.userId,
      c.entityType,
      c.entityId,
      c.actionType,
      c.description,
      JSON.stringify(c.changes)
    ])

    return [
      headers.join(','),
      ...rows.map(row => row.map(cell => `"${cell}"`).join(','))
    ].join('\n')
  }

  /**
   * Export as XLSX
   * @private
   */
  _exportXLSX(changes) {
    // Note: This would require a library like xlsx
    // For now, return JSON as fallback
    console.warn('XLSX export not implemented, falling back to JSON')
    return this._exportJSON(changes)
  }
}

// Create singleton instance
const changeTracker = new ChangeTracker()

// Change Type Enum
export const ChangeType = {
  CREATE: 'CREATE',
  UPDATE: 'UPDATE',
  DELETE: 'DELETE'
}

/**
 * Track changes between two task arrays (simplified version for tests)
 *
 * @param {Array} oldTasks - Previous task array
 * @param {Array} newTasks - New task array
 * @param {number} userId - User ID who made the changes
 * @returns {Array} Array of change entries with simplified format
 */
export function trackChanges(oldTasks, newTasks, userId = null) {
  const changes = []
  const oldMap = new Map()
  const newMap = new Map()

  // Create maps for O(1) lookup
  oldTasks?.forEach(task => oldMap.set(task.id || task._id, task))
  newTasks?.forEach(task => newMap.set(task.id || task._id, task))

  // Detect new and modified tasks
  newMap.forEach((newTask, id) => {
    const oldTask = oldMap.get(id)
    if (!oldTask) {
      // New task
      changes.push({
        entityId: id,
        type: ChangeType.CREATE,
        timestamp: Date.now(),
        changes: newTask
      })
    } else {
      // Check for modifications
      const modifiedFields = []
      const fieldChanges = {}

      // Check all fields in both old and new
      const allKeys = new Set([...Object.keys(oldTask), ...Object.keys(newTask)])
      allKeys.forEach(key => {
        if (key !== 'id' && key !== '_id') {
          const oldVal = oldTask[key]
          const newVal = newTask[key]

          if (JSON.stringify(oldVal) !== JSON.stringify(newVal)) {
            modifiedFields.push(key)
            fieldChanges[key] = {
              from: oldVal,
              to: newVal
            }
          }
        }
      })

      if (modifiedFields.length > 0) {
        changes.push({
          entityId: id,
          type: ChangeType.UPDATE,
          timestamp: Date.now(),
          changes: fieldChanges
        })
      }
    }
  })

  // Detect deleted tasks
  oldMap.forEach((oldTask, id) => {
    if (!newMap.has(id)) {
      changes.push({
        entityId: id,
        type: ChangeType.DELETE,
        timestamp: Date.now(),
        changes: oldTask
      })
    }
  })

  return changes
}

/**
 * Generate diff in various formats
 *
 * @param {Array} changes - Array of changes
 * @param {string} format - Format: 'text', 'html', 'markdown', 'json'
 * @returns {string} Formatted diff
 */
export function generateDiff(changes, format = 'text') {
  if (!changes || changes.length === 0) {
    return format === 'json' ? '[]' : 'No changes'
  }

  switch (format) {
    case 'json':
      return JSON.stringify(changes, null, 2)

    case 'html':
      return changes.map(change => {
        const lines = []
        lines.push(`<div class="change change-${change.type.toLowerCase()}">`)
        lines.push(`  <h4>Entity: ${change.entityId}</h4>`)
        lines.push(`  <p>Type: ${change.type}</p>`)

        if (change.type === ChangeType.UPDATE) {
          lines.push('  <ul>')
          Object.entries(change.changes).forEach(([field, values]) => {
            lines.push(`    <li><strong>${field}</strong>: ${JSON.stringify(values.from)} → ${JSON.stringify(values.to)}</li>`)
          })
          lines.push('  </ul>')
        }

        lines.push('</div>')
        return lines.join('\n')
      }).join('\n')

    case 'markdown':
      return changes.map(change => {
        const lines = []
        lines.push(`## Change: ${change.entityId}`)
        lines.push(`**Type:** ${change.type}`)

        if (change.type === ChangeType.UPDATE) {
          lines.push('\n**Field Changes:**')
          Object.entries(change.changes).forEach(([field, values]) => {
            lines.push(`- **${field}**: \`${JSON.stringify(values.from)}\` → \`${JSON.stringify(values.to)}\``)
          })
        }

        return lines.join('\n')
      }).join('\n\n')

    default: // text
      return changes.map(change => {
        const lines = []
        lines.push(`${change.type}: ${change.entityId}`)

        if (change.type === ChangeType.UPDATE) {
          Object.entries(change.changes).forEach(([field, values]) => {
            lines.push(`  ${field} changed from ${JSON.stringify(values.from)} to ${JSON.stringify(values.to)}`)
          })
        }

        return lines.join('\n')
      }).join('\n\n')
  }
}

/**
 * Compress changes by grouping related changes
 *
 * @param {Array} changes - Array of changes
 * @returns {Array} Compressed changes
 */
export function compressChanges(changes) {
  if (!changes || changes.length === 0) {
    return []
  }

  // Group by entity
  const grouped = new Map()

  changes.forEach(change => {
    if (!grouped.has(change.entityId)) {
      grouped.set(change.entityId, [])
    }
    grouped.get(change.entityId).push(change)
  })

  // For each entity, keep only the latest change of each type
  const compressed = []
  grouped.forEach((entityChanges, entityId) => {
    // Filter out redundant changes
    const filtered = entityChanges.filter((change, index, self) => {
      // Remove if there's a later change of same type
      const laterChanges = entityChanges.slice(index + 1)
      return !laterChanges.some(c => c.type === change.type)
    })

    compressed.push(...filtered)
  })

  return compressed.sort((a, b) => a.timestamp - b.timestamp)
}

/**
 * Calculate the impact of changes
 *
 * @param {Array} changes - Array of changes
 * @param {Array} tasks - All tasks
 * @returns {Object} Impact analysis
 */
export function calculateChangeImpact(changes, tasks) {
  if (!changes || changes.length === 0) {
    return {
      severity: 'none',
      affectedTasks: [],
      dateImpact: false,
      dependencyImpact: false,
      resourceImpact: false,
      costImpact: 0,
      estimatedDelay: 0
    }
  }

  const affectedTasks = new Set()
  let dateImpact = false
  let dependencyImpact = false
  let resourceImpact = false
  let maxDelay = 0

  changes.forEach(change => {
    affectedTasks.add(change.entityId)

    if (change.type === ChangeType.UPDATE) {
      const hasDateChange = Object.keys(change.changes).some(key =>
        key.includes('start') || key.includes('end') || key.includes('date')
      )
      const hasDependencyChange = Object.keys(change.changes).some(key =>
        key.includes('depend') || key.includes('predecessor') || key.includes('successor')
      )
      const hasResourceChange = Object.keys(change.changes).some(key =>
        key.includes('assignee') || key.includes('resource')
      )

      if (hasDateChange) {
        dateImpact = true
        // Estimate delay from date changes
        Object.entries(change.changes).forEach(([field, values]) => {
          if ((field.includes('start') || field.includes('end')) && values.from && values.to) {
            const from = new Date(values.from)
            const to = new Date(values.to)
            const delay = Math.abs((to - from) / (1000 * 60 * 60 * 24))
            maxDelay = Math.max(maxDelay, delay)
          }
        })
      }

      if (hasDependencyChange) dependencyImpact = true
      if (hasResourceChange) resourceImpact = true
    }
  })

  // Calculate severity
  let severity = 'low'
  if (dateImpact && maxDelay > 7) severity = 'high'
  else if ((dateImpact || dependencyImpact) && affectedTasks.size > 3) severity = 'medium'
  else if (affectedTasks.size > 5) severity = 'medium'

  return {
    severity,
    affectedTasks: Array.from(affectedTasks),
    dateImpact,
    dependencyImpact,
    resourceImpact,
    costImpact: maxDelay * 1000, // Rough estimate
    estimatedDelay: maxDelay
  }
}

/**
 * Import changes from external data
 *
 * @param {string} data - Data to import
 * @param {string} format - Format: 'json' or 'csv'
 * @returns {Array} Imported changes
 */
export function importChanges(data, format = 'json') {
  try {
    if (format === 'json') {
      const parsed = JSON.parse(data)
      if (Array.isArray(parsed)) {
        return parsed
      }
      return []
    } else if (format === 'csv') {
      // Simple CSV parsing
      const lines = data.split('\n')
      const headers = lines[0].split(',')
      return lines.slice(1).map(line => {
        const values = line.split(',')
        const change = {}
        headers.forEach((header, index) => {
          change[header.trim()] = values[index]?.trim().replace(/"/g, '')
        })
        return change
      }).filter(c => c.entityId)
    }
    return []
  } catch (error) {
    console.error('Import error:', error)
    return []
  }
}

// Named exports for easier testing
export const exportChanges = (changes, format = 'json', options = {}) => {
  const result = changeTracker.exportChanges(changes, format)

  if (options.includeMetadata) {
    const parsed = format === 'json' ? JSON.parse(result) : { data: result }
    return JSON.stringify({
      metadata: {
        exportedBy: options.exportedBy,
        exportDate: options.exportDate || new Date().toISOString(),
        changeCount: changes.length
      },
      changes: parsed
    }, null, 2)
  }

  return result
}

export const getChanges = () => changeTracker.getChanges()
export const clearChanges = () => changeTracker.clearChanges()
export const getStatistics = () => changeTracker.getStatistics()

export default changeTracker
