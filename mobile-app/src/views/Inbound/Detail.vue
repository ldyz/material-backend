<template>
  <div class="inbound-detail-page">
    <van-nav-bar
      title="入库单详情"
      left-arrow
      @click-left="onClickLeft"
    />

    <div v-if="loading" class="loading-container">
      <van-loading type="spinner" size="24" />
    </div>

    <div v-else-if="order" class="detail-content">
      <!-- 基本信息 -->
      <van-cell-group inset title="基本信息">
        <van-cell title="入库单号" :value="order.order_no" />
        <van-cell title="项目名称" :value="order.project_name" />
        <van-cell title="供应商" :value="order.supplier" />
        <van-cell title="创建人" :value="order.creator_name" />
        <van-cell title="入库日期" :value="formatDate(order.created_at)" />
        <van-cell title="状态">
          <template #value>
            <van-tag :type="getStatusType(order.status)">
              {{ getStatusText(order.status) }}
            </van-tag>
          </template>
        </van-cell>
        <van-cell
          v-if="order.notes"
          title="备注"
          :value="order.notes"
        />
      </van-cell-group>

      <!-- 审批流程 -->
      <van-cell-group inset title="审批流程">
        <van-steps direction="vertical" :active="getStepActive(order.status)">
          <van-step>
            <template #active-icon>
              <van-icon name="passed" />
            </template>
            <h3>提交申请</h3>
            <p>{{ order.creator_name }} 创建于 {{ formatDate(order.created_at) }}</p>
          </van-step>
          <van-step>
            <template #active-icon>
              <van-icon name="passed" />
            </template>
            <h3>待审批</h3>
            <p v-if="order.status === 'pending'">当前状态：等待审批</p>
            <p v-else>已进入审批流程</p>
          </van-step>
          <van-step v-if="order.status === 'approved' || order.status === 'completed'">
            <template #active-icon>
              <van-icon name="passed" />
            </template>
            <h3>已审批</h3>
            <p>入库单已审批通过</p>
          </van-step>
          <van-step v-if="order.status === 'rejected'">
            <template #active-icon>
              <van-icon name="close" />
            </template>
            <h3>已拒绝</h3>
            <p>入库单已被拒绝</p>
          </van-step>
          <van-step v-if="order.status === 'completed'">
            <template #active-icon>
              <van-icon name="checked" />
            </template>
            <h3>已完成</h3>
            <p>入库单已完成</p>
          </van-step>
        </van-steps>
      </van-cell-group>

      <!-- 材料明细 -->
      <van-cell-group inset title="材料明细">
        <div
          v-for="(item, index) in order.items"
          :key="item.id"
          class="material-item"
        >
          <van-cell :title="`${index + 1}. ${item.material_name}`" />
          <van-cell title="规格" :value="item.spec || '-'" />
          <van-cell title="单位" :value="item.unit" />
          <van-cell title="数量" :value="`${item.quantity} ${item.unit}`" />
          <van-cell title="单价" :value="`¥${item.unit_price}`" />
          <van-cell
            v-if="item.remark"
            title="备注"
            :value="item.remark"
          />
        </div>
        <van-empty
          v-if="!order.items || order.items.length === 0"
          description="暂无材料明细"
          image-size="80"
        />
      </van-cell-group>

      <!-- 操作按钮 -->
      <div
        v-if="canApproveInbound && order.status === 'pending'"
        class="action-buttons"
      >
        <van-button
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
import { getInboundDetail } from '@/api/inbound'
import { formatDate } from '@/utils/date'
import { INBOUND_STATUS, INBOUND_STATUS_TEXT } from '@/utils/constants'

const router = useRouter()
const route = useRoute()
const { canApproveInbound } = usePermission()

const loading = ref(true)
const order = ref(null)

// 获取状态类型
function getStatusType(status) {
  const typeMap = {
    pending: 'warning',
    approved: 'primary',
    completed: 'success',
    rejected: 'danger',
  }
  return typeMap[status] || 'default'
}

// 获取状态文本
function getStatusText(status) {
  return INBOUND_STATUS_TEXT[status] || status
}

// 获取当前激活的步骤
function getStepActive(status) {
  const stepMap = {
    pending: 1,    // 提交申请 -> 待审批
    approved: 2,   // 提交申请 -> 待审批 -> 已审批
    completed: 3,  // 提交申请 -> 待审批 -> 已审批 -> 已完成
    rejected: 2,   // 提交申请 -> 待审批 -> 已拒绝
  }
  return stepMap[status] || 0
}

// 加载详情
async function loadDetail() {
  try {
    // 适配统一响应格式
    const { data } = await getInboundDetail(route.params.id)
    order.value = data
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
  router.push(`/inbound/${route.params.id}/approve`)
}

onMounted(() => {
  loadDetail()
})
</script>

<style scoped>
.inbound-detail-page {
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
