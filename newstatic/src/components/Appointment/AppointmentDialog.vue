<template>
  <el-dialog
    v-model="dialogVisible"
    :title="dialogTitle"
    width="800px"
    :close-on-click-modal="false"
    @close="handleCloseClick"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-width="120px"
    >
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="项目" prop="project_id">
            <ProjectSelector
              v-model="formData.project_id"
              :projects="projectList"
              placeholder="请选择项目"
              style="width: 100%"
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="作业类型" prop="work_type">
            <el-select
              v-model="selectedWorkTypes"
              multiple
              collapse-tags
              collapse-tags-tooltip
              placeholder="请选择作业类型"
              style="width: 100%"
              @change="handleWorkTypeChange"
            >
              <el-option
                v-for="type in workTypeOptions"
                :key="type.value"
                :label="type.label"
                :value="type.value"
              />
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="作业日期" prop="work_date">
            <el-date-picker
              v-model="workDate"
              type="date"
              placeholder="选择日期"
              value-format="YYYY-MM-DD"
              style="width: 100%"
              :disabled-date="disabledDate"
              @panel-change="handlePanelChange"
              @visible-change="handleDatePickerVisibleChange"
            >
            </el-date-picker>
            <div class="date-legend">
              <span class="legend-item">
                <span class="legend-dot dot-available"></span>
                <span class="legend-text">空闲</span>
              </span>
              <span class="legend-item">
                <span class="legend-dot dot-light"></span>
                <span class="legend-text">有任务 (&lt;40%)</span>
              </span>
              <span class="legend-item">
                <span class="legend-dot dot-normal"></span>
                <span class="legend-text">正常 (40-75%)</span>
              </span>
              <span class="legend-item">
                <span class="legend-dot dot-busy"></span>
                <span class="legend-text">紧张 (75-100%)</span>
              </span>
              <span class="legend-item">
                <span class="legend-dot dot-overload"></span>
                <span class="legend-text">超载 (≥100%)</span>
              </span>
              <span class="legend-item">
                <span class="legend-dot dot-urgent"></span>
                <span class="legend-text">加急</span>
              </span>
            </div>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="时间段" prop="time_slot">
            <el-select v-model="formData.time_slot" placeholder="请选择" style="width: 100%">
              <el-option
                v-for="slot in timeSlotOptions"
                :key="slot.value"
                :label="slot.label"
                :value="slot.value"
                :disabled="isTimeSlotDisabled(slot.value)"
              >
                <div class="time-slot-option">
                  <span class="time-slot-label">{{ slot.label }}</span>
                  <span v-if="getTimeSlotCount(slot.value) > 0" class="time-slot-count" :class="getTimeSlotStatus(slot.value)">
                    {{ getTimeSlotCount(slot.value) }}单
                  </span>
                </div>
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>

      <el-form-item label="作业地点" prop="work_location">
        <el-input v-model="formData.work_location" placeholder="请输入作业地点" />
      </el-form-item>

      <el-form-item label="作业内容" prop="work_content">
        <el-input
          v-model="formData.work_content"
          type="textarea"
          :rows="3"
          placeholder="请输入作业内容"
        />
      </el-form-item>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="联系人" prop="contact_person">
            <el-input v-model="formData.contact_person" placeholder="请输入联系人" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="联系电话" prop="contact_phone">
            <el-input v-model="formData.contact_phone" placeholder="请输入联系电话" />
          </el-form-item>
        </el-col>
      </el-row>

      <el-form-item label="是否加急">
        <el-switch v-model="formData.is_urgent" />
      </el-form-item>

      <template v-if="formData.is_urgent">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="优先级" prop="priority">
              <el-input-number
                v-model="formData.priority"
                :min="0"
                :max="10"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="加急原因" prop="urgent_reason">
              <el-input v-model="formData.urgent_reason" placeholder="请输入加急原因" />
            </el-form-item>
          </el-col>
        </el-row>
      </template>

      <el-form-item label="作业人员">
        <div v-if="!showWorkerSelector" class="worker-selector-trigger">
          <el-button @click="handleShowWorkerSelector" :disabled="loadingWorkers">
            选择作业人员
          </el-button>
          <span v-if="selectedWorkerIds.length > 0" class="selected-worker-hint">
            已选择 {{ selectedWorkerIds.length }} 人
          </span>
        </div>
        <div v-else v-loading="loadingWorkers" class="worker-selector-container">
          <div class="worker-selector-header">
            <span class="worker-selector-title">选择作业人员（可多选）</span>
            <el-button type="text" @click="handleHideWorkerSelector">收起</el-button>
          </div>
          <div class="worker-selector-grid">
            <div
              v-for="worker in workerList"
              :key="worker.id"
              :class="['worker-selector-card', { 'selected': isWorkerSelected(worker.id) }]"
              @click="toggleWorkerSelection(worker.id)"
            >
              <div class="worker-card-avatar">
                <el-avatar
                  :size="48"
                  :src="worker.avatar || undefined"
                  :style="!worker.avatar ? { backgroundColor: getAvatarColor(worker.id) } : {}"
                >
                  {{ (worker.full_name || worker.username || '?').charAt(0) }}
                </el-avatar>
                <div v-if="isWorkerSelected(worker.id)" class="selected-badge">
                  <el-icon><Check /></el-icon>
                </div>
              </div>
              <div class="worker-card-info">
                <div class="worker-card-name">{{ worker.full_name || worker.username }}</div>
                <div v-if="worker.email" class="worker-card-email">{{ worker.email }}</div>
              </div>
            </div>
            <div v-if="workerList.length === 0 && !loadingWorkers" class="empty-workers">
              <el-empty description="暂无作业人员" :image-size="60" />
            </div>
          </div>
        </div>
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="handleCloseClick">取消</el-button>
      <el-button type="primary" :loading="submitting" @click="handleSubmit">
        {{ mode === 'create' ? '提交' : '保存' }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Check } from '@element-plus/icons-vue'
