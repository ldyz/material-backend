<template>
  <div class="plan-detail">
    <van-nav-bar title="计划详情" left-arrow @click-left="router.back()" />

    <van-loading v-if="loading" type="spinner" vertical />

    <div v-else-if="plan">
      <van-cell-group inset title="基本信息">
        <van-cell title="计划单号" :value="plan.plan_no" />
        <van-cell title="项目名称" :value="plan.project_name || '-'" />
        <van-cell title="计划日期" :value="formatDate(plan.planned_start_date)" />
        <van-cell title="状态">
          <template #value>
            <van-tag :type="getStatusType(plan.status)">
              {{ getStatusText(plan.status) }}
            </van-tag>
          </template>
        </van-cell>
      </van-cell-group>

      <van-cell-group inset title="物料明细">
        <van-empty v-if="!plan.items?.length" description="暂无物料" />
        <van-cell
          v-for="item in plan.items"
          :key="item.id"
          :title="item.material_name"
        >
          <template #label>
            <div>规格型号：{{ item.specification || '-' }}</div>
            <div>材质：{{ item.material || '-' }}</div>
            <div>数量：{{ item.planned_quantity || 0 }} {{ item.unit || '' }}</div>
          </template>
        </van-cell>
      </van-cell-group>

      <div v-if="plan.status === 'pending'" class="footer-actions">
        <van-button block type="primary" @click="router.push(`/plans/${plan.id}/approve`)">
          去审批
        </van-button>
      </div>

      <div v-if="plan.status === 'rejected'" class="footer-actions">
        <van-button
          block
          type="primary"
          :loading="resubmitting"
          @click="handleResubmit"
        >
          重新提交
        </van-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showToast, showConfirmDialog } from 'vant'
import { getPlanDetail, resubmitPlan } from '@/api/material_plan'

const router = useRouter()
const route = useRoute()

const loading = ref(true)
const resubmitting = ref(false)
const plan = ref(null)

function getStatusType(status) {
  const map = {
    draft: 'default',
    pending: 'warning',
    approved: 'success',
    rejected: 'danger',
    active: 'primary',
    completed: 'success',
    cancelled: 'info'
  }
  return map[status] || 'default'
}

function getStatusText(status) {
  const map = {
    draft: '草稿',
    pending: '待审批',
    approved: '已批准',
    rejected: '已拒绝',
    active: '进行中',
    completed: '已完成',
    cancelled: '已取消'
  }
  return map[status] || status
}

function formatDate(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('zh-CN')
}

async function loadData() {
  loading.value = true
  try {
    const response = await getPlanDetail(route.params.id)
    plan.value = response.data
  } catch (error) {
    console.error('加载失败:', error)
  } finally {
    loading.value = false
  }
}

async function handleResubmit() {
  try {
    await showConfirmDialog({
      title: '确认重新提交',
      message: '确认重新提交该计划进行审批？',
      teleport: '#app',
      confirmButtonColor: '#1989fa'
    })

    resubmitting.value = true
    try {
      await resubmitPlan(plan.value.id, {})
      showToast({ type: 'success', message: '已重新提交' })
      setTimeout(() => {
        loadData()
        router.back()
      }, 1500)
    } catch (error) {
      const errorMsg = error.error || error.message || '操作失败'
      showToast({ type: 'fail', message: errorMsg })
      throw error
    } finally {
      resubmitting.value = false
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('重新提交失败:', error)
    }
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.plan-detail {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding-bottom: 16px;
}

.van-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 200px;
}

.footer-actions {
  padding: 16px;
}
</style>
