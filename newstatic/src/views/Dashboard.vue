<template>
  <div class="dashboard">
    <!-- 欢迎区域 -->
    <div class="welcome-section">
      <div class="welcome-content">
        <h1 class="welcome-title">欢迎回来，{{ displayName }}</h1>
        <p class="welcome-subtitle">{{ currentDate }}</p>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-section">
      <el-row :gutter="20">
        <el-col :xs="12" :sm="8" :md="6" v-for="stat in statistics" :key="stat.title">
          <div class="stat-card" :class="stat.status">
            <div class="stat-icon" :style="{ background: stat.color }">
              <el-icon :size="28" color="white">
                <component :is="stat.icon" />
              </el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value" v-loading="stat.loading">{{ stat.value }}</div>
              <div class="stat-label">{{ stat.title }}</div>
              <div class="stat-trend" v-if="stat.trend">
                <el-icon><TrendCharts /></el-icon>
                <span>{{ stat.trendText }}</span>
              </div>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 主要内容区域 -->
    <el-row :gutter="20" class="content-section">
      <!-- 快捷操作 -->
      <el-col :xs="24" :lg="16">
        <div class="card quick-actions-card">
          <div class="card-header">
            <h3 class="card-title">快捷操作</h3>
            <el-button type="primary" text @click="viewAllActions">
              查看全部
            </el-button>
          </div>
          <div class="actions-grid">
            <div
              v-for="action in quickActions"
              :key="action.key"
              class="action-item"
              @click="handleAction(action)"
            >
              <div class="action-icon" :style="{ background: action.color }">
                <el-icon :size="24" color="white">
                  <component :is="action.icon" />
                </el-icon>
              </div>
              <div class="action-info">
                <div class="action-name">{{ action.name }}</div>
                <div class="action-desc">{{ action.desc }}</div>
              </div>
            </div>
          </div>
        </div>
      </el-col>

      <!-- 最近活动 -->
      <el-col :xs="24" :lg="8">
        <div class="card activity-card">
          <div class="card-header">
            <h3 class="card-title">最近活动</h3>
            <el-button type="primary" text @click="viewAllLogs">查看全部</el-button>
          </div>
          <div class="activity-list" v-loading="logsLoading">
            <div
              v-for="activity in recentActivities"
              :key="activity.id"
              class="activity-item"
            >
              <div class="activity-icon" :class="activity.type">
                <el-icon>
                  <component :is="getActivityIcon(activity.type)" />
                </el-icon>
              </div>
              <div class="activity-content">
                <div class="activity-title">{{ activity.title }}</div>
                <div class="activity-time">{{ activity.time }}</div>
              </div>
            </div>
            <el-empty
              v-if="recentActivities.length === 0 && !logsLoading"
              description="暂无活动"
              :image-size="80"
            />
          </div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import {
  Box,
  ShoppingCart,
  Document,
  DocumentCopy,
  List,
  DataAnalysis,
  Setting,
  Calendar,
  Flag,
  Bell,
  TrendCharts,
  Plus,
  Upload,
  Download,
  Edit,
  Delete,
  Management
} from '@element-plus/icons-vue'
import { systemApi } from '@/api'

const router = useRouter()
const authStore = useAuthStore()

// 当前日期
const currentDate = computed(() => {
  const now = new Date()
  const options = {
    year: 'numeric',
    month: 'long',
    weekday: 'long',
    day: 'numeric'
  }
  return now.toLocaleDateString('zh-CN', options)
})

// 显示名称
const displayName = computed(() => authStore.displayName || '用户')

// 统计数据
const statistics = ref([
  {
    title: '物资总数',
    value: '-',
    icon: Box,
    color: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    status: 'success',
    trend: 'up',
    trendText: '+12%'
  },
  {
    title: '库存预警',
    value: '-',
    icon: ShoppingCart,
    color: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
    status: 'warning',
    trend: 'down',
    trendText: '较上月'
  },
  {
    title: '待处理出库',
    value: '-',
    icon: Document,
    color: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
    status: 'info',
    trend: 'neutral',
    trendText: '待处理'
  },
  {
    title: '进行中计划',
    value: '-',
    icon: Flag,
    color: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
    status: 'primary',
    trend: 'up',
    trendText: '+5%'
  }
])

// 快捷操作配置
const quickActions = ref([
  {
    key: 'create_material',
    name: '新增物资',
    desc: '添加新的物资信息',
    icon: Plus,
    color: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    path: '/materials',
    action: 'create'
  },
  {
    key: 'create_inbound',
    name: '物资入库',
    desc: '创建入库单据',
    icon: Upload,
    color: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
    path: '/inbound',
    action: 'create'
  },
  {
    key: 'create_requisition',
    name: '创建出库单',
    desc: '创建出库申请',
    icon: DocumentCopy,
    color: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
    path: '/requisitions',
    action: 'create'
  },
  {
    key: 'create_plan',
    name: '创建计划',
    desc: '新建物资计划',
    icon: List,
    color: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
    path: '/material-plans',
    action: 'create'
  },
  {
    key: 'manage_stock',
    name: '库存管理',
    desc: '查看和调整库存',
    icon: Management,
    color: 'linear-gradient(135deg, #30cfd0 0%, #330867 100%)',
    path: '/stock',
    action: 'view'
  },
  {
    key: 'view_logs',
    name: '系统日志',
    desc: '查看系统操作日志',
    icon: Bell,
    color: 'linear-gradient(135deg, #a8edea 0%, #fed6e3 100%)',
    path: '/system',
    action: 'logs'
  }
])

// 最近活动数据
const recentActivities = ref([])
const logsLoading = ref(false)

