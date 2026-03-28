/**
 * 通用工具函数
 */

/**
 * 获取完整的资源 URL
 * 将相对路径转换为完整 URL
 * @param {string} path - 相对路径或完整 URL
 * @returns {string} 完整的资源 URL
 */
export function getAssetUrl(path) {
  // 处理空值情况
  if (!path || path === 'null' || path === '') return ''
  // 已经是完整 URL
  if (path.startsWith('http://') || path.startsWith('https://')) {
    return path
  }
  const baseURL = window.location.origin
  return path.startsWith('/') ? `${baseURL}${path}` : `${baseURL}/${path}`
}
