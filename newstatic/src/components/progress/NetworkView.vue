<template>
  <div class="network-view" ref="containerRef">
    <!-- SVG 绘图区域 -->
    <div class="network-canvas-container" ref="canvasContainerRef">
      <svg
        ref="svgRef"
        class="network-svg"
        :width="svgWidth"
        :height="svgHeight"
        @wheel.prevent="handleWheel"
        @mousedown="handleMouseDown"
        @mousemove="handleMouseMove"
        @mouseup="handleMouseUp"
        @mouseleave="handleMouseUp"
      >
        <defs>
          <!-- 任务箭头标记 - 普通任务 -->
          <marker
            id="arrowhead-task"
            markerWidth="10"
            markerHeight="7"
            refX="9"
            refY="3.5"
            orient="auto"
            markerUnits="userSpaceOnUse"
          >
            <path d="M0,0 L0,7 L9,3.5 z" fill="#64B5F6" />
          </marker>
          <!-- 任务箭头标记 - 关键任务 -->
          <marker
            id="arrowhead-critical"
            markerWidth="10"
            markerHeight="7"
            refX="9"
            refY="3.5"
            orient="auto"
            markerUnits="userSpaceOnUse"
          >
            <path d="M0,0 L0,7 L9,3.5 z" fill="#FF8A65" />
          </marker>
          <!-- 依赖关系箭头标记 -->
          <marker
            id="arrowhead-dependency"
            markerWidth="8"
            markerHeight="8"
            refX="7"
            refY="4"
            orient="auto"
            markerUnits="userSpaceOnUse"
          >
            <path d="M0,0 L0,8 L8,4 z" fill="#999" />
          </marker>
          <!-- 节点阴影 -->
          <filter id="nodeShadow" x="-50%" y="-50%" width="200%" height="200%">
            <feDropShadow dx="0" dy="2" stdDeviation="3" flood-opacity="0.2"/>
          </filter>
        </defs>

        <!-- 主绘图组（应用平移） -->
        <g :transform="`translate(${panX}, ${panY})`">
          <!-- 网格背景 -->
          <g class="grid-background" v-if="showGrid">
            <pattern
              id="gridPattern"
              :width="gridSize"
              :height="gridSize"
              patternUnits="userSpaceOnUse"
            >
              <path
                :d="`M ${gridSize} 0 L 0 0 0 ${gridSize}`"
                fill="none"
                stroke="var(--border-color-extra-light, #e4e7ed)"
                stroke-width="0.5"
              />
            </pattern>
            <rect width="100%" height="100%" fill="url(#gridPattern)" />
          </g>

          <!-- 依赖关系线（虚线） -->
          <g class="dependency-lines">
            <path
              v-for="dep in renderableDependencies"
              :key="dep.key"
              :d="dep.path"
              stroke="#999"
              stroke-width="1"
              stroke-dasharray="4,4"
              fill="none"
              :marker-end="`url(#arrowhead-dependency)`"
              class="dependency-line"
            />
          </g>

          <!-- 任务箭头（活动） -->
          <g class="task-arrows">
            <g
              v-for="task in taskArrows"
              :key="task.id"
              :class="{
                'is-critical': task.isCritical,
                'is-selected': selectedTaskId === task.id
              }"
              class="task-arrow-group"
              @click="handleTaskClick($event, task)"
              @dblclick="handleTaskDblClick($event, task)"
            >
              <!-- 任务箭头线 -->
              <path
                :d="task.path"
                :stroke="task.stroke"
                :stroke-width="task.strokeWidth"
                fill="none"
                :marker-end="`url(#${task.marker})`"
                class="task-arrow"
              />

              <!-- 任务名称标签 -->
              <text
                v-if="showTaskNames"
                :x="task.labelX"
                :y="task.labelY"
                text-anchor="middle"
                class="task-label"
                :fill="task.textColor"
              >
                {{ task.label }}
              </text>

              <!-- 工期信息 -->
              <text
                v-if="showDuration"
                :x="task.durationX"
                :y="task.durationY"
                text-anchor="middle"
                class="task-duration"
              >
                {{ task.duration }}天
              </text>
            </g>
          </g>

          <!-- 事件节点 -->
          <g class="event-nodes">
            <g
              v-for="node in eventNodes"
              :key="node.id"
              :transform="`translate(${node.x}, ${node.y})`"
              :class="{
                'is-selected': selectedNodeId === node.id,
                'is-critical': node.isCritical
              }"
              class="node-group"
              @mousedown="handleNodeMouseDown($event, node)"
              @click="handleNodeClick($event, node)"
            >
              <!-- 节点圆形 -->
              <circle
                :r="nodeRadius"
                :fill="node.fill"
                :stroke="node.stroke"
                stroke-width="2"
                filter="url(#nodeShadow)"
                class="node-circle"
              />

              <!-- 节点编号 -->
              <text
                y="5"
                text-anchor="middle"
                class="node-number"
                fill="#333"
                font-weight="bold"
              >
                {{ node.number }}
              </text>

              <!-- 时间参数 -->
              <g v-if="showTimeParams" class="node-time-params">
                <text :y="nodeRadius + 14" text-anchor="middle" class="node-time-text">
                  ES: {{ formatTime(node.ES) }}
                </text>
                <text :y="nodeRadius + 26" text-anchor="middle" class="node-time-text">
                  EF: {{ formatTime(node.EF) }}
                </text>
                <text :y="nodeRadius + 38" text-anchor="middle" class="node-time-text">
                  LS: {{ formatTime(node.LS) }}
                </text>
                <text :y="nodeRadius + 50" text-anchor="middle" class="node-time-text">
                  LF: {{ formatTime(node.LF) }}
                </text>
              </g>
            </g>
          </g>
        </g>
      </svg>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { isDummyTask } from '@/utils/ganttHelpers'

