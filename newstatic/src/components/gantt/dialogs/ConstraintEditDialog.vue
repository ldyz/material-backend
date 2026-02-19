<template>
  <el-dialog
    v-model="dialogVisible"
    :title="t('gantt.constraints.editTitle')"
    width="600px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <div v-loading="loading" class="constraint-dialog">
      <!-- Constraint Type Selection -->
      <div class="form-section">
        <label class="form-label">{{ t('gantt.constraints.type') }}</label>
        <el-select
          v-model="constraint.type"
          :placeholder="t('gantt.constraints.selectType')"
          @change="handleTypeChange"
        >
          <el-option
            v-for="type in constraintTypes"
            :key="type.value"
            :label="type.label"
            :value="type.value"
          >
            <div class="constraint-option">
              <span class="constraint-code">{{ type.value }}</span>
              <span class="constraint-label">{{ type.label }}</span>
            </div>
          </el-option>
        </el-select>

        <!-- Constraint Description -->
        <div v-if="constraint.type" class="constraint-description">
          <el-alert type="info" :closable="false">
            <template #title>
              {{ getConstraintDescription(constraint.type) }}
            </template>
          </el-alert>
        </div>
      </div>

      <!-- Constraint Date -->
      <div class="form-section">
        <label class="form-label">{{ t('gantt.constraints.date') }}</label>
        <el-date-picker
          v-model="constraint.date"
          type="date"
          :placeholder="t('gantt.constraints.selectDate')"
          :disabled="!constraint.type"
          format="YYYY-MM-DD"
          value-format="YYYY-MM-DD"
          @change="handleDateChange"
        />
      </div>

      <!-- Visual Preview -->
      <div v-if="showPreview" class="preview-section">
        <div class="preview-header">
          <span>{{ t('gantt.constraints.preview') }}</span>
        </div>
        <div class="preview-timeline">
          <!-- Timeline visualization -->
          <div class="timeline-bar">
            <div class="task-range" :style="taskRangeStyle">
              <span class="task-label">{{ task?.name }}</span>
            </div>
            <div v-if="constraintMarkerStyle" class="constraint-marker" :style="constraintMarkerStyle">
              <span class="marker-label">{{ constraint.type }}</span>
            </div>
          </div>
          <!-- Legend -->
          <div class="timeline-legend">
            <div class="legend-item">
              <span class="legend-color task-color"></span>
              <span>{{ t('gantt.constraints.currentSchedule') }}</span>
            </div>
            <div class="legend-item">
              <span class="legend-color constraint-color"></span>
              <span>{{ t('gantt.constraints.constraintDate') }}</span>
            </div>
          </div>
        </div>
        <!-- Validation Message -->
        <div v-if="validationMessage" class="validation-message">
          <el-alert :type="validationType" :closable="false">
            <template #title>
              {{ validationMessage }}
            </template>
          </el-alert>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">{{ t('common.cancel') }}</el-button>
        <el-button
          type="primary"
          :loading="saving"
          :disabled="!canSave"
          @click="handleSave"
        >
          {{ t('common.save') }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { validateConstraint, calculateConstraintImpact } from '@/utils/ganttConstraints'
import { useGanttStore } from '@/stores/ganttStore'
import { useUndoRedoStore } from '@/stores/undoRedoStore'

const { t } = useI18n()
const ganttStore = useGanttStore()
const undoRedoStore = useUndoRedoStore()

// Props
const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  task: {
    type: Object,
    default: null
  }
})

// Emits
const emit = defineEmits(['update:modelValue', 'constraintUpdated'])

// State
const dialogVisible = ref(false)
const loading = ref(false)
const saving = ref(false)
const showPreview = ref(false)

const constraint = ref({
  type: '',
  date: null
})

const originalConstraint = ref(null)

// Constraint types with labels
const constraintTypes = [
  { value: 'MSO', label: 'Must Start On - 必须开始于' },
  { value: 'MFO', label: 'Must Finish On - 必须完成于' },
  { value: 'SNET', label: 'Start No Earlier Than - 不早于开始' },
  { value: 'SNLT', label: 'Start No Later Than - 不晚于开始' },
  { value: 'FNET', label: 'Finish No Earlier Than - 不早于完成' },
  { value: 'FNLT', label: 'Finish No Later Than - 不晚于完成' }
]

// Computed
const canSave = computed(() => {
  return constraint.value.type && constraint.value.date
})

const validationMessage = computed(() => {
  if (!props.task || !constraint.value.type || !constraint.value.date) {
    return ''
  }

  const validation = validateConstraint(
    props.task,
    constraint.value.type,
    constraint.value.date
  )

  if (!validation.valid) {
    return validation.message
  }

  // Check if constraint would shift task
  const impact = calculateConstraintImpact(
    props.task,
    constraint.value.type,
    constraint.value.date
  )

  if (impact.wouldShift) {
    return t('gantt.constraints.wouldShift', {
      days: Math.abs(impact.shiftDays),
      direction: impact.shiftDays > 0 ? t('gantt.constraints.later') : t('gantt.constraints.earlier')
    })
  }

  return t('gantt.constraints.noEffect')
})

const validationType = computed(() => {
  const validation = validateConstraint(
    props.task,
    constraint.value.type,
    constraint.value.date
  )
  return validation.valid ? 'success' : 'warning'
})

