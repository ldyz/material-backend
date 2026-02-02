<template>
  <div
    class="gantt-chart"
    :class="chartClasses"
    @keydown="handleKeydown"
    tabindex="0"
    ref="ganttChartRef"
    v-bind="getAriaAttributes({ label: '甘特图，项目进度管理界面' })"
  >
    <!-- 跳过导航链接（无障碍） -->
    <SkipLink target-id="gantt-container" />

    <!-- 工具栏 -->
    <GanttToolbar
      :view-mode="viewMode"
      :current-period-text="currentPeriodText"
      :current-zoom-label="currentZoomLabel"
      :show-dependencies="showDependencies"
      :show-critical-path="showCriticalPath"
      :show-baseline="showBaseline"
      :group-mode="groupMode"
      :search-keyword="searchKeyword"
      :is-fullscreen="isFullscreen"
      :is-saving="saving"
      :has-unsaved-changes="hasUnsavedChanges"
      :is-mobile="isMobile"
      @navigate-date="navigateDate"
      @go-today="goToToday"
      @view-mode-change="handleViewModeChange"
      @zoom-in="zoomIn"
      @zoom-out="zoomOut"
      @zoom-reset="zoomReset"
      @toggle-dependencies="toggleDependencies"
      @toggle-critical-path="toggleCriticalPath"
      @toggle-baseline="toggleBaseline"
      @open-resource-management="openResourceManagement"
      @group-change="handleGroupChange"
      @search="handleSearch"
      @export-png="handleExportPNG"
      @export-pdf="handleExportPDF"
      @auto-fit="autoFitContainer"
      @add-task="handleAddTask"
      @refresh="handleRefresh"
      @toggle-fullscreen="toggleFullscreen"
      @save-all="handleSaveAll"
      @show-shortcuts="showShortcutHelp = true"
    >
      <!-- 主题切换按钮插槽 -->
      <template #actions>
        <ThemeToggle
          v-if="!isMobile"
          :default-theme="userThemePreference"
          @theme-change="handleThemeChange"
          ref="themeToggleRef"
        />

        <!-- 快捷键帮助按钮 -->
        <el-button
          :icon="QuestionFilled"
          circle
          size="small"
          @click="showShortcutHelp = true"
          title="键盘快捷键 (?)"
        />
      </template>
    </GanttToolbar>

    <!-- 统计信息 -->
    <GanttStats :stats="taskStats" />

    <!-- 桌面/平板视图 -->
    <div
      v-if="!isMobile"
      id="gantt-container"
      class="gantt-container"
      :class="{ 'is-resizing': isResizing }"
      ref="ganttContainer"
      :style="containerStyle"
    >
      <!-- 调整大小手柄 - 右下角 -->
      <div
        class="resize-handle resize-handle-corner"
        @mousedown="handleResizeStart"
        title="拖动调整大小"
        v-bind="getAriaAttributes({ label: '调整容器大小' })"
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
          :view-mode="viewMode"
          :timeline-days="timelineDays"
          :timeline-weeks="timelineWeeks"
          :timeline-months="timelineMonths"
          :timeline-quarters="timelineQuarters"
          :day-width="dayWidth"
          :today-position="todayPosition"
        />
      </div>

      <!-- 重构后的内容区域 -->
      <GanttBody
        :tasks="filteredTasks"
        :grouped-tasks="groupedTasks"
        :selected-task="selectedTask"
        :timeline-days="timelineDays"
        :timeline-weeks="timelineWeeks"
        :timeline-months="timelineMonths"
        :timeline-quarters="timelineQuarters"
        :view-mode="viewMode"
        :day-width="dayWidth"
        :row-height="rowHeight"
        :task-height="taskHeight"
        :show-dependencies="showDependencies"
        :show-critical-path="showCriticalPath"
        :show-baseline="showBaseline"
        :group-mode="groupMode"
        :collapsed-groups="collapsedGroups"
        :search-keyword="searchKeyword"
        :is-dragging="isDragging"
        :dragged-task="draggedTask"
        :tooltip-visible="tooltipVisible"
        :tooltip-position="tooltipPosition"
        :tooltip-text="tooltipText"
        :is-resizing="isResizing"
        :is-mobile="isMobile"
        :use-virtual-scroll="useVirtualScroll"
        :container-size="containerSize"
        @row-click="handleRowClick"
        @task-click="handleTaskClick"
        @task-dblclick="handleTaskDblClick"
        @task-mousedown="handleTaskMouseDown"
        @toggle-group="toggleGroup"
        @context-menu="handleContextMenu"
        @add-task="handleAddTask"
        @cell-edit="handleCellEdit"
        @dependency-create="handleDependencyCreate"
        @resize-start="handleResizeStart"
      />
    </div>

    <!-- 移动端视图 -->
    <MobileGanttView
      v-else
      :view-mode="viewMode"
      :current-period-text="currentPeriodText"
      :current-zoom-label="currentZoomLabel"
      :show-dependencies="showDependencies"
      :show-critical-path="showCriticalPath"
      :show-baseline="showBaseline"
      :group-mode="groupMode"
      :is-fullscreen="isFullscreen"
      :is-saving="saving"
      :has-unsaved-changes="hasUnsavedChanges"
      :filtered-tasks="filteredTasks"
      :grouped-tasks="groupedTasks"
      :selected-task="selectedTask"
      :collapsed-groups="collapsedGroups"
      :search-keyword="searchKeyword"
      :timeline-days="timelineDays"
      :timeline-weeks="timelineWeeks"
      :timeline-months="timelineMonths"
      :day-width="dayWidth"
      :row-height="rowHeight"
      @navigate-date="navigateDate"
      @go-today="goToToday"
      @view-mode-change="handleViewModeChange"
      @zoom-in="zoomIn"
      @zoom-out="zoomOut"
      @toggle-dependencies="toggleDependencies"
      @toggle-critical-path="toggleCriticalPath"
      @toggle-baseline="toggleBaseline"
      @row-click="handleRowClick"
      @task-click="handleTaskClick"
      @task-dblclick="handleTaskDblClick"
      @toggle-group="toggleGroup"
      @context-menu="handleContextMenu"
      @add-task="handleAddTask"
      @save-all="handleSaveAll"
      @edit-task="handleEditTaskFromDrawer"
      @duplicate-task="handleDuplicateTask"
      @delete-task="handleDeleteTask"
    />

    <!-- 图例 -->
    <GanttLegend
      :show-critical-path="showCriticalPath"
      :show-baseline="showBaseline"
    />

    <!-- 屏幕阅读器通知 -->
    <A11yAnnouncer :message="a11yAnnouncement" />

    <!-- 快捷键帮助对话框 -->
    <ShortcutHelpPanel v-model="showShortcutHelp" />

    <!-- 右键菜单 -->
    <teleport to="body">
      <GanttContextMenu
        v-model:visible="contextMenuVisible"
        :task="contextMenuTask"
        :position="contextMenuPosition"
        @add-subtask="handleAddSubtask"
        @convert-milestone="handleConvertToMilestone"
        @add-dependency="handleAddDependency"
        @view-dependencies="handleViewDependencies"
        @allocate-resources="handleAllocateResources"
        @edit="handleContextMenuEdit"
        @duplicate="handleContextMenuDuplicate"
        @delete="handleContextMenuDelete"
        ref="contextMenuRef"
      />
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

    <!-- 任务详情侧边栏 -->
    <TaskDetailDrawer
      v-model:visible="taskDetailVisible"
      :task="selectedTask"
      @edit="handleEditTaskFromDrawer"
      @duplicate="handleDuplicateTask"
      @delete="handleDeleteTask"
    />

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
import { QuestionFilled } from '@element-plus/icons-vue'

