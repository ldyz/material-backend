<template>
  <div class="gantt-chart" :class="{ 'is-fullscreen': isFullscreen }" @keydown="handleKeydown" tabindex="0" ref="ganttChartRef">
    <!-- 工具栏 -->
    <GanttToolbar
      :current-zoom-label="currentZoomLabel"
      :current-period-text="currentPeriodText"
      :show-dependencies="showDependencies"
      :show-critical-path="showCriticalPath"
      :show-baseline="showBaseline"
      :show-task-list="showTaskList"
      :group-mode="groupMode"
      :search-keyword="searchKeyword"
      :is-fullscreen="isFullscreen"
      :is-saving="saving"
      :has-unsaved-changes="hasUnsavedChanges"
      :can-undo="canUndo"
      :can-redo="canRedo"
      :selected-count="selectedTaskIds.size"
      :undo-count="undoCount"
      :command-history="commandHistory"
      :timeline-format="state.timelineFormat"
      :date-display-format="state.dateDisplayFormat"
      :pan-mode="state.panMode"
      :chart-view-mode="chartViewMode"
      :show-network-time-params="networkShowTimeParams"
      :show-network-task-names="networkShowTaskNames"
      :show-network-slack="networkShowSlack"
      :network-layout-mode="networkLayoutMode"
      @back-to-list="$emit('back-to-list')"
      @navigate-date="navigateDate"
      @go-today="goToToday"
      @zoom-in="zoomIn"
      @zoom-out="zoomOut"
      @zoom-reset="zoomReset"
      @toggle-dependencies="toggleDependencies"
      @toggle-critical-path="toggleCriticalPath"
      @toggle-baseline="toggleBaseline"
      @toggle-task-list="toggleTaskList"
      @open-resource-management="openResourceManagement"
      @group-change="handleGroupChange"
      @search="searchKeyword = $event"
      @export-png="handleExportPNG"
      @export-pdf="handleExportPDF"
      @auto-fit="autoFitContainer"
      @refresh="handleRefresh"
      @toggle-fullscreen="toggleFullscreen"
      @save-all="handleSaveAll"
      @undo="handleUndo"
      @redo="handleRedo"
      @undo-to="handleUndoTo"
      @clear-history="handleClearHistory"
      @toggle-template="templateDialogVisible = true"
      @toggle-bulk-edit="handleToggleBulkEdit"
      @timeline-format-change="handleTimelineFormatChange"
      @date-format-change="handleDateFormatChange"
      @toggle-pan-mode="togglePanMode"
      @toggle-select-mode="toggleSelectMode"
      @toggle-view-mode="handleToggleViewMode"
      @add-task="handleCreateTask"
      @toggle-network-time-params="networkShowTimeParams = !networkShowTimeParams"
      @toggle-network-task-names="networkShowTaskNames = !networkShowTaskNames"
      @toggle-network-slack="networkShowSlack = !networkShowSlack"
      @network-layout-change="networkLayoutMode = $event"
      @calculate-critical-path="handleCalculateCriticalPath"
      @analyze-node-properties="handleAnalyzeNodeProperties"
      @check-path-optimization="handleCheckPathOptimization"
      @validate-rules="handleValidateRules"
      @export-analysis-report="handleExportAnalysisReport"
    />

    <!-- 统计信息 -->
    <GanttStats v-if="chartViewMode === 'gantt'" :stats="taskStats" />
    <GanttStats v-else-if="chartViewMode === 'network'" :stats="networkStatsDisplay" />

    <!-- 甘特图容器 -->
    <div
      class="gantt-container"
      :class="{ 'is-resizing': isResizing, 'is-pan-mode': panMode }"
      ref="ganttContainer"
      :style="containerStyle"
      @mousedown="handlePanStart"
    >
      <!-- 调整大小手柄 - 右下角 -->
      <div
        class="resize-handle resize-handle-corner"
        @mousedown="handleResizeStart"
        title="拖动调整大小"
      >
        <svg width="16" height="16" viewBox="0 0 16 16">
          <path d="M12 4 L12 12 L4 12" stroke="#909399" stroke-width="2" fill="none"/>
          <path d="M9 4 L9 9 L4 9" stroke="#909399" stroke-width="2" fill="none"/>
        </svg>
      </div>

      <!-- 调整大小手柄 - 右边缘 -->
      <div
        class="resize-handle resize-handle-right"
        @mousedown="handleResizeStartRight"
      ></div>

      <!-- 调整大小手柄 - 底边缘 -->
      <div
        class="resize-handle resize-handle-bottom"
        @mousedown="handleResizeStartBottom"
      ></div>

      <!-- 表头容器（固定在顶部） -->
      <div class="gantt-header-container">
        <GanttHeader
          :timeline-format="state.timelineFormat"
          :date-display-format="state.dateDisplayFormat"
          :timeline-days="timelineDays"
          :timeline-weeks="timelineWeeks"
          :timeline-months="timelineMonths"
          :timeline-header-months="timelineHeaderMonths"
          :timeline-quarters="timelineQuarters"
          :day-width="dayWidth"
          :today-position="todayPosition"
          :pan-offset="chartViewMode === 'network' ? networkPanX : -timelinePanX"
        />
      </div>

      <!-- 主体区域：任务表格 + 视图区域 -->
      <div class="gantt-body" ref="ganttBodyRef">
        <!-- 左侧：任务表格（始终显示） -->
        <TaskTable
          :is-collapsed="!showTaskList"
          :tasks="filteredTasks"
          :grouped-tasks="groupedTasks"
          :selected-task-id="selectedTask?.id"
          :row-height="rowHeight"
          :show-critical-path="showCriticalPath"
          :group-mode="groupMode"
          :collapsed-groups="collapsedGroups"
          :empty-description="emptyDescription"
          @row-click="handleRowClick"
          @row-dblclick="handleRowDblClick"
          @toggle-group="toggleGroup"
          @context-menu="handleContextMenu"
          @cell-edit="handleCellEdit"
          @task-dragged="handleTaskDragged"
        />

        <!-- 右侧：视图区域（根据视图模式切换） -->
        <!-- 甘特图时间轴视图 -->
        <TaskTimeline
          v-if="chartViewMode === 'gantt'"
          :show-task-names="!showTaskList"
          :tasks="filteredTasks"
          :raw-tasks="formattedTasks"
          :selected-task-id="selectedTask?.id"
          :view-mode="viewMode"
          :timeline-days="timelineDays"
          :day-width="dayWidth"
          :row-height="rowHeight"
          :task-height="taskHeight"
          :show-dependencies="showDependencies"
          :show-critical-path="showCriticalPath"
          :show-baseline="showBaseline"
          :visible-dependencies="visibleDependencies"
          :is-dragging="isDragging"
          :dragged-task="draggedTask"
          :tooltip-visible="tooltipVisible"
          :tooltip-position="tooltipPosition"
          :tooltip-text="tooltipText"
          :preview-task="previewTask"
          :drag-mode="dragMode"
          :today-position="todayPosition"
          :arrow-marker-id="arrowMarkerId"
          :arrow-color="arrowColor"
          :empty-description="emptyDescription"
          :is-creating-dependency="isCreatingDependency"
          :source-task-id="dependencySourceTask?.id"
          :temp-line-end="tempLineEnd"
          @task-click="handleTaskClick"
          @task-dblclick="handleTaskDblClick"
          @task-mousedown="handleTaskMouseDown"
          @dependency-create="handleDependencyCreate"
          @dependency-contextmenu="handleDependencyContextMenu"
          @mousemove="handleTimelineMouseMove"
          @task-dragged="handleTaskDragged"
          @context-menu="handleContextMenu"
          @zoom-change="handleTimelineZoom"
          @scroll="handleTimelineScroll"
          ref="taskTimelineRef"
        />

        <!-- 网络图视图 -->
        <NetworkView
          ref="networkViewRef"
          v-if="chartViewMode === 'network'"
          :tasks="formattedTasks"
          :dependencies="formattedTasks.flatMap(t =>
            (t.predecessors || []).map(predId => ({
              task_id: t.id,
              depends_on: predId,
              type: 'FS',
              is_critical: t.is_critical
            }))
          )"
          :task-index-map="networkTaskIndexMap"
          :timeline-start="timelineDays[0]?.date || ''"
          :row-height="rowHeight"
          :day-width="dayWidth"
          :align-with-task-list="true"
          :show-critical-path="showCriticalPath"
          :show-task-names="networkShowTaskNames"
          :show-time-params="networkShowTimeParams"
          :show-slack="networkShowSlack"
          :show-duration="true"
          :tool-mode="networkToolMode"
          :svg-width="Math.max(2000, timelineWidth)"
          :svg-height="networkHeight"
          :node-radius="18"
          @node-click="handleNetworkNodeClick"
          @task-click="handleNetworkTaskClick"
          @task-dblclick="handleNetworkTaskDblClick"
          @zoom-change="handleNetworkZoom"
          @pan-change="handleNetworkPan"
          @node-contextmenu="handleNetworkNodeContextMenu"
          @task-contextmenu="handleNetworkTaskContextMenu"
          @node-time-change="handleNetworkNodeTimeChange"
          @node-dependency-create="handleNetworkDependencyCreate"
          @node-merge="handleNetworkNodeMerge"
          @node-split="handleNetworkNodeSplit"
          @task-dependency-edit="handleNetworkTaskDependencyEdit"
        />
      </div>
    </div>

    <!-- 图例 -->
    <GanttLegend
      :show-critical-path="showCriticalPath"
      :show-baseline="showBaseline"
    />

    <!-- 右键菜单 (放在全屏容器外) -->
    <teleport to="body">
      <GanttContextMenu
        v-model:visible="contextMenuVisible"
        :task="contextMenuTask"
        :position="contextMenuPosition"
        :all-tasks="formattedTasks"
        @add-subtask="handleAddSubtask"
        @convert-milestone="handleConvertToMilestone"
        @add-dependency="handleAddDependency"
        @view-dependencies="handleViewDependencies"
        @allocate-resources="handleAllocateResources"
        @create-task="handleCreateTask"
        @edit="handleContextMenuEdit"
        @duplicate="handleContextMenuDuplicate"
        @delete="handleContextMenuDelete"
        @move-up="handleMoveTaskUp"
        @move-down="handleMoveTaskDown"
        @convert-to-independent="handleConvertToIndependent"
        ref="contextMenuRef"
      />

      <!-- 网络图节点右键菜单 -->
      <NetworkNodeContextMenu
        v-model:visible="nodeContextMenuVisible"
        :node="nodeContextMenuNode"
        :selected-nodes="nodeContextMenuSelectedNodes"
        :can-merge="nodeContextMenuCanMerge"
        :position="nodeContextMenuPosition"
        @merge-nodes="(...args) => handleNodeContextMenuCommand('merge-nodes', args[0])"
        @split-node="(...args) => handleNodeContextMenuCommand('split-node', args[0])"
        @edit-tasks="(...args) => handleNodeContextMenuCommand('edit-tasks', args[0])"
      />

      <!-- 依赖关系右键菜单 -->
      <el-dropdown
        :virtual-ref="dependencyContextMenuVisible ? {
          getBoundingClientRect: () => new DOMRect(
            dependencyContextMenuPosition.x,
            dependencyContextMenuPosition.y,
            0,
            0
          )
        } : undefined"
        virtual-triggering
        trigger="contextmenu"
        @command="handleDependencyMenuCommand"
        @visible-change="handleDependencyMenuVisibleChange"
      >
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="delete-dependency" :icon="Delete">
              删除依赖关系
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </teleport>

    <!-- 资源分配对话框 -->
    <ResourceAllocationDialog
      v-model="resourceDialogVisible"
      :task-id="currentTaskForResource?.id"
      :project-id="projectId"
      @saved="handleResourceSaved"
    />

    <!-- 资源管理对话框 -->
    <ResourceManagementDialog
      v-model="resourceManagementDialogVisible"
      :project-id="projectId"
      @refresh="handleResourceRefresh"
    />

    <!-- 任务编辑对话框 -->
    <TaskEditDialog
      v-model:visible="editDialogVisible"
      :editing-task="editingTask"
      :saving="saving"
      :all-tasks="formattedTasks"
      :resource-library="resourceLibrary"
      :project-id="projectId"
      @save="handleSaveTask"
      ref="taskEditDialogRef"
    />

    <!-- 任务模板对话框 -->
    <TaskTemplatesDialog
      v-model="templateDialogVisible"
      :project-id="projectId"
      @created="handleTaskCreated"
    />

    <!-- 批量编辑对话框 -->
    <BulkEditDialog
      v-model="bulkEditDialogVisible"
      :tasks="selectedTasks"
      :project-id="projectId"
      @updated="handleBulkUpdate"
    />

    <!-- 任务详情侧边栏 -->
    <TaskDetailDrawer
      v-model:visible="taskDetailVisible"
      :task="selectedTask"
      @edit="handleEditTaskFromDrawer"
      @duplicate="handleDuplicateTask"
      @delete="handleDeleteTask"
    />

    <!-- 资源视图（已禁用） -->

    <!-- 全屏状态栏 -->
    <GanttStatusBar
      ref="statusBarRef"
      :visible="isFullscreen"
      :project-name="projectName"
      :stats="taskStats"
      :is-dragging="isDragging"
      :is-saving="saving"
    />
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ZoomIn, ZoomOut, Delete } from '@element-plus/icons-vue'
import { progressApi } from '@/api'
import { useGanttDrag } from '@/composables/useGanttDrag'
import {
  isMilestone,
  normalizePredecessors,
  getPredecessorSignature
} from '@/utils/ganttHelpers'
import { formatDate, diffDays } from '@/utils/dateFormat'
import eventBus, { GanttEvents } from '@/utils/eventBus'
import { ganttStore } from '@/stores/ganttStore'
import { useUndoRedoStore } from '@/stores/undoRedoStore'
import {
  CreateTaskCommand,
  UpdateTaskCommand,
  DeleteTaskCommand,
  MoveTaskCommand,
  BatchUpdateTasksCommand,
  DuplicateTaskCommand,
  ConvertToMilestoneCommand
} from '@/commands/taskCommands'

