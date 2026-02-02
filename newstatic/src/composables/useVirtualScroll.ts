/**
 * 虚拟滚动 Composable
 * 优化长列表渲染性能，只渲染可见区域的项
 */

import { ref, computed, Ref, watch, onMounted, onUnmounted } from 'vue'
import type { VirtualScrollOptions, VisibleRange } from '@/types/gantt'

/**
 * 虚拟滚动状态
 */
export interface VirtualScrollState<T = any> {
  // 可见范围
  visibleRange: Ref<VisibleRange>
  // 可见项
  visibleItems: Ref<T[]>
  // 总高度
  totalHeight: Ref<number>
  // 偏移量
  offsetY: Ref<number>
  // 滚动位置
  scrollTop: Ref<number>
  // 是否正在滚动
  isScrolling: Ref<boolean>
  // 滚动方向
  scrollDirection: Ref<'up' | 'down' | 'none'>
}

/**
 * 虚拟滚动配置选项
 */
export interface UseVirtualScrollOptions<T = any> extends VirtualScrollOptions {
  // 数据源
  items: Ref<T[]>
  // 容器引用
  containerRef?: Ref<HTMLElement | undefined>
  // 预渲染数量（额外渲染的项数）
  overscan?: number
  // 滚动节流延迟（毫秒）
  throttleDelay?: number
  // 滚动停止延迟（毫秒），用于检测滚动是否停止
  scrollEndDelay?: number
}

/**
 * 虚拟滚动 Hook
 * @param options 配置选项
 * @returns 虚拟滚动状态和方法
 */
export function useVirtualScroll<T = any>(
  options: UseVirtualScrollOptions<T>
): VirtualScrollState<T> & {
  // 滚动到指定项
  scrollToItem: (index: number, alignment?: 'start' | 'center' | 'end' | 'auto') => void
  // 滚动到指定位置
  scrollToPosition: (position: number) => void
  // 获取项的大小
  getItemSize: (index: number) => number
  // 获取项的位置
  getItemOffset: (index: number) => number
  // 刷新布局
  refresh: () => void
} {
  const {
    items,
    itemHeight,
    containerHeight,
    containerRef,
    overscan = 3,
    throttleDelay = 16, // ~60fps
    scrollEndDelay = 150
  } = options

  // 滚动位置
  const scrollTop = ref(0)
  const isScrolling = ref(false)
  const scrollDirection = ref<'up' | 'down' | 'none'>('none')

  // 上次滚动位置（用于检测方向）
  let lastScrollTop = 0
  let scrollTimer: ReturnType<typeof setTimeout> | null = null
  let throttleTimer: ReturnType<typeof setTimeout> | null = null

  // 计算总高度
  const totalHeight = computed(() => {
    return items.value.length * itemHeight
  })

  // 计算可见范围
  const visibleRange = computed<VisibleRange>(() => {
    if (items.value.length === 0) {
      return { start: 0, end: 0 }
    }

    const start = Math.floor(scrollTop.value / itemHeight)
    const visibleCount = Math.ceil(containerHeight / itemHeight)

    // 考虑 overscan
    const overscanStart = Math.max(0, start - overscan)
    const end = Math.min(items.value.length, start + visibleCount + overscan)

    return {
      start: overscanStart,
      end
    }
  })

  // 计算可见项
  const visibleItems = computed(() => {
    const { start, end } = visibleRange.value
    return items.value.slice(start, end)
  })

  // 计算偏移量
  const offsetY = computed(() => {
    return visibleRange.value.start * itemHeight
  })

  /**
   * 处理滚动事件
   */
  const handleScroll = (event: Event) => {
    const target = event.target as HTMLElement
    const newScrollTop = target.scrollTop

    // 检测滚动方向
    if (newScrollTop > lastScrollTop) {
      scrollDirection.value = 'down'
    } else if (newScrollTop < lastScrollTop) {
      scrollDirection.value = 'up'
    }
    lastScrollTop = newScrollTop

    // 标记正在滚动
    isScrolling.value = true

    // 清除之前的定时器
    if (scrollTimer) {
      clearTimeout(scrollTimer)
    }

    // 设置新的定时器检测滚动停止
    scrollTimer = setTimeout(() => {
      isScrolling.value = false
      scrollDirection.value = 'none'
    }, scrollEndDelay)

    // 更新滚动位置
    scrollTop.value = newScrollTop
  }

  /**
   * 节流滚动处理
   */
  const throttledHandleScroll = (event: Event) => {
    if (throttleTimer) {
      return
    }

    handleScroll(event)

    throttleTimer = setTimeout(() => {
      throttleTimer = null
    }, throttleDelay)
  }

  /**
   * 滚动到指定项
   */
  const scrollToItem = (
    index: number,
    alignment: 'start' | 'center' | 'end' | 'auto' = 'auto'
  ) => {
    if (index < 0 || index >= items.value.length) {
      return
    }

    const itemOffset = index * itemHeight
    let targetPosition = itemOffset

    switch (alignment) {
      case 'start':
        targetPosition = itemOffset
        break
      case 'center':
        targetPosition = itemOffset - containerHeight / 2 + itemHeight / 2
        break
      case 'end':
        targetPosition = itemOffset - containerHeight + itemHeight
        break
      case 'auto':
        // 如果项已在可见区域，不滚动
        if (itemOffset >= scrollTop.value && itemOffset + itemHeight <= scrollTop.value + containerHeight) {
          return
        }
        // 否则滚动到顶部
        targetPosition = itemOffset
        break
    }

    scrollToPosition(targetPosition)
  }

  /**
   * 滚动到指定位置
   */
  const scrollToPosition = (position: number) => {
    const maxPosition = Math.max(0, totalHeight.value - containerHeight)
    const targetPosition = Math.max(0, Math.min(position, maxPosition))

    if (containerRef?.value) {
      containerRef.value.scrollTop = targetPosition
    } else if (typeof window !== 'undefined') {
      window.scrollTo(0, targetPosition)
    }

    scrollTop.value = targetPosition
  }

  /**
   * 获取项的大小
   */
  const getItemSize = (index: number) => {
    return itemHeight
  }

  /**
   * 获取项的偏移量
   */
  const getItemOffset = (index: number) => {
    return index * itemHeight
  }

  /**
   * 刷新布局
   */
  const refresh = () => {
    if (containerRef?.value) {
      scrollTop.value = containerRef.value.scrollTop
    }
  }

  // 监听容器
  if (containerRef) {
    onMounted(() => {
      const container = containerRef.value
      if (container) {
        container.addEventListener('scroll', throttledHandleScroll, { passive: true })
        // 初始化滚动位置
        scrollTop.value = container.scrollTop
      }
    })

    onUnmounted(() => {
      const container = containerRef?.value
      if (container) {
        container.removeEventListener('scroll', throttledHandleScroll)
      }
      if (scrollTimer) {
        clearTimeout(scrollTimer)
      }
      if (throttleTimer) {
        clearTimeout(throttleTimer)
      }
    })
  }

  return {
    visibleRange,
    visibleItems,
    totalHeight,
    offsetY,
    scrollTop,
    isScrolling,
    scrollDirection,
    scrollToItem,
    scrollToPosition,
    getItemSize,
    getItemOffset,
    refresh
  }
}

