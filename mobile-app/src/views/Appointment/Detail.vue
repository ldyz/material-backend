<template>
  <div class="appointment-detail">
    <van-nav-bar
      title="预约单详情"
      left-arrow
      @click-left="router.back()"
    >
      <template #right>
        <van-dropdown-menu v-if="appointment && canEditOrCancel">
          <van-dropdown-item>
            <van-cell
              v-if="isEditable(appointment.status)"
              title="编辑"
              is-link
              @click="router.push(`/appointment/${appointment.id}/edit`)"
            />
            <van-cell
              v-if="isCancellable(appointment.status)"
              title="取消预约"
              is-link
              @click="showCancelDialog"
            />
            <van-cell
              v-if="canStart(appointment.status)"
              title="开始作业"
              is-link
              @click="handleStartWork"
            />
            <van-cell
              v-if="canComplete(appointment.status)"
              title="完成作业"
              is-link
              @click="showCompleteDialog"
            />
          </van-dropdown-item>
        </van-dropdown-menu>
      </template>
    </van-nav-bar>

    <van-loading v-if="loading" type="spinner" color="#1989fa" />

    <div v-else-if="appointment" class="detail-content">
      <!-- 基本信息 -->
      <DetailInfoGroup
        title="基本信息"
        :item="appointment"
        status-type="appointment"
        :fields="basicFields"
      />

      <!-- 作业信息 -->
      <DetailInfoGroup
        title="作业信息"
        :item="appointment"
        :fields="workInfoFields"
      >
        <template #field-work_date="{ field, value }">
          {{ formatDateTime(appointment.work_date, appointment.time_slot) }}
        </template>
      </DetailInfoGroup>

      <!-- 优先级信息 -->
      <DetailInfoGroup
        v-if="appointment.is_urgent || appointment.priority > 0"
        title="优先级"
        :item="appointment"
        :fields="priorityFields"
      >
        <template #field-is_urgent="{ field, value }">
          <van-tag v-if="appointment.is_urgent" type="danger">是</van-tag>
          <van-tag v-else type="default">否</van-tag>
        </template>
      </DetailInfoGroup>

      <!-- 分配信息 -->
      <DetailInfoGroup
        title="分配信息"
        :item="appointment"
        :fields="assignmentFields"
      >
        <template #field-assigned_worker_name="{ field, value }">
          {{ appointment.assigned_worker_name || '未分配' }}
        </template>
      </DetailInfoGroup>

      <!-- 时间信息 -->
      <DetailInfoGroup
        title="时间信息"
        :item="appointment"
        :fields="timeFields"
        :show-empty="false"
      >
        <template #field-created_at="{ field, value }">
          {{ formatFullDateTime(appointment.created_at) }}
        </template>
        <template #field-submitted_at="{ field, value }">
          {{ formatFullDateTime(appointment.submitted_at) }}
        </template>
        <template #field-approved_at="{ field, value }">
          {{ formatFullDateTime(appointment.approved_at) }}
        </template>
        <template #field-completed_at="{ field, value }">
          {{ formatFullDateTime(appointment.completed_at) }}
        </template>
      </DetailInfoGroup>

      <!-- 审批历史 -->
      <van-cell-group title="审批历史" inset>
        <ApprovalTimeline
          :key="`timeline-${appointment.id}-${approvalLogs.length}`"
          :approval-logs="approvalLogs"
          :workflow-config="workflowConfig"
          :current-status="appointment.status"
        />
      </van-cell-group>

      <!-- 操作按钮 -->
      <div class="action-buttons">
        <!-- 草稿状态：显示提交审批按钮 -->
        <van-button
          v-if="appointment.status === 'draft'"
          type="primary"
          block
          @click="handleSubmit"
        >
          提交审批
        </van-button>

        <!-- 待审批状态：显示审批/拒绝按钮 -->
        <div v-if="canApprove" class="approve-buttons">
          <van-button
            type="success"
            block
            @click="showApproveDialog = true"
          >
            审批通过
          </van-button>
          <van-button
            type="danger"
            block
            @click="showRejectDialog = true"
            style="margin-top: 10px;"
          >
            审批拒绝
          </van-button>
        </div>

        <!-- 进行中状态：显示完成作业按钮 -->
        <van-button
          v-else-if="appointment.status === 'in_progress'"
          type="primary"
          block
          @click="showCompleteDialog = true"
        >
          完成作业
        </van-button>
      </div>
    </div>

    <van-empty v-else description="预约单不存在" />

    <!-- 审批对话框 -->
    <van-dialog
      v-model:show="showApproveDialog"
      title="审批通过"
      :show-confirm-button="false"
    >
      <van-form @submit="handleApproveSubmit">
        <!-- 审批意见 -->
        <van-field
          v-model="approveComment"
          type="textarea"
          placeholder="请输入审批意见（可选）"
          rows="3"
          autosize
        />

        <!-- 分配作业人员开关 -->
        <van-cell-group title="分配作业人员" inset>
          <van-field title="立即分配">
            <template #input>
              <van-switch v-model="approveForm.assignNow" />
            </template>
          </van-field>
        </van-cell-group>

        <!-- 作业人员选择 -->
        <van-cell-group v-if="approveForm.assignNow" title="选择作业人员" inset>
          <van-field
            readonly
            clickable
            :value="selectedWorkerName"
            placeholder="选择作业人员"
            @click="showWorkerPicker = true"
          />
        </van-cell-group>

        <!-- 修改作业时间开关 -->
        <van-cell-group title="修改作业时间" inset>
          <van-field title="调整时间">
            <template #input>
              <van-switch v-model="approveForm.reschedule" />
            </template>
          </van-field>
        </van-cell-group>

        <!-- 新作业时间选择 -->
        <van-cell-group v-if="approveForm.reschedule" title="新作业时间" inset>
          <van-field
            readonly
            clickable
            :value="newWorkDateLabel"
            placeholder="选择日期"
            @click="showDatePicker = true"
          />
          <van-field
            readonly
            clickable
            :value="newTimeSlotLabel"
            placeholder="选择时间段"
            @click="showTimeSlotPicker = true"
          />
          <van-cell title="可用作业人员">
            <template #value>
              <van-tag :type="availableWorkersCount > 0 ? 'success' : 'danger'">
                {{ availableWorkersCount }} 人可用
              </van-tag>
            </template>
          </van-cell>
        </van-cell-group>

        <!-- 操作按钮 -->
        <div class="dialog-actions">
          <van-button plain type="default" @click="showApproveDialog = false">
            取消
          </van-button>
          <van-button type="success" native-type="submit">
            审批通过
          </van-button>
        </div>
      </van-form>
    </van-dialog>

    <!-- 作业人员选择器 -->
    <WorkerPicker
      v-model="showWorkerPicker"
      title="选择作业人员"
      :multiple="false"
      @confirm="onWorkerConfirm"
    />

    <!-- 日期选择器 -->
    <van-calendar v-model:show="showDatePicker" @confirm="onDateConfirm" />

    <!-- 时间段选择器 -->
    <van-popup v-model:show="showTimeSlotPicker" position="bottom">
      <van-picker
        :columns="timeSlotColumns"
        @confirm="onTimeSlotConfirm"
        @cancel="showTimeSlotPicker = false"
      />
    </van-popup>

    <!-- 拒绝对话框 -->
    <van-dialog
      v-model:show="showRejectDialog"
      title="审批拒绝"
      show-cancel-button
      @confirm="handleReject"
    >
      <van-field
        v-model="rejectReason"
        type="textarea"
        placeholder="请输入拒绝原因（必填）"
        rows="3"
        autosize
        :rules="[{ required: true, message: '请输入拒绝原因' }]"
      />
    </van-dialog>

    <!-- 取消对话框 -->
    <van-dialog
      v-model:show="cancelDialogVisible"
      title="取消预约"
      show-cancel-button
      @confirm="handleCancel"
    >
      <van-field
        v-model="cancelReason"
        type="textarea"
        placeholder="请输入取消原因"
        rows="3"
      />
    </van-dialog>

    <!-- 完成对话框 -->
    <van-dialog
      v-model:show="completeDialogVisible"
      title="完成作业"
      show-cancel-button
      @confirm="handleComplete"
    >
      <van-field
        v-model="completionNote"
        type="textarea"
        placeholder="请输入完成备注"
        rows="3"
      />
    </van-dialog>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showSuccessToast, showFailToast, Dialog } from 'vant'
