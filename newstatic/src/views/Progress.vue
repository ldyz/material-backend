<template>
  <div class="progress-container">
    <!-- 面包屑导航 -->
    <el-breadcrumb separator="/" class="breadcrumb-nav" v-if="currentProject">
      <el-breadcrumb-item :to="{ name: 'Projects' }">项目管理</el-breadcrumb-item>
      <el-breadcrumb-item>{{ currentProject.name }}</el-breadcrumb-item>
    </el-breadcrumb>

    <el-card shadow="never">
      <!-- 视图切换器 -->
      <ViewSwitcher v-model="currentView" @change="handleViewChange">
        <template #actions>
          <!-- AI生成计划功能暂时禁用 -->
          <!-- <el-button
            type="success"
            :icon="MagicStick"
            @click="aiGeneratorVisible = true"
            v-if="authStore.hasPermission('progress_create')"
          >
            AI生成计划
          </el-button> -->
          <el-button
            type="warning"
            @click="handleAggregateChildren"
            v-if="hasChildrenProjects"
          >
            聚合子项目计划
          </el-button>
        </template>
      </ViewSwitcher>

      <!-- 工具栏 -->
      <TableToolbar>
        <template #left>
          <el-tree-select
            v-model="searchForm.project_id"
            :data="projectTreeData"
            :props="treeProps"
            placeholder="选择项目（支持层级显示）"
            clearable
            filterable
            check-strictly
            :render-after-expand="false"
            style="width: 300px"
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
          <el-select
            v-model="searchForm.status"
            placeholder="任务状态"
            clearable
            style="width: 150px"
            @change="handleSearch"
          >
            <el-option label="全部" value="" />
            <el-option label="未开始" value="not_started" />
            <el-option label="进行中" value="in_progress" />
            <el-option label="已完成" value="completed" />
            <el-option label="已延期" value="delayed" />
          </el-select>
          <el-button :icon="Refresh" @click="handleReset">重置</el-button>
        </template>
        <template #right>
          <el-button
            type="primary"
            :icon="Plus"
            @click="handleAdd"
            v-if="authStore.hasPermission('progress_create') && (currentView === 'gantt' || currentView === 'network')"
          >
            添加任务
          </el-button>
          <el-button
            type="warning"
            :icon="Download"
            @click="handleExport"
            v-if="authStore.hasPermission('progress_export')"
          >
            导出
          </el-button>
        </template>
      </TableToolbar>

      <!-- 列表视图：显示项目和进度计划 -->
      <div v-show="currentView === 'list'" class="project-list-view">
        <ProjectScheduleList
          :projects="projectTreeData"
          :project-schedules="allProjectSchedules"
          :loading="loading"
          @view-gantt="handleViewGanttFromList"
          @view-network="handleViewNetworkFromList"
          @create-schedule="handleCreateScheduleFromList"
          @generate-plan="handleGeneratePlanFromList"
        />
      </div>

      <!-- 简单甘特图视图 -->
      <div v-show="currentView === 'gantt-simple'" class="gantt-view">
        <el-empty description="甘特图功能开发中，敬请期待" />
      </div>

      <!-- 甘特图视图 -->
      <div v-show="currentView === 'gantt'" class="gantt-view-container">
        <GanttChart
          v-if="scheduleData && Object.keys(scheduleData).length > 0"
          :project-id="currentProjectId"
          :schedule-data="scheduleData"
          @task-updated="handleTaskUpdated"
          @task-selected="handleTaskSelected"
        />

        <el-empty v-if="!scheduleData || Object.keys(scheduleData).length === 0" description="暂无进度计划，请先创建计划任务" />
      </div>

      <!-- 网络图视图 -->
      <div v-show="currentView === 'network'" class="network-view-container">
        <NetworkDiagram
          v-if="scheduleData && Object.keys(scheduleData).length > 0"
          :project-id="currentProjectId"
          :schedule-data="scheduleData"
          @node-selected="handleNodeSelected"
          @position-updated="handlePositionUpdated"
        />
        <el-empty v-else description="暂无进度计划，请先创建计划任务" />
      </div>

      <!-- 旧的甘特图视图（保留作为备用） -->
      <div v-show="currentView === 'gantt-old'" class="gantt-view">
        <div v-loading="loading" class="gantt-chart">
          <div class="gantt-header">
            <div class="gantt-tasks">任务列表</div>
            <div class="gantt-timeline">
              <div
                v-for="month in timelineMonths"
                :key="month"
                class="timeline-month"
              >
                {{ month }}
              </div>
            </div>
          </div>
          <div class="gantt-body">
            <div
              v-for="item in tableData"
              :key="item.id"
              class="gantt-row"
            >
              <div class="gantt-task-name">
                <el-icon
                  v-if="item.children && item.children.length > 0"
                  class="expand-icon"
                >
                  <FolderOpened />
                </el-icon>
                <el-icon v-else class="expand-icon">
                  <Document />
                </el-icon>
                {{ item.task_name }}
              </div>
              <div class="gantt-timeline-bar">
                <div
                  class="task-bar"
                  :style="getTaskBarStyle(item)"
                  :title="`${item.task_name} (${item.start_date} ~ ${item.end_date})`"
                >
                  <span class="task-bar-label">{{ item.progress }}%</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 分页 -->
      <el-pagination
        v-if="currentView === 'list'"
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
        class="mt-20"
      />
    </el-card>

    <!-- 创建/编辑对话框 -->
    <Dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="700px"
      :loading="dialogLoading"
      @confirm="handleSubmit"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="项目" prop="project_id">
              <el-select
                v-model="formData.project_id"
                placeholder="请选择项目"
                filterable
                style="width: 100%"
                :disabled="isViewMode || !!formData.id"
              >
                <el-option
                  v-for="item in projectOptions"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="父任务" prop="parent_id">
              <el-select
                v-model="formData.parent_id"
                placeholder="请选择父任务"
                clearable
                filterable
                style="width: 100%"
                :disabled="isViewMode"
              >
                <el-option
                  v-for="item in parentTaskOptions"
                  :key="item.id"
                  :label="item.task_name"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="任务名称" prop="task_name">
          <el-input
            v-model="formData.task_name"
            placeholder="请输入任务名称"
            maxlength="100"
            :disabled="isViewMode"
          />
        </el-form-item>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="开始日期" prop="start_date">
              <el-date-picker
                v-model="formData.start_date"
                type="date"
                placeholder="选择开始日期"
                value-format="YYYY-MM-DD"
                style="width: 100%"
                :disabled="isViewMode"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="结束日期" prop="end_date">
              <el-date-picker
                v-model="formData.end_date"
                type="date"
                placeholder="选择结束日期"
                value-format="YYYY-MM-DD"
                style="width: 100%"
                :disabled="isViewMode"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="优先级" prop="priority">
              <el-select
                v-model="formData.priority"
                placeholder="请选择优先级"
                style="width: 100%"
                :disabled="isViewMode"
              >
                <el-option label="低" value="low" />
                <el-option label="中" value="medium" />
                <el-option label="高" value="high" />
                <el-option label="紧急" value="urgent" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-select
                v-model="formData.status"
                placeholder="请选择状态"
                style="width: 100%"
                :disabled="isViewMode"
              >
                <el-option label="未开始" value="not_started" />
                <el-option label="进行中" value="in_progress" />
                <el-option label="已完成" value="completed" />
                <el-option label="已延期" value="delayed" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="完成进度" prop="progress">
              <el-input-number
                v-model="formData.progress"
                :min="0"
                :max="100"
                :step="5"
                :disabled="isViewMode"
                style="width: 100%"
              />
              <span style="margin-left: 10px">%</span>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="负责人" prop="responsible">
              <el-input
                v-model="formData.responsible"
                placeholder="请输入负责人"
                maxlength="50"
                :disabled="isViewMode"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="任务描述" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="4"
            placeholder="请输入任务描述"
            maxlength="1000"
            :disabled="isViewMode"
          />
        </el-form-item>
      </el-form>
    </Dialog>

    <!-- AI生成器对话框（暂时禁用） -->
    <!-- <el-dialog
      v-model="aiGeneratorVisible"
      title="AI智能生成进度计划"
      width="80%"
      :close-on-click-modal="false"
    >
      <AIPlanGenerator
        v-if="aiGeneratorVisible && currentProjectId"
        :project-id="currentProjectId"
        @generated="handleAIGenerated"
        @saved="handleAISaved"
      />
    </el-dialog> -->

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
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { progressApi, projectApi } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'
import { formatDate } from '@/utils/dateFormat'
import {
  Search,
  Refresh,
  Plus,
  Download,
  Edit,
  Delete,
  View,
  List,
  FolderOpened,
  Document,
  Folder,
  MagicStick
} from '@element-plus/icons-vue'
import Dialog from '@/components/common/Dialog.vue'
import TableToolbar from '@/components/common/TableToolbar.vue'
import ViewSwitcher from '@/components/progress/ViewSwitcher.vue'
import GanttChart from '@/components/progress/GanttChart.vue'
// GanttChart 现在指向重构版本 GanttChartRefactored.vue
import NetworkDiagram from '@/components/progress/NetworkDiagram.vue'
import ProjectScheduleList from '@/components/progress/ProjectScheduleList.vue'
import CreateScheduleDialog from '@/components/progress/CreateScheduleDialog.vue'
// import AIPlanGenerator from '@/components/progress/AIPlanGenerator.vue'

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
const allProjectSchedules = ref({})

