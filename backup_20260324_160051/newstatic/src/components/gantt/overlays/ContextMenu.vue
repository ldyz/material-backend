<template>
  <teleport to="body">
    <transition name="el-zoom-in-top">
      <div
        v-if="visible"
        ref="menuRef"
        class="gantt-context-menu"
        :style="menuStyle"
        @click.stop
      >
        <!-- Task Section -->
        <div v-if="task" class="menu-section">
          <div class="section-header">
            <el-icon class="task-icon"><Document /></el-icon>
            <span class="task-name">{{ task.name }}</span>
          </div>
        </div>

        <!-- Task Operations -->
        <div class="menu-section">
          <div class="section-title">{{ $t('gantt.contextMenu.taskOperations') }}</div>
          <div class="menu-item" @click="handleEdit">
            <el-icon><Edit /></el-icon>
            <span>{{ $t('common.edit') }}</span>
            <span class="shortcut">Enter</span>
          </div>
          <div class="menu-item" @click="handleDuplicate">
            <el-icon><CopyDocument /></el-icon>
            <span>{{ $t('common.duplicate') }}</span>
            <span class="shortcut">Ctrl+D</span>
          </div>
          <div class="menu-item" @click="handleDelete" class="danger">
            <el-icon><Delete /></el-icon>
            <span>{{ $t('common.delete') }}</span>
            <span class="shortcut">Del</span>
          </div>
        </div>

        <!-- Dependencies -->
        <div class="menu-section">
          <div class="section-title">{{ $t('gantt.contextMenu.dependencies') }}</div>
          <div class="menu-item" @click="handleAddDependency">
            <el-icon><Connection /></el-icon>
            <span>{{ $t('gantt.contextMenu.addDependency') }}</span>
            <span class="shortcut">Shift+Drag</span>
          </div>
          <div class="menu-item" @click="handleViewDependencies">
            <el-icon><View /></el-icon>
            <span>{{ $t('gantt.contextMenu.viewDependencies') }}</span>
          </div>
          <div class="menu-item" @click="handleAutoSchedule">
            <el-icon><MagicStick /></el-icon>
            <span>{{ $t('gantt.contextMenu.autoSchedule') }}</span>
          </div>
        </div>

        <!-- Resources -->
        <div class="menu-section">
          <div class="section-title">{{ $t('gantt.contextMenu.resources') }}</div>
          <div class="menu-item" @click="handleAssignResource">
            <el-icon><User /></el-icon>
            <span>{{ $t('gantt.contextMenu.assignResource') }}</span>
          </div>
          <div class="menu-item" @click="handleViewWorkload">
            <el-icon><DataAnalysis /></el-icon>
            <span>{{ $t('gantt.contextMenu.viewWorkload') }}</span>
          </div>
        </div>

        <!-- Templates -->
        <div class="menu-section">
          <div class="section-title">{{ $t('gantt.contextMenu.templates') }}</div>
          <div class="menu-item" @click="handleAddToTemplate">
            <el-icon><FolderAdd /></el-icon>
            <span>{{ $t('gantt.contextMenu.addToTemplate') }}</span>
          </div>
          <div class="menu-item" @click="handleCreateTemplate">
            <el-icon><DocumentAdd /></el-icon>
            <span>{{ $t('gantt.contextMenu.createTemplate') }}</span>
          </div>
        </div>

        <!-- History -->
        <div class="menu-section">
          <div class="section-title">{{ $t('gantt.contextMenu.history') }}</div>
          <div class="menu-item" @click="handleViewHistory">
            <el-icon><Clock /></el-icon>
            <span>{{ $t('gantt.contextMenu.viewHistory') }}</span>
          </div>
          <div class="menu-item" @click="handleRestoreVersion" v-if="hasPreviousVersions">
            <el-icon><RefreshLeft /></el-icon>
            <span>{{ $t('gantt.contextMenu.restoreVersion') }}</span>
          </div>
        </div>

        <!-- Comments -->
        <div class="menu-section" v-if="task">
          <div class="section-title">{{ $t('gantt.contextMenu.collaboration') }}</div>
          <div class="menu-item" @click="handleAddComment">
            <el-icon><ChatDotSquare /></el-icon>
            <span>{{ $t('gantt.contextMenu.addComment') }}</span>
            <el-badge v-if="commentCount > 0" :value="commentCount" class="comment-badge" />
          </div>
          <div class="menu-item" @click="handleViewActivity">
            <el-icon><Notification /></el-icon>
            <span>{{ $t('gantt.contextMenu.viewActivity') }}</span>
          </div>
        </div>

        <!-- Divider -->
        <div class="menu-divider"></div>

        <!-- Task Info -->
        <div class="menu-section info-section">
          <div class="info-row">
            <span class="info-label">{{ $t('gantt.task.start') }}:</span>
            <span class="info-value">{{ formatDate(task?.start) }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">{{ $t('gantt.task.duration') }}:</span>
            <span class="info-value">{{ task?.duration }} {{ $t('common.days') }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">{{ $t('gantt.task.progress') }}:</span>
            <el-progress
              :percentage="task?.progress || 0"
              :stroke-width="4"
              :show-text="false"
              style="width: 60px"
            />
            <span class="info-value">{{ task?.progress || 0 }}%</span>
          </div>
        </div>
      </div>
    </transition>

    <!-- Edit Dialog -->
    <TaskEditDialog
      v-model="editDialogVisible"
      :task="editingTask"
      @save="handleSaveTask"
    />

    <!-- Dependency Dialog -->
    <DependencyDialog
      v-model="dependencyDialogVisible"
      :task="task"
      :tasks="allTasks"
      @save="handleSaveDependency"
    />

    <!-- Resource Assignment Dialog -->
    <ResourceAssignmentDialog
      v-model="resourceDialogVisible"
      :task="task"
      :resources="resources"
      @save="handleSaveResource"
    />
  </teleport>
</template>

<script setup>
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  Document,
  Edit,
  Delete,
  CopyDocument,
  Connection,
  View,
  MagicStick,
  User,
  DataAnalysis,
  FolderAdd,
  DocumentAdd,
  Clock,
  RefreshLeft,
  ChatDotSquare,
  Notification
} from '@element-plus/icons-vue'
import { formatDate } from '@/utils/dateFormat'
import TaskEditDialog from '@/components/gantt/dialogs/TaskEditDialog.vue'
import DependencyDialog from '@/components/gantt/dialogs/DependencyDialog.vue'
import ResourceAssignmentDialog from '@/components/gantt/dialogs/ResourceAssignmentDialog.vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  position: {
    type: Object,
    default: () => ({ x: 0, y: 0 })
  },
  task: {
    type: Object,
    default: null
  },
  tasks: {
    type: Array,
    default: () => []
  },
  resources: {
    type: Array,
    default: () => []
  },
  hasPreviousVersions: {
    type: Boolean,
    default: false
  },
  commentCount: {
    type: Number,
    default: 0
  }
})

