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

          <!-- 依赖关系线（虚工作） -->
          <g class="dummy-work-lines">
            <path
              v-for="dep in renderableDependencies"
              :key="dep.key"
              :d="dep.path"
              stroke="#999"
              stroke-width="1.5"
              stroke-dasharray="6,4"
              fill="none"
              :marker-end="`url(#arrowhead-dependency)`"
              class="dummy-work-line"
            />
          </g>

          <!-- 任务箭头（实工作） -->
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
              <!-- 实工作线（实线部分） -->
              <path
                :d="task.realPath"
                :stroke="task.stroke"
                :stroke-width="task.strokeWidth"
                fill="none"
                :marker-end="`url(#${task.marker})`"
                class="task-arrow-real"
                @contextmenu="handleTaskContextMenu($event, task)"
              />

              <!-- 自由时差（波形线部分） -->
              <path
                v-if="task.slackPath"
                :d="task.slackPath"
                :stroke="task.isCritical ? '#FF8A65' : '#64B5F6'"
                stroke-width="1.5"
                fill="none"
                :marker-end="`url(#${task.marker})`"
                class="task-arrow-slack"
                stroke-dasharray="3,3"
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
                v-if="showDuration && !task.isCritical"
                :x="task.durationX"
                :y="task.durationY"
                text-anchor="middle"
                class="task-duration"
              >
                {{ task.duration }}天
              </text>

              <!-- 时差信息 -->
              <text
                v-if="showSlack && task.slack > 0"
                :x="task.slackLabelX"
                :y="task.slackLabelY"
                text-anchor="middle"
                class="task-slack-label"
                fill="#67C23A"
              >
                FF:{{ task.slack }}
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
              @contextmenu="handleNodeContextMenu($event, node)"
            >
              <!-- 节点圆形背景 -->
              <circle
                :r="nodeRadius"
                fill="#FFFFFF"
                :stroke="node.isCritical ? '#FF8A65' : '#1976D2'"
                stroke-width="2"
                filter="url(#nodeShadow)"
                class="node-circle"
              />

              <!-- 节点内分割线（水平） -->
              <line
                x1="-nodeRadius + 2"
                y1="0"
                :x2="nodeRadius - 2"
                y2="0"
                stroke="#ccc"
                stroke-width="1"
              />

              <!-- 节点内分割线（垂直） -->
              <line
                x1="0"
                y1="-nodeRadius + 2"
                :x2="0"
                :y2="nodeRadius - 2"
                stroke="#ccc"
                stroke-width="1"
              />

              <!-- 节点编号（中心） -->
              <text
                y="4"
                text-anchor="middle"
                class="node-number"
                fill="#333"
                font-weight="bold"
                font-size="11"
              >
                {{ node.number }}
              </text>

              <!-- 左上：最早开始 ES -->
              <text
                v-if="showTimeParams"
                :x="-nodeRadius/2 - 2"
                :y="-nodeRadius/2 - 2"
                text-anchor="end"
                class="node-time-mini"
                fill="#67C23A"
                font-size="8"
              >
                {{ node.ES !== undefined ? node.ES : '' }}
              </text>

              <!-- 右上：最晚开始 LS -->
              <text
                v-if="showTimeParams"
                :x="nodeRadius/2 + 2"
                :y="-nodeRadius/2 - 2"
                text-anchor="start"
                class="node-time-mini"
                fill="#F56C6C"
                font-size="8"
              >
                {{ node.LS !== undefined ? node.LS : '' }}
              </text>

              <!-- 左下：最早结束 EF -->
              <text
                v-if="showTimeParams"
                :x="-nodeRadius/2 - 2"
                :y="nodeRadius/2 + 8"
                text-anchor="end"
                class="node-time-mini"
                fill="#67C23A"
                font-size="8"
              >
                {{ node.EF !== undefined ? node.EF : '' }}
              </text>

              <!-- 右下：最晚结束 LF -->
              <text
                v-if="showTimeParams"
                :x="nodeRadius/2 + 2"
                :y="nodeRadius/2 + 8"
                text-anchor="start"
                class="node-time-mini"
                fill="#F56C6C"
                font-size="8"
              >
                {{ node.LF !== undefined ? node.LF : '' }}
              </text>
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
  'zoom-change',
  'pan-change',
  'node-contextmenu',
  'task-contextmenu'
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

