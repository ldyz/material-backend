<template>
  <div class="attendance-container">
    <el-card shadow="never">
      <!-- 标签页 -->
      <el-tabs v-model="activeTab" @tab-change="handleTabChange">
        <el-tab-pane label="打卡记录" name="records">
          <!-- 工具栏 -->
          <TableToolbar>
            <template #left>
              <el-select
                v-model="searchForm.status"
                placeholder="状态"
                clearable
                style="width: 120px"
                @change="handleSearch"
              >
                <el-option label="全部" value="" />
                <el-option label="待确认" value="pending" />
                <el-option label="已确认" value="confirmed" />
                <el-option label="已驳回" value="rejected" />
              </el-select>
              <el-select
                v-model="searchForm.attendance_type"
                placeholder="打卡类型"
                clearable
                style="width: 140px"
                @change="handleSearch"
              >
                <el-option label="全部" value="" />
                <el-option label="上午打卡" value="morning" />
                <el-option label="下午打卡" value="afternoon" />
                <el-option label="中午加班" value="noon_overtime" />
                <el-option label="晚上加班" value="night_overtime" />
              </el-select>
              <el-date-picker
                v-model="searchForm.dateRange"
                type="daterange"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                value-format="YYYY-MM-DD"
                style="width: 240px"
                @change="handleSearch"
              />
              <el-button :icon="Refresh" @click="handleReset">重置</el-button>
            </template>
            <template #right>
              <el-button
                type="primary"
                :icon="Check"
                @click="handleBatchConfirm"
                :disabled="selectedRows.length === 0"
                v-if="authStore.hasPermission('attendance_confirm')"
              >
                批量确认
              </el-button>
            </template>
          </TableToolbar>

          <!-- 统计卡片 -->
          <el-row :gutter="16" class="mb-16">
            <el-col :span="4">
              <el-card shadow="hover" class="stat-card">
                <el-statistic title="总记录" :value="stats.total_records" />
              </el-card>
            </el-col>
            <el-col :span="4">
              <el-card shadow="hover" class="stat-card stat-warning">
                <el-statistic title="待确认" :value="stats.pending_records" />
              </el-card>
            </el-col>
            <el-col :span="4">
              <el-card shadow="hover" class="stat-card stat-success">
                <el-statistic title="已确认" :value="stats.confirmed_records" />
              </el-card>
            </el-col>
            <el-col :span="4">
              <el-card shadow="hover" class="stat-card stat-danger">
                <el-statistic title="已驳回" :value="stats.rejected_records" />
              </el-card>
            </el-col>
            <el-col :span="4">
              <el-card shadow="hover" class="stat-card stat-primary">
                <el-statistic title="今日打卡" :value="stats.today_clock_ins" />
              </el-card>
            </el-col>
            <el-col :span="4">
              <el-card shadow="hover" class="stat-card stat-info">
                <el-statistic title="本月加班(h)" :value="stats.month_overtime_hours" />
              </el-card>
            </el-col>
          </el-row>

          <!-- 表格 -->
          <el-table
            v-loading="loading"
            :data="tableData"
            border
            stripe
            @selection-change="handleSelectionChange"
          >
            <el-table-column type="selection" width="50" />
            <el-table-column prop="user_name" label="姓名" width="100" />
            <el-table-column prop="attendance_type_label" label="打卡类型" width="100">
              <template #default="scope">
                <el-tag :type="getAttendanceTypeTagType(scope.row.attendance_type)" size="small">
                  {{ scope.row.attendance_type_label }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="clock_in_time" label="打卡时间" width="170" />
            <el-table-column prop="appointment_no" label="关联任务" width="160">
              <template #default="scope">
                <span v-if="scope.row.appointment_no">{{ scope.row.appointment_no }}</span>
                <span v-else class="text-muted">-</span>
              </template>
            </el-table-column>
            <el-table-column prop="work_content" label="任务内容" min-width="150" show-overflow-tooltip />
            <el-table-column prop="clock_in_location" label="打卡位置" width="150" show-overflow-tooltip />
            <el-table-column prop="overtime_hours" label="加班(h)" width="80" align="center">
              <template #default="scope">
                <span v-if="scope.row.overtime_hours > 0">{{ scope.row.overtime_hours }}</span>
                <span v-else>-</span>
              </template>
            </el-table-column>
            <el-table-column label="照片" width="80" align="center">
              <template #default="scope">
                <el-image
                  v-if="getFirstPhoto(scope.row)"
                  :src="getFirstPhoto(scope.row)"
                  :preview-src-list="getPhotoList(scope.row)"
                  fit="cover"
                  style="width: 50px; height: 50px; border-radius: 4px;"
                />
                <span v-else class="text-muted">-</span>
              </template>
            </el-table-column>
            <el-table-column prop="status_label" label="状态" width="90" align="center">
              <template #default="scope">
                <el-tag :type="getStatusTagType(scope.row.status)" size="small">
                  {{ scope.row.status_label }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="confirmed_by_name" label="确认人" width="100">
              <template #default="scope">
                <span v-if="scope.row.confirmed_by_name">{{ scope.row.confirmed_by_name }}</span>
                <span v-else>-</span>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="150" fixed="right">
              <template #default="scope">
                <el-button
                  v-if="scope.row.status === 'pending' && authStore.hasPermission('attendance_confirm')"
                  type="success"
                  size="small"
                  @click="handleConfirm(scope.row)"
                >
                  确认
                </el-button>
                <el-button
                  v-if="scope.row.status === 'pending' && authStore.hasPermission('attendance_confirm')"
                  type="danger"
                  size="small"
                  @click="handleReject(scope.row)"
                >
                  驳回
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
              @size-change="loadRecords"
              @current-change="loadRecords"
            />
          </div>
        </el-tab-pane>

        <el-tab-pane label="月度汇总" name="summary">
          <!-- 工具栏 -->
          <TableToolbar>
            <template #left>
              <el-date-picker
                v-model="summaryMonth"
                type="month"
                placeholder="选择月份"
                value-format="YYYY-MM"
                format="YYYY年MM月"
                style="width: 150px"
                @change="loadSummary"
              />
            </template>
            <template #right>
              <el-button
                type="primary"
                :icon="Refresh"
                @click="handleGenerateSummary"
                v-if="authStore.hasPermission('attendance_manage')"
              >
                生成汇总
              </el-button>
              <el-button
                type="success"
                :icon="Check"
                @click="handleBatchConfirmSummary"
                :disabled="selectedSummaries.length === 0"
                v-if="authStore.hasPermission('attendance_manage')"
              >
                批量确认
              </el-button>
            </template>
          </TableToolbar>

          <!-- 汇总表格 -->
          <el-table
            v-loading="summaryLoading"
            :data="summaryData"
            border
            stripe
            @selection-change="handleSummarySelectionChange"
          >
            <el-table-column type="selection" width="50" />
            <el-table-column prop="user_name" label="姓名" width="100" />
            <el-table-column prop="morning_count" label="上午打卡" width="100" align="center" />
            <el-table-column prop="afternoon_count" label="下午打卡" width="100" align="center" />
            <el-table-column prop="noon_overtime_hours" label="中午加班(h)" width="110" align="center" />
            <el-table-column prop="night_overtime_hours" label="晚上加班(h)" width="110" align="center" />
            <el-table-column prop="total_work_days" label="工作天数" width="100" align="center" />
            <el-table-column prop="total_overtime_hours" label="总加班(h)" width="100" align="center" />
            <el-table-column prop="status_label" label="状态" width="90" align="center">
              <template #default="scope">
                <el-tag :type="scope.row.status === 'confirmed' ? 'success' : 'warning'" size="small">
                  {{ scope.row.status_label }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100" fixed="right">
              <template #default="scope">
                <el-button
                  v-if="scope.row.status === 'draft' && authStore.hasPermission('attendance_manage')"
                  type="success"
                  size="small"
                  @click="handleConfirmSummary(scope.row)"
                >
                  确认
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- 驳回对话框 -->
    <el-dialog v-model="rejectDialogVisible" title="驳回原因" width="400px">
      <el-input
        v-model="rejectReason"
        type="textarea"
        :rows="3"
        placeholder="请输入驳回原因"
      />
      <template #footer>
        <el-button @click="rejectDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitReject">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh, Check } from '@element-plus/icons-vue'
import { attendanceApi } from '@/api'
import { useAuthStore } from '@/stores/auth'
import { getAssetUrl } from '@/utils/index'
import TableToolbar from '@/components/common/TableToolbar.vue'

const authStore = useAuthStore()

const activeTab = ref('records')
const loading = ref(false)
const summaryLoading = ref(false)
const tableData = ref([])
const summaryData = ref([])
const selectedRows = ref([])
const selectedSummaries = ref([])
const rejectDialogVisible = ref(false)
const rejectReason = ref('')
const currentRecord = ref(null)

const searchForm = reactive({
  status: '',
  attendance_type: '',
  dateRange: []
})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const stats = ref({
  total_records: 0,
  pending_records: 0,
  confirmed_records: 0,
  rejected_records: 0,
  today_clock_ins: 0,
  month_clock_ins: 0,
  month_overtime_hours: 0
})

const summaryMonth = computed({
  get() {
    const now = new Date()
    return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`
  },
  set() {}
})

// 获取打卡类型标签类型
function getAttendanceTypeTagType(type) {
  const types = {
    morning: '',
    afternoon: 'success',
    noon_overtime: 'warning',
    night_overtime: 'info'
  }
  return types[type] || ''
}

// 获取第一张照片 URL
function getFirstPhoto(row) {
  if (row.photo_urls) {
    try {
      const urls = typeof row.photo_urls === 'string' ? JSON.parse(row.photo_urls) : row.photo_urls
      // 找到第一个有效的照片 URL
      for (const url of urls) {
        if (url && url !== 'null' && url !== '') {
          return getAssetUrl(url)
        }
      }
    } catch (e) {}
  }
  // 兼容单张照片
  if (row.photo_url && row.photo_url !== 'null' && row.photo_url !== '') {
    return getAssetUrl(row.photo_url)
  }
  return ''
}

// 获取所有照片列表（用于预览）
function getPhotoList(row) {
  const photos = []
  if (row.photo_urls) {
    try {
      const urls = typeof row.photo_urls === 'string' ? JSON.parse(row.photo_urls) : row.photo_urls
      urls.forEach(url => {
        if (url && url !== 'null' && url !== '') {
          const fullUrl = getAssetUrl(url)
          if (fullUrl) photos.push(fullUrl)
        }
      })
    } catch (e) {}
  }
  // 兼容单张照片
  if (photos.length === 0 && row.photo_url && row.photo_url !== 'null' && row.photo_url !== '') {
    const fullUrl = getAssetUrl(row.photo_url)
    if (fullUrl) photos.push(fullUrl)
  }
  return photos
}

// 获取状态标签类型
function getStatusTagType(status) {
  const types = {
    pending: 'warning',
    confirmed: 'success',
    rejected: 'danger'
  }
  return types[status] || ''
}

// 加载打卡记录
async function loadRecords() {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      status: searchForm.status || undefined,
      attendance_type: searchForm.attendance_type || undefined,
      start_date: searchForm.dateRange?.[0] || undefined,
      end_date: searchForm.dateRange?.[1] || undefined
    }

    const response = await attendanceApi.getRecords(params)
    tableData.value = response.data || []
    pagination.total = response.meta?.total || 0
  } catch (error) {
    console.error('加载打卡记录失败:', error)
    ElMessage.error('加载打卡记录失败')
  } finally {
    loading.value = false
  }
}

// 加载统计数据
async function loadStats() {
  try {
    const response = await attendanceApi.getStatistics()
    stats.value = response.data || {}
  } catch (error) {
    console.error('加载统计数据失败:', error)
  }
}

// 加载月度汇总
async function loadSummary() {
  if (!summaryMonth.value) return

  summaryLoading.value = true
  try {
    const [year, month] = summaryMonth.value.split('-')
    const response = await attendanceApi.getMonthlySummary({
      year: parseInt(year),
      month: parseInt(month)
    })
    summaryData.value = response.data || []
  } catch (error) {
    console.error('加载月度汇总失败:', error)
    ElMessage.error('加载月度汇总失败')
  } finally {
    summaryLoading.value = false
  }
}

// 重置搜索
function handleReset() {
  searchForm.status = ''
  searchForm.attendance_type = ''
  searchForm.dateRange = []
  pagination.page = 1
  loadRecords()
}

// 搜索
function handleSearch() {
  pagination.page = 1
  loadRecords()
}

// 选择变更
function handleSelectionChange(selection) {
  selectedRows.value = selection.filter(row => row.status === 'pending')
}

function handleSummarySelectionChange(selection) {
  selectedSummaries.value = selection.filter(row => row.status === 'draft')
}

// 标签页切换
function handleTabChange(tab) {
  if (tab === 'records') {
    loadRecords()
    loadStats()
  } else if (tab === 'summary') {
    loadSummary()
  }
}

// 确认打卡记录
async function handleConfirm(row) {
  try {
    await ElMessageBox.confirm('确认该打卡记录？', '提示', {
      type: 'success'
    })
    await attendanceApi.confirmRecord(row.id)
    ElMessage.success('确认成功')
    loadRecords()
    loadStats()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('确认失败')
    }
  }
}

// 驳回打卡记录
function handleReject(row) {
  currentRecord.value = row
  rejectReason.value = ''
  rejectDialogVisible.value = true
}

async function submitReject() {
  if (!rejectReason.value.trim()) {
    ElMessage.warning('请输入驳回原因')
    return
  }

  try {
    await attendanceApi.rejectRecord(currentRecord.value.id, {
      reason: rejectReason.value
    })
    ElMessage.success('驳回成功')
    rejectDialogVisible.value = false
    loadRecords()
    loadStats()
  } catch (error) {
    ElMessage.error('驳回失败')
  }
}

// 批量确认
async function handleBatchConfirm() {
  if (selectedRows.value.length === 0) {
    ElMessage.warning('请选择待确认的记录')
    return
  }

  try {
    await ElMessageBox.confirm(`确认选中的 ${selectedRows.value.length} 条记录？`, '提示', {
      type: 'success'
    })

    for (const row of selectedRows.value) {
      await attendanceApi.confirmRecord(row.id)
    }

    ElMessage.success('批量确认成功')
    loadRecords()
    loadStats()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量确认失败')
    }
  }
}

// 生成月度汇总
async function handleGenerateSummary() {
  if (!summaryMonth.value) return

  try {
    const [year, month] = summaryMonth.value.split('-')
    await attendanceApi.generateMonthly({
      year: parseInt(year),
      month: parseInt(month)
    })
    ElMessage.success('生成成功')
    loadSummary()
  } catch (error) {
    ElMessage.error('生成失败')
  }
}

// 确认月度汇总
async function handleConfirmSummary(row) {
  try {
    await ElMessageBox.confirm('确认该月度汇总？确认后将无法修改。', '提示', {
      type: 'warning'
    })
    await attendanceApi.confirmMonthlySummary(row.id)
    ElMessage.success('确认成功')
    loadSummary()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('确认失败')
    }
  }
}

// 批量确认月度汇总
async function handleBatchConfirmSummary() {
  if (selectedSummaries.value.length === 0) {
    ElMessage.warning('请选择待确认的汇总')
    return
  }

  try {
    await ElMessageBox.confirm(`确认选中的 ${selectedSummaries.value.length} 条汇总？`, '提示', {
      type: 'warning'
    })

    for (const row of selectedSummaries.value) {
      await attendanceApi.confirmMonthlySummary(row.id)
    }

    ElMessage.success('批量确认成功')
    loadSummary()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量确认失败')
    }
  }
}

onMounted(() => {
  loadRecords()
  loadStats()
})
</script>

<style scoped>
.attendance-container {
  padding: 20px;
}

.mb-16 {
  margin-bottom: 16px;
}

.pagination-container {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

.stat-card {
  text-align: center;
}

.stat-card.stat-warning :deep(.el-statistic__head) {
  color: #e6a23c;
}

.stat-card.stat-success :deep(.el-statistic__head) {
  color: #67c23a;
}

.stat-card.stat-danger :deep(.el-statistic__head) {
  color: #f56c6c;
}

.stat-card.stat-primary :deep(.el-statistic__head) {
  color: #409eff;
}

.stat-card.stat-info :deep(.el-statistic__head) {
  color: #909399;
}

.text-muted {
  color: #909399;
}
</style>
