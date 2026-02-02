<template>
  <div class="gantt-body" ref="ganttBodyRef">
    <!-- 任务表格 -->
    <TaskTable
      :tasks="tasks"
      :grouped-tasks="groupedTasks"
      :selected-task-id="selectedTaskId"
      :row-height="rowHeight"
      :show-critical-path="showCriticalPath"
      :group-mode="groupMode"
      :collapsed-groups="collapsedGroups"
      :empty-description="emptyDescription"
      @row-click="handleRowClick"
      @toggle-group="toggleGroup"
      @context-menu="handleTableContextMenu"
      @cell-edit="handleCellEdit"
      @task-dragged="handleTaskDragged"
    />

    <!-- 任务时间轴 -->
    <TaskTimeline
      :tasks="tasks"
      :raw-tasks="rawTasks"
      :selected-task-id="selectedTaskId"
      :view-mode="viewMode"
      :timeline-days="timelineDays"
      :day-width="dayWidth"
      :row-height="rowHeight"
      :task-height="taskHeight"
      :show-dependencies="showDependencies"
      :show-critical-path="showCriticalPath"
      :show-baseline="showBaseline"
      :visible-dependencies="visibleDependencies"
      :is-dragging="isDragging"
      :dragged-task="draggedTask"
      :tooltip-visible="tooltipVisible"
      :tooltip-position="tooltipPosition"
      :tooltip-text="tooltipText"
      :today-position="todayPosition"
      :arrow-marker-id="arrowMarkerId"
      :arrow-color="arrowColor"
      :empty-description="emptyDescription"
      :is-creating-dependency="isCreatingDependency"
      :source-task-id="sourceTaskId"
      :temp-line-end="tempLineEnd"
      @task-click="handleTaskClick"
      @task-dblclick="handleTaskDblClick"
      @task-mousedown="handleTaskMouseDown"
      @context-menu="handleContextMenu"
      @dependency-create="handleDependencyCreate"
      @mousemove="handleMouseMove"
      @task-dragged="handleTaskDragged"
      ref="taskTimelineRef"
    />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import TaskTable from './TaskTable.vue'
import TaskTimeline from './TaskTimeline.vue'

const props = defineProps({
  tasks: {
    type: Array,
    default: () => []
  },
  rawTasks: {
    type: Array,
    default: () => []
  },
  groupedTasks: {
    type: Array,
    default: () => []
  },
  selectedTaskId: {
    type: [Number, String],
    default: null
  },
  // 依赖关系创建相关
  isCreatingDependency: {
    type: Boolean,
    default: false
  },
  sourceTaskId: {
    type: [Number, String],
    default: null
  },
  tempLineEnd: {
    type: Object,
    default: null
  },
  viewMode: {
    type: String,
    default: 'day'
  },
  timelineDays: {
    type: Array,
    default: () => []
  },
  timelineWeeks: {
    type: Array,
    default: () => []
  },
  timelineMonths: {
    type: Array,
    default: () => []
  },
  timelineQuarters: {
    type: Array,
    default: () => []
  },
  dayWidth: {
    type: Number,
    default: 40
  },
  rowHeight: {
    type: Number,
    default: 60
  },
  taskHeight: {
    type: Number,
    default: 32
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
  collapsedGroups: {
    type: Set,
    default: () => new Set()
  },
  visibleDependencies: {
    type: Array,
    default: () => []
  },
  isDragging: {
    type: Boolean,
    default: false
  },
  draggedTask: {
    type: Object,
    default: null
  },
  tooltipVisible: {
    type: Boolean,
    default: false
  },
  tooltipPosition: {
    type: Object,
    default: () => ({ x: 0, y: 0 })
  },
  tooltipText: {
    type: String,
    default: ''
  },
  searchKeyword: {
    type: String,
    default: ''
  },
  timelineWidth: {
    type: Number,
    default: 800
  },
  arrowMarkerId: {
    type: String,
    default: 'arrow-marker'
  },
  arrowColor: {
    type: String,
    default: '#909399'
  },
  emptyDescription: {
    type: String,
    default: '暂无进度计划数据'
  },
  todayPosition: {
    type: Number,
    default: null
  }
})

const emit = defineEmits([
  'row-click',
  'task-click',
  'task-dblclick',
  'task-mousedown',
  'toggle-group',
  'context-menu',
  'add-task',
  'dependency-create',
  'mousemove',
  'cell-edit',
  'task-dragged'
])

const ganttBodyRef = ref(null)
const taskTimelineRef = ref(null)

// 行点击
const handleRowClick = (task) => {
  emit('row-click', task)
}

// 任务点击
const handleTaskClick = (task) => {
  emit('task-click', task)
}

// 任务双击
const handleTaskDblClick = (task) => {
  emit('task-dblclick', task)
}

// 拖拽事件处理
const handleTaskMouseDown = (event, task) => {
  emit('task-mousedown', event, task)
}

// 分组切换
const toggleGroup = (groupName) => {
  emit('toggle-group', groupName)
}

// 右键菜单
const handleContextMenu = (event) => {
  emit('context-menu', event)
}

// 依赖关系创建
const handleDependencyCreate = (task) => {
  emit('dependency-create', task)
}

// 鼠标移动
const handleMouseMove = (position) => {
  emit('mousemove', position)
}

// 表格右键菜单
const handleTableContextMenu = ({ event, task }) => {
  emit('context-menu', { event, task })
}

// 单元格编辑
const handleCellEdit = ({ taskId, updateData }) => {
  emit('cell-edit', { taskId, updateData })
}

// 任务拖拽
const handleTaskDragged = ({ fromTask, toTask }) => {
  console.log('TaskList - handleTaskDragged 被调用:', { fromTask, toTask })
  emit('task-dragged', { fromTask, toTask })
}

defineExpose({
  ganttBodyRef,
  taskTimelineRef
})
</script>

<style scoped>
.gantt-body {
  flex: 1;
  overflow-y: auto;
  overflow-x: auto;
  min-height: 0;
  position: relative;
  display: flex;
}

/* 确保表格和时间轴同步滚动 */
.gantt-body > * {
  flex-shrink: 0;
}

.gantt-body > *:last-child {
  flex: 1;
}
</style>