// Import existing components and store
import { ganttStore } from '@/stores/ganttStore'
import { useGanttDrag } from '@/composables/useGanttDrag'

// Import NEW refactored components
import GanttBody from '@/components/gantt/core/GanttBody.vue'
import MobileGanttView from '@/components/gantt/mobile/MobileGanttView.vue'

// Import NEW accessibility and UX components
import ThemeToggle from '@/components/common/ThemeToggle.vue'
import ShortcutHelpPanel from '@/components/common/ShortcutHelpPanel.vue'
import A11yAnnouncer from '@/components/common/A11yAnnouncer.vue'
import SkipLink from '@/components/common/SkipLink.vue'

// Import existing subcomponents
import GanttToolbar from './GanttToolbar.vue'
import GanttStats from './GanttStats.vue'
import GanttHeader from './GanttHeader.vue'
import GanttLegend from './GanttLegend.vue'
import GanttContextMenu from './GanttContextMenu.vue'
import TaskDetailDrawer from './TaskDetailDrawer.vue'
import TaskEditDialog from './TaskEditDialog.vue'
import ResourceAllocationDialog from './ResourceAllocationDialog.vue'
import ResourceManagementDialog from './ResourceManagementDialog.vue'
import GanttStatusBar from './GanttStatusBar.vue'

