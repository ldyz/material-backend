<template>
  <!-- Example: Integrating all three editing components in a Gantt chart view -->
  <div class="gantt-view-example">
    <!-- Toolbar with action buttons -->
    <div class="gantt-toolbar">
      <el-button-group>
        <el-button @click="openTemplatesDialog">
          <el-icon><Plus /></el-icon>
          从模板创建
        </el-button>
        <el-button
          :disabled="selectedTasks.length === 0"
          @click="openBulkEditDialog"
        >
          <el-icon><Edit /></el-icon>
          批量编辑 ({{ selectedTasks.length }})
        </el-button>
      </el-button-group>

      <div class="toolbar-info">
        <el-tag v-if="hasUnsavedChanges" type="warning">
          有未保存的更改
        </el-tag>
      </div>
    </div>

    <!-- Gantt Table with Editable Cells -->
    <div class="gantt-table-container">
      <table class="gantt-table">
        <thead>
          <tr>
            <th width="40">
              <el-checkbox
                v-model="selectAll"
                :indeterminate="isSomeSelected"
                @change="handleSelectAll"
              />
            </th>
            <th>任务名称</th>
            <th>状态</th>
            <th>优先级</th>
            <th>进度</th>
            <th>开始日期</th>
            <th>工期</th>
            <th>负责人</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="task in tasks"
            :key="task.id"
            :class="{ 'is-selected': selectedTasks.includes(task.id) }"
          >
            <td>
              <el-checkbox
                :model-value="selectedTasks.includes(task.id)"
                @change="toggleTaskSelection(task.id)"
              />
            </td>

            <!-- Task Name - Text Input -->
            <td>
              <EditableCell
                v-model="task.name"
                type="text"
                field="task_name"
                :task-id="task.id"
                :original-data="task"
                placeholder="任务名称"
                :rules="nameRules"
                @change="handleTaskChange"
              />
            </td>

            <!-- Status - Select Input -->
            <td>
              <EditableCell
                v-model="task.status"
                type="select"
                field="status"
                :task-id="task.id"
                :original-data="task"
                :options="statusOptions"
                :display-format="formatStatus"
                :display-class="`status-${task.status}`"
                @change="handleTaskChange"
              />
            </td>

            <!-- Priority - Select Input -->
            <td>
              <EditableCell
                v-model="task.priority"
                type="select"
                field="priority"
                :task-id="task.id"
                :original-data="task"
                :options="priorityOptions"
                :display-format="formatPriority"
                @change="handleTaskChange"
              />
            </td>

            <!-- Progress - Number Input with Slider -->
            <td>
              <EditableCell
                v-model="task.progress"
                type="number"
                field="progress"
                :task-id="task.id"
                :original-data="task"
                :min="0"
                :max="100"
                :precision="0"
                :display-format="formatProgress"
                @change="handleTaskChange"
              />
            </td>

            <!-- Start Date - Date Picker -->
            <td>
              <EditableCell
                v-model="task.start"
                type="date"
                field="start_date"
                :task-id="task.id"
                :original-data="task"
                :display-format="formatDateShort"
                @change="handleTaskChange"
              />
            </td>

            <!-- Duration - Number Input -->
            <td>
              <EditableCell
                v-model="task.duration"
                type="number"
                field="duration"
                :task-id="task.id"
                :original-data="task"
                :min="0"
                :max="365"
                :step="1"
                :display-format="v => `${v} 天`"
                @change="handleTaskChange"
              />
            </td>

            <!-- Assignee - Select with Users -->
            <td>
              <EditableCell
                v-model="task.assignee_id"
                type="select"
                field="assignee_id"
                :task-id="task.id"
                :original-data="task"
                :options="userOptions"
                :display-format="v => getUserName(v)"
                clearable
                @change="handleTaskChange"
              />
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Task Templates Dialog -->
    <TaskTemplatesDialog
      v-model="templatesDialogVisible"
      :project-id="projectId"
      :start-date="defaultStartDate"
      @created="handleTaskCreated"
    />

    <!-- Bulk Edit Dialog -->
    <BulkEditDialog
      v-model="bulkEditDialogVisible"
      :tasks="selectedTaskObjects"
      :project-id="projectId"
      @updated="handleBulkUpdated"
    />

    <!-- Undo/Redo Controls -->
    <div class="undo-redo-controls">
      <el-button-group>
        <el-tooltip content="撤销 (Ctrl+Z)" placement="top">
          <el-button
            :disabled="!canUndo"
            @click="handleUndo"
          >
            <el-icon><RefreshLeft /></el-icon>
            撤销
          </el-button>
        </el-tooltip>

        <el-tooltip content="重做 (Ctrl+Y)" placement="top">
          <el-button
            :disabled="!canRedo"
            @click="handleRedo"
          >
            <el-icon><RefreshRight /></el-icon>
            重做
          </el-button>
        </el-tooltip>
      </el-button-group>

      <el-dropdown @command="handleHistoryAction">
        <el-button>
          历史记录
          <el-icon class="el-icon--right"><ArrowDown /></el-icon>
        </el-button>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="view">查看历史</el-dropdown-item>
            <el-dropdown-item command="clear" divided>清空历史</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Plus,
  Edit,
  RefreshLeft,
  RefreshRight,
  ArrowDown
} from '@element-plus/icons-vue'
import EditableCell from '@/components/gantt/table/EditableCell.vue'
import TaskTemplatesDialog from '@/components/gantt/dialogs/TaskTemplatesDialog.vue'
import BulkEditDialog from '@/components/gantt/dialogs/BulkEditDialog.vue'
import { useUndoRedoStore } from '@/stores/undoRedoStore'
import { ganttStore } from '@/stores/ganttStore'
import eventBus, { GanttEvents } from '@/utils/eventBus'

