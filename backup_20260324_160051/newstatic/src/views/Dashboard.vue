<template>
  <div class="dashboard">
    <!-- 欢迎区域 -->
    <div class="welcome-section">
      <div class="welcome-content">
        <h1 class="welcome-title">欢迎回来，{{ displayName }}</h1>
        <p class="welcome-subtitle">{{ currentDate }}</p>
      </div>
    </div>

    <!-- 统计卡片 - 根据权限过滤 -->
    <div class="stats-section" v-if="visibleStatistics.length > 0">
      <el-row :gutter="20">
        <el-col :xs="12" :sm="8" :md="6" v-for="stat in visibleStatistics" :key="stat.title">
          <div class="stat-card" :class="stat.status">
            <div class="stat-icon" :style="{ background: stat.color }">
              <el-icon :size="28" color="white">
                <component :is="stat.icon" />
              </el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value" v-loading="stat.loading">{{ stat.value }}</div>
              <div class="stat-label">{{ stat.title }}</div>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 预约统计卡片 -->
    <div class="stats-section" v-if="authStore.hasPermission('appointment_view')">
      <div class="section-header">
        <h3 class="section-title">施工预约统计</h3>
      </div>
      <el-row :gutter="20" class="mt-16">
        <el-col :xs="12" :sm="6" :md="6" v-for="stat in appointmentStats" :key="stat.title">
          <div class="stat-card" :class="stat.status" @click="handleStatClick(stat)">
            <div class="stat-icon" :style="{ background: stat.color }">
              <el-icon :size="28" color="white">
                <component :is="stat.icon" />
              </el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ stat.value }}</div>
              <div class="stat-label">{{ stat.title }}</div>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 主要内容区域 -->
    <el-row :gutter="20" class="content-section">
      <!-- 快捷操作 -->
      <el-col :xs="24" :lg="24">
        <div class="card quick-actions-card">
          <div class="card-header">
            <h3 class="card-title">快捷操作</h3>
            <el-button type="primary" text @click="viewAllActions">
              查看全部
            </el-button>
          </div>
          <div class="actions-grid" v-if="visibleQuickActions.length > 0">
            <div
              v-for="action in visibleQuickActions"
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
          <div v-else class="empty-actions">
            <p>暂无可用的快捷操作</p>
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
  Plus,
  Upload,
  Download,
  Edit,
  Delete,
  Management,
  Clock,
  User,
  Check
} from '@element-plus/icons-vue'
import { systemApi, appointmentApi } from '@/api'

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
    permissions: ['material_view']
  },
  {
    title: '项目总数',
    value: '-',
    icon: Flag,
    color: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
    status: 'primary',
    permissions: ['project_view']
  },
  {
    title: '库存预警',
    value: '-',
    icon: ShoppingCart,
    color: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
    status: 'warning',
    permissions: ['stock_view', 'stock_alerts']
  },
  {
    title: '待处理出库',
    value: '-',
    icon: Document,
    color: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
    status: 'info',
    permissions: ['requisition_view']
  },
  {
    title: '本月入库',
    value: '-',
    icon: Upload,
    color: 'linear-gradient(135deg, #30cfd0 0%, #330867 100%)',
    status: 'success',
    permissions: ['inbound_view']
  },
  {
    title: '本月出库',
    value: '-',
    icon: Download,
    color: 'linear-gradient(135deg, #a8edea 0%, #fed6e3 100%)',
    status: 'primary',
    permissions: ['requisition_view']
  }
])

// 可见统计数据（根据权限过滤）
const visibleStatistics = computed(() => {
  return statistics.value.filter(stat => {
    if (!stat.permissions || stat.permissions.length === 0) return true
    return stat.permissions.some(perm => authStore.hasPermission(perm))
  })
})

