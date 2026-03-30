<template>
  <div class="record-list-page">
    <van-nav-bar title="打卡记录" left-arrow @click-left="router.back()" />

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
            <span v-if="getClockInCount(day.date) > 0" class="clock-badge">
              {{ getClockInCount(day.date) }}
            </span>
          </div>
        </div>
      </div>

      <!-- 图例 -->
      <div class="legend">
        <div class="legend-item">
          <span class="legend-color" style="background: #f5f5f5;"></span>
          <span>无打卡</span>
        </div>
        <div class="legend-item">
          <span class="legend-color" style="background: linear-gradient(135deg, #a8d8ff, #7cc4ff);"></span>
          <span>1-2次</span>
        </div>
        <div class="legend-item">
          <span class="legend-color" style="background: linear-gradient(135deg, #4fc3f7, #29b6f6);"></span>
          <span>3-4次</span>
        </div>
        <div class="legend-item">
          <span class="legend-color" style="background: linear-gradient(135deg, #0288d1, #01579b);"></span>
          <span>5+次</span>
        </div>
      </div>
    </div>

    <!-- 选中日期的打卡记录 -->
    <div class="record-section">
      <div class="record-header">
        <div class="header-left">
          <span class="selected-date">{{ selectedDateStr }}</span>
          <span v-if="isToday(selectedDate)" class="today-tag">今天</span>
        </div>
        <span class="record-count-label">
          <van-icon name="clock-o" />
          {{ records.length }} 条记录
        </span>
      </div>

      <ListContainer
        v-model:loading="loading"
        v-model:refreshing="refreshing"
        :finished="true"
        :data="records"
        :show-empty-state="records.length === 0 && !loading"
        empty-text="当日暂无打卡记录"
        @refresh="loadRecords"
      >
        <div class="record-list">
          <div
            v-for="record in records"
            :key="record.id"
            class="record-card"
          >
            <div class="record-header-inner">
              <van-tag :color="getAttendanceTypeColor(record.attendance_type)">
                {{ getAttendanceTypeLabel(record.attendance_type) }}
              </van-tag>
              <van-tag :color="getStatusColor(record.status)" plain>
                {{ getStatusLabel(record.status) }}
              </van-tag>
            </div>
            <div class="record-body">
              <div class="record-time">
                <van-icon name="clock-o" />
                <span>{{ formatClockInTime(record.clock_in_time) }}</span>
              </div>
              <div class="record-task" v-if="record.appointment_no">
                <van-icon name="orders-o" />
                <span>{{ record.appointment_no }}</span>
              </div>
              <div class="record-content" v-if="record.work_content">
                {{ record.work_content }}
              </div>
              <div class="record-location" v-if="record.clock_in_location">
                <van-icon name="location-o" />
                <span>{{ record.clock_in_location }}</span>
              </div>
              <div class="record-overtime" v-if="record.overtime_hours > 0">
                <van-icon name="clock-o" />
                <span>加班 {{ record.overtime_hours }} 小时</span>
              </div>
              <div class="record-photos" v-if="getPhotoList(record).length > 0">
                <van-image
                  v-for="(photo, idx) in getPhotoList(record)"
                  :key="idx"
                  :src="photo"
                  width="60"
                  height="60"
                  fit="cover"
                  radius="4"
                  @click="previewImages(getPhotoList(record), idx)"
                />
              </div>
            </div>
            <div class="record-footer" v-if="record.status === 'rejected'">
              <span class="reject-reason">驳回原因: {{ record.confirmed_remark }}</span>
            </div>
          </div>
        </div>
      </ListContainer>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onActivated } from 'vue'
import { useRouter } from 'vue-router'
import { showImagePreview } from 'vant'
import {
  getMyRecords,
  getCalendarStatistics,
  getAttendanceTypeLabel,
  getAttendanceTypeColor,
  getStatusLabel,
  getStatusColor,
  formatClockInTime
} from '@/api/attendance'
import { getAssetUrl } from '@/utils/request'
import ListContainer from '@/components/common/ListContainer.vue'
import { logger } from '@/utils/logger'

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
const records = ref([])

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

