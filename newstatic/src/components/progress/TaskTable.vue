<template>
  <div class="task-table-wrapper" :class="{ 'is-collapsed': isCollapsed }">
    <!-- 拖拽提示浮层 -->
    <div v-if="draggedTask" class="drag-indicator">
      <div class="drag-indicator-content">
        <el-icon class="drag-icon"><Rank /></el-icon>
        <span class="drag-text">
          <template v-if="dropPosition === 'child'">
            <span class="keyboard-hint">按住 Ctrl</span> 放置为子任务
          </template>
          <template v-else-if="dropPosition === 'before'">
            放置到任务之前
          </template>
          <template v-else-if="dropPosition === 'after'">
            放置到任务之后
          </template>
          <template v-else>
            松开以放置任务
          </template>
        </span>
      </div>
    </div>

    <!-- 内容区域（表头在 GanttHeader 组件中统一处理） -->
    <div
      class="table-content"
      :class="{ 'is-dragging-over-root': draggedTask }"
      @dragover.prevent
      @drop="handleDropToRoot"
      :style="{ height: totalHeight + 'px' }"
    >
        <!-- 分组显示 -->
        <template v-if="groupedTasks.length > 0 && groupMode">
      <div v-for="(group, groupIndex) in groupedTasks" :key="group.name">
        <!-- 分组头部 -->
        <div class="table-group-header" @click="$emit('toggle-group', group.name)">
          <span class="group-name">{{ group.label }}</span>
          <span class="group-count">{{ group.tasks.length }} 个任务</span>
          <el-icon class="group-toggle" :class="{ 'is-collapsed': collapsedGroups.has(group.name) }">
            <ArrowDown />
          </el-icon>
        </div>

        <!-- 分组任务列表 -->
        <div v-show="!collapsedGroups.has(group.name)">
          <div
            v-for="(task, taskIndex) in group.tasks"
            :key="task.id"
            v-show="!isTaskHidden(task)"
            class="table-row"
            :class="{
              'is-selected': selectedTaskId === task.id,
              'is-critical': task.is_critical && showCriticalPath,
              'is-milestone': isMilestone(task),
              'is-dragging-over': dragOverTaskId === task.id,
              'drop-before': dragOverTaskId === task.id && dropPosition === 'before',
              'drop-after': dragOverTaskId === task.id && dropPosition === 'after',
              'drop-child': dragOverTaskId === task.id && dropPosition === 'child'
            }"
            :style="{ height: rowHeight + 'px' }"
            :data-task-id="task.id"
            :draggable="true"
            @click="$emit('row-click', task)"
            @dblclick="$emit('row-dblclick', task)"
            @contextmenu="handleContextMenu($event, task)"
            @dragstart="handleDragStart($event, task)"
            @dragover.prevent="handleDragOver($event, task)"
            @dragleave="handleDragLeave($event, task)"
            @drop="handleDrop($event, task)"
            @dragend="handleDragEnd"
          >
            <!-- 任务名称 -->
            <div
              class="row-column column-name"
              @dblclick="startEdit(task, 'name', $event)"
              :class="{ 'is-editing': editingCell?.taskId === task.id && editingCell?.field === 'name' }"
            >
              <template v-if="editingCell?.taskId === task.id && editingCell?.field === 'name'">
                <el-input
                  ref="editInput"
                  v-model="editingCell.value"
                  size="small"
                  @blur="saveEdit"
                  @keyup.enter="saveEdit"
                  @keyup.esc="cancelEdit"
                  autofocus
                />
              </template>
              <template v-else>
                <!-- 树形缩进和展开/收起按钮 -->
                <div class="task-tree-indent" :style="{ paddingLeft: (getTaskDepth(task) * 20) + 'px' }">
                  <el-icon
                    v-if="hasChildren(task)"
                    class="tree-toggle-icon"
                    :class="{ 'is-collapsed': isTaskCollapsed(task) }"
                    @click.stop="toggleTaskCollapse(task)"
                  >
                    <ArrowDown />
                  </el-icon>
                  <span v-else class="tree-toggle-placeholder"></span>
                  <el-icon v-if="isMilestone(task)" class="milestone-icon"><Star /></el-icon>
                  <div class="task-name-text" :title="task.name">{{ task.name }}</div>
                </div>
                <div class="task-progress-mini">
                  <el-progress
                    :percentage="task.progress || 0"
                    :stroke-width="4"
                    :show-text="false"
                  />
                </div>
              </template>
            </div>

            <!-- 工期 -->
            <div
              class="row-column column-duration"
              @dblclick="startEdit(task, 'duration', $event)"
              :class="{ 'is-editing': editingCell?.taskId === task.id && editingCell?.field === 'duration' }"
            >
              <template v-if="editingCell?.taskId === task.id && editingCell?.field === 'duration'">
                <el-input
                  v-model.number="editingCell.value"
                  size="small"
                  type="number"
                  @blur="saveEdit"
                  @keyup.enter="saveEdit"
                  @keyup.esc="cancelEdit"
                  autofocus
                />
              </template>
              <template v-else>
                {{ getTaskDuration(task) }} 天
              </template>
            </div>

            <!-- 起止时间 -->
            <div
              class="row-column column-dates"
              @dblclick="startEdit(task, 'dates', $event)"
              :class="{ 'is-editing': editingCell?.taskId === task.id && editingCell?.field === 'dates' }"
            >
              <template v-if="editingCell?.taskId === task.id && editingCell?.field === 'dates'">
                <el-date-picker
                  v-model="editingCell.value"
                  type="daterange"
                  size="small"
                  format="YYYY/MM/DD"
                  value-format="YYYY-MM-DD"
                  @blur="saveEdit"
                  @change="saveEdit"
                  autofocus
                />
              </template>
              <template v-else>
                <div class="date-range">
                  <span class="date-start">{{ task.start ? formatDateShort(task.start) : '-' }}</span>
                  <span class="date-separator">→</span>
                  <span class="date-end">{{ task.end ? formatDateShort(task.end) : '-' }}</span>
                </div>
              </template>
            </div>

            <!-- 资源 -->
            <div class="row-column column-resources">
              <div v-if="task.resources && task.resources.length > 0" class="resource-tags">
                <el-tag
                  v-for="(res, idx) in task.resources.slice(0, 2)"
                  :key="idx"
                  size="small"
                  :type="getResourceTagType(res.type)"
                >
                  {{ res.resource_name || res.name }}{{ res.quantity ? `×${res.quantity}` : '' }}
                </el-tag>
                <el-tag v-if="task.resources.length > 2" size="small" type="info">
                  +{{ task.resources.length - 2 }}
                </el-tag>
              </div>
              <span v-else class="no-resources">-</span>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- 无分组显示 -->
    <template v-else>
      <div
        v-for="(task, index) in visibleTasks"
        :key="task.id"
        class="table-row"
        :class="{
          'is-selected': selectedTaskId === task.id,
          'is-critical': task.is_critical && showCriticalPath,
          'is-milestone': isMilestone(task),
          'is-dragging-over': dragOverTaskId === task.id,
          'drop-before': dragOverTaskId === task.id && dropPosition === 'before',
          'drop-after': dragOverTaskId === task.id && dropPosition === 'after',
          'drop-child': dragOverTaskId === task.id && dropPosition === 'child'
        }"
        :style="{ height: rowHeight + 'px' }"
        :data-task-id="task.id"
        :draggable="true"
        @click="$emit('row-click', task)"
        @dblclick="$emit('row-dblclick', task)"
        @contextmenu="handleContextMenu($event, task)"
        @dragstart="handleDragStart($event, task)"
        @dragover.prevent="handleDragOver($event, task)"
        @dragleave="handleDragLeave($event, task)"
        @drop="handleDrop($event, task)"
      >
        <!-- 任务名称 -->
        <div
          class="row-column column-name"
          @dblclick="startEdit(task, 'name', $event)"
          :class="{ 'is-editing': editingCell?.taskId === task.id && editingCell?.field === 'name' }"
        >
          <template v-if="editingCell?.taskId === task.id && editingCell?.field === 'name'">
            <el-input
              v-model="editingCell.value"
              size="small"
              @blur="saveEdit"
              @keyup.enter="saveEdit"
              @keyup.esc="cancelEdit"
              autofocus
            />
          </template>
          <template v-else>
            <!-- 树形缩进和展开/收起按钮 -->
            <div class="task-tree-indent" :style="{ paddingLeft: (getTaskDepth(task) * 20) + 'px' }">
              <el-icon
                v-if="hasChildren(task)"
                class="tree-toggle-icon"
                :class="{ 'is-collapsed': isTaskCollapsed(task) }"
                @click.stop="toggleTaskCollapse(task)"
              >
                <ArrowDown />
              </el-icon>
              <span v-else class="tree-toggle-placeholder"></span>
              <el-icon v-if="isMilestone(task)" class="milestone-icon"><Star /></el-icon>
              <div class="task-name-text" :title="task.name">{{ task.name }}</div>
            </div>
            <div class="task-progress-mini">
              <el-progress
                :percentage="task.progress || 0"
                :stroke-width="4"
                :show-text="false"
              />
            </div>
          </template>
        </div>

        <!-- 工期 -->
        <div
          class="row-column column-duration"
          @dblclick="startEdit(task, 'duration', $event)"
          :class="{ 'is-editing': editingCell?.taskId === task.id && editingCell?.field === 'duration' }"
        >
          <template v-if="editingCell?.taskId === task.id && editingCell?.field === 'duration'">
            <el-input
              v-model.number="editingCell.value"
              size="small"
              type="number"
              @blur="saveEdit"
              @keyup.enter="saveEdit"
              @keyup.esc="cancelEdit"
              autofocus
            />
          </template>
          <template v-else>
            {{ getTaskDuration(task) }} 天
          </template>
        </div>

        <!-- 起止时间 -->
        <div
          class="row-column column-dates"
          @dblclick="startEdit(task, 'dates', $event)"
          :class="{ 'is-editing': editingCell?.taskId === task.id && editingCell?.field === 'dates' }"
        >
          <template v-if="editingCell?.taskId === task.id && editingCell?.field === 'dates'">
            <el-date-picker
              v-model="editingCell.value"
              type="daterange"
              size="small"
              format="YYYY/MM/DD"
              value-format="YYYY-MM-DD"
              @blur="saveEdit"
              @change="saveEdit"
              autofocus
            />
          </template>
          <template v-else>
            <div class="date-range">
              <span class="date-start">{{ task.start ? formatDateShort(task.start) : '-' }}</span>
              <span class="date-separator">→</span>
              <span class="date-end">{{ task.end ? formatDateShort(task.end) : '-' }}</span>
            </div>
          </template>
        </div>

        <!-- 资源 -->
        <div class="row-column column-resources">
          <div v-if="task.resources && task.resources.length > 0" class="resource-tags">
            <el-tag
              v-for="(res, idx) in task.resources.slice(0, 2)"
              :key="idx"
              size="small"
              :type="getResourceTagType(res.type)"
            >
              {{ res.resource_name || res.name }}{{ res.quantity ? `×${res.quantity}` : '' }}
            </el-tag>
            <el-tag v-if="task.resources.length > 2" size="small" type="info">
              +{{ task.resources.length - 2 }}
            </el-tag>
          </div>
          <span v-else class="no-resources">-</span>
        </div>
      </div>
    </template>

    <!-- 空白行（填充到至少10行） -->
    <template v-if="!groupMode">
      <div
        v-for="index in emptyRowCount"
        :key="'empty-' + index"
        class="table-row table-row-empty"
        :style="{ height: rowHeight + 'px' }"
        @click="handleEmptyRowClick"
        @contextmenu="handleEmptyRowContextMenu"
        @dblclick="handleEmptyRowDblClick"
        title="点击或双击添加新任务"
      >
        <div class="row-column column-name">
          <div class="task-tree-indent">
            <span class="tree-toggle-placeholder"></span>
            <el-icon class="add-task-icon"><Plus /></el-icon>
            <div class="task-name-text is-placeholder">点击添加新任务</div>
          </div>
        </div>
        <div class="row-column column-duration"></div>
        <div class="row-column column-start"></div>
        <div class="row-column column-end"></div>
        <div class="row-column column-progress"></div>
      </div>
    </template>

    <!-- 空状态 -->
    <div v-if="visibleTasks.length === 0 && props.tasks.length > 0" class="table-empty">
      <el-empty description="所有任务都已折叠" />
    </div>
    <div v-else-if="props.tasks.length === 0" class="table-empty">
      <el-empty :description="emptyDescription" />
    </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, nextTick } from 'vue'
