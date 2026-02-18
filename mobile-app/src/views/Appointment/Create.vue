<template>
  <div class="appointment-create">
    <van-nav-bar
      title="创建预约单"
      left-arrow
      @click-left="handleBack"
    />

    <van-form @submit="handleSubmit">
      <!-- 基本信息 -->
      <van-cell-group title="基本信息" inset>
        <ProjectSelector
          v-model="form.project_id"
          name="project_id"
          label="项目"
          placeholder="请选择项目"
          :required="true"
        />
        <van-field
          v-model="form.work_location"
          name="work_location"
          label="作业地点"
          placeholder="请输入作业地点"
          :rules="[{ required: true, message: '请输入作业地点' }]"
        />
        <van-field
          v-model="form.work_content"
          name="work_content"
          label="作业内容"
          type="textarea"
          placeholder="请输入作业内容"
          rows="3"
          :rules="[{ required: true, message: '请输入作业内容' }]"
        />
        <van-field
          v-model="form.work_type"
          name="work_type"
          label="作业类型"
          readonly
          is-link
          placeholder="请选择作业类型"
          @click="showWorkTypePicker = true"
          :rules="[{ required: true, message: '请选择作业类型' }]"
        >
          <template #input>
            <span :style="{ color: workTypeLabel ? '#323233' : '#c8c9cc' }">
              {{ workTypeLabel || '请选择作业类型' }}
            </span>
          </template>
        </van-field>
      </van-cell-group>

      <!-- 时间信息 -->
      <van-cell-group title="时间信息" inset>
        <van-field
          v-model="form.work_date"
          name="work_date"
          label="作业日期"
          readonly
          is-link
          placeholder="请选择作业日期"
          @click="showDatePicker = true"
          :rules="[{ required: true, message: '请选择作业日期' }]"
        >
          <template #input>
            <span :style="{ color: form.work_date ? '#323233' : '#c8c9cc' }">
              {{ form.work_date || '请选择作业日期' }}
            </span>
          </template>
        </van-field>
        <!-- 人员可用性提示 -->
        <van-cell
          v-if="form.work_date && !hasAvailableWorkers && !form.is_urgent"
          title="提示"
          icon="warning-o"
        >
          <template #value>
            <span style="color: #ee0a24">该日期所有人员已被安排，请创建加急预约单</span>
          </template>
        </van-cell>
        <van-field
          v-model="form.time_slot"
          name="time_slot"
          label="时间段"
          readonly
          is-link
          placeholder="请选择时间段"
          @click="showTimeSlotPicker = true"
          :rules="[{ required: true, message: '请选择时间段' }]"
        >
          <template #input>
            <span :style="{ color: timeSlotLabel ? '#323233' : '#c8c9cc' }">
              {{ timeSlotLabel || '请选择时间段' }}
            </span>
          </template>
        </van-field>
      </van-cell-group>

      <!-- 联系信息 -->
      <van-cell-group title="联系信息" inset>
        <van-field
          v-model="form.contact_person"
          name="contact_person"
          label="联系人"
          placeholder="请输入联系人"
        />
        <van-field
          v-model="form.contact_phone"
          name="contact_phone"
          label="联系电话"
          type="tel"
          placeholder="请输入联系电话"
        />
      </van-cell-group>

      <!-- 优先级 -->
      <van-cell-group title="优先级" inset>
        <van-field name="is_urgent" label="是否加急">
          <template #input>
            <van-switch v-model="form.is_urgent" />
          </template>
        </van-field>
        <van-field
          v-if="form.is_urgent"
          v-model="form.priority"
          name="priority"
          label="优先级"
          type="number"
          placeholder="0-10"
          :rules="[
            { required: true, message: '请输入优先级' },
            { pattern: /^\d+$/, message: '请输入数字' },
            { validator: (val) => val >= 0 && val <= 10, message: '优先级必须在0-10之间' }
          ]"
        />
        <van-field
          v-if="form.is_urgent && form.priority >= 7"
          v-model="form.urgent_reason"
          name="urgent_reason"
          label="加急原因"
          type="textarea"
          placeholder="请输入加急原因"
          rows="2"
          :rules="[{ required: true, message: '高优先级加急必须提供原因' }]"
        />
      </van-cell-group>

      <!-- 作业人员 -->
      <van-cell-group title="作业人员（可选，可多选）" inset>
        <van-field
          name="assigned_workers"
          label="作业人员"
          :value="selectedWorkersNames"
          readonly
          is-link
          placeholder="选择作业人员"
          @click="showWorkerPicker = true"
        />
        <!-- 已选择的作业人员 -->
        <div v-if="selectedWorkers.length > 0" class="selected-workers-display">
          <div
            v-for="workerId in selectedWorkers"
            :key="workerId"
            class="selected-worker-item"
          >
            <van-image
              v-if="getWorkerById(workerId)?.avatar"
              round
              width="40"
              height="40"
              :src="getAssetUrl(getWorkerById(workerId).avatar)"
            />
            <div
              v-else
              class="worker-avatar-small"
              :style="{ backgroundColor: getAvatarColor(workerId) }"
            >
              {{ (getWorkerById(workerId)?.full_name || getWorkerById(workerId)?.username || '?').charAt(0) }}
            </div>
            <span class="selected-worker-name">{{ getWorkerById(workerId)?.full_name || getWorkerById(workerId)?.username }}</span>
            <van-icon
              name="cross"
              class="remove-worker-icon"
              @click.stop="removeWorker(workerId)"
            />
          </div>
        </div>
      </van-cell-group>

      <div class="submit-bar">
        <van-button round block type="primary" native-type="submit" :loading="submitting">
          提交
        </van-button>
      </div>
    </van-form>

    <!-- 日期选择器 -->
    <van-popup v-model:show="showDatePicker" position="bottom">
      <van-date-picker
        v-model="currentDate"
        title="选择日期"
        :min-date="minDate"
        :max-date="maxDate"
        @confirm="onDateConfirm"
        @cancel="onDatePickerClose"
      />
    </van-popup>

    <!-- 时间段选择器 -->
    <van-popup v-model:show="showTimeSlotPicker" position="bottom" :style="{ height: '50%' }" round>
      <div class="time-slot-picker">
        <van-nav-bar
          title="选择时间段"
          left-text="取消"
          @click-left="showTimeSlotPicker = false"
        />
        <van-loading v-if="loadingTimeSlots" type="spinner" size="24" vertical>加载中...</van-loading>
        <van-cell-group v-else inset>
          <van-cell
            v-for="option in timeSlotOptions"
            :key="option.value"
            :title="option.text"
            is-link
            :class="{
              'time-slot-full': getTimeSlotStatus(option.value) === 'full',
              'time-slot-busy': getTimeSlotStatus(option.value) === 'busy',
              'time-slot-moderate': getTimeSlotStatus(option.value) === 'moderate',
              'time-slot-available': getTimeSlotStatus(option.value) === 'available'
            }"
            @click="selectTimeSlot(option.value)"
          >
            <template #right-icon>
              <van-tag :type="getTimeSlotStatusTagType(getTimeSlotStatus(option.value))" round>
                {{ getTimeSlotStatusText(option.value) || '请选择日期' }}
              </van-tag>
            </template>
          </van-cell>
        </van-cell-group>
      </div>
    </van-popup>

    <!-- 作业类型选择器 -->
    <van-popup v-model:show="showWorkTypePicker" position="bottom" round>
      <div class="work-type-picker">
        <van-nav-bar
          title="选择作业类型"
          left-text="取消"
          right-text="确认"
          @click-left="showWorkTypePicker = false"
          @click-right="confirmWorkType"
        />
        <van-checkbox-group v-model="selectedWorkTypes">
          <van-cell-group>
            <van-cell
              v-for="type in workTypeOptions"
              :key="type.value"
              clickable
              :title="type.label"
              @click="toggleWorkType(type.value)"
            >
              <template #right-icon>
                <van-checkbox :name="type.value" ref="checkboxes" @click.stop />
              </template>
            </van-cell>
          </van-cell-group>
        </van-checkbox-group>
      </div>
    </van-popup>

    <!-- 作业人员选择器 -->
    <van-popup v-model:show="showWorkerPicker" position="bottom" :style="{ height: '70%' }" round>
      <div class="worker-picker">
        <van-nav-bar
          title="选择作业人员"
          left-text="取消"
          right-text="确认"
          @click-left="showWorkerPicker = false"
          @click-right="confirmWorkerSelection"
        />
        <div class="worker-selection-info" v-if="selectedWorkers.length > 0">
          <span class="selected-count">已选择 {{ selectedWorkers.length }} 人</span>
        </div>
        <van-checkbox-group v-model="selectedWorkers">
          <div class="worker-grid">
            <div
              v-for="worker in workerList"
              :key="worker.id"
              :class="['worker-card', { 'selected': isWorkerSelected(worker.id) }]"
              @click="toggleWorkerSelection(worker.id)"
            >
              <div class="worker-avatar">
                <van-image
                  v-if="worker.avatar"
                  round
                  width="50"
                  height="50"
                  :src="getAssetUrl(worker.avatar)"
                />
                <div
                  v-else
                  class="worker-avatar-placeholder"
                  :style="{ backgroundColor: getAvatarColor(worker.id) }"
                >
                  {{ (worker.full_name || worker.username || '?').charAt(0) }}
                </div>
                <div v-if="isWorkerSelected(worker.id)" class="selected-badge">
                  <van-icon name="success" />
                </div>
              </div>
              <div class="worker-name">{{ worker.full_name || worker.username }}</div>
              <van-checkbox :name="worker.id" ref="workerCheckboxes" @click.stop />
            </div>
            <div v-if="workerList.length === 0" class="empty-workers">
              <van-empty description="暂无作业人员" />
            </div>
          </div>
        </van-checkbox-group>
      </div>
    </van-popup>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { showSuccessToast, showFailToast, showConfirmDialog } from 'vant'
