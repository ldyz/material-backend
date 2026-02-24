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
          <!-- 箭头路径尖端在(9, 3.5)，refX需要是圆半径减去箭头长度(9)再加偏移 -->
          <marker
            id="arrowhead-task"
            :markerWidth="10"
            :markerHeight="7"
            :refX="nodeRadius + 7"
            refY="3.5"
            orient="auto"
            markerUnits="userSpaceOnUse"
          >
            <path d="M0,0 L0,7 L9,3.5 z" fill="#64B5F6" />
          </marker>
          <!-- 任务箭头标记 - 关键任务 -->
          <marker
            id="arrowhead-critical"
            :markerWidth="10"
            :markerHeight="7"
            :refX="nodeRadius + 7"
            refY="3.5"
            orient="auto"
            markerUnits="userSpaceOnUse"
          >
            <path d="M0,0 L0,7 L9,3.5 z" fill="#FF8A65" />
          </marker>
          <!-- 依赖关系箭头标记 -->
          <marker
            id="arrowhead-dependency"
            :markerWidth="8"
            :markerHeight="8"
            :refX="nodeRadius + 6"
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
              v-for="dep in editableDependencies"
              :key="dep.key"
              :class="{ 'is-selected': selectedPathId === dep.key, 'has-custom-path': dep.hasCustomPath }"
              class="dummy-work-group"
              @click="handleDependencyClick($event, dep)"
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
              <!-- R18: 虚工作必须标注关系类型和滞后 -->
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
                {{ dep.depType }}{{ dep.lag > 0 ? '+' + dep.lag : '' }}
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
              v-for="task in editableTaskArrows"
              :key="task.id"
              :class="{
                'is-critical': task.isCritical,
                'is-selected': selectedTaskId === task.id || selectedPathId === task.id,
                'has-custom-path': task.hasCustomPath
              }"
              class="task-arrow-group"
              @click="handleTaskArrowClick($event, task)"
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

          <!-- 路径编辑弯折点 -->
          <g v-if="isEditingPath || selectedPathId" class="waypoint-editor">
            <!-- 任务箭头的弯折点 -->
            <g v-for="task in editableTaskArrows" :key="`waypoints-task-${task.id}`">
              <circle
                v-for="(wp, index) in task.waypoints"
                :key="`wp-task-${task.id}-${index}`"
                :cx="wp.x"
                :cy="wp.y"
                r="6"
                fill="#409EFF"
                stroke="#fff"
                stroke-width="2"
                class="waypoint-point"
                :class="{ 'is-selected': selectedPathId === task.id && selectedWaypointIndex === index }"
                @mousedown.stop="handleWaypointMouseDown($event, task.id, index, wp.x, wp.y)"
                @dblclick.stop="handleWaypointDoubleClick($event, task.id, index)"
              />
            </g>
            <!-- 依赖箭头的弯折点 -->
            <g v-for="dep in editableDependencies" :key="`waypoints-dep-${dep.key}`">
              <circle
                v-for="(wp, index) in dep.waypoints"
                :key="`wp-dep-${dep.key}-${index}`"
                :cx="wp.x"
                :cy="wp.y"
                r="6"
                fill="#67C23A"
                stroke="#fff"
                stroke-width="2"
                class="waypoint-point"
                :class="{ 'is-selected': selectedPathId === dep.key && selectedWaypointIndex === index }"
                @mousedown.stop="handleWaypointMouseDown($event, dep.key, index, wp.x, wp.y)"
                @dblclick.stop="handleWaypointDoubleClick($event, dep.key, index)"
              />
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

