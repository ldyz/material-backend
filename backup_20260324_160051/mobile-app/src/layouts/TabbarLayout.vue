<template>
  <div class="tabbar-layout">
    <router-view v-slot="{ Component, route }">
      <keep-alive>
        <component
          :is="Component"
          v-if="route.meta.keepAlive"
          :key="route.path"
        />
      </keep-alive>
      <component
        :is="Component"
        v-if="!route.meta.keepAlive"
        :key="route.path"
      />
    </router-view>

    <!-- 自定义底部导航栏 -->
    <div class="custom-tabbar">
      <!-- 首页按钮 -->
      <div
        class="tabbar-item"
        :class="{ active: activeTab === 'dashboard' }"
        @click="handleTabClick('dashboard')"
      >
        <van-icon name="wap-home-o" size="22" />
        <span class="tabbar-label">首页</span>
      </div>

      <!-- 中间语音按钮 -->
      <div class="voice-button-wrapper">
        <div
          class="voice-button"
          :class="{ pressing: isRecording, 'has-permission': hasRecordPermission }"
          @touchstart="startRecording"
          @touchend="stopRecording"
          @touchcancel="cancelRecording"
          @mousedown="startRecording"
          @mouseup="stopRecording"
          @mouseleave="cancelRecording"
        >
          <!-- 麦克风图标 -->
          <svg class="mic-icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 14C13.66 14 15 12.66 15 11V5C15 3.34 13.66 2 12 2C10.34 2 9 3.34 9 5V11C9 12.66 10.34 14 12 14Z" fill="currentColor"/>
            <path d="M17 11C17 13.76 14.76 16 12 16C9.24 16 7 13.76 7 11H5C5 14.53 7.61 17.43 11 17.92V21H13V17.92C16.39 17.43 19 14.53 19 11H17Z" fill="currentColor"/>
          </svg>
        </div>
        <span class="voice-label">{{ voiceLabel }}</span>
      </div>

      <!-- 我的按钮 -->
      <div
        class="tabbar-item"
        :class="{ active: activeTab === 'profile' }"
        @click="handleTabClick('profile')"
      >
        <van-icon name="user-o" size="22" />
        <span class="tabbar-label">我的</span>
      </div>
    </div>

    <!-- 语音输入动画遮罩 -->
    <van-overlay :show="isRecording" class="voice-overlay">
      <div class="voice-overlay-content">
        <div class="voice-waves">
          <div class="wave" v-for="i in 5" :key="i" :style="{ animationDelay: `${i * 0.1}s` }"></div>
        </div>
        <div class="voice-text">{{ recordingText }}</div>
        <div class="voice-hint">松开发送，上滑取消</div>
      </div>
    </van-overlay>

    <!-- AI 回复对话框 -->
    <van-dialog
      v-model:show="showAiDialog"
      title="AI 助手"
      :show-confirm-button="false"
      close-on-click-overlay
      class="ai-dialog"
    >
      <div class="ai-dialog-content">
        <div v-if="aiLoading" class="ai-loading">
          <van-loading size="24px" color="#1989fa">思考中...</van-loading>
        </div>
        <div v-else class="ai-response">
          <div class="ai-message" v-html="aiResponse"></div>
        </div>
        <van-button size="small" @click="showAiDialog = false" class="close-btn">关闭</van-button>
      </div>
    </van-dialog>

    <!-- 自动更新对话框 -->
    <AppUpdateDialog v-model:show="showUpdateDialog" />
  </div>
</template>

<script setup>
import { ref, watch, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAutoUpdate } from '@/composables/useAppUpdate'
import { useAuthStore } from '@/stores/auth'
import { showToast } from 'vant'
import AppUpdateDialog from '@/components/AppUpdateDialog.vue'
import { agentApi } from '@/api/agent'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const activeTab = ref('dashboard')
const showUpdateDialog = ref(false)

// 语音相关状态
const isRecording = ref(false)
const hasRecordPermission = ref(false)
const voiceLabel = ref('语音')
const recordingText = ref('正在录音...')
const mediaRecorder = ref(null)
const audioChunks = ref([])
const recordingStartTime = ref(0)
const touchStartY = ref(0)
const isCancelled = ref(false)

