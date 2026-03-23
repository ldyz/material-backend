<template>
  <div class="gantt-minimap" :class="{ minimized: isMinimized }">
    <!-- Header -->
    <div class="minimap-header">
      <span class="header-title">{{ $t('gantt.minimap.title') }}</span>
      <div class="header-actions">
        <el-button
          :icon="isMinimized ? ArrowDown : ArrowUp"
          circle
          text
          size="small"
          @click="toggleMinimize"
        />
        <el-button
          :icon="Close"
          circle
          text
          size="small"
          @click="$emit('close')"
        />
      </div>
    </div>

    <!-- Minimap Content -->
    <transition name="el-collapse">
      <div v-show="!isMinimized" class="minimap-content">
        <!-- Canvas for rendering minimap -->
        <canvas
          ref="canvasRef"
          class="minimap-canvas"
          @mousedown="handleMouseDown"
          @mousemove="handleMouseMove"
          @mouseup="handleMouseUp"
          @mouseleave="handleMouseLeave"
          @wheel.prevent="handleWheel"
        ></canvas>

        <!-- Viewport Indicator -->
        <div
          v-if="viewport"
          class="viewport-indicator"
          :style="viewportStyle"
        >
          <!-- Resize handles -->
          <div class="resize-handle handle-left" @mousedown.stop="startResize('left', $event)"></div>
          <div class="resize-handle handle-right" @mousedown.stop="startResize('right', $event)"></div>
          <div class="resize-handle handle-top" @mousedown.stop="startResize('top', $event)"></div>
          <div class="resize-handle handle-bottom" @mousedown.stop="startResize('bottom', $event)"></div>
        </div>

        <!-- Critical Path Overlay -->
        <div v-if="showCriticalPath" class="critical-path-overlay">
          <div
            v-for="(segment, index) in criticalPathSegments"
            :key="index"
            class="critical-segment"
            :style="segment.style"
          ></div>
        </div>

        <!-- Selection Highlight -->
        <div
          v-if="selectedTaskPosition"
          class="selection-highlight"
          :style="selectedTaskStyle"
        ></div>

        <!-- Current Time Indicator -->
        <div class="current-time-indicator" :style="currentTimeStyle"></div>

        <!-- Legend -->
        <div class="minimap-legend">
          <div class="legend-item">
            <div class="legend-color normal"></div>
            <span>{{ $t('gantt.minimap.legend.normal') }}</span>
          </div>
          <div class="legend-item">
            <div class="legend-color milestone"></div>
            <span>{{ $t('gantt.minimap.legend.milestone') }}</span>
          </div>
          <div class="legend-item">
            <div class="legend-color critical"></div>
            <span>{{ $t('gantt.minimap.legend.critical') }}</span>
          </div>
          <div class="legend-item">
            <div class="legend-color delayed"></div>
            <span>{{ $t('gantt.minimap.legend.delayed') }}</span>
          </div>
        </div>
      </div>
    </transition>

    <!-- Zoom Controls -->
    <div v-show="!isMinimized" class="minimap-zoom">
      <el-button-group size="small">
        <el-button :icon="ZoomOut" @click="zoomOut" />
        <el-button :icon="ZoomIn" @click="zoomIn" />
        <el-button @click="resetView">{{ $t('gantt.minimap.reset') }}</el-button>
      </el-button-group>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { ArrowUp, ArrowDown, Close, ZoomIn, ZoomOut } from '@element-plus/icons-vue'

