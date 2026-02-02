<template>
  <g class="timeline-background">
    <!-- 网格线 -->
    <line
      v-for="(day, index) in gridLines"
      :key="'line-' + index"
      :x1="day.x"
      :y1="0"
      :x2="day.x"
      :y2="height"
      :stroke="day.isWeekend ? '#f0f0f0' : '#e8e8e8'"
      :stroke-width="day.isToday ? 2 : 1"
      :stroke-dasharray="day.isToday ? '4' : ''"
    />

    <!-- 周末背景 -->
    <rect
      v-for="(weekend, index) in weekendRects"
      :key="'weekend-' + index"
      :x="weekend.x"
      :y="0"
      :width="weekend.width"
      :height="height"
      fill="#fafafa"
      opacity="0.5"
    />

    <!-- 今天背景 -->
    <rect
      v-if="todayRect"
      :x="todayRect.x"
      :y="0"
      :width="todayRect.width"
      :height="height"
      fill="#fff3e0"
      opacity="0.3"
    />
  </g>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  width: number
  height: number
  timelineDays: any[]
  dayWidth: number
  viewMode: string
}

const props = defineProps<Props>()

/**
 * 网格线
 */
const gridLines = computed(() => {
  return props.timelineDays.map((day, index) => ({
    x: index * props.dayWidth,
    isWeekend: day.isWeekend,
    isToday: day.isToday
  }))
})

/**
 * 周末矩形
 */
const weekendRects = computed(() => {
  const rects: Array<{ x: number; width: number }> = []
  let start = -1

  props.timelineDays.forEach((day, index) => {
    if (day.isWeekend) {
      if (start === -1) {
        start = index
      }
    } else {
      if (start !== -1) {
        rects.push({
          x: start * props.dayWidth,
          width: (index - start) * props.dayWidth
        })
        start = -1
      }
    }
  })

  // 处理末尾的周末
  if (start !== -1) {
    rects.push({
      x: start * props.dayWidth,
      width: (props.timelineDays.length - start) * props.dayWidth
    })
  }

  return rects
})

/**
 * 今天标记
 */
const todayRect = computed(() => {
  const todayIndex = props.timelineDays.findIndex(day => day.isToday)
  if (todayIndex === -1) return null

  return {
    x: todayIndex * props.dayWidth,
    width: props.dayWidth
  }
})
</script>

<style scoped>
.timeline-background {
  pointer-events: none;
}
</style>
