<template>
  <div
    class="task-bar-wrapper"
    :class="wrapperClass"
    :style="wrapperStyle"
    @click="handleClick"
    @dblclick="handleDblClick"
    @mousedown="handleMouseDown"
  >
    <!-- 任务条 -->
    <div
      class="task-bar"
      :class="taskBarClass"
      :style="taskBarStyle"
    >
      <!-- 里程碑（菱形） -->
      <template v-if="task.isMilestone">
        <div class="milestone-shape"></div>
      </template>

      <!-- 普通任务条 -->
      <template v-else>
        <!-- 进度条 -->
        <div
          v-if="task.progress > 0"
          class="task-bar__progress"
          :style="{ width: task.progress + '%' }"
        ></div>

        <!-- 任务文本 -->
        <span class="task-bar__label">{{ task.progress }}%</span>
      </template>

      <!-- 拖拽手柄 -->
      <div
        v-if="!task.isMilestone"
        class="resize-handle resize-handle--left"
      ></div>
      <div
        v-if="!task.isMilestone"
        class="resize-handle resize-handle--right"
      ></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { GanttTask } from '@/types/gantt'

interface Props {
  task: any
  rowIndex: number
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
 * 包裹容器类名
 */
const wrapperClass = computed(() => {
  return {
    'is-selected': props.selected,
    'is-critical': props.isCritical,
    'is-dragging': props.isDragging
  }
})

/**
 * 包裹容器样式
 */
const wrapperStyle = computed(() => {
  return {
    position: 'absolute',
    left: props.task.x + 'px',
    top: props.task.y + 'px',
    width: props.task.width + 'px',
    height: props.taskHeight + 'px'
  }
})

/**
 * 任务条类名
 */
const taskBarClass = computed(() => {
  const classes = ['task-bar']

  if (props.task.isMilestone) {
    classes.push('is-milestone')
  }

  switch (props.task.status) {
    case 'completed':
      classes.push('is-completed')
      break
    case 'in_progress':
      classes.push('is-in-progress')
      break
    case 'delayed':
      classes.push('is-delayed')
      break
    default:
      classes.push('is-not-started')
  }

  return classes
})

/**
 * 任务条样式
 */
const taskBarStyle = computed(() => {
  const style: Record<string, string> = {}

  if (props.isDragging) {
    style.opacity = '0.8'
  }

  return style
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
.task-bar-wrapper {
  display: flex;
  align-items: center;
  cursor: pointer;
  transition: transform 0.15s ease;
}

.task-bar-wrapper:active {
  transform: scale(0.98);
}

.task-bar-wrapper.is-selected {
  z-index: 10;
}

.task-bar {
  position: relative;
  height: 100%;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: all var(--transition-fast);
}

/* 状态颜色 */
.task-bar.is-completed {
  background: var(--task-bar-completed-bg);
}

.task-bar.is-in-progress {
  background: var(--task-bar-in-progress-bg);
}

.task-bar.is-not-started {
  background: var(--task-bar-not-started-bg);
}

.task-bar.is-delayed {
  background: var(--task-bar-delayed-bg);
}

.task-bar.is-critical {
  border: 2px solid var(--task-bar-critical);
}

/* 里程碑 */
.task-bar.is-milestone {
  background: var(--task-bar-milestone-bg);
}

.milestone-shape {
  width: 24px;
  height: 24px;
  background: var(--task-bar-milestone);
  transform: rotate(45deg);
  border-radius: 2px;
}

/* 进度条 */
.task-bar__progress {
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.3);
  transition: width var(--transition-base);
}

/* 任务文本 */
.task-bar__label {
  position: relative;
  z-index: 1;
  font-size: 12px;
  font-weight: bold;
  color: #fff;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);
  pointer-events: none;
}

/* 拖拽手柄 */
.resize-handle {
  position: absolute;
  top: 0;
  bottom: 0;
  width: 8px;
  cursor: col-resize;
  transition: background-color var(--transition-fast);
}

.resize-handle--left {
  left: 0;
}

.resize-handle--right {
  right: 0;
}

.task-bar:hover .resize-handle {
  background: rgba(0, 0, 0, 0.1);
}

/* 选中状态 */
.task-bar-wrapper.is-selected .task-bar {
  box-shadow: 0 0 0 2px var(--color-primary);
}

/* 拖拽状态 */
.task-bar-wrapper.is-dragging .task-bar {
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
  transform: scale(1.02);
}
</style>