import { appointmentApi, projectApi } from '@/api'
import ProjectSelector from '@/components/common/ProjectSelector.vue'

const projectList = ref([])

// 作业类型选项（与移动端保持一致）
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

// 时间段选项
const timeSlotOptions = [
  { value: 'morning', label: '上午 (8:00-11:30)' },
  { value: 'noon', label: '中午 (12:00-13:30)' },
  { value: 'afternoon', label: '下午 (13:30-16:30)' },
  { value: 'full_day', label: '全天' }
]

const selectedWorkTypes = ref([])

const props = defineProps({
  modelValue: Boolean,
  appointment: Object,
  mode: {
    type: String,
    default: 'create'
  }
})

const emit = defineEmits(['update:modelValue', 'success'])

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const dialogTitle = computed(() => {
  return props.mode === 'create' ? '创建预约单' : '编辑预约单'
})

const formRef = ref(null)
const submitting = ref(false)
const workDate = ref('')
const dailyAppointments = ref({}) // 存储每日预约数据，格式: { '2026-02-11': { total: 5, urgent: 2 } }
const timeSlotStatistics = ref({}) // 存储时间段统计数据，格式: { '2026-02-11': { morning: 2, noon: 1, afternoon: 3, full_day: 1 } }
const totalWorkers = ref(0) // 总作业人员数量
const currentViewDate = ref(new Date()) // 当前日历视图的日期（用于构建完整日期）

// 获取项目列表
async function fetchProjects() {
  try {
    const response = await projectApi.getList({ page: 1, page_size: 1000 })
    projectList.value = response.data?.projects || response.data || []
  } catch (error) {
    console.error('获取项目列表失败:', error)
  }
}

// 获取作业人员列表
async function fetchWorkers() {
  loadingWorkers.value = true
  try {
    const response = await appointmentApi.getWorkersList()
    workerList.value = response.data || []
  } catch (error) {
    console.error('获取作业人员列表失败:', error)
  } finally {
    loadingWorkers.value = false
  }
}