import ApprovalTimeline from '@/components/common/ApprovalTimeline.vue'
import WorkerPicker from '@/components/common/WorkerPicker.vue'
import DetailInfoGroup from '@/components/common/DetailInfoGroup.vue'
import { formatAppointmentDate, formatDateTime } from '@/composables/useDateTime'
import webSocketService from '@/utils/websocket'
import {
  getAppointmentDetail,
  submitAppointment,
  startWork,
  completeAppointment,
  cancelAppointment,
  approveAppointment,
  getApprovalHistory,
  getTimeSlotLabel,
  isEditable,
  isCancellable,
  canStart,
  canComplete,
  getAvailableWorkers
} from '@/api/appointment'

const router = useRouter()
const route = useRoute()

const appointment = ref(null)
const approvalLogs = ref([])
const loading = ref(true)
const cancelDialogVisible = ref(false)
const completeDialogVisible = ref(false)
const showApproveDialog = ref(false)

// 审批流程配置（可选，如果为空则直接使用审批记录）
const workflowConfig = []

const showRejectDialog = ref(false)
const cancelReason = ref('')
const completionNote = ref('')
const approveComment = ref('')
const rejectReason = ref('')

// 审批表单数据
const approveForm = ref({
  assignNow: false,
  workerId: null,
  reschedule: false,
  newWorkDate: '',
  newTimeSlot: ''
})

