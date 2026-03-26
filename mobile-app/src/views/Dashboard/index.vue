<template>
  <div class="dashboard">
    <!-- 顶部欢迎卡片 -->
    <div class="welcome-card">
      <div class="welcome-content">
        <div class="welcome-text">
          <div class="greeting">{{ greeting }}</div>
          <div class="username">{{ displayName }}</div>
        </div>
        <div class="welcome-icon">
          <NotificationBadge />
        </div>
      </div>
      <div class="welcome-date">{{ currentDate }}</div>
    </div>

    <!-- 统计数据卡片 -->
    <div class="stats-section">
      <div class="section-title">
        <van-icon name="bar-chart-o" />
        <span>数据概览</span>
      </div>
      <div class="stats-grid">
        <div class="stat-card" @click="goToAppointments">
          <div class="stat-icon" style="background-color: #fff3e0;">
            <van-icon name="calendar-o" size="24" color="#ff9800" />
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.todayAppointments || 0 }}</div>
            <div class="stat-label">今日任务</div>
          </div>
        </div>

        <div class="stat-card" @click="goToAppointments">
          <div class="stat-icon" style="background-color: #e3f2fd;">
            <van-icon name="todo-list-o" size="24" color="#1989fa" />
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.pendingAppointments || 0 }}</div>
            <div class="stat-label">待执行任务</div>
          </div>
        </div>

        <div class="stat-card" @click="goToClockIn">
          <div class="stat-icon" style="background-color: #e8f5e9;">
            <van-icon name="clock-o" size="24" color="#07c160" />
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.todayClockIns || 0 }}</div>
            <div class="stat-label">今日打卡</div>
          </div>
        </div>

        <div class="stat-card" @click="goToRecords">
          <div class="stat-icon" style="background-color: #f3e5f5;">
            <van-icon name="records" size="24" color="#7232dd" />
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.monthClockIns || 0 }}</div>
            <div class="stat-label">本月打卡</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 快捷功能 -->
    <div class="quick-actions">
      <div class="section-title">
        <van-icon name="apps-o" />
        <span>快捷功能</span>
      </div>
      <van-grid :column-num="4" :border="false" :gutter="12">
        <van-grid-item @click="router.push('/appointments/calendar')">
          <div class="quick-icon" style="background: linear-gradient(135deg, #ff9800 0%, #f57c00 100%);">
            <van-icon name="calendar-o" size="24" color="#fff" />
          </div>
          <div class="quick-label">预约管理</div>
        </van-grid-item>

        <van-grid-item @click="router.push('/attendance/clock-in')">
          <div class="quick-icon" style="background: linear-gradient(135deg, #07c160 0%, #06ad56 100%);">
            <van-icon name="clock-o" size="24" color="#fff" />
          </div>
          <div class="quick-label">打卡</div>
        </van-grid-item>

        <van-grid-item @click="router.push('/attendance/records')">
          <div class="quick-icon" style="background: linear-gradient(135deg, #7232dd 0%, #5a17a9 100%);">
            <van-icon name="records" size="24" color="#fff" />
          </div>
          <div class="quick-label">打卡记录</div>
        </van-grid-item>

        <!-- 物资计划 - 根据权限显示 -->
        <van-grid-item v-if="canAccessMaterialPlan" @click="router.push('/plans')">
          <div class="quick-icon" style="background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);">
            <van-icon name="todo-list-o" size="24" color="#fff" />
          </div>
          <div class="quick-label">物资计划</div>
        </van-grid-item>

        <!-- 新建计划 - 根据权限显示 -->
        <van-grid-item v-if="canCreatePlan" @click="router.push('/plans/create')">
          <div class="quick-icon" style="background: linear-gradient(135deg, #e91e63 0%, #9c27b0 100%);">
            <van-icon name="plus" size="24" color="#fff" />
          </div>
          <div class="quick-label">新建计划</div>
        </van-grid-item>

        <!-- 入库管理 - 根据权限显示 -->
        <van-grid-item v-if="canAccessInbound" @click="router.push('/inbound')">
          <div class="quick-icon" style="background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);">
            <van-icon name="logistics" size="24" color="#fff" />
          </div>
          <div class="quick-label">入库管理</div>
        </van-grid-item>

        <!-- 新建入库 - 根据权限显示 -->
        <van-grid-item v-if="canCreateInbound" @click="router.push('/inbound/create')">
          <div class="quick-icon" style="background: linear-gradient(135deg, #00c6fb 0%, #005bea 100%);">
            <van-icon name="add-o" size="24" color="#fff" />
          </div>
          <div class="quick-label">新建入库</div>
        </van-grid-item>

        <!-- 出库管理(领料) - 根据权限显示 -->
        <van-grid-item v-if="canAccessRequisition" @click="router.push('/requisition')">
          <div class="quick-icon" style="background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);">
            <van-icon name="send-gift-o" size="24" color="#fff" />
          </div>
          <div class="quick-label">领料管理</div>
        </van-grid-item>

        <!-- 新建领料 - 根据权限显示 -->
        <van-grid-item v-if="canCreateRequisition" @click="router.push('/requisition/create')">
          <div class="quick-icon" style="background: linear-gradient(135deg, #a8edea 0%, #fed6e3 100%);">
            <van-icon name="description" size="24" color="#333" />
          </div>
          <div class="quick-label">新建领料</div>
        </van-grid-item>

        <van-grid-item @click="handleRefresh">
          <div class="quick-icon" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);">
            <van-icon name="replay" size="24" color="#fff" />
          </div>
          <div class="quick-label">刷新数据</div>
        </van-grid-item>
      </van-grid>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'
