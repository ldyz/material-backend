<template>
  <div class="outbound-approve-page">
    <van-nav-bar
      title="领料审批"
      left-arrow
      @click-left="onClickLeft"
    />

    <div v-if="loading" class="loading-container">
      <van-loading type="spinner" size="24" />
    </div>

    <div v-else-if="requisition" class="approve-content">
      <!-- 基本信息 -->
      <van-cell-group inset title="基本信息">
        <van-cell title="领料单号" :value="requisition.requisition_number" />
        <van-cell title="项目名称" :value="requisition.project_name" />
        <van-cell title="申请人" :value="requisition.applicant_name" />
        <van-cell title="申请日期" :value="formatDate(requisition.requisition_date)" />
        <van-cell title="用途" :value="requisition.purpose || '-'" />
      </van-cell-group>

      <!-- 材料明细 -->
      <van-cell-group inset title="材料明细（可修改审批数量）">
        <div
          v-for="(item, index) in editableItems"
          :key="item.id"
          class="material-item"
        >
          <div class="material-header">
            <span class="material-index">{{ index + 1 }}</span>
            <span class="material-name">{{ item.material_name }}</span>
          </div>
          <div class="material-info">
            <div class="info-item full">
              <span class="label">规格:</span>
              <span class="value">{{ item.specification || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">单位:</span>
              <span class="value">{{ item.unit }}</span>
            </div>
            <div class="info-item">
              <span class="label">申请数量:</span>
              <span class="value">{{ item.requested_quantity }} {{ item.unit }}</span>
            </div>
            <div class="info-item full">
              <span class="label">审批数量:</span>
              <van-stepper
                v-model="item.approved_quantity"
                :min="0"
                :max="item.requested_quantity"
                :integer="true"
                input-width="80px"
              />
            </div>
          </div>
        </div>
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
      :confirm-button-disabled="dialogAction === 'reject' && !remark.trim()"
      @confirm="onConfirmRemark"
    >
      <van-field
        v-model="remark"
        type="textarea"
        :placeholder="dialogAction === 'reject' ? '请输入拒绝原因（必填）' : '请输入备注（可选）'"
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
import { showToast } from 'vant'
import { getRequisitionDetail, approveRequisition, rejectRequisition } from '@/api/requisition'
import { formatDate } from '@/utils/date'

const router = useRouter()
const route = useRoute()

const loading = ref(true)
const requisition = ref(null)
const editableItems = ref([])
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
    const { data } = await getRequisitionDetail(route.params.id)

    // 规范化数据格式以适配前端
    requisition.value = {
      ...data,
      requisition_number: data.requisition_no || data.requisition_number,
      applicant_name: data.applicant || data.applicant_name,
      requisition_date: data.created_at || data.requisition_date,
    }

    // 初始化可编辑的明细列表
    editableItems.value = data.items.map(item => ({
      ...item,
      material_name: item.name || item.material_name,
      specification: item.spec || item.specification,
      requested_quantity: item.quantity || item.requested_quantity,
      approved_quantity: item.approved_quantity || item.quantity || item.requested_quantity,
    }))
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
  // 检查是否所有材料都设置了审批数量
  const hasZeroQty = editableItems.value.some(item => item.approved_quantity === 0)
  if (hasZeroQty) {
    showToast('存在审批数量为0的材料，请确认')
    return
  }

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
  // 拒绝时必须填写备注，通过时备注可选
  if (dialogAction.value === 'reject' && !remark.value.trim()) {
    showToast('请输入备注')
    return
  }

  try {
    if (dialogAction.value === 'approve') {
      approving.value = true
      const items = editableItems.value.map(item => ({
        id: item.id,
        approved_quantity: item.approved_quantity,
      }))
      await approveRequisition(requisition.value.id, {
        items,
        notes: remark.value,
      })
      showToast({
        type: 'success',
        message: '审批通过',
      })
    } else {
      rejecting.value = true
      await rejectRequisition(requisition.value.id, { notes: remark.value })
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
.outbound-approve-page {
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
  align-items: center;
  justify-content: space-between;
  font-size: 14px;
}

.info-item.full {
  grid-column: 1 / -1;
}

.label {
  color: #969799;
  flex-shrink: 0;
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