// 获取每日预约统计数据
async function fetchDailyAppointments() {
  try {
    // 使用当前视图日期的月份
    const viewDate = currentViewDate.value
    const startOfMonth = new Date(viewDate.getFullYear(), viewDate.getMonth(), 1)
    const endOfMonth = new Date(viewDate.getFullYear(), viewDate.getMonth() + 1, 0)

    const startDate = startOfMonth.toISOString().split('T')[0]
    const endDate = endOfMonth.toISOString().split('T')[0]

    const response = await appointmentApi.getDailyStatistics({
      start_date: startDate,
      end_date: endDate
    })

    // 从后端获取统计数据和总作业人员数
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

// 获取时间段统计数据
async function fetchTimeSlotStatistics() {
  if (!workDate.value) return

  try {
    const response = await appointmentApi.getTimeSlotStatistics({
      date: workDate.value
    })

    const apiData = response.data
    const data = apiData?.data || apiData

    // 将统计数据存储到日期对应的对象中
    if (data?.statistics && Array.isArray(data.statistics)) {
      const slotStats = {}
      data.statistics.forEach((item) => {
        slotStats[item.time_slot] = {
          count: item.total_count,
          totalWorkers: item.total_workers
        }
      })
      timeSlotStatistics.value[workDate.value] = slotStats
    } else {
      timeSlotStatistics.value[workDate.value] = {}
    }
  } catch (error) {
    console.error('获取时间段统计数据失败:', error)
    timeSlotStatistics.value[workDate.value] = {}
  }
}

// 获取时间段的任务数量
function getTimeSlotCount(slotValue) {
  if (!workDate.value) return 0
  const slotStats = timeSlotStatistics.value[workDate.value]
  return slotStats?.[slotValue]?.count || 0
}

// 获取时间段的状态类
function getTimeSlotStatus(slotValue) {
  const count = getTimeSlotCount(slotValue)
  if (count === 0) return ''

  const slotStats = timeSlotStatistics.value[workDate.value]?.[slotValue]
  if (!slotStats) return ''

  const totalWorkers = slotStats.totalWorkers || totalWorkers.value
  if (totalWorkers === 0) return 'normal'

  const ratio = count / totalWorkers
  if (ratio >= 1) return 'overload'
  if (ratio >= 0.75) return 'busy'
  return 'normal'
}

// 判断时间段是否禁用
function isTimeSlotDisabled(slotValue) {
  const status = getTimeSlotStatus(slotValue)
  return status === 'overload'
}

// 获取日期单元格的CSS类名（已弃用，保留用于兼容性）
function getCellClassName(date) {
  // 不再使用 cell-class-name 属性，改用 DOM 操作方式
  return ''
}

// 从单元格数据构建完整日期字符串
function getDateFromCell(cell) {
  if (!cell || !cell.text) return null

  // cell.text 是日期数字，如 "15"
  const day = parseInt(cell.text, 10)
  if (isNaN(day)) return null

  // 使用当前视图的年月构建完整日期
  const year = currentViewDate.value.getFullYear()
  const month = currentViewDate.value.getMonth()

  // 构建日期字符串 (YYYY-MM-DD)
  const date = new Date(year, month, day)
  return date.toISOString().split('T')[0]
}

// 检查日期是否有预约（接受单元格数据对象）
function hasAppointmentForCell(cellData) {
  if (!cellData || !cellData.cell) return false

  const dateStr = getDateFromCell(cellData.cell)
  if (!dateStr) return false

  const info = dailyAppointments.value[dateStr]
  return info && info.total > 0
}

// 获取标记的类名（接受单元格数据对象）
function getMarkerClassForCell(cellData) {
  if (!cellData || !cellData.cell) return ''

  const dateStr = getDateFromCell(cellData.cell)
  if (!dateStr) return ''

  const info = dailyAppointments.value[dateStr]
  if (!info || info.total === 0) return ''

  if (info.urgent > 0) return 'marker-urgent'

  if (totalWorkers.value > 0) {
    const ratio = info.total / totalWorkers.value
    if (ratio >= 1) return 'marker-overload'
    if (ratio >= 0.75) return 'marker-busy'
    if (ratio >= 0.4) return 'marker-normal'
  }

  return 'marker-light'
}

// 检查日期是否有预约（保留用于其他可能的用途）
function hasAppointment(date) {
  if (!date) return false

  let dateStr
  if (typeof date === 'string') {
    dateStr = date
  } else if (date instanceof Date) {
    dateStr = date.toISOString().split('T')[0]
  } else {
    return false
  }

  const info = dailyAppointments.value[dateStr]
  return info && info.total > 0
}

// 获取标记的类名（保留用于其他可能的用途）
function getMarkerClass(date) {
  if (!date) return ''

  let dateStr
  if (typeof date === 'string') {
    dateStr = date
  } else if (date instanceof Date) {
    dateStr = date.toISOString().split('T')[0]
  } else {
    return ''
  }

  const info = dailyAppointments.value[dateStr]
  if (!info || info.total === 0) return ''

  if (info.urgent > 0) return 'marker-urgent'

  if (totalWorkers.value > 0) {
    const ratio = info.total / totalWorkers.value
    if (ratio >= 1) return 'marker-overload'
    if (ratio >= 0.75) return 'marker-busy'
    if (ratio >= 0.4) return 'marker-normal'
  }

  return 'marker-light'
}

// 禁用今天及过去的日期
function disabledDate(time) {
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  // 禁用今天和过去的日期
  return time.getTime() <= today.getTime()
}

// 面板变化时更新当前视图日期
function handlePanelChange(date) {
  if (date instanceof Date) {
    currentViewDate.value = date
  } else if (Array.isArray(date) && date[0] instanceof Date) {
    currentViewDate.value = date[0]
  }
  // 重新获取该月份的统计数据
  fetchDailyAppointments()
  // 等待DOM更新后添加标记
  setTimeout(addDateMarkers, 100)
}

// 日期选择器面板显示/隐藏时处理
function handleDatePickerVisibleChange(visible) {
  if (visible) {
    // 等待DOM完全渲染后添加标记
    setTimeout(addDateMarkers, 200)
  }
}

// 直接操作DOM添加日期标记
function addDateMarkers() {
  // 清除之前的标记
  document.querySelectorAll('.appointment-marker').forEach(el => el.remove())

  // 获取日历面板中的所有日期单元格
  const dateCells = document.querySelectorAll('.el-date-table td')

  dateCells.forEach(cell => {
    // 跳过非日期单元格（表头等）
    const cellText = cell.querySelector('.el-date-table-cell__text')
    if (!cellText) return

    const day = cellText.textContent.trim()
    if (!day || isNaN(parseInt(day))) return

    // 构建完整日期
    const year = currentViewDate.value.getFullYear()
    const month = currentViewDate.value.getMonth() + 1 // 月份从1开始
    const monthStr = month.toString().padStart(2, '0')
    const dayStr = day.toString().padStart(2, '0')
    const dateStr = `${year}-${monthStr}-${dayStr}`

    // 检查是否有预约
    const info = dailyAppointments.value[dateStr]
    if (info && info.total > 0) {
      // 确定标记颜色和大小
      let color = '#52c41a' // 绿色 - 有少量任务
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
          color = '#fa8c16' // 橙色 - 紧张
          size = 6
        } else if (ratio >= 0.4) {
          color = '#1890ff' // 蓝色 - 正常
          size = 6
        }
      }

      // 创建标记元素
      const marker = document.createElement('span')
      marker.className = 'appointment-marker'
      marker.style.cssText = `
        position: absolute;
        bottom: 2px;
        left: 50%;
        transform: translateX(-50%);
        width: ${size}px;
        height: ${size}px;
        background-color: ${color};
        border-radius: 50%;
        z-index: 10;
        pointer-events: none;
      `

      if (info.urgent > 0) {
        marker.style.boxShadow = '0 0 0 2px rgba(245, 108, 108, 0.2)'
      }

      // 添加到单元格
      const cellContent = cell.querySelector('.el-date-table-cell')
      if (cellContent) {
        cellContent.style.position = 'relative'
        cellContent.appendChild(marker)
      }
    }
  })
}