/**
 * 动态高度的虚拟滚动
 * 支持不同高度的项
 */
export function useDynamicVirtualScroll<T = any>(
  options: UseVirtualScrollOptions<T> & {
    // 获取项高度的函数
    getItemHeight: (item: T, index: number) => number
    // 缓存项高度
    estimatedItemHeight?: number
  }
) {
  const {
    items,
    containerHeight,
    containerRef,
    getItemHeight,
    estimatedItemHeight = 50,
    overscan = 3
  } = options

  // 缓存项的位置信息
  const itemOffsets = ref<Map<number, number>>(new Map())
  const itemHeights = ref<Map<number, number>>(new Map())
  const totalHeight = ref(0)

  // 计算所有项的位置
  const updateItemPositions = () => {
    let offset = 0
    const newOffsets = new Map<number, number>()
    const newHeights = new Map<number, number>()

    items.value.forEach((item, index) => {
      newOffsets.set(index, offset)
      const height = getItemHeight(item, index)
      newHeights.set(index, height)
      offset += height
    })

    itemOffsets.value = newOffsets
    itemHeights.value = newHeights
    totalHeight.value = offset
  }

  // 监听数据变化
  watch(items, updateItemPositions, { immediate: true, deep: true })

  const scrollTop = ref(0)

  // 二分查找查找第一个可见项
  const findFirstVisibleIndex = () => {
    let left = 0
    let right = items.value.length - 1

    while (left <= right) {
      const mid = Math.floor((left + right) / 2)
      const offset = itemOffsets.value.get(mid) ?? 0

      if (offset < scrollTop.value) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }

    return Math.max(0, left - 1)
  }

  // 计算可见范围
  const visibleRange = computed<VisibleRange>(() => {
    if (items.value.length === 0) {
      return { start: 0, end: 0 }
    }

    const startIndex = findFirstVisibleIndex()
    let endIndex = startIndex

    // 找到所有可见项
    for (let i = startIndex; i < items.value.length; i++) {
      const offset = itemOffsets.value.get(i) ?? 0
      const height = itemHeights.value.get(i) ?? estimatedItemHeight

      if (offset < scrollTop.value + containerHeight) {
        endIndex = i
      } else {
        break
      }
    }

    return {
      start: Math.max(0, startIndex - overscan),
      end: Math.min(items.value.length, endIndex + overscan + 1)
    }
  })

  const visibleItems = computed(() => {
    const { start, end } = visibleRange.value
    return items.value.slice(start, end)
  })

  const offsetY = computed(() => {
    const { start } = visibleRange.value
    return itemOffsets.value.get(start) ?? 0
  })

  return {
    visibleRange,
    visibleItems,
    totalHeight: computed(() => totalHeight.value),
    offsetY,
    scrollTop,
    isScrolling: ref(false),
    scrollDirection: ref<'up' | 'down' | 'none'>('none'),
    scrollToItem: (index: number) => {
      const offset = itemOffsets.value.get(index)
      if (offset !== undefined && containerRef?.value) {
        containerRef.value.scrollTop = offset
      }
    },
    scrollToPosition: (position: number) => {
      if (containerRef?.value) {
        containerRef.value.scrollTop = position
      }
    },
    getItemSize: (index: number) => itemHeights.value.get(index) ?? estimatedItemHeight,
    getItemOffset: (index: number) => itemOffsets.value.get(index) ?? 0,
    refresh: updateItemPositions
  }
}
