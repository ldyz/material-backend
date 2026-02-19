/**
 * Undo/Redo Store
 *
 * Manages command history for the Gantt chart editor with the Command pattern.
 * Supports undo/redo operations for task CRUD operations with keyboard shortcuts.
 *
 * Features:
 * - Command pattern implementation
 * - Command stack with max 50 items
 * - Task CRUD operations support
 * - Keyboard shortcuts (Ctrl+Z, Ctrl+Y, Cmd+Z, Cmd+Y)
 * - History indicators
 *
 * @module UndoRedoStore
 * @date 2025-02-18
 */

import { defineStore } from 'pinia'
import { ElMessage } from 'element-plus'
import eventBus, { GanttEvents } from '@/utils/eventBus'

/**
 * Command interface for undoable operations
 * @interface Command
 */
export class Command {
  /**
   * Execute the command
   * @abstract
   */
  execute() {
    throw new Error('execute() must be implemented')
  }

  /**
   * Undo the command
   * @abstract
   */
  undo() {
    throw new Error('undo() must be implemented')
  }

  /**
   * Get command description for display
   * @returns {string} Command description
   */
  getDescription() {
    return 'Unknown command'
  }

  /**
   * Check if command can be undone
   * @returns {boolean}
   */
  canUndo() {
    return true
  }
}

/**
 * Task Create Command
 */
export class CreateTaskCommand extends Command {
  /**
   * @param {Function} api - API function to create task
   * @param {Object} taskData - Task data
   * @param {Object} store - Gantt store reference
   */
  constructor(api, taskData, store) {
    super()
    this.api = api
    this.taskData = taskData
    this.store = store
    this.createdTaskId = null
  }

  async execute() {
    try {
      const response = await this.api(this.taskData)
      if (response.success) {
        this.createdTaskId = response.data.id
        ElMessage.success('任务已创建')
        eventBus.emit(GanttEvents.TASK_CREATED, { task: response.data })
        return response.data
      }
      throw new Error('创建任务失败')
    } catch (error) {
      console.error('CreateTaskCommand execute error:', error)
      ElMessage.error('创建任务失败')
      throw error
    }
  }

  async undo() {
    if (!this.createdTaskId) {
      console.warn('No task ID to undo')
      return
    }

    try {
      // Delete the task using the progress API
      const { progressApi } = await import('@/api')
      await progressApi.deleteTask(this.createdTaskId)

      ElMessage.info('已撤销创建任务')
      eventBus.emit(GanttEvents.TASK_DELETED, { taskId: this.createdTaskId })
    } catch (error) {
      console.error('CreateTaskCommand undo error:', error)
      ElMessage.error('撤销创建任务失败')
      throw error
    }
  }

  getDescription() {
    return `创建任务: ${this.taskData.task_name || this.taskData.name || '未命名'}`
  }
}

/**
 * Task Update Command
 */
export class UpdateTaskCommand extends Command {
  /**
   * @param {number} taskId - Task ID
   * @param {Object} updates - Updates to apply
   * @param {Object} originalData - Original task data before update
   * @param {Object} store - Gantt store reference
   */
  constructor(taskId, updates, originalData, store) {
    super()
    this.taskId = taskId
    this.updates = updates
    this.originalData = { ...originalData }
    this.store = store
  }

  async execute() {
    try {
      const { progressApi } = await import('@/api')
      const response = await progressApi.update(this.taskId, this.updates)

      if (response.success) {
        ElMessage.success('任务已更新')
        eventBus.emit(GanttEvents.TASK_UPDATED, {
          taskId: this.taskId,
          updates: this.updates
        })
        return response.data
      }
      throw new Error('更新任务失败')
    } catch (error) {
      console.error('UpdateTaskCommand execute error:', error)
      ElMessage.error('更新任务失败')
      throw error
    }
  }

  async undo() {
    try {
      const { progressApi } = await import('@/api')
      await progressApi.update(this.taskId, this.originalData)

      ElMessage.info('已撤销任务更新')
      eventBus.emit(GanttEvents.TASK_UPDATED, {
        taskId: this.taskId,
        updates: this.originalData
      })
    } catch (error) {
      console.error('UpdateTaskCommand undo error:', error)
      ElMessage.error('撤销更新失败')
      throw error
    }
  }

  getDescription() {
    return `更新任务: ${this.originalData.task_name || this.originalData.name || '未命名'}`
  }
}

/**
 * Task Delete Command
 */
export class DeleteTaskCommand extends Command {
  /**
   * @param {number} taskId - Task ID
   * @param {Object} taskData - Complete task data before deletion
   * @param {Object} store - Gantt store reference
   */
  constructor(taskId, taskData, store) {
    super()
    this.taskId = taskId
    this.taskData = { ...taskData }
    this.store = store
  }

