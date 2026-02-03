import { ref, shallowRef } from 'vue'

/**
 * useAPI - 统一的API调用处理Composable
 * 自动处理loading状态、错误处理、数据存储
 *
 * @param {Function} apiFunction - API函数
 * @param {Object} options - 配置选项
 * @returns {Object} { data, loading, error, execute, reset }
 */
export function useAPI(apiFunction, options = {}) {
  const {
    immediate = false,
    onSuccess = null,
    onError = null,
    initialData = null,
    showToast = true,
  } = options

  const data = shallowRef(initialData)
  const loading = ref(false)
  const error = ref(null)

  // 执行API调用
  async function execute(...args) {
    loading.value = true
    error.value = null

    try {
      const response = await apiFunction(...args)

      if (response.success) {
        data.value = response.data
        if (onSuccess) {
          await onSuccess(response.data, response)
        }
        return response.data
      } else {
        throw new Error(response.message || '请求失败')
      }
    } catch (err) {
      error.value = err
      if (onError) {
        await onError(err)
      }
      if (showToast) {
        showErrorToast(err)
      }
      throw err
    } finally {
      loading.value = false
    }
  }

  // 重置状态
  function reset() {
    data.value = initialData
    loading.value = false
    error.value = null
  }

  // 立即执行
  if (immediate && apiFunction) {
    execute()
  }

  return {
    data,
    loading,
    error,
    execute,
    reset,
  }
}

/**
 * 显示错误提示
 */
function showErrorToast(error) {
  // 动态导入vant避免循环依赖
  import('vant').then(({ showToast }) => {
    const message = error?.response?.data?.message ||
                   error?.message ||
                   '操作失败，请稍后重试'
    showToast({
      message,
      type: 'fail',
      duration: 2000,
    })
  })
}

/**
 * useAPIList - 用于列表数据的API调用
 * 自动处理分页、刷新、加载更多
 */
export function useAPIList(apiFunction, options = {}) {
  const {
    immediate = false,
    pageSize = 20,
    onSuccess = null,
    onError = null,
  } = options

  const list = ref([])
  const loading = ref(false)
  const loadingMore = ref(false)
  const refreshing = ref(false)
  const finished = ref(false)
  const error = ref(null)
  const currentPage = ref(1)
  const total = ref(0)

  // 加载列表数据
  async function load(params = {}, isRefresh = false, isLoadMore = false) {
    if (isRefresh) {
      refreshing.value = true
      list.value = []
      currentPage.value = 1
      finished.value = false
    } else if (isLoadMore) {
      loadingMore.value = true
    } else {
      loading.value = true
    }

    error.value = null

    try {
      const response = await apiFunction({
        page: currentPage.value,
        page_size: pageSize,
        ...params,
      })

      if (response.success) {
        const { data = [], total: totalCount = 0 } = response.data

        if (isRefresh) {
          list.value = data
        } else if (isLoadMore) {
          list.value.push(...data)
        } else {
          list.value = data
        }

        total.value = totalCount

        // 检查是否还有更多数据
        if (list.value.length >= total.value) {
          finished.value = true
        } else {
          currentPage.value++
        }

        if (onSuccess) {
          await onSuccess(list.value, response)
        }

        return list.value
      } else {
        throw new Error(response.message || '加载失败')
      }
    } catch (err) {
      error.value = err
      if (onError) {
        await onError(err)
      }
      throw err
    } finally {
      loading.value = false
      loadingMore.value = false
      refreshing.value = false
    }
  }

  // 刷新
  async function refresh(params = {}) {
    await load(params, true, false)
  }

  // 加载更多
  async function loadMore(params = {}) {
    if (loadingMore.value || finished.value) return
    await load(params, false, true)
  }

  // 重置
  function reset() {
    list.value = []
    loading.value = false
    loadingMore.value = false
    refreshing.value = false
    finished.value = false
    error.value = null
    currentPage.value = 1
    total.value = 0
  }

  // 立即执行
  if (immediate && apiFunction) {
    load()
  }

  return {
    list,
    loading,
    loadingMore,
    refreshing,
    finished,
    error,
    total,
    currentPage,
    load,
    refresh,
    loadMore,
    reset,
  }
}
