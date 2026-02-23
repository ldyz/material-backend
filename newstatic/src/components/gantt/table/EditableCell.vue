<template>
  <div
    class="editable-cell"
    :class="{ 'is-editing': isEditing, 'is-readonly': readonly }"
    @dblclick="handleDoubleClick"
  >
    <!-- Display mode -->
    <div v-if="!isEditing" class="editable-cell__display">
      <slot name="display">
        <span :class="displayClass">{{ displayValue }}</span>
      </slot>
      <el-icon v-if="!readonly" class="editable-icon">
        <Edit />
      </el-icon>
    </div>

    <!-- Edit mode -->
    <div v-else class="editable-cell__editor" @click.stop>
      <!-- Text input -->
      <el-input
        v-if="type === 'text'"
        ref="inputRef"
        v-model="editValue"
        :placeholder="placeholder"
        :maxlength="maxlength"
        size="small"
        @blur="handleBlur"
        @keyup.enter="handleEnter"
        @keyup.esc="handleEscape"
      />

      <!-- Number input -->
      <el-input-number
        v-else-if="type === 'number'"
        ref="inputRef"
        v-model="editValue"
        :min="min"
        :max="max"
        :precision="precision"
        :step="step"
        size="small"
        controls-position="right"
        @blur="handleBlur"
        @keyup.enter="handleEnter"
        @keyup.esc="handleEscape"
      />

      <!-- Date picker -->
      <el-date-picker
        v-else-if="type === 'date'"
        ref="inputRef"
        v-model="editValue"
        type="date"
        :placeholder="placeholder"
        :format="dateFormat"
        :value-format="valueFormat"
        size="small"
        :clearable="false"
        @blur="handleBlur"
        @visible-change="handleDateVisibleChange"
      />

      <!-- Select dropdown -->
      <el-select
        v-else-if="type === 'select'"
        ref="inputRef"
        v-model="editValue"
        :placeholder="placeholder"
        size="small"
        filterable
        :clearable="clearable"
        @blur="handleBlur"
        @keyup.enter="handleEnter"
        @keyup.esc="handleEscape"
      >
        <el-option
          v-for="option in options"
          :key="option.value"
          :label="option.label"
          :value="option.value"
          :disabled="option.disabled"
        >
          <slot v-if="option.slot" :name="option.slot" :option="option" />
        </el-option>
      </el-select>

      <!-- Validation message -->
      <div v-if="validationError" class="editable-cell__error">
        <el-icon><WarningFilled /></el-icon>
        <span>{{ validationError }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick, onMounted, onUnmounted } from 'vue'
import { Edit, WarningFilled } from '@element-plus/icons-vue'
import { useUndoRedoStore } from '@/stores/undoRedoStore'
import { UpdateTaskCommand } from '@/stores/undoRedoStore'
import eventBus, { GanttEvents } from '@/utils/eventBus'

/**
 * EditableCell Component
 *
 * Inline editable cell component for Gantt chart table
 * Supports text, number, date, and select input types
 * Integrates with undo/redo system
 *
 * @date 2025-02-18
 */