// 搜索表单（需要在函数声明前定义）
const searchForm = reactive({
  project_id: null,
  status: ''
})

// 表单数据（需要在computed前定义）
const formData = reactive({
  id: null,
  project_id: null,
  parent_id: null,
  task_name: '',
  start_date: '',
  end_date: '',
  priority: 'medium',
  status: 'not_started',
  progress: 0,
  responsible: '',
  description: ''
})

// 项目树数据
const projectTreeData = ref([])
const treeProps = {
  children: 'children',
  label: 'name',
  value: 'id'
}

// 项目选项
const projectOptions = ref([])

// 父任务选项
const parentTaskOptions = ref([])

// 对话框
const dialogVisible = ref(false)
const isViewMode = ref(false)
const dialogTitle = computed(() => {
  if (isViewMode.value) return '查看任务'
  return formData.id ? '编辑任务' : '添加任务'
})
const dialogLoading = ref(false)
const formRef = ref(null)

// 表单验证规则
const formRules = {
  project_id: [
    { required: true, message: '请选择项目', trigger: 'change' }
  ],
  task_name: [
    { required: true, message: '请输入任务名称', trigger: 'blur' }
  ],
  start_date: [
    { required: true, message: '请选择开始日期', trigger: 'change' }
  ],
  end_date: [
    { required: true, message: '请选择结束日期', trigger: 'change' }
  ],
  status: [
    { required: true, message: '请选择状态', trigger: 'change' }
  ]
}

