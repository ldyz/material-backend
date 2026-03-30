<template>
  <div class="list-skeleton">
    <van-skeleton
      v-for="i in count"
      :key="i"
      :row="rows"
      :title="showTitle"
      :avatar="showAvatar"
      class="skeleton-item"
    />
  </div>
</template>

<script setup>
/**
 * 列表骨架屏组件
 * 用于在数据加载时显示占位骨架
 */

defineProps({
  // 骨架项数量
  count: {
    type: Number,
    default: 5
  },
  // 每项的行数
  rows: {
    type: Number,
    default: 3
  },
  // 是否显示标题
  showTitle: {
    type: Boolean,
    default: false
  },
  // 是否显示头像
  showAvatar: {
    type: Boolean,
    default: false
  }
})
</script>

<style scoped>
.list-skeleton {
  padding: 12px;
}

.skeleton-item {
  margin-bottom: 12px;
  padding: 12px;
  background: #fff;
  border-radius: 8px;
}

.skeleton-item:last-child {
  margin-bottom: 0;
}
</style>
