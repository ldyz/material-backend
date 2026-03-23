<template>
  <div class="history-item">
    <!-- Summary -->
    <div class="history-summary">
      <div class="summary-main">
        <el-avatar :size="32" :src="user.avatar">
          {{ user.name?.charAt(0) || '?' }}
        </el-avatar>

        <div class="summary-content">
          <p class="summary-text">
            <span class="user-name">{{ user.name }}</span>
            <span class="action-text">{{ actionText }}</span>
            <span class="entity-text">{{ entityText }}</span>
          </p>
          <p v-if="change.description" class="summary-description">
            {{ change.description }}
          </p>
        </div>
      </div>

      <div class="summary-actions">
        <el-button
          :icon="View"
          size="small"
          text
          @click="$emit('toggle-diff', change.id)"
        >
          {{ showDiff ? 'Hide' : 'Show' }} Details
        </el-button>

        <el-dropdown v-if="canRollback" trigger="click">
          <el-button :icon="MoreFilled" size="small" text />
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item
                :icon="RefreshLeft"
                @click="$emit('rollback', change.id)"
              >
                Rollback Change
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>

    <!-- Diff View -->
    <el-collapse-transition>
      <div v-show="showDiff" class="history-diff">
        <DiffViewer
          :before="change.changes.before"
          :after="change.changes.after"
          :entity-type="change.entity_type"
        />
      </div>
    </el-collapse-transition>
  </div>
</template>

<script setup>
/**
 * HistoryItem.vue
 *
 * Individual history item with summary and diff view.
 */

import { computed } from 'vue'
import { View, MoreFilled, RefreshLeft } from '@element-plus/icons-vue'
import DiffViewer from './DiffViewer.vue'

// ==================== Props ====================

const props = defineProps({
  change: {
    type: Object,
    required: true
  },
  user: {
    type: Object,
    required: true
  },
  showDiff: {
    type: Boolean,
    default: false
  }
})

// ==================== Emits ====================

defineEmits(['toggle-diff', 'rollback'])

// ==================== Computed ====================

/**
 * Get action text
 */
const actionText = computed(() => {
  const actions = {
    create: 'created',
    update: 'updated',
    delete: 'deleted'
  }
  return actions[props.change.action_type] || 'modified'
})

/**
 * Get entity text
 */
const entityText = computed(() => {
  const entities = {
    task: 'task',
    dependency: 'dependency',
    resource: 'resource assignment'
  }
  const entity = entities[props.change.entity_type] || props.change.entity_type
  return `${entity} #${props.change.entity_id}`
})

/**
 * Check if can rollback
 */
const canRollback = computed(() => {
  return props.change.action_type === 'update' || props.change.action_type === 'delete'
})
</script>

<style scoped lang="scss">
.history-item {
  padding: 12px;
  background: #f5f7fa;
  border-radius: 8px;
  transition: background 0.3s;

  &:hover {
    background: #ecf5ff;
  }
}

.history-summary {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.summary-main {
  display: flex;
  gap: 12px;
  flex: 1;
}

.summary-content {
  flex: 1;
}

.summary-text {
  margin: 0 0 4px;
  font-size: 14px;
  color: #606266;
  line-height: 1.6;
}

.user-name {
  font-weight: 600;
  color: #303133;
}

.action-text {
  padding: 0 4px;
  font-weight: 500;
  color: #409eff;
}

.entity-text {
  font-family: 'Courier New', monospace;
  color: #67c23a;
}

.summary-description {
  margin: 4px 0 0;
  font-size: 12px;
  color: #909399;
}

.summary-actions {
  display: flex;
  gap: 4px;
}

.history-diff {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid #e4e7ed;
}
</style>