const formData = ref({
  project_id: null,
  work_type: '',
  work_date: '',
  time_slot: '',
  work_location: '',
  work_content: '',
  contact_person: '',
  contact_phone: '',
  is_urgent: false,
  priority: 0,
  urgent_reason: '',
  assigned_worker_id: null, // 保留用于兼容
  assigned_worker_ids: [] // 新增多选支持
})

// 选中的作业人员ID列表
const selectedWorkerIds = ref([])

const formRules = {
  work_date: [{ required: true, message: '请选择作业日期', trigger: 'change' }],
  time_slot: [{ required: true, message: '请选择时间段', trigger: 'change' }],
  work_location: [{ required: true, message: '请输入作业地点', trigger: 'blur' }],
  work_content: [{ required: true, message: '请输入作业内容', trigger: 'blur' }],
  priority: [
    { validator: (rule, value, callback) => {
      if (formData.value.is_urgent && !value) {
        callback(new Error('加急预约必须设置优先级'))
      } else {
        callback()
      }
    }, trigger: 'change' }
  ],
  urgent_reason: [
    { validator: (rule, value, callback) => {
      if (formData.value.is_urgent && formData.value.priority >= 7 && !value) {
        callback(new Error('高优先级加急必须提供原因'))
      } else {
        callback()
      }
    }, trigger: 'blur' }
  ]
}

