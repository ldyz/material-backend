<template>
  <van-popup
    v-model:show="visible"
    position="bottom"
    :style="{ height: '70%' }"
    round
    closeable
    class="ai-chat-popup-wrapper"
  >
    <div class="ai-chat-popup">
      <!-- 标题栏 -->
      <div class="chat-header">
        <span class="title">AI 助手</span>
        <div class="header-actions">
          <!-- 模型切换 -->
          <div class="model-selector" @click="showModelPicker = true">
            <span class="model-name">{{ currentModelName }}</span>
            <van-icon name="arrow-down" size="12" />
          </div>
          <!-- 语音对话模式切换 -->
          <van-icon
            :name="voiceChatMode ? 'phone-circle' : 'phone-circle-o'"
            @click="toggleVoiceChatMode"
            class="voice-mode-icon"
            :class="{ active: voiceChatMode }"
          />
          <van-icon
            :name="autoSpeak ? 'volume' : 'volume-o'"
            @click="toggleAutoSpeak"
            class="speak-icon"
            :class="{ active: autoSpeak }"
          />
          <van-icon name="delete-o" @click="clearHistory" class="clear-icon" />
        </div>
      </div>

      <!-- 消息列表 -->
      <div class="chat-messages" ref="messagesContainer">
        <div v-if="messages.length === 0 && !loading" class="empty-hint">
          <div v-if="contextHint" class="context-hint">
            {{ contextHint }}
          </div>
          <div v-else>
            {{ voiceChatMode ? '点击下方按钮开始语音对话' : '点击下方语音按钮或输入文字开始对话' }}
          </div>
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
            <!-- AI 消息播放按钮 -->
            <van-icon
              v-if="msg.role === 'assistant' && msg.content"
              :name="currentSpeakingIndex === index ? 'pause-circle-o' : 'play-circle-o'"
              class="play-btn"
              @click="toggleSpeak(msg.content, index)"
            />
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
        <!-- 语音对话模式 -->
        <template v-if="voiceChatMode">
          <div class="voice-chat-controls">
            <van-button
              :type="isListening ? 'danger' : 'primary'"
              size="small"
              @click="isListening ? stopListening() : startListening()"
              :loading="loading"
              class="voice-chat-btn"
            >
              {{ isListening ? '停止' : (loading ? '处理中' : '开始对话') }}
            </van-button>
            <span v-if="isListening" class="listening-hint">
              {{ listeningText || '正在听...' }}
            </span>
          </div>
        </template>
        <!-- 普通模式 -->
        <template v-else>
          <!-- 语音/键盘切换按钮 -->
          <div class="mode-switch-btn" @click="toggleInputMode">
            <van-icon :name="isVoiceMode ? 'edit' : 'volume-o'" size="22" />
          </div>

          <!-- 文字输入模式 -->
          <template v-if="!isVoiceMode">
            <van-field
              v-model="inputMessage"
              placeholder="输入消息..."
              :disabled="loading"
              @keydown.enter="sendMessage"
              class="input-field"
            />
            <van-button
              size="small"
              type="primary"
              @click="sendMessage"
              :loading="loading"
              :disabled="!inputMessage.trim()"
              class="send-btn"
            >
              发送
            </van-button>
          </template>

          <!-- 语音输入模式 -->
          <template v-else>
            <div
              class="voice-input-btn"
              :class="{ pressing: isRecording }"
              @touchstart.prevent="startRecording"
              @touchend="stopRecording"
              @touchcancel="cancelRecording"
              @mousedown="startRecording"
              @mouseup="stopRecording"
              @mouseleave="cancelRecording"
            >
              {{ isRecording ? '松开发送' : '按住说话' }}
            </div>
          </template>
        </template>
      </div>
    </div>

    <!-- 录音提示 -->
    <van-overlay :show="isRecording || isListening" class="recording-overlay">
      <div class="recording-content">
        <div class="recording-waves">
          <div class="wave" v-for="i in 5" :key="i" :style="{ animationDelay: `${i * 0.1}s` }"></div>
        </div>
        <div class="recording-text">{{ listeningText || '正在录音...' }}</div>
        <div class="recording-hint">{{ voiceChatMode ? '检测到静音将自动发送' : '松开发送，上滑取消' }}</div>
      </div>
    </van-overlay>
  </van-popup>

  <!-- 模型选择面板 -->
  <van-action-sheet
    v-model:show="showModelPicker"
    :actions="modelActions"
    cancel-text="取消"
    close-on-click-action
    @select="onSelectModel"
  />
</template>

<script setup>
import { ref, computed, nextTick, watch, onMounted, onUnmounted } from 'vue'
import { showToast } from 'vant'
import webSocketService from '@/utils/websocket'
import { VoiceRecorder } from '@independo/capacitor-voice-recorder'
import { Capacitor } from '@capacitor/core'
import { TextToSpeech } from '@capacitor-community/text-to-speech'
import { storage } from '@/utils/storage'

// 获取 API 基础 URL
const getApiBaseURL = () => {
  const isCapacitorEnv = typeof window !== 'undefined' && window.Capacitor
  return isCapacitorEnv ? 'https://home.mbed.org.cn:9090/api' : '/api'
}