// 时间线月份
const timelineMonths = computed(() => {
  const months = []
  const now = new Date()
  for (let i = -2; i <= 6; i++) {
    const date = new Date(now.getFullYear(), now.getMonth() + i, 1)
    months.push(`${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`)
  }
  return months
})

// AI生成器
const aiGeneratorVisible = ref(false)

// 创建计划对话框
const createScheduleVisible = ref(false)
const creatingProjectId = ref(null)

// ==================== 辅助函数 (需要在事件处理函数之前声明) ====================

// 加载所有项目的进度计划状态
const loadAllProjectSchedules = async () => {
  try {
    // 获取所有有进度计划的项目ID
    const schedules = await progressApi.getAllProjectSchedules()
    allProjectSchedules.value = schedules.data || {}
  } catch (error) {
    console.error('获取项目进度计划失败:', error)
  }
}

// 加载进度计划数据
const loadScheduleData = async () => {
  if (!currentProjectId.value) return
  try {
    const response = await progressApi.getProjectSchedule(currentProjectId.value)
    scheduleData.value = response.data || null
  } catch (error) {
    console.error('获取进度计划失败:', error)
    scheduleData.value = null
  }
}

// 加载项目信息
const loadProjectInfo = async (projectId) => {
  try {
    const response = await projectApi.getDetail(projectId)
    currentProject.value = response.data
    // 检查是否有子项目（level < 3 才能有子项目）
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

    // 同时保留原有的扁平列表（用于其他地方）
    projectOptions.value = projects

    // 构建树形结构
    projectTreeData.value = buildProjectTree(projects)
  } catch (error) {
    console.error('获取项目列表失败:', error)
  }
}

