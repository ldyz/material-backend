<template>
  <div class="construction-page">
    <!-- 顶部操作栏 -->
    <van-sticky>
      <div class="page-header">
        <div class="header-title">施工日志</div>
        <van-button
          v-if="canCreateConstructionLog"
          type="primary"
          size="small"
          icon="plus"
          @click="goToCreate"
        >
          新建日志
        </van-button>
      </div>
      <div class="search-bar">
        <van-search
          v-model="searchKeyword"
          placeholder="搜索日志内容、项目"
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
          v-model="filterDate"
          :options="dateOptions"
          @change="onFilterChange"
        />
      </van-dropdown-menu>
    </van-sticky>

    <!-- 日志列表 -->
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
          class="log-card"
          @click="goToDetail(item.id)"
        >
          <!-- 日期和天气 -->
          <div class="card-header">
            <div class="date-info">
              <van-icon name="calendar-o" />
              <span class="date-text">{{ formatDate(item.log_date, 'YYYY年MM月DD日') }}</span>
            </div>
            <div class="weather-info">
              <van-icon :name="getWeatherIcon(item.weather)" />
              <span class="weather-text">{{ item.weather_text || '晴' }}</span>
            </div>
          </div>

          <!-- 项目和标题 -->
          <div class="card-title">
            <h3>{{ item.project_name }}</h3>
            <p class="log-summary">{{ getSummary(item.content) }}</p>
          </div>

          <!-- 图片预览 -->
          <div v-if="item.images && item.images.length > 0" class="image-preview">
            <van-image
              v-for="(img, index) in item.images.slice(0, 3)"
              :key="index"
              :src="img"
              width="80"
              height="80"
              fit="cover"
              radius="4"
            />
            <div v-if="item.images.length > 3" class="more-images">
              +{{ item.images.length - 3 }}
            </div>
          </div>

          <!-- 底部信息 -->
          <div class="card-footer">
            <span class="author">
              <van-icon name="user-o" />
              {{ item.created_by_name }}
            </span>
            <span class="time">{{ formatRelativeTime(item.created_at) }}</span>
          </div>
        </div>

        <van-empty
          v-if="!loading && list.length === 0"
          description="暂无日志"
        />
      </van-list>
    </van-pull-refresh>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'
import { usePermission } from '@/composables/usePermission'
import { getConstructionLogs } from '@/api/construction'
import { formatDate, formatRelativeTime } from '@/utils/date'

const router = useRouter()
const { canCreateConstructionLog } = usePermission()

const list = ref([])
const loading = ref(false)
const finished = ref(false)
const refreshing = ref(false)
const searchKeyword = ref('')
const filterProject = ref('')
const filterDate = ref('')
const currentPage = ref(1)
const pageSize = 20

// 项目选项（模拟数据）
const projectOptions = ref([
  { text: '全部项目', value: '' },
  { text: '项目A', value: '1' },
  { text: '项目B', value: '2' },
])

// 日期选项
const dateOptions = ref([
  { text: '全部时间', value: '' },
  { text: '今天', value: 'today' },
  { text: '本周', value: 'week' },
  { text: '本月', value: 'month' },
])

// 获取天气图标
function getWeatherIcon(weather) {
  const iconMap = {
    sunny: 'sun-o',
    cloudy: 'cloud-o',
    rainy: 'gem-o',
    snowy: 'snowflake-o',
  }
  return iconMap[weather] || 'sun-o'
}

// 获取内容摘要（去除HTML标签，截取前50个字符）
function getSummary(content) {
  if (!content) return '暂无内容'
  const text = content.replace(/<[^>]*>/g, '').trim()
  return text.length > 50 ? text.substring(0, 50) + '...' : text
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
      project_id: filterProject.value || undefined,
      date_filter: filterDate.value || undefined,
      search: searchKeyword.value || undefined,
    }

    // 适配统一响应格式
    const { data, pagination } = await getConstructionLogs(params)

    if (isRefresh) {
      list.value = data || []
    } else {
      list.value.push(...(data || []))
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
  router.push(`/construction/${id}`)
}

// 跳转创建
function goToCreate() {
  router.push('/construction/create')
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.construction-page {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding-bottom: 50px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: white;
}

.header-title {
  font-size: 18px;
  font-weight: bold;
  color: #323233;
}

.search-bar {
  background: white;
  padding: 8px 0;
}

.log-card {
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

.date-info,
.weather-info {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #646566;
  font-size: 14px;
}

.date-text,
.weather-text {
  margin-left: 4px;
}

.card-title {
  padding: 12px 16px;
}

.card-title h3 {
  font-size: 16px;
  font-weight: bold;
  color: #323233;
  margin: 0 0 8px 0;
}

.log-summary {
  font-size: 14px;
  color: #646566;
  line-height: 1.6;
  margin: 0;
}

.image-preview {
  display: flex;
  gap: 8px;
  padding: 0 16px;
  margin-bottom: 12px;
}

.more-images {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 80px;
  height: 80px;
  background: #f7f8fa;
  border-radius: 4px;
  color: #969799;
  font-size: 14px;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 16px;
  border-top: 1px solid #ebedf0;
  font-size: 12px;
  color: #969799;
}

.author {
  display: flex;
  align-items: center;
  gap: 4px;
}
</style>
