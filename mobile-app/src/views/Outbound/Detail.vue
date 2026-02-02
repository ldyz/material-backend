<template>
  <div class="outbound-detail-page">
    <van-nav-bar
      title="领料单详情"
      left-arrow
      @click-left="onClickLeft"
    />

    <div v-if="loading" class="loading-container">
      <van-loading type="spinner" size="24" />
    </div>

    <div v-else-if="requisition" class="detail-content">
      <!-- 基本信息 -->
      <van-cell-group inset title="基本信息">
        <van-cell title="领料单号" :value="requisition.requisition_no" />
        <van-cell title="项目名称" :value="requisition.project_name" />
        <van-cell title="申请人" :value="requisition.applicant" />
        <van-cell title="申请日期" :value="formatDate(requisition.created_at)" />
        <van-cell title="用途" :value="requisition.purpose || '-'" />
        <van-cell title="状态">
          <template #value>
            <van-tag :type="getStatusType(requisition.status)">
              {{ getStatusText(requisition.status) }}
            </van-tag>
          </template>
        </van-cell>
        <van-cell
          v-if="requisition.remark"
          title="备注"
          :value="requisition.remark"
        />
      </van-cell-group>

      <!-- 材料明细 -->
      <van-cell-group inset title="材料明细">
        <div
          v-for="(item, index) in requisition.items"
          :key="item.id"
          class="material-item"
        >
          <div class="material-header">
            <span class="material-index">{{ index + 1 }}</span>
            <span class="material-name">{{ item.name || item.material_name }}</span>
          </div>
          <div class="material-info">
            <div class="info-item full">
              <span class="label">规格:</span>
              <span class="value">{{ item.spec || item.specification || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">材质:</span>
              <span class="value">{{ item.material || '-' }}</span>
            </div>
            <div class="info-item full">
              <span class="label">数量:</span>
              <span class="value quantity-value">{{ getDisplayQuantity(item) }}</span>
            </div>
          </div>
        </div>
        <van-empty
          v-if="!requisition.items || requisition.items.length === 0"
          description="暂无材料明细"
          image-size="80"
        />
      </van-cell-group>

      <!-- 审批流程 -->
      <van-cell-group inset title="审批流程" v-if="requisition.status !== 'pending'">
        <van-cell title="审批人" :value="requisition.approved_by || '-'" />
        <van-cell title="审批时间" :value="formatDate(requisition.approved_at)" />
        <div v-if="requisition.status === 'issued'">
          <van-cell title="发料人" :value="requisition.issued_by || '-'" />
          <van-cell title="发料时间" :value="formatDate(requisition.issued_at)" />
        </div>
      </van-cell-group>

      <!-- 操作按钮 -->
      <div class="action-buttons">
        <van-button
          v-if="canApproveRequisition && requisition.status === 'pending'"
          type="primary"
          block
          @click="goToApprove"
        >
          审批
        </van-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showToast } from 'vant'
import { usePermission } from '@/composables/usePermission'
import { getRequisitionDetail } from '@/api/requisition'
import { formatDate } from '@/utils/date'
import { REQUISITION_STATUS, REQUISITION_STATUS_TEXT } from '@/utils/constants'

const router = useRouter()
const route = useRoute()
const { canApproveRequisition } = usePermission()

const loading = ref(true)
const requisition = ref(null)

// 获取状态类型
function getStatusType(status) {
  const typeMap = {
    pending: 'warning',
    approved: 'primary',
    issued: 'success',
    rejected: 'danger',
  }
  return typeMap[status] || 'default'
}

// 获取状态文本
function getStatusText(status) {
  return REQUISITION_STATUS_TEXT[status] || status
}

// 获取显示数量
function getDisplayQuantity(item) {
  let quantity = item.requested_quantity || item.quantity || 0

  // 如果已审批或已发料，优先显示审批数量
  if (requisition.value && (requisition.value.status === 'approved' || requisition.value.status === 'issued')) {
    if (item.approved_quantity && item.approved_quantity > 0) {
      quantity = item.approved_quantity
    }
  }

  return `${quantity} ${item.unit}`
}

// 加载详情
async function loadDetail() {
  try {
    // 适配统一响应格式
    const { data } = await getRequisitionDetail(route.params.id)
    requisition.value = data
  } catch (error) {
    showToast({
      type: 'fail',
      message: '加载失败',
    })
    router.back()
  } finally {
    loading.value = false
  }
}

// 返回
function onClickLeft() {
  router.back()
}

// 去审批
function goToApprove() {
  router.push(`/outbound/${route.params.id}/approve`)
}

onMounted(() => {
  loadDetail()
})
</script>

<style scoped>
.outbound-detail-page {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.loading-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 200px;
}

.detail-content {
  padding: 16px 0 80px 0;
}

.material-item {
  margin: 8px 0;
  padding: 12px;
  background: #f7f8fa;
  border-radius: 8px;
}

.material-header {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
  padding-bottom: 8px;
  border-bottom: 1px solid #ebedf0;
}

.material-index {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  background: #1989fa;
  color: white;
  border-radius: 50%;
  font-size: 12px;
  margin-right: 8px;
}

.material-name {
  font-weight: bold;
  color: #323233;
}

.material-info {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  font-size: 14px;
}

.label {
  color: #969799;
}

.value {
  color: #323233;
}

.value.highlight {
  color: #07c160;
  font-weight: bold;
}

.action-buttons {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 12px 16px;
  background: white;
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.08);
  z-index: 100;
}
</style>
