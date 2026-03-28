<template>
  <div class="ai-chat">
    <div class="chat-header">
      <h3>AI 助手</h3>
      <el-tag v-if="loading" type="info">思考中...</el-tag>
    </div>

    <div class="chat-messages" ref="messagesContainer">
      <div
        v-for="(msg, index) in messages"
        :key="index"
        :class="['message', msg.role]"
      >
        <div class="message-content">
          <div v-if="msg.role === 'assistant'" class="ai-icon">🤖</div>
          <div v-else class="user-icon">👤</div>
          <div class="text" v-html="formatMessage(msg.content)"></div>
        </div>
        <!-- 显示工具调用信息 -->
        <div v-if="msg.toolCalls && msg.toolCalls.length" class="tool-calls">
          <el-collapse>
            <el-collapse-item title="查看执行的操作">
              <div v-for="tool in msg.toolCalls" :key="tool.id" class="tool-call">
                <el-tag size="small">{{ tool.name }}</el-tag>
                <pre>{{ JSON.stringify(tool.arguments, null, 2) }}</pre>
              </div>
            </el-collapse-item>
          </el-collapse>
        </div>
      </div>
    </div>

    <div class="chat-input">
      <el-input
        v-model="inputMessage"
        type="textarea"
        :rows="2"
        placeholder="输入消息，例如：查询库存预警的物资"
        @keydown.enter.ctrl="sendMessage"
        :disabled="loading"
      />
      <el-button type="primary" @click="sendMessage" :loading="loading">
        发送
      </el-button>
    </div>

    <div class="quick-actions">
      <el-button size="small" @click="quickAction('查询库存预警的物资')">
        库存预警
      </el-button>
      <el-button size="small" @click="quickAction('分析当前库存状况')">
        库存分析
      </el-button>
      <el-button size="small" @click="quickAction('查询待审批的工作流任务')">
        待办任务
      </el-button>
      <el-button size="small" @click="quickAction('生成库存汇总报告')">
        库存报告
      </el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick } from 'vue'
import { agentApi } from '@/api'
import { ElMessage } from 'element-plus'

const messages = ref([])
const inputMessage = ref('')
const loading = ref(false)
const conversationHistory = ref([])
const messagesContainer = ref(null)

// 发送消息
async function sendMessage() {
  const message = inputMessage.value.trim()
  if (!message || loading.value) return

  // 添加用户消息
  messages.value.push({
    role: 'user',
    content: message
  })

  inputMessage.value = ''
  loading.value = true

  try {
    const response = await agentApi.chat(message, conversationHistory.value)

    // 更新对话历史
    conversationHistory.value = response.data.conversation || []

    // 添加AI回复
    messages.value.push({
      role: 'assistant',
      content: response.data.message,
      toolCalls: response.data.tool_calls
    })

    // 滚动到底部
    await nextTick()
    scrollToBottom()
  } catch (error) {
    console.error('AI Chat error:', error)
    ElMessage.error(error.response?.data?.message || 'AI 服务暂时不可用')
    messages.value.push({
      role: 'assistant',
      content: '抱歉，我遇到了一些问题，请稍后再试。',
      isError: true
    })
  } finally {
    loading.value = false
  }
}

// 快捷操作
function quickAction(question) {
  inputMessage.value = question
  sendMessage()
}

// 格式化消息（支持简单的Markdown）
function formatMessage(content) {
  if (!content) return ''
  return content
    .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
    .replace(/\n/g, '<br>')
}

// 滚动到底部
function scrollToBottom() {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}
</script>

<style scoped>
.ai-chat {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid #eee;
}

.chat-header h3 {
  margin: 0;
  font-size: 16px;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.message {
  margin-bottom: 16px;
}

.message.user {
  text-align: right;
}

.message-content {
  display: inline-flex;
  align-items: flex-start;
  gap: 8px;
  max-width: 80%;
  text-align: left;
}

.message.user .message-content {
  flex-direction: row-reverse;
}

.ai-icon, .user-icon {
  font-size: 24px;
  flex-shrink: 0;
}

.text {
  padding: 12px 16px;
  border-radius: 12px;
  background: #f5f7fa;
  line-height: 1.6;
}

.message.user .text {
  background: #409eff;
  color: white;
}

.tool-calls {
  margin-top: 8px;
  max-width: 80%;
}

.tool-call {
  margin-bottom: 8px;
}

.tool-call pre {
  margin: 4px 0;
  padding: 8px;
  background: #f5f5f5;
  border-radius: 4px;
  font-size: 12px;
  overflow-x: auto;
}

.chat-input {
  display: flex;
  gap: 12px;
  padding: 16px;
  border-top: 1px solid #eee;
}

.chat-input .el-textarea {
  flex: 1;
}

.quick-actions {
  display: flex;
  gap: 8px;
  padding: 0 16px 16px;
  flex-wrap: wrap;
}
</style>
