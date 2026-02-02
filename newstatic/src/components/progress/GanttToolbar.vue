<template>
  <div class="gantt-toolbar" :class="{ 'is-fullscreen': isFullscreen }">
    <div class="toolbar-left">
      <!-- 日期导航 -->
      <el-button-group size="small">
        <el-button @click="$emit('navigate-date', -1)" title="向左">
          <el-icon><ArrowLeft /></el-icon>
        </el-button>
        <el-button @click="$emit('go-today')" title="今天">
          今天
        </el-button>
        <el-button @click="$emit('navigate-date', 1)" title="向右">
          <el-icon><ArrowRight /></el-icon>
        </el-button>
      </el-button-group>

      <span class="current-period">{{ currentPeriodText }}</span>

      <!-- 视图模式切换 -->
      <el-select
        :model-value="viewMode"
        @change="$emit('view-mode-change', $event)"
        size="small"
        style="width: 100px; margin-left: 12px"
      >
        <el-option label="日视图" value="day"></el-option>
        <el-option label="周视图" value="week"></el-option>
        <el-option label="月视图" value="month"></el-option>
        <el-option label="季度视图" value="quarter"></el-option>
      </el-select>

      <!-- 缩放控制 -->
      <el-button-group size="small" style="margin-left: 16px">
        <el-button @click="$emit('zoom-out')" title="缩小">
          <el-icon><ZoomOut /></el-icon>
        </el-button>
        <el-button @click="$emit('zoom-reset')" title="重置">
          {{ currentZoomLabel }}
        </el-button>
        <el-button @click="$emit('zoom-in')" title="放大">
          <el-icon><ZoomIn /></el-icon>
        </el-button>
      </el-button-group>

      <!-- 显示选项 -->
      <el-dropdown split-button type="primary" size="small" style="margin-left: 16px">
        显示选项
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item @click="$emit('toggle-dependencies')">
              <el-icon class="dropdown-icon" :class="{ 'is-active': showDependencies }"><Connection /></el-icon>
              显示依赖关系
            </el-dropdown-item>
            <el-dropdown-item @click="$emit('toggle-critical-path')">
              <el-icon class="dropdown-icon" :class="{ 'is-active': showCriticalPath }"><Flag /></el-icon>
              显示关键路径
            </el-dropdown-item>
            <el-dropdown-item @click="$emit('toggle-baseline')">
              <el-icon class="dropdown-icon" :class="{ 'is-active': showBaseline }"><Bottom /></el-icon>
              显示基线对比
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>

      <!-- 资源库管理 -->
      <el-button type="success" size="small" @click="$emit('open-resource-management')" style="margin-left: 12px">
        <el-icon style="margin-right: 4px;"><Setting /></el-icon>
        资源库管理
      </el-button>

      <!-- 分组选项 -->
      <el-select
        :model-value="groupMode"
        @change="$emit('group-change', $event)"
        size="small"
        style="width: 120px; margin-left: 12px"
        placeholder="分组方式"
        clearable
      >
        <el-option label="按状态分组" value="status"></el-option>
        <el-option label="按优先级分组" value="priority"></el-option>
      </el-select>

      <!-- 搜索框 -->
      <el-input
        :model-value="searchKeyword"
        @input="$emit('search', $event)"
        placeholder="搜索任务..."
        size="small"
        style="width: 200px; margin-left: 16px"
        clearable
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
    </div>

    <div class="toolbar-right">
      <!-- 保存按钮 -->
      <el-button
        type="success"
        size="small"
        @click="$emit('save-all')"
        :loading="isSaving"
        :disabled="!hasUnsavedChanges"
      >
        <el-icon style="margin-right: 4px;"><DocumentChecked /></el-icon>
        保存{{ hasUnsavedChanges ? ' (有未保存)' : '' }}
      </el-button>

      <el-button-group size="small">
        <el-button @click="$emit('export-png')" title="导出为PNG">
          <el-icon><Camera /></el-icon>
        </el-button>
        <el-button @click="$emit('export-pdf')" title="导出为PDF">
          <el-icon><Document /></el-icon>
        </el-button>
      </el-button-group>
      <el-button size="small" @click="$emit('auto-fit')" title="自动适应内容">
        <el-icon><FullScreen /></el-icon>
        自动适应
      </el-button>
      <el-button type="primary" size="small" @click="$emit('add-task')">
        <el-icon style="margin-right: 4px;"><Plus /></el-icon>
        添加任务
      </el-button>
      <el-button size="small" @click="$emit('refresh')">
        <el-icon><Refresh /></el-icon>
      </el-button>
      <!-- 全屏按钮 -->
      <el-button size="small" @click="$emit('toggle-fullscreen')" :title="isFullscreen ? '退出全屏' : '全屏'">
        <el-icon>
          <component :is="isFullscreen ? 'Close' : 'FullScreen'" />
        </el-icon>
      </el-button>
    </div>
  </div>
</template>

<script setup>
import {
  ZoomIn,
  ZoomOut,
  Search,
  ArrowLeft,
  ArrowRight,
  Refresh,
  Plus,
  Connection,
  Flag,
  Bottom,
  Camera,
  Document,
  DocumentChecked,
  Setting,
  FullScreen,
  Close
} from '@element-plus/icons-vue'

defineProps({
  viewMode: {
    type: String,
    default: 'day'
  },
  currentPeriodText: {
    type: String,
    default: ''
  },
  currentZoomLabel: {
    type: String,
    default: '40日'
  },
  showDependencies: {
    type: Boolean,
    default: true
  },
  showCriticalPath: {
    type: Boolean,
    default: true
  },
  showBaseline: {
    type: Boolean,
    default: false
  },
  groupMode: {
    type: String,
    default: ''
  },
  searchKeyword: {
    type: String,
    default: ''
  },
  isFullscreen: {
    type: Boolean,
    default: false
  },
  isSaving: {
    type: Boolean,
    default: false
  },
  hasUnsavedChanges: {
    type: Boolean,
    default: false
  }
})

defineEmits([
  'navigate-date',
  'go-today',
  'view-mode-change',
  'zoom-in',
  'zoom-out',
  'zoom-reset',
  'toggle-dependencies',
  'toggle-critical-path',
  'toggle-baseline',
  'open-resource-management',
  'group-change',
  'search',
  'export-png',
  'export-pdf',
  'auto-fit',
  'add-task',
  'refresh',
  'toggle-fullscreen',
  'save-all'
])
</script>

<style scoped>
.gantt-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: #fff;
  border-bottom: 1px solid #dcdfe6;
  flex-wrap: wrap;
  gap: 12px;
}

.gantt-toolbar.is-fullscreen {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 9999;
  border-radius: 0;
}

.toolbar-left,
.toolbar-right {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.current-period {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
  margin-left: 12px;
}

.dropdown-icon {
  margin-right: 8px;
  color: #909399;
}

.dropdown-icon.is-active {
  color: #409eff;
}
</style>