// 构建项目树形结构
const buildProjectTree = (projects) => {
  // 创建一个映射用于快速查找
  const projectMap = new Map()
  const rootProjects = []

  // 第一遍：初始化所有项目
  projects.forEach(project => {
    projectMap.set(project.id, {
      ...project,
      children: []
    })
  })

  // 第二遍：构建树形结构
  projects.forEach(project => {
    const node = projectMap.get(project.id)
    if (!project.parent_id) {
      // 没有父项目，作为根节点
      rootProjects.push(node)
    } else {
      // 有父项目，添加到父节点的children中
      const parent = projectMap.get(project.parent_id)
      if (parent) {
        parent.children.push(node)
      }
    }
  })

  // 返回根节点列表
  return rootProjects
}

// 获取父任务列表
const fetchParentTasks = async () => {
  if (!formData.project_id) {
    parentTaskOptions.value = []
    return
  }
  try {
    const { data } = await progressApi.getList({
      project_id: formData.project_id,
      pageSize: 1000
    })
    parentTaskOptions.value = (data || []).filter(item => item.id !== formData.id)
  } catch (error) {
    console.error('获取父任务列表失败:', error)
  }
}

// 重置表单
const resetForm = () => {
  Object.assign(formData, {
    id: null,
    project_id: null,
    parent_id: null,
    task_name: '',
    start_date: '',
    end_date: '',
    priority: 'medium',
    status: 'not_started',
    progress: 0,
    responsible: '',
    description: ''
  })
  if (formRef.value) {
    formRef.value.clearValidate()
  }
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

// 列表视图事件处理
const handleViewGanttFromList = (projectId) => {
  currentProjectId.value = projectId
  searchForm.project_id = projectId
  currentView.value = 'gantt'
  loadScheduleData()
}

const handleViewNetworkFromList = (projectId) => {
  currentProjectId.value = projectId
  searchForm.project_id = projectId
  currentView.value = 'network'
  loadScheduleData()
}

const handleCreateScheduleFromList = (projectId) => {
  currentProjectId.value = projectId
  creatingProjectId.value = projectId
  createScheduleVisible.value = true
}

// 计划创建成功后的处理
const onScheduleCreated = async (projectId) => {
  createScheduleVisible.value = false
  creatingProjectId.value = null

  // 重新加载所有项目进度计划状态
  await loadAllProjectSchedules()

  // 切换到甘特图视图
  currentProjectId.value = projectId
  searchForm.project_id = projectId
  currentView.value = 'gantt'
  await loadScheduleData()

  ElMessage.success('进度计划创建成功，您现在可以在甘特图中管理和编辑任务')
}

const handleGeneratePlanFromList = (projectId) => {
  currentProjectId.value = projectId
  // TODO: 打开AI生成计划对话框
  ElMessage.info('AI生成计划功能开发中')
}

// 项目切换处理
const handleProjectChange = (projectId) => {
  currentProjectId.value = projectId
  if (projectId) {
    // 加载选中项目的信息
    loadProjectInfo(projectId)
  } else {
    currentProject.value = null
    hasChildrenProjects.value = false
  }
  // 重新加载任务数据
  fetchData()
}

// 切换视图
const toggleView = () => {
  currentView.value = currentView.value === 'list' ? 'gantt' : 'list'
}

// 视图切换器变更
const handleViewChange = (view) => {
  console.log('视图切换:', view)
  if (view === 'gantt' || view === 'network') {
    loadScheduleData()
  }
}

// 任务更新
const handleTaskUpdated = (task) => {
  console.log('任务更新:', task)
  // 重新加载进度计划数据
  loadScheduleData()
}

// 任务选择
const handleTaskSelected = (task) => {
  console.log('选中任务:', task)
  // 可以在这里显示任务详情
}

// 节点选择
const handleNodeSelected = (node) => {
  console.log('选中节点:', node)
  // 可以在这里显示节点详情
}

// 位置更新
const handlePositionUpdated = (data) => {
  console.log('位置更新:', data)
  ElMessage.success('位置已保存')
}

// AI生成完成
const handleAIGenerated = (data) => {
  console.log('AI生成完成:', data)
  scheduleData.value = data
}

// AI保存完成
const handleAISaved = (data) => {
  console.log('AI保存完成:', data)
  aiGeneratorVisible.value = false
  scheduleData.value = data
  currentView.value = 'gantt'
  ElMessage.success('计划已保存')
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

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchData()
}

// 防抖搜索
let searchTimer = null
const debouncedSearch = () => {
  if (searchTimer) {
    clearTimeout(searchTimer)
  }
  searchTimer = setTimeout(() => {
    pagination.page = 1
    fetchData()
  }, 300)
}

// 监听搜索条件变化，实现即时搜索
watch(() => searchForm.project_id, () => {
  pagination.page = 1
  fetchData()
})

watch(() => searchForm.status, () => {
  pagination.page = 1
  fetchData()
})

// 重置
const handleReset = () => {
  searchForm.project_id = null
  searchForm.status = ''
  pagination.page = 1
  fetchData()
}

// 新增
const handleAdd = () => {
  resetForm()
  isViewMode.value = false

  // 如果工具栏已选择项目，自动填充项目字段
  if (searchForm.project_id) {
    formData.project_id = searchForm.project_id
  }

  dialogVisible.value = true
  fetchProjects()
}

// 查看
const handleView = (row) => {
  Object.assign(formData, {
    id: row.id,
    project_id: row.project_id,
    parent_id: row.parent_id,
    task_name: row.task_name,
    start_date: row.start_date,
    end_date: row.end_date,
    priority: row.priority,
    status: row.status,
    progress: row.progress || 0,
    responsible: row.responsible || '',
    description: row.description || ''
  })
  isViewMode.value = true
  dialogVisible.value = true
  fetchParentTasks()
}

// 编辑
const handleEdit = (row) => {
  Object.assign(formData, {
    id: row.id,
    project_id: row.project_id,
    parent_id: row.parent_id,
    task_name: row.task_name,
    start_date: row.start_date,
    end_date: row.end_date,
    priority: row.priority,
    status: row.status,
    progress: row.progress || 0,
    responsible: row.responsible || '',
    description: row.description || ''
  })
  isViewMode.value = false
  dialogVisible.value = true
  fetchProjects()
  fetchParentTasks()
}

// 删除
const handleDelete = (row) => {
  ElMessageBox.confirm(
    `确定要删除任务"${row.task_name}"吗？`,
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await progressApi.delete(row.id)
      ElMessage.success('删除成功')
      fetchData()
    } catch (error) {
      console.error('删除失败:', error)
      ElMessage.error(error?.message || '删除失败，请重试')
    }
  }).catch((error) => {
    // 用户取消操作，不需要提示
    if (error !== 'cancel') {
      console.error('操作失败:', error)
    }
  })
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()

    // 日期验证
    if (formData.start_date && formData.end_date) {
      if (new Date(formData.start_date) > new Date(formData.end_date)) {
        ElMessage.error('开始日期不能晚于结束日期')
        return
      }
    }

    dialogLoading.value = true

    const data = {
      project_id: formData.project_id,
      parent_id: formData.parent_id,
      task_name: formData.task_name,
      start_date: formData.start_date,
      end_date: formData.end_date,
      priority: formData.priority,
      status: formData.status,
      progress: formData.progress,
      responsible: formData.responsible,
      description: formData.description
    }

    if (formData.id) {
      await progressApi.update(formData.id, data)
      ElMessage.success('更新成功')
    } else {
      await progressApi.create(data)
      ElMessage.success('创建成功')
    }

    dialogVisible.value = false
    fetchData()
  } catch (error) {
    console.error('提交失败:', error)
  } finally {
    dialogLoading.value = false
  }
}

