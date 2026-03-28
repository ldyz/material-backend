<template>
  <div class="gantt-legend">
    <div class="legend-item">
      <span class="legend-color completed"></span>
      <span>已完成</span>
    </div>
    <div class="legend-item">
      <span class="legend-color in-progress"></span>
      <span>进行中</span>
    </div>
    <div class="legend-item">
      <span class="legend-color not-started"></span>
      <span>未开始</span>
    </div>
    <div class="legend-item">
      <span class="legend-color delayed"></span>
      <span>已延期</span>
    </div>
    <div class="legend-item" v-if="showCriticalPath">
      <span class="legend-color critical"></span>
      <span>关键路径</span>
    </div>
    <div class="legend-item">
      <span class="legend-color milestone"></span>
      <span>里程碑</span>
    </div>
    <div class="legend-item" v-if="showBaseline">
      <span class="legend-color baseline"></span>
      <span>基线</span>
    </div>
    <div class="legend-item">
      <span class="legend-color today"></span>
      <span>今天</span>
    </div>
  </div>
</template>

<script setup>
defineProps({
  showCriticalPath: {
    type: Boolean,
    default: true
  },
  showBaseline: {
    type: Boolean,
    default: false
  }
})
</script>

<style scoped>
.gantt-legend {
  display: flex;
  align-items: center;
  gap: 24px;
  padding: 12px 16px;
  background: #f5f7fa;
  border-top: 1px solid #dcdfe6;
  flex-wrap: wrap;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: #606266;
}

.legend-color {
  width: 20px;
  height: 12px;
  border-radius: 2px;
}

.legend-color.completed {
  background: #67c23a;
}

.legend-color.in-progress {
  background: #409eff;
}

.legend-color.not-started {
  background: #909399;
}

.legend-color.delayed {
  background: #f56c6c;
}

.legend-color.critical {
  background: linear-gradient(135deg, #f56c6c 0%, #c83232 100%);
}

.legend-color.milestone {
  background: #e6a23c;
  transform: rotate(45deg) scale(0.6);
}

.legend-color.baseline {
  background: #c0c4cc;
  opacity: 0.5;
}

.legend-color.today {
  background: #e6a23c;
  border: 2px dashed #e6a23c;
  background: transparent;
}
</style>
