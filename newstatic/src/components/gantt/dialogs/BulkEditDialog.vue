<template>
  <el-dialog
    v-model="dialogVisible"
    title="批量编辑"
    width="700px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <!-- Selection Summary -->
    <div class="selection-summary">
      <el-alert
        :title="`已选择 ${selectedTasks.length} 个任务`"
        type="info"
        :closable="false"
        show-icon
      >
        <template #default>
          <div class="task-list">
            <div
              v-for="task in selectedTasks.slice(0, 5)"
              :key="task.id"
              class="task-item"
            >
              <el-icon class="task-icon"><Document /></el-icon>
              <span class="task-name">{{ task.name }}</span>
            </div>
            <div v-if="selectedTasks.length > 5" class="task-more">
              还有 {{ selectedTasks.length - 5 }} 个任务...
            </div>
          </div>
        </template>
      </el-alert>
    </div>

    <!-- Edit Form -->
    <el-form
      ref="formRef"
      :model="formData"
      label-width="120px"
      class="bulk-edit-form"
    >
      <!-- Fields to Edit -->
      <el-form-item label="编辑字段">
        <el-checkbox-group v-model="selectedFields" @change="handleFieldsChange">
          <el-checkbox label="status">状态</el-checkbox>
          <el-checkbox label="priority">优先级</el-checkbox>
          <el-checkbox label="progress">进度</el-checkbox>
          <el-checkbox label="assignee">负责人</el-checkbox>
          <el-checkbox label="startDate">开始日期</el-checkbox>
          <el-checkbox label="duration">工期</el-checkbox>
          <el-checkbox label="endDate">结束日期</el-checkbox>
        </el-checkbox-group>
      </el-form-item>

      <el-divider />

      <!-- Status -->
      <el-form-item v-if="selectedFields.includes('status')" label="状态">
        <el-select v-model="formData.status" placeholder="选择状态" style="width: 100%">
          <el-option label="未开始" value="not_started" />
          <el-option label="进行中" value="in_progress" />
          <el-option label="已完成" value="completed" />
          <el-option label="已延期" value="delayed" />
        </el-select>
      </el-form-item>

      <!-- Priority -->
      <el-form-item v-if="selectedFields.includes('priority')" label="优先级">
        <el-select v-model="formData.priority" placeholder="选择优先级" style="width: 100%">
          <el-option label="高" value="high">
            <span class="priority-option priority-high">高</span>
          </el-option>
          <el-option label="中" value="medium">
            <span class="priority-option priority-medium">中</span>
          </el-option>
          <el-option label="低" value="low">
            <span class="priority-option priority-low">低</span>
          </el-option>
        </el-select>
      </el-form-item>

      <!-- Progress -->
      <el-form-item v-if="selectedFields.includes('progress')" label="进度">
        <div class="progress-input">
          <el-slider
            v-model="formData.progress"
            :marks="progressMarks"
            :step="10"
            show-input
          />
        </div>
      </el-form-item>

      <!-- Assignee -->
      <el-form-item v-if="selectedFields.includes('assignee')" label="负责人">
        <el-select
          v-model="formData.assignee"
          placeholder="选择负责人"
          filterable
          clearable
          style="width: 100%"
        >
          <el-option
            v-for="user in userList"
            :key="user.id"
            :label="user.full_name || user.username"
            :value="user.id"
          >
            <div class="user-option">
              <el-avatar :size="24" :src="user.avatar">
                {{ (user.full_name || user.username || '?').charAt(0) }}
              </el-avatar>
              <span>{{ user.full_name || user.username }}</span>
            </div>
          </el-option>
        </el-select>
      </el-form-item>

      <!-- Start Date -->
      <el-form-item v-if="selectedFields.includes('startDate')" label="开始日期">
        <el-date-picker
          v-model="formData.startDate"
          type="date"
          placeholder="选择开始日期"
          value-format="YYYY-MM-DD"
          style="width: 100%"
        />
        <el-checkbox v-model="adjustRelative" style="margin-top: 8px">
          相对调整（保持任务间隔）
        </el-checkbox>
      </el-form-item>

      <!-- Duration -->
      <el-form-item v-if="selectedFields.includes('duration')" label="工期">
        <el-input-number
          v-model="formData.duration"
          :min="0"
          :max="365"
          :step="1"
          controls-position="right"
        />
        <span class="unit-label">天</span>
      </el-form-item>

      <!-- End Date -->
      <el-form-item v-if="selectedFields.includes('endDate')" label="结束日期">
        <el-date-picker
          v-model="formData.endDate"
          type="date"
          placeholder="选择结束日期"
          value-format="YYYY-MM-DD"
          style="width: 100%"
        />
        <div class="form-tip">
          <el-icon><InfoFilled /></el-icon>
          设置结束日期将自动计算工期
        </div>
      </el-form-item>
    </el-form>

    <!-- Change Preview -->
    <div v-if="hasChanges" class="change-preview">
      <div class="preview-header">
        <h4>变更预览</h4>
        <el-button text size="small" @click="showAllChanges = !showAllChanges">
          {{ showAllChanges ? '收起' : '展开全部' }}
        </el-button>
      </div>

      <el-table
        :data="previewData"
        stripe
        border
        size="small"
        max-height="300"
        class="preview-table"
      >
        <el-table-column prop="name" label="任务名称" width="180" />
        <el-table-column label="变更内容" min-width="400">
          <template #default="{ row }">
            <div class="change-list">
              <div
                v-for="change in row.changes"
                :key="change.field"
                class="change-item"
              >
                <span class="change-field">{{ change.label }}:</span>
                <span class="change-from">{{ change.from }}</span>
                <el-icon class="change-arrow"><Right /></el-icon>
                <span class="change-to">{{ change.to }}</span>
              </div>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button
          type="primary"
          :loading="submitting"
          :disabled="!hasChanges || selectedFields.length === 0"
          @click="handleApply"
        >
          应用更改 ({{ selectedTasks.length }} 个任务)
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Document,
  InfoFilled,
  Right
} from '@element-plus/icons-vue'
import { progressApi } from '@/api'
import { useUndoRedoStore } from '@/stores/undoRedoStore'
import { BatchUpdateCommand } from '@/stores/undoRedoStore'
import eventBus, { GanttEvents } from '@/utils/eventBus'
import { addDays, diffDays } from '@/utils/dateFormat'

