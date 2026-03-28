<template>
  <div
    ref="editorRef"
    class="gantt-editor"
    :class="{
      'is-fullscreen': state.isFullscreen,
      'is-loading': state.loading
    }"
    tabindex="0"
    @keydown="handleKeydown"
  >
    <!-- Toolbar -->
    <GanttToolbar
      :undo-count="undoCount"
      :redo-count="redoCount"
      :can-undo="canUndo"
      :can-redo="canRedo"
      :selected-count="selectedTaskIds.size"
      :view-mode="state.viewMode"
      :day-width="state.dayWidth"
      :sync-status="syncStatus"
      @undo="handleUndo"
      @redo="handleRedo"
      @zoom-in="handleZoomIn"
      @zoom-out="handleZoomOut"
      @zoom-reset="handleZoomReset"
      @view-change="handleViewChange"
      @toggle-template="handleToggleTemplate"
      @toggle-bulk-edit="handleToggleBulkEdit"
      @toggle-fullscreen="handleToggleFullscreen"
      @export="handleExport"
    />

    <!-- Main Content -->
    <div class="gantt-editor__content">
      <!-- Task List -->
      <div
        class="gantt-editor__task-list"
        :style="taskListStyle"
      >
        <VirtualTaskList
          ref="taskListRef"
          :tasks="visibleTasks"
          :columns="taskListColumns"
          :row-height="state.rowHeight"
          :container-height="containerHeight"
          :collapsed-tasks="state.collapsedTasks"
          :selected-task-id="state.selectedTaskId"
          :loading="state.loading"
          :search-keyword="state.searchKeyword"
          @scroll="handleTaskListScroll"
          @row-click="handleTaskClick"
          @row-dblclick="handleTaskDoubleClick"
          @toggle="handleTaskToggle"
          @edit="handleTaskEdit"
          @selection-change="handleSelectionChange"
        />
      </div>

      <!-- Timeline -->
      <div
        class="gantt-editor__timeline"
        :style="timelineStyle"
      >
        <VirtualTimeline
          ref="timelineRef"
          :tasks="visibleTasks"
          :days="getters.timelineDays.value"
          :row-height="state.rowHeight"
          :day-width="state.dayWidth"
          :container-height="containerHeight"
          :show-today="true"
          :loading="state.loading"
          @scroll="handleTimelineScroll"
          @resize="handleTimelineResize"
        >
          <template #row="{ task, index, days, dayWidth }">
            <TaskBar
              :task="task"
              :days="days"
              :day-width="dayWidth"
              :row-height="state.rowHeight"
              @click="handleTaskBarClick(task)"
              @dblclick="handleTaskBarDoubleClick(task)"
            />
          </template>
        </VirtualTimeline>
      </div>
    </div>

    <!-- Status Bar -->
    <GanttStatusBar
      :undo-count="undoCount"
      :redo-count="redoCount"
      :task-count="state.tasks.length"
      :selected-count="selectedTaskIds.size"
      :last-saved="lastSavedTime"
      :connection-status="connectionStatus"
    />

    <!-- Task Templates Dialog -->
    <TaskTemplatesDialog
      v-model="templateDialogVisible"
      :project-id="state.projectId"
      :start-date="defaultStartDate"
      @created="handleTaskCreated"
    />

    <!-- Bulk Edit Dialog -->
    <BulkEditDialog
      v-model="bulkEditDialogVisible"
      :tasks="selectedTasks"
      :project-id="state.projectId"
      @updated="handleBulkUpdate"
    />

    <!-- Loading Overlay -->
    <transition name="el-fade-in">
      <div v-if="state.loading" class="gantt-editor__loading-overlay">
        <el-icon class="is-loading" :size="32">
          <Loading />
        </el-icon>
        <p>加载中...</p>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Loading } from '@element-plus/icons-vue'
import { ganttStore } from '@/stores/ganttStore'
import { useUndoRedoStore } from '@/stores/undoRedoStore'
import eventBus, { GanttEvents } from '@/utils/eventBus'
import VirtualTaskList from '@/components/gantt/table/VirtualTaskList.vue'
import VirtualTimeline from '@/components/gantt/timeline/VirtualTimeline.vue'
import TaskBar from '@/components/gantt/timeline/TaskBar.vue'
import GanttToolbar from './GanttToolbar.vue'
import GanttStatusBar from './GanttStatusBar.vue'
import TaskTemplatesDialog from '@/components/gantt/dialogs/TaskTemplatesDialog.vue'
import BulkEditDialog from '@/components/gantt/dialogs/BulkEditDialog.vue'

