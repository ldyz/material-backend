<template>
  <div class="network-diagram" :class="{ 'fullscreen': isFullscreen }" ref="containerRef">
    <!-- 工具栏 -->
    <div class="network-toolbar">
      <!-- 缩放控制 -->
      <el-button-group size="small">
        <el-button @click="zoomOut" title="缩小">
          <el-icon><ZoomOut /></el-icon>
        </el-button>
        <el-button @click="resetZoom" title="重置">
          {{ Math.round(zoomLevel * 100) }}%
        </el-button>
        <el-button @click="zoomIn" title="放大">
          <el-icon><ZoomIn /></el-icon>
        </el-button>
      </el-button-group>

      <!-- 视图选项 -->
      <div style="margin-left: auto; display: flex; align-items: center; gap: 12px;">
        <el-button type="primary" size="small" @click="handleAddTask">
          <el-icon style="margin-right: 4px;"><Plus /></el-icon>
          添加任务
        </el-button>
        <el-checkbox v-model="showCriticalPath" @change="render" size="small">
          关键路径
        </el-checkbox>
        <el-checkbox v-model="showTimeParams" @change="render" size="small">
          时间参数
        </el-checkbox>
        <el-checkbox v-model="showTaskNames" @change="render" size="small">
          任务名称
        </el-checkbox>

        <!-- 布局方式 -->
        <el-select
          v-model="layoutMode"
          size="small"
          style="width: 120px"
          @change="handleLayoutChange"
        >
          <el-option label="自动布局" value="auto" />
          <el-option label="从左到右" value="left-right" />
          <el-option label="从上到下" value="top-down" />
        </el-select>

        <!-- 更多操作 -->
        <el-dropdown trigger="click" @command="handleMoreAction">
          <el-button size="small">
            更多
            <el-icon style="margin-left: 4px;"><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="fit-view">
                <el-icon><FullScreen /></el-icon>
                适应视图
              </el-dropdown-item>
              <el-dropdown-item command="center-view">
                <el-icon><Aim /></el-icon>
                居中视图
              </el-dropdown-item>
              <el-dropdown-item command="export" divided>
                <el-icon><Download /></el-icon>
                导出图片
              </el-dropdown-item>
              <el-dropdown-item command="toggle-fullscreen">
                <el-icon v-if="!isFullscreen"><Crop /></el-icon>
                <el-icon v-else><Close /></el-icon>
                {{ isFullscreen ? '退出全屏' : '全屏' }}
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>

    <!-- 统计信息 -->
    <div class="network-stats" v-if="stats">
      <div class="stat-item">
        <span class="stat-label">节点数</span>
        <span class="stat-value">{{ stats.nodes }}</span>
      </div>
      <div class="stat-item">
        <span class="stat-label">活动数</span>
        <span class="stat-value">{{ stats.activities }}</span>
      </div>
      <div class="stat-item">
        <span class="stat-label">关键路径</span>
        <span class="stat-value critical">{{ stats.criticalActivities }}</span>
      </div>
      <div class="stat-item">
        <span class="stat-label">总工期</span>
        <span class="stat-value">{{ stats.totalDuration }}天</span>
      </div>
    </div>

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
            markerHeight="10"
            refX="9"
            refY="3"
            orient="auto"
            markerUnits="strokeWidth"
          >
            <path d="M0,0 L0,6 L9,3 z" fill="#909399" />
          </marker>
          <marker
            id="arrowhead-critical"
            markerWidth="10"
            markerHeight="10"
            refX="9"
            refY="3"
            orient="auto"
            markerUnits="strokeWidth"
          >
            <path d="M0,0 L0,6 L9,3 z" fill="#f56c6c" />
          </marker>
          <!-- 节点阴影 -->
          <filter id="nodeShadow" x="-50%" y="-50%" width="200%" height="200%">
            <feDropShadow dx="0" dy="2" stdDeviation="3" flood-opacity="0.2"/>
          </filter>
        </defs>

        <!-- 网格背景 -->
        <g class="grid-background" v-if="showGrid">
          <pattern
            id="grid"
            :width="gridSize * zoomLevel"
            :height="gridSize * zoomLevel"
            patternUnits="userSpaceOnUse"
          >
            <path
              :d="`M ${gridSize * zoomLevel} 0 L 0 0 0 ${gridSize * zoomLevel}`"
              fill="none"
              stroke="#f0f0f0"
              stroke-width="1"
            />
          </pattern>
          <rect
            :width="svgWidth"
            :height="svgHeight"
            fill="url(#grid)"
          />
        </g>

        <!-- 变换组 (平移和缩放) -->
        <g :transform="`translate(${panX}, ${panY}) scale(${zoomLevel})`">
          <!-- 连接线 -->
          <g class="edges">
            <line
              v-for="edge in edges"
              :key="edge.id"
              :x1="edge.x1"
              :y1="edge.y1"
              :x2="edge.x2"
              :y2="edge.y2"
              :stroke="edge.isCritical ? '#f56c6c' : '#909399'"
              :stroke-width="edge.isCritical ? 3 : 2"
              :marker-end="edge.isCritical ? 'url(#arrowhead-critical)' : 'url(#arrowhead-normal)'"
              :class="{ 'edge-critical': edge.isCritical }"
              @click="handleEdgeClick(edge)"
              style="cursor: pointer"
            />
            <!-- 活动标签 -->
            <text
              v-for="edge in edges"
              v-show="showTaskNames"
              :key="'label-' + edge.id"
              :x="(edge.x1 + edge.x2) / 2"
              :y="(edge.y1 + edge.y2) / 2 - 8"
              text-anchor="middle"
              class="edge-label"
              font-size="12"
              fill="#606266"
            >
              {{ edge.label }}
            </text>
          </g>

          <!-- 节点 -->
          <g class="nodes">
            <g
              v-for="node in nodes"
              :key="node.id"
              :transform="`translate(${node.x}, ${node.y})`"
              class="node-group"
              :class="{ 'node-critical': node.isCritical, 'node-selected': selectedNode?.id === node.id }"
              @click="handleNodeClick(node)"
              @mousedown="handleNodeMouseDown($event, node)"
              style="cursor: pointer"
            >
              <!-- 节点圆形 -->
              <circle
                :r="nodeRadius"
                :fill="node.isStart ? '#67c23a' : node.isEnd ? '#f56c6c' : '#409eff'"
                filter="url(#nodeShadow)"
                stroke="white"
                stroke-width="2"
              />
              <!-- 节点编号 -->
              <text
                y="5"
                text-anchor="middle"
                fill="white"
                font-weight="bold"
                font-size="14"
              >{{ node.number }}</text>
              <!-- 时间参数 -->
              <g v-if="showTimeParams" class="time-params">
                <text
                  y="-40"
                  text-anchor="middle"
                  class="time-param-text"
                  font-size="11"
                  fill="#67c23a"
                >ES: {{ node.ES }}</text>
                <text
                  y="50"
                  text-anchor="middle"
                  class="time-param-text"
                  font-size="11"
                  fill="#e6a23c"
                >LS: {{ node.LS }}</text>
              </g>
              <!-- 虚任务标记 -->
              <circle
                v-if="node.isDummy"
                :r="nodeRadius - 3"
                fill="none"
                stroke="white"
                stroke-width="2"
                stroke-dasharray="4 2"
              />
            </g>
          </g>
        </g>
      </svg>

      <!-- 加载状态 -->
      <div v-if="loading" class="network-loading">
        <el-icon class="is-loading" :size="40" />
        <p>加载中...</p>
      </div>

      <!-- 空状态 -->
      <div v-if="!loading && nodes.length === 0" class="network-empty">
        <el-empty description="暂无网络图数据" />
      </div>
    </div>

    <!-- 图例 -->
    <div class="network-legend">
      <div class="legend-item">
        <span class="legend-color start"></span>
        <span>起点</span>
      </div>
      <div class="legend-item">
        <span class="legend-color normal"></span>
        <span>普通节点</span>
      </div>
      <div class="legend-item">
        <span class="legend-color end"></span>
        <span>终点</span>
      </div>
      <div class="legend-item">
        <span class="legend-color critical"></span>
        <span>关键路径</span>
      </div>
    </div>

    <!-- 节点详情对话框 -->
    <el-drawer
      v-model="nodeDetailVisible"
      title="节点详情"
      direction="rtl"
      size="400px"
      :close-on-click-modal="false"
    >
      <div v-if="selectedNode" class="node-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="节点编号" :span="2">
            <el-tag type="primary" size="large">{{ selectedNode.number }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="最早开始(ES)">
            {{ selectedNode.ES }}
          </el-descriptions-item>
          <el-descriptions-item label="最早结束(EF)">
            {{ selectedNode.EF }}
          </el-descriptions-item>
          <el-descriptions-item label="最迟开始(LS)">
            {{ selectedNode.LS }}
          </el-descriptions-item>
          <el-descriptions-item label="最迟结束(LF)">
            {{ selectedNode.LF }}
          </el-descriptions-item>
          <el-descriptions-item label="总时差" :span="2">
            <el-tag :type="selectedNode.isCritical ? 'danger' : 'success'" size="small">
              {{ selectedNode.slack }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="节点类型" :span="2">
            <el-tag v-if="selectedNode.isStart" type="success" size="small">起点</el-tag>
            <el-tag v-else-if="selectedNode.isEnd" type="danger" size="small">终点</el-tag>
            <el-tag v-else-if="selectedNode.isDummy" type="info" size="small">虚节点</el-tag>
            <el-tag v-else type="primary" size="small">普通节点</el-tag>
          </el-descriptions-item>
        </el-descriptions>

        <el-divider />

        <!-- 前置节点 -->
        <div v-if="selectedNode.predecessors && selectedNode.predecessors.length > 0">
          <h4>前置节点</h4>
          <el-space wrap>
            <el-tag
              v-for="pred in selectedNode.predecessors"
              :key="pred.id"
              type="info"
              size="small"
            >
              {{ pred.number }}
            </el-tag>
          </el-space>
        </div>

        <!-- 后置节点 -->
        <div v-if="selectedNode.successors && selectedNode.successors.length > 0" style="margin-top: 16px">
          <h4>后置节点</h4>
          <el-space wrap>
            <el-tag
              v-for="succ in selectedNode.successors"
              :key="succ.id"
              type="warning"
              size="small"
            >
              {{ succ.number }}
            </el-tag>
          </el-space>
        </div>
      </div>
    </el-drawer>

    <!-- 活动详情对话框 -->
    <el-drawer
      v-model="edgeDetailVisible"
      title="活动详情"
      direction="rtl"
      size="400px"
      :close-on-click-modal="false"
    >
      <div v-if="selectedEdge" class="edge-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="活动名称" :span="2">
            {{ selectedEdge.label }}
          </el-descriptions-item>
          <el-descriptions-item label="起点">
            节点 {{ selectedEdge.from }}
          </el-descriptions-item>
          <el-descriptions-item label="终点">
            节点 {{ selectedEdge.to }}
          </el-descriptions-item>
          <el-descriptions-item label="工期">
            {{ selectedEdge.duration }} 天
          </el-descriptions-item>
          <el-descriptions-item label="类型">
            <el-tag :type="selectedEdge.isCritical ? 'danger' : 'primary'" size="small">
              {{ selectedEdge.isCritical ? '关键活动' : '普通活动' }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </el-drawer>

    <!-- 任务编辑对话框 -->
    <el-dialog
      v-model="editDialogVisible"
      :title="editingTask ? '编辑任务' : '新建任务'"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form :model="taskForm" :rules="taskFormRules" ref="taskFormRef" label-width="100px">
        <el-form-item label="任务名称" prop="name">
          <el-input v-model="taskForm.name" placeholder="请输入任务名称" />
        </el-form-item>
        <el-form-item label="开始日期" prop="start">
          <el-date-picker
            v-model="taskForm.start"
            type="date"
            placeholder="选择开始日期"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="结束日期" prop="end">
          <el-date-picker
            v-model="taskForm.end"
            type="date"
            placeholder="选择结束日期"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="进度" prop="progress">
          <el-slider v-model="taskForm.progress" :marks="{ 0: '0%', 50: '50%', 100: '100%' }" />
        </el-form-item>
        <el-form-item label="优先级" prop="priority">
          <el-radio-group v-model="taskForm.priority">
            <el-radio-button label="urgent">紧急</el-radio-button>
            <el-radio-button label="high">高</el-radio-button>
            <el-radio-button label="medium">中</el-radio-button>
            <el-radio-button label="low">低</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="taskForm.notes" type="textarea" :rows="3" placeholder="任务备注" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveTask" :loading="saving">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import {
  ZoomIn,
  ZoomOut,
  ArrowDown,
  FullScreen,
  Close,
  Download,
  Crop,
  Aim,
  Plus
} from '@element-plus/icons-vue'
import html2canvas from 'html2canvas'
import { progressApi } from '@/api'

const props = defineProps({
  projectId: {
    type: [Number, String],
    required: true
  },
  scheduleData: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['node-selected', 'position-updated', 'task-updated'])

// 视图状态
const zoomLevel = ref(1)
const panX = ref(0)
const panY = ref(0)
const isFullscreen = ref(false)
const loading = ref(false)
const showGrid = ref(true)
const gridSize = ref(20)
const nodeRadius = ref(25)

// 显示选项
const showCriticalPath = ref(true)
const showTimeParams = ref(true)
const showTaskNames = ref(true)
const layoutMode = ref('auto')

// 画布尺寸
const svgWidth = ref(2000)
const svgHeight = ref(1200)

// 拖拽状态
const isDragging = ref(false)
const isPanning = ref(false)
const dragStartPos = ref({ x: 0, y: 0 })
const panStartPos = ref({ x: 0, y: 0 })
const draggedNode = ref(null)

// 选中项
const selectedNode = ref(null)
const selectedEdge = ref(null)
const nodeDetailVisible = ref(false)
const edgeDetailVisible = ref(false)

// 任务编辑
const editDialogVisible = ref(false)
const editingTask = ref(null)
const saving = ref(false)
const taskFormRef = ref(null)
const taskForm = ref({
  name: '',
  start: null,
  end: null,
  progress: 0,
  priority: 'medium',
  notes: ''
})

const taskFormRules = {
  name: [{ required: true, message: '请输入任务名称', trigger: 'blur' }],
  start: [{ required: true, message: '请选择开始日期', trigger: 'change' }],
  end: [{ required: true, message: '请选择结束日期', trigger: 'change' }]
}

// Refs
const containerRef = ref(null)
const canvasContainerRef = ref(null)
const svgRef = ref(null)

// 节点和边数据
const nodes = ref([])
const edges = ref([])

// 统计信息
const stats = computed(() => {
  return {
    nodes: nodes.value.filter(n => !n.isDummy).length,
    activities: edges.value.filter(e => !e.isDummy).length,
    criticalActivities: edges.value.filter(e => e.isCritical).length,
    totalDuration: calculateTotalDuration()
  }
})

// 计算总工期
const calculateTotalDuration = () => {
  if (nodes.value.length === 0) return 0
  const endNode = nodes.value.find(n => n.isEnd)
  return endNode ? endNode.EF : 0
}

// 解析进度数据并构建网络图
const buildNetworkDiagram = () => {
  if (!props.scheduleData || !props.scheduleData.activities) {
    nodes.value = []
    edges.value = []
    return
  }

  const activities = props.scheduleData.activities || {}
  const nodeList = []
  const edgeList = []
  const nodeMap = new Map()

  // 构建节点
  let nodeNumber = 1
  for (const [key, activity] of Object.entries(activities)) {
    // 创建起始节点
    if (!nodeMap.has(activity.id)) {
      const node = {
        id: `node-${activity.id}`,
        number: nodeNumber++,
        ES: Math.floor(activity.earliest_start / 86400),
        EF: Math.floor(activity.earliest_finish / 86400),
        LS: Math.floor(activity.latest_start / 86400),
        LF: Math.floor(activity.latest_finish / 86400),
        slack: Math.floor((activity.latest_start - activity.earliest_start) / 86400),
        isCritical: activity.is_critical,
        isStart: activity.predecessors?.length === 0,
        isEnd: activity.successors?.length === 0,
        isDummy: activity.is_dummy || false,
        predecessors: [],
        successors: []
      }
      nodeMap.set(activity.id, node)
      nodeList.push(node)
    }
  }

  // 构建边（活动）
  for (const [key, activity] of Object.entries(activities)) {
    const node = nodeMap.get(activity.id)
    if (!node) continue

    // 添加前置依赖边
    if (activity.predecessors) {
      activity.predecessors.forEach(predId => {
        const predNode = nodeMap.get(predId)
        if (predNode) {
          edgeList.push({
            id: `edge-${predId}-${activity.id}`,
            from: predNode.number,
            to: node.number,
            x1: 0, // 将在布局时计算
            y1: 0,
            x2: 0,
            y2: 0,
            label: activity.name || `活动 ${node.number}`,
            duration: activity.duration,
            isCritical: activity.is_critical,
            isDummy: activity.is_dummy || false
          })
          predNode.successors.push(node)
          node.predecessors.push(predNode)
        }
      })
    }
  }

  // 自动布局
  layoutNodes(nodeList, edgeList)

  nodes.value = nodeList
  edges.value = edgeList
}

// 节点布局算法
const layoutNodes = (nodeList, edgeList) => {
  const width = 1800
  const height = 1000
  const nodeSpacing = 200
  const levelSpacing = 250

  // 按ES值分组（层级）
  const levels = new Map()
  nodeList.forEach(node => {
    const level = Math.floor(node.ES / 5) // 每5天一个层级
    if (!levels.has(level)) {
      levels.set(level, [])
    }
    levels.get(level).push(node)
  })

  // 分配位置
  levels.forEach((nodesInLevel, level) => {
    const totalWidth = nodesInLevel.length * nodeSpacing
    const startX = (width - totalWidth) / 2 + nodeSpacing / 2

    nodesInLevel.forEach((node, index) => {
      node.x = startX + index * nodeSpacing
      node.y = 100 + level * levelSpacing
    })
  })

  // 更新边的坐标
  edgeList.forEach(edge => {
    const fromNode = nodeList.find(n => n.number === edge.from)
    const toNode = nodeList.find(n => n.number === edge.to)
    if (fromNode && toNode) {
      edge.x1 = fromNode.x
      edge.y1 = fromNode.y
      edge.x2 = toNode.x
      edge.y2 = toNode.y
    }
  })
}

// 渲染
const render = () => {
  buildNetworkDiagram()
}

// 处理布局变化
const handleLayoutChange = () => {
  render()
}

// 缩放控制
const zoomIn = () => {
  zoomLevel.value = Math.min(zoomLevel.value + 0.1, 3)
}

const zoomOut = () => {
  zoomLevel.value = Math.max(zoomLevel.value - 0.1, 0.3)
}

const resetZoom = () => {
  zoomLevel.value = 1
  centerView()
}

// 适应视图
const fitView = () => {
  zoomLevel.value = 0.8
  centerView()
  ElMessage.success('视图已适应')
}

// 居中视图
const centerView = () => {
  if (!containerRef.value || !canvasContainerRef.value) return

  const containerRect = canvasContainerRef.value.getBoundingClientRect()
  panX.value = (containerRect.width - svgWidth.value * zoomLevel.value) / 2
  panY.value = (containerRect.height - svgHeight.value * zoomLevel.value) / 2
}

// 鼠标滚轮缩放
const handleWheel = (event) => {
  const delta = event.deltaY > 0 ? -0.1 : 0.1
  zoomLevel.value = Math.max(0.3, Math.min(3, zoomLevel.value + delta))
}

// 鼠标拖拽
const handleMouseDown = (event) => {
  // 只在空白区域才允许平移
  if (event.target.tagName === 'svg') {
    isPanning.value = true
    panStartPos.value = {
      x: event.clientX - panX.value,
      y: event.clientY - panY.value
    }
  }
}

const handleMouseMove = (event) => {
  if (isPanning.value) {
    panX.value = event.clientX - panStartPos.value.x
    panY.value = event.clientY - panStartPos.value.y
  }

  if (isDragging.value && draggedNode.value) {
    const node = draggedNode.value
    const scale = zoomLevel.value
    node.x = (event.clientX - panX.value) / scale
    node.y = (event.clientY - panY.value) / scale

    // 更新相关边的坐标
    edges.value.forEach(edge => {
      if (edge.from === node.number) {
        edge.x1 = node.x
        edge.y1 = node.y
      }
      if (edge.to === node.number) {
        edge.x2 = node.x
        edge.y2 = node.y
      }
    })
  }
}

const handleMouseUp = () => {
  if (isDragging.value) {
    // 保存位置
    emit('position-updated', {
      nodes: nodes.value.map(n => ({
        id: n.id,
        x: n.x,
        y: n.y
      }))
    })
  }
  isDragging.value = false
  isPanning.value = false
  draggedNode.value = null
}

// 节点拖拽
const handleNodeMouseDown = (event, node) => {
  event.stopPropagation()
  isDragging.value = true
  draggedNode.value = node
  dragStartPos.value = {
    x: event.clientX,
    y: event.clientY
  }
}

// 节点点击
const handleNodeClick = (node) => {
  selectedNode.value = node
  nodeDetailVisible.value = true
  emit('node-selected', node)
}

// 边点击
const handleEdgeClick = (edge) => {
  selectedEdge.value = edge
  edgeDetailVisible.value = true
}

// 更多操作
const handleMoreAction = (command) => {
  switch (command) {
    case 'fit-view':
      fitView()
      break
    case 'center-view':
      centerView()
      ElMessage.success('视图已居中')
      break
    case 'export':
      handleExportImage()
      break
    case 'toggle-fullscreen':
      toggleFullscreen()
      break
  }
}

// 导出图片
const handleExportImage = async () => {
  if (!svgRef.value) return

  try {
    const canvas = await html2canvas(svgRef.value.parentElement, {
      backgroundColor: '#ffffff',
      logging: false,
      useCORS: true,
      scale: 2
    })

    canvas.toBlob((blob) => {
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `网络图_${props.projectId}_${new Date().getTime()}.png`
      a.click()
      URL.revokeObjectURL(url)
      ElMessage.success('导出成功')
    })
  } catch (error) {
    console.error('导出失败:', error)
    ElMessage.error('导出失败，请重试')
  }
}

// 全屏切换
const toggleFullscreen = () => {
  isFullscreen.value = !isFullscreen.value

  if (isFullscreen.value) {
    document.body.classList.add('network-fullscreen')
  } else {
    document.body.classList.remove('network-fullscreen')
  }

  nextTick(() => {
    centerView()
  })
}

// 监听数据变化
watch(
  () => props.scheduleData,
  () => {
    render()
  },
  { deep: true }
)

onMounted(() => {
  render()
  centerView()
})

onUnmounted(() => {
  document.body.classList.remove('network-fullscreen')
})

// 添加任务
const handleAddTask = () => {
  editingTask.value = null
  taskForm.value = {
    name: '',
    start: new Date(),
    end: new Date(),
    progress: 0,
    priority: 'medium',
    notes: ''
  }
  editDialogVisible.value = true
}

// 保存任务
const handleSaveTask = async () => {
  if (!taskFormRef.value) return

  try {
    await taskFormRef.value.validate()

    saving.value = true

    const taskData = {
      project_id: props.projectId,
      name: taskForm.value.name,
      start_date: formatDate(taskForm.value.start),
      end_date: formatDate(taskForm.value.end),
      progress: taskForm.value.progress,
      priority: taskForm.value.priority,
      description: taskForm.value.notes
    }

    if (editingTask.value) {
      // 更新现有任务
      await progressApi.update(editingTask.value.id, taskData)
      ElMessage.success('任务更新成功')
      emit('task-updated', { ...editingTask.value, ...taskData })
    } else {
      // 创建新任务
      await progressApi.create(taskData)
      ElMessage.success('任务创建成功')
      emit('task-updated', taskData)
    }

    editDialogVisible.value = false
    editingTask.value = null
  } catch (error) {
    if (error !== false) {
      console.error('保存任务失败:', error)
      ElMessage.error('保存任务失败')
    }
  } finally {
    saving.value = false
  }
}

// 格式化日期
const formatDate = (date) => {
  const d = new Date(date)
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}
</script>

<style scoped>
.network-diagram {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #fff;
  outline: none;
}

.network-diagram.fullscreen {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 9999;
  background: #fff;
}

.network-toolbar {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  background: #fff;
  border-bottom: 1px solid #dcdfe6;
  flex-wrap: wrap;
  gap: 8px;
}

.network-stats {
  display: flex;
  align-items: center;
  gap: 24px;
  padding: 12px 16px;
  background: #f5f7fa;
  border-bottom: 1px solid #dcdfe6;
  flex-wrap: wrap;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
}

.stat-label {
  color: #909399;
}

.stat-value {
  font-weight: bold;
  color: #303133;
}

.stat-value.critical {
  color: #f56c6c;
}

.network-canvas-container {
  flex: 1;
  position: relative;
  overflow: hidden;
  background: #fafafa;
}

.network-svg {
  display: block;
  background: white;
}

.network-loading,
.network-empty {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  color: #909399;
}

.node-group {
  transition: all 0.3s;
}

.node-group:hover {
  filter: brightness(1.1);
}

.node-group.node-selected circle {
  stroke: #409eff;
  stroke-width: 4;
}

.node-critical circle {
  stroke: #f56c6c;
  stroke-width: 3;
}

.time-param-text {
  font-weight: 500;
}

.edge-critical {
  stroke-width: 3 !important;
}

.edge-label {
  pointer-events: none;
  font-weight: 500;
}

.network-legend {
  display: flex;
  align-items: center;
  gap: 24px;
  padding: 12px 16px;
  background: #f5f7fa;
  border-top: 1px solid #dcdfe6;
  flex-wrap: wrap;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #606266;
}

.legend-color {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  border: 2px solid;
}

.legend-color.start {
  background: #67c23a;
  border-color: #67c23a;
}

.legend-color.normal {
  background: #409eff;
  border-color: #409eff;
}

.legend-color.end {
  background: #f56c6c;
  border-color: #f56c6c;
}

.legend-color.critical {
  background: #f56c6c;
  border-color: #f56c6c;
}

/* 节点详情 */
.node-detail,
.edge-detail {
  padding: 0 20px;
}

/* 全屏样式 */
:deep(.network-fullscreen) {
  overflow: hidden;
}

/* 打印样式 */
@media print {
  .network-toolbar,
  .network-stats,
  .network-legend {
    display: none !important;
  }

  .network-diagram {
    position: static !important;
    height: auto !important;
  }

  .network-canvas-container {
    overflow: visible !important;
  }
}
</style>