// 预约统计数据
const appointmentStats = ref([
  {
    title: '全部预约',
    value: '-',
    icon: Calendar,
    color: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    status: 'primary'
  },
  {
    title: '待审批',
    value: '-',
    icon: Clock,
    color: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
    status: 'warning'
  },
  {
    title: '已排期',
    value: '-',
    icon: Check,
    color: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
    status: 'success'
  },
  {
    title: '进行中',
    value: '-',
    icon: User,
    color: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
    status: 'info'
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
    action: 'create',
    permissions: ['material_create']
  },
  {
    key: 'create_inbound',
    name: '物资入库',
    desc: '创建入库单据',
    icon: Upload,
    color: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
    path: '/inbound',
    action: 'create',
    permissions: ['inbound_create']
  },
  {
    key: 'create_requisition',
    name: '创建出库单',
    desc: '创建出库申请',
    icon: DocumentCopy,
    color: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
    path: '/requisitions',
    action: 'create',
    permissions: ['requisition_create']
  },
  {
    key: 'create_plan',
    name: '创建计划',
    desc: '新建物资计划',
    icon: List,
    color: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
    path: '/material-plans',
    action: 'create',
    permissions: ['material_plan_create']
  },
  {
    key: 'manage_stock',
    name: '库存管理',
    desc: '查看和调整库存',
    icon: Management,
    color: 'linear-gradient(135deg, #30cfd0 0%, #330867 100%)',
    path: '/stock',
    action: 'view',
    permissions: ['stock_view']
  },
  {
    key: 'view_logs',
    name: '系统日志',
    desc: '查看系统操作日志',
    icon: Bell,
    color: 'linear-gradient(135deg, #a8edea 0%, #fed6e3 100%)',
    path: '/operation-logs',
    action: 'view',
    permissions: ['audit_view']
  }
])

// 可见快捷操作（根据权限过滤）
const visibleQuickActions = computed(() => {
  return quickActions.value.filter(action => {
    if (!action.permissions || action.permissions.length === 0) return true
    return action.permissions.some(perm => authStore.hasPermission(perm))
  })
})

// 获取统计数据
const fetchStatistics = async () => {
  try {
    const response = await systemApi.getStats()
    const data = response.data || response

    // 更新统计数据 - 按照新的6个统计卡片顺序
    if (data.total_materials !== undefined) {
      statistics.value[0].value = data.total_materials.toLocaleString()
    }
    if (data.total_projects !== undefined) {
      statistics.value[1].value = data.total_projects.toLocaleString()
    }
    if (data.low_stock_count !== undefined) {
      statistics.value[2].value = data.low_stock_count.toLocaleString()
    }
    if (data.pending_requisitions !== undefined) {
      statistics.value[3].value = data.pending_requisitions.toLocaleString()
    }
    if (data.monthly_inbound !== undefined) {
      statistics.value[4].value = data.monthly_inbound.toLocaleString()
    }
    if (data.monthly_requisitions !== undefined) {
      statistics.value[5].value = data.monthly_requisitions.toLocaleString()
    }
  } catch (error) {
    console.error('获取统计数据失败:', error)
    // 保持默认值
  }
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

// 获取预约统计数据
const fetchAppointmentStats = async () => {
  try {
    const response = await appointmentApi.getStats()
    const data = response.data || {}

    // 更新预约统计数据
    if (data.total !== undefined) {
      appointmentStats.value[0].value = data.total.toLocaleString()
    }
    if (data.pending !== undefined) {
      appointmentStats.value[1].value = data.pending.toLocaleString()
    }
    if (data.scheduled !== undefined) {
      appointmentStats.value[2].value = data.scheduled.toLocaleString()
    }
    if (data.in_progress !== undefined) {
      appointmentStats.value[3].value = data.in_progress.toLocaleString()
    }
  } catch (error) {
    console.error('获取预约统计失败:', error)
  }
}

// 处理统计卡片点击
const handleStatClick = (stat) => {
  router.push({
    path: '/appointments',
    query: { status: stat.title === '待审批' ? 'pending' : stat.title === '已排期' ? 'scheduled' : stat.title === '进行中' ? 'in_progress' : '' }
  })
}

onMounted(() => {
  fetchStatistics()
  if (authStore.hasPermission('appointment_view')) {
    fetchAppointmentStats()
  }
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

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin: 0;
}

.mt-16 {
  margin-top: 16px;
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

/* 空状态 */
.empty-actions {
  padding: 40px 20px;
  text-align: center;
  color: #909399;
}

.empty-actions p {
  margin: 0;
  font-size: 14px;
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
