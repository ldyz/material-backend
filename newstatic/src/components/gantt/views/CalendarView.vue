<template>
  <div class="calendar-view">
    <!-- Calendar Toolbar -->
    <div class="calendar-view__toolbar">
      <div class="toolbar-left">
        <el-button-group>
          <el-button
            :type="viewMode === 'month' ? 'primary' : ''"
            @click="setViewMode('month')"
          >
            Month
          </el-button>
          <el-button
            :type="viewMode === 'week' ? 'primary' : ''"
            @click="setViewMode('week')"
          >
            Week
          </el-button>
          <el-button
            :type="viewMode === 'day' ? 'primary' : ''"
            @click="setViewMode('day')"
          >
            Day
          </el-button>
        </el-button-group>

        <el-divider direction="vertical" />

        <el-button-group>
          <el-button @click="goToToday">Today</el-button>
          <el-button @click="navigate(-1)">
            <el-icon><ArrowLeft /></el-icon>
          </el-button>
          <el-button @click="navigate(1)">
            <el-icon><ArrowRight /></el-icon>
          </el-button>
        </el-button-group>

        <div class="current-date">
          {{ currentDateLabel }}
        </div>
      </div>

      <div class="toolbar-right">
        <el-button
          :icon="Refresh"
          @click="refreshData"
          :loading="loading"
        >
          Refresh
        </el-button>

        <el-dropdown @command="handleViewCommand">
          <el-button :icon="Setting">
            View Options<el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="show-weekends">
                <el-checkbox v-model="showWeekends">Show Weekends</el-checkbox>
              </el-dropdown-item>
              <el-dropdown-item command="show-milestones">
                <el-checkbox v-model="showMilestones">Show Milestones</el-checkbox>
              </el-dropdown-item>
              <el-dropdown-item command="show-deadlines">
                <el-checkbox v-model="showDeadlines">Show Deadlines</el-checkbox>
              </el-dropdown-item>
              <el-dropdown-item command="show-dependencies" divided>
                <el-checkbox v-model="showDependencies">Show Dependencies</el-checkbox>
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>

    <!-- Calendar Grid -->
    <div
      class="calendar-view__container"
      v-loading="loading"
      element-loading-text="Loading calendar..."
    >
      <el-scrollbar height="calc(100vh - 250px)">
        <div class="calendar-grid" :class="`calendar-grid--${viewMode}`">
          <!-- Month View -->
          <div v-if="viewMode === 'month'" class="month-view">
            <div class="weekdays-header">
              <div
                v-for="day in weekdays"
                :key="day"
                class="weekday-cell"
                :class="{ 'is-weekend': [0, 6].includes(day.index) }"
              >
                {{ day.label }}
              </div>
            </div>

            <div class="month-grid">
              <div
                v-for="(day, index) in monthDays"
                :key="index"
                class="day-cell"
                :class="{
                  'is-other-month': !isSameMonth(day, currentDate),
                  'is-today': isToday(day),
                  'is-weekend': !showWeekends && isWeekend(day)
                }"
                @drop="handleDrop(day, $event)"
                @dragover.prevent
                @dragenter.prevent
              >
                <div class="day-number">
                  {{ formatDate(day, 'd') }}
                </div>

                <div class="day-events">
                  <div
                    v-for="event in getEventsForDay(day)"
                    :key="event.id"
                    class="event-item"
                    :class="{
                      'is-milestone': event.isMilestone,
                      'has-deadline': event.hasDeadline
                    }"
                    :style="{ backgroundColor: event.color }"
                    draggable="true"
                    @dragstart="handleDragStart(event, $event)"
                    @click="handleEventClick(event)"
                  >
                    <div class="event-name">{{ event.name }}</div>

                    <!-- Milestone marker -->
                    <el-icon v-if="event.isMilestone && showMilestones" class="milestone-icon">
                      <StarFilled />
                    </el-icon>

                    <!-- Deadline indicator -->
                    <div
                      v-if="event.hasDeadline && showDeadlines"
                      class="deadline-indicator"
                    />

                    <!-- Dependency indicators -->
                    <div v-if="showDependencies && event.hasDependencies" class="dependencies-badge">
                      <el-tag size="small" type="info">
                        {{ event.dependencyCount }} deps
                      </el-tag>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Week View -->
          <div v-else-if="viewMode === 'week'" class="week-view">
            <div class="week-header">
              <div class="time-column"></div>
              <div
                v-for="day in weekDays"
                :key="day.key"
                class="day-header"
                :class="{ 'is-today': isToday(day.date) }"
              >
                <div class="day-name">{{ day.name }}</div>
                <div class="day-number">{{ formatDate(day.date, 'd') }}</div>
              </div>
            </div>

            <div class="week-grid">
              <div class="time-column">
                <div v-for="hour in 24" :key="hour" class="hour-cell">
                  {{ hour - 1 }}:00
                </div>
              </div>

              <div
                v-for="day in weekDays"
                :key="day.key"
                class="day-column"
                :class="{ 'is-today': isToday(day.date) }"
                @drop="handleDrop(day.date, $event)"
                @dragover.prevent
                @dragenter.prevent
              >
                <div
                  v-for="event in getEventsForDay(day.date)"
                  :key="event.id"
                  class="week-event"
                  :style="{
                    top: `${event.startHour * 60 + event.startMinute}px`,
                    height: `${event.durationMinutes}px`,
                    backgroundColor: event.color
                  }"
                  draggable="true"
                  @dragstart="handleDragStart(event, $event)"
                  @click="handleEventClick(event)"
                >
                  <div class="event-time">
                    {{ formatTime(event.startHour, event.startMinute) }}
                  </div>
                  <div class="event-name">{{ event.name }}</div>

                  <el-icon v-if="event.isMilestone" class="milestone-icon">
                    <StarFilled />
                  </el-icon>
                </div>
              </div>
            </div>
          </div>

          <!-- Day View -->
          <div v-else-if="viewMode === 'day'" class="day-view">
            <div class="day-header">
              <div class="time-column"></div>
              <div class="day-title">
                <h2>{{ formatDate(currentDate, 'MMMM d, yyyy') }}</h2>
              </div>
            </div>

            <div class="day-grid">
              <div class="time-column">
                <div v-for="hour in 24" :key="hour" class="hour-cell">
                  {{ hour - 1 }}:00
                </div>
              </div>

              <div
                class="day-content"
                @drop="handleDrop(currentDate, $event)"
                @dragover.prevent
                @dragenter.prevent
              >
                <div
                  v-for="event in getEventsForDay(currentDate)"
                  :key="event.id"
                  class="day-event"
                  :style="{
                    top: `${event.startHour * 60 + event.startMinute}px`,
                    height: `${event.durationMinutes}px`,
                    backgroundColor: event.color
                  }"
                  draggable="true"
                  @dragstart="handleDragStart(event, $event)"
                  @click="handleEventClick(event)"
                >
                  <div class="event-header">
                    <span class="event-name">{{ event.name }}</span>
                    <el-icon v-if="event.isMilestone" class="milestone-icon">
                      <StarFilled />
                    </el-icon>
                  </div>

                  <div class="event-time">
                    {{ formatTime(event.startHour, event.startMinute) }} -
                    {{ formatTime(event.endHour, event.endMinute) }}
                  </div>

                  <div class="event-details">
                    <el-tag size="small" :type="getProgressType(event.progress)">
                      {{ event.progress }}%
                    </el-tag>
                    <span v-if="event.assignee" class="event-assignee">
                      {{ event.assignee }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </el-scrollbar>
    </div>

    <!-- Task Detail Dialog -->
    <el-dialog
      v-model="showTaskDialog"
      :title="selectedEvent?.name"
      width="600px"
    >
      <div v-if="selectedEvent" class="task-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="Task ID">
            {{ selectedEvent.id }}
          </el-descriptions-item>
          <el-descriptions-item label="Status">
            <el-tag :type="getStatusType(selectedEvent.status)">
              {{ selectedEvent.status }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="Start Date">
            {{ formatDate(selectedEvent.startDate, 'MMM d, yyyy') }}
          </el-descriptions-item>
          <el-descriptions-item label="End Date">
            {{ formatDate(selectedEvent.endDate, 'MMM d, yyyy') }}
          </el-descriptions-item>
          <el-descriptions-item label="Progress">
            <el-progress :percentage="selectedEvent.progress" />
          </el-descriptions-item>
          <el-descriptions-item label="Assignee">
            {{ selectedEvent.assignee || 'Unassigned' }}
          </el-descriptions-item>
          <el-descriptions-item label="Duration" :span="2">
            {{ selectedEvent.duration }} days
          </el-descriptions-item>
        </el-descriptions>

        <div v-if="selectedEvent.dependencies?.length" class="dependencies-section">
          <h4>Dependencies</h4>
          <el-tag
            v-for="dep in selectedEvent.dependencies"
            :key="dep.id"
            class="dependency-tag"
          >
            {{ dep.name }}
          </el-tag>
        </div>
      </div>

      <template #footer>
        <el-button @click="showTaskDialog = false">Close</el-button>
        <el-button type="primary" @click="editTask(selectedEvent)">
          Edit Task
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import {
  ArrowLeft,
  ArrowRight,
  ArrowDown,
  Refresh,
  Setting,
  StarFilled
} from '@element-plus/icons-vue'
import { format, addMonths, addWeeks, addDays, startOfMonth, endOfMonth, startOfWeek, endOfWeek, eachDayOfInterval, isSameMonth, isToday, isWeekend, parseISO } from 'date-fns'

/**
 * CalendarView Component
 * Displays Gantt tasks in calendar format with month, week, and day views
 *
 * @props {Array} tasks - Array of task objects
 * @props {Date} initialDate - Initial date to display
 * @props {string} initialViewMode - Initial view mode (month/week/day)
 *
 * @emits {Date} date-change - Emitted when current date changes
 * @emits {string} view-change - Emitted when view mode changes
 * @emits {Object} task-click - Emitted when a task is clicked
 * @emits {Object} task-reschedule - Emitted when a task is dragged to new date
 */

const props = defineProps({
  tasks: {
    type: Array,
    default: () => []
  },
  initialDate: {
    type: Date,
    default: () => new Date()
  },
  initialViewMode: {
    type: String,
    default: 'month',
    validator: (value) => ['month', 'week', 'day'].includes(value)
  }
})

const emit = defineEmits(['date-change', 'view-change', 'task-click', 'task-reschedule'])

// State
const loading = ref(false)
const viewMode = ref(props.initialViewMode)
const currentDate = ref(props.initialDate)
const showWeekends = ref(true)
const showMilestones = ref(true)
const showDeadlines = ref(true)
const showDependencies = ref(true)
const showTaskDialog = ref(false)
const selectedEvent = ref(null)
const draggedEvent = ref(null)

// Computed
const weekdays = computed(() => [
  { index: 0, label: 'Sun' },
  { index: 1, label: 'Mon' },
  { index: 2, label: 'Tue' },
  { index: 3, label: 'Wed' },
  { index: 4, label: 'Thu' },
  { index: 5, label: 'Fri' },
  { index: 6, label: 'Sat' }
])

const currentDateLabel = computed(() => {
  if (viewMode.value === 'month') {
    return format(currentDate.value, 'MMMM yyyy')
  } else if (viewMode.value === 'week') {
    const start = startOfWeek(currentDate.value)
    const end = endOfWeek(currentDate.value)
    return `${format(start, 'MMM d')} - ${format(end, 'MMM d, yyyy')}`
  } else {
    return format(currentDate.value, 'MMMM d, yyyy')
  }
})

const monthDays = computed(() => {
  const start = startOfWeek(startOfMonth(currentDate.value))
  const end = endOfWeek(endOfMonth(currentDate.value))
  return eachDayOfInterval({ start, end })
})

const weekDays = computed(() => {
  const start = startOfWeek(currentDate.value)
  const end = endOfWeek(currentDate.value)
  const days = eachDayOfInterval({ start, end })

  return days.map((day, index) => ({
    key: index,
    date: day,
    name: format(day, 'EEE'),
    number: format(day, 'd')
  }))
})

// Methods
const setViewMode = (mode) => {
  viewMode.value = mode
  emit('view-change', mode)
}

const goToToday = () => {
  currentDate.value = new Date()
  emit('date-change', currentDate.value)
}

const navigate = (direction) => {
  if (viewMode.value === 'month') {
    currentDate.value = addMonths(currentDate.value, direction)
  } else if (viewMode.value === 'week') {
    currentDate.value = addWeeks(currentDate.value, direction)
  } else {
    currentDate.value = addDays(currentDate.value, direction)
  }
  emit('date-change', currentDate.value)
}

const formatDate = (date, formatStr) => {
  return format(date, formatStr)
}

const formatTime = (hour, minute) => {
  return `${hour.toString().padStart(2, '0')}:${minute.toString().padStart(2, '0')}`
}

const getEventsForDay = (day) => {
  const dayStr = format(day, 'yyyy-MM-dd')

  return props.tasks
    .filter(task => {
      const startDate = format(parseISO(task.startDate), 'yyyy-MM-dd')
      const endDate = format(parseISO(task.endDate), 'yyyy-MM-dd')
      return dayStr >= startDate && dayStr <= endDate
    })
    .map(task => {
      const startDate = parseISO(task.startDate)
      const endDate = parseISO(task.endDate)

      return {
        id: task.id,
        name: task.name,
        startDate: task.startDate,
        endDate: task.endDate,
        color: task.color || '#409EFF',
        progress: task.progress || 0,
        status: task.status,
        assignee: task.assignee,
        duration: task.duration || 1,
        isMilestone: task.isMilestone || false,
        hasDeadline: task.hasDeadline || false,
        hasDependencies: task.dependencies?.length > 0,
        dependencyCount: task.dependencies?.length || 0,
        dependencies: task.dependencies || [],
        // Week/Day view specific
        startHour: task.startHour || 9,
        startMinute: task.startMinute || 0,
        endHour: task.endHour || 17,
        endMinute: task.endMinute || 0,
        durationMinutes: task.durationMinutes || 480
      }
    })
}

const handleDragStart = (event, dragEvent) => {
  draggedEvent.value = event
  dragEvent.dataTransfer.effectAllowed = 'move'
}

const handleDrop = (date, dropEvent) => {
  if (!draggedEvent.value) return

  const newDate = format(date, 'yyyy-MM-dd')
  emit('task-reschedule', {
    taskId: draggedEvent.value.id,
    newStartDate: newDate,
    task: draggedEvent.value
  })

  draggedEvent.value = null
}

const handleEventClick = (event) => {
  selectedEvent.value = event
  showTaskDialog.value = true
  emit('task-click', event)
}

const editTask = (task) => {
  showTaskDialog.value = false
  // Emit event to open edit dialog
  emit('task-click', { ...task, action: 'edit' })
}

const handleViewCommand = (command) => {
  // Handle dropdown commands if needed
}

const refreshData = async () => {
  loading.value = true
  try {
    // Emit event to refresh data
    emit('date-change', currentDate.value)
    await new Promise(resolve => setTimeout(resolve, 500))
  } finally {
    loading.value = false
  }
}

const getProgressType = (progress) => {
  if (progress >= 100) return 'success'
  if (progress >= 50) return 'warning'
  return 'danger'
}

const getStatusType = (status) => {
  const statusMap = {
    completed: 'success',
    in_progress: 'warning',
    not_started: 'info',
    delayed: 'danger'
  }
  return statusMap[status] || 'info'
}

onMounted(() => {
  emit('date-change', currentDate.value)
  emit('view-change', viewMode.value)
})
</script>

<style scoped lang="scss">
.calendar-view {
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

    .toolbar-left,
    .toolbar-right {
      display: flex;
      align-items: center;
      gap: 12px;
    }

    .current-date {
      font-size: 16px;
      font-weight: 600;
      color: #303133;
      min-width: 200px;
      text-align: center;
    }
  }

  &__container {
    flex: 1;
    overflow: hidden;
  }
}

.calendar-grid {
  &--month,
  &--week,
  &--day {
    height: 100%;
  }
}

/* Month View */
.month-view {
  .weekdays-header {
    display: grid;
    grid-template-columns: repeat(7, 1fr);
    background: #f5f7fa;
    border-bottom: 1px solid #dcdfe6;

    .weekday-cell {
      padding: 8px;
      text-align: center;
      font-weight: 600;
      color: #606266;

      &.is-weekend {
        color: #909399;
      }
    }
  }

  .month-grid {
    display: grid;
    grid-template-columns: repeat(7, 1fr);
    grid-auto-rows: minmax(100px, 1fr);
  }

  .day-cell {
    border: 1px solid #ebeef5;
    padding: 4px;
    min-height: 100px;
    position: relative;

    &.is-today {
      background: #ecf5ff;
    }

    &.is-weekend {
      background: #fafafa;
    }

    &.is-other-month {
      background: #fafafa;
      opacity: 0.6;
    }

    .day-number {
      font-size: 14px;
      font-weight: 600;
      color: #606266;
      margin-bottom: 4px;
    }

    .day-events {
      display: flex;
      flex-direction: column;
      gap: 2px;
    }

    .event-item {
      padding: 4px 6px;
      border-radius: 4px;
      font-size: 12px;
      color: #fff;
      cursor: pointer;
      position: relative;

      &:hover {
        opacity: 0.9;
        transform: scale(1.02);
      }

      &.is-milestone {
        border: 2px solid #f59e0b;
      }

      .event-name {
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }

      .milestone-icon {
        position: absolute;
        top: 2px;
        right: 2px;
        color: #f59e0b;
      }

      .deadline-indicator {
        position: absolute;
        bottom: 0;
        left: 0;
        right: 0;
        height: 3px;
        background: #f56c6c;
        border-radius: 0 0 4px 4px;
      }

      .dependencies-badge {
        position: absolute;
        top: 2px;
        right: 2px;
      }
    }
  }
}

/* Week View */
.week-view {
  .week-header {
    display: grid;
    grid-template-columns: 60px repeat(7, 1fr);
    background: #f5f7fa;
    border-bottom: 1px solid #dcdfe6;

    .day-header {
      padding: 12px;
      text-align: center;
      border-left: 1px solid #ebeef5;

      &.is-today {
        background: #ecf5ff;
      }

      .day-name {
        font-size: 12px;
        color: #909399;
      }

      .day-number {
        font-size: 18px;
        font-weight: 600;
        color: #303133;
      }
    }
  }

  .week-grid {
    display: grid;
    grid-template-columns: 60px repeat(7, 1fr);
    position: relative;
  }

  .time-column {
    background: #fafafa;
    border-right: 1px solid #dcdfe6;

    .hour-cell {
      height: 60px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 11px;
      color: #909399;
      border-bottom: 1px dashed #ebeef5;
    }
  }

  .day-column {
    position: relative;
    height: calc(24 * 60px);
    border-left: 1px solid #ebeef5;
    background-image: linear-gradient(to bottom, #ebeef5 1px, transparent 1px);
    background-size: 100% 60px;

    &.is-today {
      background-color: #ecf5ff;
    }

    .week-event {
      position: absolute;
      left: 4px;
      right: 4px;
      padding: 4px 8px;
      border-radius: 4px;
      font-size: 12px;
      color: #fff;
      cursor: pointer;

      &:hover {
        opacity: 0.9;
        z-index: 10;
      }

      .event-time {
        font-size: 11px;
        opacity: 0.9;
      }

      .event-name {
        margin-top: 2px;
        font-weight: 600;
      }

      .milestone-icon {
        position: absolute;
        top: 4px;
        right: 4px;
        color: #f59e0b;
      }
    }
  }
}

/* Day View */
.day-view {
  .day-header {
    display: grid;
    grid-template-columns: 60px 1fr;
    background: #f5f7fa;
    border-bottom: 1px solid #dcdfe6;

    .day-title {
      padding: 12px 16px;

      h2 {
        margin: 0;
        font-size: 18px;
        font-weight: 600;
        color: #303133;
      }
    }
  }

  .day-grid {
    display: grid;
    grid-template-columns: 60px 1fr;
  }

  .time-column {
    background: #fafafa;
    border-right: 1px solid #dcdfe6;

    .hour-cell {
      height: 60px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 11px;
      color: #909399;
      border-bottom: 1px dashed #ebeef5;
    }
  }

  .day-content {
    position: relative;
    height: calc(24 * 60px);
    background-image: linear-gradient(to bottom, #ebeef5 1px, transparent 1px);
    background-size: 100% 60px;
  }

  .day-event {
    position: absolute;
    left: 8px;
    right: 8px;
    padding: 8px 12px;
    border-radius: 4px;
    color: #fff;
    cursor: pointer;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

    &:hover {
      opacity: 0.9;
      z-index: 10;
      transform: scale(1.01);
    }

    .event-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 4px;

      .event-name {
        font-weight: 600;
        font-size: 14px;
      }

      .milestone-icon {
        color: #f59e0b;
      }
    }

    .event-time {
      font-size: 12px;
      opacity: 0.9;
      margin-bottom: 8px;
    }

    .event-details {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .event-assignee {
        font-size: 12px;
        opacity: 0.9;
      }
    }
  }
}

/* Task Detail Dialog */
.task-detail {
  .dependencies-section {
    margin-top: 16px;

    h4 {
      margin: 0 0 8px 0;
      font-size: 14px;
      color: #606266;
    }

    .dependency-tag {
      margin-right: 8px;
      margin-bottom: 8px;
    }
  }
}

/* Responsive Design */
@media (max-width: 768px) {
  .calendar-view__toolbar {
    flex-direction: column;
    gap: 12px;

    .toolbar-left,
    .toolbar-right {
      width: 100%;
      justify-content: center;
      flex-wrap: wrap;
    }
  }

  .month-view .month-grid {
    grid-auto-rows: minmax(80px, 1fr);
  }

  .day-cell {
    min-height: 80px !important;

    .event-item {
      font-size: 11px;
      padding: 2px 4px;
    }
  }

  .week-view .week-header,
  .week-view .week-grid {
    grid-template-columns: 50px repeat(7, 1fr);
  }

  .day-view .day-header,
  .day-view .day-grid {
    grid-template-columns: 50px 1fr;
  }
}
</style>
