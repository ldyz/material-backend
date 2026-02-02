/**
 * 甘特图键盘快捷键 Composable
 * 处理键盘导航和快捷操作
 */

import { ref, onMounted, onUnmounted, Ref } from 'vue'
import type { GanttTask } from '@/types/gantt'

/**
 * 键盘快捷键配置
 */
export interface KeyboardShortcut {
  key: string
  ctrl?: boolean
  shift?: boolean
  alt?: boolean
  description: string
  action: () => void
}

/**
 * 键盘导航方向
 */
export type NavigationDirection = 'up' | 'down' | 'left' | 'right' | 'home' | 'end' | 'pageup' | 'pagedown'

/**
 * 使用甘特图键盘功能
 * @param options 配置选项
 * @returns 键盘状态和方法
 */
export function useGanttKeyboard(options: {
  tasks: Ref<GanttTask[]>
  selectedTask: Ref<GanttTask | null>
  onTaskSelect?: (task: GanttTask) => void
  onTaskEdit?: (task: GanttTask) => void
  onTaskDelete?: (task: GanttTask) => void
  onCopy?: () => void
  onPaste?: () => void
  onUndo?: () => void
  onRedo?: () => void
  onSave?: () => void
  onZoomIn?: () => void
  onZoomOut?: () => void
  onNavigateDate?: (direction: 'prev' | 'next' | 'today') => void
  onToggleDependencies?: () => void
  onToggleCriticalPath?: () => void
  onViewModeChange?: (mode: string) => void
}) {
  const {
    tasks,
    selectedTask,
    onTaskSelect,
    onTaskEdit,
    onTaskDelete,
    onCopy,
    onPaste,
    onUndo,
    onRedo,
    onSave,
    onZoomIn,
    onZoomOut,
    onNavigateDate,
    onToggleDependencies,
    onToggleCriticalPath,
    onViewModeChange
  } = options

  // 是否显示快捷键帮助面板
  const showShortcutHelp = ref(false)

  // 所有可用的快捷键
  const shortcuts: KeyboardShortcut[] = [
    {
      key: 'ArrowUp',
      description: '选择上一个任务',
      action: () => navigate('up')
    },
    {
      key: 'ArrowDown',
      description: '选择下一个任务',
      action: () => navigate('down')
    },
    {
      key: 'ArrowLeft',
      description: '选择前一天的任务（依赖关系）',
      action: () => navigate('left')
    },
    {
      key: 'ArrowRight',
      description: '选择后一天的任务（依赖关系）',
      action: () => navigate('right')
    },
    {
      key: 'Enter',
      description: '编辑选中的任务',
      action: () => handleEdit()
    },
    {
      key: 'Delete',
      description: '删除选中的任务',
      action: () => handleDelete()
    },
    {
      key: 'Escape',
      description: '取消选择/关闭对话框',
      action: () => handleEscape()
    },
    {
      key: 'c',
      ctrl: true,
      description: '复制任务',
      action: () => onCopy?.()
    },
    {
      key: 'v',
      ctrl: true,
      description: '粘贴任务',
      action: () => onPaste?.()
    },
    {
      key: 'z',
      ctrl: true,
      description: '撤销',
      action: () => onUndo?.()
    },
    {
      key: 'z',
      ctrl: true,
      shift: true,
      description: '重做',
      action: () => onRedo?.()
    },
    {
      key: 's',
      ctrl: true,
      description: '保存',
      action: () => onSave?.()
    },
    {
      key: '=',
      ctrl: true,
      description: '放大时间轴',
      action: () => onZoomIn?.()
    },
    {
      key: '-',
      ctrl: true,
      description: '缩小时间轴',
      action: () => onZoomOut?.()
    },
    {
      key: 'd',
      alt: true,
      description: '切换依赖关系显示',
      action: () => onToggleDependencies?.()
    },
    {
      key: 'p',
      alt: true,
      description: '切换关键路径显示',
      action: () => onToggleCriticalPath?.()
    },
    {
      key: '?',
      description: '显示/隐藏快捷键帮助',
      action: () => toggleShortcutHelp()
    }
  ]

  /**
   * 处理键盘事件
   */
  function handleKeydown(event: KeyboardEvent) {
    // 如果在输入框中，不处理快捷键
    const target = event.target as HTMLElement
    if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA' || target.isContentEditable) {
      return
    }

    // 查找匹配的快捷键
    const matchingShortcut = shortcuts.find(shortcut => {
      const keyMatch = event.key === shortcut.key
      const ctrlMatch = shortcut.ctrl ? event.ctrlKey || event.metaKey : !event.ctrlKey && !event.metaKey
      const shiftMatch = shortcut.shift ? event.shiftKey : !event.shiftKey
      const altMatch = shortcut.alt ? event.altKey : !event.altKey

      return keyMatch && ctrlMatch && shiftMatch && altMatch
    })

    if (matchingShortcut) {
      event.preventDefault()
      matchingShortcut.action()
    }
  }

  /**
   * 键盘导航
   */
  function navigate(direction: NavigationDirection) {
    if (!selectedTask.value || tasks.value.length === 0) return

    const currentIndex = tasks.value.findIndex(t => t.id === selectedTask.value!.id)
    if (currentIndex === -1) return

    let nextIndex = currentIndex
    let nextTask: GanttTask | null = null

    switch (direction) {
      case 'up':
        nextIndex = Math.max(0, currentIndex - 1)
        nextTask = tasks.value[nextIndex]
        break
      case 'down':
        nextIndex = Math.min(tasks.value.length - 1, currentIndex + 1)
        nextTask = tasks.value[nextIndex]
        break
      case 'left':
        // 查找前驱任务
        if (selectedTask.value.predecessors && selectedTask.value.predecessors.length > 0) {
          const predecessorId = selectedTask.value.predecessors[0].predecessor_id
          nextTask = tasks.value.find(t => t.id === predecessorId) || null
        }
        break
      case 'right':
        // 查找后继任务
        if (selectedTask.value.successors && selectedTask.value.successors.length > 0) {
          const successorId = selectedTask.value.successors[0].successor_id
          nextTask = tasks.value.find(t => t.id === successorId) || null
        }
        break
      case 'home':
        nextTask = tasks.value[0]
        break
      case 'end':
        nextTask = tasks.value[tasks.value.length - 1]
        break
      case 'pageup':
        nextIndex = Math.max(0, currentIndex - 10)
        nextTask = tasks.value[nextIndex]
        break
      case 'pagedown':
        nextIndex = Math.min(tasks.value.length - 1, currentIndex + 10)
        nextTask = tasks.value[nextIndex]
        break
    }

    if (nextTask && onTaskSelect) {
      onTaskSelect(nextTask)
    }
  }

  /**
   * 编辑任务
   */
  function handleEdit() {
    if (selectedTask.value && onTaskEdit) {
      onTaskEdit(selectedTask.value)
    }
  }

  /**
   * 删除任务
   */
  function handleDelete() {
    if (selectedTask.value && onTaskDelete) {
      onTaskDelete(selectedTask.value)
    }
  }

  /**
   * 取消/关闭
   */
  function handleEscape() {
    if (showShortcutHelp.value) {
      showShortcutHelp.value = false
    } else if (selectedTask.value) {
      // 取消选择
      // 这里需要与主组件的状态管理集成
    }
  }

  /**
   * 切换快捷键帮助显示
   */
  function toggleShortcutHelp() {
    showShortcutHelp.value = !showShortcutHelp.value
  }

  /**
   * 获取快捷键文本表示
   */
  function getShortcutText(shortcut: KeyboardShortcut): string {
    const parts: string[] = []
    if (shortcut.ctrl) parts.push('Ctrl')
    if (shortcut.shift) parts.push('Shift')
    if (shortcut.alt) parts.push('Alt')
    parts.push(shortcut.key)
    return parts.join(' + ')
  }

  // 绑定/解绑事件
  onMounted(() => {
    document.addEventListener('keydown', handleKeydown)
  })

  onUnmounted(() => {
    document.removeEventListener('keydown', handleKeydown)
  })

  return {
    showShortcutHelp,
    shortcuts,
    getShortcutText,
    toggleShortcutHelp,
    navigate
  }
}