const props = defineProps({
  show: Boolean,
  context: {
    type: String,
    default: '' // 'inbound', 'requisition', 'dashboard' 等
  },
  contextData: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['update:show'])

const visible = ref(props.show)
const messages = ref([])
const inputMessage = ref('')
const loading = ref(false)
const messagesContainer = ref(null)

// 历史记录存储键 - 按用户区分
const getHistoryStorageKey = () => {
  const userStr = localStorage.getItem('user_info')
  if (userStr) {
    try {
      const user = JSON.parse(userStr)
      if (user.id) {
        return `ai_chat_history_${user.id}`
      }
    } catch (e) {
      console.error('解析用户信息失败:', e)
    }
  }
  return 'ai_chat_history_default'
}
const MAX_HISTORY_MESSAGES = 50 // 最多保存50条消息

// TTS 语音播报相关
const autoSpeak = ref(true) // 自动播放开关
const currentSpeakingIndex = ref(-1) // 当前正在播放的消息索引
const isSpeaking = ref(false)

// 语音对话模式相关
const voiceChatMode = ref(false) // 语音对话模式开关
const isListening = ref(false) // 是否正在监听
const listeningText = ref('') // 监听中的文字
let audioContext = null
let analyser = null
let mediaStream = null
let silenceTimer = null
let silenceStartTime = 0
const SILENCE_THRESHOLD = 1.5 // 静音 1.5 秒自动发送

// 模型切换相关
const showModelPicker = ref(false)
const providers = ref([])
const currentProvider = ref('')
const currentModelName = computed(() => {
  const provider = providers.value.find(p => p.id === currentProvider.value)
  return provider ? provider.name : 'AI'
})
const modelActions = computed(() => {
  return providers.value.map(p => ({
    name: p.name,
    value: p.id,
    className: p.id === currentProvider.value ? 'active-model' : ''
  }))
})

// 获取模型列表
async function fetchProviders() {
  try {
    const baseURL = getApiBaseURL()
    const token = storage.getToken()
    const res = await fetch(`${baseURL}/agent/providers`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    const data = await res.json()
    console.log('[AiChatPopup] fetchProviders response:', data)
    if (data.data) {
      providers.value = data.data.providers || []
      currentProvider.value = data.data.current_provider || ''
      console.log('[AiChatPopup] providers loaded:', providers.value.length, 'current:', currentProvider.value)
    }
  } catch (error) {
    console.error('获取模型列表失败:', error)
  }
}

// 切换模型
async function onSelectModel(action) {
  if (action.value === currentProvider.value) {
    showModelPicker.value = false
    return
  }

  try {
    const baseURL = getApiBaseURL()
    const token = storage.getToken()
    const res = await fetch(`${baseURL}/agent/providers/switch`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({ provider: action.value })
    })
    const data = await res.json()
    if (data.data) {
      currentProvider.value = data.data.current_provider
      showToast(`已切换到 ${action.name}`)
    }
  } catch (error) {
    console.error('切换模型失败:', error)
    showToast('切换模型失败')
  }
  showModelPicker.value = false
}

// 录音相关状态
const isRecording = ref(false)
const isStarting = ref(false)
const mediaRecorder = ref(null)
const audioChunks = ref([])
const recordingStartTime = ref(0)
const touchStartY = ref(0)
const isCancelled = ref(false)
const isTouchDevice = ref(false)

// 输入模式：false=文字, true=语音
const isVoiceMode = ref(false)

// AI 响应回调移除函数
let removeAiCallback = null

// 当前会话ID（用于匹配流式响应）
let currentSessionId = ref('')
let isProcessing = ref(false) // 请求处理锁

// 超时保护定时器
let responseTimeoutTimer = null
const RESPONSE_TIMEOUT = 90000 // 90秒超时保护

// 打字机效果相关
let typewriterBuffer = '' // 待显示的文本缓冲区
let typewriterTimer = null // 打字机定时器
const TYPEWRITER_SPEED = 50 // 每个字符显示间隔（毫秒）

// 上下文提示信息
const contextHints = {
  inbound: '您好！我是入库助手。\n• 我可以帮您查询物资库存\n• 查询供应商信息\n• 推荐合适的物资计划\n• 解答入库流程问题\n\n请问有什么可以帮您的？',
  requisition: '您好！我是出库助手。\n• 我可以帮您查询物资库存\n• 查询可用的物资计划\n• 检查出库流程\n• 解答领料问题\n\n请问有什么可以帮您的？',
  dashboard: '您好！我是AI助手。\n• 查询今日/明日任务安排\n• 查看库存预警\n• 查看待审批事项\n• 考勤打卡查询\n\n请问有什么可以帮您的？'
}

// 默认功能提示（无特定上下文时显示）
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

// 计算当前上下文的提示信息
const contextHint = computed(() => {
  if (props.context && contextHints[props.context]) {
    return contextHints[props.context]
  }
  return defaultHint
})

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
  stopTypewriterAndFlush()
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
    showToast('响应超时，请重试')
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
    // 使用 WebSocket 发送消息
    if (!webSocketService.isConnected()) {
      showToast('连接中断，正在重连...')
      webSocketService.reconnect()
      await new Promise(resolve => setTimeout(resolve, 2000))
      if (!webSocketService.isConnected()) {
        showToast('连接失败，请重试')
        resetProcessingState()
        return
      }
    }

    // 注册 AI 回调
    registerAiCallback()

    // 发送消息（带上历史记录）
    const history = getHistoryForAPI().slice(0, -1) // 排除刚添加的用户消息
    console.log('[AiChatPopup] Sending chat message, sessionId:', currentSessionId.value)
    webSocketService.sendChat(message, history)
  } catch (error) {
    console.error('Chat error:', error)
    showToast('发送失败: ' + (error.message || '请稍后重试'))
    addMessage('assistant', '抱歉，我遇到了一些问题，请稍后再试。')
    resetProcessingState()
  }
}

