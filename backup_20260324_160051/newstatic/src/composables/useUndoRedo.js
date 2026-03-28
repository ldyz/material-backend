/**
 * Undo/Redo Composable
 *
 * Provides a convenient interface for undo/redo operations in Vue components.
 * Exposes command factory functions and keyboard shortcuts support.
 *
 * @module useUndoRedo
 * @date 2025-02-18
 */

import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useUndoRedoStore } from '@/stores/undoRedoStore'
import {
  CreateTaskCommand,
  UpdateTaskCommand,
  DeleteTaskCommand,
  BatchUpdateCommand,
  MacroCommand
} from '@/stores/undoRedoStore'
import { ElMessage } from 'element-plus'

/**
 * Undo/Redo Composable
 * @param {Object} options - Configuration options
 * @param {Object} options.store - Gantt store reference (optional)
 * @param {boolean} options.enableKeyboard - Enable keyboard shortcuts (default: true)
 * @param {Function} options.onUndo - Callback when undo is performed
 * @param {Function} options.onRedo - Callback when redo is performed
 * @returns {Object} Undo/Redo interface
 */
export function useUndoRedo(options = {}) {
  const {
    store = null,
    enableKeyboard = true,
    onUndo = null,
    onRedo = null
  } = options

  // Store reference
  const undoRedoStore = useUndoRedoStore()

  // Local state
  const isProcessing = ref(false)

  /**
   * Check if currently processing undo/redo
   */
  const isUndoing = computed(() => undoRedoStore.isUndoing || isProcessing.value)
  const isRedoing = computed(() => undoRedoStore.isRedoing || isProcessing.value)

  /**
   * Undo availability
   */
  const canUndo = computed(() => undoRedoStore.canUndo)

  /**
   * Redo availability
   */
  const canRedo = computed(() => undoRedoStore.canRedo)

  /**
   * Last command description
   */
  const lastCommandDescription = computed(() => undoRedoStore.lastCommandDescription)

  /**
   * Next command description
   */
  const nextCommandDescription = computed(() => undoRedoStore.nextCommandDescription)

  /**
   * History indicators
   */
  const historyCount = computed(() => undoRedoStore.undoStack.length)
  const redoCount = computed(() => undoRedoStore.redoStack.length)

  /**
   * Undo operation
   */
  const undo = async () => {
    if (!canUndo.value || isUndoing.value) {
      return null
    }

    isProcessing.value = true

    try {
      const result = await undoRedoStore.undo()

      if (onUndo) {
        onUndo(result)
      }

      return result
    } catch (error) {
      console.error('Undo failed:', error)
      throw error
    } finally {
      isProcessing.value = false
    }
  }

  /**
   * Redo operation
   */
  const redo = async () => {
    if (!canRedo.value || isRedoing.value) {
      return null
    }

    isProcessing.value = true

    try {
      const result = await undoRedoStore.redo()

      if (onRedo) {
        onRedo(result)
      }

      return result
    } catch (error) {
      console.error('Redo failed:', error)
      throw error
    } finally {
      isProcessing.value = false
    }
  }

  /**
   * Execute a command
   * @param {import('@/stores/undoRedoStore').Command} command - Command to execute
   * @returns {Promise<any>} Result of command execution
   */
  const executeCommand = async (command) => {
    if (isUndoing.value || isRedoing.value) {
      console.warn('Cannot execute command during undo/redo')
      return null
    }

    return await undoRedoStore.executeCommand(command)
  }

  // ==================== Command Factory Functions ====================

  /**
   * Create a task creation command
   * @param {Object} taskData - Task data to create
   * @returns {CreateTaskCommand} Command instance
   */
  const createCreateTaskCommand = (taskData) => {
    return new CreateTaskCommand(
      async (data) => {
        const { progressApi } = await import('@/api')
        return await progressApi.create(data)
      },
      taskData,
      store
    )
  }

  /**
   * Create a task update command
   * @param {number} taskId - Task ID
   * @param {Object} updates - Updates to apply
   * @param {Object} originalData - Original task data
   * @returns {UpdateTaskCommand} Command instance
   */
  const createUpdateTaskCommand = (taskId, updates, originalData) => {
    return new UpdateTaskCommand(taskId, updates, originalData, store)
  }

  /**
   * Create a task delete command
   * @param {number} taskId - Task ID
   * @param {Object} taskData - Complete task data before deletion
   * @returns {DeleteTaskCommand} Command instance
   */
  const createDeleteTaskCommand = (taskId, taskData) => {
    return new DeleteTaskCommand(taskId, taskData, store)
  }

  /**
   * Create a batch update command
   * @param {Array} updates - Array of {taskId, updates, originalData}
   * @returns {BatchUpdateCommand} Command instance
   */
  const createBatchUpdateCommand = (updates) => {
    return new BatchUpdateCommand(updates, store)
  }

  /**
   * Create a macro command for grouping multiple commands
   * @param {Array<Command>} commands - Array of commands
   * @returns {MacroCommand} Command instance
   */
  const createMacroCommand = (commands) => {
    return new MacroCommand(commands)
  }

  /**
   * Execute task creation
   * @param {Object} taskData - Task data to create
   * @returns {Promise<any>} Result of creation
   */
  const executeCreateTask = async (taskData) => {
    const command = createCreateTaskCommand(taskData)
    return await executeCommand(command)
  }

  /**
   * Execute task update
   * @param {number} taskId - Task ID
   * @param {Object} updates - Updates to apply
   * @param {Object} originalData - Original task data
   * @returns {Promise<any>} Result of update
   */
  const executeUpdateTask = async (taskId, updates, originalData) => {
    const command = createUpdateTaskCommand(taskId, updates, originalData)
    return await executeCommand(command)
  }

  /**
   * Execute task deletion
   * @param {number} taskId - Task ID
   * @param {Object} taskData - Complete task data
   * @returns {Promise<any>} Result of deletion
   */
  const executeDeleteTask = async (taskId, taskData) => {
    const command = createDeleteTaskCommand(taskId, taskData)
    return await executeCommand(command)
  }

  /**
   * Execute batch update
   * @param {Array} updates - Array of {taskId, updates, originalData}
   * @returns {Promise<any>} Result of batch update
   */
  const executeBatchUpdate = async (updates) => {
    const command = createBatchUpdateCommand(updates)
    return await executeCommand(command)
  }

  /**
   * Execute macro command
   * @param {Array<Command>} commands - Array of commands
   * @returns {Promise<any>} Result of macro command
   */
  const executeMacro = async (commands) => {
    const command = createMacroCommand(commands)
    return await executeCommand(command)
  }

  /**
   * Clear all history
   */
  const clearHistory = () => {
    undoRedoStore.clearHistory()
    ElMessage.info('历史记录已清空')
  }

  /**
   * Get history snapshot
   * @returns {Array} Array of command descriptions
   */
  const getHistorySnapshot = () => {
    return undoRedoStore.getHistorySnapshot()
  }

  /**
   * Get redo snapshot
   * @returns {Array} Array of command descriptions
   */
  const getRedoSnapshot = () => {
    return undoRedoStore.getRedoSnapshot()
  }

  /**
   * Set maximum stack size
   * @param {number} size - New maximum size
   */
  const setMaxStackSize = (size) => {
    undoRedoStore.setMaxStackSize(size)
  }

  // ==================== Keyboard Shortcuts ====================

  /**
   * Handle keyboard events
   */
  const handleKeyDown = (event) => {
    // Ctrl+Z or Cmd+Z - Undo
    if ((event.ctrlKey || event.metaKey) && event.key === 'z' && !event.shiftKey) {
      event.preventDefault()
      if (canUndo.value) {
        undo()
      }
      return
    }

    // Ctrl+Y or Ctrl+Shift+Z or Cmd+Shift+Z - Redo
    if ((event.ctrlKey || event.metaKey) &&
        (event.key === 'y' || (event.key === 'z' && event.shiftKey))) {
      event.preventDefault()
      if (canRedo.value) {
        redo()
      }
      return
    }
  }

  // Setup keyboard shortcuts
  if (enableKeyboard) {
    onMounted(() => {
      window.addEventListener('keydown', handleKeyDown)
    })

    onUnmounted(() => {
      window.removeEventListener('keydown', handleKeyDown)
    })
  }

  // Return public interface
  return {
    // State
    canUndo,
    canRedo,
    isUndoing,
    isRedoing,
    lastCommandDescription,
    nextCommandDescription,
    historyCount,
    redoCount,

    // Actions
    undo,
    redo,
    executeCommand,
    clearHistory,
    getHistorySnapshot,
    getRedoSnapshot,
    setMaxStackSize,

    // Command factories
    createCreateTaskCommand,
    createUpdateTaskCommand,
    createDeleteTaskCommand,
    createBatchUpdateCommand,
    createMacroCommand,

    // Convenience methods
    executeCreateTask,
    executeUpdateTask,
    executeDeleteTask,
    executeBatchUpdate,
    executeMacro
  }
}

export default useUndoRedo
