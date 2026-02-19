import { describe, it, expect, beforeEach, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useUndoRedoStore } from '../undoRedoStore.js'

describe('Undo/Redo Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  describe('Command Execution', () => {
    it('should execute a command', () => {
      const store = useUndoRedoStore()
      const command = {
        execute: vi.fn(() => ({ success: true })),
        undo: vi.fn(),
      }

      const result = store.execute(command)

      expect(command.execute).toHaveBeenCalled()
      expect(result.success).toBe(true)
    })

    it('should add command to stack after execution', () => {
      const store = useUndoRedoStore()
      const command = {
        execute: vi.fn(() => ({ success: true })),
        undo: vi.fn(),
      }

      store.execute(command)

      expect(store.commandStack.length).toBe(1)
      expect(store.commandStack[0]).toBe(command)
    })

    it('should set canUndo to true after command execution', () => {
      const store = useUndoRedoStore()
      const command = {
        execute: vi.fn(() => ({ success: true })),
        undo: vi.fn(),
      }

      store.execute(command)

      expect(store.canUndo).toBe(true)
      expect(store.canRedo).toBe(false)
    })

    it('should clear redo stack when new command is executed', () => {
      const store = useUndoRedoStore()
      const command1 = {
        execute: vi.fn(() => ({ success: true })),
        undo: vi.fn(),
      }
      const command2 = {
        execute: vi.fn(() => ({ success: true })),
        undo: vi.fn(),
      }

      store.execute(command1)
      store.undo()
      store.execute(command2)

      expect(store.canRedo).toBe(false)
    })

    it('should not add failed commands to stack', () => {
      const store = useUndoRedoStore()
      const command = {
        execute: vi.fn(() => ({ success: false, error: 'Failed' })),
        undo: vi.fn(),
      }

      store.execute(command)

      expect(store.commandStack.length).toBe(0)
      expect(store.canUndo).toBe(false)
    })

    it('should handle command execution errors', () => {
      const store = useUndoRedoStore()
      const command = {
        execute: vi.fn(() => {
          throw new Error('Execution error')
        }),
        undo: vi.fn(),
      }

      expect(() => store.execute(command)).not.toThrow()
    })
  })

  describe('Undo Operations', () => {
    it('should undo the last command', () => {
      const store = useUndoRedoStore()
      const command = {
        execute: vi.fn(() => ({ success: true })),
        undo: vi.fn(),
      }

      store.execute(command)
      store.undo()

      expect(command.undo).toHaveBeenCalled()
    })

    it('should move command to redo stack after undo', () => {
      const store = useUndoRedoStore()
      const command = {
        execute: vi.fn(() => ({ success: true })),
        undo: vi.fn(),
      }

      store.execute(command)
      store.undo()

      expect(store.canUndo).toBe(false)
      expect(store.canRedo).toBe(true)
    })

    it('should undo multiple commands', () => {
      const store = useUndoRedoStore()
      const command1 = {
        execute: vi.fn(() => ({ success: true })),
        undo: vi.fn(),
      }
      const command2 = {
        execute: vi.fn(() => ({ success: true })),
        undo: vi.fn(),
      }

      store.execute(command1)
      store.execute(command2)
      store.undo()
      store.undo()

      expect(command1.undo).toHaveBeenCalled()
      expect(command2.undo).toHaveBeenCalled()
      expect(store.canUndo).toBe(false)
    })

    it('should handle undo when stack is empty', () => {
      const store = useUndoRedoStore()

      expect(() => store.undo()).not.toThrow()
      expect(store.canUndo).toBe(false)
    })

    it('should call undo callback if provided', () => {
      const store = useUndoRedoStore()
      const callback = vi.fn()
      const command = {
        execute: vi.fn(() => ({ success: true })),
        undo: vi.fn(),
      }

      store.execute(command)
      store.undo(callback)

      expect(callback).toHaveBeenCalled()
    })
  })

  describe('Redo Operations', () => {
    it('should redo the last undone command', () => {
      const store = useUndoRedoStore()
      const command = {
        execute: vi.fn(() => ({ success: true })),
        undo: vi.fn(),
      }

      store.execute(command)
      store.undo()
      store.redo()

      expect(command.execute).toHaveBeenCalledTimes(2)
    })

    it('should move command back to command stack after redo', () => {
      const store = useUndoRedoStore()
      const command = {
        execute: vi.fn(() => ({ success: true })),
        undo: vi.fn(),
      }

      store.execute(command)
      store.undo()
      store.redo()

      expect(store.canUndo).toBe(true)
      expect(store.canRedo).toBe(false)
    })

    it('should redo multiple commands', () => {
      const store = useUndoRedoStore()
      const command1 = {
        execute: vi.fn(() => ({ success: true })),
        undo: vi.fn(),
      }
      const command2 = {
        execute: vi.fn(() => ({ success: true })),
        undo: vi.fn(),
      }

      store.execute(command1)
      store.execute(command2)
      store.undo()
      store.undo()
      store.redo()
      store.redo()

      expect(command1.execute).toHaveBeenCalledTimes(2)
      expect(command2.execute).toHaveBeenCalledTimes(2)
    })

    it('should handle redo when stack is empty', () => {
      const store = useUndoRedoStore()

      expect(() => store.redo()).not.toThrow()
      expect(store.canRedo).toBe(false)
    })

    it('should call redo callback if provided', () => {
      const store = useUndoRedoStore()
      const callback = vi.fn()
      const command = {
        execute: vi.fn(() => ({ success: true })),
        undo: vi.fn(),
      }

      store.execute(command)
      store.undo()
      store.redo(callback)

      expect(callback).toHaveBeenCalled()
    })
  })

  describe('Command Stack Limits', () => {
    it('should respect max stack size', () => {
      const store = useUndoRedoStore()
      store.setMaxStackSize(5)

      for (let i = 0; i < 10; i++) {
        const command = {
          execute: vi.fn(() => ({ success: true })),
          undo: vi.fn(),
        }
        store.execute(command)
      }

      expect(store.commandStack.length).toBe(5)
    })

    it('should remove oldest commands when limit is reached', () => {
      const store = useUndoRedoStore()
      store.setMaxStackSize(3)

      const command1 = {
        execute: vi.fn(() => ({ success: true, data: 1 })),
        undo: vi.fn(),
      }
      const command2 = {
        execute: vi.fn(() => ({ success: true, data: 2 })),
        undo: vi.fn(),
      }
      const command3 = {
        execute: vi.fn(() => ({ success: true, data: 3 })),
        undo: vi.fn(),
      }
      const command4 = {
        execute: vi.fn(() => ({ success: true, data: 4 })),
        undo: vi.fn(),
      }

      store.execute(command1)
      store.execute(command2)
      store.execute(command3)
      store.execute(command4)

      expect(store.commandStack.length).toBe(3)
      expect(store.commandStack[0]).toBe(command2)
    })

    it('should clear redo stack when limit is reduced', () => {
      const store = useUndoRedoStore()

      for (let i = 0; i < 5; i++) {
        const command = {
          execute: vi.fn(() => ({ success: true })),
          undo: vi.fn(),
        }
        store.execute(command)
      }

      store.undo()
      store.undo()

      store.setMaxStackSize(2)

      expect(store.canRedo).toBe(false)
    })
  })

  describe('Macro Commands', () => {
    it('should execute macro command', () => {
      const store = useUndoRedoStore()
      const subCommand1 = {
        execute: vi.fn(() => ({ success: true })),
        undo: vi.fn(),
      }
      const subCommand2 = {
        execute: vi.fn(() => ({ success: true })),
        undo: vi.fn(),
      }

      const macro = store.createMacro([subCommand1, subCommand2])
      store.execute(macro)

      expect(subCommand1.execute).toHaveBeenCalled()
      expect(subCommand2.execute).toHaveBeenCalled()
    })

    it('should undo macro command', () => {
      const store = useUndoRedoStore()
      const subCommand1 = {
        execute: vi.fn(() => ({ success: true })),
        undo: vi.fn(),
      }
      const subCommand2 = {
        execute: vi.fn(() => ({ success: true })),
        undo: vi.fn(),
      }

      const macro = store.createMacro([subCommand1, subCommand2])
      store.execute(macro)
      store.undo()

      expect(subCommand1.undo).toHaveBeenCalled()
      expect(subCommand2.undo).toHaveBeenCalled()
    })

    it('should stop macro execution on first failure', () => {
      const store = useUndoRedoStore()
      const subCommand1 = {
        execute: vi.fn(() => ({ success: true })),
        undo: vi.fn(),
      }
      const subCommand2 = {
        execute: vi.fn(() => ({ success: false })),
        undo: vi.fn(),
      }
      const subCommand3 = {
        execute: vi.fn(() => ({ success: true })),
        undo: vi.fn(),
      }

      const macro = store.createMacro([subCommand1, subCommand2, subCommand3])
      const result = store.execute(macro)

      expect(subCommand1.execute).toHaveBeenCalled()
      expect(subCommand2.execute).toHaveBeenCalled()
      expect(subCommand3.execute).not.toHaveBeenCalled()
      expect(result.success).toBe(false)
    })

    it('should undo all executed sub-commands on macro failure', () => {
      const store = useUndoRedoStore()
      const subCommand1 = {
        execute: vi.fn(() => ({ success: true })),
        undo: vi.fn(),
      }
      const subCommand2 = {
        execute: vi.fn(() => ({ success: false })),
        undo: vi.fn(),
      }

      const macro = store.createMacro([subCommand1, subCommand2])
      store.execute(macro)

      expect(subCommand1.undo).toHaveBeenCalled()
    })
  })

  describe('Batch Operations', () => {
    it('should execute batch commands', () => {
      const store = useUndoRedoStore()
      const commands = Array.from({ length: 5 }, (_, i) => ({
        execute: vi.fn(() => ({ success: true })),
        undo: vi.fn(),
      }))

      store.executeBatch(commands)

      commands.forEach(cmd => {
        expect(cmd.execute).toHaveBeenCalled()
      })
    })

    it('should treat batch as single undo/redo operation', () => {
      const store = useUndoRedoStore()
      const commands = Array.from({ length: 3 }, () => ({
        execute: vi.fn(() => ({ success: true })),
        undo: vi.fn(),
      }))

      store.executeBatch(commands)

      expect(store.commandStack.length).toBe(1)
      expect(store.canUndo).toBe(true)

      store.undo()

      commands.forEach(cmd => {
        expect(cmd.undo).toHaveBeenCalled()
      })
    })

    it('should handle partial batch failures', () => {
      const store = useUndoRedoStore()
      const commands = [
        {
          execute: vi.fn(() => ({ success: true })),
          undo: vi.fn(),
        },
        {
          execute: vi.fn(() => ({ success: false })),
          undo: vi.fn(),
        },
        {
          execute: vi.fn(() => ({ success: true })),
          undo: vi.fn(),
        },
      ]

      const result = store.executeBatch(commands)

      expect(result.success).toBe(false)
      expect(commands[0].undo).toHaveBeenCalled()
    })
  })

  describe('Clear Operations', () => {
    it('should clear command stack', () => {
      const store = useUndoRedoStore()

      for (let i = 0; i < 5; i++) {
        const command = {
          execute: vi.fn(() => ({ success: true })),
          undo: vi.fn(),
        }
        store.execute(command)
      }

      store.clear()

      expect(store.commandStack.length).toBe(0)
      expect(store.canUndo).toBe(false)
      expect(store.canRedo).toBe(false)
    })

    it('should clear only undo stack', () => {
      const store = useUndoRedoStore()
      const command = {
        execute: vi.fn(() => ({ success: true })),
        undo: vi.fn(),
      }

      store.execute(command)
      store.undo()
      store.clearUndo()

      expect(store.canUndo).toBe(false)
      expect(store.canRedo).toBe(true)
    })

    it('should clear only redo stack', () => {
      const store = useUndoRedoStore()
      const command = {
        execute: vi.fn(() => ({ success: true })),
        undo: vi.fn(),
      }

      store.execute(command)
      store.undo()
      store.clearRedo()

      expect(store.canUndo).toBe(true)
      expect(store.canRedo).toBe(false)
    })
  })

  describe('State Queries', () => {
    it('should provide current position', () => {
      const store = useUndoRedoStore()
      const command = {
        execute: vi.fn(() => ({ success: true })),
        undo: vi.fn(),
      }

      expect(store.currentPosition).toBe(-1)

      store.execute(command)

      expect(store.currentPosition).toBe(0)
    })

    it('should provide stack size', () => {
      const store = useUndoRedoStore()

      expect(store.stackSize).toBe(0)

      for (let i = 0; i < 3; i++) {
        const command = {
          execute: vi.fn(() => ({ success: true })),
          undo: vi.fn(),
        }
        store.execute(command)
      }

      expect(store.stackSize).toBe(3)
    })

    it('should peek at current command', () => {
      const store = useUndoRedoStore()
      const command = {
        execute: vi.fn(() => ({ success: true })),
        undo: vi.fn(),
      }

      store.execute(command)

      expect(store.peek()).toBe(command)
    })

    it('should return null when peeking empty stack', () => {
      const store = useUndoRedoStore()

      expect(store.peek()).toBeNull()
    })
  })

  describe('Command Metadata', () => {
    it('should store command description', () => {
      const store = useUndoRedoStore()
      const command = {
        execute: vi.fn(() => ({ success: true })),
        undo: vi.fn(),
        description: 'Create task',
      }

      store.execute(command)

      expect(store.lastCommandDescription).toBe('Create task')
    })

    it('should provide command history', () => {
      const store = useUndoRedoStore()
      const command1 = {
        execute: vi.fn(() => ({ success: true })),
        undo: vi.fn(),
        description: 'Command 1',
      }
      const command2 = {
        execute: vi.fn(() => ({ success: true })),
        undo: vi.fn(),
        description: 'Command 2',
      }

      store.execute(command1)
      store.execute(command2)

      const history = store.getCommandHistory()
      expect(history).toHaveLength(2)
      expect(history[0].description).toBe('Command 1')
      expect(history[1].description).toBe('Command 2')
    })
  })

  describe('Edge Cases', () => {
    it('should handle command without undo function', () => {
      const store = useUndoRedoStore()
      const command = {
        execute: vi.fn(() => ({ success: true })),
      }

      store.execute(command)

      expect(() => store.undo()).not.toThrow()
    })

    it('should handle null commands', () => {
      const store = useUndoRedoStore()

      expect(() => store.execute(null)).not.toThrow()
    })

    it('should handle very large stack size', () => {
      const store = useUndoRedoStore()
      store.setMaxStackSize(10000)

      for (let i = 0; i < 100; i++) {
        const command = {
          execute: vi.fn(() => ({ success: true })),
          undo: vi.fn(),
        }
        store.execute(command)
      }

      expect(store.commandStack.length).toBe(100)
    })
  })
})
