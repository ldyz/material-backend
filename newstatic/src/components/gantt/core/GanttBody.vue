<template>
  <div
    class="gantt-body"
    :class="{ 'is-resizing': isResizing }"
    ref="bodyRef"
  >
    <!-- 任务列表和时间轴容器 -->
    <div class="gantt-content" :style="contentStyle">
      <!-- 任务表格 -->
      <TaskTable
        :tasks="tasks"
        :grouped-tasks="groupedTasks"
        :selected-task-id="selectedTask?.id"
        :view-mode="viewMode"
        :row-height="rowHeight"
        :show-dependencies="showDependencies"
        :show-critical-path="showCriticalPath"
        :group-mode="groupMode"
        :collapsed-groups="collapsedGroups"
        :search-keyword="searchKeyword"
        :is-mobile="isMobile"
        @row-click="handleRowClick"
        @task-click="handleTaskClick"
        @toggle-group="toggleGroup"
        @context-menu="handleContextMenu"
        @add-task="handleAddTask"
        @cell-edit="handleCellEdit"
      />

      <!-- 时间轴 -->
      <TaskTimeline
        :tasks="tasks"
        :view-mode="viewMode"
        :timeline-days="timelineDays"
        :timeline-weeks="timelineWeeks"
        :timeline-months="timelineMonths"
        :timeline-quarters="timelineQuarters"
        :day-width="dayWidth"
        :row-height="rowHeight"
        :task-height="taskHeight"
        :show-dependencies="showDependencies"
        :show-critical-path="showCriticalPath"
        :show-baseline="showBaseline"
        :selected-task-id="selectedTask?.id"
        :is-dragging="isDragging"
        :dragged-task="draggedTask"
        :preview-task="previewTask"
        :drag-mode="dragMode"
        :use-virtual-scroll="useVirtualScroll"
        @task-click="handleTaskClick"
        @task-dblclick="handleTaskDblClick"
        @task-mousedown="handleTaskMouseDown"
        @dependency-create="handleDependencyCreate"
      />
    </div>

    <!-- 调整大小手柄 -->
    <template v-if="!isMobile">
      <div
        class="resize-handle resize-handle-corner"
        @mousedown="handleResizeStart"
        title="拖动调整大小"
      >
        <svg width="16" height="16" viewBox="0 0 16 16">
          <path d="M12 4 L12 12 L4 12" stroke="#909399" stroke-width="2" fill="none"/>
          <path d="M9 4 L9 9 L4 9" stroke="#909399" stroke-width="2" fill="none"/>
        </svg>
      </div>

      <div
        class="resize-handle resize-handle-right"
        @mousedown="handleResizeStartRight"
      ></div>

      <div
        class="resize-handle resize-handle-bottom"
        @mousedown="handleResizeStartBottom"
      ></div>
    </template>

    <!-- 拖拽提示框 -->
    <Transition name="fade">
      <div
        v-if="tooltipVisible && tooltipText"
        class="gantt-drag-tooltip"
        :style="{ left: tooltipPosition.x + 'px', top: tooltipPosition.y + 'px' }"
      >
        {{ tooltipText }}
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, PropType } from 'vue'
import TaskTable from '@/components/progress/TaskTable.vue'
import TaskTimeline from '@/components/progress/TaskTimeline.vue'
import type { GanttTask, TooltipPosition } from '@/types/gantt'

/**
 * 甘特图主体组件
 * 包含任务表格和时间轴
 */
interface Props {
  // 数据
  tasks: any[]
  groupedTasks: any
  selectedTask: any
  timelineDays: any[]
  timelineWeeks: any[]
  timelineMonths: any[]
  timelineQuarters: any[]
  collapsedGroups: Set<string>

  // 尺寸
  containerSize: { width: number | null; height: number | null }
  rowHeight: number
  taskHeight: number
  dayWidth: number

  // 显示选项
  viewMode: string
  groupMode: string
  showDependencies: boolean
  showCriticalPath: boolean
  showBaseline: boolean
  searchKeyword: string

  // 拖拽状态
  isDragging: boolean
  draggedTask: any
  previewTask: any
  dragMode: string
  tooltipVisible: boolean
  tooltipPosition: TooltipPosition
  tooltipText: string

