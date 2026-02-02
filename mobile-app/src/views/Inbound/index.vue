<template>
  <div class="inbound-page">
    <!-- 搜索和筛选 -->
    <van-sticky>
      <div class="search-bar">
        <van-search
          v-model="searchKeyword"
          placeholder="搜索入库单号、供应商"
          @search="onSearch"
        />
      </div>
      <van-dropdown-menu>
        <van-dropdown-item
          v-model="filterStatus"
          :options="statusOptions"
          @change="onFilterChange"
        />
        <van-dropdown-item
          v-model="filterProject"
          :options="projectOptions"
          @change="onFilterChange"
        />
      </van-dropdown-menu>
    </van-sticky>

    <!-- 新建按钮 -->
    <div class="create-section">
      <van-button
        type="primary"
        icon="plus"
        round
        block
        @click="goToCreate"
      >
        新建入库单
      </van-button>
    </div>

    <!-- 入库单列表 -->
    <van-pull-refresh
      v-model="refreshing"
      @refresh="onRefresh"
    >
      <van-list
        v-model:loading="loading"
        :finished="finished"
        finished-text="没有更多了"
        @load="onLoad"
      >
        <div
          v-for="item in list"
          :key="item.id"
          class="inbound-card"
          @click="goToDetail(item.id)"
        >
          <div class="card-header">
            <span class="order-number">{{ item.order_no }}</span>
            <van-tag :type="getStatusType(item.status)">
              {{ getStatusText(item.status) }}
            </van-tag>
          </div>
          <div class="card-body">
            <div class="info-row">
              <span class="label">项目:</span>
              <span class="value">{{ item.project_name }}</span>
            </div>
            <div class="info-row">
              <span class="label">供应商:</span>
              <span class="value">{{ item.supplier }}</span>
            </div>
            <div class="info-row">
              <span class="label">创建人:</span>
              <span class="value">{{ item.creator_name }}</span>
            </div>
            <div class="info-row">
              <span class="label">日期:</span>
              <span class="value">{{ formatDate(item.created_at) }}</span>
            </div>
          </div>
          <div
            v-if="canApproveInbound && item.status === 'pending'"
            class="card-footer"
          >
            <van-button
              type="primary"
              size="small"
              @click.stop="goToApprove(item.id)"
            >
              审批
            </van-button>
          </div>
        </div>

        <van-empty
          v-if="!loading && list.length === 0"
          description="暂无数据"
        />
      </van-list>
    </van-pull-refresh>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showToast } from 'vant'
import { usePermission } from '@/composables/usePermission'
import { getInboundOrders } from '@/api/inbound'
import { getProjects } from '@/api/project'
import { formatDate } from '@/utils/date'
import { INBOUND_STATUS, INBOUND_STATUS_TEXT } from '@/utils/constants'

const router = useRouter()
const route = useRoute()
const { canApproveInbound } = usePermission()

const list = ref([])
const loading = ref(false)
const finished = ref(false)
const refreshing = ref(false)
const searchKeyword = ref('')
const filterStatus = ref('')
const filterProject = ref('')
const currentPage = ref(1)
const pageSize = 20

// 状态选项
const statusOptions = [
  { text: '全部状态', value: '' },
  { text: '待审批', value: 'pending' },
  { text: '已审批', value: 'approved' },
  { text: '已完成', value: 'completed' },
  { text: '已拒绝', value: 'rejected' },
]

// 项目选项
const projectOptions = ref([
  { text: '全部项目', value: '' },
])

// 加载项目列表
async function loadProjects() {
  try {
    // 适配统一响应格式
    const { data } = await getProjects({ per_page: 1000 })
    const projects = data || []
    projectOptions.value = [
      { text: '全部项目', value: '' },
      ...projects.map(p => ({
        text: p.name,
        value: p.id.toString()
      }))
    ]
  } catch (error) {
    console.error('加载项目列表失败:', error)
  }
}

// 获取状态类型
function getStatusType(status) {
  const typeMap = {
    pending: 'warning',
    approved: 'primary',
    completed: 'success',
    rejected: 'danger',
  }
  return typeMap[status] || 'default'
}

// 获取状态文本
function getStatusText(status) {
  return INBOUND_STATUS_TEXT[status] || status
}

// 加载数据
async function loadData(isRefresh = false) {
  if (isRefresh) {
    currentPage.value = 1
    finished.value = false
  }

  if (finished.value && !isRefresh) return

  loading.value = true

  try {
    const params = {
      page: currentPage.value,
      per_page: pageSize,
      status: filterStatus.value || undefined,
      project_id: filterProject.value || undefined,
      search: searchKeyword.value || undefined,
    }

    // 适配统一响应格式
    const { data, pagination } = await getInboundOrders(params)
    const items = data || []

    if (isRefresh) {
      list.value = items
    } else {
      list.value.push(...items)
    }

    // 从分页信息中获取总数
    const total = pagination?.total || 0
    if (list.value.length >= total) {
      finished.value = true
    } else {
      currentPage.value++
    }
  } catch (error) {
    showToast({
      type: 'fail',
      message: '加载失败',
    })
  } finally {
    loading.value = false
    refreshing.value = false
  }
}

// 下拉刷新
function onRefresh() {
  refreshing.value = true
  loadData(true)
}

// 上拉加载
function onLoad() {
  loadData()
}

// 搜索
function onSearch() {
  loadData(true)
}

// 筛选改变
function onFilterChange() {
  loadData(true)
}

// 跳转详情
function goToDetail(id) {
  router.push(`/inbound/${id}`)
}

// 跳转审批
function goToApprove(id) {
  router.push(`/inbound/${id}/approve`)
}

// 跳转新建
function goToCreate() {
  router.push('/inbound/create')
}

onMounted(async () => {
  await loadProjects()
  loadData()
})

// 监听路由变化，当有 refresh 参数时刷新数据
watch(() => route.query.refresh, (newVal) => {
  if (newVal) {
    loadData(true)
  }
})
</script>

<style scoped>
.inbound-page {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding-bottom: 50px;
}

.search-bar {
  background: white;
  padding: 8px 0;
}

.inbound-card {
  margin: 12px 16px;
  background: white;
  border-radius: 8px;
  overflow: hidden;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid #ebedf0;
}

.order-number {
  font-size: 14px;
  font-weight: bold;
  color: #323233;
}

.card-body {
  padding: 12px 16px;
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
}

.value {
  color: #323233;
}

.value.amount {
  color: #ee0a24;
  font-weight: bold;
}

.card-footer {
  padding: 8px 16px;
  border-top: 1px solid #ebedf0;
  text-align: right;
}

.create-section {
  padding: 12px 16px;
  background: white;
  margin-bottom: 8px;
}
</style>
