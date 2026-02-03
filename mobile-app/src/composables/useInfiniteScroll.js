import { ref, onMounted, onUnmounted } from 'vue'
import { useDebounceFn } from './useDebounce'

/**
 * useInfiniteScroll - 无限滚动Composable
 *
 * @param {Function} loadMore - 加载更多的回调函数
 * @param {Object} options - 配置选项
 * @returns {Object} { hasMore, loading, loadMore, reset }
 */
export function useInfiniteScroll(loadMore, options = {}) {
  const {
    threshold = 100, // 距离底部多少像素触发加载
    immediate = true,
    disabled = false,
  } = options

  const hasMore = ref(true)
  const loading = ref(false)
  const elementRef = ref(null)

  // 检查是否接近底部
  const checkBottom = useDebounceFn(() => {
    if (disabled || loading.value || !hasMore.value || !elementRef.value) {
      return
    }

    const { scrollTop, scrollHeight, clientHeight } = elementRef.value
    const distanceToBottom = scrollHeight - scrollTop - clientHeight

    if (distanceToBottom <= threshold) {
      handleLoadMore()
    }
  }, 100)

  // 处理加载更多
  async function handleLoadMore() {
    if (loading.value || !hasMore.value || disabled) {
      return
    }

    loading.value = true
    try {
      const result = await loadMore()
      // 如果返回的数据少于页面大小，说明没有更多数据了
      if (result === false || (Array.isArray(result) && result.length === 0)) {
        hasMore.value = false
      }
    } catch (error) {
      console.error('加载更多失败:', error)
      hasMore.value = false
    } finally {
      loading.value = false
    }
  }

  // 绑定滚动事件
  function bindScroll() {
    if (elementRef.value) {
      elementRef.value.addEventListener('scroll', checkBottom, { passive: true })
    }
  }

  // 解绑滚动事件
  function unbindScroll() {
    if (elementRef.value) {
      elementRef.value.removeEventListener('scroll', checkBottom)
    }
  }

  // 重置状态
  function reset() {
    hasMore.value = true
    loading.value = false
  }

  // 设置元素引用
  function setElement(el) {
    if (elementRef.value) {
      unbindScroll()
    }
    elementRef.value = el
    if (el) {
      bindScroll()
    }
  }

  onMounted(() => {
    if (elementRef.value && immediate) {
      bindScroll()
    }
  })

  onUnmounted(() => {
    unbindScroll()
  })

  return {
    hasMore,
    loading,
    elementRef: setElement,
    loadMore: handleLoadMore,
    reset,
  }
}
