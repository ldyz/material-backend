<template>
  <div class="mobile-gantt-view" :class="{ 'is-fullscreen': isFullscreen }">
    <!-- 紧凑工具栏 -->
    <GanttToolbar
      :view-mode="viewMode"
      :current-period-text="currentPeriodText"
      :current-zoom-label="currentZoomLabel"
      :show-dependencies="showDependencies"
      :show-critical-path="showCriticalPath"
      :show-baseline="showBaseline"
      :is-fullscreen="isFullscreen"
      :is-saving="saving"
      :has-unsaved-changes="hasUnsavedChanges"
      :compact="true"
      @navigate-date="navigateDate"
      @go-today="goToToday"
      @view-mode-change="handleViewModeChange"
      @zoom-in="zoomIn"
      @zoom-out="zoomOut"
      @toggle-dependencies="toggleDependencies"
      @toggle-critical-path="toggleCriticalPath"
      @toggle-baseline="toggleBaseline"
      @add-task="handleAddTask"
      @save-all="handleSaveAll"
    />

    <!-- 视图切换器 -->
    <div class="mobile-view-switcher">
      <button
        class="view-switcher__btn"
        :class="{ 'is-active: activeView === 'list' }"
        @click="activeView = 'list'"
      >
        <i class="el-icon-s-order"></i>
        <span>任务列表</span>
      </button>
      <button
        class="view-switcher__btn"
        :class="{ 'is-active': activeView === 'timeline' }"
        @click="activeView = 'timeline'"
      >
        <i class="el-icon-date"></i>
        <span>时间轴</span>
      </button>
    </div>

    <!-- 内容区域 -->
    <div class="mobile-gantt-content">
      <!-- 任务列表视图 -->
      <Transition name="fade" mode="out-in">
        <TaskListView
          v-if="activeView === 'list'"
          :tasks="filteredTasks"
          :grouped-tasks="groupedTasks"
          :selected-task-id="selectedTask?.id"
          :group-mode="groupMode"
          :collapsed-groups="collapsedGroups"
          :search-keyword="searchKeyword"
          @row-click="handleRowClick"
          @task-click="handleTaskClick"
          @toggle-group="toggleGroup"
          @context-menu="handleContextMenu"
          @add-task="handleAddTask"
        />
      </Transition>

      <!-- 时间轴视图 -->
      <Transition name="fade" mode="out-in">
        <TimelineSwipeView
          v-if="activeView === 'timeline'"
          :tasks="filteredTasks"
          :timeline-days="timelineDays"
          :timeline-weeks="timelineWeeks"
          :timeline-months="timelineMonths"
          :view-mode="viewMode"
          :day-width="dayWidth"
          :row-height="rowHeight"
          :show-dependencies="showDependencies"
          :show-critical-path="showCriticalPath"
          :show-baseline="showBaseline"
          @task-click="handleTaskClick"
          @task-dblclick="handleTaskDblClick"
        />
      </Transition>
    </div>

    <!-- 浮动操作按钮 -->
    <div class="mobile-fab">
      <button
        class="fab-button"
        @click="handleAddTask"
        aria-label="添加任务"
      >
        <i class="el-icon-plus"></i>
      </button>
    </div>

    <!-- 任务详情抽屉（移动端全屏） -->
    <TaskDetailDrawer
      v-model:visible="taskDetailVisible"
      :task="selectedTask"
      :fullscreen="true"
      @edit="handleEditTaskFromDrawer"
      @duplicate="handleDuplicateTask"
      @delete="handleDeleteTask"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import GanttToolbar from '@/components/progress/GanttToolbar.vue'
import TaskListView from './TaskListView.vue'
import TimelineSwipeView from './TimelineSwipeView.vue'
import TaskDetailDrawer from '@/components/progress/TaskDetailDrawer.vue'

/**
 * 移动端甘特图视图
 * 专为移动设备优化的甘特图界面
 */