// TODO: 获取作业人员列表
const workerList = ref([])
const loadingWorkers = ref(false)
const showWorkerSelector = ref(false)

// 显示人员选择器
async function handleShowWorkerSelector() {
  showWorkerSelector.value = true
  if (workerList.value.length === 0) {
    await fetchWorkers()
  }
}

// 隐藏人员选择器
function handleHideWorkerSelector() {
  showWorkerSelector.value = false
}

// 获取已选择的人员名称
function getSelectedWorkerName() {
  if (selectedWorkerIds.value.length === 0) return ''
  const workers = workerList.value.filter(w => selectedWorkerIds.value.includes(w.id))
  return workers.map(w => w.full_name || w.username).join('、')
}

// 检查作业人员是否被选中
function isWorkerSelected(workerId) {
  return selectedWorkerIds.value.includes(workerId)
}

// 切换作业人员选择状态（支持多选）
function toggleWorkerSelection(workerId) {
  const index = selectedWorkerIds.value.indexOf(workerId)
  if (index > -1) {
    // 取消选中
    selectedWorkerIds.value.splice(index, 1)
  } else {
    // 选中
    selectedWorkerIds.value.push(workerId)
  }
  // 同步到 formData
  formData.value.assigned_worker_ids = [...selectedWorkerIds.value]
  // 保留兼容性：如果是单选，第一个值作为 assigned_worker_id
  formData.value.assigned_worker_id = selectedWorkerIds.value.length > 0 ? selectedWorkerIds.value[0] : null
}

// 根据用户ID生成头像背景色
function getAvatarColor(userId) {
  const colors = [
    '#f56c6c', '#e6a23c', '#409eff', '#67c23a',
    '#909399', '#f75da8', '#fb8c00', '#9c27b0'
  ]
  return colors[userId % colors.length]
}

// 处理作业类型变化
function handleWorkTypeChange(values) {
  formData.value.work_type = values.join(',')
}

// 编辑时加载作业类型
watch(() => props.appointment, (val) => {
  if (val && props.mode === 'edit') {
    Object.assign(formData.value, val)
    workDate.value = val.work_date
    // 解析作业类型字符串为数组
    if (val.work_type) {
      selectedWorkTypes.value = val.work_type.split(',').filter(v => v)
    }
    // 加载已选择的作业人员
    // assigned_worker_ids 可能是 JSON 字符串格式 "[1,2,3]" 或数组
    if (val.assigned_worker_ids) {
      if (typeof val.assigned_worker_ids === 'string') {
        try {
          selectedWorkerIds.value = JSON.parse(val.assigned_worker_ids)
        } catch (e) {
          console.error('解析 assigned_worker_ids 失败:', e)
          selectedWorkerIds.value = []
        }
      } else if (Array.isArray(val.assigned_worker_ids)) {
        selectedWorkerIds.value = [...val.assigned_worker_ids]
      } else {
        selectedWorkerIds.value = []
      }
    } else if (val.assigned_worker_id) {
      selectedWorkerIds.value = [val.assigned_worker_id]
    } else {
      selectedWorkerIds.value = []
    }
  } else {
    resetForm()
  }
}, { immediate: true })

