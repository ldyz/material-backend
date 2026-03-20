<template>
  <div class="calendar-board">
    <van-nav-bar title="施工预约日历" left-arrow @click-left="handleBack" />

    <!-- 日历区域 -->
    <div class="calendar-wrapper">
      <!-- 自定义日历头部 -->
      <div class="calendar-header">
        <van-icon name="arrow-left" class="nav-icon" @click="prevMonth" />
        <span class="month-title">{{ currentMonthStr }}</span>
        <van-icon name="arrow-right" class="nav-icon" @click="nextMonth" />
      </div>

      <!-- 星期头部 -->
      <div class="weekdays">
        <span v-for="day in weekdays" :key="day" class="weekday">{{ day }}</span>
      </div>

      <!-- 日期网格 -->
      <div class="days-grid">
        <div
          v-for="(day, index) in calendarDays"
          :key="index"
          class="day-item"
          :class="{
            'other-month': day.otherMonth,
            'selected': isSelected(day.date),
            'today': isToday(day.date)
          }"
          @click="onDateSelect(day.date)"
        >
          <div class="day-content" :style="getDayStyle(day.date)">
            <span class="day-number">{{ day.date.getDate() }}</span>
            <span v-if="getTaskCount(day.date) > 0" class="task-badge">
              {{ getTaskCount(day.date) }}
            </span>
          </div>
        </div>
      </div>

      <!-- 新建预约按钮 -->
      <div class="create-btn-wrapper">
        <van-button
          type="primary"
          icon="plus"
          round
          size="large"
          @click="goToCreate"
          class="create-btn"
        >
          新建预约
        </van-button>
      </div>

      <!-- 图例 -->
      <div class="legend">
        <div class="legend-item">
          <span class="legend-color" style="background: #f5f5f5;"></span>
          <span>无任务</span>
        </div>
        <div class="legend-item">
          <span class="legend-color" style="background: linear-gradient(135deg, #a8e6cf, #88d4ab);"></span>
          <span>1-2个</span>
        </div>
        <div class="legend-item">
          <span class="legend-color" style="background: linear-gradient(135deg, #56c596, #3db88a);"></span>
          <span>3-5个</span>
        </div>
        <div class="legend-item">
          <span class="legend-color" style="background: linear-gradient(135deg, #2e7d32, #1b5e20);"></span>
          <span>6+个</span>
        </div>
      </div>
    </div>

    <!-- 选中日期的任务列表 -->
    <div class="task-section">
      <div class="task-header">
        <div class="header-left">
          <span class="selected-date">{{ selectedDateStr }}</span>
          <span v-if="isToday(selectedDate)" class="today-tag">今天</span>
        </div>
        <span class="task-count-label">
          <van-icon name="todo-list-o" />
          {{ tasks.length }} 个任务
        </span>
      </div>

      <ListContainer
        v-model:loading="loading"
        v-model:refreshing="refreshing"
        :finished="true"
        :data="tasks"
        :show-empty-state="tasks.length === 0 && !loading"
        empty-text="当日暂无施工预约"
        @refresh="loadTasks"
      >
        <div class="task-list">
          <ListItemCard
            v-for="task in tasks"
            :key="task.id"
            :item="task"
            type="appointment"
            @click="goToDetail"
          />
        </div>
      </ListContainer>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { getDailyStatistics, getAppointments } from '@/api/appointment'
import { formatDate } from '@/composables/useDateTime'
import ListContainer from '@/components/common/ListContainer.vue'
import ListItemCard from '@/components/common/ListItemCard.vue'

const router = useRouter()

// 星期标题
const weekdays = ['日', '一', '二', '三', '四', '五', '六']

// 当前显示的月份
const currentYear = ref(new Date().getFullYear())
const currentMonth = ref(new Date().getMonth())

// 当前选中日期
const selectedDate = ref(new Date())

// 统计数据
const statistics = ref([])
const loading = ref(false)
const refreshing = ref(false)
const tasks = ref([])

// 当前月份字符串
const currentMonthStr = computed(() => {
  return `${currentYear.value}年${currentMonth.value + 1}月`
})

// 选中日期字符串
const selectedDateStr = computed(() => {
  const d = selectedDate.value
  return `${d.getFullYear()}年${d.getMonth() + 1}月${d.getDate()}日`
})