// 路径编辑状态
const isEditingPath = ref(false)  // 是否处于路径编辑模式
const selectedPathId = ref(null)  // 当前选中的路径ID
const selectedWaypointIndex = ref(null)  // 当前选中的弯折点索引
const draggedWaypoint = ref(null)  // 正在拖动的弯折点 {pathId, index, x, y}
const customPaths = ref(new Map())  // 自定义路径存储: Map<pathId, {waypoints: [{x, y}, ...], autoRoute: boolean}

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

  // ==================== 第二步：按紧前/紧后条件创建共享节点 ====================
  // R11: 相同紧前条件 → 共享开始节点
  // R11: 相同紧后条件 → 共享结束节点

  const nodeStates = new Map()  // key: 节点签名, value: { date, x, tasks: {start:[], end:[]}, taskYPositions: Map }

  // 生成节点签名的函数
  const getStartNodeSignature = (date, predecessors) => {
    const sortedPred = [...predecessors].sort((a, b) => a - b).join(',')
    return `START_${date}_${sortedPred}`
  }

  const getEndNodeSignature = (date, successors) => {
    const sortedSucc = [...successors].sort((a, b) => a - b).join(',')
    return `END_${date}_${sortedSucc}`
  }

  // 为每个任务的开始创建/查找节点
  realTasks.value.forEach(task => {
    const taskStart = new Date(task.start)
    const startDateKey = formatDate(taskStart)
    const startDays = Math.ceil((taskStart - timelineStartDate.value) / (1000 * 60 * 60 * 24))
    const predecessors = task.predecessors || []

    // R11: 相同紧前条件 → 共享同一个开始节点
    const signature = getStartNodeSignature(startDateKey, predecessors)

    if (!nodeStates.has(signature)) {
      nodeStates.set(signature, {
        date: startDateKey,
        x: startDays * props.dayWidth,
        tasks: { start: [], end: [] },
        taskYPositions: new Map() // 记录每个任务的目标Y坐标
      })
    }

    const nodeState = nodeStates.get(signature)
    nodeState.tasks.start.push(task.id)

    // 记录该任务的目标Y坐标（用于后续计算引出方向）
    if (props.taskIndexMap[task.id] !== undefined) {
      nodeState.taskYPositions.set(task.id, props.taskIndexMap[task.id] * props.rowHeight + props.rowHeight / 2)
    }
  })

  // 为每个任务的结束创建/查找节点
  realTasks.value.forEach(task => {
    const taskEnd = new Date(task.end)
    const endDateKey = formatDate(taskEnd)
    const endDays = Math.ceil((taskEnd - timelineStartDate.value) / (1000 * 60 * 60 * 24))
    const successors = task.successors || []

    // R11: 相同紧后条件 → 共享同一个结束节点
    const signature = getEndNodeSignature(endDateKey, successors)

    if (!nodeStates.has(signature)) {
      nodeStates.set(signature, {
        date: endDateKey,
        x: endDays * props.dayWidth,
        tasks: { start: [], end: [] },
        taskYPositions: new Map()
      })
    }

    const nodeState = nodeStates.get(signature)
    nodeState.tasks.end.push(task.id)

    // 记录该任务的目标Y坐标
    if (props.taskIndexMap[task.id] !== undefined) {
      nodeState.taskYPositions.set(task.id, props.taskIndexMap[task.id] * props.rowHeight + props.rowHeight / 2)
    }
  })

  // ==================== 第三步：R13/R14 - 创建唯一的开始和结束节点 ====================
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

  // 按X坐标排序节点（时间顺序）
  statesArray.sort((a, b) => {
    // 共享节点模式下，主要按X坐标（时间）排序
    if (a.x !== b.x) {
      return a.x - b.x
    }
    // X坐标相同时，按日期字符串排序保证稳定性
    return (a.date || '').localeCompare(b.date || '')
  })

  // 为每个节点分配编号
  const taskStartNodes = new Map() // taskId -> node
  const taskEndNodes = new Map()   // taskId -> node

  statesArray.forEach(state => {
    // 收集该节点关联的所有任务ID
    const allTaskIds = [...state.tasks.start, ...state.tasks.end]

    // 找到第一个任务用于获取时间参数
    const firstTaskId = allTaskIds[0]
    const firstTask = realTasks.value.find(t => t.id === firstTaskId)

    // 计算节点的Y坐标：使用所有关联任务Y坐标的中位数
    const yPositions = Array.from(state.taskYPositions.values()).sort((a, b) => a - b)
    const nodeY = yPositions.length > 0 ? yPositions[Math.floor(yPositions.length / 2)] : 100

    const node = {
      id: `node-${nodeNumber}`,
      number: nodeNumber++,
      x: state.x,
      baseX: state.x,  // 保存原始X坐标
      y: nodeY, // 使用中位数Y坐标
      baseY: nodeY,
      date: state.date,
      tasks: state.tasks,
      taskYPositions: state.taskYPositions, // 保存每个任务的Y坐标用于渲染
      isMerged: allTaskIds.length > 1,
      // 从第一个任务获取时间参数
      ES: state.tasks.start.length > 0 ? (firstTask ? firstTask.early_start : undefined) : undefined,
      EF: state.tasks.start.length > 0 ? (firstTask ? firstTask.early_finish : undefined) : undefined,
      LS: state.tasks.end.length > 0 ? (firstTask ? firstTask.late_start : undefined) : undefined,
      LF: state.tasks.end.length > 0 ? (firstTask ? firstTask.late_finish : undefined) : undefined
    }
    numberedNodes.push(node)

    // 记录每个任务的开始和结束节点
    state.tasks.start.forEach(taskId => {
      taskStartNodes.set(taskId, node)
    })
    state.tasks.end.forEach(taskId => {
      taskEndNodes.set(taskId, node)
    })
  })

  // ==================== 第五步：验证R4规则 ====================
  // R4: 箭线方向必须满足：起点节点编号 < 终点节点编号
  // 在共享节点模式下，检查每个任务的开始节点编号是否小于结束节点编号
  realTasks.value.forEach(task => {
    const startNode = taskStartNodes.get(task.id)
    const endNode = taskEndNodes.get(task.id)

    if (startNode && endNode && startNode.number >= endNode.number) {
      console.warn(`❌ R4违规: 任务${task.id} 起点(${startNode.number}) >= 终点(${endNode.number})`)
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
    const allTaskIds = [...node.tasks.start, ...node.tasks.end]
    const isCritical = allTaskIds.some(taskId => {
      const task = realTasks.value.find(t => t.id === taskId)
      return task && task.is_critical
    })
    node.isCritical = isCritical

    // 应用拖动偏移
    const offset = nodeOffsets.value.get(node.id) || { x: 0, y: 0 }
    // 保存原始X坐标到baseX（如果还没设置）
    if (node.baseX === undefined) {
      node.baseX = node.x
    }
    node.x = node.x + offset.x
    node.y = node.y + offset.y
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

// 计算多方向引出的箭头路径
// 从共享节点引出，向上/中/下三个方向，然后到达目标任务行后90度转向
// 任务箭头（连接起点和终点节点）
const taskArrows = computed(() => {
  if (!realTasks.value.length || !timelineStartDate.value) return []

  return realTasks.value.map(task => {
    // 查找该任务的开始和结束节点（共享节点）
    const startNode = eventNodes.value.find(node =>
      node.tasks.start.includes(task.id)
    )
    const endNode = eventNodes.value.find(node =>
      node.tasks.end.includes(task.id)
    )

    // 计算位置
    let startX, startY, endX, endY

    // 使用节点的实际位置（包含拖动偏移）
    // 定义任务开始和结束时间，用于计算工期
    const taskStartDate = new Date(task.start)
    const taskEndDate = new Date(task.end)

    // 计算任务的实际开始时间X坐标
    const taskStartDays = Math.ceil((taskStartDate - timelineStartDate.value) / (1000 * 60 * 60 * 24))
    const taskActualStartX = taskStartDays * props.dayWidth

    if (startNode) {
      startX = startNode.x
      startY = startNode.y
    } else {
      // 如果找不到节点，使用时间计算位置
      startX = taskActualStartX
      startY = 100
    }

    if (endNode) {
      endX = endNode.x
      endY = endNode.y
    } else {
      // 如果找不到节点，使用时间计算位置
      const endDays = Math.ceil((taskEndDate - timelineStartDate.value) / (1000 * 60 * 60 * 24))
      endX = endDays * props.dayWidth
      endY = 100
    }

    // 获取当前任务在任务列表中的Y坐标
    const taskTargetY = props.taskIndexMap[task.id] !== undefined
      ? props.taskIndexMap[task.id] * props.rowHeight + props.rowHeight / 2
      : startY

    // 实工作线的起点：从节点位置开始
    const realStartX = startX
    const realStartY = startY

    // 检查开始节点是否被多个任务共享
    const isSharedStartNode = startNode && startNode.isMerged && startNode.tasks.start.length > 1

    // 计算标签Y偏移（避免共享节点的任务标签重叠）
    let labelYOffset = 0
    if (isSharedStartNode) {
      // 找到当前任务在共享节点的tasks.start中的索引
      const taskIndex = startNode.tasks.start.indexOf(task.id)
      if (taskIndex >= 0) {
        // 每个任务向上偏移16px，避免重叠
        labelYOffset = taskIndex * 16
      }
    }

    // R19: 验证箭线只能向右（x2 ≥ x1）
    if (realStartX > endX) {
      console.warn(`❌ R19违规: 任务"${task.name}" 箭线向左 startX(${realStartX.toFixed(0)}) > endX(${endX.toFixed(0)})`)
      endX = Math.max(realStartX, endX)
    }

    // 计算路径：使用优化的最短路径算法，同时获取标签位置
    const pathResult = calculateOptimizedPath(
      realStartX, realStartY, endX, endY, taskTargetY, props.nodeRadius, false
    )
    const realPath = pathResult.path
    const labelBaseX = pathResult.labelX
    const labelBaseY = pathResult.labelY

    // 计算总时差
    const duration = Math.ceil((taskEndDate - taskStartDate) / (1000 * 60 * 60 * 24))
    const slack = task.slack || 0

    // 如果有时差，绘制时差线（波形线）
    let slackPath = null
    let slackLabelX = 0
    let slackLabelY = 0

    if (slack > 0 && !task.is_critical) {
      // 时差线：从圆右侧边缘开始，向右延伸
      // 起点在圆边缘，终点在时差结束位置
      const slackStartX = endX + props.nodeRadius - 2  // 圆右边缘附近
      const slackEndX = endX + slack * props.dayWidth
      // 终点坐标设为结束位置，箭头会自动指向该位置
      slackPath = calculateWavePath(slackStartX, endY, slackEndX, endY)
      slackLabelX = slackStartX + (slackEndX - slackStartX) / 2
      slackLabelY = endY - 12
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
      labelX: labelBaseX,
      labelY: labelBaseY - 8 - labelYOffset,
      durationX: labelBaseX,
      durationY: labelBaseY + 18 - labelYOffset,
      slackLabelX,
      slackLabelY
    }
  })
})

// 带编辑信息的任务箭头（包含自定义路径和弯折点）
const editableTaskArrows = computed(() => {
  return taskArrows.value.map(task => {
    const customPath = customPaths.value.get(task.id)
    return {
      ...task,
      waypoints: customPath?.waypoints || [],
      hasCustomPath: !!customPath
    }
  })
})

// 带编辑信息的依赖箭头（包含自定义路径和弯折点）
const editableDependencies = computed(() => {
  return renderableDependencies.value.map(dep => {
    const customPath = customPaths.value.get(dep.key)
    return {
      ...dep,
      waypoints: customPath?.waypoints || [],
      hasCustomPath: !!customPath
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

// 跟踪每个节点的出线方向使用情况（用于避免线条重叠）
const nodeExitDirectionUsage = computed(() => {
  const usage = {} // { nodeId: { right: count, top: count, bottom: count } }

  // 初始化所有节点的使用计数
  eventNodes.value.forEach(node => {
    usage[node.id] = { right: 0, top: 0, bottom: 0 }
  })

  // 分析任务箭头的出线方向
  taskArrows.value.forEach(task => {
    const path = task.realPath
    const commands = path.match(/[MLH][^,]*/g)
    if (!commands || commands.length < 2) return

    // 解析路径的起点和第二个点，确定出线方向
    const startCmd = commands[0]
    const startCoords = startCmd.slice(1).split(/[, ]+/).filter(s => s).map(Number)
    const startX = startCoords[0]
    const startY = startCoords[1]

    if (commands.length > 1) {
      const secondCmd = commands[1]
      const secondCoords = secondCmd.slice(1).split(/[, ]+/).filter(s => s).map(Number)
      const secondX = secondCoords[0]
      const secondY = secondCoords[1]

      // 找到包含这个起点的节点
      const node = eventNodes.value.find(n => {
        const dx = Math.abs(n.x - startX)
        const dy = Math.abs(n.y - startY)
        return dx < 5 && dy < 5
      })

      if (node) {
        // 判断出线方向
        if (Math.abs(secondY - startY) < 5 && secondX > startX) {
          usage[node.id].right++
        } else if (secondY < startY) {
          usage[node.id].top++
        } else if (secondY > startY) {
          usage[node.id].bottom++
        }
      }
    }
  })

  return usage
})

// 依赖关系（虚工作）- 按照双代号时标网络图规范
// 支持四种关系类型：FS/SS/FF/SF
const renderableDependencies = computed(() => {
  if (!props.dependencies.length) return []

  return props.dependencies.map((dep, index) => {
    const fromTask = props.tasks.find(t => t.id === dep.depends_on)
    const toTask = props.tasks.find(t => t.id === dep.task_id)

    if (!fromTask || !toTask) return null

    // 获取关系类型和滞后时间
    const depType = dep.type || 'FS'
    const lag = dep.lag || 0

    // 根据关系类型确定连接的节点
    let fromNode, toNode

    // 查找节点（现在是共享节点）
    if (depType === 'FS') {
      // Finish-to-Start: 前置结束 → 后续开始
      fromNode = eventNodes.value.find(node => node.tasks.end.includes(fromTask.id))
      toNode = eventNodes.value.find(node => node.tasks.start.includes(toTask.id))
    } else if (depType === 'SS') {
      // Start-to-Start: 前置开始 → 后续开始
      fromNode = eventNodes.value.find(node => node.tasks.start.includes(fromTask.id))
      toNode = eventNodes.value.find(node => node.tasks.start.includes(toTask.id))
    } else if (depType === 'FF') {
      // Finish-to-Finish: 前置结束 → 后续结束
      fromNode = eventNodes.value.find(node => node.tasks.end.includes(fromTask.id))
      toNode = eventNodes.value.find(node => node.tasks.end.includes(toTask.id))
    } else if (depType === 'SF') {
      // Start-to-Finish: 前置开始 → 后续结束
      fromNode = eventNodes.value.find(node => node.tasks.start.includes(fromTask.id))
      toNode = eventNodes.value.find(node => node.tasks.end.includes(toTask.id))
    }

    if (!fromNode || !toNode) return null

    // 计算起点和终点坐标
    let fromX = fromNode.x
    let fromY = fromNode.y
    let toX = toNode.x
    let toY = toNode.y

    // 处理时间滞后（lag）
    if (lag > 0) {
      // 滞后时间增加水平偏移
      toX += lag * props.dayWidth
    }

    // 计算路径 - 使用多方向进出的智能路由
    let path
    if (depType === 'FS' && fromNode.id === toNode.id) {
      // FS且同一节点：使用垂直虚线（从上侧出，垂直向下经过圆心附近）
      const startY = fromNode.y - props.nodeRadius
      // 让路径向下延伸超过圆心，然后箭头标记会让箭头尖端停在合适位置
      const endY = toNode.y + props.nodeRadius / 2
      path = `M ${fromX} ${startY} L ${toX} ${endY}`
    } else {
      // 使用多方向路径算法，传入出线方向使用情况以避免重叠
      const exitUsage = nodeExitDirectionUsage.value[fromNode.id] || {right: 0, top: 0, bottom: 0}
      path = calculateMultiDirectionPath(fromX, fromY, toX, toY, props.nodeRadius, exitUsage)
    }

    return {
      key: `dep-${dep.depends_on}-${dep.task_id}-${index}`,
      path,
      isDummy: true,
      depType,  // 关系类型
      lag,      // 滞后时间
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
// 约束：所有转角必须为90度，终点在圆心
function calculateOrthogonalPath(startX, startY, endX, endY) {
  const dx = Math.abs(endX - startX)
  const dy = Math.abs(endY - startY)

  // 如果起点和终点在同一行或同一列，直接画直线到圆心
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

  // 终点在圆心(endX, endY)
  return `M ${startX} ${startY} L ${midX} ${startY} L ${midX} ${endY} L ${endX} ${endY}`
}

// ==================== 优化的路径算法（类A*算法）====================
// 重要约束：所有路径必须使用正交路径（只有水平和垂直线段），所有转角必须是90度

/**
 * 计算最短正交路径（支持多方向进出）
 * 约束：只使用水平和垂直线段，所有转角必须为90度
 * @param {number} fromX - 起点X（节点中心）
 * @param {number} fromY - 起点Y（节点中心）
 * @param {number} toX - 终点X（节点中心）
 * @param {number} toY - 终点Y（节点中心）
 * @param {number} targetY - 任务行Y坐标（中间必须经过的点）
 * @param {number} radius - 圆半径
 * @param {boolean} isFromStartLine - 是否从起点线引出（已废弃，保留兼容性）
 * @returns {{path: string, labelX: number, labelY: number}} 路径和标签位置
 */
function calculateOptimizedPath(fromX, fromY, toX, toY, targetY, radius, isFromStartLine = false) {
  const dy = toY - fromY
  const dx = toX - fromX

  // 策略0: 终点在起点右侧且Y坐标相近 - 直接水平连接（0弯折）
  if (dx > 0 && Math.abs(dy) < 10) {
    const startX = fromX + radius
    const labelX = (startX + toX) / 2
    const labelY = fromY
    return {
      path: `M ${startX} ${fromY} L ${toX} ${toY}`,
      labelX,
      labelY
    }
  }

  // 策略1: 终点在右侧但Y坐标相差较大 - 优先选择经过任务行的路径
  if (dx > 0) {
    const startX = fromX + radius
    const horizontalLength = toX - startX

    // 如果任务行的水平线段足够长，使用经过任务行的路径
    if (horizontalLength > 50) {
      const labelX = (startX + toX) / 2
      const labelY = targetY
      return {
        path: `M ${startX} ${fromY} L ${startX} ${targetY} L ${toX} ${targetY} L ${toX} ${toY}`,
        labelX,
        labelY
      }
    } else {
      // 直接垂直连接（1弯折）
      const bendX = fromX + dx * 0.6
      const labelX = (startX + bendX) / 2
      const labelY = fromY
      return {
        path: `M ${startX} ${fromY} L ${bendX} ${fromY} L ${bendX} ${toY} L ${toX} ${toY}`,
        labelX,
        labelY
      }
    }
  }

  // 策略2: 终点在左侧或垂直方向
  // 如果终点在上方
  if (dy < -10) {
    // 方案1: 右侧出线 → 水平 → 向下经过任务行 → 向上到圆心
    if (dx > -radius * 2) {
      const startX = fromX + radius
      const bendX = Math.max(fromX + radius, toX + radius * 0.5)
      // 标签在任务行的水平线段上
      const labelX = (bendX + toX) / 2
      const labelY = targetY
      return {
        path: `M ${startX} ${fromY} L ${bendX} ${fromY} L ${bendX} ${targetY} L ${toX} ${targetY} L ${toX} ${toY}`,
        labelX,
        labelY
      }
    } else {
      // 空间不足，从上方出线
      const startX = fromX
      const startY = fromY - radius
      const midY = fromY + dy / 2
      const labelX = (startX + toX) / 2
      const labelY = midY
      return {
        path: `M ${startX} ${startY} L ${startX} ${midY} L ${toX} ${midY} L ${toX} ${toY}`,
        labelX,
        labelY
      }
    }
  }

  // 如果终点在下方
  if (dy > 10) {
    // 方案1: 右侧出线 → 水平 → 向上经过任务行 → 向下到圆心
    if (dx > -radius * 2) {
      const startX = fromX + radius
      const bendX = Math.max(fromX + radius, toX + radius * 0.5)
      // 标签在任务行的水平线段上
      const labelX = (bendX + toX) / 2
      const labelY = targetY
      return {
        path: `M ${startX} ${fromY} L ${bendX} ${fromY} L ${bendX} ${targetY} L ${toX} ${targetY} L ${toX} ${toY}`,
        labelX,
        labelY
      }
    } else {
      // 空间不足，从下方出线
      const startX = fromX
      const startY = fromY + radius
      const midY = fromY + dy / 2
      const labelX = (startX + toX) / 2
      const labelY = midY
      return {
        path: `M ${startX} ${startY} L ${startX} ${midY} L ${toX} ${midY} L ${toX} ${toY}`,
        labelX,
        labelY
      }
    }
  }

  // 策略3: 从起点线引出时（兼容旧逻辑）
  if (isFromStartLine) {
    const startX = fromX + radius
    const labelX = (startX + toX) / 2
    const labelY = fromY
    return {
      path: `M ${startX} ${fromY} L ${startX + 15} ${fromY} L ${startX + 15} ${targetY} L ${toX} ${targetY} L ${toX} ${toY}`,
      labelX,
      labelY
    }
  }

  // 默认 - 使用正交路径
  const startX = fromX + radius
  const midX = startX + (toX - startX) / 2
  const labelX = (startX + toX) / 2
  const labelY = fromY
  return {
    path: `M ${startX} ${fromY} L ${midX} ${fromY} L ${midX} ${toY} L ${toX} ${toY}`,
    labelX,
    labelY
  }
}
/**
 * 计算支持多方向进出的正交路径（专门用于依赖关系/虚工作）
 * 约束：只使用水平和垂直线段，所有转角必须为90度
 * 支持从节点的上/中（右）/下三个方向出，从左/上/下三个方向进
 * 优先级：弯折最少 > 路径最短 > 方向优先级（各方向较平衡）
 * 注意：右侧出线不是强制优先级，各方向优先级较接近
 * @param {number} fromX - 起点X（节点中心）
 * @param {number} fromY - 起点Y（节点中心）
 * @param {number} toX - 终点X（节点中心）
 * @param {number} toY - 终点Y（节点中心）
 * @param {number} radius - 圆半径
 * @param {object} exitUsage - 出线方向使用情况 {right: count, top: count, bottom: count}
 * @returns {string} SVG路径字符串（正交路径，所有转角为90度）
 */
function calculateMultiDirectionPath(fromX, fromY, toX, toY, radius, exitUsage = {right: 0, top: 0, bottom: 0}) {
  const dy = toY - fromY
  const dx = toX - fromX

  // 定义所有可能的路径方案
  // 重要：终点都在圆心(toX, toY)，箭头标记的refX会处理偏移
  // 根据出线方向使用情况动态调整优先级
  const pathOptions = []

  // 计算每个方向的优先级惩罚（已使用的方向降低优先级）
  const rightPenalty = exitUsage.right * 30  // 每个已有的右侧出线降低30优先级
  const topPenalty = exitUsage.top * 30     // 每个已有的上侧出线降低30优先级
  const bottomPenalty = exitUsage.bottom * 30  // 每个已有的下侧出线降低30优先级

  // 方案1: 水平对齐时的直线连接（0弯折）
  if (Math.abs(dy) < 10) {
    // 终点在右侧：从右侧出线
    if (dx > 0) {
      pathOptions.push({
        name: 'right-to-center',
        bends: 0,
        startX: fromX + radius,
        startY: fromY,
        endX: toX,
        endY: toY,
        midPoints: [],
        length: Math.abs(dx),
        priority: 100 - rightPenalty
      })
    }
    // 终点在左侧：从左侧出线（0弯折）
    if (dx < 0) {
      pathOptions.push({
        name: 'left-to-center',
        bends: 0,
        startX: fromX - radius,
        startY: fromY,
        endX: toX,
        endY: toY,
        midPoints: [],
        length: Math.abs(dx),
        priority: 100 - rightPenalty  // 使用同样的优先级
      })
    }
  }

  // 方案2: 右出到圆心（1弯折：水平→垂直）
  if (dx > 0 && Math.abs(dy) >= 10) {
    const bendX = fromX + dx * 0.6
    pathOptions.push({
      name: 'right-to-center-1bend',
      bends: 1,
      startX: fromX + radius,
      startY: fromY,
      endX: toX,
      endY: toY,
      midPoints: [
        { x: bendX, y: fromY },
        { x: bendX, y: toY }
      ],
      length: Math.abs(bendX - (fromX + radius)) + Math.abs(fromY - toY),
      priority: 90 - rightPenalty  // 基础优先级 - 右侧使用惩罚
    })
  }

  // 方案3: 右侧出线，向左回绕（2弯折）
  if (dx <= 0 && Math.abs(dy) >= 10) {
    const bendX = fromX + radius + Math.max(50, Math.abs(dx))
    pathOptions.push({
      name: 'right-to-left-around',
      bends: 2,
      startX: fromX + radius,
      startY: fromY,
      endX: toX,
      endY: toY,
      midPoints: [
        { x: bendX, y: fromY },
        { x: bendX, y: toY }
      ],
      length: Math.abs(bendX - (fromX + radius)) + Math.abs(fromY - toY) + Math.abs(bendX - toX),
      priority: 80 - rightPenalty  // 基础优先级 - 右侧使用惩罚
    })
  }

  // 方案4: 上出到圆心（1弯折：垂直-水平）
  if (dy < -10 && dx > 0) {
    pathOptions.push({
      name: 'top-to-center',
      bends: 1,
      startX: fromX,
      startY: fromY - radius,
      endX: toX,
      endY: toY,
      midPoints: [
        { x: fromX, y: fromY - radius },
        { x: toX, y: fromY - radius }
      ],
      length: Math.abs(fromY - radius - toY) + Math.abs(dx),
      priority: 95 - topPenalty  // 基础优先级 - 上侧使用惩罚
    })
  }

  // 方案5: 下出到圆心（1弯折：垂直-水平）
  if (dy > 10 && dx > 0) {
    pathOptions.push({
      name: 'bottom-to-center',
      bends: 1,
      startX: fromX,
      startY: fromY + radius,
      endX: toX,
      endY: toY,
      midPoints: [
        { x: fromX, y: fromY + radius },
        { x: toX, y: fromY + radius }
      ],
      length: Math.abs(fromY + radius - toY) + Math.abs(dx),
      priority: 95 - bottomPenalty  // 基础优先级 - 下侧使用惩罚
    })
  }

  // 方案6: 上出到圆心（2弯折：垂直-水平-垂直）
  if (Math.abs(dy) >= 10) {  // 只有dy足够大时才使用
    const midY = fromY + dy / 2
    pathOptions.push({
      name: 'top-to-center-2bend',
      bends: 2,
      startX: fromX,
      startY: fromY - radius,
      endX: toX,
      endY: toY,
      midPoints: [
        { x: fromX, y: midY },
        { x: toX, y: midY }
      ],
      length: Math.abs(dy) + Math.abs(dx),
      priority: 70 - topPenalty  // 基础优先级 - 上侧使用惩罚
    })
  }

  // 方案7: 下出到圆心（2弯折：垂直-水平-垂直）
  if (Math.abs(dy) >= 10) {  // 只有dy足够大时才使用
    const midY = fromY + dy / 2
    pathOptions.push({
      name: 'bottom-to-center-2bend',
      bends: 2,
      startX: fromX,
      startY: fromY + radius,
      endX: toX,
      endY: toY,
      midPoints: [
        { x: fromX, y: midY },
        { x: toX, y: midY }
      ],
      length: Math.abs(dy) + Math.abs(dx),
      priority: 70 - bottomPenalty  // 基础优先级 - 下侧使用惩罚
    })
  }

  // 选择最优路径：
  // 1. 优先选择弯折最少的（弯折数权重最高）
  // 2. 弯折相同时，选择优先级最高的（各方向优先级较平衡）
  // 3. 弯折和优先级都相同时，选择路径最短的
  // 使用加权排序：弯折数 × 10000 + priority × 100 - length
  // 这样确保弯折数的影响远大于priority和length
  pathOptions.sort((a, b) => {
    const scoreA = a.bends * 10000 - a.priority * 100 - a.length
    const scoreB = b.bends * 10000 - b.priority * 100 - b.length
    return scoreA - scoreB  // 分数低的优先（弯折少、priority高、length短）
  })

  const selected = pathOptions[0] || {
    name: 'default',
    bends: 2,
    startX: fromX + radius,
    startY: fromY,
    endX: toX,
    endY: toY,
    midPoints: [  // 使用正交路径：水平→垂直
      { x: fromX + radius + (toX - fromX) / 2, y: fromY },
      { x: fromX + radius + (toX - fromX) / 2, y: toY }
    ]
  }

  // 构建路径
  let path = `M ${selected.startX} ${selected.startY}`

  // 添加中间点
  if (selected.midPoints) {
    selected.midPoints.forEach(point => {
      path += ` L ${point.x} ${point.y}`
    })
  }

  // 添加终点（圆心）
  path += ` L ${selected.endX} ${selected.endY}`

  return path
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
    // 取消选中路径
    if (selectedPathId.value) {
      selectedPathId.value = null
      selectedWaypointIndex.value = null
      return
    }
    isPanning.value = true
    dragStartPos.value = { x: event.clientX - panX.value, y: event.clientY - panY.value }
  }
}

// 弯折点鼠标按下 - 开始拖动
function handleWaypointMouseDown(event, pathId, index, x, y) {
  event.stopPropagation()
  draggedWaypoint.value = {
    pathId,
    index,
    startX: x,
    startY: y,
    mouseX: event.clientX,
    mouseY: event.clientY
  }
  selectedPathId.value = pathId
  selectedWaypointIndex.value = index
}

// 弯折点双击 - 删除弯折点
function handleWaypointDoubleClick(event, pathId, index) {
  event.stopPropagation()
  const customPath = customPaths.value.get(pathId)
  if (customPath && customPath.waypoints.length > 0) {
    customPath.waypoints.splice(index, 1)
    // 如果没有弯折点了，删除自定义路径
    if (customPath.waypoints.length === 0) {
      customPaths.value.delete(pathId)
    }
    // 自动保存
    saveCustomPaths()
  }
}

// 任务箭头点击 - 选中路径或添加弯折点
function handleTaskArrowClick(event, task) {
  if (props.toolMode !== 'select') return

  event.stopPropagation()

  // 如果点击的是已选中的路径，尝试添加弯折点
  if (selectedPathId.value === task.id) {
    // 获取点击位置相对于SVG的坐标
    const svgRect = svgRef.value.getBoundingClientRect()
    const clickX = event.clientX - svgRect.left - panX.value
    const clickY = event.clientY - svgRect.top - panY.value

    addWaypointToPath(task.id, clickX, clickY)
  } else {
    // 选中路径
    selectedPathId.value = task.id
    selectedWaypointIndex.value = null
  }
}

// 依赖箭头点击 - 选中路径或添加弯折点
function handleDependencyClick(event, dep) {
  if (props.toolMode !== 'select') return

  event.stopPropagation()

  // 如果点击的是已选中的路径，尝试添加弯折点
  if (selectedPathId.value === dep.key) {
    // 获取点击位置相对于SVG的坐标
    const svgRect = svgRef.value.getBoundingClientRect()
    const clickX = event.clientX - svgRect.left - panX.value
    const clickY = event.clientY - svgRect.top - panY.value

    addWaypointToPath(dep.key, clickX, clickY)
  } else {
    // 选中路径
    selectedPathId.value = dep.key
    selectedWaypointIndex.value = null
  }
}

// 添加弯折点到路径
function addWaypointToPath(pathId, x, y) {
  if (!customPaths.value.has(pathId)) {
    customPaths.value.set(pathId, { waypoints: [], autoRoute: false })
  }
  const customPath = customPaths.value.get(pathId)

  // 找到合适的插入位置（按X坐标排序）
  let insertIndex = customPath.waypoints.findIndex(wp => wp.x > x)
  if (insertIndex === -1) {
    insertIndex = customPath.waypoints.length
  }

  customPath.waypoints.splice(insertIndex, 0, { x, y })
  saveCustomPaths()
}

// 保存自定义路径到localStorage
function saveCustomPaths() {
  const data = Array.from(customPaths.value.entries()).map(([key, value]) => [key, value])
  localStorage.setItem('network-custom-paths', JSON.stringify(data))
}

// 从localStorage加载自定义路径
function loadCustomPaths() {
  const saved = localStorage.getItem('network-custom-paths')
  if (saved) {
    try {
      const data = JSON.parse(saved)
      customPaths.value = new Map(data)
    } catch (e) {
      console.error('Failed to load custom paths:', e)
      customPaths.value = new Map()
    }
  }
}

// 使用弯折点重新计算路径
function buildPathWithWaypoints(fromX, fromY, toX, toY, waypoints) {
  if (!waypoints || waypoints.length === 0) {
    return null
  }

  // 构建路径：起点 → 弯折点1 → 弯折点2 → ... → 终点
  let path = `M ${fromX} ${fromY}`

  // 按X坐标排序弯折点
  const sortedWaypoints = [...waypoints].sort((a, b) => a.x - b.x)

  for (const wp of sortedWaypoints) {
    path += ` L ${wp.x} ${wp.y}`
  }

  path += ` L ${toX} ${toY}`
  return path
}

// 导出自定义路径数据供下载
function exportCustomPaths() {
  const data = Array.from(customPaths.value.entries()).map(([key, value]) => [key, value])
  return JSON.stringify(data, null, 2)
}

// 导入自定义路径数据
function importCustomPaths(jsonData) {
  try {
    const data = JSON.parse(jsonData)
    customPaths.value = new Map(data)
    saveCustomPaths()
    return true
  } catch (e) {
    console.error('Failed to import custom paths:', e)
    return false
  }
}

// 组件挂载时加载自定义路径
onMounted(() => {
  loadCustomPaths()
})

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

  // 拖动弯折点
  if (draggedWaypoint.value) {
    const dx = event.clientX - draggedWaypoint.value.mouseX
    const dy = event.clientY - draggedWaypoint.value.mouseY
    const newX = draggedWaypoint.value.startX + dx
    const newY = draggedWaypoint.value.startY + dy

    const customPath = customPaths.value.get(draggedWaypoint.value.pathId)
    if (customPath && customPath.waypoints[draggedWaypoint.value.index]) {
      customPath.waypoints[draggedWaypoint.value.index] = { x: newX, y: newY }
    }
    return
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
  // 弯折点拖动结束
  if (draggedWaypoint.value) {
    draggedWaypoint.value = null
    saveCustomPaths()
    return
  }

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

// 导出供父组件访问
defineExpose({
  eventNodes,
  taskArrows,
  exportCustomPaths,
  importCustomPaths,
  clearCustomPaths: () => {
    customPaths.value.clear()
    saveCustomPaths()
  }
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

/* 路径编辑相关样式 */
.waypoint-point {
  cursor: move;
  transition: all 0.2s;
  filter: drop-shadow(0 1px 3px rgba(0, 0, 0, 0.2));
}

.waypoint-point:hover {
  r: 8;
  filter: drop-shadow(0 2px 5px rgba(0, 0, 0, 0.3));
}

.waypoint-point.is-selected {
  r: 9;
  stroke-width: 3;
  filter: drop-shadow(0 2px 6px rgba(64, 158, 255, 0.5));
}

.dummy-work-group {
  cursor: pointer;
  transition: opacity 0.2s;
}

.dummy-work-group:hover {
  opacity: 0.8;
}

.dummy-work-group.is-selected .dummy-work-line {
  stroke: #409EFF;
  stroke-width: 2.5;
}

.task-arrow-group.has-custom-path .task-arrow-real,
.dummy-work-group.has-custom-path .dummy-work-line {
  stroke-dasharray: none;
  stroke-opacity: 0.8;
}

.waypoint-editor {
  pointer-events: none;
}

.waypoint-editor > * {
  pointer-events: auto;
}
</style>
