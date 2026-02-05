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
              <van-tag :type="getStatusType(plan.workflow_status)">
                {{ getStatusText(plan.workflow_status) }}
              </van-tag>
              <span class="plan-number">{{ plan.plan_number }}</span>
            </div>
          </template>
          <template #label>
            <div>项目：{{ plan.project_name || '-' }}</div>
            <div>日期：{{ formatDate(plan.plan_date) }}</div>
          </template>
        </van-cell>
      </van-list>
    </van-pull-refresh>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { getPlans } from '@/api/material_plan'

const router = useRouter()

const plans = ref([])
const loading = ref(false)
const finished = ref(false)
const refreshing = ref(false)
const currentPage = ref(1)
const pageSize = 20

function getStatusType(status) {
  const map = {
    draft: 'default',
    pending_approval: 'warning',
    approved: 'success',
    rejected: 'danger'
  }
  return map[status] || 'default'
}

function getStatusText(status) {
  const map = {
    draft: '草稿',
    pending_approval: '待审批',
    approved: '已批准',
    rejected: '已拒绝'
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
    const response = await getPlans({
      page: currentPage.value,
      page_size: pageSize
    })
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
