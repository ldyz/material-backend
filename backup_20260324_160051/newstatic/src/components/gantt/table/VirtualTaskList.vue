<template>
  <div
    ref="taskListRef"
    class="virtual-task-list"
    :style="containerStyle"
    @scroll.passive="handleScroll"
  >
    <!-- Task list container -->
    <div
      class="virtual-task-list__content"
      :style="contentStyle"
    >
      <!-- Table header -->
      <div
        class="virtual-task-list__header"
        :style="headerStyle"
      >
        <slot
          name="header"
          :columns="columns"
        >
          <div
            v-for="column in columns"
            :key="column.key"
            class="virtual-task-list__header-cell"
            :style="{
              width: column.width ? `${column.width}px` : 'auto',
              minWidth: column.minWidth ? `${column.minWidth}px` : 'auto',
              textAlign: column.align || 'left'
            }"
          >
            {{ column.label }}
          </div>
        </slot>
      </div>

      <!-- Virtual scrolling task rows -->
      <div
        class="virtual-task-list__body"
        :style="bodyStyle"
      >
        <RecycleScroller
          ref="scrollerRef"
          class="virtual-task-list__scroller"
          :items="filteredTasks"
          :item-size="rowHeight"
          :buffer="buffer"
          key-field="id"
          :emit-update="true"
          @update="handleScrollerUpdate"
          @resize="handleScrollerResize"
          v-slot="{ item: task, index }"
        >
          <div
            class="virtual-task-list__row"
            :class="{
              'is-selected': selectedTaskId === task.id,
              'is-collapsed': isTaskCollapsed(task.id),
              'is-editing': editingTaskId === task.id
            }"
            :style="{
              height: `${rowHeight}px`
            }"
            @click="handleRowClick(task, index)"
            @dblclick="handleRowDoubleClick(task, index)"
            @mouseenter="handleRowMouseEnter(task, index)"
            @mouseleave="handleRowMouseLeave(task, index)"
          >
            <!-- Tree toggle -->
            <div
              v-if="hasChildren(task)"
              class="virtual-task-list__toggle"
              @click.stop="handleToggle(task)"
            >
              <ElIcon :size="14">
                <component :is="isTaskCollapsed(task.id) ? 'ArrowRight' : 'ArrowDown'" />
              </ElIcon>
            </div>

            <!-- Task cells -->
            <div
              v-for="column in columns"
              :key="column.key"
              class="virtual-task-list__cell"
              :style="{
                width: column.width ? `${column.width}px` : 'auto',
                minWidth: column.minWidth ? `${column.minWidth}px` : 'auto',
                textAlign: column.align || 'left',
                paddingLeft: `${(getTaskLevel(task) * 20) + 10}px`
              }"
            >
              <!-- Inline editing -->
              <template v-if="editingTaskId === task.id && column.editable">
                <ElInput
                  v-if="column.type === 'text' || column.type === 'input'"
                  :model-value="task[column.key]"
                  size="small"
                  @blur="handleEditBlur(task, column, $event)"
                  @keyup.enter="handleEditEnter(task, column, $event)"
                  @keyup.esc="handleEditEscape"
                  :ref="el => setInputRef(el, task.id, column.key)"
                />
                <ElInputNumber
                  v-else-if="column.type === 'number'"
                  :model-value="task[column.key]"
                  size="small"
                  :controls="false"
                  @blur="handleEditBlur(task, column, $event)"
                  @keyup.enter="handleEditEnter(task, column, $event)"
                  @keyup.esc="handleEditEscape"
                />
                <ElSelect
                  v-else-if="column.type === 'select'"
                  :model-value="task[column.key]"
                  size="small"
                  @change="handleSelectChange(task, column, $event)"
                  @blur="handleEditBlur(task, column, $event)"
                >
                  <ElOption
                    v-for="option in column.options"
                    :key="option.value"
                    :label="option.label"
                    :value="option.value"
                  />
                </ElSelect>
              </template>

              <!-- Normal display -->
              <template v-else>
                <slot
                  :name="`column-${column.key}`"
                  :task="task"
                  :column="column"
                  :index="index"
                >
                  <span v-if="column.key === 'name'">
                    {{ task.name || task.task_name || '未命名' }}
                  </span>
                  <span v-else-if="column.key === 'progress'">
                    {{ task.progress || 0 }}%
                  </span>
                  <span v-else-if="column.key === 'duration'">
                    {{ task.duration || 1 }}天
                  </span>
                  <span v-else-if="column.key === 'status'">
                    <ElTag
                      :type="getStatusType(task.status)"
                      size="small"
                    >
                      {{ getStatusLabel(task.status) }}
                    </ElTag>
                  </span>
                  <span v-else-if="column.key === 'priority'">
                    <ElTag
                      :type="getPriorityType(task.priority)"
                      size="small"
                    >
                      {{ getPriorityLabel(task.priority) }}
                    </ElTag>
                  </span>
                  <span v-else>
                    {{ task[column.key] }}
                  </span>
                </slot>
              </template>
            </div>
          </div>
        </RecycleScroller>
      </div>
    </div>

    <!-- Loading indicator -->
    <div
      v-if="loading"
      class="virtual-task-list__loading"
    >
      <ElIcon class="is-loading">
        <Loading />
      </ElIcon>
    </div>

    <!-- Empty state -->
    <div
      v-if="!loading && filteredTasks.length === 0"
      class="virtual-task-list__empty"
    >
      <ElEmpty
        :description="emptyDescription"
        :image-size="80"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { RecycleScroller } from 'vue-virtual-scroller'