  async execute() {
    try {
      const { progressApi } = await import('@/api')
      const response = await progressApi.deleteTask(this.taskId)

      if (response.success) {
        ElMessage.success('任务已删除')
        eventBus.emit(GanttEvents.TASK_DELETED, { taskId: this.taskId })
        return response
      }
      throw new Error('删除任务失败')
    } catch (error) {
      console.error('DeleteTaskCommand execute error:', error)
      ElMessage.error('删除任务失败')
      throw error
    }
  }

  async undo() {
    try {
      // Recreate the task with original data
      const { progressApi } = await import('@/api')
      const response = await progressApi.create({
        project_id: this.taskData.project_id,
        ...this.taskData
      })

      if (response.success) {
        ElMessage.info('已撤销删除任务')
        eventBus.emit(GanttEvents.TASK_CREATED, { task: response.data })
        return response.data
      }
      throw new Error('撤销删除失败')
    } catch (error) {
      console.error('DeleteTaskCommand undo error:', error)
      ElMessage.error('撤销删除失败')
      throw error
    }
  }

  getDescription() {
    return `删除任务: ${this.taskData.task_name || this.taskData.name || '未命名'}`
  }
}

/**
 * Batch Update Command (for multiple task operations)
 */
export class BatchUpdateCommand extends Command {
  /**
   * @param {Array} updates - Array of {taskId, updates, originalData}
   * @param {Object} store - Gantt store reference
   */
  constructor(updates, store) {
    super()
    this.updates = updates
    this.store = store
  }

  async execute() {
    try {
      const { progressApi } = await import('@/api')
      const promises = this.updates.map(({ taskId, updates }) =>
        progressApi.update(taskId, updates)
      )

      await Promise.all(promises)

      ElMessage.success(`已更新 ${this.updates.length} 个任务`)

      this.updates.forEach(({ taskId, updates }) => {
        eventBus.emit(GanttEvents.TASK_UPDATED, { taskId, updates })
      })

      return true
    } catch (error) {
      console.error('BatchUpdateCommand execute error:', error)
      ElMessage.error('批量更新失败')
      throw error
    }
  }

  async undo() {
    try {
      const { progressApi } = await import('@/api')
      const promises = this.updates.map(({ taskId, originalData }) =>
        progressApi.update(taskId, originalData)
      )

      await Promise.all(promises)

      ElMessage.info(`已撤销更新 ${this.updates.length} 个任务`)

      this.updates.forEach(({ taskId, originalData }) => {
        eventBus.emit(GanttEvents.TASK_UPDATED, { taskId, updates: originalData })
      })

      return true
    } catch (error) {
      console.error('BatchUpdateCommand undo error:', error)
      ElMessage.error('撤销批量更新失败')
      throw error
    }
  }

  getDescription() {
    return `批量更新 ${this.updates.length} 个任务`
  }
}

/**
 * Macro Command (for grouping multiple commands)
 */
export class MacroCommand extends Command {
  /**
   * @param {Array<Command>} commands - Array of commands to execute
   */
  constructor(commands = []) {
    super()
    this.commands = commands
  }

  addCommand(command) {
    this.commands.push(command)
  }

  async execute() {
    const results = []
    for (const command of this.commands) {
      try {
        const result = await command.execute()
        results.push(result)
      } catch (error) {
        console.error('MacroCommand execute error:', error)
        // Rollback: undo all previously executed commands
        await this._rollback(results.length - 1)
        throw error
      }
    }
    return results
  }

  async undo() {
    // Undo in reverse order
    for (let i = this.commands.length - 1; i >= 0; i--) {
      try {
        await this.commands[i].undo()
      } catch (error) {
        console.error(`MacroCommand undo error at index ${i}:`, error)
      }
    }
  }

  async _rollback(executedCount) {
    // Undo executed commands in reverse order
    for (let i = executedCount - 1; i >= 0; i--) {
      try {
        await this.commands[i].undo()
      } catch (error) {
        console.error(`MacroCommand rollback error at index ${i}:`, error)
      }
    }
  }

  getDescription() {
    if (this.commands.length === 0) return '空操作'
    if (this.commands.length === 1) return this.commands[0].getDescription()
    return `${this.commands.length} 个操作`
  }

  canUndo() {
    return this.commands.length > 0 && this.commands.every(cmd => cmd.canUndo())
  }
}

/**
 * Undo/Redo State Management Store
 */