// Import new subcomponents
import GanttToolbar from './GanttToolbar.vue'
import GanttStats from './GanttStats.vue'
import GanttHeader from './GanttHeader.vue'
import TaskList from './TaskList.vue'
import TaskTable from './TaskTable.vue'
import TaskTimeline from './TaskTimeline.vue'
import NetworkView from './NetworkView.vue'
import GanttLegend from './GanttLegend.vue'
import GanttContextMenu from './GanttContextMenu.vue'
import NetworkNodeContextMenu from './NetworkNodeContextMenu.vue'
import TaskDetailDrawer from './TaskDetailDrawer.vue'
import TaskEditDialog from './TaskEditDialog.vue'
import ResourceAllocationDialog from './ResourceAllocationDialog.vue'
import ResourceManagementDialog from './ResourceManagementDialog.vue'
import GanttStatusBar from './GanttStatusBar.vue'
import TaskTemplatesDialog from '@/components/gantt/dialogs/TaskTemplatesDialog.vue'
import BulkEditDialog from '@/components/gantt/dialogs/BulkEditDialog.vue'

const props = defineProps({
  projectId: {
    type: [Number, String],
    required: true
  },
  projectName: {
    type: String,
    default: '未命名项目'
  },
  scheduleData: {
    type: Object,
    default: () => ({})
  },
  dataVersion: {
    type: Number,
    default: 0
  }
})

const emit = defineEmits(['task-updated', 'task-selected', 'back-to-list'])

// ==================== 初始化 Store ====================
const store = ganttStore
const { state, getters, actions } = store

// 设置项目信息
actions.setProject(props.projectId, props.projectName)

// ==================== Undo/Redo Store ====================
const undoRedoStore = useUndoRedoStore()

// 初始化项目上下文
undoRedoStore.clear()

// 多选状态
const selectedTaskIds = ref(new Set())

// ==================== 组件引用 ====================
const ganttChartRef = ref(null)
const ganttContainer = ref(null)
const ganttBodyRef = ref(null)
const taskTimelineRef = ref(null)
const networkViewRef = ref(null)
const statusBarRef = ref(null)
const contextMenuRef = ref(null)

// ==================== 清理引用 ====================
let resizeObserver = null
let resizeTimer = null
let unsubscribeDependencyError = null

// 容器调整大小状态（保留在组件中）
const resizeDirection = ref(null)
const resizeStart = ref({ x: 0, y: 0, width: 0, height: 0 })

// 平移状态
const panStart = ref({ x: 0, scrollLeft: 0 })
const isPanning = ref(false)

// 使用 store 中的状态（通过解构简化访问）
const loading = computed(() => state.loading)
const hasUnsavedChanges = computed(() => state.hasUnsavedChanges)
const saving = computed(() => state.saving)
const isFullscreen = computed(() => state.isFullscreen)
const viewMode = computed(() => state.viewMode)
const groupMode = computed(() => state.groupMode)
const collapsedGroups = computed(() => state.collapsedGroups)
const dayWidth = computed(() => state.dayWidth)
const rowHeight = computed(() => state.rowHeight)
const taskHeight = computed(() => state.taskHeight)
const showDependencies = computed(() => state.showDependencies)
const showCriticalPath = computed(() => state.showCriticalPath)
const showBaseline = computed(() => state.showBaseline)
const showTaskList = computed(() => state.showTaskList)
const panMode = computed(() => state.panMode)
const searchKeyword = computed(() => state.searchKeyword)
const taskDetailVisible = computed({
  get: () => state.taskDetailVisible,
  set: (val) => actions.closeTaskDetail()
})
const selectedTask = computed(() => getters.selectedTask)
const tempLineEnd = computed(() => state.tempLineEnd)
const isResizing = computed(() => state.isResizing)
const containerSize = computed(() => state.containerSize)
const editDialogVisible = computed({
  get: () => state.editDialogVisible,
  set: (val) => { if (!val) actions.closeEditDialog() }
})
const editingTask = computed(() => state.editingTask)
const contextMenuVisible = computed({
  get: () => state.contextMenuVisible,
  set: (val) => { if (!val) actions.hideContextMenu() }
})
const contextMenuTask = computed(() => state.contextMenuTask)
const contextMenuPosition = computed(() => state.contextMenuPosition)

// 依赖关系右键菜单状态（独立于任务右键菜单）
const dependencyContextMenuVisible = ref(false)
const dependencyContextMenuDependency = ref(null)
const dependencyContextMenuPosition = ref({ x: 0, y: 0 })

// 网络图节点右键菜单状态
const nodeContextMenuVisible = ref(false)
const nodeContextMenuNode = ref(null)
const nodeContextMenuSelectedNodes = ref([])
const nodeContextMenuCanMerge = ref(false)
const nodeContextMenuPosition = ref({ x: 0, y: 0 })

const resourceDialogVisible = ref(false)
const currentTaskForResource = ref(null)
const resourceManagementDialogVisible = ref(false)
const resourceLibrary = computed(() => state.resourceLibrary)

// ==================== 对话框状态 ====================
const templateDialogVisible = ref(false)
const bulkEditDialogVisible = ref(false)

// ==================== 视图模式 ====================
const chartViewMode = ref('gantt') // 'gantt' or 'network'

// 网络图视图状态
// 网络图工具模式与全局平移模式同步（pan时为移动，否则为选择）
const networkToolMode = computed(() => panMode.value ? 'pan' : 'select')
const networkShowTimeParams = ref(false)
const networkShowTaskNames = ref(true)
const networkShowSlack = ref(false)
const networkLayoutMode = ref('auto')
const networkPanX = ref(0) // 网络图水平平移偏移量
const timelinePanX = ref(0) // 甘特图滚动偏移量

// 网络图统计信息
const networkStats = computed(() => {
  const tasks = formattedTasks.value || []
  const activities = tasks.length
  const criticalActivities = tasks.filter(t => t.is_critical).length

  // 计算节点数（去重的任务ID）
  const nodes = new Set()
  tasks.forEach(t => {
    nodes.add(t.id)
    if (t.predecessors) {
      t.predecessors.forEach(p => nodes.add(p))
    }
  })

  // 计算总工期
  const totalDuration = tasks.length > 0
    ? Math.max(...tasks.map(t => new Date(t.end).getTime())) -
      Math.min(...tasks.map(t => new Date(t.start).getTime()))
    : 0
  const totalDays = totalDuration / (1000 * 60 * 60 * 24)

  return {
    nodes: nodes.size,
    activities,
    criticalActivities,
    totalDuration: totalDays
  }
})

// 网络图统计信息显示格式（适配 GanttStats 组件）
const networkStatsDisplay = computed(() => {
  const stats = networkStats.value
  const completed = formattedTasks.value.filter(t => t.status === 'completed').length
  const inProgress = formattedTasks.value.filter(t => t.status === 'in_progress').length

  return {
    total: stats.activities,
    completed,
    inProgress,
    notStarted: stats.activities - completed - inProgress,
    delayed: 0,
    critical: stats.criticalActivities,
    progressRate: stats.activities > 0 ? Math.round((completed / stats.activities) * 100) : 0
  }
})

// 网络图任务索引映射（用于对齐任务列表）
const networkTaskIndexMap = computed(() => {
  const map = {}
  const visibleTasks = filteredTasks.value || []
  visibleTasks.forEach((task, index) => {
    map[task.id] = index
  })
  return map
})

// 网络图高度（根据节点层级计算，确保所有内容可见）
const networkHeight = computed(() => {
  // 使用 formattedTasks 以包含所有展开的子任务
  const allTasks = formattedTasks.value || []
  // 计算任务总高度（考虑分组标题可能会增加额外高度，这里简化处理）
  const taskListHeight = allTasks.length * rowHeight.value
  // 添加额外的空间用于网络图层级布局
  // 网络图可能有多个层级（每150px一个层级），确保足够的高度
  const networkLevelHeight = Math.ceil(allTasks.length / 5) * 150
  // 设置最小高度，取任务列表高度和网络图层级高度的最大值
  return Math.max(600, taskListHeight, networkLevelHeight + 200)
})

// 甘特图内容高度（任务列表高度）
const ganttContentHeight = computed(() => {
  const allTasks = formattedTasks.value || []
  return Math.max(600, allTasks.length * rowHeight.value)
})

// ==================== Undo/Redo 状态 ====================
const canUndo = computed(() => undoRedoStore.canUndo)
const canRedo = computed(() => undoRedoStore.canRedo)
const undoCount = computed(() => undoRedoStore.stackSize)
const selectedTasks = computed(() => {
  return state.filteredTasks.filter(t => selectedTaskIds.value.has(t.id))
})
const commandHistory = computed(() => undoRedoStore.getCommandHistory())

// 拖拽状态（使用 useGanttDrag composable 返回的值）
// 这些会在下面的 useGanttDrag 调用中定义
let isDragging, draggedTask, tooltipVisible, tooltipPosition, tooltipText, startDrag, cancelDrag

// VIEW_CONFIG（从 store 获取）
const VIEW_CONFIG = state.VIEW_CONFIG