// 打字机效果：开始
function startTypewriter() {
  if (typewriterTimer) {
    clearInterval(typewriterTimer)
  }
  typewriterBuffer = ''

  typewriterTimer = setInterval(() => {
    if (typewriterBuffer.length > 0) {
      // 取出第一个字符
      const char = typewriterBuffer[0]
      typewriterBuffer = typewriterBuffer.slice(1)

      // 添加到最后一条 AI 消息
      if (messages.value.length > 0 && messages.value[messages.value.length - 1].role === 'assistant') {
        messages.value[messages.value.length - 1].content += char
        nextTick(() => scrollToBottom())
      }
    }
  }, TYPEWRITER_SPEED)
}

// 打字机效果：添加文本到缓冲区
function appendToTypewriter(text) {
  typewriterBuffer += text
}

// 打字机效果：停止并立即显示剩余内容
function stopTypewriterAndFlush() {
  if (typewriterTimer) {
    clearInterval(typewriterTimer)
    typewriterTimer = null
  }
  // 立即显示剩余内容
  if (typewriterBuffer.length > 0) {
    if (messages.value.length > 0 && messages.value[messages.value.length - 1].role === 'assistant') {
      messages.value[messages.value.length - 1].content += typewriterBuffer
    }
    typewriterBuffer = ''
  }
}

// 注册 AI 回调
function registerAiCallback() {
  // 移除旧的回调
  if (removeAiCallback) {
    removeAiCallback()
  }

  // 保存当前会话ID的快照
  const sessionId = currentSessionId.value

  // 注册新的回调
  removeAiCallback = webSocketService.onAiResponse((message) => {
    console.log('[AiChatPopup] AI message received:', message)

    // 检查会话ID匹配（如果有）
    if (message.session_id && message.session_id !== sessionId) {
      console.log('[AiChatPopup] Ignoring message from different session:', message.session_id, 'current:', sessionId)
      return
    }

    switch (message.type) {
      case 'voice_processing':
        // 语音正在处理
        console.log('[AiChatPopup] Voice processing...')
        break

      case 'voice_transcript':
        // 语音识别结果 - 更新占位消息内容
        console.log('[AiChatPopup] Voice transcript received:', message.text)
        if (message.text) {
          // 查找最后一条用户消息（可能是占位符）
          let lastUserIndex = -1
          for (let i = messages.value.length - 1; i >= 0; i--) {
            if (messages.value[i].role === 'user') {
              lastUserIndex = i
              break
            }
          }
          console.log('[AiChatPopup] lastUserIndex:', lastUserIndex, 'messages count:', messages.value.length)
          if (lastUserIndex >= 0) {
            // 检查是否是占位符消息（包含 🎤 或是刚添加的）
            const lastMsg = messages.value[lastUserIndex]
            console.log('[AiChatPopup] Last user message content:', lastMsg.content)
            if (lastMsg.content.includes('🎤') || lastMsg.content.includes('语音消息')) {
              // 使用 splice 替换确保 Vue 响应式更新
              messages.value.splice(lastUserIndex, 1, {
                role: 'user',
                content: message.text
              })
              saveHistory()
              nextTick(() => scrollToBottom())
              console.log('[AiChatPopup] Updated placeholder to:', message.text)
            } else {
              // 最后一条用户消息不是占位符，添加新消息
              console.log('[AiChatPopup] Last user message is not a placeholder, adding new message')
              addMessage('user', message.text)
            }
          } else {
            // 没有用户消息，添加新消息
            console.log('[AiChatPopup] No user message found, adding new message')
            addMessage('user', message.text)
          }
        }
        break

      case 'voice_transcript_partial':
        // 实时识别部分结果
        if (message.text && voiceChatMode.value) {
          listeningText.value = message.text
        }
        break

      case 'ai_response_start':
        // AI 回复开始
        loading.value = true
        // 添加一个空的assistant消息，用于后续填充
        if (messages.value.length === 0 || messages.value[messages.value.length - 1].role !== 'assistant') {
          addMessage('assistant', '')
        }
        // 启动打字机效果
        startTypewriter()
        break

      case 'ai_response_chunk':
        // AI 回复片段（流式）- 添加到打字机缓冲区
        if (message.content) {
          appendToTypewriter(message.content)
        }
        break

      case 'ai_response_done':
        // AI 回复完成
        console.log('[AiChatPopup] ai_response_done received, message:', message.message?.substring(0, 50))
        // 清除超时定时器
        if (responseTimeoutTimer) {
          clearTimeout(responseTimeoutTimer)
          responseTimeoutTimer = null
        }
        loading.value = false
        isProcessing.value = false // 释放处理锁
        // 停止打字机并显示剩余内容
        stopTypewriterAndFlush()
        // 如果有完整消息，确保内容正确
        if (message.message) {
          // 更新最后一条消息为完整内容
          if (messages.value.length > 0 && messages.value[messages.value.length - 1].role === 'assistant') {
            messages.value[messages.value.length - 1].content = message.message
          } else {
            addMessage('assistant', message.message)
          }
        }
        saveHistory() // 保存历史
        nextTick(() => scrollToBottom())
        // 自动播放 AI 回复（异步执行，避免阻塞）
        setTimeout(() => autoSpeakLastMessage(), 100)
        break

      case 'error':
        // 错误
        console.log('[AiChatPopup] error received:', message.message)
        // 清除超时定时器
        if (responseTimeoutTimer) {
          clearTimeout(responseTimeoutTimer)
          responseTimeoutTimer = null
        }
        loading.value = false
        isProcessing.value = false // 释放处理锁
        // 停止打字机
        stopTypewriterAndFlush()
        showToast(message.message || '处理失败')
        break
    }
  })
}

