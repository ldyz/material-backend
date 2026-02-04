import request from '@/utils/request'

/**
 * AI 自然语言查询
 * @param {string} question - 自然语言问题
 * @returns {Promise}
 */
export function analyzeQuestion(question) {
  return request.post('/system/ai/analyze', { question })
}

/**
 * 获取数据洞察
 * @param {Object} params - 查询参数
 * @returns {Promise}
 */
export function getInsights(params) {
  return request.get('/system/ai/insights', { params })
}

/**
 * 获取智能推荐 (后端路由是 /system/ai/suggestions)
 * @param {string} type - 推荐类型
 * @returns {Promise}
 */
export function getRecommendations(type) {
  return request.get('/system/ai/suggestions', { params: { type } })
}

/**
 * 获取 AI 查询历史
 * @param {Object} params - 查询参数
 * @returns {Promise}
 */
export function getAIHistory(params) {
  return request.get('/system/ai/history', { params })
}
