<template>
  <div class="network-view" ref="containerRef">
    <!-- SVG 绘图区域 -->
    <div class="network-canvas-container" ref="canvasContainerRef">
      <svg
        ref="svgRef"
        class="network-svg"
        :style="{ cursor: toolMode === 'pan' ? (isPanning ? 'grabbing' : 'grab') : 'default' }"
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
            <g
              v-for="dep in renderableDependencies"
              :key="dep.key"
              class="dummy-work-group"
            >
              <!-- R17: 虚工作必须用虚线 -->
              <path
                :d="dep.path"
                stroke="#999"
                stroke-width="1.5"
                stroke-dasharray="6,4"
                fill="none"
                :marker-end="`url(#arrowhead-dependency)`"
                class="dummy-work-line"
              />
              <!-- R18: 虚工作必须标注为"虚" -->
              <text
                v-if="dep.isDummy"
                :x="dep.labelX"
                :y="dep.labelY - 5"
                text-anchor="middle"
                class="dummy-work-label"
                fill="#999"
                font-size="10"
                font-weight="bold"
              >
                虚
              </text>
            </g>
          </g>

          <!-- 临时连线（创建依赖关系时） -->
          <g v-if="isCreatingDependency && tempLinePath" class="temp-dependency-line">
            <path
              :d="tempLinePath"
              stroke="#409EFF"
              stroke-width="2"
              stroke-dasharray="4,2"
              fill="none"
              :marker-end="`url(#arrowhead-dependency)`"
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

          <!-- 线条交叉点标记 -->
          <g class="arrow-intersections">
            <circle
              v-for="(point, index) in arrowIntersectionPoints"
              :key="`intersection-${index}`"
              :cx="point.x"
              :cy="point.y"
              r="3"
              fill="#FF5722"
              stroke="#fff"
              stroke-width="1"
              class="intersection-point"
            />
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
                :x1="-nodeRadius + 2"
                y1="0"
                :x2="nodeRadius - 2"
                y2="0"
                stroke="#ccc"
                stroke-width="1"
              />

              <!-- 节点内分割线（垂直） -->
              <line
                x1="0"
                :y1="-nodeRadius + 2"
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

    <!-- 状态栏：显示选中节点的详细信息 -->
    <div class="network-status-bar" v-if="selectedNode && showTimeParams">
      <div class="status-info">
        <div class="status-item">
          <span class="status-label">节点编号:</span>
          <span class="status-value">{{ selectedNode.number }}</span>
        </div>
        <div class="status-divider"></div>
        <div class="status-item">
          <span class="status-label">日期:</span>
          <span class="status-value">{{ selectedNode.date }}</span>
        </div>
        <div class="status-divider"></div>
        <div class="status-item">
          <span class="status-label">ES (最早开始):</span>
          <span class="status-value status-es">{{ selectedNode.ES }}</span>
        </div>
        <div class="status-divider"></div>
        <div class="status-item">
          <span class="status-label">EF (最早结束):</span>
          <span class="status-value status-ef">{{ selectedNode.EF }}</span>
        </div>
        <div class="status-divider"></div>
        <div class="status-item">
          <span class="status-label">LS (最迟开始):</span>
          <span class="status-value status-ls">{{ selectedNode.LS }}</span>
        </div>
        <div class="status-divider"></div>
        <div class="status-item">
          <span class="status-label">LF (最迟结束):</span>
          <span class="status-value status-lf">{{ selectedNode.LF }}</span>
        </div>
        <div class="status-divider"></div>
        <div class="status-item">
          <span class="status-label">总时差:</span>
          <span class="status-value status-total-slack">
            {{ selectedNode.totalSlack !== undefined ? selectedNode.totalSlack : '-' }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { isDummyTask } from '@/utils/ganttHelpers'
import { formatDate } from '@/utils/dateFormat'

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
  'task-contextmenu',
  'node-time-change', // 节点拖动改变任务时间
  'node-dependency-create' // 节点间创建依赖关系
])

// 鼠标状态追踪
const mouseState = ref({
  isCtrlPressed: false
})

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
const dragStartPos = ref({ x: 0, y: 0 })

// 节点位置偏移（用于拖动）
const nodeOffsets = ref(new Map())

// 拖动开始时的节点位置
const dragStartNodePos = ref({
  nodeId: null,
  offsetX: 0,
  offsetY: 0
})

