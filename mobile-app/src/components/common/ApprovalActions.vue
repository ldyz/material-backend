<template>
  <div class="approval-actions">
    <van-button
      round
      block
      type="danger"
      plain
      :loading="rejecting"
      :disabled="disabled"
      @click="handleReject"
    >
      {{ rejectText }}
    </van-button>
    <van-button
      round
      block
      type="success"
      :loading="approving"
      :disabled="disabled"
      @click="handleApprove"
    >
      {{ approveText }}
    </van-button>
  </div>
</template>

<script setup>
const props = defineProps({
  /**
   * 审批中状态
   */
  approving: {
    type: Boolean,
    default: false
  },
  /**
   * 拒绝中状态
   */
  rejecting: {
    type: Boolean,
    default: false
  },
  /**
   * 是否禁用按钮
   */
  disabled: {
    type: Boolean,
    default: false
  },
  /**
   * 审批按钮文本
   */
  approveText: {
    type: String,
    default: '通过'
  },
  /**
   * 拒绝按钮文本
   */
  rejectText: {
    type: String,
    default: '拒绝'
  },
  /**
   * 按钮布局方向
   */
  direction: {
    type: String,
    default: 'horizontal',
    validator: (value) => ['horizontal', 'vertical'].includes(value)
  }
})

const emit = defineEmits(['approve', 'reject'])

function handleApprove() {
  emit('approve')
}

function handleReject() {
  emit('reject')
}
</script>

<style scoped>
.approval-actions {
  display: flex;
  gap: 12px;
  padding: 16px;
  background-color: #fff;
  border-top: 1px solid #ebedf0;
}

.approval-actions.vertical {
  flex-direction: column;
}
</style>