const props = defineProps({
  /**
   * All tasks to display
   */
  tasks: {
    type: Array,
    default: () => []
  },
  /**
   * Critical path task IDs
   */
  criticalPath: {
    type: Array,
    default: () => []
  },
  /**
   * Selected task ID
   */
  selectedTaskId: {
    type: [String, Number],
    default: null
  },
  /**
   * Viewport information
   */
  viewport: {
    type: Object,
    default: () => ({ x: 0, y: 0, width: 0, height: 0 })
  },
  /**
   * Timeline range
   */
  timelineRange: {
    type: Object,
    default: () => ({ start: 0, end: 100 })
  },
  /**
   * Show critical path overlay
   */
  showCriticalPath: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['viewport-change', 'navigate-to-task', 'close'])

const { t } = useI18n()

const canvasRef = ref(null)
const isMinimized = ref(false)
const isDragging = ref(false)
const isResizing = ref(false)
const resizeDirection = ref(null)
const dragStart = ref({ x: 0, y: 0 })
const viewportStart = ref({ x: 0, y: 0, width: 0, height: 0 })

const ctx = ref(null)
const animationFrameId = ref(null)

/**
 * Toggle minimize
 */
function toggleMinimize() {
  isMinimized.value = !isMinimized.value
  if (!isMinimized.value) {
    nextTick(() => {
      renderMinimap()
    })
  }
}

/**
 * Get viewport style
 */
const viewportStyle = computed(() => {
  if (!props.viewport) return {}

  const canvasWidth = canvasRef.value?.width || 0
  const canvasHeight = canvasRef.value?.height || 0
  const totalDuration = props.timelineRange.end - props.timelineRange.start

  return {
    left: `${((props.viewport.x - props.timelineRange.start) / totalDuration) * canvasWidth}px`,
    top: `${(props.viewport.y / props.tasks.length) * canvasHeight}px`,
    width: `${(props.viewport.width / totalDuration) * canvasWidth}px`,
    height: `${(props.viewport.height / props.tasks.length) * canvasHeight}px`
  }
})

/**
 * Get selected task style
 */
const selectedTaskStyle = computed(() => {
  if (!props.selectedTaskId) return {}

  const task = props.tasks.find(t => t.id === props.selectedTaskId)
  if (!task) return {}

  const canvasWidth = canvasRef.value?.width || 0
  const canvasHeight = canvasRef.value?.height || 0
  const totalDuration = props.timelineRange.end - props.timelineRange.start
  const taskIndex = props.tasks.findIndex(t => t.id === props.selectedTaskId)

  const rowHeight = canvasHeight / props.tasks.length
  const dayWidth = canvasWidth / totalDuration

  return {
    left: `${(task.start / totalDuration) * canvasWidth}px`,
    top: `${taskIndex * rowHeight}px`,
    width: `${task.duration * dayWidth}px`,
    height: `${rowHeight - 2}px`
  }
})

/**
 * Get current time style
 */
const currentTimeStyle = computed(() => {
  const canvasWidth = canvasRef.value?.width || 0
  const totalDuration = props.timelineRange.end - props.timelineRange.start
  const currentDay = Math.floor(Date.now() / (1000 * 60 * 60 * 24)) - props.timelineRange.start

  return {
    left: `${(currentDay / totalDuration) * canvasWidth}px`
  }
})

/**
 * Get critical path segments
 */
const criticalPathSegments = computed(() => {
  if (!props.criticalPath || props.criticalPath.length === 0) return []

  const canvasWidth = canvasRef.value?.width || 0
  const canvasHeight = canvasRef.value?.height || 0
  const totalDuration = props.timelineRange.end - props.timelineRange.start
  const rowHeight = canvasHeight / props.tasks.length
  const dayWidth = canvasWidth / totalDuration

  return props.criticalPath.map((taskId, index) => {
    const task = props.tasks.find(t => t.id === taskId)
    if (!task) return null

    const taskIndex = props.tasks.findIndex(t => t.id === taskId)

    return {
      style: {
        left: `${(task.start / totalDuration) * canvasWidth}px`,
        top: `${taskIndex * rowHeight + rowHeight / 2}px`,
        width: `${task.duration * dayWidth}px`,
        height: '2px'
      }
    }
  }).filter(Boolean)
})

/**
 * Render minimap
 */
function renderMinimap() {
  if (!canvasRef.value || !ctx.value) return

  const canvas = canvasRef.value
  const context = ctx.value
  const width = canvas.width
  const height = canvas.height

  // Clear canvas
  context.clearRect(0, 0, width, height)

  // Draw background
  context.fillStyle = '#f5f7fa'
  context.fillRect(0, 0, width, height)

  if (props.tasks.length === 0) return

  const rowHeight = height / props.tasks.length
  const totalDuration = props.timelineRange.end - props.timelineRange.start
  const dayWidth = width / totalDuration

  // Draw tasks
  props.tasks.forEach((task, index) => {
    const y = index * rowHeight
    const x = (task.start / totalDuration) * width
    const taskWidth = task.duration * dayWidth

    // Determine task color
    let fillColor = '#409EFF' // Default blue
    let strokeColor = '#337ECC'

    if (task.isMilestone) {
      fillColor = '#E6A23C' // Orange for milestones
      strokeColor = '#CF8E2D'
    } else if (props.criticalPath.includes(task.id)) {
      fillColor = '#F56C6C' // Red for critical path
      strokeColor = '#D64542'
    } else if (task.delayed) {
      fillColor = '#909399' // Gray for delayed
      strokeColor = '#73767A'
    }

    // Draw task bar
    context.fillStyle = fillColor
    context.fillRect(x + 1, y + 1, taskWidth - 2, rowHeight - 2)

    // Draw task border
    context.strokeStyle = strokeColor
    context.lineWidth = 1
    context.strokeRect(x + 1, y + 1, taskWidth - 2, rowHeight - 2)

    // Draw progress indicator
    if (task.progress > 0 && !task.isMilestone) {
      context.fillStyle = 'rgba(255, 255, 255, 0.3)'
      context.fillRect(x + 1, y + 1, taskWidth * (task.progress / 100) - 2, rowHeight - 2)
    }
  })

  // Draw grid lines
  context.strokeStyle = 'rgba(0, 0, 0, 0.05)'
  context.lineWidth = 1

  // Vertical lines (time)
  const gridInterval = Math.ceil(totalDuration / 10)
  for (let day = props.timelineRange.start; day <= props.timelineRange.end; day += gridInterval) {
    const x = ((day - props.timelineRange.start) / totalDuration) * width
    context.beginPath()
    context.moveTo(x, 0)
    context.lineTo(x, height)
    context.stroke()
  }

  // Horizontal lines (tasks)
  for (let i = 0; i <= props.tasks.length; i++) {
    const y = i * rowHeight
    context.beginPath()
    context.moveTo(0, y)
    context.lineTo(width, y)
    context.stroke()
  }
}

/**
 * Handle mouse down
 */
function handleMouseDown(event) {
  isDragging.value = true
  dragStart.value = { x: event.clientX, y: event.clientY }
  viewportStart.value = { ...props.viewport }
}

/**
 * Handle mouse move
 */
function handleMouseMove(event) {
  if (!isDragging.value || isResizing.value) return

  const dx = event.clientX - dragStart.value.x
  const dy = event.clientY - dragStart.value.y

  const canvasWidth = canvasRef.value?.width || 0
  const canvasHeight = canvasRef.value?.height || 0
  const totalDuration = props.timelineRange.end - props.timelineRange.start

  // Convert pixel delta to time/task delta
  const timeDelta = (dx / canvasWidth) * totalDuration
  const taskDelta = (dy / canvasHeight) * props.tasks.length

  emit('viewport-change', {
    x: viewportStart.value.x + timeDelta,
    y: viewportStart.value.y + taskDelta,
    width: viewportStart.value.width,
    height: viewportStart.value.height
  })
}

/**
 * Handle mouse up
 */
function handleMouseUp() {
  isDragging.value = false
  isResizing.value = false
  resizeDirection.value = null
}

/**
 * Handle mouse leave
 */
function handleMouseLeave() {
  isDragging.value = false
  isResizing.value = false
  resizeDirection.value = null
}

/**
 * Start resize
 */
function startResize(direction, event) {
  isResizing.value = true
  resizeDirection.value = direction
  dragStart.value = { x: event.clientX, y: event.clientY }
  viewportStart.value = { ...props.viewport }
}

/**
 * Handle wheel
 */
function handleWheel(event) {
  const zoomFactor = event.deltaY > 0 ? 1.1 : 0.9
  const newWidth = props.viewport.width * zoomFactor
  const newHeight = props.viewport.height * zoomFactor

  emit('viewport-change', {
    x: props.viewport.x,
    y: props.viewport.y,
    width: Math.max(10, Math.min(newWidth, props.timelineRange.end - props.timelineRange.start)),
    height: Math.max(5, Math.min(newHeight, props.tasks.length))
  })
}

/**
 * Zoom in
 */
function zoomIn() {
  emit('viewport-change', {
    x: props.viewport.x,
    y: props.viewport.y,
    width: props.viewport.width * 0.8,
    height: props.viewport.height * 0.8
  })
}

/**
 * Zoom out
 */
function zoomOut() {
  emit('viewport-change', {
    x: props.viewport.x,
    y: props.viewport.y,
    width: props.viewport.width * 1.2,
    height: props.viewport.height * 1.2
  })
}

/**
 * Reset view
 */
function resetView() {
  emit('viewport-change', {
    x: props.timelineRange.start,
    y: 0,
    width: props.timelineRange.end - props.timelineRange.start,
    height: props.tasks.length
  })
}

/**
 * Resize canvas
 */
function resizeCanvas() {
  if (!canvasRef.value) return

  const container = canvasRef.value.parentElement
  const width = container.clientWidth - 32
  const height = Math.min(300, props.tasks.length * 3)

  canvasRef.value.width = width
  canvasRef.value.height = height

  renderMinimap()
}

// Watch for data changes
watch(() => [props.tasks, props.criticalPath, props.timelineRange], () => {
  if (animationFrameId.value) {
    cancelAnimationFrame(animationFrameId.value)
  }
  animationFrameId.value = requestAnimationFrame(renderMinimap)
}, { deep: true })

// Lifecycle
onMounted(() => {
  ctx.value = canvasRef.value?.getContext('2d')
  resizeCanvas()
  window.addEventListener('resize', resizeCanvas)
})

onBeforeUnmount(() => {
  if (animationFrameId.value) {
    cancelAnimationFrame(animationFrameId.value)
  }
  window.removeEventListener('resize', resizeCanvas)
})
</script>

<script>
export default {
  name: 'GanttMinimap'
}
</script>

<style scoped>
.gantt-minimap {
  background: white;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  overflow: hidden;
}

/* Header */
.minimap-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #f5f7fa;
  border-bottom: 1px solid #e4e7ed;
}