interface Props {
  // 从主组件传递的 props
  viewMode: string
  currentPeriodText: string
  currentZoomLabel: string
  showDependencies: boolean
  showCriticalPath: boolean
  showBaseline: boolean
  groupMode: string
  isFullscreen: boolean
  saving: boolean
  hasUnsavedChanges: boolean
  filteredTasks: any[]
  groupedTasks: any
  selectedTask: any
  collapsedGroups: Set<string>
  searchKeyword: string
  timelineDays: any[]
  timelineWeeks: any[]
  timelineMonths: any[]
  dayWidth: number
  rowHeight: number
}

const props = defineProps<Props>()

const emit = defineEmits<{
  navigateDate: [direction: string]
  goToday: []
  viewModeChange: [mode: string]
  zoomIn: []
  zoomOut: []
  toggleDependencies: []
  toggleCriticalPath: []
  toggleBaseline: []
  rowClick: [task: any]
  taskClick: [task: any, event: Event]
  taskDblClick: [task: any]
  toggleGroup: [groupKey: string]
  contextMenu: [task: any, position: { x: number; y: number }]
  addTask: [parentId?: string | number]
  saveAll: []
  editTask: [task: any]
  duplicateTask: [task: any]
  deleteTask: [task: any]
}>()

// 当前激活的视图
const activeView = ref<'list' | 'timeline'>('list')
const taskDetailVisible = ref(false)

// 事件转发
function navigateDate(direction: string) {
  emit('navigateDate', direction)
}

function goToToday() {
  emit('goToday')
}

function handleViewModeChange(mode: string) {
  emit('viewModeChange', mode)
}

function zoomIn() {
  emit('zoomIn')
}

function zoomOut() {
  emit('zoomOut')
}

function toggleDependencies() {
  emit('toggleDependencies')
}

function toggleCriticalPath() {
  emit('toggleCriticalPath')
}

function toggleBaseline() {
  emit('toggleBaseline')
}

function handleRowClick(task: any) {
  emit('rowClick', task)
}

function handleTaskClick(task: any, event: Event) {
  emit('taskClick', task, event)
}

function handleTaskDblClick(task: any) {
  emit('taskDblClick', task)
}

function toggleGroup(groupKey: string) {
  emit('toggleGroup', groupKey)
}

function handleContextMenu(task: any, position: { x: number; y: number }) {
  emit('contextMenu', task, position)
}

function handleAddTask(parentId?: string | number) {
  emit('addTask', parentId)
}

function handleSaveAll() {
  emit('saveAll')
}

function handleEditTaskFromDrawer(task: any) {
  emit('editTask', task)
}

function handleDuplicateTask(task: any) {
  emit('duplicateTask', task)
}

function handleDeleteTask(task: any) {
  emit('deleteTask', task)
}
</script>

<style scoped>
.mobile-gantt-view {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: var(--color-bg-primary);
}

.mobile-view-switcher {
  display: flex;
  gap: 8px;
  padding: 12px 16px;
  background: var(--color-bg-secondary);
  border-bottom: 1px solid var(--color-border-light);
}

.view-switcher__btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 10px 16px;
  border: 1px solid var(--color-border-base);
  border-radius: var(--radius-md);
  background: var(--color-bg-primary);
  color: var(--color-text-regular);
  font-size: var(--font-size-sm);
  transition: all var(--transition-fast);
}

.view-switcher__btn:active {
  transform: scale(0.98);
}

.view-switcher__btn.is-active {
  background: var(--color-primary);
  border-color: var(--color-primary);
  color: #fff;
}

.mobile-gantt-content {
  flex: 1;
  overflow: hidden;
  position: relative;
}

.mobile-fab {
  position: fixed;
  right: 16px;
  bottom: 16px;
  z-index: var(--z-index-fixed);
}

.fab-button {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: var(--color-primary);
  color: #fff;
  border: none;
  box-shadow: var(--shadow-lg);
  font-size: 24px;
  cursor: pointer;
  transition: all var(--transition-base);
}

.fab-button:active {
  transform: scale(0.95);
}

/* 视图切换动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity var(--transition-base);
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 全屏模式 */
.mobile-gantt-view.is-fullscreen {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: var(--z-index-modal);
}

/* 响应式调整 */
@media (min-width: 768px) {
  .mobile-gantt-view {
    display: none;
  }
}
</style>
