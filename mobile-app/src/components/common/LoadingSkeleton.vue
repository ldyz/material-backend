<template>
  <div class="skeleton-wrapper">
    <!-- 卡片骨架屏 -->
    <template v-if="type === 'card'">
      <div v-for="i in count" :key="`card-${i}`" class="skeleton-card">
        <div class="skeleton-title"></div>
        <div class="skeleton-content">
          <div class="skeleton-line short"></div>
          <div class="skeleton-line"></div>
          <div class="skeleton-line medium"></div>
        </div>
      </div>
    </template>

    <!-- 列表骨架屏 -->
    <template v-else-if="type === 'list'">
      <div v-for="i in count" :key="`list-${i}`" class="skeleton-list-item">
        <div class="skeleton-avatar"></div>
        <div class="skeleton-list-content">
          <div class="skeleton-line long"></div>
          <div class="skeleton-line short"></div>
        </div>
      </div>
    </template>

    <!-- 详情骨架屏 -->
    <div v-else-if="type === 'detail'" class="skeleton-detail">
      <div class="skeleton-header">
        <div class="skeleton-title large"></div>
        <div class="skeleton-line long"></div>
      </div>
      <div class="skeleton-body">
        <div class="skeleton-section">
          <div class="skeleton-line"></div>
          <div class="skeleton-line"></div>
          <div class="skeleton-line"></div>
        </div>
        <div class="skeleton-section">
          <div class="skeleton-line medium"></div>
          <div class="skeleton-line"></div>
          <div class="skeleton-line long"></div>
        </div>
      </div>
    </div>

    <!-- 表格骨架屏 -->
    <div v-else-if="type === 'table'" class="skeleton-table">
      <div class="skeleton-table-header">
        <div v-for="col in columns" :key="`th-${col}`" class="skeleton-th"></div>
      </div>
      <div v-for="row in rows" :key="`row-${row}`" class="skeleton-table-row">
        <div v-for="col in columns" :key="`td-${col}-${row}`" class="skeleton-td"></div>
      </div>
    </div>

    <!-- 默认基础骨架屏 -->
    <template v-else>
      <div v-for="i in count" :key="`basic-${i}`" class="skeleton-basic">
        <div class="skeleton-line"></div>
        <div class="skeleton-line short"></div>
      </div>
    </template>
  </div>
</template>

<script setup>
defineProps({
  type: {
    type: String,
    default: 'basic', // basic, card, list, detail, table
  },
  count: {
    type: Number,
    default: 3,
  },
  columns: {
    type: Number,
    default: 4,
  },
  rows: {
    type: Number,
    default: 5,
  },
})
</script>

<style scoped>
.skeleton-wrapper {
  padding: 16px;
}

/* 基础骨架 */
.skeleton-basic {
  margin-bottom: 12px;
}

/* 卡片骨架 */
.skeleton-card {
  background: white;
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 12px;
}

.skeleton-title {
  height: 16px;
  background: linear-gradient(90deg, #f2f2f2 25%, #e6e6e6 50%, #f2f2f2 75%);
  background-size: 200% 100%;
  animation: loading 1.5s ease-in-out infinite;
  border-radius: 4px;
  margin-bottom: 12px;
}

.skeleton-title.large {
  height: 24px;
  width: 60%;
}

.skeleton-content {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.skeleton-line {
  height: 14px;
  background: linear-gradient(90deg, #f2f2f2 25%, #e6e6e6 50%, #f2f2f2 75%);
  background-size: 200% 100%;
  animation: loading 1.5s ease-in-out infinite;
  border-radius: 4px;
}

.skeleton-line.short {
  width: 40%;
}

.skeleton-line.medium {
  width: 60%;
}

.skeleton-line.long {
  width: 80%;
}

/* 列表骨架 */
.skeleton-list-item {
  display: flex;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #f5f5f5;
}

.skeleton-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(90deg, #f2f2f2 25%, #e6e6e6 50%, #f2f2f2 75%);
  background-size: 200% 100%;
  animation: loading 1.5s ease-in-out infinite;
  margin-right: 12px;
  flex-shrink: 0;
}

.skeleton-list-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

/* 详情骨架 */
.skeleton-detail {
  background: white;
  border-radius: 12px;
  padding: 16px;
}

.skeleton-header {
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f5f5f5;
}

.skeleton-body {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.skeleton-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

/* 表格骨架 */
.skeleton-table {
  background: white;
  border-radius: 12px;
  overflow: hidden;
}

.skeleton-table-header {
  display: flex;
  padding: 12px 16px;
  background: #f7f8fa;
  border-bottom: 1px solid #ebedf0;
}

.skeleton-th {
  flex: 1;
  height: 16px;
  background: linear-gradient(90deg, #e6e6e6 25%, #d9d9d9 50%, #e6e6e6 75%);
  background-size: 200% 100%;
  animation: loading 1.5s ease-in-out infinite;
  border-radius: 4px;
  margin: 0 4px;
}

.skeleton-table-row {
  display: flex;
  padding: 12px 16px;
  border-bottom: 1px solid #f5f5f5;
}

.skeleton-td {
  flex: 1;
  height: 14px;
  background: linear-gradient(90deg, #f2f2f2 25%, #e6e6e6 50%, #f2f2f2 75%);
  background-size: 200% 100%;
  animation: loading 1.5s ease-in-out infinite;
  border-radius: 4px;
  margin: 0 4px;
}

@keyframes loading {
  0% {
    background-position: 100% 50%;
  }
  100% {
    background-position: 0 50%;
  }
}
</style>
