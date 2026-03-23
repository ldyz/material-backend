<template>
  <div class="workflow-editor" :class="{ 'fullscreen': isFullscreen }">
    <div class="editor-toolbar">
      <el-button-group>
        <el-button :icon="Plus" @click="addNode('approval')">添加审批节点</el-button>
        <el-button :icon="Plus" @click="addNode('parallel')">添加并行节点</el-button>
        <el-button :icon="Plus" @click="addNode('merge')">添加合并节点</el-button>
      </el-button-group>
      <div class="toolbar-spacer"></div>
      <el-button :icon="isFullscreen ? Rank : FullScreen" @click="toggleFullscreen">
        {{ isFullscreen ? '窗口模式' : '全屏模式' }}
      </el-button>
      <el-button :icon="Check" type="success" @click="handleSave">保存工作流</el-button>
      <el-button :icon="Refresh" @click="handleReset">重置</el-button>
    </div>

    <div class="editor-main">
      <!-- 左侧画布区域 -->
      <div class="editor-container" ref="containerRef">
        <!-- SVG画布 -->
        <svg
          class="workflow-canvas"
          ref="canvasRef"
          @mousedown="handleCanvasMouseDown"
          @mousemove="handleCanvasMouseMove"
          @mouseup="handleCanvasMouseUp"
          @contextmenu.prevent
        >
          <!-- 网格背景 -->
          <defs>
            <pattern id="grid" width="20" height="20" patternUnits="userSpaceOnUse">
              <path d="M 20 0 L 0 0 0 20" fill="none" stroke="#e0e0e0" stroke-width="0.5"/>
            </pattern>
          </defs>
          <rect width="100%" height="100%" fill="url(#grid)" />

          <!-- 连接线 -->
          <g class="edges">
            <g
              v-for="edge in edges"
              :key="edge.id"
              :class="['edge-group', { 'edge-selected': selectedEdge === edge.id }]"
            >
              <!-- 连线路径 -->
              <path
                :d="getEdgePath(edge)"
                :class="['edge', { 'edge-selected': selectedEdge === edge.id }]"
                :marker-end="selectedEdge === edge.id ? 'url(#arrowhead-selected)' : 'url(#arrowhead)'"
                @click="selectEdge(edge.id)"
                @contextmenu.prevent="showEdgeContextMenu($event, edge.id)"
              />

              <!-- 连线操作按钮（选中时显示） -->
              <g v-if="selectedEdge === edge.id" class="edge-actions">
                <circle
                  :cx="getEdgeMidPoint(edge).x"
                  :cy="getEdgeMidPoint(edge).y"
                  r="10"
                  fill="#f56c6c"
                  class="edge-delete-btn"
                  @click.stop="deleteEdge(edge.id)"
                />
                <text
                  :x="getEdgeMidPoint(edge).x"
                  :y="getEdgeMidPoint(edge).y + 4"
                  text-anchor="middle"
                  fill="white"
                  font-size="12"
                  font-weight="bold"
                  pointer-events="none"
                >×</text>
              </g>
            </g>

            <!-- 箭头 -->
            <defs>
              <marker id="arrowhead" markerWidth="10" markerHeight="7" refX="9" refY="3.5" orient="auto">
                <polygon points="0 0, 10 3.5, 0 7" fill="#409eff" />
              </marker>
              <marker id="arrowhead-selected" markerWidth="10" markerHeight="7" refX="9" refY="3.5" orient="auto">
                <polygon points="0 0, 10 3.5, 0 7" fill="#f56c6c" />
              </marker>
            </defs>
          </g>

          <!-- 节点 -->
          <g class="nodes">
            <g
              v-for="node in nodes"
              :key="node.id"
              :class="['node', `node-${node.node_type}`, { 'node-selected': selectedNode === node.id }]"
              :transform="`translate(${node.x}, ${node.y})`"
              @mousedown.stop="handleNodeMouseDown($event, node)"
            >
              <!-- 节点形状 -->
              <rect
                :width="getNodeWidth(node)"
                :height="getNodeHeight(node)"
                :class="['node-shape', `shape-${node.node_type}`]"
                rx="6"
                ry="6"
              />

              <!-- 节点图标 -->
              <text :x="getNodeWidth(node) / 2" y="20" text-anchor="middle" class="node-icon">
                {{ getNodeIcon(node.node_type) }}
              </text>

              <!-- 节点名称 -->
              <text :x="getNodeWidth(node) / 2" y="40" text-anchor="middle" class="node-name">
                {{ node.node_name }}
              </text>

              <!-- 节点类型标签 -->
              <text :x="getNodeWidth(node) / 2" y="58" text-anchor="middle" class="node-type">
                {{ getNodeTypeText(node.node_type) }}
              </text>

              <!-- 输入连接点（左侧） -->
              <circle
                v-if="node.node_type !== 'start'"
                cx="0"
                :cy="getNodeHeight(node) / 2"
                r="6"
                class="connector connector-in"
                @mouseup.stop="completeConnectToInput($event, node)"
                @mouseover="highlightInput(node.id, true)"
                @mouseout="highlightInput(node.id, false)"
              />

              <!-- 输出连接点（右侧） -->
              <circle
                v-if="node.node_type !== 'end'"
                :cx="getNodeWidth(node)"
                :cy="getNodeHeight(node) / 2"
                r="6"
                class="connector connector-out"
                @mousedown.stop="startConnect($event, node)"
                @mouseup.stop="completeConnectFromOutput($event, node)"
              />

              <!-- 连线端点拖拽（选中连线时显示源节点和目标节点的拖拽点） -->
              <circle
                v-if="selectedEdge && isEdgeEndpoint(selectedEdge, node.id, 'source')"
                :cx="getNodeWidth(node)"
                :cy="getNodeHeight(node) / 2"
                r="8"
                class="connector connector-drag"
                @mousedown.stop="startEditEdge($event, selectedEdge, 'source')"
              />
              <circle
                v-if="selectedEdge && isEdgeEndpoint(selectedEdge, node.id, 'target')"
                cx="0"
                :cy="getNodeHeight(node) / 2"
                r="8"
                class="connector connector-drag"
                @mousedown.stop="startEditEdge($event, selectedEdge, 'target')"
              />
            </g>
          </g>

          <!-- 临时连线（拖拽时显示） -->
          <path
            v-if="tempEdge"
            :d="tempEdge"
            class="edge temp-edge"
            stroke-dasharray="5,5"
          />
        </svg>
      </div>

      <!-- 右侧属性面板 -->
      <div class="property-panel" v-if="selectedNode || selectedEdge">
        <div class="panel-header">
          <h3>{{ panelTitle }}</h3>
          <el-button :icon="Close" circle size="small" @click="clearSelection" />
        </div>

        <div class="panel-content">
          <!-- 节点属性 -->
          <template v-if="selectedNodeData">
            <el-form :model="nodeForm" label-width="90px" size="small">
              <el-form-item label="节点名称">
                <el-input v-model="nodeForm.node_name" placeholder="请输入节点名称" @input="updateNodeData" />
              </el-form-item>

              <el-form-item label="节点类型">
                <el-select v-model="nodeForm.node_type" placeholder="节点类型" :disabled="nodeForm.node_type === 'start' || nodeForm.node_type === 'end'" @change="updateNodeData">
                  <el-option label="开始节点" value="start" />
                  <el-option label="审批节点" value="approval" />
                  <el-option label="并行节点" value="parallel" />
                  <el-option label="合并节点" value="merge" />
                  <el-option label="结束节点" value="end" />
                </el-select>
              </el-form-item>

              <el-form-item label="节点描述">
                <el-input v-model="nodeForm.description" type="textarea" :rows="2" placeholder="请输入描述" @input="updateNodeData" />
              </el-form-item>

              <template v-if="nodeForm.node_type === 'approval'">
                <el-divider content-position="left">审批配置</el-divider>

                <el-form-item label="审批类型">
                  <el-select v-model="nodeForm.approval_type" placeholder="请选择审批类型" @change="updateNodeData">
                    <el-option label="顺序审批" value="sequential" />
                    <el-option label="并行审批" value="parallel" />
                    <el-option label="任一审批" value="any" />
                  </el-select>
                </el-form-item>

                <el-form-item label="超时时间">
                  <el-input-number v-model="nodeForm.timeout_hours" :min="0" :max="720" @change="updateNodeData" />
                  <span style="margin-left: 8px; color: #909399; font-size: 12px">小时</span>
                </el-form-item>

                <el-form-item label="自动通过">
                  <el-switch v-model="nodeForm.auto_approve" @change="updateNodeData" />
                </el-form-item>

                <el-form-item label="必需节点">
                  <el-switch v-model="nodeForm.is_required" @change="updateNodeData" />
                </el-form-item>

                <el-divider content-position="left">审批人配置</el-divider>

                <div class="approvers-list">
                  <div v-if="!nodeForm.approvers || nodeForm.approvers.length === 0" class="empty-approvers">
                    <el-empty description="暂无审批人" :image-size="60" />
                  </div>
                  <div v-else>
                    <div v-for="(approver, index) in nodeForm.approvers" :key="index" class="approver-item">
                      <div class="approver-info">
                        <el-tag :type="approver.approver_type === 'user' ? 'primary' : 'success'" size="small">
                          {{ approver.approver_type === 'user' ? '用户' : '角色' }}
                        </el-tag>
                        <span class="approver-name">{{ getApproverName(approver) }}</span>
                        <el-tag size="small" type="info">顺序: {{ approver.sequence }}</el-tag>
                      </div>
                      <div class="approver-actions">
                        <el-button type="danger" size="small" :icon="Delete" circle @click="removeApprover(index)" />
                      </div>
                    </div>
                  </div>

                  <el-button
                    type="primary"
                    :icon="Plus"
                    size="small"
                    style="width: 100%; margin-top: 12px"
                    @click="showAddApproverDialog"
                  >
                    添加审批人
                  </el-button>
                </div>
              </template>

              <template v-else-if="nodeForm.node_type === 'parallel'">
                <el-divider content-position="left">并行配置</el-divider>
                <el-form-item label="分支数量">
                  <el-input-number v-model="parallelBranches" :min="2" :max="10" :disabled="true" />
                  <span style="margin-left: 8px; color: #909399; font-size: 12px">通过连线配置</span>
                </el-form-item>
              </template>

              <template v-else-if="nodeForm.node_type === 'merge'">
                <el-divider content-position="left">合并配置</el-divider>
                <el-form-item label="合并方式">
                  <el-select v-model="nodeForm.approval_type" placeholder="请选择合并方式" @change="updateNodeData">
                    <el-option label="全部完成" value="sequential" />
                    <el-option label="任一完成" value="any" />
                  </el-select>
                </el-form-item>
              </template>

              <el-divider content-position="left">节点操作</el-divider>

              <div class="node-actions">
                <el-button type="danger" :icon="Delete" @click="deleteSelectedNode" style="width: 100%">
                  删除节点
                </el-button>
              </div>
            </el-form>
          </template>

          <!-- 连线属性 -->
          <template v-else-if="selectedEdgeData">
            <el-form label-width="90px" size="small">
              <el-form-item label="来源节点">
                <div class="edge-info">{{ getSelectedEdgeFromNode() }}</div>
              </el-form-item>
              <el-form-item label="目标节点">
                <div class="edge-info">{{ getSelectedEdgeToNode() }}</div>
              </el-form-item>
              <el-form-item label="条件表达式">
                <el-input
                  v-model="edgeCondition"
                  type="textarea"
                  :rows="3"
                  placeholder="可选：设置流转条件表达式"
                  @input="updateEdgeCondition"
                />
              </el-form-item>
              <el-divider content-position="left">连线操作</el-divider>
              <div class="edge-actions">
                <el-button type="danger" :icon="Delete" @click="deleteSelectedEdge" style="width: 100%; margin-bottom: 8px">
                  删除连线
                </el-button>
                <p style="color: #909399; font-size: 12px; margin: 8px 0;">
                  提示：点击连线中点的×按钮也可删除
                </p>
              </div>
            </el-form>
          </template>
        </div>
      </div>
    </div>

    <!-- 添加审批人对话框 -->
    <el-dialog
      v-model="addApproverDialogVisible"
      title="添加审批人"
      width="500px"
    >
      <el-form :model="newApprover" label-width="90px">
        <el-form-item label="审批类型" required>
          <el-radio-group v-model="newApprover.approver_type">
            <el-radio label="user">指定用户</el-radio>
            <el-radio label="role">指定角色</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item :label="newApprover.approver_type === 'user' ? '选择用户' : '选择角色'" required>
          <el-select
            v-if="newApprover.approver_type === 'user'"
            v-model="newApprover.approver_id"
            placeholder="请选择用户"
            filterable
            style="width: 100%"
          >
            <el-option
              v-for="user in userList"
              :key="user.id"
              :label="`${user.username} - ${user.real_name || user.username}`"
              :value="user.id"
            />
          </el-select>
          <el-select
            v-else-if="newApprover.approver_type === 'role'"
            v-model="newApprover.approver_id"
            placeholder="请选择角色"
            filterable
            style="width: 100%"
          >
            <el-option
              v-for="role in roleList"
              :key="role.id"
              :label="role.name"
              :value="role.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="审批顺序">
          <el-input-number v-model="newApprover.sequence" :min="0" :max="100" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="addApproverDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="addApprover">确定</el-button>
      </template>
    </el-dialog>

    <!-- 浮动半透明操作提示 -->
    <div class="tips-overlay" v-if="showTips" @click.self="showTips = false">
      <div class="tips-content">
        <div class="tips-header">
          <strong>操作提示</strong>
          <el-icon class="close-btn" @click="showTips = false"><Close /></el-icon>
        </div>
        <ul class="tips-list">
          <li><strong>拖拽节点</strong>：按住节点可移动位置</li>
          <li><strong>创建连线</strong>：从节点右侧蓝点拖拽到另一个节点左侧绿点</li>
          <li><strong>连接规则</strong>：
            <ul class="tips-sub-list">
              <li>开始节点（绿色）不能被连接，只能从它发出连线</li>
              <li>结束节点（红色）不能发出连线，只能被连接</li>
              <li>其他节点可以相互连接</li>
            </ul>
          </li>
          <li><strong>选择节点/连线</strong>：点击节点或连线可选中并显示属性</li>
          <li><strong>修改连线目标</strong>：选中连线后，拖拽橙色圆点到新目标节点</li>
          <li><strong>删除连线</strong>：右键点击连线选择删除，或选中连线后点击删除按钮/连线中点的×按钮</li>
          <li><strong>删除节点</strong>：选中节点后在属性面板中点击删除</li>
        </ul>
        <div class="tips-footer">
          <el-link type="primary" @click="showTips = false">不再显示</el-link>
        </div>
      </div>
    </div>
    <div class="tips-toggle" v-else>
      <el-button type="primary" link @click="showTips = true">
        <el-icon><QuestionFilled /></el-icon>
        显示操作提示
      </el-button>
    </div>

    <!-- 右键菜单 -->
    <teleport to="body">
      <div
        v-if="contextMenuVisible"
        class="edge-context-menu"
        :style="{ left: contextMenuPosition.x + 'px', top: contextMenuPosition.y + 'px' }"
        @click.stop
      >
        <div class="context-menu-item" @click="deleteEdgeFromContext">
          <el-icon :size="14"><Delete /></el-icon>
          <span>删除连线</span>
        </div>
      </div>
    </teleport>

    <!-- 点击遮罩层关闭右键菜单 -->
    <div
      v-if="contextMenuVisible"
      class="context-menu-overlay"
      @click="hideContextMenu"
      @contextmenu.prevent="hideContextMenu"
    ></div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Check, Refresh, Delete, Close, QuestionFilled, FullScreen, Rank } from '@element-plus/icons-vue'

