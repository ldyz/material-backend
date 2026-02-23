/**
 * 任务操作命令类
 *
 * 使用命令模式实现任务操作的撤销/重做功能
 *
 * @module taskCommands
 * @date 2026-02-21
 */

import { Command } from '@/stores/undoRedoStore'
import { progressApi } from '@/api'

/**
 * 创建任务命令
 */
export class CreateTaskCommand extends Command {
  /**
   * @param {number} projectId - 项目ID
   * @param {Object} taskData - 任务数据
   * @param {Function} onSuccess - 成功回调
   * @param {Function} onError - 错误回调
   */
  constructor(projectId, taskData, onSuccess = null, onError = null) {
    super()
    this.projectId = projectId
    this.taskData = taskData
    this.createdTask = null
    this.onSuccess = onSuccess
    this.onError = onError
    this.description = `创建任务: ${taskData.name || '未命名'}`
  }

  async execute() {
    try {
      const response = await progressApi.createTask(this.projectId, this.taskData)
      this.createdTask = response.data

      if (this.onSuccess) {
        this.onSuccess(this.createdTask)
      }

      return { success: true, data: this.createdTask }
    } catch (error) {
      if (this.onError) {
        this.onError(error)
      }
      return { success: false, error }
    }
  }

  async undo() {
    if (this.createdTask) {
      try {
        await progressApi.deleteTask(this.createdTask.id)
        return { success: true }
      } catch (error) {
        console.error('撤销创建任务失败:', error)
        return { success: false, error }
      }
    }
    return { success: true }
  }

  getDescription() {
    return this.description
  }
}

/**
 * 更新任务命令
 */
export class UpdateTaskCommand extends Command {
  /**
   * @param {number} taskId - 任务ID
   * @param {Object} originalData - 原始数据
   * @param {Object} newData - 新数据
   * @param {Function} onSuccess - 成功回调
   * @param {Function} onError - 错误回调
   */
  constructor(taskId, originalData, newData, onSuccess = null, onError = null) {
    super()
    this.taskId = taskId
    this.originalData = { ...originalData }
    this.newData = { ...newData }
    this.onSuccess = onSuccess
    this.onError = onError
    this.description = `更新任务: ${newData.name || originalData.name || '未命名'}`
  }

  async execute() {
    try {
      const response = await progressApi.updateTask(this.taskId, this.newData)

      if (this.onSuccess) {
        this.onSuccess(response.data)
      }

      return { success: true, data: response.data }
    } catch (error) {
      if (this.onError) {
        this.onError(error)
      }
      return { success: false, error }
    }
  }

  async undo() {
    try {
      await progressApi.updateTask(this.taskId, this.originalData)
      return { success: true }
    } catch (error) {
      console.error('撤销更新任务失败:', error)
      return { success: false, error }
    }
  }

  getDescription() {
    return this.description
  }
}

/**
 * 删除任务命令
 */
export class DeleteTaskCommand extends Command {
  /**
   * @param {Object} task - 任务对象
   * @param {number} projectId - 项目ID（用于撤销时重新创建）
   * @param {Function} onSuccess - 成功回调
   * @param {Function} onError - 错误回调
   */
  constructor(task, projectId, onSuccess = null, onError = null) {
    super()
    this.task = { ...task }
    this.projectId = projectId
    this.childTasks = []
    this.onSuccess = onSuccess
    this.onError = onError
    this.description = `删除任务: ${task.name || '未命名'}`
  }

  async execute() {
    try {
      // 先删除所有子任务
      if (this.childTasks.length > 0) {
        for (const child of this.childTasks) {
          await progressApi.deleteTask(child.id)
        }
      }

      // 删除主任务
      await progressApi.deleteTask(this.task.id)

      if (this.onSuccess) {
        this.onSuccess(this.task)
      }

      return { success: true }
    } catch (error) {
      if (this.onError) {
        this.onError(error)
      }
      return { success: false, error }
    }
  }