// Props
const props = defineProps({
  // 任务数据
  tasks: {
    type: Array,
    default: () => []
  },
  // 依赖关系
  dependencies: {
    type: Array,
    default: () => []
  },
  // 任务索引映射（用于对齐任务列表）
  taskIndexMap: {
    type: Object,
    default: () => ({})
  },
  // 时间轴数据（用于计算X坐标）
  timelineStart: {
    type: String,
    default: ''
  },
  // 行高（用于计算Y坐标）
  rowHeight: {
    type: Number,
    default: 40
  },
  // 每天宽度（用于计算X坐标）
  dayWidth: {
    type: Number,
    default: 40
  },
  // 是否对齐任务列表
  alignWithTaskList: {
    type: Boolean,
    default: true
  },
  // 视图选项
  showCriticalPath: {
    type: Boolean,
    default: true
  },
  showTaskNames: {
    type: Boolean,
    default: true
  },
  showTimeParams: {
    type: Boolean,
    default: false
  },
  showSlack: {
    type: Boolean,
    default: false
  },
  showDuration: {
    type: Boolean,
    default: true
  },
  // 视图控制
  toolMode: {
    type: String,
    default: 'select'
  },
  // 样式选项
  nodeRadius: {
    type: Number,
    default: 18
  },
  svgWidth: {
    type: Number,
    default: 2000
  },
  svgHeight: {
    type: Number,
    default: 1200
  }
})

// Emits
const emit = defineEmits([
  'node-click',
  'task-click',
  'task-dblclick',
  'zoom-change'
])

// Refs
const containerRef = ref(null)
const canvasContainerRef = ref(null)
const svgRef = ref(null)

// 状态
const showGrid = ref(true)
const gridSize = ref(20)
const panX = ref(0)
const panY = ref(0)

// 拖拽状态
const isDragging = ref(false)
const isPanning = ref(false)
const draggedNode = ref(null)
const dragStartPos = ref({ x: 0, y: 0 })

// 选中状态
const selectedNodeId = ref(null)
const selectedTaskId = ref(null)

