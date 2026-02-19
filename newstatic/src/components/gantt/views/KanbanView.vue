<template>
  <div class="kanban-view">
    <!-- Kanban Toolbar -->
    <div class="kanban-view__toolbar">
      <div class="toolbar-left">
        <el-input
          v-model="searchKeyword"
          placeholder="Search tasks..."
          :prefix-icon="Search"
          clearable
          style="width: 250px"
          @input="handleSearch"
        />

        <el-select
          v-model="filterAssignee"
          placeholder="Filter by assignee"
          clearable
          style="width: 180px"
          @change="handleFilterChange"
        >
          <el-option
            v-for="assignee in uniqueAssignees"
            :key="assignee"
            :label="assignee"
            :value="assignee"
          />
        </el-select>

        <el-select
          v-model="filterPriority"
          placeholder="Filter by priority"
          clearable
          style="width: 150px"
          @change="handleFilterChange"
        >
          <el-option label="High" value="high" />
          <el-option label="Medium" value="medium" />
          <el-option label="Low" value="low" />
        </el-select>

        <el-button @click="clearFilters">Clear Filters</el-button>
      </div>

      <div class="toolbar-right">
        <el-button-group>
          <el-button
            :type="cardView === 'compact' ? 'primary' : ''"
            @click="cardView = 'compact'"
          >
            <el-icon><List /></el-icon>
            Compact
          </el-button>
          <el-button
            :type="cardView === 'detailed' ? 'primary' : ''"
            @click="cardView = 'detailed'"
          >
            <el-icon><Grid /></el-icon>
            Detailed
          </el-button>
        </el-button-group>

        <el-dropdown @command="handleViewCommand">
          <el-button :icon="Setting">
            Options<el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="show-wip-limits">
                <el-checkbox v-model="showWipLimits">Show WIP Limits</el-checkbox>
              </el-dropdown-item>
              <el-dropdown-item command="enable-swimlanes">
                <el-checkbox v-model="enableSwimlanes">Enable Swimlanes</el-checkbox>
              </el-dropdown-item>
              <el-dropdown-item command="group-by-assignee" divided>
                <el-checkbox v-model="groupByAssignee">Group by Assignee</el-checkbox>
              </el-dropdown-item>
              <el-dropdown-item command="show-avatars">
                <el-checkbox v-model="showAvatars">Show Avatars</el-checkbox>
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>

        <el-button
          :icon="Refresh"
          @click="refreshData"
          :loading="loading"
        >
          Refresh
        </el-button>
      </div>
    </div>

    <!-- Kanban Board -->
    <div
      class="kanban-view__board"
      v-loading="loading"
      element-loading-text="Loading board..."
    >
      <el-scrollbar>
        <div class="kanban-columns" :class="{ 'has-swimlanes': enableSwimlanes }">
          <!-- Regular Columns -->
          <template v-if="!enableSwimlanes">
            <KanbanColumn
              v-for="column in filteredColumns"
              :key="column.id"
              :column="column"
              :card-view="cardView"
              :show-wip-limit="showWipLimits"
              :show-avatars="showAvatars"
              :tasks="getTasksForColumn(column.id)"
              :wip-limit="column.wipLimit"
              @task-click="handleTaskClick"
              @task-move="handleTaskMove"
              @task-delete="handleTaskDelete"
            />
          </template>

          <!-- Swimlane View -->
          <template v-else>
            <div
              v-for="swimlane in swimlanes"
              :key="swimlane.id"
              class="swimlane"
            >
              <div class="swimlane-header">
                <h4>{{ swimlane.name }}</h4>
                <el-tag size="small">{{ swimlane.taskCount }} tasks</el-tag>
              </div>

              <div class="swimlane-columns">
                <KanbanColumn
                  v-for="column in filteredColumns"
                  :key="`${swimlane.id}-${column.id}`"
                  :column="column"
                  :card-view="cardView"
                  :show-wip-limit="showWipLimits"
                  :show-avatars="showAvatars"
                  :tasks="getTasksForSwimlaneColumn(swimlane.id, column.id)"
                  :wip-limit="column.wipLimit"
                  :swimlane="swimlane"
                  @task-click="handleTaskClick"
                  @task-move="handleTaskMove"
                  @task-delete="handleTaskDelete"
                />
              </div>
            </div>
          </template>
        </div>
      </el-scrollbar>
    </div>

    <!-- Add Task Dialog -->
    <el-dialog
      v-model="showAddTaskDialog"
      title="Add New Task"
      width="600px"
    >
      <el-form :model="newTask" label-width="100px">
        <el-form-item label="Task Name">
          <el-input v-model="newTask.name" placeholder="Enter task name" />
        </el-form-item>

        <el-form-item label="Status">
          <el-select v-model="newTask.status" style="width: 100%">
            <el-option
              v-for="column in columns"
              :key="column.id"
              :label="column.title"
              :value="column.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="Assignee">
          <el-select v-model="newTask.assignee" style="width: 100%">
            <el-option
              v-for="assignee in uniqueAssignees"
              :key="assignee"
              :label="assignee"
              :value="assignee"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="Priority">
          <el-select v-model="newTask.priority" style="width: 100%">
            <el-option label="High" value="high" />
            <el-option label="Medium" value="medium" />
            <el-option label="Low" value="low" />
          </el-select>
        </el-form-item>

        <el-form-item label="Due Date">
          <el-date-picker
            v-model="newTask.dueDate"
            type="date"
            placeholder="Select due date"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="showAddTaskDialog = false">Cancel</el-button>
        <el-button type="primary" @click="addTask">Add Task</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import {
  Search,
  List,
  Grid,
  Setting,
  ArrowDown,
  Refresh
} from '@element-plus/icons-vue'
import KanbanColumn from './KanbanColumn.vue'