const props = defineProps({
  workflowId: {
    type: Number,
    default: null
  },
  module: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['save', 'change'])

const containerRef = ref(null)
const canvasRef = ref(null)

// 节点和边
const nodes = ref([])
const edges = ref([])

// 选择状态
const selectedNode = ref(null)
const selectedEdge = ref(null)

// 拖拽状态
const draggingNode = ref(null)
const dragOffset = ref({ x: 0, y: 0 })
const isDraggingNode = ref(false)

// 连线状态
const connectingFrom = ref(null)
const isConnectingEdge = ref(false)
const tempEdge = ref('')

// 编辑连线状态
const editingEdge = ref(null)
const editingEndpoint = ref(null) // 'source' or 'target'
const highlightedInput = ref(null)

// 节点表单（属性面板）
const nodeForm = reactive({
  node_key: '',
  node_name: '',
  node_type: 'approval',
  description: '',
  approval_type: 'sequential',
  timeout_hours: 24,
  auto_approve: false,
  is_required: true,
  approvers: []
})

// 边条件
const edgeCondition = ref('')

// 并行分支数
const parallelBranches = ref(2)

// 添加审批人对话框
const addApproverDialogVisible = ref(false)
const newApprover = reactive({
  approver_type: 'role',
  approver_id: null,
  approver_name: '',
  sequence: 0
})

// 用户和角色列表
const userList = ref([])
const roleList = ref([])

// 操作提示显示状态
const showTips = ref(true)

// 全屏状态
const isFullscreen = ref(false)

// 右键菜单状态
const contextMenuVisible = ref(false)
const contextMenuPosition = ref({ x: 0, y: 0 })
const contextMenuEdgeId = ref(null)

// 节点ID计数器
let nodeIdCounter = 1
let edgeIdCounter = 1

// 计算属性
const selectedNodeData = computed(() => {
  if (!selectedNode.value) return null
  return nodes.value.find(n => n.id === selectedNode.value)
})

const selectedEdgeData = computed(() => {
  if (!selectedEdge.value) return null
  return edges.value.find(e => e.id === selectedEdge.value)
})

const panelTitle = computed(() => {
  if (selectedNode.value) {
    const node = selectedNodeData.value
    return `节点属性 - ${node?.node_name || ''}`
  }
  if (selectedEdge.value) {
    return '连线属性'
  }
  return '属性面板'
})

// 初始化工作流
const initWorkflow = async () => {
  if (props.workflowId) {
    const { workflowApi } = await import('@/api')
    try {
      const response = await workflowApi.getDetail(props.workflowId)
      loadWorkflowData(response.data)
    } catch (error) {
      console.error('加载工作流失败:', error)
    }
  } else {
    createDefaultNodes()
  }
}

// 创建默认节点
const createDefaultNodes = () => {
  nodes.value = [
    {
      id: 'start',
      node_key: 'start',
      node_type: 'start',
      node_name: '开始',
      description: '流程开始',
      x: 100,
      y: 200,
      approval_type: 'sequential',
      auto_approve: false,
      is_required: true,
      approvers: []
    },
    {
      id: 'end',
      node_key: 'end',
      node_type: 'end',
      node_name: '结束',
      description: '流程结束',
      x: 500,
      y: 200,
      approval_type: 'sequential',
      auto_approve: false,
      is_required: true,
      approvers: []
    }
  ]
  edges.value = []
  nodeIdCounter = 1
  edgeIdCounter = 1
}

// 加载工作流数据
const loadWorkflowData = (data) => {
  nodes.value = data.nodes?.map(node => ({
    ...node,
    // 使用 node_key 作为 id，因为边的引用使用的是 node_key
    id: node.node_key || node.id,
    approvers: node.approvers || []
  })) || []
  edges.value = data.edges?.map(edge => ({
    ...edge,
    id: `edge_${edge.id}`
  })) || []

  // 计算最大节点 ID（用于生成新节点）
  const maxNodeId = nodes.value.reduce((max, node) => {
    const num = parseInt(node.id.toString().replace(/\D/g, ''))
    return num > max ? num : max
  }, 0)
  nodeIdCounter = maxNodeId + 1

  // 计算最大边 ID（用于生成新边）
  const maxEdgeId = edges.value.reduce((max, edge) => {
    const num = parseInt(edge.id.toString().replace(/\D/g, ''))
    return num > max ? num : max
  }, 0)
  edgeIdCounter = maxEdgeId + 1
}

// 节点操作
// 节点操作
const addNode = (type) => {
  const nodeKey = `node_${nodeIdCounter++}`
  const node = {
    id: nodeKey,           // 使用相同的值
    node_key: nodeKey,     // 使用相同的值
    node_type: type,
    node_name: '新节点',
    description: '',
    x: 300 + Math.random() * 100,
    y: 150 + Math.random() * 100,
    approval_type: 'sequential',
    timeout_hours: 24,
    auto_approve: false,
    is_required: true,
    approvers: []
  }

  nodes.value.push(node)
  selectNode(node.id)
}

const deleteNode = (nodeId) => {
  const index = nodes.value.findIndex(n => n.id === nodeId)
  if (index !== -1) {
    nodes.value.splice(index, 1)
  }

  // 删除相关连线
  edges.value = edges.value.filter(e => e.from_node !== nodeId && e.to_node !== nodeId)

  clearSelection()
  emitChange()
}

// 节点鼠标按下处理
const handleNodeMouseDown = (event, node) => {
  // 如果正在编辑连线
  if (isConnectingEdge.value && editingEdge.value) {
    completeEditEdge(node.id)
    return
  }

  // 如果是开始节点，不能被连线
  if (node.node_type === 'start' && connectingFrom.value && !editingEdge.value) {
    ElMessage.warning('不能连接到开始节点')
    connectingFrom.value = null
    tempEdge.value = ''
    return
  }

  // 如果正在创建连线，且点击的是节点输入点
  if (connectingFrom.value && !editingEdge.value) {
    completeConnect(node)
    return
  }

  // 先选择节点（显示属性面板）
  selectNode(node.id)

  // 检查是否点击在节点区域内（拖拽节点）
  const rect = canvasRef.value.getBoundingClientRect()
  const clickX = event.clientX - rect.left
  const clickY = event.clientY - rect.top

  // 在节点范围内，开始拖拽节点
  if (clickX >= node.x && clickX <= node.x + getNodeWidth(node) &&
      clickY >= node.y && clickY <= node.y + getNodeHeight(node)) {
    draggingNode.value = node
    dragOffset.value = {
      x: event.clientX - rect.left - node.x,
      y: event.clientY - rect.top - node.y
    }
    isDraggingNode.value = true
  }
}

// 连线操作
const startConnect = (event, node) => {
  // 结束节点不能发起连线
  if (node.node_type === 'end') {
    ElMessage.warning('结束节点不能连接其他节点')
    return
  }

  connectingFrom.value = node
  tempEdge.value = ''
}

const completeConnect = (toNode) => {
  if (!connectingFrom.value) return

  // 开始节点不能被连接
  if (toNode.node_type === 'start') {
    ElMessage.warning('不能连接到开始节点')
    connectingFrom.value = null
    tempEdge.value = ''
    return
  }

  // 不能连接到自己
  if (connectingFrom.value.id === toNode.id) {
    ElMessage.warning('不能连接到自己')
    connectingFrom.value = null
    tempEdge.value = ''
    return
  }

  // 检查是否已存在相同连线
  const exists = edges.value.some(
    e => e.from_node === connectingFrom.value.id && e.to_node === toNode.id
  )
  if (exists) {
    ElMessage.warning('连线已存在')
    connectingFrom.value = null
    tempEdge.value = ''
    return
  }

  // 添加连线
  edges.value.push({
    id: `edge_${edgeIdCounter++}`,
    from_node: connectingFrom.value.id,
    to_node: toNode.id,
    workflow_id: props.workflowId,
    condition_expression: ''
  })

  connectingFrom.value = null
  tempEdge.value = ''
  emitChange()
}

// 从输出点完成连线
const completeConnectFromOutput = (event, fromNode) => {
  // 查找鼠标位置下的节点
  const rect = canvasRef.value.getBoundingClientRect()
  const x = event.clientX - rect.left
  const y = event.clientY - rect.top

  // 查找目标节点
  const targetNode = nodes.value.find(node => {
    const nodeX = node.x
    const nodeY = node.y
    const nodeWidth = getNodeWidth(node)
    const nodeHeight = getNodeHeight(node)

    // 检查是否在节点左侧输入点附近（给一些容差）
    const inputX = nodeX
    const inputY = nodeY + nodeHeight / 2
    const distance = Math.sqrt(Math.pow(x - inputX, 2) + Math.pow(y - inputY, 2))

    return distance <= 15 && node.id !== fromNode.id
  })

  if (targetNode) {
    completeConnect(targetNode)
  } else {
    connectingFrom.value = null
    tempEdge.value = ''
  }
}

// 拖拽到输入点完成连线
const completeConnectToInput = (event, toNode) => {
  if (!connectingFrom.value) return

  // 开始节点不能被连接
  if (toNode.node_type === 'start') {
    ElMessage.warning('不能连接到开始节点')
    connectingFrom.value = null
    tempEdge.value = ''
    return
  }

  // 不能连接到自己
  if (connectingFrom.value.id === toNode.id) {
    ElMessage.warning('不能连接到自己')
    connectingFrom.value = null
    tempEdge.value = ''
    return
  }

  // 检查是否已存在相同连线
  const exists = edges.value.some(
    e => e.from_node === connectingFrom.value.id && e.to_node === toNode.id
  )
  if (exists) {
    ElMessage.warning('连线已存在')
    connectingFrom.value = null
    tempEdge.value = ''
    return
  }

  // 添加连线
  edges.value.push({
    id: `edge_${edgeIdCounter++}`,
    from_node: connectingFrom.value.id,
    to_node: toNode.id,
    workflow_id: props.workflowId,
    condition_expression: ''
  })

  connectingFrom.value = null
  tempEdge.value = ''
  emitChange()
}

// 高亮输入连接点
const highlightInput = (nodeId, highlight) => {
  highlightedInput.value = highlight ? nodeId : null
}

// 判断节点是否是连线的端点
const isEdgeEndpoint = (edgeId, nodeId, endpointType) => {
  const edge = edges.value.find(e => e.id === edgeId)
  if (!edge) return false

  if (endpointType === 'source') {
    return edge.from_node === nodeId
  } else if (endpointType === 'target') {
    return edge.to_node === nodeId
  }
  return false
}

// 开始编辑连线
const startEditEdge = (event, edgeId, endpointType) => {
  event.stopPropagation()
  editingEdge.value = edgeId
  editingEndpoint.value = endpointType
  isConnectingEdge.value = true
  tempEdge.value = ''

  const edge = edges.value.find(e => e.id === edgeId)
  if (edge) {
    const node = nodes.value.find(n => n.id === (endpointType === 'source' ? edge.from_node : edge.to_node))
    if (node) {
      connectingFrom.value = node
    }
  }
}

// 完成编辑连线
const completeEditEdge = (targetNodeId) => {
  if (!editingEdge.value || !editingEndpoint.value) return

  const edgeIndex = edges.value.findIndex(e => e.id === editingEdge.value)
  if (edgeIndex === -1) return

  const edge = edges.value[edgeIndex]
  const targetNode = nodes.value.find(n => n.id === targetNodeId)

  if (!targetNode) {
    cancelEditEdge()
    return
  }

  // 验证规则
  if (editingEndpoint.value === 'source') {
    // 不能从结束节点发起连线
    if (targetNode.node_type === 'end') {
      ElMessage.warning('结束节点不能连接其他节点')
      cancelEditEdge()
      return
    }
    // 检查是否已存在相同连线
    const exists = edges.value.some(
      e => e.id !== editingEdge.value && e.from_node === targetNodeId && e.to_node === edge.to_node
    )
    if (exists) {
      ElMessage.warning('连线已存在')
      cancelEditEdge()
      return
    }
    edge.from_node = targetNodeId
  } else {
    // 目标节点不能是开始节点
    if (targetNode.node_type === 'start') {
      ElMessage.warning('不能连接到开始节点')
      cancelEditEdge()
      return
    }
    // 检查是否已存在相同连线
    const exists = edges.value.some(
      e => e.id !== editingEdge.value && e.from_node === edge.from_node && e.to_node === targetNodeId
    )
    if (exists) {
      ElMessage.warning('连线已存在')
      cancelEditEdge()
      return
    }
    edge.to_node = targetNodeId
  }

  cancelEditEdge()
  emitChange()
}

// 取消编辑连线
const cancelEditEdge = () => {
  editingEdge.value = null
  editingEndpoint.value = null
  isConnectingEdge.value = false
  connectingFrom.value = null
  tempEdge.value = ''
}

// 画布鼠标移动
const handleCanvasMouseMove = (event) => {
  const rect = canvasRef.value.getBoundingClientRect()

  // 拖拽节点
  if (isDraggingNode.value && draggingNode.value) {
    draggingNode.value.x = event.clientX - rect.left - dragOffset.value.x
    draggingNode.value.y = event.clientY - rect.top - dragOffset.value.y
  }

  // 创建连线
  if (connectingFrom.value) {
    const fromX = connectingFrom.value.x + getNodeWidth(connectingFrom.value)
    const fromY = connectingFrom.value.y + getNodeHeight(connectingFrom.value) / 2
    const toX = event.clientX - rect.left
    const toY = event.clientY - rect.top

    // 使用贝塞尔曲线
    const midX = (fromX + toX) / 2
    tempEdge.value = `M ${fromX} ${fromY} C ${midX} ${fromY}, ${midX} ${toY}, ${toX} ${toY}`
  }

  // 重新计算临时连线（拖拽边时）
  if (isConnectingEdge.value && selectedEdge.value) {
    // TODO: 实现拖拽修改连线目标的功能
  }
}

// 画布鼠标释放
const handleCanvasMouseUp = (event) => {
  // 如果正在编辑连线，点击空白处取消
  if (isConnectingEdge.value && editingEdge.value) {
    cancelEditEdge()
    return
  }

  if (isConnectingEdge.value) {
    isConnectingEdge.value = false
    emitChange()
  }

  if (isDraggingNode.value) {
    isDraggingNode.value = false
    draggingNode.value = null
    emitChange()
  }

  // 连线创建
  if (connectingFrom.value && !editingEdge.value) {
    // completeConnectFromOutput 已经处理了
    // 如果还在连接状态但没有完成，取消连接
    connectingFrom.value = null
    tempEdge.value = ''
  }
}

// 画布鼠标按下
const handleCanvasMouseDown = (event) => {
  // 点击画布空白区域，清除选择
  if (event.target === canvasRef.value) {
    clearSelection()
  }
}

// 选择
const selectNode = (nodeId) => {
  selectedNode.value = nodeId
  selectedEdge.value = null

  const node = selectedNodeData.value
  if (node) {
    Object.assign(nodeForm, {
      node_key: node.node_key,
      node_name: node.node_name,
      node_type: node.node_type,
      description: node.description || '',
      approval_type: node.approval_type || 'sequential',
      timeout_hours: node.timeout_hours || 24,
      auto_approve: node.auto_approve || false,
      is_required: node.is_required !== false,
      approvers: node.approvers ? [...node.approvers] : []
    })
  }
}

const selectEdge = (edgeId) => {
  selectedEdge.value = edgeId
  selectedNode.value = null

  const edge = selectedEdgeData.value
  if (edge) {
    edgeCondition.value = edge.condition_expression || ''
  }
}

const clearSelection = () => {
  selectedNode.value = null
  selectedEdge.value = null
}

// 右键菜单功能
const showEdgeContextMenu = (event, edgeId) => {
  contextMenuEdgeId.value = edgeId
  contextMenuPosition.value = {
    x: event.clientX,
    y: event.clientY
  }
  contextMenuVisible.value = true
}

const hideContextMenu = () => {
  contextMenuVisible.value = false
  contextMenuEdgeId.value = null
}

const deleteEdgeFromContext = () => {
  if (contextMenuEdgeId.value) {
    deleteEdge(contextMenuEdgeId.value)
  }
  hideContextMenu()
}

// 更新节点数据
const updateNodeData = () => {
  const node = selectedNodeData.value
  if (node) {
    Object.assign(node, {
      node_name: nodeForm.node_name,
      node_type: nodeForm.node_type,
      description: nodeForm.description,
      approval_type: nodeForm.approval_type,
      timeout_hours: nodeForm.timeout_hours,
      auto_approve: nodeForm.auto_approve,
      is_required: nodeForm.is_required,
      approvers: nodeForm.approvers
    })
    emitChange()
  }
}

// 更新边条件
const updateEdgeCondition = () => {
  const edge = selectedEdgeData.value
  if (edge) {
    edge.condition_expression = edgeCondition.value
    emitChange()
  }
}

// 删除节点
const deleteSelectedNode = () => {
  if (selectedNode.value) {
    const node = selectedNodeData.value
    if (node) {
      if (node.node_type === 'start' || node.node_type === 'end') {
        ElMessage.warning('开始和结束节点不能删除')
        return
      }
      deleteNode(selectedNode.value)
    }
  }
}

// 删除边
const deleteEdge = (edgeId) => {
  if (edgeId) {
    const index = edges.value.findIndex(e => e.id === edgeId)
    if (index !== -1) {
      edges.value.splice(index, 1)
      emitChange()
    }
    if (selectedEdge.value === edgeId) {
      selectedEdge.value = null
    }
  }
}

const deleteSelectedEdge = () => {
  if (selectedEdge.value) {
    deleteEdge(selectedEdge.value)
  }
}

// 审批人操作
const showAddApproverDialog = () => {
  Object.assign(newApprover, {
    approver_type: 'role',
    approver_id: null,
    approver_name: '',
    sequence: (nodeForm.approvers?.length || 0)
  })
  addApproverDialogVisible.value = true
}

const addApprover = () => {
  if (!newApprover.approver_id) {
    ElMessage.warning('请选择审批人')
    return
  }

  const approverList = newApprover.approver_type === 'user' ? userList.value : roleList.value
  const approver = approverList.find(u => u.id === newApprover.approver_id)

  if (!nodeForm.approvers) {
    nodeForm.approvers = []
  }

  nodeForm.approvers.push({
    approver_type: newApprover.approver_type,
    approver_id: newApprover.approver_id,
    approver_name: approver?.name || approver?.username || approver?.real_name || '',
    sequence: newApprover.sequence
  })

  addApproverDialogVisible.value = false
  updateNodeData()
}

const removeApprover = (index) => {
  nodeForm.approvers.splice(index, 1)
  updateNodeData()
}

const getApproverName = (approver) => {
  if (approver.approver_type === 'user') {
    const user = userList.value.find(u => u.id === approver.approver_id)
    return user ? `${user.username} - ${user.real_name || user.username}` : approver.approver_name
  } else {
    const role = roleList.value.find(r => r.id === approver.approver_id)
    return role ? role.name : approver.approver_name
  }
}

// 获取选中边的源节点名称
const getSelectedEdgeFromNode = () => {
  if (!selectedEdgeData.value) return '-'
  const node = nodes.value.find(n => n.id === selectedEdgeData.value.from_node)
  return node?.node_name || '-'
}

// 获取选中边的目标节点名称
const getSelectedEdgeToNode = () => {
  if (!selectedEdgeData.value) return '-'
  const node = nodes.value.find(n => n.id === selectedEdgeData.value.to_node)
  return node?.node_name || '-'
}

// 获取边的中点（用于显示删除按钮）
const getEdgeMidPoint = (edge) => {
  const fromNode = nodes.value.find(n => n.id === edge.from_node)
  const toNode = nodes.value.find(n => n.id === edge.to_node)

  if (!fromNode || !toNode) return { x: 0, y: 0 }

  const fromX = fromNode.x + getNodeWidth(fromNode)
  const fromY = fromNode.y + getNodeHeight(fromNode) / 2
  const toX = toNode.x
  const toY = toNode.y + getNodeHeight(toNode) / 2

  return {
    x: (fromX + toX) / 2,
    y: (fromY + toY) / 2
  }
}

// 工具函数
const getNodeWidth = (node) => {
  return node.node_type === 'start' || node.node_type === 'end' ? 100 : 140
}

const getNodeHeight = (node) => {
  return 70
}

const getNodeIcon = (type) => {
  const icons = {
    start: '▶',
    end: '■',
    approval: '✓',
    parallel: '∥',
    merge: '⇄'
  }
  return icons[type] || '?'
}

const getNodeTypeText = (type) => {
  const texts = {
    start: '开始',
    end: '结束',
    approval: '审批',
    parallel: '并行',
    merge: '合并'
  }
  return texts[type] || type
}

const getEdgePath = (edge) => {
  const fromNode = nodes.value.find(n => n.id === edge.from_node)
  const toNode = nodes.value.find(n => n.id === edge.to_node)

  if (!fromNode || !toNode) return ''

  const fromX = fromNode.x + getNodeWidth(fromNode)
  const fromY = fromNode.y + getNodeHeight(fromNode) / 2
  const toX = toNode.x
  const toY = toNode.y + getNodeHeight(toNode) / 2

  const midX = (fromX + toX) / 2
  return `M ${fromX} ${fromY} C ${midX} ${fromY}, ${midX} ${toY}, ${toX} ${toY}`
}

// 保存工作流
const handleSave = () => {
  if (nodes.value.length < 2) {
    ElMessage.error('工作流至少需要开始和结束节点')
    return
  }

  const hasStart = nodes.value.some(n => n.node_type === 'start')
  const hasEnd = nodes.value.some(n => n.node_type === 'end')

  if (!hasStart || !hasEnd) {
    ElMessage.error('工作流必须包含开始和结束节点')
    return
  }

  // 验证必须有从 start 节点出发的边
  const startNode = nodes.value.find(n => n.node_type === 'start')
  const hasStartEdge = edges.value.some(e => e.from_node === startNode?.id)

  if (!hasStartEdge) {
    ElMessage.error('开始节点必须连接到下一个节点（请拖拽连接线）')
    return
  }

  // 验证所有审批节点都必须有连接
  const approvalNodes = nodes.value.filter(n => n.node_type === 'approval')
  for (const node of approvalNodes) {
    const hasIncomingEdge = edges.value.some(e => e.to_node === node.id)
    const hasOutgoingEdge = edges.value.some(e => e.from_node === node.id)

    if (!hasIncomingEdge) {
      ElMessage.error(`节点"${node.node_name}"没有输入连接（必须有前序节点连接到它）`)
      return
    }

    if (!hasOutgoingEdge) {
      ElMessage.error(`节点"${node.node_name}"没有输出连接（必须连接到下一个节点）`)
      return
    }
  }

  // 验证结束节点必须有输入连接
  const endNode = nodes.value.find(n => n.node_type === 'end')
  const hasEndEdge = edges.value.some(e => e.to_node === endNode?.id)

  if (!hasEndEdge) {
    ElMessage.error('必须有一个节点连接到结束节点')
    return
  }

  const workflowData = {
    nodes: nodes.value.map(node => ({
      node_key: node.node_key,
      node_type: node.node_type,
      node_name: node.node_name,
      description: node.description,
      approval_type: node.approval_type,
      timeout_hours: node.timeout_hours,
      auto_approve: node.auto_approve,
      is_required: node.is_required,
      x: Math.round(node.x),
      y: Math.round(node.y),
      approvers: node.approvers || []
    })),
    edges: edges.value.map(edge => ({
      from_node: edge.from_node,
      to_node: edge.to_node,
      condition_expression: edge.condition_expression
    }))
  }

  emit('save', workflowData)
}

// 重置
const handleReset = () => {
  ElMessageBox.confirm('确定要重置工作流吗？所有未保存的更改将丢失。', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    initWorkflow()
    clearSelection()
  }).catch(() => {})
}