// 依赖关系创建状态（使用 store）
const isCreatingDependency = computed(() => state.isCreatingDependency)
const dependencySourceTask = computed(() => state.dependencySourceTask)

// 依赖关系创建（使用 store actions）
const startDependencyCreation = (task) => actions.startDependencyCreation(task)
const cancelDependencyCreation = () => actions.cancelDependencyCreation()
const handleDependencyTargetClick = (targetTask) => actions.completeDependencyCreation(targetTask)

// ==================== 拖拽事件处理 ====================
const handleDragChange = (preview) => {
  // 拖拽过程中的回调，可以用于实时预览
}

const handleDragEnd = async (newTask, originalTask) => {
  try {
    // 使用 store 的 endDrag 方法
    actions.endDrag(newTask, originalTask)
    statusBarRef.value?.showStatus('任务位置已更改，记得保存', 'info', 2000)
  } catch (error) {
    console.error('处理任务拖拽失败:', error)
    ElMessage.error('处理任务拖拽失败')
  }
}

// ==================== 未保存更改管理 ====================
// 使用 store 的方法
const markAsUnsaved = () => actions.markUnsaved()
const handleSaveAll = () => actions.saveAll()

// ==================== Undo/Redo 处理 ====================
const handleUndo = async () => {
  if (!canUndo.value) return

  try {
    await undoRedoStore.undo((command) => {
      ElMessage.success(`已撤销: ${command.getDescription()}`)
      emit('task-updated', null)
    })
  } catch (error) {
    console.error('撤销失败:', error)
    ElMessage.error('撤销失败')
  }
}

const handleRedo = async () => {
  if (!canRedo.value) return

  try {
    await undoRedoStore.redo((command) => {
      ElMessage.success(`已重做: ${command.getDescription()}`)
      emit('task-updated', null)
    })
  } catch (error) {
    console.error('重做失败:', error)
    ElMessage.error('重做失败')
  }
}

// 撤销到指定位置
const handleUndoTo = async (targetIndex) => {
  try {
    await undoRedoStore.undoTo(targetIndex)
    ElMessage.success(`已撤销到指定位置`)
    emit('task-updated', null)
  } catch (error) {
    console.error('撤销到指定位置失败:', error)
    ElMessage.error('撤销失败')
  }
}

// 清空历史
const handleClearHistory = () => {
  try {
    undoRedoStore.clear()
    ElMessage.success('已清空操作历史')
  } catch (error) {
    console.error('清空历史失败:', error)
    ElMessage.error('清空历史失败')
  }
}

// ==================== 多选功能 ====================
const handleDeleteSelected = async () => {
  if (selectedTaskIds.value.size === 0) return

  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedTaskIds.value.size} 个任务吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    statusBarRef.value?.showStatus(`正在删除 ${selectedTaskIds.value.size} 个任务...`, 'loading')

    // 批量删除
    for (const taskId of selectedTaskIds.value) {
      const task = state.filteredTasks.find(t => t.id === taskId)
      if (task) {
        const command = new DeleteTaskCommand(task, props.projectId)
        await undoRedoStore.execute(command)
      }
    }

    selectedTaskIds.value.clear()
    ElMessage.success('删除成功')
    statusBarRef.value?.showStatus('删除成功', 'success', 2000)
    emit('task-updated', null)
  } catch (error) {
    if (error !== 'cancel') {
      console.error('批量删除失败:', error)
      ElMessage.error('批量删除失败')
      statusBarRef.value?.showStatus('批量删除失败', 'error', 2000)
    }
  }
}

// ==================== 拖拽功能 ====================
const dragResult = useGanttDrag({
  dayWidth,
  timelineDays: computed(() => getters.timelineDays.value),
  onDragEnd: handleDragEnd,
  onDragChange: handleDragChange
})

// 解构拖拽状态
isDragging = dragResult.isDragging
const dragMode = dragResult.dragMode
draggedTask = dragResult.draggedTask
const previewTask = dragResult.previewTask
tooltipPosition = dragResult.tooltipPosition
tooltipVisible = dragResult.tooltipVisible
tooltipText = dragResult.tooltipText
startDrag = dragResult.startDrag
cancelDrag = dragResult.cancelDrag

// ==================== 使用 Store Getters ====================
// 直接使用 store 中的计算属性
const formattedTasks = computed(() => state.tasks)
const filteredTasks = computed(() => state.filteredTasks)
const groupedTasks = computed(() => getters.groupedTasks.value)
const taskStats = computed(() => getters.taskStats.value)
const timelineDays = computed(() => getters.timelineDays.value)
const timelineWeeks = computed(() => getters.timelineWeeks.value)
const timelineMonths = computed(() => getters.timelineMonths.value)
const timelineHeaderMonths = computed(() => getters.timelineHeaderMonths.value)
const timelineQuarters = computed(() => getters.timelineQuarters.value)
const timelineWidth = computed(() => getters.timelineWidth.value)
const todayPosition = computed(() => getters.todayPosition.value)
const currentPeriodText = computed(() => getters.currentPeriodText.value)
const currentZoomLabel = computed(() => getters.currentZoomLabel.value)
const emptyDescription = computed(() => getters.emptyDescription.value)
const containerStyle = computed(() => getters.containerStyle.value)
const arrowColor = computed(() => getters.arrowColor.value)
const arrowMarkerId = computed(() => getters.arrowMarkerId.value)
const visibleDependencies = computed(() => getters.visibleDependencies.value)

// 获取活动状态（保留用于向后兼容）
const getActivityStatus = (activity) => actions.getActivityStatus(activity)

// 自动适应内容尺寸
const autoFitContainer = () => actions.autoFit()

// ==================== 任务操作 ====================
// 行点击（支持多选）
const handleRowClick = (task) => {
  if (!task) return

  // 检查是否按下了 Ctrl/Cmd 键进行多选
  const event = window.event
  if (event && (event.ctrlKey || event.metaKey)) {
    if (selectedTaskIds.value.has(task.id)) {
      selectedTaskIds.value.delete(task.id)
    } else {
      selectedTaskIds.value.add(task.id)
    }
  } else {
    // 单选模式
    selectedTaskIds.value.clear()
    actions.selectTask(task.id)
  }
  emit('task-selected', task)
}

// 任务点击（支持多选）
const handleTaskClick = (task) => {
  if (!task) return

  // 检查是否按下了 Ctrl/Cmd 键进行多选
  const event = window.event
  if (event && (event.ctrlKey || event.metaKey)) {
    if (selectedTaskIds.value.has(task.id)) {
      selectedTaskIds.value.delete(task.id)
    } else {
      selectedTaskIds.value.add(task.id)
    }
  } else {
    // 单选模式
    selectedTaskIds.value.clear()
    actions.selectTask(task.id)
  }
  emit('task-selected', task)
}

// 任务双击
const handleTaskDblClick = (task) => {
  if (!task) return
  actions.selectTask(task.id)
  actions.openEditDialog(task)
}

// 行双击 - 打开编辑对话框
const handleRowDblClick = (task) => {
  if (!task) return
  actions.selectTask(task.id)
  actions.openEditDialog(task)
}

// 拖拽事件处理
const handleTaskMouseDown = (event, taskOrId, taskBarElement) => {
  // 只在左键点击时启动拖拽
  if (event.button !== 0) return

  // TaskTimeline may pass either a task object or just an id
  const taskId = typeof taskOrId === 'object' ? taskOrId.id : taskOrId
  const task = filteredTasks.value.find(t => t.id === taskId)
  if (!task) return

  // 如果有taskBarElement，使用它来检测拖拽模式（边缘调整）
  // 否则默认为移动模式
  if (taskBarElement) {
    startDrag(event, task, taskBarElement)
  } else {
    startDrag(event, task, event.target)
  }
}

// 依赖关系创建处理
const handleDependencyCreate = (taskOrId) => {
  const taskId = typeof taskOrId === 'object' ? taskOrId.id : taskOrId
  const task = filteredTasks.value.find(t => t.id === taskId)

  if (!task) return

  if (!isCreatingDependency.value) {
    // 开始创建依赖关系
    startDependencyCreation(task)
  } else {
    // 点击目标任务，完成依赖创建
    handleDependencyTargetClick(task)
  }
}

// 鼠标移动处理（更新临时连线位置）
const handleTimelineMouseMove = (position) => {
  if (state.isCreatingDependency) {
    actions.updateTempLineEnd(position)
  }
}

// 单元格编辑处理 - 仍然直接调用 API（暂时保持不变）
const handleCellEdit = async ({ taskId, updateData }) => {
  try {
    const task = state.filteredTasks.find(t => t.id === taskId)
    if (!task) return

    console.log('单元格编辑:', { taskId, updateData, task })

    statusBarRef.value?.showStatus('正在保存...', 'loading')

    const taskData = {
      project_id: props.projectId,
      ...updateData
    }

    await progressApi.update(taskId, taskData)
    ElMessage.success('保存成功')
    statusBarRef.value?.showStatus('保存成功', 'success', 1500)
    emit('task-updated', { ...task, ...taskData })
  } catch (error) {
    console.error('保存编辑失败:', error)
    const errorMsg = error.response?.data?.error || error.response?.data?.message || error.message || '保存失败'
    ElMessage.error(errorMsg)
    statusBarRef.value?.showStatus('保存失败', 'error', 2000)
  }
}

// 任务拖拽处理（更新父子关系和排序）- 使用 store
const handleTaskDragged = async ({ fromTask, toTask, position = 'child' }) => {
  console.log('GanttChart - handleTaskDragged 被调用:', { fromTask, toTask, position })

  try {
    actions.reorderTask(fromTask.id, toTask?.id || null, position)

    const messageMap = {
      before: `已将 "${fromTask.name}" 移动到 "${toTask?.name || '根'}" 之前`,
      after: `已将 "${fromTask.name}" 移动到 "${toTask?.name || '根'}" 之后`,
      child: toTask === null ? `已将 "${fromTask.name}" 移动到根级别` : `已将 "${fromTask.name}" 移动到 "${toTask.name}" 下`
    }

    statusBarRef.value?.showStatus(messageMap[position] || messageMap.child, 'info', 2000)
  } catch (error) {
    console.error('处理任务拖拽失败:', error)
    ElMessage.error('处理任务拖拽失败')
  }
}

// ==================== 容器调整大小 ====================
const handleResizeStart = (event) => {
  console.log('handleResizeStart 触发')
  event.preventDefault()
  event.stopPropagation()
  actions.setResizing(true)
  resizeDirection.value = 'both'

  const rect = ganttContainer.value.getBoundingClientRect()
  resizeStart.value = {
    x: event.clientX,
    y: event.clientY,
    width: rect.width,
    height: rect.height
  }

  document.addEventListener('mousemove', handleResizeMove)
  document.addEventListener('mouseup', handleResizeEnd)

  document.body.style.cursor = 'nwse-resize'
  document.body.style.userSelect = 'none'
}

const handleResizeStartRight = (event) => {
  console.log('handleResizeStartRight 触发')
  event.preventDefault()
  event.stopPropagation()
  actions.setResizing(true)
  resizeDirection.value = 'horizontal'

  const rect = ganttContainer.value.getBoundingClientRect()
  resizeStart.value = {
    x: event.clientX,
    y: 0,
    width: rect.width,
    height: 0
  }

  document.addEventListener('mousemove', handleResizeMove)
  document.addEventListener('mouseup', handleResizeEnd)

  document.body.style.cursor = 'ew-resize'
  document.body.style.userSelect = 'none'
}