// AI 对话相关
const showAiDialog = ref(false)
const aiLoading = ref(false)
const aiResponse = ref('')

// 使用自动更新检测
const { hasUpdate, forceUpdate, performCheck } = useAutoUpdate({
  autoCheck: true,
  checkOnMount: true,
  checkInterval: 24 * 60 * 60 * 1000
})

// 处理标签点击
function handleTabClick(name) {
  if (name === activeTab.value) return
  activeTab.value = name
  const targetPath = name === 'dashboard' ? '/' : `/${name}`
  if (route.path !== targetPath) {
    router.push(targetPath)
  }
}

// 检查录音权限
async function checkRecordPermission() {
  try {
    const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
    stream.getTracks().forEach(track => track.stop())
    hasRecordPermission.value = true
    return true
  } catch (error) {
    hasRecordPermission.value = false
    showToast('请授予麦克风权限')
    return false
  }
}

// 开始录音
async function startRecording(event) {
  if (event.type === 'touchstart') {
    touchStartY.value = event.touches[0].clientY
  }

  if (!await checkRecordPermission()) return

  try {
    const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
    mediaRecorder.value = new MediaRecorder(stream, {
      mimeType: 'audio/webm'
    })
    audioChunks.value = []

    mediaRecorder.value.ondataavailable = (event) => {
      if (event.data.size > 0) {
        audioChunks.value.push(event.data)
      }
    }

    mediaRecorder.value.start()
    isRecording.value = true
    isCancelled.value = false
    recordingStartTime.value = Date.now()
    voiceLabel.value = '录音中'
    recordingText.value = '正在录音...'
  } catch (error) {
    console.error('录音启动失败:', error)
    showToast('录音启动失败')
  }
}

// 停止录音
async function stopRecording(event) {
  if (!isRecording.value || !mediaRecorder.value) return

  // 检查是否上滑取消
  if (event.type === 'touchend') {
    const touchEndY = event.changedTouches[0].clientY
    if (touchStartY.value - touchEndY > 100) {
      isCancelled.value = true
    }
  }

  if (isCancelled.value) {
    cancelRecording()
    return
  }

  return new Promise((resolve) => {
    mediaRecorder.value.onstop = async () => {
      const audioBlob = new Blob(audioChunks.value, { type: 'audio/webm' })
      const duration = Date.now() - recordingStartTime.value

      // 停止所有音轨
      mediaRecorder.value.stream.getTracks().forEach(track => track.stop())

      isRecording.value = false
      voiceLabel.value = '语音'

      // 录音时间太短
      if (duration < 500) {
        showToast('录音时间太短')
        resolve()
        return
      }

      // 发送语音进行识别
      await sendVoiceToServer(audioBlob)
      resolve()
    }

    mediaRecorder.value.stop()
  })
}

// 取消录音
function cancelRecording() {
  if (mediaRecorder.value && isRecording.value) {
    mediaRecorder.value.stream.getTracks().forEach(track => track.stop())
  }
  isRecording.value = false
  voiceLabel.value = '语音'
  isCancelled.value = false
}

// 发送语音到服务器
async function sendVoiceToServer(audioBlob) {
  try {
    aiLoading.value = true
    showAiDialog.value = true
    aiResponse.value = ''

    // 创建 FormData
    const formData = new FormData()
    formData.append('audio', audioBlob, 'recording.webm')

    // 调用语音识别和 AI 接口
    const response = await agentApi.voiceChat(formData)

    if (response.data) {
      aiResponse.value = formatAiResponse(response.data.response || response.data.message || '处理完成')
    } else {
      aiResponse.value = '未能识别语音内容，请重试'
    }
  } catch (error) {
    console.error('语音处理失败:', error)
    aiResponse.value = '处理失败: ' + (error.message || '请稍后重试')
  } finally {
    aiLoading.value = false
  }
}

