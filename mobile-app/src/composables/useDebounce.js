import { ref, watch, onUnmounted } from 'vue'

/**
 * useDebounce - 防抖值Composable
 *
 * @param {any} value - 需要防抖的值
 * @param {number} delay - 延迟时间（毫秒）
 * @returns {Object} { debouncedValue, flush }
 */
export function useDebounce(value, delay = 300) {
  const debouncedValue = ref(value.value)
  let timeout = null

  const flush = () => {
    if (timeout) {
      clearTimeout(timeout)
      timeout = null
    }
    debouncedValue.value = value.value
  }

  watch(value, (newValue) => {
    if (timeout) {
      clearTimeout(timeout)
    }

    timeout = setTimeout(() => {
      debouncedValue.value = newValue
    }, delay)
  }, { immediate: true })

  onUnmounted(() => {
    if (timeout) {
      clearTimeout(timeout)
    }
  })

  return { debouncedValue, flush }
}

/**
 * useDebounceFn - 防抖函数Composable
 *
 * @param {Function} fn - 需要防抖的函数
 * @param {number} delay - 延迟时间（毫秒）
 * @returns {Function} 防抖后的函数
 */
export function useDebounceFn(fn, delay = 300) {
  let timeout = null

  const debouncedFn = (...args) => {
    if (timeout) {
      clearTimeout(timeout)
    }

    return new Promise((resolve) => {
      timeout = setTimeout(() => {
        const result = fn(...args)
        resolve(result)
      }, delay)
    })
  }

  // 取消防抖
  const cancel = () => {
    if (timeout) {
      clearTimeout(timeout)
      timeout = null
    }
  }

  // 立即执行
  const flush = (...args) => {
    cancel()
    return fn(...args)
  }

  onUnmounted(() => {
    cancel()
  })

  debouncedFn.cancel = cancel
  debouncedFn.flush = flush

  return debouncedFn
}
