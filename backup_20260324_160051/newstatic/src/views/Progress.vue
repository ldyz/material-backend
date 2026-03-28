<template>
  <div class="progress-container">
    <!-- 主标题栏 -->
    <div class="page-header" v-if="currentView === 'list'">
      <div class="page-title">
        <h1>项目管理</h1>
        <p class="page-subtitle">查看和管理项目进度计划，使用甘特图可视化项目时间线</p>
      </div>
    </div>

    <!-- 项目导航条（仅在网络图视图显示） -->
    <div class="project-nav" v-if="currentProject && currentView === 'network'">
      <el-breadcrumb separator="/">
        <el-breadcrumb-item @click="backToList" class="clickable">
          <el-icon><Back /></el-icon>
          返回项目列表
        </el-breadcrumb-item>
        <el-breadcrumb-item>{{ currentProject.name }}</el-breadcrumb-item>
      </el-breadcrumb>
      <div class="project-actions">
        <el-tag v-if="hasChildrenProjects" type="info" size="large">
          <el-icon><Folder /></el-icon>
          {{ currentProject.level === 1 ? '一级项目' : currentProject.level === 2 ? '二级项目' : '任务级项目' }}
        </el-tag>
        <el-button v-if="hasChildrenProjects" type="warning" @click="handleAggregateChildren">
          <el-icon><Operation /></el-icon>
          聚合子项目计划
        </el-button>
      </div>
    </div>

    <!-- 视图切换器（仅在列表视图显示） -->
    <div class="view-controls" v-if="currentView === 'list'">
      <div class="search-section">
        <el-tree-select
          v-model="searchForm.project_id"
          :data="projectTreeData"
          :props="treeProps"
          placeholder="搜索项目..."
          clearable
          filterable
          check-strictly
          :render-after-expand="false"
          class="project-search"
          @change="handleProjectChange"
        >
          <template #default="{ node, data }">
            <span class="tree-node">
              <el-icon
                v-if="data.level > 0"
                :size="14"
                style="margin-right: 4px; color: #909399"
              >
                <component :is="getLevelIcon(data.level)" />
              </el-icon>
              {{ node.label }}
              <el-tag
                v-if="data.progress_percentage !== undefined"
                size="small"
                :type="data.progress_percentage >= 100 ? 'success' : 'info'"
                style="margin-left: 8px"
              >
                {{ data.progress_percentage }}%
              </el-tag>
            </span>
          </template>
        </el-tree-select>
      </div>
    </div>

    <!-- 列表视图：显示项目和进度计划 -->
    <div v-show="currentView === 'list'" class="project-list-view">
      <ProjectScheduleList
        :projects="projectTreeData"
        :project-schedules="allProjectSchedules"
        :loading="loading"
        @view-gantt="handleViewGanttFromList"
        @view-network="handleViewNetworkFromList"
        @create-schedule="handleCreateScheduleFromList"
        @delete-schedule="handleDeleteScheduleFromList"
        @generate-plan="handleGeneratePlanFromList"
      />
    </div>

    <!-- 甘特图视图 -->
    <div v-if="currentView === 'gantt'" class="gantt-view-container">
      <GanttChart
        v-if="!scheduleLoading && scheduleData && Object.keys(scheduleData).length > 0"
        :key="`gantt-${currentProjectId}-${scheduleData?.updated || Date.now()}`"
        :project-id="currentProjectId"
        :schedule-data="scheduleData"
        :data-version="scheduleVersion"
        @task-updated="handleTaskUpdated"
        @task-selected="handleTaskSelected"
        @back-to-list="backToList"
      />
      <div v-else-if="scheduleLoading" class="loading-schedule">
        <el-skeleton :rows="5" animated />
      </div>
      <div v-else class="empty-schedule">
        <el-empty description="该项目暂无进度计划">
          <template #image>
            <el-icon :size="100" color="#dcdfe6"><Histogram /></el-icon>
          </template>
          <el-button type="primary" @click="handleCreateSchedule">
            <el-icon><Plus /></el-icon>
            创建进度计划
          </el-button>
        </el-empty>
      </div>
    </div>

    <!-- 网络图视图 -->
    <div v-if="currentView === 'network'" class="network-view-container">
      <NetworkView
        v-if="!scheduleLoading && scheduleData && Object.keys(scheduleData).length > 0"
        :key="`network-${currentProjectId}-${scheduleData?.updated || Date.now()}`"
        :project-id="currentProjectId"
        :schedule-data="scheduleData"
        @node-selected="handleNodeSelected"
        @position-updated="handlePositionUpdated"
      />
      <div v-else-if="scheduleLoading" class="loading-schedule">
        <el-skeleton :rows="5" animated />
      </div>
      <div v-else class="empty-schedule">
        <el-empty description="该项目暂无进度计划">
          <template #image>
            <el-icon :size="100" color="#dcdfe6"><Share /></el-icon>
          </template>
          <el-button type="primary" @click="handleCreateSchedule">
            <el-icon><Plus /></el-icon>
            创建进度计划
          </el-button>
        </el-empty>
      </div>
    </div>

    <!-- 创建进度计划对话框 -->
    <CreateScheduleDialog
      v-model="createScheduleVisible"
      :project-id="creatingProjectId"
      :project-name="getProjectName(creatingProjectId)"
      @created="onScheduleCreated"
    />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { progressApi, projectApi } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Back,
  Folder,
  Operation,
  Plus,
  Histogram,
  Share,
  FolderOpened,
  Document
} from '@element-plus/icons-vue'
import { defineAsyncComponent } from 'vue'

