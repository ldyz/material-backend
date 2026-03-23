<template>
  <div class="diff-field" :class="fieldClass">
    <div class="field-label">
      {{ formattedLabel }}
    </div>

    <div v-if="isUnchanged" class="field-value unchanged">
      {{ displayValue }}
    </div>

    <div v-else-if="type === 'unified'" class="field-value unified">
      <div class="unified-before">
        <el-tag size="small" type="danger">-</el-tag>
        <span class="value-text">{{ beforeDisplay }}</span>
      </div>
      <div class="unified-after">
        <el-tag size="small" type="success">+</el-tag>
        <span class="value-text">{{ afterDisplay }}</span>
      </div>
    </div>

    <div v-else class="field-value changed">
      {{ displayValue }}
    </div>
  </div>
</template>

<script setup>
/**
 * DiffField.vue
 *
 * Individual field in diff view with change highlighting.
 */

import { computed } from 'vue'

// ==================== Props ====================

const props = defineProps({
  field: {
    type: String,
    required: true
  },
  value: {
    type: [String, Number, Boolean, Object, Array],
    default: null
  },
  beforeValue: {
    type: [String, Number, Boolean, Object, Array, null],
    default: null
  },
  afterValue: {
    type: [String, Number, Boolean, Object, Array, null],
    default: null
  },
  type: {
    type: String,
    default: 'unified',
    validator: (value) => ['before', 'after', 'unified'].includes(value)
  }
})

// ==================== Computed ====================

/**
 * Format field label
 */
const formattedLabel = computed(() => {
  return props.field
    .split('_')
    .map(word => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ')
})

/**
 * Check if field is unchanged
 */
const isUnchanged = computed(() => {
  return JSON.stringify(props.beforeValue) === JSON.stringify(props.afterValue)
})

/**
 * Get field CSS class
 */
const fieldClass = computed(() => {
  if (isUnchanged.value) return 'is-unchanged'
  if (props.type === 'before') return 'is-before'
  if (props.type === 'after') return 'is-after'
  return 'is-changed'
})

/**
 * Get display value
 */
const displayValue = computed(() => {
  return formatValue(props.value, props.field)
})

/**
 * Get before display value
 */
const beforeDisplay = computed(() => {
  return formatValue(props.beforeValue, props.field)
})

/**
 * Get after display value
 */
const afterDisplay = computed(() => {
  return formatValue(props.afterValue, props.field)
})

/**
 * Format value for display
 */
function formatValue(value, field) {
  // Handle null/undefined
  if (value === null || value === undefined) {
    return '-'
  }

  // Handle dates
  if (field.toLowerCase().includes('date') || field.toLowerCase().includes('time')) {
    try {
      const date = new Date(value)
      if (!isNaN(date.getTime())) {
        return date.toLocaleString()
      }
    } catch {
      // Invalid date, return as-is
    }
  }

  // Handle booleans
  if (typeof value === 'boolean') {
    return value ? 'Yes' : 'No'
  }

  // Handle arrays
  if (Array.isArray(value)) {
    if (value.length === 0) return 'Empty'
    return value.join(', ')
  }

  // Handle objects
  if (typeof value === 'object') {
    const keys = Object.keys(value)
    if (keys.length === 0) return 'Empty'
    return keys.map(k => `${k}: ${value[k]}`).join(', ')
  }

  // Handle numbers
  if (typeof value === 'number') {
    return value.toLocaleString()
  }

  // Handle strings
  return String(value)
}
</script>

<style scoped lang="scss">
.diff-field {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 8px;
  border-radius: 4px;
  transition: background 0.3s;

  &.is-before {
    background: #fef0f0;
  }

  &.is-after {
    background: #f0f9ff;
  }

  &.is-changed {
    background: #fdf6ec;
  }

  &.is-unchanged {
    background: transparent;
  }
}

.field-label {
  font-size: 12px;
  font-weight: 600;
  color: #606266;
  text-transform: capitalize;
}

.field-value {
  font-size: 14px;
  color: #303133;
  word-break: break-word;

  &.unchanged {
    color: #909399;
  }

  &.changed {
    font-weight: 500;
  }
}

.unified-before,
.unified-after {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 0;

  .value-text {
    flex: 1;
  }
}

.unified-before {
  color: #f56c6c;
  text-decoration: line-through;
}

.unified-after {
  color: #67c23a;
}
</style>
