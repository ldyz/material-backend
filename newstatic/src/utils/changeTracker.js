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

export default changeTracker