const props = defineProps({
  // Current value
  modelValue: {
    type: [String, Number, Date, Object],
    default: ''
  },
  // Cell type
  type: {
    type: String,
    default: 'text',
    validator: (value) => ['text', 'number', 'date', 'select'].includes(value)
  },
  // Field name (for API updates)
  field: {
    type: String,
    required: true
  },
  // Task ID (for API updates)
  taskId: {
    type: [Number, String],
    required: true
  },
  // Original task data (for undo)
  originalData: {
    type: Object,
    default: () => ({})
  },
  // Display format
  displayFormat: {
    type: Function,
    default: null
  },
  // Display class
  displayClass: {
    type: String,
    default: ''
  },
  // Placeholder text
  placeholder: {
    type: String,
    default: '请输入'
  },
  // Readonly mode
  readonly: {
    type: Boolean,
    default: false
  },
  // Number input props
  min: {
    type: Number,
    default: -Infinity
  },
  max: {
    type: Number,
    default: Infinity
  },
  precision: {
    type: Number,
    default: 0
  },
  step: {
    type: Number,
    default: 1
  },
  // Text input props
  maxlength: {
    type: Number,
    default: 255
  },
  // Date picker props
  dateFormat: {
    type: String,
    default: 'YYYY-MM-DD'
  },
  valueFormat: {
    type: String,
    default: 'YYYY-MM-DD'
  },
  // Select options
  options: {
    type: Array,
    default: () => []
  },
  // Select clearable
  clearable: {
    type: Boolean,
    default: false
  },
  // Validation rules
  rules: {
    type: Array,
    default: () => []
  },
  // Auto-save on blur
  autoSave: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['update:modelValue', 'change', 'edit', 'cancel'])

// Store
const undoRedoStore = useUndoRedoStore()

// State
const isEditing = ref(false)
const editValue = ref(props.modelValue)
const inputRef = ref(null)
const validationError = ref('')
const isSaving = ref(false)

// 定时器引用
let selectTimeout = null

// Computed
const displayValue = computed(() => {
  if (props.displayFormat) {
    return props.displayFormat(props.modelValue)
  }

  // Default display formats by type
  switch (props.type) {
    case 'date':
      return props.modelValue || '-'
    case 'number':
      return props.modelValue !== undefined && props.modelValue !== null
        ? props.modelValue
        : '-'
    case 'select':
      const option = props.options.find(opt => opt.value === props.modelValue)
      return option ? option.label : props.modelValue || '-'
    default:
      return props.modelValue || '-'
  }
})

// Watch for external value changes
watch(() => props.modelValue, (newValue) => {
  if (!isEditing.value) {
    editValue.value = newValue
  }
})

// Methods
/**
 * Start editing
 */
function startEditing() {
  if (props.readonly) return

  isEditing.value = true
  editValue.value = props.modelValue
  validationError.value = ''

  emit('edit', props.field)

  nextTick(() => {
    if (inputRef.value) {
      inputRef.value.focus()
      // For select, we need to focus and then show the dropdown
      if (props.type === 'select' && inputRef.value.focus) {
        // 清理之前的定时器
        if (selectTimeout) {
          clearTimeout(selectTimeout)
        }
        selectTimeout = setTimeout(() => {
          if (inputRef.value?.blur) {
            // The select dropdown will show automatically on focus
          }
          selectTimeout = null
        }, 100)
      }
    }
  })
}

/**
 * Handle double click to edit
 */
function handleDoubleClick() {
  if (!isEditing.value && !props.readonly) {
    startEditing()
  }
}

/**
 * Validate current value
 */
function validate() {
  if (!props.rules || props.rules.length === 0) {
    return true
  }

  for (const rule of props.rules) {
    if (rule.required && (editValue.value === '' || editValue.value === null || editValue.value === undefined)) {
      validationError.value = rule.message || '该字段为必填项'
      return false
    }

    if (rule.validator && typeof rule.validator === 'function') {
      const result = rule.validator(editValue.value)
      if (result !== true) {
        validationError.value = result || rule.message || '验证失败'
        return false
      }
    }

    if (rule.min !== undefined && editValue.value < rule.min) {
      validationError.value = rule.message || `不能小于 ${rule.min}`
      return false
    }

    if (rule.max !== undefined && editValue.value > rule.max) {
      validationError.value = rule.message || `不能大于 ${rule.max}`
      return false
    }
  }

  validationError.value = ''
  return true
}

/**
 * Handle Enter key
 */
async function handleEnter() {
  if (props.type === 'select') {
    // For select, don't save on enter - let user select
    return
  }

  await saveValue()
}

/**
 * Handle Escape key
 */
function handleEscape() {
  cancelEditing()
}

/**
 * Handle blur event
 */
async function handleBlur() {
  // For date picker and select, blur happens when selecting
  // Don't save immediately, give time for selection to complete
  if (props.type === 'date' || props.type === 'select') {
    setTimeout(async () => {
      if (isEditing.value) {
        await saveValue()
      }
    }, 200)
  } else {
    await saveValue()
  }
}

/**
 * Handle date picker visibility change
 */
function handleDateVisibleChange(visible) {
  if (!visible && isEditing.value) {
    // Date picker closed - save will be handled by blur
  }
}

/**
 * Save value
 */
async function saveValue() {
  if (!isEditing.value) return
  if (isSaving.value) return

  // Check if value changed
  if (editValue.value === props.modelValue) {
    cancelEditing()
    return
  }

  // Validate
  if (!validate()) {
    return
  }

  try {
    isSaving.value = true

    // Create update command
    const updates = {
      [props.field]: editValue.value
    }

    const command = new UpdateTaskCommand(
      props.taskId,
      updates,
      props.originalData
    )

    // Execute command (integrates with undo/redo)
    await undoRedoStore.executeCommand(command)

    // Emit change event
    emit('update:modelValue', editValue.value)
    emit('change', {
      field: props.field,
      value: editValue.value,
      taskId: props.taskId
    })

    // Emit gantt event
    eventBus.emit(GanttEvents.TASK_UPDATED, {
      taskId: props.taskId,
      updates
    })

    isEditing.value = false
  } catch (error) {
    console.error('Save failed:', error)
    validationError.value = error.message || '保存失败'
  } finally {
    isSaving.value = false
  }
}

/**
 * Cancel editing
 */
function cancelEditing() {
  isEditing.value = false
  editValue.value = props.modelValue
  validationError.value = ''
  emit('cancel', props.field)
}

/**
 * Expose methods for parent component
 */
defineExpose({
  startEditing,
  cancelEditing,
  saveValue,
  isEditing
})

// Keyboard shortcuts
function handleKeydown(event) {
  if (!isEditing.value) {
    // F2 to start editing
    if (event.key === 'F2') {
      event.preventDefault()
      startEditing()
    }
  } else {
    // Enter to save, Escape to cancel (handled by inline handlers)
    if (event.key === 'Enter' && props.type !== 'select') {
      event.preventDefault()
    }
    if (event.key === 'Escape') {
      event.preventDefault()
    }
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
  // 清理 select 定时器
  if (selectTimeout) {
    clearTimeout(selectTimeout)
    selectTimeout = null
  }
})
</script>

<style scoped>
.editable-cell {
  position: relative;
  width: 100%;
  height: 100%;
  min-height: 32px;
  display: flex;
  align-items: center;
  cursor: pointer;
  transition: background-color var(--transition-fast);
}

.editable-cell:hover:not(.is-readonly) {
  background-color: var(--row-hover-bg, #f5f7fa);
}

.editable-cell.is-editing {
  background-color: #fff;
  cursor: default;
}

.editable-cell.is-readonly {
  cursor: default;
}

.editable-cell__display {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 0 8px;
  gap: 8px;
}

.editable-cell__display span {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.editable-icon {
  opacity: 0;
  font-size: 14px;
  color: var(--color-primary, #409eff);
  transition: opacity var(--transition-fast);
  flex-shrink: 0;
}

.editable-cell:hover:not(.is-readonly) .editable-icon {
  opacity: 1;
}

.editable-cell__editor {
  width: 100%;
  padding: 0;
  position: relative;
}

.editable-cell__editor :deep(.el-input),
.editable-cell__editor :deep(.el-input-number),
.editable-cell__editor :deep(.el-select),
.editable-cell__editor :deep(.el-date-picker) {
  width: 100%;
}

.editable-cell__editor :deep(.el-input__wrapper),
.editable-cell__editor :deep(.el-input-number__wrapper) {
  box-shadow: 0 0 0 1px var(--color-primary, #409eff) inset;
}

.editable-cell__error {
  position: absolute;
  bottom: 100%;
  left: 0;
  right: 0;
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  background-color: #fef0f0;
  color: #f56c6c;
  font-size: 12px;
  border-radius: 4px 4px 0 0;
  z-index: 1000;
  white-space: nowrap;
}

.editable-cell__error .el-icon {
  font-size: 14px;
  flex-shrink: 0;
}

/* Status-specific styles */
.editable-cell :deep(.status-text) {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.editable-cell :deep(.status-dot) {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.editable-cell :deep(.status-dot--not_started) {
  background: var(--color-info, #909399);
}

.editable-cell :deep(.status-dot--in_progress) {
  background: var(--color-primary, #409eff);
}

.editable-cell :deep(.status-dot--completed) {
  background: var(--color-success, #67c23a);
}

.editable-cell :deep(.status-dot--delayed) {
  background: var(--color-danger, #f56c6c);
}

/* Priority-specific styles */
.editable-cell :deep(.priority-high) {
  color: #f56c6c;
  font-weight: 500;
}

.editable-cell :deep(.priority-medium) {
  color: #e6a23c;
  font-weight: 500;
}

.editable-cell :deep(.priority-low) {
  color: #67c23a;
  font-weight: 500;
}

/* Progress bar style */
.editable-cell :deep(.progress-display) {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
}

.editable-cell :deep(.progress-bar-bg) {
  flex: 1;
  height: 6px;
  background: #ebeef5;
  border-radius: 3px;
  overflow: hidden;
}

.editable-cell :deep(.progress-bar-fill) {
  height: 100%;
  background: var(--color-primary, #409eff);
  transition: width var(--transition-base);
}

.editable-cell :deep(.progress-text) {
  font-size: 12px;
  color: var(--color-text-secondary, #909399);
  min-width: 35px;
  text-align: right;
}
</style>