// 检查录音权限
async function checkRecordPermission() {
  try {
    if (Capacitor.isNativePlatform()) {
      const status = await VoiceRecorder.requestAudioRecordingPermission()
      if (status.value) {
        return true
      } else {
        showToast('请授予麦克风权限')
        return false
      }
    } else {
      // Web 环境
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
      stream.getTracks().forEach(track => track.stop())
      return true
    }
  } catch (error) {
    console.error('权限请求失败:', error)
    showToast('请授予麦克风权限')
    return false
  }
}

// 开始录音
async function startRecording(event) {
  // 防止鼠标和触摸事件重复触发
  if (event.type === 'touchstart') {
    isTouchDevice.value = true
    touchStartY.value = event.touches[0].clientY
  } else if (event.type === 'mousedown' && isTouchDevice.value) {
    return
  }

  if (isRecording.value || isStarting.value) return

  isStarting.value = true

  try {
    if (!await checkRecordPermission()) {
      isStarting.value = false
      return
    }

    if (Capacitor.isNativePlatform()) {
      await VoiceRecorder.startRecording()
      isRecording.value = true
      isCancelled.value = false
      recordingStartTime.value = Date.now()
    } else {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true })

      let mimeType = 'audio/webm'
      if (!MediaRecorder.isTypeSupported(mimeType)) {
        mimeType = 'audio/webm;codecs=opus'
        if (!MediaRecorder.isTypeSupported(mimeType)) {
          mimeType = 'audio/mp4'
          if (!MediaRecorder.isTypeSupported(mimeType)) {
            mimeType = ''
          }
        }
      }

      const options = mimeType ? { mimeType } : {}
      mediaRecorder.value = new MediaRecorder(stream, options)
      audioChunks.value = []

      mediaRecorder.value.ondataavailable = (event) => {
        if (event.data.size > 0) {
          audioChunks.value.push(event.data)
        }
      }

      mediaRecorder.value.onerror = (event) => {
        console.error('MediaRecorder error:', event.error)
        showToast('录音出错: ' + event.error.message)
        stopRecordingInternal()
      }

      mediaRecorder.value.start()
      isRecording.value = true
      isCancelled.value = false
      recordingStartTime.value = Date.now()
    }
  } catch (error) {
    console.error('录音启动失败:', error)
    showToast('录音启动失败: ' + (error.message || '请检查麦克风权限'))
  } finally {
    isStarting.value = false
  }
}

// 停止录音
async function stopRecording(event) {
  if (event.type === 'touchend') {
    const touchEndY = event.changedTouches[0].clientY
    if (touchStartY.value - touchEndY > 100) {
      isCancelled.value = true
    }
  } else if (event.type === 'mouseup' && isTouchDevice.value) {
    return
  }

  if (isStarting.value) {
    await new Promise(resolve => {
      const check = setInterval(() => {
        if (!isStarting.value) {
          clearInterval(check)
          resolve()
        }
      }, 50)
    })
  }

  if (!isRecording.value) return

  if (isCancelled.value) {
    cancelRecording()
    return
  }

  try {
    if (Capacitor.isNativePlatform()) {
      const result = await VoiceRecorder.stopRecording()
      const audioBlob = base64ToBlob(result.value.recordDataBase64, result.value.mimeType)
      const duration = result.value.msDuration

      isRecording.value = false

      if (duration < 500) {
        showToast('录音时间太短')
        return
      }

      await sendVoiceViaWebSocket(audioBlob, result.value.mimeType)
    } else {
      if (!mediaRecorder.value) return

      return new Promise((resolve) => {
        mediaRecorder.value.onstop = async () => {
          const mimeType = mediaRecorder.value.mimeType || 'audio/webm'
          const audioBlob = new Blob(audioChunks.value, { type: mimeType })
          const duration = Date.now() - recordingStartTime.value

          if (mediaRecorder.value.stream) {
            mediaRecorder.value.stream.getTracks().forEach(track => track.stop())
          }

          isRecording.value = false

          if (duration < 500) {
            showToast('录音时间太短')
            resolve()
            return
          }

          await sendVoiceViaWebSocket(audioBlob, mimeType)
          resolve()
        }

        mediaRecorder.value.stop()
      })
    }
  } catch (error) {
    console.error('停止录音失败:', error)
    isRecording.value = false
  }
}