// 计算日历天数
const calendarDays = computed(() => {
  const year = currentYear.value
  const month = currentMonth.value
  const days = []

  // 当月第一天
  const firstDay = new Date(year, month, 1)
  const firstDayWeek = firstDay.getDay()

  // 上个月的天数填充
  const prevMonth = new Date(year, month, 0)
  const prevMonthDays = prevMonth.getDate()
  for (let i = firstDayWeek - 1; i >= 0; i--) {
    const date = new Date(year, month - 1, prevMonthDays - i)
    days.push({ date, otherMonth: true })
  }

  // 当月天数
  const currentMonthDays = new Date(year, month + 1, 0).getDate()
  for (let i = 1; i <= currentMonthDays; i++) {
    const date = new Date(year, month, i)
    days.push({ date, otherMonth: false })
  }

  // 下个月的天数填充（补齐6行）
  const remainingDays = 42 - days.length
  for (let i = 1; i <= remainingDays; i++) {
    const date = new Date(year, month + 1, i)
    days.push({ date, otherMonth: true })
  }

  return days
})

// 判断是否是今天
function isToday(date) {
  const today = new Date()
  return date.getFullYear() === today.getFullYear() &&
         date.getMonth() === today.getMonth() &&
         date.getDate() === today.getDate()
}

// 判断是否选中
function isSelected(date) {
  return date.getFullYear() === selectedDate.value.getFullYear() &&
         date.getMonth() === selectedDate.value.getMonth() &&
         date.getDate() === selectedDate.value.getDate()
}

// 获取日期的任务数量
function getTaskCount(date) {
  const dateStr = formatDateToISO(date)
  const stat = statistics.value.find(s => s.date === dateStr)
  return stat?.total_count || 0
}

// 获取日期的样式（根据任务数量）
function getDayStyle(date) {
  const count = getTaskCount(date)

  if (count === 0) {
    return {}
  } else if (count <= 2) {
    return {
      background: 'linear-gradient(135deg, #a8e6cf 0%, #88d4ab 100%)',
      color: '#1b5e20',
      boxShadow: '0 2px 8px rgba(168, 230, 207, 0.5)'
    }
  } else if (count <= 5) {
    return {
      background: 'linear-gradient(135deg, #56c596 0%, #3db88a 100%)',
      color: '#ffffff',
      boxShadow: '0 2px 8px rgba(86, 197, 150, 0.5)'
    }
  } else {
    return {
      background: 'linear-gradient(135deg, #2e7d32 0%, #1b5e20 100%)',
      color: '#ffffff',
      boxShadow: '0 2px 8px rgba(46, 125, 50, 0.5)'
    }
  }
}

// 格式化日期为 ISO 格式 (YYYY-MM-DD)
function formatDateToISO(date) {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

// 获取当月统计数据
async function loadStatistics() {
  try {
    const startDate = new Date(currentYear.value, currentMonth.value, 1)
    const endDate = new Date(currentYear.value, currentMonth.value + 1, 0)

    const response = await getDailyStatistics({
      start_date: formatDateToISO(startDate),
      end_date: formatDateToISO(endDate)
    })

    statistics.value = response.data?.statistics || response.statistics || []
  } catch (error) {
    console.error('加载统计数据失败:', error)
    statistics.value = []
  }
}

// 加载选中日期的任务
async function loadTasks() {
  loading.value = true
  try {
    const dateStr = formatDateToISO(selectedDate.value)
    const response = await getAppointments({
      start_date: dateStr,
      end_date: dateStr,
      page: 1,
      page_size: 100
    })

    tasks.value = response.data || []
  } catch (error) {
    console.error('加载任务列表失败:', error)
    tasks.value = []
  } finally {
    loading.value = false
    refreshing.value = false
  }
}

// 日期选择
function onDateSelect(date) {
  selectedDate.value = date
  loadTasks()
}

// 跳转到新建预约页面
function goToCreate() {
  const dateStr = formatDateToISO(selectedDate.value)
  router.push({
    path: '/appointment/create',
    query: { work_date: dateStr }
  })
}

// 上一月
function prevMonth() {
  if (currentMonth.value === 0) {
    currentMonth.value = 11
    currentYear.value--
  } else {
    currentMonth.value--
  }
  loadStatistics()
}

// 下一月
function nextMonth() {
  if (currentMonth.value === 11) {
    currentMonth.value = 0
    currentYear.value++
  } else {
    currentMonth.value++
  }
  loadStatistics()
}

// 跳转到详情页
function goToDetail(task) {
  router.push(`/appointment/${task.id}`)
}

// 返回处理 - 直接返回首页，避免返回两次
function handleBack() {
  // 检查是否能返回到非 tabbar 页面，如果可以就返回，否则直接跳转到首页
  if (window.history.length > 1) {
    const prevRoute = window.history.state
    // 如果上一个路由是 tabbar 页面(根路径 /)，直接返回
    if (prevRoute && prevRoute.fullPath === '/') {
      router.back()
      return
    }
  }
  // 其他情况直接返回首页
  router.replace('/')
}

onMounted(() => {
  loadStatistics()
  loadTasks()
})
</script>

<style scoped>
.calendar-board {
  min-height: 100vh;
  background: linear-gradient(180deg, #f8fffe 0%, #f5f5f5 100%);
  display: flex;
  flex-direction: column;
}

.calendar-wrapper {
  background: #ffffff;
  margin: 12px;
  border-radius: 16px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.06);
  overflow: hidden;
}

.calendar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: linear-gradient(135deg, #4caf50 0%, #2e7d32 100%);
  color: #ffffff;
}

.month-title {
  font-size: 18px;
  font-weight: 600;
  letter-spacing: 1px;
}

.nav-icon {
  font-size: 20px;
  padding: 8px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.15);
  cursor: pointer;
  transition: all 0.2s;
}

