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

    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <van-list
        v-model:loading="loading"
        :finished="finished"
        finished-text="没有更多了"
        @load="onLoad"
      >
        <van-empty v-if="!appointments.length && !loading" description="暂无数据" />

        <van-cell
          v-for="apt in appointments"
          :key="apt.id"
          class="appointment-item"
          is-link
          @click="goToDetail(apt.id)"
        >
          <template #title>
            <div class="appointment-title">
              <van-tag :type="getStatusColor(apt.status)">
                {{ getStatusLabel(apt.status) }}
              </van-tag>
              <span class="appointment-no">{{ apt.appointment_no }}</span>
              <van-tag v-if="apt.is_urgent" type="danger" size="small">加急</van-tag>
            </div>
          </template>
          <template #label>
            <div class="appointment-info">
              <div>作业时间：{{ formatDateTime(apt.work_date, apt.time_slot) }}</div>
              <div>作业地点：{{ apt.work_location }}</div>
              <div>作业内容：{{ apt.work_content }}</div>
              <div v-if="apt.assigned_worker_name">作业人员：{{ apt.assigned_worker_name }}</div>
            </div>
          </template>
        </van-cell>
      </van-list>
    </van-pull-refresh>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getAppointments, getTimeSlotLabel, getStatusLabel, getStatusColor } from '@/api/appointment'

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

    const { data } = await getAppointments(params)
    const items = data.data || []

    if (currentPage.value === 1) {
      appointments.value = items
    } else {
      appointments.value.push(...items)
    }

    const meta = data.meta || {}
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

function goToDetail(id) {
  router.push(`/appointment/${id}`)
}

function formatDateTime(dateStr, timeSlot) {
  const date = new Date(dateStr)
  const dateStr2 = date.toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })
  const slot = getTimeSlotLabel(timeSlot)
  return `${dateStr2} ${slot}`
}
</script>

<style scoped>
.appointment-list {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.appointment-item {
  margin-bottom: 8px;
}

.appointment-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.appointment-no {
  font-weight: bold;
  font-size: 15px;
}

.appointment-info {
  margin-top: 8px;
  font-size: 13px;
  color: #666;
}

.appointment-info > div {
  margin-bottom: 4px;
}
</style>
