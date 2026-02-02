/**
 * 触摸手势 Composable
 * 为移动端提供触摸手势支持
 */

import { ref, onMounted, onUnmounted, Ref } from 'vue'
import type { TouchGestureHandlers, TouchPoint, TouchGestureType } from '@/types/gantt'

/**
 * 触摸手势配置选项
 */
export interface UseTouchGesturesOptions {
  // 元素引用
  elementRef: Ref<HTMLElement | undefined>
  // 手势处理器
  handlers: TouchGestureHandlers
  // 滑动阈值（像素）
  swipeThreshold?: number
  // 长按延迟（毫秒）
  longPressDelay?: number
  // 双指缩放灵敏度
  pinchSensitivity?: number
  // 是否启用 passive 事件监听
  passive?: boolean
}

/**
 * 手势状态
 */
interface GestureState {
  startPoint: TouchPoint | null
  currentPoint: TouchPoint | null
  startTime: number
  longPressTimer: ReturnType<typeof setTimeout> | null
  initialDistance: number
  lastScale: number
}

/**
 * 触摸手势 Hook
 * @param options 配置选项
 * @returns 解绑函数
 */
export function useTouchGestures(options: UseTouchGesturesOptions) {
  const {
    elementRef,
    handlers,
    swipeThreshold = 50,
    longPressDelay = 500,
    pinchSensitivity = 0.001,
    passive = true
  } = options

  const state: GestureState = {
    startPoint: null,
    currentPoint: null,
    startTime: 0,
    longPressTimer: null,
    initialDistance: 0,
    lastScale: 1
  }

  /**
   * 获取触摸点坐标
   */
  function getTouchPoint(touch: Touch): TouchPoint {
    return {
      x: touch.clientX,
      y: touch.clientY
    }
  }

  /**
   * 获取两点之间的距离
   */
  function getDistance(touch1: Touch, touch2: Touch): number {
    const dx = touch1.clientX - touch2.clientX
    const dy = touch1.clientY - touch2.clientY
    return Math.sqrt(dx * dx + dy * dy)
  }

  /**
   * 处理触摸开始
   */
  function handleTouchStart(event: TouchEvent) {
    if (!event.touches || event.touches.length === 0) return

    // 单指手势
    if (event.touches.length === 1) {
      const touch = event.touches[0]
      state.startPoint = getTouchPoint(touch)
      state.currentPoint = state.startPoint
      state.startTime = Date.now()

      // 设置长按定时器
      state.longPressTimer = setTimeout(() => {
        if (handlers.onLongPress) {
          handlers.onLongPress()
        }
      }, longPressDelay)
    }
    // 双指缩放
    else if (event.touches.length === 2) {
      state.initialDistance = getDistance(event.touches[0], event.touches[1])
      state.lastScale = 1
    }
  }

  /**
   * 处理触摸移动
   */
  function handleTouchMove(event: TouchEvent) {
    if (!event.touches || event.touches.length === 0) return

    // 清除长按定时器
    if (state.longPressTimer) {
      clearTimeout(state.longPressTimer)
      state.longPressTimer = null
    }

    // 单指手势
    if (event.touches.length === 1 && state.startPoint) {
      const touch = event.touches[0]
      state.currentPoint = getTouchPoint(touch)
    }
    // 双指缩放
    else if (event.touches.length === 2 && handlers.onPinch) {
      const currentDistance = getDistance(event.touches[0], event.touches[1])
      const scale = currentDistance / state.initialDistance

      // 防抖：只有缩放变化超过阈值才触发
      if (Math.abs(scale - state.lastScale) > pinchSensitivity) {
        handlers.onPinch(scale)
        state.lastScale = scale
      }
    }
  }

  /**
   * 处理触摸结束
   */
  function handleTouchEnd(event: TouchEvent) {
    // 清除长按定时器
    if (state.longPressTimer) {
      clearTimeout(state.longPressTimer)
      state.longPressTimer = null
    }

    if (!state.startPoint || !state.currentPoint) {
      resetState()
      return
    }

    const deltaX = state.currentPoint.x - state.startPoint.x
    const deltaY = state.currentPoint.y - state.startPoint.y
    const deltaTime = Date.now() - state.startTime
    const absDeltaX = Math.abs(deltaX)
    const absDeltaY = Math.abs(deltaY)

    // 检测滑动
    if (absDeltaX > swipeThreshold || absDeltaY > swipeThreshold) {
      // 判断是水平滑动还是垂直滑动
      if (absDeltaX > absDeltaY) {
        // 水平滑动
        if (deltaX > 0 && handlers.onSwipeRight) {
          handlers.onSwipeRight()
        } else if (deltaX < 0 && handlers.onSwipeLeft) {
          handlers.onSwipeLeft()
        }
      } else {
        // 垂直滑动
        if (deltaY > 0 && handlers.onSwipeDown) {
          handlers.onSwipeDown()
        } else if (deltaY < 0 && handlers.onSwipeUp) {
          handlers.onSwipeUp()
        }
      }
    }
    // 检测点击（移动距离小且时间短）
    else if (absDeltaX < 10 && absDeltaY < 10 && deltaTime < 300) {
      if (handlers.onTap) {
        handlers.onTap()
      }
    }

    resetState()
  }

  /**
   * 重置状态
   */
  function resetState() {
    state.startPoint = null
    state.currentPoint = null
    state.startTime = 0
    if (state.longPressTimer) {
      clearTimeout(state.longPressTimer)
      state.longPressTimer = null
    }
    state.initialDistance = 0
    state.lastScale = 1
  }

  /**
   * 绑定事件监听
   */
  let isBound = false

  function bind() {
    if (!elementRef.value || isBound) return

    const element = elementRef.value
    element.addEventListener('touchstart', handleTouchStart, { passive })
    element.addEventListener('touchmove', handleTouchMove, { passive })
    element.addEventListener('touchend', handleTouchEnd, { passive })
    element.addEventListener('touchcancel', resetState, { passive })

    isBound = true
  }

  /**
   * 解绑事件监听
   */
  function unbind() {
    if (!elementRef.value || !isBound) return

    const element = elementRef.value
    element.removeEventListener('touchstart', handleTouchStart)
    element.removeEventListener('touchmove', handleTouchMove)
    element.removeEventListener('touchend', handleTouchEnd)
    element.removeEventListener('touchcancel', resetState)

    resetState()
    isBound = false
  }

  onMounted(() => {
    bind()
  })

  onUnmounted(() => {
    unbind()
  })

  return {
    bind,
    unbind
  }
}

