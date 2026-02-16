<template>
  <van-pull-refresh v-model="refreshingProxy" @refresh="handleRefresh">
    <van-list
      v-model:loading="loadingProxy"
      :finished="finished"
      :finished-text="finishedText"
      :immediate-check="immediateCheck"
      @load="handleLoad"
    >
      <van-empty
        v-if="showEmpty"
        :description="emptyText"
        :image="emptyImage"
      />

      <slot></slot>
    </van-list>
  </van-pull-refresh>
</template>

<script setup>
import { ref, computed, watch } from 'vue'

const props = defineProps({
  /**
   * 加载状态
   */
  loading: {
    type: Boolean,
    default: false
  },
  /**
   * 下拉刷新状态
   */
  refreshing: {
    type: Boolean,
    default: false
  },
  /**
   * 是否已加载完成
   */
  finished: {
    type: Boolean,
    default: false
  },
  /**
   * 数据数组（用于判断是否显示空状态）
   */
  data: {
    type: Array,
    default: () => []
  },
  /**
   * 空状态提示文本
   */
  emptyText: {
    type: String,
    default: '暂无数据'
  },
  /**
   * 空状态图片
   * 可选值: 'default', 'error', 'network', 'search'
   */
  emptyImage: {
    type: String,
    default: 'default'
  },
  /**
   * 加载完成提示文本
   */
  finishedText: {
    type: String,
    default: '没有更多了'
  },
  /**
   * 是否在初始化时立即执行滚动加载
   */
  immediateCheck: {
    type: Boolean,
    default: true
  },
  /**
   * 是否显示空状态（设置为 false 可隐藏空状态）
   */
  showEmptyState: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['update:loading', 'update:refreshing', 'load', 'refresh'])

// 代理状态，用于双向绑定
const loadingProxy = computed({
  get: () => props.loading,
  set: (value) => emit('update:loading', value)
})

const refreshingProxy = computed({
  get: () => props.refreshing,
  set: (value) => emit('update:refreshing', value)
})

// 是否显示空状态
const showEmpty = computed(() => {
  return props.showEmptyState && !props.loading && props.data.length === 0
})

// 处理加载事件
function handleLoad() {
  emit('load')
}

// 处理刷新事件
function handleRefresh() {
  emit('refresh')
}
</script>

<style scoped>
/* 可以添加自定义样式 */
</style>
