<template>
  <div class="plan-list">
    <van-nav-bar title="物资计划" />

    <ListContainer
      v-model:loading="loading"
      v-model:refreshing="refreshing"
      :finished="finished"
      :data="plans"
      @load="onLoad"
      @refresh="onRefresh"
    >
      <ListItemCard
        v-for="plan in plans"
        :key="plan.id"
        :item="plan"
        type="plan"
        @click="goToDetail"
      />
    </ListContainer>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getPlans } from '@/api/material_plan'
import ListContainer from '@/components/common/ListContainer.vue'
import ListItemCard from '@/components/common/ListItemCard.vue'
import { logger } from '@/utils/logger'

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

function goToDetail(plan) {
  router.push(`/plans/${plan.id}`)
}
</script>

<style scoped>
.plan-list {
  min-height: 100vh;
  background-color: #f7f8fa;
}
</style>
