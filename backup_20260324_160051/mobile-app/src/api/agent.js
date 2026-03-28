import request from '@/utils/request'

/**
 * AI 智能体 API
 */
export const agentApi = {
  /**
   * 语音对话接口
   * 发送语音文件，返回语音识别结果和 AI 回复
   *
   * @param {FormData} formData - 包含 audio 文件的 FormData
   * @returns {Promise} 返回 { response: string, transcript?: string }
   */
  voiceChat(formData) {
    return request({
      url: '/agent/voice-chat',
      method: 'POST',
      data: formData,
      headers: {
        'Content-Type': 'multipart/form-data'
      },
      timeout: 60000 // 60秒超时
    })
  },

  /**
   * 文本对话接口
   * @param {string} message - 用户消息
   * @param {Array} history - 对话历史
   * @returns {Promise}
   */
  chat(message, history = []) {
    return request({
      url: '/agent/chat',
      method: 'POST',
      data: {
        message,
        conversation_history: history
      },
      timeout: 60000
    })
  }
}

export default agentApi
