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
          <!-- 箭头标记 -->
          <marker
            id="arrowhead-normal"
            markerWidth="10"
            markerHeight="7"
            refX="9"
            refY="3.5"
            orient="auto"
            markerUnits="userSpaceOnUse"
          >
            <path d="M0,0 L0,7 L9,3.5 z" fill="#64B5F6" />
          </marker>
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
          <!-- 节点阴影 -->
          <filter id="nodeShadow" x="-50%" y="-50%" width="200%" height="200%">
            <feDropShadow dx="0" dy="2" stdDeviation="3" flood-opacity="0.2"/>
          </filter>
        </defs>

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

        <!-- 依赖关系线 -->
        <g class="dependency-lines">
          <path
            v-for="edge in renderableEdges"
            :key="edge.key"
            :d="edge.path"
            :stroke="edge.stroke"
            :stroke-width="edge.strokeWidth"
            fill="none"
            :marker-end="`url(#${edge.marker})`"
            :class="{ 'is-critical': edge.isCritical }"
            class="dependency-line"
            @click="handleEdgeClick(edge)"
          />
        </g>

        <!-- 节点 -->
        <g class="nodes">
          <g
            v-for="node in nodes"
            :key="node.id"
            :transform="`translate(${node.x}, ${node.y})`"
            :class="{ 'is-selected': selectedNode?.id === node.id, 'is-critical': node.isCritical }"
            class="node-group"
            @mousedown="handleNodeMouseDown($event, node)"
            @click="handleNodeClick($event, node)"
            @dblclick="handleNodeDblClick($event, node)"
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

            <!-- 节点标签 -->
            <text
              v-if="showTaskNames"
              :y="nodeRadius + 16"
              text-anchor="middle"
              class="node-label"
              :fill="node.textColor"
            >
              {{ node.label }}
            </text>

            <!-- 时间参数 -->
            <g v-if="showTimeParams" class="node-time-params">
              <text :y="nodeRadius + 28" text-anchor="middle" class="node-time-text">
                ES: {{ node.timeParams?.ES || '-' }}
              </text>
              <text :y="nodeRadius + 40" text-anchor="middle" class="node-time-text">
                EF: {{ node.timeParams?.EF || '-' }}
              </text>
            </g>

            <!-- 时差信息 -->
            <text
              v-if="showSlack && node.slack !== undefined"
              :y="nodeRadius + 52"
              text-anchor="middle"
              class="node-slack-text"
              :fill="node.slackColor"
            >
              时差: {{ node.slack }}天
            </text>
          </g>
        </g>

        <!-- 连线时的临时线 -->
        <line
          v-if="isConnecting && connectStartNode && connectEndPosition"
          :x1="connectStartNode.x"
          :y1="connectStartNode.y"
          :x2="connectEndPosition.x"
          :y2="connectEndPosition.y"
          stroke="#64B5F6"
          stroke-width="2"
          stroke-dasharray="5,5"
        />
      </svg>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { findOrthogonalPath } from '@/utils/ganttHelpers'

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
  'node-dblclick',
  'edge-click',
  'task-updated'
])

// Refs
const containerRef = ref(null)
const canvasContainerRef = ref(null)
const svgRef = ref(null)

// 状态
const showGrid = ref(true)
const gridSize = ref(20)
const zoomLevel = ref(1)
const panX = ref(0)
const panY = ref(0)

// 拖拽状态
const isDragging = ref(false)
const isPanning = ref(false)
const draggedNode = ref(null)
const dragStartPos = ref({ x: 0, y: 0 })

// 连线状态
const isConnecting = ref(false)
const connectStartNode = ref(null)
const connectEndPosition = ref({ x: 0, y: 0 })

// 选中状态
const selectedNode = ref(null)
const selectedEdge = ref(null)

// 节点数据
const nodes = computed(() => {
  return props.tasks.map(task => ({
    id: task.id,
    x: task.position_x || 100 + (task.id % 10) * 150,
    y: task.position_y || 100 + Math.floor(task.id / 10) * 100,
    label: task.name,
    fill: getNodeFill(task),
    stroke: getNodeStroke(task),
    textColor: '#606266',
    isCritical: task.is_critical || false,
    slack: task.slack,
    slackColor: task.slack > 0 ? '#67C23A' : '#F56C6C',
    timeParams: {
      ES: task.early_start,
      EF: task.early_finish,
      LS: task.late_start,
      LF: task.late_finish
    },
    taskId: task.id,
    originalX: task.position_x,
    originalY: task.position_y
  }))
})

