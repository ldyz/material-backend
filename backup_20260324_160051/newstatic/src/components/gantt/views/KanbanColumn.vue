<template>
  <div
    class="kanban-column"
    :class="{
      'is-over': isDraggingOver,
      'is-over-limit': isOverLimit,
      'has-wip-limit': showWipLimit && wipLimit
    }"
    @drop="handleDrop"
    @dragover="handleDragOver"
    @dragenter="handleDragEnter"
    @dragleave="handleDragLeave"
  >
    <!-- Column Header -->
    <div class="kanban-column__header">
      <div class="column-info">
        <div
          class="column-color"
          :style="{ backgroundColor: column.color }"
        />
        <h3 class="column-title">{{ column.title }}</h3>
        <el-badge
          :value="tasks.length"
          :max="99"
          :type="getBadgeType"
          class="column-count"
        />
      </div>

      <div class="column-actions">
        <!-- WIP Limit Indicator -->
        <el-tooltip
          v-if="showWipLimit && wipLimit"
          :content="`Work in Progress Limit: ${tasks.length}/${wipLimit}`"
          placement="top"
        >
          <div
            class="wip-indicator"
            :class="{
              'is-warning': tasks.length >= wipLimit * 0.8,
              'is-danger': tasks.length >= wipLimit
            }"
          >
            <el-icon><Warning /></el-icon>
            <span>{{ tasks.length }}/{{ wipLimit }}</span>
          </div>
        </el-tooltip>

        <el-dropdown trigger="click" @command="handleCommand">
          <el-button :icon="MoreFilled" circle size="small" text />
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="add-task">
                <el-icon><Plus /></el-icon>
                Add Task
              </el-dropdown-item>
              <el-dropdown-item command="sort-by-priority">
                <el-icon><Sort /></el-icon>
                Sort by Priority
              </el-dropdown-item>
              <el-dropdown-item command="sort-by-date">
                <el-icon><Calendar /></el-icon>
                Sort by Date
              </el-dropdown-item>
              <el-dropdown-item command="sort-by-name" divided>
                <el-icon><Sort /></el-icon>
                Sort by Name
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>

    <!-- Task Drop Zone -->
    <div
      class="kanban-column__tasks"
      ref="tasksContainer"
    >
      <transition-group name="task-list" tag="div">
        <KanbanCard
          v-for="task in sortedTasks"
          :key="task.id"
          :task="task"
          :card-view="cardView"
          :show-avatar="showAvatars"
          :column-color="column.color"
          draggable="true"
          @dragstart="handleDragStart(task, $event)"
          @dragend="handleDragEnd"
          @click="handleClick(task)"
        />
      </transition-group>

      <!-- Empty State -->
      <div v-if="tasks.length === 0" class="empty-state">
        <el-empty
          :image-size="60"
          description="No tasks"
        />
        <el-button
          v-if="canAddTasks"
          type="primary"
          size="small"
          @click="handleAddTask"
        >
          <el-icon><Plus /></el-icon>
          Add Task
        </el-button>
      </div>

      <!-- WIP Limit Warning -->
      <el-alert
        v-if="showWipLimit && wipLimit && tasks.length >= wipLimit"
        type="warning"
        :closable="false"
        show-icon
        class="wip-warning"
      >
        WIP limit reached ({{ tasks.length }}/{{ wipLimit }})
      </el-alert>
    </div>

    <!-- Quick Add Task -->
    <div v-if="canAddTasks" class="kanban-column__footer">
      <el-button
        class="add-task-btn"
        @click="handleAddTask"
        :icon="Plus"
        text
      >
        Add Task
      </el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import {
  Warning,
  MoreFilled,
  Plus,
  Sort,
  Calendar
} from '@element-plus/icons-vue'
import KanbanCard from './KanbanCard.vue'

/**
 * KanbanColumn Component
 * Individual column in the Kanban board
 *
 * @props {Object} column - Column definition with id, title, color, wipLimit
 * @props {Array} tasks - Array of tasks in this column
 * @props {Number} wipLimit - Work in progress limit
 * @props {Boolean} showWipLimit - Whether to show WIP limit indicator
 * @props {Boolean} showAvatars - Whether to show assignee avatars
 * @props {String} cardView - Card display mode (compact/detailed)
 * @props {Object} swimlane - Optional swimlane info
 *
 * @emits {Object} task-move - Emitted when a task is moved to this column
 * @emits {Object} task-click - Emitted when a task is clicked
 * @emits {String} command - Emitted when column command is triggered
 */