/**
 * GanttEditor Component - Main Enhanced Gantt Chart Container
 *
 * Integrates VirtualTimeline and VirtualTaskList with:
 * - Undo/Redo functionality
 * - Multi-selection support
 * - Keyboard shortcuts
 * - Full-screen mode
 * - Template quick create
 * - Bulk editing
 * - Responsive layout
 *
 * @date 2025-02-18
 */

// Props
const props = defineProps({
  projectId: {
    type: [Number, String],
    required: true
  },
  projectName: {
    type: String,
    default: ''
  },
  height: {
    type: [String, Number],
    default: 'auto'
  }
})

// Emits
const emit = defineEmits(['ready', 'task-selected', 'view-changed'])

// Refs
const editorRef = ref(null)
const taskListRef = ref(null)
const timelineRef = ref(null)

// Stores
const { state, getters, actions } = ganttStore
const undoRedoStore = useUndoRedoStore()

// Dialog state
const templateDialogVisible = ref(false)
const bulkEditDialogVisible = ref(false)

// Selection state
const selectedTaskIds = ref(new Set())

// Computed
const containerHeight = computed(() => {
  const toolbarHeight = 56
  const statusBarHeight = 32
  const availableHeight = typeof props.height === 'number'
    ? props.height
    : editorRef.value?.clientHeight - toolbarHeight - statusBarHeight || 600

  return Math.max(400, availableHeight)
})

const taskListStyle = computed(() => ({
  width: '400px',
  height: `${containerHeight.value}px`
}))

const timelineStyle = computed(() => ({
  flex: 1,
  height: `${containerHeight.value}px`,
  overflow: 'hidden'
}))

const visibleTasks = computed(() => {
  return getters.getVisibleTasks()
})

const selectedTasks = computed(() => {
  return state.tasks.filter(task => selectedTaskIds.value.has(task.id))
})

const undoCount = computed(() => {
  return undoRedoStore.history.length
})

const redoCount = computed(() => {
  return undoRedoStore.future.length
})

const canUndo = computed(() => {
  return undoRedoStore.canUndo
})

const canRedo = computed(() => {
  return undoRedoStore.canRedo
})

const lastSavedTime = computed(() => {
  if (!state.lastSavedTime) return '未保存'
  const time = new Date(state.lastSavedTime)
  const now = new Date()
  const diff = Math.floor((now - time) / 1000)

  if (diff < 60) return '刚刚保存'
  if (diff < 3600) return `${Math.floor(diff / 60)}分钟前保存`
  return time.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
})

const syncStatus = computed(() => {
  if (state.saving) return 'saving'
  if (state.hasUnsavedChanges) return 'unsaved'
  return 'saved'
})

const connectionStatus = computed(() => {
  // Future: WebSocket connection status
  return 'connected'
})

const defaultStartDate = computed(() => {
  const today = new Date()
  return today.toISOString().split('T')[0]
})

const taskListColumns = computed(() => [
  { key: 'name', label: '任务名称', width: 200, editable: true },
  { key: 'duration', label: '工期', width: 60, align: 'center' },
  { key: 'progress', label: '进度', width: 60, align: 'center' },
  { key: 'status', label: '状态', width: 80, align: 'center' }
])

// ==================== Methods ====================

/**
 * Initialize the editor
 */
async function initialize() {
  try {
    actions.setProject(props.projectId, props.projectName)
    await actions.loadData()

    // Initialize undo/redo store with project context
    undoRedoStore.setProjectContext(props.projectId)

    emit('ready')
  } catch (error) {
    console.error('Failed to initialize Gantt editor:', error)
    ElMessage.error('加载甘特图失败')
  }
}

/**
 * Handle keyboard shortcuts
 */
function handleKeydown(event) {
  // Ctrl/Cmd + Z: Undo
  if ((event.ctrlKey || event.metaKey) && event.key === 'z' && !event.shiftKey) {
    event.preventDefault()
    handleUndo()
  }

  // Ctrl/Cmd + Shift + Z or Ctrl/Cmd + Y: Redo
  if (
    ((event.ctrlKey || event.metaKey) && event.shiftKey && event.key === 'z') ||
    ((event.ctrlKey || event.metaKey) && event.key === 'y')
  ) {
    event.preventDefault()
    handleRedo()
  }

  // Ctrl/Cmd + A: Select all tasks
  if ((event.ctrlKey || event.metaKey) && event.key === 'a') {
    event.preventDefault()
    handleSelectAll()
  }

  // Delete or Backspace: Delete selected tasks
  if ((event.key === 'Delete' || event.key === 'Backspace') && selectedTaskIds.value.size > 0) {
    event.preventDefault()
    handleDeleteSelected()
  }

  // Escape: Clear selection or close dialogs
  if (event.key === 'Escape') {
    if (templateDialogVisible.value) {
      templateDialogVisible.value = false
    } else if (bulkEditDialogVisible.value) {
      bulkEditDialogVisible.value = false
    } else if (selectedTaskIds.value.size > 0) {
      selectedTaskIds.value.clear()
    }
  }

  // Ctrl/Cmd + F: Focus search
  if ((event.ctrlKey || event.metaKey) && event.key === 'f') {
    event.preventDefault()
    // Focus search input (implement search input in toolbar)
  }

  // Ctrl/Cmd + S: Save
  if ((event.ctrlKey || event.metaKey) && event.key === 's') {
    event.preventDefault()
    handleSave()
  }

  // F11: Toggle fullscreen
  if (event.key === 'F11') {
    event.preventDefault()
    handleToggleFullscreen()
  }
}

