<template>
  <div class="project-schedule-list">
    <!-- 项目和进度计划列表 -->
    <el-table
      v-loading="loading"
      :data="paginatedProjects"
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
              编辑进度计划
            </el-button>



            <!-- 删除进度计划 - 只在已创建进度计划时显示 -->
            <el-button
              v-if="scope.row.has_schedule"
              type="danger"
              size="small"
              @click="handleDeleteSchedule(scope.row)"
            >
              删除进度计划
            </el-button>

            <!-- 创建进度计划 - 只在未创建进度计划时显示 -->
            <el-button
              v-if="!scope.row.has_schedule && authStore.hasPermission('progress_create')"
              type="warning"
              size="small"
              @click="handleCreateSchedule(scope.row)"
            >
              创建进度计划
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

    <!-- 分页 -->
    <div class="pagination-container" v-if="totalItems > 0">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="totalItems"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
      />
    </div>

    <!-- 空状态 -->
    <el-empty v-if="totalItems === 0 && !loading" description="暂无项目数据" />
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

// 分页状态
const currentPage = ref(1)
const pageSize = ref(20)

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

// 展平项目树以进行分页（包括所有层级的项目）
const flattenedProjects = computed(() => {
  const result = []

  const flatten = (projects) => {
    for (const project of projects) {
      result.push(project)
      if (project.children && project.children.length > 0) {
        flatten(project.children)
      }
    }
  }

  flatten(projectsWithSchedule.value)
  return result
})

// 总项目数
const totalItems = computed(() => flattenedProjects.value.length)

// 当前页的项目（保持树形结构）
const paginatedProjects = computed(() => {
  if (pageSize.value >= totalItems.value) {
    return projectsWithSchedule.value
  }

  // 计算当前页应该显示的根项目范围
  let count = 0
  const startIndex = (currentPage.value - 1) * pageSize.value
  const endIndex = startIndex + pageSize.value

  const result = []

  const collectPageProjects = (projects) => {
    for (const project of projects) {
      // 计算该项目及其所有子项目的总数
      const projectAndChildrenCount = countProjectAndChildren(project)

      // 检查是否应该包含此项目
      const shouldInclude = (
        count >= startIndex && count < endIndex || // 当前项目在范围内
        count + projectAndChildrenCount > startIndex && count < endIndex || // 跨越范围开始
        count >= startIndex && count < endIndex // 跨越范围结束
      )

      if (shouldInclude) {
        result.push(project)
      }

      count += projectAndChildrenCount

      // 如果已经达到范围末尾，停止处理
      if (count >= endIndex) {
        return true
      }
    }
    return false
  }

  collectPageProjects(projectsWithSchedule.value)
  return result
})

// 计算项目及其所有子项目的总数
const countProjectAndChildren = (project) => {
  let count = 1
  if (project.children && project.children.length > 0) {
    for (const child of project.children) {
      count += countProjectAndChildren(child)
    }
  }
  return count
}

// 分页事件处理
const handlePageChange = (page) => {
  currentPage.value = page
}

const handleSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
}

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
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.project-name-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  padding: 16px 0;
  background: #fff;
}

:deep(.el-table__body-wrapper) {
  max-height: calc(100vh - 400px);
  overflow-y: auto;
}
</style>
