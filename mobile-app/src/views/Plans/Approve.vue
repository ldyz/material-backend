<template>
  <div class="plan-approve-page">
    <van-nav-bar
      title="审批计划"
      left-text="返回"
      @click-left="onClickLeft"
    />

    <div v-if="loading" class="loading-container">
      <van-loading size="24" vertical>加载中...</van-loading>
    </div>

    <div v-else-if="plan">
      <!-- 计划基本信息 -->
      <div class="plan-summary">
        <van-tag :type="getStatusType(plan.status)" size="large">
          {{ getStatusText(plan.status) }}
        </van-tag>
        <h2 class="plan-no">{{ plan.plan_no }}</h2>
        <p class="project-name">{{ plan.project_name }}</p>
        <div class="plan-stats">
          <div class="stat-item">
            <span class="stat-label">物资数量</span>
            <span class="stat-value">{{ plan.item_count || 0 }} 项</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">预算金额</span>
            <span class="stat-value">¥{{ formatMoney(plan.total_budget) }}</span>
          </div>
        </div>
      </div>

      <!-- 物资明细预览 -->
      <van-cell-group inset title="物资明细">
        <van-collapse v-model="activeNames" accordion>
          <van-collapse-item
            v-for="(item, index) in plan.items"
            :key="index"
            :title="`${index + 1}. ${item.material_name || item.name}`"
          >
            <template #default>
              <div class="collapse-content">
                <div class="content-row">
                  <span class="label">规格:</span>
                  <span class="value">{{ item.specification || '-' }}</span>
                </div>
                <div class="content-row">
                  <span class="label">单位:</span>
                  <span class="value">{{ item.unit }}</span>
                </div>
                <div class="content-row">
                  <span class="label">数量:</span>
                  <span class="value">{{ item.quantity }}</span>
                </div>
                <div class="content-row">
                  <span class="label">单价:</span>
                  <span class="value">¥{{ formatMoney(item.unit_price) }}</span>
                </div>
                <div class="content-row highlight">
                  <span class="label">小计:</span>
                  <span class="value">¥{{ formatMoney(item.quantity * item.unit_price) }}</span>
                </div>
                <div v-if="item.remark" class="content-row">
                  <span class="label">备注:</span>
                  <span class="value">{{ item.remark }}</span>
                </div>
              </div>
            </template>
          </van-collapse-item>
        </van-collapse>
      </van-cell-group>

      <!-- 审批表单 -->
      <van-cell-group inset title="审批操作" class="approve-section">
        <van-field
          v-model="remark"
          type="textarea"
          label="审批意见"
          placeholder="请输入审批意见"
          rows="3"
          autosize
        />
      </van-cell-group>

      <!-- 操作按钮 -->
      <div class="action-buttons">
        <van-button
          type="success"
          size="large"
          block
          :loading="submitting"
          @click="handleApprove(true)"
        >
          通过
        </van-button>
        <van-button
          type="danger"
          size="large"
          block
          :loading="submitting"
          @click="handleReject"
        >
          拒绝
        </van-button>
      </div>
    </div>

    <van-empty v-else description="计划不存在" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showToast, showConfirmDialog } from 'vant'
import { getPlanDetail, approvePlan, rejectPlan } from '@/api/material_plan'

const router = useRouter()
const route = useRoute()

const plan = ref(null)
const loading = ref(true)
const submitting = ref(false)
const remark = ref('')
const activeNames = ref(['0'])

async function loadDetail() {
  loading.value = true
  try {
    const id = route.params.id
    const response = await getPlanDetail(id)
    if (response.success) {
      plan.value = response.data
      if (plan.value.items && plan.value.items.length > 0) {
        activeNames.value = ['0']
      }
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

async function handleApprove(approved) {
  if (!remark.value && approved) {
    showToast('请输入审批意见')
    return
  }

  submitting.value = true
  try {
    const id = route.params.id
    let response
    if (approved) {
      response = await approvePlan(id, { remark: remark.value })
    } else {
      response = await rejectPlan(id, { remark: remark.value })
    }

    if (response.success) {
      showToast(approved ? '已通过' : '已拒绝')
      setTimeout(() => {
        router.back()
      }, 1000)
    } else {
      showToast(response.message || '操作失败')
    }
  } catch (error) {
    console.error('审批失败:', error)
    showToast('操作失败')
  } finally {
    submitting.value = false
  }
}

function handleReject() {
  showConfirmDialog({
    title: '确认拒绝',
    message: '拒绝后计划将被退回，是否继续？',
  })
    .then(() => {
      handleApprove(false)
    })
    .catch(() => {
      // 用户取消
    })
}

function onClickLeft() {
  router.back()
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

function formatMoney(amount) {
  if (!amount) return '0.00'
  return Number(amount).toLocaleString('zh-CN', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  })
}

onMounted(() => {
  loadDetail()
})
</script>

<style scoped>
.plan-approve-page {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding-bottom: 80px;
}

.loading-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 50vh;
}

.plan-summary {
  background: white;
  margin: 16px;
  padding: 20px;
  border-radius: 12px;
  text-align: center;
}

.plan-no {
  font-size: 20px;
  font-weight: bold;
  color: #323233;
  margin: 12px 0 8px 0;
}

.project-name {
  color: #646566;
  margin: 0 0 16px 0;
}

.plan-stats {
  display: flex;
  justify-content: center;
  gap: 32px;
  padding-top: 16px;
  border-top: 1px solid #ebedf0;
}

.stat-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-label {
  font-size: 12px;
  color: #969799;
}

.stat-value {
  font-size: 16px;
  font-weight: 600;
  color: #323233;
}

.collapse-content {
  padding: 8px 0;
}

.content-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  font-size: 14px;
}

.content-row.highlight {
  padding-top: 8px;
  border-top: 1px dashed #ebedf0;
  margin-top: 8px;
}

.content-row .label {
  color: #969799;
}

.content-row .value {
  color: #323233;
  font-weight: 500;
}

.content-row.highlight .value {
  color: #1989fa;
  font-weight: 600;
}

.approve-section {
  margin: 16px;
}

.action-buttons {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 12px 16px;
  background: white;
  box-shadow: 0 -2px 12px rgba(0, 0, 0, 0.08);
  display: flex;
  gap: 12px;
  padding-bottom: calc(12px + env(safe-area-inset-bottom));
}
</style>