import { createAppointment, submitAppointment, getTimeSlotOptions, getDailyStatistics, getWorkersList, getTimeSlotStatistics, getAvailableWorkers } from '@/api/appointment'
import { getAssetUrl } from '@/utils/request'
import ProjectSelector from '@/components/common/ProjectSelector.vue'

const router = useRouter()

const form = ref({
  project_id: null,
  work_location: '',
  work_content: '',
  work_type: '',
  work_date: '',
  time_slot: '',
  contact_person: '',
  contact_phone: '',
  is_urgent: false,
  priority: 0,
  urgent_reason: ''
})

const submitting = ref(false)
const showDatePicker = ref(false)
const showTimeSlotPicker = ref(false)
const showWorkerPicker = ref(false)
const showWorkTypePicker = ref(false)
const currentDate = ref([new Date().getFullYear(), new Date().getMonth() + 1, new Date().getDate()])

// 计算最小可选日期（明天）
const minDate = computed(() => {
  const tomorrow = new Date()
  tomorrow.setDate(tomorrow.getDate() + 1)
  tomorrow.setHours(0, 0, 0, 0)
  return tomorrow
})

// 计算最大可选日期（加急可以选择之后1天，普通无限制）
const maxDate = computed(() => {
  // 不设置最大日期限制，让用户可以预约未来的任意日期
  return undefined
})

