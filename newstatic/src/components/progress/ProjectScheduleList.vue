<template>
  <div class="project-schedule-list">
    <!-- 项目和进度计划列表 -->
    <el-table
      v-loading="loading"
      :data="projectsWithSchedule"
      border
      stripe
      style="width: 100%"
      :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
      row-key="id"
    >
      <el-table-column prop="name" label="项目名称" min-width="250" fixed="left">
        <template #default="scope">
          <div class="project-name-cell">
            <el-icon
              v-if="scope.row.level > 0"
              :size="14"
              style="margin-right: 4px; color: #909399"
            >
              <component :is="getLevelIcon(scope.row.level)" />
            </el-icon>
            {{ scope.row.name }}
          </div>
        </template>
      </el-table-column>

      <el-table-column prop="status" label="项目状态" width="120">
        <template #default="scope">
          <el-tag
            :type="getProjectStatusTagType(scope.row.status)"
            size="small"
          >
            {{ getProjectStatusText(scope.row.status) }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column prop="progress_percentage" label="项目进度" width="150">
        <template #default="scope">
          <el-progress
            v-if="scope.row.progress_percentage !== undefined"
            :percentage="scope.row.progress_percentage || 0"
            :status="getProgressStatus(scope.row.progress_percentage)"
            :stroke-width="12"
          />
          <span v-else>-</span>
        </template>
      </el-table-column>

      <el-table-column label="进度计划" width="150">
        <template #default="scope">
          <el-tag
            v-if="scope.row.has_schedule"
            type="success"
            size="small"
          >
            已创建
          </el-tag>
          <el-tag
            v-else
            type="info"
            size="small"
          >
            未创建
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column label="任务数" width="100">
        <template #default="scope">
          {{ scope.row.task_count || 0 }}
        </template>
      </el-table-column>

      <el-table-column label="操作" width="380" fixed="right">
        <template #default="scope">
          <!-- 只对有子项目权限的项目显示操作 -->
          <template v-if="!scope.row.children || scope.row.children.length === 0">
            <!-- 查看甘特图 - 只在已创建进度计划时显示 -->
            <el-button
              v-if="scope.row.has_schedule"
              type="primary"
              size="small"
              @click="handleViewGantt(scope.row)"
            >
              甘特图
            </el-button>

            <!-- 查看网络图 - 只在已创建进度计划时显示 -->
            <el-button
              v-if="scope.row.has_schedule"
              type="success"
              size="small"
              @click="handleViewNetwork(scope.row)"
            >
              网络图
            </el-button>

            <!-- 删除进度计划 - 只在已创建进度计划时显示 -->
            <el-button
              v-if="scope.row.has_schedule"
              type="danger"
              size="small"
              @click="handleDeleteSchedule(scope.row)"
            >
              删除计划
            </el-button>

            <!-- 创建进度计划 - 只在未创建进度计划时显示 -->
            <el-button
              v-if="!scope.row.has_schedule && authStore.hasPermission('progress_create')"
              type="warning"
              size="small"
              @click="handleCreateSchedule(scope.row)"
            >
              创建计划
            </el-button>

            <!-- AI生成计划 -->
            <!-- <el-button
              v-if="!scope.row.has_schedule && authStore.hasPermission('progress_create')"
              type="success"
              size="small"
              @click="handleGeneratePlan(scope.row)"
            >
              AI生成
            </el-button> -->
          </template>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import { Folder, FolderOpened, Document } from '@element-plus/icons-vue'

const props = defineProps({
  projects: {
    type: Array,
    default: () => []
  },
  projectSchedules: {
    type: Object,
    default: () => ({})
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['view-gantt', 'view-network', 'create-schedule', 'delete-schedule', 'generate-plan'])

const authStore = useAuthStore()
const router = useRouter()

// 项目和进度计划数据
const projectsWithSchedule = computed(() => {
  const schedules = props.projectSchedules || {}

  return props.projects.map(project => {
    // 确保 project.id 转换为字符串来匹配 schedules 的 key
    const projectKey = String(project.id)
    const schedule = schedules[projectKey]
    // 计算任务数
    const taskCount = schedule?.activities ?
      Object.values(schedule.activities).filter(a => !a.is_dummy).length : 0

    return {
      ...project,
      has_schedule: !!schedule,
      task_count: taskCount
    }
  })
})

// 查看甘特图
const handleViewGantt = (project) => {
  emit('view-gantt', project.id)
}

// 查看网络图
const handleViewNetwork = (project) => {
  emit('view-network', project.id)
}

// 创建进度计划
const handleCreateSchedule = (project) => {
  emit('create-schedule', project.id)
}

// 删除进度计划
const handleDeleteSchedule = (project) => {
  emit('delete-schedule', project.id)
}

// 生成计划
const handleGeneratePlan = (project) => {
  emit('generate-plan', project.id)
}

// 获取项目状态标签类型
const getProjectStatusTagType = (status) => {
  const types = {
    planning: 'info',
    active: 'primary',
    suspended: 'warning',
    completed: 'success',
    cancelled: 'danger'
  }
  return types[status] || 'info'
}

// 获取项目状态文本
const getProjectStatusText = (status) => {
  const texts = {
    planning: '规划中',
    active: '进行中',
    suspended: '已暂停',
    completed: '已完成',
    cancelled: '已取消'
  }
  return texts[status] || status
}

// 获取进度条状态
const getProgressStatus = (progress) => {
  if (progress >= 100) return 'success'
  if (progress >= 80) return 'warning'
  if (progress > 0) return null
  return 'exception'
}

// 获取层级图标
const getLevelIcon = (level) => {
  const icons = {
    0: null,
    1: Folder,
    2: FolderOpened,
    3: Document
  }
  return icons[level] || Document
}
</script>

<style scoped>
.project-schedule-list {
  width: 100%;
}

.project-name-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

:deep(.el-table__body-wrapper) {
  max-height: calc(100vh - 300px);
  overflow-y: auto;
}
</style>