// Dynamic imports for large components
const GanttChart = defineAsyncComponent(() => import('@/components/progress/GanttChart.vue'))
const NetworkView = defineAsyncComponent(() => import('@/components/progress/NetworkView.vue'))
const ProjectScheduleList = defineAsyncComponent(() => import('@/components/progress/ProjectScheduleList.vue'))
const CreateScheduleDialog = defineAsyncComponent(() => import('@/components/progress/CreateScheduleDialog.vue'))

const authStore = useAuthStore()
const route = useRoute()

// 列表数据
const loading = ref(false)
const tableData = ref([])
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 视图模式
const currentView = ref('list')

// 项目相关
const currentProjectId = ref(null)
const currentProject = ref(null)
const hasChildrenProjects = ref(false)
const scheduleData = ref(null)
const scheduleLoading = ref(false)
const scheduleVersion = ref(0)  // 添加版本号，用于强制触发更新
const allProjectSchedules = ref({})

// 搜索表单
const searchForm = reactive({
  project_id: null,
  status: ''
})

// 项目树数据
const projectTreeData = ref([])
const treeProps = {
  children: 'children',
  label: 'name',
  value: 'id'
}

// 创建计划对话框
const createScheduleVisible = ref(false)
const creatingProjectId = ref(null)

// ==================== 辅助函数 ====================

// 加载所有项目的进度计划状态
const loadAllProjectSchedules = async () => {
  try {
    const schedules = await progressApi.getAllProjectSchedules()
    allProjectSchedules.value = schedules.data || {}
  } catch (error) {
    console.error('获取项目进度计划失败:', error)
  }
}

// 加载进度计划数据
const loadScheduleData = async () => {
  if (!currentProjectId.value) return
  scheduleLoading.value = true
  try {
    const response = await progressApi.getProjectSchedule(currentProjectId.value)
    const rawData = response.data || null

    // 统一处理和清理数据
    const cleanedData = rawData ? cleanScheduleData(rawData) : null

    // 强制创建新的对象引用，确保 Vue 能检测到变化
    scheduleData.value = cleanedData ? { ...cleanedData } : null

    console.log('Progress - loadScheduleData completed:', {
      hasData: !!cleanedData,
      activitiesCount: cleanedData?.activities ? Object.keys(cleanedData.activities).length : 0,
      timestamp: cleanedData?.updated
    })
  } catch (error) {
    console.error('获取进度计划失败:', error)
    scheduleData.value = null
  } finally {
    scheduleLoading.value = false
  }
}

// 统一处理和清理进度计划数据
const cleanScheduleData = (data) => {
  if (!data || !data.activities) {
    return null
  }

  const cleanedData = {
    ...data,
    activities: {},
    updated: Date.now()
  }

  for (const [key, activity] of Object.entries(data.activities)) {
    if (!activity.earliest_start || activity.earliest_start <= 0 ||
        !activity.earliest_finish || activity.earliest_finish <= 0) {
      console.warn('Progress - 过滤无效活动:', activity.name || key, {
        earliest_start: activity.earliest_start,
        earliest_finish: activity.earliest_finish
      })
      continue
    }

    cleanedData.activities[key] = {
      ...activity,
      task_id: activity.task_id || activity.id,
      duration: activity.duration || 1,
      progress: activity.progress || 0,
      predecessors: Array.isArray(activity.predecessors) ? activity.predecessors : [],
      successors: Array.isArray(activity.successors) ? activity.successors : []
    }
  }

  return cleanedData
}

// 加载项目信息
const loadProjectInfo = async (projectId) => {
  try {
    const response = await projectApi.getDetail(projectId)
    currentProject.value = response.data
    hasChildrenProjects.value = (response.data.level || 0) < 3
  } catch (error) {
    console.error('获取项目信息失败:', error)
  }
}