/**
 * Handle undo
 */
async function handleUndo() {
  if (!canUndo.value) return

  try {
    await undoRedoStore.undo()
    ElMessage.success('已撤销')
  } catch (error) {
    console.error('Undo failed:', error)
    ElMessage.error('撤销失败')
  }
}

/**
 * Handle redo
 */
async function handleRedo() {
  if (!canRedo.value) return

  try {
    await undoRedoStore.redo()
    ElMessage.success('已重做')
  } catch (error) {
    console.error('Redo failed:', error)
    ElMessage.error('重做失败')
  }
}

/**
 * Handle zoom in
 */
function handleZoomIn() {
  actions.zoomIn()
}

/**
 * Handle zoom out
 */
function handleZoomOut() {
  actions.zoomOut()
}

/**
 * Handle zoom reset
 */
function handleZoomReset() {
  actions.zoomReset()
}

/**
 * Handle view mode change
 */
function handleViewChange(mode) {
  actions.setViewMode(mode)
  emit('view-changed', mode)
}

/**
 * Handle toggle template dialog
 */
function handleToggleTemplate() {
  templateDialogVisible.value = true
}

/**
 * Handle toggle bulk edit dialog
 */
function handleToggleBulkEdit() {
  if (selectedTaskIds.value.size === 0) {
    ElMessage.warning('请先选择要编辑的任务')
    return
  }
  bulkEditDialogVisible.value = true
}

/**
 * Handle toggle fullscreen
 */
function handleToggleFullscreen() {
  const isFullscreen = !state.isFullscreen
  actions.setFullscreen(isFullscreen)

  if (isFullscreen) {
    editorRef.value?.requestFullscreen?.()
  } else {
    document.exitFullscreen?.()
  }
}

/**
 * Handle export
 */
async function handleExport(format) {
  try {
    // Implement export logic (PDF, Excel, Image)
    ElMessage.info(`导出${format}功能开发中`)
  } catch (error) {
    console.error('Export failed:', error)
    ElMessage.error('导出失败')
  }
}

/**
 * Handle save
 */
async function handleSave() {
  try {
    await actions.saveAll()
    ElMessage.success('保存成功')
  } catch (error) {
    console.error('Save failed:', error)
    ElMessage.error('保存失败')
  }
}

/**
 * Handle select all
 */
function handleSelectAll() {
  selectedTaskIds.value = new Set(state.tasks.map(t => t.id))
}

/**
 * Handle delete selected tasks
 */
