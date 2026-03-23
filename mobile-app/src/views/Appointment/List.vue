<template>
  <div class="appointment-list">
    <van-nav-bar title="施工预约">
      <template #right>
        <van-icon name="plus" size="18" @click="router.push('/appointment/create')" />
      </template>
    </van-nav-bar>

    <!-- 筛选器 -->
    <van-dropdown-menu>
      <van-dropdown-item v-model="statusFilter" :options="statusOptions" @change="onFilterChange" />
      <van-dropdown-item v-model="urgentFilter" :options="urgentOptions" @change="onFilterChange" />
    </van-dropdown-menu>

    <ListContainer
      v-model:loading="loading"
      v-model:refreshing="refreshing"
      :finished="finished"
      :data="appointments"
      @load="onLoad"
      @refresh="onRefresh"
    >
      <ListItemCard
        v-for="apt in appointments"
        :key="apt.id"
        :item="apt"
        type="appointment"
        @click="goToDetail"
      />
    </ListContainer>
  </div>
</template>

<script setup>
import { ref, watch, onMounted, onActivated } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getAppointments } from '@/api/appointment'
import ListContainer from '@/components/common/ListContainer.vue'
import ListItemCard from '@/components/common/ListItemCard.vue'

const router = useRouter()
const route = useRoute()

const appointments = ref([])
const loading = ref(false)
const finished = ref(false)
const refreshing = ref(false)
const currentPage = ref(1)
const pageSize = 20

const statusFilter = ref(route.query.status || '')
const urgentFilter = ref('')

const statusOptions = [
  { text: '全部状态', value: '' },
  { text: '草稿', value: 'draft' },
  { text: '待审批', value: 'pending' },
  { text: '已排期', value: 'scheduled' },
  { text: '进行中', value: 'in_progress' },
  { text: '已完成', value: 'completed' },
  { text: '已取消', value: 'cancelled' },
  { text: '已拒绝', value: 'rejected' }
]

const urgentOptions = [
  { text: '全部', value: '' },
  { text: '加急', value: 'urgent' },
  { text: '普通', value: 'normal' }
]

onMounted(() => {
  loadData()
})

// 页面激活时刷新数据（从创建页返回时）
onActivated(() => {
  // 从其他页面返回时，总是刷新列表
  console.log('List activated, refreshing data...')
  onFilterChange()
})

// 监听路由变化，确保返回时刷新
watch(
  () => route.path,
  (newPath) => {
    if (newPath === '/appointments') {
      onFilterChange()
    }
  }
)

function onFilterChange() {
  appointments.value = []
  currentPage.value = 1
  finished.value = false
  loadData()
}

function onRefresh() {
  appointments.value = []
  currentPage.value = 1
  finished.value = false
  loadData()
}

async function onLoad() {
  if (refreshing.value) {
    appointments.value = []
    refreshing.value = false
  }
  await loadData()
}

async function loadData() {
  if (loading.value) return

  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize,
      status: statusFilter.value || undefined
    }

    if (urgentFilter.value === 'urgent') {
      params.is_urgent = true
    } else if (urgentFilter.value === 'normal') {
      params.is_urgent = false
    }

    console.log('Loading appointments with params:', params)
    const response = await getAppointments(params)
    console.log('API response:', response)
    console.log('Response.data:', response.data)
    console.log('Response.meta:', response.meta)

    // response 直接就是拦截器返回的 { success: true, data: [...], meta: {...} }
    // response.data 才是数据数组
    const items = response.data || []
    console.log('Parsed items:', items)

    if (currentPage.value === 1) {
      appointments.value = items
    } else {
      appointments.value.push(...items)
    }

    const meta = response.meta || {}
    const total = meta.total || 0
    finished.value = appointments.value.length >= total

    if (!finished.value) {
      currentPage.value++
    }
  } catch (error) {
    console.error('加载预约单失败:', error)
    finished.value = true
  } finally {
    loading.value = false
  }
}

function goToDetail(apt) {
  router.push(`/appointment/${apt.id}`)
}
</script>

<style scoped>
.appointment-list {
  min-height: 100vh;
  background-color: #f5f5f5;
}
</style>
