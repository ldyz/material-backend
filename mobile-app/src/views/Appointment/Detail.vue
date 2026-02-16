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
      <van-cell-group title="基本信息" inset>
        <van-cell title="预约单号" :value="appointment.appointment_no" />
        <van-cell title="状态">
          <template #value>
            <van-tag :type="getStatusColor(appointment.status)">
              {{ getStatusLabel(appointment.status) }}
            </van-tag>
          </template>
        </van-cell>
        <van-cell title="申请人" :value="appointment.applicant_name" />
        <van-cell title="联系电话" :value="appointment.contact_phone || '-'" />
        <van-cell title="联系人" :value="appointment.contact_person || '-'" />
      </van-cell-group>

      <!-- 作业信息 -->
      <van-cell-group title="作业信息" inset>
        <van-cell title="作业时间" :value="formatDateTime(appointment.work_date, appointment.time_slot)" />
        <van-cell title="作业地点" :value="appointment.work_location" />
        <van-cell title="作业内容" :value="appointment.work_content" />
        <van-cell title="作业类型" :value="appointment.work_type || '-'" />
      </van-cell-group>

      <!-- 优先级信息 -->
      <van-cell-group v-if="appointment.is_urgent || appointment.priority > 0" title="优先级" inset>
        <van-cell title="是否加急">
          <template #value>
            <van-tag v-if="appointment.is_urgent" type="danger">是</van-tag>
            <van-tag v-else type="default">否</van-tag>
          </template>
        </van-cell>
        <van-cell title="优先级" :value="appointment.priority" />
        <van-cell v-if="appointment.urgent_reason" title="加急原因" :value="appointment.urgent_reason" />
      </van-cell-group>

      <!-- 分配信息 -->
      <van-cell-group title="分配信息" inset>
        <van-cell title="作业人员" :value="appointment.assigned_worker_name || '未分配'" />
      </van-cell-group>

      <!-- 时间信息 -->
      <van-cell-group title="时间信息" inset>
        <van-cell title="创建时间" :value="formatFullDateTime(appointment.created_at)" />
        <van-cell v-if="appointment.submitted_at" title="提交时间" :value="formatFullDateTime(appointment.submitted_at)" />
        <van-cell v-if="appointment.approved_at" title="审批时间" :value="formatFullDateTime(appointment.approved_at)" />
        <van-cell v-if="appointment.completed_at" title="完成时间" :value="formatFullDateTime(appointment.completed_at)" />
      </van-cell-group>

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
      show-cancel-button
      @confirm="handleApprove"
    >
      <van-field
        v-model="approveComment"
        type="textarea"
        placeholder="请输入审批意见（可选）"
        rows="3"
        autosize
      />
    </van-dialog>

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
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showSuccessToast, showFailToast, Dialog } from 'vant'
import ApprovalTimeline from '@/components/common/ApprovalTimeline.vue'
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
  getStatusLabel,
  getStatusColor,
  isEditable,
  isCancellable,
  canStart,
  canComplete
} from '@/api/appointment'

const router = useRouter()
const route = useRoute()

const appointment = ref(null)
const approvalLogs = ref([])
const loading = ref(true)
const cancelDialogVisible = ref(false)
const completeDialogVisible = ref(false)
const showApproveDialog = ref(false)

// 审批流程配置（根据实际业务配置）
const workflowConfig = [
  {
    role: 'project_manager',
    title: '项目经理审批',
    order: 1,
    role_name: '项目经理'
  },
  {
    role: 'supervisor',
    title: '主管审批',
    order: 2,
    role_name: '主管'
  }
]
const showRejectDialog = ref(false)
const cancelReason = ref('')
const completionNote = ref('')
const approveComment = ref('')
const rejectReason = ref('')

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
    // response 直接就是 { success: true, data: {...}, meta: {...} }
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

// 审批通过
async function handleApprove() {
  try {
    console.log('[审批] 开始审批通过流程')
    const result = await approveAppointment(route.params.id, {
      action: 'approve',
      comment: approveComment.value
    })
    console.log('[审批] 审批响应:', result)

    showSuccessToast('审批通过')
    showApproveDialog.value = false
    approveComment.value = ''

    // 重新加载数据
    console.log('[审批] 重新加载详情...')
    await loadDetail()

    console.log('[审批] 重新加载审批历史...')
    await loadApprovalHistory()

    console.log('[审批] 所有数据更新完成')
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
    console.log('[审批] 开始审批拒绝流程')
    const result = await approveAppointment(route.params.id, {
      action: 'reject',
      comment: rejectReason.value
    })
    console.log('[审批] 拒绝响应:', result)

    showSuccessToast('已拒绝')
    showRejectDialog.value = false
    rejectReason.value = ''

    // 重新加载数据
    console.log('[审批] 重新加载详情...')
    await loadDetail()

    console.log('[审批] 重新加载审批历史...')
    await loadApprovalHistory()

    console.log('[审批] 所有数据更新完成')
  } catch (error) {
    console.error('[审批] 拒绝失败:', error)
    showFailToast(error.message || '拒绝失败')
  }
}

function showCompleteDialog() {
  completionNote.value = ''
  completeDialogVisible.value = true
}

function formatDateTime(dateStr, timeSlot) {
  const date = new Date(dateStr)
  const dateStr2 = date.toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })
  const slot = getTimeSlotLabel(timeSlot)
  return `${dateStr2} ${slot}`
}

function formatFullDateTime(dateStr) {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

// 获取审批动作的中文标签
function getActionLabel(action) {
  const actionLabels = {
    'start': '提交',
    'approve': '通过',
    'reject': '拒绝',
    'return': '退回',
    'comment': '评论',
    'cancel': '取消',
    'submit': '提交',
    'resubmit': '重新提交'
  }
  return actionLabels[action] || action
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

.approval-remark {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.remark-text {
  color: #646566;
  font-size: 12px;
  margin-top: 4px;
}
</style>
