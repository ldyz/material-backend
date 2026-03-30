<template>
  <div class="appointments-container">
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
            placeholder="搜索单号、地点、内容"
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
            <el-option label="待审批" value="pending" />
            <el-option label="已排期" value="scheduled" />
            <el-option label="进行中" value="in_progress" />
            <el-option label="已完成" value="completed" />
            <el-option label="已取消" value="cancelled" />
            <el-option label="已拒绝" value="rejected" />
          </el-select>
          <el-select
            v-model="searchForm.is_urgent"
            placeholder="优先级"
            clearable
            style="width: 120px"
          >
            <el-option label="全部" value="" />
            <el-option label="加急" value="true" />
            <el-option label="普通" value="false" />
          </el-select>
          <el-date-picker
            v-model="searchForm.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
            style="width: 240px"
          />
          <el-button :icon="Refresh" @click="handleReset">重置</el-button>
        </template>
        <template #right>
          <el-button
            type="primary"
            :icon="Plus"
            @click="handleAdd"
            v-if="authStore.hasPermission('appointment_create')"
          >
            创建预约单
          </el-button>
          <el-button
            type="success"
            :icon="Calendar"
            @click="handleCalendarView"
          >
            日历视图
          </el-button>
          <el-button
            type="warning"
            :icon="Download"
            @click="handleExport"
            v-if="authStore.hasPermission('appointment_export')"
          >
            导出
          </el-button>
        </template>
      </TableToolbar>

      <!-- 统计卡片 -->
      <el-row :gutter="16" class="mb-16">
        <el-col :span="4">
          <el-card shadow="hover" class="stat-card">
            <el-statistic title="全部" :value="stats.total" />
          </el-card>
        </el-col>
        <el-col :span="4">
          <el-card shadow="hover" class="stat-card stat-warning">
            <el-statistic title="待审批" :value="stats.pending" />
          </el-card>
        </el-col>
        <el-col :span="4">
          <el-card shadow="hover" class="stat-card stat-primary">
            <el-statistic title="已排期" :value="stats.scheduled" />
          </el-card>
        </el-col>
        <el-col :span="4">
          <el-card shadow="hover" class="stat-card stat-info">
            <el-statistic title="进行中" :value="stats.in_progress" />
          </el-card>
        </el-col>
        <el-col :span="4">
          <el-card shadow="hover" class="stat-card stat-success">
            <el-statistic title="已完成" :value="stats.completed" />
          </el-card>
        </el-col>
        <el-col :span="4">
          <el-card shadow="hover" class="stat-card stat-danger">
            <el-statistic title="加急" :value="stats.urgent" />
          </el-card>
        </el-col>
      </el-row>

      <!-- 表格 -->
      <el-table
        v-loading="loading"
        :data="tableData"
        border
        stripe
        style="width: 100%"
      >
        <el-table-column prop="appointment_no" label="预约单号" width="160" fixed="left">
          <template #default="scope">
            <el-link type="primary" @click="handleView(scope.row)">
              {{ scope.row.appointment_no }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column prop="applicant_name" label="申请人" width="100" />
        <el-table-column prop="work_date" label="作业时间" width="140">
          <template #default="scope">
            {{ formatDateTime(scope.row.work_date, scope.row.time_slot) }}
          </template>
        </el-table-column>
        <el-table-column prop="work_location" label="作业地点" min-width="150" show-overflow-tooltip />
        <el-table-column prop="work_content" label="作业内容" min-width="200" show-overflow-tooltip />
        <el-table-column prop="assigned_worker_name" label="作业人员" width="120">
          <template #default="scope">
            {{ scope.row.assigned_worker_names || scope.row.assigned_worker_name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="is_urgent" label="优先级" width="80" align="center">
          <template #default="scope">
            <el-tag v-if="scope.row.is_urgent" type="danger" size="small">加急</el-tag>
            <el-tag v-else-if="scope.row.priority >= 5" type="warning" size="small">重要</el-tag>
            <el-tag v-else type="info" size="small">普通</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="90">
          <template #default="scope">
            <el-tag :type="getStatusTagType(scope.row.status)" size="small">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160" />
        <el-table-column label="操作" width="340" fixed="right">
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
              审批
            </el-button>
            <el-button
              type="warning"
              size="small"
              :icon="User"
              @click="handleAssign(scope.row)"
              v-if="canAssign(scope.row)"
            >
              分配
            </el-button>
            <el-button
              type="warning"
              size="small"
              :icon="Close"
              @click="handleCancel(scope.row)"
              v-if="canCancel(scope.row)"
            >
              取消
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
    <AppointmentDialog
      v-model="dialogVisible"
      :appointment="currentAppointment"
      :mode="dialogMode"
      @success="handleDialogSuccess"
    />

    <!-- 详情对话框 -->
    <AppointmentDetailDialog
      v-model="detailVisible"
      :appointment-id="currentAppointmentId"
      @approve="handleApproveFromDetail"
    />

    <!-- 审批对话框 -->
    <AppointmentApproveDialog
      v-model="approveVisible"
      :appointment-id="currentAppointmentId"
      @success="handleApproveSuccess"
    />

    <!-- 分配对话框 -->
    <AppointmentAssignDialog
      v-model="assignVisible"
      :appointment="currentAppointment"
      @success="handleAssignSuccess"
    />

    <!-- 日历视图对话框 -->
    <AppointmentCalendarDialog
      v-model="calendarVisible"
    />
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Search, Refresh, Plus, Download, View, Edit, Delete, Check,
  Calendar, User, Close
} from '@element-plus/icons-vue'
import { defineAsyncComponent } from 'vue'
import { useAuthStore } from '@/stores/auth'
import TableToolbar from '@/components/common/TableToolbar.vue'
import ProjectSelector from '@/components/common/ProjectSelector.vue'

// Dynamic imports for dialog components
const AppointmentDialog = defineAsyncComponent(() => import('@/components/Appointment/AppointmentDialog.vue'))
const AppointmentDetailDialog = defineAsyncComponent(() => import('@/components/Appointment/AppointmentDetailDialog.vue'))
const AppointmentApproveDialog = defineAsyncComponent(() => import('@/components/Appointment/AppointmentApproveDialog.vue'))
const AppointmentAssignDialog = defineAsyncComponent(() => import('@/components/Appointment/AppointmentAssignDialog.vue'))
const AppointmentCalendarDialog = defineAsyncComponent(() => import('@/components/Appointment/AppointmentCalendarDialog.vue'))

import { appointmentApi, projectApi } from '@/api'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const loading = ref(false)
const tableData = ref([])
const dialogVisible = ref(false)
const detailVisible = ref(false)
const approveVisible = ref(false)
const assignVisible = ref(false)
const calendarVisible = ref(false)
const dialogMode = ref('create')
const currentAppointment = ref(null)
const currentAppointmentId = ref(null)
const projectList = ref([])

const searchForm = reactive({
  project_id: '',
  keyword: '',
  status: '',
  is_urgent: '',
  dateRange: []
})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const stats = reactive({
  total: 0,
  draft: 0,
  pending: 0,
  scheduled: 0,
  in_progress: 0,
  completed: 0,
  cancelled: 0,
  rejected: 0,
  urgent: 0
})

onMounted(() => {
  loadData()
  loadStats()
  fetchProjects()
})

// 监听路由参数，自动打开详情弹窗
watch(() => route.query, (query) => {
  const appointmentNo = query.appointment_no
  const appointmentId = query.id

  if (appointmentNo || appointmentId) {
    // 延迟执行，确保数据已加载
    setTimeout(() => {
      let targetAppointment = null

      if (appointmentNo) {
        targetAppointment = tableData.value.find(item => item.appointment_no === appointmentNo)
      } else if (appointmentId) {
        targetAppointment = tableData.value.find(item => item.id === parseInt(appointmentId))
      }

      if (targetAppointment) {
        handleView(targetAppointment)
        // 清除查询参数，避免重复触发
        router.replace({ query: {} })
      } else {
        // 如果当前页没有找到，可能数据在其他页，尝试重新加载数据
        loadData().then(() => {
          setTimeout(() => {
            let retryTarget = null
            if (appointmentNo) {
              retryTarget = tableData.value.find(item => item.appointment_no === appointmentNo)
            } else if (appointmentId) {
              retryTarget = tableData.value.find(item => item.id === parseInt(appointmentId))
            }

            if (retryTarget) {
              handleView(retryTarget)
              router.replace({ query: {} })
            }
          }, 100)
        })
      }
    }, 300)
  }
}, { immediate: true })

async function fetchProjects() {
  try {
    const response = await projectApi.getList({ page: 1, page_size: 1000 })
    projectList.value = response.data?.projects || response.data || []
  } catch (error) {
    console.error('获取项目列表失败:', error)
  }
}

async function loadData() {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      status: searchForm.status || undefined,
      keyword: searchForm.keyword || undefined,
      project_id: searchForm.project_id || undefined,
      is_urgent: searchForm.is_urgent ? searchForm.is_urgent === 'true' : undefined
    }

    if (searchForm.dateRange && searchForm.dateRange.length === 2) {
      params.start_date = searchForm.dateRange[0]
      params.end_date = searchForm.dateRange[1]
    }

    const response = await appointmentApi.getList(params)
    // 后端返回格式: { success: true, data: [...], meta: {...} }
    tableData.value = response.data || []
    pagination.total = response.meta?.total || 0
  } catch (error) {
    ElMessage.error('加载数据失败')
  } finally {
    loading.value = false
  }
}