/**
 * Example Integration: Gantt Chart with Enhanced Editing
 *
 * This example demonstrates how to integrate all three editing components
 * with the Gantt chart table and undo/redo system.
 */

// Props
const props = defineProps({
  projectId: {
    type: [Number, String],
    required: true
  }
})

// Stores
const undoRedoStore = useUndoRedoStore()

// State
const tasks = ref([])
const selectedTasks = ref([])
const templatesDialogVisible = ref(false)
const bulkEditDialogVisible = ref(false)
const defaultStartDate = ref(new Date().toISOString().split('T')[0])

// Computed
const selectAll = computed({
  get: () => selectedTasks.value.length === tasks.value.length && tasks.value.length > 0,
  set: (value) => {
    if (value) {
      selectedTasks.value = tasks.value.map(t => t.id)
    } else {
      selectedTasks.value = []
    }
  }
})

const isSomeSelected = computed(() => {
  return selectedTasks.value.length > 0 && selectedTasks.value.length < tasks.value.length
})

const selectedTaskObjects = computed(() => {
  return tasks.value.filter(t => selectedTasks.value.includes(t.id))
})

const hasUnsavedChanges = computed(() => {
  return ganttStore.state.hasUnsavedChanges
})

const canUndo = computed(() => undoRedoStore.canUndo)
const canRedo = computed(() => undoRedoStore.canRedo)

// Options for EditableCell components
const statusOptions = [
  { value: 'not_started', label: '未开始' },
  { value: 'in_progress', label: '进行中' },
  { value: 'completed', label: '已完成' },
  { value: 'delayed', label: '已延期' }
]

const priorityOptions = [
  { value: 'high', label: '高' },
  { value: 'medium', label: '中' },
  { value: 'low', label: '低' }
]

const userOptions = ref([])

// Validation rules
const nameRules = [
  { required: true, message: '请输入任务名称', trigger: 'blur' },
  { min: 2, max: 100, message: '长度在 2 到 100 个字符', trigger: 'blur' }
]

// Methods
/**
 * Load tasks
 */
