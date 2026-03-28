<template>
  <div class="gantt-toolbar" :class="{ 'is-fullscreen': isFullscreen }">
    <div class="toolbar-left">
      <!-- 返回按钮 -->
      <el-button type="info" size="small" @click="$emit('back-to-list')" class="back-button">
        <el-icon><ArrowLeft /></el-icon>
        返回项目列表
      </el-button>

      <el-divider direction="vertical" style="margin: 0 12px;" />

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

      <!-- 时间轴格式 -->
      <el-select
        :model-value="timelineFormat"
        @change="$emit('timeline-format-change', $event)"
        size="small"
        style="width: 130px; margin-left: 12px"
      >
        <el-option label="日期" value="day"></el-option>
        <el-option label="月/日" value="month-day"></el-option>
        <el-option label="年/月" value="year-month"></el-option>
        <el-option label="年/月/日" value="year-month-day"></el-option>
        <el-option label="周" value="week"></el-option>
        <el-option label="月" value="month"></el-option>
        <el-option label="季度" value="quarter"></el-option>
      </el-select>

      <!-- 日期显示格式 -->
      <el-select
        v-if="showDateFormat"
        :model-value="dateDisplayFormat"
        @change="$emit('date-format-change', $event)"
        size="small"
        style="width: 130px; margin-left: 8px"
      >
        <el-option label="全部日期" value="all"></el-option>
        <el-option label="奇数日期" value="odd"></el-option>
        <el-option label="间隔3天" value="interval3"></el-option>
        <el-option label="间隔5天" value="interval5"></el-option>
        <el-option label="每月1号" value="first"></el-option>
      </el-select>

      <!-- 视图切换（甘特图/网络图） -->
      <el-button-group size="small" style="margin-left: 12px" title="视图模式">
        <el-button
          @click="$emit('toggle-view-mode', 'gantt')"
          :type="chartViewMode === 'gantt' ? 'primary' : 'default'"
          title="甘特图视图"
        >
          <el-icon>
            <svg viewBox="0 0 1024 1024" width="14" height="14">
              <path d="M128 128h768v768H128z" fill="none" stroke="currentColor" stroke-width="80"/>
              <rect x="160" y="400" width="120" height="40" fill="currentColor"/>
              <rect x="320" y="320" width="120" height="40" fill="currentColor"/>
              <rect x="480" y="240" width="120" height="40" fill="currentColor"/>
            </svg>
          </el-icon>
          甘特图
        </el-button>
        <el-button
          @click="$emit('toggle-view-mode', 'network')"
          :type="chartViewMode === 'network' ? 'primary' : 'default'"
          title="网络图视图"
        >
          <el-icon>
            <svg viewBox="0 0 1024 1024" width="14" height="14">
              <circle cx="200" cy="512" r="60" fill="currentColor"/>
              <circle cx="512" cy="200" r="60" fill="currentColor"/>
              <circle cx="824" cy="512" r="60" fill="currentColor"/>
              <circle cx="512" cy="824" r="60" fill="currentColor"/>
              <path d="M200 512 L512 200 L824 512 L512 824 Z" fill="none" stroke="currentColor" stroke-width="40"/>
            </svg>
          </el-icon>
          网络图
        </el-button>
      </el-button-group>

      <!-- 工具模式切换 -->
      <el-button-group size="small" style="margin-left: 12px" title="工具模式">
        <!-- 箭头选择工具 -->
        <el-button
          @click="$emit('toggle-select-mode')"
          :type="!panMode ? 'primary' : 'default'"
          title="选择/拖拽模式"
        >
          <el-icon>
            <svg viewBox="0 0 1024 1024" width="14" height="14">
              <path d="M896 448H560V112c0-17.7-14.3-32-32-32s-32 14.3-32 32v336H160c-17.7 0-32 14.3-32 32s14.3 32 32 32h336v336c0 17.7 14.3 32 32 32s32-14.3 32-32V512h336c17.7 0 32-14.3 32-32s-14.3-32-32-32z" fill="currentColor"></path>
            </svg>
          </el-icon>
        </el-button>
        <!-- 手形平移工具 -->
        <el-button
          @click="$emit('toggle-pan-mode')"
          :type="panMode ? 'primary' : 'default'"
          title="平移工具"
        >
          <el-icon>
            <svg viewBox="0 0 1024 1024" width="14" height="14">
              <path d="M768 256h-64V128c0-17.7-14.3-32-32-32s-32 14.3-32 32v128H384V128c0-17.7-14.3-32-32-32s-32 14.3-32 32v128h-64c-17.7 0-32 14.3-32 32s14.3 32 32 32h64v416c0 17.7 14.3 32 32 32h256v128c0 17.7 14.3 32 32 32s32-14.3 32-32V768h64c17.7 0 32-14.3 32-32s-14.3-32-32-32h-64V320h64c17.7 0 32-14.3 32-32s-14.3-32-32-32z" fill="currentColor"></path>
            </svg>
          </el-icon>
        </el-button>
      </el-button-group>

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
            <el-dropdown-item v-if="chartViewMode === 'gantt'" @click="$emit('toggle-dependencies')">
              <el-icon class="dropdown-icon" :class="{ 'is-active': showDependencies }"><Connection /></el-icon>
              显示依赖关系
            </el-dropdown-item>
            <el-dropdown-item @click="$emit('toggle-critical-path')">
              <el-icon class="dropdown-icon" :class="{ 'is-active': showCriticalPath }"><Flag /></el-icon>
              显示关键路径
            </el-dropdown-item>
            <el-dropdown-item v-if="chartViewMode === 'gantt'" @click="$emit('toggle-baseline')">
              <el-icon class="dropdown-icon" :class="{ 'is-active': showBaseline }"><Bottom /></el-icon>
              显示基线对比
            </el-dropdown-item>
            <el-dropdown-item v-if="chartViewMode === 'network'" @click="$emit('toggle-network-time-params')">
              <el-icon class="dropdown-icon" :class="{ 'is-active': $attrs.showNetworkTimeParams }"><Timer /></el-icon>
              显示时间参数
            </el-dropdown-item>
            <el-dropdown-item v-if="chartViewMode === 'network'" @click="$emit('toggle-network-task-names')">
              <el-icon class="dropdown-icon" :class="{ 'is-active': $attrs.showNetworkTaskNames }"><Document /></el-icon>
              显示任务名称
            </el-dropdown-item>
            <el-dropdown-item v-if="chartViewMode === 'network'" @click="$emit('toggle-network-slack')">
              <el-icon class="dropdown-icon" :class="{ 'is-active': $attrs.showNetworkSlack }"><Clock /></el-icon>
              显示时差信息
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>

      <!-- 网络图布局选项 -->
      <el-select
        v-if="chartViewMode === 'network'"
        :model-value="$attrs.networkLayoutMode"
        @change="$emit('network-layout-change', $event)"
        size="small"
        style="width: 120px; margin-left: 12px"
        placeholder="布局方式"
      >
        <el-option label="自动布局" value="auto"></el-option>
        <el-option label="从左到右" value="left-right"></el-option>
        <el-option label="从上到下" value="top-down"></el-option>
      </el-select>

      <!-- AOA图分析工具 -->
      <el-dropdown split-button type="warning" size="small" style="margin-left: 12px" v-if="chartViewMode === 'network'">
        <el-icon><Operation /></el-icon>
        AOA分析
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item @click="$emit('calculate-critical-path')">
              <el-icon><Flag /></el-icon>
              计算关键路径
            </el-dropdown-item>
            <el-dropdown-item @click="$emit('analyze-node-properties')">
              <el-icon><DataAnalysis /></el-icon>
              节点属性分析
            </el-dropdown-item>
            <el-dropdown-item @click="$emit('check-path-optimization')">
              <el-icon><TrendCharts /></el-icon>
              路径优化检查
            </el-dropdown-item>
            <el-dropdown-item divided @click="$emit('validate-rules')">
              <el-icon><CircleCheck /></el-icon>
              规则验证（R4/R11）
            </el-dropdown-item>
            <el-dropdown-item @click="$emit('export-analysis-report')">
              <el-icon><Download /></el-icon>
              导出分析报告
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
  Close,
  Timer,
  Clock,
  Operation,
  DataAnalysis,
  TrendCharts,
  CircleCheck,
  Download
} from '@element-plus/icons-vue'
import { ElDivider } from 'element-plus'

import { computed } from 'vue'

const props = defineProps({
  currentZoomLabel: {
    type: String,
    default: '40日'
  },
  currentPeriodText: {
    type: String,
    default: ''
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
  },
  timelineFormat: {
    type: String,
    default: 'month-day'
  },
  dateDisplayFormat: {
    type: String,
    default: 'all'
  },
  panMode: {
    type: Boolean,
    default: false
  },
  chartViewMode: {
    type: String,
    default: 'gantt' // 'gantt' or 'network'
  }
})

const showDateFormat = computed(() => {
  const format = props.timelineFormat
  return format === 'day' || format === 'month-day' || format === 'year-month-day'
})

defineEmits([
  'back-to-list',
  'navigate-date',
  'go-today',
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
  'save-all',
  'timeline-format-change',
  'date-format-change',
  'toggle-pan-mode',
  'toggle-select-mode',
  'toggle-view-mode'
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

.back-button {
  display: flex;
  align-items: center;
  gap: 4px;
}

.back-button .el-icon {
  font-size: 16px;
}
</style>
