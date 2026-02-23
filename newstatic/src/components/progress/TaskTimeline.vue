<template>
  <div
    class="task-timeline"
    :class="{ 'is-grabbing': isCanvasDragging }"
    ref="timelineRef"
    :style="{ height: totalHeight + 'px', cursor: currentCursor }"
    @contextmenu="$emit('context-menu', $event)"
    @mousedown="handleCanvasMouseDown"
    @dblclick="handleCanvasDblClick"
    @mousemove="handleAllMouseMove($event)"
    @mouseup="handleCanvasMouseUp"
    @mouseleave="handleCanvasMouseUp"
  >
    <svg
      :width="svgWidth"
      :height="totalHeight"
      :style="{ background: 'linear-gradient(to right, rgba(0, 0, 0, 0.02) 1px, rgba(0, 0, 0, 0.02) 1px)', backgroundSize: dayWidth + 'px 100%' }"
    >
      <!-- 定义箭头标记 -->
      <defs>
        <!-- 普通依赖关系箭头 - 更小 -->
        <marker
          :id="arrowMarkerId"
          markerWidth="8"
          markerHeight="8"
          refX="7"
          refY="2.5"
          orient="0"
          markerUnits="strokeWidth"
        >
          <path d="M0,0 L0,5 L7,2.5 z" :fill="arrowColor" />
        </marker>

        <!-- 临时连线箭头 - 灰色更小 -->
        <marker
          id="temp-arrow"
          markerWidth="8"
          markerHeight="8"
          refX="7"
          refY="2.5"
          orient="auto"
          markerUnits="strokeWidth"
        >
          <path d="M0,0 L0,5 L7,2.5 z" fill="#999999" />
        </marker>
      </defs>

      <!-- 基线层 (在任务条下方) -->
      <g v-if="showBaseline" class="baseline-layer">
        <rect
          v-for="task in baselineTasks"
          :key="`baseline-${task.id}`"
          :x="task.baselineX"
          :y="task.baselineY"
          :width="task.baselineWidth"
          :height="8"
          fill="#c0c4cc"
          opacity="0.5"
          rx="2"
        />
      </g>

      <!-- 任务条层 -->
      <g class="task-bars-layer">
        <!-- 任务条 -->
        <g
          v-for="task in renderTasks"
          :key="`${task.id}-${task.x}-${task.y}`"
          :transform="`translate(${task.x}, ${task.y})`"
          class="task-bar-group"
          :class="{
            'is-selected': selectedTaskId === task.id,
            'is-critical': task.is_critical && showCriticalPath,
            'is-dragging': isDragging && draggedTask?.id === task.id,
            'is-preview': task.isPreview,
            'is-creating-dependency': isCreatingDependency && sourceTaskId === task.id,
            'is-dependency-target': isCreatingDependency && hoveredTargetTaskId === task.id,
            'is-dragging-over': timelineDragOverTaskId === task.id,
            'drop-before': timelineDragOverTaskId === task.id && timelineDropPosition === 'before',
            'drop-after': timelineDragOverTaskId === task.id && timelineDropPosition === 'after',
            'drop-child': timelineDragOverTaskId === task.id && timelineDropPosition === 'child'
          }"
          @click="handleTaskClick($event, task)"
          @dblclick="handleTaskDblClick($event, task)"
          @mouseenter="handleTaskMouseEnter"
          @mouseleave="handleTaskMouseLeave"
          @mousedown="handleTaskMouseDown($event, task)"
          @contextmenu.stop
          :draggable="!isCreatingDependency"
          @dragstart="handleDragStart($event, task)"
          @dragover.prevent="handleDragOver($event, task)"
          @dragleave="handleDragLeave($event, task)"
          @drop="handleDrop($event, task)"
          :data-task-id="task.id"
        >
          <!-- 任务条矩形 -->
          <rect
            v-if="!task.isMilestone"
            :x="0"
            :y="0"
            :width="task.width"
            :height="taskHeight"
            :fill="getTaskBarColor(task)"
            :stroke="task.is_critical && showCriticalPath ? '#f56c6c' : 'none'"
            :stroke-width="task.is_critical && showCriticalPath ? 2 : 0"
            :rx="4"
            class="task-bar-rect"
            style="filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1)); pointer-events: all;"
          />

          <!-- 里程碑 (菱形) -->
          <g v-else class="milestone-group">
            <rect
              :x="task.width / 2 - task.height / 2"
              :y="-task.height / 2"
              :width="task.height"
              :height="task.height"
              :fill="getTaskBarColor(task)"
              :transform="`rotate(45, ${task.width / 2}, 0)`"
              style="filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1)); pointer-events: all;"
            />
          </g>

          <!-- 进度条 (内部填充) -->
          <rect
            v-if="!task.isMilestone && task.progress > 0"
            :x="0"
            :y="0"
            :width="task.width * (Number(task.progress) / 100)"
            :height="taskHeight"
            :fill="getProgressBarColor(task)"
            :rx="4"
            opacity="0.3"
            style="pointer-events: none;"
          />

          <!-- 任务文本 -->
          <text
            :x="task.width / 2"
            :y="taskHeight / 2 + 3"
            text-anchor="middle"
            fill="white"
            font-size="10"
            font-weight="bold"
            style="pointer-events: none; text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);"
          >
            {{ task.isMilestone ? '' : (showTaskNames ? task.name : task.progress + '%') }}
          </text>

          <!-- 拖拽手柄 (左侧) -->
          <rect
            v-if="!task.isMilestone"
            :x="0"
            :y="0"
            :width="8"
            :height="taskHeight"
            fill="transparent"
            class="resize-handle-left"
            style="cursor: col-resize; pointer-events: all;"
          />

          <!-- 拖拽手柄 (右侧) -->
          <rect
            v-if="!task.isMilestone"
            :x="task.width - 8"
            :y="0"
            :width="8"
            :height="taskHeight"
            fill="transparent"
            class="resize-handle-right"
            style="cursor: col-resize; pointer-events: all;"
          />
        </g>
      </g>

      <!-- 临时连线层（用于创建依赖关系）- 灰色 -->
      <g v-if="isCreatingDependency && tempLinePath" class="temp-line-layer">
        <!-- 主连线 - 灰色虚线 -->
        <path
          :d="tempLinePath"
          stroke="#999999"
          stroke-width="1.5"
          fill="none"
          stroke-dasharray="6,3"
          marker-end="url(#temp-arrow)"
          class="temp-line-path"
        />
      </g>

      <!-- 依赖关系箭头层 - 黑色细线 -->
      <g
        v-if="showDependencies"
        class="dependencies-layer"
        style="pointer-events: none;"
      >
        <path
          v-for="dep in visibleDependencies"
          :key="dep.id"
          :d="dep.path"
          :stroke="dep.isCritical ? '#f56c6c' : '#303133'"
          :stroke-width="dep.isCritical ? 1.5 : 1"
          fill="none"
          :marker-end="`url(#${arrowMarkerId})`"
          class="dependency-path"
        />
      </g>

      <!-- 今天标记线 -->
      <line
        v-if="todayPosition !== null"
        :x1="todayPosition"
        :y1="0"
        :x2="todayPosition"
        :y2="totalHeight"
        :stroke="'#e6a23c'"
        :stroke-width="2"
        class="today-line"
        style="pointer-events: none;"
      />
    </svg>

    <!-- 拖拽提示框 -->
    <div
      v-if="isDragging && tooltipVisible"
      class="drag-tooltip"
      :style="{ left: tooltipPosition.x + 'px', top: tooltipPosition.y + 'px' }"
    >
      {{ tooltipText }}
    </div>

    <!-- 空状态 -->
    <div v-if="tasks.length === 0" class="timeline-empty">
      <el-empty :description="emptyDescription" />
    </div>
  </div>
