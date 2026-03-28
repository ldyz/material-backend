<template>
  <div
    class="kanban-card"
    :class="{
      'is-compact': cardView === 'compact',
      'is-detailed': cardView === 'detailed',
      'is-overdue': isOverdue,
      'is-due-soon': isDueSoon
    }"
    @click="handleClick"
  >
    <!-- Card Header -->
    <div class="kanban-card__header">
      <div class="card-id">#{{ task.id }}</div>
      <div class="card-priority">
        <el-tag
          :type="priorityType"
          size="small"
          effect="plain"
        >
          {{ priorityLabel }}
        </el-tag>
      </div>
    </div>

    <!-- Card Title -->
    <div class="kanban-card__title">
      {{ task.name }}
    </div>

    <!-- Detailed View Content -->
    <template v-if="cardView === 'detailed'">
      <!-- Assignee -->
      <div v-if="showAvatar && task.assignee" class="kanban-card__assignee">
        <el-avatar :size="24" :src="getAvatarUrl(task.assignee)">
          {{ task.assignee.charAt(0) }}
        </el-avatar>
        <span class="assignee-name">{{ task.assignee }}</span>
      </div>

      <!-- Due Date -->
      <div v-if="task.dueDate" class="kanban-card__due-date">
        <el-icon><Calendar /></el-icon>
        <span :class="{ 'text-danger': isOverdue, 'text-warning': isDueSoon }">
          {{ formattedDueDate }}
        </span>
        <el-tag
          v-if="isOverdue"
          type="danger"
          size="small"
          effect="plain"
          class="due-badge"
        >
          Overdue
        </el-tag>
      </div>

      <!-- Progress Bar -->
      <div v-if="task.progress !== undefined" class="kanban-card__progress">
        <el-progress
          :percentage="task.progress"
          :status="progressStatus"
          :stroke-width="6"
          :show-text="false"
        />
        <span class="progress-text">{{ task.progress }}%</span>
      </div>

      <!-- Dependencies -->
      <div v-if="task.dependencies?.length" class="kanban-card__dependencies">
        <el-tooltip
          :content="`${task.dependencies.length} dependencies`"
          placement="top"
        >
          <el-tag size="small" type="info" effect="plain">
            <el-icon><Link /></el-icon>
            {{ task.dependencies.length }}
          </el-tag>
        </el-tooltip>
      </div>

      <!-- Milestone Marker -->
      <div v-if="task.isMilestone" class="kanban-card__milestone">
        <el-icon color="#f59e0b"><StarFilled /></el-icon>
        <span>Milestone</span>
      </div>

      <!-- Description Preview -->
      <div
        v-if="task.description"
        class="kanban-card__description"
      >
        {{ truncatedDescription }}
      </div>

      <!-- Task Labels/Tags -->
      <div v-if="task.labels?.length" class="kanban-card__labels">
        <el-tag
          v-for="label in task.labels.slice(0, 3)"
          :key="label"
          size="small"
          effect="plain"
          class="label-tag"
        >
          {{ label }}
        </el-tag>
        <el-tag
          v-if="task.labels.length > 3"
          size="small"
          type="info"
          effect="plain"
        >
          +{{ task.labels.length - 3 }}
        </el-tag>
      </div>
    </template>

    <!-- Compact View Content -->
    <template v-else>
      <div class="kanban-card__compact">
        <div class="compact-info">
          <!-- Assignee Avatar -->
          <el-avatar
            v-if="showAvatar && task.assignee"
            :size="20"
            :src="getAvatarUrl(task.assignee)"
          >
            {{ task.assignee.charAt(0) }}
          </el-avatar>

          <!-- Due Date -->
          <span v-if="task.dueDate" class="compact-due" :class="{ 'text-danger': isOverdue }">
            {{ formattedDueDateShort }}
          </span>
        </div>

        <!-- Progress Indicator -->
        <div v-if="task.progress !== undefined" class="compact-progress">
          <div
            class="progress-bar"
            :style="{ width: `${task.progress}%`, backgroundColor: columnColor }"
          />
        </div>
      </div>
    </template>

    <!-- Card Actions -->
    <div class="kanban-card__actions">
      <el-button-group>
        <el-button
          :icon="Edit"
          size="small"
          text
          @click.stop="handleEdit"
        />
        <el-button
          :icon="Delete"
          size="small"
          text
          type="danger"
          @click.stop="handleDelete"
        />
      </el-button-group>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import {
  Calendar,
  Edit,
  Delete,
  Link,
  StarFilled
} from '@element-plus/icons-vue'
import { format, parseISO, isPast, isWithinInterval, subDays } from 'date-fns'

/**
 * KanbanCard Component
 * Individual task card in the Kanban board
 *
 * @props {Object} task - Task object with all task details
 * @props {String} cardView - Display mode (compact/detailed)
 * @props {Boolean} showAvatar - Whether to show assignee avatar
 * @props {String} columnColor - Color of the parent column
 *
 * @emits {Object} click - Emitted when card is clicked
 * @emits {Object} edit - Emitted when edit button is clicked
 * @emits {Object} delete - Emitted when delete button is clicked
 */

const props = defineProps({
  task: {
    type: Object,
    required: true
  },
  cardView: {
    type: String,
    default: 'detailed',
    validator: (value) => ['compact', 'detailed'].includes(value)
  },
  showAvatar: {
    type: Boolean,
    default: true
  },
  columnColor: {
    type: String,
    default: '#409EFF'
  }
})

const emit = defineEmits(['click', 'edit', 'delete'])

// Computed
const priorityType = computed(() => {
  const types = {
    high: 'danger',
    medium: 'warning',
    low: 'info'
  }
  return types[props.task.priority] || 'info'
})