async function loadTasks() {
  try {
    await ganttStore.actions.setProject(props.projectId, 'Project Name')
    await ganttStore.actions.loadData()
    tasks.value = [...ganttStore.state.tasks]
  } catch (error) {
    ElMessage.error('加载任务失败')
    console.error(error)
  }
}

/**
 * Load users for assignee selection
 */
async function loadUsers() {
  try {
    const { progressApi } = await import('@/api')
    const response = await progressApi.getProjectUsers(props.projectId)
    userOptions.value = (response.data || []).map(user => ({
      value: user.id,
      label: user.full_name || user.username,
      ...user
    }))
  } catch (error) {
    console.error('Failed to load users:', error)
  }
}

/**
 * Toggle task selection
 */
function toggleTaskSelection(taskId) {
  const index = selectedTasks.value.indexOf(taskId)
  if (index > -1) {
    selectedTasks.value.splice(index, 1)
  } else {
    selectedTasks.value.push(taskId)
  }
}

/**
 * Handle select all checkbox
 */
function handleSelectAll(checked) {
  if (checked) {
    selectedTasks.value = tasks.value.map(t => t.id)
  } else {
    selectedTasks.value = []
  }
}

/**
 * Open templates dialog
 */
function openTemplatesDialog() {
  templatesDialogVisible.value = true
}

/**
 * Open bulk edit dialog
 */
function openBulkEditDialog() {
  if (selectedTasks.value.length === 0) {
    ElMessage.warning('请先选择要编辑的任务')
    return
  }
  bulkEditDialogVisible.value = true
}

/**
 * Handle task created from template
 */
function handleTaskCreated(task) {
  ElMessage.success('任务已创建')
  loadTasks() // Reload tasks
}

/**
 * Handle bulk update completed
 */
function handleBulkUpdated({ count, fields, changes }) {
  ElMessage.success(`已更新 ${count} 个任务`)
  loadTasks() // Reload tasks
  selectedTasks.value = [] // Clear selection
}

/**
 * Handle task cell change
 */
function handleTaskChange({ field, value, taskId }) {
  console.log(`Task ${taskId} field ${field} changed to ${value}`)

  // Reload tasks to reflect changes
  setTimeout(() => {
    loadTasks()
  }, 100)
}

/**
 * Handle undo
 */
async function handleUndo() {
  try {
    await undoRedoStore.undo()
    await loadTasks()
    ElMessage.info('已撤销')
  } catch (error) {
    console.error('Undo failed:', error)
  }
}

/**
 * Handle redo
 */
async function handleRedo() {
  try {
    await undoRedoStore.redo()
    await loadTasks()
    ElMessage.info('已重做')
  } catch (error) {
    console.error('Redo failed:', error)
  }
}

/**
 * Handle history actions
 */
function handleHistoryAction(command) {
  switch (command) {
    case 'view':
      ElMessageBox.alert(
        getHistorySnapshot(),
        '历史记录',
        {
          customClass: 'history-dialog',
          dangerouslyUseHTMLString: true
        }
      )
      break
    case 'clear':
      ElMessageBox.confirm('确定要清空所有历史记录吗？', '提示', {
        type: 'warning'
      }).then(() => {
        undoRedoStore.clearHistory()
        ElMessage.success('历史记录已清空')
      }).catch(() => {})
      break
  }
}

/**
 * Get history snapshot for display
 */
function getHistorySnapshot() {
  const undoStack = undoRedoStore.getHistorySnapshot()
  const redoStack = undoRedoStore.getRedoSnapshot()

  let html = '<div style="max-height: 400px; overflow-y: auto;">'

  if (undoStack.length > 0) {
    html += '<h4>撤销栈 (最近 → 最远):</h4><ul>'
    undoStack.slice().reverse().forEach((item, index) => {
      html += `<li>${index + 1}. ${item.description}</li>`
    })
    html += '</ul>'
  } else {
    html += '<p>暂无撤销历史</p>'
  }

  if (redoStack.length > 0) {
    html += '<h4>重做栈 (最近 → 最远):</h4><ul>'
    redoStack.forEach((item, index) => {
      html += `<li>${index + 1}. ${item.description}</li>`
    })
    html += '</ul>'
  } else {
    html += '<p>暂无重做历史</p>'
  }

  html += '</div>'
  return html
}