// 计算时间轴起始日期
const timelineStartDate = computed(() => {
  return props.timelineStart ? new Date(props.timelineStart) : null
})

// 过滤掉虚任务（有子任务的父任务）
const realTasks = computed(() => {
  return props.tasks.filter(task => !isDummyTask(task, props.tasks))
})

// 创建事件节点和任务箭头
const eventNodes = computed(() => {
  if (!realTasks.value.length || !timelineStartDate.value) return []

  const nodes = []
  const nodeMap = new Map() // 用于合并相同位置的节点
  let nodeNumber = 1

  // 为每个任务创建起点和终点节点
  realTasks.value.forEach(task => {
    const taskStart = new Date(task.start)
    const taskEnd = new Date(task.end)

    const startDays = Math.ceil((taskStart - timelineStartDate.value) / (1000 * 60 * 60 * 24))
    const endDays = Math.ceil((taskEnd - timelineStartDate.value) / (1000 * 60 * 60 * 24))

    // 创建起点节点
    const startKey = `${startDays}`
    if (!nodeMap.has(startKey)) {
      const startNode = {
        id: `start-${task.id}`,
        number: nodeNumber++,
        x: startDays * props.dayWidth,
        y: 0, // 稍后计算
        ES: task.early_start,
        EF: task.early_finish,
        LS: task.late_start,
        LF: task.late_finish,
        isCritical: task.is_critical,
        fill: getEventNodeFill(task),
        stroke: getEventNodeStroke(task),
        tasks: [task.id]
      }
      nodeMap.set(startKey, startNode)
      nodes.push(startNode)
    } else {
      const existingNode = nodeMap.get(startKey)
      existingNode.tasks.push(task.id)
    }

    // 创建终点节点
    const endKey = `${endDays}`
    if (!nodeMap.has(endKey)) {
      const endNode = {
        id: `end-${task.id}`,
        number: nodeNumber++,
        x: endDays * props.dayWidth,
        y: 0, // 稍后计算
        ES: task.early_start,
        EF: task.early_finish,
        LS: task.late_start,
        LF: task.late_finish,
        isCritical: task.is_critical,
        fill: getEventNodeFill(task),
        stroke: getEventNodeStroke(task),
        tasks: [task.id]
      }
      nodeMap.set(endKey, endNode)
      nodes.push(endNode)
    } else {
      const existingNode = nodeMap.get(endKey)
      existingNode.tasks.push(task.id)
    }
  })

  // 计算节点的Y位置（基于关联任务的行索引）
  nodes.forEach(node => {
    if (node.tasks.length > 0 && props.alignWithTaskList) {
      // 使用第一个关联任务的行索引
      const firstTaskId = node.tasks[0]
      const taskIndex = props.taskIndexMap[firstTaskId]
      if (taskIndex !== undefined) {
        node.y = taskIndex * props.rowHeight + props.rowHeight / 2
      } else {
        node.y = 100
      }
    } else {
      node.y = 100 + node.number * 60
    }
  })

  return nodes
})

