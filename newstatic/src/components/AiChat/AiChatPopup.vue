<template>
  <el-drawer
    v-model="visible"
    direction="rtl"
    size="400px"
    title="AI 助手"
    :append-to-body="true"
    class="ai-chat-drawer"
  >
    <template #header>
      <div class="chat-header">
        <span class="title">AI 助手</span>
        <div class="header-actions">
          <!-- 模型切换 -->
          <el-dropdown @command="onSelectModel" trigger="click">
            <div class="model-selector">
              <span class="model-name">{{ currentModelName }}</span>
              <el-icon><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item
                  v-for="p in providers"
                  :key="p.id"
                  :command="p.id"
                  :class="{ 'is-active': p.id === currentProvider }"
                >
                  {{ p.name }}
                  <el-icon v-if="p.id === currentProvider" class="check-icon"><Check /></el-icon>
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
          <el-tooltip content="清空对话" placement="top">
            <el-icon class="clear-icon" @click="clearHistory"><Delete /></el-icon>
          </el-tooltip>
        </div>
      </div>
    </template>

    <!-- 消息列表 -->
    <div class="chat-messages" ref="messagesContainer">
      <div v-if="messages.length === 0 && !loading" class="empty-hint">
        <div class="context-hint">{{ defaultHint }}</div>
      </div>
      <div
        v-for="(msg, index) in messages"
        :key="index"
        :class="['message', msg.role]"
      >
        <div class="message-content">
          <div v-if="msg.role === 'assistant'" class="avatar ai">AI</div>
          <div v-else class="avatar user">我</div>
          <div class="text" v-html="formatMessage(msg.content)"></div>
        </div>
      </div>
      <div v-if="loading" class="message assistant">
        <div class="message-content">
          <div class="avatar ai">AI</div>
          <div class="thinking-bubble">
            <span class="thinking-text">思考中</span>
            <span class="thinking-dots">
              <span class="dot">.</span>
              <span class="dot">.</span>
              <span class="dot">.</span>
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- 输入区域 -->
    <div class="chat-input">
      <el-input
        v-model="inputMessage"
        placeholder="输入消息..."
        :disabled="loading"
        @keydown.enter="sendMessage"
        class="input-field"
      >
        <template #append>
          <el-button
            type="primary"
            @click="sendMessage"
            :loading="loading"
            :disabled="!inputMessage.trim()"
          >
            发送
          </el-button>
        </template>
      </el-input>
    </div>
  </el-drawer>
</template>

<script setup>
import { ref, computed, nextTick, watch, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { ArrowDown, Delete, Check } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { useNotificationStore } from '@/stores/notificationStore'

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:show'])

const authStore = useAuthStore()
const notificationStore = useNotificationStore()

const visible = ref(props.show)
const messages = ref([])
const inputMessage = ref('')
const loading = ref(false)
const messagesContainer = ref(null)

// 历史记录存储键 - 按用户区分
const getHistoryStorageKey = () => {
  const userId = authStore.user?.id
  if (userId) {
    return `ai_chat_history_${userId}`
  }
  return 'ai_chat_history_default'
}
const MAX_HISTORY_MESSAGES = 50 // 最多保存50条消息

// 模型切换相关
const providers = ref([])
const currentProvider = ref('')
const currentModelName = computed(() => {
  const provider = providers.value.find(p => p.id === currentProvider.value)
  return provider ? provider.name : 'AI'
})

// 获取模型列表
async function fetchProviders() {
  try {
    const res = await fetch('/api/agent/providers', {
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      }
    })
    const data = await res.json()
    if (data.data) {
      providers.value = data.data.providers || []
      currentProvider.value = data.data.current_provider || ''
    }
  } catch (error) {
    console.error('获取模型列表失败:', error)
  }
}