const handleResizeStartBottom = (event) => {
  console.log('handleResizeStartBottom 触发')
  event.preventDefault()
  event.stopPropagation()
  actions.setResizing(true)
  resizeDirection.value = 'vertical'

  const rect = ganttContainer.value.getBoundingClientRect()
  resizeStart.value = {
    x: 0,
    y: event.clientY,
    width: 0,
    height: rect.height
  }

  document.addEventListener('mousemove', handleResizeMove)
  document.addEventListener('mouseup', handleResizeEnd)

  document.body.style.cursor = 'ns-resize'
  document.body.style.userSelect = 'none'
}

const handleResizeMove = (event) => {
  if (!state.isResizing) return

  const deltaX = event.clientX - resizeStart.value.x
  const deltaY = event.clientY - resizeStart.value.y

  if (resizeDirection.value === 'both' || resizeDirection.value === 'horizontal') {
    const newWidth = Math.max(600, resizeStart.value.width + deltaX)
    actions.setContainerSize(newWidth, state.containerSize.height)
  }

  if (resizeDirection.value === 'both' || resizeDirection.value === 'vertical') {
    const newHeight = Math.max(300, resizeStart.value.height + deltaY)
    actions.setContainerSize(state.containerSize.width, newHeight)
  }
}

const handleResizeEnd = () => {
  actions.setResizing(false)
  resizeDirection.value = null

  document.removeEventListener('mousemove', handleResizeMove)
  document.removeEventListener('mouseup', handleResizeEnd)

  document.body.style.cursor = ''
  document.body.style.userSelect = ''
}

// ==================== 平移功能 ====================
// 平移开始
const handlePanStart = (event) => {
  if (!state.panMode) return

  event.preventDefault()
  isPanning.value = true

  // 获取时间轴滚动区域（ganttBodyRef）
  const timelineScrollArea = ganttBodyRef?.value
  const currentScrollLeft = timelineScrollArea?.scrollLeft || 0

  panStart.value = {
    x: event.clientX,
    scrollLeft: currentScrollLeft
  }

  document.addEventListener('mousemove', handlePanMove)
  document.addEventListener('mouseup', handlePanEnd)
  document.body.style.cursor = 'grabbing'
  document.body.style.userSelect = 'none'
}

// 平移移动
const handlePanMove = (event) => {
  if (!isPanning.value) return

  const deltaX = event.clientX - panStart.value.x
  const newScrollLeft = panStart.value.scrollLeft - deltaX

  // 只滚动时间轴区域（ganttBodyRef），不影响任务列表
  const timelineScrollArea = ganttBodyRef?.value
  if (timelineScrollArea) {
    timelineScrollArea.scrollLeft = newScrollLeft
  }
  actions.setScrollLeft(newScrollLeft)
}

// 平移结束
const handlePanEnd = () => {
  isPanning.value = false

  document.removeEventListener('mousemove', handlePanMove)
  document.removeEventListener('mouseup', handlePanEnd)

  if (state.panMode) {
    document.body.style.cursor = 'grab'
  } else {
    document.body.style.cursor = ''
  }
  document.body.style.userSelect = ''
}

// ==================== 视图控制 ====================
// 日期导航（滚动时间轴视图）
const navigateDate = (direction) => {
  // 根据时间轴格式决定滚动的天数
  const format = state.timelineFormat
  let daysToScroll = 7 // 默认滚动一周

  if (format === 'day' || format === 'month-day' || format === 'year-month-day') {
    daysToScroll = 7
  } else if (format === 'week') {
    daysToScroll = 1 // 滚动一周
  } else if (format === 'month') {
    daysToScroll = 30 // 滚动一个月
  } else if (format === 'quarter' || format === 'year-month') {
    daysToScroll = 90 // 滚动一个季度
  }

  // TODO: 实现时间轴滚动功能
  // 这里需要配合时间轴的 scrollLeft 值来实现滚动
  // 暂时先使用简单的 dayWidth 调整作为临时方案
  const offset = direction * daysToScroll
  console.log('navigateDate:', { direction, daysToScroll, offset })
}

const goToToday = () => {
  ElMessage.success('已回到今天')
}

// 缩放控制
const zoomIn = () => actions.zoomIn()
const zoomOut = () => actions.zoomOut()
const zoomReset = () => actions.zoomReset()

// 视图模式切换
const handleViewModeChange = (newMode) => actions.setViewMode(newMode)

// 时间轴格式切换
const handleTimelineFormatChange = (newFormat) => actions.setTimelineFormat(newFormat)

// 日期显示格式切换
const handleDateFormatChange = (newFormat) => actions.setDateDisplayFormat(newFormat)

// 显示选项切换
const toggleDependencies = () => actions.toggleDependencies()
const toggleCriticalPath = () => actions.toggleCriticalPath()
const toggleBaseline = () => actions.toggleBaseline()
const toggleTaskList = () => actions.toggleTaskList()

// 平移模式切换
const togglePanMode = () => actions.togglePanMode()

// 选择模式切换（关闭平移模式，返回正常拖拽模式）
const toggleSelectMode = () => {
  if (state.panMode) {
    actions.togglePanMode() // 关闭平移模式
  }
}

// 打开资源管理对话框
const openResourceManagement = () => actions.openResourceManagementDialog()

// 分组控制
const handleGroupChange = (newGroup) => actions.setGroupMode(newGroup)

// 图表视图切换（甘特图/网络图）
const handleToggleViewMode = (mode) => {
  chartViewMode.value = mode
}

// ==================== 网络图视图控制 ====================
// 网络图缩放控制（操作 dayWidth，与时间标尺保持一致）
const networkZoomIn = () => {
  state.dayWidth = Math.min(100, state.dayWidth + 5)
  eventBus.emit(GanttEvents.ZOOM_CHANGED, { width: state.dayWidth })
}

const networkZoomOut = () => {
  state.dayWidth = Math.max(10, state.dayWidth - 5)
  eventBus.emit(GanttEvents.ZOOM_CHANGED, { width: state.dayWidth })
}

const networkZoomReset = () => {
  // 重置为当前视图模式的默认值
  const format = state.timelineFormat
  const config = state.VIEW_CONFIG[format] || state.VIEW_CONFIG['month-day']
  state.dayWidth = config.default
  eventBus.emit(GanttEvents.ZOOM_CHANGED, { width: state.dayWidth })
}

// ==================== 网络图事件处理 ====================
// 网络图节点点击
const handleNetworkNodeClick = (node) => {
  // 在 AOA 模式中，节点代表事件，点击节点显示信息
  console.log('Network node clicked:', node)
}

// 网络图任务箭头点击
const handleNetworkTaskClick = (task) => {
  actions.selectTask(task.id)
}

// 网络图任务箭头双击（编辑任务）
const handleNetworkTaskDblClick = (task) => {
  actions.selectTask(task.id)
  actions.openEditDialog(state.tasks.find(t => t.id === task.id))
}

// 网络图节点右键菜单
const handleNetworkNodeContextMenu = ({ event, node, selectedNodes, selectedIds, canMerge, x, y }) => {
  // 显示节点右键菜单
  nodeContextMenuNode.value = node
  nodeContextMenuSelectedNodes.value = selectedNodes || [node]
  nodeContextMenuCanMerge.value = canMerge || false
  nodeContextMenuPosition.value = { x, y }
  nodeContextMenuVisible.value = true
}

// 处理节点菜单命令
const handleNodeContextMenuCommand = async (command, data) => {
  switch (command) {
    case 'merge-nodes':
      await handleNetworkNodeMerge({
        nodes: data.selectedNodes,
        nodeIds: data.selectedNodes.map(n => n.id)
      })
      break
    case 'split-node':
      await handleNetworkNodeSplit({
        node: data.node,
        nodeId: data.nodeId
      })
      break
    case 'edit-tasks':
      // 找到节点关联的任务并打开编辑对话框
      if (data.node?.tasks) {
        const taskIds = [
          ...(data.node.tasks.start || []),
          ...(data.node.tasks.end || [])
        ]
        if (taskIds.length > 0) {
          const firstTask = state.tasks.find(t => t.id === taskIds[0])
          if (firstTask) {
            actions.selectTask(firstTask.id)
            actions.openEditDialog(firstTask)
          }
        }
      }
      break
  }
  nodeContextMenuVisible.value = false
}

// 网络图任务连线右键菜单
const handleNetworkTaskContextMenu = ({ event, task, x, y }) => {
  // 设置右键菜单任务
  contextMenuTask.value = state.tasks.find(t => t.id === task.id)
  contextMenuPosition.value = { x, y }
  contextMenuVisible.value = true
}

// 网络图节点拖动改变任务时间
const handleNetworkNodeTimeChange = ({ taskId, nodeType, newDate, daysDelta }) => {
  const task = state.tasks.find(t => t.id === taskId)
  if (!task) return

  // 计算新的开始/结束日期
  let updateData = {
    id: taskId,
    name: task.name,
    start_date: task.start,
    end_date: task.end,
    progress: task.progress || 0,
    priority: task.priority || 1
  }

  if (nodeType === 'start') {
    // 改变开始时间，保持结束时间不变（工期会改变）
    updateData.start_date = newDate
    // end_date 保持原值不变（已经在上面设置为 task.end）
  } else {
    // 改变结束时间
    updateData.end_date = newDate
  }

  // 更新本地状态
  Object.assign(task, {
    start: updateData.start_date,
    end: updateData.end_date,
    startDate: new Date(updateData.start_date),
    endDate: new Date(updateData.end_date)
  })

  // 记录待保存的更新
  state.pendingTaskUpdates.set(taskId, updateData)
  actions.markUnsaved()

  // 不显示提示，避免批量更新时重复显示
  // ElMessage.success(`任务时间已更新`)
}

// 网络图节点间创建依赖关系
const handleNetworkDependencyCreate = async ({ fromTaskIds, toTaskIds }) => {
  // 新节点系统：从节点对象中直接获取任务ID
  // fromTaskIds: 源节点结束的任务ID数组
  // toTaskIds: 目标节点开始的任野ID数组

  // 取第一个任务ID创建依赖（如果有多个任务，用户可以后续手动添加）
  const fromTaskId = fromTaskIds[0]
  const toTaskId = toTaskIds[0]

  // 验证：不能创建到自身的依赖
  if (fromTaskId === toTaskId) {
    ElMessage.warning('不能创建任务自身的依赖关系')
    return
  }

  // 获取源任务和目标任务
  const sourceTask = state.tasks.find(t => t.id === fromTaskId)
  const targetTask = state.tasks.find(t => t.id === toTaskId)

  if (!sourceTask || !targetTask) {
    ElMessage.error('任务不存在')
    return
  }

  // 检查是否是父子关系
  const isParentChild = actions.checkParentChildRelation(fromTaskId, toTaskId)
  if (isParentChild) {
    ElMessage.warning('不能在父子任务之间创建依赖关系')
    return
  }

  // 检查是否是从右到左（源任务的结束晚于目标任务的开始）
  const sourceEndDate = new Date(sourceTask.end)
  const targetStartDate = new Date(targetTask.start)
  if (sourceEndDate > targetStartDate) {
    ElMessage.warning('不能创建从右到左的依赖关系（前置任务必须先于后置任务）')
    return
  }

  // 调用API创建依赖关系
  try {
    await progressApi.createDependencyVisual(fromTaskId, toTaskId, { type: 'FS', lag: 0 })

    // 立即更新本地状态，添加前置依赖关系
    if (!targetTask.predecessors) {
      targetTask.predecessors = []
    }
    if (!targetTask.predecessors.includes(fromTaskId)) {
      targetTask.predecessors.push(fromTaskId)
    }

    // 同时更新源任务的后置依赖
    if (!sourceTask.successors) {
      sourceTask.successors = []
    }
    if (!sourceTask.successors.includes(toTaskId)) {
      sourceTask.successors.push(toTaskId)
    }

    ElMessage.success('依赖关系创建成功')

    // 异步重新加载数据以获取最新的完整状态
    actions.loadData().catch(err => {
      console.warn('后台重新加载数据失败:', err)
    })
  } catch (error) {
    console.error('创建依赖关系失败:', error)
    const errorMsg = error.response?.data?.error || error.response?.data?.message || error.message || '创建依赖关系失败'
    ElMessage.error(errorMsg)
  }
}

