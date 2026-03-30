<template>
  <el-dialog
    v-model="dialogVisible"
    title="预约日历"
    width="1000px"
    fullscreen
    @close="handleClose"
  >
    <div class="calendar-toolbar">
      <el-radio-group v-model="viewMode" @change="handleViewChange">
        <el-radio-button label="day">日视图</el-radio-button>
        <el-radio-button label="week">周视图</el-radio-button>
        <el-radio-button label="month">月视图</el-radio-button>
      </el-radio-group>

      <div class="date-navigator">
        <el-button :icon="ArrowLeft" @click="previousPeriod">上一个</el-button>
        <span class="current-date">{{ currentDateRange }}</span>
        <el-button @click="nextPeriod">下一个 <el-icon><ArrowRight /></el-icon></el-button>
        <el-button style="margin-left: 8px" @click="goToToday">今天</el-button>
      </div>

      <el-select
        v-model="selectedWorkerId"
        filterable
        clearable
        placeholder="全部作业人员"
        style="width: 200px"
        @change="loadData"
      >
        <el-option
          v-for="worker in workerList"
          :key="worker.id"
          :label="worker.name"
          :value="worker.id"
        />
      </el-select>
    </div>

    <div class="calendar-container" v-loading="loading">
      <!-- 日视图 -->
      <div v-if="viewMode === 'day'" class="day-view">
        <div class="day-header-large">
          <h3>{{ formatDayHeader(currentDate) }}</h3>
          <span class="day-stats">{{ getDayStats(currentDate) }}</span>
        </div>
        <div class="day-slots">
          <div
            v-for="slot in dayTimeSlots"
            :key="slot.time_slot"
            class="day-slot-card"
            :class="{ 'has-appointment': slot.appointments && slot.appointments.length > 0 }"
          >
            <div class="slot-header">
              <span class="slot-title">{{ slot.label }}</span>
              <span class="slot-count" v-if="slot.appointments && slot.appointments.length > 0">
                {{ slot.appointments.length }} 个任务
              </span>
            </div>
            <div class="slot-content">
              <div v-if="slot.appointments && slot.appointments.length > 0">
                <div
                  v-for="apt in slot.appointments"
                  :key="apt.id"
                  class="appointment-card"
                  :class="{ 'urgent': apt.is_urgent }"
                  @click="showAppointment(apt)"
                >
                  <div class="apt-header">
                    <el-tag :type="getStatusType(apt.status)" size="small">
                      {{ getStatusLabel(apt.status) }}
                    </el-tag>
                    <el-tag v-if="apt.is_urgent" type="danger" size="small">加急</el-tag>
                    <el-tag v-if="apt.time_slot === 'full_day'" type="info" size="small">全天</el-tag>
                  </div>
                  <div class="apt-content">
                    <div class="apt-location">{{ apt.work_location }}</div>
                    <div class="apt-detail">{{ apt.work_content }}</div>
                  </div>
                  <div class="apt-footer">
                    <span v-if="apt.assigned_worker_names || apt.assigned_worker_name" class="apt-worker">
                      <el-icon><User /></el-icon>
                      {{ apt.assigned_worker_names || apt.assigned_worker_name }}
                    </span>
                    <span class="apt-no">{{ apt.appointment_no }}</span>
                  </div>
                </div>
              </div>
              <div v-else class="slot-empty-large">
                <el-empty description="暂无任务" :image-size="60" />
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 周视图 -->
      <div v-if="viewMode === 'week'" class="week-view">
        <el-row :gutter="10">
          <el-col
            v-for="(day, index) in weekDays"
            :key="index"
            :span="3"
            class="day-column"
          >
            <div class="day-header">
              <span class="day-name">{{ day.dayName }}</span>
              <span class="day-date">{{ day.date }}</span>
            </div>
            <div class="day-content">
              <div
                v-for="slot in day.slots"
                :key="slot.time_slot"
                class="time-slot"
                :class="{ 'has-appointment': slot.appointments && slot.appointments.length > 0 }"
              >
                <div class="slot-label">{{ slot.label }}</div>
                <div v-if="slot.appointments && slot.appointments.length > 0">
                  <div
                    v-for="apt in slot.appointments"
                    :key="apt.id"
                    class="appointment-item-simple"
                    @click="showAppointment(apt)"
                  >
                    <el-tag :type="getStatusType(apt.status)" size="small">
                      {{ getTimeSlotShortLabel(slot.time_slot) }}
                    </el-tag>
                    <span class="appointment-summary">{{ apt.work_location }}</span>
                    <span v-if="apt.assigned_worker_name" class="appointment-worker-small">
                      {{ apt.assigned_worker_name }}
                    </span>
                    <el-tag v-if="apt.is_urgent" type="danger" size="small">急</el-tag>
                  </div>
                </div>
                <div v-else class="slot-empty">
                  <span class="slot-empty-text">空闲</span>
                </div>
              </div>
            </div>
          </el-col>
        </el-row>
      </div>

      <!-- 月视图 -->
      <div v-else-if="viewMode === 'month'" class="month-view">
        <!-- 任务状态图例 -->
        <div class="legend-section">
          <div class="legend-group">
            <div class="legend-title">任务状态：</div>
            <div class="legend-items">
              <span class="legend-item">
                <span class="legend-box" style="background: linear-gradient(135deg, #fff1f0 0%, #ffccc7 100%); border: 1px solid #ffa39e;"></span>
                加急任务
              </span>
              <span class="legend-item">
                <span class="legend-box" style="background: linear-gradient(135deg, #fffbe6 0%, #ffe58f 100%); border: 1px solid #ffd666;"></span>
                待审批
              </span>
              <span class="legend-item">
                <span class="legend-box" style="background: linear-gradient(135deg, #e6f7ff 0%, #bae7ff 100%); border: 1px solid #91d5ff;"></span>
                进行中
              </span>
              <span class="legend-item">
                <span class="legend-box" style="background: linear-gradient(135deg, #f6ffed 0%, #d9f7be 100%); border: 1px solid #b7eb8f;"></span>
                已排期
              </span>
              <span class="legend-item">
                <span class="legend-box" style="background: linear-gradient(135deg, #f5f5f5 0%, #e8e8e8 100%); border: 1px solid #d9d9d9;"></span>
                已完成
              </span>
            </div>
          </div>

          <!-- 人力状态图例 -->
          <div class="legend-group">
            <div class="legend-title">人力状态：</div>
            <div class="legend-items">
              <span class="legend-item">
                <span class="legend-box legend-available"></span>
                空闲 (任务 &lt; 30%人力)
              </span>
              <span class="legend-item">
                <span class="legend-box legend-normal"></span>
                正常 (30%-75%人力)
              </span>
              <span class="legend-item">
                <span class="legend-box legend-busy"></span>
                紧张 (75%-100%人力)
              </span>
              <span class="legend-item">
                <span class="legend-box legend-overload"></span>
                超载 (任务 ≥ 100%人力)
              </span>
            </div>
          </div>
        </div>
        <el-calendar v-model="currentDate">
          <template #date-cell="{ data }">
            <div
              class="calendar-day"
              :class="[getDayClass(data.day), getWorkloadClass(data.day)]"
              @click="handleDayClick(data)"
            >
              <div class="day-number">
                {{ data.day.split('-').slice(-1)[0] }}
                <span v-if="getDayWorkloadInfo(data.day)" class="workload-indicator" :class="getWorkloadClass(data.day)">
                  {{ getDayWorkloadInfo(data.day)?.text }}
                </span>
              </div>
              <div class="day-appointments">
                <div
                  v-for="apt in getAppointmentsForDay(data.day)"
                  :key="apt.id"
                  class="day-appointment-item"
                  @click.stop="showAppointment(apt)"
                >
                  <el-tag :type="getStatusType(apt.status)" size="small">
                    {{ apt.time_slot }}
                  </el-tag>
                  <span class="appointment-summary">{{ apt.work_location }}</span>
                  <span v-if="apt.assigned_worker_name" class="appointment-worker-small">
                    {{ apt.assigned_worker_name }}
                  </span>
                  <el-tag v-if="apt.is_urgent" type="danger" size="small">急</el-tag>
                </div>
              </div>
            </div>
          </template>
        </el-calendar>
      </div>
    </div>

    <!-- 预约详情对话框 -->
    <AppointmentDetailDialog
      v-model="detailVisible"
      :appointment-id="selectedAppointmentId"
    />

    <template #footer>
      <el-button @click="handleClose">关闭</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { ArrowLeft, ArrowRight, User } from '@element-plus/icons-vue'
