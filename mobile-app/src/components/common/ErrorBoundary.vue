<template>
  <div class="error-boundary">
    <slot v-if="!hasError" />

    <div v-else class="error-fallback">
      <van-empty
        image="error"
        :description="errorMessage"
      >
        <template #description>
          <p class="error-message">{{ errorMessage }}</p>
          <p v-if="showDetails && errorDetails" class="error-details">
            {{ errorDetails }}
          </p>
        </template>

        <van-button type="primary" size="small" @click="handleRetry">
          重试
        </van-button>
        <van-button v-if="onReset" size="small" @click="handleReset">
          返回首页
        </van-button>
      </van-empty>
    </div>
  </div>
</template>

<script setup>
import { ref, onErrorCaptured } from 'vue'

const props = defineProps({
  // 是否显示错误详情
  showDetails: {
    type: Boolean,
    default: false,
  },
  // 自定义错误消息
  fallbackMessage: {
    type: String,
    default: '出错了',
  },
})

const emit = defineEmits(['error', 'retry', 'reset'])

const hasError = ref(false)
const errorMessage = ref('')
const errorDetails = ref('')

// 捕获子组件错误
onErrorCaptured((error, instance, info) => {
  console.error('[Error Boundary]', error, info)

  hasError.value = true
  errorMessage.value = error?.message || props.fallbackMessage
  errorDetails.value = error?.stack || ''

  emit('error', error, info)

  // 阻止错误继续传播
  return false
})

// 重试
function handleRetry() {
  hasError.value = false
  errorMessage.value = ''
  errorDetails.value = ''
  emit('retry')
}

// 重置
function handleReset() {
  hasError.value = false
  errorMessage.value = ''
  errorDetails.value = ''
  emit('reset')
}

// 暴露方法
defineExpose({
  handleRetry,
  handleReset,
})
</script>

<style scoped>
.error-boundary {
  min-height: 100%;
}

.error-fallback {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  padding: 20px;
}

.error-message {
  color: #323233;
  font-size: 16px;
  font-weight: 500;
  margin: 12px 0;
}

.error-details {
  color: #969799;
  font-size: 12px;
  max-width: 300px;
  word-break: break-all;
  margin: 8px 0 16px;
}

.error-fallback :deep(.van-button) {
  margin: 0 4px;
}
</style>
