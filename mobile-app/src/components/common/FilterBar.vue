<template>
  <div class="filter-bar">
    <van-dropdown-menu>
      <!-- 筛选项 -->
      <van-dropdown-item
        v-for="filter in filters"
        :key="filter.key"
        v-model="filterValues[filter.key]"
        :options="filter.options"
        :title="filter.title"
        @change="handleFilterChange(filter.key, $event)"
      />
    </van-dropdown-menu>

    <!-- 已选筛选标签 -->
    <div v-if="hasActiveFilters" class="active-filters">
      <van-tag
        v-for="filter in activeFilters"
        :key="filter.key"
        closeable
        type="primary"
        size="medium"
        @close="clearFilter(filter.key)"
      >
        {{ filter.label }}
      </van-tag>
      <van-button
        type="default"
        size="mini"
        @click="clearAllFilters"
      >
        清空
      </van-button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'

const props = defineProps({
  filters: {
    type: Array,
    default: () => [],
    // 每个筛选项的结构：
    // {
    //   key: 'status',        // 筛选键名
    //   title: '状态',        // 显示标题
    //   options: [            // 选项列表
    //     { text: '全部', value: '' },
    //     { text: '待审批', value: 'pending' },
    //   ]
    // }
  },
  modelValue: {
    type: Object,
    default: () => ({}),
  },
})

const emit = defineEmits(['update:modelValue', 'change', 'filter', 'clear'])

const filterValues = ref({})

// 初始化筛选值
function initFilterValues() {
  const values = {}
  props.filters.forEach(filter => {
    values[filter.key] = props.modelValue[filter.key] || ''
  })
  filterValues.value = values
}

// 监听外部值变化
watch(() => props.modelValue, (newVal) => {
  Object.keys(newVal).forEach(key => {
    if (filterValues.value.hasOwnProperty(key)) {
      filterValues.value[key] = newVal[key]
    }
  })
}, { deep: true })

// 监听筛选项变化
watch(() => props.filters, () => {
  initFilterValues()
}, { immediate: true })

// 筛选变化处理
function handleFilterChange(key, value) {
  emit('update:modelValue', { ...filterValues.value })
  emit('change', { key, value })
  emit('filter', { ...filterValues.value })
}

// 清空单个筛选
function clearFilter(key) {
  filterValues.value[key] = ''
  emit('update:modelValue', { ...filterValues.value })
  emit('clear', key)
  emit('filter', { ...filterValues.value })
}

// 清空所有筛选
function clearAllFilters() {
  Object.keys(filterValues.value).forEach(key => {
    filterValues.value[key] = ''
  })
  emit('update:modelValue', { ...filterValues.value })
  emit('clear', null)
  emit('filter', { ...filterValues.value })
}

// 获取当前激活的筛选
const activeFilters = computed(() => {
  const filters = []
  props.filters.forEach(filter => {
    const value = filterValues.value[filter.key]
    if (value) {
      const option = filter.options.find(opt => opt.value === value)
      if (option) {
        filters.push({
          key: filter.key,
          label: option.text,
        })
      }
    }
  })
  return filters
})

// 是否有激活的筛选
const hasActiveFilters = computed(() => {
  return activeFilters.value.length > 0
})

// 暴露方法
defineExpose({
  clearAllFilters,
  clearFilter,
  getFilters: () => filterValues.value,
  setFilter: (key, value) => {
    if (filterValues.value.hasOwnProperty(key)) {
      filterValues.value[key] = value
      emit('update:modelValue', { ...filterValues.value })
    }
  },
})
</script>

<style scoped>
.filter-bar {
  background: white;
  position: sticky;
  top: 0;
  z-index: 10;
}

.active-filters {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 12px 16px;
  background: #f7f8fa;
  border-top: 1px solid #ebedf0;
}

.active-filters .van-tag {
  margin: 0;
}
</style>
