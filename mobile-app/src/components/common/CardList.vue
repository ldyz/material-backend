<template>
  <div class="card-list-wrapper">
    <!-- 下拉刷新 + 列表 -->
    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <van-list
        v-model:loading="loading"
        :finished="finished"
        :finished-text="finishedText"
        :error="error"
        :error-text="errorText"
        @load="onLoad"
      >
        <!-- 列表项插槽 -->
        <template v-for="(item, index) in list" :key="item[idField] || index">
          <slot name="item" :item="item" :index="index">
            <!-- 默认卡片 -->
            <div class="default-card" @click="handleClick(item)">
              <slot name="card" :item="item">
                <pre>{{ JSON.stringify(item, null, 2) }}</pre>
              </slot>
            </div>
          </template>
        </template>
      </van-list>
    </van-pull-refresh>

    <!-- 空状态 -->
    <van-empty
      v-if="!loading && list.length === 0 && !refreshing"
      :description="emptyText"
      :image="emptyImage"
    >
      <template v-if="showEmptyAction" #default>
        <van-button type="primary" size="small" @click="handleEmptyAction">
          {{ emptyActionText }}
        </van-button>
      </template>
    </van-empty>
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { useAPIList } from '@/composables/useAPI'

const props = defineProps({
  // API函数
  apiFunction: {
    type: Function,
    required: true,
  },
  // 额外的API参数
  apiParams: {
    type: Object,
    default: () => ({}),
  },
  // ID字段名
  idField: {
    type: String,
    default: 'id',
  },
  // 每页数量
  pageSize: {
    type: Number,
    default: 20,
  },
  // 完成文本
  finishedText: {
    type: String,
    default: '没有更多了',
  },
  // 错误文本
  errorText: {
    type: String,
    default: '请求失败，点击重试',
  },
  // 空状态文本
  emptyText: {
    type: String,
    default: '暂无数据',
  },
  // 空状态图片
  emptyImage: {
    type: String,
    default: 'default',
  },
  // 是否显示空状态操作按钮
  showEmptyAction: {
    type: Boolean,
    default: false,
  },
  // 空状态操作按钮文本
  emptyActionText: {
    type: String,
    default: '去添加',
  },
  // 是否立即加载
  immediate: {
    type: Boolean,
    default: true,
  },
})

const emit = defineEmits(['item-click', 'empty-action', 'refresh', 'load'])

// 使用useAPIList
const {
  list,
  loading,
  refreshing,
  finished,
  error,
  load,
  refresh,
  loadMore,
  reset,
} = useAPIList(props.apiFunction, {
  immediate: props.immediate,
  pageSize: props.pageSize,
})

// 加载数据
async function onLoad() {
  await loadMore(props.apiParams)
}

// 刷新数据
async function onRefresh() {
  await refresh(props.apiParams)
  emit('refresh')
}

// 项目点击
function handleClick(item) {
  emit('item-click', item)
}

// 空状态操作
function handleEmptyAction() {
  emit('empty-action')
}

// 监听API参数变化，重新加载
watch(() => props.apiParams, () => {
  reset()
  onLoad()
}, { deep: true })

// 暴露方法
defineExpose({
  refresh: onRefresh,
  reset,
  list,
})
</script>

<style scoped>
.card-list-wrapper {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.default-card {
  background: white;
  margin: 12px 16px;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}
</style>
