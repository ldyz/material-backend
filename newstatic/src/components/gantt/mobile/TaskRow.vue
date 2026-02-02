<template>
  <div
    class="task-row"
    :class="{ 'is-selected': isSelected }"
    @click="$emit('click')"
  >
    <!-- 任务图标 -->
    <div class="task-row__icon">
      <i
        v-if="task.is_milestone"
        class="el-icon-s-flag"
        :style="{ color: '#f39c12' }"
      />
      <i
        v-else-if="task.is_critical"
        class="el-icon-warning"
        :style="{ color: '#f56c6c' }"
      />
      <i v-else class="el-icon-s-order" />
    </div>

    <!-- 任务信息 -->
    <div class="task-row__content">
      <div class="task-row__name">{{ task.name }}</div>
      <div class="task-row__meta">
        <span class="task-row__dates">
          {{ formatDate(task.start) }} - {{ formatDate(task.end) }}
        </span>
        <span class="task-row__duration">{{ task.duration }}天</span>
      </div>
    </div>

    <!-- 状态标签 -->
    <div class="task-row__status">
      <span
        class="status-badge"
        :class="`status-badge--${task.status}`"
      >
        {{ statusText }}
      </span>
    </div>

    <!-- 进度 -->
    <div class="task-row__progress">
      <div class="progress-bar">
        <div
          class="progress-bar__fill"
          :style="{ width: task.progress + '%' }"
        ></div>
      </div>
      <span class="progress-text">{{ task.progress }}%</span>
    </div>

    <!-- 右箭头 -->
    <div class="task-row__arrow">
      <i class="el-icon-arrow-right" />
    </div>
  </div>
</template>

<script setup lang="ts">
/**
 * 移动端任务行组件
 */
interface Task {
  id: string | number
  name: string
  start: string
  end: string
  duration: number
  progress: number
  status: string
  is_milestone: boolean
  is_critical: boolean
}

interface Props {
  task: Task
  isSelected?: boolean
}

defineProps<Props>()

defineEmits<{
  click: []
}>()

function formatDate(dateStr: string): string {
  const date = new Date(dateStr)
  return `${date.getMonth() + 1}/${date.getDate()}`
}

const statusMap: Record<string, string> = {
  not_started: '未开始',
  in_progress: '进行中',
  completed: '已完成',
  delayed: '延期'
}

function statusText(status: string): string {
  return statusMap[status] || status
}
</script>

<style scoped>
.task-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: var(--color-bg-primary);
  border-bottom: 1px solid var(--color-border-lighter);
  cursor: pointer;
  transition: background var(--transition-fast);
}

.task-row:active {
  background: var(--color-bg-secondary);
}

.task-row.is-selected {
  background: var(--color-primary-lighter);
  border-left: 3px solid var(--color-primary);
}

.task-row__icon {
  font-size: 20px;
  flex-shrink: 0;
}

.task-row__content {
  flex: 1;
  min-width: 0;
}

.task-row__name {
  font-size: var(--font-size-md);
  font-weight: var(--font-weight-medium);
  color: var(--color-text-primary);
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.task-row__meta {
  display: flex;
  gap: 8px;
  font-size: var(--font-size-xs);
  color: var(--color-text-secondary);
}

.task-row__dates {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.task-row__duration {
  flex-shrink: 0;
}

.task-row__status {
  flex-shrink: 0;
}

.status-badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: var(--radius-full);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-medium);
}

.status-badge--not_started {
  background: var(--color-info-lighter);
  color: var(--color-info);
}

.status-badge--in_progress {
  background: var(--color-primary-lighter);
  color: var(--color-primary);
}

.status-badge--completed {
  background: var(--color-success-lighter);
  color: var(--color-success);
}

.status-badge--delayed {
  background: var(--color-danger-lighter);
  color: var(--color-danger);
}

.task-row__progress {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  width: 60px;
  flex-shrink: 0;
}

.progress-bar {
  width: 100%;
  height: 4px;
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-full);
  overflow: hidden;
}

.progress-bar__fill {
  height: 100%;
  background: var(--color-primary);
  transition: width var(--transition-base);
}

.progress-text {
  font-size: var(--font-size-xs);
  color: var(--color-text-secondary);
}

.task-row__arrow {
  flex-shrink: 0;
  color: var(--color-text-placeholder);
}
</style>
