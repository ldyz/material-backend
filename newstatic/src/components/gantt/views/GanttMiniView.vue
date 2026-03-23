<template>
  <div class="gantt-mini-view" :style="{ height: height + 'px' }">
    <svg :width="width" :height="height">
      <!-- Background -->
      <rect
        :x="0"
        :y="0"
        :width="width"
        :height="height"
        fill="#f5f7fa"
      />

      <!-- Grid lines -->
      <g class="grid">
        <line
          v-for="(line, index) in gridLines"
          :key="'grid-' + index"
          :x1="line.x"
          :y1="0"
          :x2="line.x"
          :y2="height"
          :stroke="index % 2 === 0 ? '#e4e7ed' : '#ebeef5'"
          stroke-width="1"
        />
      </g>

      <!-- Task bars -->
      <g class="tasks">
        <rect
          v-for="task in visibleTasks"
          :key="task.id"
          :x="task.x"
          :y="task.y"
          :width="task.width"
          :height="task.height"
          :fill="task.fill"
          :stroke="task.stroke"
          stroke-width="1"
          rx="2"
        />
      </g>

      <!-- Conflict indicators -->
      <g v-if="showConflicts" class="conflicts">
        <rect
          v-for="conflict in conflictZones"
          :key="'conflict-' + conflict.index"
          :x="conflict.x"
          :y="0"
          :width="conflict.width"
          :height="height"
          fill="rgba(245, 108, 108, 0.1)"
        />
      </g>
    </svg>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  tasks: {
    type: Array,
    default: () => []
  },
  conflicts: {
    type: Array,
    default: () => []
  },
  height: {
    type: Number,
    default: 150
  }
})

const width = computed(() => 600) // Fixed width for mini view

const dateRange = computed(() => {
  if (!props.tasks || props.tasks.length === 0) {
    return { start: new Date(), end: new Date() }
  }

  const startDates = props.tasks.map(t => new Date(t.start_date).getTime())
  const endDates = props.tasks.map(t => new Date(t.end_date).getTime())

  return {
    start: new Date(Math.min(...startDates)),
    end: new Date(Math.max(...endDates))
  }
})

const totalDays = computed(() => {
  return Math.ceil((dateRange.value.end - dateRange.value.start) / (1000 * 60 * 60 * 24))
})

const dayWidth = computed(() => {
  return width.value / totalDays.value
})

const gridLines = computed(() => {
  const lines = []
  for (let i = 0; i <= totalDays.value; i++) {
    lines.push({ x: i * dayWidth.value })
  }
  return lines
})

const visibleTasks = computed(() => {
  if (!props.tasks || props.tasks.length === 0) {
    return []
  }

  const rowHeight = 20
  const barHeight = 14

  return props.tasks.slice(0, 10).map((task, index) => {
    const startDate = new Date(task.start_date)
    const endDate = new Date(task.end_date)

    const x = ((startDate - dateRange.value.start) / (1000 * 60 * 60 * 24)) * dayWidth.value
    const taskWidth = ((endDate - startDate) / (1000 * 60 * 60 * 24)) * dayWidth.value

    let fill = '#409EFF'
    let stroke = 'transparent'

    if (task.leveled) {
      fill = '#E6A23C'
    }

    if (task.dependencyAdjusted) {
      stroke = '#F56C6C'
    }

    return {
      id: task.id,
      x: Math.max(0, x),
      y: index * rowHeight + (rowHeight - barHeight) / 2,
      width: Math.max(taskWidth, 2),
      height: barHeight,
      fill,
      stroke
    }
  })
})

const showConflicts = computed(() => {
  return props.conflicts && props.conflicts.length > 0
})

const conflictZones = computed(() => {
  if (!showConflicts.value) {
    return []
  }

  return props.conflicts.map((conflict, index) => {
    const conflictDate = new Date(conflict.date)
    const x = ((conflictDate - dateRange.value.start) / (1000 * 60 * 60 * 24)) * dayWidth.value

    return {
      index,
      x: Math.max(0, x),
      width: dayWidth.value
    }
  })
})
</script>

<style scoped>
.gantt-mini-view {
  width: 100%;
  overflow: hidden;
}
</style>