import { appointmentApi } from '@/api'
import AppointmentDetailDialog from './AppointmentDetailDialog.vue'

const props = defineProps({
  modelValue: Boolean
})

const emit = defineEmits(['update:modelValue'])

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const viewMode = ref('month')
const currentDate = ref(new Date())
const selectedWorkerId = ref(null)
const loading = ref(false)
const appointments = ref([])
const weekDays = ref([])
const workerList = ref([])

const detailVisible = ref(false)
const selectedAppointmentId = ref(null)

const canCreate = ref(false) // TODO: 根据权限设置

// 人力状态相关
const totalWorkers = ref(0) // 总作业人员数量
const dailyWorkload = ref({}) // 每日工作量统计

const currentDateRange = computed(() => {
  const date = currentDate.value
  if (viewMode.value === 'day') {
    return formatDayHeader(date)
  } else if (viewMode.value === 'week') {
    // 计算周范围
    const start = new Date(date)
    const day = start.getDay()
    start.setDate(start.getDate() - (day === 0 ? 6 : day - 1))
    const end = new Date(start)
    end.setDate(end.getDate() + 6)
    return `${start.toLocaleDateString('zh-CN')} - ${end.toLocaleDateString('zh-CN')}`
  } else if (viewMode.value === 'month') {
    return `${date.getFullYear()}年 ${date.getMonth() + 1}月`
  }
  return ''
})

