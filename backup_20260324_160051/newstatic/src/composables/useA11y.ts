/**
 * 无障碍访问 Composable
 * 提供 ARIA 属性、焦点管理和屏幕阅读器支持
 */

import { ref, computed, watch, onMounted, onUnmounted, Ref } from 'vue'

/**
 * 焦点 trap 配置
 */
interface FocusTrapConfig {
  enabled: Ref<boolean>
  container: Ref<HTMLElement | undefined>
  onEscape?: () => void
}

/**
 * ARIA 属性生成器
 */
export function useAria() {
  /**
   * 生成标准的 ARIA 标签
   */
  function getAriaAttributes(props: {
    label?: string
    labelledBy?: string
    describedBy?: string
    description?: string
    expanded?: boolean
    hidden?: boolean
    disabled?: boolean
    required?: boolean
    invalid?: boolean
    live?: 'polite' | 'assertive' | 'off'
  }) {
    const attrs: Record<string, any> = {}

    if (props.label) attrs['aria-label'] = props.label
    if (props.labelledBy) attrs['aria-labelledby'] = props.labelledBy
    if (props.describedBy) attrs['aria-describedby'] = props.describedBy
    if (props.description) attrs['aria-description'] = props.description
    if (props.expanded !== undefined) attrs['aria-expanded'] = props.expanded
    if (props.hidden) attrs['aria-hidden'] = 'true'
    if (props.disabled) attrs['aria-disabled'] = 'true'
    if (props.required) attrs['aria-required'] = 'true'
    if (props.invalid) attrs['aria-invalid'] = 'true'
    if (props.live) attrs['aria-live'] = props.live

    return attrs
  }

  /**
   * 生成任务条的 ARIA 属性
   */
  function getTaskBarAria(task: any) {
    return getAriaAttributes({
      label: `任务: ${task.name}, 工期: ${task.duration}天, 进度: ${task.progress}%, 状态: ${getTaskStatusText(task.status)}`,
      describedBy: `task-${task.id}-desc`
    })
  }

  /**
   * 生成按钮的 ARIA 属性
   */
  function getButtonAria(props: {
    label?: string
    pressed?: boolean
    expanded?: boolean
    disabled?: boolean
    description?: string
  }) {
    return getAriaAttributes({
      ...props,
      label: props.label || props.description
    })
  }

  /**
   * 生成对话框的 ARIA 属性
   */
  function getDialogAria(props: {
    label: string
    describedBy?: string
    modal?: boolean
  }) {
    const attrs = getAriaAttributes({
      label: props.label,
      describedBy: props.describedBy
    })

    attrs['role'] = props.modal ? 'dialog' : 'alertdialog'
    attrs['aria-modal'] = props.modal ? 'true' : 'false'

    return attrs
  }

  /**
   * 获取任务状态文本
   */
  function getTaskStatusText(status: string): string {
    const statusMap: Record<string, string> = {
      not_started: '未开始',
      in_progress: '进行中',
      completed: '已完成',
      delayed: '延期'
    }
    return statusMap[status] || status
  }

  return {
    getAriaAttributes,
    getTaskBarAria,
    getButtonAria,
    getDialogAria,
    getTaskStatusText
  }
}

/**
 * 焦点管理 Composable
 */
export function useFocusManager() {
  // 焦点历史栈（用于对话框关闭后恢复焦点）
  const focusHistory = ref<HTMLElement[]>([])

  /**
   * 保存当前焦点元素
   */
  function saveFocus() {
    const activeElement = document.activeElement as HTMLElement
    if (activeElement && activeElement !== document.body) {
      focusHistory.value.push(activeElement)
    }
  }

  /**
   * 恢复上次焦点
   */
  function restoreFocus() {
    const lastFocused = focusHistory.value.pop()
    if (lastFocused) {
      lastFocused.focus()
    }
  }

  /**
   * 设置焦点到指定元素
   */
  function setFocus(element: HTMLElement | undefined) {
    if (element) {
      // 延迟一帧，确保 DOM 已更新
      requestAnimationFrame(() => {
        element.focus()
      })
    }
  }

  /**
   * 焦点 trap（将焦点限制在容器内）
   */
  function trapFocus(config: FocusTrapConfig) {
    if (!config.enabled.value || !config.container.value) return

    const container = config.container.value
    const focusableElements = container.querySelectorAll(
      'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
    ) as NodeListOf<HTMLElement>

    const firstFocusable = focusableElements[0]
    const lastFocusable = focusableElements[focusableElements.length - 1]

    /**
     * 处理 Tab 键
     */
    function handleTabKey(e: KeyboardEvent) {
      if (e.key !== 'Tab') return

      if (e.shiftKey) {
        // Shift + Tab
        if (document.activeElement === firstFocusable) {
          e.preventDefault()
          lastFocusable.focus()
        }
      } else {
        // Tab
        if (document.activeElement === lastFocusable) {
          e.preventDefault()
          firstFocusable.focus()
        }
      }
    }

    /**
     * 处理 Escape 键
     */
    function handleEscapeKey(e: KeyboardEvent) {
      if (e.key === 'Escape' && config.onEscape) {
        config.onEscape()
      }
    }

    // 添加事件监听
    container.addEventListener('keydown', handleTabKey)
    document.addEventListener('keydown', handleEscapeKey)

    // 焦点到第一个元素
    if (firstFocusable) {
      firstFocusable.focus()
    }

    // 返回清理函数
    return () => {
      container.removeEventListener('keydown', handleTabKey)
      document.removeEventListener('keydown', handleEscapeKey)
    }
  }

  /**
   * 获取元素的焦点索引
   */
  function getFocusIndex(element: HTMLElement): number {
    const focusable = Array.from(
      document.querySelectorAll(
        'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
      )
    ) as HTMLElement[]
    return focusable.indexOf(element)
  }

  /**
   * 移动焦点到下一个/上一个元素
   */
  function moveFocus(direction: 'next' | 'prev' | 'first' | 'last') {
    const focusable = Array.from(
      document.querySelectorAll(
        'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
      )
    ) as HTMLElement[]

    if (focusable.length === 0) return

    const currentIndex = getFocusIndex(document.activeElement as HTMLElement)

    let nextIndex = 0

    switch (direction) {
      case 'next':
        nextIndex = (currentIndex + 1) % focusable.length
        break
      case 'prev':
        nextIndex = currentIndex <= 0 ? focusable.length - 1 : currentIndex - 1
        break
      case 'first':
        nextIndex = 0
        break
      case 'last':
        nextIndex = focusable.length - 1
        break
    }

    focusable[nextIndex].focus()
  }

  return {
    focusHistory,
    saveFocus,
    restoreFocus,
    setFocus,
    trapFocus,
    getFocusIndex,
    moveFocus
  }
}