const emitChange = () => {
  emit('change', {
    nodes: nodes.value,
    edges: edges.value
  })
}

// 切换全屏
const toggleFullscreen = () => {
  isFullscreen.value = !isFullscreen.value
}

// 加载用户和角色列表
const loadUsersAndRoles = async () => {
  try {
    const { userApi, roleApi } = await import('@/api')
    const [usersRes, rolesRes] = await Promise.all([
      userApi.getList({ pageSize: 1000 }),
      roleApi.getList()
    ])
    userList.value = usersRes.data || []
    roleList.value = rolesRes.data || []
  } catch (error) {
    console.error('加载用户和角色失败:', error)
  }
}

// 监听 workflowId 变化，重新初始化工作流
watch(() => props.workflowId, (newId, oldId) => {
  // 当 workflowId 变化时，重新初始化
  if (newId !== oldId) {
    clearSelection()
    initWorkflow()
  }
})

onMounted(() => {
  initWorkflow()
  loadUsersAndRoles()
})

defineExpose({
  handleSave,
  getWorkflowData: () => ({
    nodes: nodes.value,
    edges: edges.value
  }),
  setWorkflowData: loadWorkflowData
})
</script>

<style scoped>
.workflow-editor {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #f5f5f5;
  position: relative;
}

.workflow-editor.fullscreen {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 10000;
}