const emit = defineEmits([
  'close',
  'edit',
  'duplicate',
  'delete',
  'add-dependency',
  'view-dependencies',
  'auto-schedule',
  'assign-resource',
  'view-workload',
  'add-to-template',
  'create-template',
  'view-history',
  'restore-version',
  'add-comment',
  'view-activity'
])

const { t } = useI18n()

const menuRef = ref(null)
const editDialogVisible = ref(false)
const editingTask = ref(null)
const dependencyDialogVisible = ref(false)
const resourceDialogVisible = ref(false)

// 事件监听器状态跟踪
let clickOutsideListener = null
let keydownListener = null
let isListenersAdded = false

const allTasks = computed(() => props.tasks)

const menuStyle = computed(() => ({
  position: 'fixed',
  left: `${props.position.x}px`,
  top: `${props.position.y}px`,
  zIndex: 9999
}))

/**
 * Handle edit
 */
function handleEdit() {
  editingTask.value = { ...props.task }
  editDialogVisible.value = true
  emit('close')
}

/**
 * Handle save task
 */
function handleSaveTask(updatedTask) {
  emit('edit', updatedTask)
  editDialogVisible.value = false
}

/**
 * Handle duplicate
 */
function handleDuplicate() {
  emit('duplicate', props.task)
  emit('close')
}

/**
 * Handle delete
 */
function handleDelete() {
  emit('delete', props.task)
  emit('close')
}

/**
 * Handle add dependency
 */
function handleAddDependency() {
  dependencyDialogVisible.value = true
  emit('close')
}

/**
 * Handle save dependency
 */
function handleSaveDependency(dependency) {
  emit('add-dependency', dependency)
  dependencyDialogVisible.value = false
}

/**
 * Handle view dependencies
 */
function handleViewDependencies() {
  emit('view-dependencies', props.task)
  emit('close')
}

/**
 * Handle auto schedule
 */
function handleAutoSchedule() {
  emit('auto-schedule', props.task)
  emit('close')
}

/**
 * Handle assign resource
 */
function handleAssignResource() {
  resourceDialogVisible.value = true
  emit('close')
}

/**
 * Handle save resource
 */
function handleSaveResource(assignment) {
  emit('assign-resource', assignment)
  resourceDialogVisible.value = false
}

/**
 * Handle view workload
 */
function handleViewWorkload() {
  emit('view-workload', props.task)
  emit('close')
}

/**
 * Handle add to template
 */
function handleAddToTemplate() {
  emit('add-to-template', props.task)
  emit('close')
}