// 创建依赖关系模式
const isCreatingDependency = ref(false)
const dependencySourceNode = ref(null)
const tempLineEnd = ref(null)

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

// 创建事件节点和任务箭头（按照双代号时标网络图标准规范）
const eventNodes = computed(() => {
  if (!realTasks.value.length || !timelineStartDate.value) return []

  // ==================== 第一步：收集所有任务信息 ====================
  const taskPredecessors = new Map()
  const taskSuccessors = new Map()
  realTasks.value.forEach(task => {
    taskPredecessors.set(task.id, task.predecessors || [])
    taskSuccessors.set(task.id, task.successors || [])
  })

  // ==================== 第二步：为每个任务创建独立的节点 ====================
  // 为了确保每个任务都能在任务列表对应行上显示节点
  // 每个任务有独立的开始和结束节点

  const nodeStates = new Map()  // key: task_id + '_start' 或 task_id + '_end', value: { date, x, taskId, type }

  // 为每个任务的开始创建节点
  realTasks.value.forEach(task => {
    const taskStart = new Date(task.start)
    const startDateKey = formatDate(taskStart)
    const startDays = Math.ceil((taskStart - timelineStartDate.value) / (1000 * 60 * 60 * 24))

    const startSignature = `${task.id}_start`
    nodeStates.set(startSignature, {
      date: startDateKey,
      x: startDays * props.dayWidth,
      taskId: task.id,
      type: 'start'
    })
  })

  // 为每个任务的结束创建节点
  realTasks.value.forEach(task => {
    const taskEnd = new Date(task.end)
    const endDateKey = formatDate(taskEnd)
    const endDays = Math.ceil((taskEnd - timelineStartDate.value) / (1000 * 60 * 60 * 24))

    const endSignature = `${task.id}_end`
    nodeStates.set(endSignature, {
      date: endDateKey,
      x: endDays * props.dayWidth,
      taskId: task.id,
      type: 'end'
    })
  })

  // ==================== 第三步：合并相同日期和位置的节点 ====================
  // 如果多个任务的开始/结束节点在相同日期且没有前置/后置依赖差异，可以合并
  // 但为了确保可视化对齐，我们保持每个任务有独立的节点显示

  // ==================== 第三步：合并无前置/无后续任务的节点 ====================

  // 合并所有没有前置任务的开始节点（可选，为了保持图的简洁性）
  // 注意：这会违反严格的AOA规范，但能提供更好的可视化效果
  const noPredecessorTasks = realTasks.value.filter(t => !t.predecessors || t.predecessors.length === 0)
  if (noPredecessorTasks.length > 1 && false) { // 设为false以禁用合并，保持每个任务独立
    // 如果需要合并，可以在这里实现
  }

  // 合并所有没有后续任务的结束节点（可选）
  const noSuccessorTasks = realTasks.value.filter(t => !t.successors || t.successors.length === 0)
  if (noSuccessorTasks.length > 1 && false) { // 设为false以禁用合并，保持每个任务独立
    // 如果需要合并，可以在这里实现
  }

  // ==================== 第四步：按拓扑顺序编号 ====================
  // R4: 起点编号 < 终点编号
  // 使用拓扑排序确保节点编号符合依赖关系

  const statesArray = Array.from(nodeStates.values())
  const numberedNodes = []
  let nodeNumber = 1

  // 计算节点的拓扑顺序
  // 规则：如果任务B依赖任务A，则：
  //   - 任务A的开始节点 < 任务A的结束节点 < 任务B的开始节点 < 任务B的结束节点

  // 构建依赖图
  const taskDependencies = new Map() // taskId -> 依赖的任务ID列表
  realTasks.value.forEach(task => {
    taskDependencies.set(task.id, task.predecessors || [])
  })

  // 拓扑排序：按依赖关系对任务排序
  const sortedTasks = []
  const visited = new Set()
  const visiting = new Set()

  function dfsVisit(taskId) {
    if (visited.has(taskId)) return
    if (visiting.has(taskId)) {
      // 检测到循环依赖，跳过
      console.warn(`检测到循环依赖，涉及任务 ${taskId}`)
      return
    }

    visiting.add(taskId)

    // 先访问所有前置任务
    const deps = taskDependencies.get(taskId) || []
    deps.forEach(depId => {
      if (realTasks.value.find(t => t.id === depId)) {
        dfsVisit(depId)
      }
    })

    visiting.delete(taskId)
    visited.add(taskId)
    sortedTasks.push(taskId)
  }

  // 对所有任务执行DFS
  realTasks.value.forEach(task => {
    if (!visited.has(task.id)) {
      dfsVisit(task.id)
    }
  })

  // 按拓扑顺序排序节点
  const taskOrder = new Map()
  sortedTasks.forEach((taskId, index) => {
    taskOrder.set(taskId, index)
  })

  statesArray.sort((a, b) => {
    // 先按拓扑顺序排序（确保依赖的任务排在前面）
    const orderA = taskOrder.get(a.taskId) ?? 999999
    const orderB = taskOrder.get(b.taskId) ?? 999999
    if (orderA !== orderB) {
      return orderA - orderB
    }
    // 同一任务，开始节点优先
    if (a.type === 'start' && b.type === 'end') return -1
    if (a.type === 'end' && b.type === 'start') return 1
    return 0
  })

  // 为每个节点分配编号
  const taskStartNodes = new Map() // taskId -> node
  const taskEndNodes = new Map()   // taskId -> node

  statesArray.forEach(state => {
    const task = realTasks.value.find(t => t.id === state.taskId)
    const node = {
      id: `node-${state.taskId}-${state.type}`,
      number: nodeNumber++,
      x: state.x,
      y: 0,
      date: state.date,
      taskId: state.taskId,
      type: state.type,
      tasks: {
        start: state.type === 'start' ? [state.taskId] : [],
        end: state.type === 'end' ? [state.taskId] : []
      },
      isMerged: false, // 每个任务独立节点，不合并
      ES: state.type === 'start' ? (task ? task.early_start : undefined) : undefined,
      EF: state.type === 'start' ? (task ? task.early_finish : undefined) : undefined,
      LS: state.type === 'end' ? (task ? task.late_start : undefined) : undefined,
      LF: state.type === 'end' ? (task ? task.late_finish : undefined) : undefined
    }
    numberedNodes.push(node)

    // 记录任务的开始和结束节点
    if (state.type === 'start') {
      taskStartNodes.set(state.taskId, node)
    } else {
      taskEndNodes.set(state.taskId, node)
    }
  })

  // ==================== 第五步：验证R4规则 ====================
  // R4: 箭线方向必须满足：起点节点编号 < 终点节点编号
  numberedNodes.forEach(node => {
    if (node.type === 'start') {
      const startNode = node
      const endNode = taskEndNodes.get(node.taskId)

      if (endNode && startNode.number >= endNode.number) {
        console.warn(`❌ R4违规: 任务${node.taskId} 起点(${startNode.number}) >= 终点(${endNode.number})`)
      }
    }
  })

  // 验证依赖关系也满足R4规则
  realTasks.value.forEach(task => {
    const predecessors = task.predecessors || []
    predecessors.forEach(predId => {
      const predEndNode = taskEndNodes.get(predId)
      const taskStartNode = taskStartNodes.get(task.id)

      if (predEndNode && taskStartNode) {
        if (predEndNode.number >= taskStartNode.number) {
          console.warn(`❌ R4违规: 依赖关系 任务${predId} -> 任务${task.id} 终点(${predEndNode.number}) >= 起点(${taskStartNode.number})`)
        }
      }
    })
  })

  // ==================== 第六步：Y坐标调整 ====================

  numberedNodes.forEach((node) => {
    // 判断是否为关键节点
    const task = realTasks.value.find(t => t.id === node.taskId)
    node.isCritical = task && task.is_critical

    // 计算Y坐标，确保与任务列表对齐
    const offset = nodeOffsets.value.get(node.id) || { x: 0, y: 0 }

    // 使用taskIndexMap对齐任务列表
    if (props.taskIndexMap[node.taskId] !== undefined) {
      const baseTaskIndex = props.taskIndexMap[node.taskId]
      node.baseY = baseTaskIndex * props.rowHeight + props.rowHeight / 2
      node.y = node.baseY + offset.y
    } else {
      // 如果taskIndexMap中没有该任务，从realTasks中查找
      const taskIndex = realTasks.value.findIndex(t => t.id === node.taskId)
      if (taskIndex !== -1) {
        node.baseY = taskIndex * props.rowHeight + props.rowHeight / 2
        node.y = node.baseY + offset.y
      } else {
        // 找不到，使用默认位置
        node.baseY = 100
        node.y = node.baseY + offset.y
      }
    }

    node.baseX = node.x
    node.x = node.x + offset.x
  })

  return numberedNodes
})