const props = defineProps({
  column: {
    type: Object,
    required: true
  },
  tasks: {
    type: Array,
    default: () => []
  },
  wipLimit: {
    type: Number,
    default: null
  },
  showWipLimit: {
    type: Boolean,
    default: true
  },
  showAvatars: {
    type: Boolean,
    default: true
  },
  cardView: {
    type: String,
    default: 'detailed',
    validator: (value) => ['compact', 'detailed'].includes(value)
  },
  swimlane: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['task-move', 'task-click', 'task-add', 'sort', 'command'])

// State
const isDraggingOver = ref(false)
const draggedTask = ref(null)
const sortBy = ref(null)

// Computed
const isOverLimit = computed(() => {
  return props.showWipLimit && props.wipLimit && props.tasks.length >= props.wipLimit
})

const canAddTasks = computed(() => {
  return !props.showWipLimit || !props.wipLimit || props.tasks.length < props.wipLimit
})

const getBadgeType = computed(() => {
  if (isOverLimit.value) return 'danger'
  if (props.wipLimit && props.tasks.length >= props.wipLimit * 0.8) return 'warning'
  return 'info'
})

const sortedTasks = computed(() => {
  let sorted = [...props.tasks]

  if (sortBy.value === 'priority') {
    const priorityOrder = { high: 0, medium: 1, low: 2 }
    sorted.sort((a, b) => priorityOrder[a.priority] - priorityOrder[b.priority])
  } else if (sortBy.value === 'date') {
    sorted.sort((a, b) => new Date(a.dueDate) - new Date(b.dueDate))
  } else if (sortBy.value === 'name') {
    sorted.sort((a, b) => a.name.localeCompare(b.name))
  }

  return sorted
})

// Methods
const handleDragStart = (task, event) => {
  draggedTask.value = task
  event.dataTransfer.effectAllowed = 'move'
  event.dataTransfer.setData('text/plain', task.id)

  // Add drag styling
  event.target.classList.add('is-dragging')
}

const handleDragEnd = (event) => {
  event.target.classList.remove('is-dragging')
  isDraggingOver.value = false
  draggedTask.value = null
}

const handleDragOver = (event) => {
  event.preventDefault()
  event.dataTransfer.dropEffect = 'move'
}

const handleDragEnter = (event) => {
  event.preventDefault()
  isDraggingOver.value = true
}

const handleDragLeave = (event) => {
  // Only set false if we're actually leaving the column
  const rect = event.currentTarget.getBoundingClientRect()
  const x = event.clientX
  const y = event.clientY

  if (x < rect.left || x >= rect.right || y < rect.top || y >= rect.bottom) {
    isDraggingOver.value = false
  }
}

const handleDrop = (event) => {
  event.preventDefault()
  isDraggingOver.value = false

  if (!draggedTask.value) return

  const taskId = event.dataTransfer.getData('text/plain')
  const newColumnId = props.column.id

  emit('task-move', {
    taskId: parseInt(taskId),
    fromColumn: draggedTask.value.status,
    toColumn: newColumnId,
    swimlane: props.swimlane
  })

  draggedTask.value = null
}

const handleClick = (task) => {
  emit('task-click', task)
}

const handleAddTask = () => {
  emit('task-add', {
    columnId: props.column.id,
    swimlane: props.swimlane
  })
}

const handleCommand = (command) => {
  if (command === 'sort-by-priority') {
    sortBy.value = sortBy.value === 'priority' ? null : 'priority'
    emit('sort', { columnId: props.column.id, sortBy: sortBy.value })
  } else if (command === 'sort-by-date') {
    sortBy.value = sortBy.value === 'date' ? null : 'date'
    emit('sort', { columnId: props.column.id, sortBy: sortBy.value })
  } else if (command === 'sort-by-name') {
    sortBy.value = sortBy.value === 'name' ? null : 'name'
    emit('sort', { columnId: props.column.id, sortBy: sortBy.value })
  } else {
    emit('command', { command, columnId: props.column.id })
  }
}
</script>

<style scoped lang="scss">
.kanban-column {
  background: #fff;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  min-width: 280px;
  max-width: 400px;
  height: 100%;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: all 0.3s;

  &.is-over {
    background: #ecf5ff;
    box-shadow: 0 0 0 2px #409eff;
  }

  &.is-over-limit {
    border: 2px solid #f56c6c;
  }

  &__header {
    padding: 12px 16px;
    border-bottom: 1px solid #ebeef5;
    background: #fafafa;
    border-radius: 8px 8px 0 0;

    .column-info {
      display: flex;
      align-items: center;
      gap: 8px;
      margin-bottom: 8px;

      .column-color {
        width: 4px;
        height: 20px;
        border-radius: 2px;
      }

      .column-title {
        margin: 0;
        font-size: 14px;
        font-weight: 600;
        color: #303133;
        flex: 1;
      }

      .column-count {
        margin-left: auto;
      }
    }

    .column-actions {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 8px;

      .wip-indicator {
        display: flex;
        align-items: center;
        gap: 4px;
        padding: 4px 8px;
        border-radius: 4px;
        background: #e8f4ff;
        color: #409eff;
        font-size: 12px;
        font-weight: 600;

        &.is-warning {
          background: #fef0e8;
          color: #e6a23c;
        }

        &.is-danger {
          background: #fee;
          color: #f56c6c;
        }
      }
    }
  }

  &__tasks {
    flex: 1;
    padding: 12px;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
    gap: 8px;
    min-height: 100px;

    .empty-state {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 24px;
      gap: 12px;
      min-height: 150px;
    }

    .wip-warning {
      margin-top: auto;
    }
  }

  &__footer {
    padding: 8px 12px;
    border-top: 1px solid #ebeef5;

    .add-task-btn {
      width: 100%;
      color: #909399;

      &:hover {
        color: #409eff;
        background: #ecf5ff;
      }
    }
  }
}

/* Task List Transitions */
.task-list-move,
.task-list-enter-active,
.task-list-leave-active {
  transition: all 0.3s ease;
}

.task-list-enter-from,
.task-list-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}

.task-list-leave-active {
  position: absolute;
  width: 100%;
}

/* Dragging State */
:deep(.is-dragging) {
  opacity: 0.5;
  cursor: move;
}

/* Responsive Design */
@media (max-width: 768px) {
  .kanban-column {
    min-width: 250px;
    max-width: none;
    width: 100%;

    &__header {
      .column-title {
        font-size: 13px;
      }
    }

    &__tasks {
      padding: 8px;
    }
  }
}
</style>
