<template>
  <div class="stock-page">
    <!-- 顶部统计 -->
    <div class="stats-cards">
      <div class="stat-card" @click="showAlertList">
        <div class="stat-icon alert">
          <van-icon name="warning-o" size="20" />
        </div>
        <div class="stat-info">
          <p class="stat-label">库存预警</p>
          <p class="stat-value">{{ alertCount }}</p>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon total">
          <van-icon name="box-o" size="20" />
        </div>
        <div class="stat-info">
          <p class="stat-label">总库存数</p>
          <p class="stat-value">{{ totalCount }}</p>
        </div>
      </div>
    </div>

    <!-- 搜索和筛选 -->
    <van-sticky>
      <div class="search-bar">
        <van-search
          v-model="searchKeyword"
          placeholder="搜索物资名称/规格/材质"
          @search="onSearch"
        />
      </div>
      <van-dropdown-menu>
        <van-dropdown-item
          v-model="filterProject"
          :options="projectOptions"
          @change="onFilterChange"
        />
        <van-dropdown-item
          v-model="filterStatus"
          :options="statusOptions"
          @change="onFilterChange"
        />
      </van-dropdown-menu>
    </van-sticky>

    <!-- 库存列表 -->
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
          class="stock-card"
          @click="goToDetail(item.id)"
        >
          <div class="card-header">
            <h3 class="material-name">{{ item.material_name || item.code || item.material_code || '未知材料' }}</h3>
            <van-tag
              :type="getStockStatusType(item)"
            >
              {{ getStockStatusText(item) }}
            </van-tag>
          </div>
          <div class="card-body">
            <div class="info-row">
              <span class="label">规格:</span>
              <span class="value">{{ item.specification || '-' }}</span>
            </div>
            <div class="info-row">
              <span class="label">单位:</span>
              <span class="value">{{ item.unit }}</span>
            </div>
            <div class="info-row">
              <span class="label">当前库存:</span>
              <span
                class="value quantity"
                :class="{ warning: item.quantity <= item.min_stock }"
              >
                {{ item.quantity }} {{ item.unit }}
              </span>
            </div>
            <div class="info-row">
              <span class="label">最小库存:</span>
              <span class="value">{{ item.min_stock || 0 }} {{ item.unit }}</span>
            </div>
          </div>

          <!-- 进度条 -->
          <div class="stock-progress">
            <van-progress
              :percentage="getStockPercentage(item)"
              :color="getStockProgressColor(item)"
            />
          </div>
        </div>

        <van-empty
          v-if="!loading && list.length === 0"
          description="暂无数据"
        />
      </van-list>
    </van-pull-refresh>

    <!-- 预警列表弹窗 -->
    <van-popup
      v-model:show="showAlerts"
      position="right"
      :style="{ width: '100%', height: '100%' }"
    >
      <div class="alert-list-page">
        <van-nav-bar
          title="库存预警"
          left-arrow
          @click-left="showAlerts = false"
        />
        <div class="alert-list">
          <div
            v-for="item in alertList"
            :key="item.id"
            class="alert-item"
          >
            <van-icon name="warning-o" color="#ee0a24" />
            <div class="alert-info">
              <p class="material-name">{{ item.material_name || item.code || item.material_code || '未知材料' }}</p>
              <p class="alert-text">
                当前库存: {{ item.quantity }} {{ item.unit }}
                <span class="min-stock">（最小: {{ item.min_stock || 0 }} {{ item.unit }}）</span>
              </p>
            </div>
          </div>
          <van-empty
            v-if="alertList.length === 0"
            description="暂无预警"
          />
        </div>
      </div>
    </van-popup>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'
import { getStocks, getStockAlerts } from '@/api/stock'
import { getProjects } from '@/api/project'

const router = useRouter()

const list = ref([])
const alertList = ref([])
const loading = ref(false)
const finished = ref(false)
const refreshing = ref(false)
const searchKeyword = ref('')
const filterProject = ref(0)  // 0表示全部项目
const filterStatus = ref('')  // 空字符串表示全部状态
const currentPage = ref(1)
const pageSize = 20
const showAlerts = ref(false)

// 项目选项
const projectOptions = ref([
  { text: '全部项目', value: 0 },
])

