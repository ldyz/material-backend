<template>
  <div class="home-page">
    <!-- 头部欢迎 -->
    <div class="home-header">
      <div class="header-content">
        <div class="greeting">
          <p class="greeting-text">{{ greetingText }}</p>
          <h2 class="username">{{ displayName }}</h2>
        </div>
        <div class="header-icon">
          <van-icon name="user-circle-o" size="48" color="#666" />
        </div>
      </div>
    </div>

    <!-- 数据卡片 -->
    <div class="data-cards">
      <div
        v-if="canViewInbound"
        class="data-card"
        @click="goToInbound"
      >
        <div class="card-icon inbound">
          <van-icon name="logistics" size="24" />
        </div>
        <div class="card-content">
          <p class="card-label">待审批入库</p>
          <p class="card-value">{{ pendingInboundCount }}</p>
        </div>
      </div>

      <div
        v-if="canViewRequisition"
        class="data-card"
        @click="goToOutbound"
      >
        <div class="card-icon outbound">
          <van-icon name="send-gift-o" size="24" />
        </div>
        <div class="card-content">
          <p class="card-label">待审批出库</p>
          <p class="card-value">{{ pendingRequisitionCount }}</p>
        </div>
      </div>
    </div>

    <!-- 快捷功能 -->
    <van-cell-group inset title="快捷功能" class="quick-actions">
      <van-cell
        v-if="canViewInbound"
        title="入库管理"
        icon="logistics"
        is-link
        @click="goToInbound"
      />
      <van-cell
        v-if="canViewRequisition"
        title="出库管理"
        icon="send-gift-o"
        is-link
        @click="goToOutbound"
      />
      <van-cell
        v-if="canViewStock"
        title="库存查询"
        icon="shop-o"
        is-link
        @click="goToStock"
      />
      <van-cell
        v-if="canViewConstructionLog"
        title="施工日志"
        icon="notes-o"
        is-link
        @click="goToConstructionLog"
      />
      <van-cell
        v-if="canUseAI"
        title="AI 分析"
        icon="aim"
        is-link
        @click="goToAI"
      />
    </van-cell-group>

    <!-- 系统信息 -->
    <van-cell-group inset title="系统信息">
      <van-cell title="系统版本" :value="systemVersion" />
      <van-cell title="用户角色" :value="userRoles" />
    </van-cell-group>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'
import { usePermission } from '@/composables/usePermission'
import { useUserStore } from '@/stores/user'
import { getPendingInboundCount } from '@/api/inbound'
import { getPendingRequisitionCount } from '@/api/requisition'

const router = useRouter()
const { canViewInbound, canViewRequisition, canViewStock, canViewConstructionLog, canUseAI } = usePermission()
const userStore = useUserStore()

const pendingInboundCount = ref(0)
const pendingRequisitionCount = ref(0)
const systemVersion = ref('1.0.0')

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
  return userStore.userInfo?.username || '用户'
})

// 用户角色
const userRoles = computed(() => {
  const roles = userStore.roles
  if (roles.includes('admin')) return '管理员'
  if (roles.includes('project_manager')) return '项目经理'
  if (roles.includes('warehouse_manager')) return '仓库管理员'
  if (roles.includes('worker')) return '施工人员'
  return '普通用户'
})

// 加载待审批数量
async function loadPendingCounts() {
  try {
    if (canViewInbound.value) {
      // 适配统一响应格式
      const { data } = await getPendingInboundCount()
      pendingInboundCount.value = data?.count || 0
    }
    if (canViewRequisition.value) {
      // 适配统一响应格式
      const { data } = await getPendingRequisitionCount()
      pendingRequisitionCount.value = data?.count || 0
    }
  } catch (error) {
    console.error('加载待审批数量失败:', error)
  }
}

// 导航方法
function goToInbound() {
  router.push('/inbound')
}

function goToOutbound() {
  router.push('/outbound')
}

function goToStock() {
  router.push('/stock')
}

function goToConstructionLog() {
  router.push('/construction')
}

function goToAI() {
  router.push('/ai')
}

onMounted(() => {
  loadPendingCounts()
})
</script>

<style scoped>
.home-page {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding: 0 0 60px 0;
}

.home-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 30px 20px;
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
  font-size: 24px;
  font-weight: bold;
  margin: 0;
}

.data-cards {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  padding: 16px;
  margin-top: -20px;
}

.data-card {
  background: white;
  border-radius: 12px;
  padding: 16px;
  display: flex;
  align-items: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  cursor: pointer;
  transition: all 0.3s;
}

.data-card:active {
  transform: scale(0.98);
}

.card-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12px;
}

.card-icon.inbound {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  color: white;
}

.card-icon.outbound {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  color: white;
}

.card-content {
  flex: 1;
}

.card-label {
  font-size: 12px;
  color: #969799;
  margin: 0 0 4px 0;
}

.card-value {
  font-size: 20px;
  font-weight: bold;
  color: #323233;
  margin: 0;
}

.quick-actions {
  margin: 16px;
}
</style>