// 内部停止录音
async function stopRecordingInternal() {
  if (Capacitor.isNativePlatform()) {
    try {
      await VoiceRecorder.stopRecording()
    } catch (e) {
      // 忽略错误
    }
  } else if (mediaRecorder.value && mediaRecorder.value.stream) {
    mediaRecorder.value.stream.getTracks().forEach(track => track.stop())
  }
  isRecording.value = false
}

// 取消录音
async function cancelRecording() {
  if (isStarting.value) {
    return
  }

  if (isRecording.value) {
    if (Capacitor.isNativePlatform()) {
      try {
        await VoiceRecorder.stopRecording()
      } catch (e) {
        // 忽略错误
      }
    } else if (mediaRecorder.value && mediaRecorder.value.stream) {
      mediaRecorder.value.stream.getTracks().forEach(track => track.stop())
    }
  }
  isRecording.value = false
  isCancelled.value = false
}

// Base64 转 Blob
function base64ToBlob(base64, mimeType) {
  const byteCharacters = atob(base64)
  const byteNumbers = new Array(byteCharacters.length)
  for (let i = 0; i < byteCharacters.length; i++) {
    byteNumbers[i] = byteCharacters.charCodeAt(i)
  }
  const byteArray = new Uint8Array(byteNumbers)
  return new Blob([byteArray], { type: mimeType })
}

// Blob 转 Base64
function blobToBase64(blob) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onloadend = () => {
      const base64 = reader.result.split(',')[1]
      resolve(base64)
    }
    reader.onerror = reject
    reader.readAsDataURL(blob)
  })
}

// 通过 WebSocket 发送语音
async function sendVoiceViaWebSocket(audioBlob, mimeType) {
  try {
    // 检查是否有请求正在处理
    if (isProcessing.value) {
      showToast('请等待当前请求完成')
      return
    }

    console.log('[AiChatPopup] Sending voice message, size:', audioBlob.size, 'mimeType:', mimeType)
    isProcessing.value = true
    loading.value = true

    // 启动超时保护（语音处理需要更长时间）
    startResponseTimeout()

    if (!webSocketService.isConnected()) {
      showToast('WebSocket 未连接，正在重连...')
      webSocketService.reconnect()
      await new Promise(resolve => setTimeout(resolve, 2000))
      if (!webSocketService.isConnected()) {
        showToast('连接失败，请重试')
        resetProcessingState()
        return
      }
    }

    const base64Audio = await blobToBase64(audioBlob)

    // 生成新的会话ID
    currentSessionId.value = 'session_' + Date.now()

    // 立即添加用户消息占位符
    console.log('[AiChatPopup] Adding placeholder message')
    addMessage('user', '🎤 语音消息...')

    registerAiCallback()

    // 发送语音时带上历史记录（排除刚添加的占位消息）
    const history = getHistoryForAPI().slice(0, -1)
    console.log('[AiChatPopup] Sending voice via WebSocket, history length:', history.length)
    webSocketService.sendVoice(base64Audio, mimeType || 'audio/webm', history)
  } catch (error) {
    console.error('[AiChatPopup] 语音处理失败:', error)
    showToast('处理失败: ' + (error.message || '请稍后重试'))
    resetProcessingState()
    // 出错时移除占位消息
    const lastIndex = messages.value.length - 1
    if (lastIndex >= 0 && messages.value[lastIndex].content.includes('🎤')) {
      messages.value.pop()
    }
  }
}

// 清空历史
function clearHistory() {
  messages.value = []
  localStorage.removeItem(getHistoryStorageKey())
  // 同时调用后端清除API
  const baseURL = getApiBaseURL()
  const token = storage.getToken()
  fetch(`${baseURL}/agent/conversation-history`, {
    method: 'DELETE',
    headers: {
      'Authorization': `Bearer ${token}`
    }
  }).catch(e => console.error('清除服务器历史失败:', e))
  showToast('对话已清空')
}