// Import NEW composables
import { useTheme } from '@/composables/useTheme'
import { useGanttKeyboard } from '@/composables/useGanttKeyboard'
import { useGanttTooltip } from '@/composables/useGanttTooltip'
import { useGanttSelection } from '@/composables/useGanttSelection'
import { useBreakpoint } from '@/composables/useBreakpoint'
import { useAria } from '@/composables/useA11y'
import { useScreenReader } from '@/composables/useA11y'

// Utils
import { progressApi } from '@/api'
import eventBus, { GanttEvents } from '@/utils/eventBus'

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
  }
})

const emit = defineEmits(['task-updated', 'task-selected'])

// ==================== 初始化 Store ====================
const store = ganttStore
const { state, getters, actions } = store
actions.setProject(props.projectId, props.projectName)

// ==================== NEW: Theme Management ====================
const { mode: themeMode, setTheme: setThemeMode } = useTheme({ defaultTheme: 'light' })
const userThemePreference = ref('light')
const showShortcutHelp = ref(false)
const themeToggleRef = ref(null)

function handleThemeChange(newMode) {
  userThemePreference.value = newMode
  // Announce theme change for screen readers
  announce(`主题已切换为${newMode === 'dark' ? '暗色' : newMode === 'light' ? '亮色' : '跟随系统'}`)
}

// ==================== NEW: Accessibility ====================
const { getAriaAttributes } = useAria()
const { announce } = useScreenReader()
const a11yAnnouncement = ref('')

// ==================== NEW: Responsive Detection ====================
const { isMobile, isTablet } = useBreakpoint()

// ==================== NEW: Virtual Scroll Toggle ====================
const useVirtualScroll = computed(() => {
  // Enable virtual scroll for 50+ tasks or on mobile
  return state.filteredTasks.length > 50 || isMobile.value
})

// ==================== NEW: Keyboard Shortcuts ====================
const keyboardHandlers = {
  onTaskSelect: (task) => actions.selectTask(task.id),
  onTaskEdit: (task) => {
    actions.selectTask(task.id)
    actions.openEditDialog(task)
  },
  onTaskDelete: (task) => {
    actions.selectTask(task.id)
    handleDeleteTask()
  },
  onCopy: () => {
    const task = getters.selectedTask.value
    if (task) handleContextMenuDuplicate(task)
  },
  onPaste: () => {
    // TODO: Implement paste
  },
  onUndo: () => {
    // TODO: Implement undo
  },
  onRedo: () => {
    // TODO: Implement redo
  },
  onSave: () => handleSaveAll(),
  onZoomIn: () => actions.zoomIn(),
  onZoomOut: () => actions.zoomOut(),
  onNavigateDate: (direction) => navigateDate(direction === 'prev' ? -1 : 1),
  onToggleDependencies: () => actions.toggleDependencies(),
  onToggleCriticalPath: () => actions.toggleCriticalPath(),
  onViewModeChange: (mode) => actions.setViewMode(mode)
}

// ==================== NEW: Selection Management ====================
const {
  isSelected,
  selectTask,
  clearSelection
} = useGanttSelection({
  tasks: computed(() => state.filteredTasks),
  mode: ref('single'),
  onSelectionChange: (selectedTasks) => {
    if (selectedTasks.length > 0) {
      emit('task-selected', selectedTasks[0])
    }
  }
})

// ==================== 组件引用 ====================
const ganttChartRef = ref(null)
const ganttContainer = ref(null)
const statusBarRef = ref(null)
const contextMenuRef = ref(null)

// 容器调整大小状态
const resizeDirection = ref(null)
const resizeStart = ref({ x: 0, y: 0, width: 0, height: 0 })

