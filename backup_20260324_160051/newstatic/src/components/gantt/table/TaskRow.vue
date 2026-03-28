<template>
  <div
    class="task-row"
    :class="rowClass"
    :style="{ height: rowHeight + 'px' }"
    @click="handleClick"
    @contextmenu="handleContextMenu"
  >
    <!-- 缩进 -->
    <div
      class="task-row__indent"
      :style="{ width: (depth * 20) + 'px' }"
    ></div>

    <!-- 展开/收起图标 -->
    <div
      v-if="hasChildren"
      class="task-row__toggle"
      @click.stop="handleToggle"
    >
      <el-icon :class="{ 'is-collapsed': isCollapsed }">
        <ArrowDown />
      </el-icon>
    </div>
    <div v-else class="task-row__toggle-placeholder"></div>

    <!-- 里程碑图标 -->
    <div class="task-row__icon">
      <el-icon v-if="isMilestone" class="milestone-icon">
        <Star />
      </el-icon>
    </div>

    <!-- 任务名称 -->
    <div class="task-row__name" :title="task.name">
      {{ task.name }}
    </div>

    <!-- 迷你进度条 -->
    <div class="task-row__progress-mini">
      <div
        class="progress-bar"
        :style="{ width: task.progress + '%' }"
      ></div>
    </div>

    <!-- 工期 -->
    <div class="task-row__duration">
      {{ task.duration }} 天
    </div>

    <!-- 起止时间 -->
    <div class="task-row__dates">
      {{ formatDateShort(task.start) }} → {{ formatDateShort(task.end) }}
    </div>

    <!-- 状态（仅平板显示） -->
    <div v-if="isTablet" class="task-row__status">
      <span
        class="status-dot"
        :class="`status-dot--${task.status}`"
      ></span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { ArrowDown, Star } from '@element-plus/icons-vue'
import type { GanttTask } from '@/types/gantt'

interface Props {
  task: GanttTask
  rowHeight: number
  isSelected: boolean
  isCritical: boolean
  isMilestone: boolean
  depth: number
  hasChildren: boolean
  isCollapsed: boolean
  isTablet: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
  click: []
  contextMenu: [event: MouseEvent]
  toggleCollapse: []
  startEdit: [field: string, event: Event]
}>()

/**
 * 行类名
 */
const rowClass = computed(() => {
  return {
    'is-selected': props.isSelected,
    'is-critical': props.isCritical,
    'is-milestone': props.isMilestone
  }
})

/**
 * 格式化短日期
 */
function formatDateShort(dateStr: string): string {
  const date = new Date(dateStr)
  return `${date.getMonth() + 1}/${date.getDate()}`
}

/**
 * 处理点击
 */
function handleClick() {
  emit('click')
}

/**
 * 处理右键菜单
 */
function handleContextMenu(event: MouseEvent) {
  emit('contextMenu', event)
}

/**
 * 处理切换折叠
 */
function handleToggle() {
  emit('toggleCollapse')
}
</script>

<style scoped>
.task-row {
  display: flex;
  align-items: center;
  padding: 0 12px;
  border-bottom: 1px solid var(--color-border-lighter);
  cursor: pointer;
  transition: background-color var(--transition-fast);
  user-select: none;
}

.task-row:hover {
  background-color: var(--row-hover-bg);
}

.task-row.is-selected {
  background-color: var(--row-selected-bg);
  border-left: 3px solid var(--color-primary);
  padding-left: 9px;
}

.task-row.is-critical {
  background-color: rgba(245, 108, 108, 0.05);
}

.task-row.is-milestone {
  font-weight: var(--font-weight-medium);
}

.task-row__indent {
  flex-shrink: 0;
}

.task-row__toggle {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  cursor: pointer;
  transition: transform var(--transition-fast);
}

.task-row__toggle.is-collapsed {
  transform: rotate(-90deg);
}

.task-row__toggle-placeholder {
  width: 20px;
  flex-shrink: 0;
}

.task-row__icon {
  width: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  margin-right: 4px;
}

.milestone-icon {
  color: var(--color-warning);
  font-size: 16px;
}

.task-row__name {
  flex: 1;
  min-width: 150px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: var(--color-text-primary);
}

.task-row__progress-mini {
  width: 60px;
  height: 4px;
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-full);
  overflow: hidden;
  flex-shrink: 0;
  margin-right: 12px;
}

.progress-bar {
  height: 100%;
  background: var(--color-primary);
  transition: width var(--transition-base);
}

.task-row__duration {
  width: 60px;
  text-align: center;
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  flex-shrink: 0;
}

.task-row__dates {
  width: 120px;
  font-size: var(--font-size-xs);
  color: var(--color-text-secondary);
  flex-shrink: 0;
}

.task-row__status {
  width: 40px;
  display: flex;
  justify-content: center;
  flex-shrink: 0;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.status-dot--not_started {
  background: var(--color-info);
}

.status-dot--in_progress {
  background: var(--color-primary);
}

.status-dot--completed {
  background: var(--color-success);
}

.status-dot--delayed {
  background: var(--color-danger);
}

/* 响应式 */
@media (max-width: 1023px) {
  .task-row__dates {
    display: none;
  }
}

@media (max-width: 768px) {
  .task-row__duration {
    display: none;
  }
}
</style>
