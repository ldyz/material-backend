<template>
  <div class="search-bar-wrapper">
    <van-search
      v-model="searchValue"
      :placeholder="placeholder"
      :shape="shape"
      :background="background"
      :show-action="showAction"
      :clearable="clearable"
      @update:model-value="onInput"
      @search="onSearch"
      @cancel="onCancel"
      @focus="onFocus"
      @blur="onBlur"
    >
      <template v-if="$slots.left" #left>
        <slot name="left"></slot>
      </template>
      <template v-if="$slots.action" #action>
        <slot name="action"></slot>
      </template>
    </van-search>

    <!-- 快速筛选标签 -->
    <div v-if="showFilterTags && filterTags.length > 0" class="filter-tags">
      <van-tag
        v-for="tag in filterTags"
        :key="tag.value"
        :type="activeFilter === tag.value ? 'primary' : 'default'"
        :plain="activeFilter !== tag.value"
        size="medium"
        @click="onFilterTagClick(tag.value)"
      >
        {{ tag.label }}
      </van-tag>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useDebounceFn } from '@/composables/useDebounce'

const props = defineProps({
  modelValue: {
    type: String,
    default: '',
  },
  placeholder: {
    type: String,
    default: '请输入搜索关键词',
  },
  shape: {
    type: String,
    default: 'round', // round, square
  },
  background: {
    type: String,
    default: '#f7f8fa',
  },
  showAction: {
    type: Boolean,
    default: false,
  },
  clearable: {
    type: Boolean,
    default: true,
  },
  debounceDelay: {
    type: Number,
    default: 300,
  },
  showFilterTags: {
    type: Boolean,
    default: false,
  },
  filterTags: {
    type: Array,
    default: () => [],
  },
})

const emit = defineEmits(['update:modelValue', 'search', 'cancel', 'focus', 'blur', 'filter'])

const searchValue = ref(props.modelValue)
const activeFilter = ref(null)

// 监听外部值变化
watch(() => props.modelValue, (newVal) => {
  searchValue.value = newVal
})

// 防抖搜索
const debouncedSearch = useDebounceFn((value) => {
  emit('search', value)
}, props.debounceDelay)

// 输入处理
function onInput(value) {
  emit('update:modelValue', value)
  debouncedSearch(value)
}

// 搜索确认
function onSearch(value) {
  emit('search', value)
}

// 取消搜索
function onCancel() {
  searchValue.value = ''
  emit('update:modelValue', '')
  emit('cancel')
}

// 获得焦点
function onFocus(event) {
  emit('focus', event)
}

// 失去焦点
function onBlur(event) {
  emit('blur', event)
}

// 筛选标签点击
function onFilterTagClick(value) {
  if (activeFilter.value === value) {
    activeFilter.value = null
    emit('filter', null)
  } else {
    activeFilter.value = value
    emit('filter', value)
  }
}

// 清空筛选
function clearFilter() {
  activeFilter.value = null
}

defineExpose({
  clearFilter,
})
</script>

<style scoped>
.search-bar-wrapper {
  background: #f7f8fa;
  position: sticky;
  top: 0;
  z-index: 10;
}

.filter-tags {
  display: flex;
  gap: 8px;
  padding: 12px 16px;
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.filter-tags::-webkit-scrollbar {
  display: none;
}

.filter-tags .van-tag {
  flex-shrink: 0;
  cursor: pointer;
}
</style>
