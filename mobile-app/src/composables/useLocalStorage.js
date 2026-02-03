import { ref, watch } from 'vue'
import { storage } from '@/utils/storage'

/**
 * useLocalStorage - 响应式LocalStorage Composable
 *
 * @param {string} key - 存储键名
 * @param {any} defaultValue - 默认值
 * @returns {Object} { value, remove, update }
 */
export function useLocalStorage(key, defaultValue = null) {
  const storedValue = storage.get(key)

  const value = ref(storedValue !== null ? storedValue : defaultValue)

  // 监听值变化，自动保存
  watch(
    value,
    (newValue) => {
      if (newValue === null || newValue === undefined) {
        storage.remove(key)
      } else {
        storage.set(key, newValue)
      }
    },
    { deep: true }
  )

  // 更新值
  const update = (newValue) => {
    value.value = newValue
  }

  // 删除值
  const remove = () => {
    value.value = defaultValue
    storage.remove(key)
  }

  return {
    value,
    update,
    remove,
  }
}

/**
 * useSessionStorage - 响应式SessionStorage Composable
 *
 * @param {string} key - 存储键名
 * @param {any} defaultValue - 默认值
 * @returns {Object} { value, remove, update }
 */
export function useSessionStorage(key, defaultValue = null) {
  const getStoredValue = () => {
    try {
      const item = sessionStorage.getItem(key)
      return item ? JSON.parse(item) : defaultValue
    } catch {
      return defaultValue
    }
  }

  const value = ref(getStoredValue())

  // 监听值变化，自动保存
  watch(
    value,
    (newValue) => {
      if (newValue === null || newValue === undefined) {
        sessionStorage.removeItem(key)
      } else {
        sessionStorage.setItem(key, JSON.stringify(newValue))
      }
    },
    { deep: true }
  )

  // 更新值
  const update = (newValue) => {
    value.value = newValue
  }

  // 删除值
  const remove = () => {
    value.value = defaultValue
    sessionStorage.removeItem(key)
  }

  return {
    value,
    update,
    remove,
  }
}