import 'vue-virtual-scroller/dist/vue-virtual-scroller.css'
import { Loading, ArrowRight, ArrowDown } from '@element-plus/icons-vue'
import { ElIcon, ElInput, ElInputNumber, ElSelect, ElOption, ElTag, ElEmpty } from 'element-plus'

/**
 * Virtual Scrolling Task List Component
 *
 * High-performance task list with virtual scrolling support.
 * Uses vue-virtual-scroller's RecycleScroller for efficient rendering.
 *
 * Features:
 * - Virtual scrolling with dynamic row height
 * - Tree structure support (expand/collapse)
 * - Inline editing support
 * - Sync scrolling with timeline
 *
 * @props
 * @param {Array} tasks - Task list to display
 * @param {Array} columns - Column definitions
 * @param {number} rowHeight - Height of each row (default: 60)
 * @param {number} buffer - Buffer size in pixels (default: 500)
 * @param {number} containerHeight - Container height in pixels
 * @param {Set} collapsedTasks - Set of collapsed task IDs
 * @param {number} selectedTaskId - Selected task ID
 * @param {boolean} loading - Loading state
 * @param {string} searchKeyword - Search keyword
 */

const props = withDefaults(defineProps({
  tasks: {
    type: Array,
    default: () => []
  },
  columns: {
    type: Array,
    default: () => [
      { key: 'name', label: '任务名称', width: 300, editable: true },
      { key: 'duration', label: '工期', width: 80, align: 'center' },
      { key: 'progress', label: '进度', width: 80, align: 'center' },
      { key: 'status', label: '状态', width: 100, align: 'center' },
      { key: 'priority', label: '优先级', width: 100, align: 'center' }
    ]
  },
  rowHeight: {
    type: Number,
    default: 60
  },
  buffer: {
    type: Number,
    default: 500
  },
  containerHeight: {
    type: Number,
    default: 600
  },
  collapsedTasks: {
    type: Set,
    default: () => new Set()
  },
  selectedTaskId: {
    type: Number,
    default: null
  },
  loading: {
    type: Boolean,
    default: false
  },
  searchKeyword: {
    type: String,
    default: ''
  }
}), {
  tasks: () => [],
  columns: () => [],
  rowHeight: 60,
  buffer: 500,
  containerHeight: 600,
  collapsedTasks: () => new Set(),
  selectedTaskId: null,
  loading: false,
  searchKeyword: ''
})

const emit = defineEmits([
  'scroll',
  'resize',
  'row-click',
  'row-dblclick',
  'row-mouseenter',
  'row-mouseleave',
  'toggle',
  'edit',
  'selection-change'
])

// Refs
const taskListRef = ref(null)
const scrollerRef = ref(null)
const inputRefs = ref(new Map())

// State
const scrollTop = ref(0)
const editingTaskId = ref(null)
const editingColumnKey = ref(null)

// ==================== Computed Properties ====================

/**
 * Container style
 */
const containerStyle = computed(() => ({
  height: `${props.containerHeight}px`,
  position: 'relative',
  overflow: 'auto',
  '-webkit-overflow-scrolling': 'touch'
}))

/**
 * Content style
 */
const contentStyle = computed(() => ({
  height: `${totalHeight.value}px`,
  position: 'relative',
  minHeight: '100%'
}))

/**
 * Total height
 */
const totalHeight = computed(() => {
  return filteredTasks.value.length * props.rowHeight
})

/**
 * Header style
 */
