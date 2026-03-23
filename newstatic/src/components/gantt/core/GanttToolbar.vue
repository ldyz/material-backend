<template>
  <div class="gantt-toolbar">
    <!-- Left Section: Undo/Redo & Actions -->
    <div class="gantt-toolbar__left">
      <!-- Undo/Redo -->
      <div class="toolbar-group">
        <el-tooltip content="撤销 (Ctrl+Z)" placement="bottom">
          <el-button
            :icon="RefreshLeft"
            :disabled="!canUndo"
            circle
            size="small"
            @click="handleUndo"
          />
        </el-tooltip>

        <el-tooltip content="重做 (Ctrl+Y)" placement="bottom">
          <el-button
            :icon="RefreshRight"
            :disabled="!canRedo"
            circle
            size="small"
            @click="handleRedo"
          />
        </el-tooltip>

        <el-dropdown @command="handleHistoryCommand" trigger="click">
          <el-button size="small" text>
            <el-icon class="history-icon"><Clock /></el-icon>
            <span class="history-text">{{ undoCount }} / {{ redoCount }}</span>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item disabled>
                <div class="history-info">
                  <span>可撤销: {{ undoCount }} 项</span>
                  <span>可重做: {{ redoCount }} 项</span>
                </div>
              </el-dropdown-item>
              <el-dropdown-item divided command="clear-history">
                清空历史
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>

      <el-divider direction="vertical" />

      <!-- Quick Actions -->
      <div class="toolbar-group">
        <el-tooltip content="从模板创建" placement="bottom">
          <el-button
            :icon="DocumentAdd"
            size="small"
            @click="handleToggleTemplate"
          >
            模板
          </el-button>
        </el-tooltip>

        <el-tooltip :content="selectedCount > 0 ? `批量编辑 (${selectedCount}个任务)` : '批量编辑'" placement="bottom">
          <el-button
            :icon="Edit"
            :disabled="selectedCount === 0"
            size="small"
            @click="handleToggleBulkEdit"
          >
            批量编辑
          </el-button>
        </el-tooltip>
      </div>

      <el-divider direction="vertical" />

      <!-- View Mode -->
      <div class="toolbar-group">
        <el-button-group>
          <el-button
            :type="viewMode === 'day' ? 'primary' : ''"
            size="small"
            @click="handleViewChange('day')"
          >
            日视图
          </el-button>
          <el-button
            :type="viewMode === 'week' ? 'primary' : ''"
            size="small"
            @click="handleViewChange('week')"
          >
            周视图
          </el-button>
          <el-button
            :type="viewMode === 'month' ? 'primary' : ''"
            size="small"
            @click="handleViewChange('month')"
          >
            月视图
          </el-button>
          <el-button
            :type="viewMode === 'quarter' ? 'primary' : ''"
            size="small"
            @click="handleViewChange('quarter')"
          >
            季度视图
          </el-button>
        </el-button-group>
      </div>
    </div>

    <!-- Right Section: Zoom & Export -->
    <div class="gantt-toolbar__right">
      <!-- Zoom Controls -->
      <div class="toolbar-group">
        <el-tooltip content="缩小 (Ctrl+-)" placement="bottom">
          <el-button
            :icon="ZoomOut"
            circle
            size="small"
            @click="handleZoomOut"
          />
        </el-tooltip>

        <el-tooltip :content="`缩放: ${zoomLevel}`" placement="bottom">
          <el-button
            size="small"
            @click="handleZoomReset"
          >
            {{ zoomLevel }}
          </el-button>
        </el-tooltip>

        <el-tooltip content="放大 (Ctrl++)" placement="bottom">
          <el-button
            :icon="ZoomIn"
            circle
            size="small"
            @click="handleZoomIn"
          />
        </el-tooltip>
      </div>

      <el-divider direction="vertical" />

      <!-- Export & More -->
      <div class="toolbar-group">
        <el-dropdown @command="handleExport" trigger="click">
          <el-button :icon="Download" size="small">
            导出
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="pdf">导出为 PDF</el-dropdown-item>
              <el-dropdown-item command="excel">导出为 Excel</el-dropdown-item>
              <el-dropdown-item command="image">导出为图片</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>

        <el-tooltip :content="isFullscreen ? '退出全屏 (F11)' : '全屏 (F11)'" placement="bottom">
          <el-button
            :icon="isFullscreen ? ExitFullScreen : FullScreen"
            circle
            size="small"
            @click="handleToggleFullscreen"
          />
        </el-tooltip>

        <!-- Sync Status Indicator -->
        <div class="sync-status" :class="`sync-status--${syncStatus}`">
          <el-icon class="sync-icon">
            <component :is="syncStatusIcon" />
          </el-icon>
          <span class="sync-text">{{ syncStatusText }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { ElMessageBox } from 'element-plus'
import {
  RefreshLeft,
  RefreshRight,
  Clock,
  DocumentAdd,
  Edit,
  ZoomIn,
  ZoomOut,
  Download,
  FullScreen,
  ExitFullScreen,
  Check,
  Loading,
  Warning
} from '@element-plus/icons-vue'

/**
 * GanttToolbar Component - Enhanced Gantt Chart Toolbar
 *
 * Features:
 * - Undo/Redo controls with history dropdown
 * - Template quick create button
 * - Bulk edit button with selection count
 * - View mode toggle (Day/Week/Month/Quarter)
 * - Zoom controls with level indicator
 * - Export options (PDF/Excel/Image)
 * - Full-screen toggle
 * - Sync status indicator
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
  canUndo: {
    type: Boolean,
    default: false
  },
  canRedo: {
    type: Boolean,
    default: false
  },
  selectedCount: {
    type: Number,
    default: 0
  },
  viewMode: {
    type: String,
    default: 'day'
  },
  dayWidth: {
    type: Number,
    default: 40
  },
  syncStatus: {
    type: String,
    default: 'saved' // saved, unsaved, saving
  },
  isFullscreen: {
    type: Boolean,
    default: false
  }
})

// Emits
const emit = defineEmits([
  'undo',
  'redo',
  'clear-history',
  'zoom-in',
  'zoom-out',
  'zoom-reset',
  'view-change',
  'toggle-template',
  'toggle-bulk-edit',
  'toggle-fullscreen',
  'export'
])

// Computed
const zoomLevel = computed(() => {
  const labels = {
    day: '日',
    week: '周',
    month: '月',
    quarter: '季'
  }
  return `${props.dayWidth}${labels[props.viewMode] || ''}`
})

const syncStatusIcon = computed(() => {
  const icons = {
    saved: Check,
    unsaved: Warning,
    saving: Loading
  }
  return icons[props.syncStatus] || Check
})

const syncStatusText = computed(() => {
  const texts = {
    saved: '已保存',
    unsaved: '未保存',
    saving: '保存中'
  }
  return texts[props.syncStatus] || '未知'
})

// Methods
/**
 * Handle undo
 */
function handleUndo() {
  emit('undo')
}

/**
 * Handle redo
 */
function handleRedo() {
  emit('redo')
}

/**
 * Handle history dropdown command
 */
async function handleHistoryCommand(command) {
  if (command === 'clear-history') {
    try {
      await ElMessageBox.confirm(
        '确定要清空所有历史记录吗？此操作不可撤销。',
        '确认清空',
        {
          type: 'warning',
          confirmButtonText: '确定',
          cancelButtonText: '取消'
        }
      )
      emit('clear-history')
    } catch {
      // User cancelled
    }
  }
}

/**
 * Handle zoom in
 */
function handleZoomIn() {
  emit('zoom-in')
}

/**
 * Handle zoom out
 */
function handleZoomOut() {
  emit('zoom-out')
}

/**
 * Handle zoom reset
 */
function handleZoomReset() {
  emit('zoom-reset')
}

/**
 * Handle view mode change
 */
function handleViewChange(mode) {
  emit('view-change', mode)
}

/**
 * Handle toggle template dialog
 */
function handleToggleTemplate() {
  emit('toggle-template')
}

/**
 * Handle toggle bulk edit dialog
 */
function handleToggleBulkEdit() {
  emit('toggle-bulk-edit')
}

/**
 * Handle toggle fullscreen
 */
function handleToggleFullscreen() {
  emit('toggle-fullscreen')
}

/**
 * Handle export
 */
function handleExport(format) {
  emit('export', format)
}
</script>

<style scoped lang="scss">
.gantt-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 16px;
  background-color: #fff;
  border-bottom: 1px solid #e4e7ed;
  min-height: 56px;
  gap: 16px;
}

