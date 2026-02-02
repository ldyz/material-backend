/**
 * 响应式断点检测 Composable
 * 提供响应式的断点状态和设备检测
 */

import { ref, computed, onMounted, onUnmounted, Ref } from 'vue'
import { BREAKPOINT_VALUES, type Breakpoint, DeviceType } from '@/config/breakpoints'

/**
 * 断点检测 Composable 返回值
 */
export interface UseBreakpointReturn {
  // 当前断点
  current: Ref<Breakpoint>
  // 屏幕宽度
  width: Ref<number>
  // 屏幕高度
  height: Ref<number>
  // 设备类型
  isMobile: Ref<boolean>
  isTablet: Ref<boolean>
  isDesktop: Ref<boolean>
  // 触摸支持
  isTouch: Ref<boolean>
  // 高DPI
  isHighDPI: Ref<boolean>
  // 断点检测方法
  isGreater: (bp: Breakpoint) => boolean
  isLess: (bp: Breakpoint) => boolean
  isBetween: (min: Breakpoint, max: Breakpoint) => boolean
  // 当前是否匹配指定断点
  is: (bp: Breakpoint) => boolean
}

/**
 * 响应式断点检测 Hook
 * @returns 断点状态和检测方法
 */
export function useBreakpoint(): UseBreakpointReturn {
  // 屏幕尺寸
  const width = ref(0)
  const height = ref(0)

  // 当前断点
  const current = ref<Breakpoint>('md')

  // 设备类型
  const isMobile = ref(false)
  const isTablet = ref(false)
  const isDesktop = ref(false)

  // 特性检测
  const isTouch = ref(false)
  const isHighDPI = ref(false)

  /**
   * 更新屏幕尺寸和断点
   */
  const updateBreakpoint = () => {
    if (typeof window === 'undefined') return

    width.value = window.innerWidth
    height.value = window.innerHeight

    // 更新当前断点
    current.value = getCurrentBreakpoint(width.value)

    // 更新设备类型
    isMobile.value = width.value < BREAKPOINT_VALUES.md
    isTablet.value = width.value >= BREAKPOINT_VALUES.md && width.value < BREAKPOINT_VALUES.lg
    isDesktop.value = width.value >= BREAKPOINT_VALUES.lg

    // 更新特性
    isTouch.value = 'ontouchstart' in window || navigator.maxTouchPoints > 0
    isHighDPI.value = window.devicePixelRatio > 1
  }

  /**
   * 根据宽度获取当前断点
   */
  function getCurrentBreakpoint(w: number): Breakpoint {
    if (w < BREAKPOINT_VALUES.sm) return 'xs'
    if (w < BREAKPOINT_VALUES.md) return 'sm'
    if (w < BREAKPOINT_VALUES.lg) return 'md'
    if (w < BREAKPOINT_VALUES.xl) return 'lg'
    if (w < BREAKPOINT_VALUES['2xl']) return 'xl'
    return '2xl'
  }

  /**
   * 检测当前宽度是否大于指定断点
   */
  function isGreater(bp: Breakpoint): boolean {
    return width.value > BREAKPOINT_VALUES[bp]
  }

  /**
   * 检测当前宽度是否小于指定断点
   */
  function isLess(bp: Breakpoint): boolean {
    return width.value < BREAKPOINT_VALUES[bp]
  }

  /**
   * 检测当前宽度是否在两个断点之间
   */
  function isBetween(min: Breakpoint, max: Breakpoint): boolean {
    return width.value >= BREAKPOINT_VALUES[min] && width.value < BREAKPOINT_VALUES[max]
  }

  /**
   * 检测当前是否匹配指定断点
   */
  function is(bp: Breakpoint): boolean {
    return current.value === bp
  }

  // 监听窗口大小变化
  let resizeObserver: ResizeObserver | null = null

  onMounted(() => {
    updateBreakpoint()

    // 使用 ResizeObserver 监听更准确
    if (typeof window !== 'undefined' && window.ResizeObserver) {
      resizeObserver = new ResizeObserver(() => {
        updateBreakpoint()
      })
      resizeObserver.observe(document.documentElement)
    } else {
      // 降级使用 resize 事件
      window.addEventListener('resize', updateBreakpoint)
      window.addEventListener('orientationchange', updateBreakpoint)
    }
  })

  onUnmounted(() => {
    if (resizeObserver) {
      resizeObserver.disconnect()
    } else {
      window.removeEventListener('resize', updateBreakpoint)
      window.removeEventListener('orientationchange', updateBreakpoint)
    }
  })

  return {
    current,
    width,
    height,
    isMobile,
    isTablet,
    isDesktop,
    isTouch,
    isHighDPI,
    isGreater,
    isLess,
    isBetween,
    is
  }
}

