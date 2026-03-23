<template>
  <div class="gantt-status-bar">
    <!-- Left: Undo/Redo Status -->
    <div class="status-bar__left">
      <div class="status-item status-item--history">
        <el-icon class="status-icon">
          <Clock />
        </el-icon>
        <span class="status-text">
          <template v-if="undoCount > 0 || redoCount > 0">
            {{ undoCount }} 项可撤销 · {{ redoCount }} 项可重做
          </template>
          <template v-else>
            无历史记录
          </template>
        </span>
      </div>

      <el-divider direction="vertical" />

      <div class="status-item status-item--connection">
        <el-icon
          class="status-icon connection-icon"
          :class="`connection-icon--${connectionStatus}`"
        >
          <component :is="connectionIcon" />
        </el-icon>
        <span class="status-text">{{ connectionText }}</span>
      </div>
    </div>

    <!-- Center: Task Counts -->
    <div class="status-bar__center">
      <div class="status-item status-item--tasks">
        <span class="status-text">
          总计 <strong>{{ taskCount }}</strong> 个任务
          <template v-if="selectedCount > 0">
            · 已选择 <strong>{{ selectedCount }}</strong> 个
          </template>
        </span>
      </div>
    </div>

    <!-- Right: Last Saved Time -->
    <div class="status-bar__right">
      <div class="status-item status-item--saved">
        <el-icon class="status-icon" :class="{ 'is-spinning': isSaving }">
          <component :is="saveIcon" />
        </el-icon>
        <span class="status-text">{{ lastSaved }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { Clock, CircleCheck, CircleClose, Loading, Finished } from '@element-plus/icons-vue'

/**
 * GanttStatusBar Component - Gantt Chart Status Bar
 *
 * Features:
 * - Undo/Redo status (X changes to undo, Y changes to redo)
 * - Connection status for future collaboration
 * - Task count with selection count
 * - Last saved time with loading state
 *
 * @date 2025-02-18
 */

// Props
const props = defineProps({
  undoCount: {
    type: Number,
    default: 0
  },
  redoCount: {
    type: Number,
    default: 0
  },
  taskCount: {
    type: Number,
    default: 0
  },
  selectedCount: {
    type: Number,
    default: 0
  },
  lastSaved: {
    type: String,
    default: '未保存'
  },
  connectionStatus: {
    type: String,
    default: 'connected' // connected, disconnected, syncing
  },
  isSaving: {
    type: Boolean,
    default: false
  }
})

// Computed
const connectionIcon = computed(() => {
  const icons = {
    connected: CircleCheck,
    disconnected: CircleClose,
    syncing: Loading
  }
  return icons[props.connectionStatus] || CircleCheck
})

const connectionText = computed(() => {
  const texts = {
    connected: '已连接',
    disconnected: '未连接',
    syncing: '同步中'
  }
  return texts[props.connectionStatus] || '未知'
})

const saveIcon = computed(() => {
  return props.isSaving ? Loading : Finished
})
</script>

<style scoped lang="scss">
.gantt-status-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 6px 16px;
  background-color: #f5f7fa;
  border-top: 1px solid #e4e7ed;
  min-height: 32px;
  font-size: 12px;
  gap: 16px;
}

.status-bar__left,
.status-bar__right {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
}

.status-bar__center {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  min-width: 0;
}

.status-item {
  display: flex;
  align-items: center;
  gap: 6px;
  color: var(--el-text-color-regular, #606266);
}

.status-icon {
  font-size: 14px;
  flex-shrink: 0;

  &.is-spinning {
    animation: rotating 2s linear infinite;
  }

  &.connection-icon {
    &--connected {
      color: var(--el-color-success, #67c23a);
    }

    &--disconnected {
      color: var(--el-color-danger, #f56c6c);
    }

    &--syncing {
      color: var(--el-color-warning, #e6a23c);
      animation: rotating 2s linear infinite;
    }
  }
}

.status-text {
  font-size: 12px;
  line-height: 1.5;
  white-space: nowrap;

  strong {
    font-weight: 500;
    color: var(--el-text-color-primary, #303133);
  }
}

.el-divider--vertical {
  height: 16px;
  margin: 0;
}

@keyframes rotating {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

/* Responsive */
@media (max-width: 768px) {
  .gantt-status-bar {
    flex-wrap: wrap;
    padding: 8px 12px;
    gap: 8px;
  }

  .status-bar__left,
  .status-bar__right {
    flex: 1;
    min-width: 0;
  }

  .status-bar__center {
    order: 3;
    width: 100%;
    margin-top: 4px;
  }

  .status-item--connection {
    display: none;
  }

  .status-text {
    font-size: 11px;
  }
}

@media (max-width: 480px) {
  .status-bar__left {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }

  .status-bar__right {
    flex-direction: column;
    align-items: flex-end;
    gap: 4px;
  }

  .el-divider--vertical {
    display: none;
  }
}
</style>
