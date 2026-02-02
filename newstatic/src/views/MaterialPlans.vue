<template>
  <div class="material-plans-container">
    <el-card shadow="never">
      <!-- 工具栏 -->
      <TableToolbar>
        <template #left>
          <ProjectSelector
            v-model="searchForm.project_id"
            :projects="projectList"
            placeholder="选择项目"
            width="250px"
            @change="handleSearch"
          />
          <el-select
            v-model="searchForm.status"
            placeholder="状态筛选"
            clearable
            style="width: 130px"
            @change="handleSearch"
          >
            <el-option label="全部" value="" />
            <el-option label="草稿" value="draft" />
            <el-option label="待审批" value="pending" />
            <el-option label="已批准" value="approved" />
            <el-option label="进行中" value="active" />
            <el-option label="已完成" value="completed" />
            <el-option label="已取消" value="cancelled" />
          </el-select>
          <el-select
            v-model="searchForm.plan_type"
            placeholder="计划类型"
            clearable
            style="width: 130px"
            @change="handleSearch"
          >
            <el-option label="全部" value="" />
            <el-option label="采购计划" value="procurement" />
            <el-option label="使用计划" value="usage" />
            <el-option label="混合计划" value="mixed" />
          </el-select>
          <el-select
            v-model="searchForm.priority"
            placeholder="优先级"
            clearable
            style="width: 120px"
            @change="handleSearch"
          >
            <el-option label="全部" value="" />
            <el-option label="紧急" value="urgent" />
            <el-option label="高" value="high" />
            <el-option label="普通" value="normal" />
            <el-option label="低" value="low" />
          </el-select>
          <el-input
            v-model="searchForm.search"
            placeholder="搜索计划编号、名称"
            clearable
            style="width: 250px"
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-button :icon="Refresh" @click="handleReset">重置</el-button>
        </template>
        <template #right>
          <el-button
            type="primary"
            :icon="Plus"
            @click="handleAdd"
            v-if="authStore.hasPermission('material_plan_create')"
          >
            新建计划
          </el-button>
          <el-button
            type="info"
            :icon="DataAnalysis"
            @click="handleStatistics"
          >
            统计分析
          </el-button>
        </template>
      </TableToolbar>

      <!-- 表格 -->
      <el-table
        v-loading="planStore.loading"
        :data="planStore.plans"
        border
        stripe
        style="width: 100%"
      >
        <el-table-column prop="plan_no" label="计划编号" width="150" fixed="left">
          <template #default="scope">
            <el-link type="primary" @click="handleView(scope.row)">
              {{ scope.row.plan_no }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column prop="plan_name" label="计划名称" min-width="200" show-overflow-tooltip />
        <el-table-column prop="project_name" label="所属项目" width="180" show-overflow-tooltip />
        <el-table-column prop="plan_type" label="计划类型" width="110" align="center">
          <template #default="scope">
            <el-tag :type="getPlanTypeTagType(scope.row.plan_type)" size="small">
              {{ getPlanTypeLabel(scope.row.plan_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="scope">
            <el-tag :type="getStatusTagType(scope.row.status)" size="small">
              {{ getStatusLabel(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="priority" label="优先级" width="90" align="center">
          <template #default="scope">
            <el-tag :type="getPriorityTagType(scope.row.priority)" size="small">
              {{ getPriorityLabel(scope.row.priority) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="items_count" label="项目数" width="80" align="center">
          <template #default="scope">
            {{ scope.row.items_count || 0 }}
          </template>
        </el-table-column>
        <el-table-column label="进度" width="180" align="center">
          <template #default="scope">
            <div v-if="scope.row.progress" class="progress-info">
              <el-progress
                :percentage="Math.round(scope.row.progress.overall_progress || 0)"
                :stroke-width="8"
              />
            </div>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="total_budget" label="总预算" width="120" align="right">
          <template #default="scope">
            ¥{{ formatAmount(scope.row.total_budget) }}
          </template>
        </el-table-column>
        <el-table-column prop="creator_name" label="创建人" width="100" />
        <el-table-column prop="created_at" label="创建时间" width="110">
          <template #default="scope">
            {{ formatDate(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right">
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
              v-if="canEdit(scope.row)"
              type="primary"
              size="small"
              :icon="Edit"
              @click="handleEdit(scope.row)"
            >
              编辑
            </el-button>
            <el-dropdown @command="(cmd) => handleActionCommand(cmd, scope.row)">
              <el-button type="default" size="small">
                更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item
                    v-if="canSubmit(scope.row)"
                    command="submit"
                    :icon="Promotion"
                  >
                    提交审批
                  </el-dropdown-item>
                  <el-dropdown-item
                    v-if="canApprove(scope.row)"
                    command="approve"
                    :icon="Select"
                  >
                    批准
                  </el-dropdown-item>
                  <el-dropdown-item
                    v-if="canApprove(scope.row)"
                    command="reject"
                    :icon="Close"
                  >
                    拒绝
                  </el-dropdown-item>
                  <el-dropdown-item
                    v-if="canResubmit(scope.row)"
                    command="resubmit"
                    :icon="RefreshRight"
                  >
                    重新提交
                  </el-dropdown-item>
                  <el-dropdown-item
                    v-if="canCancel(scope.row)"
                    command="cancel"
                    :icon="Remove"
                  >
                    取消
                  </el-dropdown-item>
                  <el-dropdown-item
                    v-if="canDelete(scope.row)"
                    command="delete"
                    :icon="Delete"
                    divided
                  >
                    删除
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
        v-model:current-page="planStore.pagination.page"
        v-model:page-size="planStore.pagination.page_size"
        :page-sizes="[10, 20, 50, 100]"
        :total="planStore.pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
        class="mt-20"
      />
    </el-card>

    <!-- 新增/编辑对话框 -->
    <Dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="900px"
      :loading="dialogLoading"
      @confirm="handleSubmit"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="110px"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="计划名称" prop="plan_name">
              <el-input
                v-model="formData.plan_name"
                placeholder="请输入计划名称"
                maxlength="200"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="所属项目" prop="project_id">
              <ProjectSelector
                v-model="formData.project_id"
                :projects="projectList"
                placeholder="请选择项目"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="计划类型" prop="plan_type">
              <el-select v-model="formData.plan_type" placeholder="请选择" style="width: 100%">
                <el-option label="采购计划" value="procurement" />
                <el-option label="使用计划" value="usage" />
                <el-option label="混合计划" value="mixed" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="优先级" prop="priority">
              <el-select v-model="formData.priority" placeholder="请选择" style="width: 100%">
                <el-option label="紧急" value="urgent" />
                <el-option label="高" value="high" />
                <el-option label="普通" value="normal" />
                <el-option label="低" value="low" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="总预算(元)" prop="total_budget">
              <el-input-number
                v-model="formData.total_budget"
                :min="0"
                :precision="2"
                :controls="false"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="计划开始日期">
              <el-date-picker
                v-model="formData.planned_start_date"
                type="date"
                placeholder="选择日期"
                format="YYYY-MM-DD"
                value-format="YYYY-MM-DD"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="计划结束日期">
              <el-date-picker
                v-model="formData.planned_end_date"
                type="date"
                placeholder="选择日期"
                format="YYYY-MM-DD"
                value-format="YYYY-MM-DD"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="描述">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="2"
            placeholder="请输入计划描述"
          />
        </el-form-item>

        <el-form-item label="备注">
          <el-input
            v-model="formData.remark"
            type="textarea"
            :rows="2"
            placeholder="请输入备注"
          />
        </el-form-item>

        <!-- 计划项目 -->
        <el-divider content-position="left">计划项目</el-divider>
        <el-button type="primary" size="small" :icon="Plus" @click="handleAddItem" class="mb-10">
          添加项目
        </el-button>
        <el-button type="success" size="small" :icon="Upload" @click="handleImportItems" class="mb-10 ml-10">
          导入项目
        </el-button>
        <el-table
          :data="formData.items"
          border
          size="small"
          max-height="300"
        >
          <el-table-column prop="material_name" label="物资名称" min-width="150" />
          <el-table-column prop="specification" label="规格型号" width="120" />
          <el-table-column prop="category" label="分类" width="100" />
          <el-table-column prop="unit" label="单位" width="70" />
          <el-table-column prop="planned_quantity" label="计划数量" width="90" align="right" />
          <el-table-column prop="unit_price" label="单价" width="90" align="right">
            <template #default="scope">
              {{ scope.row.unit_price ? '¥' + scope.row.unit_price : '-' }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="80" fixed="right">
            <template #default="scope">
              <el-button
                type="danger"
                size="small"
                :icon="Delete"
                link
                @click="handleDeleteItem(scope.$index)"
              >
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-form>
    </Dialog>

    <!-- 项目编辑对话框 -->
    <el-dialog
      v-model="itemDialogVisible"
      title="添加计划项目"
      width="600px"
    >
      <el-form
        ref="itemFormRef"
        :model="itemForm"
        :rules="itemFormRules"
        label-width="100px"
      >
        <el-form-item label="物资名称" prop="material_name">
          <el-input
            v-model="itemForm.material_name"
            placeholder="请输入物资名称"
            maxlength="200"
          />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="物资编码">
              <el-input
                v-model="itemForm.material_code"
                placeholder="请输入物资编码"
                maxlength="50"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="分类">
              <el-input
                v-model="itemForm.category"
                placeholder="请输入分类"
                maxlength="100"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="规格型号">
          <el-input
            v-model="itemForm.specification"
            placeholder="请输入规格型号"
            maxlength="200"
          />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="单位" prop="unit">
              <el-input
                v-model="itemForm.unit"
                placeholder="单位"
                maxlength="20"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="计划数量" prop="planned_quantity">
              <el-input-number
                v-model="itemForm.planned_quantity"
                :min="1"
                :precision="0"
                :controls="false"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="单价(元)">
              <el-input-number
                v-model="itemForm.unit_price"
                :min="0"
                :precision="2"
                :controls="false"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="需求日期">
              <el-date-picker
                v-model="itemForm.required_date"
                type="date"
                placeholder="选择日期"
                format="YYYY-MM-DD"
                value-format="YYYY-MM-DD"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="优先级">
              <el-select v-model="itemForm.priority" placeholder="请选择" style="width: 100%">
                <el-option label="紧急" value="urgent" />
                <el-option label="高" value="high" />
                <el-option label="普通" value="normal" />
                <el-option label="低" value="low" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="备注">
          <el-input
            v-model="itemForm.remark"
            type="textarea"
            :rows="2"
            placeholder="请输入备注"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="itemDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveItem">确定</el-button>
      </template>
    </el-dialog>

    <!-- 导入项目对话框 -->
    <el-dialog
      v-model="importDialogVisible"
      title="导入计划项目"
      width="900px"
      :close-on-click-modal="false"
      draggable
    >
      <ExcelImportMapper
        :fields="planItemFields"
        :projects="[]"
        :show-project-step="false"
        :import-api="handleBatchImportItems"
        @success="handleImportItemsSuccess"
        @close="importDialogVisible = false"
      />
    </el-dialog>

    <!-- 查看详情对话框 -->
    <PlanDetailDialog
      v-model="detailDialogVisible"
      :plan-id="viewingPlanId"
    />
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { usePlanStore } from '@/stores/planStore'
import { projectApi } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Search,
  Refresh,
  Plus,
  View,
  Edit,
  Delete,
  Promotion,
  Select,
  Close,
  Check,
  RefreshRight,
  Remove,
  ArrowDown,
  DataAnalysis,
  Upload
} from '@element-plus/icons-vue'
import Dialog from '@/components/common/Dialog.vue'
import TableToolbar from '@/components/common/TableToolbar.vue'
import PlanDetailDialog from './PlanDetailDialog.vue'
import ExcelImportMapper from '@/components/common/ExcelImportMapper.vue'
import ProjectSelector from '@/components/common/ProjectSelector.vue'

const authStore = useAuthStore()
const planStore = usePlanStore()
const router = useRouter()

// 项目列表
const projectList = ref([])

// 搜索表单
const searchForm = reactive({
  status: '',
  plan_type: '',
  priority: '',
  search: '',
  project_id: null
})

// 对话框
const dialogVisible = ref(false)
const dialogTitle = computed(() => formData.id ? '编辑物资计划' : '新建物资计划')
const dialogLoading = ref(false)
const formRef = ref(null)

// 表单数据
const formData = reactive({
  id: null,
  plan_name: '',
  project_id: null,
  plan_type: 'procurement',
  priority: 'normal',
  planned_start_date: '',
  planned_end_date: '',
  total_budget: null,
  description: '',
  remark: '',
  items: []
})

// 表单验证规则
const formRules = {
  plan_name: [
    { required: true, message: '请输入计划名称', trigger: 'blur' }
  ],
  project_id: [
    { required: true, message: '请选择所属项目', trigger: 'change' }
  ]
}

// 项目对话框
const itemDialogVisible = ref(false)
const itemFormRef = ref(null)
const itemForm = reactive({
  material_name: '',
  material_code: '',
  specification: '',
  category: '',
  unit: '',
  planned_quantity: 1,
  unit_price: null,
  required_date: '',
  priority: 'normal',
  remark: ''
})

const itemFormRules = {
  material_name: [
    { required: true, message: '请输入物资名称', trigger: 'blur' }
  ],
  planned_quantity: [
    { required: true, message: '请输入计划数量', trigger: 'blur' }
  ],
  unit: [
    { required: true, message: '请输入单位', trigger: 'blur' }
  ]
}

let editingItemIndex = -1

// 详情对话框
const detailDialogVisible = ref(false)
const viewingPlanId = ref(null)

// 导入对话框
const importDialogVisible = ref(false)

// 计划项目字段定义
const planItemFields = [
  {
    field: 'material_name',
    label: '物资名称',
    required: true,
    hint: '物资的名称（必填）'
  },
  {
    field: 'material_code',
    label: '物资编码',
    required: false,
    hint: '物资的编码'
  },
  {
    field: 'specification',
    label: '规格型号',
    required: false,
    hint: '规格型号说明'
  },
  {
    field: 'category',
    label: '分类',
    required: false,
    hint: '物资分类'
  },
  {
    field: 'unit',
    label: '单位',
    required: true,
    hint: '计量单位（必填）'
  },
  {
    field: 'planned_quantity',
    label: '计划数量',
    required: true,
    hint: '计划数量'
  },
  {
    field: 'unit_price',
    label: '单价',
    required: false,
    hint: '单价（元）'
  },
  {
    field: 'priority',
    label: '优先级',
    required: false,
    hint: '紧急/高/普通/低'
  },
  {
    field: 'required_date',
    label: '需求日期',
    required: false,
    hint: 'YYYY-MM-DD格式'
  },
  {
    field: 'remark',
    label: '备注',
    required: false,
    hint: '其他说明'
  }
]

// 获取项目列表
const fetchProjects = async () => {
  try {
    const response = await projectApi.getList({ pageSize: 1000 })
    if (response.success) {
      projectList.value = response.data?.projects || response.data || []
    }
  } catch (error) {
    console.error('获取项目列表失败:', error)
  }
}

// 搜索
const handleSearch = () => {
  planStore.setFilters(searchForm)
  planStore.fetchPlans()
}

// 重置
const handleReset = () => {
  Object.assign(searchForm, {
    status: '',
    plan_type: '',
    priority: '',
    search: '',
    project_id: null
  })
  planStore.resetFilters()
  planStore.fetchPlans()
}

// 分页
const handlePageChange = (page) => {
  planStore.setPagination({ page })
  planStore.fetchPlans()
}

const handleSizeChange = (size) => {
  planStore.setPagination({ page_size: size, page: 1 })
  planStore.fetchPlans()
}

// 新增
const handleAdd = () => {
  Object.assign(formData, {
    id: null,
    plan_name: '',
    project_id: searchForm.project_id,
    plan_type: 'procurement',
    priority: 'normal',
    planned_start_date: '',
    planned_end_date: '',
    total_budget: null,
    description: '',
    remark: '',
    items: []
  })
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row) => {
  Object.assign(formData, {
    ...row,
    planned_start_date: row.planned_start_date ? row.planned_start_date.substring(0, 10) : '',
    planned_end_date: row.planned_end_date ? row.planned_end_date.substring(0, 10) : ''
  })
  dialogVisible.value = true
}

// 查看
const handleView = (row) => {
  viewingPlanId.value = row.id
  detailDialogVisible.value = true
}

// 提交表单
const handleSubmit = async () => {
  try {
    await formRef.value.validate()

    if (formData.items.length === 0) {
      ElMessage.warning('请至少添加一个计划项目')
      return
    }

    dialogLoading.value = true

    if (formData.id) {
      await planStore.updatePlan(formData.id, formData)
    } else {
      await planStore.createPlan(formData)
    }

    dialogVisible.value = false
  } catch (error) {
    console.error('提交失败:', error)
  } finally {
    dialogLoading.value = false
  }
}

// 权限检查
const canEdit = (row) => {
  return authStore.hasPermission('material_plan_edit') && row.status === 'draft'
}

const canDelete = (row) => {
  return authStore.hasPermission('material_plan_delete') && row.status === 'draft'
}

const canSubmit = (row) => {
  return authStore.hasPermission('material_plan_edit') && row.status === 'draft'
}

const canApprove = (row) => {
  return authStore.hasPermission('material_plan_approve') && row.status === 'pending'
}


const canResubmit = (row) => {
  return authStore.hasPermission('material_plan_edit') && row.status === 'rejected'
}

const canCancel = (row) => {
  return authStore.hasPermission('material_plan_approve') &&
    (row.status === 'draft' || row.status === 'pending')
}

// 操作命令
const handleActionCommand = async (command, row) => {
  switch (command) {
    case 'submit':
      await handleSubmitPlan(row)
      break
    case 'approve':
      await handleApprove(row)
      break
    case 'reject':
      await handleReject(row)
      break
    case 'resubmit':
      await handleResubmit(row)
      break
    case 'cancel':
      await handleCancel(row)
      break
    case 'delete':
      await handleDelete(row)
      break
  }
}

const handleSubmitPlan = async (row) => {
  try {
    await ElMessageBox.confirm('确认提交该计划进行审批？', '提示', {
      type: 'warning'
    })
    await planStore.submitPlan(row.id)
  } catch (error) {
    // 用户取消
  }
}

const handleApprove = async (row) => {
  try {
    await ElMessageBox.confirm('确认批准该计划？', '提示', {
      type: 'warning'
    })
    await planStore.approvePlan(row.id)
  } catch (error) {
    // 用户取消
  }
}

const handleReject = async (row) => {
  try {
    const { value } = await ElMessageBox.prompt('请输入拒绝原因', '拒绝计划', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputPattern: /.+/,
      inputErrorMessage: '请输入拒绝原因'
    })
    await planStore.rejectPlan(row.id, { remark: value })
  } catch (error) {
    // 用户取消
  }
}


const handleResubmit = async (row) => {
  try {
    await ElMessageBox.confirm('确认重新提交该计划？', '提示', {
      type: 'warning'
    })
    await planStore.resubmitPlan(row.id)
  } catch (error) {
    // 用户取消
  }
}

const handleCancel = async (row) => {
  try {
    const { value } = await ElMessageBox.prompt('请输入取消原因', '取消计划', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputPattern: /.+/,
      inputErrorMessage: '请输入取消原因'
    })
    await planStore.cancelPlan(row.id, { reason: value })
  } catch (error) {
    // 用户取消
  }
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确认删除该计划？删除后不可恢复。', '提示', {
      type: 'warning'
    })
    await planStore.deletePlan(row.id)
  } catch (error) {
    // 用户取消
  }
}

// 项目管理
const handleAddItem = () => {
  Object.assign(itemForm, {
    material_name: '',
    material_code: '',
    specification: '',
    category: '',
    unit: '',
    planned_quantity: 1,
    unit_price: null,
    required_date: '',
    priority: 'normal',
    remark: ''
  })
  editingItemIndex = -1
  itemDialogVisible.value = true
}

const handleSaveItem = async () => {
  try {
    await itemFormRef.value.validate()

    const item = { ...itemForm }
    if (editingItemIndex >= 0) {
      formData.items[editingItemIndex] = item
    } else {
      formData.items.push(item)
    }

    itemDialogVisible.value = false
  } catch (error) {
    console.error('保存项目失败:', error)
  }
}

const handleDeleteItem = (index) => {
  formData.items.splice(index, 1)
}

/**
 * 打开导入对话框
 */
const handleImportItems = () => {
  importDialogVisible.value = true
}

/**
 * 处理导入数据
 * ExcelImportMapper 组件完成字段映射后会调用此函数
 * 这是一个纯前端处理，不涉及后端 API 调用
 * 数据将直接添加到计划项目列表中
 */
const handleBatchImportItems = async (data) => {
  // 直接返回数据，不做任何 API 调用
  // 调用端会将数据添加到表单中
  return {
    success: true,
    data: {
      total: data.items?.length || 0,
      success: data.items?.length || 0,
      failed: 0,
      items: data.items || []
    }
  }
}

/**
 * 导入成功回调
 * 将导入的数据添加到表单
 */
const handleImportItemsSuccess = (result) => {
  if (result.success && result.data?.items) {
    const importedItems = result.data.items.map(item => ({
      material_name: item.material_name || '',
      material_code: item.material_code || '',
      specification: item.specification || '',
      category: item.category || '',
      unit: item.unit || '',
      planned_quantity: parseInt(item.planned_quantity) || 1,
      unit_price: parseFloat(item.unit_price) || null,
      required_date: item.required_date || '',
      priority: item.priority || 'normal',
      remark: item.remark || ''
    }))

    // 合并到现有项目列表
    formData.items = [...formData.items, ...importedItems]
    ElMessage.success(`成功导入 ${importedItems.length} 个项目`)
  }
}

// 统计分析
const handleStatistics = () => {
  router.push('/plan-statistics')
}

// 辅助函数
const getStatusLabel = (status) => {
  const labels = {
    draft: '草稿',
    pending: '待审批',
    approved: '已批准',
    active: '进行中',
    completed: '已完成',
    cancelled: '已取消',
    rejected: '已拒绝'
  }
  return labels[status] || status
}

const getStatusTagType = (status) => {
  const types = {
    draft: 'info',
    pending: 'warning',
    approved: 'success',
    active: 'primary',
    completed: 'success',
    cancelled: 'danger',
    rejected: 'danger'
  }
  return types[status] || 'info'
}

const getPlanTypeLabel = (type) => {
  const labels = {
    procurement: '采购',
    usage: '使用',
    mixed: '混合'
  }
  return labels[type] || type
}

const getPlanTypeTagType = (type) => {
  const types = {
    procurement: 'primary',
    usage: 'success',
    mixed: 'warning'
  }
  return types[type] || 'info'
}

const getPriorityLabel = (priority) => {
  const labels = {
    urgent: '紧急',
    high: '高',
    normal: '普通',
    low: '低'
  }
  return labels[priority] || priority
}

const getPriorityTagType = (priority) => {
  const types = {
    urgent: 'danger',
    high: 'warning',
    normal: 'info',
    low: ''
  }
  return types[priority] || 'info'
}

const formatAmount = (amount) => {
  return Number(amount || 0).toLocaleString('zh-CN', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  })
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN')
}

// 初始化
onMounted(() => {
  fetchProjects()
  planStore.fetchPlans()
})
</script>

<style scoped>
.material-plans-container {
  padding: 20px;
}

.mt-20 {
  margin-top: 20px;
}

.mb-10 {
  margin-bottom: 10px;
}

.ml-10 {
  margin-left: 10px;
}

.progress-info {
  padding: 0 10px;
}

.tip-list {
  margin: 0;
  padding-left: 20px;
}

.tip-list li {
  margin-bottom: 5px;
}

.mb-20 {
  margin-bottom: 20px;
}
</style>