// 选中的节点详细信息
const selectedNode = computed(() => {
  if (!selectedNodeId.value) return null
  const node = eventNodes.value.find(n => n.id === selectedNodeId.value)
  if (!node) return null

  // 计算总时差
  // 总时差 = LS - ES 或 LF - EF
  const ES = node.ES
  const EF = node.EF
  const LS = node.LS
  const LF = node.LF

  let totalSlack
  if (ES !== undefined && LS !== undefined) {
    // 将时间戳转换为天数，然后计算差值
    const esDate = new Date(ES * 1000)
    const lsDate = new Date(LS * 1000)
    totalSlack = Math.ceil((lsDate - esDate) / (1000 * 60 * 60 * 24))
  } else if (EF !== undefined && LF !== undefined) {
    const efDate = new Date(EF * 1000)
    const lfDate = new Date(LF * 1000)
    totalSlack = Math.ceil((lfDate - efDate) / (1000 * 60 * 60 * 24))
  }

  return {
    ...node,
    totalSlack
  }
})

// 计算从圆边缘到圆边缘的箭头路径
const calculateArrowPathFromCircle = (fromX, fromY, toX, toY, radius) => {
  // 计算从起点到终点的方向向量
  const dx = toX - fromX
  const dy = toY - fromY
  const distance = Math.sqrt(dx * dx + dy * dy)

  // 如果起点和终点重合，返回空路径
  if (distance < 0.1) {
    return ''
  }

  // 计算单位方向向量
  const ux = dx / distance
  const uy = dy / distance

  // 计算箭头起点和终点（在圆边缘上）
  const arrowStartX = fromX + ux * radius
  const arrowStartY = fromY + uy * radius
  const arrowEndX = toX - ux * radius
  const arrowEndY = toY - uy * radius

  // 计算水平和垂直距离
  const horizontalDist = Math.abs(arrowEndX - arrowStartX)
  const verticalDist = Math.abs(arrowEndY - arrowStartY)

  // 如果垂直距离很小（小于10像素），使用直线
  if (verticalDist < 10) {
    return `M ${arrowStartX} ${arrowStartY} L ${arrowEndX} ${arrowEndY}`
  }

  // 如果水平距离很小（小于10像素），使用直线
  if (horizontalDist < 10) {
    return `M ${arrowStartX} ${arrowStartY} L ${arrowEndX} ${arrowEndY}`
  }

  // 根据垂直距离选择连接方式
  // 垂直距离较大时使用斜45度，较小时使用直角
  if (verticalDist > 30) {
    // 使用斜45度连接
    return calculateDiagonalPath(arrowStartX, arrowStartY, arrowEndX, arrowEndY)
  } else {
    // 使用直角连接
    return calculateOrthogonalPath(arrowStartX, arrowStartY, arrowEndX, arrowEndY)
  }
}

