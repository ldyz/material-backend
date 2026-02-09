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
        <van-cell
          v-for="(log, index) in approvalLogs"
          :key="index"
          :title="log.node_name"
          :value="log.action === 'approve' ? '通过' : log.action === 'reject' ? '拒绝' : log.action"
          :label="`${log.approver_name} - ${formatFullDateTime(log.created_at)}`"
        />
        <van-empty v-if="!approvalLogs.length" description="暂无审批记录" />
      </van-cell-group>

      <!-- 操作按钮 -->
      <div class="action-buttons">
        <van-button
          v-if="appointment.status === 'draft'"
          type="primary"
          block
          @click="handleSubmit"
        >
          提交审批
        </van-button>
      </div>
    </div>

    <van-empty v-else description="预约单不存在" />

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
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Toast, Dialog } from 'vant'
import {
  getAppointmentDetail,
  submitAppointment,
  startWork,
  completeAppointment,
  cancelAppointment,
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
const cancelReason = ref('')
const completionNote = ref('')

const canEditOrCancel = computed(() => {
  if (!appointment.value) return false
  const status = appointment.value.status
  return status === 'draft' || status === 'pending' || status === 'scheduled'
})

onMounted(async () => {
  await loadDetail()
  await loadApprovalHistory()
})

async function loadDetail() {
  try {
    const { data } = await getAppointmentDetail(route.params.id)
    appointment.value = data.data
  } catch (error) {
    Toast.fail('加载预约单详情失败')
  } finally {
    loading.value = false
  }
}

async function loadApprovalHistory() {
  try {
    const { data } = await getApprovalHistory(route.params.id)
    approvalLogs.value = data.data || []
  } catch (error) {
    console.error('加载审批历史失败:', error)
  }
}

async function handleSubmit() {
  try {
    await submitAppointment(route.params.id)
    Toast.success('提交成功')
    await loadDetail()
  } catch (error) {
    Toast.fail(error.message || '提交失败')
  }
}

async function handleStartWork() {
  try {
    await startWork(route.params.id)
    Toast.success('操作成功')
    await loadDetail()
  } catch (error) {
    Toast.fail(error.message || '操作失败')
  }
}

async function handleComplete() {
  try {
    await completeAppointment(route.params.id, {
      completion_note: completionNote.value
    })
    Toast.success('操作成功')
    completeDialogVisible.value = false
    await loadDetail()
  } catch (error) {
    Toast.fail(error.message || '操作失败')
  }
}

function showCancelDialog() {
  cancelReason.value = ''
  cancelDialogVisible.value = true
}

async function handleCancel() {
  if (!cancelReason.value.trim()) {
    Toast.fail('请输入取消原因')
    return
  }
  try {
    await cancelAppointment(route.params.id, {
      reason: cancelReason.value
    })
    Toast.success('取消成功')
    cancelDialogVisible.value = false
    await loadDetail()
  } catch (error) {
    Toast.fail(error.message || '取消失败')
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
</style>