.nav-icon:active {
  background: rgba(255, 255, 255, 0.3);
  transform: scale(0.95);
}

.weekdays {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  padding: 12px 8px;
  background: #fafafa;
  border-bottom: 1px solid #f0f0f0;
}

.weekday {
  text-align: center;
  font-size: 12px;
  font-weight: 500;
  color: #969799;
}

.days-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 4px;
  padding: 8px;
}

.day-item {
  aspect-ratio: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
}

.day-item.other-month {
  opacity: 0.3;
}

.day-item.today .day-content {
  border: 2px solid #4caf50;
}

.day-item.selected .day-content {
  transform: scale(1.1);
  border: 2px solid #1976d2 !important;
  box-shadow: 0 4px 12px rgba(25, 118, 210, 0.3) !important;
}

.day-content {
  width: 40px;
  height: 40px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
  background: #f5f5f5;
  transition: all 0.2s;
  position: relative;
}

.day-number {
  font-size: 14px;
  font-weight: 500;
  line-height: 1;
}

.task-badge {
  position: absolute;
  top: -4px;
  right: -4px;
  min-width: 16px;
  height: 16px;
  padding: 0 4px;
  font-size: 10px;
  font-weight: 600;
  color: #fff;
  background: linear-gradient(135deg, #ff6b6b, #ee5a5a);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 4px rgba(238, 90, 90, 0.4);
}

.legend {
  display: flex;
  justify-content: center;
  gap: 16px;
  padding: 12px 16px;
  background: #fafafa;
  border-top: 1px solid #f0f0f0;
}

.create-btn-wrapper {
  padding: 12px 16px;
  background: #fafafa;
}

.create-btn {
  width: 100%;
  height: 44px;
  font-size: 15px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 11px;
  color: #646566;
}

.legend-color {
  width: 12px;
  height: 12px;
  border-radius: 4px;
}

.task-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: #ffffff;
  margin: 0 12px 12px;
  border-radius: 16px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.06);
  overflow: hidden;
}

.task-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 14px 16px;
  background: linear-gradient(135deg, #e8f5e9 0%, #c8e6c9 100%);
  border-bottom: 1px solid #e0e0e0;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.selected-date {
  font-size: 15px;
  font-weight: 600;
  color: #2e7d32;
}

.today-tag {
  font-size: 11px;
  padding: 2px 8px;
  background: linear-gradient(135deg, #4caf50, #2e7d32);
  color: #fff;
  border-radius: 10px;
}

.task-count-label {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #56c596;
  font-weight: 500;
}

.task-list {
  padding: 8px;
  background: #fafafa;
  min-height: 200px;
}

/* 列表项间距 */
.task-list :deep(.list-item-card) {
  margin-bottom: 8px;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

/* 空状态样式 */
.task-list :deep(.van-empty) {
  padding: 40px 0;
}

.task-list :deep(.van-empty__description) {
  color: #969799;
}
</style>