// 计算带斜45度过渡的路径（水平 -> 斜45度 -> 水平）
// 斜45度只用于中间过渡段，不能直接斜连接
function calculateDiagonalPath(startX, startY, endX, endY) {
  const dx = endX - startX
  const dy = endY - startY
  const absDx = Math.abs(dx)
  const absDy = Math.abs(dy)

  // 确保有足够的空间进行三段式过渡
  const minTransitionLength = 30 // 最小过渡长度
  if (absDx < minTransitionLength * 2) {
    // 空间不足，使用直角
    return calculateOrthogonalPath(startX, startY, endX, endY)
  }

  // 三段式路径：水平 -> 斜45度 -> 水平
  const transitionLength = Math.min(minTransitionLength, absDx / 3)

  // 第一段：水平线（从起点）
  const firstEndX = startX + transitionLength
  const firstEndY = startY

  // 第二段：斜45度线（垂直偏移）
  const diagonalLength = absDy
  const secondEndX = firstEndX + diagonalLength
  const secondEndY = endY

  // 第三段：水平线（到终点）
  const path = `M ${startX} ${startY} L ${firstEndX} ${firstEndY} L ${secondEndX} ${secondEndY} L ${endX} ${endY}`

  return path
}

// 任务箭头（连接起点和终点节点）
const taskArrows = computed(() => {
  if (!realTasks.value.length || !timelineStartDate.value) return []

  return realTasks.value.map(task => {
    // 直接通过节点ID查找该任务的开始和结束节点
    const startNodeId = `node-${task.id}-start`
    const endNodeId = `node-${task.id}-end`

    const startNode = eventNodes.value.find(node => node.id === startNodeId)
    const endNode = eventNodes.value.find(node => node.id === endNodeId)

    // 计算位置
    let startX, startY, endX, endY

    // 使用节点的实际位置（包含拖动偏移）
    if (startNode) {
      startX = startNode.x
      startY = startNode.y
    } else {
      // 如果找不到节点，使用时间计算位置
      const taskStart = new Date(task.start)
      const startDays = Math.ceil((taskStart - timelineStartDate.value) / (1000 * 60 * 60 * 24))
      startX = startDays * props.dayWidth
      startY = 100
    }

    if (endNode) {
      endX = endNode.x
      endY = endNode.y
    } else {
      // 如果找不到节点，使用时间计算位置
      const taskEnd = new Date(task.end)
      const endDays = Math.ceil((taskEnd - timelineStartDate.value) / (1000 * 60 * 60 * 24))
      endX = endDays * props.dayWidth
      endY = 100
    }

    // R19: 验证箭线只能向右（x2 ≥ x1）
    if (startX > endX) {
      console.warn(`❌ R19违规: 任务"${task.name}" 箭线向左 startX(${startX.toFixed(0)}) > endX(${endX.toFixed(0)})`)
      // 强制修正：确保箭线向右
      endX = Math.max(startX, endX)
    }

    // 计算实工作路径（从开始到结束）- 使用圆边缘到圆边缘的路径
    const realPath = calculateArrowPathFromCircle(startX, startY, endX, endY, props.nodeRadius)

    // 计算总时差
    const duration = Math.ceil((taskEnd - taskStart) / (1000 * 60 * 60 * 24))
    const slack = task.slack || 0

    // 如果有时差，绘制时差线（波形线）
    let slackPath = null
    let slackLabelX = 0
    let slackLabelY = 0

    if (slack > 0 && !task.is_critical) {
      // 时差线：从圆边缘开始，到圆边缘结束
      const slackEndX = endX + slack * props.dayWidth
      // 计算时差线的起点（从圆边缘开始）
      const dx = slackEndX - endX
      const distance = Math.abs(dx)
      if (distance > 0) {
        const ux = dx / distance
        const arrowStartX = endX + ux * props.nodeRadius
        const arrowStartY = endY
        const arrowEndX = slackEndX
        const arrowEndY = endY
        slackPath = calculateWavePath(arrowStartX, arrowStartY, arrowEndX, arrowEndY)
        slackLabelX = arrowStartX + (arrowEndX - arrowStartX) / 2
        slackLabelY = arrowStartY - 12
      }
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

// 检测箭头线条的交叉点
const arrowIntersectionPoints = computed(() => {
  if (!taskArrows.value.length) return []

  const intersections = []

  // 解析每个任务的路径，获取线段
  const taskSegments = taskArrows.value.map(task => {
    const pathData = task.realPath
    const segments = []

    // 解析SVG路径命令
    const commands = pathData.match(/[MLH][^,]*/g)
    let currentX = 0
    let currentY = 0

    for (let i = 0; i < commands.length; i++) {
      const cmd = commands[i][0]
      const coords = commands[i].slice(1).split(/[, ]+/).filter(s => s).map(Number)

      if (cmd === 'M') {
        currentX = coords[0]
        currentY = coords[1]
      } else if (cmd === 'L') {
        const startX = currentX
        const startY = currentY
        currentX = coords[0]
        currentY = coords[1]
        segments.push({
          x1: startX,
          y1: startY,
          x2: currentX,
          y2: currentY,
          taskId: task.id
        })
      }
    }

    return segments
  })

  // 检测所有线段对之间的交叉
  for (let i = 0; i < taskSegments.length; i++) {
    for (let j = i + 1; j < taskSegments.length; j++) {
      const seg1 = taskSegments[i]
      const seg2 = taskSegments[j]

      // 跳过来自同一个任务的线段
      if (seg1.taskId === seg2.taskId) continue

      // 检测两条线段是否交叉
      const intersection = getLineIntersection(
        seg1.x1, seg1.y1,
        seg1.x2, seg1.y2,
        seg2.x1, seg2.y1,
        seg2.x2, seg2.y2
      )

      if (intersection) {
        intersections.push({
          x: intersection.x,
          y: intersection.y,
          task1Id: seg1.taskId,
          task2Id: seg2.taskId
        })
      }
    }
  }

  // 去重：距离很近的交叉点只保留一个
  const uniqueIntersections = []
  const mergedThreshold = 10 // 合并阈值（像素）

  intersections.forEach(point => {
    const isDuplicate = uniqueIntersections.some(existing => {
      const dist = Math.sqrt(
        Math.pow(existing.x - point.x, 2) +
        Math.pow(existing.y - point.y, 2)
      )
      return dist < mergedThreshold
    })

    if (!isDuplicate) {
      uniqueIntersections.push(point)
    }
  })

  return uniqueIntersections
})

// 计算两条线段的交点
function getLineIntersection(x1, y1, x2, y2, x3, y3, x4, y4) {
  const denom = (x1 - x2) * (y3 - y4) - (y1 - y2) * (x3 - x4)

  // 如果线段平行或重合，无交点
  if (Math.abs(denom) < 0.001) {
    return null
  }

  const t = ((x1 - x3) * (y3 - y4) - (y1 - y3) * (x3 - x4)) / denom
  const u = -((x1 - x2) * (y1 - y3) - (y1 - y2) * (x3 - x4)) / denom

  // 检查交点是否在两条线段上
  if (t >= 0 && t <= 1 && u >= 0 && u <= 1) {
    return {
      x: x1 + t * (x2 - x1),
      y: y1 + t * (y2 - y1)
    }
  }

  return null
}

// 依赖关系（虚工作）- 按照双代号时标网络图规范
const renderableDependencies = computed(() => {
  if (!props.dependencies.length) return []

  return props.dependencies.map((dep, index) => {
    const fromTask = props.tasks.find(t => t.id === dep.depends_on)
    const toTask = props.tasks.find(t => t.id === dep.task_id)

    if (!fromTask || !toTask) return null

    // 直接通过节点ID查找
    const fromNodeId = `node-${fromTask.id}-end`
    const toNodeId = `node-${toTask.id}-start`

    const fromNode = eventNodes.value.find(node => node.id === fromNodeId)
    const toNode = eventNodes.value.find(node => node.id === toNodeId)

    if (!fromNode || !toNode) return null

    // 计算从圆边缘到圆边缘的路径
    const fromX = fromNode.x
    const fromY = fromNode.y
    const toX = toNode.x
    const toY = toNode.y

    // 使用虚线路径连接两个节点
    const path = calculateArrowPathFromCircle(fromX, fromY, toX, toY, props.nodeRadius)

    return {
      key: `dep-${dep.depends_on}-${dep.task_id}-${index}`,
      path,
      isDummy: true,
      labelX: (fromX + toX) / 2,
      labelY: (fromY + toY) / 2 - 5
    }
  }).filter(Boolean)
})

// 临时连线路径（创建依赖关系时）
const tempLinePath = computed(() => {
  if (!isCreatingDependency.value || !dependencySourceNode.value || !tempLineEnd.value) {
    return ''
  }

  const sourceNode = dependencySourceNode.value
  const startX = sourceNode.x
  const startY = sourceNode.y
  const endX = tempLineEnd.value.x - panX.value
  const endY = tempLineEnd.value.y - panY.value

  // 使用正交路径
  return calculateOrthogonalPath(startX, startY, endX, endY)
})

// 计算正交路径（只包含水平和垂直线段）
function calculateOrthogonalPath(startX, startY, endX, endY) {
  const dx = Math.abs(endX - startX)
  const dy = Math.abs(endY - startY)

  // 如果起点和终点在同一行或同一列，直接画直线
  if (dx < 5) {
    return `M ${startX} ${startY} L ${endX} ${endY}`
  }
  if (dy < 5) {
    return `M ${startX} ${startY} L ${endX} ${endY}`
  }

  // 正交路径：先水平再垂直（或先垂直再水平）
  // 使用一个转折点，路径形状取决于起点和终点的相对位置
  let midX

  if (endX > startX) {
    // 目标在右侧
    midX = startX + dx * 0.5
  } else {
    // 目标在左侧
    midX = endX + dx * 0.5
  }

  return `M ${startX} ${startY} L ${midX} ${startY} L ${midX} ${endY} L ${endX} ${endY}`
}

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

  // 检测是否按住Ctrl键（创建依赖关系模式）
  if (event.ctrlKey || event.metaKey) {
    isCreatingDependency.value = true
    dependencySourceNode.value = node
    dragStartPos.value = { x: event.clientX, y: event.clientY }
    tempLineEnd.value = { x: event.clientX, y: event.clientY }
    return
  }

  // 普通拖动模式（只能水平拖动）
  isDragging.value = true
  dragStartPos.value = { x: event.clientX, y: event.clientY }

  // 记录拖动开始时的节点原始位置和关联任务
  const offset = nodeOffsets.value.get(node.id) || { x: 0, y: 0 }

  // 确定节点类型和要修改的任务
  // 新节点系统：节点包含 tasks.start 和 tasks.end 数组
  let nodeType = null
  let taskId = null

  if (node.tasks) {
    // 如果节点只关联一个任务的开始，则改变该任务的开始时间
    if (node.tasks.start.length === 1 && node.tasks.end.length === 0) {
      nodeType = 'start'
      taskId = node.tasks.start[0]
    }
    // 如果节点只关联一个任务的结束，则改变该任务的结束时间
    else if (node.tasks.end.length === 1 && node.tasks.start.length === 0) {
      nodeType = 'end'
      taskId = node.tasks.end[0]
    }
    // 如果节点既有关联的开始又有关联的结束（合并节点），默认处理结束
    else if (node.tasks.end.length > 0) {
      nodeType = 'end'
      taskId = node.tasks.end[0]
    }
    // 否则处理开始
    else if (node.tasks.start.length > 0) {
      nodeType = 'start'
      taskId = node.tasks.start[0]
    }
  }

  dragStartNodePos.value = {
    nodeId: node.id,
    offsetX: offset.x,
    offsetY: 0, // 不记录Y偏移
    nodeType,
    taskId,
    originalX: node.x,
    mergedNodeIds: node.mergedNodeIds || [] // 记录被合并的节点ID
  }
}

// 鼠标移动
function handleMouseMove(event) {
  // 更新Ctrl键状态
  mouseState.value.isCtrlPressed = event.ctrlKey || event.metaKey

  if (isPanning.value) {
    panX.value = event.clientX - dragStartPos.value.x
    panY.value = event.clientY - dragStartPos.value.y
  }

  // 创建依赖关系模式：更新临时连线终点
  if (isCreatingDependency.value && tempLineEnd.value) {
    const rect = svgRef.value.getBoundingClientRect()
    tempLineEnd.value = {
      x: event.clientX - rect.left,
      y: event.clientY - rect.top
    }
    return
  }

  // 普通拖动模式（只能水平拖动）
  if (isDragging.value && dragStartNodePos.value.nodeId && props.toolMode === 'select') {
    const dx = event.clientX - dragStartPos.value.x

    // 只应用X偏移（水平方向移动）
    // 拖动过程中只更新视觉位置，不触发时间变更
    const newOffset = {
      x: dragStartNodePos.value.offsetX + dx,
      y: 0  // Y方向不偏移
    }

    // 只更新当前拖动的合并节点的偏移量
    // 合并节点作为一个整体移动，不需要分别更新原始节点
    nodeOffsets.value.set(dragStartNodePos.value.nodeId, newOffset)
  }
}

// 鼠标释放
function handleMouseUp(event) {
  // 创建依赖关系模式：检测是否在另一个节点上释放
  if (isCreatingDependency.value && dependencySourceNode.value) {
    // 查找鼠标位置下的节点
    const rect = svgRef.value.getBoundingClientRect()
    const mouseX = event.clientX - rect.left
    const mouseY = event.clientY - rect.top

    // 查找目标节点
    const targetNode = eventNodes.value.find(node => {
      const dx = node.x - mouseX
      const dy = node.y - mouseY
      return Math.sqrt(dx * dx + dy * dy) < props.nodeRadius + 10
    })

    if (targetNode && targetNode.id !== dependencySourceNode.value.id) {
      // 从源节点的 tasks.end 中获取源任务ID
      const sourceTaskIds = dependencySourceNode.value.tasks?.end || []
      // 从目标节点的 tasks.start 中获取目标任务ID
      const targetTaskIds = targetNode.tasks?.start || []

      if (sourceTaskIds.length === 0) {
        ElMessage.warning('源节点不是结束节点，无法创建依赖关系')
        return
      }
      if (targetTaskIds.length === 0) {
        ElMessage.warning('目标节点不是开始节点，无法创建依赖关系')
        return
      }

      // 触发依赖关系创建事件（传递任务ID而不是节点ID）
      emit('node-dependency-create', {
        fromTaskIds: sourceTaskIds,
        toTaskIds: targetTaskIds,
        fromNodeId: dependencySourceNode.value.id,
        toNodeId: targetNode.id
      })
    }

    // 重置创建依赖关系状态
    isCreatingDependency.value = false
    dependencySourceNode.value = null
    tempLineEnd.value = null
    return
  }

  // 普通拖动模式
  if (isDragging.value && dragStartNodePos.value.nodeId) {
    // 计算拖动的总位移
    const currentOffset = nodeOffsets.value.get(dragStartNodePos.value.nodeId) || { x: 0, y: 0 }
    const totalDx = currentOffset.x - dragStartNodePos.value.offsetX

    // 计算天数变化
    const daysDelta = Math.round(totalDx / props.dayWidth)

    // 只有在拖动距离超过阈值时才触发时间变更
    if (daysDelta !== 0 && dragStartNodePos.value.taskId) {
      // 如果是合并节点，需要更新所有相关任务的时间
      const node = eventNodes.value.find(n => n.id === dragStartNodePos.value.nodeId)

      if (node && node.isMerged) {
        // 合并节点：更新所有相关任务
        const allTaskIds = [...node.tasks.start, ...node.tasks.end]

        allTaskIds.forEach(taskId => {
          const newDate = new Date(timelineStartDate.value)
          const originalDays = Math.round(dragStartNodePos.value.originalX / props.dayWidth)
          const newDays = originalDays + daysDelta
          newDate.setDate(newDate.getDate() + newDays)

          // 确定节点类型
          let nodeType = 'start'
          if (node.tasks.end.includes(taskId)) {
            nodeType = 'end'
          }

          emit('node-time-change', {
            taskId,
            nodeType,
            newDate: newDate.toISOString().split('T')[0],
            daysDelta
          })
        })

        // 显示统一的提示
        ElMessage.success(`已更新 ${allTaskIds.length} 个任务的时间`)
      } else {
        // 单个任务节点
        const newDate = new Date(timelineStartDate.value)
        const originalDays = Math.round(dragStartNodePos.value.originalX / props.dayWidth)
        const newDays = originalDays + daysDelta
        newDate.setDate(newDate.getDate() + newDays)

        emit('node-time-change', {
          taskId: dragStartNodePos.value.taskId,
          nodeType: dragStartNodePos.value.nodeType,
          newDate: newDate.toISOString().split('T')[0],
          daysDelta
        })

        // 单个任务的提示
        ElMessage.success('任务时间已更新')
      }

      // 重置节点的偏移量
      nodeOffsets.value.set(dragStartNodePos.value.nodeId, { x: 0, y: 0 })
    } else {
      // 没有实际拖动，重置偏移量
      nodeOffsets.value.delete(dragStartNodePos.value.nodeId)
    }
  }

  isDragging.value = false
  isPanning.value = false
  dragStartNodePos.value = {
    nodeId: null,
    offsetX: 0,
    offsetY: 0,
    nodeType: null,
    taskId: null,
    originalX: 0
  }
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
  /* cursor 通过内联样式动态设置 */
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

.dummy-work-label {
  font-size: 10px;
  fill: #999;
  font-weight: bold;
  pointer-events: none;
  text-shadow: 0 0 2px rgba(255, 255, 255, 0.8);
}

.intersection-point {
  pointer-events: none;
  opacity: 0.8;
}

.is-critical .node-circle {
  stroke: #FF8A65;
}

/* 状态栏样式 */
.network-status-bar {
  position: sticky;
  bottom: 0;
  left: 0;
  right: 0;
  background: #fff;
  border-top: 1px solid #dcdfe6;
  padding: 12px 16px;
  z-index: 100;
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.1);
}

.status-info {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.status-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
}

.status-label {
  color: #909399;
  font-weight: 500;
}

.status-value {
  color: #303133;
  font-weight: 600;
  font-family: 'Courier New', monospace;
}

.status-value.status-es,
.status-value.status-ef {
  color: #67C23A;
}

.status-value.status-ls,
.status-value.status-lf {
  color: #F56C6C;
}

.status-value.status-total-slack {
  color: #409EFF;
}

.status-divider {
  width: 1px;
  height: 16px;
  background: #dcdfe6;
}
</style>
