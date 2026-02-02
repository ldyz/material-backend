/**
 * 甘特图工具提示 Composable
 * 处理任务条的工具提示显示
 */

import { ref, computed, Ref } from 'vue'
import type { GanttTask, TooltipPosition } from '@/types/gantt'

/**
 * 工具提示状态
 */
export interface TooltipState {
  visible: boolean
  position: TooltipPosition
  task: GanttTask | null
  content: string
}

/**
 * 使用甘特图工具提示
 * @returns 工具提示状态和方法
 */
export function useGanttTooltip() {
  const visible = ref(false)
  const position = ref<TooltipPosition>({ x: 0, y: 0 })
  const task = ref<GanttTask | null>(null)
  const content = ref('')

  /**
   * 显示工具提示
   */
  function show(
    targetTask: GanttTask,
    event: MouseEvent | { x: number; y: number },
    customContent?: string
  ) {
    task.value = targetTask
    position.value = {
      x: event.x + 15,
      y: event.y + 15
    }
    content.value = customContent || generateTooltipContent(targetTask)
    visible.value = true
  }

  /**
   * 隐藏工具提示
   */
  function hide() {
    visible.value = false
  }

  /**
   * 更新工具提示位置
   */
  function update(x: number, y: number) {
    position.value = { x: x + 15, y: y + 15 }
  }

  /**
   * 生成工具提示内容
   */
  function generateTooltipContent(targetTask: GanttTask): string {
    const lines: string[] = []

    // 任务名称
    lines.push(`<strong>${targetTask.name}</strong>`)

    // 时间信息
    lines.push(`<div style="margin-top: 4px;">`)
    lines.push(`<span>开始:</span> ${targetTask.start}<br>`)
    lines.push(`<span>结束:</span> ${targetTask.end}<br>`)
    lines.push(`<span>工期:</span> ${targetTask.duration} 天`)
    lines.push(`</div>`)

    // 进度信息
    lines.push(`<div style="margin-top: 4px;">`)
    lines.push(`<span>进度:</span> ${targetTask.progress}%`)
    lines.push(`</div>`)

    // 状态
    const statusText = getStatusText(targetTask.status)
    lines.push(`<div style="margin-top: 4px;">`)
    lines.push(`<span>状态:</span> ${statusText}`)
    lines.push(`</div>`)

    // 关键路径标记
    if (targetTask.is_critical) {
      lines.push(`<div style="margin-top: 4px; color: #f56c6c;">`)
      lines.push(`<span>⚠ 关键路径任务</span>`)
      lines.push(`</div>`)
    }

    // 里程碑标记
    if (targetTask.is_milestone) {
      lines.push(`<div style="margin-top: 4px; color: #f39c12;">`)
      lines.push(`<span>★ 里程碑</span>`)
      lines.push(`</div>`)
    }

    return lines.join('')
  }

  /**
   * 获取状态文本
   */
  function getStatusText(status: string): string {
    const statusMap: Record<string, string> = {
      not_started: '未开始',
      in_progress: '进行中',
      completed: '已完成',
      delayed: '延期'
    }
    return statusMap[status] || status
  }

  /**
   * 工具提示位置计算（智能定位，避免超出视口）
   */
  function calculatePosition(
    event: MouseEvent,
    tooltipWidth: number = 200,
    tooltipHeight: number = 150
  ): TooltipPosition {
    const padding = 10
    let x = event.clientX + padding
    let y = event.clientY + padding

    // 检查右边界
    if (x + tooltipWidth > window.innerWidth) {
      x = event.clientX - tooltipWidth - padding
    }

    // 检查下边界
    if (y + tooltipHeight > window.innerHeight) {
      y = event.clientY - tooltipHeight - padding
    }

    return { x, y }
  }

  return {
    visible,
    position,
    task,
    content,
    show,
    hide,
    update,
    calculatePosition,
    generateTooltipContent
  }
}

/**
 * 使用工具提示控制器（高级版本）
 * 支持延迟显示、自动隐藏等高级功能
 */
export function useGanttTooltipAdvanced(options: {
  showDelay?: number
  hideDelay?: number
  maxWidth?: number
} = {}) {
  const {
    showDelay = 300,
    hideDelay = 200,
    maxWidth = 300
  } = options

  const base = useGanttTooltip()
  const isPendingShow = ref(false)
  const isPendingHide = ref(false)

  let showTimer: ReturnType<typeof setTimeout> | null = null
  let hideTimer: ReturnType<typeof setTimeout> | null = null

  /**
   * 延迟显示工具提示
   */
  function showWithDelay(
    targetTask: GanttTask,
    event: MouseEvent,
    customContent?: string
  ) {
    // 清除待执行的隐藏操作
    if (hideTimer) {
      clearTimeout(hideTimer)
      hideTimer = null
    }

    // 如果已经显示，直接更新
    if (base.visible.value) {
      base.show(targetTask, event, customContent)
      return
    }

    isPendingShow.value = true

    if (showTimer) {
      clearTimeout(showTimer)
    }

    showTimer = setTimeout(() => {
      base.show(targetTask, event, customContent)
      isPendingShow.value = false
      showTimer = null
    }, showDelay)
  }

  /**
   * 延迟隐藏工具提示
   */
  function hideWithDelay() {
    // 清除待执行的显示操作
    if (showTimer) {
      clearTimeout(showTimer)
      showTimer = null
      isPendingShow.value = false
    }

    if (!base.visible.value) {
      return
    }

    isPendingHide.value = true

    if (hideTimer) {
      clearTimeout(hideTimer)
    }

    hideTimer = setTimeout(() => {
      base.hide()
      isPendingHide.value = false
      hideTimer = null
    }, hideDelay)
  }

  /**
   * 立即显示（取消延迟）
   */
  function showImmediately(targetTask: GanttTask, event: MouseEvent, customContent?: string) {
    if (showTimer) {
      clearTimeout(showTimer)
      showTimer = null
    }
    if (hideTimer) {
      clearTimeout(hideTimer)
      hideTimer = null
    }
    isPendingShow.value = false
    isPendingHide.value = false
    base.show(targetTask, event, customContent)
  }

  /**
   * 立即隐藏（取消延迟）
   */
  function hideImmediately() {
    if (showTimer) {
      clearTimeout(showTimer)
      showTimer = null
    }
    if (hideTimer) {
      clearTimeout(hideTimer)
      hideTimer = null
    }
    isPendingShow.value = false
    isPendingHide.value = false
    base.hide()
  }

  return {
    ...base,
    isPendingShow,
    isPendingHide,
    showWithDelay,
    hideWithDelay,
    showImmediately,
    hideImmediately,
    maxWidth
  }
}
