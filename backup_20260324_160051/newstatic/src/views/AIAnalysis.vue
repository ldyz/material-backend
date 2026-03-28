<template>
  <div class="ai-dashboard">
    <!-- 页面标题栏 -->
    <div class="dashboard-header">
      <h1>
        <el-icon><TrendCharts /></el-icon>
        数据分析看板
      </h1>
      <div class="header-actions">
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          value-format="YYYY-MM-DD"
          size="default"
          @change="handleDateChange"
        />
        <el-button :icon="Refresh" @click="refreshData">刷新</el-button>
      </div>
    </div>

    <!-- 关键指标卡片 -->
    <el-row :gutter="20" class="metrics-row">
      <el-col :span="6">
        <div class="metric-card blue">
          <div class="metric-icon">
            <el-icon :size="32"><Box /></el-icon>
          </div>
          <div class="metric-content">
            <div class="metric-label">物资总数</div>
            <div class="metric-value">{{ metrics.totalMaterials }}</div>
          </div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="metric-card green">
          <div class="metric-icon">
            <el-icon :size="32"><Goods /></el-icon>
          </div>
          <div class="metric-content">
            <div class="metric-label">库存总量</div>
            <div class="metric-value">{{ metrics.totalStock }}</div>
          </div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="metric-card orange">
          <div class="metric-icon">
            <el-icon :size="32"><Upload /></el-icon>
          </div>
          <div class="metric-content">
            <div class="metric-label">本月入库</div>
            <div class="metric-value">{{ metrics.monthlyInbound }}</div>
          </div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="metric-card red">
          <div class="metric-icon">
            <el-icon :size="32"><Download /></el-icon>
          </div>
          <div class="metric-content">
            <div class="metric-label">本月出库</div>
            <div class="metric-value">{{ metrics.monthlyOutbound }}</div>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 主要内容区域 -->
    <el-row :gutter="20" style="margin-top: 20px">
      <!-- 左侧：库存预警 -->
      <el-col :span="12">
        <el-card class="dashboard-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><Warning /></el-icon>
                库存预警
              </span>
              <el-tag v-if="warnings.length > 0" type="danger" size="small">
                {{ warnings.length }} 项
              </el-tag>
              <el-tag v-else type="success" size="small">正常</el-tag>
            </div>
          </template>
          <div v-if="warnings.length > 0" class="warning-list">
            <div v-for="(item, index) in warnings" :key="index" class="warning-item" :class="'warning-' + item.level">
              <el-icon v-if="item.level === 'danger'"><CircleClose /></el-icon>
              <el-icon v-else><WarningFilled /></el-icon>
              <div class="warning-content">
                <div class="warning-title">{{ item.name }}</div>
                <div class="warning-detail">当前库存：{{ item.stock }} {{ item.unit }}</div>
              </div>
              <el-tag :type="item.level === 'danger' ? 'danger' : 'warning'" size="small">
                {{ item.levelText }}
              </el-tag>
            </div>
          </div>
          <el-empty v-else description="暂无库存预警" :image-size="100" />
        </el-card>
      </el-col>

      <!-- 右侧：AI 智能推荐 -->
      <el-col :span="12">
        <el-card class="dashboard-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><MagicStick /></el-icon>
                AI 智能洞察
              </span>
            </div>
          </template>
          <div v-if="insights.length > 0" class="insight-list">
            <div v-for="(item, index) in insights" :key="index" class="insight-item" :class="'insight-' + item.type">
              <div class="insight-icon">
                <el-icon v-if="item.type === 'warning'"><Warning /></el-icon>
                <el-icon v-else-if="item.type === 'success'"><CircleCheck /></el-icon>
                <el-icon v-else><InfoFilled /></el-icon>
              </div>
              <div class="insight-content">
                <div class="insight-title">{{ item.title }}</div>
                <div class="insight-desc">{{ item.description }}</div>
              </div>
            </div>
          </div>
          <el-empty v-else description="暂无AI洞察" :image-size="100" />
        </el-card>
      </el-col>
    </el-row>

    <!-- 物资分类占比 -->
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card class="dashboard-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><PieChart /></el-icon>
                物资分类占比
              </span>
            </div>
          </template>
          <div v-if="categoryData.length > 0" class="category-list">
            <div v-for="(item, index) in categoryData" :key="index" class="category-item">
              <div class="category-info">
                <span class="category-dot" :style="{ backgroundColor: item.color }"></span>
                <span class="category-name">{{ item.name }}</span>
              </div>
              <div class="category-stats">
                <span class="category-value">{{ item.value }}%</span>
                <el-progress
                  :percentage="item.value"
                  :show-text="false"
                  :stroke-width="8"
                  :color="item.color"
                />
              </div>
            </div>
          </div>
          <el-empty v-else description="暂无分类数据" :image-size="100" />
        </el-card>
      </el-col>

      <!-- 项目进度 Top 5 -->
      <el-col :span="12">
        <el-card class="dashboard-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><TrendCharts /></el-icon>
                项目进度 Top 5
              </span>
            </div>
          </template>
          <div v-if="projectRanking.length > 0" class="ranking-list">
            <div v-for="(item, index) in projectRanking" :key="index" class="ranking-item">
              <div class="ranking-badge" :class="'rank-' + (index + 1)">{{ index + 1 }}</div>
              <div class="ranking-content">
                <div class="ranking-name">{{ item.name }}</div>
                <el-progress
                  :percentage="item.progress"
                  :status="getProgressStatus(item.progress)"
                  :stroke-width="10"
                />
              </div>
            </div>
          </div>
          <el-empty v-else description="暂无项目数据" :image-size="100" />
        </el-card>
      </el-col>
    </el-row>

    <!-- 快速操作 -->
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="24">
        <el-card class="dashboard-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><Operation /></el-icon>
                快速操作
              </span>
            </div>
          </template>
          <div class="quick-actions">
            <el-button type="primary" :icon="DocumentAdd" @click="$router.push('/materials')">
              物资管理
            </el-button>
            <el-button type="success" :icon="Box" @click="$router.push('/stock')">
              库存查询
            </el-button>
            <el-button type="warning" :icon="ShoppingCart" @click="$router.push('/requisitions')">
              出库单
            </el-button>
            <el-button type="info" :icon="Upload" @click="$router.push('/inbound')">
              入库单
            </el-button>
            <el-button :icon="DataAnalysis" @click="generateReport">
              生成报告
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { projectApi, materialApi, stockApi, systemApi, aiApi } from '@/api'
import { ElMessage } from 'element-plus'
import {
  TrendCharts,
  Refresh,
  Box,
  Goods,
  Upload,
  Download,
  Warning,
  WarningFilled,
  CircleClose,
  MagicStick,
  CircleCheck,
  InfoFilled,
  PieChart,
  Operation,
  DocumentAdd,
  ShoppingCart,
  DataAnalysis
} from '@element-plus/icons-vue'

