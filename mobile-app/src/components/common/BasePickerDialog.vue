<template>
  <van-popup
    :model-value="modelValue"
    :position="position"
    :style="{ height: height }"
    @update:model-value="handleClose"
  >
    <div class="base-picker-dialog">
      <!-- 标题栏 -->
      <div class="dialog-header">
        <van-button
          type="default"
          size="small"
          plain
          @click="handleCancel"
        >
          取消
        </van-button>
        <div class="dialog-title">{{ title }}</div>
        <van-button
          type="primary"
          size="small"
          :disabled="confirmDisabled"
          @click="handleConfirm"
        >
          {{ confirmText }}
        </van-button>
      </div>

      <!-- 内容区插槽 -->
      <div class="dialog-content">
        <slot :close="handleClose">
          <!-- 默认内容 -->
          <van-empty description="请提供内容" />
        </slot>
      </div>
    </div>
  </van-popup>
</template>

<script setup>
const props = defineProps({
  /**
   * 显示状态（v-model）
   */
  modelValue: {
    type: Boolean,
    default: false
  },
  /**
   * 标题
   */
  title: {
    type: String,
    default: '请选择'
  },
  /**
   * 弹窗位置
   */
  position: {
    type: String,
    default: 'bottom'
  },
  /**
   * 弹窗高度
   */
  height: {
    type: String,
    default: '60%'
  },
  /**
   * 确认按钮文本
   */
  confirmText: {
    type: String,
    default: '确认'
  },
  /**
   * 是否禁用确认按钮
   */
  confirmDisabled: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'confirm', 'cancel', 'close'])

function handleClose(value) {
  emit('update:modelValue', value)
  emit('close')
}

function handleCancel() {
  emit('cancel')
  handleClose(false)
}

function handleConfirm() {
  emit('confirm')
  // 注意：不自动关闭弹窗，让父组件决定何时关闭
  // 如果需要自动关闭，父组件可以在 confirm 事件处理中设置 modelValue = false
}
</script>

<style scoped>
.base-picker-dialog {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: #f7f8fa;
}

.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background-color: #fff;
  border-bottom: 1px solid #ebedf0;
  flex-shrink: 0;
}

.dialog-title {
  font-size: 16px;
  font-weight: 500;
  color: #323233;
}

.dialog-content {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
}
</style>