/**
 * KanbanView Component
 * Displays Gantt tasks in Kanban board format with drag-and-drop
 *
 * @props {Array} tasks - Array of task objects
 * @props {Array} columns - Array of column definitions
 * @props {Object} wipLimits - Work in progress limits per column
 *
 * @emits {Object} task-click - Emitted when a task is clicked
 * @emits {Object} task-move - Emitted when a task is moved to another column
 * @emits {Object} task-delete - Emitted when a task is deleted
 * @emits {Object} task-add - Emitted when a new task is added
 */

const props = defineProps({
  tasks: {
    type: Array,
    default: () => []
  },
  columns: {
    type: Array,
    default: () => [
      { id: 'not_started', title: 'Not Started', color: '#909399', wipLimit: null },
      { id: 'in_progress', title: 'In Progress', color: '#409EFF', wipLimit: 5 },
      { id: 'review', title: 'Review', color: '#E6A23C', wipLimit: 3 },
      { id: 'completed', title: 'Completed', color: '#67C23A', wipLimit: null }
    ]
  },
  wipLimits: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['task-click', 'task-move', 'task-delete', 'task-add', 'refresh'])

// State
const loading = ref(false)
const searchKeyword = ref('')
const filterAssignee = ref('')
const filterPriority = ref('')
const cardView = ref('detailed')
const showWipLimits = ref(true)
const enableSwimlanes = ref(false)
const groupByAssignee = ref(false)
const showAvatars = ref(true)
const showAddTaskDialog = ref(false)
const newTask = ref({
  name: '',
  status: 'not_started',
  assignee: '',
  priority: 'medium',
  dueDate: null
})

// Computed
const columns = computed(() => {
  return props.columns.map(col => ({
    ...col,
    wipLimit: props.wipLimits[col.id] || col.wipLimit
  }))
})

const uniqueAssignees = computed(() => {
  const assignees = new Set(props.tasks.map(task => task.assignee).filter(Boolean))
  return Array.from(assignees).sort()
})

const filteredTasks = computed(() => {
  return props.tasks.filter(task => {
    const matchesSearch = !searchKeyword.value ||
      task.name.toLowerCase().includes(searchKeyword.value.toLowerCase()) ||
      task.id?.toString().includes(searchKeyword.value)

    const matchesAssignee = !filterAssignee.value || task.assignee === filterAssignee.value
    const matchesPriority = !filterPriority.value || task.priority === filterPriority.value

    return matchesSearch && matchesAssignee && matchesPriority
  })
})

const filteredColumns = computed(() => {
  return columns.value.map(column => ({
    ...column,
    taskCount: getTasksForColumn(column.id).length
  }))
})

const swimlanes = computed(() => {
  if (!enableSwimlanes.value || !groupByAssignee.value) {
    return []
  }

  const lanes = {}
  filteredTasks.value.forEach(task => {
    const key = task.assignee || 'Unassigned'
    if (!lanes[key]) {
      lanes[key] = {
        id: key,
        name: key,
        taskCount: 0
      }
    }
    lanes[key].taskCount++
  })

  return Object.values(lanes)
})

// Methods
const getTasksForColumn = (columnId) => {
  return filteredTasks.value.filter(task => {
    // Map task status to column ID
    const statusMap = {
      not_started: 'not_started',
      in_progress: 'in_progress',
      review: 'review',
      completed: 'completed'
    }
    return statusMap[task.status] === columnId
  })
}

const getTasksForSwimlaneColumn = (swimlaneId, columnId) => {
  return filteredTasks.value.filter(task => {
    const statusMap = {
      not_started: 'not_started',
      in_progress: 'in_progress',
      review: 'review',
      completed: 'completed'
    }
    const taskAssignee = task.assignee || 'Unassigned'
    return statusMap[task.status] === columnId && taskAssignee === swimlaneId
  })
}

const handleSearch = () => {
  // Triggered by search input
}

const handleFilterChange = () => {
  // Triggered by filter changes
}

const clearFilters = () => {
  searchKeyword.value = ''
  filterAssignee.value = ''
  filterPriority.value = ''
}

const handleTaskClick = (task) => {
  emit('task-click', task)
}

const handleTaskMove = (data) => {
  emit('task-move', data)
}

const handleTaskDelete = (task) => {
  emit('task-delete', task)
}

const handleViewCommand = (command) => {
  // Handle dropdown commands if needed
}

const addTask = () => {
  emit('task-add', newTask.value)
  showAddTaskDialog.value = false
  newTask.value = {
    name: '',
    status: 'not_started',
    assignee: '',
    priority: 'medium',
    dueDate: null
  }
}

const refreshData = async () => {
  loading.value = true
  try {
    emit('refresh')
    await new Promise(resolve => setTimeout(resolve, 500))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  // Initialize board
})
</script>

<style scoped lang="scss">
.kanban-view {
  height: 100%;
  display: flex;
  flex-direction: column;

  &__toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 16px;
    background: #fff;
    border-bottom: 1px solid #dcdfe6;
    flex-wrap: wrap;
    gap: 12px;

    .toolbar-left,
    .toolbar-right {
      display: flex;
      align-items: center;
      gap: 12px;
      flex-wrap: wrap;
    }
  }

  &__board {
    flex: 1;
    overflow: hidden;
    background: #f5f7fa;
  }
}

.kanban-columns {
  display: flex;
  gap: 16px;
  padding: 16px;
  height: 100%;
  min-width: min-content;

  &.has-swimlanes {
    flex-direction: column;
    gap: 24px;
  }
}

.swimlane {
  background: #fff;
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

  &-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    padding-bottom: 12px;
    border-bottom: 2px solid #ebeef5;

    h4 {
      margin: 0;
      font-size: 16px;
      font-weight: 600;
      color: #303133;
    }
  }

  &-columns {
    display: flex;
    gap: 16px;
    overflow-x: auto;
  }
}

/* Responsive Design */
@media (max-width: 1200px) {
  .kanban-columns {
    overflow-x: auto;
  }
}

@media (max-width: 768px) {
  .kanban-view__toolbar {
    flex-direction: column;

    .toolbar-left,
    .toolbar-right {
      width: 100%;
      justify-content: center;
    }
  }

  .kanban-columns {
    flex-direction: column;
    overflow-x: hidden;

    &.has-swimlanes {
      gap: 16px;
    }
  }

  .swimlane-columns {
    flex-direction: column;
  }
}
</style>