const router = useRouter()

// 日期范围
const dateRange = ref([])

// 关键指标
const metrics = reactive({
  totalMaterials: 0,
  totalStock: 0,
  monthlyInbound: 0,
  monthlyOutbound: 0
})

// 库存预警
const warnings = ref([])

// AI 洞察
const insights = ref([])

// 物资分类数据
const categoryData = ref([])

// 项目进度排行
const projectRanking = ref([])

// 加载状态
const loading = ref(false)

// 刷新数据
const refreshData = async () => {
  loading.value = true
  try {
    await Promise.all([
      fetchMetrics(),
      fetchWarnings(),
      fetchInsights(),
      fetchCategoryData(),
      fetchProjectRanking()
    ])
    ElMessage.success('数据已刷新')
  } catch (error) {
    console.error('刷新数据失败:', error)
  } finally {
    loading.value = false
  }
}

// 日期变化
const handleDateChange = () => {
  refreshData()
}

// 获取关键指标
const fetchMetrics = async () => {
  try {
    const { data } = await systemApi.getStats()
    Object.assign(metrics, {
      totalMaterials: data.total_materials || 0,
      totalStock: data.total_stock || 0,
      monthlyInbound: data.monthly_inbound || 0,
      monthlyOutbound: data.monthly_outbound || 0
    })
  } catch (error) {
    console.error('获取指标失败:', error)
  }
}