// 格式化日期为 ISO 格式 (YYYY-MM-DD)
function formatDateToISO(date) {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

// 获取日期的打卡次数
function getClockInCount(date) {
  const dateStr = formatDateToISO(date)
  const stat = statistics.value.find(s => s.date === dateStr)
  return stat?.total_count || 0
}

// 获取日期的样式（根据打卡次数）
function getDayStyle(date) {
  const count = getClockInCount(date)

  if (count === 0) {
    return {}
  } else if (count <= 2) {
    return {
      background: 'linear-gradient(135deg, #a8d8ff 0%, #7cc4ff 100%)',
      color: '#01579b',
      boxShadow: '0 2px 8px rgba(168, 216, 255, 0.5)'
    }
  } else if (count <= 4) {
    return {
      background: 'linear-gradient(135deg, #4fc3f7 0%, #29b6f6 100%)',
      color: '#ffffff',
      boxShadow: '0 2px 8px rgba(79, 195, 247, 0.5)'
    }
  } else {
    return {
      background: 'linear-gradient(135deg, #0288d1 0%, #01579b 100%)',
      color: '#ffffff',
      boxShadow: '0 2px 8px rgba(2, 136, 209, 0.5)'
    }
  }
}

// 获取当月统计数据
async function loadStatistics() {
  try {
    const startDate = new Date(currentYear.value, currentMonth.value, 1)
    const endDate = new Date(currentYear.value, currentMonth.value + 1, 0)

    const response = await getCalendarStatistics({
      start_date: formatDateToISO(startDate),
      end_date: formatDateToISO(endDate)
    })

    statistics.value = response.data || []
  } catch (error) {
    logger.error('加载统计数据失败:', error)
    statistics.value = []
  }
}

// 加载选中日期的打卡记录
async function loadRecords() {
  loading.value = true
  try {
    const dateStr = formatDateToISO(selectedDate.value)
    const response = await getMyRecords({
      start_date: dateStr,
      end_date: dateStr,
      page: 1,
      page_size: 100
    })

    records.value = response.data || []
  } catch (error) {
    logger.error('加载打卡记录失败:', error)
    records.value = []
  } finally {
    loading.value = false
    refreshing.value = false
  }
}

// 日期选择
function onDateSelect(date) {
  selectedDate.value = date
  loadRecords()
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

// 预览图片
function previewImage(url) {
  const fullUrl = getAssetUrl(url)
  showImagePreview([fullUrl])
}

// 获取照片列表
function getPhotoList(record) {
  const photos = []
  if (record.photo_urls) {
    try {
      const urls = typeof record.photo_urls === 'string' ? JSON.parse(record.photo_urls) : record.photo_urls
      if (Array.isArray(urls)) {
        urls.forEach(url => {
          if (url && url !== 'null' && url !== '') {
            // 直接拼接完整 URL，不使用 getAssetUrl（避免返回默认头像）
            if (url.startsWith('http')) {
              photos.push(url)
            } else {
              const isCapacitorEnv = typeof window !== 'undefined' && window.Capacitor
              photos.push(isCapacitorEnv ? 'https://home.mbed.org.cn:9090' + url : url)
            }
          }
        })
      }
    } catch (e) {
      console.error('解析 photo_urls 失败:', e)
    }
  }
  // 兼容单张照片
  if (photos.length === 0 && record.photo_url && record.photo_url !== 'null' && record.photo_url !== '') {
    const url = record.photo_url
    if (url.startsWith('http')) {
      photos.push(url)
    } else {
      const isCapacitorEnv = typeof window !== 'undefined' && window.Capacitor
      photos.push(isCapacitorEnv ? 'https://home.mbed.org.cn:9090' + url : url)
    }
  }
  return photos
}

// 预览多张图片（从指定位置开始）
function previewImages(images, startIndex = 0) {
  showImagePreview({
    images,
    startPosition: startIndex
  })
}

onMounted(() => {
  loadStatistics()
  loadRecords()
})

onActivated(() => {
  loadStatistics()
  loadRecords()
})
</script>

<style scoped>
.record-list-page {
  min-height: 100vh;
  background: linear-gradient(180deg, #f0f7ff 0%, #f5f5f5 100%);
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
  background: linear-gradient(135deg, #1976d2 0%, #0d47a1 100%);
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
  border: 2px solid #1976d2;
}

.day-item.selected .day-content {
  transform: scale(1.1);
  border: 2px solid #ff9800 !important;
  box-shadow: 0 4px 12px rgba(255, 152, 0, 0.3) !important;
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

.clock-badge {
  position: absolute;
  top: -4px;
  right: -4px;
  min-width: 16px;
  height: 16px;
  padding: 0 4px;
  font-size: 10px;
  font-weight: 600;
  color: #fff;
  background: linear-gradient(135deg, #ff9800, #f57c00);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 4px rgba(255, 152, 0, 0.4);
}

.legend {
  display: flex;
  justify-content: center;
  gap: 16px;
  padding: 12px 16px;
  background: #fafafa;
  border-top: 1px solid #f0f0f0;
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

.record-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: #ffffff;
  margin: 0 12px 12px;
  border-radius: 16px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.06);
  overflow: hidden;
}

.record-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 14px 16px;
  background: linear-gradient(135deg, #e3f2fd 0%, #bbdefb 100%);
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
  color: #0d47a1;
}

.today-tag {
  font-size: 11px;
  padding: 2px 8px;
  background: linear-gradient(135deg, #1976d2, #0d47a1);
  color: #fff;
  border-radius: 10px;
}

.record-count-label {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #1976d2;
  font-weight: 500;
}

.record-list {
  padding: 8px;
  background: #fafafa;
  min-height: 200px;
}

.record-card {
  margin-bottom: 10px;
  padding: 12px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.record-header-inner {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
}

.record-body {
  font-size: 14px;
  color: #323233;
}

.record-time,
.record-task,
.record-location,
.record-overtime {
  display: flex;
  align-items: center;
  gap: 5px;
  margin-bottom: 5px;
  font-size: 13px;
  color: #646566;
}

.record-content {
  font-size: 14px;
  color: #323233;
  margin-bottom: 5px;
  line-height: 1.5;
}

.record-photos {
  margin-top: 8px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.record-footer {
  margin-top: 10px;
  padding-top: 10px;
  border-top: 1px solid #ebedf0;
}

.reject-reason {
  font-size: 12px;
  color: #ee0a24;
}

/* 空状态样式 */
.record-list :deep(.van-empty) {
  padding: 40px 0;
}

.record-list :deep(.van-empty__description) {
  color: #969799;
}
</style>