// 该日期是否有空闲人员
const hasAvailableWorkers = ref(true)
const selectedWorkTypes = ref([])
const selectedWorkers = ref([]) // 选中的作业人员ID列表
const workerList = ref([]) // 作业人员列表

// 每日预约数据
const dailyAppointments = ref({})
const totalWorkers = ref(0)

// 时间段统计数据
const timeSlotStatistics = ref([])
const loadingTimeSlots = ref(false)

// 检查日期是否有空闲人员
async function checkAvailabilityForDate(date) {
  try {
    // 检查上午时段的可用性作为参考
    const response = await getAvailableWorkers({
      work_date: date,
      time_slot: 'morning'
    })
    const workers = response.data || []
    hasAvailableWorkers.value = workers.length > 0

    if (!hasAvailableWorkers.value && !form.value.is_urgent) {
      showFailToast('该日期所有人员已被安排，请选择其他日期或创建加急预约单')
    }
  } catch (error) {
    console.error('检查可用人员失败:', error)
    hasAvailableWorkers.value = true
  }
}

// 获取作业人员列表
async function fetchWorkers() {
  try {
    const response = await getWorkersList()
    workerList.value = response.data || []
  } catch (error) {
    console.error('获取作业人员列表失败:', error)
  }
}

// 获取时间段统计数据
async function fetchTimeSlotStatistics() {
  if (!form.value.work_date) return

  loadingTimeSlots.value = true
  try {
    const response = await getTimeSlotStatistics(form.value.work_date)
    timeSlotStatistics.value = response.data.statistics || []
    totalWorkers.value = response.data.total_workers || 0
  } catch (error) {
    console.error('获取时间段统计数据失败:', error)
    timeSlotStatistics.value = []
  } finally {
    loadingTimeSlots.value = false
  }
}

