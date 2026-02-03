<template>
  <div class="home-page">
    <!-- 头部欢迎 -->
    <div class="home-header">
      <div class="header-content">
        <div class="greeting">
          <p class="greeting-text">{{ greetingText }}</p>
          <h2 class="username">{{ displayName }}</h2>
          <p class="role-text">{{ userRoles }}</p>
        </div>
        <div class="header-icon" @click="goToProfile">
          <van-icon name="user-circle-o" size="48" color="#fff" />
        </div>
      </div>
    </div>

    <!-- 加载状态 -->
    <LoadingSkeleton v-if="loading" type="card" :count="4" />

    <!-- 内容区域 -->
    <template v-else>
      <!-- 统计卡片 -->
      <div class="stats-cards">
        <div class="stat-card" @click="goToTasks">
          <div class="stat-icon" style="background: linear-gradient(135deg, #ff6b6b 0%, #ee5a6f 100%);">
            <van-icon name="todo-list-o" size="24" color="white" />
          </div>
          <div class="stat-content">
            <p class="stat-value">{{ stats.pendingTasks }}</p>
            <p class="stat-label">待办任务</p>
          </div>
        </div>

        <div class="stat-card" @click="goToPlans">
          <div class="stat-icon" style="background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);">
            <van-icon name="orders-o" size="24" color="white" />
          </div>
          <div class="stat-content">
            <p class="stat-value">{{ stats.pendingPlans }}</p>
            <p class="stat-label">待审计划</p>
          </div>
        </div>

        <div class="stat-card" @click="goToInbound">
          <div class="stat-icon" style="background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);">
            <van-icon name="logistics" size="24" color="white" />
          </div>
          <div class="stat-content">
            <p class="stat-value">{{ stats.pendingInbound }}</p>
            <p class="stat-label">待审入库</p>
          </div>
        </div>

        <div class="stat-card" @click="goToOutbound">
          <div class="stat-icon" style="background: linear-gradient(135deg, #30cfd0 0%, #330867 100%);">
            <van-icon name="send-gift-o" size="24" color="white" />
          </div>
          <div class="stat-content">
            <p class="stat-value">{{ stats.pendingOutbound }}</p>
            <p class="stat-label">待审出库</p>
          </div>
        </div>
      </div>

      <!-- 图表区域 -->
      <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
        <div class="chart-section">
          <div class="section-header">
            <h3 class="section-title">数据概览</h3>
          </div>

          <!-- 物资计划状态饼图 -->
          <div class="chart-card">
            <h4 class="chart-title">物资计划状态分布</h4>
            <div class="pie-chart-container">
              <svg viewBox="0 0 200 200" class="pie-chart">
                <circle
                  v-for="(segment, index) in planStatusSegments"
                  :key="index"
                  cx="100"
                  cy="100"
                  :r="segment.radius"
                  :stroke="segment.color"
                  :stroke-width="30"
                  :stroke-dasharray="segment.dashArray"
                  :stroke-dashoffset="segment.offset"
                  fill="none"
                  transform="rotate(-90 100 100)"
                  class="pie-segment"
                />
                <text x="100" y="95" text-anchor="middle" class="pie-total">{{ stats.totalPlans }}</text>
                <text x="100" y="115" text-anchor="middle" class="pie-label">总数</text>
              </svg>
              <div class="pie-legend">
                <div v-for="(item, index) in planStatusLegend" :key="index" class="legend-item">
                  <span class="legend-color" :style="{ background: item.color }"></span>
                  <span class="legend-text">{{ item.label }}: {{ item.value }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- 本月物资计划趋势柱状图 -->
          <div class="chart-card">
            <h4 class="chart-title">本月物资计划趋势</h4>
            <div class="bar-chart-container">
              <svg viewBox="0 0 300 150" class="bar-chart">
                <!-- Y轴网格线 -->
                <line v-for="i in 5" :key="'grid-' + i"
                  :x1="30" :y1="20 + i * 20"
                  :x2="290" :y2="20 + i * 20"
                  stroke="#e0e0e0" stroke-width="1"
                />
                <!-- Y轴标签 -->
                <text v-for="i in 6" :key="'label-' + i"
                  :x="20" :y="25 + (6 - i) * 20"
                  text-anchor="end"
                  class="bar-label"
                >{{ (i - 1) * 2 }}</text>
                <!-- 柱状图 -->
                <rect
                  v-for="(bar, index) in monthlyTrend"
                  :key="index"
                  :x="35 + index * 40"
                  :y="120 - bar.value * 10"
                  width="30"
                  :height="bar.value * 10"
                  :fill="bar.color"
                  rx="4"
                  class="bar-rect"
                />
                <!-- X轴标签 -->
                <text
                  v-for="(bar, index) in monthlyTrend"
                  :key="'x-' + index"
                  :x="50 + index * 40"
                  :y="140"
                  text-anchor="middle"
                  class="bar-label"
                >{{ bar.label }}</text>
              </svg>
            </div>
          </div>

          <!-- 库存预警 -->
          <div class="chart-card">
            <h4 class="chart-title">库存预警</h4>
            <div class="alert-list">
              <div
                v-for="(alert, index) in stockAlerts"
                :key="index"
                class="alert-item"
                @click="goToStock"
              >
                <div class="alert-icon" :class="alert.level">
                  <van-icon :name="alert.icon" size="20" />
                </div>
                <div class="alert-content">
                  <p class="alert-name">{{ alert.name }}</p>
                  <p class="alert-info">库存: {{ alert.current }} / 安全库存: {{ alert.safety }}</p>
                </div>
                <van-icon name="arrow" color="#969799" />
              </div>
              <EmptyState
                v-if="stockAlerts.length === 0"
                text="库存状态良好"
                :show-action="false"
              />
            </div>
          </div>
        </div>
      </van-pull-refresh>

      <!-- 快捷功能 -->
      <div class="quick-actions">
        <div class="section-header">
          <h3 class="section-title">快捷功能</h3>
        </div>
        <div class="action-grid">
          <div class="action-item" @click="goToPlans">
            <div class="action-icon" style="background: #e8f4ff;">
              <van-icon name="orders-o" size="28" color="#1989fa" />
            </div>
            <span class="action-label">物资计划</span>
          </div>
          <div class="action-item" @click="goToInbound">
            <div class="action-icon" style="background: #fff0f0;">
              <van-icon name="logistics" size="28" color="#ee0a24" />
            </div>
            <span class="action-label">入库管理</span>
          </div>
          <div class="action-item" @click="goToOutbound">
            <div class="action-icon" style="background: #f0f9ff;">
              <van-icon name="send-gift-o" size="28" color="#ff976a" />
            </div>
            <span class="action-label">出库管理</span>
          </div>
          <div class="action-item" @click="goToMaterials">
            <div class="action-icon" style="background: #fff7e6;">
              <van-icon name="apps-o" size="28" color="#ff976a" />
            </div>
            <span class="action-label">物资浏览</span>
          </div>
          <div class="action-item" @click="goToStock">
            <div class="action-icon" style="background: #f0fff4;">
              <van-icon name="shop-o" size="28" color="#07c160" />
            </div>
            <span class="action-label">库存查询</span>
          </div>
          <div class="action-item" @click="goToTasks">
            <div class="action-icon" style="background: #fff0ff;">
              <van-icon name="todo-list-o" size="28" color="#7232dd" />
            </div>
            <span class="action-label">待办任务</span>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'
import { useUserStore } from '@/stores/user'
import { useErrorStore } from '@/stores/error'
import { getDashboardData } from '@/api/system'
import { useAPI } from '@/composables/useAPI'
import LoadingSkeleton from '@/components/common/LoadingSkeleton.vue'
import EmptyState from '@/components/common/EmptyState.vue'

const router = useRouter()
const userStore = useUserStore()
const errorStore = useErrorStore()

const stats = ref({
  pendingTasks: 0,
  pendingPlans: 0,
  pendingInbound: 0,
  pendingOutbound: 0,
  totalPlans: 0,
  planStatus: {
    draft: 0,
    pending: 0,
    active: 0,
    completed: 0
  }
})

const monthlyTrend = ref([])
const stockAlerts = ref([])
const loading = ref(true)
const refreshing = ref(false)

// 问候语
const greetingText = computed(() => {
  const hour = new Date().getHours()
  if (hour < 6) return '夜深了'
  if (hour < 9) return '早上好'
  if (hour < 12) return '上午好'
  if (hour < 14) return '中午好'
  if (hour < 18) return '下午好'
  if (hour < 22) return '晚上好'
  return '夜深了'
})

// 显示名称
const displayName = computed(() => {
  return userStore.userInfo?.username || userStore.userInfo?.full_name || '用户'
})

// 用户角色
const userRoles = computed(() => {
  const roles = userStore.roles || []
  if (roles.includes('admin')) return '管理员'
  if (roles.includes('project_manager')) return '项目经理'
  if (roles.includes('warehouse_manager')) return '仓库管理员'
  if (roles.includes('worker')) return '施工人员'
  return '普通用户'
})

// 物资计划状态饼图数据
const planStatusSegments = computed(() => {
  const data = stats.value.planStatus
  const total = data.draft + data.pending + data.active + data.completed
  if (total === 0) return []

  const colors = ['#969799', '#ff976a', '#1989fa', '#07c160']
  const labels = ['draft', 'pending', 'active', 'completed']
  let offset = 0
  const segments = []

  labels.forEach((key, index) => {
    const value = data[key]
    if (value > 0) {
      const percentage = (value / total) * 100
      const circumference = 2 * Math.PI * 70
      const dashArray = `${(percentage / 100) * circumference} ${circumference}`
      segments.push({
        radius: 70,
        color: colors[index],
        dashArray,
        offset: -offset
      })
      offset += (percentage / 100) * circumference
    }
  })

  return segments
})

// 饼图图例
const planStatusLegend = computed(() => {
  const data = stats.value.planStatus
  return [
    { label: '草稿', value: data.draft, color: '#969799' },
    { label: '待审批', value: data.pending, color: '#ff976a' },
    { label: '进行中', value: data.active, color: '#1989fa' },
    { label: '已完成', value: data.completed, color: '#07c160' }
  ].filter(item => item.value > 0)
})

// 使用useAPI加载数据
const { execute: loadDashboardData, loading: apiLoading } = useAPI(
  getDashboardData,
  {
    immediate: false,
    showToast: false,
    onSuccess: (data) => {
      stats.value = {
        pendingTasks: data.pending_tasks || 0,
        pendingPlans: data.pending_plans || 0,
        pendingInbound: data.pending_inbound || 0,
        pendingOutbound: data.pending_outbound || 0,
        totalPlans: data.total_plans || 0,
        planStatus: {
          draft: data.plan_status?.draft || 0,
          pending: data.plan_status?.pending || 0,
          active: data.plan_status?.active || 0,
          completed: data.plan_status?.completed || 0
        }
      }

      monthlyTrend.value = (data.monthly_trend || []).map((item, index) => ({
        label: `${index * 5 + 1}-${index * 5 + 5}日`,
        value: item.count,
        color: '#1989fa'
      }))

      stockAlerts.value = (data.stock_alerts || []).map(item => ({
        name: item.material_name,
        current: item.current_quantity,
        safety: item.safety_stock,
        level: item.current_quantity < item.safety_stock * 0.5 ? 'high' : 'medium',
        icon: item.current_quantity < item.safety_stock * 0.5 ? 'warning-o' : 'info-o'
      }))

      loading.value = false
    },
    onError: (error) => {
      errorStore.handleApiError(error, { context: 'Home Dashboard' })
      loading.value = false
      // 使用模拟数据作为降级方案
      loadMockData()
    }
  }
)

// 加载模拟数据（降级方案）
function loadMockData() {
  stats.value = {
    pendingTasks: 5,
    pendingPlans: 3,
    pendingInbound: 2,
    pendingOutbound: 4,
    totalPlans: 24,
    planStatus: {
      draft: 2,
      pending: 3,
      active: 8,
      completed: 11
    }
  }

  monthlyTrend.value = [
    { label: '1-5', value: 3, color: '#1989fa' },
    { label: '6-10', value: 5, color: '#1989fa' },
    { label: '11-15', value: 2, color: '#1989fa' },
    { label: '16-20', value: 6, color: '#1989fa' },
    { label: '21-25', value: 4, color: '#1989fa' },
    { label: '26-31', value: 4, color: '#1989fa' }
  ]

  stockAlerts.value = [
    { name: '钢筋 Φ12', current: 15, safety: 50, level: 'high', icon: 'warning-o' },
    { name: '水泥 42.5', current: 30, safety: 100, level: 'high', icon: 'warning-o' },
    { name: '砂石', current: 80, safety: 200, level: 'medium', icon: 'info-o' }
  ]
}

// 下拉刷新
async function onRefresh() {
  try {
    await loadDashboardData()
    showToast('刷新成功')
  } finally {
    refreshing.value = false
  }
}

// 导航方法
function goToProfile() { router.push('/profile') }
function goToTasks() { router.push('/tasks') }
function goToPlans() { router.push('/plans') }
function goToInbound() { router.push('/inbound') }
function goToOutbound() { router.push('/outbound') }
function goToMaterials() { router.push('/materials') }
function goToStock() { router.push('/stock') }

let refreshTimer = null

onMounted(async () => {
  await loadDashboardData()
  // 每30秒自动刷新数据
  refreshTimer = setInterval(() => {
    loadDashboardData()
  }, 30000)
})

onUnmounted(() => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
  }
})
</script>