// 网络图节点合并（实现R11：相同紧前/紧后条件的任务共享节点）
const handleNetworkNodeMerge = async ({ nodes, nodeIds }) => {
  if (!nodes || nodes.length < 2) {
    ElMessage.warning('请选择至少两个节点进行合并')
    return
  }

  // R11规范：合并节点的条件
  // 1. 合并所有开始节点（相同紧前条件） - 所有任务应该具有相同的predecessors
  // 2. 合并所有结束节点（相同紧后条件） - 所有任务应该具有相同的successors

  // 分析节点类型（开始节点 vs 结束节点）
  const startNodes = nodes.filter(n => n.tasks && n.tasks.start.length > 0)
  const endNodes = nodes.filter(n => n.tasks && n.tasks.end.length > 0)

  // 获取所有涉及的任务ID
  const allTaskIds = new Set()
  nodes.forEach(node => {
    if (node.tasks) {
      node.tasks.start?.forEach(id => allTaskIds.add(id))
      node.tasks.end?.forEach(id => allTaskIds.add(id))
    }
  })

  // 检查是否可以合并（开始节点和结束节点不能混合合并）
  if (startNodes.length > 0 && endNodes.length > 0) {
    ElMessage.warning('开始节点和结束节点不能一起合并，请选择相同类型的节点')
    return
  }

  const isStartNodeMerge = startNodes.length > 0

  try {
    // 获取目标参考日期（使用最早节点的日期）
    const targetDays = Math.min(...nodes.map(n => n.days))

    // 对于开始节点合并：统一所有任务的开始日期
    // 对于结束节点合并：统一所有任务的结束日期
    const targetDate = timelineDays.value[targetDays]?.date
    if (!targetDate) {
      ElMessage.error('无法确定目标日期')
      return
    }

    // 批量更新任务日期
    const updatePromises = []
    for (const taskId of allTaskIds) {
      const task = state.tasks.find(t => t.id === taskId)
      if (!task) continue

      let updateData = {
        id: taskId,
        name: task.name,
        start_date: task.start,
        end_date: task.end,
        progress: task.progress || 0,
        priority: task.priority || 1
      }

      if (isStartNodeMerge) {
        // 合并开始节点：统一开始日期
        updateData.start_date = targetDate
      } else {
        // 合并结束节点：统一结束日期
        updateData.end_date = targetDate
      }

      // 验证日期有效性
      const newStart = new Date(updateData.start_date)
      const newEnd = new Date(updateData.end_date)
      if (newStart > newEnd) {
        if (isStartNodeMerge) {
          updateData.end_date = updateData.start_date
        } else {
          updateData.start_date = updateData.end_date
        }
      }

      // 更新本地状态
      Object.assign(task, {
        start: updateData.start_date,
        end: updateData.end_date,
        startDate: new Date(updateData.start_date),
        endDate: new Date(updateData.end_date)
      })

      state.pendingTaskUpdates.set(taskId, updateData)
      updatePromises.push(
        progressApi.updateTask(taskId, updateData).catch(err => {
          console.error(`更新任务 ${taskId} 失败:`, err)
        })
      )
    }

    await Promise.all(updatePromises)

    ElMessage.success(`成功合并 ${allTaskIds.size} 个任务的${isStartNodeMerge ? '开始' : '结束'}节点`)
    actions.markUnsaved()

    // 重新加载数据以刷新网络图
    await actions.loadData()

  } catch (error) {
    console.error('合并节点失败:', error)
    ElMessage.error('合并节点失败: ' + (error.message || '未知错误'))
  }
}

// 网络图节点拆分
const handleNetworkNodeSplit = async ({ node, nodeId }) => {
  if (!node) {
    ElMessage.warning('无效的节点')
    return
  }

  // 检查节点是否可以被拆分
  // 只有共享节点（关联多个任务）才能被拆分
  const relatedTaskCount = (node.tasks?.start?.length || 0) + (node.tasks?.end?.length || 0)

  if (relatedTaskCount <= 1) {
    ElMessage.info('该节点只关联一个任务，无需拆分')
    return
  }

  // 显示拆分选项对话框
  ElMessageBox.confirm(
    `该节点关联 ${relatedTaskCount} 个任务。拆分后，每个任务将拥有独立的节点。是否继续？`,
    '拆分节点',
    {
      confirmButtonText: '确定拆分',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      // 拆分逻辑：为每个任务创建稍微偏移的日期
      // 这样它们就会自动获得独立的节点
      const taskIds = [
        ...(node.tasks?.start || []),
        ...(node.tasks?.end || [])
      ]

      let offsetDays = 0
      const updatePromises = []

      for (const taskId of taskIds) {
        const task = state.tasks.find(t => t.id === taskId)
        if (!task) continue

        // 为后续任务添加1天偏移，确保它们获得独立的节点
        if (offsetDays > 0) {
          let updateData = {
            id: taskId,
            name: task.name,
            progress: task.progress || 0,
            priority: task.priority || 1
          }

          // 根据节点类型调整相应日期
          if (node.tasks?.start?.includes(taskId)) {
            const newStart = addDays(new Date(task.start), offsetDays)
            updateData.start_date = formatDate(newStart)
            updateData.end_date = task.end

            Object.assign(task, {
              start: updateData.start_date,
              startDate: newStart
            })
          } else if (node.tasks?.end?.includes(taskId)) {
            const newEnd = addDays(new Date(task.end), offsetDays)
            updateData.start_date = task.start
            updateData.end_date = formatDate(newEnd)

            Object.assign(task, {
              end: updateData.end_date,
              endDate: newEnd
            })
          }

          state.pendingTaskUpdates.set(taskId, updateData)
          updatePromises.push(
            progressApi.updateTask(taskId, updateData).catch(err => {
              console.error(`更新任务 ${taskId} 失败:`, err)
            })
          )
        }

        offsetDays++
      }

      if (updatePromises.length > 0) {
        await Promise.all(updatePromises)
        ElMessage.success(`成功拆分节点，${updatePromises.length} 个任务已调整为独立节点`)
        actions.markUnsaved()

        // 重新加载数据以刷新网络图
        await actions.loadData()
      } else {
        ElMessage.info('无需拆分')
      }

    } catch (error) {
      if (error !== 'cancel') {
        console.error('拆分节点失败:', error)
        ElMessage.error('拆分节点失败: ' + (error.message || '未知错误'))
      }
    }
  }).catch(() => {
    // 用户取消
  })
}

// 网络图任务依赖编辑
const handleNetworkTaskDependencyEdit = ({ task }) => {
  // 打开任务编辑对话框，用户可以在其中修改依赖关系
  actions.selectTask(task.id)
  actions.openEditDialog(state.tasks.find(t => t.id === task.id))
}

const toggleGroup = (groupName) => actions.toggleGroup(groupName)

// 刷新
const handleRefresh = () => emit('task-updated', null)

// 全屏切换（使用 CSS 模拟，避免弹出层问题）
const toggleFullscreen = () => {
  actions.setFullscreen(!state.isFullscreen)
}

// ==================== AOA图分析功能 ====================

// 计算关键路径
const handleCalculateCriticalPath = () => {
  const tasks = formattedTasks.value
  const criticalTasks = tasks.filter(t => t.is_critical)

  if (criticalTasks.length === 0) {
    ElMessage.warning('没有找到关键路径任务')
    return
  }

  // 计算关键路径的总工期
  const criticalPathDuration = criticalTasks.reduce((max, task) => {
    return Math.max(max, task.duration || 0)
  }, 0)

  ElMessageBox.alert(`
    <div style="text-align: left;">
      <p><strong>关键路径分析结果：</strong></p>
      <p>• 关键任务数量：${criticalTasks.length}</p>
      <p>• 关键路径总工期：${criticalPathDuration} 天</p>
      <p><strong>关键任务列表：</strong></p>
      ${criticalTasks.map(t => `<p>• ${t.name} (工期: ${t.duration}天)</p>`).join('')}
    </div>
  `, '关键路径', {
    dangerouslyUseHTMLString: true,
    confirmButtonText: '确定'
  })
}

// 节点属性分析
const handleAnalyzeNodeProperties = () => {
  const nodes = networkViewRef.value?.eventNodes || []

  if (nodes.length === 0) {
    ElMessage.warning('没有节点可分析')
    return
  }

  // 分析节点属性
  const mergedNodes = nodes.filter(n => n.isMerged)
  const startNodes = nodes.filter(n => n.tasks.start.length > 0)
  const endNodes = nodes.filter(n => n.tasks.end.length > 0)

  ElMessageBox.alert(`
    <div style="text-align: left;">
      <p><strong>节点属性分析：</strong></p>
      <p>• 总节点数：${nodes.length}</p>
      <p>• 合并节点数（共享节点）：${mergedNodes.length}</p>
      <p>• 起始节点数：${startNodes.length}</p>
      <p>• 结束节点数：${endNodes.length}</p>
      <p><strong>R11规范符合度：${((mergedNodes.length / nodes.length) * 100).toFixed(1)}%</strong></p>
    </div>
  `, '节点属性分析', {
    dangerouslyUseHTMLString: true,
    confirmButtonText: '确定'
  })
}

// 路径优化检查
const handleCheckPathOptimization = () => {
  const tasks = formattedTasks.value
  const tasksWithSlack = tasks.filter(t => t.slack > 0)

  if (tasksWithSlack.length === 0) {
    ElMessage.info('所有任务都在关键路径上，无可优化空间')
    return
  }

  // 找出优化潜力最大的任务
  const optimizableTasks = tasksWithSlack
    .sort((a, b) => b.slack - a.slack)
    .slice(0, 5)

  const totalSlack = tasksWithSlack.reduce((sum, t) => sum + t.slack, 0)

  ElMessageBox.alert(`
    <div style="text-align: left;">
      <p><strong>路径优化分析：</strong></p>
      <p>• 可优化任务数：${tasksWithSlack.length}</p>
      <p>• 总时差（优化潜力）：${totalSlack} 天</p>
      <p><strong>优化潜力最大的5个任务：</strong></p>
      ${optimizableTasks.map(t => `<p>• ${t.name} (时差: ${t.slack}天)</p>`).join('')}
      <p style="color: #909399; font-size: 12px;">建议：可以考虑延长这些任务的工期或调整资源配置</p>
    </div>
  `, '路径优化检查', {
    dangerouslyUseHTMLString: true,
    confirmButtonText: '确定'
  })
}