// 对话框打开时加载数据
watch(dialogVisible, (val) => {
  if (val) {
    fetchProjects()
    fetchDailyAppointments() // 获取每日预约数据用于日期标注
    // 不自动加载人员列表，等待用户点击选择按钮
  }
})

watch(workDate, (val) => {
  formData.value.work_date = val
  // 当选择日期时，获取时间段统计数据
  if (val) {
    fetchTimeSlotStatistics()
  }
})

function resetForm() {
  formData.value = {
    project_id: null,
    work_type: '',
    work_date: '',
    time_slot: '',
    work_location: '',
    work_content: '',
    contact_person: '',
    contact_phone: '',
    is_urgent: false,
    priority: 0,
    urgent_reason: '',
    assigned_worker_id: null,
    assigned_worker_ids: []
  }
  workDate.value = ''
  selectedWorkTypes.value = []
  selectedWorkerIds.value = []
  showWorkerSelector.value = false
  formRef.value?.clearValidate()
}

async function handleSubmit() {
  try {
    await formRef.value.validate()
    submitting.value = true

    // 准备提交数据，确保包含多选的作业人员
    // 获取选中的作业人员名称（逗号分隔）
    const selectedWorkers = workerList.value.filter(w => selectedWorkerIds.value.includes(w.id))
    const workerNames = selectedWorkers.map(w => w.full_name || w.username).join(',')

    const submitData = {
      ...formData.value,
      // 后端期望 assigned_worker_ids 是 JSON 字符串格式 "[1,2,3]"
      assigned_worker_ids: JSON.stringify(selectedWorkerIds.value),
      // 作业人员名称是逗号分隔的字符串
      assigned_worker_names: workerNames,
      // 保留兼容性
      assigned_worker_id: selectedWorkerIds.value.length > 0 ? selectedWorkerIds.value[0] : null
    }

    if (props.mode === 'create') {
      // 创建时直接提交
      const createResult = await appointmentApi.create(submitData)
      // 创建成功后提交
      const appointmentId = createResult.data?.id || createResult.id
      if (appointmentId) {
        await appointmentApi.submit(appointmentId)
      }
      ElMessage.success('提交成功')
    } else {
      await appointmentApi.update(props.appointment.id, submitData)
      ElMessage.success('保存成功')
    }

    emit('success')
    handleClose()
  } catch (error) {
    if (error !== false) { // 表单验证失败时不显示错误
      ElMessage.error(error.message || '操作失败')
    }
  } finally {
    submitting.value = false
  }
}

// 检查表单是否有数据
function hasFormData() {
  const data = formData.value
  return !!(data.work_location || data.work_content ||
    (selectedWorkTypes.value && selectedWorkTypes.value.length > 0) ||
    (data.work_date) || (data.contact_person) ||
    (data.contact_phone) || (data.project_id))
}

// 处理关闭按钮点击
async function handleCloseClick() {
  if (hasFormData()) {
    try {
      await ElMessageBox.confirm(
        '表单有未保存的数据，是否保存为草稿？',
        '提示',
        {
          confirmButtonText: '保存草稿',
          cancelButtonText: '放弃',
          type: 'warning'
        }
      )
      // 用户选择保存草稿
      await saveDraft()
    } catch {
      // 用户选择放弃，直接关闭
      handleClose()
    }
  } else {
    handleClose()
  }
}

