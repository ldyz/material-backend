<template>
  <div class="plan-statistics-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <h3>物资计划统计分析</h3>
          <el-button :icon="Refresh" @click="refreshData">刷新</el-button>
        </div>
      </template>

      <el-skeleton v-if="loading" :rows="10" animated />
      <div v-else-if="statistics">
        <!-- 统计概览卡片 -->
        <el-row :gutter="20" class="stats-cards">
          <el-col :span="6">
            <el-card shadow="hover" class="stat-card">
              <div class="stat-content">
                <div class="stat-icon total">
                  <el-icon><Document /></el-icon>
                </div>
                <div class="stat-info">
                  <div class="stat-value">{{ statistics.total_plans || 0 }}</div>
                  <div class="stat-label">计划总数</div>
                </div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover" class="stat-card">
              <div class="stat-content">
                <div class="stat-icon active">
                  <el-icon><Clock /></el-icon>
                </div>
                <div class="stat-info">
                  <div class="stat-value">
                    {{ (statistics.by_status?.active || 0) + (statistics.by_status?.pending || 0) }}
                  </div>
                  <div class="stat-label">进行中/待审批</div>
                </div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover" class="stat-card">
              <div class="stat-content">
                <div class="stat-icon budget">
                  <el-icon><Coin /></el-icon>
                </div>
                <div class="stat-info">
                  <div class="stat-value">¥{{ formatAmount(statistics.total_budget) }}</div>
                  <div class="stat-label">计划总预算</div>
                </div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover" class="stat-card">
              <div class="stat-content">
                <div class="stat-icon progress">
                  <el-icon><TrendCharts /></el-icon>
                </div>
                <div class="stat-info">
                  <div class="stat-value">{{ statistics.completion_rate?.toFixed(1) || 0 }}%</div>
                  <div class="stat-label">完成率</div>
                </div>
              </div>
            </el-card>
          </el-col>
        </el-row>

        <!-- 图表区域 -->
        <el-row :gutter="20" class="mt-20">
          <el-col :span="12">
            <el-card shadow="hover">
              <template #header>
                <h4>计划状态分布</h4>
              </template>
              <div class="chart-container">
                <Doughnut
                  :data="statusChartData"
                  :options="chartOptions"
                />
              </div>
            </el-card>
          </el-col>
          <el-col :span="12">
            <el-card shadow="hover">
              <template #header>
                <h4>项目状态分布</h4>
              </template>
              <div class="chart-container">
                <Doughnut
                  :data="itemChartData"
                  :options="chartOptions"
                />
              </div>
            </el-card>
          </el-col>
        </el-row>

        <!-- 详细统计表格 -->
        <el-card shadow="hover" class="mt-20">
          <template #header>
            <h4>按状态统计详情</h4>
          </template>
          <el-table :data="statusDetailData" border stripe>
            <el-table-column prop="status" label="状态" width="120">
              <template #default="scope">
                <el-tag :type="getStatusTagType(scope.row.status)" size="small">
                  {{ getStatusLabel(scope.row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="count" label="计划数量" width="120" align="right" />
            <el-table-column prop="percentage" label="占比" width="120" align="right">
              <template #default="scope">
                {{ scope.row.percentage }}%
              </template>
            </el-table-column>
            <el-table-column prop="items" label="项目数量" width="120" align="right" />
            <el-table-column label="进度" min-width="200">
              <template #default="scope">
                <el-progress
                  :percentage="parseFloat(scope.row.percentage)"
                  :stroke-width="10"
                  :show-text="true"
                />
              </template>
            </el-table-column>
          </el-table>
        </el-card>

        <!-- 项目状态表格 -->
        <el-card shadow="hover" class="mt-20">
          <template #header>
            <h4>项目完成情况统计</h4>
          </template>
          <el-table :data="itemStatusData" border stripe>
            <el-table-column prop="status" label="项目状态" width="120">
              <template #default="scope">
                <el-tag :type="getItemStatusTagType(scope.row.status)" size="small">
                  {{ getItemStatusLabel(scope.row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="count" label="项目数量" width="120" align="right" />
            <el-table-column prop="percentage" label="占比" width="120" align="right">
              <template #default="scope">
                {{ scope.row.percentage }}%
              </template>
            </el-table-column>
            <el-table-column label="占比分布" min-width="300">
              <template #default="scope">
                <el-progress
                  :percentage="parseFloat(scope.row.percentage)"
                  :stroke-width="10"
                  :color="getItemStatusColor(scope.row.status)"
                />
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { Doughnut } from 'vue-chartjs'
import { Chart as ChartJS, Title, Tooltip, Legend, ArcElement, CategoryScale, LinearScale } from 'chart.js'
import { materialPlanApi } from '@/api'
import { ElMessage } from 'element-plus'
import {
  Refresh,
  Document,
  Clock,
  Coin,
  TrendCharts
} from '@element-plus/icons-vue'

// Register Chart.js components
ChartJS.register(Title, Tooltip, Legend, ArcElement, CategoryScale, LinearScale)

const loading = ref(false)
const statistics = ref(null)

// Chart.js options
const chartOptions = ref({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'left'
    },
    tooltip: {
      callbacks: {
        label: function(context) {
          const label = context.label || ''
          const value = context.parsed || 0
          const total = context.dataset.data.reduce((a, b) => a + b, 0)
          const percentage = ((value / total) * 100).toFixed(1)
          return `${label}: ${value} (${percentage}%)`
        }
      }
    }
  }
})

// 计划状态图表数据
const statusChartData = computed(() => {
  if (!statistics.value) return { labels: [], datasets: [] }

  const byStatus = statistics.value.by_status || {}
  const labels = []
  const data = []
  const colors = []

  const statusConfig = [
    { key: 'draft', label: '草稿', color: '#909399' },
    { key: 'pending', label: '待审批', color: '#e6a23c' },
    { key: 'approved', label: '已批准', color: '#67c23a' },
    { key: 'active', label: '进行中', color: '#409eff' },
    { key: 'completed', label: '已完成', color: '#67c23a' },
    { key: 'cancelled', label: '已取消', color: '#f56c6c' },
    { key: 'rejected', label: '已拒绝', color: '#f56c6c' }
  ]

  for (const config of statusConfig) {
    const count = byStatus[config.key] || 0
    if (count > 0) {
      labels.push(config.label)
      data.push(count)
      colors.push(config.color)
    }
  }

  return {
    labels,
    datasets: [{
      data,
      backgroundColor: colors
    }]
  }
})

// 项目状态图表数据
const itemChartData = computed(() => {
  if (!statistics.value) return { labels: [], datasets: [] }

  const itemsByStatus = statistics.value.items_by_status || {}
  const labels = []
  const data = []
  const colors = []

  const statusConfig = [
    { key: 'pending', label: '待处理', color: '#909399' },
    { key: 'partial', label: '部分完成', color: '#e6a23c' },
    { key: 'completed', label: '已完成', color: '#67c23a' },
    { key: 'cancelled', label: '已取消', color: '#f56c6c' }
  ]

  for (const config of statusConfig) {
    const count = itemsByStatus[config.key] || 0
    if (count > 0) {
      labels.push(config.label)
      data.push(count)
      colors.push(config.color)
    }
  }

  return {
    labels,
    datasets: [{
      data,
      backgroundColor: colors
    }]
  }
})

// 获取统计数据
const fetchStatistics = async () => {
  loading.value = true
  try {
    const response = await materialPlanApi.getStatistics()
    if (response.success) {
      statistics.value = response.data
    } else {
      ElMessage.error('获取统计数据失败')
    }
  } catch (error) {
    console.error('获取统计数据失败:', error)
    ElMessage.error('获取统计数据失败')
  } finally {
    loading.value = false
  }
}

// 刷新数据
const refreshData = () => {
  fetchStatistics()
}

// 状态详情数据
const statusDetailData = computed(() => {
  if (!statistics.value) return []

  const byStatus = statistics.value.by_status || {}
  const total = statistics.value.total_plans || 1

  return [
    { status: 'draft', count: byStatus.draft || 0, items: 0 },
    { status: 'pending', count: byStatus.pending || 0, items: 0 },
    { status: 'approved', count: byStatus.approved || 0, items: 0 },
    { status: 'active', count: byStatus.active || 0, items: 0 },
    { status: 'completed', count: byStatus.completed || 0, items: 0 },
    { status: 'cancelled', count: byStatus.cancelled || 0, items: 0 },
    { status: 'rejected', count: byStatus.rejected || 0, items: 0 }
  ]
    .filter(item => item.count > 0)
    .map(item => ({
      ...item,
      percentage: ((item.count / total) * 100).toFixed(1)
    }))
})

// 项目状态数据
const itemStatusData = computed(() => {
  if (!statistics.value) return []

  const itemsByStatus = statistics.value.items_by_status || {}
  const total = statistics.value.total_items || 1

  return [
    { status: 'pending', count: itemsByStatus.pending || 0 },
    { status: 'partial', count: itemsByStatus.partial || 0 },
    { status: 'completed', count: itemsByStatus.completed || 0 },
    { status: 'cancelled', count: itemsByStatus.cancelled || 0 }
  ]
    .filter(item => item.count > 0)
    .map(item => ({
      ...item,
      percentage: ((item.count / total) * 100).toFixed(1)
    }))
})

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

const getItemStatusLabel = (status) => {
  const labels = {
    pending: '待处理',
    partial: '部分完成',
    completed: '已完成',
    cancelled: '已取消'
  }
  return labels[status] || status
}

const getItemStatusTagType = (status) => {
  const types = {
    pending: 'info',
    partial: 'warning',
    completed: 'success',
    cancelled: 'danger'
  }
  return types[status] || 'info'
}

const getItemStatusColor = (status) => {
  const colors = {
    pending: '#909399',
    partial: '#e6a23c',
    completed: '#67c23a',
    cancelled: '#f56c6c'
  }
  return colors[status] || '#909399'
}

const formatAmount = (amount) => {
  return Number(amount || 0).toLocaleString('zh-CN', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  })
}

onMounted(() => {
  fetchStatistics()
})
</script>

<style scoped>
.plan-statistics-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.stats-cards {
  margin-bottom: 20px;
}

.stat-card {
  border-radius: 8px;
  overflow: hidden;
}

.stat-content {
  display: flex;
  align-items: center;
  padding: 10px;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 15px;
}

.stat-icon :deep(.el-icon) {
  font-size: 28px;
  color: white;
}

.stat-icon.total {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-icon.active {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.stat-icon.budget {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-icon.progress {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #303133;
  margin-bottom: 5px;
}

.stat-label {
  font-size: 14px;
  color: #909399;
}

.chart-container {
  height: 300px;
  position: relative;
}

.mt-20 {
  margin-top: 20px;
}
</style>