// Display formatters
function formatStatus(value) {
  const map = {
    not_started: '未开始',
    in_progress: '进行中',
    completed: '已完成',
    delayed: '已延期'
  }
  return map[value] || value
}

function formatPriority(value) {
  const map = {
    high: '高',
    medium: '中',
    low: '低'
  }
  return map[value] || value
}

function formatProgress(value) {
  return `${value || 0}%`
}

function formatDateShort(value) {
  if (!value) return '-'
  const date = new Date(value)
  return `${date.getMonth() + 1}/${date.getDate()}`
}

function getUserName(userId) {
  if (!userId) return '未分配'
  const user = userOptions.value.find(u => u.value === userId)
  return user ? user.label : '未知'
}

// Keyboard shortcuts
function handleKeydown(event) {
  // Ctrl+Z or Cmd+Z - Undo
  if ((event.ctrlKey || event.metaKey) && event.key === 'z' && !event.shiftKey) {
    event.preventDefault()
    if (canUndo.value) {
      handleUndo()
    }
  }

  // Ctrl+Y or Cmd+Shift+Z or Ctrl+Shift+Z - Redo
  if (
    ((event.ctrlKey || event.metaKey) && event.key === 'y') ||
    ((event.ctrlKey || event.metaKey) && event.shiftKey && event.key === 'z')
  ) {
    event.preventDefault()
    if (canRedo.value) {
      handleRedo()
    }
  }
}

// Lifecycle
onMounted(() => {
  loadTasks()
  loadUsers()
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
})
</script>

<style scoped>
.gantt-view-example {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 16px;
}

.gantt-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.toolbar-info {
  display: flex;
  gap: 12px;
}

.gantt-table-container {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  overflow-x: auto;
}

.gantt-table {
  width: 100%;
  border-collapse: collapse;
}

.gantt-table th,
.gantt-table td {
  padding: 0;
  text-align: left;
  border-bottom: 1px solid #ebeef5;
}

.gantt-table th {
  background: #f5f7fa;
  padding: 12px 16px;
  font-weight: 500;
  color: #303133;
  white-space: nowrap;
}

.gantt-table td {
  padding: 8px;
}

.gantt-table tbody tr:hover {
  background-color: #f5f7fa;
}

.gantt-table tbody tr.is-selected {
  background-color: #ecf5ff;
}

.undo-redo-controls {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 16px;
  padding: 12px 16px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

/* Status cell styles */
:deep(.status-not_started) {
  color: #909399;
}

:deep(.status-in_progress) {
  color: #409eff;
}

:deep(.status-completed) {
  color: #67c23a;
}

:deep(.status-delayed) {
  color: #f56c6c;
}

/* History dialog */
:deep(.history-dialog) {
  width: 600px;
}

:deep(.history-dialog .el-message-box__content) {
  padding: 20px;
}

:deep(.history-dialog h4) {
  margin: 0 0 12px;
  font-size: 14px;
  font-weight: 500;
  color: #303133;
}

:deep(.history-dialog ul) {
  margin: 0;
  padding-left: 20px;
}

:deep(.history-dialog li) {
  margin: 4px 0;
  font-size: 13px;
  color: #606266;
}

:deep(.history-dialog p) {
  margin: 0;
  font-size: 13px;
  color: #909399;
}

/* Responsive */
@media (max-width: 1024px) {
  .gantt-table {
    font-size: 12px;
  }

  .gantt-table th,
  .gantt-table td {
    padding: 6px 4px;
  }
}

@media (max-width: 768px) {
  .gantt-toolbar {
    flex-direction: column;
    gap: 12px;
  }

  .undo-redo-controls {
    flex-direction: column;
    gap: 8px;
  }
}
</style>
