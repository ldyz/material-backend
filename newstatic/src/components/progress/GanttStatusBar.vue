<template>
  <div class="gantt-status-bar" v-if="visible">
    <div class="status-left">
      <div class="status-item">
        <el-icon class="status-icon"><Folder /></el-icon>
        <span class="status-label">项目:</span>
        <span class="status-value">{{ projectName }}</span>
      </div>
      <div class="status-divider"></div>
      <div class="status-item">
        <el-icon class="status-icon"><List /></el-icon>
        <span class="status-label">任务:</span>
        <span class="status-value">{{ stats.total }} 个</span>
      </div>
      <div class="status-divider"></div>
      <div class="status-item">
        <el-icon class="status-icon" style="color: #67c23a;"><CircleCheck /></el-icon>
        <span class="status-label">已完成:</span>
        <span class="status-value success">{{ stats.completed }}</span>
      </div>
      <div class="status-divider"></div>
      <div class="status-item">
        <el-icon class="status-icon" style="color: #409eff;"><Clock /></el-icon>
        <span class="status-label">进行中:</span>
        <span class="status-value primary">{{ stats.inProgress }}</span>
      </div>
      <div class="status-divider"></div>
      <div class="status-item" v-if="stats.critical > 0">
        <el-icon class="status-icon" style="color: #f56c6c;"><Flag /></el-icon>
        <span class="status-label">关键路径:</span>
        <span class="status-value critical">{{ stats.critical }}</span>
      </div>
    </div>

    <div class="status-center">
      <transition name="fade">
        <div v-if="operationStatus.text" class="operation-status" :class="operationStatus.type">
          <el-icon v-if="operationStatus.type === 'loading'" class="status-icon-spin">
            <Loading />
          </el-icon>
          <el-icon v-else-if="operationStatus.type === 'success'">
            <CircleCheck />
          </el-icon>
          <el-icon v-else-if="operationStatus.type === 'error'">
            <CircleClose />
          </el-icon>
          <span>{{ operationStatus.text }}</span>
        </div>
      </transition>
    </div>

    <div class="status-right">
      <div class="status-item">
        <span class="status-label">进度:</span>
        <el-progress
          :percentage="stats.progressRate"
          :stroke-width="8"
          :show-text="true"
          style="width: 150px; margin-left: 8px;"
        />
      </div>
      <div class="status-divider"></div>
      <div class="status-item">
        <span class="status-label">更新时间:</span>
        <span class="status-value">{{ lastUpdateTime }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, onUnmounted } from 'vue'
import {
  Folder,
  List,
  CircleCheck,
  Clock,
  Flag,
  Loading,
  CircleClose
} from '@element-plus/icons-vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  projectName: {
    type: String,
    default: '未命名项目'
  },
  stats: {
    type: Object,
    default: () => ({
      total: 0,
      completed: 0,
      inProgress: 0,
      notStarted: 0,
      critical: 0,
      progressRate: 0
    })
  },
  isDragging: {
    type: Boolean,
    default: false
  },
  isSaving: {
    type: Boolean,
    default: false
  }
})

// 操作状态
const operationStatus = ref({
  text: '',
  type: '' // 'loading' | 'success' | 'error'
})

let statusTimeout = null

// 显示操作状态
const showStatus = (text, type = 'loading', duration = 0) => {
  operationStatus.value = { text, type }

  if (duration > 0) {
    if (statusTimeout) clearTimeout(statusTimeout)
    statusTimeout = setTimeout(() => {
      operationStatus.value = { text: '', type: '' }
    }, duration)
  }
}

// 隐藏状态
const hideStatus = () => {
  operationStatus.value = { text: '', type: '' }
}

// 更新时间
const lastUpdateTime = ref('')

const updateLastTime = () => {
  const now = new Date()
  lastUpdateTime.value = `${String(now.getHours()).padStart(2, '0')}:${String(now.getMinutes()).padStart(2, '0')}:${String(now.getSeconds()).padStart(2, '0')}`
}

// 每秒更新时间
const timeInterval = setInterval(updateLastTime, 1000)
updateLastTime()

onUnmounted(() => {
  if (statusTimeout) clearTimeout(statusTimeout)
  if (timeInterval) clearInterval(timeInterval)
})

// 监听拖拽状态
const draggingText = computed(() => {
  if (props.isDragging) {
    return '正在拖拽任务...'
  }
  return ''
})

// 监听保存状态
const savingText = computed(() => {
  if (props.isSaving) {
    return '正在保存...'
  }
  return ''
})

// 暴露方法
defineExpose({
  showStatus,
  hideStatus
})
</script>

<style scoped>
.gantt-status-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 48px;
  background: #fff;
  border-top: 1px solid #dcdfe6;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  z-index: 9999;
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.1);
}

.status-left,
.status-center,
.status-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.status-center {
  flex: 1;
  justify-content: center;
}

.status-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.status-icon {
  font-size: 16px;
  color: #606266;
}

.status-icon-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.status-label {
  font-size: 12px;
  color: #909399;
}

.status-value {
  font-size: 13px;
  font-weight: 500;
  color: #303133;
}

.status-value.success {
  color: #67c23a;
}

.status-value.primary {
  color: #409eff;
}

.status-value.critical {
  color: #f56c6c;
}

.status-divider {
  width: 1px;
  height: 20px;
  background: #dcdfe6;
}

.operation-status {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 16px;
  border-radius: 4px;
  font-size: 13px;
}

.operation-status.loading {
  background: #ecf5ff;
  color: #409eff;
}

.operation-status.success {
  background: #f0f9ff;
  color: #67c23a;
}

.operation-status.error {
  background: #fef0f0;
  color: #f56c6c;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s, transform 0.3s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
