<template>
  <div class="clock-in-page">
    <van-nav-bar title="打卡" left-arrow @click-left="router.back()">
      <template #right>
        <van-icon name="calendar-o" size="20" @click="showBackfillCalendar = true" />
      </template>
    </van-nav-bar>

    <!-- 当前时间和位置 -->
    <div class="time-location-card">
      <div class="current-time">
        <div class="time">{{ currentTime }}</div>
        <div class="date">{{ currentDate }}</div>
        <!-- 补卡模式提示 -->
        <div v-if="isBackfillMode" class="backfill-badge">
          <van-icon name="clock-o" />
          补卡: {{ backfillDate }}
        </div>
      </div>
      <div class="location" v-if="location">
        <van-icon name="location-o" />
        <span>{{ location }}</span>
      </div>
      <div class="location" v-else-if="locationLoading">
        <van-loading size="14" />
        <span>正在获取位置...</span>
      </div>
      <div class="location error" v-else-if="locationError">
        <van-icon name="warning-o" />
        <span>{{ locationError }}</span>
      </div>
    </div>

    <!-- 打卡类型 -->
    <div class="section-title">打卡类型</div>
    <div class="clock-type-grid">
      <div
        v-for="option in clockTypeOptions"
        :key="option.value"
        class="clock-type-btn"
        :class="{ active: selectedClockType === option.value, done: hasClockedInSelectedDate(option.value) }"
        :style="{ borderColor: selectedClockType === option.value ? option.color : (hasClockedInSelectedDate(option.value) ? '#07c160' : '#ebedf0') }"
        @click="selectClockType(option.value)"
      >
        <van-icon
          :name="hasClockedInSelectedDate(option.value) ? 'passed' : option.icon"
          :color="hasClockedInSelectedDate(option.value) ? '#07c160' : (selectedClockType === option.value ? option.color : '#969799')"
          size="24"
        />
        <span class="label">{{ option.label }}</span>
        <span v-if="hasClockedInSelectedDate(option.value)" class="done-text">已打卡</span>
      </div>
    </div>

    <!-- 任务选择 / 工作内容 -->
    <div class="section-title">工作内容</div>
    <div class="work-content-section">
      <!-- 任务列表 -->
      <van-loading v-if="loading" class="loading-center" />
      <template v-else>
        <div v-if="appointments.length > 0" class="task-list">
          <div
            v-for="apt in appointments"
            :key="apt.id"
            class="task-card"
            :class="{ active: selectedAppointment?.id === apt.id }"
            @click="selectAppointment(apt)"
          >
            <div class="task-header">
              <span class="task-no">{{ apt.appointment_no }}</span>
              <van-tag :type="apt.is_urgent ? 'danger' : 'primary'" size="small">
                {{ apt.is_urgent ? '紧急' : '普通' }}
              </van-tag>
            </div>
            <div class="task-content">{{ apt.work_content }}</div>
            <div class="task-info">
              <span><van-icon name="location-o" /> {{ apt.work_location }}</span>
            </div>
          </div>
        </div>

        <!-- 手动填写 -->
        <div class="manual-input">
          <div class="manual-label" @click="toggleManualInput">
            <van-icon :name="showManualInput ? 'arrow-down' : 'arrow'" />
            <span>{{ showManualInput ? '收起手动填写' : '没有任务？点击手动填写' }}</span>
          </div>
          <van-field
            v-if="showManualInput"
            v-model="manualWorkContent"
            type="textarea"
            rows="2"
            autosize
            placeholder="请输入工作内容"
            show-word-limit
            maxlength="200"
          />
        </div>
      </template>
    </div>

    <!-- 加班小时数 -->
    <template v-if="selectedClockType && isOvertimeType(selectedClockType)">
      <div class="section-title">加班小时数</div>
      <div class="overtime-section">
        <van-stepper v-model="overtimeHours" min="0.5" max="24" step="0.5" />
        <span class="hours-text">{{ overtimeHours }} 小时</span>
      </div>
    </template>

    <!-- 照片上传 -->
    <div class="section-title">打卡照片（可选，最多9张）</div>
    <div class="photo-section">
      <van-uploader
        v-model="photoList"
        :max-count="9"
        multiple
        :after-read="handlePhotoAfterRead"
        @delete="handlePhotoDelete"
        :disabled="photoUploading"
      />
      <div class="photo-tip">
        <van-icon name="info-o" />
        <span>可拍照或从相册选择照片，最多9张</span>
      </div>
    </div>

    <!-- 打卡按钮 -->
    <div class="clock-in-footer">
      <van-button
        type="primary"
        size="large"
        :disabled="!canSubmit"
        :loading="submitting"
        @click="handleClockIn"
      >
        {{ submitButtonText }}
      </van-button>
      <van-button
        v-if="isBackfillMode"
        size="large"
        @click="cancelBackfill"
        style="margin-top: 10px;"
      >
        取消补卡
      </van-button>
    </div>

    <!-- 打卡成功提示 -->
    <van-dialog
      v-model:show="showSuccessDialog"
      title="打卡成功"
      :show-confirm-button="true"
    >
      <div class="success-content">
        <van-icon name="passed" color="#07c160" size="48" />
        <div class="success-time">{{ lastClockInTime }}</div>
        <div class="success-type">{{ getAttendanceTypeLabel(lastClockInType) }}</div>
      </div>
    </van-dialog>

    <!-- 补卡日历弹窗 -->
    <van-calendar
      v-model:show="showBackfillCalendar"
      :min-date="minBackfillDate"
      :max-date="maxBackfillDate"
      :default-date="backfillDefaultDate"
      @confirm="onBackfillDateConfirm"
      title="选择补卡日期"
      :show-confirm="true"
      confirm-text="确定补卡"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showLoadingToast, closeToast } from 'vant'