// 格式化日期为 YYYY-MM-DD（与 NetworkView.vue 中的 formatDate 完全一致）
const formatLocalDate = (date) => {
  if (!date) return ''
  const d = date instanceof Date ? date : new Date(date)
  if (isNaN(d.getTime())) return ''
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

// 规则验证（R4/R11）
const handleValidateRules = () => {
  const nodes = networkViewRef.value?.eventNodes || []
  const tasks = formattedTasks.value
  const violations = []

  // R4: 起点节点编号 < 终点节点编号
  tasks.forEach(task => {
    const startNode = nodes.find(n => n.tasks.start.includes(task.id))
    const endNode = nodes.find(n => n.tasks.end.includes(task.id))

    if (startNode && endNode && startNode.number >= endNode.number) {
      violations.push({
        rule: 'R4',
        task: task.name,
        message: `起点(${startNode.number}) >= 终点(${endNode.number})`
      })
    }
  })

  // R11: 共享节点检查
  // R11: 相同开始日期 + 相同前置条件的任务应该共享同一个开始节点
  const unmergedNodes = tasks.filter(task => {
    const predecessors = normalizePredecessors(task.predecessors)
    const predecessorSignature = getPredecessorSignature(predecessors)
    // 使用与 NetworkView.vue 完全相同的日期格式化方式
    const startDateKey = formatLocalDate(task.start)

    // 查找有相同前置条件 AND 相同开始日期的其他任务
    const sameConditionTasks = tasks.filter(t => {
      if (t.id === task.id) return false
      const tPreds = normalizePredecessors(t.predecessors)
      const tSignature = getPredecessorSignature(tPreds)
      // 使用与 NetworkView.vue 完全相同的日期格式化方式
      const tStartDateKey = formatLocalDate(t.start)

      // 使用签名比较，确保与 NetworkView.vue 中的逻辑一致
      return tSignature === predecessorSignature && tStartDateKey === startDateKey
    })

    // 如果没有其他任务有相同条件，不需要共享节点
    if (sameConditionTasks.length === 0) return false

    // 检查这些任务是否实际上共享了同一个开始节点
    const taskStartNode = nodes.find(n => n.tasks.start.includes(task.id))
    if (!taskStartNode) return false

    // 检查所有相同条件的任务是否都在同一个开始节点中
    const allInSameNode = sameConditionTasks.every(t =>
      taskStartNode.tasks.start.includes(t.id)
    )

    // 如果不在同一个节点，则是R11违规
    return !allInSameNode
  })

  unmergedNodes.forEach(task => {
    violations.push({
      rule: 'R11',
      task: task.name,
      message: '应与其他任务共享开始节点'
    })
  })

  if (violations.length === 0) {
    ElMessage.success('规则验证通过！未发现违规项')
    return
  }

  ElMessageBox.alert(`
    <div style="text-align: left;">
      <p><strong>规则验证结果：</strong></p>
      <p style="color: #F56C6C;">发现 ${violations.length} 项违规</p>
      ${violations.slice(0, 10).map(v => `
        <p style="color: #F56C6C;">
          <strong>${v.rule} 违规</strong> - ${v.task}<br/>
          <span style="font-size: 12px;">${v.message}</span>
        </p>
      `).join('')}
      ${violations.length > 10 ? `<p style="color: #909399;">...还有 ${violations.length - 10} 项</p>` : ''}
    </div>
  `, '规则验证', {
    dangerouslyUseHTMLString: true,
    confirmButtonText: '确定'
  })
}

// 导出分析报告
const handleExportAnalysisReport = () => {
  const tasks = formattedTasks.value
  const nodes = networkViewRef.value?.eventNodes || []
  const criticalTasks = tasks.filter(t => t.is_critical)

  const report = {
    generatedAt: new Date().toISOString(),
    summary: {
      totalTasks: tasks.length,
      criticalTasks: criticalTasks.length,
      totalNodes: nodes.length,
      totalDuration: tasks.reduce((max, t) => Math.max(max, t.duration || 0), 0)
    },
    criticalPath: {
      tasks: criticalTasks.map(t => ({
        id: t.id,
        name: t.name,
        duration: t.duration,
        earlyStart: t.early_start,
        earlyFinish: t.early_finish,
        lateStart: t.late_start,
        lateFinish: t.late_finish
      }))
    },
    allTasks: tasks.map(t => ({
      id: t.id,
      name: t.name,
      duration: t.duration,
      slack: t.slack,
      isCritical: t.is_critical,
      predecessors: t.predecessors
    }))
  }

  // 导出为JSON文件
  const blob = new Blob([JSON.stringify(report, null, 2)], { type: 'application/json' })
  const link = document.createElement('a')
  link.href = URL.createObjectURL(blob)
  link.download = `aoa-analysis-${Date.now()}.json`
  link.click()
  URL.revokeObjectURL(link.href)

  ElMessage.success('分析报告已导出')
}

// ==================== 导出功能 ====================
const handleExportPNG = async () => {
  try {
    const html2canvas = (await import('html2canvas')).default
    const canvas = await html2canvas(ganttContainer.value, {
      backgroundColor: '#ffffff',
      scale: 2
    })

    const link = document.createElement('a')
    link.download = `gantt-chart-${Date.now()}.png`
    link.href = canvas.toDataURL('image/png')
    link.click()

    ElMessage.success('导出PNG成功')
  } catch (error) {
    console.error('导出PNG失败:', error)
    ElMessage.error('导出PNG失败')
  }
}

const handleExportPDF = async () => {
  try {
    const html2canvas = (await import('html2canvas')).default
    const { default: jsPDF } = await import('jspdf')

    const canvas = await html2canvas(ganttContainer.value, {
      backgroundColor: '#ffffff',
      scale: 2
    })

    const imgData = canvas.toDataURL('image/png')
    const pdf = new jsPDF({
      orientation: canvas.width > canvas.height ? 'landscape' : 'portrait',
      unit: 'px',
      format: [canvas.width, canvas.height]
    })

    pdf.addImage(imgData, 'PNG', 0, 0, canvas.width, canvas.height)
    pdf.save(`gantt-chart-${Date.now()}.pdf`)

    ElMessage.success('导出PDF成功')
  } catch (error) {
    console.error('导出PDF失败:', error)
    ElMessage.error('导出PDF失败')
  }
}

// ==================== 右键菜单 ====================
const handleContextMenu = (eventOrData) => {
  // 支持多种参数格式
  let event, task, type, action

  if (eventOrData && typeof eventOrData === 'object' && 'event' in eventOrData) {
    // 来自 TaskTable 的数据对象格式
    event = eventOrData.event
    task = eventOrData.task
    type = eventOrData.type
    action = eventOrData.action
  } else if (eventOrData && eventOrData.target) {
    // 来自 TaskTimeline 的原始事件格式
    event = eventOrData
    task = null
    type = 'task'

    // 查找被点击的任务行
    const taskRow = event.target.closest('.table-row')
    if (taskRow) {
      // 通过索引或属性找到对应的任务
      const taskId = taskRow.getAttribute('data-task-id')
      if (taskId) {
        task = state.filteredTasks.find(t => t.id == taskId)
      }
      if (!task) {
        // 获取任务索引
        const ganttBodyRef = ganttBodyRef
        if (ganttBodyRef) {
          const taskRows = Array.from(ganttBodyRef.value?.querySelectorAll('.table-row') || [])
          const taskIndex = taskRows.indexOf(taskRow)
          const tasks = state.filteredTasks
          if (taskIndex >= 0 && tasks[taskIndex]) {
            task = tasks[taskIndex]
          }
        }
      }
    }
  } else {
    // 没有事件对象，可能是直接调用
    event = null
    task = eventOrData?.task || null
    type = eventOrData?.type || 'task'
    action = eventOrData?.action || null
  }

  // 处理空白行的新建任务
  if (type === 'new-task') {
    if (action === 'create-immediate') {
      // 双击直接创建
      actions.openEditDialog(null)
    } else if (action === 'context-menu' && event) {
      // 右键显示菜单
      event.preventDefault()
      event.stopPropagation()
      actions.showContextMenu(null, { x: event.clientX, y: event.clientY })
      nextTick(() => {
        contextMenuRef.value?.open()
      })
    } else {
      // 单击也直接创建
      actions.openEditDialog(null)
    }
    return
  }

  // 处理已有任务的右键菜单
  if (!task) return

  if (event) {
    event.preventDefault()
    event.stopPropagation()
  }

  const position = event ? { x: event.clientX, y: event.clientY } : { x: 0, y: 0 }
  actions.showContextMenu(task, position)

  nextTick(() => {
    contextMenuRef.value?.open()
  })
}

// 时间轴缩放处理
const handleTimelineZoom = (newDayWidth) => {
  // 直接更新 state.dayWidth
  state.dayWidth = newDayWidth
  eventBus.emit(GanttEvents.ZOOM_CHANGED, { width: newDayWidth })
}

// 网络图缩放处理（与时间轴缩放保持一致）
const handleNetworkZoom = (newDayWidth) => {
  // 直接更新 state.dayWidth，这样时间标尺也会跟着缩放
  state.dayWidth = newDayWidth
  eventBus.emit(GanttEvents.ZOOM_CHANGED, { width: newDayWidth })
}

// 网络图平移处理（同步时间轴）
const handleNetworkPan = (pan) => {
  networkPanX.value = pan.x
}

// 甘特图滚动处理（同步时间轴）
const handleTimelineScroll = (scroll) => {
  timelinePanX.value = scroll.scrollLeft
}

// 右键菜单命令处理
const handleContextMenuEdit = (task) => handleTaskDblClick(task)

const handleContextMenuDuplicate = async (task) => {
  try {
    const newTask = {
      project_id: props.projectId,
      name: `${task.name} (副本)`,
      start_date: task.start,
      end_date: task.end,
      progress: 0,
      priority: task.priority || 'medium',
      description: ''
    }

    await progressApi.create(newTask)
    ElMessage.success('任务已复制')
    emit('task-updated', newTask)
  } catch (error) {
    console.error('复制任务失败:', error)
    ElMessage.error('复制任务失败')
  }
}

const handleContextMenuDelete = async (task) => {
  actions.selectTask(task.id)
  await handleDeleteTask()
}

// ==================== 依赖关系操作 ====================
// 依赖关系右键菜单
const handleDependencyContextMenu = ({ event, dependency, x, y }) => {
  event.preventDefault()
  event.stopPropagation()

  console.log('依赖关系右键菜单:', dependency)

  // 设置依赖关系上下文菜单状态
  dependencyContextMenuDependency.value = dependency
  dependencyContextMenuPosition.value = { x, y }
  dependencyContextMenuVisible.value = true
}

// 删除依赖关系
const handleDeleteDependency = async (dependency) => {
  if (!dependency) return

  console.log('删除依赖关系:', dependency)

  try {
    // 首先需要获取任务的依赖关系列表，找到对应的依赖ID
    const response = await progressApi.getDependencies(dependency.toId)
    const dependencies = response.data || []

    // 查找匹配的依赖关系（depends_on == fromId）
    const targetDep = dependencies.find(dep => dep.depends_on === dependency.fromId)

    if (!targetDep) {
      ElMessage.error('未找到对应的依赖关系')
      return
    }

    // 调用删除API
    await progressApi.removeDependency(targetDep.id)
    ElMessage.success('依赖关系已删除')

    // 刷新数据
    emit('task-updated', null)
  } catch (error) {
    console.error('删除依赖关系失败:', error)
    const errorMsg = error.response?.data?.error || error.response?.data?.message || error.message || '删除依赖关系失败'
    ElMessage.error(errorMsg)
  }
}

// 依赖关系菜单命令处理
const handleDependencyMenuCommand = async (command) => {
  if (command === 'delete-dependency') {
    await handleDeleteDependency(dependencyContextMenuDependency.value)
  }
}

// 依赖关系菜单可见性变化
const handleDependencyMenuVisibleChange = (visible) => {
  if (!visible) {
    dependencyContextMenuVisible.value = false
  }
}

// ==================== 任务层级操作 ====================
// 上移任务
const handleMoveTaskUp = async (task) => {
  try {
    await actions.moveTaskUp(task.id)
    ElMessage.success('任务已上移')
    emit('task-updated', null)
  } catch (error) {
    console.error('上移任务失败:', error)
    ElMessage.error('上移任务失败')
  }
}

// 下移任务
const handleMoveTaskDown = async (task) => {
  try {
    await actions.moveTaskDown(task.id)
    ElMessage.success('任务已下移')
    emit('task-updated', null)
  } catch (error) {
    console.error('下移任务失败:', error)
    ElMessage.error('下移任务失败')
  }
}

// 转为独立任务（解除父子关系）
const handleConvertToIndependent = async (task) => {
  try {
    await actions.convertToIndependentTask(task.id)
    ElMessage.success('任务已转为独立任务')
    emit('task-updated', null)
  } catch (error) {
    console.error('转为独立任务失败:', error)
    ElMessage.error('转为独立任务失败')
  }
}

// ==================== 新建任务 ====================
// 从右键菜单新建任务
const handleCreateTask = () => {
  actions.openEditDialog(null)
}

// ==================== 子任务管理 ====================
// 添加子任务
const handleAddSubtask = async (parentTask) => {
  try {
    console.log('添加子任务，父任务:', parentTask)

    const newTask = {
      project_id: props.projectId,
      name: `${parentTask.name} - 子任务`,
      start_date: parentTask.start,
      end_date: parentTask.end,
      progress: 0,
      priority: parentTask.priority || 'medium',
      description: ''
    }

    // 只有当父任务有ID时才设置parent_id
    if (parentTask.id) {
      newTask.parent_id = parentTask.id
    }

    console.log('创建子任务数据:', newTask)

    const response = await progressApi.create(newTask)
    console.log('子任务创建响应:', response)

    ElMessage.success('子任务已添加')

    // 更新父任务进度
    if (parentTask.id) {
      try {
        await progressApi.updateParentProgress(parentTask.id)
      } catch (updateError) {
        console.error('更新父任务进度失败:', updateError)
        // 不阻塞主流程，只记录错误
      }
    }

    // 刷新任务列表以显示新创建的子任务
    emit('task-updated', response.data)
  } catch (error) {
    console.error('添加子任务失败:', error)
    const errorMsg = error.response?.data?.error || error.response?.data?.message || error.message || '添加子任务失败'
    ElMessage.error(errorMsg)
  }
}

// 转换为里程碑
const handleConvertToMilestone = async (task) => {
  try {
    await ElMessageBox.confirm(
      '转换为里程碑会将开始日期设置为结束日期，是否继续？',
      '转换为里程碑',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    // 将开始日期设置为结束日期
    await progressApi.updateTask(task.id, {
      start_date: task.end,
      end_date: task.end,
      is_milestone: true,
      duration: 0
    })

    ElMessage.success('已转换为里程碑')
    // 里程碑转换后不需要刷新整个数据
    // emit('task-updated')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('转换里程碑失败:', error)
      ElMessage.error('转换失败')
    }
  }
}

// ==================== 依赖关系管理 ====================
// 添加前置任务（进入可视化创建模式）
const handleAddDependency = (task) => {
  startDependencyCreation(task)
}

// 查看依赖关系
const handleViewDependencies = async (task) => {
  try {
    const response = await progressApi.getDependencies(task.id)
    const deps = response.data || []

    if (deps.length === 0) {
      ElMessage.info('该任务暂无依赖关系')
      return
    }

    // 构建依赖关系描述
    const depDescriptions = deps.map(dep => {
      const depTask = state.tasks.find(t => t.id === dep.depends_on)
      if (!depTask) return ''

      const typeNames = {
        FS: '完成-开始',
        FF: '完成-完成',
        SS: '开始-开始',
        SF: '开始-完成'
      }

      return `${depTask.name} → ${task.name} (${typeNames[dep.type] || dep.type}${dep.lag ? ` +${dep.lag}天` : ''})`
    }).filter(Boolean)

    ElMessageBox.alert(
      depDescriptions.join('<br>'),
      `${task.name} 的依赖关系`,
      {
        dangerouslyUseHTMLString: true,
        confirmButtonText: '关闭'
      }
    )
  } catch (error) {
    console.error('查看依赖失败:', error)
    ElMessage.error('查看依赖失败')
  }
}

// ==================== 资源管理 ====================
// 加载资源库
const loadResources = async () => {
  await actions.loadResources()
}

// 分配资源
const handleAllocateResources = (task) => {
  actions.openResourceDialog(task)
}

// 资源分配保存后回调
const handleResourceSaved = () => {
  console.log('资源分配已保存，刷新数据')
  // 触发数据刷新以获取最新的资源信息
  emit('task-updated', null)
}

// 资源库刷新处理
const handleResourceRefresh = () => {
  // 重新加载资源库
  loadResources()
  // 只更新资源库，不刷新任务列表
  // emit('task-updated')
}

// ==================== 添加/编辑/删除任务 ====================
// 保存任务（使用命令模式支持撤销/重做）
const handleSaveTask = async (formData) => {
  try {
    console.log('GanttChart - handleSaveTask 收到的数据:', formData)
    console.log('GanttChart - formData.name:', formData.name)

    const isEdit = !!state.editingTask
    const statusMsg = isEdit ? '正在更新任务...' : '正在创建任务...'
    statusBarRef.value?.showStatus(statusMsg, 'loading')

    const taskData = {
      project_id: props.projectId,
      name: formData.name,
      start_date: formData.start,
      end_date: formData.end,
      progress: formData.progress,
      priority: formData.priority,
      description: formData.notes,
      parent_id: formData.parent_id || null,
      // 资源分配
      resources: formData.resources?.map(r => ({
        resource_id: r.resource_id || r.id,
        quantity: r.quantity,
        cost: r.cost,
        type: r.type
      })) || [],
      // 紧前任务
      predecessor_ids: formData.predecessor_ids || [],
      // 紧后任务
      successor_ids: formData.successor_ids || []
    }

    console.log('GanttChart - 准备发送到API的数据:', taskData)

    let result
    if (isEdit) {
      // 使用更新命令
      const command = new UpdateTaskCommand(
        state.editingTask.id,
        state.editingTask,
        taskData,
        (updatedTask) => {
          ElMessage.success('任务更新成功')
          statusBarRef.value?.showStatus('任务更新成功', 'success', 2000)
          emit('task-updated', updatedTask)
        },
        (error) => {
          console.error('保存任务失败:', error)
          ElMessage.error('保存任务失败')
          statusBarRef.value?.showStatus('保存任务失败', 'error', 2000)
        }
      )
      result = await undoRedoStore.execute(command)
    } else {
      // 使用创建命令
      const command = new CreateTaskCommand(
        props.projectId,
        taskData,
        async (createdTask) => {
          console.log('任务创建响应:', createdTask)
          ElMessage.success('任务创建成功')
          statusBarRef.value?.showStatus('任务创建成功', 'success', 2000)

          // 等待一小段时间，确保后端完成 CPM 计算
          await new Promise(resolve => setTimeout(resolve, 500))

          // 重新加载任务数据以获取最新的调度信息
          await actions.loadData()

          // 通知父组件刷新
          emit('task-updated', createdTask)
        },
        (error) => {
          console.error('保存任务失败:', error)
          ElMessage.error('保存任务失败')
          statusBarRef.value?.showStatus('保存任务失败', 'error', 2000)
        }
      )
      result = await undoRedoStore.execute(command)
    }

    if (result.success) {
      actions.closeEditDialog()
    }
  } catch (error) {
    console.error('保存任务失败:', error)
    ElMessage.error('保存任务失败')
    statusBarRef.value?.showStatus('保存任务失败', 'error', 2000)
  }
}

// 从详情抽屉编辑任务
const handleEditTaskFromDrawer = () => {
  const task = getters.selectedTask.value
  if (!task) return
  actions.closeTaskDetail()
  actions.openEditDialog(task)
}

// 复制任务
const handleDuplicateTask = async () => {
  const task = getters.selectedTask.value
  if (!task) return

  try {
    const newTask = {
      project_id: props.projectId,
      name: `${task.name} (副本)`,
      start_date: task.start,
      end_date: task.end,
      progress: 0,
      priority: task.priority || 'medium',
      description: ''
    }

    await progressApi.create(newTask)
    ElMessage.success('任务已复制')
    emit('task-updated', newTask)
  } catch (error) {
    console.error('复制任务失败:', error)
    ElMessage.error('复制任务失败')
  }
}

// 删除任务（包括所有子任务）
const handleDeleteTask = async () => {
  const task = getters.selectedTask.value
  if (!task) return

  try {
    // 查找所有子任务
    const findAllChildren = (parentTask) => {
      const children = state.tasks.filter(t => t.parent_id === parentTask.id)
      let allChildren = [...children]
      children.forEach(child => {
        allChildren = allChildren.concat(findAllChildren(child))
      })
      return allChildren
    }

    const childTasks = findAllChildren(task)
    const hasChildren = childTasks.length > 0

    // 构建确认消息
    let confirmMessage = `确定要删除任务"${task.name}"吗？`
    if (hasChildren) {
      confirmMessage = `确定要删除任务"${task.name}"吗？\n\n该任务包含 ${childTasks.length} 个子任务，删除父任务将同时删除所有子任务。`
    }

    await ElMessageBox.confirm(
      confirmMessage,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
        dangerouslyUseHTMLString: false
      }
    )

    statusBarRef.value?.showStatus('正在删除任务...', 'loading')

    // 删除所有子任务
    if (hasChildren) {
      for (const child of childTasks) {
        try {
          await progressApi.delete(child.id)
        } catch (error) {
          console.error(`删除子任务 ${child.name} (ID: ${child.id}) 失败:`, error)
        }
      }
    }

    // 删除父任务
    await progressApi.delete(task.id)

    ElMessage.success(hasChildren
      ? `已删除任务"${task.name}"及其 ${childTasks.length} 个子任务`
      : '任务已删除')
    statusBarRef.value?.showStatus('任务已删除', 'success', 2000)
    actions.closeTaskDetail()

    // 重新加载数据以刷新图形显示
    await actions.loadData()
    // 发送更新事件（不需要传递具体任务数据）
    emit('task-updated')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除任务失败:', error)
      ElMessage.error('删除任务失败')
      statusBarRef.value?.showStatus('删除任务失败', 'error', 2000)
    }
  }
}