async function loadStats() {
  try {
    const response = await appointmentApi.getStats()
    // 后端返回格式: { success: true, data: {...} }
    Object.assign(stats, response.data || {})
  } catch (error) {
    console.error('加载统计数据失败:', error)
  }
}

function handleReset() {
  Object.assign(searchForm, {
    project_id: '',
    keyword: '',
    status: '',
    is_urgent: '',
    dateRange: []
  })
  pagination.page = 1
  loadData()
}

function handleAdd() {
  currentAppointment.value = null
  dialogMode.value = 'create'
  dialogVisible.value = true
}

function handleView(row) {
  currentAppointmentId.value = row.id
  detailVisible.value = true
}

function handleEdit(row) {
  currentAppointment.value = { ...row }
  dialogMode.value = 'edit'
  dialogVisible.value = true
}

function handleDelete(row) {
  ElMessageBox.confirm('确认删除此预约单？', '提示', {
    type: 'warning'
  }).then(async () => {
    try {
      await appointmentApi.delete(row.id)
      ElMessage.success('删除成功')
      loadData()
      loadStats()
    } catch (error) {
      ElMessage.error(error.message || '删除失败')
    }
  })
}

function handleApprove(row) {
  currentAppointmentId.value = row.id
  approveVisible.value = true
}