// 获取列表数据
const fetchData = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      project_id: searchForm.project_id || undefined,
      status: searchForm.status || undefined
    }
    const { data, pagination: pag } = await progressApi.getList(params)
    tableData.value = data || []
    pagination.total = pag?.total || 0
  } catch (error) {
    console.error('获取进度列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 获取项目列表
const fetchProjects = async () => {
  try {
    const { data } = await projectApi.getList({ pageSize: 1000 })
    const projects = data || []
    projectTreeData.value = buildProjectTree(projects)
  } catch (error) {
    console.error('获取项目列表失败:', error)
  }
}

// 构建项目树形结构
const buildProjectTree = (projects) => {
  const projectMap = new Map()
  const rootProjects = []

  projects.forEach(project => {
    projectMap.set(project.id, {
      ...project,
      children: []
    })
  })

  projects.forEach(project => {
    const node = projectMap.get(project.id)
    if (!project.parent_id) {
      rootProjects.push(node)
    } else {
      const parent = projectMap.get(project.parent_id)
      if (parent) {
        parent.children.push(node)
      }
    }
  })

  return rootProjects
}

// 获取项目名称
const getProjectName = (projectId) => {
  if (!projectId) return ''

  const findProject = (projects, id) => {
    for (const project of projects) {
      if (project.id === id) return project
      if (project.children) {
        const found = findProject(project.children, id)
        if (found) return found
      }
    }
    return null
  }

  const project = findProject(projectTreeData.value, projectId)
  return project ? project.name : ''
}

// 获取层级图标
const getLevelIcon = (level) => {
  const levelIcons = {
    0: null,
    1: Folder,
    2: FolderOpened,
    3: Document
  }
  return levelIcons[level] || Document
}

// ==================== 事件处理函数 ====================

// 返回列表
const backToList = () => {
  currentView.value = 'list'
  currentProjectId.value = null
  currentProject.value = null
  scheduleData.value = null
  scheduleLoading.value = false
}

// 创建进度计划
const handleCreateSchedule = () => {
  creatingProjectId.value = currentProjectId.value
  createScheduleVisible.value = true
}

// 列表视图事件处理
const handleViewGanttFromList = (projectId) => {
  currentProjectId.value = projectId
  searchForm.project_id = projectId
  currentView.value = 'gantt'
  scheduleData.value = null
  loadScheduleData()
}

const handleViewNetworkFromList = (projectId) => {
  currentProjectId.value = projectId
  searchForm.project_id = projectId
  currentView.value = 'network'
  scheduleData.value = null
  loadScheduleData()
}

const handleCreateScheduleFromList = (projectId) => {
  currentProjectId.value = projectId
  creatingProjectId.value = projectId
  createScheduleVisible.value = true
}

const handleDeleteScheduleFromList = async (projectId) => {
  const project = projectTreeData.value.find(p => p.id === projectId)
    || projectTreeData.value.flatMap(p => p.children || []).find(c => c.id === projectId)

  const projectName = project?.name || `项目 ${projectId}`

  ElMessageBox.confirm(
    `确定要删除项目"${projectName}"的进度计划吗？此操作将删除该项目的所有任务和依赖关系，且无法恢复。`,
    '删除进度计划',
    {
      confirmButtonText: '确定删除',
      cancelButtonText: '取消',
      type: 'warning',
      distinguishCancelAndClose: true
    }
  ).then(async () => {
    try {
      await progressApi.deleteProjectSchedule(projectId)
      ElMessage.success('进度计划删除成功')
      await fetchProjects()
      await loadAllProjectSchedules()
    } catch (error) {
      console.error('删除进度计划失败:', error)
      const errorMsg = error?.response?.data?.error || error?.response?.data?.message || error?.message || '删除失败，请重试'
      ElMessage.error(errorMsg)
    }
  }).catch((error) => {
    if (error !== 'cancel' && error !== 'close') {
      console.error('操作失败:', error)
    }
  })
}

const onScheduleCreated = async (projectId) => {
  createScheduleVisible.value = false
  creatingProjectId.value = null
  await loadAllProjectSchedules()
  currentProjectId.value = projectId
  searchForm.project_id = projectId
  currentView.value = 'gantt'
  await loadScheduleData()
  ElMessage.success('进度计划创建成功，您现在可以在甘特图中管理和编辑任务')
}

const handleGeneratePlanFromList = (projectId) => {
  currentProjectId.value = projectId
  ElMessage.info('AI生成计划功能开发中')
}

// 项目切换处理
const handleProjectChange = (projectId) => {
  if (projectId) {
    loadProjectInfo(projectId)
  }
  fetchData()
}

// 任务更新
const handleTaskUpdated = (task) => {
  console.log('任务更新:', task)

  // 如果没有传递具体任务，需要完全刷新
  if (!task) {
    loadScheduleData()
    return
  }

  // 如果传递了具体任务，只更新该任务，避免完全刷新导致组件重新渲染
  if (scheduleData.value && scheduleData.value.activities) {
    const activityId = String(task.id)

    // 创建新的activities对象，确保Vue能检测到变化
    const newActivities = {
      ...scheduleData.value.activities
    }

    if (scheduleData.value.activities[activityId]) {
      // 更新现有任务
      newActivities[activityId] = {
        ...scheduleData.value.activities[activityId],
        ...task
      }
      console.log('已更新任务:', activityId, newActivities[activityId])
    } else {
      // 新任务，添加到数据中
      newActivities[activityId] = {
        task_id: task.id,
        duration: task.duration || 1,
        progress: task.progress || 0,
        predecessors: Array.isArray(task.predecessors) ? task.predecessors : [],
        successors: Array.isArray(task.successors) ? task.successors : [],
        ...task
      }
      console.log('已添加新任务:', activityId)
    }

    // 创建新的scheduleData对象，触发响应式更新
    // 保持 updated 字段不变，避免组件重新创建
    scheduleData.value = {
      ...scheduleData.value,
      activities: newActivities,
      updated: scheduleData.value.updated  // 保持原值，不更新 key
    }

    // 递增版本号，强制 GanttChart 组件重新渲染
    scheduleVersion.value++
  }
}

// 任务选择
const handleTaskSelected = (task) => {
  console.log('选中任务:', task)
}

// 节点选择
const handleNodeSelected = (node) => {
  console.log('选中节点:', node)
}

// 位置更新
const handlePositionUpdated = (data) => {
  console.log('位置更新:', data)
  ElMessage.success('位置已保存')
}

// 聚合子项目计划
const handleAggregateChildren = async () => {
  if (!currentProjectId.value) return
  try {
    await ElMessageBox.confirm(
      '确定要聚合所有子项目的进度计划吗？这将覆盖当前项目的进度计划。',
      '确认操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    await progressApi.aggregateChildren(currentProjectId.value)
    ElMessage.success('聚合成功')
    loadScheduleData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('聚合失败:', error)
      ElMessage.error('聚合失败')
    }
  }
}

// 监听搜索条件变化
watch(() => searchForm.project_id, () => {
  pagination.page = 1
  fetchData()
})

watch(() => searchForm.status, () => {
  pagination.page = 1
  fetchData()
})

onMounted(async () => {
  fetchData()
  fetchProjects()
  loadAllProjectSchedules()

  const projectId = route.query.projectId
  if (projectId) {
    currentProjectId.value = parseInt(projectId)
    searchForm.project_id = parseInt(projectId)
    try {
      const project = await projectApi.getDetail(projectId)
      currentProject.value = project.data
      hasChildrenProjects.value = (project.data.level || 0) < 3
    } catch (error) {
      console.error('获取项目信息失败:', error)
    }
  }
})
</script>

<style scoped>
.progress-container {
  padding: 0;
  display: flex;
  flex-direction: column;
  height: 100vh;
  min-height: 0;
  overflow: hidden;
  background: #f5f7fa;
}

/* 页面标题栏 */
.page-header {
  padding: 24px 32px;
  background: #fff;
  border-bottom: 1px solid #ebeef5;
  flex-shrink: 0;
}

.page-title h1 {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin: 0 0 8px 0;
}

.page-subtitle {
  font-size: 14px;
  color: #909399;
  margin: 0;
}

/* 项目导航条 */
.project-nav {
  padding: 16px 24px;
  background: #fff;
  border-bottom: 1px solid #ebeef5;
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-shrink: 0;
}

.project-nav :deep(.el-breadcrumb) {
  flex: 1;
}

.project-nav .clickable {
  cursor: pointer;
  color: #409eff;
}

.project-nav .clickable:hover {
  color: #66b1ff;
}

.project-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.project-actions .el-tag {
  display: flex;
  align-items: center;
  gap: 4px;
}

/* 视图控制区 */
.view-controls {
  padding: 16px 24px;
  background: #fff;
  border-bottom: 1px solid #ebeef5;
  flex-shrink: 0;
}

.search-section {
  display: flex;
  gap: 12px;
  align-items: center;
}

.project-search {
  width: 400px;
}

/* 项目树节点样式 */
.tree-node {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
}

/* 视图容器 */
.project-list-view,
.gantt-view-container,
.network-view-container {
  flex: 1;
  min-height: 0;
  overflow: visible;
  background: #fff;
  margin: 0;
  border-radius: 4px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
}

.project-list-view {
  overflow: visible;
  background: transparent;
  box-shadow: none;
  margin: 0;
  padding: 0 24px;
}

/* 空状态 */
.empty-schedule {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  min-height: 400px;
}

/* 加载状态 */
.loading-schedule {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  min-height: 400px;
  padding: 40px;
}

:deep(.el-icon) {
  vertical-align: middle;
}
</style>
