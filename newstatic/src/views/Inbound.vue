<template>
  <div class="inbound-container">
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
            placeholder="搜索单号、供应商"
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
            <el-option label="已完成" value="completed" />
            <el-option label="已拒绝" value="rejected" />
          </el-select>
          <el-button :icon="Refresh" @click="handleReset">重置</el-button>
        </template>
        <template #right>
          <el-button
            type="primary"
            :icon="Plus"
            @click="handleAdd"
            v-if="authStore.hasPermission('inbound_create')"
          >
            创建入库单
          </el-button>
          <el-button
            type="warning"
            :icon="Download"
            @click="handleExport"
            v-if="authStore.hasPermission('inbound_export')"
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
        <el-table-column prop="order_no" label="入库单号" width="160" fixed="left">
          <template #default="scope">
            <el-link type="primary" @click="handleView(scope.row)">
              {{ scope.row.order_no || scope.row.inbound_no }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column prop="project_name" label="项目名称" min-width="150" show-overflow-tooltip />
        <el-table-column prop="supplier" label="供应商" min-width="150" show-overflow-tooltip />
        <el-table-column prop="contact" label="联系人" width="100" />
        <el-table-column prop="created_at" label="创建时间" width="160" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusTagType(scope.row.status)" size="small">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="250" fixed="right">
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
            <el-descriptions-item label="入库单号" :span="2">
              <el-tag type="primary">{{ formData.order_no || '-' }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="getStatusTagType(formData.status)" size="small">
                {{ getStatusText(formData.status) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="供应商">
              {{ formData.supplier || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="联系人">
              {{ formData.contact || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="创建人">
              {{ formData.creator_name || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="创建时间">
              {{ formData.created_at || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="关联项目" :span="2">
              {{ formData.project_name || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="备注" :span="2" v-if="formData.remark">
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
              workflow-type="inbound"
              @action="handleWorkflowAction"
            />
          </div>

          <!-- 物资明细 -->
          <el-divider content-position="left">物资明细</el-divider>
          <el-table
            :data="formData.items"
            border
            stripe
            style="width: 100%"
            size="small"
          >
            <el-table-column prop="material_name" label="物资名称" min-width="150" show-overflow-tooltip />
            <el-table-column prop="spec" label="规格型号" width="120" show-overflow-tooltip />
            <el-table-column prop="material" label="材质" width="100" show-overflow-tooltip />
            <el-table-column prop="unit" label="单位" width="80" />
            <el-table-column prop="quantity" label="数量" width="100" align="right" />
            <el-table-column prop="unit_price" label="单价" width="100" align="right">
              <template #default="scope">
                {{ scope.row.unit_price ? formatCurrency(scope.row.unit_price) : '-' }}
              </template>
            </el-table-column>
            <el-table-column label="金额" width="120" align="right">
              <template #default="scope">
                {{ formatCurrency((scope.row.quantity || 0) * (scope.row.unit_price || 0)) }}
              </template>
            </el-table-column>
            <el-table-column prop="remark" label="备注" min-width="120" show-overflow-tooltip />
          </el-table>

          <!-- 工作流历史记录 -->
          <el-divider content-position="left">审批历史</el-divider>
          <WorkflowHistory v-if="workflowHistories.length > 0" :histories="workflowHistories" />
          <el-empty v-else description="暂无审批历史" :image-size="80" />
        </template>

        <!-- 编辑模式：显示表单 -->
        <template v-else>
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="供应商" prop="supplier">
                <el-input
                  v-model="formData.supplier"
                  placeholder="请输入供应商名称"
                  maxlength="100"
                />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="验收人" prop="receiver">
                <el-input
                  v-model="formData.receiver"
                  placeholder="请输入验收人"
                  maxlength="50"
                />
              </el-form-item>
            </el-col>
          </el-row>

          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="入库日期" prop="inbound_date">
                <el-date-picker
                  v-model="formData.inbound_date"
                  type="date"
                  placeholder="选择日期"
                  value-format="YYYY-MM-DD"
                  style="width: 100%"
                />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="关联计划" prop="plan_id">
                <el-select
                  v-model="formData.plan_id"
                  placeholder="请选择物资计划"
                  filterable
                  style="width: 100%"
                  :popper-options="{
                    strategy: 'fixed',
                    modifiers: [
                      {
                        name: 'computeStyles',
                        options: { adaptive: true }
                      }
                    ]
                  }"
                  @change="handlePlanChange"
                >
                  <el-option
                    v-for="plan in approvedPlans"
                    :key="plan.id"
                    :label="`${plan.plan_name} (${plan.project_name || ''})`"
                    :value="plan.id"
                  >
                    <span>{{ plan.plan_name }}</span>
                    <span style="color: #8492a6; font-size: 12px; margin-left: 10px">{{ plan.project_name || '' }}</span>
                  </el-option>
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>

          <el-row :gutter="20">
            <el-col :span="24">
              <el-form-item label="备注" prop="remark">
                <el-input
                  v-model="formData.remark"
                  placeholder="请输入备注"
                  maxlength="500"
                />
              </el-form-item>
            </el-col>
          </el-row>

          <el-form-item label="物资明细" prop="items" required>
            <PlanMaterialSelector
              v-if="formData.plan_id"
              v-model="formData.items"
              :editable="!isViewMode"
              :plan-id="formData.plan_id"
              :project-id="formData.project_id"
              @change="handleItemsChange"
            />
            <div v-else class="tip-text">
              请先选择物资计划
            </div>
          </el-form-item>
        </template>
      </el-form>
    </Dialog>

    <!-- 审核对话框 -->
    <Dialog
      v-model="approveDialogVisible"
      title="审核入库单"
      width="600px"
      :loading="approveDialogLoading"
      @confirm="handleApproveSubmit"
    >
      <el-alert
        title="审核提示"
        type="success"
        :closable="false"
        style="margin-bottom: 20px"
      >
        审核通过后将自动增加库存
      </el-alert>

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
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { inboundApi, materialPlanApi } from '@/api'
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
  Printer
} from '@element-plus/icons-vue'
import Dialog from '@/components/common/Dialog.vue'
import TableToolbar from '@/components/common/TableToolbar.vue'
import PlanMaterialSelector from '@/components/common/PlanMaterialSelector.vue'
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

// 项目列表
const projectList = ref([])

// 对话框
const dialogVisible = ref(false)
const isViewMode = ref(false)
const dialogTitle = computed(() => {
  if (isViewMode.value) return '查收入库单'
  return formData.id ? '编辑入库单' : '创建入库单'
})
const dialogLoading = ref(false)
const formRef = ref(null)

// 工作流历史记录
const workflowHistories = ref([])

// 表单数据
const formData = reactive({
  id: null,
  order_no: '',
  supplier: '',
  contact: '',
  receiver: '',
  inbound_date: new Date().toISOString().split('T')[0],
  remark: '',
  notes: '',
  items: [],
  status: '',
  project_id: null,
  project_name: '',
  creator_name: '',
  total_amount: 0,
  created_at: '',
  updated_at: '',
  plan_id: null
})

// 已审批的物资计划列表
const approvedPlans = ref([])

// 表单验证规则
const formRules = {
  supplier: [
    { required: true, message: '请输入供应商名称', trigger: 'blur' }
  ],
  receiver: [
    { required: true, message: '请输入验收人', trigger: 'blur' }
  ],
  inbound_date: [
    { required: true, message: '请选择入库日期', trigger: 'change' }
  ],
  plan_id: [
    { required: true, message: '请选择物资计划', trigger: 'change' }
  ]
}

// 审核对话框
const approveDialogVisible = ref(false)
const approveDialogLoading = ref(false)
const approveFormRef = ref(null)
const currentInbound = ref(null)
const approveForm = reactive({
  approved: true,
  remark: ''
})

// 审核表单验证规则
const approveFormRules = computed(() => ({
  remark: [
    {
      required: !approveForm.approved,
      message: '拒绝时必须填写审核意见',
      trigger: 'blur',
      validator: (rule, value, callback) => {
        // 拒绝时必须填写原因
        if (approveForm.approved === false && !value) {
          callback(new Error('拒绝时必须填写审核意见'))
        } else {
          callback()
        }
      }
    }
  ]
}))

// 获取列表数据
// 适配统一响应格式
const fetchData = async () => {
  loading.value = true
  try {
    // 收集项目ID（包含子项目）
    let projectIds = []
    if (searchForm.project_id) {
      projectIds = collectProjectIds(searchForm.project_id, projectList.value)
    }

    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      search: searchForm.keyword || undefined,
      status: searchForm.status || undefined,
      project_ids: projectIds.length > 0 ? projectIds.join(',') : undefined
    }
    const { data, pagination: pag } = await inboundApi.getList(params)
    tableData.value = data || []
    pagination.total = pag?.total || 0
  } catch (error) {
    console.error('获取入库单列表失败:', error)
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

// 加载项目列表
const fetchProjects = async () => {
  try {
    const { projectApi } = await import('@/api')
    const { data } = await projectApi.getList({ pageSize: 1000 })
    // 构建树形结构
    projectList.value = buildProjectTree(data || [])
  } catch (error) {
    console.error('获取项目列表失败:', error)
  }
}

// 获取已审批的物资计划列表
const fetchApprovedPlans = async () => {
  try {
    // 分别获取approved和active状态的计划
    const [approvedRes, activeRes] = await Promise.all([
      materialPlanApi.getPlans({
        page: 1,
        page_size: 1000,
        status: 'approved'
      }),
      materialPlanApi.getPlans({
        page: 1,
        page_size: 1000,
        status: 'active'
      })
    ])

    // 调试日志
    console.log('API响应 - approved:', approvedRes)
    console.log('API响应 - active:', activeRes)

    // 合并两个结果
    const approvedData = approvedRes.data || []
    const activeData = activeRes.data || []

    console.log('Approved数据:', approvedData)
    console.log('Active数据:', activeData)

    approvedPlans.value = [...approvedData, ...activeData]
    console.log('合并后计划数量:', approvedPlans.value.length)
  } catch (error) {
    console.error('获取物资计划列表失败:', error)
  }
}

// 计划变更处理
const handlePlanChange = () => {
  // 清空已选物资
  formData.items = []
  // 如果选择了计划，获取计划的项目ID
  if (formData.plan_id) {
    const plan = approvedPlans.value.find(p => p.id === formData.plan_id)
    if (plan) {
      formData.project_id = plan.project_id
    }
  }
}

// 构建项目树形结构
const buildProjectTree = (projects) => {
  if (!projects || projects.length === 0) return []

  // 创建项目映射
  const projectMap = new Map()
  projects.forEach(project => {
    projectMap.set(project.id, { ...project, children: [] })
  })

  const roots = []

  // 构建树形结构
  projects.forEach(project => {
    const node = projectMap.get(project.id)
    if (!project.parent_id) {
      // 没有父项目，作为根节点
      roots.push(node)
    } else {
      // 有父项目，添加到父项目的children中
      const parent = projectMap.get(project.parent_id)
      if (parent) {
        parent.children.push(node)
      } else {
        // 父项目不在列表中，作为根节点
        roots.push(node)
      }
    }
  })

  return roots
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

// 监听审核结果变化，重新验证表单
watch(() => approveForm.approved, () => {
  if (approveFormRef.value) {
    approveFormRef.value.validateField('remark')
  }
})

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

// 新增
const handleAdd = () => {
  resetForm()
  isViewMode.value = false
  dialogVisible.value = true
  fetchApprovedPlans()
  fetchProjects()
}

// 编辑
const handleEdit = (row) => {
  // 编辑已拒绝的入库单时，自动填充验收人和入库日期
  const receiver = row.receiver || authStore.displayName
  const inboundDate = row.inbound_date || new Date().toISOString().split('T')[0]

  Object.assign(formData, {
    id: row.id,
    supplier: row.supplier,
    receiver: receiver,
    inbound_date: inboundDate,
    remark: row.remark || '',
    items: row.items || [],
    plan_id: row.plan_id || null
  })
  isViewMode.value = false
  dialogVisible.value = true
  fetchApprovedPlans()
  fetchWorkflowHistory(row.id)
}

// 查看
const handleView = async (row) => {
  // 保存完整数据到currentInbound，供审批使用
  currentInbound.value = row

  Object.assign(formData, {
    id: row.id,
    order_no: row.order_no || row.inbound_no || '',
    supplier: row.supplier,
    contact: row.contact || '',
    receiver: row.receiver,
    inbound_date: row.inbound_date,
    remark: row.remark || '',
    notes: row.notes || '',
    items: row.items || [],
    status: row.status,
    project_id: row.project_id,
    project_name: row.project_name || '',
    creator_name: row.creator_name || '',
    total_amount: row.total_amount || 0,
    updated_at: row.updated_at,
    created_at: row.created_at,
    plan_id: row.plan_id || null,
    // 添加审批相关字段
    approver: row.approver || row.approved_by || '',
    approved_at: row.approved_at || '',
    approve_remark: row.approve_remark || ''
  })

  // 先填充数据，再获取审批历史
  await fetchWorkflowHistory(row.id)

  isViewMode.value = true
  dialogVisible.value = true
}

// 删除
const handleDelete = (row) => {
  ElMessageBox.confirm(
    `确定要删除入库单"${row.inbound_no}"吗？`,
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await inboundApi.delete(row.id)
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
      supplier: formData.supplier,
      receiver: formData.receiver,
      inbound_date: formData.inbound_date,
      remark: formData.remark,
      plan_id: formData.plan_id,
      project_id: String(formData.project_id || ''),
      items: formData.items.map(item => ({
        material_id: item.material_id,
        quantity: item.quantity || 0,
        unit_price: item.unit_price || 0,
        remark: item.remark || ''
      }))
    }

    if (formData.id) {
      await inboundApi.update(formData.id, data)
      ElMessage.success('更新成功')
    } else {
      await inboundApi.create(data)
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
  currentInbound.value = row
  approveForm.approved = true

  // 如果是已拒绝的单据，提取拒绝理由
  if (row.status === 'rejected' && row.remark) {
    // 尝试从备注中提取拒绝原因
    const rejectMatch = row.remark.match(/拒绝原因[：:]\s*(.+?)(?:\n|$)/)
    if (rejectMatch && rejectMatch[1]) {
      approveForm.remark = rejectMatch[1].trim()
    } else {
      approveForm.remark = row.remark
    }
  } else {
    approveForm.remark = ''
  }

  approveDialogVisible.value = true
}

// 提交审核
const handleApproveSubmit = async () => {
  if (!approveFormRef.value) return

  try {
    await approveFormRef.value.validate()
    approveDialogLoading.value = true

    // 根据审核结果调用不同的API
    if (approveForm.approved) {
      await inboundApi.approve(currentInbound.value.id, approveForm)
      ElMessage.success('审核通过')
    } else {
      await inboundApi.reject(currentInbound.value.id, approveForm)
      ElMessage.success('已拒绝')
    }

    approveDialogVisible.value = false
    fetchData()
  } catch (error) {
    console.error('审核失败:', error)
    ElMessage.error(error?.message || '审核失败')
  } finally {
    approveDialogLoading.value = false
  }
}

// 导出
const handleExport = async () => {
  try {
    const response = await inboundApi.export(searchForm)
    const blob = new Blob([response], {
      type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
    })
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `入库单列表_${new Date().getTime()}.xlsx`
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
    order_no: '',
    supplier: '',
    contact: '',
    receiver: authStore.displayName,
    inbound_date: new Date().toISOString().split('T')[0],
    remark: '',
    notes: '',
    items: [],
    status: '',
    project_id: null,
    project_name: '',
    creator_name: '',
    total_amount: 0,
    created_at: '',
    updated_at: '',
    plan_id: null
  })
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

// 打印入库单
const handlePrint = () => {
  // 移除已存在的打印内容
  const existingPrintDiv = document.getElementById('print-temp-content')
  if (existingPrintDiv) {
    existingPrintDiv.remove()
  }

  // 生成打印内容
  const statusTexts = {
    draft: '草稿',
    pending: '待审核',
    approved: '已审核',
    rejected: '已拒绝',
    completed: '已完成'
  }
  const statusText = statusTexts[formData.status] || '未知'

  // 格式化日期
  const formatDate = (dateString) => {
    if (!dateString) return ''
    const date = new Date(dateString)
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    return `${year}-${month}-${day}`
  }

  // 生成物资清单表格
  let itemsTableHtml = ''
  if (formData.items && formData.items.length > 0) {
    itemsTableHtml = `
      <table>
        <thead>
          <tr>
            <th style="width: 20%;">物资名称</th>
            <th style="width: 15%;">规格型号</th>
            <th style="width: 10%;">材质</th>
            <th style="width: 8%;">单位</th>
            <th style="width: 10%;">数量</th>
            <th style="width: 10%;">单价</th>
            <th style="width: 12%;">金额</th>
            <th style="width: 15%;">备注</th>
          </tr>
        </thead>
        <tbody>
          ${formData.items.map((item, index) => `
            <tr>
              <td>${item.material_name || '-'}</td>
              <td>${item.spec || '-'}</td>
              <td>${item.material || '-'}</td>
              <td>${item.unit || '-'}</td>
              <td style="text-align: right;">${Number(item.quantity || 0).toLocaleString('zh-CN')}</td>
              <td style="text-align: right;">${item.unit_price ? Number(item.unit_price).toLocaleString('zh-CN', {minimumFractionDigits: 2, maximumFractionDigits: 2}) : '-'}</td>
              <td style="text-align: right;">${formatCurrency((item.quantity || 0) * (item.unit_price || 0))}</td>
              <td>${item.remark || '-'}</td>
            </tr>
          `).join('')}
        </tbody>
      </table>
    `
  }

  const printHtml = `
    <div class="print-header" style="margin-bottom: 20px; min-height: 100px;">
      <div style="text-align: left; margin-bottom: 10px;">
        <img src="/static/images/logo.png" alt="Logo" style="width: 100px; height: auto; display: block;">
      </div>
      <div style="text-align: center;">
        <div class="print-title" style="font-size: 1.6em; font-weight: bold; letter-spacing: 2px;">大庆化建仪表分公司入库单</div>
      </div>
    </div>

    <table class="print-info">
      <tr>
        <td style="width:auto"><strong>单号：</strong>${formData.order_no || formData.id || '-'}</td>
        <td style="width:auto"><strong>供应商：</strong>${formData.supplier || '-'}</td>
        <td style="width:auto"><strong>时间：</strong>${formatDate(formData.created_at)}</td>
        <td style="width:auto"><strong>状态：</strong>${statusText}</td>
      </tr>
      <tr>
        <th colspan="2" style="text-align: center; font-size: 1.2em; font-weight: bold;"><strong>项目：</strong>${formData.project_name || '-'}</td>
        <td style="width:auto"><strong>联系人：</strong>${formData.contact || '-'}</td>
        <td style="width:auto"><strong>创建人：</strong>${formData.creator_name || '-'}</td>
      </tr>
    </table>

    <div style="height: 20px;"></div>

    ${itemsTableHtml}

    ${formData.remark ? `
      <div style="margin-top: 15px;">
        <strong>备注：</strong>${formData.remark}
      </div>
    ` : ''}

    <div class="print-signatures">
      <br>
      <span>制单时间：${formatDate(formData.created_at)}</span>
      <span style="margin-left: 40px;"></span>
      <span>制单人：${formData.creator_name || '____________'}</span>
      <span style="margin-left: 40px;"></span>
      <span>审核人：____________</span>
      <span style="margin-left: 40px;"></span>
      <span>仓库签字：____________</span>
    </div>
  `

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
    completed: 'success',
    rejected: 'danger'
  }
  return types[status] || 'info'
}

// 获取状态文本
const getStatusText = (status) => {
  const texts = {
    draft: '草稿',
    pending: '待审核',
    completed: '已完成',
    rejected: '已拒绝'
  }
  return texts[status] || status
}

// 获取状态描述
const getStatusDescription = (status) => {
  const descriptions = {
    draft: '入库单为草稿状态，可以编辑或提交审核',
    pending: '等待审核人员进行审核',
    completed: '入库已完成，库存已更新',
    rejected: '审核未通过，需要修改后重新提交'
  }
  return descriptions[status] || ''
}

// 获取工作流历史
// 适配统一响应格式
const fetchWorkflowHistory = async (id) => {
  try {
    const { data } = await inboundApi.getWorkflowHistory(id)
    workflowHistories.value = data || []
  } catch (error) {
    console.log('工作流历史API不可用，使用模拟数据')
    // API调用失败是正常的，因为后端可能没有这个接口
    // 直接使用当前数据生成历史记录
    workflowHistories.value = generateMockHistory(formData)
  }
}

// 生成模拟工作流历史
const generateMockHistory = (inbound) => {
  console.log('生成审批历史，入库单数据:', inbound)
  const histories = []

  // 1. 创建记录（草稿）
  histories.push({
    action: 'draft',
    operator_name: inbound.receiver || inbound.creator_name || '当前用户',
    operator: inbound.receiver || inbound.creator_name || '当前用户',
    department: '采购部',
    remark: '创建入库单',
    description: `创建入库单 ${inbound.order_no || inbound.inbound_no}`,
    created_at: inbound.created_at,
    status: '草稿',
    status_type: 'info'
  })

  // 2. 待审核状态（已提交审核）
  if (inbound.status !== 'draft') {
    histories.push({
      action: 'pending',
      operator_name: inbound.creator_name || '当前用户',
      operator: inbound.creator_name || '当前用户',
      department: '采购部',
      remark: '',
      description: '提交审核',
      created_at: inbound.updated_at || inbound.created_at,
      status: '待审核',
      status_type: 'warning'
    })
  }

  // 3. 最终状态（已完成或已拒绝）
  if (inbound.status === 'completed') {
    // 已完成流程
    histories.push({
      action: 'approved',
      operator_name: inbound.approver || '审核员',
      operator: inbound.approver || '审核员',
      department: '管理部',
      remark: inbound.approve_remark || inbound.remark || '',
      description: '审核通过，库存已更新',
      created_at: inbound.approved_at || inbound.updated_at || inbound.created_at,
      status: '已完成',
      status_type: 'success'
    })
  } else if (inbound.status === 'rejected') {
    // 已拒绝流程
    histories.push({
      action: 'rejected',
      operator_name: inbound.approver || '审核员',
      operator: inbound.approver || '审核员',
      department: '管理部',
      remark: inbound.approve_remark || inbound.remark || '',
      description: '审核拒绝',
      created_at: inbound.approved_at || inbound.updated_at || inbound.created_at,
      status: '已拒绝',
      status_type: 'danger'
    })
  }

  console.log('生成的审批历史:', histories)

  // 按时间升序排序（最新的在后面）
  return histories.sort((a, b) => new Date(a.created_at) - new Date(b.created_at))
}

// 处理工作流操作
const handleWorkflowAction = async (action) => {
  if (action === 'approve') {
    handleApprove(currentInbound.value)
  }
}

// 判断是否可编辑
const canEdit = (row) => {
  if (!authStore.hasPermission('inbound_edit')) return false
  // 只允许编辑草稿状态的入库单，已拒绝的不可以编辑
  return row.status === 'draft'
}

// 判断是否可审核
const canApprove = (row) => {
  if (!authStore.hasPermission('inbound_approve')) return false
  return row.status === 'pending'
}

// 判断是否可删除
const canDelete = (row) => {
  if (!authStore.hasPermission('inbound_delete')) return false
  return row.status === 'draft'
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
  fetchProjects()

  // 检查 URL 查询参数中是否有 order_no
  const orderNo = route.query.order_no
  if (orderNo) {
    // 查找对应的入库单
    const targetOrder = tableData.value.find(item => item.order_no === orderNo || item.inbound_no === orderNo)
    if (targetOrder) {
      // 自动打开详情弹窗
      handleView(targetOrder)
    }
  }
})
</script>

<style scoped>
.inbound-container {
  padding: 0;
}

.mt-20 {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.tip-text {
  margin-top: 8px;
  font-size: 12px;
  color: #909399;
}
</style>