// 获取库存预警
const fetchWarnings = async () => {
  try {
    const { data } = await aiApi.getInsights({ type: 'stock_warnings' })
    if (data && data.alerts) {
      warnings.value = data.alerts
    }
  } catch (error) {
    console.error('获取预警失败:', error)
  }
}

// 获取 AI 洞察
const fetchInsights = async () => {
  try {
    const { recommendations } = await aiApi.getSuggestions({ type: 'insights' })
    if (recommendations && recommendations.length > 0) {
      insights.value = recommendations.map(rec => {
        let type = 'info'
        if (rec.type === 'alert') type = 'warning'
        else if (rec.type === 'success') type = 'success'

        return {
          type,
          title: rec.title,
          description: rec.description
        }
      })
    } else {
      // 默认示例
      insights.value = [
        {
          type: 'info',
          title: '数据分析',
          description: '本月物资消耗较上月下降 8%，整体运行平稳'
        },
        {
          type: 'success',
          title: '库存充足',
          description: '主要物资库存充足，可满足未来两周需求'
        }
      ]
    }
  } catch (error) {
    console.error('获取洞察失败:', error)
    // 默认示例
    insights.value = [
      {
        type: 'info',
        title: '数据分析',
        description: '本月物资消耗较上月下降 8%，整体运行平稳'
      }
    ]
  }
}

// 获取物资分类占比
const fetchCategoryData = async () => {
  try {
    const response = await materialApi.getCategories()
    const categories = response.data || []

    // 模拟分类占比数据
    const colors = ['#409eff', '#67c23a', '#e6a23c', '#f56c6c', '#909399']
    categoryData.value = categories.slice(0, 5).map((cat, index) => ({
      name: cat.name,
      value: Math.floor(Math.random() * 30) + 10,
      color: colors[index % colors.length]
    }))
  } catch (error) {
    console.error('获取分类失败:', error)
  }
}

// 获取项目进度排行
const fetchProjectRanking = async () => {
  try {
    const { data } = await projectApi.getList({ pageSize: 100 })
    const projects = data || []

    // 取前5个项目并生成随机进度
    projectRanking.value = projects.slice(0, 5).map(project => ({
      name: project.name,
      progress: Math.floor(Math.random() * 60) + 40
    })).sort((a, b) => b.progress - a.progress)
  } catch (error) {
    console.error('获取项目进度失败:', error)
  }
}

// 获取进度状态
const getProgressStatus = (progress) => {
  if (progress >= 80) return 'success'
  if (progress >= 60) return ''
  if (progress >= 40) return 'warning'
  return 'exception'
}

// 生成报告
const generateReport = () => {
  ElMessage.info('正在生成分析报告...')
}

// 初始化
onMounted(() => {
  refreshData()
})
</script>

<style scoped>
.ai-dashboard {
  padding: 0;
}

/* 页面标题 */
.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 0 0 20px 0;
  border-bottom: 1px solid #ebeef5;
}

.dashboard-header h1 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  display: flex;
  align-items: center;
  gap: 10px;
}

.dashboard-header .el-icon {
  color: #409eff;
}

.header-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

/* 指标卡片 */
.metrics-row {
  margin-bottom: 0;
}

.metric-card {
  background: white;
  border-radius: 12px;
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  transition: all 0.3s;
  cursor: pointer;
}

.metric-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
}