<style scoped>
.home-page {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding-bottom: 60px;
}

.home-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 40px 20px 30px;
  color: white;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.greeting-text {
  font-size: 14px;
  opacity: 0.9;
  margin: 0 0 8px 0;
}

.username {
  font-size: 26px;
  font-weight: bold;
  margin: 0 0 4px 0;
}

.role-text {
  font-size: 14px;
  opacity: 0.8;
  margin: 0;
}

.header-icon {
  cursor: pointer;
  transition: transform 0.2s;
}

.header-icon:active {
  transform: scale(0.95);
}

.stats-cards {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  padding: 16px;
  margin-top: -20px;
}

.stat-card {
  background: white;
  border-radius: 12px;
  padding: 16px;
  display: flex;
  align-items: center;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  cursor: pointer;
  transition: all 0.3s;
}

.stat-card:active {
  transform: scale(0.98);
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12px;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 22px;
  font-weight: bold;
  color: #323233;
  margin: 0 0 2px 0;
}

.stat-label {
  font-size: 12px;
  color: #969799;
  margin: 0;
}

.chart-section {
  padding: 16px;
}

.section-header {
  margin-bottom: 12px;
}

.section-title {
  font-size: 18px;
  font-weight: bold;
  color: #323233;
  margin: 0;
}