.header-title {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
}

.header-actions {
  display: flex;
  gap: 4px;
}

/* Content */
.minimap-content {
  position: relative;
  padding: 16px;
}

.minimap-canvas {
  width: 100%;
  height: 200px;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  cursor: grab;
}

.minimap-canvas:active {
  cursor: grabbing;
}

/* Viewport Indicator */
.viewport-indicator {
  position: absolute;
  border: 2px solid #409EFF;
  background: rgba(64, 158, 255, 0.1);
  pointer-events: none;
  transition: all 0.1s ease;
}

.resize-handle {
  position: absolute;
  background: #409EFF;
  pointer-events: auto;
}

.handle-left,
.handle-right {
  top: 0;
  bottom: 0;
  width: 4px;
  cursor: ew-resize;
}

.handle-left {
  left: -2px;
}

.handle-right {
  right: -2px;
}

.handle-top,
.handle-bottom {
  left: 0;
  right: 0;
  height: 4px;
  cursor: ns-resize;
}

.handle-top {
  top: -2px;
}

.handle-bottom {
  bottom: -2px;
}

/* Critical Path Overlay */
.critical-path-overlay {
  position: absolute;
  top: 16px;
  left: 16px;
  right: 16px;
  bottom: 16px;
  pointer-events: none;
}

