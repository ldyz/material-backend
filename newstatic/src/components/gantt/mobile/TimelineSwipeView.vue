<template>
  <div
    ref="timelineRef"
    class="timeline-swipe-view"
    :style="{ height: containerHeight + 'px' }"
  >
    <!-- 时间轴头部（可水平滚动） -->
    <div
      class="timeline-swipe-view__header"
      :style="{ width: totalWidth + 'px' }"
    >
      <div
        v-for="(day, index) in visibleDays"
        :key="index"
        class="timeline-header__day"
        :class="{ 'is-weekend': day.isWeekend, 'is-today': day.isToday }"
        :style="{ width: dayWidth + 'px' }"
      >
        <div class="day__weekday">{{ day.weekday }}</div>
        <div class="day__date">{{ day.day }}</div>
      </div>
    </div>

    <!-- 任务条区域（支持触摸手势） -->
    <div
      ref="contentRef"
      class="timeline-swipe-view__content"
      @touchstart="handleTouchStart"
      @touchmove="handleTouchMove"
      @touchend="handleTouchEnd"
    >
      <!-- 背景网格 -->
      <div
        class="timeline-grid"
        :style="{ width: totalWidth + 'px' }"
      >
        <div
          v-for="(day, index) in visibleDays"
          :key="'grid-' + index"
          class="timeline-grid__cell"
          :class="{ 'is-weekend': day.isWeekend }"
          :style="{
            width: dayWidth + 'px',
            left: index * dayWidth + 'px'
          }"
        ></div>
      </div>

      <!-- 任务条 -->
      <div
        v-for="(task, index) in tasks"
        :key="task.id"
        class="timeline-task-row"
        :style="{
          height: rowHeight + 'px',
          top: index * rowHeight + 'px'
        }"
        @click="handleTaskClick(task, $event)"
      >
        <div
          class="task-bar"
          :class="getTaskBarClass(task)"
          :style="getTaskBarStyle(task)"
        >
          <span class="task-bar__label">{{ task.name }}</span>
          <!-- 进度条 -->
          <div
            class="task-bar__progress"
            :style="{ width: task.progress + '%' }"
          ></div>
        </div>

        <!-- 依赖关系线（简化版） -->
        <svg
          v-if="showDependencies && index < tasks.length - 1"
          class="dependency-lines"
          :style="{ width: totalWidth + 'px', height: rowHeight + 'px' }"
        >
          <line
            v-for="dep in getVisibleDependencies(task)"
            :key="dep.id"
            :x1="getDependencyX1(task, dep)"
            :y1="rowHeight / 2"
            :x2="getDependencyX2(task, dep)"
            :y2="rowHeight / 2 + rowHeight"
            :stroke="dep.is_critical ? '#f56c6c' : '#909399'"
            stroke-width="2"
          />
        </svg>
      </div>
    </div>

    <!-- 缩放提示 -->
    <div v-if="isPinching" class="pinch-indicator">
      {{ Math.round(pinchScale * 100) }}%
    </div>

    <!-- 页面指示器 -->
    <div class="page-indicator">
      {{ currentPage }} / {{ totalPages }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useTouchGestures } from '@/composables/useTouchGestures'

/**
 * 移动端时间轴滑动视图
 * 支持触摸滑动和缩放
 */
interface Props {
  tasks: any[]
  timelineDays: any[]
  timelineWeeks?: any[]
  timelineMonths?: any[]
  viewMode: string
  dayWidth: number
  rowHeight: number
  showDependencies?: boolean
  showCriticalPath?: boolean
  showBaseline?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  timelineWeeks: () => [],
  timelineMonths: () => [],
  showDependencies: true,
  showCriticalPath: true,
  showBaseline: false
})

const emit = defineEmits<{
  taskClick: [task: any, event: Event]
  taskDblClick: [task: any]
}>()

const timelineRef = ref<HTMLElement>()
const contentRef = ref<HTMLElement>()
const containerHeight = ref(600)

// 滚动位置
const scrollLeft = ref(0)
const currentPage = ref(1)

// 缩放状态
const isPinching = ref(false)
const pinchScale = ref(1)

// 计算总宽度
const totalWidth = computed(() => {
  return props.timelineDays.length * props.dayWidth
})

// 计算可见天数（简化版，显示所有天）
const visibleDays = computed(() => {
  return props.timelineDays
})

// 计算总页数
const totalPages = computed(() => {
  const viewportWidth = containerHeight.value * 2 // 假设 2:1 宽高比
  const daysPerPage = Math.ceil(viewportWidth / props.dayWidth)
  return Math.ceil(props.timelineDays.length / daysPerPage)
})

/**
 * 使用触摸手势
 */
useTouchGestures({
  elementRef: contentRef,
  handlers: {
    onSwipeLeft: () => {
      // 向左滑动，显示未来日期
      const viewportWidth = timelineRef.value?.clientWidth || 0
      scrollLeft.value = Math.min(
        scrollLeft.value + viewportWidth / 2,
        totalWidth.value - viewportWidth
      )
      updatePage()
    },
    onSwipeRight: () => {
      // 向右滑动，显示过去日期
      const viewportWidth = timelineRef.value?.clientWidth || 0
      scrollLeft.value = Math.max(0, scrollLeft.value - viewportWidth / 2)
      updatePage()
    },
    onPinch: (scale: number) => {
      isPinching.value = true
      pinchScale.value = scale
      // TODO: 实现缩放逻辑
    }
  }
})

/**
 * 触摸事件处理（备用方案）
 */
let touchStartX = 0
let touchStartY = 0

function handleTouchStart(event: TouchEvent) {
  if (event.touches.length === 1) {
    touchStartX = event.touches[0].clientX
    touchStartY = event.touches[0].clientY
  }
}