// ==================== 模板和批量编辑 ====================
const handleTaskCreated = (result) => {
  ElMessage.success('任务已创建')
  emit('task-updated', result.data)
  templateDialogVisible.value = false
}

const handleToggleBulkEdit = () => {
  if (selectedTaskIds.value.size === 0) {
    ElMessage.warning('请先选择要编辑的任务')
    return
  }
  bulkEditDialogVisible.value = true
}

const handleBulkUpdate = (result) => {
  ElMessage.success(`已更新 ${result.count} 个任务`)
  selectedTaskIds.value.clear()
  bulkEditDialogVisible.value = false
  emit('task-updated', null)
}

// ==================== 键盘快捷键 ====================
const handleKeydown = (event) => {
  // 如果在编辑对话框中，不处理快捷键
  if (state.editDialogVisible) return

  // Ctrl/Cmd + Z: 撤销
  if ((event.ctrlKey || event.metaKey) && event.key === 'z' && !event.shiftKey) {
    event.preventDefault()
    handleUndo()
  }

  // Ctrl/Cmd + Y 或 Ctrl/Cmd + Shift + Z: 重做
  if (
    ((event.ctrlKey || event.metaKey) && event.key === 'y') ||
    ((event.ctrlKey || event.metaKey) && event.shiftKey && event.key === 'z')
  ) {
    event.preventDefault()
    handleRedo()
  }

  // Ctrl/Cmd + A: 全选任务
  if ((event.ctrlKey || event.metaKey) && event.key === 'a') {
    event.preventDefault()
    selectedTaskIds.value = new Set(state.filteredTasks.map(t => t.id))
    ElMessage.info(`已选中 ${selectedTaskIds.value.size} 个任务`)
  }

  // Ctrl/Cmd + N: 新建任务
  if ((event.ctrlKey || event.metaKey) && event.key === 'n') {
    event.preventDefault()
    handleAddTask()
  }

  // Ctrl/Cmd + D: 复制任务
  if ((event.ctrlKey || event.metaKey) && event.key === 'd' && getters.selectedTask.value) {
    event.preventDefault()
    handleDuplicateTask()
  }

  // Delete: 删除任务（单选或多选）
  if (event.key === 'Delete') {
    event.preventDefault()
    if (selectedTaskIds.value.size > 0) {
      // TODO: 批量删除
      handleDeleteSelected()
    } else if (getters.selectedTask.value) {
      handleDeleteTask()
    }
  }

  // ESC: 取消依赖连线绘制或关闭对话框
  if (event.key === 'Escape') {
    event.preventDefault()
    if (isCreatingDependency.value) {
      cancelDependencyCreation()
      ElMessage.info('已取消依赖关系绘制')
    }
  }

  // Enter: 编辑任务
  if (event.key === 'Enter' && getters.selectedTask.value) {
    event.preventDefault()
    const task = getters.selectedTask.value
    actions.selectTask(task.id)
    actions.openEditDialog(task)
  }

  // Escape: 取消选择
  if (event.key === 'Escape') {
    event.preventDefault()
    actions.selectTask(null)
    actions.closeTaskDetail()
    selectedTaskIds.value.clear()
    if (isDragging.value) {
      cancelDrag()
    }
  }

  // Ctrl/Cmd + F: 聚焦搜索框
  if ((event.ctrlKey || event.metaKey) && event.key === 'f') {
    event.preventDefault()
    document.querySelector('.gantt-toolbar input')?.focus()
  }

  // ← →: 导航日期
  if (event.key === 'ArrowLeft' && !event.ctrlKey && !event.metaKey) {
    navigateDate(-1)
  }
  if (event.key === 'ArrowRight' && !event.ctrlKey && !event.metaKey) {
    navigateDate(1)
  }

  // + -: 缩放
  if (event.key === '+' || event.key === '=') {
    actions.zoomIn()
  }
  if (event.key === '-' || event.key === '_') {
    actions.zoomOut()
  }
}

