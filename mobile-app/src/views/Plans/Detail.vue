<template>
  <div class="plan-detail-page">
    <van-nav-bar
      title="计划详情"
      left-text="返回"
      @click-left="onClickLeft"
    >
      <template #right>
        <van-button
          v-if="canApprove && plan.status === 'pending'"
          type="primary"
          size="small"
          @click="goToApprove"
        >
          审批
        </van-button>
      </template>
    </van-nav-bar>

    <div v-if="loading" class="loading-container">
      <van-loading size="24" vertical>加载中...</van-loading>
    </div>

    <div v-else-if="plan">
      <!-- 状态卡片 -->
      <div class="status-card" :class="'status-' + plan.status">
        <div class="status-content">
          <van-tag :type="getStatusType(plan.status)" size="large">
            {{ getStatusText(plan.status) }}
          </van-tag>
          <h2 class="plan-no">{{ plan.plan_no }}</h2>
          <p class="plan-description">{{ plan.description || plan.remark || '无备注' }}</p>
        </div>
      </div>

      <!-- 计划信息 -->
      <van-cell-group inset title="计划信息" class="info-section">
        <van-cell title="项目名称" :value="plan.project_name" />
        <van-cell title="计划类型" :value="getPlanTypeText(plan.plan_type)" />
        <van-cell title="优先级" :value="getPriorityText(plan.priority)" />
        <van-cell title="计划日期" :value="formatDate(plan.plan_date)" />
        <van-cell
          title="预算金额"
          :value="'¥' + formatMoney(plan.total_budget)"
        />
        <van-cell title="创建人" :value="plan.creator_name" />
        <van-cell title="创建时间" :value="formatDateTime(plan.created_at)" />
      </van-cell-group>

      <!-- 物资明细 -->
      <van-cell-group inset title="物资明细" class="items-section">
        <div
          v-for="(item, index) in plan.items"
          :key="index"
          class="item-card"
        >
          <div class="item-header">
            <span class="item-index">{{ index + 1 }}</span>
            <span class="item-name">{{ item.material_name || item.name }}</span>
          </div>
          <div class="item-details">
            <div class="detail-row">
              <span class="detail-label">规格:</span>
              <span class="detail-value">{{ item.specification || '-' }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">单位:</span>
              <span class="detail-value">{{ item.unit }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">数量:</span>
              <span class="detail-value highlight">{{ item.quantity }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">单价:</span>
              <span class="detail-value">¥{{ formatMoney(item.unit_price) }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">小计:</span>
              <span class="detail-value highlight">¥{{ formatMoney(item.quantity * item.unit_price) }}</span>
            </div>
            <div v-if="item.remark" class="detail-row full-width">
              <span class="detail-label">备注:</span>
              <span class="detail-value">{{ item.remark }}</span>
            </div>
          </div>
        </div>
      </van-cell-group>

      <!-- 审批记录 -->
      <van-cell-group inset title="审批记录" class="approvals-section">
        <van-steps direction="vertical" :active="currentStep">
          <van-step v-for="(approval, index) in approvals" :key="index">
            <template #title>
              <div class="approval-title">
                {{ approval.approver_name }}
                <van-tag
                  :type="getApprovalStatusType(approval.status)"
                  size="small"
                >
                  {{ getApprovalStatusText(approval.status) }}
                </van-tag>
              </div>
            </template>
            <template #default>
              <div class="approval-content">
                <p v-if="approval.comment" class="approval-comment">{{ approval.comment }}</p>
                <p class="approval-time">{{ formatDateTime(approval.created_at) }}</p>
              </div>
            </template>
          </van-step>
        </van-steps>
        <van-empty v-if="approvals.length === 0" description="暂无审批记录" />
      </van-cell-group>
    </div>

    <van-empty v-else description="计划不存在" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showToast } from 'vant'
import { getPlanDetail } from '@/api/material_plan'
import { usePermission } from '@/composables/usePermission'

const router = useRouter()
const route = useRoute()
const { canApproveMaterialPlan } = usePermission()

const canApprove = computed(() => canApproveMaterialPlan.value)

const plan = ref(null)
const approvals = ref([])
const loading = ref(true)

const currentStep = computed(() => {
  // 计算当前审批步骤
  const approvedIndex = approvals.value.findIndex(a => a.status === 'approved')
  return approvedIndex >= 0 ? approvedIndex : -1
})

async function loadDetail() {
  loading.value = true
  try {
    const id = route.params.id
    const response = await getPlanDetail(id)
    if (response.success) {
      plan.value = response.data
      approvals.value = response.data.approvals || []
    } else {
      showToast('加载失败')
    }
  } catch (error) {
    console.error('加载计划详情失败:', error)
    showToast('加载失败')
  } finally {
    loading.value = false
  }
}

function onClickLeft() {
  router.back()
}

function goToApprove() {
  router.push(`/plans/${route.params.id}/approve`)
}

function getStatusType(status) {
  const types = {
    draft: 'default',
    pending: 'warning',
    active: 'primary',
    completed: 'success',
    cancelled: 'danger',
  }
  return types[status] || 'default'
}

function getStatusText(status) {
  const texts = {
    draft: '草稿',
    pending: '待审批',
    active: '进行中',
    completed: '已完成',
    cancelled: '已取消',
  }
  return texts[status] || status
}

function getPlanTypeText(type) {
  const texts = {
    monthly: '月度计划',
    weekly: '周计划',
    temporary: '临时计划',
    project: '项目计划',
  }
  return texts[type] || type
}

function getPriorityText(priority) {
  const texts = {
    low: '低',
    medium: '中',
    high: '高',
    urgent: '紧急',
  }
  return texts[priority] || priority
}

function getApprovalStatusType(status) {
  const types = {
    pending: 'warning',
    approved: 'success',
    rejected: 'danger',
  }
  return types[status] || 'default'
}

function getApprovalStatusText(status) {
  const texts = {
    pending: '待审批',
    approved: '已通过',
    rejected: '已拒绝',
  }
  return texts[status] || status
}

function formatMoney(amount) {
  if (!amount) return '0.00'
  return Number(amount).toLocaleString('zh-CN', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  })
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

function formatDateTime(dateStr) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return `${date.getMonth() + 1}月${date.getDate()}日 ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
}

onMounted(() => {
  loadDetail()
})
</script>

<style scoped>
.plan-detail-page {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.loading-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 50vh;
}

.status-card {
  padding: 24px 20px;
  margin: 16px;
  border-radius: 12px;
  text-align: center;
}

.status-card.status-pending {
  background: linear-gradient(135deg, #fff7e6 0%, #ffebd9 100%);
}

.status-card.status-active {
  background: linear-gradient(135deg, #e8f4ff 0%, #d0e9ff 100%);
}

.status-card.status-completed {
  background: linear-gradient(135deg, #f0fff4 0%, #d4f8e4 100%);
}

.status-card.status-cancelled {
  background: linear-gradient(135deg, #fff0f0 0%, #ffd6d6 100%);
}

.plan-no {
  font-size: 20px;
  font-weight: bold;
  color: #323233;
  margin: 12px 0 8px 0;
}

.plan-description {
  color: #646566;
  margin: 0;
}

.info-section,
.items-section,
.approvals-section {
  margin: 16px 0;
}

.item-card {
  padding: 12px;
  margin: 8px 16px;
  background: #f7f8fa;
  border-radius: 8px;
}

.item-header {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
}

.item-index {
  width: 24px;
  height: 24px;
  background: #1989fa;
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  margin-right: 8px;
}

.item-name {
  font-size: 15px;
  font-weight: 600;
  color: #323233;
}

.item-details {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.detail-row {
  display: flex;
  font-size: 13px;
  min-width: 45%;
}

.detail-row.full-width {
  min-width: 100%;
  margin-top: 8px;
  padding-top: 8px;
  border-top: 1px dashed #ebedf0;
}

.detail-label {
  color: #969799;
  margin-right: 4px;
}

.detail-value {
  color: #323233;
}

.detail-value.highlight {
  font-weight: 600;
  color: #1989fa;
}

.approval-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 15px;
}

.approval-content {
  padding-left: 4px;
}

.approval-comment {
  margin: 4px 0;
  font-size: 14px;
  color: #323233;
}

.approval-time {
  margin: 0;
  font-size: 12px;
  color: #969799;
}
</style>