// 使用 store 中的状态
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
const searchKeyword = computed(() => state.searchKeyword)
const taskDetailVisible = computed(() => state.taskDetailVisible)
const selectedTask = computed(() => getters.selectedTask)
const tempLineEnd = computed(() => state.tempLineEnd)
const isResizing = computed(() => state.isResizing)
const containerSize = computed(() => state.containerSize)
const editDialogVisible = computed(() => state.editDialogVisible)
const editingTask = computed(() => state.editingTask)
const contextMenuVisible = computed(() => state.contextMenuVisible)
const contextMenuTask = computed(() => state.contextMenuTask)
const contextMenuPosition = computed(() => state.contextMenuPosition)
const resourceDialogVisible = computed(() => state.resourceDialogVisible)
const currentTaskForResource = computed(() => state.currentTaskForResource)
const resourceManagementDialogVisible = computed(() => state.resourceManagementDialogVisible)
const resourceLibrary = computed(() => state.resourceLibrary)

// 拖拽状态
let isDragging, draggedTask, tooltipVisible, tooltipPosition, tooltipText, startDrag, cancelDrag

// ==================== 拖拽事件处理 ====================
const handleDragChange = (preview) => {
  // 拖拽过程中的回调
}

const handleDragEnd = async (newTask, originalTask) => {
  try {
    actions.endDrag(newTask, originalTask)
    announce('任务位置已更改，记得保存')
    statusBarRef.value?.showStatus('任务位置已更改，记得保存', 'info', 2000)
  } catch (error) {
    console.error('处理任务拖拽失败:', error)
    ElMessage.error('处理任务拖拽失败')
  }
}

// 拖拽功能
const dragResult = useGanttDrag({
  dayWidth,
  timelineDays: computed(() => getters.timelineDays.value),
  onDragEnd: handleDragEnd,
  onDragChange: handleDragChange,
  enableRAF: true,
  throttleMs: 16
})

isDragging = dragResult.isDragging
draggedTask = dragResult.draggedTask
tooltipPosition = dragResult.tooltipPosition
tooltipVisible = dragResult.tooltipVisible
tooltipText = dragResult.tooltipText
startDrag = dragResult.startDrag
cancelDrag = dragResult.cancelDrag

// ==================== 使用 Store Getters ====================
const formattedTasks = computed(() => state.tasks)
const filteredTasks = computed(() => state.filteredTasks)
const groupedTasks = computed(() => getters.groupedTasks.value)
const taskStats = computed(() => getters.taskStats.value)
const timelineDays = computed(() => getters.timelineDays.value)
const timelineWeeks = computed(() => getters.timelineWeeks.value)
const timelineMonths = computed(() => getters.timelineMonths.value)
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

// ==================== 计算属性 ====================
const chartClasses = computed(() => ({
  'is-fullscreen': isFullscreen.value,
  'is-mobile': isMobile.value,
  'is-tablet': isTablet.value,
  [`theme-${themeMode.value}`]: true
}))

// ==================== 任务操作 ====================
const handleRowClick = (task) => {
  actions.selectTask(task.id)
  emit('task-selected', task)
}

const handleTaskClick = (task) => {
  if (!task) return
  actions.selectTask(task.id)
  emit('task-selected', task)
}

const handleTaskDblClick = (task) => {
  if (!task) return
  actions.selectTask(task.id)
  actions.openEditDialog(task)
}

const handleTaskMouseDown = (event, taskOrId) => {
  if (event.button !== 0) return
  const taskId = typeof taskOrId === 'object' ? taskOrId.id : taskOrId
  const task = filteredTasks.value.find(t => t.id === taskId)
  if (!task) return
  startDrag(event, task)
}

const handleSearch = (value) => {
  state.searchKeyword = value
  actions.filterTasks(value)
}

const handleDependencyCreate = (taskOrId) => {
  const taskId = typeof taskOrId === 'object' ? taskOrId.id : taskOrId
  const task = filteredTasks.value.find(t => t.id === taskId)

  if (!task) return

  if (!state.isCreatingDependency) {
    actions.startDependencyCreation(task)
  } else {
    actions.completeDependencyCreation(task)
  }
}

const handleCellEdit = async ({ taskId, updateData }) => {
  try {
    const task = state.filteredTasks.find(t => t.id === taskId)
    if (!task) return

    statusBarRef.value?.showStatus('正在保存...', 'loading')

    const taskData = {
      project_id: props.projectId,
      ...updateData
    }

    await progressApi.update(taskId, taskData)
    ElMessage.success('保存成功')
    statusBarRef.value?.showStatus('保存成功', 'success', 1500)
    emit('task-updated', { ...task, ...taskData })
    announce('任务已保存')
  } catch (error) {
    console.error('保存编辑失败:', error)
    const errorMsg = error.response?.data?.error || error.response?.data?.message || error.message || '保存失败'
    ElMessage.error(errorMsg)
    statusBarRef.value?.showStatus('保存失败', 'error', 2000)
  }
}