import {
  getTodayAppointments,
  clockIn,
  getAttendanceTypeLabel,
  isOvertimeType,
  uploadImage,
  getMyRecords
} from '@/api/attendance'
import { getAppointments } from '@/api/appointment'
import { getAssetUrl } from '@/utils/request'

const router = useRouter()

const loading = ref(true)
const appointments = ref([])
const selectedAppointment = ref(null)
const selectedClockType = ref('')
const currentTime = ref('')
const currentDate = ref('')
const location = ref('')
const locationLoading = ref(false)
const locationError = ref('')
const latitude = ref(0)
const longitude = ref(0)

const showManualInput = ref(false)
const manualWorkContent = ref('')
const overtimeHours = ref(1)

// 照片相关
const photoList = ref([])
const photoUploading = ref(false)
const photoUrls = ref([])  // 多张照片的URL数组

const submitting = ref(false)
const todayClockedTypes = ref([])

const showSuccessDialog = ref(false)
const lastClockInTime = ref('')
const lastClockInType = ref('')

// 补卡相关
const showBackfillCalendar = ref(false)
const isBackfillMode = ref(false)
const backfillDate = ref('')
const backfillDateClockedTypes = ref([])

let timeInterval = null

// 补卡日期范围：最近30天
const minBackfillDate = computed(() => {
  const date = new Date()
  date.setDate(date.getDate() - 30)
  return date
})
const maxBackfillDate = computed(() => new Date())
const backfillDefaultDate = computed(() => {
  if (backfillDate.value) {
    const parts = backfillDate.value.split('-')
    return new Date(parts[0], parts[1] - 1, parts[2])
  }
  return new Date()
})

const clockTypeOptions = [
  { value: 'morning', label: '上午打卡', icon: 'sun-o', color: '#1989fa' },
  { value: 'afternoon', label: '下午打卡', icon: 'sun-o', color: '#07c160' },
  { value: 'noon_overtime', label: '中午加班', icon: 'clock-o', color: '#ff976a' },
  { value: 'night_overtime', label: '晚上加班', icon: 'moon-o', color: '#7232dd' }
]

// 是否可以提交
const canSubmit = computed(() => {
  if (!selectedClockType.value) return false
  if (hasClockedInSelectedDate(selectedClockType.value)) return false
  if (!selectedAppointment.value && !manualWorkContent.value.trim()) return false
  if (isOvertimeType(selectedClockType.value) && overtimeHours.value <= 0) return false
  return true
})

