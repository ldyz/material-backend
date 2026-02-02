/**
 * 依赖关系创建 Composable
 * 处理鼠标拖拽创建任务依赖关系
 */

import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { progressApi } from '@/api'

/**
 * 使用依赖关系创建功能
 * @param {Object} options - 配置选项
 * @param {Number|Ref} options.projectId - 项目ID
 * @returns {Object} 依赖创建相关的状态和方法
 */
export function useDependencyCreation(options = {}) {
  const {
    onDependencyCreated = null,
    onCancelled = null,
    projectId = null
  } = options

  // 获取项目ID的值（可能是 ref 或普通值）
  const getProjectId = () => {
    if (projectId && typeof projectId === 'object' && 'value' in projectId) {
      return projectId.value
    }
    return projectId
  }

  // 依赖创建状态
  const isCreating = ref(false)
  const sourceTask = ref(null)
  const targetTask = ref(null)
  const mousePosition = ref({ x: 0, y: 0 })

  // 临时连线数据
  const tempLine = ref(null)

  /**
   * 开始创建依赖关系
   * @param {Object} task - 源任务（前置任务）
   */
  const startCreating = (task) => {
    isCreating.value = true
    sourceTask.value = task
    targetTask.value = null
    tempLine.value = null

    ElMessage.info({
      message: '请点击要连接的后置任务，按ESC取消',
      duration: 3000,
      grouping: true
    })

    // 添加ESC键监听
    document.addEventListener('keydown', handleEscKey)
  }

  /**
   * 处理ESC键取消创建
   */
  const handleEscKey = (event) => {
    if (event.key === 'Escape' && isCreating.value) {
      cancelCreating()
    }
  }

  /**
   * 取消创建依赖关系
   */
  const cancelCreating = () => {
    isCreating.value = false
    sourceTask.value = null
    targetTask.value = null
    tempLine.value = null

    document.removeEventListener('keydown', handleEscKey)

    if (onCancelled) {
      onCancelled()
    }
  }

  /**
   * 处理鼠标移动（用于绘制临时连线）
   * @param {MouseEvent} event - 鼠标事件
   * @param {HTMLElement} container - 甘特图容器
   */
  const handleMouseMove = (event, container) => {
    if (!isCreating.value || !container) return

    const rect = container.getBoundingClientRect()
    mousePosition.value = {
      x: event.clientX - rect.left,
      y: event.clientY - rect.top
    }

    // 计算源任务的结束位置（需要从任务条元素获取）
    // 这里简化处理，实际应该通过ref获取任务条元素
    const sourceTaskBar = document.querySelector(`[data-task-id="${sourceTask.value?.id}"]`)
    if (sourceTaskBar) {
      const barRect = sourceTaskBar.getBoundingClientRect()
      tempLine.value = {
        x1: barRect.right - rect.left,
        y1: barRect.top + barRect.height / 2 - rect.top,
        x2: mousePosition.value.x,
        y2: mousePosition.value.y
      }
    }
  }

  /**
   * 处理目标任务点击
   * @param {Object} task - 目标任务（后置任务）
   */
  const handleTargetClick = async (task) => {
    if (!isCreating.value || !sourceTask.value) return

    // 避免自引用
    if (task.id === sourceTask.value.id) {
      ElMessage.warning('不能创建自引用依赖')
      return
    }

    try {
      targetTask.value = task

      // 调用API创建依赖
      await progressApi.createDependencyVisual(
        sourceTask.value.id,
        task.id,
        {
          project_id: getProjectId(),
          type: 'FS',
          lag: 0
        }
      )

      ElMessage.success('依赖关系创建成功')

      // 清理状态
      const createdSource = sourceTask.value
      const createdTarget = targetTask.value
      cancelCreating()

      // 触发回调
      if (onDependencyCreated) {
        onDependencyCreated({
          from: createdSource,
          to: createdTarget,
          type: 'FS'
        })
      }
    } catch (error) {
      console.error('创建依赖失败:', error)
      const errorMsg = error.response?.data?.error || error.response?.data?.message || error.message || '创建依赖失败'
      ElMessage.error(errorMsg)
      cancelCreating()
    }
  }

  /**
   * 检查是否正在创建依赖
   */
  const isDependencyCreatingMode = computed(() => isCreating.value)

  /**
   * 获取源任务
   */
  const getSourceTask = computed(() => sourceTask.value)

  /**
   * 获取临时连线SVG路径
   */
  const getTempLinePath = computed(() => {
    if (!tempLine.value) return ''

    const { x1, y1, x2, y2 } = tempLine.value

    // 使用虚线
    return `M ${x1} ${y1} L ${x2} ${y2}`
  })

  return {
    // 状态
    isCreating,
    sourceTask,
    targetTask,
    mousePosition,
    tempLine,

    // 计算属性
    isDependencyCreatingMode,
    getSourceTask,
    getTempLinePath,

    // 方法
    startCreating,
    cancelCreating,
    handleMouseMove,
    handleTargetClick
  }
}