// ==================== 生命周期 ====================
// 页面离开前的警告
const handleBeforeUnload = (e) => {
  if (state.hasUnsavedChanges) {
    e.preventDefault()
    e.returnValue = '您有未保存的更改，确定要离开吗？'
    return '您有未保存的更改，确定要离开吗？'
  }
}

// 更新容器尺寸
const updateContainerSize = () => {
  if (ganttContainer.value) {
    const rect = ganttContainer.value.getBoundingClientRect()
    // 只设置宽度，让高度通过 flex 自适应
    actions.setContainerSize(rect.width, null)
  }
}

// 窗口大小改变时的处理
const handleWindowResize = () => {
  // 使用防抖来避免频繁更新
  if (resizeTimer) {
    clearTimeout(resizeTimer)
  }
  resizeTimer = setTimeout(() => {
    updateContainerSize()
  }, 200)
}

// ==================== 监听 props 变化 ====================
// 监听 scheduleData 变化（当任务更新时，父组件会重新加载 scheduleData）
watch(() => props.scheduleData, (newData) => {
  console.log('GanttChart - scheduleData changed, reformatting tasks')
  if (newData) {
    state.scheduleData = newData
    actions.formatTasks()
  }
}, { deep: true })

// 监听 dataVersion 变化（当任务更新时，强制重新格式化）
watch(() => props.dataVersion, (newVersion, oldVersion) => {
  if (newVersion !== oldVersion && newVersion > 0) {
    console.log('GanttChart - dataVersion changed to', newVersion, ', reformatting tasks')
    actions.formatTasks()
  }
})

onMounted(() => {
  console.log('GanttChart - mounted with refactored components and store')
  // 同步 scheduleData 到 store 并格式化任务
  state.scheduleData = props.scheduleData
  actions.formatTasks()

  // 加载资源库
  actions.loadResources()

  // 加载列宽配置
  actions.loadColumnWidths()
  actions.loadTaskListWidth()
  actions.loadTaskListHeight()

  // 初始化容器尺寸
  updateContainerSize()

  // 添加容器尺寸监听
  resizeObserver = new ResizeObserver(() => {
    updateContainerSize()
  })
  if (ganttContainer.value) {
    resizeObserver.observe(ganttContainer.value)
  }

  // 添加窗口大小监听
  window.addEventListener('resize', handleWindowResize)

  // 添加页面离开前的警告
  window.addEventListener('beforeunload', handleBeforeUnload)

  // 监听依赖关系错误事件 - 保存取消订阅函数
  unsubscribeDependencyError = eventBus.on(GanttEvents.DEPENDENCY_ERROR, ({ message }) => {
    ElMessage.warning(message)
  })
})

onUnmounted(() => {
  // 清理拖拽事件监听
  if (isDragging.value) {
    cancelDrag()
  }
  // 清理调整大小事件监听
  if (state.isResizing) {
    handleResizeEnd()
  }
  // 清理 ResizeObserver
  if (resizeObserver) {
    resizeObserver.disconnect()
    resizeObserver = null
  }
  // 移除窗口大小监听
  window.removeEventListener('resize', handleWindowResize)
  if (resizeTimer) {
    clearTimeout(resizeTimer)
    resizeTimer = null
  }
  // 移除页面离开前的警告
  window.removeEventListener('beforeunload', handleBeforeUnload)
  // 移除依赖关系错误监听 - 使用取消订阅函数
  if (unsubscribeDependencyError) {
    unsubscribeDependencyError()
    unsubscribeDependencyError = null
  }
  // 清理自动保存计时器
  actions.stopAutoSave()
})

// 监听scheduleData变化 - 同步到 store
watch(
  () => props.scheduleData,
  (newVal) => {
    console.log('GanttChart - scheduleData changed, syncing to store')
    state.scheduleData = newVal
    actions.formatTasks()

    // 初次加载时自动适配视图
    const newActivitiesCount = Object.keys(newVal.activities || {}).length
    if (newActivitiesCount > 0 && state.dataVersion === 0) {
      state.dataVersion++
      nextTick(() => {
        actions.autoFit()
      })
    }
  },
  { deep: true }
)

// 监听搜索关键词变化
watch(
  () => state.searchKeyword,
  (newKeyword) => {
    actions.filterTasks(newKeyword)
  }
)

// 监听 store 中的资源对话框状态变化
watch(
  () => state.resourceDialogVisible,
  (val) => {
    resourceDialogVisible.value = val
  }
)

// 双向同步：本地 ref → store
watch(resourceDialogVisible, (val) => {
  state.resourceDialogVisible = val
})

watch(
  () => state.resourceManagementDialogVisible,
  (val) => {
    resourceManagementDialogVisible.value = val
  }
)

// 双向同步：本地 ref → store
watch(resourceManagementDialogVisible, (val) => {
  state.resourceManagementDialogVisible = val
})

watch(
  () => state.currentTaskForResource,
  (val) => {
    currentTaskForResource.value = val
  }
)

// 监听视图模式切换，重置平移状态
watch(chartViewMode, (newMode) => {
  if (newMode === 'gantt') {
    // 切换到甘特图模式时，重置网络图平移偏移
    networkPanX.value = 0
  } else if (newMode === 'network') {
    // 切换到网络图模式时，重置甘特图滚动偏移
    timelinePanX.value = 0
  }
})
</script>

<style scoped>
.gantt-chart {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: #fff;
  border: none;
  border-radius: 0;
  overflow: visible;
  outline: none;
  min-height: 0;
}

/* 确保工具栏和统计信息不被压缩 */
.gantt-chart :deep(.gantt-toolbar),
.gantt-chart :deep(.gantt-stats) {
  flex-shrink: 0;
}

/* 甘特图容器 */
.gantt-container {
  flex: 1;
  overflow: visible;
  display: flex;
  flex-direction: column;
  position: relative;
  min-width: 600px;
  min-height: 0;
}

/* 表头容器 - 固定在顶部 */
.gantt-header-container {
  flex-shrink: 0;
  position: sticky;
  top: 0;
  z-index: 100;
  background: #fff;
  border-bottom: 1px solid #dcdfe6;
}

.gantt-container.is-resizing .gantt-header-container {
  pointer-events: none;
}

/* 平移模式 */
.gantt-container.is-pan-mode {
  cursor: grab;
}

.gantt-container.is-pan-mode:active {
  cursor: grabbing;
}

/* 调整大小手柄 */
.resize-handle {
  position: absolute;
  z-index: 9999;
  transition: background-color 0.2s;
  pointer-events: auto;
}

.resize-handle-corner {
  right: 0;
  bottom: 0;
  width: 24px;
  height: 24px;
  cursor: nwse-resize;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(144, 147, 153, 0.1);
  border-radius: 4px 0 0 0;
}

.resize-handle-corner:hover {
  background: rgba(64, 158, 255, 0.3);
}

.resize-handle-corner svg {
  pointer-events: none;
}

.resize-handle-right {
  right: 0;
  top: 0;
  bottom: 24px;
  width: 8px;
  cursor: ew-resize;
}

.resize-handle-right:hover {
  background: rgba(64, 158, 255, 0.3);
}

.resize-handle-bottom {
  right: 24px;
  bottom: 0;
  left: 0;
  height: 8px;
  cursor: ns-resize;
}

.resize-handle-bottom:hover {
  background: rgba(64, 158, 255, 0.3);
}

.gantt-container.is-resizing {
  user-select: none;
}

/* 确保调整手柄始终可见 */
.gantt-container:hover .resize-handle-corner {
  opacity: 1;
}

/* 全屏样式 */
.gantt-chart.is-fullscreen {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  width: 100vw;
  height: 100vh;
  z-index: 1000; /* 降低 z-index，避免覆盖 Element Plus 弹出层 */
  border-radius: 0;
  overflow: hidden;
}

.gantt-chart.is-fullscreen .gantt-container {
  /* 全屏时也不限制最大高度，让内容撑开 */
}
/* 主体区域：任务表格 + 视图区域 */
.gantt-body {
  display: flex;
  flex-direction: row;
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  min-height: 0;
  max-height: calc(100vh - 460px);
  /* 防止子元素被拉伸 */
  align-items: flex-start;
}

/* 任务表格区域 - 固定宽度 */
.gantt-body > *:first-child {
  flex-shrink: 0;
}

/* 视图区域 - 填充剩余空间 */
.gantt-body > .task-timeline,
.gantt-body > .network-view {
  flex: 1 1 auto;
  min-width: 0;
  overflow: hidden;
}

/* 自定义滚动条样式 */
.gantt-body::-webkit-scrollbar {
  width: 10px;
}

.gantt-body::-webkit-scrollbar-track {
  background: #f0f0f0;
  border-radius: 5px;
}

.gantt-body::-webkit-scrollbar-thumb {
  background: #c0c0c0;
  border-radius: 5px;
  transition: background 0.2s;
}

.gantt-body::-webkit-scrollbar-thumb:hover {
  background: #a0a0a0;
}

/* Firefox 滚动条样式 */
.gantt-body {
  scrollbar-width: thin;
  scrollbar-color: #c0c0c0 #f0f0f0;
}
</style>
