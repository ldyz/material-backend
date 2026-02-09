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
        <el-radio-button label="week">周视图</el-radio-button>
        <el-radio-button label="month">月视图</el-radio-button>
        <el-radio-button label="season">季度视图</el-radio-button>
      </el-radio-group>

      <div class="date-navigator">
        <el-button :icon="ArrowLeft" @click="previousPeriod">上一段</el-button>
        <span class="current-date">{{ currentDateRange }}</span>
        <el-button @click="nextPeriod">下一段 <el-icon><ArrowRight /></el-icon></el-button>
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
                :class="{ 'has-appointment': slot.appointment }"
              >
                <div class="slot-label">{{ slot.label }}</div>
                <div v-if="slot.appointment" class="appointment-item" @click="showAppointment(slot.appointment)">
                  <div class="appointment-no">{{ slot.appointment.appointment_no }}</div>
                  <div class="appointment-location">{{ slot.appointment.work_location }}</div>
                  <div class="appointment-content">{{ slot.appointment.work_content }}</div>
                  <el-tag
                    v-if="slot.appointment.is_urgent"
                    type="danger"
                    size="small"
                    class="appointment-urgent"
                  >加急</el-tag>
                </div>
                <div v-else class="slot-empty">
                  <el-button
                    v-if="canCreate"
                    type="primary"
                    size="small"
                    link
                    @click="createAppointment(day.date, slot.time_slot)"
                  >
                    + 预约
                  </el-button>
                </div>
              </div>
            </div>
          </el-col>
        </el-row>
      </div>

      <!-- 月视图 -->
      <div v-else-if="viewMode === 'month'" class="month-view">
        <el-calendar v-model="currentDate">
          <template #date-cell="{ data }">
            <div class="calendar-day" @click="handleDayClick(data)">
              <div class="day-number">{{ data.day.split('-').slice(-1)[0] }}</div>
              <div class="day-appointments">
                <div
                  v-for="apt in getAppointmentsForDay(data.day)"
                  :key="apt.id"
                  class="day-appointment-item"
                  @click.stop="showAppointment(apt)"
                >
                  <el-tag :type="getStatusColor(apt.status)" size="small">
                    {{ apt.time_slot }}
                  </el-tag>
                  <span class="appointment-summary">{{ apt.work_location }}</span>
                  <el-tag v-if="apt.is_urgent" type="danger" size="small">急</el-tag>
                </div>
              </div>
            </div>
          </template>
        </el-calendar>
      </div>

      <!-- 季度视图 - 统计列表 -->
      <div v-else class="season-view">
        <el-table :data="seasonData" border style="width: 100%">
          <el-table-column prop="date" label="日期" width="120" />
          <el-table-column prop="time_slot" label="时间段" width="100">
            <template #default="scope">
              {{ getTimeSlotLabel(scope.row.time_slot) }}
            </template>
          </el-table-column>
          <el-table-column prop="appointment_no" label="预约单号" width="140" />
          <el-table-column prop="work_location" label="作业地点" />
          <el-table-column prop="work_content" label="作业内容" show-overflow-tooltip />
          <el-table-column label="状态" width="90">
            <template #default="scope">
              <el-tag :type="getStatusType(scope.row.status)" size="small">
                {{ getStatusLabel(scope.row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="80">
            <template #default="scope">
              <el-button type="primary" link @click="showAppointment(scope.row)">
                查看
              </el-button>
            </template>
          </el-table-column>
        </el-table>
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
import { ArrowLeft, ArrowRight } from '@element-plus/icons-vue'
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
const seasonData = ref([])
const workerList = ref([])

const detailVisible = ref(false)
const selectedAppointmentId = ref(null)

const canCreate = ref(false) // TODO: 根据权限设置

const currentDateRange = computed(() => {
  const date = currentDate.value
  if (viewMode.value === 'week') {
    // 计算周范围
    const start = new Date(date)
    const day = start.getDay()
    start.setDate(start.getDate() - (day === 0 ? 6 : day - 1))
    const end = new Date(start)
    end.setDate(end.getDate() + 6)
    return `${start.toLocaleDateString('zh-CN')} - ${end.toLocaleDateString('zh-CN')}`
  } else if (viewMode.value === 'month') {
    return `${date.getFullYear()}年 ${date.getMonth() + 1}月`
  } else {
    // 季度
    const quarter = Math.floor(date.getMonth() / 3)
    const start = new Date(date.getFullYear(), quarter * 3, 1)
    const end = new Date(date.getFullYear(), (quarter + 1) * 3, 0)
    return `${start.toLocaleDateString('zh-CN')} - ${end.toLocaleDateString('zh-CN')}`
  }
})

watch(() => props.modelValue, (val) => {
  if (val) {
    loadData()
  }
})

watch(viewMode, () => {
  loadData()
})

async function loadData() {
  loading.value = true
  try {
    // 计算日期范围
    const { startDate, endDate } = getDateRange()

    const params = {
      start_date: startDate,
      end_date: endDate
    }

    if (selectedWorkerId.value) {
      params.worker_id = selectedWorkerId.value
    }

    const { data } = await appointmentApi.getList(params)
    appointments.value = data.data || []

    if (viewMode.value === 'week') {
      buildWeekData()
    } else if (viewMode.value === 'season') {
      buildSeasonData()
    }
  } catch (error) {
    ElMessage.error('加载数据失败')
  } finally {
    loading.value = false
  }
}

function getDateRange() {
  const date = currentDate.value
  let startDate, endDate

  if (viewMode.value === 'week') {
    const start = new Date(date)
    const day = start.getDay()
    start.setDate(start.getDate() - (day === 0 ? 6 : day - 1))
    startDate = start.toISOString().split('T')[0]

    const end = new Date(start)
    end.setDate(end.getDate() + 6)
    endDate = end.toISOString().split('T')[0]
  } else if (viewMode.value === 'month') {
    startDate = new Date(date.getFullYear(), date.getMonth(), 1).toISOString().split('T')[0]
    endDate = new Date(date.getFullYear(), date.getMonth() + 1, 0).toISOString().split('T')[0]
  } else {
    // 季度
    const quarter = Math.floor(date.getMonth() / 3)
    startDate = new Date(date.getFullYear(), quarter * 3, 1).toISOString().split('T')[0]
    endDate = new Date(date.getFullYear(), (quarter + 1) * 3, 0).toISOString().split('T')[0]
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
    { time_slot: 'afternoon', label: '下午' },
    { time_slot: 'evening', label: '晚上' }
  ]

  for (let i = 0; i < 7; i++) {
    const d = new Date(start)
    d.setDate(d.getDate() + i)
    const dateStr = d.toISOString().split('T')[0]

    const slots = timeSlots.map(slot => {
      const apt = appointments.value.find(a => {
        const aptDate = a.work_date.split(' ')[0]
        return aptDate === dateStr && a.time_slot === slot.time_slot
      })

      return {
        ...slot,
        appointment: apt || null
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

function buildSeasonData() {
  // 按日期排序
  seasonData.value = [...appointments.value].sort((a, b) => {
    const dateA = new Date(a.work_date)
    const dateB = new Date(b.work_date)
    return dateA - dateB
  })
}

function getAppointmentsForDay(day) {
  return appointments.value.filter(apt => {
    const aptDate = apt.work_date.split(' ')[0]
    return aptDate === day
  })
}

function handleViewChange() {
  currentDate.value = new Date()
  loadData()
}

function previousPeriod() {
  const date = new Date(currentDate.value)
  if (viewMode.value === 'week') {
    date.setDate(date.getDate() - 7)
  } else if (viewMode.value === 'month') {
    date.setMonth(date.getMonth() - 1)
  } else {
    date.setMonth(date.getMonth() - 3)
  }
  currentDate.value = date
  loadData()
}

function nextPeriod() {
  const date = new Date(currentDate.value)
  if (viewMode.value === 'week') {
    date.setDate(date.getDate() + 7)
  } else if (viewMode.value === 'month') {
    date.setMonth(date.getMonth() + 1)
  } else {
    date.setMonth(date.getMonth() + 3)
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

function handleClose() {
  emit('update:modelValue', false)
}
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
  max-height: 500px;
  overflow-y: auto;
}

.time-slot {
  border-bottom: 1px solid #eee;
  min-height: 80px;
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

.appointment-item {
  padding: 12px;
  cursor: pointer;
  transition: background 0.2s;
}

.appointment-item:hover {
  background: #f5f7fa;
}

.appointment-no {
  font-weight: 500;
  color: #333;
  margin-bottom: 4px;
}

.appointment-location {
  font-size: 13px;
  color: #666;
  margin-bottom: 4px;
}

.appointment-content {
  font-size: 12px;
  color: #999;
  margin-bottom: 4px;
}

.appointment-urgent {
  margin-top: 4px;
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

.season-view {
  padding: 10px;
}
</style>