.metric-card.blue {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.metric-card.green {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  color: white;
}

.metric-card.orange {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  color: white;
}

.metric-card.red {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
  color: white;
}

.metric-icon {
  width: 64px;
  height: 64px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
}

.metric-content {
  flex: 1;
}

.metric-label {
  font-size: 14px;
  opacity: 0.9;
  margin-bottom: 8px;
}

.metric-value {
  font-size: 32px;
  font-weight: bold;
  line-height: 1;
}

/* 卡片通用样式 */
.dashboard-card {
  border-radius: 12px;
  margin-bottom: 20px;
}

.dashboard-card :deep(.el-card__header) {
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
}

.dashboard-card :deep(.el-card__body) {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
  font-size: 16px;
  color: #303133;
}

.card-header .el-icon {
  margin-right: 8px;
  color: #409eff;
}

/* 库存预警列表 */
.warning-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.warning-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: #fef0f0;
  border-radius: 8px;
  border-left: 4px solid #f56c6c;
}

.warning-warning {
  background: #fdf6ec;
  border-left-color: #e6a23c;
}

.warning-item .el-icon {
  font-size: 24px;
  color: #f56c6c;
}

.warning-warning .el-icon {
  color: #e6a23c;
}

.warning-content {
  flex: 1;
}

.warning-title {
  font-weight: 600;
  font-size: 14px;
  color: #303133;
  margin-bottom: 4px;
}

.warning-detail {
  font-size: 13px;
  color: #606266;
}

/* AI 洞察列表 */
.insight-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  max-height: 320px;
  overflow-y: auto;
}

.insight-item {
  display: flex;
  gap: 12px;
  padding: 16px;
  background: #f5f7fa;
  border-radius: 8px;
  border-left: 4px solid #909399;
}

.insight-warning {
  border-left-color: #e6a23c;
  background: #fdf6ec;
}

.insight-success {
  border-left-color: #67c23a;
  background: #f0f9ff;
}

.insight-info {
  border-left-color: #409eff;
  background: #ecf5ff;
}

.insight-icon {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.insight-warning .insight-icon {
  background: #faecd8;
  color: #e6a23c;
}

.insight-success .insight-icon {
  background: #e1f3d8;
  color: #67c23a;
}

.insight-info .insight-icon {
  background: #d9ecff;
  color: #409eff;
}

.insight-content {
  flex: 1;
}

.insight-title {
  font-weight: 600;
  font-size: 14px;
  color: #303133;
  margin-bottom: 6px;
}

.insight-desc {
  font-size: 13px;
  color: #606266;
  line-height: 1.6;
}

/* 分类占比列表 */
.category-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.category-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.category-info {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 120px;
}

.category-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.category-name {
  font-size: 14px;
  color: #606266;
}

.category-stats {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 12px;
}

.category-value {
  font-weight: 600;
  font-size: 14px;
  color: #303133;
  width: 50px;
  text-align: right;
}

.category-stats .el-progress {
  flex: 1;
}

/* 项目进度排行 */
.ranking-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.ranking-item {
  display: flex;
  align-items: center;
  gap: 16px;
}

.ranking-badge {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: 16px;
  background: #f5f7fa;
  color: #909399;
  flex-shrink: 0;
}

.ranking-badge.rank-1 {
  background: linear-gradient(135deg, #ffd666 0%, #ffec8b 100%);
  color: #d46b08;
}

.ranking-badge.rank-2 {
  background: linear-gradient(135deg, #c0c4cc 0%, #d3d4d6 100%);
  color: #606266;
}

.ranking-badge.rank-3 {
  background: linear-gradient(135deg, #e6a23c 0%, #f3d19e 100%);
  color: #ffffff;
}

.ranking-content {
  flex: 1;
}

.ranking-name {
  font-size: 14px;
  color: #303133;
  margin-bottom: 8px;
  font-weight: 500;
}

/* 快速操作 */
.quick-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.quick-actions .el-button {
  flex: 0 0 auto;
}
</style>