.editor-toolbar {
  padding: 12px 16px;
  background: white;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  align-items: center;
  gap: 12px;
}

.toolbar-spacer {
  flex: 1;
}

.editor-main {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.editor-container {
  flex: 1;
  overflow: auto;
  position: relative;
}

.workflow-canvas {
  width: 100%;
  height: 100%;
  min-width: 1200px;
  min-height: 600px;
  background: white;
  cursor: default;
}

.property-panel {
  width: 360px;
  background: white;
  border-left: 1px solid #e0e0e0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.panel-header {
  padding: 16px;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.panel-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 500;
  color: #303133;
}

.panel-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.node {
  cursor: move;
  transition: filter 0.2s;
}

.node:hover {
  filter: brightness(0.95);
}

.node-selected .node-shape {
  stroke: #409eff;
  stroke-width: 2;
}

.node-shape {
  fill: white;
  stroke: #d0d0d0;
  stroke-width: 1;
  transition: all 0.2s;
}

.shape-start {
  fill: #e8f5e9;
  stroke: #4caf50;
}

.shape-end {
  fill: #ffebee;
  stroke: #f44336;
}

.shape-approval {
  fill: #e3f2fd;
  stroke: #2196f3;
}

.shape-parallel {
  fill: #fff3e0;
  stroke: #ff9800;
}

.shape-merge {
  fill: #f3e5f5;
  stroke: #9c27b0;
}

.node-icon {
  font-size: 16px;
  font-weight: bold;
  fill: #333;
  pointer-events: none;
}

.node-name {
  font-size: 14px;
  font-weight: 500;
  fill: #333;
  pointer-events: none;
}

.node-type {
  font-size: 12px;
  fill: #666;
  pointer-events: none;
}

.connector {
  fill: #409eff;
  cursor: crosshair;
  transition: r 0.2s, fill 0.2s;
}

.connector-in {
  fill: #67c23a;
}

.connector-in:hover {
  r: 8;
  fill: #85ce61;
}

.connector-out {
  fill: #409eff;
}

.connector-out:hover {
  r: 8;
  fill: #66b1ff;
}

.connector-drag {
  fill: #e6a23c;
  stroke: #fff;
  stroke-width: 2;
}

.connector-drag:hover {
  r: 10;
  fill: #f0b857;
}

.edge-group {
  cursor: pointer;
}

.edge {
  fill: none;
  stroke: #409eff;
  stroke-width: 2;
  transition: stroke-width 0.2s;
}

.edge-selected {
  stroke: #f56c6c;
  stroke-width: 3;
}

.edge-actions {
  /* Allow pointer events on children but not on the group itself */
  pointer-events: none;
}

.edge-delete-btn {
  cursor: pointer;
  pointer-events: auto;
  transition: r 0.2s;
}

.edge-delete-btn:hover {
  r: 12;
}

.temp-edge {
  stroke-dasharray: 5, 5;
  pointer-events: none;
}

.approvers-list {
  margin-top: 8px;
}

.empty-approvers {
  padding: 20px 0;
}

.approver-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  background: #f5f7fa;
  border-radius: 4px;
  margin-bottom: 8px;
}