/**
 * 滑动手势 Hook（简化版）
 * 只处理左右滑动
 */
export function useSwipeGestures(
  elementRef: Ref<HTMLElement | undefined>,
  onSwipeLeft: () => void,
  onSwipeRight: () => void,
  threshold: number = 50
) {
  return useTouchGestures({
    elementRef,
    handlers: {
      onSwipeLeft,
      onSwipeRight
    },
    swipeThreshold: threshold
  })
}

/**
 * 缩放手势 Hook（简化版）
 * 只处理双指缩放
 */
export function usePinchGesture(
  elementRef: Ref<HTMLElement | undefined>,
  onPinch: (scale: number) => void,
  sensitivity: number = 0.001
) {
  return useTouchGestures({
    elementRef,
    handlers: {
      onPinch
    },
    pinchSensitivity: sensitivity
  })
}

/**
 * 长按手势 Hook（简化版）
 */
export function useLongPress(
  elementRef: Ref<HTMLElement | undefined>,
  onLongPress: () => void,
  delay: number = 500
) {
  return useTouchGestures({
    elementRef,
    handlers: {
      onLongPress
    },
    longPressDelay: delay
  })
}

/**
 * 拖拽手势 Hook
 * 用于处理触摸拖拽
 */
export function useTouchDrag(
  elementRef: Ref<HTMLElement | undefined>,
  options: {
    onDragStart?: (point: TouchPoint) => void
    onDragMove?: (delta: { x: number; y: number }) => void
    onDragEnd?: (delta: { x: number; y: number }) => void
    threshold?: number
  }
) {
  const { onDragStart, onDragMove, onDragEnd, threshold = 5 } = options

  const isDragging = ref(false)
  const startPoint = ref<TouchPoint | null>(null)
  const currentPoint = ref<TouchPoint | null>(null)

  function handleTouchStart(event: TouchEvent) {
    if (!event.touches || event.touches.length === 0) return

    const touch = event.touches[0]
    startPoint.value = getTouchPoint(touch)
    currentPoint.value = startPoint.value
  }

  function handleTouchMove(event: TouchEvent) {
    if (!startPoint.value || !event.touches || event.touches.length === 0) return

    const touch = event.touches[0]
    currentPoint.value = getTouchPoint(touch)

    const deltaX = currentPoint.value.x - startPoint.value.x
    const deltaY = currentPoint.value.y - startPoint.value.y
    const distance = Math.sqrt(deltaX * deltaX + deltaY * deltaY)

    // 超过阈值才开始拖拽
    if (!isDragging.value && distance > threshold) {
      isDragging.value = true
      if (onDragStart) {
        onDragStart(startPoint.value)
      }
    }

    if (isDragging.value && onDragMove) {
      onDragMove({ x: deltaX, y: deltaY })
    }
  }

  function handleTouchEnd(event: TouchEvent) {
    if (!isDragging.value || !startPoint.value || !currentPoint.value) {
      startPoint.value = null
      currentPoint.value = null
      return
    }

    const deltaX = currentPoint.value.x - startPoint.value.x
    const deltaY = currentPoint.value.y - startPoint.value.y

    if (onDragEnd) {
      onDragEnd({ x: deltaX, y: deltaY })
    }

    isDragging.value = false
    startPoint.value = null
    currentPoint.value = null
  }

  function getTouchPoint(touch: Touch): TouchPoint {
    return {
      x: touch.clientX,
      y: touch.clientY
    }
  }

  onMounted(() => {
    const element = elementRef.value
    if (!element) return

    element.addEventListener('touchstart', handleTouchStart, { passive: true })
    element.addEventListener('touchmove', handleTouchMove, { passive: true })
    element.addEventListener('touchend', handleTouchEnd, { passive: true })
    element.addEventListener('touchcancel', handleTouchEnd, { passive: true })
  })

  onUnmounted(() => {
    const element = elementRef.value
    if (!element) return

    element.removeEventListener('touchstart', handleTouchStart)
    element.removeEventListener('touchmove', handleTouchMove)
    element.removeEventListener('touchend', handleTouchEnd)
    element.removeEventListener('touchcancel', handleTouchEnd)
  })

  return {
    isDragging
  }
}

/**
 * 防止页面滚动（用于全屏手势）
 */
export function usePreventScroll(
  elementRef: Ref<HTMLElement | undefined>,
  enabled: Ref<boolean> = ref(true)
) {
  function handleTouchMove(event: TouchEvent) {
    if (enabled.value) {
      event.preventDefault()
    }
  }

  onMounted(() => {
    const element = elementRef.value
    if (!element) return

    element.addEventListener('touchmove', handleTouchMove, { passive: false })
  })

  onUnmounted(() => {
    const element = elementRef.value
    if (!element) return

    element.removeEventListener('touchmove', handleTouchMove)
  })
}