/**
 * BulkEditDialog Component
 *
 * Batch edit multiple tasks at once
 * Integrates with undo/redo system for batch operations
 *
 * @date 2025-02-18
 */

const props = defineProps({
  modelValue: Boolean,
  tasks: {
    type: Array,
    default: () => []
  },
  projectId: {
    type: [Number, String],
    required: true
  }
})

const emit = defineEmits(['update:modelValue', 'updated'])

// Store
const undoRedoStore = useUndoRedoStore()

// State
const dialogVisible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const selectedTasks = ref([])
const selectedFields = ref([])
const formData = ref({
  status: '',
  priority: '',
  progress: 0,
  assignee: null,
  startDate: '',
  duration: null,
  endDate: ''
})

const adjustRelative = ref(false)
const showAllChanges = ref(false)
const submitting = ref(false)
const userList = ref([])

const progressMarks = {
  0: '0%',
  50: '50%',
  100: '100%'
}

// Computed
const hasChanges = computed(() => {
  return selectedFields.value.length > 0
})

const previewData = computed(() => {
  if (!hasChanges.value) return []

  return selectedTasks.value.slice(0, showAllChanges.value ? undefined : 5).map(task => {
    const changes = []

    if (selectedFields.value.includes('status') && formData.value.status) {
      changes.push({
        field: 'status',
        label: '状态',
        from: statusLabel(task.status),
        to: statusLabel(formData.value.status)
      })
    }

    if (selectedFields.value.includes('priority') && formData.value.priority) {
      changes.push({
        field: 'priority',
        label: '优先级',
        from: priorityLabel(task.priority),
        to: priorityLabel(formData.value.priority)
      })
    }

    if (selectedFields.value.includes('progress')) {
      changes.push({
        field: 'progress',
        label: '进度',
        from: `${task.progress || 0}%`,
        to: `${formData.value.progress}%`
      })
    }

    if (selectedFields.value.includes('assignee') && formData.value.assignee) {
      const fromUser = userList.value.find(u => u.id === task.assignee_id)
      const toUser = userList.value.find(u => u.id === formData.value.assignee)
      changes.push({
        field: 'assignee',
        label: '负责人',
        from: fromUser ? (fromUser.full_name || fromUser.username) : '未分配',
        to: toUser ? (toUser.full_name || toUser.username) : '未分配'
      })
    }

    if (selectedFields.value.includes('startDate') && formData.value.startDate) {
      changes.push({
        field: 'startDate',
        label: '开始日期',
        from: task.start || '-',
        to: formData.value.startDate
      })
    }

    if (selectedFields.value.includes('duration') && formData.value.duration !== null) {
      changes.push({
        field: 'duration',
        label: '工期',
        from: `${task.duration || 0} 天`,
        to: `${formData.value.duration} 天`
      })
    }

    if (selectedFields.value.includes('endDate') && formData.value.endDate) {
      changes.push({
        field: 'endDate',
        label: '结束日期',
        from: task.end || '-',
        to: formData.value.endDate
      })
    }

    return {
      name: task.name,
      changes
    }
  })
})