// 保存为草稿
async function saveDraft() {
  try {
    submitting.value = true
    if (props.mode === 'create') {
      await appointmentApi.create(formData.value)
      ElMessage.success('草稿已保存')
    } else {
      await appointmentApi.update(props.appointment.id, formData.value)
      ElMessage.success('草稿已更新')
    }
    emit('success')
    handleClose()
  } catch (error) {
    ElMessage.error(error.message || '保存草稿失败')
  } finally {
    submitting.value = false
  }
}

function handleClose() {
  resetForm()
  emit('update:modelValue', false)
}
</script>

<style scoped>
.worker-selector-trigger {
  display: flex;
  align-items: center;
  gap: 12px;
}

.selected-worker-hint {
  color: #409eff;
  font-size: 14px;
}

.worker-selector-container {
  width: 100%;
}

.worker-selector-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding: 8px 0;
}

.worker-selector-title {
  font-weight: 500;
  color: #303133;
}

.worker-selector-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 12px;
  padding: 12px;
  background: #f5f7fa;
  border-radius: 8px;
  min-height: 100px;
}

.worker-selector-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 12px;
  background: white;
  border: 2px solid #dcdfe6;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.worker-selector-card:hover {
  border-color: #409eff;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.2);
}

.worker-selector-card.selected {
  border-color: #409eff;
  background: #ecf5ff;
}

.worker-card-avatar {
  position: relative;
  margin-bottom: 8px;
}

.selected-badge {
  position: absolute;
  top: -4px;
  right: -4px;
  width: 20px;
  height: 20px;
  background: #409eff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 12px;
}

.worker-card-info {
  text-align: center;
}

.worker-card-name {
  font-size: 13px;
  font-weight: 500;
  color: #303133;
  margin-bottom: 2px;
}

.worker-card-email {
  font-size: 11px;
  color: #909399;
}

.empty-workers {
  grid-column: 1 / -1;
  padding: 20px 0;
}

/* 日期图例样式 */
.date-legend {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-top: 8px;
  padding: 8px 12px;
  background: #fafafa;
  border-radius: 4px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #606266;
}

.legend-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  display: inline-block;
}

.dot-available {
  background: #f0f9ff;
  border: 1px solid #b3d8ff;
}

.dot-light {
  background: #f6ffed;
  border: 1px solid #b7eb8f;
}

.dot-normal {
  background: #e6f7ff;
  border: 1px solid #91d5ff;
}

.dot-busy {
  background: #fff7e6;
  border: 1px solid #ffd591;
}

.dot-overload {
  background: #fff1f0;
  border: 1px solid #ffa39e;
}

.dot-urgent {
  background: #fef0f0;
  border: 1px solid #f56c6c;
}

/* 日期选择器中的自定义单元格样式 - 使用圆圈标记 */
.el-date-table-cell {
  position: relative;
}

.date-marker {
  position: absolute;
  bottom: 3px;
  left: 50%;
  transform: translateX(-50%);
  border-radius: 50%;
  z-index: 1;
}

.date-marker.marker-light,
.date-marker.marker-normal,
.date-marker.marker-busy {
  width: 6px;
  height: 6px;
}

.date-marker.marker-overload,
.date-marker.marker-urgent {
  width: 8px;
  height: 8px;
}

.date-marker.marker-light {
  background-color: #52c41a;
}

.date-marker.marker-normal {
  background-color: #1890ff;
}

.date-marker.marker-busy {
  background-color: #fa8c16;
}

.date-marker.marker-overload {
  background-color: #f5222d;
}

.date-marker.marker-urgent {
  background-color: #f56c6c;
  box-shadow: 0 0 0 2px rgba(245, 108, 108, 0.2);
}

/* 日期单元格标记样式 - 使用伪元素显示圆圈 */
:deep(.el-date-table td.date-cell-light .el-date-table-cell),
:deep(.el-date-table td.disabled.date-cell-light .el-date-table-cell) {
  position: relative !important;
}