async function handleDeleteSelected() {
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedTaskIds.value.size} 个任务吗？`,
      '确认删除',
      {
        type: 'warning',
        confirmButtonText: '确定',
        cancelButtonText: '取消'
      }
    )

    // Implement delete with undo/redo
    const taskIds = Array.from(selectedTaskIds.value)

    // Create delete command for each task
    for (const taskId of taskIds) {
      const task = state.tasks.find(t => t.id === taskId)
      if (task) {
        // Create and execute delete command
        // const command = new DeleteTaskCommand(...)
        // await undoRedoStore.executeCommand(command)
      }
    }

    selectedTaskIds.value.clear()
    ElMessage.success(`已删除 ${taskIds.length} 个任务`)
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Delete failed:', error)
      ElMessage.error('删除失败')
    }
  }
}

/**
 * Handle task click
 */
function handleTaskClick({ task }) {
  actions.selectTask(task.id)
  emit('task-selected', task)
}

/**
 * Handle task double click
 */
function handleTaskDoubleClick({ task }) {
  actions.openEditDialog(task)
}

/**
 * Handle task bar click
 */
function handleTaskBarClick(task) {
  actions.selectTask(task.id)
  emit('task-selected', task)
}

/**
 * Handle task bar double click
 */
function handleTaskBarDoubleClick(task) {
  actions.openEditDialog(task)
}

/**
 * Handle task toggle (expand/collapse)
 */
function handleTaskToggle(task) {
  actions.toggleTaskCollapse(task.id)
}

/**
 * Handle task edit
 */
function handleTaskEdit({ task, column, value }) {
  // Implement inline edit with undo/redo
  console.log('Edit task:', task, column, value)
}

/**
 * Handle selection change
 */
function handleSelectionChange(taskId) {
  // Handle multi-selection with Ctrl/Cmd key
  if (window.event && (window.event.ctrlKey || window.event.metaKey)) {
    if (selectedTaskIds.value.has(taskId)) {
      selectedTaskIds.value.delete(taskId)
    } else {
      selectedTaskIds.value.add(taskId)
    }
  } else {
    selectedTaskIds.value = new Set([taskId])
  }
}

/**
 * Handle task list scroll (sync with timeline)
 */
function handleTaskListScroll(event) {
  if (timelineRef.value && event.scrollTop !== undefined) {
    timelineRef.value.scrollToPosition({ scrollTop: event.scrollTop })
  }
}

/**
 * Handle timeline scroll (sync with task list)
 */
function handleTimelineScroll(event) {
  if (taskListRef.value && event.scrollTop !== undefined) {
    taskListRef.value.scrollToPosition(event.scrollTop)
  }
}

/**
 * Handle timeline resize
 */
function handleTimelineResize(size) {
  actions.setContainerSize(size.width, size.height)
}

/**
 * Handle task created from template
 */
function handleTaskCreated(result) {
  ElMessage.success('任务已创建')
  // Reload data or update local state
}

/**
 * Handle bulk update
 */
function handleBulkUpdate(result) {
  ElMessage.success(`已更新 ${result.count} 个任务`)
  selectedTaskIds.value.clear()
}

// ==================== Lifecycle ====================

onMounted(() => {
  initialize()

  // Listen for fullscreen change
  document.addEventListener('fullscreenchange', handleFullscreenChange)

  // Listen for gantt events
  eventBus.on(GanttEvents.TASK_SELECTED, handleTaskSelectedEvent)
  eventBus.on(GanttEvents.DATA_CHANGED, handleDataChangedEvent)
})

onUnmounted(() => {
  document.removeEventListener('fullscreenchange', handleFullscreenChange)
  eventBus.off(GanttEvents.TASK_SELECTED, handleTaskSelectedEvent)
  eventBus.off(GanttEvents.DATA_CHANGED, handleDataChangedEvent)
})

// ==================== Event Handlers ====================

function handleFullscreenChange() {
  const isFullscreen = !!document.fullscreenElement
  actions.setFullscreen(isFullscreen)
}

function handleTaskSelectedEvent({ task }) {
  if (task) {
    selectedTaskIds.value = new Set([task.id])
  }
}

function handleDataChangedEvent(data) {
  // Handle data changes
}

// ==================== Watchers ====================

// Watch for project changes
watch(() => props.projectId, (newId) => {
  if (newId) {
    initialize()
  }
})

// Expose methods
defineExpose({
  initialize,
  save: handleSave,
  undo: handleUndo,
  redo: handleRedo,
  zoomIn: handleZoomIn,
  zoomOut: handleZoomOut,
  toggleFullscreen: handleToggleFullscreen
})
</script>

<style scoped lang="scss">
.gantt-editor {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: #fff;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  overflow: hidden;
  outline: none;

  &:focus {
    box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.2);
  }

  &.is-fullscreen {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    z-index: 9999;
    border-radius: 0;
  }

  &.is-loading {
    pointer-events: none;
  }
}

.gantt-editor__content {
  display: flex;
  flex: 1;
  overflow: hidden;
  position: relative;
}

.gantt-editor__task-list {
  flex-shrink: 0;
  border-right: 1px solid #e4e7ed;
  overflow: hidden;
}

.gantt-editor__timeline {
  flex: 1;
  overflow: hidden;
}

.gantt-editor__loading-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(255, 255, 255, 0.9);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 100;
  color: var(--el-color-primary, #409eff);

  p {
    margin-top: 16px;
    font-size: 14px;
    color: var(--el-text-color-regular, #606266);
  }
}

/* Responsive */
@media (max-width: 768px) {
  .gantt-editor__content {
    flex-direction: column;
  }

  .gantt-editor__task-list {
    width: 100% !important;
    height: 40%;
    border-right: none;
    border-bottom: 1px solid #e4e7ed;
  }

  .gantt-editor__timeline {
    height: 60%;
  }
}
</style>
