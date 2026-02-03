<template>
  <div class="plans-page">
    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <!-- 筛选栏 -->
      <van-sticky>
        <div class="filter-bar">
          <van-dropdown-menu>
            <van-dropdown-item v-model="statusFilter" :options="statusOptions" @change="onFilterChange" />
            <van-dropdown-item v-model="projectFilter" :options="projectOptions" @change="onFilterChange" />
          </van-dropdown-menu>
        </div>
      </van-sticky>

      <!-- 计划列表 -->
      <van-list
        v-model:loading="loading"
        :finished="finished"
        finished-text="没有更多了"
        @load="onLoad"
      >
        <div
          v-for="plan in plans"
          :key="plan.id"
          class="plan-card"
          @click="goToDetail(plan.id)"
        >
          <div class="plan-header">
            <div class="plan-title">
              <van-tag :type="getStatusType(plan.status)" size="medium">
                {{ getStatusText(plan.status) }}
              </van-tag>
              <span class="plan-no">{{ plan.plan_no }}</span>
            </div>
            <van-icon name="arrow" color="#969799" />
          </div>

          <div class="plan-info">
            <div class="info-row">
              <span class="label">项目:</span>
              <span class="value">{{ plan.project_name }}</span>
            </div>
            <div class="info-row">
              <span class="label">计划类型:</span>
              <span class="value">{{ getPlanTypeText(plan.plan_type) }}</span>
            </div>
            <div class="info-row">
              <span class="label">预算:</span>
              <span class="value">¥{{ formatMoney(plan.total_budget) }}</span>
            </div>
            <div class="info-row">
              <span class="label">物资数量:</span>
              <span class="value">{{ plan.item_count || 0 }} 项</span>
            </div>
            <div class="info-row">
              <span class="label">创建时间:</span>
              <span class="value">{{ formatDate(plan.created_at) }}</span>
            </div>
          </div>

          <!-- 操作按钮 -->
          <div class="plan-actions" v-if="plan.status === 'pending' && canApprove">
            <van-button
              type="primary"
              size="small"
              @click.stop="goToApprove(plan.id)"
            >
              审批
            </van-button>
          </div>
        </div>
      </van-list>
    </van-pull-refresh>

    <!-- 悬浮新建按钮 -->
    <van-floating-bubble
      v-if="canCreate"
      axis="xy"
      icon="plus"
      @click="goToCreate"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showConfirmDialog } from 'vant'
import { getPlans, approvePlan } from '@/api/material_plan'
import { usePermission } from '@/composables/usePermission'

const router = useRouter()
const { canViewMaterialPlan, canApproveMaterialPlan, canCreateMaterialPlan } = usePermission()

const canApprove = computed(() => canApproveMaterialPlan.value)
const canCreate = computed(() => canCreateMaterialPlan.value)

const plans = ref([])
const loading = ref(false)
const finished = ref(false)
const refreshing = ref(false)
const currentPage = ref(1)
const pageSize = 20

const statusFilter = ref('')
const projectFilter = ref('')

const statusOptions = [
  { text: '全部状态', value: '' },
  { text: '草稿', value: 'draft' },
  { text: '待审批', value: 'pending' },
  { text: '进行中', value: 'active' },
  { text: '已完成', value: 'completed' },
  { text: '已取消', value: 'cancelled' },
]

const projectOptions = ref([
  { text: '全部项目', value: '' },
])

// 加载计划列表
async function onLoad() {
  if (refreshing.value) {
    plans.value = []
    currentPage.value = 1
    finished.value = false
  }

  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize,
      status: statusFilter.value || undefined,
      project_id: projectFilter.value || undefined,
    }

    const response = await getPlans(params)
    if (response.success) {
      const { data, total } = response.data

      if (refreshing.value) {
        plans.value = data || []
      } else {
        plans.value.push(...(data || []))
      }

      // 检查是否还有更多数据
      if (plans.value.length >= total) {
        finished.value = true
      } else {
        currentPage.value++
      }
    }
  } catch (error) {
    console.error('加载计划列表失败:', error)
    showToast('加载失败')
  } finally {
    loading.value = false
    refreshing.value = false
  }
}

// 下拉刷新
async function onRefresh() {
  finished.value = false
  await onLoad()
}

// 筛选变化
function onFilterChange() {
  plans.value = []
  currentPage.value = 1
  finished.value = false
  onLoad()
}

// 跳转详情
function goToDetail(id) {
  router.push(`/plans/${id}`)
}

// 跳转审批
function goToApprove(id) {
  router.push(`/plans/${id}/approve`)
}

// 跳转新建
function goToCreate() {
  router.push('/plans/create')
}

// 获取状态类型
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

// 获取状态文本
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

// 获取计划类型文本
function getPlanTypeText(type) {
  const texts = {
    monthly: '月度计划',
    weekly: '周计划',
    temporary: '临时计划',
    project: '项目计划',
  }
  return texts[type] || type
}

// 格式化金额
function formatMoney(amount) {
  if (!amount) return '0.00'
  return Number(amount).toLocaleString('zh-CN', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  })
}

// 格式化日期
function formatDate(dateStr) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return `${date.getMonth() + 1}/${date.getDate()}`
}
</script>

<style scoped>
.plans-page {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding-bottom: 60px;
}

.filter-bar {
  background: white;
}

.plan-card {
  background: white;
  margin: 12px 16px;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.plan-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.plan-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.plan-no {
  font-size: 16px;
  font-weight: 600;
  color: #323233;
}

.plan-info {
  margin-bottom: 12px;
}

.info-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  font-size: 14px;
}

.info-row:last-child {
  margin-bottom: 0;
}

.label {
  color: #969799;
  margin-right: 8px;
}

.value {
  color: #323233;
  font-weight: 500;
}

.plan-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid #ebedf0;
}
</style>
