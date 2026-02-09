<template>
  <div class="inbound-list">
    <van-nav-bar title="入库单">
      <template #right>
        <van-icon name="plus" size="18" @click="router.push('/inbound/create')" />
      </template>
    </van-nav-bar>

    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <van-list
        v-model:loading="loading"
        :finished="finished"
        finished-text="没有更多了"
        @load="onLoad"
      >
        <van-empty v-if="!orders.length && !loading" description="暂无数据" />

        <van-cell
          v-for="order in orders"
          :key="order.id"
          class="order-item"
          is-link
          @click="goToDetail(order.id)"
        >
          <template #title>
            <div class="order-title">
              <van-tag :type="getStatusType(order.status)">
                {{ getStatusText(order.status) }}
              </van-tag>
              <span class="order-number">{{ order.order_number }}</span>
            </div>
          </template>
          <template #label>
            <div>项目：{{ order.project_name || '-' }}</div>
            <div>日期：{{ formatDate(order.inbound_date) }}</div>
          </template>
        </van-cell>
      </van-list>
    </van-pull-refresh>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getInboundOrders } from '@/api/inbound'

const router = useRouter()
const route = useRoute()

const orders = ref([])
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
  orders.value = []
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
    completed: 'primary'
  }
  return map[status] || 'default'
}

function getStatusText(status) {
  const map = {
    draft: '草稿',
    pending: '待审批',
    approved: '已批准',
    rejected: '已拒绝',
    completed: '已完成'
  }
  return map[status] || status
}

function formatDate(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('zh-CN')
}

async function loadData() {
  if (refreshing.value) {
    orders.value = []
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

    const response = await getInboundOrders(params)
    const data = response.data || []

    if (refreshing.value) {
      orders.value = data
    } else {
      orders.value.push(...data)
    }

    const total = response.pagination?.total || 0
    finished.value = orders.value.length >= total
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
  router.push(`/inbound/${id}`)
}
</script>

<style scoped>
.inbound-list {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.order-item {
  margin: 12px 16px;
  border-radius: 8px;
  overflow: hidden;
}

.order-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.order-number {
  font-weight: bold;
}
</style>