</template>

<script setup>
import { computed, ref, onMounted, onUnmounted, watch } from 'vue'
import { isMilestone } from '@/utils/ganttHelpers'
import { ganttStore } from '@/stores/ganttStore'
import { useGanttDrag, DragMode } from '@/composables/useGanttDrag'
import { formatDate, addDays } from '@/utils/dateFormat'

const props = defineProps({
  tasks: {
    type: Array,
    default: () => []
  },
  rawTasks: {
    type: Array,
    default: () => []
  },
  selectedTaskId: {
    type: [Number, String],
    default: null
  },
  // 依赖关系创建相关
  isCreatingDependency: {
    type: Boolean,
    default: false
  },
  sourceTaskId: {
    type: [Number, String],
    default: null
  },
  tempLineEnd: {
    type: Object,
    default: null
  },
  viewMode: {
    type: String,
    default: 'day'
  },
  timelineDays: {
    type: Array,
    default: () => []
  },
  dayWidth: {
    type: Number,
    default: 40
  },
  rowHeight: {
    type: Number,
    default: 50
  },
  taskHeight: {
    type: Number,
    default: 24
  },
  showDependencies: {
    type: Boolean,
    default: true
  },
  showCriticalPath: {
    type: Boolean,
    default: true
  },
  showBaseline: {
    type: Boolean,
    default: false
  },
  visibleDependencies: {
    type: Array,
    default: () => []
  },
  isDragging: {
    type: Boolean,
    default: false
  },
  draggedTask: {
    type: Object,
    default: null
  },
  tooltipVisible: {
    type: Boolean,
    default: false
  },
  tooltipPosition: {
    type: Object,
    default: () => ({ x: 0, y: 0 })
  },
  tooltipText: {
    type: String,
    default: ''
  },
  previewTask: {
    type: Object,
    default: null
  },
  dragMode: {
    type: String,
    default: 'none'
  },
  todayPosition: {
    type: Number,
    default: null
  },
  arrowMarkerId: {
    type: String,
    default: 'arrow-marker'
  },
  arrowColor: {
    type: String,
    default: '#909399'
  },
  timelineWidth: {
    type: Number,
    default: 800
  },
  emptyDescription: {
    type: String,
    default: '暂无进度计划数据'
  },
  showTaskNames: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits([
  'task-click',
  'task-dblclick',
  'task-mousedown',
  'context-menu',
  'dependency-create',
  'mousemove',
  'task-dragged'
])

const timelineRef = ref(null)

// ==================== 鼠标状态追踪 ====================
const mouseState = ref({
  isRightButtonDown: false,
  isCtrlPressed: false,
  isOverTask: false
})

// ==================== 拖拽排序状态 ====================
const timelineDraggedTask = ref(null)
const timelineDragOverTaskId = ref(null)
const timelineDropPosition = ref(null) // 'before', 'after', or 'child'

// 当前应该显示的 cursor 样式
const currentCursor = computed(() => {
  if (props.isCreatingDependency) {
    return 'crosshair' // 创建依赖关系时显示十字
  }
  if (mouseState.value.isRightButtonDown || isCanvasDragging.value) {
    return 'grabbing' // 右键拖动或画布拖动时显示抓取
  }
  if (mouseState.value.isCtrlPressed && mouseState.value.isOverTask) {
    return 'crosshair' // Ctrl + 任务上显示十字
  }
  return 'default' // 默认箭头
})

// ==================== 画布拖动状态 ====================
const isCanvasDragging = ref(false)
const canvasDragStart = ref({ x: 0, y: 0 })
const canvasScrollStart = ref({ left: 0, top: 0 })

// 获取父容器的滚动元素
const getScrollParent = () => {
  return timelineRef.value?.parentElement
}

// 画布鼠标按下
const handleCanvasMouseDown = (event) => {
  // 如果点击的是任务条或其子元素，不触发画布拖动
  if (event.target.closest('.task-bar-group')) {
    return
  }

  // 右键拖动画布
  if (event.button === 2) {
    mouseState.value.isRightButtonDown = true
    isCanvasDragging.value = true
    canvasDragStart.value = { x: event.clientX, y: event.clientY }

    const scrollParent = getScrollParent()
    if (scrollParent) {
      canvasScrollStart.value = {
        left: scrollParent.scrollLeft,
        top: scrollParent.scrollTop
      }
    }
  }
}

// 画布双击 - 快速添加任务
const handleCanvasDblClick = (event) => {
  // 如果双击的是任务条，不触发
  if (event.target.closest('.task-bar-group')) {
    return
  }

  // 发送新建任务事件
  emit('context-menu', {
    event,
    type: 'new-task',
    action: 'create-immediate'
  })
}

// 画布鼠标移动
const handleCanvasMouseMove = (event) => {
  if (!isCanvasDragging.value) return

  const scrollParent = getScrollParent()
  if (!scrollParent) return

  const deltaX = event.clientX - canvasDragStart.value.x
  const deltaY = event.clientY - canvasDragStart.value.y

  // 更新滚动位置（反向移动）
  scrollParent.scrollLeft = canvasScrollStart.value.left - deltaX
  scrollParent.scrollTop = canvasScrollStart.value.top - deltaY
}

// 鼠标移动（用于画线连接）
const handleMouseMove = (event) => {
  // 通知父组件鼠标位置，用于绘制临时连线
  if (props.isCreatingDependency && timelineRef.value) {
    const rect = timelineRef.value.getBoundingClientRect()
    emit('mousemove', {
      x: event.clientX - rect.left,
      y: event.clientY - rect.top
    })
  }
}

// 统一处理鼠标移动
const handleAllMouseMove = (event) => {
  // 更新 Ctrl 键状态
  mouseState.value.isCtrlPressed = event.ctrlKey || event.metaKey
  handleCanvasMouseMove(event)
  handleMouseMove(event)
}

// 任务条鼠标进入
const handleTaskMouseEnter = () => {
  mouseState.value.isOverTask = true
}

// 任务条鼠标离开
const handleTaskMouseLeave = () => {
  mouseState.value.isOverTask = false
}

// 画布鼠标松开
const handleCanvasMouseUp = () => {
  isCanvasDragging.value = false
  mouseState.value.isRightButtonDown = false
}

// 全局鼠标松开（防止鼠标移出元素后松开）
const handleGlobalMouseUp = () => {
  isCanvasDragging.value = false
  mouseState.value.isRightButtonDown = false
}

// 全局键盘按下（追踪 Ctrl 键）
const handleGlobalKeyDown = (event) => {
  if (event.ctrlKey || event.metaKey) {
    mouseState.value.isCtrlPressed = true
  }
}

// 全局键盘松开（追踪 Ctrl 键）
const handleGlobalKeyUp = (event) => {
  if (!event.ctrlKey && !event.metaKey) {
    mouseState.value.isCtrlPressed = false
  }
}

onMounted(() => {
  document.addEventListener('mouseup', handleGlobalMouseUp)
  document.addEventListener('keydown', handleGlobalKeyDown)
  document.addEventListener('keyup', handleGlobalKeyUp)

  // 初始化容器宽度
  updateContainerWidth()

  // 监听容器宽度变化
  const resizeObserver = new ResizeObserver(() => {
    updateContainerWidth()
  })

  if (timelineRef.value?.parentElement) {
    resizeObserver.observe(timelineRef.value.parentElement)
  }

  // 保存observer引用以便清理
  timelineRef.value._resizeObserver = resizeObserver
})

onUnmounted(() => {
  document.removeEventListener('mouseup', handleGlobalMouseUp)
  document.removeEventListener('keydown', handleGlobalKeyDown)
  document.removeEventListener('keyup', handleGlobalKeyUp)

  // 清理ResizeObserver
  if (timelineRef.value?._resizeObserver) {
    timelineRef.value._resizeObserver.disconnect()
  }
})

// ==================== 计算属性 ====================
// 容器宽度监听
const containerWidth = ref(800)

// SVG宽度：取容器宽度和计算出的时间轴宽度中的较大值
const svgWidth = computed(() => {
  const calculatedWidth = props.timelineWidth || 800
  return Math.max(containerWidth.value, calculatedWidth)
})

// 监听容器宽度变化
const updateContainerWidth = () => {
  if (timelineRef.value?.parentElement) {
    containerWidth.value = timelineRef.value.parentElement.clientWidth
  }
}

// 获取可见任务（过滤掉被折叠隐藏的任务）
const { actions } = ganttStore
const visibleTasks = computed(() => {
  return props.tasks.filter(task => !actions.isTaskHidden(task.id, props.tasks))
})

// 计算总高度
const totalHeight = computed(() => {
  return visibleTasks.value.length * props.rowHeight
})

// 准备渲染的任务数据
const renderTasks = computed(() => {
  const timelineStart = props.timelineDays[0]?.date
  if (!timelineStart) return []

  return visibleTasks.value.map((task, index) => {
    // 如果正在拖拽此任务且有预览数据，使用预览数据
    const isDraggingThisTask = props.isDragging && props.draggedTask?.id === task.id
    const taskData = isDraggingThisTask && props.previewTask ? props.previewTask : task

    const taskStart = new Date(taskData.start)
    const taskEnd = new Date(taskData.end)
    const timelineStartDate = new Date(timelineStart)

    // 计算位置和宽度
    const daysDiff = Math.ceil((taskStart - timelineStartDate) / (1000 * 60 * 60 * 24))
    const duration = Math.ceil((taskEnd - taskStart) / (1000 * 60 * 60 * 24))

    const x = daysDiff * props.dayWidth
    const width = duration * props.dayWidth

    // Y位置：居中于行
    const y = index * props.rowHeight + (props.rowHeight - props.taskHeight) / 2

    const taskIsMilestone = isMilestone(taskData)

    return {
      id: task.id,
      name: taskData.name,
      x,
      y,
      width,
      isMilestone: taskIsMilestone,
      height: taskIsMilestone ? props.taskHeight : props.taskHeight,
      progress: taskData.progress || 0,
      status: taskData.status,
      is_critical: taskData.is_critical,
      isPreview: isDraggingThisTask
    }
  })
})

// 准备基线数据
const baselineTasks = computed(() => {
  if (!props.showBaseline) return []

  const timelineStart = props.timelineDays[0]?.date
  if (!timelineStart) return []

  return visibleTasks.value.map((task, index) => {
    if (!task.baseline_start || !task.baseline_end) return null

    const taskStart = new Date(task.baseline_start)
    const taskEnd = new Date(task.baseline_end)
    const timelineStartDate = new Date(timelineStart)

    const daysDiff = Math.ceil((taskStart - timelineStartDate) / (1000 * 60 * 60 * 24))
    const duration = Math.ceil((taskEnd - taskStart) / (1000 * 60 * 60 * 24))

    const x = daysDiff * props.dayWidth
    const width = duration * props.dayWidth
    const y = index * props.rowHeight + 12

    return {
      id: task.id,
      baselineX: x,
      baselineY: y,
      baselineWidth: width
    }
  }).filter(Boolean)
})

// ==================== 任务操作 ====================
// 当前鼠标悬停的任务（用于依赖连线创建时的视觉反馈）
const hoveredTargetTaskId = computed(() => {
  if (!props.isCreatingDependency || !props.tempLineEnd) {
    return null
  }

  const mouseX = props.tempLineEnd.x
  const mouseY = props.tempLineEnd.y

  // 检查鼠标是否在某个任务范围内
  for (const task of renderTasks.value) {
    if (task.id === props.sourceTaskId) continue // 跳过源任务本身

    const taskLeft = task.x
    const taskRight = task.x + task.width
    const taskTop = task.y
    const taskBottom = task.y + props.taskHeight

    // 扩大检测范围，让用户更容易选中目标任务
    const padding = 10

    if (
      mouseX >= taskLeft - padding &&
      mouseX <= taskRight + padding &&
      mouseY >= taskTop - padding &&
      mouseY <= taskBottom + padding
    ) {
      return task.id
    }
  }

  return null
})

// 临时连线路径（使用平滑的贝塞尔曲线）
const tempLinePath = computed(() => {
  if (!props.isCreatingDependency || !props.sourceTaskId || !props.tempLineEnd) {
    return ''
  }

  // 找到源任务的位置
  const sourceTask = renderTasks.value.find(t => t.id === props.sourceTaskId)
  if (!sourceTask) return ''

  // 源任务的结束位置（右侧中心）
  const x1 = sourceTask.x + sourceTask.width
  const y1 = sourceTask.y + props.taskHeight / 2

  // 当前鼠标位置
  const x2 = props.tempLineEnd.x
  const y2 = props.tempLineEnd.y

  // 计算控制点偏移量
  const curveOffset = Math.min(Math.abs(x2 - x1) * 0.5, 100) // 最大偏移 100px

  if (x1 < x2) {
    // 目标在源任务右侧：使用平滑的 S 型曲线
    // 控制点让线条先水平延伸，再平滑过渡到目标点
    const cp1x = x1 + curveOffset * 0.6
    const cp1y = y1
    const cp2x = x2 - curveOffset * 0.6
    const cp2y = y2

    return `M ${x1} ${y1} C ${cp1x} ${cp1y}, ${cp2x} ${cp2y}, ${x2} ${y2}`
  } else {
    // 目标在源任务左侧：使用 U 型曲线
    const midX = x1 + curveOffset + 40
    const cp1x = x1 + curveOffset * 0.5
    const cp1y = y1
    const cp2x = midX
    const cp2y = y1
    const cp3x = midX
    const cp3y = y2
    const cp4x = x2 + curveOffset * 0.5
    const cp4y = y2

    return `M ${x1} ${y1} C ${cp1x} ${cp1y}, ${cp2x} ${cp2y}, ${midX} ${(y1 + y2) / 2} S ${cp4x} ${cp4y}, ${x2} ${y2}`
  }
})

// 获取任务条颜色
const getTaskBarColor = (task) => {
  if (task.isMilestone) return '#e6a23c'

  const colors = {
    completed: '#67c23a',
    in_progress: '#409eff',
    not_started: '#909399',
    delayed: '#f56c6c'
  }

  if (task.is_critical && props.showCriticalPath) {
    const criticalColors = {
      completed: '#e85d04',
      in_progress: '#f56c6c',
      not_started: '#e85d75',
      delayed: '#f56c6c'
    }
    return criticalColors[task.status] || colors[task.status]
  }

  return colors[task.status] || '#909399'
}

// 获取进度条颜色
const getProgressBarColor = (task) => {
  return getTaskBarColor(task)
}

// 任务点击
const handleTaskClick = (event, task) => {
  event.stopPropagation()

  // 如果正在创建依赖关系，触发依赖创建
  if (props.isCreatingDependency) {
    emit('dependency-create', task.id)
    return
  }

  // 获取原始任务数据
  const rawTask = props.rawTasks.find(t => t.id === task.id)
  emit('task-click', rawTask || { id: task.id })
}

// 任务双击
const handleTaskDblClick = (event, task) => {
  event.stopPropagation()
  // 获取原始任务数据
  const rawTask = props.rawTasks.find(t => t.id === task.id)
  emit('task-dblclick', rawTask || { id: task.id })
}

// 任务鼠标按下（开始拖拽或画线连接）
const handleTaskMouseDown = (event, task) => {
  event.stopPropagation()
  if (event.button !== 0) return // 只响应左键

  // 如果正在创建依赖关系，不做任何处理（等待点击目标任务）
  if (props.isCreatingDependency) {
    return
  }

  // 检查是否按住了Shift键或Ctrl键，用于创建依赖关系
  if (event.shiftKey || event.ctrlKey) {
    // 开始创建依赖关系
    const rawTask = props.rawTasks.find(t => t.id === task.id)
    emit('dependency-create', rawTask || task)
    return
  }

  // 获取任务条DOM元素（用于检测拖拽边缘）
  const taskBarElement = event.target.closest('.task-bar-group')

  // 获取原始任务数据
  const rawTask = props.rawTasks.find(t => t.id === task.id)
  emit('task-mousedown', event, rawTask || task, taskBarElement)
}

// ==================== 拖拽排序处理 ====================
const handleDragStart = (event, task) => {
  if (props.isCreatingDependency) {
    event.preventDefault()
    return
  }

  timelineDraggedTask.value = task
  timelineDragOverTaskId.value = null
  timelineDropPosition.value = null

  event.dataTransfer.effectAllowed = 'move'
  event.dataTransfer.setData('text/plain', task.id.toString())

  console.log('[TaskTimeline] 开始拖拽任务:', task.name)
}

const handleDragOver = (event, task) => {
  event.preventDefault()

  if (!timelineDraggedTask.value || timelineDraggedTask.value.id === task.id) {
    return
  }

  timelineDragOverTaskId.value = task.id

  // 获取任务条的位置信息
  const taskElement = event.currentTarget
  const rect = taskElement.getBoundingClientRect()
  const relativeY = event.clientY - rect.top
  const rowHeight = props.rowHeight

  // 上部 1/3: 插入之前
  // 中部 1/3: 作为子任务
  // 下部 1/3: 插入之后
  if (relativeY < rowHeight / 3) {
    timelineDropPosition.value = 'before'
  } else if (relativeY > rowHeight * 2 / 3) {
    timelineDropPosition.value = 'after'
  } else {
    timelineDropPosition.value = 'child'
  }
}

const handleDragLeave = (event, task) => {
  if (timelineDragOverTaskId.value === task.id) {
    timelineDragOverTaskId.value = null
    timelineDropPosition.value = null
  }
}

const handleDrop = (event, targetTask) => {
  event.preventDefault()
  event.stopPropagation()

  if (!timelineDraggedTask.value) return
  if (timelineDraggedTask.value.id === targetTask.id) {
    timelineDraggedTask.value = null
    timelineDragOverTaskId.value = null
    timelineDropPosition.value = null
    return
  }

  console.log('[TaskTimeline] 拖拽完成:', {
    from: timelineDraggedTask.value.name,
    to: targetTask.name,
    position: timelineDropPosition.value
  })

  // 获取原始任务对象
  const fromRawTask = props.rawTasks.find(t => t.id === timelineDraggedTask.value.id)
  const toRawTask = props.rawTasks.find(t => t.id === targetTask.id)

  emit('task-dragged', {
    fromTask: fromRawTask || timelineDraggedTask.value,
    toTask: toRawTask || targetTask,
    position: timelineDropPosition.value || 'child'
  })

  // 重置状态
  timelineDraggedTask.value = null
  timelineDragOverTaskId.value = null
  timelineDropPosition.value = null
}

defineExpose({
  timelineRef
})
</script>

<style scoped>
.task-timeline {
  flex: 1;
  position: relative;
  min-width: 800px;
  overflow: hidden;
  user-select: none;
}

.task-bar-group {
  transition: filter 0.3s;
}

.task-bar-group:hover {
  filter: brightness(1.1);
}

.task-bar-group.is-selected {
  filter: brightness(1.05);
}

.task-bar-group.is-dragging {
  opacity: 0.8;
}

.task-bar-group.is-preview {
  opacity: 0.7;
}

.task-bar-group.is-creating-dependency {
  filter: brightness(1.2) drop-shadow(0 0 8px rgba(64, 158, 255, 0.6));
}

/* 依赖连线目标高亮 */
.task-bar-group.is-dependency-target {
  filter: brightness(1.3) drop-shadow(0 0 12px rgba(64, 158, 255, 0.9));
  animation: pulse-target 1s ease-in-out infinite;
}

@keyframes pulse-target {
  0%, 100% {
    filter: brightness(1.2) drop-shadow(0 0 10px rgba(64, 158, 255, 0.8));
  }
  50% {
    filter: brightness(1.4) drop-shadow(0 0 16px rgba(64, 158, 255, 1));
  }
}

/* 拖拽排序视觉反馈 */
.task-bar-group.is-dragging-over {
  filter: brightness(1.15);
}

.task-bar-group.drop-before {
  filter: brightness(1.15) drop-shadow(0 -4px 0 #409eff);
}

.task-bar-group.drop-after {
  filter: brightness(1.15) drop-shadow(0 4px 0 #409eff);
}

.task-bar-group.drop-child {
  filter: brightness(1.2) drop-shadow(0 0 6px rgba(64, 158, 255, 0.8));
}

/* 临时连线样式 */
.temp-line-layer {
  pointer-events: none;
}

.temp-line-path {
  /* 简单的灰色虚线，无动画 */
}

.dependency-path {
  transition: stroke-width 0.2s;
}

.dependency-path:hover {
  stroke-width: 2 !important;
}

.drag-tooltip {
  position: fixed;
  background: rgba(0, 0, 0, 0.8);
  color: white;
  padding: 8px 12px;
  border-radius: 4px;
  font-size: 12px;
  z-index: 9999;
  pointer-events: none;
  white-space: nowrap;
}

.timeline-empty {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}
</style>
