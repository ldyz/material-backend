<template>
  <div class="outbound-issue-page">
    <van-nav-bar
      title="发料"
      left-arrow
      @click-left="onClickLeft"
    />

    <div v-if="loading" class="loading-container">
      <van-loading type="spinner" size="24" />
    </div>

    <div v-else-if="requisition" class="issue-content">
      <!-- 基本信息 -->
      <van-cell-group inset title="领料单信息">
        <van-cell title="领料单号" :value="requisition.requisition_no" />
        <van-cell title="项目名称" :value="requisition.project_name" />
        <van-cell title="申请人" :value="requisition.applicant" />
        <van-cell title="申请日期" :value="formatDate(requisition.created_at)" />
        <van-cell title="用途" :value="requisition.purpose || '-'" />
      </van-cell-group>

      <!-- 材料明细 -->
      <van-cell-group inset title="材料明细（可修改发料数量）">
        <div
          v-for="(item, index) in editableItems"
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
              <span class="label">单位:</span>
              <span class="value">{{ item.unit }}</span>
            </div>
            <div class="info-item">
              <span class="label">申请数量:</span>
              <span class="value">{{ item.approved_quantity || item.requested_quantity }} {{ item.unit }}</span>
            </div>
            <div class="info-item">
              <span class="label">库存数量:</span>
              <span class="value">{{ item.stock_quantity || 0 }} {{ item.unit }}</span>
            </div>
            <div class="info-item full">
              <span class="label">发料数量:</span>
              <van-stepper
                v-model="item.actual_quantity"
                :min="0"
                :max="item.approved_quantity || item.requested_quantity"
                :integer="true"
                input-width="80px"
              />
            </div>
          </div>
        </div>
      </van-cell-group>

      <!-- 备注 -->
      <van-cell-group inset title="备注">
        <van-field
          v-model="notes"
          type="textarea"
          placeholder="请输入发料备注"
          rows="3"
          maxlength="200"
          show-word-limit
        />
      </van-cell-group>

      <!-- 发料操作 -->
      <div class="issue-actions">
        <van-button
          type="primary"
          size="large"
          block
          :loading="issuing"
          @click="onIssue"
        >
          确认发料
        </van-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showToast, showDialog } from 'vant'
import { getRequisitionDetail, issueRequisition } from '@/api/requisition'
import { formatDate } from '@/utils/date'

const router = useRouter()
const route = useRoute()

const loading = ref(true)
const requisition = ref(null)
const editableItems = ref([])
const notes = ref('')
const issuing = ref(false)

// 加载详情
async function loadDetail() {
  try {
    // 适配统一响应格式
    const { data } = await getRequisitionDetail(route.params.id)
    requisition.value = data

    // 初始化可编辑的明细列表
    editableItems.value = (data.items || []).map(item => ({
      ...item,
      approved_quantity: item.approved_quantity || item.requested_quantity || item.quantity,
      actual_quantity: item.actual_quantity || item.approved_quantity || item.requested_quantity || item.quantity,
      stock_quantity: item.stock_quantity || 0,
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

// 发料
async function onIssue() {
  // 检查是否所有材料都设置了发料数量
  const hasZeroQty = editableItems.value.some(item => item.actual_quantity === 0)
  if (hasZeroQty) {
    showToast('存在发料数量为0的材料，请确认')
    return
  }

  await showDialog({
    title: '确认发料',
    message: `确认发料？共 ${editableItems.value.length} 项物资`,
    showCancelButton: true
  })

  issuing.value = true

  try {
    const items = editableItems.value.map(item => ({
      id: item.id,
      actual_quantity: item.actual_quantity,
    }))

    await issueRequisition(requisition.value.id, {
      items,
      notes: notes.value,
    })

    showToast({
      type: 'success',
      message: '发料成功',
    })

    // 延迟跳转
    setTimeout(() => {
      router.back()
    }, 1000)
  } catch (error) {
    showToast({
      type: 'fail',
      message: error.message || '发料失败',
    })
  } finally {
    issuing.value = false
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
.outbound-issue-page {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.loading-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 200px;
}

.issue-content {
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

.issue-actions {
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