.gantt-toolbar__left,
.gantt-toolbar__right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.toolbar-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.el-divider--vertical {
  height: 24px;
  margin: 0;
}

.history-icon,
.history-text {
  vertical-align: middle;
}

.history-text {
  margin-left: 4px;
  font-size: 12px;
  color: var(--el-text-color-regular, #606266);
}

.history-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 12px;
  color: var(--el-text-color-regular, #606266);
}

.sync-status {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: 16px;
  font-size: 12px;
  font-weight: 500;
  transition: all 0.3s;

  &.sync-status--saved {
    background-color: #f0f9ff;
    color: #67c23a;
  }

  &.sync-status--unsaved {
    background-color: #fef0f0;
    color: #f56c6c;
  }

  &.sync-status--saving {
    background-color: #fdf6ec;
    color: #e6a23c;
  }
}

.sync-icon {
  font-size: 14px;
}

.sync-text {
  font-size: 12px;
}

/* Responsive */
@media (max-width: 1200px) {
  .gantt-toolbar {
    flex-wrap: wrap;
    min-height: auto;
    padding: 12px 16px;
  }

  .gantt-toolbar__left,
  .gantt-toolbar__right {
    flex-wrap: wrap;
  }
}

@media (max-width: 768px) {
  .gantt-toolbar__left {
    width: 100%;
    justify-content: space-between;
    margin-bottom: 8px;
  }

  .gantt-toolbar__right {
    width: 100%;
    justify-content: flex-end;
  }

  .toolbar-group .el-button-group {
    display: none;
  }

  .sync-text {
    display: none;
  }
}

/* Button group responsive */
@media (max-width: 992px) {
  .toolbar-group .el-button-group .el-button {
    padding: 8px 12px;
    font-size: 12px;
  }
}

@media (max-width: 768px) {
  .toolbar-group .el-button-group {
    .el-button:not(:first-child):not(:last-child) {
      display: none;
    }

    .el-button:first-child {
      border-top-right-radius: 4px;
      border-bottom-right-radius: 4px;
    }

    .el-button:last-child {
      display: block;
      border-top-left-radius: 4px;
      border-bottom-left-radius: 4px;
    }
  }
}
</style>
