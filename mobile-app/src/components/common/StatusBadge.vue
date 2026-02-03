<template>
  <van-tag
    :type="tagType"
    :size="size"
    :plain="plain"
    :round="round"
    :class="['status-badge', `status-${status}`]"
  >
    <slot>{{ displayText }}</slot>
  </van-tag>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  status: {
    type: String,
    required: true,
  },
  size: {
    type: String,
    default: 'medium', // large, medium, small
  },
  plain: {
    type: Boolean,
    default: false,
  },
  round: {
    type: Boolean,
    default: false,
  },
  // 自定义状态文本映射
  statusText: {
    type: Object,
    default: () => ({}),
  },
  // 自定义状态类型映射
  statusType: {
    type: Object,
    default: () => ({}),
  },
})

// 预定义状态文本
const defaultStatusText = {
  // 通用状态
  draft: '草稿',
  pending: '待审批',
  approved: '已通过',
  rejected: '已拒绝',
  processing: '处理中',
  completed: '已完成',
  cancelled: '已取消',
  active: '进行中',
  inactive: '已停用',

  // 入库状态
  inbound_pending: '待入库',
  inbound_completed: '已入库',
  inbound_cancelled: '已取消',

  // 出库状态
  outbound_pending: '待出库',
  outbound_approved: '已审批',
  outbound_issued: '已发料',
  outbound_cancelled: '已取消',

  // 物资计划状态
  plan_draft: '草稿',
  plan_pending: '待审批',
  plan_approved: '已审批',
  plan_active: '执行中',
  plan_completed: '已完成',

  // 支付状态
  unpaid: '待支付',
  paid: '已支付',
  refunded: '已退款',

  // 发货状态
  unshipped: '待发货',
  shipped: '已发货',
  delivered: '已签收',

  // 库存状态
  in_stock: '有货',
  low_stock: '库存不足',
  out_of_stock: '缺货',
}

// 预定义状态类型（Vant Tag type）
const defaultStatusType = {
  // 成功/正常状态
  approved: 'success',
  completed: 'success',
  active: 'success',
  paid: 'success',
  delivered: 'success',
  in_stock: 'success',
  plan_active: 'success',
  inbound_completed: 'success',

  // 警告/待处理状态
  draft: 'default',
  pending: 'warning',
  processing: 'warning',
  unpaid: 'warning',
  unshipped: 'warning',
  plan_pending: 'warning',
  outbound_pending: 'warning',
  inbound_pending: 'warning',

  // 危险/异常状态
  rejected: 'danger',
  cancelled: 'danger',
  inactive: 'danger',
  refunded: 'danger',
  out_of_stock: 'danger',
  outbound_cancelled: 'danger',
  inbound_cancelled: 'danger',

  // 主要状态
  plan_approved: 'primary',
  outbound_approved: 'primary',
  plan_completed: 'success',
}

// 显示文本
const displayText = computed(() => {
  return props.statusText[props.status] || defaultStatusText[props.status] || props.status
})

// Tag类型
const tagType = computed(() => {
  return props.statusType[props.status] || defaultStatusType[props.status] || 'default'
})
</script>

<style scoped>
.status-badge {
  font-weight: 500;
}

/* 状态特定样式 */
.status-draft {
  opacity: 0.8;
}

.status-pending {
  animation: pulse 2s infinite;
}

.status-active {
  position: relative;
}

.status-active::after {
  content: '';
  position: absolute;
  right: -4px;
  top: -4px;
  width: 6px;
  height: 6px;
  background: currentColor;
  border-radius: 50%;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.7;
  }
}
</style>