  // 其他
  isResizing: boolean
  isMobile: boolean
  useVirtualScroll: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
  rowClick: [task: any]
  taskClick: [task: any, event: Event]
  taskDblClick: [task: any]
  taskMousedown: [task: any, event: MouseEvent, taskBar: HTMLElement]
  toggleGroup: [groupKey: string]
  contextMenu: [task: any, position: { x: number; y: number }]
  addTask: [parentId?: string | number]
  cellEdit: [task: any, field: string, value: any]
  dependencyCreate: [data: any]
  resizeStart: [event: MouseEvent]
}>()

const bodyRef = ref<HTMLElement>()

/**
 * 内容区域样式
 */
const contentStyle = computed(() => {
  const style: Record<string, string> = {}

  if (props.containerSize.width) {
    style.width = props.containerSize.width + 'px'
  }

  if (props.containerSize.height) {
    style.height = props.containerSize.height + 'px'
  }

  return style
})

/**
 * 处理行点击
 */
function handleRowClick(task: any) {
  emit('rowClick', task)
}

/**
 * 处理任务点击
 */
function handleTaskClick(task: any, event: Event) {
  emit('taskClick', task, event)
}

/**
 * 处理任务双击
 */
function handleTaskDblClick(task: any) {
  emit('taskDblClick', task)
}

/**
 * 处理任务鼠标按下（开始拖拽）
 */
function handleTaskMouseDown(task: any, event: MouseEvent, taskBar: HTMLElement) {
  // 注意：TaskTimeline emit 的顺序是 event, task, taskBar
  emit('taskMousedown', event, task, taskBar)
}

/**
 * 切换分组折叠状态
 */
function toggleGroup(groupKey: string) {
  emit('toggleGroup', groupKey)
}

/**
 * 处理右键菜单
 */
function handleContextMenu(task: any, position: { x: number; y: number }) {
  emit('contextMenu', task, position)
}

/**
 * 添加任务
 */
function handleAddTask(parentId?: string | number) {
  emit('addTask', parentId)
}

/**
 * 单元格编辑
 */
function handleCellEdit(task: any, field: string, value: any) {
  emit('cellEdit', task, field, value)
}

/**
 * 创建依赖关系
 */
function handleDependencyCreate(data: any) {
  emit('dependencyCreate', data)
}

/**
 * 开始调整大小（右下角）
 */
function handleResizeStart(event: MouseEvent) {
  emit('resizeStart', event)
}

/**
 * 开始调整大小（右边）
 */
function handleResizeStartRight(event: MouseEvent) {
  event.preventDefault()
  emit('resizeStart', event)
}

/**
 * 开始调整大小（底边）
 */
function handleResizeStartBottom(event: MouseEvent) {
  event.preventDefault()
  emit('resizeStart', event)
}
</script>

<style scoped>
.gantt-body {
  position: relative;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.gantt-content {
  display: flex;
  flex: 1;
  overflow: hidden;
  position: relative;
}

/* 调整大小手柄 */
.resize-handle {
  position: absolute;
  z-index: 100;
  transition: background-color var(--transition-fast);
}

.resize-handle-corner {
  right: 0;
  bottom: 0;
  width: 20px;
  height: 20px;
  cursor: nwse-resize;
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0.5;
}

.resize-handle-corner:hover {
  opacity: 1;
}

.resize-handle-right {
  right: 0;
  top: 0;
  bottom: 20px;
  width: 5px;
  cursor: ew-resize;
  opacity: 0;
  transition: opacity var(--transition-fast);
}

.gantt-body:hover .resize-handle-right {
  opacity: 0.5;
}

.resize-handle-right:hover {
  opacity: 1;
}

.resize-handle-bottom {
  right: 20px;
  bottom: 0;
  left: 0;
  height: 5px;
  cursor: ns-resize;
  opacity: 0;
  transition: opacity var(--transition-fast);
}

.gantt-body:hover .resize-handle-bottom {
  opacity: 0.5;
}

.resize-handle-bottom:hover {
  opacity: 1;
}

.is-resizing .resize-handle {
  background: var(--color-primary);
  opacity: 0.3;
}

/* 拖拽提示框 */
.gantt-drag-tooltip {
  position: fixed;
  padding: 8px 12px;
  background: rgba(0, 0, 0, 0.8);
  color: #fff;
  font-size: 12px;
  border-radius: var(--radius-sm);
  pointer-events: none;
  z-index: var(--z-index-tooltip);
  white-space: nowrap;
}

/* 淡入淡出动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity var(--transition-fast);
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 响应式 */
@media (max-width: 768px) {
  .gantt-content {
    flex-direction: column;
  }
}
</style>