// 监听平移变化，通知父组件
watch(panX, (newX) => {
  emit('pan-change', { x: newX, y: panY.value })
})

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
      // 对齐模式：Y坐标基于任务列表
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

    // 计算实工作路径（从开始到结束）
    const realPath = `M ${startX} ${startY} L ${endX} ${endY}`

    // 计算总时差
    const duration = Math.ceil((taskEnd - taskStart) / (1000 * 60 * 60 * 24))
    const slack = task.slack || 0

    // 如果有时差，绘制时差线（波形线）
    let slackPath = null
    let slackLabelX = 0
    let slackLabelY = 0

    if (slack > 0 && !task.is_critical) {
      // 时差线：从任务结束到结束+时差
      const slackEndX = endX + slack * props.dayWidth
      // 使用波浪线路径
      slackPath = calculateWavePath(endX, endY, slackEndX, endY)
      slackLabelX = endX + (slack * props.dayWidth) / 2
      slackLabelY = startY - 12
    }

    return {
      id: task.id,
      realPath, // 实工作线
      slackPath, // 时差线（波形）
      path: realPath, // 保持向后兼容
      label: task.name,
      stroke: task.is_critical ? '#FF8A65' : '#64B5F6',
      strokeWidth: task.is_critical ? 3 : 2,
      marker: task.is_critical ? 'arrowhead-critical' : 'arrowhead-task',
      textColor: '#606266',
      isCritical: task.is_critical,
      duration,
      slack,
      labelX: (startX + endX) / 2,
      labelY: startY - 8,
      durationX: (startX + endX) / 2,
      durationY: startY + 18,
      slackLabelX,
      slackLabelY
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

// 计算波形线路径（用于表示自由时差）
function calculateWavePath(startX, startY, endX, endY) {
  const distance = endX - startX
  const amplitude = 4 // 波浪振幅
  const frequency = 8 // 波浪频率

  if (distance <= frequency) {
    // 距离太短，使用直线
    return `M ${startX} ${startY} L ${endX} ${endY}`
  }

  let path = `M ${startX} ${startY}`
  const numWaves = Math.floor(distance / frequency)

  for (let i = 0; i < numWaves; i++) {
    const x1 = startX + i * frequency
    const x2 = startX + (i + 0.5) * frequency
    const x3 = startX + (i + 1) * frequency

    if (i === 0) {
      path += ` Q ${x1 + frequency / 4} ${startY - amplitude}, ${x2} ${startY}`
      path += ` Q ${x2 + frequency / 4} ${startY + amplitude}, ${x3} ${startY}`
    } else {
      path += ` Q ${x1 + frequency / 4} ${startY + amplitude}, ${x2} ${startY}`
      path += ` Q ${x2 + frequency / 4} ${startY - amplitude}, ${x3} ${startY}`
    }
  }

  // 处理剩余部分
  const lastX = startX + numWaves * frequency
  if (lastX < endX) {
    path += ` L ${endX} ${endY}`
  }

  return path
}

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

// 格式化时间显示（转换为天数）
function formatTime(timestamp) {
  if (!timestamp) return '-'
  // 如果是Unix时间戳（秒），转换为天数
  if (typeof timestamp === 'number' && timestamp > 100000) {
    const days = Math.floor(timestamp / (24 * 60 * 60))
    return days.toString()
  }
  return timestamp.toString()
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

// 节点右键菜单
function handleNodeContextMenu(event, node) {
  event.preventDefault()
  event.stopPropagation()
  selectedNodeId.value = node.id
  selectedTaskId.value = null

  // 触发右键菜单事件，传递节点和鼠标位置
  emit('node-contextmenu', {
    event,
    node,
    x: event.clientX,
    y: event.clientY
  })
}

// 任务点击
function handleTaskClick(event, task) {
  event.stopPropagation()
  selectedTaskId.value = task.id
  selectedNodeId.value = null
  emit('task-click', task)
}

// 任务右键菜单
function handleTaskContextMenu(event, task) {
  event.preventDefault()
  event.stopPropagation()
  selectedTaskId.value = task.id
  selectedNodeId.value = null

  // 触发右键菜单事件，传递任务和鼠标位置
  emit('task-contextmenu', {
    event,
    task,
    x: event.clientX,
    y: event.clientY
  })
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

.node-time-mini {
  font-family: monospace;
  font-weight: 600;
}

.task-arrow-group {
  cursor: pointer;
  transition: opacity 0.2s;
}

.task-arrow-group:hover {
  opacity: 0.8;
}

.task-arrow-group.is-selected .task-arrow-real {
  stroke-width: 4 !important;
}

.task-arrow-real {
  transition: stroke-width 0.15s ease;
}

.task-arrow-real:hover {
  stroke-width: 3 !important;
}

.task-arrow-slack {
  opacity: 0.7;
}

.task-arrow.is-critical {
  stroke: #FF8A65;
  stroke-width: 2.5;
}

.task-label {
  font-size: 12px;
  fill: #606266;
  pointer-events: none;
  font-weight: 500;
}

.task-duration {
  font-size: 10px;
  fill: #909399;
  pointer-events: none;
}

.task-slack-label {
  font-size: 9px;
  font-weight: bold;
  pointer-events: none;
}

.dummy-work-line {
  cursor: pointer;
  pointer-events: none;
}

.is-critical .node-circle {
  stroke: #FF8A65;
}
</style>
