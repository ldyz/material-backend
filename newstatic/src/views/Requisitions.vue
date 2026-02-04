<template>
  <div class="requisitions-container">
    <el-card shadow="never">
      <!-- 工具栏 -->
      <TableToolbar>
        <template #left>
          <ProjectSelector
            v-model="searchForm.project_id"
            :projects="projectList"
            placeholder="全部项目"
            width="200px"
          />
          <el-input
            v-model="searchForm.keyword"
            placeholder="搜索单号、项目"
            clearable
            style="width: 250px"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-select
            v-model="searchForm.status"
            placeholder="单据状态"
            clearable
            style="width: 150px"
          >
            <el-option label="全部" value="" />
            <el-option label="草稿" value="draft" />
            <el-option label="待审核" value="pending" />
            <el-option label="已审核" value="approved" />
            <el-option label="已发货" value="issued" />
            <el-option label="已拒绝" value="rejected" />
          </el-select>
          <el-button :icon="Refresh" @click="handleReset">重置</el-button>
        </template>
        <template #right>
          <el-button
            type="primary"
            :icon="Plus"
            @click="handleAdd"
            v-if="authStore.hasPermission('requisition_create')"
          >
            创建出库单
          </el-button>
          <el-button
            type="warning"
            :icon="Download"
            @click="handleExport"
            v-if="authStore.hasPermission('requisition_export')"
          >
            导出
          </el-button>
        </template>
      </TableToolbar>

      <!-- 表格 -->
      <el-table
        v-loading="loading"
        :data="tableData"
        border
        stripe
        style="width: 100%"
      >
        <el-table-column prop="requisition_no" label="出库单号" width="160" fixed="left">
          <template #default="scope">
            <el-link type="primary" @click="handleView(scope.row)">
              {{ scope.row.requisition_no || scope.row.requisition_number }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column prop="project_name" label="项目名称" min-width="150" show-overflow-tooltip />
        <el-table-column prop="applicant_name" label="申请人" width="100" />
        <el-table-column prop="requisition_date" label="申请日期" width="110" />
        <el-table-column prop="purpose" label="用途" min-width="120" show-overflow-tooltip />
        <el-table-column prop="urgent" label="紧急" width="80" align="center">
          <template #default="scope">
            <el-tag v-if="scope.row.urgent" type="danger" size="small">紧急</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusTagType(scope.row.status)" size="small">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="150" show-overflow-tooltip />
        <el-table-column label="操作" width="320" fixed="right">
          <template #default="scope">
            <el-button
              type="info"
              size="small"
              :icon="View"
              @click="handleView(scope.row)"
            >
              查看
            </el-button>
            <el-button
              type="primary"
              size="small"
              :icon="Edit"
              @click="handleEdit(scope.row)"
              v-if="canEdit(scope.row)"
            >
              编辑
            </el-button>
            <el-button
              type="success"
              size="small"
              :icon="Check"
              @click="handleApprove(scope.row)"
              v-if="canApprove(scope.row)"
            >
              审核
            </el-button>
            <el-button
              type="warning"
              size="small"
              :icon="Van"
              @click="handleIssue(scope.row)"
              v-if="canIssue(scope.row)"
            >
              发货
            </el-button>
            <el-button
              type="danger"
              size="small"
              :icon="Delete"
              @click="handleDelete(scope.row)"
              v-if="canDelete(scope.row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
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
      width="900px"
      :loading="dialogLoading"
      @confirm="handleSubmit"
    >
      <template #extra v-if="isViewMode">
        <el-button
          type="primary"
          :icon="Printer"
          @click="handlePrint"
          size="small"
        >
          打印
        </el-button>
      </template>

      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
      >
        <!-- 查看模式：显示详细信息 -->
        <template v-if="isViewMode">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="出库单号" :span="2">
              <el-tag type="primary">{{ formData.requisition_no || '-' }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="getStatusTagType(formData.status)" size="small">
                {{ getStatusText(formData.status) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="紧急">
              <el-tag v-if="formData.urgent" type="danger" size="small">紧急</el-tag>
              <span v-else>否</span>
            </el-descriptions-item>
            <el-descriptions-item label="项目名称" :span="2">
              {{ formData.project_name || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="申请人">
              {{ formData.applicant_name || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="部门">
              {{ formData.department || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="申请日期">
              {{ formData.requisition_date || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="用途">
              {{ formData.purpose || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="创建时间">
              {{ formData.created_at || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="审核人" v-if="formData.approved_by">
              {{ formData.approved_by || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="审核时间" v-if="formData.approved_at">
              {{ formData.approved_at || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="发货人" v-if="formData.issued_by">
              {{ formData.issued_by || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="发货时间" v-if="formData.issued_at">
              {{ formData.issued_at || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="备注" :span="2">
              {{ formData.remark || '-' }}
            </el-descriptions-item>
          </el-descriptions>

          <!-- 工作流状态显示 -->
          <div style="margin-top: 20px;">
            <WorkflowStatus
              v-if="formData.status"
              :status="formData.status"
              :status-time="formData.updated_at || formData.created_at"
              :status-description="getStatusDescription(formData.status)"
              workflow-type="requisition"
              @action="handleWorkflowAction"
            />
          </div>

          <!-- 物资明细 -->
          <el-divider content-position="left">物资明细 ({{ formData.items_count || formData.items?.length || 0 }})</el-divider>
          <el-table
            :data="formData.items"
            border
            stripe
            style="width: 100%"
            size="small"
          >
            <!-- 序号列已移除 -->
            <el-table-column prop="material" label="材质" width="100" show-overflow-tooltip />
            <el-table-column prop="material_name" label="物资名称" min-width="150" show-overflow-tooltip />
            <el-table-column prop="specification" label="规格型号" min-width="150" show-overflow-tooltip />
            <el-table-column prop="unit" label="单位" width="80" />
            <el-table-column prop="requested_quantity" label="申请数量" width="100" align="right" />
            <el-table-column prop="approved_quantity" label="批准数量" width="100" align="right" />
          </el-table>

          <!-- 工作流历史记录 -->
          <template v-if="workflowHistories.length > 0">
            <el-divider content-position="left">审批历史</el-divider>
            <WorkflowHistory :histories="workflowHistories" />
          </template>
        </template>

        <!-- 编辑模式：显示表单 -->
        <template v-else>
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="项目" prop="project_id">
                <ProjectSelector
                  v-model="formData.project_id"
                  :projects="projectList"
                  placeholder="请选择项目"
                  style="width: 100%"
                  :loading="projectsLoading"
                />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="申请人" prop="applicant">
                <el-input
                  v-model="formData.applicant"
                  placeholder="请输入申请人"
                  maxlength="50"
                />
              </el-form-item>
            </el-col>
          </el-row>

          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="申请日期" prop="requisition_date">
                <el-date-picker
                  v-model="formData.requisition_date"
                  type="date"
                  placeholder="选择日期"
                  value-format="YYYY-MM-DD"
                  style="width: 100%"
                />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="用途" prop="purpose">
                <el-input
                  v-model="formData.purpose"
                  placeholder="请输入用途"
                  maxlength="200"
                />
              </el-form-item>
            </el-col>
          </el-row>

          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="紧急">
                <el-switch v-model="formData.urgent" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="备注" prop="remark">
                <el-input
                  v-model="formData.remark"
                  type="textarea"
                  :rows="2"
                  placeholder="请输入备注"
                />
              </el-form-item>
            </el-col>
          </el-row>

          <el-form-item label="物资明细" prop="items" required>
            <StockMaterialSelector
              v-model="formData.items"
              :editable="!isViewMode"
              :project-id="formData.project_id"
              @change="handleItemsChange"
            />
          </el-form-item>
        </template>
      </el-form>
    </Dialog>

    <!-- 审核对话框 -->
    <Dialog
      v-model="approveDialogVisible"
      title="审核出库单"
      width="600px"
      :loading="approveDialogLoading"
      @confirm="handleApproveSubmit"
    >
      <el-form
        ref="approveFormRef"
        :model="approveForm"
        :rules="approveFormRules"
        label-width="100px"
      >
        <el-form-item label="审核结果" prop="approved">
          <el-radio-group v-model="approveForm.approved">
            <el-radio :label="true">通过</el-radio>
            <el-radio :label="false">拒绝</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="审核意见" prop="remark">
          <el-input
            v-model="approveForm.remark"
            type="textarea"
            :rows="4"
            placeholder="请输入审核意见"
            maxlength="500"
          />
        </el-form-item>
      </el-form>
    </Dialog>

    <!-- 发货对话框 -->
    <Dialog
      v-model="issueDialogVisible"
      title="出库发货"
      width="600px"
      :loading="issueDialogLoading"
      @confirm="handleIssueSubmit"
    >
      <el-alert
        title="发货提示"
        type="warning"
        :closable="false"
        style="margin-bottom: 20px"
      >
        发货后将自动扣减库存，请确认物资信息无误
      </el-alert>

      <el-form
        ref="issueFormRef"
        :model="issueForm"
        :rules="issueFormRules"
        label-width="100px"
      >
        <el-form-item label="发货日期" prop="issue_date">
          <el-date-picker
            v-model="issueForm.issue_date"
            type="date"
            placeholder="选择日期"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="备注" prop="remark">
          <el-input
            v-model="issueForm.remark"
            type="textarea"
            :rows="3"
            placeholder="请输入备注"
            maxlength="500"
          />
        </el-form-item>
      </el-form>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { requisitionApi, projectApi } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Search,
  Refresh,
  Plus,
  Download,
  Edit,
  Delete,
  View,
  Check,
  Van,
  Printer
} from '@element-plus/icons-vue'
import Dialog from '@/components/common/Dialog.vue'
import TableToolbar from '@/components/common/TableToolbar.vue'
import MaterialSelector from '@/components/common/MaterialSelector.vue'
import StockMaterialSelector from '@/components/common/StockMaterialSelector.vue'
import WorkflowHistory from '@/components/common/WorkflowHistory.vue'
import WorkflowStatus from '@/components/common/WorkflowStatus.vue'
import ProjectSelector from '@/components/common/ProjectSelector.vue'

const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()

// 列表数据
const loading = ref(false)
const tableData = ref([])
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 搜索表单
const searchForm = reactive({
  keyword: '',
  status: '',
  project_id: null
})

// 项目列表（扁平，用于 ProjectSelector）
const projectList = ref([])
// 项目数据是否已加载
const projectsLoaded = ref(false)
// 项目数据加载中
const projectsLoading = ref(false)

// 构建树形结构的辅助函数（用于 collectProjectIds）
const buildTree = (projects) => {
  if (!projects || projects.length === 0) return []
  const projectMap = new Map()
  projects.forEach(project => {
    projectMap.set(project.id, { ...project, children: [] })
  })
  const roots = []
  projects.forEach(project => {
    const node = projectMap.get(project.id)
    if (!project.parent_id) {
      roots.push(node)
    } else {
      const parent = projectMap.get(project.parent_id)
      if (parent) {
        parent.children.push(node)
      } else {
        roots.push(node)
      }
    }
  })
  return roots
}

// 项目树（用于 collectProjectIds）
const projectTree = computed(() => buildTree(projectList.value))

// 对话框
const dialogVisible = ref(false)
const isViewMode = ref(false)
const dialogTitle = computed(() => {
  if (isViewMode.value) return '查看出库单'
  return formData.id ? '编辑出库单' : '创建出库单'
})
const dialogLoading = ref(false)
const formRef = ref(null)

// 表单数据
const formData = reactive({
  id: null,
  requisition_no: '',
  project_id: null,
  project_name: '',
  applicant: '',
  applicant_name: '',
  department: '',
  requisition_date: new Date().toISOString().split('T')[0],
  purpose: '',
  urgent: false,
  remark: '',
  items: [],
  status: '',
  creator_name: '',
  created_at: '',
  updated_at: '',
  approver: '',
  approved_at: '',
  issued_at: '',
  issuer: ''
})

// 项目选项
const projectOptions = ref([])

// 工作流历史记录
const workflowHistories = ref([])

// 表单验证规则
const formRules = {
  project_id: [
    { required: true, message: '请选择项目', trigger: 'change' }
  ],
  applicant: [
    { required: true, message: '请输入申请人', trigger: 'blur' }
  ],
  requisition_date: [
    { required: true, message: '请选择申请日期', trigger: 'change' }
  ]
}

// 审核对话框
const approveDialogVisible = ref(false)
const approveDialogLoading = ref(false)
const approveFormRef = ref(null)
const currentRequisition = ref(null)
const approveForm = reactive({
  approved: true,
  remark: ''
})
const approveFormRules = {
  remark: [
    {
      required: true,
      message: '请输入审核意见',
      trigger: 'blur',
      validator: (rule, value, callback) => {
        // 审核通过时不需要填写原因，拒绝时必须填写
        if (approveForm.approved === false && !value) {
          callback(new Error('拒绝时必须填写审核意见'))
        } else {
          callback()
        }
      }
    }
  ]
}

// 发货对话框
const issueDialogVisible = ref(false)
const issueDialogLoading = ref(false)
const issueFormRef = ref(null)
const issueForm = reactive({
  issue_date: new Date().toISOString().split('T')[0],
  remark: ''
})
const issueFormRules = {
  issue_date: [
    { required: true, message: '请选择发货日期', trigger: 'change' }
  ]
}

// 获取列表数据
// 适配统一响应格式
const fetchData = async () => {
  loading.value = true
  try {
    // 收集项目ID（包含子项目）
    let projectIds = []
    if (searchForm.project_id) {
      projectIds = collectProjectIds(searchForm.project_id, projectTree.value)
    }

    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      search: searchForm.keyword || undefined,
      status: searchForm.status || undefined,
      project_ids: projectIds.length > 0 ? projectIds.join(',') : undefined
    }
    const { data, pagination: pag } = await requisitionApi.getList(params)
    tableData.value = data || []
    pagination.total = pag?.total || 0
  } catch (error) {
    console.error('获取出库单列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 收集项目及其所有子项目的ID
const collectProjectIds = (projectId, projectTree) => {
  const ids = [projectId]

  // 在树形结构中查找项目并收集子项目ID
  const findAndCollectChildren = (nodes) => {
    for (const node of nodes) {
      if (node.id === projectId) {
        // 找到目标项目，递归收集所有子项目ID
        const collectAllChildren = (project) => {
          if (project.children && project.children.length > 0) {
            for (const child of project.children) {
              ids.push(child.id)
              collectAllChildren(child)
            }
          }
        }
        collectAllChildren(node)
        return true
      }
      if (node.children && node.children.length > 0) {
        if (findAndCollectChildren(node.children)) {
          return true
        }
      }
    }
    return false
  }

  findAndCollectChildren(projectTree)
  return ids
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchData()
}

// 重置
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.status = ''
  searchForm.project_id = null
  pagination.page = 1
  fetchData()
}

// 加载项目列表
const fetchProjects = async () => {
  // 如果正在加载，避免重复请求
  if (projectsLoading.value) {
    console.log('[fetchProjects] 项目数据正在加载中，跳过')
    return
  }

  // 如果数据已加载，直接返回
  if (projectsLoaded.value && projectList.value.length > 0) {
    console.log('[fetchProjects] 项目数据已缓存，跳过加载')
    return
  }

  try {
    projectsLoading.value = true
    console.log('[fetchProjects] ===== 开始加载项目数据 =====')
    console.log('[fetchProjects] 当前 projectList.value:', projectList.value)
    console.log('[fetchProjects] 当前 projectList.value.length:', projectList.value.length)

    const response = await projectApi.getList({ pageSize: 1000 })
    console.log('[fetchProjects] API响应:', response)
    console.log('[fetchProjects] response.success:', response?.success)
    console.log('[fetchProjects] response.data:', response?.data)
    console.log('[fetchProjects] response.data 类型:', typeof response?.data)

    // 兼容不同的响应格式
    let projects = []

    if (response && response.success) {
      // 尝试从不同位置获取项目数据
      if (Array.isArray(response.data)) {
        projects = response.data
        console.log('[fetchProjects] ✓ 数据来源: response.data (数组)，长度:', projects.length)
      } else if (response.data && response.data.projects && Array.isArray(response.data.projects)) {
        projects = response.data.projects
        console.log('[fetchProjects] ✓ 数据来源: response.data.projects (数组)，长度:', projects.length)
      } else if (response.projects && Array.isArray(response.projects)) {
        projects = response.projects
        console.log('[fetchProjects] ✓ 数据来源: response.projects (数组)，长度:', projects.length)
      } else {
        console.error('[fetchProjects] ✗ 无法识别的数据格式')
        console.error('[fetchProjects] response.data 的键:', response.data ? Object.keys(response.data) : 'response.data 为空')
      }
    } else {
      console.error('[fetchProjects] ✗ API响应格式错误或success为false')
    }

    console.log('[fetchProjects] 解析后的项目列表:', projects)
    console.log('[fetchProjects] 项目列表长度:', projects.length)

    // 直接赋值，不判断是否为空
    projectList.value = [...projects]
    projectsLoaded.value = true

    console.log('[fetchProjects] ✓ 已赋值给 projectList.value')

    // 验证数据是否已绑定
    setTimeout(() => {
      console.log('[fetchProjects] 验证绑定 - projectList.value.length:', projectList.value.length)
      console.log('[fetchProjects] 验证绑定 - projectList.value:', projectList.value)
      console.log('[fetchProjects] ===== 加载完成 =====')
    }, 100)
  } catch (error) {
    console.error('[fetchProjects] ✗ 获取项目列表失败:', error)
    projectsLoaded.value = false
  } finally {
    projectsLoading.value = false
  }
}

// buildProjectTree 函数已移除，ProjectSelector 组件内置了树形结构构建功能

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
watch(() => searchForm.keyword, () => {
  debouncedSearch()
})

watch(() => searchForm.status, () => {
  pagination.page = 1
  fetchData()
})

watch(() => searchForm.project_id, () => {
  pagination.page = 1
  fetchData()
})

// 新增
const handleAdd = async () => {
  console.log('[handleAdd] 开始创建出库单')
  console.log('[handleAdd] projectsLoaded:', projectsLoaded.value)
  console.log('[handleAdd] projectsLoading:', projectsLoading.value)
  console.log('[handleAdd] projectList.length:', projectList.value.length)

  // 确保项目数据已加载
  if (!projectsLoaded.value) {
    console.log('[handleAdd] 项目数据未加载，开始加载...')
    await fetchProjects()
    console.log('[handleAdd] fetchProjects 完成，projectList.length:', projectList.value.length)
  } else {
    console.log('[handleAdd] 项目数据已加载，跳过')
  }

  resetForm()
  isViewMode.value = false
  dialogVisible.value = true
  console.log('[handleAdd] 对话框已打开')
}

// 编辑
const handleEdit = async (row) => {
  console.log('[handleEdit] 开始编辑出库单')

  // 确保项目数据已加载
  if (!projectsLoaded.value) {
    console.log('[handleEdit] 项目数据未加载，开始加载...')
    await fetchProjects()
  }

  Object.assign(formData, {
    id: row.id,
    requisition_no: row.requisition_no || '',
    project_id: row.project_id,
    project_name: row.project_name || '',
    applicant: row.applicant || '',
    applicant_name: row.applicant_name || '',
    department: row.department || '',
    requisition_date: row.requisition_date,
    purpose: row.purpose || '',
    urgent: row.urgent || false,
    remark: row.remark || '',
    items: row.items || []
  })
  isViewMode.value = false
  dialogVisible.value = true
  fetchWorkflowHistory(row.id)
}

// 查看
const handleView = (row) => {
  Object.assign(formData, {
    id: row.id,
    requisition_no: row.requisition_no || '',
    project_id: row.project_id,
    project_name: row.project_name || '',
    applicant: row.applicant || '',
    applicant_name: row.applicant_name || '',
    department: row.department || '',
    requisition_date: row.requisition_date,
    purpose: row.purpose || '',
    urgent: row.urgent || false,
    remark: row.remark || '',
    items: row.items || [],
    status: row.status,
    creator_name: row.creator_name || '',
    created_at: row.created_at,
    updated_at: row.updated_at,
    approver: row.approver || '',
    approved_at: row.approved_at || '',
    issued_at: row.issued_at || '',
    issuer: row.issuer || ''
  })
  isViewMode.value = true
  dialogVisible.value = true
  fetchWorkflowHistory(row.id)
}

// 删除
const handleDelete = (row) => {
  ElMessageBox.confirm(
    `确定要删除出库单"${row.requisition_no}"吗？`,
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await requisitionApi.delete(row.id)
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

    if (formData.items.length === 0) {
      ElMessage.warning('请添加物资明细')
      return
    }

    dialogLoading.value = true

    const data = {
      project_id: formData.project_id,
      applicant: formData.applicant,
      requisition_date: formData.requisition_date,
      purpose: formData.purpose,
      urgent: formData.urgent,
      remark: formData.remark,
      items: formData.items.map(item => ({
        stock_id: item.stock_id,
        material_id: item.material_id,
        requested_quantity: item.quantity || 0,
        remark: item.remark || ''
      }))
    }

    if (formData.id) {
      await requisitionApi.update(formData.id, data)
      ElMessage.success('更新成功')
    } else {
      await requisitionApi.create(data)
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

// 审核
const handleApprove = (row) => {
  if (!row || !row.id) {
    ElMessage.error('出库单信息不完整')
    return
  }
  currentRequisition.value = row
  approveForm.approved = true
  approveForm.remark = ''
  approveDialogVisible.value = true
}

// 提交审核
const handleApproveSubmit = async () => {
  if (!approveFormRef.value) return

  try {
    await approveFormRef.value.validate()
    approveDialogLoading.value = true

    await requisitionApi.approve(currentRequisition.value.id, approveForm)
    ElMessage.success(approveForm.approved ? '审核通过' : '已拒绝')

    approveDialogVisible.value = false
    fetchData()
  } catch (error) {
    console.error('审核失败:', error)
  } finally {
    approveDialogLoading.value = false
  }
}

// 发货
const handleIssue = (row) => {
  if (!row || !row.id) {
    ElMessage.error('出库单信息不完整')
    return
  }
  currentRequisition.value = row
  issueForm.issue_date = new Date().toISOString().split('T')[0]
  issueForm.remark = ''
  issueDialogVisible.value = true
}

// 提交发货
const handleIssueSubmit = async () => {
  if (!issueFormRef.value) return

  try {
    await issueFormRef.value.validate()

    if (!currentRequisition.value || !currentRequisition.value.id) {
      ElMessage.error('出库单信息不完整，请重试')
      return
    }

    issueDialogLoading.value = true

    await requisitionApi.issue(currentRequisition.value.id, issueForm)
    ElMessage.success('发货成功')

    issueDialogVisible.value = false
    fetchData()
  } catch (error) {
    console.error('发货失败:', error)
    ElMessage.error(error?.message || '发货失败')
  } finally {
    issueDialogLoading.value = false
  }
}

// 导出
const handleExport = async () => {
  try {
    const response = await requisitionApi.export(searchForm)
    const blob = new Blob([response], {
      type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
    })
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `出库单列表_${new Date().getTime()}.xlsx`
    a.click()
    window.URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
  } catch (error) {
    console.error('导出失败:', error)
  }
}

// 物资明细变化
const handleItemsChange = (items) => {
  formData.items = items
}

// 重置表单
const resetForm = () => {
  Object.assign(formData, {
    id: null,
    requisition_no: '',
    project_id: null,
    project_name: '',
    applicant: authStore.displayName,
    applicant_name: authStore.displayName,
    department: '',
    requisition_date: new Date().toISOString().split('T')[0],
    purpose: '',
    urgent: false,
    remark: '',
    items: [],
    status: '',
    creator_name: '',
    created_at: '',
    updated_at: '',
    approver: '',
    approved_at: '',
    issued_at: '',
    issuer: ''
  })
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

// 打印出库单
const handlePrint = () => {
  // 移除已存在的打印内容
  const existingPrintDiv = document.getElementById('print-temp-content')
  if (existingPrintDiv) {
    existingPrintDiv.remove()
  }

  // 生成打印内容
  const statusTexts = {
    draft: '草稿',
    pending: '待审批',
    approved: '已批准',
    issued: '已发放',
    rejected: '已拒绝'
  }
  const statusText = statusTexts[formData.status] || '未知'

  // 格式化日期
  const formatDate = (dateString) => {
    if (!dateString) return ''
    const date = new Date(dateString)
    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit'
    })
  }

  // 打印表头
  const printHeader = `
    <div class="print-header" style="margin-bottom: 20px; min-height: 100px;">
      <div style="text-align: left; margin-bottom: 10px;">
        <img src="/static/images/logo.png" alt="Logo" style="width: 100px; height: auto; display: block;">
      </div>
      <div style="text-align: center;">
        <div class="print-title" style="font-size: 1.6em; font-weight: bold; letter-spacing: 2px;">大庆化建仪表分公司出库单</div>
      </div>
    </div>
    <table class="print-info">
      <tr>
        <td>出库单号：${formData.requisition_no || '-'}</td>
        <td>项目名称：${formData.project_name || '-'}</td>
        <td>申请人：${formData.applicant_name || '-'}</td>
      </tr>
      <tr>
        <td>申请日期：${formData.requisition_date || formatDate(formData.created_at)}</td>
        <td>状态：${statusText}</td>
        <td>${formData.urgent ? '紧急' : ''}</td>
      </tr>
      ${formData.purpose ? `<tr><td colspan="3">用途：${formData.purpose}</td></tr>` : ''}
    </table>
    ${formData.remark ? `<div style="margin-bottom:10px;"><strong>备注说明：</strong>${formData.remark}</div>` : ''}
  `

  // 签字栏
  const printSignatures = `
    <div class="print-signatures" style="display: flex; justify-content: space-between; margin-top: 20mm;">
      <div style="flex: 1; text-align: left;">申请人签字：_______________</div>
      <div style="flex: 1; text-align: center;">部门负责人签字：_______________</div>
      <div style="flex: 1; text-align: right;">库管员签字：_______________</div>
    </div>
  `

  // 生成分页的物资表格
  let pagedHtml = ''
  if (formData.items && formData.items.length > 0) {
    const itemsPerPage = 10 // 每页显示10条物资
    const totalItems = formData.items.length
    const totalPages = Math.ceil(totalItems / itemsPerPage)

    for (let page = 0; page < totalPages; page++) {
      const pageItems = formData.items.slice(page * itemsPerPage, (page + 1) * itemsPerPage)

      pagedHtml += `
        ${printHeader}
        <table>
          <thead>
            <tr>
              <th style="width: 8%;">序号</th>
              <th style="width: 15%;">材质</th>
              <th style="width: 25%;">物资名称</th>
              <th style="width: 20%;">规格型号</th>
              <th style="width: 10%;">单位</th>
              <th style="width: 11%;">申请数量</th>
              <th style="width: 11%;">批准数量</th>
            </tr>
          </thead>
          <tbody>
            ${pageItems.map((item, index) => `
              <tr>
                <td style="text-align: center;">${page * itemsPerPage + index + 1}</td>
                <td>${item.material || '-'}</td>
                <td>${item.material_name || '-'}</td>
                <td>${item.specification || '-'}</td>
                <td>${item.unit || '-'}</td>
                <td style="text-align: right;">${Number(item.requested_quantity || 0).toLocaleString('zh-CN')}</td>
                <td style="text-align: right;">${Number(item.approved_quantity || 0).toLocaleString('zh-CN')}</td>
              </tr>
            `).join('')}
          </tbody>
        </table>
        ${printSignatures}
        ${page < totalPages - 1 ? `<div class="page-break"></div>` : ''}
      `
    }
  } else {
    pagedHtml = printHeader + '<div style="text-align:center;">暂无物资明细</div>' + printSignatures
  }

  const printHtml = pagedHtml

  // 创建打印样式和内容
  const printStyles = `
    <style>
      /* 打印样式 (优化版) */
      @media print {
        /* 1. 重置全局打印设置 */
        @page {
          margin: 0 !important;
          size: A4 portrait;
        }

        /* 2. 隐藏所有非打印内容 */
        body > *:not(#print-temp-content) {
          display: none !important;
        }

        /* 3. 打印容器设置 (关键修复) */
        #print-temp-content {
          all: initial !important;
          display: block !important;
          position: static !important;
          width: 100% !important;
          max-width: 210mm !important;
          margin: 0 auto !important;
          padding: 10mm !important;
          box-sizing: border-box;
          background: white !important;
          color: black !important;
          font-family: "SimSun", Arial, sans-serif !important;
          line-height: 1.5;
          page-break-inside: avoid;
          page-break-after: avoid;
        }

        /* 4. 表格优化 */
        #print-temp-content table {
          width: 100% !important;
          border-collapse: collapse !important;
          margin: 5mm 0 !important;
          page-break-inside: avoid;
        }

        #print-temp-content th,
        #print-temp-content td {
          border: 1px solid #000 !important;
          padding: 2mm 3mm !important;
          text-align: left !important;
          font-size: 12pt !important;
        }

        /* 5. 特定元素优化 */
        .print-title {
          font-size: 16pt !important;
          text-align: center !important;
          margin-bottom: 8mm !important;
        }

        .print-signatures > div {
          display: inline-block;
          width: 55%;
          margin-top: 15mm;
        }

        /* 6. 强制分页控制 */
        .page-break {
          page-break-after: always !important;
          visibility: hidden !important;
        }
      }

      /* 屏幕样式 (不影响打印) */
      #print-temp-content {
        display: none;
      }
    </style>
  `

  // 将样式和内容添加到页面
  const printDiv = document.createElement('div')
  printDiv.id = 'print-temp-content'
  printDiv.innerHTML = printStyles + printHtml
  document.body.appendChild(printDiv)

  // 执行打印
  setTimeout(() => {
    window.print()
    // 打印完成后移除临时元素
    setTimeout(() => {
      document.body.removeChild(printDiv)
    }, 100)
  }, 100)
}

// 格式化货币
const formatCurrency = (value) => {
  return Number(value).toLocaleString('zh-CN', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  })
}

// 获取状态标签类型
const getStatusTagType = (status) => {
  const types = {
    draft: 'info',
    pending: 'warning',
    approved: 'success',
    issued: 'primary',
    rejected: 'danger'
  }
  return types[status] || 'info'
}

// 获取状态文本
const getStatusText = (status) => {
  const texts = {
    draft: '草稿',
    pending: '待审核',
    approved: '已审核',
    issued: '已发货',
    rejected: '已拒绝'
  }
  return texts[status] || status
}

// 获取状态描述
const getStatusDescription = (status) => {
  const descriptions = {
    draft: '出库单为草稿状态，可以编辑或提交审核',
    pending: '等待审核人员进行审核',
    approved: '审核已通过，等待发货',
    issued: '物资已发货出库',
    rejected: '审核未通过，需要修改后重新提交'
  }
  return descriptions[status] || ''
}

// 获取工作流历史
// 适配统一响应格式
const fetchWorkflowHistory = async (id) => {
  try {
    const { data } = await requisitionApi.getWorkflowHistory(id)
    workflowHistories.value = data || []
  } catch (error) {
    console.error('获取工作流历史失败:', error)
    // 如果没有专门的接口，使用审核历史模拟
    workflowHistories.value = generateMockHistory(formData.value)
  }
}

// 生成模拟工作流历史
const generateMockHistory = (requisition) => {
  const histories = []

  // 创建记录
  histories.push({
    action: 'draft',
    operator_name: requisition.applicant_name || '当前用户',
    operator: requisition.applicant || '当前用户',
    department: '项目部',
    remark: '创建出库单',
    description: `创建出库单 ${requisition.requisition_no}`,
    created_at: requisition.created_at
  })

  // 审核记录
  if (requisition.status === 'approved' || requisition.status === 'rejected' || requisition.status === 'issued') {
    histories.push({
      action: 'pending',
      operator_name: '审核员',
      operator: '审核员',
      department: '管理部',
      remark: '',
      description: '提交审核',
      created_at: requisition.created_at
    })

    if (requisition.approver) {
      histories.push({
        action: requisition.status,
        operator_name: requisition.approver || '审核员',
        operator: requisition.approver || '审核员',
        department: '管理部',
        remark: requisition.approve_remark || '',
        description: requisition.status === 'approved' ? '审核通过' : '审核拒绝',
        created_at: requisition.approved_at || requisition.updated_at
      })
    }
  }

  // 发货记录
  if (requisition.status === 'issued' && requisition.issuer) {
    histories.push({
      action: 'approved',
      operator_name: '仓管员',
      operator: '仓管员',
      department: '仓储部',
      remark: '准备发货',
      description: '审核通过，转入发货流程',
      created_at: requisition.approved_at || requisition.updated_at
    })

    histories.push({
      action: 'issued',
      operator_name: requisition.issuer || '仓管员',
      operator: requisition.issuer || '仓管员',
      department: '仓储部',
      remark: requisition.issue_remark || '',
      description: '完成发货',
      created_at: requisition.issued_at || requisition.updated_at
    })
  }

  return histories.sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
}

// 处理工作流操作
const handleWorkflowAction = async (action) => {
  if (action === 'approve') {
    handleApprove(currentRequisition.value)
  } else if (action === 'reject') {
    handleApprove(currentRequisition.value)
  } else if (action === 'issue') {
    handleIssue(currentRequisition.value)
  }
}

// 判断是否可编辑
const canEdit = (row) => {
  if (!authStore.hasPermission('requisition_edit')) return false
  return row.status === 'draft' || row.status === 'rejected'
}

// 判断是否可审核
const canApprove = (row) => {
  if (!authStore.hasPermission('requisition_approve')) return false
  return row.status === 'pending'
}

// 判断是否可发货
const canIssue = (row) => {
  if (!authStore.hasPermission('requisition_issue')) return false
  return row.status === 'approved'
}

// 判断是否可删除
const canDelete = (row) => {
  if (!authStore.hasPermission('requisition_delete')) return false
  return row.status === 'draft'
}

// 监听审核结果变化，重新验证表单
watch(() => approveForm.approved, () => {
  if (approveFormRef.value) {
    approveFormRef.value.validateField('remark')
  }
})

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

// 监听路由参数，自动触发操作
watch(() => route.query, (query) => {
  if (query.action === 'create') {
    // 重置查询参数，避免重复触发
    router.replace({ query: {} })
    // 稍微延迟，等待页面加载完成
    setTimeout(() => {
      handleAdd()
    }, 100)
  }
}, { immediate: true })

onMounted(async () => {
  await fetchData()
  // 不在页面加载时获取项目列表，而是在点击项目选择器时才加载

  // 检查 URL 查询参数中是否有 requisition_no
  const requisitionNo = route.query.requisition_no
  if (requisitionNo) {
    // 查找对应的出库单（领料单）
    const targetRequisition = tableData.value.find(item => item.requisition_no === requisitionNo)
    if (targetRequisition) {
      // 自动打开详情弹窗
      handleView(targetRequisition)
    }
  }
})
</script>

<style scoped>
.requisitions-container {
  padding: 0;
}

.mt-20 {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>