export const useUndoRedoStore = defineStore('undoRedo', {
  state: () => ({
    // Command stacks
    undoStack: [],
    redoStack: [],

    // Maximum stack size
    maxStackSize: 50,

    // Current operation state
    isUndoing: false,
    isRedoing: false,

    // History statistics
    historyCount: 0,
    totalOperations: 0
  }),

  getters: {
    /**
     * Check if undo is available
     */
    canUndo: (state) => {
      return state.undoStack.length > 0 && !state.isUndoing
    },

    /**
     * Check if redo is available
     */
    canRedo: (state) => {
      return state.redoStack.length > 0 && !state.isRedoing
    },

    /**
     * Get current undo stack size
     */
    undoStackSize: (state) => {
      return state.undoStack.length
    },

    /**
     * Get current redo stack size
     */
    redoStackSize: (state) => {
      return state.redoStack.length
    },

    /**
     * Get last command description
     */
    lastCommandDescription: (state) => {
      const lastCommand = state.undoStack[state.undoStack.length - 1]
      return lastCommand ? lastCommand.getDescription() : null
    },

    /**
     * Get next command description
     */
    nextCommandDescription: (state) => {
      const nextCommand = state.redoStack[state.redoStack.length - 1]
      return nextCommand ? nextCommand.getDescription() : null
    }
  },

  actions: {
    /**
     * Execute a command and add it to the undo stack
     * @param {Command} command - Command to execute
     * @returns {Promise<any>} Result of command execution
     */
    async executeCommand(command) {
      if (!(command instanceof Command)) {
        console.error('executeCommand requires a Command instance')
        return null
      }

      try {
        // Execute the command
        const result = await command.execute()

        // Add to undo stack if it can be undone
        if (command.canUndo()) {
          this.undoStack.push(command)

          // Maintain max stack size
          if (this.undoStack.length > this.maxStackSize) {
            this.undoStack.shift()
          }

          // Clear redo stack on new command
          this.redoStack = []

          // Update statistics
          this.historyCount++
          this.totalOperations++

          // Emit event
          eventBus.emit(GanttEvents.DATA_CHANGED, {
            hasUnsavedChanges: true,
            commandDescription: command.getDescription()
          })
        }

        return result
      } catch (error) {
        console.error('executeCommand error:', error)
        throw error
      }
    },

    /**
     * Undo the last command
     * @returns {Promise<any>} Result of undo operation
     */
    async undo() {
      if (!this.canUndo) {
        console.warn('Nothing to undo')
        return null
      }

      this.isUndoing = true

      try {
        const command = this.undoStack.pop()
        const result = await command.undo()

        // Add to redo stack
        this.redoStack.push(command)

        // Maintain redo stack size
        if (this.redoStack.length > this.maxStackSize) {
          this.redoStack.shift()
        }

        ElMessage.info(`已撤销: ${command.getDescription()}`)

        // Emit event
        eventBus.emit('gantt:undo', {
          command: command.getDescription(),
          canRedo: true
        })

        return result
      } catch (error) {
        console.error('Undo error:', error)
        ElMessage.error('撤销操作失败')

        // Put command back in undo stack on error
        if (command) {
          this.undoStack.push(command)
        }

        throw error
      } finally {
        this.isUndoing = false
      }
    },

    /**
     * Redo the last undone command
     * @returns {Promise<any>} Result of redo operation
     */
    async redo() {
      if (!this.canRedo) {
        console.warn('Nothing to redo')
        return null
      }

      this.isRedoing = true

      try {
        const command = this.redoStack.pop()
        const result = await command.execute()

        // Add back to undo stack
        this.undoStack.push(command)

        this.totalOperations++

        ElMessage.info(`已重做: ${command.getDescription()}`)

        // Emit event
        eventBus.emit('gantt:redo', {
          command: command.getDescription(),
          canUndo: true
        })

        return result
      } catch (error) {
        console.error('Redo error:', error)
        ElMessage.error('重做操作失败')

        // Put command back in redo stack on error
        if (command) {
          this.redoStack.push(command)
        }

        throw error
      } finally {
        this.isRedoing = false
      }
    },

    /**
     * Clear all history
     */
    clearHistory() {
      this.undoStack = []
      this.redoStack = []
      this.historyCount = 0
      this.totalOperations = 0

      eventBus.emit('gantt:history-cleared')
    },

    /**
     * Get history snapshot for display
     * @returns {Array} Array of command descriptions
     */
    getHistorySnapshot() {
      return this.undoStack.map(cmd => ({
        description: cmd.getDescription(),
        canUndo: cmd.canUndo()
      }))
    },

    /**
     * Get redo snapshot for display
     * @returns {Array} Array of command descriptions
     */
    getRedoSnapshot() {
      return [...this.redoStack].reverse().map(cmd => ({
        description: cmd.getDescription(),
        canUndo: cmd.canUndo()
      }))
    },

    /**
     * Trim history to specific size
     * @param {number} size - Target size
     */
    trimHistory(size) {
      while (this.undoStack.length > size) {
        this.undoStack.shift()
      }
      while (this.redoStack.length > size) {
        this.redoStack.shift()
      }
    },

    /**
     * Set maximum stack size
     * @param {number} size - New maximum size
     */
    setMaxStackSize(size) {
      this.maxStackSize = Math.max(1, Math.min(100, size))
      this.trimHistory(this.maxStackSize)
    }
  }
})

export default useUndoRedoStore