/**
 * 媒体查询 Hook
 * 监听特定的媒体查询条件
 * @param query 媒体查询字符串
 * @returns 是否匹配
 */
export function useMediaQuery(query: string): Ref<boolean> {
  const matches = ref(false)

  let mediaQuery: MediaQueryList | null = null

  const updateMatches = () => {
    if (mediaQuery) {
      matches.value = mediaQuery.matches
    }
  }

  onMounted(() => {
    if (typeof window !== 'undefined' && window.matchMedia) {
      mediaQuery = window.matchMedia(query)
      updateMatches()

      // 现代浏览器使用 addEventListener
      if (mediaQuery.addEventListener) {
        mediaQuery.addEventListener('change', updateMatches)
      } else {
        // 旧版浏览器降级
        mediaQuery.addListener(updateMatches)
      }
    }
  })

  onUnmounted(() => {
    if (mediaQuery) {
      if (mediaQuery.removeEventListener) {
        mediaQuery.removeEventListener('change', updateMatches)
      } else {
        mediaQuery.removeListener(updateMatches)
      }
    }
  })

  return matches
}

/**
 * 响应式值 Hook
 * 根据断点返回不同的值
 * @param values 断点值映射
 * @param defaultValue 默认值
 * @returns 当前断点对应的值
 */
export function useResponsiveValue<T>(
  values: Ref<Partial<Record<Breakpoint, T>>>,
  defaultValue: Ref<T> | T
): Ref<T> {
  const { current } = useBreakpoint()

  return computed(() => {
    const breakpointOrder: Breakpoint[] = ['2xl', 'xl', 'lg', 'md', 'sm', 'xs']

    // 从当前断点向下查找第一个有值的
    const currentIndex = breakpointOrder.indexOf(current.value)
    for (let i = currentIndex; i < breakpointOrder.length; i++) {
      const bp = breakpointOrder[i]
      const value = values.value[bp]
      if (value !== undefined) {
        return value
      }
    }

    // 返回默认值
    return defaultValue instanceof Ref ? defaultValue.value : defaultValue
  })
}

/**
 * 防抖函数
 * 用于优化频繁触发的事件
 */
function debounce<T extends (...args: any[]) => any>(
  fn: T,
  delay: number
): (...args: Parameters<T>) => void {
  let timeoutId: ReturnType<typeof setTimeout> | null = null

  return function (this: any, ...args: Parameters<T>) {
    if (timeoutId) {
      clearTimeout(timeoutId)
    }

    timeoutId = setTimeout(() => {
      fn.apply(this, args)
      timeoutId = null
    }, delay)
  }
}

/**
 * 节流函数
 * 用于限制函数的执行频率
 */
function throttle<T extends (...args: any[]) => any>(
  fn: T,
  limit: number
): (...args: Parameters<T>) => void {
  let inThrottle = false

  return function (this: any, ...args: Parameters<T>) {
    if (!inThrottle) {
      fn.apply(this, args)
      inThrottle = true
      setTimeout(() => {
        inThrottle = false
      }, limit)
    }
  }
}

// 导出工具函数
export const breakpointUtils = {
  debounce,
  throttle,
  getCurrentBreakpoint: DeviceType.getCurrentBreakpoint.bind(DeviceType),
  isMobile: DeviceType.isMobile.bind(DeviceType),
  isTablet: DeviceType.isTablet.bind(DeviceType),
  isDesktop: DeviceType.isDesktop.bind(DeviceType),
  isTouchDevice: DeviceType.isTouchDevice.bind(DeviceType),
  isHighDPI: DeviceType.isHighDPI.bind(DeviceType)
}