// Methods
/**
 * Load users list
 */
async function loadUsers() {
  try {
    const response = await progressApi.getProjectUsers(props.projectId)
    userList.value = response.data || []
  } catch (error) {
    console.error('Failed to load users:', error)
    userList.value = []
  }
}

/**
 * Handle fields selection change
 */
function handleFieldsChange(fields) {
  // Reset unselected field values
  if (!fields.includes('status')) {
    formData.value.status = ''
  }
  if (!fields.includes('priority')) {
    formData.value.priority = ''
  }
  if (!fields.includes('progress')) {
    formData.value.progress = 0
  }
  if (!fields.includes('assignee')) {
    formData.value.assignee = null
  }
  if (!fields.includes('startDate')) {
    formData.value.startDate = ''
  }
  if (!fields.includes('duration')) {
    formData.value.duration = null
  }
  if (!fields.includes('endDate')) {
    formData.value.endDate = ''
  }
}

/**
 * Apply bulk changes
 */
async function handleApply() {
  if (!hasChanges.value) {
    ElMessage.warning('请选择要编辑的字段')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要批量更新 ${selectedTasks.value.length} 个任务吗？此操作可以撤销。`,
      '确认批量编辑',
      {
        type: 'warning',
        confirmButtonText: '确定',
        cancelButtonText: '取消'
      }
    )

    submitting.value = true

    // Prepare updates for each task
    const updates = selectedTasks.value.map(task => {
      const updateData = {
        taskId: task.id,
        updates: {},
        originalData: {
          id: task.id,
          task_name: task.name,
          status: task.status,
          priority: task.priority,
          progress: task.progress || 0,
          assignee_id: task.assignee_id,
          start_date: task.start,
          duration: task.duration,
          end_date: task.end
        }
      }

      // Apply selected field changes
      if (selectedFields.value.includes('status') && formData.value.status) {
        updateData.updates.status = formData.value.status
      }

      if (selectedFields.value.includes('priority') && formData.value.priority) {
        updateData.updates.priority = formData.value.priority
      }

      if (selectedFields.value.includes('progress')) {
        updateData.updates.progress = formData.value.progress
      }

      if (selectedFields.value.includes('assignee') && formData.value.assignee) {
        updateData.updates.assignee_id = formData.value.assignee
      }

      if (selectedFields.value.includes('startDate') && formData.value.startDate) {
        if (adjustRelative.value) {
          // Calculate offset from original start date
          const originalStart = new Date(task.start)
          const newStart = new Date(formData.value.startDate)
          const offset = diffDays(task.start, formData.value.startDate)

          updateData.updates.start_date = formData.value.startDate

          // Adjust end date by same offset
          if (task.end) {
            const newEnd = addDays(task.end, offset)
            updateData.updates.end_date = newEnd
          }
        } else {
          updateData.updates.start_date = formData.value.startDate
        }
      }

      if (selectedFields.value.includes('duration') && formData.value.duration !== null) {
        updateData.updates.duration = formData.value.duration

        // Recalculate end date if needed
        if (task.start) {
          const endDate = addDays(task.start, formData.value.duration)
          updateData.updates.end_date = endDate
        }
      }

      if (selectedFields.value.includes('endDate') && formData.value.endDate) {
        updateData.updates.end_date = formData.value.endDate

        // Recalculate duration
        if (task.start) {
          const duration = diffDays(task.start, formData.value.endDate)
          updateData.updates.duration = Math.max(duration, 0)
        }
      }

      return updateData
    })

    // Create batch update command
    const command = new BatchUpdateCommand(updates, null)

    // Execute command (integrates with undo/redo)
    await undoRedoStore.executeCommand(command)

    ElMessage.success(`已成功更新 ${updates.length} 个任务`)

    emit('updated', {
      count: updates.length,
      fields: selectedFields.value,
      changes: formData.value
    })

    // Emit gantt event
    eventBus.emit(GanttEvents.TASKS_UPDATED, {
      taskIds: updates.map(u => u.taskId),
      fields: selectedFields.value
    })

    handleClose()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Bulk update failed:', error)
      ElMessage.error(error.message || '批量更新失败')
    }
  } finally {
    submitting.value = false
  }
}

/**
 * Close dialog
 */
function handleClose() {
  selectedFields.value = []
  formData.value = {
    status: '',
    priority: '',
    progress: 0,
    assignee: null,
    startDate: '',
    duration: null,
    endDate: ''
  }
  adjustRelative.value = false
  showAllChanges.value = false
  emit('update:modelValue', false)
}

/**
 * Helper functions
 */
function statusLabel(status) {
  const map = {
    not_started: '未开始',
    in_progress: '进行中',
    completed: '已完成',
    delayed: '已延期'
  }
  return map[status] || status || '-'
}

function priorityLabel(priority) {
  const map = {
    high: '高',
    medium: '中',
    low: '低'
  }
  return map[priority] || priority || '-'
}

// Watch for tasks prop changes
watch(() => props.tasks, (tasks) => {
  selectedTasks.value = [...tasks]
}, { immediate: true })

// Watch for dialog open
watch(() => props.modelValue, (val) => {
  if (val) {
    loadUsers()
    selectedTasks.value = [...props.tasks]
  }
})
</script>

<style scoped>
.selection-summary {
  margin-bottom: 24px;
}

.task-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 12px;
}

.task-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 0;
}

.task-icon {
  color: var(--el-color-primary, #409eff);
  font-size: 16px;
  flex-shrink: 0;
}

.task-name {
  font-size: 13px;
  color: var(--el-text-color-regular, #606266);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.task-more {
  font-size: 12px;
  color: var(--el-text-color-secondary, #909399);
  padding: 4px 0;
}

.bulk-edit-form {
  margin-bottom: 24px;
}

.progress-input {
  width: 100%;
  padding-right: 80px;
}

.unit-label {
  margin-left: 8px;
  font-size: 12px;
  color: var(--el-text-color-secondary, #909399);
}

.form-tip {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-top: 8px;
  font-size: 12px;
  color: var(--el-text-color-secondary, #909399);
}

.priority-option {
  display: inline-flex;
  align-items: center;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.priority-option.priority-high {
  background: #fef0f0;
  color: #f56c6c;
}

.priority-option.priority-medium {
  background: #fdf6ec;
  color: #e6a23c;
}

.priority-option.priority-low {
  background: #f0f9ff;
  color: #67c23a;
}

.user-option {
  display: flex;
  align-items: center;
  gap: 8px;
}

.change-preview {
  border: 1px solid var(--el-border-color, #dcdfe6);
  border-radius: 8px;
  padding: 16px;
  background: var(--el-fill-color-light, #f5f7fa);
}

.preview-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.preview-header h4 {
  margin: 0;
  font-size: 14px;
  font-weight: 500;
  color: var(--el-text-color-primary, #303133);
}

.preview-table {
  font-size: 12px;
}

.preview-table :deep(.el-table__cell) {
  padding: 6px 0;
}

.change-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.change-item {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.change-field {
  font-weight: 500;
  color: var(--el-text-color-primary, #303133);
  min-width: 50px;
}

.change-from {
  color: var(--el-text-color-secondary, #909399);
  text-decoration: line-through;
}

.change-arrow {
  color: var(--el-color-primary, #409eff);
  font-size: 14px;
}

.change-to {
  color: var(--el-color-success, #67c23a);
  font-weight: 500;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  width: 100%;
}

/* Responsive */
@media (max-width: 768px) {
  .change-item {
    font-size: 11px;
  }

  .preview-table {
    max-height: 200px;
  }
}
</style>