// 导出
const handleExport = async () => {
  try {
    const response = await progressApi.export(searchForm)
    const blob = new Blob([response], {
      type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
    })
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `进度计划_${new Date().getTime()}.xlsx`
    a.click()
    window.URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
  } catch (error) {
    console.error('导出失败:', error)
  }
}

// 获取任务条样式
const getTaskBarStyle = (item) => {
  const startDate = new Date(item.start_date)
  const endDate = new Date(item.end_date)
  const timelineStart = new Date(timelineMonths.value[0])
  const totalDays = (endDate - startDate) / (1000 * 60 * 60 * 24)
  const offsetDays = (startDate - timelineStart) / (1000 * 60 * 60 * 24)
  const totalTimelineDays = timelineMonths.value.length * 30

  const left = (offsetDays / totalTimelineDays) * 100
  const width = (totalDays / totalTimelineDays) * 100

  const colors = {
    not_started: '#909399',
    in_progress: '#409eff',
    completed: '#67c23a',
    delayed: '#f56c6c'
  }

  return {
    left: `${Math.max(0, left)}%`,
    width: `${Math.min(100, width)}%`,
    backgroundColor: colors[item.status] || '#409eff'
  }
}

// 获取状态标签类型
const getStatusTagType = (status) => {
  const types = {
    not_started: 'info',
    in_progress: 'primary',
    completed: 'success',
    delayed: 'danger'
  }
  return types[status] || 'info'
}