// 作业人员选择器
const showWorkerPicker = ref(false)
const selectedWorkerName = ref('')

// 日期选择器
const showDatePicker = ref(false)
const newWorkDateLabel = computed(() => {
  if (!approveForm.value.newWorkDate) return ''
  const date = new Date(approveForm.value.newWorkDate)
  return date.toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })
})

// 时间段选择器
const showTimeSlotPicker = ref(false)
const timeSlotColumns = [
  { text: '上午 (8:00-11:30)', value: 'morning' },
  { text: '中午 (12:00-13:30)', value: 'noon' },
  { text: '下午 (13:30-16:30)', value: 'afternoon' },
  { text: '全天', value: 'full_day' }
]
const newTimeSlotLabel = computed(() => {
  if (!approveForm.value.newTimeSlot) return ''
  return getTimeSlotLabel(approveForm.value.newTimeSlot)
})

// 可用作业人员数量
const availableWorkersCount = ref(0)

// DetailInfoGroup 字段配置
const basicFields = computed(() => [
  { key: 'appointment_no', label: '预约单号' },
  { key: 'status', label: '状态', type: 'status' },
  { key: 'applicant_name', label: '申请人' },
  { key: 'contact_phone', label: '联系电话', defaultValue: '-' },
  { key: 'contact_person', label: '联系人', defaultValue: '-' }
])

const workInfoFields = computed(() => [
  { key: 'work_date', label: '作业时间', type: 'custom' },
  { key: 'work_location', label: '作业地点' },
  { key: 'work_content', label: '作业内容' },
  { key: 'work_type', label: '作业类型', defaultValue: '-' }
])

const priorityFields = computed(() => [
  { key: 'is_urgent', label: '是否加急', type: 'custom' },
  { key: 'priority', label: '优先级' },
  { key: 'urgent_reason', label: '加急原因', defaultValue: '-' }
])

const assignmentFields = computed(() => [
  { key: 'assigned_worker_name', label: '作业人员', type: 'custom' }
])

