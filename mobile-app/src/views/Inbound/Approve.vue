<template>
  <div class="inbound-approve-page">
    <van-nav-bar
      title="入库审批"
      left-arrow
      @click-left="onClickLeft"
    />

    <div v-if="loading" class="loading-container">
      <van-loading type="spinner" size="24" />
    </div>

    <div v-else-if="order" class="approve-content">
      <!-- 基本信息 -->
      <van-cell-group inset title="基本信息">
        <van-cell title="入库单号" :value="order.order_number" />
        <van-cell title="项目名称" :value="order.project_name" />
        <van-cell title="供应商" :value="order.supplier_name" />
        <van-cell title="入库日期" :value="formatDate(order.inbound_date)" />
        <van-cell title="总金额" :value="`¥${order.total_amount}`" />
      </van-cell-group>

      <!-- 材料明细 -->
      <van-cell-group inset title="材料明细">
        <div
          v-for="(item, index) in order.items"
          :key="item.id"
          class="material-item"
        >
          <div class="material-header">
            <span class="material-index">{{ index + 1 }}</span>
            <span class="material-name">{{ item.material_name }}</span>
          </div>
          <div class="material-info">
            <div class="info-item">
              <span class="label">规格:</span>
              <span class="value">{{ item.specification || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">数量:</span>
              <span class="value">{{ item.quantity }} {{ item.unit }}</span>
            </div>
            <div class="info-item">
              <span class="label">单价:</span>
              <span class="value">¥{{ item.unit_price }}</span>
            </div>
            <div class="info-item total">
              <span class="label">金额:</span>
              <span class="value">¥{{ item.total_price }}</span>
            </div>
          </div>
        </div>
      </van-cell-group>

      <!-- 备注信息 -->
      <van-cell-group
        v-if="order.notes"
        inset
        title="备注"
      >
        <van-cell :value="order.notes" />
      </van-cell-group>

      <!-- 审批操作 -->
      <div class="approve-actions">
        <van-button
          type="danger"
          size="large"
          :loading="rejecting"
          @click="onReject"
        >
          拒绝
        </van-button>
        <van-button
          type="primary"
          size="large"
          :loading="approving"
          @click="onApprove"
        >
          通过
        </van-button>
      </div>
    </div>

    <!-- 备注输入弹窗 -->
    <van-dialog
      v-model:show="showRemarkDialog"
      :title="dialogTitle"
      show-cancel-button
      @confirm="onConfirmRemark"
    >
      <van-field
        v-model="remark"
        type="textarea"
        placeholder="请输入备注"
        rows="3"
        maxlength="200"
        show-word-limit
      />
    </van-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showToast, showConfirmDialog } from 'vant'
import { getInboundDetail, approveInbound, rejectInbound } from '@/api/inbound'
import { formatDate } from '@/utils/date'

const router = useRouter()
const route = useRoute()

const loading = ref(true)
const order = ref(null)
const approving = ref(false)
const rejecting = ref(false)
const showRemarkDialog = ref(false)
const remark = ref('')
const dialogAction = ref('') // 'approve' or 'reject'

const dialogTitle = computed(() =>
  dialogAction.value === 'approve' ? '审批通过' : '审批拒绝'
)

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

// 审批通过
function onApprove() {
  dialogAction.value = 'approve'
  remark.value = ''
  showRemarkDialog.value = true
}

// 审批拒绝
function onReject() {
  dialogAction.value = 'reject'
  remark.value = ''
  showRemarkDialog.value = true
}

// 确认备注
async function onConfirmRemark() {
  if (!remark.value.trim()) {
    showToast('请输入备注')
    return
  }

  try {
    if (dialogAction.value === 'approve') {
      approving.value = true
      await approveInbound(order.value.id, { notes: remark.value })
      showToast({
        type: 'success',
        message: '审批通过',
      })
    } else {
      rejecting.value = true
      await rejectInbound(order.value.id, { notes: remark.value })
      showToast({
        type: 'success',
        message: '已拒绝',
      })
    }
    router.back()
  } catch (error) {
    showToast({
      type: 'fail',
      message: error.message || '操作失败',
    })
  } finally {
    approving.value = false
    rejecting.value = false
  }
}

// 返回
function onClickLeft() {
  router.back()
}

onMounted(() => {
  loadDetail()
})
</script>

<style scoped>
.inbound-approve-page {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.loading-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 200px;
}

.approve-content {
  padding: 16px 0 100px 0;
}

.material-item {
  margin: 8px 0;
  padding: 12px;
  background: white;
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

.info-item.total {
  grid-column: 1 / -1;
  padding-top: 8px;
  border-top: 1px dashed #ebedf0;
  font-weight: bold;
}

.label {
  color: #969799;
}

.value {
  color: #323233;
}

.approve-actions {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  padding: 12px 16px;
  background: white;
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.08);
  z-index: 100;
}
</style>