:deep(.el-date-table td.date-cell-light .el-date-table-cell::after),
:deep(.el-date-table td.disabled.date-cell-light .el-date-table-cell::after) {
  content: '' !important;
  position: absolute !important;
  bottom: 2px !important;
  left: 50% !important;
  transform: translateX(-50%) !important;
  width: 6px !important;
  height: 6px !important;
  background-color: #52c41a !important;
  border-radius: 50% !important;
  z-index: 10 !important;
  pointer-events: none !important;
}

:deep(.el-date-table td.date-cell-normal .el-date-table-cell),
:deep(.el-date-table td.disabled.date-cell-normal .el-date-table-cell) {
  position: relative !important;
}

:deep(.el-date-table td.date-cell-normal .el-date-table-cell::after),
:deep(.el-date-table td.disabled.date-cell-normal .el-date-table-cell::after) {
  content: '' !important;
  position: absolute !important;
  bottom: 2px !important;
  left: 50% !important;
  transform: translateX(-50%) !important;
  width: 6px !important;
  height: 6px !important;
  background-color: #1890ff !important;
  border-radius: 50% !important;
  z-index: 10 !important;
  pointer-events: none !important;
}

:deep(.el-date-table td.date-cell-busy .el-date-table-cell),
:deep(.el-date-table td.disabled.date-cell-busy .el-date-table-cell) {
  position: relative !important;
}

:deep(.el-date-table td.date-cell-busy .el-date-table-cell::after),
:deep(.el-date-table td.disabled.date-cell-busy .el-date-table-cell::after) {
  content: '' !important;
  position: absolute !important;
  bottom: 2px !important;
  left: 50% !important;
  transform: translateX(-50%) !important;
  width: 6px !important;
  height: 6px !important;
  background-color: #fa8c16 !important;
  border-radius: 50% !important;
  z-index: 10 !important;
  pointer-events: none !important;
}

:deep(.el-date-table td.date-cell-overload .el-date-table-cell),
:deep(.el-date-table td.disabled.date-cell-overload .el-date-table-cell) {
  position: relative !important;
}

:deep(.el-date-table td.date-cell-overload .el-date-table-cell::after),
:deep(.el-date-table td.disabled.date-cell-overload .el-date-table-cell::after) {
  content: '' !important;
  position: absolute !important;
  bottom: 2px !important;
  left: 50% !important;
  transform: translateX(-50%) !important;
  width: 8px !important;
  height: 8px !important;
  background-color: #f5222d !important;
  border-radius: 50% !important;
  z-index: 10 !important;
  pointer-events: none !important;
}

:deep(.el-date-table td.date-cell-urgent .el-date-table-cell),
:deep(.el-date-table td.disabled.date-cell-urgent .el-date-table-cell) {
  position: relative !important;
}

:deep(.el-date-table td.date-cell-urgent .el-date-table-cell::after),
:deep(.el-date-table td.disabled.date-cell-urgent .el-date-table-cell::after) {
  content: '' !important;
  position: absolute !important;
  bottom: 2px !important;
  left: 50% !important;
  transform: translateX(-50%) !important;
  width: 8px !important;
  height: 8px !important;
  background-color: #f56c6c !important;
  border-radius: 50% !important;
  box-shadow: 0 0 0 2px rgba(245, 108, 108, 0.2) !important;
  z-index: 10 !important;
  pointer-events: none !important;
}

:deep(.el-date-table td.today .el-date-table-cell__text) {
  color: #409eff;
  font-weight: 600;
}

:deep(.el-date-table td.available .el-date-table-cell__text) {
  color: #606266;
}

/* 时间段选择器样式 */
.time-slot-option {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.time-slot-label {
  flex: 1;
}

.time-slot-count {
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 12px;
  font-weight: 500;
}

.time-slot-count.normal {
  background-color: #e6f7ff;
  color: #1890ff;
  border: 1px solid #91d5ff;
}

.time-slot-count.busy {
  background-color: #fff7e6;
  color: #fa8c16;
  border: 1px solid #ffd591;
}

.time-slot-count.overload {
  background-color: #fff1f0;
  color: #f5222d;
  border: 1px solid #ffa39e;
}
</style>
