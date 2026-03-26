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
  if (!path) return ''
  if (path.startsWith('http://') || path.startsWith('https://')) {
    return path
  }
  const baseURL = window.location.origin
  return path.startsWith('/') ? `${baseURL}${path}` : `${baseURL}/${path}`
}