// 边数据
const edges = computed(() => {
  return props.dependencies.map((dep, index) => {
    const fromNode = nodes.value.find(n => n.taskId === dep.depends_on)
    const toNode = nodes.value.find(n => n.taskId === dep.task_id)

    if (!fromNode || !toNode) return null

    return {
      key: `${dep.depends_on}-${dep.task_id}-${index}`,
      from: fromNode,
      to: toNode,
      fromId: dep.depends_on,
      toId: dep.task_id,
      type: dep.type || 'FS',
      isCritical: dep.is_critical || false,
      taskId: dep.task_id
    }
  }).filter(Boolean)
})

// 可渲染的边（包含路径）
const renderableEdges = computed(() => {
  return edges.value.map(edge => {
    const path = calculateEdgePath(edge.from, edge.to)
    return {
      ...edge,
      path,
      stroke: edge.isCritical ? '#FF8A65' : '#64B5F6',
      strokeWidth: edge.isCritical ? 2 : 1.5,
      marker: edge.isCritical ? 'arrowhead-critical' : 'arrowhead-normal'
    }
  })
})

// 获取节点填充色
function getNodeFill(task) {
  if (task.is_milestone) return '#FFF9C4'
  if (task.status === 'completed') return '#C8E6C9'
  if (task.status === 'in_progress') return '#BBDEFB'
  if (task.status === 'delayed') return '#FFCDD2'
  return '#E3F2FD'
}

// 获取节点边框色
function getNodeStroke(task) {
  if (task.is_critical) return '#FF8A65'
  if (task.is_milestone) return '#FBC02D'
  return '#1976D2'
}

// 计算边路径
function calculateEdgePath(fromNode, toNode) {
  const startX = fromNode.x + props.nodeRadius
  const startY = fromNode.y
  const endX = toNode.x - props.nodeRadius - 1
  const endY = toNode.y

  // 使用 A* 路径规划算法
  const path = findOrthogonalPath(startX, startY, endX, endY, [])
  return path || `M ${startX} ${startY} L ${endX} ${endY}`
}

// 鼠标滚轮缩放
function handleWheel(event) {
  const delta = event.deltaY > 0 ? -0.1 : 0.1
  zoomLevel.value = Math.max(0.1, Math.min(3, zoomLevel.value + delta))
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

  if (isDragging.value && draggedNode.value) {
    const rect = svgRef.value.getBoundingClientRect()
    draggedNode.value.x = (event.clientX - rect.left - dragStartPos.value.x) / zoomLevel.value
    draggedNode.value.y = (event.clientY - rect.top - dragStartPos.value.y) / zoomLevel.value
  }

  if (isConnecting.value) {
    const rect = svgRef.value.getBoundingClientRect()
    connectEndPosition.value = {
      x: (event.clientX - rect.left) / zoomLevel.value,
      y: (event.clientY - rect.top) / zoomLevel.value
    }
  }
}

// 鼠标释放
function handleMouseUp() {
  if (isDragging.value && draggedNode.value) {
    // 发送节点位置更新
    emit('task-updated', {
      taskId: draggedNode.value.taskId,
      position_x: draggedNode.value.x,
      position_y: draggedNode.value.y
    })
  }

  isDragging.value = false
  isPanning.value = false
  draggedNode.value = null
}

// 节点点击
function handleNodeClick(event, node) {
  event.stopPropagation()
  selectedNode.value = node
  selectedEdge.value = null
  emit('node-click', node)
}

// 节点双击
function handleNodeDblClick(event, node) {
  event.stopPropagation()
  emit('node-dblclick', node)
}

// 边点击
function handleEdgeClick(edge) {
  selectedEdge.value = edge
  selectedNode.value = null
  emit('edge-click', edge)
}

// 键盘事件
function handleKeydown(event) {
  if (event.key === 'Escape') {
    isConnecting.value = false
    connectStartNode.value = null
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

.node-label {
  font-size: 12px;
  fill: #606266;
  pointer-events: none;
}

.node-time-params {
  font-size: 10px;
  fill: #909399;
  pointer-events: none;
}

.node-time-text {
  font-size: 10px;
  fill: #909399;
}

.node-slack-text {
  font-size: 10px;
  font-weight: bold;
}

.dependency-line {
  cursor: pointer;
  transition: stroke-width 0.15s ease;
}

.dependency-line:hover {
  stroke-width: 3 !important;
}

.dependency-line.is-critical {
  stroke: #FF8A65;
  stroke-width: 2;
}

.is-critical .node-circle {
  stroke: #FF8A65;
}
</style>