// 获取统计数据
const fetchStatistics = async () => {
  try {
    const response = await systemApi.getStats()
    const data = response.data || response

    // 更新统计数据
    if (data.total_materials !== undefined) {
      statistics.value[0].value = data.total_materials.toLocaleString()
    }
    if (data.low_stock_count !== undefined) {
      statistics.value[1].value = data.low_stock_count.toLocaleString()
    }
    if (data.pending_requisitions !== undefined) {
      statistics.value[2].value = data.pending_requisitions.toLocaleString()
    }
    if (data.total_users !== undefined) {
      statistics.value[3].value = data.total_users.toLocaleString()
    }
  } catch (error) {
    console.error('获取统计数据失败:', error)
    // 保持默认值
  }
}

// 获取最近活动
const fetchRecentActivities = async () => {
  logsLoading.value = true
  try {
    const { data } = await systemApi.getLogs({
      page: 1,
      page_size: 5
    })

    const logs = data || []

    recentActivities.value = logs.slice(0, 5).map((log, index) => ({
      id: index,
      type: log.status === 'success' || log.result === '成功' ? 'success' : 'error',
      title: log.action || log.description || log.operation || '系统操作',
      time: formatTime(log.created_at || log.timestamp || ''),
      raw: log
    }))
  } catch (error) {
    console.error('获取最近活动失败:', error)
  } finally {
    logsLoading.value = false
  }
}

// 格式化时间
const formatTime = (timeStr) => {
  if (!timeStr) return ''

  try {
    const date = new Date(timeStr)
    const now = new Date()
    const diff = now - date

    if (diff < 60000) { // 1分钟内
      return '刚刚'
    } else if (diff < 3600000) { // 1小时内
      const minutes = Math.floor(diff / 60000)
      return `${minutes}分钟前`
    } else if (diff < 86400000) { // 24小时内
      const hours = Math.floor(diff / 3600000)
      return `${hours}小时前`
    } else {
      return date.toLocaleDateString('zh-CN')
    }
  } catch (error) {
    return timeStr
  }
}

// 获取活动图标
const getActivityIcon = (type) => {
  const iconMap = {
    success: 'CircleCheck',
    error: 'CircleClose',
    info: 'InfoFilled',
    warning: 'Warning'
  }
  return iconMap[type] || 'InfoFilled'
}

// 处理快捷操作
const handleAction = (action) => {
  if (action.path) {
    router.push({
      path: action.path,
      query: action.action ? { action: action.action } : undefined
    })
  }
}

// 查看全部快捷操作
const viewAllActions = () => {
  // 可以打开一个对话框显示所有快捷操作
  console.log('查看全部快捷操作')
}

// 查看全部日志
const viewAllLogs = () => {
  router.push('/system')
}

onMounted(() => {
  fetchStatistics()
  fetchRecentActivities()
})
</script>

<style scoped>
.dashboard {
  padding: 0;
}

/* 欢迎区域 */
.welcome-section {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  padding: 30px;
  margin-bottom: 24px;
  color: white;
}

.welcome-title {
  font-size: 24px;
  font-weight: 600;
  margin: 0 0 8px 0;
}

.welcome-subtitle {
  font-size: 14px;
  opacity: 0.9;
  margin: 0;
}

/* 统计区域 */
.stats-section {
  margin-bottom: 24px;
}

.stat-card {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
  border-left: 4px solid transparent;
  cursor: pointer;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.12);
}

.stat-card.success {
  border-left-color: #67c23a;
}

.stat-card.warning {
  border-left-color: #e6a23c;
}

.stat-card.info {
  border-left-color: #409eff;
}

.stat-card.primary {
  border-left-color: #667eea;
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16px;
}

.stat-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  color: #303133;
  margin-right: 16px;
}

.stat-label {
  font-size: 14px;
  color: #606266;
  margin-right: auto;
}

.stat-trend {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #909399;
}

.stat-trend .el-icon {
  font-size: 16px;
}

/* 主内容区域 */
.content-section {
  margin-bottom: 24px;
}

/* 卡片通用样式 */
.card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  height: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin: 0;
}

/* 快捷操作卡片 */
.quick-actions-card {
  height: 100%;
}

.actions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
  padding: 20px 0;
}

.action-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.action-item:hover {
  background-color: #f5f7fa;
}

.action-icon {
  width: 48px;
  height: 48px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.action-info {
  flex: 1;
}

.action-name {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
  margin-bottom: 4px;
}

.action-desc {
  font-size: 12px;
  color: #909399;
}

/* 最近活动卡片 */
.activity-card {
  height: 100%;
}

.activity-list {
  padding: 20px 0;
  max-height: 400px;
  overflow-y: auto;
}

.activity-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px 0;
  border-bottom: 1px solid #f0f0f0;
}

.activity-item:last-child {
  border-bottom: none;
}

.activity-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.activity-icon.success {
  background-color: #f0f9ff;
  color: #67c23a;
}

.activity-icon.error {
  background-color: #fef0f0;
  color: #f56c6c;
}

.activity-icon.info {
  background-color: #e1f3ff;
  color: #409eff;
}

.activity-content {
  flex: 1;
}

.activity-title {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
  margin-bottom: 4px;
}

.activity-time {
  font-size: 12px;
  color: #909399;
}

/* 响应式 */
@media (max-width: 768px) {
  .welcome-section {
    padding: 20px;
  }

  .welcome-title {
    font-size: 20px;
  }

  .actions-grid {
    grid-template-columns: 1fr;
  }

  .stat-card {
    padding: 20px;
  }
}
</style>