// 获取状态文本
const getStatusText = (status) => {
  const texts = {
    not_started: '未开始',
    in_progress: '进行中',
    completed: '已完成',
    delayed: '已延期'
  }
  return texts[status] || status
}

// 获取优先级标签类型
const getPriorityTagType = (priority) => {
  const types = {
    low: 'info',
    medium: '',
    high: 'warning',
    urgent: 'danger'
  }
  return types[priority] || ''
}

// 获取优先级文本
const getPriorityText = (priority) => {
  const texts = {
    low: '低',
    medium: '中',
    high: '高',
    urgent: '紧急'
  }
  return texts[priority] || priority
}

// 获取进度条状态
const getProgressStatus = (progress, status) => {
  if (status === 'completed') return 'success'
  if (status === 'delayed') return 'exception'
  if (progress >= 100) return 'success'
  if (progress >= 80) return 'warning'
  return null
}

// 判断是否可编辑
const canEdit = (row) => {
  if (!authStore.hasPermission('progress_edit')) return false
  return row.status !== 'completed'
}

// 判断是否可删除
const canDelete = (row) => {
  if (!authStore.hasPermission('progress_delete')) return false
  return row.status === 'not_started'
}

// 分页处理
const handlePageChange = (page) => {
  pagination.page = page
  fetchData()
}

const handleSizeChange = (size) => {
  pagination.pageSize = size
  pagination.page = 1
  fetchData()
}

onMounted(async () => {
  fetchData()
  fetchProjects()
  loadAllProjectSchedules()

  // 从路由查询参数获取projectId
  const projectId = route.query.projectId
  if (projectId) {
    currentProjectId.value = parseInt(projectId)
    searchForm.project_id = parseInt(projectId)
    // 获取项目信息
    try {
      const project = await projectApi.getDetail(projectId)
      currentProject.value = project.data
      // 检查是否有子项目
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
}

.breadcrumb-nav {
  padding: 12px 16px;
  background: #f5f7fa;
  border-radius: 4px;
  margin-bottom: 16px;
}

.mt-20 {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

/* 新视图容器样式 */
.gantt-view-container,
.network-view-container {
  height: 600px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  overflow: hidden;
}

/* 甘特图样式 */
.gantt-view {
  padding: 20px 0;
}

.gantt-chart {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  overflow: hidden;
}

.gantt-header {
  display: flex;
  background: #f5f7fa;
  border-bottom: 1px solid #dcdfe6;
}

.gantt-tasks {
  width: 250px;
  padding: 12px;
  font-weight: bold;
  border-right: 1px solid #dcdfe6;
  flex-shrink: 0;
}

.gantt-timeline {
  flex: 1;
  display: flex;
  overflow-x: auto;
}

.timeline-month {
  min-width: 100px;
  padding: 12px;
  text-align: center;
  border-right: 1px solid #dcdfe6;
  font-size: 14px;
  font-weight: bold;
}

.gantt-body {
  max-height: 600px;
  overflow-y: auto;
}

.gantt-row {
  display: flex;
  border-bottom: 1px solid #ebeef5;
  transition: background 0.3s;
}

.gantt-row:hover {
  background: #f5f7fa;
}

.gantt-task-name {
  width: 250px;
  padding: 12px;
  border-right: 1px solid #dcdfe6;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
}

.expand-icon {
  color: #909399;
}

/* 树形节点样式 */
.tree-node {
  display: flex;
  align-items: center;
  gap: 8px;
}

.tree-node .level-icon {
  font-size: 14px;
  color: #909399;
}

.tree-node .progress-tag {
  margin-left: auto;
  font-size: 12px;
}

.gantt-timeline-bar {
  flex: 1;
  position: relative;
  min-height: 50px;
  padding: 10px 0;
}

.task-bar {
  position: absolute;
  height: 30px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 12px;
  font-weight: bold;
  cursor: pointer;
  transition: all 0.3s;
  top: 50%;
  transform: translateY(-50%);
}

.task-bar:hover {
  filter: brightness(1.1);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.task-bar-label {
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);
}
</style>