/**
 * 快捷键帮助面板组件所需的数据
 */
export function useShortcutPanel() {
  const shortcutGroups = [
    {
      title: '导航',
      shortcuts: [
        { key: '↑', description: '选择上一个任务' },
        { key: '↓', description: '选择下一个任务' },
        { key: '←', description: '选择前驱任务' },
        { key: '→', description: '选择后继任务' },
        { key: 'Home', description: '跳到第一个任务' },
        { key: 'End', description: '跳到最后一个任务' }
      ]
    },
    {
      title: '编辑',
      shortcuts: [
        { key: 'Enter', description: '编辑任务' },
        { key: 'Delete', description: '删除任务' },
        { key: 'Escape', description: '取消操作' },
        { key: 'Ctrl+C', description: '复制' },
        { key: 'Ctrl+V', description: '粘贴' },
        { key: 'Ctrl+Z', description: '撤销' },
        { key: 'Ctrl+Shift+Z', description: '重做' },
        { key: 'Ctrl+S', description: '保存' }
      ]
    },
    {
      title: '视图',
      shortcuts: [
        { key: 'Ctrl++', description: '放大时间轴' },
        { key: 'Ctrl+-', description: '缩小时间轴' },
        { key: 'Alt+D', description: '切换依赖关系' },
        { key: 'Alt+P', description: '切换关键路径' },
        { key: '?', description: '显示此帮助' }
      ]
    }
  ]

  return {
    shortcutGroups
  }
}