// 计算时间段的繁忙程度
function getTimeSlotStatus(timeSlot) {
  const stat = timeSlotStatistics.value.find(s => s.time_slot === timeSlot)
  if (!stat) return 'unknown'

  const totalWorkersCount = totalWorkers.value || 1
  const ratio = stat.total_count / totalWorkersCount

  if (ratio >= 1) return 'full'
  if (ratio >= 0.75) return 'busy'
  if (ratio >= 0.4) return 'moderate'
  return 'available'
}

// 获取时间段状态文本
function getTimeSlotStatusText(timeSlot) {
  const stat = timeSlotStatistics.value.find(s => s.time_slot === timeSlot)
  if (!stat) return ''

  const totalWorkersCount = totalWorkers.value || 1
  const ratio = stat.total_count / totalWorkersCount

  if (ratio >= 1) return '已满'
  if (ratio >= 0.75) return '繁忙'
  if (ratio >= 0.4) return '适中'
  return '空闲'
}

// 计算已选择作业人员的名称
const selectedWorkersNames = computed(() => {
  if (selectedWorkers.value.length === 0) return ''
  const names = selectedWorkers.value.map(id => {
    const worker = workerList.value.find(w => w.id === id)
    return worker ? (worker.full_name || worker.username) : ''
  }).filter(name => name)
  return names.join('、')
})

// 判断作业人员是否被选中
function isWorkerSelected(workerId) {
  return selectedWorkers.value.includes(workerId)
}

// 切换作业人员选择状态
function toggleWorkerSelection(workerId) {
  const index = selectedWorkers.value.indexOf(workerId)
  if (index !== -1) {
    selectedWorkers.value.splice(index, 1)
  } else {
    selectedWorkers.value.push(workerId)
  }
}

// 确认作业人员选择
function confirmWorkerSelection() {
  showWorkerPicker.value = false
}

// 根据ID获取作业人员
function getWorkerById(workerId) {
  return workerList.value.find(w => w.id === workerId)
}

// 移除已选择的作业人员
function removeWorker(workerId) {
  const index = selectedWorkers.value.indexOf(workerId)
  if (index !== -1) {
    selectedWorkers.value.splice(index, 1)
  }
}

// 组件挂载时获取作业人员列表
onMounted(() => {
  fetchWorkers()
})

const timeSlotOptions = getTimeSlotOptions().map(opt => ({
  text: opt.label,
  value: opt.value
}))