import { useAuthStore } from '@/stores/auth'
import { useNotificationStore } from '@/stores/notification'
import { getAppointmentStats } from '@/api/appointment'
import { getCalendarStatistics } from '@/api/attendance'
import NotificationBadge from '@/components/NotificationBadge.vue'

const router = useRouter()
const authStore = useAuthStore()
const notificationStore = useNotificationStore()

// 显示名称
const displayName = computed(() => {
  if (authStore.user) {
    return authStore.user.full_name || authStore.user.username || '用户'
  }
  return '游客'
})

// 问候语
const greeting = computed(() => {
  const hour = new Date().getHours()
  if (hour < 6) return '凌晨好'
  if (hour < 9) return '早上好'
  if (hour < 12) return '上午好'
  if (hour < 14) return '中午好'
  if (hour < 18) return '下午好'
  if (hour < 22) return '晚上好'
  return '夜深了'
})

// 当前日期
const currentDate = computed(() => {
  const now = new Date()
  const options = { year: 'numeric', month: 'long', day: 'numeric', weekday: 'long' }
  return now.toLocaleDateString('zh-CN', options)
})

// 统计数据
const stats = ref({
  todayAppointments: 0,
  pendingAppointments: 0,
  todayClockIns: 0,
  monthClockIns: 0
})

// 权限控制
const canAccessMaterialPlan = computed(() => {
  return authStore.hasPermission('material_plan_view')
})

const canCreatePlan = computed(() => {
  return authStore.hasPermission('material_plan_create')
})

const canAccessInbound = computed(() => {
  return authStore.hasPermission('inbound_view')
})

const canCreateInbound = computed(() => {
  return authStore.hasPermission('inbound_create')
})

const canAccessRequisition = computed(() => {
  return authStore.hasPermission('requisition_view')
})

const canCreateRequisition = computed(() => {
  return authStore.hasPermission('requisition_create')
})

// 页面加载
onMounted(async () => {
  // 刷新用户权限
  await authStore.initAuth()
  await loadStats()
  // 初始化通知
  notificationStore.fetchUnreadCount()
})

// 加载统计数据
async function loadStats() {
  try {
    // 加载预约统计
    const aptRes = await getAppointmentStats()
    if (aptRes.data) {
      stats.value.todayAppointments = aptRes.data.today_count || 0
      stats.value.pendingAppointments = aptRes.data.pending_count || 0
    }

    // 加载打卡统计
    const today = new Date().toISOString().split('T')[0]
    const clockRes = await getCalendarStatistics({
      start_date: today,
      end_date: today
    })
    if (clockRes.data && clockRes.data.length > 0) {
      stats.value.todayClockIns = clockRes.data[0].total_count || 0
    }

    // 本月打卡统计
    const now = new Date()
    const monthStart = new Date(now.getFullYear(), now.getMonth(), 1).toISOString().split('T')[0]
    const monthEnd = new Date(now.getFullYear(), now.getMonth() + 1, 0).toISOString().split('T')[0]
    const monthRes = await getCalendarStatistics({
      start_date: monthStart,
      end_date: monthEnd
    })
    if (monthRes.data) {
      stats.value.monthClockIns = monthRes.data.reduce((sum, item) => sum + (item.total_count || 0), 0)
    }
  } catch (error) {
    console.error('加载统计数据失败:', error)
  }
}

// 导航方法
function goToAppointments() {
  router.push('/appointments/calendar')
}

function goToClockIn() {
  router.push('/attendance/clock-in')
}

function goToRecords() {
  router.push('/attendance/records')
}

// 刷新数据
async function handleRefresh() {
  showToast({ type: 'loading', message: '刷新中...', forbidClick: true })
  await authStore.initAuth()
  await loadStats()
  showToast({ type: 'success', message: '刷新成功' })
}
</script>

<style scoped>
.dashboard {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding-bottom: 20px;
}

/* 欢迎卡片 */
.welcome-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 24px 16px;
  color: #fff;
}

.welcome-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
}

.welcome-text {
  flex: 1;
}

.greeting {
  font-size: 14px;
  opacity: 0.9;
  margin-bottom: 4px;
}

.username {
  font-size: 24px;
  font-weight: bold;
}

.welcome-icon {
  position: relative;
}

.welcome-date {
  font-size: 13px;
  opacity: 0.85;
}

/* 区块标题 */
.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 20px 16px 12px;
  font-size: 16px;
  font-weight: 600;
  color: #323233;
}

.section-title .van-icon {
  color: #1989fa;
}

/* 统计数据 */
.stats-section {
  margin: 0 16px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.stat-card {
  background: #fff;
  border-radius: 12px;
  padding: 16px;
  display: flex;
  align-items: center;
  gap: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  cursor: pointer;
  transition: transform 0.2s;
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
  flex-shrink: 0;
}

.stat-content {
  flex: 1;
  min-width: 0;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #323233;
  line-height: 1.2;
}

.stat-label {
  font-size: 13px;
  color: #969799;
  margin-top: 4px;
}

/* 快捷功能 */
.quick-actions {
  margin: 0 16px;
}

.quick-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.quick-label {
  font-size: 12px;
  color: #646566;
  text-align: center;
}
</style>
