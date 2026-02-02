/**
 * 响应式断点配置
 * 用于检测当前设备类型和屏幕尺寸
 */

/**
 * 断点配置
 * 使用 Tailwind CSS 标准断点值
 */
export const BREAKPOINTS = {
  xs: '320px',   // 超小屏手机
  sm: '375px',   // 小屏手机
  md: '768px',   // 平板竖屏
  lg: '1024px',  // 平板横屏
  xl: '1280px',  // 桌面
  '2xl': '1536px' // 大屏桌面
} as const

/**
 * 断点类型
 */
export type Breakpoint = keyof typeof BREAKPOINTS

/**
 * 断点像素值（数值类型）
 */
export const BREAKPOINT_VALUES: Record<Breakpoint, number> = {
  xs: 320,
  sm: 375,
  md: 768,
  lg: 1024,
  xl: 1280,
  '2xl': 1536
}

/**
 * 断点媒体查询字符串
 */
export const BREAKPOINT_QUERIES: Record<Breakpoint, string> = {
  xs: '(min-width: 320px)',
  sm: '(min-width: 375px)',
  md: '(min-width: 768px)',
  lg: '(min-width: 1024px)',
  xl: '(min-width: 1280px)',
  '2xl': '(min-width: 1536px)'
}

/**
 * 设备类型检测工具
 */
export const DeviceType = {
  /**
   * 检测是否为移动设备
   */
  isMobile(): boolean {
    if (typeof window === 'undefined') return false
    return window.innerWidth < BREAKPOINT_VALUES.md
  },

  /**
   * 检测是否为平板设备
   */
  isTablet(): boolean {
    if (typeof window === 'undefined') return false
    const width = window.innerWidth
    return width >= BREAKPOINT_VALUES.md && width < BREAKPOINT_VALUES.lg
  },

  /**
   * 检测是否为桌面设备
   */
  isDesktop(): boolean {
    if (typeof window === 'undefined') return false
    return window.innerWidth >= BREAKPOINT_VALUES.lg
  },

  /**
   * 获取当前断点
   */
  getCurrentBreakpoint(): Breakpoint {
    if (typeof window === 'undefined') return 'md'

    const width = window.innerWidth
    if (width < BREAKPOINT_VALUES.sm) return 'xs'
    if (width < BREAKPOINT_VALUES.md) return 'sm'
    if (width < BREAKPOINT_VALUES.lg) return 'md'
    if (width < BREAKPOINT_VALUES.xl) return 'lg'
    if (width < BREAKPOINT_VALUES['2xl']) return 'xl'
    return '2xl'
  },

  /**
   * 检测是否支持触摸
   */
  isTouchDevice(): boolean {
    if (typeof window === 'undefined') return false
    return 'ontouchstart' in window || navigator.maxTouchPoints > 0
  },

  /**
   * 检测是否为高 DPI 设备
   */
  isHighDPI(): boolean {
    if (typeof window === 'undefined') return false
    return window.devicePixelRatio && window.devicePixelRatio > 1
  }
}

/**
 * 响应式值配置
 * 根据断点返回不同的值
 */
export function getResponsiveValue<T>(values: Partial<Record<Breakpoint, T>>, defaultValue: T): T {
  const currentBreakpoint = DeviceType.getCurrentBreakpoint()
  const breakpointOrder: Breakpoint[] = ['2xl', 'xl', 'lg', 'md', 'sm', 'xs']

  // 从当前断点向下查找第一个有值的
  for (let i = breakpointOrder.indexOf(currentBreakpoint); i < breakpointOrder.length; i++) {
    const bp = breakpointOrder[i]
    if (values[bp] !== undefined) {
      return values[bp]!
    }
  }

  return defaultValue
}