// 状态选项
const statusOptions = ref([
  { text: '全部状态', value: '' },
  { text: '库存正常', value: 'normal' },
  { text: '库存偏低', value: 'low' },
  { text: '库存不足', value: 'shortage' },
])

// 统计数据
const alertCount = computed(() => alertList.value.length)
const totalCount = computed(() => list.value.reduce((sum, item) => sum + item.quantity, 0))

// 获取库存状态类型
function getStockStatusType(item) {
  if (item.quantity === 0) return 'danger'
  if (item.quantity <= item.min_stock) return 'warning'
  return 'success'
}

// 获取库存状态文本
function getStockStatusText(item) {
  if (item.quantity === 0) return '缺货'
  if (item.quantity <= item.min_stock) return '预警'
  return '正常'
}

// 获取库存百分比
function getStockPercentage(item) {
  if (!item.max_stock) return 50
  const percentage = (item.quantity / item.max_stock) * 100
  return Math.min(100, Math.max(0, percentage))
}

// 获取进度条颜色
function getStockProgressColor(item) {
  if (item.quantity === 0) return '#ee0a24'
  if (item.quantity <= item.min_stock) return '#ff976a'
  return '#07c160'
}

// 加载库存数据
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
      search: searchKeyword.value || undefined,
    }

    // 添加项目筛选（0表示全部项目，不传参数）
    if (filterProject.value !== 0) {
      params.project_id = filterProject.value
    }

    // 添加状态筛选（空字符串表示全部状态，不传参数）
    if (filterStatus.value !== '') {
      params.status = filterStatus.value
    }

    // 适配统一响应格式
    const { data, pagination } = await getStocks(params)
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

// 加载项目列表
async function loadProjects() {
  try {
    // 适配统一响应格式
    const { data } = await getProjects()
    const projects = data || []
    projectOptions.value = [
      { text: '全部项目', value: 0 },
      ...projects.map(p => ({
        text: p.name,
        value: p.id
      }))
    ]
  } catch (error) {
    console.error('加载项目失败:', error)
  }
}

// 加载预警数据
async function loadAlerts() {
  try {
    // 适配统一响应格式
    const { data } = await getStockAlerts()
    alertList.value = data || []
  } catch (error) {
    console.error('加载预警失败:', error)
  }
}

// 下拉刷新
function onRefresh() {
  loadData(true)
  loadAlerts()
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

// 显示预警列表
function showAlertList() {
  showAlerts.value = true
}

// 跳转详情
function goToDetail(id) {
  router.push(`/stock/${id}`)
}

onMounted(() => {
  loadProjects()
  loadData()
  loadAlerts()
})
</script>

<style scoped>
.stock-page {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding-bottom: 50px;
}

.stats-cards {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  padding: 16px;
  background: white;
}

.stat-card {
  display: flex;
  align-items: center;
  padding: 12px;
  background: #f7f8fa;
  border-radius: 8px;
  cursor: pointer;
}

.stat-icon {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12px;
  color: white;
}

.stat-icon.alert {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.stat-icon.total {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-info {
  flex: 1;
}

.stat-label {
  font-size: 12px;
  color: #969799;
  margin: 0 0 4px 0;
}

.stat-value {
  font-size: 20px;
  font-weight: bold;
  color: #323233;
  margin: 0;
}

.search-bar {
  background: white;
  padding: 8px 0;
}

.stock-card {
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

.material-name {
  font-size: 16px;
  font-weight: bold;
  color: #323233;
  margin: 0;
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

.value.quantity {
  font-weight: bold;
}

.value.quantity.warning {
  color: #ee0a24;
}

.stock-progress {
  padding: 0 16px 12px 16px;
}

.alert-list-page {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.alert-list {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.alert-item {
  display: flex;
  align-items: flex-start;
  padding: 12px;
  background: white;
  border-radius: 8px;
  margin-bottom: 12px;
}

.alert-item .van-icon {
  margin-right: 12px;
  margin-top: 2px;
}

.alert-info {
  flex: 1;
}

.alert-info .material-name {
  font-size: 15px;
  font-weight: bold;
  color: #323233;
  margin: 0 0 4px 0;
}

.alert-info .alert-text {
  font-size: 13px;
  color: #646566;
  margin: 0;
}

.min-stock {
  color: #969799;
  margin-left: 8px;
}
</style>