// 按钮文字
const submitButtonText = computed(() => {
  if (!selectedClockType.value) return '请选择打卡类型'
  if (hasClockedInSelectedDate(selectedClockType.value)) return '已完成该类型打卡'
  if (!selectedAppointment.value && !manualWorkContent.value.trim()) return '请选择任务或填写工作内容'
  const label = getAttendanceTypeLabel(selectedClockType.value)
  return isBackfillMode.value ? `补卡: ${label}` : label
})

onMounted(() => {
  updateTime()
  timeInterval = setInterval(updateTime, 1000)
  getLocation()
  loadData()
})

onUnmounted(() => {
  if (timeInterval) {
    clearInterval(timeInterval)
  }
})

function updateTime() {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
  currentDate.value = now.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    weekday: 'long'
  })
}

function getLocation() {
  locationLoading.value = true
  locationError.value = ''

  if (!navigator.geolocation) {
    locationLoading.value = false
    locationError.value = '浏览器不支持定位'
    return
  }

  navigator.geolocation.getCurrentPosition(
    (position) => {
      locationLoading.value = false
      latitude.value = position.coords.latitude
      longitude.value = position.coords.longitude
      location.value = `${position.coords.latitude.toFixed(6)}, ${position.coords.longitude.toFixed(6)}`
    },
    (error) => {
      locationLoading.value = false
      locationError.value = '无法获取位置'
      console.error('定位失败:', error)
    },
    {
      enableHighAccuracy: true,
      timeout: 10000,
      maximumAge: 60000
    }
  )
}

async function loadData(targetDate = null) {
  loading.value = true
  try {
    let data = []

    if (targetDate) {
      // 加载指定日期的任务
      const response = await getAppointments({
        start_date: targetDate,
        end_date: targetDate,
        page: 1,
        page_size: 100
      })
      data = response.data || []
    } else {
      // 加载今天的任务
      const response = await getTodayAppointments()
      data = response.data || []
    }

    appointments.value = data

    // 提取已打卡类型
    const types = new Set()
    appointments.value.forEach(apt => {
      if (apt.clocked_types) {
        apt.clocked_types.forEach(t => types.add(t))
      }
    })

    // 更新对应的已打卡类型
    if (targetDate) {
      backfillDateClockedTypes.value = Array.from(types)
    } else {
      todayClockedTypes.value = Array.from(types)
    }

    // 默认选中第一个未打卡的类型
    selectFirstAvailableType()

    // 默认选中第一个任务
    if (appointments.value.length > 0) {
      selectedAppointment.value = appointments.value[0]
    } else {
      selectedAppointment.value = null
    }
  } catch (error) {
    console.error('加载数据失败:', error)
    showToast('加载数据失败')
  } finally {
    loading.value = false
  }
}

// 选择第一个可用的打卡类型
function selectFirstAvailableType() {
  const clockedTypes = isBackfillMode.value ? backfillDateClockedTypes.value : todayClockedTypes.value
  for (const option of clockTypeOptions) {
    if (!clockedTypes.includes(option.value)) {
      selectedClockType.value = option.value
      return
    }
  }
  selectedClockType.value = ''
}

// 检查选定日期是否已打卡
function hasClockedInSelectedDate(type) {
  const clockedTypes = isBackfillMode.value ? backfillDateClockedTypes.value : todayClockedTypes.value
  return clockedTypes.includes(type)
}

function selectClockType(type) {
  if (hasClockedInSelectedDate(type)) {
    const dateStr = isBackfillMode.value ? backfillDate.value : '今天'
    showToast(`${dateStr}已完成该类型打卡`)
    return
  }
  selectedClockType.value = type
}

function selectAppointment(apt) {
  selectedAppointment.value = apt
  manualWorkContent.value = ''
  showManualInput.value = false
}

function toggleManualInput() {
  showManualInput.value = !showManualInput.value
  if (showManualInput.value) {
    selectedAppointment.value = null
  }
}