// 格式化 AI 响应
function formatAiResponse(text) {
  if (!text) return ''
  // 简单的 Markdown 转换
  return text
    .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
    .replace(/\n/g, '<br>')
    .replace(/- (.*)/g, '<li>$1</li>')
}

// 监听是否有更新
watch([hasUpdate, forceUpdate], ([hasUpdate, forceUpdate]) => {
  if (hasUpdate) {
    if (forceUpdate) {
      showUpdateDialog.value = true
    } else {
      setTimeout(() => {
        showUpdateDialog.value = true
      }, 3000)
    }
  }
})

watch(
  () => route.path,
  (path) => {
    if (path === '/' || path.startsWith('/dashboard')) {
      activeTab.value = 'dashboard'
    } else if (path.startsWith('/profile')) {
      activeTab.value = 'profile'
    }
  },
  { immediate: true }
)

onMounted(() => {
  window.checkAppUpdate = performCheck
  // 预检查录音权限
  checkRecordPermission()
})
</script>

<style scoped>
.tabbar-layout {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding-bottom: 70px;
  padding-top: var(--capacitor-status-bar-height, 0px);
}

/* 自定义底部导航栏 */
.custom-tabbar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 60px;
  background: #fff;
  display: flex;
  justify-content: space-around;
  align-items: center;
  box-shadow: 0 -2px 12px rgba(0, 0, 0, 0.08);
  padding-bottom: env(safe-area-inset-bottom);
  z-index: 100;
}

.tabbar-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 8px 0;
  cursor: pointer;
  color: #969799;
  transition: color 0.2s;
}

.tabbar-item.active {
  color: #1989fa;
}

.tabbar-label {
  font-size: 12px;
  margin-top: 4px;
}

/* 中间语音按钮 */
.voice-button-wrapper {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: -25px;
}

.voice-button {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: linear-gradient(135deg, #1989fa 0%, #0d7ce9 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 12px rgba(25, 137, 250, 0.4);
  cursor: pointer;
  transition: all 0.2s;
  border: 3px solid #fff;
}

.voice-button:active,
.voice-button.pressing {
  transform: scale(1.1);
  box-shadow: 0 6px 20px rgba(25, 137, 250, 0.5);
}

.voice-button.pressing {
  background: linear-gradient(135deg, #07c160 0%, #06ad56 100%);
  box-shadow: 0 6px 20px rgba(7, 193, 96, 0.5);
}

.mic-icon {
  width: 28px;
  height: 28px;
  color: #fff;
}

.voice-label {
  font-size: 10px;
  color: #969799;
  margin-top: 4px;
}

/* 语音输入遮罩 */
.voice-overlay {
  display: flex;
  align-items: center;
  justify-content: center;
}

.voice-overlay-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 30px;
  background: rgba(0, 0, 0, 0.7);
  border-radius: 16px;
}

.voice-waves {
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

.voice-text {
  color: #fff;
  font-size: 18px;
  margin-bottom: 10px;
}

.voice-hint {
  color: rgba(255, 255, 255, 0.7);
  font-size: 14px;
}

/* AI 对话框 */
.ai-dialog :deep(.van-dialog__header) {
  background: linear-gradient(135deg, #1989fa 0%, #0d7ce9 100%);
  color: #fff;
  padding-top: 16px;
  padding-bottom: 16px;
}

.ai-dialog-content {
  padding: 16px;
  max-height: 60vh;
  overflow-y: auto;
}

.ai-loading {
  display: flex;
  justify-content: center;
  padding: 40px 0;
}

.ai-message {
  font-size: 14px;
  line-height: 1.8;
  color: #333;
  word-break: break-word;
}

.ai-message :deep(li) {
  margin-left: 16px;
  list-style: disc;
}

.close-btn {
  margin-top: 16px;
  width: 100%;
}

/* Capacitor 状态栏高度变量 */
:global(html) {
  --capacitor-status-bar-height: 0px;
}

@supports (padding: max(0px)) {
  .tabbar-layout {
    padding-top: max(var(--capacitor-status-bar-height), env(safe-area-inset-top));
  }
}
</style>