const taskRangeStyle = computed(() => {
  if (!props.task) return {}

  const startDate = new Date(props.task.start_date)
  const endDate = new Date(props.task.end_date)
  const dayWidth = 20
  const days = Math.ceil((endDate - startDate) / (1000 * 60 * 60 * 24))

  return {
    left: '0px',
    width: `${days * dayWidth}px`,
    backgroundColor: '#409EFF'
  }
})

const constraintMarkerStyle = computed(() => {
  if (!props.task || !constraint.value.date) return null

  const startDate = new Date(props.task.start_date)
  const constraintDate = new Date(constraint.value.date)
  const dayWidth = 20
  const offset = Math.ceil((constraintDate - startDate) / (1000 * 60 * 60 * 24))

  return {
    left: `${offset * dayWidth}px`,
    backgroundColor: '#E6A23C'
  }
})

// Methods
const getConstraintDescription = (type) => {
  const descriptions = {
    MSO: t('gantt.constraints.descMSO'),
    MFO: t('gantt.constraints.descMFO'),
    SNET: t('gantt.constraints.descSNET'),
    SNLT: t('gantt.constraints.descSNLT'),
    FNET: t('gantt.constraints.descFNET'),
    FNLT: t('gantt.constraints.descFNLT')
  }
  return descriptions[type] || ''
}

const handleTypeChange = () => {
  showPreview.value = true
  validateConstraint()
}

const handleDateChange = () => {
  showPreview.value = true
  validateConstraint()
}

const validateConstraint = () => {
  // Validation is computed
}

const handleSave = async () => {
  if (!canSave.value) return

  saving.value = true

  try {
    // Create snapshot for undo/redo
    const snapshot = {
      type: 'constraint',
      taskId: props.task.id,
      before: { ...originalConstraint.value },
      after: { ...constraint.value }
    }

    // Update task with constraint
    await ganttStore.updateTaskConstraint(props.task.id, constraint.value)

    // Add to undo/redo stack
    undoRedoStore.addOperation({
      type: 'updateConstraint',
      taskId: props.task.id,
      before: originalConstraint.value,
      after: { ...constraint.value },
      timestamp: new Date(),
      description: t('gantt.constraints.undoDescription', {
        task: props.task.name,
        constraint: constraint.value.type
      })
    })

    ElMessage.success(t('gantt.constraints.saveSuccess'))
    emit('constraintUpdated', { task: props.task, constraint: constraint.value })
    handleClose()
  } catch (error) {
    console.error('Failed to save constraint:', error)
    ElMessage.error(t('gantt.constraints.saveError'))
  } finally {
    saving.value = false
  }
}

const handleClose = () => {
  dialogVisible.value = false
  emit('update:modelValue', false)
  // Reset form
  nextTick(() => {
    constraint.value = {
      type: '',
      date: null
    }
    showPreview.value = false
    originalConstraint.value = null
  })
}

// Watch modelValue changes
watch(() => props.modelValue, (newVal) => {
  dialogVisible.value = newVal
  if (newVal && props.task) {
    // Load existing constraint if any
    if (props.task.constraint) {
      constraint.value = {
        type: props.task.constraint.type || '',
        date: props.task.constraint.date || null
      }
      originalConstraint.value = { ...constraint.value }
    } else {
      constraint.value = {
        type: '',
        date: null
      }
      originalConstraint.value = null
    }
    showPreview.value = false
  }
})

watch(dialogVisible, (newVal) => {
  emit('update:modelValue', newVal)
})
</script>

<style scoped>
.constraint-dialog {
  padding: 10px 0;
}

.form-section {
  margin-bottom: 24px;
}

.form-label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #606266;
}

.constraint-option {
  display: flex;
  align-items: center;
  gap: 12px;
}

.constraint-code {
  font-family: 'Courier New', monospace;
  font-weight: bold;
  color: #409EFF;
  background: #ecf5ff;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.constraint-label {
  flex: 1;
}

.constraint-description {
  margin-top: 12px;
}

.preview-section {
  margin-top: 24px;
  padding: 16px;
  background: #f5f7fa;
  border-radius: 4px;
}

.preview-header {
  font-weight: 500;
  margin-bottom: 12px;
  color: #606266;
}

.preview-timeline {
  margin: 16px 0;
}

.timeline-bar {
  position: relative;
  height: 40px;
  background: #e4e7ed;
  border-radius: 4px;
  margin-bottom: 12px;
}

.task-range {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  height: 24px;
  background: #409EFF;
  border-radius: 3px;
  display: flex;
  align-items: center;
  padding: 0 8px;
  color: white;
  font-size: 12px;
  white-space: nowrap;
  overflow: hidden;
}

.constraint-marker {
  position: absolute;
  top: 0;
  bottom: 0;
  width: 2px;
  background: #E6A23C;
}

.constraint-marker::before {
  content: '';
  position: absolute;
  top: -4px;
  left: -4px;
  width: 0;
  height: 0;
  border-left: 5px solid transparent;
  border-right: 5px solid transparent;
  border-top: 5px solid #E6A23C;
}

.marker-label {
  position: absolute;
  top: -20px;
  left: 50%;
  transform: translateX(-50%);
  font-size: 10px;
  color: #E6A23C;
  white-space: nowrap;
  font-weight: bold;
}

.timeline-legend {
  display: flex;
  gap: 16px;
  font-size: 12px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.legend-color {
  width: 12px;
  height: 12px;
  border-radius: 2px;
}

.legend-color.task-color {
  background: #409EFF;
}

.legend-color.constraint-color {
  background: #E6A23C;
}

.validation-message {
  margin-top: 12px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