const workTypeOptions = [
  { value: 'general', label: '一般作业' },
  { value: 'hot_work', label: '动火作业' },
  { value: 'high_work', label: '高处作业' },
  { value: 'excavation', label: '动土作业' },
  { value: 'confined_space', label: '受限空间' },
  { value: 'electrical', label: '临时用电' },
  { value: 'lifting', label: '吊装作业' },
  { value: 'blind_plate', label: '盲板抽堵' }
]

// 日期显示标签
const dateLabel = computed(() => {
  if (!form.value.work_date) return ''
  return form.value.work_date
})

// 时间段标签
const timeSlotLabel = computed(() => {
  if (!form.value.time_slot) return ''
  const option = timeSlotOptions.find(opt => opt.value === form.value.time_slot)
  return option ? option.text : ''
})

// 作业类型标签
const workTypeLabel = computed(() => {
  if (!form.value.work_type) return ''
  const types = form.value.work_type.split(',').map(t => {
    const option = workTypeOptions.find(opt => opt.value === t)
    return option ? option.label : t
  })
  return types.join('、')
})

async function handleSubmit() {
  try {
    submitting.value = true

    // 验证项目是否已选择
    if (!form.value.project_id) {
      showFailToast('请选择项目')
      return
    }

    // 检查是否有空闲人员
    if (!form.value.is_urgent && !hasAvailableWorkers.value) {
      showFailToast('该日期所有人员已被安排，请选择其他日期或创建加急预约单')
      return
    }

    // 准备提交数据
    const submitData = { ...form.value }

    // 处理多选作业人员
    if (selectedWorkers.value.length > 0) {
      // 为了兼容性，第一个工人作为主作业人员
      submitData.assigned_worker_id = selectedWorkers.value[0]
      // 多选工人以 JSON 数组字符串格式发送
      submitData.assigned_worker_ids = JSON.stringify(selectedWorkers.value)
      // 工人姓名以逗号分隔格式发送
      const workerNames = selectedWorkers.value.map(id => {
        const worker = workerList.value.find(w => w.id === id)
        return worker ? (worker.full_name || worker.username) : ''
      }).filter(name => name)
      submitData.assigned_worker_names = workerNames.join(',')
    }

    console.log('Creating appointment with data:', submitData)
    const result = await createAppointment(submitData)
    console.log('Create appointment result:', result)

    // 创建成功后提交
    const appointmentId = result.data?.id || result.id
    if (appointmentId) {
      await submitAppointment(appointmentId)
    }

    showSuccessToast('提交成功')
    router.replace('/appointments')
  } catch (error) {
    console.error('Create appointment error:', error)
    showFailToast(error.message || error.error?.message || '提交失败')
  } finally {
    submitting.value = false
  }
}

// 检查表单是否有数据
function hasFormData() {
  const data = form.value
  return !!(data.project_id || data.work_location || data.work_content ||
    (selectedWorkTypes.value && selectedWorkTypes.value.length > 0) ||
    (data.work_date) || (data.contact_person) ||
    (data.contact_phone))
}

// 处理返回按钮
async function handleBack() {
  if (hasFormData()) {
    try {
      await showConfirmDialog({
        title: '提示',
        message: '表单有未保存的数据，是否保存为草稿？',
        confirmButtonText: '保存草稿',
        cancelButtonText: '放弃'
      })
      // 用户选择保存草稿
      await saveDraft()
    } catch {
      // 用户选择放弃，直接返回
      router.back()
    }
  } else {
    router.back()
  }
}

// 保存为草稿
async function saveDraft() {
  try {
    submitting.value = true
    const result = await createAppointment(form.value)
    showSuccessToast('草稿已保存')
    router.replace('/appointments')
  } catch (error) {
    showFailToast(error.message || '保存草稿失败')
  } finally {
    submitting.value = false
  }
}

