<template>
  <g
    :transform="`translate(${task.x}, ${task.y})`"
    class="task-bar-group"
    :class="groupClass"
    @click.stop="handleClick"
    @dblclick.stop="handleDblClick"
    @mousedown.stop="handleMouseDown"
  >
    <!-- 任务条矩形 -->
    <rect
      v-if="!task.isMilestone"
      :x="0"
      :y="0"
      :width="task.width"
      :height="taskHeight"
      :fill="barColor"
      :stroke="isCritical ? '#f56c6c' : 'none'"
      :stroke-width="isCritical ? 2 : 0"
      :rx="4"
      class="task-bar-rect"
    />

    <!-- 里程碑（菱形） -->
    <g v-else class="milestone-group">
      <rect
        :x="task.width / 2 - taskHeight / 2"
        :y="-taskHeight / 2"
        :width="taskHeight"
        :height="taskHeight"
        :fill="barColor"
        :transform="`rotate(45, ${task.width / 2}, 0)`"
      />
    </g>

    <!-- 进度条 -->
    <rect
      v-if="!task.isMilestone && task.progress > 0"
      :x="0"
      :y="0"
      :width="task.width * (task.progress / 100)"
      :height="taskHeight"
      :fill="progressColor"
      :rx="4"
      opacity="0.3"
    />

    <!-- 任务文本 -->
    <text
      :x="task.width / 2"
      :y="taskHeight / 2 + 4"
      text-anchor="middle"
      fill="white"
      font-size="12"
      font-weight="bold"
    >
      {{ task.isMilestone ? '' : task.progress + '%' }}
    </text>

    <!-- 拖拽手柄 -->
    <rect
      v-if="!task.isMilestone"
      :x="0"
      :y="0"
      :width="8"
      :height="taskHeight"
      fill="transparent"
      class="resize-handle-left"
    />
    <rect
      v-if="!task.isMilestone"
      :x="task.width - 8"
      :y="0"
      :width="8"
      :height="taskHeight"
      fill="transparent"
      class="resize-handle-right"
    />
  </g>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  task: any
  taskHeight: number
  selected: boolean
  isCritical: boolean
  isDragging: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
  click: [event: Event]
  dblClick: [event: Event]
  mousedown: [event: MouseEvent]
}>()

/**
 * 组类名
 */
const groupClass = computed(() => {
  return {
    'is-selected': props.selected,
    'is-critical': props.isCritical,
    'is-dragging': props.isDragging
  }
})

/**
 * 任务条颜色
 */
const barColor = computed(() => {
  if (props.task.isMilestone) {
    return '#f39c12'
  }

  switch (props.task.status) {
    case 'completed':
      return '#67c23a'
    case 'in_progress':
      return '#409eff'
    case 'delayed':
      return '#f56c6c'
    default:
      return '#909399'
  }
})

/**
 * 进度条颜色
 */
const progressColor = computed(() => {
  return barColor.value
})

/**
 * 处理点击
 */
function handleClick(event: Event) {
  emit('click', event)
}

/**
 * 处理双击
 */
function handleDblClick(event: Event) {
  emit('dblClick', event)
}

/**
 * 处理鼠标按下
 */
function handleMouseDown(event: MouseEvent) {
  emit('mousedown', event)
}
</script>

<style scoped>
.task-bar-group {
  cursor: pointer;
  transition: opacity 0.15s ease;
}

.task-bar-group:active {
  opacity: 0.8;
}

.task-bar-group.is-selected .task-bar-rect {
  filter: drop-shadow(0 0 4px rgba(64, 158, 255, 0.5));
}

.task-bar-group.is-dragging {
  cursor: grabbing;
}

.resize-handle-left {
  cursor: col-resize;
}

.resize-handle-right {
  cursor: col-resize;
}
</style>
