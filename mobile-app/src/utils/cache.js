/**
 * 数据缓存工具
 * 用于减少 API 请求，提升用户体验
 */

// 内存缓存
const cache = new Map()

// 默认缓存过期时间（5分钟）
const DEFAULT_TTL = 5 * 60 * 1000

/**
 * 获取缓存数据
 * @param {string} key - 缓存键
 * @param {Function} fetcher - 数据获取函数
 * @param {number} ttl - 缓存过期时间（毫秒）
 * @returns {Promise} 缓存数据或新获取的数据
 */
export function getCachedData(key, fetcher, ttl = DEFAULT_TTL) {
  const cached = cache.get(key)

  // 检查缓存是否存在且未过期
  if (cached && Date.now() - cached.time < ttl) {
    return Promise.resolve(cached.data)
  }

  // 获取新数据并缓存
  return fetcher().then(data => {
    cache.set(key, { data, time: Date.now() })
    return data
  })
}

/**
 * 设置缓存数据
 * @param {string} key - 缓存键
 * @param {any} data - 要缓存的数据
 */
export function setCache(key, data) {
  cache.set(key, { data, time: Date.now() })
}

/**
 * 清除指定缓存
 * @param {string} key - 缓存键，如果不传则清除所有缓存
 */
export function clearCache(key) {
  if (key) {
    cache.delete(key)
  } else {
    cache.clear()
  }
}

/**
 * 检查缓存是否存在且有效
 * @param {string} key - 缓存键
 * @param {number} ttl - 缓存过期时间（毫秒）
 * @returns {boolean}
 */
export function hasCache(key, ttl = DEFAULT_TTL) {
  const cached = cache.get(key)
  return cached && Date.now() - cached.time < ttl
}

/**
 * 获取缓存的原始数据（不检查过期）
 * @param {string} key - 缓存键
 * @returns {any|null}
 */
export function getCacheRaw(key) {
  const cached = cache.get(key)
  return cached ? cached.data : null
}

/**
 * 生成缓存键
 * @param {string} prefix - 前缀
 * @param {object} params - 参数对象
 * @returns {string} 缓存键
 */
export function getCacheKey(prefix, params = {}) {
  const sortedParams = Object.keys(params)
    .sort()
    .map(key => `${key}=${JSON.stringify(params[key])}`)
    .join('&')
  return sortedParams ? `${prefix}:${sortedParams}` : prefix
}

/**
 * 批量清除以指定前缀开头的缓存
 * @param {string} prefix - 缓存键前缀
 */
export function clearCacheByPrefix(prefix) {
  for (const key of cache.keys()) {
    if (key.startsWith(prefix)) {
      cache.delete(key)
    }
  }
}

/**
 * 获取缓存统计信息
 * @returns {object} 缓存统计
 */
export function getCacheStats() {
  return {
    size: cache.size,
    keys: Array.from(cache.keys())
  }
}

// 导出缓存实例（高级用法）
export { cache }

export default {
  getCachedData,
  setCache,
  clearCache,
  hasCache,
  getCacheRaw,
  getCacheKey,
  clearCacheByPrefix,
  getCacheStats
}
