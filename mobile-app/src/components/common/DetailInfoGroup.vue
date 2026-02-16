<template>
  <van-cell-group inset :title="title">
    <van-cell
      v-for="field in renderedFields"
      :key="field.key"
      :title="field.label"
      :value="getFieldValue(field)"
    >
      <template v-if="field.type === 'status'" #value>
        <StatusTag :status="getFieldValue(field)" :type="statusType" />
      </template>
      <template v-else-if="field.type === 'tag'" #value>
        <van-tag :type="field.tagType || 'primary'">
          {{ getFieldValue(field) }}
        </van-tag>
      </template>
      <template v-else-if="field.type === 'date'" #value>
        {{ formatDateValue(getFieldValue(field)) }}
      </template>
      <template v-else-if="field.type === 'datetime'" #value>
        {{ formatDateTimeValue(getFieldValue(field)) }}
      </template>
      <template v-else-if="field.type === 'custom'" #value>
        <slot :name="`field-${field.key}`" :field="field" :value="getFieldValue(field)">
          {{ getFieldValue(field) }}
        </slot>
      </template>
    </van-cell>
  </van-cell-group>
</template>

<script setup>
import { computed } from 'vue'
import StatusTag from './StatusTag.vue'
import { formatDate, formatDateTime, formatAppointmentDate } from '@/composables/useDateTime'

const props = defineProps({
  /**
   * 分组标题
   */
  title: {
    type: String,
    required: true
  },
  /**
   * 数据对象
   */
  item: {
    type: Object,
    required: true
  },
  /**
   * 字段配置
   * 示例: [
   *   { key: 'order_number', label: '入库单号' },
   *   { key: 'status', label: '状态', type: 'status' },
   *   { key: 'inbound_date', label: '入库日期', type: 'date' }
   * ]
   */
  fields: {
    type: Array,
    required: true
  },
  /**
   * 状态类型（用于状态标签）
   */
  statusType: {
    type: String,
    default: 'inbound'
  },
  /**
   * 是否显示空值字段
   */
  showEmpty: {
    type: Boolean,
    default: true
  }
})

// 过滤后的字段列表
const renderedFields = computed(() => {
  if (props.showEmpty) {
    return props.fields
  }

  return props.fields.filter(field => {
    const value = getFieldValue(field)
    return value !== null && value !== undefined && value !== ''
  })
})

// 获取字段值
function getFieldValue(field) {
  const value = props.item[field.key]

  // 如果有格式化函数，使用它
  if (field.formatter && typeof field.formatter === 'function') {
    return field.formatter(value, props.item)
  }

  // 返回原始值
  return value || field.defaultValue || '-'
}

// 格式化日期
function formatDateValue(value) {
  if (!value || value === '-') return '-'
  return formatDate(value)
}

// 格式化日期时间
function formatDateTimeValue(value) {
  if (!value || value === '-') return '-'
  return formatDateTime(value)
}
</script>

<style scoped>
/* 可以添加自定义样式 */
</style>