  async undo() {
    try {
      // 重新创建主任务
      const { id, ...taskData } = this.task
      const response = await progressApi.createTask(this.projectId, taskData)
      const restoredTask = response.data

      // 重新创建子任务
      if (this.childTasks.length > 0) {
        for (const child of this.childTasks) {
          const { id: childId, ...childData } = child
          childData.parent_id = restoredTask.id
          await progressApi.createTask(this.projectId, childData)
        }
      }

      return { success: true, data: restoredTask }
    } catch (error) {
      console.error('撤销删除任务失败:', error)
      return { success: false, error }
    }
  }

  /**
   * 设置子任务列表（用于撤销时恢复）
   */
  setChildTasks(children) {
    this.childTasks = children
  }

  getDescription() {
    const total = 1 + this.childTasks.length
    return total > 1 ? `删除任务及 ${this.childTasks.length} 个子任务` : this.description
  }
}

/**
 * 移动任务命令（调整父子关系和排序）
 */
export class MoveTaskCommand extends Command {
  /**
   * @param {number} taskId - 任务ID
   * @param {number|null} fromParentId - 原父任务ID
   * @param {number|null} toParentId - 新父任务ID
   * @param {string} position - 位置: 'before', 'after', 'child'
   * @param {Object} referenceTask - 参考任务
   * @param {Function} onSuccess - 成功回调
   * @param {Function} onError - 错误回调
   */
  constructor(taskId, fromParentId, toParentId, position, referenceTask, onSuccess = null, onError = null) {
    super()
    this.taskId = taskId
    this.fromParentId = fromParentId
    this.toParentId = toParentId
    this.position = position
    this.referenceTask = referenceTask
    this.fromPosition = null
    this.toPosition = null
    this.onSuccess = onSuccess
    this.onError = onError
    this.description = '移动任务'
  }

  async execute() {
    try {
      // 构建移动数据
      const moveData = {
        parent_id: this.toParentId
      }

      // 根据位置设置排序
      if (this.position === 'child') {
        this.description = `移动到 "${this.referenceTask?.name || '根'}" 下`
      } else if (this.position === 'before') {
        this.description = `移动到 "${this.referenceTask?.name || '根'}" 之前`
      } else if (this.position === 'after') {
        this.description = `移动到 "${this.referenceTask?.name || '根'}" 之后`
      }

      const response = await progressApi.updateTask(this.taskId, moveData)

      if (this.onSuccess) {
        this.onSuccess(response.data)
      }

      return { success: true, data: response.data }
    } catch (error) {
      if (this.onError) {
        this.onError(error)
      }
      return { success: false, error }
    }
  }

  async undo() {
    try {
      await progressApi.updateTask(this.taskId, {
        parent_id: this.fromParentId
      })
      return { success: true }
    } catch (error) {
      console.error('撤销移动任务失败:', error)
      return { success: false, error }
    }
  }

  getDescription() {
    return this.description
  }
}

/**
 * 批量更新任务命令
 */
export class BatchUpdateTasksCommand extends Command {
  /**
   * @param {Array} tasks - 任务数组
   * @param {Object} updates - 更新数据
   * @param {number} projectId - 项目ID
   * @param {Function} onSuccess - 成功回调
   * @param {Function} onError - 错误回调
   */
  constructor(tasks, updates, projectId, onSuccess = null, onError = null) {
    super()
    this.tasks = tasks.map(t => ({ ...t }))
    this.updates = { ...updates }
    this.projectId = projectId
    this.originalData = []
    this.onSuccess = onSuccess
    this.onError = onError
    this.description = `批量更新 ${tasks.length} 个任务`
  }

  async execute() {
    try {
      const results = []

      for (const task of this.tasks) {
        // 保存原始数据
        this.originalData.push({ id: task.id, data: { ...task } })

        const response = await progressApi.updateTask(task.id, {
          ...this.updates,
          project_id: this.projectId
        })
        results.push(response.data)
      }

      if (this.onSuccess) {
        this.onSuccess(results)
      }

      return { success: true, data: results }
    } catch (error) {
      if (this.onError) {
        this.onError(error)
      }
      return { success: false, error }
    }
  }