const timeFields = computed(() => [
  { key: 'created_at', label: '创建时间', type: 'custom' },
  { key: 'submitted_at', label: '提交时间', type: 'custom' },
  { key: 'approved_at', label: '审批时间', type: 'custom' },
  { key: 'completed_at', label: '完成时间', type: 'custom' }
])

// 当开启分配时，显示选择提示
watch(() => approveForm.value.assignNow, async (val) => {
  if (!val) {
    selectedWorkerName.value = ''
    approveForm.value.workerId = null
  }
})

// 当开启修改时间或日期/时间改变时，加载可用作业人员数量
watch([() => approveForm.value.reschedule, () => approveForm.value.newWorkDate, () => approveForm.value.newTimeSlot], async () => {
  if (approveForm.value.reschedule && approveForm.value.newWorkDate && approveForm.value.newTimeSlot) {
    await loadAvailableWorkersForNewTime()
  }
})

const canEditOrCancel = computed(() => {
  if (!appointment.value) return false
  const status = appointment.value.status
  return status === 'draft' || status === 'pending' || status === 'scheduled'
})

// 检查是否有审批权限（不能审批自己的预约单）
const canApprove = computed(() => {
  if (!appointment.value) return false
  // 如果是草稿或已完成状态，不能审批
  if (appointment.value.status !== 'pending') return false

  // 从 localStorage 获取当前用户ID
  const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}')
  const currentUserId = userInfo.id || userInfo.user_id

  // 不能审批自己的预约单
  if (appointment.value.applicant_id === currentUserId) return false

  return true
})

onMounted(async () => {
  await loadDetail()
  await loadApprovalHistory()

  // 注册 WebSocket 消息监听
  setupWebSocketListener()

  // 注册页面可见性监听
  setupVisibilityListener()
})

// WebSocket 消息处理函数
function handleWebSocketMessage(data) {
  console.log('[WebSocket] 收到消息:', data)

  // 处理审批更新消息
  if (data.type === 'appointment_approval_update') {
    const updateData = data.data
    console.log('[WebSocket] 审批更新消息:', updateData)

    // 检查是否是当前预约单的更新
    if (updateData.appointment_id === parseInt(route.params.id)) {
      console.log('[WebSocket] 当前预约单有更新，正在刷新审批历史...')

      // 播放提示音（可选）
      playNotificationSound()

      // 显示 Toast 提示
      showSuccessToast('审批状态已更新')

      // 刷新预约单详情（状态可能已改变）
      loadDetail()

      // 刷新审批历史
      loadApprovalHistory()
    }
  }
}

// 播放提示音
function playNotificationSound() {
  try {
    // 使用 Web Audio API 生成简单的"叮"声
    if ('AudioContext' in window || 'webkitAudioContext' in window) {
      const AudioContext = window.AudioContext || window.webkitAudioContext
      const audioCtx = new AudioContext()

      // 创建振荡器（发声源）
      const oscillator = audioCtx.createOscillator()
      const gainNode = audioCtx.createGain()

      // 设置音调（800Hz）和类型（正弦波）
      oscillator.type = 'sine'
      oscillator.frequency.value = 800

      // 设置音量（开始时 0.3）
      gainNode.gain.value = 0.3

      // 连接节点
      oscillator.connect(gainNode)
      gainNode.connect(audioCtx.destination)

      // 播放 200ms 的提示音
      oscillator.start()
      setTimeout(() => {
        oscillator.stop()
        audioCtx.close()
      }, 200)

      console.log('[提示音] 播放成功')
    } else {
      console.log('[提示音] 浏览器不支持 Web Audio API')
    }
  } catch (error) {
    console.log('[提示音] 播放失败:', error)
  }
}