/**
 * Handle create template
 */
function handleCreateTemplate() {
  emit('create-template', props.task)
  emit('close')
}

/**
 * Handle view history
 */
function handleViewHistory() {
  emit('view-history', props.task)
  emit('close')
}

/**
 * Handle restore version
 */
function handleRestoreVersion() {
  emit('restore-version', props.task)
  emit('close')
}

/**
 * Handle add comment
 */
function handleAddComment() {
  emit('add-comment', props.task)
  emit('close')
}

/**
 * Handle view activity
 */
function handleViewActivity() {
  emit('view-activity', props.task)
  emit('close')
}

/**
 * Close menu on click outside
 */
function handleClickOutside(event) {
  if (menuRef.value && !menuRef.value.contains(event.target)) {
    emit('close')
  }
}

/**
 * Close menu on escape
 */
function handleEscape(event) {
  if (event.key === 'Escape') {
    emit('close')
  }
}

// Watch for visibility changes
watch(() => props.visible, (visible) => {
  if (visible) {
    nextTick(() => {
      // 创建绑定的监听器函数
      clickOutsideListener = (event) => {
        if (menuRef.value && !menuRef.value.contains(event.target)) {
          emit('close')
        }
      }
      keydownListener = (event) => {
        if (event.key === 'Escape') {
          emit('close')
        }
      }

      document.addEventListener('click', clickOutsideListener)
      document.addEventListener('keydown', keydownListener)
      isListenersAdded = true
    })
  } else {
    // 移除监听器
    if (clickOutsideListener && isListenersAdded) {
      document.removeEventListener('click', clickOutsideListener)
      clickOutsideListener = null
    }
    if (keydownListener && isListenersAdded) {
      document.removeEventListener('keydown', keydownListener)
      keydownListener = null
    }
    isListenersAdded = false
  }
})

// Cleanup
onBeforeUnmount(() => {
  // 确保监听器被移除
  if (clickOutsideListener) {
    document.removeEventListener('click', clickOutsideListener)
    clickOutsideListener = null
  }
  if (keydownListener) {
    document.removeEventListener('keydown', keydownListener)
    keydownListener = null
  }
  isListenersAdded = false
})
</script>

<script>
export default {
  name: 'GanttContextMenu'
}
</script>

<style scoped>
.gantt-context-menu {
  background: white;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  min-width: 240px;
  max-width: 320px;
  max-height: 80vh;
  overflow-y: auto;
}

/* Menu Section */
.menu-section {
  padding: 8px 0;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: #f5f7fa;
  border-bottom: 1px solid #e4e7ed;
}

.task-icon {
  font-size: 18px;
  color: #409EFF;
}

.task-name {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.section-title {
  padding: 4px 16px;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: #909399;
}

/* Menu Item */
.menu-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 16px;
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
}

.menu-item:hover {
  background: #f5f7fa;
  color: #409EFF;
}

.menu-item.danger:hover {
  background: #fef0f0;
  color: #F56C6C;
}

.menu-item .el-icon {
  font-size: 16px;
  flex-shrink: 0;
}

.menu-item span {
  font-size: 14px;
  flex: 1;
}

.shortcut {
  font-size: 11px;
  color: #C0C4CC;
  margin-left: auto;
  background: #f5f7fa;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: monospace;
}

.comment-badge {
  margin-left: auto;
}

/* Divider */
.menu-divider {
  height: 1px;
  background: #e4e7ed;
  margin: 0;
}

/* Info Section */
.info-section {
  background: #fafafa;
}

.info-row {
  display: flex;
  align-items: center;
  padding: 8px 16px;
  gap: 12px;
}

.info-label {
  font-size: 12px;
  color: #909399;
  min-width: 70px;
}

.info-value {
  font-size: 13px;
  color: #303133;
  font-weight: 500;
}

/* Scrollbar */
.gantt-context-menu::-webkit-scrollbar {
  width: 6px;
}

.gantt-context-menu::-webkit-scrollbar-track {
  background: #f5f7fa;
}

.gantt-context-menu::-webkit-scrollbar-thumb {
  background: #dcdfe6;
  border-radius: 3px;
}

.gantt-context-menu::-webkit-scrollbar-thumb:hover {
  background: #c0c4cc;
}

/* Animations */
.el-zoom-in-top-enter-from,
.el-zoom-in-top-leave-to {
  opacity: 0;
  transform: scaleY(0.8) translateY(-10px);
}

.el-zoom-in-top-enter-to,
.el-zoom-in-top-leave-from {
  opacity: 1;
  transform: scaleY(1) translateY(0);
}

.el-zoom-in-top-enter-active,
.el-zoom-in-top-leave-active {
  transition: all 0.2s ease;
  transform-origin: top left;
}
</style>
