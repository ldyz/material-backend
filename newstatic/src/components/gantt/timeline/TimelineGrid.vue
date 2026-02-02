<template>
  <div class="timeline-grid">
    <div
      v-for="(day, index) in days"
      :key="index"
      class="timeline-grid__cell"
      :class="{ 'is-weekend': day.isWeekend, 'is-today': day.isToday }"
      :style="cellStyle(day, index)"
    ></div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  days: any[]
  dayWidth: number
}

const props = defineProps<Props>()

/**
 * 单元格样式
 */
function cellStyle(day: any, index: number) {
  return {
    position: 'absolute',
    left: index * props.dayWidth + 'px',
    width: props.dayWidth + 'px',
    top: 0,
    bottom: 0
  }
}
</script>

<style scoped>
.timeline-grid {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
}

.timeline-grid__cell {
  border-right: 1px solid var(--color-border-lighter);
}

.timeline-grid__cell.is-weekend {
  background: var(--color-bg-secondary);
}

.timeline-grid__cell.is-today {
  background: var(--color-primary-lighter);
}
</style>