function handleAssign(row) {
  currentAppointment.value = { ...row }
  assignVisible.value = true
}

function handleCalendarView() {
  calendarVisible.value = true
}

function handleExport() {
  // TODO: 实现导出功能
  ElMessage.info('导出功能开发中')
}

function handleDialogSuccess() {
  loadData()
  loadStats()
}

function handleApproveSuccess() {
  loadData()
  loadStats()
}

function handleApproveFromDetail(row) {
  currentAppointmentId.value = row.id
  approveVisible.value = true
}

function handleAssignSuccess() {
  loadData()
}

function handlePageChange(page) {
  pagination.page = page
  loadData()
}

function handleSizeChange(size) {
  pagination.pageSize = size
  pagination.page = 1
  loadData()
}

function canEdit(row) {
  // 草稿和待审批状态都可以编辑
  return (row.status === 'draft' || row.status === 'pending') && authStore.hasPermission('appointment_edit')
}

function canDelete(row) {
  // 只有草稿和已取消状态的预约单可以删除
  return (row.status === 'draft' || row.status === 'cancelled') && authStore.hasPermission('appointment_delete')
}

function canCancel(row) {
  // 草稿、待审批、已排期状态可以取消
  return ['draft', 'pending', 'scheduled'].includes(row.status) && authStore.hasPermission('appointment_cancel')
}

function handleCancel(row) {
  ElMessageBox.confirm('确认取消此预约单？取消后状态将变为"已取消"。', '提示', {
    type: 'warning'
  }).then(async () => {
    try {
      await appointmentApi.cancel(row.id)
      ElMessage.success('预约已取消')
      loadData()
      loadStats()
    } catch (error) {
      ElMessage.error(error.message || '取消失败')
    }
  })
}

function canApprove(row) {
  return row.status === 'pending' && authStore.hasPermission('appointment_approve')
}

function canAssign(row) {
  return (row.status === 'pending' || row.status === 'scheduled') &&
         authStore.hasPermission('appointment_assign')
}

function getStatusTagType(status) {
  const types = {
    draft: '',
    pending: 'warning',
    scheduled: 'primary',
    in_progress: 'info',
    completed: 'success',
    cancelled: 'info',
    rejected: 'danger'
  }
  return types[status] || ''
}

function getStatusText(status) {
  const texts = {
    draft: '草稿',
    pending: '待审批',
    scheduled: '已排期',
    in_progress: '进行中',
    completed: '已完成',
    cancelled: '已取消',
    rejected: '已拒绝'
  }
  return texts[status] || status
}

function formatDateTime(dateStr, timeSlot) {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  const dateStr2 = date.toLocaleDateString('zh-CN')
  const slots = { morning: '上午', afternoon: '下午', evening: '晚上', full_day: '全天' }
  return `${dateStr2} ${slots[timeSlot] || timeSlot}`
}
</script>

<style scoped>
.appointments-container {
  padding: 20px;
}

.stat-card {
  text-align: center;
  border: none;
}

.stat-warning :deep(.el-statistic__number) {
  color: #e6a23c;
}

.stat-primary :deep(.el-statistic__number) {
  color: #409eff;
}

.stat-info :deep(.el-statistic__number) {
  color: #909399;
}

.stat-success :deep(.el-statistic__number) {
  color: #67c23a;
}

.stat-danger :deep(.el-statistic__number) {
  color: #f56c6c;
}

.mb-16 {
  margin-bottom: 16px;
}

.mt-20 {
  margin-top: 20px;
}
</style>