// 设置 WebSocket 监听
function setupWebSocketListener() {
  if (webSocketService && webSocketService.isConnected()) {
    // 保存原始的 onmessage 处理器
    const originalHandler = webSocketService.ws.onmessage

    // 监听消息
    webSocketService.ws.onmessage = (event) => {
      // 先调用原始处理（通知等）
      if (originalHandler) {
        originalHandler.call(webSocketService.ws, event)
      }

      // 然后处理我们的消息
      try {
        const data = JSON.parse(event.data)
        handleWebSocketMessage(data)
      } catch (error) {
        console.error('[WebSocket] 解析消息失败:', error)
      }
    }

    // 监听连接打开（包括重连成功）
    const originalOnOpen = webSocketService.ws.onopen
    webSocketService.ws.onopen = (event) => {
      if (originalOnOpen) {
        originalOnOpen.call(webSocketService.ws, event)
      }
      console.log('[WebSocket] 连接已建立（包括重连），刷新数据...')
      // 连接建立后刷新数据
      refreshDataIfNeeded()
    }

    console.log('[WebSocket] 已注册审批更新监听')
  } else {
    console.warn('[WebSocket] 未连接，无法注册监听')
  }
}

// 页面可见性检测
function setupVisibilityListener() {
  // 监听页面可见性变化
  document.addEventListener('visibilitychange', handleVisibilityChange)

  // 监听页面获得焦点
  window.addEventListener('focus', handlePageFocus)
}

// 处理页面可见性变化
function handleVisibilityChange() {
  if (!document.hidden && appointment.value) {
    console.log('[页面可见] 页面重新可见，刷新数据...')
    refreshDataIfNeeded()
  }
}

// 处理页面获得焦点
function handlePageFocus() {
  console.log('[页面焦点] 页面获得焦点，刷新数据...')
  refreshDataIfNeeded()
}

// 刷新数据（如果需要）
function refreshDataIfNeeded() {
  // 只在有预约单数据时刷新
  if (appointment.value && route.params.id) {
    loadDetail()
    loadApprovalHistory()
  }
}

// 清理监听器
onUnmounted(() => {
  // 移除页面可见性监听
  document.removeEventListener('visibilitychange', handleVisibilityChange)
  window.removeEventListener('focus', handlePageFocus)

  console.log('[清理] 已移除所有监听器')
})

async function loadDetail() {
  try {
    const response = await getAppointmentDetail(route.params.id)
    console.log('Appointment detail response:', response)
    appointment.value = response.data
  } catch (error) {
    console.error('加载预约单详情失败:', error)
    showFailToast('加载预约单详情失败')
  } finally {
    loading.value = false
  }
}

async function loadApprovalHistory() {
  try {
    const response = await getApprovalHistory(route.params.id)
    console.log('[审批历史] 原始响应:', response)

    // 处理不同的响应格式
    if (response && response.data) {
      approvalLogs.value = Array.isArray(response.data) ? response.data : []
      console.log('[审批历史] 更新后的记录数:', approvalLogs.value.length)
      console.log('[审批历史] 记录详情:', approvalLogs.value)
    } else if (Array.isArray(response)) {
      // 如果直接返回数组
      approvalLogs.value = response
      console.log('[审批历史] 直接数组格式，记录数:', approvalLogs.value.length)
    } else {
      approvalLogs.value = []
      console.warn('[审批历史] 未知响应格式:', response)
    }
  } catch (error) {
    console.error('[审批历史] 加载失败:', error)
    showFailToast('加载审批历史失败')
  }
}

async function handleSubmit() {
  try {
    await submitAppointment(route.params.id)
    showSuccessToast('提交成功')
    await loadDetail()
  } catch (error) {
    showFailToast(error.message || '提交失败')
  }
}

async function handleStartWork() {
  try {
    await startWork(route.params.id)
    showSuccessToast('操作成功')
    await loadDetail()
  } catch (error) {
    showFailToast(error.message || '操作失败')
  }
}

async function handleComplete() {
  try {
    await completeAppointment(route.params.id, {
      completion_note: completionNote.value
    })
    showSuccessToast('操作成功')
    completeDialogVisible.value = false
    await loadDetail()
  } catch (error) {
    showFailToast(error.message || '操作失败')
  }
}

function showCancelDialog() {
  cancelReason.value = ''
  cancelDialogVisible.value = true
}

