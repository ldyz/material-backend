/**
 * Progress/Gantt components index
 * Centralized export of all gantt-related components
 */

// Core components - 使用重构版本作为默认导出
export { default as GanttChart } from './GanttChartRefactored.vue'
export { default as GanttChartRefactored } from './GanttChartRefactored.vue'

// 保留旧版本导出以保持向后兼容（可选）
export { default as GanttChartLegacy } from './GanttChart.vue'
export { default as GanttToolbar } from './GanttToolbar.vue'
export { default as GanttStats } from './GanttStats.vue'
export { default as GanttHeader } from './GanttHeader.vue'
export { default as GanttLegend } from './GanttLegend.vue'
export { default as GanttStatusBar } from './GanttStatusBar.vue'

// Task components
export { default as TaskList } from './TaskList.vue'
export { default as TaskTable } from './TaskTable.vue'
export { default as TaskTimeline } from './TaskTimeline.vue'

// Dialog components
export { default as GanttContextMenu } from './GanttContextMenu.vue'
export { default as TaskDetailDrawer } from './TaskDetailDrawer.vue'
export { default as TaskEditDialog } from './TaskEditDialog.vue'
export { default as ResourceAllocationDialog } from './ResourceAllocationDialog.vue'
export { default as ResourceManagementDialog } from './ResourceManagementDialog.vue'

// View components
export { default as NetworkDiagram } from './NetworkDiagram.vue'
export { default as ResourceView } from './ResourceView.vue'
export { default as SimpleGantt } from './SimpleGantt.vue'
export { default as ViewSwitcher } from './ViewSwitcher.vue'
