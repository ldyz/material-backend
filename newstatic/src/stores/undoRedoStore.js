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
import { ref, computed } from 'vue'

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
 * Macro Command (for grouping multiple commands)
 */
export class MacroCommand extends Command {
  /**
   * @param {Array<Command>} commands - Array of commands to execute
   */
  constructor(commands = []) {
    super()
    this.commands = commands
    this.description = 'Macro command'
  }

  addCommand(command) {
    this.commands.push(command)
  }

  execute() {
    const results = []
    for (const command of this.commands) {
      const result = command.execute()

      // Handle async execute
      if (result && typeof result.then === 'function') {
        return result.then((r) => {
          if (r && r.success === false) {
            // Undo all previously executed commands
            for (let i = results.length - 1; i >= 0; i--) {
              if (this.commands[i].undo) {
                this.commands[i].undo()
              }
            }
            return r
          }
          results.push(r)
          return { success: true, data: results }
        })
      }

      // Handle sync execute
      if (result && result.success === false) {
        // Undo all previously executed commands
        for (let i = results.length - 1; i >= 0; i--) {
          if (this.commands[i].undo) {
            this.commands[i].undo()
          }
        }
        return result
      }
      results.push(result)
    }
    return { success: true, data: results }
  }

  undo() {
    // Undo in reverse order
    for (let i = this.commands.length - 1; i >= 0; i--) {
      if (this.commands[i].undo) {
        this.commands[i].undo()
      }
    }
  }

  getDescription() {
    if (this.commands.length === 0) return 'Empty macro'
    return this.description
  }

  canUndo() {
    return this.commands.length > 0 && this.commands.every(cmd => cmd.canUndo ? cmd.canUndo() : true)
  }
}

/**
 * Undo/Redo State Management Store
 */
export const useUndoRedoStore = defineStore('undoRedo', () => {
  // State
  const commandStack = ref([])
  const redoStack = ref([])
  const maxStackSize = ref(50)
  const currentPosition = ref(-1)

  // Computed
  const canUndo = computed(() => commandStack.value.length > 0)
  const canRedo = computed(() => redoStack.value.length > 0)
  const stackSize = computed(() => commandStack.value.length)

  const lastCommandDescription = computed(() => {
    const lastCommand = commandStack.value[commandStack.value.length - 1]
    if (!lastCommand) return null
    return lastCommand.description || (lastCommand.getDescription ? lastCommand.getDescription() : 'Unknown')
  })

  // Actions
  const execute = (command) => {
    if (!command) {
      return { success: false, error: 'Command is null' }
    }

    try {
      const result = command.execute()

      // Handle promise results
      if (result && typeof result.then === 'function') {
        return result.then((r) => {
          // Only add to stack if execution was successful
          if (r && r.success === false) {
            return r
          }

          commandStack.value.push(command)
          currentPosition.value = commandStack.value.length - 1

          // Maintain max stack size
          if (commandStack.value.length > maxStackSize.value) {
            commandStack.value.shift()
            currentPosition.value--
          }

          // Clear redo stack on new command
          redoStack.value = []

          return r || { success: true }
        }).catch((error) => {
          console.error('Execute error:', error)
          return { success: false, error: error.message }
        })
      }

      // Handle synchronous results
      // Only add to stack if execution was successful
      if (result && result.success === false) {
        return result
      }

      commandStack.value.push(command)
      currentPosition.value = commandStack.value.length - 1

      // Maintain max stack size
      if (commandStack.value.length > maxStackSize.value) {
        commandStack.value.shift()
        currentPosition.value--
      }

      // Clear redo stack on new command
      redoStack.value = []

      return result || { success: true }
    } catch (error) {
      // Handle execution errors gracefully
      console.error('Execute error:', error)
      return { success: false, error: error.message }
    }
  }

  const undo = (callback) => {
    if (!canUndo.value) {
      return null
    }

    const command = commandStack.value.pop()
    if (!command) {
      return null
    }

    currentPosition.value--

    try {
      if (command.undo) {
        const result = command.undo()

        // Handle async undo
        if (result && typeof result.then === 'function') {
          return result.then(() => {
            // Add to redo stack
            redoStack.value.push(command)

            if (callback) {
              callback(command)
            }

            return { success: true }
          }).catch((error) => {
            console.error('Undo error:', error)
            // Put command back on error
            commandStack.value.push(command)
            currentPosition.value++
            throw error
          })
        }
      }

      // Add to redo stack
      redoStack.value.push(command)

      if (callback) {
        callback(command)
      }

      return { success: true }
    } catch (error) {
      console.error('Undo error:', error)
      // Put command back on error
      commandStack.value.push(command)
      currentPosition.value++
      throw error
    }
  }

  const redo = (callback) => {
    if (!canRedo.value) {
      return null
    }

    const command = redoStack.value.pop()
    if (!command) {
      return null
    }

    try {
      const result = command.execute()

      // Handle async execute
      if (result && typeof result.then === 'function') {
        return result.then(() => {
          // Add back to command stack
          commandStack.value.push(command)
          currentPosition.value++

          if (callback) {
            callback(command)
          }

          return { success: true }
        }).catch((error) => {
          console.error('Redo error:', error)
          // Put command back on error
          redoStack.value.push(command)
          throw error
        })
      }

      // Add back to command stack
      commandStack.value.push(command)
      currentPosition.value++

      if (callback) {
        callback(command)
      }

      return { success: true }
    } catch (error) {
      console.error('Redo error:', error)
      // Put command back on error
      redoStack.value.push(command)
      throw error
    }
  }

  const clear = () => {
    commandStack.value = []
    redoStack.value = []
    currentPosition.value = -1
  }

  const clearUndo = () => {
    commandStack.value = []
    currentPosition.value = -1
  }

  const clearRedo = () => {
    redoStack.value = []
  }

  const setMaxStackSize = (size) => {
    const newSize = Math.max(1, size)
    const oldSize = maxStackSize.value
    maxStackSize.value = newSize

    // If reducing size, clear redo stack to maintain consistency
    if (newSize < oldSize) {
      redoStack.value = []
    }

    // Trim existing stacks to new size
    while (commandStack.value.length > maxStackSize.value) {
      commandStack.value.shift()
    }
    while (redoStack.value.length > maxStackSize.value) {
      redoStack.value.shift()
    }
    currentPosition.value = commandStack.value.length - 1
  }

  const createMacro = (commands) => {
    const macro = new MacroCommand(commands)
    macro.description = `${commands.length} commands`
    return macro
  }

  const executeBatch = (commands) => {
    const macro = createMacro(commands)
    return execute(macro)
  }

  const peek = () => {
    return commandStack.value.length > 0
      ? commandStack.value[commandStack.value.length - 1]
      : null
  }

  const getCommandHistory = () => {
    return commandStack.value.map(cmd => ({
      description: cmd.description || (cmd.getDescription ? cmd.getDescription() : 'Unknown'),
      canUndo: cmd.canUndo ? cmd.canUndo() : true
    }))
  }

  return {
    // State
    commandStack,
    redoStack,
    maxStackSize,
    currentPosition,

    // Computed
    canUndo,
    canRedo,
    stackSize,
    lastCommandDescription,

    // Actions
    execute,
    undo,
    redo,
    clear,
    clearUndo,
    clearRedo,
    setMaxStackSize,
    createMacro,
    executeBatch,
    peek,
    getCommandHistory
  }
})

export default useUndoRedoStore
