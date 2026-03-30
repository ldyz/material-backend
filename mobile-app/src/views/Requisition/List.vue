<template>
  <div class="req-list">
    <van-nav-bar title="出库单">
      <template #right>
        <van-icon name="plus" size="18" @click="router.push('/requisition/create')" />
      </template>
    </van-nav-bar>

    <ListContainer
      v-model:loading="loading"
      v-model:refreshing="refreshing"
      :finished="finished"
      :data="requisitions"
      @load="onLoad"
      @refresh="onRefresh"
    >
      <ListItemCard
        v-for="req in requisitions"
        :key="req.id"
        :item="req"
        type="requisition"
        @click="goToDetail"
      />
    </ListContainer>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getRequisitions } from '@/api/requisition'
import ListContainer from '@/components/common/ListContainer.vue'
import ListItemCard from '@/components/common/ListItemCard.vue'
import { logger } from '@/utils/logger'

const router = useRouter()
const route = useRoute()

const requisitions = ref([])
const loading = ref(false)
const finished = ref(false)
const refreshing = ref(false)
const currentPage = ref(1)
const pageSize = 20
const statusFilter = ref(route.query.status || '')
const fullyIssuedFilter = ref(route.query.fully_issued === 'false' ? false : null)
const isInitialLoad = ref(true)

// 监听路由变化，重新加载数据
watch(() => [route.query.status, route.query.fully_issued], () => {
  statusFilter.value = route.query.status || ''
  fullyIssuedFilter.value = route.query.fully_issued === 'false' ? false : null
  // 不是首次加载时才重置数据
  if (!isInitialLoad.value) {
    refreshData()
  }
})

function refreshData() {
  requisitions.value = []
  currentPage.value = 1
  finished.value = false
  loadData()
}

async function loadData() {
  if (refreshing.value) {
    requisitions.value = []
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

    if (fullyIssuedFilter.value !== null) {
      params.fully_issued = fullyIssuedFilter.value
    }

    const response = await getRequisitions(params)
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
    logger.error('加载失败:', error)
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

function goToDetail(req) {
  router.push(`/requisition/${req.id}`)
}
</script>

<style scoped>
.req-list {
  min-height: 100vh;
  background-color: #f7f8fa;
}
</style>