  async undo() {
    try {
      for (const original of this.originalData) {
        await progressApi.updateTask(original.id, original.data)
      }
      return { success: true }
    } catch (error) {
      console.error('撤销批量更新失败:', error)
      return { success: false, error }
    }
  }

  getDescription() {
    return this.description
  }
}

/**
 * 复制任务命令
 */
export class DuplicateTaskCommand extends Command {
  /**
   * @param {Object} task - 原任务对象
   * @param {number} projectId - 项目ID
   * @param {Function} onSuccess - 成功回调
   * @param {Function} onError - 错误回调
   */
  constructor(task, projectId, onSuccess = null, onError = null) {
    super()
    this.task = { ...task }
    this.projectId = projectId
    this.duplicatedTask = null
    this.onSuccess = onSuccess
    this.onError = onError
    this.description = `复制任务: ${task.name || '未命名'}`
  }

  async execute() {
    try {
      // 移除 ID 和可能需要重置的字段
      const { id, created_at, updated_at, ...taskData } = this.task
      taskData.name = `${taskData.name} (副本)`
      taskData.progress = 0

      const response = await progressApi.createTask(this.projectId, taskData)
      this.duplicatedTask = response.data

      if (this.onSuccess) {
        this.onSuccess(this.duplicatedTask)
      }

      return { success: true, data: this.duplicatedTask }
    } catch (error) {
      if (this.onError) {
        this.onError(error)
      }
      return { success: false, error }
    }
  }

  async undo() {
    if (this.duplicatedTask) {
      try {
        await progressApi.deleteTask(this.duplicatedTask.id)
        return { success: true }
      } catch (error) {
        console.error('撤销复制任务失败:', error)
        return { success: false, error }
      }
    }
    return { success: true }
  }

  getDescription() {
    return this.description
  }
}

/**
 * 转换为里程碑命令
 */
export class ConvertToMilestoneCommand extends Command {
  /**
   * @param {Object} task - 任务对象
   * @param {boolean} toMilestone - 是否转换为里程碑
   * @param {Function} onSuccess - 成功回调
   * @param {Function} onError - 错误回调
   */
  constructor(task, toMilestone = true, onSuccess = null, onError = null) {
    super()
    this.task = { ...task }
    this.toMilestone = toMilestone
    this.originalData = {
      start_date: task.start_date || task.start,
      end_date: task.end_date || task.end,
      is_milestone: task.is_milestone || false,
      duration: task.duration || 0
    }
    this.onSuccess = onSuccess
    this.onError = onError
    this.description = toMilestone ? `转换为里程碑: ${task.name || '未命名'}` : `转换为任务: ${task.name || '未命名'}`
  }

  async execute() {
    try {
      let updateData

      if (this.toMilestone) {
        // 转换为里程碑：将开始日期设置为结束日期
        const endDate = this.task.end_date || this.task.end
        updateData = {
          start_date: endDate,
          end_date: endDate,
          is_milestone: true,
          duration: 0
        }
      } else {
        // 转换为普通任务：恢复原始数据
        updateData = { ...this.originalData }
      }

      const response = await progressApi.updateTask(this.task.id, updateData)

      if (this.onSuccess) {
        this.onSuccess(response.data)
      }

      return { success: true, data: response.data }
    } catch (error) {
      if (this.onError) {
        this.onError(error)
      }
      return { success: false, error }
    }
  }

  async undo() {
    try {
      await progressApi.updateTask(this.task.id, this.originalData)
      return { success: true }
    } catch (error) {
      console.error('撤销转换里程碑失败:', error)
      return { success: false, error }
    }
  }

  getDescription() {
    return this.description
  }
}

export default {
  CreateTaskCommand,
  UpdateTaskCommand,
  DeleteTaskCommand,
  MoveTaskCommand,
  BatchUpdateTasksCommand,
  DuplicateTaskCommand,
  ConvertToMilestoneCommand
}
