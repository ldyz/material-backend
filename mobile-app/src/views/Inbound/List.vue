<template>
  <div class="inbound-list">
    <van-nav-bar title="入库单">
      <template #right>
        <van-icon name="plus" size="18" @click="router.push('/inbound/create')" />
      </template>
    </van-nav-bar>

    <ListContainer
      v-model:loading="loading"
      v-model:refreshing="refreshing"
      :finished="finished"
      :data="orders"
      @load="onLoad"
      @refresh="onRefresh"
    >
      <ListItemCard
        v-for="order in orders"
        :key="order.id"
        :item="order"
        type="inbound"
        @click="goToDetail"
      />
    </ListContainer>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getInboundOrders } from '@/api/inbound'
import ListContainer from '@/components/common/ListContainer.vue'
import ListItemCard from '@/components/common/ListItemCard.vue'

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

function goToDetail(order) {
  router.push(`/inbound/${order.id}`)
}
</script>

<style scoped>
.inbound-list {
  min-height: 100vh;
  background-color: #f7f8fa;
}
</style>