// 切换输入模式（文字/语音）
function toggleInputMode() {
  isVoiceMode.value = !isVoiceMode.value
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

// ============ TTS 语音播报功能 ============

// 初始化 TTS（检查是否可用）
async function initTTS() {
  try {
    if (Capacitor.isNativePlatform()) {
      // 原生平台使用 Capacitor TTS
      console.log('[TTS] Using native Capacitor TTS')
    } else if ('speechSynthesis' in window) {
      // Web 平台使用 Web Speech API
      console.log('[TTS] Using Web Speech API')
      // 预加载语音列表
      window.speechSynthesis.getVoices()
    } else {
      console.log('[TTS] TTS not available')
    }
  } catch (error) {
    console.error('[TTS] Init error:', error)
  }
}

// 切换自动播放
function toggleAutoSpeak() {
  autoSpeak.value = !autoSpeak.value
  if (!autoSpeak.value && isSpeaking.value) {
    stopSpeaking()
  }
  showToast(autoSpeak.value ? '已开启自动播放' : '已关闭自动播放')
}

// 播放语音（支持原生和Web）
async function speakText(text) {
  // 清理 HTML 标签
  const plainText = text.replace(/<[^>]*>/g, '').replace(/\*\*/g, '')

  if (!plainText.trim()) {
    return false
  }

  try {
    // 先停止当前播放
    await stopSpeaking()

    const isNative = Capacitor.isNativePlatform()
    console.log('[TTS] isNativePlatform:', isNative)

    if (isNative) {
      // 使用 Capacitor 原生 TTS
      console.log('[TTS] Speaking with native TTS:', plainText.substring(0, 50) + '...')

      try {
        await TextToSpeech.speak({
          text: plainText,
          lang: 'zh-CN',
          rate: 1.0,
          pitch: 1.0,
          volume: 1.0,
          category: 'ambient',
        })
        isSpeaking.value = true
        console.log('[TTS] Native TTS started successfully')
        return true
      } catch (nativeError) {
        console.error('[TTS] Native TTS error:', nativeError)
        // 如果原生 TTS 失败，尝试 Web Speech API 作为后备
        if ('speechSynthesis' in window) {
          console.log('[TTS] Falling back to Web Speech API')
          return speakWithWebAPI(plainText)
        }
        showToast('语音播报失败: ' + (nativeError.message || 'TTS不可用'))
        return false
      }
    } else if ('speechSynthesis' in window) {
      // 使用 Web Speech API（浏览器环境）
      return speakWithWebAPI(plainText)
    } else {
      console.log('[TTS] No TTS available, isNative:', isNative)
      showToast('您的设备不支持语音播报')
      return false
    }
  } catch (error) {
    console.error('[TTS] Speak error:', error)
    showToast('语音播报失败: ' + (error.message || '未知错误'))
    isSpeaking.value = false
    return false
  }
}

// 使用 Web Speech API 播放语音
function speakWithWebAPI(plainText) {
  console.log('[TTS] Speaking with Web Speech API:', plainText.substring(0, 50) + '...')

  const utterance = new SpeechSynthesisUtterance(plainText)
  utterance.lang = 'zh-CN'
  utterance.rate = 1.0
  utterance.pitch = 1.0
  utterance.volume = 1.0

  // 尝试找到中文语音
  const voices = window.speechSynthesis.getVoices()
  const zhVoice = voices.find(v => v.lang.includes('zh') || v.lang.includes('CN'))
  if (zhVoice) {
    utterance.voice = zhVoice
  }

  utterance.onstart = () => {
    isSpeaking.value = true
    console.log('[TTS] Web TTS started')
  }

  utterance.onend = () => {
    isSpeaking.value = false
    currentSpeakingIndex.value = -1
    console.log('[TTS] Web TTS ended')
  }

  utterance.onerror = (event) => {
    console.error('[TTS] Web TTS Error:', event.error)
    if (event.error !== 'interrupted' && event.error !== 'canceled') {
      showToast('语音播报失败: ' + event.error)
    }
    isSpeaking.value = false
    currentSpeakingIndex.value = -1
  }

  // 恢复暂停的 speechSynthesis
  if (window.speechSynthesis.paused) {
    window.speechSynthesis.resume()
  }

  window.speechSynthesis.speak(utterance)
  isSpeaking.value = true
  return true
}

// 停止播放
async function stopSpeaking() {
  try {
    if (Capacitor.isNativePlatform()) {
      await TextToSpeech.stop()
    } else if ('speechSynthesis' in window) {
      window.speechSynthesis.cancel()
    }
  } catch (error) {
    console.error('[TTS] Stop error:', error)
  }
  isSpeaking.value = false
  currentSpeakingIndex.value = -1
}

// 切换播放/暂停
function toggleSpeak(text, index) {
  if (currentSpeakingIndex.value === index && isSpeaking.value) {
    // 当前正在播放这条消息，停止播放
    stopSpeaking()
  } else {
    // 播放这条消息
    stopSpeaking()
    if (speakText(text)) {
      currentSpeakingIndex.value = index
    }
  }
}

// 自动播放最后一条 AI 消息
function autoSpeakLastMessage() {
  if (!autoSpeak.value) return

  // 找到最后一条 AI 消息
  for (let i = messages.value.length - 1; i >= 0; i--) {
    if (messages.value[i].role === 'assistant' && messages.value[i].content) {
      // 延迟一点播放，确保消息已渲染
      setTimeout(() => {
        if (autoSpeak.value) {
          speakText(messages.value[i].content)
          currentSpeakingIndex.value = i
        }
      }, 300)
      break
    }
  }
}

// ============ 语音对话模式功能 ============

// 切换语音对话模式
function toggleVoiceChatMode() {
  voiceChatMode.value = !voiceChatMode.value
  if (!voiceChatMode.value && isListening.value) {
    stopListening()
  }
  // 切换到语音对话模式时自动开启自动播放
  if (voiceChatMode.value) {
    autoSpeak.value = true
  }
  showToast(voiceChatMode.value ? '已切换到语音对话模式' : '已切换到文字模式')
}

// 开始监听（语音对话模式）
async function startListening() {
  if (loading.value) return

  try {
    // 获取麦克风权限
    mediaStream = await navigator.mediaDevices.getUserMedia({ audio: true })

    // 创建音频上下文用于静音检测
    audioContext = new (window.AudioContext || window.webkitAudioContext)()
    analyser = audioContext.createAnalyser()
    const source = audioContext.createMediaStreamSource(mediaStream)
    source.connect(analyser)
    analyser.fftSize = 256

    isListening.value = true
    listeningText.value = ''

    // 开始分片录音
    await startStreamingRecording()

    // 开始静音检测
    startSilenceDetection()

  } catch (error) {
    console.error('启动语音监听失败:', error)
    showToast('启动失败: ' + (error.message || '请检查麦克风权限'))
    stopListening()
  }
}

// 开始分片录音（实时识别）
async function startStreamingRecording() {
  let mimeType = 'audio/webm'
  if (!MediaRecorder.isTypeSupported(mimeType)) {
    mimeType = 'audio/webm;codecs=opus'
    if (!MediaRecorder.isTypeSupported(mimeType)) {
      mimeType = 'audio/mp4'
    }
  }

  const options = mimeType ? { mimeType } : {}
  mediaRecorder.value = new MediaRecorder(mediaStream, options)
  audioChunks.value = []

  // 每 1.5 秒发送一个片段进行实时识别
  mediaRecorder.value.ondataavailable = async (event) => {
    if (event.data.size > 0) {
      audioChunks.value.push(event.data)
      // 实时发送片段进行识别
      if (audioChunks.value.length > 0) {
        const partialBlob = new Blob([event.data], { type: mimeType })
        await sendPartialVoice(partialBlob, mimeType)
      }
    }
  }

  // 每 1.5 秒触发一次 ondataavailable
  mediaRecorder.value.start(1500)
  recordingStartTime.value = Date.now()
}

// 发送部分语音进行实时识别
async function sendPartialVoice(audioBlob, mimeType) {
  try {
    if (!webSocketService.isConnected()) return

    const base64Audio = await blobToBase64(audioBlob)

    // 发送流式语音数据
    webSocketService.send({
      type: 'voice_stream_chunk',
      data: base64Audio,
      mimeType: mimeType
    })
  } catch (error) {
    console.error('发送语音片段失败:', error)
  }
}

// 开始语音对话模式的录音
async function startVoiceChatRecording() {
  if (Capacitor.isNativePlatform()) {
    await VoiceRecorder.startRecording()
  } else {
    let mimeType = 'audio/webm'
    if (!MediaRecorder.isTypeSupported(mimeType)) {
      mimeType = 'audio/webm;codecs=opus'
      if (!MediaRecorder.isTypeSupported(mimeType)) {
        mimeType = 'audio/mp4'
      }
    }

    const options = mimeType ? { mimeType } : {}
    mediaRecorder.value = new MediaRecorder(mediaStream, options)
    audioChunks.value = []

    mediaRecorder.value.ondataavailable = (event) => {
      if (event.data.size > 0) {
        audioChunks.value.push(event.data)
      }
    }

    mediaRecorder.value.start()
  }
  recordingStartTime.value = Date.now()
}

// 静音检测
function startSilenceDetection() {
  const dataArray = new Uint8Array(analyser.frequencyBinCount)
  let wasSpeaking = false

  const checkAudio = () => {
    if (!isListening.value) return

    analyser.getByteFrequencyData(dataArray)

    // 计算平均音量
    const average = dataArray.reduce((a, b) => a + b, 0) / dataArray.length

    const isSpeakingNow = average > 20 // 音量阈值

    if (isSpeakingNow) {
      // 正在说话
      wasSpeaking = true
      silenceStartTime = 0
      listeningText.value = '正在说话...'
    } else if (wasSpeaking) {
      // 曾经说话，现在静音
      if (silenceStartTime === 0) {
        silenceStartTime = Date.now()
      }

      const silenceDuration = (Date.now() - silenceStartTime) / 1000

      if (silenceDuration >= SILENCE_THRESHOLD) {
        // 静音超过阈值，自动发送
        listeningText.value = '检测到静音，发送中...'
        stopListeningAndSend()
        return
      } else {
        listeningText.value = `检测到静音 (${SILENCE_THRESHOLD - silenceDuration.toFixed(1)}s)`
      }
    }

    requestAnimationFrame(checkAudio)
  }

  requestAnimationFrame(checkAudio)
}

// 停止监听并发送
async function stopListeningAndSend() {
  if (!isListening.value) return

  isListening.value = false

  // 停止静音检测计时器
  if (silenceTimer) {
    clearTimeout(silenceTimer)
    silenceTimer = null
  }

  try {
    let audioBlob
    let mimeType = 'audio/webm'

    if (Capacitor.isNativePlatform()) {
      const result = await VoiceRecorder.stopRecording()
      audioBlob = base64ToBlob(result.value.recordDataBase64, result.value.mimeType)
      mimeType = result.value.mimeType
    } else if (mediaRecorder.value) {
      audioBlob = await new Promise((resolve) => {
        mediaRecorder.value.onstop = () => {
          mimeType = mediaRecorder.value.mimeType || 'audio/webm'
          resolve(new Blob(audioChunks.value, { type: mimeType }))
        }
        mediaRecorder.value.stop()
      })
    }

    // 检查录音时长
    const duration = Date.now() - recordingStartTime.value
    if (duration < 500) {
      showToast('说话时间太短')
      cleanupListening()
      return
    }

    // 发送语音
    await sendVoiceViaWebSocket(audioBlob, mimeType)

  } catch (error) {
    console.error('停止录音失败:', error)
    showToast('处理失败')
  }

  cleanupListening()
}

// 停止监听（手动）
function stopListening() {
  if (!isListening.value) return

  isListening.value = false

  // 停止录音
  if (Capacitor.isNativePlatform()) {
    VoiceRecorder.stopRecording().catch(() => {})
  } else if (mediaRecorder.value && mediaRecorder.value.state !== 'inactive') {
    mediaRecorder.value.stop()
  }

  cleanupListening()
}

// 清理监听资源
function cleanupListening() {
  if (mediaStream) {
    mediaStream.getTracks().forEach(track => track.stop())
    mediaStream = null
  }
  if (audioContext) {
    audioContext.close()
    audioContext = null
  }
  analyser = null
  silenceStartTime = 0
  listeningText.value = ''
}

onMounted(() => {
  loadHistory() // 加载历史记录
  registerAiCallback()
  initTTS() // 初始化 TTS
  fetchProviders() // 加载模型列表
})

onUnmounted(() => {
  // 清理 AI 回调
  if (removeAiCallback) {
    removeAiCallback()
    removeAiCallback = null
  }
  // 清理超时定时器
  if (responseTimeoutTimer) {
    clearTimeout(responseTimeoutTimer)
    responseTimeoutTimer = null
  }
  // 停止打字机效果
  if (typewriterTimer) {
    clearInterval(typewriterTimer)
    typewriterTimer = null
  }
  // 停止语音播放
  stopSpeaking()
  // 清理语音对话模式资源
  cleanupListening()
})

// 监听弹窗关闭，停止播放和监听
watch(visible, (val) => {
  if (!val) {
    stopSpeaking()
    if (isListening.value) {
      stopListening()
    }
    // 重置处理状态，确保下次打开时不会卡住
    resetProcessingState()
  }
})
</script>

<style scoped>
.ai-chat-popup-wrapper :deep(.van-popup__close-icon) {
  color: #969799;
}

.ai-chat-popup {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #f7f8fa;
}

.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
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

.speak-icon {
  font-size: 20px;
  color: #969799;
  cursor: pointer;
  transition: color 0.2s;
}

.speak-icon.active {
  color: #1989fa;
}

.voice-mode-icon {
  font-size: 20px;
  color: #969799;
  cursor: pointer;
  transition: color 0.2s;
}

.voice-mode-icon.active {
  color: #07c160;
}

.clear-icon {
  font-size: 20px;
  color: #969799;
  cursor: pointer;
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
  background: #f7f8fa;
  border-radius: 12px;
  padding: 16px;
  font-size: 14px;
  color: #323233;
  line-height: 1.8;
  white-space: pre-line;
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
  background: linear-gradient(135deg, #1989fa 0%, #0d7ce9 100%);
  color: white;
}

.avatar.user {
  background: linear-gradient(135deg, #07c160 0%, #06ad56 100%);
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
  background: linear-gradient(135deg, #1989fa 0%, #0d7ce9 100%);
  color: white;
}

.play-btn {
  font-size: 22px;
  color: #1989fa;
  cursor: pointer;
  flex-shrink: 0;
  margin-left: 4px;
}

.play-btn:active {
  transform: scale(0.9);
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
  color: #1989fa;
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

/* 语音对话模式 */
.voice-chat-controls {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  width: 100%;
}

.voice-chat-btn {
  min-width: 120px;
}

.listening-hint {
  font-size: 14px;
  color: #646566;
}

.chat-input {
  display: flex;
  gap: 8px;
  padding: 12px;
  background: #fff;
  border-top: 1px solid #eee;
  align-items: center;
}

.input-field {
  flex: 1;
  background: #f5f7fa;
  border-radius: 20px;
  padding: 0 12px;
}

.input-field :deep(.van-field__control) {
  background: transparent;
}

.send-btn {
  border-radius: 16px;
  padding: 0 12px;
}

/* 语音/键盘切换按钮 */
.mode-switch-btn {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: #f5f7fa;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
  color: #646566;
  flex-shrink: 0;
}

.mode-switch-btn:active {
  background: #e8e8e8;
  transform: scale(0.95);
}

/* 语音输入按钮（类似微信） */
.voice-input-btn {
  flex: 1;
  height: 40px;
  background: #fff;
  border: 1px solid #e8e8e8;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 15px;
  color: #323233;
  cursor: pointer;
  user-select: none;
  transition: all 0.15s;
}

.voice-input-btn:active,
.voice-input-btn.pressing {
  background: #c6c6c6;
  border-color: #c6c6c6;
}

/* 录音遮罩 */
.recording-overlay {
  display: flex;
  align-items: center;
  justify-content: center;
}

.recording-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 30px;
  background: rgba(0, 0, 0, 0.7);
  border-radius: 16px;
}

.recording-waves {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  height: 60px;
  margin-bottom: 20px;
}

.wave {
  width: 4px;
  height: 20px;
  background: #1989fa;
  border-radius: 2px;
  animation: wave 0.6s ease-in-out infinite;
}

@keyframes wave {
  0%, 100% {
    height: 20px;
  }
  50% {
    height: 50px;
  }
}

.recording-text {
  color: #fff;
  font-size: 18px;
  margin-bottom: 10px;
}

.recording-hint {
  color: rgba(255, 255, 255, 0.7);
  font-size: 14px;
}

/* 当前选中模型样式 */
:deep(.active-model) {
  color: #1989fa;
  font-weight: 500;
}

:deep(.active-model .van-action-sheet__name) {
  color: #1989fa;
}
</style>
