<template>
  <div class="operation-logs-container">
    <el-card shadow="never">
      <!-- 高级筛选区域 -->
      <div class="filter-section">
        <el-collapse v-model="activeFilterPanels" class="filter-collapse">
          <el-collapse-item title="筛选条件" name="filters">
            <template #title>
              <div class="filter-header">
                <el-icon><Filter /></el-icon>
                <span class="filter-title">筛选条件</span>
                <el-tag v-if="activeFilterCount > 0" type="primary" size="small" style="margin-left: 10px">
                  {{ activeFilterCount }} 个条件
                </el-tag>
              </div>
            </template>

            <el-form :model="filters" label-width="90px" class="filter-form">
              <!-- 第一行：关键字、模块、操作类型 -->
              <el-row :gutter="16">
                <el-col :xs="24" :sm="12" :md="8" :lg="6">
                  <el-form-item label="关键字搜索">
                    <el-input
                      v-model="filters.keyword"
                      placeholder="编号/用户名"
                      clearable
                      @clear="handleFilterChange"
                    >
                      <template #prefix>
                        <el-icon><Search /></el-icon>
                      </template>
                    </el-input>
                  </el-form-item>
                </el-col>
                <el-col :xs="24" :sm="12" :md="8" :lg="6">
                  <el-form-item label="模块">
                    <el-select
                      v-model="filters.module"
                      placeholder="选择模块"
                      clearable
                      @change="handleFilterChange"
                      style="width: 100%"
                    >
                      <el-option
                        v-for="item in moduleOptions"
                        :key="item.value"
                        :label="item.label"
                        :value="item.value"
                      />
                    </el-select>
                  </el-form-item>
                </el-col>
                <el-col :xs="24" :sm="12" :md="8" :lg="6">
                  <el-form-item label="操作类型">
                    <el-select
                      v-model="filters.operation"
                      placeholder="选择操作"
                      clearable
                      @change="handleFilterChange"
                      style="width: 100%"
                    >
                      <el-option
                        v-for="item in operationOptions"
                        :key="item.value"
                        :label="item.label"
                        :value="item.value"
                      />
                    </el-select>
                  </el-form-item>
                </el-col>
                <el-col :xs="24" :sm="12" :md="8" :lg="6">
                  <el-form-item label="状态">
                    <el-select
                      v-model="filters.status"
                      placeholder="选择状态"
                      clearable
                      @change="handleFilterChange"
                      style="width: 100%"
                    >
                      <el-option label="成功" value="success" />
                      <el-option label="失败" value="error" />
                    </el-select>
                  </el-form-item>
                </el-col>
              </el-row>

              <!-- 第二行：日期范围、资源类型 -->
              <el-row :gutter="16">
                <el-col :xs="24" :sm="12" :md="8" :lg="8">
                  <el-form-item label="日期范围">
                    <el-date-picker
                      v-model="dateRange"
                      type="daterange"
                      range-separator="至"
                      start-placeholder="开始日期"
                      end-placeholder="结束日期"
                      value-format="YYYY-MM-DD"
                      unlink-panels
                      @change="handleDateRangeChange"
                      style="width: 100%"
                    />
                  </el-form-item>
                </el-col>
                <el-col :xs="24" :sm="12" :md="8" :lg="8">
                  <el-form-item label="资源类型">
                    <el-select
                      v-model="filters.resource_type"
                      placeholder="选择资源类型"
                      clearable
                      @change="handleFilterChange"
                      style="width: 100%"
                    >
                      <el-option
                        v-for="item in resourceTypeOptions"
                        :key="item.value"
                        :label="item.label"
                        :value="item.value"
                      />
                    </el-select>
                  </el-form-item>
                </el-col>
                <el-col :xs="24" :sm="12" :md="8" :lg="8">
                  <el-form-item label="资源编号">
                    <el-input
                      v-model="filters.resource_no"
                      placeholder="输入资源编号"
                      clearable
                      @clear="handleFilterChange"
                    />
                  </el-form-item>
                </el-col>
              </el-row>

              <!-- 快捷时间选择 -->
              <el-row :gutter="16">
                <el-col :span="24">
                  <div class="quick-filters">
                    <span class="quick-filter-label">快捷选择：</span>
                    <el-button
                      v-for="item in quickDateOptions"
                      :key="item.value"
                      :type="isQuickDateActive(item.value) ? 'primary' : 'default'"
                      size="small"
                      @click="handleQuickDate(item.value)"
                    >
                      {{ item.label }}
                    </el-button>
                  </div>
                </el-col>
              </el-row>

              <!-- 操作按钮 -->
              <el-row :gutter="16" style="margin-top: 10px">
                <el-col :span="24" class="filter-actions">
                  <el-button type="primary" :icon="Search" @click="handleSearch">
                    查询
                  </el-button>
                  <el-button :icon="Refresh" @click="handleReset">
                    重置
                  </el-button>
                  <el-button
                    :icon="Download"
                    @click="handleExport"
                    :loading="exporting"
                  >
                    导出
                  </el-button>
                  <el-divider direction="vertical" />
                  <el-tag v-if="logs.length > 0" type="info" size="small">
                    共 {{ pagination.total }} 条记录
                  </el-tag>
                </el-col>
              </el-row>
            </el-form>
          </el-collapse-item>
        </el-collapse>
      </div>

      <!-- 统计信息 -->
      <div v-if="statistics" class="statistics-section">
        <el-row :gutter="16">
          <el-col :xs="12" :sm="6">
            <div class="stat-card stat-primary">
              <div class="stat-value">{{ statistics.total_operations || 0 }}</div>
              <div class="stat-label">总操作数</div>
            </div>
          </el-col>
          <el-col :xs="12" :sm="6">
            <div class="stat-card stat-success">
              <div class="stat-value">{{ statistics.success_count || 0 }}</div>
              <div class="stat-label">成功</div>
            </div>
          </el-col>
          <el-col :xs="12" :sm="6">
            <div class="stat-card stat-danger">
              <div class="stat-value">{{ statistics.error_count || 0 }}</div>
              <div class="stat-label">失败</div>
            </div>
          </el-col>
          <el-col :xs="12" :sm="6">
            <div class="stat-card stat-info">
              <div class="stat-value">{{ statistics.period_days || 7 }}天</div>
              <div class="stat-label">统计周期</div>
            </div>
          </el-col>
        </el-row>
      </div>

      <!-- 数据表格 -->
      <el-table
        v-loading="loading"
        :data="logs"
        border
        stripe
        style="width: 100%; margin-top: 20px"
        :default-sort="{ prop: 'created_at', order: 'descending' }"
      >
        <el-table-column type="index" label="#" width="50" />
        <el-table-column prop="username" label="操作人" width="100" show-overflow-tooltip />
        <el-table-column prop="operation" label="操作类型" width="100">
          <template #default="scope">
            <el-tag :type="getOperationTagType(scope.row.operation)" size="small">
              {{ getOperationText(scope.row.operation) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="module" label="模块" width="100">
          <template #default="scope">
            {{ getModuleText(scope.row.module) }}
          </template>
        </el-table-column>
        <el-table-column prop="resource_type" label="资源类型" width="120">
          <template #default="scope">
            {{ getResourceTypeText(scope.row.resource_type) }}
          </template>
        </el-table-column>
        <el-table-column prop="resource_no" label="资源编号" width="140" show-overflow-tooltip />
        <el-table-column prop="ip_address" label="IP地址" width="130" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'success' ? 'success' : 'danger'" size="small">
              {{ scope.row.status === 'success' ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="操作时间" width="160" sortable />
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="scope">
            <el-button
              type="primary"
              size="small"
              link
              :icon="View"
              @click="handleViewDetail(scope.row)"
            >
              详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="操作日志详情"
      width="700px"
      :close-on-click-modal="false"
    >
      <div v-if="currentLog" class="log-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="操作人">{{ currentLog.username }}</el-descriptions-item>
          <el-descriptions-item label="操作类型">
            <el-tag :type="getOperationTagType(currentLog.operation)" size="small">
              {{ getOperationText(currentLog.operation) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="模块">{{ getModuleText(currentLog.module) }}</el-descriptions-item>
          <el-descriptions-item label="资源类型">{{ getResourceTypeText(currentLog.resource_type) }}</el-descriptions-item>
          <el-descriptions-item label="资源编号">{{ currentLog.resource_no || '-' }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="currentLog.status === 'success' ? 'success' : 'danger'" size="small">
              {{ currentLog.status === 'success' ? '成功' : '失败' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="IP地址" :span="2">{{ currentLog.ip_address || '-' }}</el-descriptions-item>
          <el-descriptions-item label="请求路径" :span="2">{{ currentLog.request_path || '-' }}</el-descriptions-item>
          <el-descriptions-item label="操作时间" :span="2">{{ currentLog.created_at }}</el-descriptions-item>
        </el-descriptions>

        <!-- 错误信息 -->
        <div v-if="currentLog.error_message" class="error-section">
          <div class="section-title">错误信息</div>
          <el-alert type="error" :closable="false">
            {{ currentLog.error_message }}
          </el-alert>
        </div>

        <!-- 请求参数 -->
        <div v-if="currentLog.request_params" class="params-section">
          <div class="section-title">请求参数</div>
          <pre class="json-content">{{ formatJson(currentLog.request_params) }}</pre>
        </div>

        <!-- 变更内容 -->
        <div v-if="currentLog.changes" class="changes-section">
          <div class="section-title">变更内容</div>
          <pre class="json-content">{{ formatJson(currentLog.changes) }}</pre>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Search,
  Refresh,
  Download,
  Filter,
  View
} from '@element-plus/icons-vue'
import request from '@/api/request'

// 筛选面板展开状态
const activeFilterPanels = ref(['filters'])

// 加载状态
const loading = ref(false)
const exporting = ref(false)

// 数据
const logs = ref([])
const statistics = ref(null)

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 筛选条件
const filters = reactive({
  keyword: '',
  module: '',
  operation: '',
  resource_type: '',
  resource_no: '',
  status: '',
  start_date: '',
  end_date: ''
})

// 日期范围
const dateRange = ref([])

// 详情对话框
const detailDialogVisible = ref(false)
const currentLog = ref(null)

// 快捷日期选项
const quickDateOptions = [
  { label: '今天', value: 'today' },
  { label: '昨天', value: 'yesterday' },
  { label: '最近7天', value: 'week' },
  { label: '最近30天', value: 'month' },
  { label: '本月', value: 'thisMonth' }
]

// 当前快捷日期选择
const currentQuickDate = ref('week')

// 模块选项
const moduleOptions = [
  { label: '物资计划', value: 'material_plan' },
  { label: '入库单', value: 'inbound' },
  { label: '出库单', value: 'outbound' },
  { label: '领料单', value: 'requisition' },
  { label: '库存', value: 'stock' },
  { label: '物资', value: 'material' },
  { label: '工作流', value: 'workflow' },
  { label: '项目', value: 'project' },
  { label: '施工日志', value: 'construction' },
  { label: '系统', value: 'system' }
]

// 操作类型选项
const operationOptions = [
  { label: '创建', value: 'create' },
  { label: '更新', value: 'update' },
  { label: '删除', value: 'delete' },
  { label: '审批通过', value: 'approve' },
  { label: '审批拒绝', value: 'reject' },
  { label: '提交', value: 'submit' },
  { label: '取消', value: 'cancel' },
  { label: '激活', value: 'activate' },
  { label: '完成', value: 'complete' },
  { label: '发放', value: 'issue' },
  { label: '接收', value: 'receive' },
  { label: '调整', value: 'adjust' },
  { label: '登录', value: 'login' },
  { label: '登出', value: 'logout' }
]

// 资源类型选项
const resourceTypeOptions = [
  { label: '物资计划', value: 'MaterialPlan' },
  { label: '入库单', value: 'InboundOrder' },
  { label: '领料单', value: 'Requisition' },
  { label: '库存', value: 'Stock' },
  { label: '物资', value: 'Material' },
  { label: '工作流任务', value: 'WorkflowTask' },
  { label: '工作流实例', value: 'WorkflowInstance' }
]

// 计算激活的筛选条件数量
const activeFilterCount = computed(() => {
  let count = 0
  if (filters.keyword) count++
  if (filters.module) count++
  if (filters.operation) count++
  if (filters.resource_type) count++
  if (filters.resource_no) count++
  if (filters.status) count++
  if (filters.start_date || filters.end_date) count++
  return count
})

// 获取操作日志列表
const fetchLogs = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      keyword: filters.keyword || undefined,
      module: filters.module || undefined,
      operation: filters.operation || undefined,
      resource_type: filters.resource_type || undefined,
      resource_no: filters.resource_no || undefined,
      status: filters.status || undefined,
      start_date: filters.start_date || undefined,
      end_date: filters.end_date || undefined
    }

    const response = await request.get('/audit/operation-logs', { params })
    if (response.data?.success) {
      const result = response.data.data
      logs.value = result.data || []
      pagination.total = result.total || 0
    }
  } catch (error) {
    console.error('获取操作日志失败:', error)
    ElMessage.error('获取操作日志失败')
  } finally {
    loading.value = false
  }
}

// 获取统计数据
const fetchStatistics = async () => {
  try {
    const response = await request.get('/audit/operation-logs/statistics', {
      params: { days: 7 }
    })
    if (response.data?.success) {
      statistics.value = response.data.data
    }
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

// 查询
const handleSearch = () => {
  pagination.page = 1
  fetchLogs()
}

// 重置筛选条件
const handleReset = () => {
  filters.keyword = ''
  filters.module = ''
  filters.operation = ''
  filters.resource_type = ''
  filters.resource_no = ''
  filters.status = ''
  filters.start_date = ''
  filters.end_date = ''
  dateRange.value = []
  currentQuickDate.value = 'week'
  handleQuickDate('week')
  handleSearch()
}

// 处理日期范围变化
const handleDateRangeChange = (value) => {
  if (value && value.length === 2) {
    filters.start_date = value[0]
    filters.end_date = value[1]
    currentQuickDate.value = null
  } else {
    filters.start_date = ''
    filters.end_date = ''
  }
}

// 快捷日期选择
const handleQuickDate = (type) => {
  currentQuickDate.value = type
  const today = new Date()
  let startDate = new Date()
  let endDate = new Date()

  switch (type) {
    case 'today':
      startDate = new Date(today.setHours(0, 0, 0, 0))
      endDate = new Date()
      break
    case 'yesterday':
      const yesterday = new Date(today)
      yesterday.setDate(yesterday.getDate() - 1)
      startDate = new Date(yesterday.setHours(0, 0, 0, 0))
      endDate = new Date(yesterday.setHours(23, 59, 59, 999))
      break
    case 'week':
      startDate = new Date(today)
      startDate.setDate(startDate.getDate() - 7)
      break
    case 'month':
      startDate = new Date(today)
      startDate.setDate(startDate.getDate() - 30)
      break
    case 'thisMonth':
      startDate = new Date(today.getFullYear(), today.getMonth(), 1)
      endDate = new Date()
      break
  }

  const formatDate = (date) => {
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    return `${year}-${month}-${day}`
  }

  filters.start_date = formatDate(startDate)
  filters.end_date = formatDate(endDate)
  dateRange.value = [filters.start_date, filters.end_date]
  handleFilterChange()
}

// 判断快捷日期是否激活
const isQuickDateActive = (type) => {
  return currentQuickDate.value === type
}

// 筛选条件变化时自动查询
const handleFilterChange = () => {
  pagination.page = 1
  fetchLogs()
}

// 分页变化
const handlePageChange = (page) => {
  pagination.page = page
  fetchLogs()
}

const handleSizeChange = (size) => {
  pagination.pageSize = size
  pagination.page = 1
  fetchLogs()
}

// 查看详情
const handleViewDetail = async (row) => {
  try {
    const response = await request.get(`/audit/operation-logs/${row.id}`)
    if (response.data?.success) {
      currentLog.value = response.data.data
      detailDialogVisible.value = true
    }
  } catch (error) {
    console.error('获取详情失败:', error)
    ElMessage.error('获取详情失败')
  }
}

// 导出
const handleExport = async () => {
  exporting.value = true
  try {
    const params = {
      keyword: filters.keyword || undefined,
      module: filters.module || undefined,
      operation: filters.operation || undefined,
      resource_type: filters.resource_type || undefined,
      status: filters.status || undefined,
      start_date: filters.start_date || undefined,
      end_date: filters.end_date || undefined
    }

    const response = await request.get('/audit/operation-logs/export', {
      params,
      responseType: 'blob'
    })

    // 创建下载链接
    const blob = new Blob([response.data], { type: 'application/json' })
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `operation-logs-${new Date().getTime()}.json`
    a.click()
    window.URL.revokeObjectURL(url)

    ElMessage.success('导出成功')
  } catch (error) {
    console.error('导出失败:', error)
    ElMessage.error('导出失败')
  } finally {
    exporting.value = false
  }
}

// 格式化 JSON
const formatJson = (data) => {
  try {
    if (typeof data === 'string') {
      return JSON.stringify(JSON.parse(data), null, 2)
    }
    return JSON.stringify(data, null, 2)
  } catch (error) {
    return data
  }
}

// 获取操作类型文本
const getOperationText = (operation) => {
  const texts = {
    create: '创建',
    update: '更新',
    delete: '删除',
    approve: '审批通过',
    reject: '审批拒绝',
    submit: '提交',
    cancel: '取消',
    activate: '激活',
    complete: '完成',
    issue: '发放',
    receive: '接收',
    adjust: '调整',
    transfer: '转移',
    login: '登录',
    logout: '登出'
  }
  return texts[operation] || operation
}

// 获取操作类型标签颜色
const getOperationTagType = (operation) => {
  const types = {
    create: 'success',
    update: 'primary',
    delete: 'danger',
    approve: 'success',
    reject: 'warning',
    submit: 'info',
    cancel: 'info',
    activate: 'success',
    complete: 'success',
    issue: 'primary',
    receive: 'primary',
    adjust: 'warning',
    transfer: 'info',
    login: 'info',
    logout: 'info'
  }
  return types[operation] || 'info'
}

// 获取模块文本
const getModuleText = (module) => {
  const texts = {
    material_plan: '物资计划',
    inbound: '入库单',
    outbound: '出库单',
    requisition: '领料单',
    stock: '库存',
    material: '物资',
    workflow: '工作流',
    project: '项目',
    construction: '施工日志',
    system: '系统',
    auth: '认证'
  }
  return texts[module] || module
}

// 获取资源类型文本
const getResourceTypeText = (type) => {
  const texts = {
    MaterialPlan: '物资计划',
    InboundOrder: '入库单',
    Requisition: '领料单',
    Stock: '库存',
    Material: '物资',
    WorkflowTask: '工作流任务',
    WorkflowInstance: '工作流实例'
  }
  return texts[type] || type || '-'
}

onMounted(() => {
  handleQuickDate('week')
  fetchLogs()
  fetchStatistics()
})
</script>

<style scoped>
.operation-logs-container {
  padding: 0;
}

/* 筛选区域 */
.filter-section {
  margin-bottom: 20px;
}

.filter-collapse {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
}

.filter-header {
  display: flex;
  align-items: center;
}

.filter-title {
  margin-left: 8px;
  font-weight: 500;
}

.filter-form {
  padding: 16px;
  padding-top: 8px;
}

.filter-actions {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}

/* 快捷筛选 */
.quick-filters {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
  padding: 8px 0;
}

.quick-filter-label {
  font-size: 14px;
  color: #606266;
  font-weight: 500;
}

/* 统计卡片 */
.statistics-section {
  margin-bottom: 20px;
}

.stat-card {
  padding: 16px;
  border-radius: 8px;
  text-align: center;
  background: #fff;
  border: 1px solid #ebeef5;
  transition: all 0.3s;
}

.stat-card:hover {
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

.stat-card.stat-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  border: none;
}

.stat-card.stat-success {
  background: linear-gradient(135deg, #84fab0 0%, #8fd3f4 100%);
  color: #fff;
  border: none;
}

.stat-card.stat-danger {
  background: linear-gradient(135deg, #ff9a9e 0%, #fecfef 99%, #fecfef 100%);
  color: #fff;
  border: none;
}

.stat-card.stat-info {
  background: linear-gradient(135deg, #a1c4fd 0%, #c2e9fb 100%);
  color: #fff;
  border: none;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 8px;
}

.stat-label {
  font-size: 14px;
  opacity: 0.9;
}

/* 分页 */
.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

/* 日志详情 */
.log-detail {
  padding: 10px 0;
}

.section-title {
  font-weight: 600;
  margin-bottom: 10px;
  color: #303133;
  font-size: 14px;
}

.error-section,
.params-section,
.changes-section {
  margin-top: 16px;
}

.json-content {
  background: #f5f7fa;
  padding: 12px;
  border-radius: 4px;
  font-size: 12px;
  line-height: 1.6;
  max-height: 300px;
  overflow-y: auto;
  font-family: 'Courier New', monospace;
}

/* 响应式调整 */
@media (max-width: 768px) {
  .filter-actions {
    flex-direction: column;
    align-items: stretch;
  }

  .filter-actions .el-button {
    width: 100%;
  }

  .quick-filters {
    flex-direction: column;
    align-items: flex-start;
  }

  .stat-value {
    font-size: 20px;
  }

  .stat-label {
    font-size: 12px;
  }
}

/* Element Plus 样式覆盖 */
:deep(.el-collapse-item__header) {
  height: 48px;
  line-height: 48px;
  padding: 0 16px;
  background-color: #f5f7fa;
}

:deep(.el-collapse-item__content) {
  padding-bottom: 0;
}

:deep(.el-table) {
  font-size: 13px;
}

:deep(.el-descriptions__label) {
  font-weight: 600;
}
</style>