// 任务箭头（连接起点和终点节点）
const taskArrows = computed(() => {
  if (!realTasks.value.length || !timelineStartDate.value) return []

  return realTasks.value.map(task => {
    const taskStart = new Date(task.start)
    const taskEnd = new Date(task.end)

    const startDays = Math.ceil((taskStart - timelineStartDate.value) / (1000 * 60 * 60 * 24))
    const endDays = Math.ceil((taskEnd - timelineStartDate.value) / (1000 * 60 * 60 * 24))

    let startX, startY, endX, endY

    if (props.alignWithTaskList && props.taskIndexMap[task.id] !== undefined) {
      // 对齐模式
      const taskY = props.taskIndexMap[task.id] * props.rowHeight + props.rowHeight / 2
      startX = startDays * props.dayWidth
      startY = taskY
      endX = endDays * props.dayWidth
      endY = taskY
    } else {
      // 自由模式
      startX = startDays * props.dayWidth
      startY = 100 + props.taskIndexMap[task.id] * props.rowHeight
      endX = endDays * props.dayWidth
      endY = startY
    }

    // 计算箭头路径（带弯曲度，避免重叠）
    const path = calculateTaskArrowPath(startX, startY, endX, endY)

    const duration = Math.ceil((taskEnd - taskStart) / (1000 * 60 * 60 * 24))

    return {
      id: task.id,
      path,
      label: task.name,
      stroke: task.is_critical ? '#FF8A65' : '#64B5F6',
      strokeWidth: task.is_critical ? 2.5 : 2,
      marker: task.is_critical ? 'arrowhead-critical' : 'arrowhead-task',
      textColor: '#606266',
      isCritical: task.is_critical,
      duration,
      labelX: (startX + endX) / 2,
      labelY: startY - 8,
      durationX: (startX + endX) / 2,
      durationY: startY + 18
    }
  })
})

// 依赖关系（用于可视化前置任务关系）
const renderableDependencies = computed(() => {
  if (!props.dependencies.length) return []

  return props.dependencies.map((dep, index) => {
    const fromTask = props.tasks.find(t => t.id === dep.depends_on)
    const toTask = props.tasks.find(t => t.id === dep.task_id)

    if (!fromTask || !toTask) return null

    const fromTaskStart = new Date(fromTask.end)
    const toTaskStart = new Date(toTask.start)

    const fromDays = Math.ceil((fromTaskStart - timelineStartDate.value) / (1000 * 60 * 60 * 24))
    const toDays = Math.ceil((toTaskStart - timelineStartDate.value) / (1000 * 60 * 60 * 24))

    let fromX, fromY, toX, toY

    if (props.alignWithTaskList) {
      const fromTaskIndex = props.taskIndexMap[fromTask.id] ?? 0
      const toTaskIndex = props.taskIndexMap[toTask.id] ?? 0
      fromX = fromDays * props.dayWidth
      fromY = fromTaskIndex * props.rowHeight + props.rowHeight / 2
      toX = toDays * props.dayWidth
      toY = toTaskIndex * props.rowHeight + props.rowHeight / 2
    } else {
      fromX = fromDays * props.dayWidth
      fromY = 100
      toX = toDays * props.dayWidth
      toY = 100
    }

    // 计算依赖线路径（正交路径）
    const midX = (fromX + toX) / 2
    const path = `M ${fromX} ${fromY} L ${midX} ${fromY} L ${midX} ${toY} L ${toX} ${toY}`

    return {
      key: `${dep.depends_on}-${dep.task_id}-${index}`,
      path
    }
  }).filter(Boolean)
})

// 计算任务箭头路径（带轻微弧度）
function calculateTaskArrowPath(startX, startY, endX, endY) {
  const midX = (startX + endX) / 2
  const curvature = 0  // 可以设置弯曲度
  const controlY = startY + curvature

  if (curvature === 0) {
    // 直线
    return `M ${startX} ${startY} L ${endX} ${endY}`
  } else {
    // 贝塞尔曲线
    return `M ${startX} ${startY} Q ${midX} ${controlY} ${endX} ${endY}`
  }
}

// 获取事件节点填充色
function getEventNodeFill(task) {
  if (task.is_critical) return '#FFF3E0'
  return '#FFFFFF'
}

// 获取事件节点边框色
function getEventNodeStroke(task) {
  if (task.is_critical) return '#FF8A65'
  return '#1976D2'
}

// 格式化时间显示
function formatTime(timestamp) {
  if (!timestamp) return '-'
  const date = new Date(timestamp * 1000)
  return date.toLocaleDateString('zh-CN', { month: '2-digit', day: '2-digit' })
}

// 鼠标滚轮缩放
function handleWheel(event) {
  const delta = event.deltaY > 0 ? -2 : 2
  // 直接返回新的 dayWidth 值，与甘特图保持一致
  emit('zoom-change', Math.max(10, Math.min(100, props.dayWidth + delta)))
}

