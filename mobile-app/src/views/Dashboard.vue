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
          <van-icon name="chat-o" size="24" @click="goToNotifications" />
          <van-badge
            v-if="unreadCount > 0"
            :content="unreadCount > 99 ? '99+' : unreadCount"
            style="position: absolute; top: -5px; right: -5px;"
          />
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
        <div class="stat-card" @click="goToPendingApprovals">
          <div class="stat-icon" style="background-color: #fff3e0;">
            <van-icon name="clock-o" size="24" color="#ff976a" />
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.pendingCount || 0 }}</div>
            <div class="stat-label">待审批</div>
          </div>
        </div>

        <div class="stat-card" @click="goToMyAppointments">
          <div class="stat-icon" style="background-color: #e3f2fd;">
            <van-icon name="calendar-o" size="24" color="#1989fa" />
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.myAppointments || 0 }}</div>
            <div class="stat-label">我的预约</div>
          </div>
        </div>

        <div class="stat-card" v-if="canAccessInbound" @click="goToInbound">
          <div class="stat-icon" style="background-color: #f3e5f5;">
            <van-icon name="logistics" size="24" color="#9c27b0" />
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.inboundCount || 0 }}</div>
            <div class="stat-label">待入库</div>
          </div>
        </div>

        <div class="stat-card" @click="goToCompleted">
          <div class="stat-icon" style="background-color: #e8f5e9;">
            <van-icon name="checked" size="24" color="#07c160" />
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.completedCount || 0 }}</div>
            <div class="stat-label">已完成</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 快捷功能 -->
    <div class="quick-actions">
      <div class="section-title">
        <van-icon name="apps-o" />
        <span>快捷功能</span>
        <span style="font-size:12px;color:#999;margin-left:10px;">
          权限: {{ authStore.permissions.length }} |
          物资计划: {{ authStore.hasPermission('material_plan_view') ? '是' : '否' }}
        </span>
      </div>
      <van-grid :column-num="4" :border="false" :gutter="12">
        <van-grid-item
          v-if="canCreateAppointments"
          @click="goToCreateAppointment"
        >
          <div class="quick-icon" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);">
            <van-icon name="plus" size="24" color="#fff" />
          </div>
          <div class="quick-label">新建预约</div>
        </van-grid-item>

        <van-grid-item
          v-if="canAccessApprovals"
          @click="goToPendingApprovals"
        >
          <div class="quick-icon" style="background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);">
            <van-icon name="todo-list-o" size="24" color="#fff" />
          </div>
          <div class="quick-label">待审批</div>
        </van-grid-item>

        <van-grid-item @click="goToMyAppointments">
          <div class="quick-icon" style="background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);">
            <van-icon name="orders-o" size="24" color="#fff" />
          </div>
          <div class="quick-label">我的预约</div>
        </van-grid-item>

        <van-grid-item
          v-if="canAccessInbound"
          @click="goToInbound"
        >
          <div class="quick-icon" style="background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);">
            <van-icon name="logistics" size="24" color="#fff" />
          </div>
          <div class="quick-label">入库管理</div>
        </van-grid-item>

        <van-grid-item
          v-if="canAccessMaterialPlan"
          @click="goToPlans"
        >
          <div class="quick-icon" style="background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);">
            <van-icon name="todo-list-o" size="24" color="#fff" />
          </div>
          <div class="quick-label">物资计划</div>
        </van-grid-item>

        <van-grid-item @click="goToProfile">
          <div class="quick-icon" style="background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);">
            <van-icon name="user-o" size="24" color="#fff" />
          </div>
          <div class="quick-label">个人中心</div>
        </van-grid-item>

        <van-grid-item @click="handleRefresh">
          <div class="quick-icon" style="background: linear-gradient(135deg, #30cfd0 0%, #330867 100%);">
            <van-icon name="replay" size="24" color="#fff" />
          </div>
          <div class="quick-label">刷新数据</div>
        </van-grid-item>
      </van-grid>
    </div>

    <!-- 最近预约 -->
    <div class="recent-section" v-if="recentAppointments.length > 0">
      <div class="section-title">
        <van-icon name="clock-o" />
        <span>最近预约</span>
        <van-button type="primary" size="mini" plain @click="goToMyAppointments">
          查看全部
        </van-button>
      </div>
      <div class="recent-list">
        <div
          v-for="item in recentAppointments"
          :key="item.id"
          class="recent-item"
          @click="goToDetail(item.id)"
        >
          <div class="recent-header">
            <StatusTag :status="item.status" type="appointment" />
            <span class="recent-date">{{ formatAppointmentDate(item.appointment_date, item.time_slot) }}</span>
          </div>
          <div class="recent-title">{{ item.project_name || '未命名项目' }}</div>
          <div class="recent-info">
            <van-icon name="contact" size="14" />
            <span>{{ item.work_type || item.applicant_name || '-' }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 系统公告 -->
    <div class="notice-section" v-if="systemNotice">
      <div class="section-title">
        <van-icon name="volume-o" />
        <span>系统公告</span>
      </div>
      <div class="notice-card">
        <div class="notice-content">{{ systemNotice }}</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'
import { useAuthStore } from '@/stores/auth'
import { getPendingApprovalCount } from '@/api/appointment'
import StatusTag from '@/components/common/StatusTag.vue'
import { formatAppointmentDate } from '@/composables/useDateTime'
import { logger } from '@/utils/logger'

const router = useRouter()
const authStore = useAuthStore()

// 显示名称
const displayName = computed(() => {
  return authStore.user?.full_name || authStore.user?.username || '用户'
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
  pendingCount: 0,
  myAppointments: 0,
  inboundCount: 0,
  completedCount: 0
})

// 未读消息数
const unreadCount = ref(0)

// 最近预约
const recentAppointments = ref([])

// 系统公告
const systemNotice = ref('欢迎使用化建仪表移动端！如有问题请联系管理员。')

// 权限控制 - 基于权限系统
const canAccessApprovals = computed(() => {
  logger.log('[Dashboard] canAccessApprovals 计算, permissions:', authStore.permissions.length)
  return authStore.hasPermission('appointment_approve')
})

const canCreateAppointments = computed(() => {
  return authStore.hasPermission('appointment_create')
})

const canAccessInbound = computed(() => {
  return authStore.hasPermission('inbound_view')
})

const canAccessMaterialPlan = computed(() => {
  const perms = authStore.permissions
  const hasPerms = perms.includes('material_plan_view')
  logger.log('[Dashboard] canAccessMaterialPlan 计算:', {
    permissionsCount: perms.length,
    hasMaterialPlanView: hasPerms,
    allPerms: perms
  })
  return hasPerms
})

// 页面加载
const permissionsReady = ref(false)

// 监听权限变化
watch(
  () => authStore.permissions,
  (newPerms) => {
    logger.log('[Dashboard] 权限变化，数量:', newPerms.length)
    logger.log('[Dashboard] 是否包含 material_plan_view:', newPerms.includes('material_plan_view'))
  },
  { immediate: true, deep: true }
)

onMounted(async () => {
  // 确保用户权限是最新的
  await authStore.initAuth()
  permissionsReady.value = true
  logger.log('[Dashboard] initAuth 完成，当前权限:', authStore.permissions.length)
  logger.log('[Dashboard] hasPermission(material_plan_view):', authStore.hasPermission('material_plan_view'))
  await loadStats()
  await loadRecentAppointments()
})

// 加载统计数据
async function loadStats() {
  try {
    // 获取待审批数量
    const response = await getPendingApprovalCount()
    if (response.success) {
      stats.value.pendingCount = response.data || 0
    }

    // 获取我的预约数量（模拟数据，实际应从后端获取）
    stats.value.myAppointments = Math.floor(Math.random() * 10)

    // 入库数量（模拟数据）
    if (canAccessInbound.value) {
      stats.value.inboundCount = Math.floor(Math.random() * 5)
    }

    // 完成数量（模拟数据）
    stats.value.completedCount = Math.floor(Math.random() * 20)
  } catch (error) {
    logger.error('加载统计数据失败:', error)
  }
}

// 加载最近预约
async function loadRecentAppointments() {
  try {
    // 这里应该调用实际的API获取最近预约
    // 暂时使用模拟数据
    recentAppointments.value = [
      {
        id: 1,
        project_name: '测试项目A',
        work_type: '施工作业',
        appointment_date: new Date().toISOString(),
        time_slot: 'morning',
        status: 'pending'
      },
      {
        id: 2,
        project_name: '测试项目B',
        work_type: '材料运输',
        appointment_date: new Date().toISOString(),
        time_slot: 'afternoon',
        status: 'approved'
      }
    ]
  } catch (error) {
    logger.error('加载最近预约失败:', error)
  }
}

// 导航方法
function goToPendingApprovals() {
  router.push('/appointments/pending')
}

function goToMyAppointments() {
  router.push('/appointments')
}

function goToCreateAppointment() {
  router.push('/appointment/create')
}

function goToInbound() {
  router.push('/inbound')
}

function goToPlans() {
  router.push('/plans')
}

function goToCompleted() {
  router.push('/appointments?status=completed')
}

function goToProfile() {
  router.push('/profile')
}

function goToNotifications() {
  router.push('/notifications')
}

function goToDetail(id) {
  router.push(`/appointment/${id}`)
}

// 刷新数据
async function handleRefresh() {
  showToast({ type: 'loading', message: '刷新中...', forbidClick: true })
  await loadStats()
  await loadRecentAppointments()
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
  padding: 8px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 50%;
  cursor: pointer;
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

/* 最近预约 */
.recent-section {
  margin: 0 16px;
}

.recent-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.recent-item {
  background: #fff;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  cursor: pointer;
  transition: transform 0.2s;
}

.recent-item:active {
  transform: scale(0.98);
}

.recent-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.recent-date {
  font-size: 12px;
  color: #969799;
}

.recent-title {
  font-size: 15px;
  font-weight: 500;
  color: #323233;
  margin-bottom: 8px;
}

.recent-info {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #646566;
}

/* 系统公告 */
.notice-section {
  margin: 0 16px;
}

.notice-card {
  background: linear-gradient(135deg, #ffeaa7 0%, #fdcb6e 100%);
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.notice-content {
  font-size: 14px;
  color: #2d3436;
  line-height: 1.6;
}
</style>
