<template>
  <van-tag :type="tagType">
    {{ displayText }}
  </van-tag>
</template>

<script setup>
import { computed } from 'vue'
import { getStatusType, getStatusText } from '@/composables/useStatus'

const props = defineProps({
  /**
   * 状态值
   */
  status: {
    type: String,
    default: ''
  },
  /**
   * 业务类型
   * 可选值: 'inbound', 'plan', 'requisition', 'appointment'
   */
  type: {
    type: String,
    default: 'inbound',
    validator: (value) => ['inbound', 'plan', 'requisition', 'appointment'].includes(value)
  },
  /**
   * 自定义显示文本（优先级高于自动映射）
   */
  text: {
    type: String,
    default: ''
  },
  /**
   * 自定义标签类型（优先级高于自动映射）
   * 可选值: 'default', 'primary', 'success', 'warning', 'danger'
   */
  tagType: {
    type: String,
    default: '' // 空字符串表示使用自动映射
  }
})

const computedType = computed(() => props.tagType || getStatusType(props.status, props.type))
const displayText = computed(() => props.text || getStatusText(props.status, props.type))
</script>

<style scoped>
/* 可以添加自定义样式 */
</style>