// 日视图时间段数据
const dayTimeSlots = computed(() => {
  const dateStr = formatDateToLocal(new Date(currentDate.value))
  const timeSlots = [
    { time_slot: 'morning', label: '上午' },
    { time_slot: 'noon', label: '中午' },
    { time_slot: 'afternoon', label: '下午' }
  ]

  return timeSlots.map(slot => {
    const slotAppointments = appointments.value.filter(a => {
      const aptDate = a.work_date.split(' ')[0] || a.work_date.split('T')[0]
      // 全天任务(full_day)在上午、中午、下午三个时间段都显示
      if (a.time_slot === 'full_day') {
        return aptDate === dateStr
      }
      return aptDate === dateStr && a.time_slot === slot.time_slot
    })

    return {
      ...slot,
      appointments: slotAppointments
    }
  })
})

// 格式化日期头部显示
function formatDayHeader(date) {
  const d = new Date(date)
  const weekDays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
  return `${d.getFullYear()}年${d.getMonth() + 1}月${d.getDate()}日 ${weekDays[d.getDay()]}`
}

// 获取某天的统计信息
function getDayStats(date) {
  const dateStr = formatDateToLocal(new Date(date))
  const dayAppointments = appointments.value.filter(a => {
    const aptDate = a.work_date.split(' ')[0] || a.work_date.split('T')[0]
    return aptDate === dateStr
  })
  return `共 ${dayAppointments.length} 个预约任务`
}

watch(() => props.modelValue, (val) => {
  if (val) {
    loadData()
  }
})

watch(viewMode, () => {
  loadData()
})

// 获取作业人员总数
async function fetchTotalWorkers() {
  try {
    const response = await appointmentApi.getWorkersList()
    const workers = response.data || []
    totalWorkers.value = workers.length
  } catch (error) {
    console.error('获取作业人员总数失败:', error)
    totalWorkers.value = 10 // 设置默认值
  }
}

