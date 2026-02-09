<template>
  <div class="plan-list">
    <van-nav-bar title="物资计划" />

    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <van-list
        v-model:loading="loading"
        :finished="finished"
        finished-text="没有更多了"
        @load="onLoad"
      >
        <van-empty v-if="!plans.length && !loading" description="暂无数据" />

        <van-cell
          v-for="plan in plans"
          :key="plan.id"
          class="plan-item"
          is-link
          @click="goToDetail(plan.id)"
        >
          <template #title>
            <div class="plan-title">
              <van-tag :type="getStatusType(plan.status)">
                {{ getStatusText(plan.status) }}
              </van-tag>
              <span class="plan-number">{{ plan.plan_no }}</span>
            </div>
          </template>
          <template #label>
            <div>项目：{{ plan.project_name || '-' }}</div>
            <div>日期：{{ formatDate(plan.planned_start_date) }}</div>
          </template>
        </van-cell>
      </van-list>
    </van-pull-refresh>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getPlans } from '@/api/material_plan'

const router = useRouter()
const route = useRoute()

const plans = ref([])
const loading = ref(false)
const finished = ref(false)
const refreshing = ref(false)
const currentPage = ref(1)
const pageSize = 20
const statusFilter = ref(route.query.status || '')
const isInitialLoad = ref(true)

// 监听路由变化，重新加载数据
watch(() => route.query.status, () => {
  statusFilter.value = route.query.status || ''
  // 不是首次加载时才重置数据
  if (!isInitialLoad.value) {
    refreshData()
  }
})

function refreshData() {
  plans.value = []
  currentPage.value = 1
  finished.value = false
  loadData()
}

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
  if (refreshing.value) {
    plans.value = []
    currentPage.value = 1
    finished.value = false
  }

  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize
    }

    if (statusFilter.value) {
      params.status = statusFilter.value
    }

    const response = await getPlans(params)
    const data = response.data || []

    if (refreshing.value) {
      plans.value = data
    } else {
      plans.value.push(...data)
    }

    const total = response.pagination?.total || 0
    finished.value = plans.value.length >= total
    currentPage.value++
  } catch (error) {
    console.error('加载失败:', error)
  } finally {
    loading.value = false
    refreshing.value = false
    isInitialLoad.value = false
  }
}

function onLoad() {
  loadData()
}

function onRefresh() {
  refreshing.value = true
  loadData()
}

function goToDetail(id) {
  router.push(`/plans/${id}`)
}
</script>

<style scoped>
.plan-list {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.plan-item {
  margin: 12px 16px;
  border-radius: 8px;
  overflow: hidden;
}

.plan-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.plan-number {
  font-weight: bold;
}
</style>