// ==================== 容器调整大小 ====================
const handleResizeStart = (event) => {
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

// ==================== 视图控制 ====================
const navigateDate = (direction) => {
  const offset = direction * (state.viewMode === 'day' ? 7 : state.viewMode === 'week' ? 4 : 1)
  const newDayWidth = Math.max(state.dayWidth + offset, state.VIEW_CONFIG[state.viewMode].minWidth)
  state.dayWidth = newDayWidth
}

const goToToday = () => {
  ElMessage.success('已回到今天')
  announce('已回到今天')
}

const zoomIn = () => actions.zoomIn()
const zoomOut = () => actions.zoomOut()
const zoomReset = () => actions.zoomReset()

const handleViewModeChange = (newMode) => actions.setViewMode(newMode)
const toggleDependencies = () => actions.toggleDependencies()
const toggleCriticalPath = () => actions.toggleCriticalPath()
const toggleBaseline = () => actions.toggleBaseline()
const openResourceManagement = () => actions.openResourceManagementDialog()
const handleGroupChange = (newGroup) => actions.setGroupMode(newGroup)
const toggleGroup = (groupName) => actions.toggleGroup(groupName)
const handleRefresh = () => emit('task-updated', null)

// 全屏切换
const toggleFullscreen = () => {
  if (!document.fullscreenElement) {
    if (ganttChartRef.value?.requestFullscreen) {
      ganttChartRef.value.requestFullscreen()
    }
    actions.setFullscreen(true)
    announce('已进入全屏模式')
  } else {
    if (document.exitFullscreen) {
      document.exitFullscreen()
    }
    actions.setFullscreen(false)
    announce('已退出全屏模式')
  }
}

const handleFullscreenChange = () => {
  actions.setFullscreen(!!document.fullscreenElement)
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
    announce('PNG图片已导出')
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
    announce('PDF文件已导出')
  } catch (error) {
    console.error('导出PDF失败:', error)
    ElMessage.error('导出PDF失败')
  }
}

// ==================== 右键菜单 ====================
const handleContextMenu = (eventOrData) => {
  let event, task

  if (eventOrData.event && eventOrData.task) {
    event = eventOrData.event
    task = eventOrData.task
  } else {
    event = eventOrData
    task = null
    const taskRow = event.target.closest('.table-row')
    if (taskRow) {
      const taskId = taskRow.getAttribute('data-task-id')
      if (taskId) {
        task = state.filteredTasks.find(t => t.id == taskId)
      }
    }
  }

  if (!task) return

  event.preventDefault()
  event.stopPropagation()

  actions.showContextMenu(task, { x: event.clientX, y: event.clientY })

  nextTick(() => {
    contextMenuRef.value?.open()
  })
}

const handleContextMenuEdit = (task) => handleTaskDblClick(task)

const handleContextMenuDuplicate = async (task) => {
  try {
    const newTask = {
      project_id: props.projectId,
      task_name: `${task.name} (副本)`,
      start_date: task.start,
      end_date: task.end,
      progress: 0,
      priority: task.priority || 'medium',
      description: ''
    }

    await progressApi.create(newTask)
    ElMessage.success('任务已复制')
    announce('任务已复制')
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

// ==================== 子任务管理 ====================
const handleAddSubtask = async (parentTask) => {
  try {
    const newTask = {
      project_id: props.projectId,
      task_name: `${parentTask.name} - 子任务`,
      start_date: parentTask.start,
      end_date: parentTask.end,
      progress: 0,
      priority: parentTask.priority || 'medium',
      description: ''
    }

    if (parentTask.id) {
      newTask.parent_id = parentTask.id
    }

    const response = await progressApi.create(newTask)
    ElMessage.success('子任务已添加')
    announce('子任务已添加')

    if (parentTask.id) {
      try {
        await progressApi.updateParentProgress(parentTask.id)
      } catch (updateError) {
        console.error('更新父任务进度失败:', updateError)
      }
    }
  } catch (error) {
    console.error('添加子任务失败:', error)
    const errorMsg = error.response?.data?.error || error.response?.data?.message || error.message || '添加子任务失败'
    ElMessage.error(errorMsg)
  }
}

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

    await progressApi.updateTask(task.id, {
      start_date: task.end,
      end_date: task.end,
      is_milestone: true,
      duration: 0
    })

    ElMessage.success('已转换为里程碑')
    announce('已转换为里程碑')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('转换里程碑失败:', error)
      ElMessage.error('转换失败')
    }
  }
}