function onDateConfirm(result) {
  let dateArray
  if (Array.isArray(result)) {
    dateArray = result
  } else if (result && result.selectedValues) {
    dateArray = result.selectedValues
  } else if (result && Array.isArray(result.value)) {
    dateArray = result.value
  } else {
    dateArray = currentDate.value
  }

  const [year, month, day] = dateArray
  const monthStr = month.toString().padStart(2, '0')
  const dayStr = day.toString().padStart(2, '0')
  form.value.work_date = `${year}-${monthStr}-${dayStr}`

  // 清空已选时间段
  form.value.time_slot = ''

  // 检查该日期是否有空闲人员
  checkAvailabilityForDate(form.value.work_date)

  // 获取该日期的时间段统计
  fetchTimeSlotStatistics()

  showDatePicker.value = false
}

function onTimeSlotConfirm({ selectedOptions }) {
  form.value.time_slot = selectedOptions[0].value
  showTimeSlotPicker.value = false
}

// 选择时间段
function selectTimeSlot(value) {
  form.value.time_slot = value
  showTimeSlotPicker.value = false
}

// 获取时间段状态标签类型
function getTimeSlotStatusTagType(status) {
  switch (status) {
    case 'full': return 'danger'
    case 'busy': return 'warning'
    case 'moderate': return 'primary'
    case 'available': return 'success'
    default: return 'default'
  }
}

function toggleWorkType(value) {
  const index = selectedWorkTypes.value.indexOf(value)
  if (index !== -1) {
    selectedWorkTypes.value.splice(index, 1)
  } else {
    selectedWorkTypes.value.push(value)
  }
}

function confirmWorkType() {
  form.value.work_type = selectedWorkTypes.value.join(',')
  showWorkTypePicker.value = false
}

// 根据用户ID生成头像背景色
function getAvatarColor(userId) {
  const colors = [
    '#f56c6c', '#e6a23c', '#409eff', '#67c23a',
    '#909399', '#f75da8', '#fb8c00', '#9c27b0'
  ]
  return colors[userId % colors.length]
}

// 获取每日预约统计数据
async function fetchDailyAppointments() {
  try {
    const now = new Date()
    const startOfMonth = new Date(now.getFullYear(), now.getMonth(), 1)
    const endOfMonth = new Date(now.getFullYear(), now.getMonth() + 1, 0)

    const response = await getDailyStatistics({
      start_date: startOfMonth.toISOString().split('T')[0],
      end_date: endOfMonth.toISOString().split('T')[0]
    })

    const apiData = response.data
    const data = apiData?.data || apiData

    totalWorkers.value = data?.total_workers || 0

    // 将统计数据转换为日期映射
    const stats = {}
    if (data?.statistics && Array.isArray(data.statistics)) {
      data.statistics.forEach((item) => {
        stats[item.date] = {
          total: item.total_count,
          urgent: item.urgent_count
        }
      })
    }

    dailyAppointments.value = stats
  } catch (error) {
    console.error('获取每日预约数据失败:', error)
  }
}

// 添加日期标记到 DOM
function addDateMarkers() {
  // 清除之前的标记
  document.querySelectorAll('.appointment-marker-mobile').forEach(el => el.remove())

  const [year, month] = currentDate.value

  // Vant DatePicker 的 DOM 结构
  const dateItems = document.querySelectorAll('.van-picker-column__item')

  dateItems.forEach(item => {
    const dayText = item.textContent.trim()
    const day = parseInt(dayText)
    if (isNaN(day) || day < 1 || day > 31) return

    const monthStr = month.toString().padStart(2, '0')
    const dayStr = day.toString().padStart(2, '0')
    const dateStr = `${year}-${monthStr}-${dayStr}`

    // 检查是否有预约
    const info = dailyAppointments.value[dateStr]
    if (info && info.total > 0) {
      // 确定标记颜色和大小
      let color = '#52c41a' // 绿色
      let size = 6

      if (info.urgent > 0) {
        color = '#f56c6c' // 红色 - 加急
        size = 8
      } else if (totalWorkers.value > 0) {
        const ratio = info.total / totalWorkers.value
        if (ratio >= 1) {
          color = '#f5222d' // 红色 - 超载
          size = 8
        } else if (ratio >= 0.75) {
          color = '#fa8c16' // 橙色
          size = 6
        } else if (ratio >= 0.4) {
          color = '#1890ff' // 蓝色
          size = 6
        }
      }

      // 创建标记元素
      const marker = document.createElement('span')
      marker.className = 'appointment-marker-mobile'
      marker.style.cssText = `
        position: absolute;
        top: 8px;
        right: 8px;
        width: ${size}px;
        height: ${size}px;
        background-color: ${color};
        border-radius: 50%;
        z-index: 10;
      `

      if (info.urgent > 0) {
        marker.style.boxShadow = '0 0 0 2px rgba(245, 108, 108, 0.2)'
      }

      // 添加到日期项
      item.style.position = 'relative'
      item.appendChild(marker)
    }
  })
}