function handleTouchMove(event: TouchEvent) {
  if (event.touches.length === 1) {
    const deltaX = event.touches[0].clientX - touchStartX
    const deltaY = event.touches[0].clientY - touchStartY

    // 水平滑动
    if (Math.abs(deltaX) > Math.abs(deltaY)) {
      event.preventDefault()
    }
  }
}

function handleTouchEnd(event: TouchEvent) {
  isPinching.value = false
  pinchScale.value = 1
}

function updatePage() {
  const viewportWidth = timelineRef.value?.clientWidth || 0
  const daysPerPage = Math.ceil(viewportWidth / props.dayWidth)
  currentPage.value = Math.floor(scrollLeft.value / (daysPerPage * props.dayWidth)) + 1
}

/**
 * 获取任务条样式类
 */
function getTaskBarClass(task: any): string {
  const classes = ['task-bar']

  if (task.is_milestone) {
    classes.push('is-milestone')
  } else if (task.is_critical) {
    classes.push('is-critical')
  }

  switch (task.status) {
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

  return classes.join(' ')
}

/**
 * 获取任务条位置和宽度
 */
function getTaskBarStyle(task: any) {
  const startDate = new Date(task.start)
  const endDate = new Date(task.end)
  const timelineStart = new Date(props.timelineDays[0]?.date)

  const daysDiff = Math.ceil((startDate - timelineStart) / (1000 * 60 * 60 * 24))
  const duration = Math.ceil((endDate - startDate) / (1000 * 60 * 60 * 24))

  const left = daysDiff * props.dayWidth
  const width = duration * props.dayWidth

  return {
    left: left + 'px',
    width: Math.max(width - 8, 20) + 'px'
  }
}

/**
 * 获取可见的依赖关系
 */
function getVisibleDependencies(task: any) {
  return task.successors || []
}

/**
 * 计算依赖线起点 X 坐标
 */
function getDependencyX1(task: any, dep: any) {
  // 简化实现
  const endDate = new Date(task.end)
  const timelineStart = new Date(props.timelineDays[0]?.date)
  const daysDiff = Math.ceil((endDate - timelineStart) / (1000 * 60 * 60 * 24))
  return daysDiff * props.dayWidth
}

/**
 * 计算依赖线终点 X 坐标
 */
function getDependencyX2(task: any, dep: any) {
  // 简化实现
  const successor = props.tasks.find((t: any) => t.id === dep.successor_id)
  if (!successor) return 0

  const startDate = new Date(successor.start)
  const timelineStart = new Date(props.timelineDays[0]?.date)
  const daysDiff = Math.ceil((startDate - timelineStart) / (1000 * 60 * 60 * 24))
  return daysDiff * props.dayWidth
}

/**
 * 处理任务点击
 */
function handleTaskClick(task: any, event: Event) {
  emit('taskClick', task, event)
}

onMounted(() => {
  if (timelineRef.value) {
    containerHeight.value = timelineRef.value.clientHeight - 60 // 减去头部高度
  }
})
</script>

<style scoped>
.timeline-swipe-view {
  position: relative;
  overflow: hidden;
  background: var(--color-bg-secondary);
  user-select: none;
}

.timeline-swipe-view__header {
  display: flex;
  height: 50px;
  background: var(--color-bg-primary);
  border-bottom: 1px solid var(--color-border-light);
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.timeline-header__day {
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border-right: 1px solid var(--color-border-lighter);
}

.timeline-header__day.is-weekend {
  background: var(--color-bg-secondary);
}

.timeline-header__day.is-today {
  background: var(--color-primary-lighter);
}

.day__weekday {
  font-size: 10px;
  color: var(--color-text-secondary);
}

.day__date {
  font-size: 14px;
  font-weight: var(--font-weight-medium);
  color: var(--color-text-primary);
}

.timeline-swipe-view__content {
  position: relative;
  overflow: auto;
  -webkit-overflow-scrolling: touch;
}

.timeline-grid {
  position: absolute;
  top: 0;
  left: 0;
  height: 100%;
}

.timeline-grid__cell {
  position: absolute;
  top: 0;
  bottom: 0;
  border-right: 1px solid var(--color-border-lighter);
}

.timeline-grid__cell.is-weekend {
  background: var(--color-bg-secondary);
}

.timeline-task-row {
  position: absolute;
  left: 0;
  right: 0;
  border-bottom: 1px solid var(--color-border-lighter);
  display: flex;
  align-items: center;
}

.task-bar {
  position: absolute;
  height: 28px;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  padding: 0 8px;
  font-size: 12px;
  color: #fff;
  overflow: hidden;
  cursor: pointer;
  transition: transform var(--transition-fast);
}

.task-bar:active {
  transform: scale(0.98);
}

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
  background: var(--task-bar-critical-bg);
}

.task-bar.is-milestone {
  width: 28px !important;
  border-radius: 50%;
  background: var(--task-bar-milestone-bg);
}

.task-bar__label {
  position: relative;
  z-index: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.task-bar__progress {
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.3);
}

.dependency-lines {
  position: absolute;
  top: 0;
  left: 0;
  pointer-events: none;
}

.pinch-indicator {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: rgba(0, 0, 0, 0.8);
  color: #fff;
  padding: 12px 24px;
  border-radius: var(--radius-lg);
  font-size: var(--font-size-lg);
  pointer-events: none;
  z-index: var(--z-index-tooltip);
}

.page-indicator {
  position: fixed;
  bottom: 16px;
  left: 50%;
  transform: translateX(-50%);
  background: rgba(0, 0, 0, 0.6);
  color: #fff;
  padding: 6px 12px;
  border-radius: var(--radius-full);
  font-size: var(--font-size-sm);
  pointer-events: none;
}
</style>