// ==================== 依赖关系管理 ====================
const handleAddDependency = (task) => {
  actions.startDependencyCreation(task)
}

const handleViewDependencies = async (task) => {
  try {
    const response = await progressApi.getDependencies(task.id)
    const deps = response.data || []

    if (deps.length === 0) {
      ElMessage.info('该任务暂无依赖关系')
      return
    }

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
const handleAllocateResources = (task) => {
  actions.openResourceDialog(task)
}

const handleResourceSaved = () => {
  emit('task-updated', null)
}

const handleResourceRefresh = () => {
  actions.loadResources()
}

// ==================== 添加/编辑/删除任务 ====================
const handleAddTask = () => {
  actions.openEditDialog(null)
}

const handleSaveTask = async (formData) => {
  try {
    const isEdit = !!state.editingTask
    const statusMsg = isEdit ? '正在更新任务...' : '正在创建任务...'
    statusBarRef.value?.showStatus(statusMsg, 'loading')

    const taskData = {
      project_id: props.projectId,
      task_name: formData.name,
      start_date: formData.start,
      end_date: formData.end,
      progress: formData.progress,
      priority: formData.priority,
      description: formData.notes,
      resources: formData.resources?.map(r => ({
        resource_id: r.resource_id || r.id,
        quantity: r.quantity,
        cost: r.cost,
        type: r.type
      })) || [],
      predecessor_ids: formData.predecessor_ids || [],
      successor_ids: formData.successor_ids || []
    }

    if (isEdit) {
      await progressApi.update(state.editingTask.id, taskData)
      ElMessage.success('任务更新成功')
      announce('任务已更新')
      statusBarRef.value?.showStatus('任务更新成功', 'success', 2000)
      emit('task-updated', { ...state.editingTask, ...taskData })
    } else {
      await progressApi.create(taskData)
      ElMessage.success('任务创建成功')
      announce('任务已创建')
      statusBarRef.value?.showStatus('任务创建成功', 'success', 2000)
      emit('task-updated', taskData)
    }

    actions.closeEditDialog()
  } catch (error) {
    console.error('保存任务失败:', error)
    ElMessage.error('保存任务失败')
    statusBarRef.value?.showStatus('保存任务失败', 'error', 2000)
  }
}

const handleEditTaskFromDrawer = () => {
  const task = getters.selectedTask.value
  if (!task) return
  actions.closeTaskDetail()
  actions.openEditDialog(task)
}

const handleDuplicateTask = async () => {
  const task = getters.selectedTask.value
  if (!task) return

  try {
    const newTask = {
      project_id: props.projectId,
      task_name: `${task.name} (副本)`,
      start_date: task.start,
      end_date: task.end,
      progress: 0,
      priority: task.priority || 'medium',
      description: ''
    }

    await progressApi.create(newTask)
    ElMessage.success('任务已复制')
    announce('任务已复制')
    emit('task-updated', newTask)
  } catch (error) {
    console.error('复制任务失败:', error)
    ElMessage.error('复制任务失败')
  }
}

const handleDeleteTask = async () => {
  const task = getters.selectedTask.value
  if (!task) return

  try {
    await ElMessageBox.confirm(
      `确定要删除任务"${task.name}"吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    statusBarRef.value?.showStatus('正在删除任务...', 'loading')
    await progressApi.delete(task.id)
    ElMessage.success('任务已删除')
    announce('任务已删除')
    statusBarRef.value?.showStatus('任务已删除', 'success', 2000)
    actions.closeTaskDetail()
    emit('task-updated', task)
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除任务失败:', error)
      ElMessage.error('删除任务失败')
      statusBarRef.value?.showStatus('删除任务失败', 'error', 2000)
    }
  }
}

const autoFitContainer = () => actions.autoFit()

// ==================== 键盘快捷键 ====================
const handleKeydown = (event) => {
  if (state.editDialogVisible) return

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

  // Delete: 删除任务
  if (event.key === 'Delete' && getters.selectedTask.value) {
    event.preventDefault()
    handleDeleteTask()
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
    if (isDragging.value) {
      cancelDrag()
    }
    if (showShortcutHelp.value) {
      showShortcutHelp.value = false
    }
  }

  // Ctrl/Cmd + F: 聚焦搜索框
  if ((event.ctrlKey || event.metaKey) && event.key === 'f') {
    event.preventDefault()
    document.querySelector('.gantt-toolbar input')?.focus()
  }

  // ?: 显示快捷键帮助
  if (event.key === '?' || (event.key === '/' && (event.ctrlKey || event.metaKey))) {
    event.preventDefault()
    showShortcutHelp.value = !showShortcutHelp.value
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
const handleBeforeUnload = (e) => {
  if (state.hasUnsavedChanges) {
    e.preventDefault()
    e.returnValue = '您有未保存的更改，确定要离开吗？'
    return '您有未保存的更改，确定要离开吗？'
  }
}

onMounted(() => {
  state.scheduleData = props.scheduleData
  actions.formatTasks()
  actions.loadResources()

  // 添加全屏状态变化监听
  document.addEventListener('fullscreenchange', handleFullscreenChange)
  document.addEventListener('webkitfullscreenchange', handleFullscreenChange)
  document.addEventListener('mozfullscreenchange', handleFullscreenChange)
  document.addEventListener('MSFullscreenChange', handleFullscreenChange)

  // 添加页面离开前的警告
  window.addEventListener('beforeunload', handleBeforeUnload)

  // 加载用户主题偏好
  try {
    const savedTheme = localStorage.getItem('gantt-theme-preference')
    if (savedTheme) {
      userThemePreference.value = savedTheme
      setThemeMode(savedTheme)
    }
  } catch (e) {
    console.warn('Failed to load theme preference:', e)
  }
})

onUnmounted(() => {
  if (isDragging.value) {
    cancelDrag()
  }
  if (state.isResizing) {
    handleResizeEnd()
  }

  document.removeEventListener('fullscreenchange', handleFullscreenChange)
  document.removeEventListener('webkitfullscreenchange', handleFullscreenChange)
  document.removeEventListener('mozfullscreenchange', handleFullscreenChange)
  document.removeEventListener('MSFullscreenChange', handleFullscreenChange)
  window.removeEventListener('beforeunload', handleBeforeUnload)
  actions.stopAutoSave()
})

// 监听scheduleData变化
watch(
  () => props.scheduleData,
  (newVal) => {
    state.scheduleData = newVal
    actions.formatTasks()

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
</script>

<style scoped>
.gantt-chart {
  height: 100vh;
  max-height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--color-bg-primary, #fff);
  border: 1px solid var(--color-border-base, #dcdfe6);
  border-radius: var(--radius-md, 4px);
  overflow: hidden;
  outline: none;
}

/* 甘特图容器 */
.gantt-container {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  position: relative;
  min-width: 600px;
  min-height: 0;
}

/* 表头容器 */
.gantt-header-container {
  flex-shrink: 0;
  position: sticky;
  top: 0;
  z-index: 100;
  background: var(--color-bg-primary, #fff);
  border-bottom: 1px solid var(--color-border-base, #dcdfe6);
}

/* 调整大小手柄 */
.resize-handle {
  position: absolute;
  z-index: 9999;
  transition: background-color var(--transition-fast, 0.2s);
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

/* 全屏样式 */
.gantt-chart.is-fullscreen {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  width: 100vw;
  height: 100vh;
  z-index: 9998;
  border-radius: 0;
  overflow: hidden;
}

/* 响应式 */
@media (max-width: 768px) {
  .gantt-chart {
    border: none;
    border-radius: 0;
  }

  .gantt-container {
    min-width: 100%;
  }
}

/* 主题样式 */
.theme-dark {
  --color-bg-primary: #1a1a1a;
  --color-bg-secondary: #2a2a2a;
  --color-text-primary: #ffffff;
  --color-border-base: #4c4d4f;
}

/* 减少动画模式支持 */
@media (prefers-reduced-motion: reduce) {
  .resize-handle {
    transition: none;
  }
}
</style>