const headerStyle = computed(() => ({
  position: 'sticky',
  top: 0,
  left: 0,
  right: 0,
  height: '50px',
  backgroundColor: '#fff',
  zIndex: 10,
  borderBottom: '1px solid #e4e7ed'
}))

/**
 * Body style
 */
const bodyStyle = computed(() => ({
  position: 'absolute',
  top: '50px',
  left: 0,
  right: 0,
  bottom: 0
}))

/**
 * Filter tasks based on search keyword and collapse state
 */
const filteredTasks = computed(() => {
  let tasks = [...props.tasks]

  // Apply search filter
  if (props.searchKeyword) {
    const keyword = props.searchKeyword.toLowerCase()
    tasks = tasks.filter(task =>
      (task.name || task.task_name || '').toLowerCase().includes(keyword)
    )
  }

  // Filter out hidden tasks (children of collapsed parents)
  return tasks.filter(task => !isTaskHidden(task.id, tasks))
})

/**
 * Empty state description
 */
const emptyDescription = computed(() => {
  return props.searchKeyword ? '未找到匹配的任务' : '暂无任务'
})

// ==================== Methods ====================

/**
 * Check if task has children
 */
const hasChildren = (task) => {
  return props.tasks.some(t => t.parent_id === task.id)
}

/**
 * Get task level in tree
 */
const getTaskLevel = (task) => {
  let level = 0
  let currentTask = task
  while (currentTask.parent_id) {
    level++
    currentTask = props.tasks.find(t => t.id === currentTask.parent_id)
    if (!currentTask) break
  }
  return level
}

/**
 * Check if task is collapsed
 */
const isTaskCollapsed = (taskId) => {
  return props.collapsedTasks.has(taskId)
}

/**
 * Check if task is hidden (parent is collapsed)
 */
const isTaskHidden = (taskId, tasks) => {
  const task = tasks.find(t => t.id === taskId)
  if (!task || !task.parent_id) return false

  const parent = tasks.find(t => t.id === task.parent_id)
  if (!parent) return false

  if (props.collapsedTasks.has(parent.id)) return true
  return isTaskHidden(parent.id, tasks)
}

/**
 * Handle scroll event
 */
const handleScroll = (event) => {
  const target = event.target
  scrollTop.value = target.scrollTop

  emit('scroll', {
    scrollTop: scrollTop.value,
    scrollLeft: target.scrollLeft
  })
}

/**
 * Handle scroller update
 */
const handleScrollerUpdate = () => {
  nextTick(() => {
    if (taskListRef.value) {
      // Refresh viewport dimensions
    }
  })
}

/**
 * Handle scroller resize
 */
const handleScrollerResize = () => {
  if (taskListRef.value) {
    emit('resize', {
      width: taskListRef.value.clientWidth,
      height: taskListRef.value.clientHeight
    })
  }
}

/**
 * Handle row click
 */
const handleRowClick = (task, index) => {
  emit('row-click', { task, index })
  emit('selection-change', task.id)
}

/**
 * Handle row double click
 */
const handleRowDoubleClick = (task, index) => {
  // Start inline editing
  startEditing(task, props.columns[0])
  emit('row-dblclick', { task, index })
}

/**
 * Handle row mouse enter
 */
const handleRowMouseEnter = (task, index) => {
  emit('row-mouseenter', { task, index })
}

/**
 * Handle row mouse leave
 */
const handleRowMouseLeave = (task, index) => {
  emit('row-mouseleave', { task, index })
}

/**
 * Handle toggle click
 */
const handleToggle = (task) => {
  emit('toggle', task)
}

/**
 * Start inline editing
 */
const startEditing = (task, column) => {
  if (!column.editable) return

  editingTaskId.value = task.id
  editingColumnKey.value = column.key

  nextTick(() => {
    const refKey = `${task.id}-${column.key}`
    const inputRef = inputRefs.value.get(refKey)
    if (inputRef && inputRef.focus) {
      inputRef.focus()
    }
  })
}

/**
 * Stop inline editing
 */
const stopEditing = () => {
  editingTaskId.value = null
  editingColumnKey.value = null
}

/**
 * Handle edit blur
 */
const handleEditBlur = (task, column, event) => {
  const value = event.target.value

  emit('edit', {
    task,
    column,
    value
  })

  stopEditing()
}

/**
 * Handle edit enter key
 */
const handleEditEnter = (task, column, event) => {
  const value = event.target.value

  emit('edit', {
    task,
    column,
    value
  })

  stopEditing()
}