async function loadData() {
  loading.value = true
  try {
    // 计算日期范围
    const { startDate, endDate } = getDateRange()

    const params = {
      start_date: startDate,
      end_date: endDate,
      page: 1,
      page_size: 1000 // 获取更多数据，避免分页问题
    }

    if (selectedWorkerId.value) {
      params.worker_id = selectedWorkerId.value
    }

    const response = await appointmentApi.getList(params)
    // 处理不同的响应格式
    let data = response.data
    if (data && data.data && Array.isArray(data.data)) {
      appointments.value = data.data
    } else if (data && Array.isArray(data)) {
      appointments.value = data
    } else if (response.meta && response.meta.data && Array.isArray(response.meta.data)) {
      appointments.value = response.meta.data
    } else {
      appointments.value = []
    }

    // 计算工作量
    calculateDailyWorkload()

    if (viewMode.value === 'week') {
      buildWeekData()
    }
  } catch (error) {
    console.error('加载数据失败:', error)
    ElMessage.error('加载数据失败')
    appointments.value = []
  } finally {
    loading.value = false
  }
}

// 格式化日期为 YYYY-MM-DD 格式（避免时区问题）
function formatDateToLocal(date) {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

function getDateRange() {
  const date = currentDate.value
  let startDate, endDate

  if (viewMode.value === 'day') {
    startDate = formatDateToLocal(new Date(date))
    endDate = startDate
  } else if (viewMode.value === 'week') {
    const start = new Date(date)
    const day = start.getDay()
    start.setDate(start.getDate() - (day === 0 ? 6 : day - 1))
    startDate = formatDateToLocal(start)

    const end = new Date(start)
    end.setDate(end.getDate() + 6)
    endDate = formatDateToLocal(end)
  } else if (viewMode.value === 'month') {
    // 当月第一天
    const firstDay = new Date(date.getFullYear(), date.getMonth(), 1)
    startDate = formatDateToLocal(firstDay)
    // 当月最后一天（下个月的第0天）
    const lastDay = new Date(date.getFullYear(), date.getMonth() + 1, 0)
    endDate = formatDateToLocal(lastDay)
  }

  return { startDate, endDate }
}

function buildWeekData() {
  const date = currentDate.value
  const start = new Date(date)
  const day = start.getDay()
  start.setDate(start.getDate() - (day === 0 ? 6 : day - 1))

  const days = []
  const dayNames = ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
  const timeSlots = [
    { time_slot: 'morning', label: '上午' },
    { time_slot: 'noon', label: '中午' },
    { time_slot: 'afternoon', label: '下午' }
  ]

  for (let i = 0; i < 7; i++) {
    const d = new Date(start)
    d.setDate(d.getDate() + i)
    const dateStr = formatDateToLocal(d)

    const slots = timeSlots.map(slot => {
      const slotAppointments = appointments.value.filter(a => {
        const aptDate = a.work_date.split(' ')[0]
        // 全天任务(full_day)在上午、中午、下午三个时间段都显示
        if (a.time_slot === 'full_day') {
          return aptDate === dateStr
        }
        return aptDate === dateStr && a.time_slot === slot.time_slot
      })

      return {
        ...slot,
        appointments: slotAppointments // 改为数组，支持多个预约
      }
    })

    days.push({
      dayName: dayNames[i],
      date: dateStr,
      slots
    })
  }

  weekDays.value = days
}

function getAppointmentsForDay(day) {
  return appointments.value.filter(apt => {
    // 处理不同格式的日期
    let aptDate = apt.work_date
    if (aptDate.includes(' ')) {
      aptDate = aptDate.split(' ')[0]
    } else if (aptDate.includes('T')) {
      aptDate = aptDate.split('T')[0]
    }
    return aptDate === day
  })
}

function handleViewChange() {
  currentDate.value = new Date()
  loadData()
}

function previousPeriod() {
  const date = new Date(currentDate.value)
  if (viewMode.value === 'day') {
    date.setDate(date.getDate() - 1)
  } else if (viewMode.value === 'week') {
    date.setDate(date.getDate() - 7)
  } else if (viewMode.value === 'month') {
    date.setMonth(date.getMonth() - 1)
  }
  currentDate.value = date
  loadData()
}

function nextPeriod() {
  const date = new Date(currentDate.value)
  if (viewMode.value === 'day') {
    date.setDate(date.getDate() + 1)
  } else if (viewMode.value === 'week') {
    date.setDate(date.getDate() + 7)
  } else if (viewMode.value === 'month') {
    date.setMonth(date.getMonth() + 1)
  }
  currentDate.value = date
  loadData()
}

function goToToday() {
  currentDate.value = new Date()
  loadData()
}

function showAppointment(apt) {
  selectedAppointmentId.value = apt.id
  detailVisible.value = true
}

function createAppointment(date, timeSlot) {
  ElMessage.info('创建预约功能开发中')
}

function handleDayClick(data) {
  // 月视图点击日期的处理
  console.log('Clicked day:', data.day)
}

function getTimeSlotLabel(timeSlot) {
  return appointmentApi.getTimeSlotLabel(timeSlot)
}

function getStatusLabel(status) {
  return appointmentApi.getStatusLabel(status)
}

function getStatusType(status) {
  return appointmentApi.getStatusType(status)
}

function getTimeSlotShortLabel(timeSlot) {
  const labels = {
    morning: '上午',
    noon: '中午',
    afternoon: '下午',
    full_day: '全天'
  }
  return labels[timeSlot] || timeSlot
}

function getDayClass(day) {
  const dayAppointments = getAppointmentsForDay(day)

  if (dayAppointments.length === 0) {
    return ''
  }

  // 检查是否有加急任务
  const hasUrgent = dayAppointments.some(apt => apt.is_urgent)
  if (hasUrgent) {
    return 'day-has-urgent'
  }

  // 检查状态优先级: pending > in_progress > scheduled > completed > other
  const statusPriority = ['pending', 'in_progress', 'scheduled', 'completed']
  for (const status of statusPriority) {
    if (dayAppointments.some(apt => apt.status === status)) {
      return `day-has-${status}`
    }
  }

  return 'day-has-appointments'
}

// 计算每日工作量
function calculateDailyWorkload() {
  const workload = {}

  appointments.value.forEach(apt => {
    const date = apt.work_date.split(' ')[0] || apt.work_date.split('T')[0]

    if (!workload[date]) {
      workload[date] = {
        total: 0,
        assigned: 0,
        unassigned: 0
      }
    }

    workload[date].total++

    if (apt.assigned_worker_id) {
      workload[date].assigned++
    } else {
      workload[date].unassigned++
    }
  })

  dailyWorkload.value = workload
}

// 获取某天的人力状态类名
function getWorkloadClass(day) {
  const info = getDayWorkloadInfo(day)
  if (!info) return ''

  if (info.ratio >= 1) return 'workload-overload'
  if (info.ratio >= 0.75) return 'workload-busy'
  if (info.ratio >= 0.3) return 'workload-normal'
  return 'workload-available'
}

// 获取某天的人力状态信息
function getDayWorkloadInfo(day) {
  if (!dailyWorkload.value[day] || totalWorkers.value === 0) {
    return null
  }

  const workload = dailyWorkload.value[day]
  const ratio = workload.total / totalWorkers.value
  let status = ''
  let text = `${workload.total}/${totalWorkers.value}任务`

  if (ratio >= 1) {
    status = '超载'
  } else if (ratio >= 0.75) {
    status = '紧张'
  } else if (ratio >= 0.3) {
    status = '正常'
  } else {
    status = '空闲'
  }

  return {
    ratio,
    status,
    text: `${text}`,
    total: workload.total,
    workers: totalWorkers.value
  }
}

function handleClose() {
  emit('update:modelValue', false)
}

// 初始化
onMounted(() => {
  fetchTotalWorkers()
})
</script>

<style scoped>
.calendar-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid #eee;
}

