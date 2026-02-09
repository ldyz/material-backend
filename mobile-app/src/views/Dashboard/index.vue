<template>
  <div class="dashboard">
    <van-nav-bar title="首页">
      <template #right>
        <NotificationBadge />
      </template>
    </van-nav-bar>

    <van-cell-group inset class="user-card">
      <van-cell center>
        <template #title>
          <span class="user-name">{{ authStore.token ? '已登录' : '游客' }}</span>
        </template>
        <template #icon>
          <van-icon name="user-o" size="40" color="#1989fa" />
        </template>
      </van-cell>
    </van-cell-group>

    <van-grid column-num="2" :border="false" class="stats-grid">
      <van-grid-item
        icon="orders-o"
        text="待审批计划"
        :badge="stats.pendingPlans || 0"
        @click="goToPendingPlans"
      />
      <van-grid-item
        icon="logistics"
        text="待审批入库"
        :badge="stats.pendingInbound || 0"
        @click="goToPendingInbound"
      />
      <van-grid-item
        icon="send-gift-o"
        text="待审批出库"
        :badge="stats.pendingRequisition || 0"
        @click="goToPendingRequisition"
      />
      <van-grid-item
        icon="bag-o"
        text="待发放"
        :badge="stats.pendingIssue || 0"
        @click="goToPendingIssue"
      />
    </van-grid>

    <van-cell-group inset title="快速入口">
      <van-cell title="物资计划" icon="orders-o" is-link @click="router.push('/plans')" />
      <van-cell title="入库单" icon="logistics" is-link @click="router.push('/inbound')" />
      <van-cell title="新建入库" icon="add-o" is-link @click="router.push('/inbound/create')" />
      <van-cell title="出库单" icon="send-gift-o" is-link @click="router.push('/requisition')" />
      <van-cell title="新建出库" icon="add-o" is-link @click="router.push('/requisition/create')" />
    </van-cell-group>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useNotificationStore } from '@/stores/notification'
import { getPlans } from '@/api/material_plan'
import { getInboundOrders } from '@/api/inbound'
import { getRequisitions } from '@/api/requisition'
import NotificationBadge from '@/components/NotificationBadge.vue'

const router = useRouter()
const authStore = useAuthStore()
const notificationStore = useNotificationStore()

const stats = ref({
  pendingPlans: 0,
  pendingInbound: 0,
  pendingRequisition: 0,
  pendingIssue: 0
})

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

onMounted(() => {
  loadStats()
  // 初始化通知
  notificationStore.fetchUnreadCount()
})
</script>

<style scoped>
.dashboard {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding: 16px 0;
}

.user-card {
  margin-bottom: 16px;
}

.user-name {
  font-size: 18px;
  font-weight: bold;
  margin-left: 12px;
}

.stats-grid {
  margin-bottom: 16px;
  background: #fff;
  padding: 16px 0;
}

.stats-grid :deep(.van-grid-item) {
  cursor: pointer;
  transition: background-color 0.2s;
}

.stats-grid :deep(.van-grid-item:active) {
  background-color: #f2f3f5;
}
</style>