// 补卡日期确认
async function onBackfillDateConfirm(date) {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  backfillDate.value = `${year}-${month}-${day}`

  showBackfillCalendar.value = false
  isBackfillMode.value = true

  // 加载该日期的打卡统计
  await loadBackfillStatistics()

  // 加载补卡日期当天的任务
  await loadData(backfillDate.value)

  // 选择第一个可用的打卡类型
  selectFirstAvailableType()

  showToast(`已选择 ${backfillDate.value} 进行补卡`)
}

// 加载补卡日期的打卡统计
async function loadBackfillStatistics() {
  try {
    // 获取补卡日期的打卡记录
    const response = await getMyRecords({
      start_date: backfillDate.value,
      end_date: backfillDate.value,
      page: 1,
      page_size: 100
    })

    // 提取已打卡的类型
    const types = new Set()
    const records = response.data || []
    records.forEach(record => {
      if (record.attendance_type) {
        types.add(record.attendance_type)
      }
    })

    backfillDateClockedTypes.value = Array.from(types)
  } catch (error) {
    console.error('加载打卡记录失败:', error)
    backfillDateClockedTypes.value = []
  }
}

// 取消补卡模式
async function cancelBackfill() {
  isBackfillMode.value = false
  backfillDate.value = ''
  backfillDateClockedTypes.value = []

  // 重新加载今天的任务
  await loadData()

  selectFirstAvailableType()
}

// 处理照片选择后上传
async function handlePhotoAfterRead(file) {
  // 支持多选，file可能是数组
  const files = Array.isArray(file) ? file : [file]

  photoUploading.value = true

  for (const f of files) {
    const uploadFile = f.file || f

    f.status = 'uploading'
    f.message = '上传中...'

    try {
      const response = await uploadImage(uploadFile)
      console.log('上传响应:', response)

      if (response.data && response.data.url) {
        photoUrls.value.push(response.data.url)
        const fullUrl = getAssetUrl(response.data.url)
        f.status = 'done'
        f.message = '上传成功'
        f.url = fullUrl
      } else {
        throw new Error('上传返回数据格式错误')
      }
    } catch (error) {
      console.error('上传照片失败:', error)
      showToast('上传照片失败: ' + (error.message || '未知错误'))
      f.status = 'failed'
      f.message = '上传失败'
    }
  }

  photoUploading.value = false
}

// 删除照片
function handlePhotoDelete(file, detail) {
  // detail.index 是被删除照片的索引
  const index = detail.index
  if (index >= 0 && index < photoUrls.value.length) {
    photoUrls.value.splice(index, 1)
  }
}

async function handleClockIn() {
  if (!canSubmit.value) return

  submitting.value = true
  showLoadingToast({
    message: isBackfillMode.value ? '补卡中...' : '打卡中...',
    forbidClick: true,
    duration: 0
  })

  try {
    const data = {
      attendance_type: selectedClockType.value,
      clock_in_location: location.value,
      clock_in_latitude: latitude.value,
      clock_in_longitude: longitude.value,
      overtime_hours: isOvertimeType(selectedClockType.value) ? overtimeHours.value : 0,
      photo_urls: photoUrls.value  // 多张照片URL数组
    }

    // 补卡模式下设置打卡时间
    if (isBackfillMode.value && backfillDate.value) {
      // 根据打卡类型设置默认时间
      let hour = 9
      if (selectedClockType.value === 'afternoon') hour = 14
      else if (selectedClockType.value === 'noon_overtime') hour = 12
      else if (selectedClockType.value === 'night_overtime') hour = 19

      data.clock_in_time = `${backfillDate.value} ${String(hour).padStart(2, '0')}:00:00`
    }

    console.log('打卡数据:', data)

    if (selectedAppointment.value) {
      data.appointment_id = selectedAppointment.value.id
    } else {
      data.work_content = manualWorkContent.value
    }

    await clockIn(data)
    closeToast()

    // 更新状态
    if (isBackfillMode.value) {
      backfillDateClockedTypes.value.push(selectedClockType.value)
    } else {
      todayClockedTypes.value.push(selectedClockType.value)
    }

    // 选择下一个未打卡类型
    selectFirstAvailableType()

    // 显示成功
    lastClockInTime.value = isBackfillMode.value
      ? `${backfillDate.value}`
      : new Date().toLocaleTimeString('zh-CN', {
          hour: '2-digit',
          minute: '2-digit',
          second: '2-digit'
        })
    lastClockInType.value = selectedClockType.value
    showSuccessDialog.value = true

    // 清空照片
    photoList.value = []
    photoUrls.value = []

  } catch (error) {
    closeToast()
    showToast(error.response?.data?.message || error.message || '打卡失败')
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.clock-in-page {
  min-height: 100vh;
  background: #f7f8fa;
  padding-bottom: 120px;
}

.time-location-card {
  background: linear-gradient(135deg, #1989fa 0%, #07c160 100%);
  padding: 20px;
  color: white;
  text-align: center;
}

.current-time .time {
  font-size: 32px;
  font-weight: bold;
  letter-spacing: 2px;
}

.current-time .date {
  font-size: 14px;
  margin-top: 5px;
  opacity: 0.9;
}

.backfill-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  margin-top: 8px;
  padding: 4px 12px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  font-size: 13px;
}

.location {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 10px;
  font-size: 12px;
  opacity: 0.9;
}

.location .van-icon {
  margin-right: 4px;
}

.location.error {
  color: #ffcd42;
}

.section-title {
  padding: 15px 15px 10px;
  font-size: 14px;
  color: #969799;
  background: #f7f8fa;
}

.clock-type-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 10px;
  padding: 0 10px;
  background: white;
}

.clock-type-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 15px 5px;
  border-radius: 8px;
  border: 2px solid;
  background: white;
}

