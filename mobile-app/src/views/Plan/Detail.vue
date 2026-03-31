<template>
  <div class="plan-detail">
    <van-nav-bar title="计划详情" left-arrow @click-left="router.back()" />

    <van-loading v-if="loading" type="spinner" vertical />

    <div v-else-if="plan">
      <DetailInfoGroup
        title="基本信息"
        :item="plan"
        status-type="plan"
        :fields="basicFields"
      />

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
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showToast, showConfirmDialog } from 'vant'
import { getPlanDetail, resubmitPlan } from '@/api/material_plan'
import DetailInfoGroup from '@/components/common/DetailInfoGroup.vue'
import { logger } from '@/utils/logger'

const router = useRouter()
const route = useRoute()

const loading = ref(true)
const resubmitting = ref(false)
const plan = ref(null)

const basicFields = computed(() => [
  { key: 'plan_no', label: '计划单号' },
  { key: 'project_name', label: '项目名称' },
  { key: 'planned_start_date', label: '计划日期', type: 'date' },
  { key: 'status', label: '状态', type: 'status' }
])

async function loadData() {
  loading.value = true
  try {
    const response = await getPlanDetail(route.params.id)
    plan.value = response.data
  } catch (error) {
    logger.error('加载失败:', error)
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
  } catch (cancel) {
    return // 用户取消
  }

  resubmitting.value = true
  try {
    await resubmitPlan(plan.value.id, {})
    showToast({ type: 'success', message: '已重新提交' })
    setTimeout(() => router.back(), 1500)
  } catch (error) {
    const errorMsg = error.error || error.message || '操作失败'
    showToast({ type: 'fail', message: errorMsg })
    logger.error('重新提交失败:', error)
  } finally {
    resubmitting.value = false
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