import { Star, ArrowDown, Rank, Plus } from '@element-plus/icons-vue'
import { diffDays, formatDate } from '@/utils/dateFormat'
import { isMilestone } from '@/utils/ganttHelpers'
import { ganttStore } from '@/stores/ganttStore'

const props = defineProps({
  tasks: {
    type: Array,
    default: () => []
  },
  groupedTasks: {
    type: Array,
    default: () => []
  },
  selectedTaskId: {
    type: [Number, String],
    default: null
  },
  rowHeight: {
    type: Number,
    default: 60
  },
  showCriticalPath: {
    type: Boolean,
    default: true
  },
  groupMode: {
    type: String,
    default: ''
  },
  collapsedGroups: {
    type: Set,
    default: () => new Set()
  },
  emptyDescription: {
    type: String,
    default: '暂无进度计划数据'
  },
  isCollapsed: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits([
  'row-click',
  'row-dblclick',
  'toggle-group',
  'context-menu',
  'cell-edit',
  'task-dragged'
])

// 编辑状态
const editingCell = ref(null)
const editInput = ref(null)

// 拖拽状态
const draggedTask = ref(null)
const dragOverTaskId = ref(null)
const dropPosition = ref(null) // 'before', 'after', or 'child'

// 使用 store 中的折叠状态
const { state, actions } = ganttStore

// 获取任务层级（深度）
const getTaskDepth = (task) => {
  let depth = 0
  let currentTask = task
  while (currentTask.parent_id) {
    depth++
    const parent = props.tasks.find(t => t.id === currentTask.parent_id)
    if (!parent) break
    currentTask = parent
  }
  return depth
}

// 获取任务的子任务
const getTaskChildren = (taskId) => {
  return props.tasks.filter(t => t.parent_id === taskId)
}

// 检查任务是否有子任务
const hasChildren = (task) => {
  return props.tasks.some(t => t.parent_id === task.id)
}

// 切换任务折叠状态
const toggleTaskCollapse = (task) => {
  actions.toggleTaskCollapse(task.id)
}

// 检查任务是否被折叠
const isTaskCollapsed = (task) => {
  return actions.isTaskCollapsed(task.id)
}

// 检查任务是否应该被隐藏（因为父任务被折叠）
const isTaskHidden = (task) => {
  return actions.isTaskHidden(task.id, props.tasks)
}

// 获取应该显示的任务列表（考虑折叠状态）
const visibleTasks = computed(() => {
  return props.tasks.filter(task => !isTaskHidden(task))
})

// 开始编辑
const startEdit = (task, field, event) => {
  event.stopPropagation()

  let initialValue
  if (field === 'name') {
    initialValue = task.name
  } else if (field === 'duration') {
    initialValue = getTaskDuration(task)
  } else if (field === 'dates') {
    initialValue = [task.start, task.end]
  }

  editingCell.value = {
    taskId: task.id,
    field,
    value: initialValue,
    originalTask: { ...task }
  }

  // 聚焦输入框
  nextTick(() => {
    if (field === 'name' || field === 'duration') {
      editInput.value?.focus?.()
      editInput.value?.select?.()
    }
  })
}

// 保存编辑
const saveEdit = async () => {
  if (!editingCell.value) return

  const { taskId, field, value, originalTask } = editingCell.value

  try {
    let updateData = {}

    if (field === 'name') {
      updateData = {
        name: value
      }
    } else if (field === 'duration') {
      // 根据工期计算新的结束日期
      const startDate = new Date(originalTask.start)
      const endDate = new Date(startDate)
      endDate.setDate(startDate.getDate() + value)
      updateData = {
        start_date: formatDate(startDate),
        end_date: formatDate(endDate)
      }
    } else if (field === 'dates') {
      if (value && value.length === 2) {
        updateData = {
          start_date: value[0],
          end_date: value[1]
        }
      }
    }

    emit('cell-edit', { taskId, updateData })
  } catch (error) {
    console.error('保存编辑失败:', error)
  } finally {
    editingCell.value = null
  }
}

// 取消编辑
const cancelEdit = () => {
  editingCell.value = null
}

// 处理右键菜单（已有任务）
const handleContextMenu = (event, task) => {
  // 如果正在编辑或拖拽，不显示右键菜单
  if (editingCell.value || draggedTask.value) return

  event.preventDefault()
  event.stopPropagation()
  // 传递任务，让父组件判断是添加子任务还是其他操作
  emit('context-menu', { event, task, type: 'task' })
}

// 空白行点击 - 创建新任务
const handleEmptyRowClick = () => {
  emit('context-menu', {
    type: 'new-task',
    action: 'create'
  })
}

// 空白行右键 - 显示新建任务菜单
const handleEmptyRowContextMenu = (event) => {
  event.preventDefault()
  event.stopPropagation()
  emit('context-menu', {
    event,
    type: 'new-task',
    action: 'context-menu'
  })
}

// 空白行双击 - 直接创建新任务
const handleEmptyRowDblClick = () => {
  emit('context-menu', {
    type: 'new-task',
    action: 'create-immediate'
  })
}

// ==================== 拖拽处理 ====================
// 开始拖拽
const handleDragStart = (event, task) => {
  // 如果正在编辑，不允许拖拽
  if (editingCell.value) {
    event.preventDefault()
    return
  }

  draggedTask.value = task
  dragOverTaskId.value = null

  event.dataTransfer.effectAllowed = 'move'
  event.dataTransfer.setData('text/plain', task.id.toString())

  // 添加键盘事件监听（Ctrl/Cmd 切换模式）
  document.addEventListener('keydown', handleDragKeyDown)
  document.addEventListener('keyup', handleDragKeyUp)

  console.log('开始拖拽任务:', task.name)
}

// 拖拽时处理按键
const handleDragKeyDown = (event) => {
  if ((event.ctrlKey || event.metaKey) && dragOverTaskId.value) {
    dropPosition.value = 'child'
  }
}

const handleDragKeyUp = (event) => {
  if (dragOverTaskId.value && !event.ctrlKey && !event.metaKey) {
    // 释放 Ctrl 后恢复为 before/after 模式
    const targetRow = document.querySelector(`[data-task-id="${dragOverTaskId.value}"]`)
    if (targetRow) {
      const rect = targetRow.getBoundingClientRect()
      const relativeY = event.clientY - rect.top
      dropPosition.value = relativeY < rect.height / 2 ? 'before' : 'after'
    }
  }
}

// 拖拽结束（无论是否成功放置）
const handleDragEnd = () => {
  // 清理状态
  draggedTask.value = null
  dragOverTaskId.value = null
  dropPosition.value = null

  // 移除键盘事件监听
  document.removeEventListener('keydown', handleDragKeyDown)
  document.removeEventListener('keyup', handleDragKeyUp)
}

// 拖拽经过
const handleDragOver = (event, task) => {
  event.preventDefault()

  // 不能拖到自己身上
  if (draggedTask.value?.id === task.id) return

  dragOverTaskId.value = task.id

  // 方案2：使用键盘修饰键区分行为
  // 按住 Ctrl/Cmd → 成为子任务
  // 直接拖拽 → 插入到任务前/后（调整顺序）
  if (event.ctrlKey || event.metaKey) {
    // 按住 Ctrl/Cmd：强制成为子任务
    dropPosition.value = 'child'
  } else {
    // 直接拖拽：根据垂直位置判断 before/after
    const targetRow = event.currentTarget
    const rect = targetRow.getBoundingClientRect()
    const relativeY = event.clientY - rect.top
    const rowHeight = rect.height

    // 上半部分：插入之前
    // 下半部分：插入之后
    dropPosition.value = relativeY < rowHeight / 2 ? 'before' : 'after'
  }
}

// 拖拽离开
const handleDragLeave = (event, task) => {
  if (dragOverTaskId.value === task.id) {
    dragOverTaskId.value = null
    dropPosition.value = null
  }
}

// 放置
const handleDrop = async (event, targetTask) => {
  event.preventDefault()
  event.stopPropagation()

  if (!draggedTask.value) return
  if (draggedTask.value.id === targetTask.id) {
    // 拖到自己身上，不做任何操作
    draggedTask.value = null
    dragOverTaskId.value = null
    dropPosition.value = null
    return
  }

  console.log('拖拽完成:', {
    from: draggedTask.value.name,
    to: targetTask.name,
    position: dropPosition.value
  })

  // 发射拖拽事件到父组件处理
  emit('task-dragged', {
    fromTask: draggedTask.value,
    toTask: targetTask,
    position: dropPosition.value || 'child'
  })

  // 重置状态
  draggedTask.value = null
  dragOverTaskId.value = null
  dropPosition.value = null

  // 移除键盘事件监听
  document.removeEventListener('keydown', handleDragKeyDown)
  document.removeEventListener('keyup', handleDragKeyUp)
}

// 拖拽到根级别（脱离父任务）
const handleDropToRoot = (event) => {
  if (!draggedTask.value) return

  // 检查是否拖拽到了任务行上
  const targetRow = event.target.closest('.table-row')
  if (targetRow) {
    // 如果是拖到任务行上，让任务行的 drop 处理器处理
    return
  }

  event.preventDefault()
  event.stopPropagation()

  console.log('拖拽到根级别:', draggedTask.value.name)

  // 发射拖拽事件，toTask 为 null 表示根级别
  emit('task-dragged', {
    fromTask: draggedTask.value,
    toTask: null  // null 表示移动到根级别
  })

  // 重置状态
  draggedTask.value = null
  dragOverTaskId.value = null

  // 移除键盘事件监听
  document.removeEventListener('keydown', handleDragKeyDown)
  document.removeEventListener('keyup', handleDragKeyUp)
}

// 计算总高度（用于同步滚动，至少显示10行）
const totalHeight = computed(() => {
  let count = 0
  const minRows = 10 // 最小显示行数

  if (props.groupedTasks.length > 0 && props.groupMode) {
    // 分组模式：计算分组中可见的任务数量
    props.groupedTasks.forEach(group => {
      if (!props.collapsedGroups.has(group.name)) {
        count += group.tasks.length
      }
    })
  } else {
    // 无分组模式：使用可见任务数量（考虑树形折叠状态）
    count = visibleTasks.value.length
  }

  // 确保至少显示minRows行
  return Math.max(count, minRows) * props.rowHeight
})

// 计算空白行数量
const emptyRowCount = computed(() => {
  let count = 0
  const minRows = 10

  if (props.groupedTasks.length > 0 && props.groupMode) {
    props.groupedTasks.forEach(group => {
      if (!props.collapsedGroups.has(group.name)) {
        count += group.tasks.length
      }
    })
  } else {
    count = visibleTasks.value.length
  }

  return Math.max(0, minRows - count)
})

// 获取任务工期
const getTaskDuration = (task) => {
  try {
    if (!task.start || !task.end) return 0
    return diffDays(task.start, task.end)
  } catch (e) {
    console.error('Error calculating task duration:', e)
    return 0
  }
}

// 格式化日期为短格式 (MM/DD)
const formatDateShort = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  if (isNaN(date.getTime())) return ''
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${month}/${day}`
}

// 获取资源标签类型
const getResourceTagType = (type) => {
  const types = {
    labor: 'primary',
    equipment: 'success',
    material: 'warning'
  }
  return types[type] || 'info'
}
</script>

<style scoped>
.task-table-wrapper {
  width: 550px;
  border-right: 1px solid #dcdfe6;
  flex-shrink: 0;
  position: sticky;
  left: 0;
  z-index: 10;
  background: #fff;
  box-shadow: 2px 0 4px rgba(0, 0, 0, 0.05);
  display: flex;
  flex-direction: column;
  /* 高度由内容决定 */
  transition: width 0.3s ease, opacity 0.3s ease, transform 0.3s ease, min-width 0.3s ease;
}

/* 折叠状态 */
.task-table-wrapper.is-collapsed {
  width: 0 !important;
  min-width: 0 !important;
  overflow: hidden;
  opacity: 0;
  border-right: none;
  box-shadow: none;
}

/* 可滚动内容区域 */
.table-content {
  /* 不独立滚动，由父容器控制 */
  flex-shrink: 0;
}

/* 分组头部 */
.table-group-header {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  background: #f5f7fa;
  border-bottom: 1px solid #dcdfe6;
  cursor: pointer;
  user-select: none;
  transition: background 0.2s;
}

.table-group-header:hover {
  background: #ecf5ff;
}

.group-name {
  font-weight: bold;
  color: #303133;
  margin-right: 12px;
}

.group-count {
  font-size: 12px;
  color: #909399;
  margin-right: auto;
}

.group-toggle {
  transition: transform 0.3s;
}

.group-toggle.is-collapsed {
  transform: rotate(-90deg);
}

/* 任务行 */
.table-row {
  display: flex;
  border-bottom: 1px solid #ebeef5;
  transition: background 0.3s;
  cursor: pointer;
  /* 性能优化 */
  contain: layout style paint;
  content-visibility: auto;
  will-change: background;
}

.table-row:hover {
  background: #f5f7fa;
}

.table-row.is-selected {
  background: #ecf5ff;
}

.table-row.is-critical {
  background: #fef0f0;
}

/* 空白行 */
.table-row-empty {
  border-bottom: 1px dashed #e4e7ed;
  background: #fafafa;
  cursor: pointer;
  transition: all 0.2s;
}

.table-row-empty:hover {
  background: #ecf5ff;
  border-color: #409eff;
}

.table-row-empty:hover .add-task-icon {
  color: #409eff;
  transform: scale(1.1);
}

.table-row-empty:hover .task-name-text.is-placeholder {
  color: #409eff;
}

.table-row-empty .add-task-icon {
  margin-right: 8px;
  color: #c0c4cc;
  transition: all 0.2s;
}

.table-row-empty .task-name-text.is-placeholder {
  color: #909399;
  font-style: normal;
}

.row-column {
  padding: 8px;
  display: flex;
  align-items: center;
  border-right: 1px solid #ebeef5;
  font-size: 12px;
}

.row-column:last-child {
  border-right: none;
}

.row-column.column-name {
  flex: 0 0 200px;
  flex-direction: column;
  justify-content: center;
  gap: 4px;
  padding: 8px 12px;
}

.row-column.column-duration {
  flex: 0 0 70px;
  justify-content: center;
  color: #606266;
  font-weight: 500;
}

.row-column.column-dates {
  flex: 0 0 150px;
  justify-content: center;
}

.row-column.column-resources {
  flex: 1;
  justify-content: flex-start;
  padding-left: 12px;
}

.date-range {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 11px;
  color: #606266;
}

.date-start,
.date-end {
  font-weight: 500;
}

.date-separator {
  color: #c0c4cc;
}

.resource-tags {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
}

.no-resources {
  color: #c0c4cc;
  font-size: 12px;
}

.milestone-icon {
  color: #e6a23c;
  font-size: 14px;
  margin-right: 4px;
}

.task-name-text {
  font-size: 13px;
  color: #303133;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-weight: 500;
  display: flex;
  align-items: center;
}

.task-progress-mini {
  width: 100%;
}

/* 树形结构样式 */
.task-tree-indent {
  display: flex;
  align-items: center;
  width: 100%;
  gap: 4px;
}

.tree-toggle-icon {
  color: #909399;
  font-size: 12px;
  cursor: pointer;
  transition: transform 0.2s, color 0.2s;
  flex-shrink: 0;
  width: 16px;
  height: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.tree-toggle-icon:hover {
  color: #409eff;
}

.tree-toggle-icon.is-collapsed {
  transform: rotate(-90deg);
}

.tree-toggle-placeholder {
  width: 16px;
  flex-shrink: 0;
}

/* 空状态 */
.table-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 300px;
}

/* 编辑状态 */
.row-column.is-editing {
  background: #fff;
  padding: 4px 8px;
}

.row-column.is-editing .el-input,
.row-column.is-editing .el-date-picker {
  width: 100%;
}

.row-column:has(.is-editing) {
  background: #e8f4ff;
}

/* 拖拽状态 */
.table-row.is-dragging-over {
  background: #e8f4ff !important;
}

.table-row.drop-before {
  border-top: 3px solid #409eff;
  background: linear-gradient(to bottom, rgba(64, 158, 255, 0.15) 0%, rgba(64, 158, 255, 0.05) 30%, transparent 100%) !important;
}

.table-row.drop-after {
  border-bottom: 3px solid #409eff;
  background: linear-gradient(to top, rgba(64, 158, 255, 0.15) 0%, rgba(64, 158, 255, 0.05) 30%, transparent 100%) !important;
}

.table-row.drop-child {
  background: rgba(64, 158, 255, 0.15) !important;
  border: 2px dashed #409eff;
}

.table-row[draggable="true"] {
  cursor: grab;
}

.table-row[draggable="true"]:active {
  cursor: grabbing;
}

/* 拖拽到根级别时的视觉反馈 */
.table-content.is-dragging-over-root {
  background: linear-gradient(135deg, rgba(64, 158, 255, 0.05) 0%, rgba(64, 158, 255, 0.02) 100%);
}

/* 拖拽提示浮层 */
.drag-indicator {
  position: fixed;
  top: 80px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 1000;
  background: rgba(0, 0, 0, 0.8);
  color: white;
  padding: 10px 20px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  pointer-events: none;
  animation: fadeIn 0.2s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateX(-50%) translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateX(-50%) translateY(0);
  }
}

.drag-indicator-content {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 14px;
}

.drag-icon {
  font-size: 18px;
}

.drag-text {
  display: flex;
  align-items: center;
  gap: 8px;
}

.keyboard-hint {
  background: rgba(64, 158, 255, 0.2);
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  border: 1px solid rgba(64, 158, 255, 0.4);
}
</style>