.critical-segment {
  position: absolute;
  background: #F56C6C;
  border-radius: 1px;
}

/* Selection Highlight */
.selection-highlight {
  position: absolute;
  border: 2px solid #67C23A;
  background: rgba(103, 194, 58, 0.1);
  pointer-events: none;
}

/* Current Time Indicator */
.current-time-indicator {
  position: absolute;
  top: 16px;
  bottom: 16px;
  width: 2px;
  background: #F56C6C;
  pointer-events: none;
}

.current-time-indicator::before {
  content: '';
  position: absolute;
  top: -6px;
  left: -4px;
  width: 10px;
  height: 10px;
  background: #F56C6C;
  border-radius: 50%;
}

/* Legend */
.minimap-legend {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  margin-top: 12px;
  padding: 12px;
  background: #fafafa;
  border-radius: 4px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #606266;
}

.legend-color {
  width: 16px;
  height: 16px;
  border-radius: 3px;
  border: 1px solid rgba(0, 0, 0, 0.1);
}

.legend-color.normal {
  background: #409EFF;
}

.legend-color.milestone {
  background: #E6A23C;
}

.legend-color.critical {
  background: #F56C6C;
}

.legend-color.delayed {
  background: #909399;
}

/* Zoom Controls */
.minimap-zoom {
  padding: 12px 16px;
  border-top: 1px solid #e4e7ed;
  display: flex;
  justify-content: center;
}

/* Minimized state */
.gantt-minimap.minimized .minimap-content,
.gantt-minimap.minimized .minimap-zoom {
  display: none;
}

/* Transition */
.el-collapse-enter-active,
.el-collapse-leave-active {
  transition: all 0.3s ease;
}

.el-collapse-enter-from,
.el-collapse-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

/* Responsive */
@media (max-width: 768px) {
  .minimap-legend {
    flex-direction: column;
    gap: 8px;
  }

  .minimap-canvas {
    height: 150px;
  }
}
</style>