/**
 * 屏幕阅读器通知 Composable
 */
export function useScreenReader() {
  /**
   * 创建屏幕阅读器通知区域
   */
  const announcer = ref<HTMLDivElement | null>(null)

  /**
   * 初始化通知区域
   */
  function initAnnouncer() {
    if (announcer.value) return

    announcer.value = document.createElement('div')
    announcer.value.setAttribute('role', 'status')
    announcer.value.setAttribute('aria-live', 'polite')
    announcer.value.setAttribute('aria-atomic', 'true')
    announcer.value.style.position = 'absolute'
    announcer.value.style.left = '-10000px'
    announcer.value.style.width = '1px'
    announcer.value.style.height = '1px'
    announcer.value.style.overflow = 'hidden'

    document.body.appendChild(announcer.value)
  }

  /**
   * 通知屏幕阅读器
   */
  function announce(message: string, priority: 'polite' | 'assertive' = 'polite') {
    initAnnouncer()

    if (!announcer.value) return

    announcer.value.setAttribute('aria-live', priority)
    announcer.value.textContent = ''

    // 延迟清空后设置新内容（确保屏幕阅读器能检测到变化）
    requestAnimationFrame(() => {
      requestAnimationFrame(() => {
        if (announcer.value) {
          announcer.value.textContent = message
        }
      })
    })
  }

  /**
   * 清理通知区域
   */
  function destroyAnnouncer() {
    if (announcer.value && announcer.value.parentNode) {
      announcer.value.parentNode.removeChild(announcer.value)
      announcer.value = null
    }
  }

  // 组件卸载时清理
  onUnmounted(() => {
    destroyAnnouncer()
  })

  return {
    announce,
    initAnnouncer,
    destroyAnnouncer
  }
}

/**
 * 键盘导航 Composable
 */
export function useKeyboardNav(options: {
  onArrowUp?: () => void
  onArrowDown?: () => void
  onArrowLeft?: () => void
  onArrowRight?: () => void
  onEnter?: () => void
  onEscape?: () => void
  onHome?: () => void
  onEnd?: () => void
  onPageUp?: () => void
  onPageDown?: () => void
  onSpace?: () => void
}) {
  /**
   * 处理键盘事件
   */
  function handleKeydown(event: KeyboardEvent, callback?: (e: KeyboardEvent) => boolean) {
    // 如果在输入框中，不处理
    const target = event.target as HTMLElement
    if (
      target.tagName === 'INPUT' ||
      target.tagName === 'TEXTAREA' ||
      target.isContentEditable
    ) {
      return
    }

    let handled = true

    switch (event.key) {
      case 'ArrowUp':
        options.onArrowUp?.()
        break
      case 'ArrowDown':
        options.onArrowDown?.()
        break
      case 'ArrowLeft':
        options.onArrowLeft?.()
        break
      case 'ArrowRight':
        options.onArrowRight?.()
        break
      case 'Enter':
        options.onEnter?.()
        break
      case 'Escape':
        options.onEscape?.()
        break
      case 'Home':
        options.onHome?.()
        break
      case 'End':
        options.onEnd?.()
        break
      case 'PageUp':
        options.onPageUp?.()
        break
      case 'PageDown':
        options.onPageDown?.()
        break
      case ' ':
        options.onSpace?.()
        break
      default:
        handled = false
    }

    if (handled) {
      event.preventDefault()
    }

    // 调用自定义回调
    if (callback && !callback(event)) {
      event.preventDefault()
    }
  }

  return {
    handleKeydown
  }
}

/**
 * 高对比度模式检测
 */
export function useHighContrast() {
  const isHighContrast = ref(false)

  function checkHighContrast() {
    const query = window.matchMedia('(prefers-contrast: high)')
    isHighContrast.value = query.matches

    query.addEventListener('change', (e) => {
      isHighContrast.value = e.matches
    })
  }

  onMounted(() => {
    checkHighContrast()
  })

  return {
    isHighContrast
  }
}

/**
 * 减少动画模式检测
 */
export function useReducedMotion() {
  const prefersReducedMotion = ref(false)

  function checkReducedMotion() {
    const query = window.matchMedia('(prefers-reduced-motion: reduce)')
    prefersReducedMotion.value = query.matches

    query.addEventListener('change', (e) => {
      prefersReducedMotion.value = e.matches
    })
  }

  onMounted(() => {
    checkReducedMotion()
  })

  return {
    prefersReducedMotion
  }
}