// 切换模型
async function onSelectModel(providerId) {
  if (providerId === currentProvider.value) {
    return
  }

  try {
    const res = await fetch('/api/agent/providers/switch', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authStore.token}`
      },
      body: JSON.stringify({ provider: providerId })
    })
    const data = await res.json()
    if (data.data) {
      currentProvider.value = data.data.current_provider
      ElMessage.success(`已切换到 ${providers.value.find(p => p.id === providerId)?.name || 'AI'}`)
    }
  } catch (error) {
    console.error('切换模型失败:', error)
    ElMessage.error('切换模型失败')
  }
}

// AI 响应回调移除函数
let removeAiCallback = null

// 当前会话ID（用于匹配流式响应）
let currentSessionId = ref('')
let isProcessing = ref(false) // 请求处理锁

// 超时保护定时器
let responseTimeoutTimer = null
const RESPONSE_TIMEOUT = 90000 // 90秒超时保护

// 默认功能提示
const defaultHint = `您好！我是AI助手，我可以帮您：

📋 任务管理
• 查询今日/明日任务安排
• 创建、修改施工预约
• 查看待审批事项

📦 库存管理
• 查询物资库存和预警
• 入库/出库操作
• 查询物资计划和领用单

📊 其他功能
• 查询考勤打卡记录
• 查询项目列表
• 施工日志查询

请问有什么可以帮您的？`

// 监听显示状态
watch(() => props.show, (val) => {
  visible.value = val
  if (val) {
    nextTick(() => scrollToBottom())
  }
})

watch(visible, (val) => {
  emit('update:show', val)
})

// 保存历史到本地存储
function saveHistory() {
  try {
    const historyToSave = messages.value.slice(-MAX_HISTORY_MESSAGES)
    localStorage.setItem(getHistoryStorageKey(), JSON.stringify(historyToSave))
  } catch (e) {
    console.error('保存聊天历史失败:', e)
  }
}

// 从本地存储加载历史
function loadHistory() {
  try {
    const saved = localStorage.getItem(getHistoryStorageKey())
    if (saved) {
      messages.value = JSON.parse(saved)
    }
  } catch (e) {
    console.error('加载聊天历史失败:', e)
  }
}

// 获取用于发送的历史记录（格式转换）
function getHistoryForAPI() {
  return messages.value.map(msg => ({
    role: msg.role,
    content: msg.content
  }))
}

// 添加消息
function addMessage(role, content) {
  messages.value.push({ role, content })
  saveHistory()
  nextTick(() => scrollToBottom())
}

// 重置处理状态（用于超时或错误恢复）
function resetProcessingState() {
  loading.value = false
  isProcessing.value = false
  if (responseTimeoutTimer) {
    clearTimeout(responseTimeoutTimer)
    responseTimeoutTimer = null
  }
}

// 启动超时保护
function startResponseTimeout() {
  // 清除之前的定时器
  if (responseTimeoutTimer) {
    clearTimeout(responseTimeoutTimer)
  }
  responseTimeoutTimer = setTimeout(() => {
    console.log('[AiChatPopup] Response timeout, resetting state')
    resetProcessingState()
    ElMessage.warning('响应超时，请重试')
    // 添加超时提示消息
    if (messages.value.length > 0 && messages.value[messages.value.length - 1].role === 'assistant' && !messages.value[messages.value.length - 1].content) {
      messages.value.pop() // 移除空的 assistant 消息
    }
  }, RESPONSE_TIMEOUT)
}

// 发送消息
async function sendMessage() {
  const message = inputMessage.value.trim()
  if (!message || loading.value || isProcessing.value) return

  // 设置处理锁
  isProcessing.value = true
  addMessage('user', message)
  inputMessage.value = ''
  loading.value = true

  // 生成新的会话ID
  currentSessionId.value = 'session_' + Date.now()

  // 启动超时保护
  startResponseTimeout()

  try {
    // 使用 HTTP 发送消息
    const res = await fetch('/api/agent/chat', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authStore.token}`
      },
      body: JSON.stringify({
        message: message,
        conversation_history: getHistoryForAPI().slice(0, -1) // 排除刚添加的用户消息
      })
    })

    if (!res.ok) {
      throw new Error(`HTTP error! status: ${res.status}`)
    }

    const result = await res.json()

    // 清除超时定时器
    if (responseTimeoutTimer) {
      clearTimeout(responseTimeoutTimer)
      responseTimeoutTimer = null
    }
    loading.value = false
    isProcessing.value = false

    if (result.data && result.data.message) {
      addMessage('assistant', result.data.message)
    } else if (result.message) {
      addMessage('assistant', result.message)
    } else {
      addMessage('assistant', '抱歉，我没有理解您的问题，请换一种方式提问。')
    }

    saveHistory()
  } catch (error) {
    console.error('Chat error:', error)
    ElMessage.error('发送失败: ' + (error.message || '请稍后重试'))
    addMessage('assistant', '抱歉，我遇到了一些问题，请稍后再试。')
    resetProcessingState()
  }
}

