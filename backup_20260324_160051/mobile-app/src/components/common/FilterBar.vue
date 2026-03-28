<template>
  <van-dropdown-menu>
    <van-dropdown-item
      v-for="filter in filters"
      :key="filter.key"
      v-model="filterValues[filter.key]"
      :options="filter.options"
      @change="handleFilterChange(filter.key, $event)"
    >
      <!-- 支持自定义插槽内容 -->
      <template v-if="filter.custom" #title>
        <slot :name="`filter-${filter.key}`" :filter="filter" :value="filterValues[filter.key]">
          {{ getFilterLabel(filter) }}
        </slot>
      </template>
    </van-dropdown-item>
  </van-dropdown-menu>
</template>

<script setup>
import { ref, watch, computed } from 'vue'

const props = defineProps({
  /**
   * 筛选器配置
   * 示例: [
   *   { key: 'status', options: [{ text: '全部', value: '' }] },
   *   { key: 'type', options: [{ text: '类型A', value: 'a' }] }
   * ]
   */
  filters: {
    type: Array,
    required: true,
    default: () => []
  },
  /**
   * 初始筛选值
   * 示例: { status: '', type: 'a' }
   */
  modelValue: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['update:modelValue', 'filter-change'])

// 筛选值
const filterValues = ref({})

// 初始化筛选值
const initFilterValues = () => {
  const values = {}
  props.filters.forEach(filter => {
    values[filter.key] = props.modelValue?.[filter.key] !== undefined
      ? props.modelValue[filter.key]
      : (filter.defaultValue !== undefined ? filter.defaultValue : '')
  })
  filterValues.value = values
}

// 监听 modelValue 变化
watch(() => props.modelValue, (newVal) => {
  if (newVal) {
    Object.keys(newVal).forEach(key => {
      if (filterValues.value.hasOwnProperty(key)) {
        filterValues.value[key] = newVal[key]
      }
    })
  }
}, { deep: true })

// 监听 filters 变化，重新初始化
watch(() => props.filters, () => {
  initFilterValues()
}, { immediate: true })

// 获取筛选器显示标签
function getFilterLabel(filter) {
  const value = filterValues.value[filter.key]
  const option = filter.options?.find(opt => opt.value === value)
  return option?.text || filter.title || filter.key
}

// 处理筛选器变化
function handleFilterChange(key, value) {
  emit('update:modelValue', { ...filterValues.value })
  emit('filter-change', { key, value, allValues: { ...filterValues.value } })
}
</script>

<style scoped>
/* 可以添加自定义样式 */
</style>
