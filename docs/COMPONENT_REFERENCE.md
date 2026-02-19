# Gantt Chart Editor - Component Reference

**Version:** 1.0.0
**Last Updated:** 2026-02-19
**Author:** Material Management System Team

---

## Table of Contents

1. [Core Components](#core-components)
2. [Timeline Components](#timeline-components)
3. [Table Components](#table-components)
4. [View Components](#view-components)
5. [Dashboard Components](#dashboard-components)
6. [Panel Components](#panel-components)
7. [Dialog Components](#dialog-components)
8. [Overlay Components](#overlay-components)
9. [Mobile Components](#mobile-components)
10. [Utility Components](#utility-components)

---

## Core Components

### GanttEditor

The main container component that orchestrates all Gantt chart functionality.

**Location:** `/src/components/gantt/core/GanttEditor.vue`

#### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `projectId` | `string \| number` | **required** | Unique identifier for the project |
| `tasks` | `Task[]` | `[]` | Array of task objects to display |
| `dependencies` | `Dependency[]` | `[]` | Array of task dependencies |
| `readOnly` | `boolean` | `false` | Whether the editor is read-only |
| `height` | `string \| number` | `'600px'` | Height of the Gantt editor |
| `initialViewMode` | `ViewMode` | `'gantt'` | Initial view mode (gantt/kanban/calendar/dashboard) |
| `initialZoom` | `number` | `40` | Initial zoom level (pixels per day) |
| `showMinimap` | `boolean` | `true` | Whether to show the minimap |
| `showGuidedTour` | `boolean` | `true` | Whether to show the guided tour for first-time users |
| `enableCollaboration` | `boolean` | `true` | Enable real-time collaboration via WebSocket |
| `autoSave` | `boolean` | `true` | Enable automatic saving |
| `autoSaveInterval` | `number` | `30000` | Auto-save interval in milliseconds |

#### Events

| Event | Payload | Description |
|-------|---------|-------------|
| `update:tasks` | `Task[]` | Emitted when tasks are updated |
| `update:dependencies` | `Dependency[]` | Emitted when dependencies are updated |
| `task-select` | `Task` | Emitted when a task is selected |
| `task-update` | `{ taskId: string, updates: Partial<Task> }` | Emitted when a task is updated |
| `task-create` | `Task` | Emitted when a new task is created |
| `task-delete` | `string` | Emitted when a task is deleted (taskId) |
| `dependency-create` | `Dependency` | Emitted when a dependency is created |
| `dependency-delete` | `string` | Emitted when a dependency is deleted (dependencyId) |
| `view-change` | `ViewMode` | Emitted when view mode changes |
| `zoom-change` | `number` | Emitted when zoom level changes |
| `export` | `{ format: ExportFormat, data: any }` | Emitted when data is exported |
| `save` | `{ tasks: Task[], dependencies: Dependency[] }` | Emitted when data is saved |

#### Methods

| Method | Parameters | Returns | Description |
|--------|-----------|---------|-------------|
| `loadData()` | `{ tasks?: Task[], dependencies?: Dependency[] }` | `Promise<void>` | Load data into the editor |
| `saveData()` | - | `Promise<void>` | Manually trigger save |
| `exportData()` | `format: ExportFormat` | `Promise<Blob>` | Export data in specified format |
| `zoomIn()` | - | `void` | Increase zoom level |
| `zoomOut()` | - | `void` | Decrease zoom level |
| `resetZoom()` | - | `void` | Reset zoom to default |
| `setViewMode()` | `mode: ViewMode` | `void` | Change view mode |
| `addTask()` | `task: Partial<Task>` | `string` | Add a new task, returns taskId |
| `updateTask()` | `taskId: string, updates: Partial<Task>` | `Promise<void>` | Update a task |
| `deleteTask()` | `taskId: string` | `Promise<void>` | Delete a task |
| `addDependency()` | `dependency: Partial<Dependency>` | `string` | Add a dependency, returns dependencyId |
| `deleteDependency()` | `dependencyId: string` | `Promise<void>` | Delete a dependency |
| `undo()` | - | `void` | Undo last action |
| `redo()` | - | `void` | Redo last undone action |
| `startTour()` | `tourType?: TourType` | `void` | Start guided tour |
| `showBulkEdit()` | `taskIds: string[]` | `void` | Show bulk edit dialog |
| `showResourceLeveling()` | - | `void` | Show resource leveling dialog |
| `calculateCriticalPath()` | - | `string[]` | Calculate and return critical path task IDs |

#### Slots

| Slot | Props | Description |
|------|-------|-------------|
| `toolbar-left` | - | Custom content for left side of toolbar |
| `toolbar-right` | - | Custom content for right side of toolbar |
| `task-bar` | `{ task: Task, days: Day[], dayWidth: number }` | Custom task bar rendering |
| `task-list-cell` | `{ task: Task, column: Column }` | Custom task list cell rendering |
| `status-bar-left` | - | Custom content for left side of status bar |
| `status-bar-right` | - | Custom content for right side of status bar |

#### Usage Example

```vue
<template>
  <GanttEditor
    ref="ganttRef"
    :project-id="projectId"
    :tasks="tasks"
    :dependencies="dependencies"
    :read-only="false"
    height="800px"
    :initial-view-mode="'gantt'"
    @update:tasks="handleTasksUpdate"
    @task-select="handleTaskSelect"
    @save="handleSave"
  >
    <template #toolbar-right>
      <el-button @click="customAction">Custom Action</el-button>
    </template>

    <template #task-bar="{ task, days, dayWidth }">
      <CustomTaskBar :task="task" :days="days" :day-width="dayWidth" />
    </template>
  </GanttEditor>
</template>

<script setup>
import { ref } from 'vue'
import { GanttEditor } from '@/components/gantt'

const ganttRef = ref(null)
const projectId = ref('proj-123')
const tasks = ref([])
const dependencies = ref([])

const handleTasksUpdate = (newTasks) => {
  tasks.value = newTasks
}

const handleTaskSelect = (task) => {
  console.log('Selected task:', task)
}

const handleSave = (data) => {
  // Save to backend
  api.saveProjectData(projectId.value, data)
}

// Call methods
const exportToPDF = async () => {
  const blob = await ganttRef.value.exportData('pdf')
  // Download file
}
</script>
```

---

### GanttToolbar

Toolbar component with undo/redo, zoom, view switching, and other controls.

**Location:** `/src/components/gantt/core/GanttToolbar.vue`

#### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `undoCount` | `number` | `0` | Number of available undo actions |
| `redoCount` | `number` | `0` | Number of available redo actions |
| `canUndo` | `boolean` | `false` | Whether undo is available |
| `canRedo` | `boolean` | `false` | Whether redo is available |
| `selectedCount` | `number` | `0` | Number of selected tasks |
| `viewMode` | `ViewMode` | `'gantt'` | Current view mode |
| `dayWidth` | `number` | `40` | Current day width in pixels |
| `syncStatus` | `'synced' \| 'syncing' \| 'unsaved'` | `'synced'` | Synchronization status |
| `showTemplates` | `boolean` | `true` | Show templates button |
| `showBulkEdit` | `boolean` | `true` | Show bulk edit button |
| `showExport` | `boolean` | `true` | Show export button |

#### Events

| Event | Payload | Description |
|-------|---------|-------------|
| `undo` | - | User clicked undo |
| `redo` | - | User clicked redo |
| `zoom-in` | - | User clicked zoom in |
| `zoom-out` | - | User clicked zoom out |
| `zoom-reset` | - | User clicked reset zoom |
| `view-change` | `ViewMode` | User changed view mode |
| `toggle-template` | - | User clicked templates |
| `toggle-bulk-edit` | - | User clicked bulk edit |
| `toggle-fullscreen` | - | User clicked fullscreen |
| `export` | `ExportFormat` | User clicked export |

#### Usage Example

```vue
<template>
  <GanttToolbar
    :can-undo="canUndo"
    :can-redo="canRedo"
    :selected-count="selectedCount"
    :view-mode="viewMode"
    :day-width="dayWidth"
    :sync-status="syncStatus"
    @undo="handleUndo"
    @redo="handleRedo"
    @zoom-in="handleZoomIn"
    @view-change="handleViewChange"
  />
</template>
```

---

### GanttStatusBar

Status bar displaying task counts, selection info, and connection status.

**Location:** `/src/components/gantt/core/GanttStatusBar.vue`

#### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `undoCount` | `number` | `0` | Number of undo actions available |
| `redoCount` | `number` | `0` | Number of redo actions available |
| `taskCount` | `number` | `0` | Total number of tasks |
| `selectedCount` | `number` | `0` | Number of selected tasks |
| `lastSaved` | `Date \| null` | `null` | Last saved timestamp |
| `connectionStatus` | `'connected' \| 'disconnected' \| 'connecting'` | `'disconnected'` | WebSocket connection status |
| `saving` | `boolean` | `false` | Currently saving |

#### Events

| Event | Payload | Description |
|-------|---------|-------------|
| `show-history` | - | User clicked to show history panel |

#### Usage Example

```vue
<template>
  <GanttStatusBar
    :task-count="tasks.length"
    :selected-count="selectedTaskIds.size"
    :last-saved="lastSavedTime"
    :connection-status="connectionStatus"
  />
</template>
```

---

## Timeline Components

### VirtualTimeline

High-performance timeline with virtual scrolling support.

**Location:** `/src/components/gantt/timeline/VirtualTimeline.vue`

#### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `tasks` | `Task[]` | **required** | Array of tasks to display |
| `days` | `Day[]` | **required** | Array of days to display |
| `rowHeight` | `number` | `50` | Height of each row in pixels |
| `dayWidth` | `number` | `40` | Width of each day in pixels |
| `containerHeight` | `string \| number` | `'600px'` | Height of the container |
| `showToday` | `boolean` | `true` | Show today marker line |
| `loading` | `boolean` | `false` | Show loading state |
| `bufferSize` | `number` | `200` | Buffer size for virtual scrolling |

#### Events

| Event | Payload | Description |
|-------|---------|-------------|
| `scroll` | `{ scrollTop: number, scrollLeft: number }` | Timeline scrolled |
| `resize` | `{ width: number, height: number }` | Timeline resized |
| `task-click` | `Task` | Task clicked |
| `task-dblclick` | `Task` | Task double-clicked |

#### Slots

| Slot | Props | Description |
|------|-------|-------------|
| `row` | `{ task: Task, index: number, days: Day[], dayWidth: number }` | Custom row rendering |
| `background` | `{ days: Day[], dayWidth: number }` | Custom background rendering |

#### Usage Example

```vue
<template>
  <VirtualTimeline
    :tasks="tasks"
    :days="timelineDays"
    :row-height="50"
    :day-width="dayWidth"
    :container-height="containerHeight"
    @scroll="handleScroll"
  >
    <template #row="{ task, days, dayWidth }">
      <TaskBar :task="task" :days="days" :day-width="dayWidth" />
    </template>
  </VirtualTimeline>
</template>
```

---

### TaskBar

Individual task bar component on the timeline.

**Location:** `/src/components/gantt/timeline/TaskBar.vue`

#### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `task` | `Task` | **required** | Task object to display |
| `days` | `Day[]` | **required** | Days array for positioning |
| `dayWidth` | `number` | **required** | Width of each day in pixels |
| `rowHeight` | `number` | `50` | Height of the row |
| `showProgress` | `boolean` | `true` | Show progress indicator |
| `showDependencyHandles` | `boolean` | `true` | Show handles for creating dependencies |
| `draggable` | `boolean` | `true` | Allow dragging to move |
| `resizable` | `boolean` | `true` | Allow resizing |

#### Events

| Event | Payload | Description |
|-------|---------|-------------|
| `click` | `Task` | Task bar clicked |
| `dblclick` | `Task` | Task bar double-clicked |
| `drag-start` | `{ task: Task, event: DragEvent }` | Drag started |
| `drag` | `{ task: Task, deltaX: number, deltaY: number }` | Dragging |
| `drag-end` | `{ task: Task, newStartDate: Date, newEndDate: Date }` | Drag ended |
| `resize-start` | `{ task: Task, edge: 'left' \| 'right', event: MouseEvent }` | Resize started |
| `resize` | `{ task: Task, deltaDays: number }` | Resizing |
| `resize-end` | `{ task: Task, newStartDate: Date, newEndDate: Date }` | Resize ended |
| `dependency-create` | `{ fromTask: Task, toTask: Task }` | Dependency creation requested |

#### Usage Example

```vue
<template>
  <TaskBar
    :task="task"
    :days="days"
    :day-width="dayWidth"
    :show-progress="true"
    :draggable="!readOnly"
    :resizable="!readOnly"
    @drag-end="handleTaskMove"
    @resize-end="handleTaskResize"
    @dependency-create="handleDependencyCreate"
  />
</template>
```

---

### DependencyLines

SVG component for rendering dependency lines between tasks.

**Location:** `/src/components/gantt/timeline/DependencyLines.vue`

#### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `dependencies` | `Dependency[]` | **required** | Array of dependencies |
| `tasks` | `Task[]` | **required** | Array of tasks |
| `days` | `Day[]` | **required** | Days array |
| `dayWidth` | `number` | **required** | Width of each day |
| `rowHeight` | `number` | `50` | Height of each row |
| `lineColor` | `string` | `'#409EFF'` | Color of dependency lines |
| `lineWidth` | `number` | `2` | Width of dependency lines |
| `arrowSize` | `number` | `6` | Size of arrow heads |
| `curved` | `boolean` | `true` | Use curved lines instead of straight |

#### Events

| Event | Payload | Description |
|-------|---------|-------------|
| `dependency-click` | `Dependency` | Dependency line clicked |
| `dependency-delete` | `string` | Delete dependency requested (dependencyId) |

#### Usage Example

```vue
<template>
  <DependencyLines
    :dependencies="dependencies"
    :tasks="tasks"
    :days="days"
    :day-width="dayWidth"
    :line-color="'#409EFF'"
    :curved="true"
    @dependency-click="handleDependencyClick"
    @dependency-delete="handleDependencyDelete"
  />
</template>
```

---

## Table Components

### VirtualTaskList

Virtual scrolling task list for performance.

**Location:** `/src/components/gantt/table/VirtualTaskList.vue`

#### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `tasks` | `Task[]` | **required** | Array of tasks to display |
| `columns` | `Column[]` | **required** | Column definitions |
| `rowHeight` | `number` | `50` | Height of each row in pixels |
| `containerHeight` | `string \| number` | `'600px'` | Height of the container |
| `collapsedTasks` | `Set<string>` | `new Set()` | Set of collapsed task IDs |
| `selectedTaskId` | `string \| null` | `null` | Currently selected task ID |
| `loading` | `boolean` | `false` | Show loading state |
| `searchKeyword` | `string` | `''` | Search keyword for filtering |

#### Events

| Event | Payload | Description |
|-------|---------|-------------|
| `scroll` | `{ scrollTop: number }` | Task list scrolled |
| `row-click` | `Task` | Row clicked |
| `row-dblclick` | `Task` | Row double-clicked |
| `toggle` | `string` | Toggle task collapse (taskId) |
| `edit` | `{ taskId: string, field: string, value: any }` | Cell edited |
| `selection-change` | `string[]` | Selection changed (array of taskIds) |

#### Slots

| Slot | Props | Description |
|------|-------|-------------|
| `cell` | `{ task: Task, column: Column, value: any }` | Custom cell rendering |

#### Usage Example

```vue
<template>
  <VirtualTaskList
    :tasks="tasks"
    :columns="columns"
    :row-height="50"
    :container-height="600"
    :collapsed-tasks="collapsedTasks"
    :selected-task-id="selectedTaskId"
    @row-click="handleTaskClick"
    @edit="handleCellEdit"
  >
    <template #cell="{ task, column, value }">
      <template v-if="column.prop === 'status'">
        <el-tag :type="getStatusType(value)">{{ value }}</el-tag>
      </template>
      <template v-else>
        {{ value }}
      </template>
    </template>
  </VirtualTaskList>
</template>
```

---

### EditableCell

Editable cell component for task list.

**Location:** `/src/components/gantt/table/EditableCell.vue`

#### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `value` | `any` | **required** | Cell value |
| `task` | `Task` | **required** | Task object |
| `column` | `Column` | **required** | Column definition |
| `editable` | `boolean` | `true` | Whether cell is editable |
| `type` | `'text' \| 'number' \| 'date' \| 'select'` | `'text'` | Input type |
| `options` | `any[]` | `[]` | Options for select type |

#### Events

| Event | Payload | Description |
|-------|---------|-------------|
| `update:value` | `any` | Value changed |
| `edit` | `{ task: Task, field: string, value: any }` | Edit committed |
| `cancel` | - | Edit cancelled |

#### Usage Example

```vue
<template>
  <EditableCell
    :value="task.status"
    :task="task"
    :column="{ prop: 'status', label: 'Status' }"
    type="select"
    :options="['Not Started', 'In Progress', 'Completed']"
    @edit="handleCellEdit"
  />
</template>
```

---

## View Components

### KanbanView

Kanban board view of tasks.

**Location:** `/src/components/gantt/views/KanbanView.vue`

#### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `tasks` | `Task[]` | **required** | Array of tasks |
| `groupBy` | `'status' \| 'priority' \| 'assignee'` | `'status'` | Field to group by |
| `cardView` | `'compact' \| 'detailed'` | `'compact'` | Card display style |
| `showWipLimits` | `boolean` | `false` | Show WIP limits |
| `enableSwimlanes` | `boolean` | `false` | Enable swimlanes |
| `showAvatars` | `boolean` | `true` | Show assignee avatars |
| `draggable` | `boolean` | `true` | Allow dragging cards |

#### Events

| Event | Payload | Description |
|-------|---------|-------------|
| `card-click` | `Task` | Card clicked |
| `card-move` | `{ task: Task, fromColumn: string, toColumn: string }` | Card moved |
| `status-change` | `{ taskId: string, newStatus: string }` | Status changed |

#### Usage Example

```vue
<template>
  <KanbanView
    :tasks="tasks"
    group-by="status"
    :card-view="'detailed'"
    :show-wip-limits="true"
    :draggable="!readOnly"
    @card-move="handleCardMove"
  />
</template>
```

---

### CalendarView

Calendar view of tasks.

**Location:** `/src/components/gantt/views/CalendarView.vue`

#### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `tasks` | `Task[]` | **required** | Array of tasks |
| `currentDate` | `Date` | `new Date()` | Current month to display |
| `view` | `'month' \| 'week' \| 'day'` | `'month'` | Calendar view type |
| `showWeekends` | `boolean` | `true` | Highlight weekends |
| `showToday` | `boolean` | `true` | Highlight today |

#### Events

| Event | Payload | Description |
|-------|---------|-------------|
| `task-click` | `Task` | Task clicked |
| `date-click` | `Date` | Date clicked |
| `date-change` | `Date` | Displayed month changed |

#### Usage Example

```vue
<template>
  <CalendarView
    :tasks="tasks"
    :current-date="currentMonth"
    view="month"
    @task-click="handleTaskClick"
    @date-change="handleMonthChange"
  />
</template>
```

---

### DashboardView

Dashboard with statistics and charts.

**Location:** `/src/components/gantt/views/DashboardView.vue`

#### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `tasks` | `Task[]` | **required** | Array of tasks |
| `refreshInterval` | `number` | `60000` | Auto-refresh interval (ms) |
| `showBurndown` | `boolean` | `true` | Show burndown chart |
| `showEarnedValue` | `boolean` | `true` | Show earned value chart |
| `showMilestones` | `boolean` | `true` | Show milestone tracker |
| `showResources` | `boolean` | `true` | Show resource utilization |

#### Events

| Event | Payload | Description |
|-------|---------|-------------|
| `refresh` | - | User clicked refresh |

#### Usage Example

```vue
<template>
  <DashboardView
    :tasks="tasks"
    :refresh-interval="30000"
    :show-burndown="true"
    :show-earned-value="true"
  />
</template>
```

---

## Panel Components

### CommentsPanel

Panel for viewing and adding task comments.

**Location:** `/src/components/gantt/panels/CommentsPanel.vue`

#### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `taskId` | `string` | **required** | Task ID to show comments for |
| `comments` | `Comment[]` | `[]` | Array of comments |
| `loading` | `boolean` | `false` | Loading state |
| `readonly` | `boolean` | `false` | Prevent adding comments |

#### Events

| Event | Payload | Description |
|-------|---------|-------------|
| `add-comment` | `{ taskId: string, content: string }` | Comment added |
| `delete-comment` | `{ commentId: string }` | Comment deleted |

#### Usage Example

```vue
<template>
  <CommentsPanel
    :task-id="selectedTaskId"
    :comments="comments"
    @add-comment="handleAddComment"
  />
</template>
```

---

### HistoryPanel

Panel showing task change history.

**Location:** `/src/components/gantt/panels/HistoryPanel.vue`

#### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `taskId` | `string` | **required** | Task ID to show history for |
| `history` | `HistoryEntry[]` | `[]` | Array of history entries |
| `loading` | `boolean` | `false` | Loading state |
| `showDiff` | `boolean` | `true` | Show diff viewer |

#### Events

| Event | Payload | Description |
|-------|---------|-------------|
| `restore` | `{ historyId: string }` | Restore to this version |

#### Usage Example

```vue
<template>
  <HistoryPanel
    :task-id="selectedTaskId"
    :history="history"
    :show-diff="true"
    @restore="handleRestore"
  />
</template>
```

---

### SmartSuggestionsPanel

Panel displaying AI-powered suggestions.

**Location:** `/src/components/gantt/panels/SmartSuggestionsPanel.vue`

#### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `suggestions` | `Suggestion[]` | `[]` | Array of AI suggestions |
| `loading` | `boolean` | `false` | Loading state |
| `autoAnalyze` | `boolean` | `true` | Auto-analyze on load |

#### Events

| Event | Payload | Description |
|-------|---------|-------------|
| `accept` | `Suggestion` | Suggestion accepted |
| `dismiss` | `string` | Suggestion dismissed (suggestionId) |
| `refresh` | - | Request new suggestions |

#### Usage Example

```vue
<template>
  <SmartSuggestionsPanel
    :suggestions="suggestions"
    :auto-analyze="true"
    @accept="handleAcceptSuggestion"
    @dismiss="handleDismissSuggestion"
  />
</template>
```

---

## Dialog Components

### BulkEditDialog

Dialog for editing multiple tasks at once.

**Location:** `/src/components/gantt/dialogs/BulkEditDialog.vue`

#### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `visible` | `boolean` | **required** | Dialog visibility |
| `selectedTasks` | `Task[]` | **required** | Tasks to edit |
| `fields` | `string[]` | `['status', 'priority', 'assignee']` | Available fields to edit |

#### Events

| Event | Payload | Description |
|-------|---------|-------------|
| `update:visible` | `boolean` | Visibility changed |
| `confirm` | `{ taskIds: string[], updates: Partial<Task> }` | Bulk edit confirmed |
| `cancel` | - | Dialog cancelled |

#### Usage Example

```vue
<template>
  <BulkEditDialog
    v-model:visible="showBulkEdit"
    :selected-tasks="selectedTasks"
    @confirm="handleBulkEdit"
  />
</template>
```

---

### ResourceLevelingDialog

Dialog for resource leveling optimization.

**Location:** `/src/components/gantt/dialogs/ResourceLevelingDialog.vue`

#### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `visible` | `boolean` | **required** | Dialog visibility |
| `tasks` | `Task[]` | **required** | All tasks |
| `resources` | `Resource[]` | **required** | Available resources |
| `options` | `LevelingOptions` | `{}` | Leveling options |

#### Events

| Event | Payload | Description |
|-------|---------|-------------|
| `update:visible` | `boolean` | Visibility changed |
| `apply` | `{ leveledTasks: Task[], changes: Change[] }` | Apply leveling |
| `cancel` | - | Dialog cancelled |

#### Usage Example

```vue
<template>
  <ResourceLevelingDialog
    v-model:visible="showLeveling"
    :tasks="tasks"
    :resources="resources"
    @apply="handleApplyLeveling"
  />
</template>
```

---

### ReportBuilderDialog

Dialog for building custom reports.

**Location:** `/src/components/gantt/dialogs/ReportBuilderDialog.vue`

#### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `visible` | `boolean` | **required** | Dialog visibility |
| `tasks` | `Task[]` | **required** | Tasks to report on |
| `template` | `ReportTemplate \| null` | `null` | Report template |

#### Events

| Event | Payload | Description |
|-------|---------|-------------|
| `update:visible` | `boolean` | Visibility changed |
| `generate` | `{ config: ReportConfig, format: ExportFormat }` | Generate report |
| `save-template` | `ReportTemplate` | Save report template |

#### Usage Example

```vue
<template>
  <ReportBuilderDialog
    v-model:visible="showReportBuilder"
    :tasks="tasks"
    @generate="handleGenerateReport"
  />
</template>
```

---

## Overlay Components

### ContextMenu

Context menu for quick actions.

**Location:** `/src/components/gantt/overlays/ContextMenu.vue`

#### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `visible` | `boolean` | **required** | Menu visibility |
| `x` | `number` | **required** | X coordinate |
| `y` | `number` | **required** | Y coordinate |
| `items` | `MenuItem[]` | **required** | Menu items |
| `task` | `Task \| null` | `null` | Related task |

#### Events

| Event | Payload | Description |
|-------|---------|-------------|
| `select` | `MenuItem` | Item selected |
| `close` | - | Menu closed |

#### Usage Example

```vue
<template>
  <ContextMenu
    :visible="showContextMenu"
    :x="contextMenuX"
    :y="contextMenuY"
    :items="contextMenuItems"
    :task="contextMenuTask"
    @select="handleContextMenuSelect"
  />
</template>
```

---

### GuidedTour

Interactive guided tour for first-time users.

**Location:** `/src/components/gantt/overlays/GuidedTour.vue`

#### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `visible` | `boolean` | **required** | Tour visibility |
| `steps` | `TourStep[]` | **required** | Tour steps |
| `currentStep` | `number` | `0` | Current step index |

#### Events

| Event | Payload | Description |
|-------|---------|-------------|
| `next` | - | Next step |
| `prev` | - | Previous step |
| `skip` | - | Skip tour |
| `finish` | - | Finish tour |

#### Usage Example

```vue
<template>
  <GuidedTour
    :visible="showTour"
    :steps="tourSteps"
    :current-step="currentStep"
    @finish="handleTourFinish"
  />
</template>
```

---

### Minimap

Miniature overview map of the timeline.

**Location:** `/src/components/gantt/overlays/Minimap.vue`

#### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `tasks` | `Task[]` | **required** | Tasks to display |
| `viewportStart` | `Date` | **required** | Viewport start date |
| `viewportEnd` | `Date` | **required** | Viewport end date |
| `timelineStart` | `Date` | **required** | Timeline start date |
| `timelineEnd` | `Date` | **required** | Timeline end date |
| `width` | `number` | `200` | Minimap width in pixels |
| `height` | `number` | `100` | Minimap height in pixels |

#### Events

| Event | Payload | Description |
|-------|---------|-------------|
| `viewport-change` | `{ start: Date, end: Date }` | Viewport changed via minimap |

#### Usage Example

```vue
<template>
  <Minimap
    :tasks="tasks"
    :viewport-start="viewportStartDate"
    :viewport-end="viewportEndDate"
    :timeline-start="timelineStartDate"
    :timeline-end="timelineEndDate"
    @viewport-change="handleViewportChange"
  />
</template>
```

---

## Mobile Components

### MobileGanttView

Mobile-optimized Gantt view with touch gestures.

**Location:** `/src/components/gantt/mobile/MobileGanttView.vue`

#### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `tasks` | `Task[]` | **required** | Tasks to display |
| `dayWidth` | `number` | `60` | Day width in pixels |
| `enableSwipe` | `boolean` | `true` | Enable swipe gestures |
| `enablePinchZoom` | `boolean` | `true` | Enable pinch-to-zoom |

#### Events

| Event | Payload | Description |
|-------|---------|-------------|
| `task-click` | `Task` | Task clicked |
| `swipe-left` | - | Swiped left |
| `swipe-right` | - | Swiped right |
| `pinch-zoom` | `number` | Zoom level |

#### Usage Example

```vue
<template>
  <MobileGanttView
    :tasks="tasks"
    :day-width="60"
    @task-click="handleTaskClick"
  />
</template>
```

---

## Dashboard Components

### BurndownChart

Burndown chart component using Chart.js.

**Location:** `/src/components/gantt/dashboard/BurndownChart.vue`

#### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `tasks` | `Task[]` | **required** | Tasks for calculation |
| `startDate` | `Date` | **required** | Sprint start date |
| `endDate` | `Date` | **required** | Sprint end date |
| `type` | `'task' \| 'story-point' \| 'hour'` | `'task'` | Burndown type |
| `height` | `number` | `300` | Chart height in pixels |

#### Usage Example

```vue
<template>
  <BurndownChart
    :tasks="sprintTasks"
    :start-date="sprintStart"
    :end-date="sprintEnd"
    type="story-point"
    :height="350"
  />
</template>
```

---

### ResourceUtilization

Resource utilization chart.

**Location:** `/src/components/gantt/dashboard/ResourceUtilization.vue`

#### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `tasks` | `Task[]` | **required** | Tasks for calculation |
| `resources` | `Resource[]` | **required** | Resources |
| `startDate` | `Date` | **required** | Start date |
| `endDate` | `Date` | **required** | End date |
| `height` | `number` | `300` | Chart height |

#### Usage Example

```vue
<template>
  <ResourceUtilization
    :tasks="tasks"
    :resources="resources"
    :start-date="projectStart"
    :end-date="projectEnd"
  />
</template>
```

---

## Data Types Reference

### Task

```typescript
interface Task {
  id: string
  projectId: string
  name: string
  description?: string
  startDate: Date
  endDate: Date
  duration: number // in days
  progress: number // 0-100
  status: 'not_started' | 'in_progress' | 'completed' | 'delayed' | 'blocked'
  priority: 'low' | 'medium' | 'high' | 'critical'
  assignee?: string
  parentId?: string // for subtasks
  position: number // sort order
  color?: string
  milestone: boolean
  constraint?: {
    type: 'start-no-earlier-than' | 'finish-no-later-than' | 'must-start-on' | 'must-finish-on'
    date: Date
  }
  customFields?: Record<string, any>
}
```

### Dependency

```typescript
interface Dependency {
  id: string
  fromTaskId: string
  toTaskId: string
  type: 'finish-to-start' | 'start-to-start' | 'finish-to-finish' | 'start-to-finish'
  lag?: number // days
}
```

### Comment

```typescript
interface Comment {
  id: string
  taskId: string
  userId: string
  userName: string
  content: string
  createdAt: Date
  updatedAt: Date
}
```

### HistoryEntry

```typescript
interface HistoryEntry {
  id: string
  taskId: string
  userId: string
  userName: string
  timestamp: Date
  changes: {
    field: string
    oldValue: any
    newValue: any
  }[]
}
```

### Suggestion

```typescript
interface Suggestion {
  id: string
  type: 'schedule' | 'resource' | 'risk' | 'optimization'
  title: string
  description: string
  priority: 'low' | 'medium' | 'high'
  impact: string
  effort: string
  actions: SuggestionAction[]
}
```

---

## Component Best Practices

### 1. Performance

- **Always use virtual scrolling** for lists with 50+ items
- **Debounce expensive operations** like search and auto-save
- **Use computed properties** for derived data
- **Avoid unnecessary reactivity** for large datasets

### 2. Accessibility

- **Use semantic HTML** elements
- **Provide keyboard shortcuts** for common actions
- **Include ARIA labels** on interactive elements
- **Ensure color contrast** meets WCAG standards

### 3. Error Handling

- **Validate props** in component setup
- **Provide loading states** for async operations
- **Show user-friendly error messages**
- **Log errors to console** in development

### 4. Testing

- **Unit test** component logic in isolation
- **Integration test** component interactions
- **E2E test** critical user flows
- **Performance test** with large datasets

---

**Document Version:** 1.0.0
**Last Updated:** 2026-02-19
