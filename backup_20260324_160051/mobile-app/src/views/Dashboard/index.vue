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
        <div class="stat-card" @click="goToPendingPlans">
          <div class="stat-icon" style="background-color: #fff3e0;">
            <van-icon name="clock-o" size="24" color="#ff976a" />
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.pendingPlans || 0 }}</div>
            <div class="stat-label">待审批计划</div>
          </div>
        </div>

        <div class="stat-card" @click="goToPendingInbound">
          <div class="stat-icon" style="background-color: #e3f2fd;">
            <van-icon name="logistics" size="24" color="#1989fa" />
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.pendingInbound || 0 }}</div>
            <div class="stat-label">待入库</div>
          </div>
        </div>

        <div class="stat-card" @click="goToPendingRequisition">
          <div class="stat-icon" style="background-color: #f3e5f5;">
            <van-icon name="send-gift-o" size="24" color="#9c27b0" />
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.pendingRequisition || 0 }}</div>
            <div class="stat-label">待出库</div>
          </div>
        </div>

        <div class="stat-card" @click="goToPendingIssue">
          <div class="stat-icon" style="background-color: #e8f5e9;">
            <van-icon name="bag-o" size="24" color="#07c160" />
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.pendingIssue || 0 }}</div>
            <div class="stat-label">待发放</div>
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
        <van-grid-item @click="router.push('/plans/create')">
          <div class="quick-icon" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);">
            <van-icon name="plus" size="24" color="#fff" />
          </div>
          <div class="quick-label">新建计划</div>
        </van-grid-item>

        <van-grid-item @click="router.push('/inbound/create')">
          <div class="quick-icon" style="background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);">
            <van-icon name="logistics" size="24" color="#fff" />
          </div>
          <div class="quick-label">新建入库</div>
        </van-grid-item>

        <van-grid-item @click="router.push('/requisition/create')">
          <div class="quick-icon" style="background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);">
            <van-icon name="send-gift-o" size="24" color="#fff" />
          </div>
          <div class="quick-label">新建出库</div>
        </van-grid-item>

        <van-grid-item @click="router.push('/plans')">
          <div class="quick-icon" style="background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);">
            <van-icon name="orders-o" size="24" color="#fff" />
          </div>
          <div class="quick-label">计划管理</div>
        </van-grid-item>

        <van-grid-item @click="router.push('/inbound')">
          <div class="quick-icon" style="background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);">
            <van-icon name="bag-o" size="24" color="#fff" />
          </div>
          <div class="quick-label">入库管理</div>
        </van-grid-item>

        <van-grid-item @click="router.push('/requisition')">
          <div class="quick-icon" style="background: linear-gradient(135deg, #30cfd0 0%, #330867 100%);">
            <van-icon name="todo-list-o" size="24" color="#fff" />
          </div>
          <div class="quick-label">出库管理</div>
        </van-grid-item>

        <van-grid-item @click="router.push('/appointments/calendar')">
          <div class="quick-icon" style="background: linear-gradient(135deg, #ff9800 0%, #f57c00 100%);">
            <van-icon name="calendar-o" size="24" color="#fff" />
          </div>
          <div class="quick-label">预约管理</div>
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
import { getPlans } from '@/api/material_plan'
import { getInboundOrders } from '@/api/inbound'
import { getRequisitions } from '@/api/requisition'
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
  pendingPlans: 0,
  pendingInbound: 0,
  pendingRequisition: 0,
  pendingIssue: 0
})

// 页面加载
onMounted(async () => {
  await loadStats()
  // 初始化通知
  notificationStore.fetchUnreadCount()
})

// 加载统计数据
async function loadStats() {
  try {
    const [planRes, inboundRes, reqRes] = await Promise.all([
      getPlans({ status: 'pending', page: 1, page_size: 1 }),
      getInboundOrders({ status: 'pending', page: 1, page_size: 1 }),
      getRequisitions({ status: 'pending', page: 1, page_size: 1 })
    ])

    stats.value.pendingPlans = planRes.pagination?.total || 0
    stats.value.pendingInbound = inboundRes.pagination?.total || 0
    stats.value.pendingRequisition = reqRes.pagination?.total || 0

    // 待发放数量 = 已批准但未发放完成的领料单
    const issueRes = await getRequisitions({
      status: 'approved',
      page: 1,
      page_size: 1,
      fully_issued: false
    })
    stats.value.pendingIssue = issueRes.pagination?.total || 0
  } catch (error) {
    console.error('加载统计数据失败:', error)
  }
}

// 导航方法
function goToPendingPlans() {
  router.push({ path: '/plans', query: { status: 'pending' } })
}

function goToPendingInbound() {
  router.push({ path: '/inbound', query: { status: 'pending' } })
}

function goToPendingRequisition() {
  router.push({ path: '/requisition', query: { status: 'pending' } })
}

function goToPendingIssue() {
  router.push({ path: '/requisition', query: { status: 'approved', fully_issued: 'false' } })
}

// 刷新数据
async function handleRefresh() {
  showToast({ type: 'loading', message: '刷新中...', forbidClick: true })
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
