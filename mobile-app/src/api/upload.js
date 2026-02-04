import request from '@/utils/request'

/**
 * 文件上传 API (移动端)
 *
 * 提供图片、文件上传等功能
 *
 * @module Upload API
 */

/**
 * 上传单张图片
 * @param {FormData} formData - 包含file字段的表单数据
 * @returns {Promise} 返回上传后的图片URL
 */
export function uploadImage(formData) {
  return request.post('/upload/image', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

/**
 * 批量上传图片
 * @param {FormData} formData - 包含多个file字段的表单数据
 * @returns {Promise} 返回上传后的图片URL数组
 */
export function uploadImages(formData) {
  return request.post('/upload/images', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}