.date-navigator {
  display: flex;
  align-items: center;
  gap: 8px;
}

.current-date {
  padding: 0 16px;
  font-weight: 500;
  color: #333;
}

/* 日视图样式 */
.day-view {
  padding: 20px;
}

.day-header-large {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 16px 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  color: white;
}

.day-header-large h3 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
}

.day-stats {
  font-size: 14px;
  opacity: 0.9;
}

.day-slots {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.day-slot-card {
  background: #fff;
  border: 1px solid #ebeef5;
  border-radius: 12px;
  overflow: hidden;
  transition: box-shadow 0.3s;
}

.day-slot-card.has-appointment {
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.slot-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #f5f7fa;
  border-bottom: 1px solid #ebeef5;
}

.slot-title {
  font-weight: 600;
  font-size: 15px;
  color: #303133;
}

.slot-count {
  font-size: 13px;
  color: #909399;
}

.slot-content {
  padding: 12px;
  min-height: 100px;
}

.slot-empty-large {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100px;
}

.appointment-card {
  padding: 14px;
  margin-bottom: 10px;
  background: #fafafa;
  border-radius: 8px;
  border: 1px solid #ebeef5;
  cursor: pointer;
  transition: all 0.2s;
}

.appointment-card:last-child {
  margin-bottom: 0;
}

.appointment-card:hover {
  background: #f0f7ff;
  border-color: #409eff;
}

.appointment-card.urgent {
  background: linear-gradient(135deg, #fff1f0 0%, #ffccc7 100%);
  border-color: #ffa39e;
}

.apt-header {
  display: flex;
  gap: 6px;
  margin-bottom: 10px;
  flex-wrap: wrap;
}

.apt-content {
  margin-bottom: 10px;
}

.apt-location {
  font-size: 15px;
  font-weight: 500;
  color: #303133;
  margin-bottom: 4px;
}

.apt-detail {
  font-size: 13px;
  color: #606266;
  line-height: 1.5;
}

.apt-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: #909399;
}

.apt-worker {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #409eff;
}

.apt-no {
  color: #c0c4cc;
}

.week-view {
  padding: 10px;
}

.day-column {
  border: 1px solid #eee;
  border-radius: 4px;
  overflow: hidden;
}

.day-header {
  display: flex;
  justify-content: space-between;
  padding: 12px;
  background: #f5f7fa;
  border-bottom: 1px solid #eee;
  font-weight: 500;
}

.day-name {
  color: #333;
}

.day-date {
  color: #666;
  font-size: 12px;
}

.day-content {
  /* 使用 flexbox 使所有时间段均匀拉伸 */
  display: flex;
  flex-direction: column;
  /* 设置一个合理的高度，让内容自适应 */
  height: 500px;
  overflow-y: auto;
}

.time-slot {
  border-bottom: 1px solid #eee;
  /* 使所有时间段均匀拉伸 */
  flex: 1;
  /* 保持最小高度以确保内容可读 */
  min-height: 0;
}

.time-slot:last-child {
  border-bottom: none;
}

.slot-label {
  padding: 8px 12px;
  font-size: 12px;
  color: #999;
  background: #fafafa;
}

.slot-empty {
  padding: 16px 12px;
  text-align: center;
  color: #ddd;
}

.slot-empty-text {
  color: #c0c4cc;
  font-size: 13px;
}

.appointment-item-simple {
  padding: 10px 12px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  transition: background 0.2s;
}

.appointment-item-simple:hover {
  background: #f5f7fa;
}

.appointment-summary {
  font-size: 13px;
  color: #606266;
  flex: 1;
}

.appointment-worker-small {
  font-size: 12px;
  color: #909399;
}

.month-view {
  padding: 20px;
}

.calendar-day {
  min-height: 80px;
  padding: 4px;
}

.day-number {
  text-align: center;
  font-weight: 500;
  color: #333;
  margin-bottom: 4px;
}

.day-appointments {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.day-appointment-item {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  background: #f5f7fa;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
}

.day-appointment-item:hover {
  background: #e6f7ff;
}

.appointment-summary {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: #666;
}

.appointment-worker-small {
  font-size: 11px;
  color: #409eff;
  padding: 0 4px;
  border-radius: 2px;
  background: #ecf5ff;
}

.text-muted {
  color: #999;
  font-style: italic;
}

/* 日历视图天样式 - 根据任务状态显示不同背景色 */
.month-view :deep(.el-calendar-day) {
  height: 100%;
  min-height: 85px;
  padding: 0;
}

.month-view :deep(.el-calendar-table__td) {
  border: 1px solid #ebeef5;
  padding: 0;
}

.month-view :deep(.el-calendar-table) {
  border: none;
}

/* 隐藏 el-calendar 自带的英文头部 */
.month-view :deep(.el-calendar__header) {
  display: none;
}

.month-view :deep(.el-calendar-table td.is-selected) {
  background: transparent;
}

/* 我们的 calendar-day div 需要占满整个单元格 */
.month-view :deep(.el-calendar-day .calendar-day) {
  width: 100%;
  height: 100%;
  min-height: 80px;
  padding: 4px;
  box-sizing: border-box;
}

.month-view :deep(.el-calendar-day .calendar-day.day-has-urgent) {
  background: linear-gradient(135deg, #fff1f0 0%, #ffccc7 100%) !important;
  border-radius: 4px;
  border: 1px solid #ffa39e;
}

.month-view :deep(.el-calendar-day .calendar-day.day-has-pending) {
  background: linear-gradient(135deg, #fffbe6 0%, #ffe58f 100%) !important;
  border-radius: 4px;
  border: 1px solid #ffd666;
}

.month-view :deep(.el-calendar-day .calendar-day.day-has-in_progress) {
  background: linear-gradient(135deg, #e6f7ff 0%, #bae7ff 100%) !important;
  border-radius: 4px;
  border: 1px solid #91d5ff;
}

.month-view :deep(.el-calendar-day .calendar-day.day-has-scheduled) {
  background: linear-gradient(135deg, #f6ffed 0%, #d9f7be 100%) !important;
  border-radius: 4px;
  border: 1px solid #b7eb8f;
}

.month-view :deep(.el-calendar-day .calendar-day.day-has-completed) {
  background: linear-gradient(135deg, #f5f5f5 0%, #e8e8e8 100%) !important;
  border-radius: 4px;
  border: 1px solid #d9d9d9;
}

.month-view :deep(.el-calendar-day .calendar-day.day-has-appointments) {
  background: linear-gradient(135deg, #f9f9f9 0%, #f0f0f0 100%) !important;
  border-radius: 4px;
  border: 1px solid #e0e0e0;
}

/* 人力状态指示器 */
.legend-section {
  padding: 16px;
  margin-bottom: 16px;
  background: #fafafa;
  border-radius: 8px;
  border: 1px solid #ebeef5;
}

.legend-group {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  margin-bottom: 12px;
}

.legend-group:last-child {
  margin-bottom: 0;
}

.legend-title {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
  margin-right: 12px;
  min-width: 80px;
}

.legend-items {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #606266;
}

.legend-box {
  width: 20px;
  height: 20px;
  border-radius: 4px;
  border: 1px solid #dcdfe6;
  flex-shrink: 0;
}

.legend-box.legend-available {
  background: #67c23a;
}

.legend-box.legend-normal {
  background: #409eff;
}

.legend-box.legend-busy {
  background: #e6a23c;
}

.legend-box.legend-overload {
  background: #f56c6c;
}

.month-view :deep(.day-number) {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.workload-indicator {
  font-size: 10px;
  padding: 2px 4px;
  border-radius: 3px;
  font-weight: 500;
}

.month-view :deep(.workload-indicator.workload-available) {
  background: #67c23a;
  color: white;
}

.month-view :deep(.workload-indicator.workload-normal) {
  background: #409eff;
  color: white;
}

.month-view :deep(.workload-indicator.workload-busy) {
  background: #e6a23c;
  color: white;
}

.month-view :deep(.workload-indicator.workload-overload) {
  background: #f56c6c;
  color: white;
}
</style>