const priorityLabel = computed(() => {
  const labels = {
    high: 'High',
    medium: 'Med',
    low: 'Low'
  }
  return labels[props.task.priority] || 'Med'
})

const isOverdue = computed(() => {
  if (!props.task.dueDate) return false
  if (props.task.status === 'completed') return false
  return isPast(parseISO(props.task.dueDate))
})

const isDueSoon = computed(() => {
  if (!props.task.dueDate) return false
  if (props.task.status === 'completed') return false
  if (isOverdue.value) return false
  const dueDate = parseISO(props.task.dueDate)
  const now = new Date()
  const threeDaysFromNow = subDays(now, -3)
  return isWithinInterval(dueDate, { start: now, end: threeDaysFromNow })
})

const formattedDueDate = computed(() => {
  if (!props.task.dueDate) return ''
  return format(parseISO(props.task.dueDate), 'MMM d, yyyy')
})

const formattedDueDateShort = computed(() => {
  if (!props.task.dueDate) return ''
  return format(parseISO(props.task.dueDate), 'M/d')
})

const progressStatus = computed(() => {
  if (props.task.progress >= 100) return 'success'
  if (props.task.progress >= 50) return undefined
  return 'exception'
})

const truncatedDescription = computed(() => {
  if (!props.task.description) return ''
  return props.task.description.length > 80
    ? props.task.description.substring(0, 80) + '...'
    : props.task.description
})

// Methods
const getAvatarUrl = (assignee) => {
  // Generate avatar URL based on assignee name
  return `https://ui-avatars.com/api/?name=${encodeURIComponent(assignee)}&background=random`
}

const handleClick = () => {
  emit('click', props.task)
}

const handleEdit = () => {
  emit('edit', props.task)
}

const handleDelete = () => {
  emit('delete', props.task)
}
</script>

<style scoped lang="scss">
.kanban-card {
  background: #fff;
  border: 1px solid #dcdfe6;
  border-radius: 6px;
  padding: 12px;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;

  &:hover {
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
    border-color: #409eff;
    transform: translateY(-2px);
  }

  &.is-overdue {
    border-left: 3px solid #f56c6c;
  }

  &.is-due-soon {
    border-left: 3px solid #e6a23c;
  }

  &__header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;

    .card-id {
      font-size: 11px;
      color: #909399;
      font-weight: 500;
    }

    .card-priority {
      :deep(.el-tag) {
        height: 20px;
        line-height: 18px;
        padding: 0 6px;
        font-size: 11px;
      }
    }
  }

  &__title {
    font-size: 14px;
    font-weight: 600;
    color: #303133;
    margin-bottom: 8px;
    line-height: 1.4;
  }

  &__assignee {
    display: flex;
    align-items: center;
    gap: 6px;
    margin-bottom: 8px;

    .assignee-name {
      font-size: 12px;
      color: #606266;
    }
  }

  &__due-date {
    display: flex;
    align-items: center;
    gap: 4px;
    margin-bottom: 8px;
    font-size: 12px;
    color: #606266;

    .text-danger {
      color: #f56c6c;
      font-weight: 600;
    }

    .text-warning {
      color: #e6a23c;
      font-weight: 600;
    }

    .due-badge {
      margin-left: auto;
      :deep(.el-tag__content) {
        display: flex;
        align-items: center;
        gap: 2px;
      }
    }
  }

  &__progress {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 8px;

    .el-progress {
      flex: 1;
    }

    .progress-text {
      font-size: 12px;
      color: #606266;
      font-weight: 600;
      min-width: 35px;
      text-align: right;
    }
  }

  &__dependencies {
    margin-bottom: 8px;

    :deep(.el-tag) {
      display: inline-flex;
      align-items: center;
      gap: 2px;
    }
  }

  &__milestone {
    display: flex;
    align-items: center;
    gap: 4px;
    padding: 4px 8px;
    background: #fef0e8;
    border-radius: 4px;
    font-size: 12px;
    color: #e6a23c;
    margin-bottom: 8px;
    font-weight: 600;
  }

  &__description {
    font-size: 12px;
    color: #909399;
    line-height: 1.5;
    margin-bottom: 8px;
  }

  &__labels {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
    margin-bottom: 8px;

    .label-tag {
      :deep(.el-tag__content) {
        font-size: 11px;
      }
    }
  }

  &__compact {
    display: flex;
    flex-direction: column;
    gap: 8px;

    .compact-info {
      display: flex;
      align-items: center;
      gap: 8px;

      .compact-due {
        font-size: 11px;
        color: #606266;
        margin-left: auto;

        &.text-danger {
          color: #f56c6c;
          font-weight: 600;
        }
      }
    }

    .compact-progress {
      height: 3px;
      background: #ebeef5;
      border-radius: 2px;
      overflow: hidden;

      .progress-bar {
        height: 100%;
        transition: width 0.3s;
      }
    }
  }

  &__actions {
    position: absolute;
    top: 8px;
    right: 8px;
    opacity: 0;
    transition: opacity 0.2s;

    :deep(.el-button-group) {
      display: flex;
      background: rgba(255, 255, 255, 0.9);
      border-radius: 4px;
      padding: 2px;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    }
  }

  &:hover &__actions {
    opacity: 1;
  }

  &.is-compact {
    padding: 8px;

    &__header {
      margin-bottom: 4px;
    }

    &__title {
      font-size: 13px;
      margin-bottom: 6px;
    }
  }
}

/* Responsive Design */
@media (max-width: 768px) {
  .kanban-card {
    &__actions {
      opacity: 1;
    }
  }
}
</style>