/**
 * Handle edit escape key
 */
const handleEditEscape = () => {
  stopEditing()
}

/**
 * Handle select change
 */
const handleSelectChange = (task, column, value) => {
  emit('edit', {
    task,
    column,
    value
  })
}

/**
 * Set input ref
 */
const setInputRef = (el, taskId, columnKey) => {
  if (el) {
    const refKey = `${taskId}-${columnKey}`
    inputRefs.value.set(refKey, el)
  }
}

/**
 * Get status type for tag
 */
const getStatusType = (status) => {
  const types = {
    completed: 'success',
    in_progress: 'primary',
    not_started: 'info',
    delayed: 'danger'
  }
  return types[status] || 'info'
}

/**
 * Get status label
 */
const getStatusLabel = (status) => {
  const labels = {
    completed: '已完成',
    in_progress: '进行中',
    not_started: '未开始',
    delayed: '延期'
  }
  return labels[status] || status
}

/**
 * Get priority type for tag
 */
const getPriorityType = (priority) => {
  const types = {
    high: 'danger',
    medium: 'warning',
    low: 'info'
  }
  return types[priority] || 'info'
}

/**
 * Get priority label
 */
const getPriorityLabel = (priority) => {
  const labels = {
    high: '高',
    medium: '中',
    low: '低'
  }
  return labels[priority] || priority
}

/**
 * Scroll to task
 * @param {number} taskIndex - Task index
 * @param {string} alignment - Alignment (start, center, end, auto)
 */
const scrollToTask = (taskIndex, alignment = 'auto') => {
  if (!scrollerRef.value || taskIndex < 0 || taskIndex >= filteredTasks.value.length) {
    return
  }

  scrollerRef.value.scrollToItem(taskIndex, alignment)
}

/**
 * Scroll to position
 * @param {number} scrollTop - Scroll top position
 */
const scrollToPosition = (scrollTop) => {
  if (!taskListRef.value) {
    return
  }

  taskListRef.value.scrollTop = scrollTop
}

/**
 * Refresh layout
 */
const refresh = () => {
  if (scrollerRef.value) {
    scrollerRef.value.forceUpdate()
  }
}

// ==================== Lifecycle ====================

onMounted(() => {
  // Initialize
})

onUnmounted(() => {
  // Cleanup
  inputRefs.value.clear()
})

// ==================== Watchers ====================

// Watch for tasks changes
watch(() => props.tasks, () => {
  refresh()
}, { deep: true })

// Expose methods
defineExpose({
  scrollToTask,
  scrollToPosition,
  refresh,
  startEditing,
  stopEditing
})
</script>

<style scoped lang="scss">
.virtual-task-list {
  position: relative;
  overflow: auto;
  -webkit-overflow-scrolling: touch;
  user-select: none;
  background-color: #fff;

  &__content {
    position: relative;
    min-width: 100%;
    min-height: 100%;
  }

  &__header {
    display: flex;
    align-items: center;
    padding: 0 16px;
    background-color: #f5f7fa;
    border-bottom: 1px solid #e4e7ed;
    font-weight: 500;
    color: #606266;
    will-change: transform;
  }

  &__header-cell {
    padding: 12px 8px;
    font-size: 14px;
    flex-shrink: 0;
  }

  &__body {
    position: relative;
  }

  &__scroller {
    height: 100%;

    :deep(.vue-recycle-scroller__item-wrapper) {
      position: relative;
    }
  }

  &__row {
    position: relative;
    display: flex;
    align-items: center;
    border-bottom: 1px solid #f0f0f0;
    cursor: pointer;
    transition: background-color 0.2s;

    &:hover {
      background-color: #f5f7fa;
    }

    &.is-selected {
      background-color: #ecf5ff;
    }

    &.is-editing {
      background-color: #fff;
    }
  }

  &__toggle {
    position: absolute;
    left: 0;
    width: 20px;
    height: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    color: #909399;
    transition: color 0.2s;

    &:hover {
      color: var(--el-color-primary);
    }
  }

  &__cell {
    padding: 8px;
    font-size: 14px;
    color: #606266;
    flex-shrink: 0;
    display: flex;
    align-items: center;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;

    :deep(.el-input) {
      width: 100%;
    }

    :deep(.el-input-number) {
      width: 100%;
    }

    :deep(.el-select) {
      width: 100%;
    }
  }

  &__loading {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    z-index: 100;
    font-size: 24px;
    color: var(--el-color-primary);
  }

  &__empty {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 100%;
    text-align: center;
  }
}
</style>