async function handleCancel() {
  if (!cancelReason.value.trim()) {
    showFailToast('请输入取消原因')
    return
  }
  try {
    await cancelAppointment(route.params.id, {
      reason: cancelReason.value
    })
    showSuccessToast('取消成功')
    cancelDialogVisible.value = false
    await loadDetail()
  } catch (error) {
    showFailToast(error.message || '取消失败')
  }
}

// 加载新时间段的可用作业人员数量
async function loadAvailableWorkersForNewTime() {
  try {
    const response = await getAvailableWorkers({
      work_date: approveForm.value.newWorkDate,
      time_slot: approveForm.value.newTimeSlot
    })
    availableWorkersCount.value = (response.data || []).length
  } catch (error) {
    console.error('获取可用作业人员数量失败:', error)
  }
}

// 作业人员选择确认
function onWorkerConfirm(workerId, worker) {
  approveForm.value.workerId = workerId
  selectedWorkerName.value = worker ? (worker.full_name || worker.username) : ''
  showWorkerPicker.value = false
}

// 日期选择确认
function onDateConfirm(value) {
  approveForm.value.newWorkDate = value.toISOString().split('T')[0]
  showDatePicker.value = false
}

// 时间段选择确认
function onTimeSlotConfirm({ selectedOptions }) {
  approveForm.value.newTimeSlot = selectedOptions[0].value
  showTimeSlotPicker.value = false
}

// 审批通过表单提交
async function handleApproveSubmit() {
  // 验证：如果修改了时间但可用作业人员为0，阻止提交
  if (approveForm.value.reschedule && availableWorkersCount.value === 0) {
    showFailToast('所选时间段没有可用作业人员，请选择其他时间')
    return
  }

  // 验证：如果开启分配但未选择作业人员
  if (approveForm.value.assignNow && !approveForm.value.workerId) {
    showFailToast('请选择作业人员')
    return
  }

  try {
    console.log('[审批] 开始审批通过流程')
    const payload = {
      action: 'approve',
      comment: approveComment.value,
      assign_now: approveForm.value.assignNow,
      worker_id: approveForm.value.workerId,
      reschedule: approveForm.value.reschedule,
      new_work_date: approveForm.value.newWorkDate,
      new_time_slot: approveForm.value.newTimeSlot
    }
    const result = await approveAppointment(route.params.id, payload)
    console.log('[审批] 审批响应:', result)

    showSuccessToast('审批通过')
    showApproveDialog.value = false

    // 重置表单
    approveComment.value = ''
    approveForm.value = {
      assignNow: false,
      workerId: null,
      reschedule: false,
      newWorkDate: '',
      newTimeSlot: ''
    }
    availableWorkersCount.value = 0

    // 重新加载数据
    await loadDetail()
    await loadApprovalHistory()
  } catch (error) {
    console.error('[审批] 审批失败:', error)
    showFailToast(error.message || '审批失败')
  }
}

// 审批拒绝
async function handleReject() {
  if (!rejectReason.value.trim()) {
    showFailToast('请输入拒绝原因')
    return
  }
  try {
    await approveAppointment(route.params.id, {
      action: 'reject',
      comment: rejectReason.value
    })

    showSuccessToast('已拒绝')
    showRejectDialog.value = false
    rejectReason.value = ''

    // 重新加载数据
    await loadDetail()
    await loadApprovalHistory()
  } catch (error) {
    console.error('[审批] 拒绝失败:', error)
    showFailToast(error.message || '拒绝失败')
  }
}

function showCompleteDialog() {
  completionNote.value = ''
  completeDialogVisible.value = true
}

function formatFullDateTime(dateStr) {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}
</script>

<style scoped>
.appointment-detail {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.detail-content {
  padding-bottom: 16px;
}

.van-cell-group {
  margin-bottom: 12px;
}

.action-buttons {
  padding: 16px;
}

.action-buttons .van-button {
  margin-bottom: 8px;
}

.approve-buttons {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.dialog-actions {
  display: flex;
  gap: 12px;
  padding: 16px;
}

.dialog-actions .van-button {
  flex: 1;
}
</style>
