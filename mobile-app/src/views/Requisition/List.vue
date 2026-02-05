<template>
  <div class="req-list">
    <van-nav-bar title="出库单">
      <template #right>
        <van-icon name="plus" size="18" @click="router.push('/requisition/create')" />
      </template>
    </van-nav-bar>

    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <van-list
        v-model:loading="loading"
        :finished="finished"
        finished-text="没有更多了"
        @load="onLoad"
      >
        <van-empty v-if="!requisitions.length && !loading" description="暂无数据" />

        <van-cell
          v-for="req in requisitions"
          :key="req.id"
          class="req-item"
          is-link
          @click="goToDetail(req.id)"
        >
          <template #title>
            <div class="req-title">
              <van-tag :type="getStatusType(req.workflow_status)">
                {{ getStatusText(req.workflow_status) }}
              </van-tag>
              <span class="req-number">{{ req.requisition_number }}</span>
            </div>
          </template>
          <template #label>
            <div>项目：{{ req.project_name || '-' }}</div>
            <div>日期：{{ formatDate(req.requisition_date) }}</div>
          </template>
        </van-cell>
      </van-list>
    </van-pull-refresh>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { getRequisitions } from '@/api/requisition'

const router = useRouter()

const requisitions = ref([])
const loading = ref(false)
const finished = ref(false)
const refreshing = ref(false)
const currentPage = ref(1)
const pageSize = 20

function getStatusType(status) {
  const map = {
    draft: 'default',
    pending_approval: 'warning',
    approved: 'primary',
    rejected: 'danger',
    issued: 'success'
  }
  return map[status] || 'default'
}

function getStatusText(status) {
  const map = {
    draft: '草稿',
    pending_approval: '待审批',
    approved: '已批准',
    rejected: '已拒绝',
    issued: '已发放'
  }
  return map[status] || status
}

function formatDate(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('zh-CN')
}

async function loadData() {
  if (refreshing.value) {
    requisitions.value = []
    currentPage.value = 1
    finished.value = false
  }

  loading.value = true
  try {
    const response = await getRequisitions({
      page: currentPage.value,
      page_size: pageSize
    })
    const data = response.data || []

    if (refreshing.value) {
      requisitions.value = data
    } else {
      requisitions.value.push(...data)
    }

    const total = response.pagination?.total || 0
    finished.value = requisitions.value.length >= total
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
  router.push(`/requisition/${id}`)
}
</script>

<style scoped>
.req-list {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.req-item {
  margin: 12px 16px;
  border-radius: 8px;
  overflow: hidden;
}

.req-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.req-number {
  font-weight: bold;
}
</style>
