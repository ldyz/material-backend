<template>
  <div class="diff-viewer">
    <el-tabs v-model="activeTab" type="border-card">
      <!-- Side by Side View -->
      <el-tab-pane label="Side by Side" name="side-by-side">
        <div class="diff-side-by-side">
          <div class="diff-column">
            <div class="diff-header before">Before</div>
            <div class="diff-content">
              <DiffField
                v-for="(value, key) in before"
                :key="`before-${key}`"
                :field="key"
                :value="value"
                :after-value="after[key]"
                type="before"
              />
            </div>
          </div>

          <div class="diff-divider" />

          <div class="diff-column">
            <div class="diff-header after">After</div>
            <div class="diff-content">
              <DiffField
                v-for="(value, key) in after"
                :key="`after-${key}`"
                :field="key"
                :value="value"
                :before-value="before[key]"
                type="after"
              />
            </div>
          </div>
        </div>
      </el-tab-pane>

      <!-- Unified View -->
      <el-tab-pane label="Unified" name="unified">
        <div class="diff-unified">
          <DiffField
            v-for="(value, key) in allKeys"
            :key="`unified-${key}`"
            :field="key"
            :value="before[key]"
            :before-value="before[key]"
            :after-value="after[key]"
            type="unified"
          />
        </div>
      </el-tab-pane>

      <!-- JSON View -->
      <el-tab-pane label="JSON" name="json">
        <div class="diff-json">
          <el-button
            :icon="CopyDocument"
            size="small"
            @click="copyJSON"
          >
            Copy JSON
          </el-button>
          <pre><code>{{ formattedJSON }}</code></pre>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
/**
 * DiffViewer.vue
 *
 * Visual diff viewer with side-by-side, unified, and JSON views.
 */

import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { CopyDocument } from '@element-plus/icons-vue'
import DiffField from './DiffField.vue'

// ==================== Props ====================

const props = defineProps({
  before: {
    type: Object,
    default: () => ({})
  },
  after: {
    type: Object,
    default: () => ({})
  },
  entityType: {
    type: String,
    default: ''
  }
})

// ==================== State ====================

const activeTab = ref('side-by-side')

// ==================== Computed ====================

/**
 * Get all unique keys
 */
const allKeys = computed(() => {
  const keys = new Set([
    ...Object.keys(props.before),
    ...Object.keys(props.after)
  ])
  return Array.from(keys)
})

/**
 * Format JSON for display
 */
const formattedJSON = computed(() => {
  return JSON.stringify(
    {
      before: props.before,
      after: props.after
    },
    null,
    2
  )
})

// ==================== Methods ====================

/**
 * Copy JSON to clipboard
 */
function copyJSON() {
  navigator.clipboard.writeText(formattedJSON.value)
    .then(() => {
      ElMessage.success('JSON copied to clipboard')
    })
    .catch(() => {
      ElMessage.error('Failed to copy JSON')
    })
}

/**
 * Format field value for display
 */
function formatValue(value, field) {
  if (value === null || value === undefined) {
    return '-'
  }

  // Date fields
  if (field.includes('date') || field.includes('time')) {
    try {
      return new Date(value).toLocaleString()
    } catch {
      return value
    }
  }

  // Array fields
  if (Array.isArray(value)) {
    return value.join(', ')
  }

  // Objects
  if (typeof value === 'object') {
    return JSON.stringify(value)
  }

  return String(value)
}
</script>

<style scoped lang="scss">
.diff-viewer {
  background: #fff;
  border-radius: 4px;
}

.diff-side-by-side {
  display: flex;
  gap: 16px;
  min-height: 200px;
}

.diff-column {
  flex: 1;
  display: flex;
  flex-direction: column;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  overflow: hidden;
}

.diff-header {
  padding: 8px 12px;
  font-weight: 600;
  text-align: center;
  color: #fff;

  &.before {
    background: #f56c6c;
  }

  &.after {
    background: #67c23a;
  }
}

.diff-content {
  flex: 1;
  padding: 12px;
  overflow-y: auto;
}

.diff-divider {
  width: 1px;
  background: #e4e7ed;
}

.diff-unified {
  min-height: 200px;
  padding: 12px;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
}

.diff-json {
  position: relative;
  min-height: 200px;
  padding: 12px;
  background: #f5f7fa;
  border: 1px solid #e4e7ed;
  border-radius: 4px;

  .el-button {
    position: absolute;
    top: 8px;
    right: 8px;
  }

  pre {
    margin: 0;
    padding: 12px;
    background: #fff;
    border-radius: 4px;
    overflow-x: auto;

    code {
      font-family: 'Courier New', monospace;
      font-size: 12px;
      line-height: 1.6;
      color: #303133;
    }
  }
}
</style>