// 日期选择器打开时
watch(showDatePicker, async (newVal) => {
  if (newVal) {
    await fetchDailyAppointments()
    // 等待 DOM 更新后添加标记
    setTimeout(addDateMarkers, 100)
  }
})

// 日期改变后重新添加标记
watch(currentDate, () => {
  if (showDatePicker.value) {
    setTimeout(addDateMarkers, 100)
  }
})

// 日期选择器关闭处理
function onDatePickerClose() {
  showDatePicker.value = false
}
</script>

<style scoped>
.appointment-create {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 140px;
}

.van-cell-group {
  margin-bottom: 12px;
}

.submit-bar {
  position: fixed;
  bottom: 50px;
  left: 0;
  right: 0;
  padding: 16px;
  background-color: #fff;
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.1);
  z-index: 999;
}

.work-type-picker {
  max-height: 60vh;
  overflow-y: auto;
}

.worker-picker {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.worker-selection-info {
  padding: 8px 16px;
  background: #f5f5f5;
  text-align: center;
}

.selected-count {
  color: #1989fa;
  font-size: 14px;
}

.worker-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  padding: 16px;
  overflow-y: auto;
  flex: 1;
}

.worker-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 12px 8px;
  background: white;
  border: 2px solid #e5e5e5;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.worker-card:active {
  transform: scale(0.95);
}

.worker-card.selected {
  border-color: #1989fa;
  background: #ecf5ff;
}

.worker-avatar {
  position: relative;
  margin-bottom: 8px;
}

.worker-avatar-placeholder {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 20px;
  font-weight: bold;
}

.selected-badge {
  position: absolute;
  top: -4px;
  right: -4px;
  width: 18px;
  height: 18px;
  background: #1989fa;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 12px;
}

.worker-name {
  font-size: 12px;
  color: #323233;
  text-align: center;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  width: 100%;
}

.worker-card .van-checkbox {
  position: absolute;
  top: 8px;
  right: 8px;
  opacity: 0;
}

.empty-workers {
  grid-column: 1 / -1;
  padding: 40px 0;
}

.selected-workers-display {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  padding: 12px 16px;
  background: #f5f5f5;
}

.selected-worker-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 10px;
  background: white;
  border-radius: 20px;
  border: 1px solid #e5e5e5;
}

.worker-avatar-small {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 16px;
  font-weight: bold;
  flex-shrink: 0;
}

.selected-worker-name {
  font-size: 13px;
  color: #323233;
  max-width: 100px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.remove-worker-icon {
  color: #969799;
  font-size: 14px;
  padding: 4px;
  cursor: pointer;
}

.remove-worker-icon:active {
  color: #323233;
}

/* 时间段选择器样式 */
.time-slot-picker {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.time-slot-picker .van-loading {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.time-slot-picker .van-cell-group {
  flex: 1;
  overflow-y: auto;
}

.time-slot-full {
  opacity: 0.5;
}

.time-slot-busy {
  background-color: #fff7e6;
}

.time-slot-available {
  background-color: #f0f9ff;
}
</style>
