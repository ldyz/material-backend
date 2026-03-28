import request from '@/utils/request'

/**
 * 获取物资列表
 * @param {Object} params - 查询参数
 * @param {string} params.search - 搜索关键词
 * @param {string} params.category - 分类筛选
 * @param {number} params.page - 页码
 * @param {number} params.page_size - 每页数量
 */
export function getMaterials(params) {
  return request({
    url: '/material/materials',
    method: 'GET',
    params
  })
}

/**
 * 获取物资详情
 * @param {number} id - 物资ID
 */
export function getMaterialDetail(id) {
  return request({
    url: `/material/materials/${id}`,
    method: 'GET'
  })
}