// 清空历史
function clearHistory() {
  messages.value = []
  localStorage.removeItem(getHistoryStorageKey())
  // 同时调用后端清除API
  fetch('/api/agent/conversation-history', {
    method: 'DELETE',
    headers: {
      'Authorization': `Bearer ${authStore.token}`
    }
  }).catch(e => console.error('清除服务器历史失败:', e))
  ElMessage.success('对话已清空')
}

// 滚动到底部
function scrollToBottom() {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

// 格式化消息
function formatMessage(content) {
  if (!content) return ''
  return content
    .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
    .replace(/\n/g, '<br>')
}

// 清除当前用户的聊天历史（用于登出时调用）
function clearCurrentUserHistory() {
  localStorage.removeItem(getHistoryStorageKey())
}

// 暴露方法供外部调用
defineExpose({
  clearCurrentUserHistory
})

onMounted(() => {
  loadHistory() // 加载历史记录
  fetchProviders() // 加载模型列表
})

onUnmounted(() => {
  // 清理超时定时器
  if (responseTimeoutTimer) {
    clearTimeout(responseTimeoutTimer)
    responseTimeoutTimer = null
  }
})

// 监听弹窗关闭
watch(visible, (val) => {
  if (!val) {
    // 重置处理状态，确保下次打开时不会卡住
    resetProcessingState()
  }
})
</script>

<style scoped>
.ai-chat-drawer :deep(.el-drawer__header) {
  margin-bottom: 0;
  padding: 0;
}

.ai-chat-drawer :deep(.el-drawer__body) {
  padding: 0;
  display: flex;
  flex-direction: column;
  height: calc(100vh - 40px);
}

.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: #fff;
  border-bottom: 1px solid #eee;
}

.title {
  font-size: 16px;
  font-weight: bold;
  color: #323233;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.model-selector {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  background: #f5f7fa;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.model-selector:hover {
  background: #e8e8e8;
}

.model-name {
  font-size: 12px;
  color: #323233;
  max-width: 80px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.clear-icon {
  font-size: 20px;
  color: #969799;
  cursor: pointer;
  transition: color 0.2s;
}

.clear-icon:hover {
  color: #409eff;
}

.check-icon {
  margin-left: 8px;
  color: #409eff;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  background: #f7f8fa;
}

.empty-hint {
  text-align: center;
  color: #969799;
  padding: 40px 20px;
  font-size: 14px;
}

.context-hint {
  text-align: left;
  background: #fff;
  border-radius: 12px;
  padding: 16px;
  font-size: 14px;
  color: #323233;
  line-height: 1.8;
  white-space: pre-line;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
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
  gap: 10px;
  max-width: 85%;
  text-align: left;
}

.message.user .message-content {
  flex-direction: row-reverse;
}

.avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: bold;
  flex-shrink: 0;
}

.avatar.ai {
  background: linear-gradient(135deg, #409eff 0%, #66b1ff 100%);
  color: white;
}

.avatar.user {
  background: linear-gradient(135deg, #67c23a 0%, #85ce61 100%);
  color: white;
}

.text {
  padding: 12px 16px;
  border-radius: 16px;
  background: #fff;
  line-height: 1.6;
  font-size: 14px;
  word-break: break-word;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.message.user .text {
  background: linear-gradient(135deg, #409eff 0%, #66b1ff 100%);
  color: white;
}

/* 思考中动画 */
.thinking-bubble {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  border-radius: 16px;
  background: #fff;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.thinking-text {
  font-size: 14px;
  color: #646566;
  margin-right: 2px;
}

.thinking-dots {
  display: inline-flex;
}

.thinking-dots .dot {
  font-size: 14px;
  color: #409eff;
  animation: bounce 1.4s infinite ease-in-out both;
}

.thinking-dots .dot:nth-child(1) {
  animation-delay: -0.32s;
}

.thinking-dots .dot:nth-child(2) {
  animation-delay: -0.16s;
}

.thinking-dots .dot:nth-child(3) {
  animation-delay: 0s;
}

@keyframes bounce {
  0%, 80%, 100% {
    transform: translateY(0);
  }
  40% {
    transform: translateY(-8px);
  }
}

.chat-input {
  padding: 12px 16px;
  background: #fff;
  border-top: 1px solid #eee;
}

.input-field {
  --el-input-border-radius: 20px;
}

.input-field :deep(.el-input-group__append) {
  border-radius: 0 20px 20px 0;
  background: transparent;
  padding: 0;
}

.input-field :deep(.el-input-group__append .el-button) {
  border-radius: 0 20px 20px 0;
}
</style>
