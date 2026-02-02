<template>
  <div class="gantt-stats" v-if="stats">
    <div class="stat-item">
      <span class="stat-label">总任务</span>
      <span class="stat-value">{{ stats.total }}</span>
    </div>
    <div class="stat-item">
      <span class="stat-label">已完成</span>
      <span class="stat-value success">{{ stats.completed }}</span>
    </div>
    <div class="stat-item">
      <span class="stat-label">进行中</span>
      <span class="stat-value primary">{{ stats.inProgress }}</span>
    </div>
    <div class="stat-item">
      <span class="stat-label">关键路径</span>
      <span class="stat-value critical">{{ stats.critical }}</span>
    </div>
    <div class="stat-item">
      <span class="stat-label">完成率</span>
      <span class="stat-value">{{ stats.progressRate }}%</span>
    </div>
  </div>
</template>

<script setup>
defineProps({
  stats: {
    type: Object,
    default: () => ({
      total: 0,
      completed: 0,
      inProgress: 0,
      notStarted: 0,
      delayed: 0,
      critical: 0,
      progressRate: 0
    })
  }
})
</script>

<style scoped>
.gantt-stats {
  display: flex;
  align-items: center;
  gap: 24px;
  padding: 12px 16px;
  background: #f5f7fa;
  border-bottom: 1px solid #dcdfe6;
  flex-wrap: wrap;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.stat-label {
  font-size: 12px;
  color: #909399;
}

.stat-value {
  font-size: 16px;
  font-weight: bold;
  color: #303133;
}

.stat-value.success {
  color: #67c23a;
}

.stat-value.primary {
  color: #409eff;
}

.stat-value.danger {
  color: #f56c6c;
}

.stat-value.critical {
  color: #f56c6c;
}
</style>
