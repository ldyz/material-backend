<template>
  <div :class="['network-node', 'node-' + data.type]">
    <div class="node-number" v-if="data.number">{{ data.number }}</div>
    <div class="node-content" v-if="showLabel">
      <div class="node-label">{{ data.label }}</div>
      <div class="node-times" v-if="data.earliestTime !== undefined">
        ES: {{ data.earliestTime }} / LS: {{ data.latestTime }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  data: {
    type: Object,
    required: true
  },
  showLabel: {
    type: Boolean,
    default: true
  }
})
</script>

<style scoped>
.network-node {
  padding: 12px 16px;
  border-radius: 8px;
  border: 2px solid #409eff;
  background: white;
  min-width: 120px;
  text-align: center;
  position: relative;
  transition: all 0.3s;
}

.network-node:hover {
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
  transform: translateY(-2px);
}

.node-start {
  border-color: #67c23a;
  background: linear-gradient(135deg, #f0f9ff 0%, #e8f5e9 100%);
}

.node-event {
  border-color: #409eff;
  background: linear-gradient(135deg, #f0f9ff 0%, #e3f2fd 100%);
}

.node-end {
  border-color: #f56c6c;
  background: linear-gradient(135deg, #f0f9ff 0%, #ffebee 100%);
}

.node-number {
  position: absolute;
  top: -8px;
  left: -8px;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: #409eff;
  color: white;
  font-size: 12px;
  font-weight: bold;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px solid white;
}

.node-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.node-label {
  font-weight: 500;
  color: #303133;
  font-size: 14px;
}

.node-times {
  font-size: 11px;
  color: #909399;
  margin-top: 4px;
}
</style>