// 鼠标按下
function handleMouseDown(event) {
  // 如果点击的是背景，开始平移
  if (event.target.tagName === 'svg' || event.target.tagName === 'rect') {
    isPanning.value = true
    dragStartPos.value = { x: event.clientX - panX.value, y: event.clientY - panY.value }
  }
}

// 节点鼠标按下
function handleNodeMouseDown(event, node) {
  event.stopPropagation()
  if (props.toolMode === 'pan') return
  isDragging.value = true
  draggedNode.value = node
  dragStartPos.value = { x: event.clientX - node.x, y: event.clientY - node.y }
}

// 鼠标移动
function handleMouseMove(event) {
  if (isPanning.value) {
    panX.value = event.clientX - dragStartPos.value.x
    panY.value = event.clientY - dragStartPos.value.y
  }

  if (isDragging.value && draggedNode.value && props.toolMode === 'select') {
    const rect = svgRef.value.getBoundingClientRect()
    draggedNode.value.x = (event.clientX - rect.left - dragStartPos.value.x)
    draggedNode.value.y = (event.clientY - rect.top - dragStartPos.value.y)
  }
}

// 鼠标释放
function handleMouseUp() {
  isDragging.value = false
  isPanning.value = false
  draggedNode.value = null
}

// 节点点击
function handleNodeClick(event, node) {
  event.stopPropagation()
  selectedNodeId.value = node.id
  selectedTaskId.value = null
  emit('node-click', node)
}

// 任务点击
function handleTaskClick(event, task) {
  event.stopPropagation()
  selectedTaskId.value = task.id
  selectedNodeId.value = null
  emit('task-click', task)
}

// 任务双击
function handleTaskDblClick(event, task) {
  event.stopPropagation()
  emit('task-dblclick', task)
}

// 键盘事件
function handleKeydown(event) {
  if (event.key === 'Escape') {
    selectedNodeId.value = null
    selectedTaskId.value = null
  }
}

// 生命周期
onMounted(() => {
  window.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
})
</script>

<style scoped>
.network-view {
  width: 100%;
  height: 100%;
  overflow: hidden;
  position: relative;
}

.network-canvas-container {
  width: 100%;
  height: 100%;
  overflow: auto;
  background: var(--bg-color-page, #f5f7fa);
}

.network-svg {
  display: block;
  background: white;
  cursor: grab;
}

.network-svg:active {
  cursor: grabbing;
}

.node-group {
  cursor: pointer;
  transition: opacity 0.2s;
}

.node-group:hover {
  opacity: 0.8;
}

.node-group.is-selected .node-circle {
  stroke: #409EFF;
  stroke-width: 3;
}

.node-circle {
  transition: all 0.2s;
}

.node-number {
  font-size: 12px;
  font-weight: bold;
  pointer-events: none;
}

.node-time-params {
  font-size: 9px;
  fill: #909399;
  pointer-events: none;
}

.node-time-text {
  font-size: 9px;
  fill: #909399;
}

.task-arrow-group {
  cursor: pointer;
  transition: opacity 0.2s;
}

.task-arrow-group:hover {
  opacity: 0.8;
}

.task-arrow-group.is-selected .task-arrow {
  stroke-width: 4 !important;
}

.task-arrow {
  transition: stroke-width 0.15s ease;
}

.task-arrow:hover {
  stroke-width: 3.5 !important;
}

.task-arrow.is-critical {
  stroke: #FF8A65;
  stroke-width: 2.5;
}

.task-label {
  font-size: 12px;
  fill: #606266;
  pointer-events: none;
}

.task-duration {
  font-size: 10px;
  fill: #909399;
  pointer-events: none;
}

.dependency-line {
  cursor: pointer;
  pointer-events: none;
}

.is-critical .node-circle {
  stroke: #FF8A65;
}
</style>