.chart-card {
  background: white;
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.chart-title {
  font-size: 15px;
  font-weight: 600;
  color: #323233;
  margin: 0 0 16px 0;
}

/* 饼图样式 */
.pie-chart-container {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.pie-chart {
  width: 160px;
  height: 160px;
}

.pie-segment {
  transition: opacity 0.3s;
}

.pie-segment:hover {
  opacity: 0.8;
}

.pie-total {
  font-size: 28px;
  font-weight: bold;
  fill: #323233;
}

.pie-label {
  font-size: 12px;
  fill: #969799;
}

.pie-legend {
  margin-top: 16px;
  width: 100%;
}

.legend-item {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
  font-size: 13px;
}

.legend-color {
  width: 12px;
  height: 12px;
  border-radius: 2px;
  margin-right: 8px;
}

.legend-text {
  color: #646566;
}

/* 柱状图样式 */
.bar-chart-container {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.bar-chart {
  width: 100%;
  height: 150px;
}

.bar-rect {
  transition: opacity 0.3s;
}

.bar-rect:hover {
  opacity: 0.8;
}

.bar-label {
  font-size: 10px;
  fill: #969799;
}

/* 库存预警 */
.alert-list {
  max-height: 200px;
  overflow-y: auto;
}

.alert-item {
  display: flex;
  align-items: center;
  padding: 12px;
  background: #f7f8fa;
  border-radius: 8px;
  margin-bottom: 8px;
  cursor: pointer;
  transition: background 0.2s;
}

.alert-item:active {
  background: #f0f0f0;
}

.alert-item:last-child {
  margin-bottom: 0;
}

.alert-icon {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12px;
  flex-shrink: 0;
}

.alert-icon.high {
  background: #ffe0e0;
  color: #ee0a24;
}

.alert-icon.medium {
  background: #fff7e6;
  color: #ff976a;
}

.alert-content {
  flex: 1;
  min-width: 0;
}

.alert-name {
  font-size: 14px;
  font-weight: 500;
  color: #323233;
  margin: 0 0 2px 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.alert-info {
  font-size: 12px;
  color: #969799;
  margin: 0;
}

/* 快捷功能 */
.quick-actions {
  padding: 0 16px 16px;
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}

.action-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px 8px;
  background: white;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s;
}

.action-item:active {
  transform: scale(0.95);
}

.action-icon {
  width: 52px;
  height: 52px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 8px;
}

.action-label {
  font-size: 12px;
  color: #646566;
}
</style>