.clock-type-btn.active {
  background: #f0f7ff;
}

.clock-type-btn.done {
  background: #f0fff0;
}

.clock-type-btn .label {
  margin-top: 8px;
  font-size: 12px;
  color: #323233;
}

.clock-type-btn .done-text {
  font-size: 10px;
  color: #07c160;
  margin-top: 2px;
}

.work-content-section {
  background: white;
  padding: 10px;
}

.loading-center {
  display: flex;
  justify-content: center;
  padding: 20px;
}

.task-list {
  margin-bottom: 10px;
}

.task-card {
  padding: 12px;
  margin-bottom: 8px;
  border-radius: 8px;
  border: 2px solid #ebedf0;
  background: white;
}

.task-card.active {
  border-color: #1989fa;
  background: #f0f7ff;
}

.task-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.task-no {
  font-size: 12px;
  color: #969799;
}

.task-content {
  font-size: 14px;
  color: #323233;
  margin-bottom: 6px;
  line-height: 1.4;
}

.task-info {
  font-size: 12px;
  color: #969799;
}

.task-info .van-icon {
  vertical-align: middle;
  margin-right: 2px;
}

.manual-input {
  border-top: 1px solid #ebedf0;
  padding-top: 10px;
}

.manual-label {
  display: flex;
  align-items: center;
  color: #1989fa;
  font-size: 14px;
  padding: 5px 0;
}

.manual-label .van-icon {
  margin-right: 5px;
}

.overtime-section {
  background: white;
  padding: 15px;
  display: flex;
  align-items: center;
  gap: 15px;
}

.hours-text {
  font-size: 14px;
  color: #323233;
}

.photo-section {
  background: white;
  padding: 15px;
}

.photo-section :deep(.van-uploader) {
  display: block;
}

.photo-section :deep(.van-uploader__wrapper) {
  justify-content: flex-start;
}

.photo-tip {
  margin-top: 10px;
  font-size: 12px;
  color: #969799;
  display: flex;
  align-items: center;
  gap: 4px;
}

.preview-cover {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  background: rgba(0, 0, 0, 0.5);
  color: white;
  font-size: 12px;
  text-align: center;
  padding: 4px 0;
}

.clock-in-footer {
  position: fixed;
  bottom: 50px;
  left: 0;
  right: 0;
  padding: 10px 15px;
  background: white;
  box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.05);
}

.success-content {
  padding: 20px;
  text-align: center;
}

.success-time {
  margin-top: 15px;
  font-size: 18px;
  font-weight: bold;
}

.success-type {
  margin-top: 5px;
  font-size: 14px;
  color: #969799;
}
</style>