.approver-info {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.approver-name {
  font-size: 13px;
  color: #303133;
}

.edge-info {
  font-size: 13px;
  color: #606266;
  padding: 8px 12px;
  background: #f5f7fa;
  border-radius: 4px;
}

.edge-actions {
  margin-top: 16px;
  text-align: right;
}

.node-actions {
  margin-top: 16px;
}

/* 浮动半透明操作提示 */
.tips-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  backdrop-filter: blur(2px);
}

.tips-content {
  background: rgba(255, 255, 255, 0.95);
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  max-width: 400px;
  width: 90%;
  max-height: 80vh;
  overflow-y: auto;
  backdrop-filter: blur(10px);
}

.tips-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid #e0e0e0;
  font-size: 16px;
  color: #303133;
}

.close-btn {
  cursor: pointer;
  transition: all 0.2s;
}

.close-btn:hover {
  color: #f56c6c;
}

.tips-list {
  margin: 16px 20px;
  padding-left: 20px;
  line-height: 1.8;
  color: #606266;
}

.tips-list > li {
  margin-bottom: 12px;
}

.tips-sub-list {
  margin: 8px 0 4px 20px;
  padding-left: 20px;
}

.tips-sub-list li {
  margin-bottom: 4px;
  font-size: 13px;
  color: #909399;
}

.tips-footer {
  padding: 12px 20px;
  border-top: 1px solid #e0e0e0;
  text-align: right;
  background: rgba(245, 247, 250, 0.8);
  border-radius: 0 0 8px 8px;
}

.tips-toggle {
  padding: 8px 16px;
  background: #fff;
  border-top: 1px solid #e0e0e0;
  text-align: center;
}

/* 连线右键菜单 */
.edge-context-menu {
  position: fixed;
  background: white;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  z-index: 10001;
  min-width: 120px;
  padding: 4px 0;
}

.context-menu-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  cursor: pointer;
  font-size: 14px;
  color: #303133;
  transition: background 0.2s;
}

.context-menu-item:hover {
  background: #f5f7fa;
  color: #f56c6c;
}

.context-menu-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 10000;
  background: transparent;
}
</style>
