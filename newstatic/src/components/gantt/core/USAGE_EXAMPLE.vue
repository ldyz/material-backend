<template>
  <div class="gantt-example-page">
    <!-- Page Header -->
    <div class="page-header">
      <h1>项目进度管理</h1>
      <div class="header-actions">
        <el-button @click="handleRefresh">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
        <el-button type="primary" @click="handleSave">
          <el-icon><Check /></el-icon>
          保存
        </el-button>
      </div>
    </div>

    <!-- Gantt Editor -->
    <GanttEditor
      ref="ganttEditorRef"
      :project-id="projectId"
      :project-name="projectName"
      :height="editorHeight"
      @ready="handleEditorReady"
      @task-selected="handleTaskSelected"
      @view-changed="handleViewChanged"
    />

    <!-- Selected Task Detail Panel -->
    <transition name="el-fade-in-linear">
      <div v-if="selectedTask" class="task-detail-panel">
        <div class="detail-header">
          <h3>{{ selectedTask.name }}</h3>
          <el-button text @click="selectedTask = null">
            <el-icon><Close /></el-icon>
          </el-button>
        </div>

        <div class="detail-content">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="任务ID">
              {{ selectedTask.id }}
            </el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="getStatusType(selectedTask.status)">
                {{ getStatusLabel(selectedTask.status) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="开始日期">
              {{ selectedTask.start }}
            </el-descriptions-item>
            <el-descriptions-item label="结束日期">
              {{ selectedTask.end }}
            </el-descriptions-item>
            <el-descriptions-item label="工期">
              {{ selectedTask.duration }} 天
            </el-descriptions-item>
            <el-descriptions-item label="进度">
              <el-progress
                :percentage="selectedTask.progress || 0"
                :color="getProgressColor(selectedTask.progress)"
              />
            </el-descriptions-item>
            <el-descriptions-item label="优先级">
              <el-tag :type="getPriorityType(selectedTask.priority)">
                {{ getPriorityLabel(selectedTask.priority) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="关键路径">
              <el-tag v-if="selectedTask.is_critical" type="danger" size="small">
                关键任务
              </el-tag>
              <span v-else>普通任务</span>
            </el-descriptions-item>
          </el-descriptions>

          <div class="detail-actions">
            <el-button type="primary" @click="handleEditTask">
              <el-icon><Edit /></el-icon>
              编辑任务
            </el-button>
            <el-button @click="handleViewDependencies">
              <el-icon><Connection /></el-icon>
              查看依赖
            </el-button>
            <el-button @click="handleAllocateResources">
              <el-icon><User /></el-icon>
              分配资源
            </el-button>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Refresh,
  Check,
  Close,
  Edit,
  Connection,
  User
} from '@element-plus/icons-vue'
import GanttEditor from './GanttEditor.vue'

/**
 * GanttEditor Usage Example
 *
 * This example demonstrates how to integrate and use the GanttEditor component
 * with task details panel and action buttons.
 */

// Props (usually from route)
const projectId = ref(123)
const projectName = ref('示例项目')

// Refs
const ganttEditorRef = ref(null)
const selectedTask = ref(null)

// Computed
const editorHeight = computed(() => {
  // Calculate available height
  const headerHeight = 80
  const detailPanelHeight = selectedTask.value ? 300 : 0
  const padding = 40
  const windowHeight = window.innerHeight

  return windowHeight - headerHeight - detailPanelHeight - padding
})

// ==================== Event Handlers ====================

/**
 * Handle editor ready
 */
function handleEditorReady() {
  console.log('Gantt editor initialized')
  ElMessage.success('甘特图加载完成')
}

/**
 * Handle task selected
 */
function handleTaskSelected(task) {
  selectedTask.value = task
  console.log('Selected task:', task)
}

/**
 * Handle view changed
 */
function handleViewChanged(mode) {
  console.log('View mode changed to:', mode)
  ElMessage.info(`切换到${getViewModeLabel(mode)}`)
}

/**
 * Handle refresh
 */
async function handleRefresh() {
  try {
    await ganttEditorRef.value?.initialize()
    ElMessage.success('已刷新')
  } catch (error) {
    console.error('Refresh failed:', error)
    ElMessage.error('刷新失败')
  }
}

/**
 * Handle save
 */
async function handleSave() {
  try {
    await ganttEditorRef.value?.save()
    ElMessage.success('保存成功')
  } catch (error) {
    console.error('Save failed:', error)
    ElMessage.error('保存失败')
  }
}

/**
 * Handle edit task
 */
function handleEditTask() {
  if (!selectedTask.value) return

  // Open edit dialog
  ElMessage.info('编辑功能开发中')
  // ganttStore.actions.openEditDialog(selectedTask.value)
}

/**
 * Handle view dependencies
 */
function handleViewDependencies() {
  if (!selectedTask.value) return

  ElMessage.info('依赖关系查看功能开发中')
}

/**
 * Handle allocate resources
 */
function handleAllocateResources() {
  if (!selectedTask.value) return

  ElMessage.info('资源分配功能开发中')
}

// ==================== Helper Functions ====================

function getStatusType(status) {
  const types = {
    completed: 'success',
    in_progress: 'primary',
    not_started: 'info',
    delayed: 'danger'
  }
  return types[status] || 'info'
}

function getStatusLabel(status) {
  const labels = {
    completed: '已完成',
    in_progress: '进行中',
    not_started: '未开始',
    delayed: '延期'
  }
  return labels[status] || status
}

function getPriorityType(priority) {
  const types = {
    high: 'danger',
    medium: 'warning',
    low: 'info'
  }
  return types[priority] || 'info'
}

function getPriorityLabel(priority) {
  const labels = {
    high: '高',
    medium: '中',
    low: '低'
  }
  return labels[priority] || priority
}

function getProgressColor(progress) {
  if (progress >= 80) return '#67c23a'
  if (progress >= 50) return '#409eff'
  if (progress >= 20) return '#e6a23c'
  return '#f56c6c'
}

function getViewModeLabel(mode) {
  const labels = {
    day: '日视图',
    week: '周视图',
    month: '月视图',
    quarter: '季度视图'
  }
  return labels[mode] || mode
}

// ==================== Lifecycle ====================

onMounted(() => {
  // Initialize with data from route or API
  console.log('Component mounted')
})
</script>

<style scoped lang="scss">
.gantt-example-page {
  display: flex;
  flex-direction: column;
  height: 100vh;
  padding: 20px;
  background-color: #f5f7fa;
  gap: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);

  h1 {
    margin: 0;
    font-size: 24px;
    font-weight: 500;
    color: var(--el-text-color-primary, #303133);
  }

  .header-actions {
    display: flex;
    gap: 12px;
  }
}

.task-detail-panel {
  position: fixed;
  bottom: 20px;
  right: 20px;
  width: 600px;
  max-height: 400px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  overflow: hidden;
  z-index: 1000;
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background-color: #f5f7fa;
  border-bottom: 1px solid #e4e7ed;

  h3 {
    margin: 0;
    font-size: 16px;
    font-weight: 500;
    color: var(--el-text-color-primary, #303133);
  }
}

.detail-content {
  padding: 20px;
  max-height: 320px;
  overflow-y: auto;
}

.detail-actions {
  display: flex;
  gap: 12px;
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #e4e7ed;
}

/* Responsive */
@media (max-width: 768px) {
  .gantt-example-page {
    padding: 12px;
    gap: 12px;
  }

  .page-header {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
    padding: 12px 16px;

    h1 {
      font-size: 18px;
    }

    .header-actions {
      justify-content: stretch;

      .el-button {
        flex: 1;
      }
    }
  }

  .task-detail-panel {
    right: 10px;
    bottom: 10px;
    left: 10px;
    width: auto;
    max-height: 50vh;
  }

  .detail-actions {
    flex-direction: column;

    .el-button {
      width: 100%;
    }
  }
}
</style>
