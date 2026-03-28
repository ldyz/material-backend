<template>
  <div class="appointment-approve">
    <van-nav-bar
      title="审批预约单"
      left-arrow
      @click-left="router.back()"
    />

    <van-loading v-if="loading" type="spinner" color="#1989fa" />

    <div v-else-if="appointment" class="approve-content">
      <!-- 预约单信息 -->
      <van-cell-group title="预约单信息" inset>
        <van-cell title="预约单号" :value="appointment.appointment_no" />
        <van-cell title="申请人" :value="appointment.applicant_name" />
        <van-cell title="作业时间" :value="formatDateTime(appointment.work_date, appointment.time_slot)" />
        <van-cell title="作业地点" :value="appointment.work_location" />
        <van-cell title="作业内容" :value="appointment.work_content" />
        <van-cell v-if="appointment.is_urgent" title="优先级">
          <template #value>
            <van-tag type="danger">加急 ({{ appointment.priority }})</van-tag>
          </template>
        </van-cell>
        <van-cell v-if="appointment.urgent_reason" title="加急原因" :value="appointment.urgent_reason" />
      </van-cell-group>

      <!-- 当前审批节点 -->
      <van-cell-group v-if="currentApproval" title="当前审批" inset>
        <van-cell title="审批节点" :value="currentApproval.node_name" />
        <van-cell title="审批人">
          <template #value>
            <van-tag v-for="approver in currentApproval.approvers" :key="approver.id" type="primary" size="small">
              {{ approver.name }}
            </van-tag>
          </template>
        </van-cell>
      </van-cell-group>

      <!-- 审批历史 -->
      <van-cell-group title="审批历史" inset>
        <van-cell
          v-for="(log, index) in approvalLogs"
          :key="index"
          :title="log.node_name"
          :value="getActionText(log.action)"
          :label="`${log.approver_name} - ${formatFullDateTime(log.created_at)}`"
        />
        <van-empty v-if="!approvalLogs.length" description="暂无审批记录" />
      </van-cell-group>

      <!-- 分配作业人员 -->
      <van-cell-group v-if="canAssignWorker" title="分配作业人员" inset>
        <van-field
          name="worker_id"
          label="作业人员"
          :value="workerName"
          readonly
          is-link
          placeholder="选择作业人员"
          @click="showWorkerPicker = true"
        />
        <van-field name="assign_now" label="立即分配">
          <template #input>
            <van-switch v-model="assignNow" />
          </template>
        </van-field>
      </van-cell-group>

      <!-- 审批操作 -->
      <van-cell-group title="审批操作" inset>
        <van-field
          v-model="comment"
          name="comment"
          label="审批意见"
          type="textarea"
          placeholder="请输入审批意见"
          rows="3"
        />
      </van-cell-group>

      <!-- 操作按钮 -->
      <div class="action-buttons">
        <van-button
          type="success"
          block
          @click="handleApprove"
          :loading="submitting"
        >
          同意
        </van-button>
        <van-button
          type="danger"
          block
          @click="handleReject"
          :loading="submitting"
        >
          拒绝
        </van-button>
      </div>
    </div>

    <van-empty v-else description="预约单不存在" />

    <!-- 作业人员选择器 -->
    <WorkerPicker
      v-model="showWorkerPicker"
      title="选择作业人员"
      :multiple="false"
      @confirm="onWorkerConfirm"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showSuccessToast, showFailToast, Dialog } from 'vant'
import WorkerPicker from '@/components/common/WorkerPicker.vue'
import {
  getAppointmentDetail,
  approveAppointment,
  getApprovalHistory,
  getCurrentApproval,
  getAvailableWorkers,
  getTimeSlotLabel
} from '@/api/appointment'

const router = useRouter()
const route = useRoute()

const appointment = ref(null)
const approvalLogs = ref([])
const currentApproval = ref(null)
const loading = ref(true)
const submitting = ref(false)
const comment = ref('')
const assignNow = ref(false)
const selectedWorkerId = ref(null)
const workerName = ref('')
const showWorkerPicker = ref(false)

const canAssignWorker = computed(() => {
  return appointment.value && appointment.value.status === 'pending'
})

onMounted(async () => {
  await loadDetail()
  await loadApprovalHistory()
  await loadCurrentApproval()
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
    console.log('Approval history response:', response)
    approvalLogs.value = response.data || []
  } catch (error) {
    console.error('加载审批历史失败:', error)
  }
}

async function loadCurrentApproval() {
  try {
    const response = await getCurrentApproval(route.params.id)
    console.log('Current approval response:', response)
    currentApproval.value = response.data
  } catch (error) {
    console.error('加载当前审批节点失败:', error)
  }
}

// 作业人员选择确认
function onWorkerConfirm(workerId, worker) {
  selectedWorkerId.value = workerId
  workerName.value = worker ? (worker.full_name || worker.username) : ''
  showWorkerPicker.value = false
}

async function handleApprove() {
  if (assignNow.value && !selectedWorkerId.value) {
    showFailToast('请选择作业人员')
    return
  }

  try {
    submitting.value = true
    await approveAppointment(route.params.id, {
      action: 'approve',
      comment: comment.value,
      assign_now: assignNow.value,
      worker_id: selectedWorkerId.value
    })
    showSuccessToast('审批成功')
    router.back()
  } catch (error) {
    showFailToast(error.message || '审批失败')
  } finally {
    submitting.value = false
  }
}

async function handleReject() {
  try {
    submitting.value = true
    await approveAppointment(route.params.id, {
      action: 'reject',
      comment: comment.value
    })
    showSuccessToast('已拒绝')
    router.back()
  } catch (error) {
    showFailToast(error.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

function getActionText(action) {
  const map = {
    approve: '通过',
    reject: '拒绝',
    start: '开始',
    comment: '评论'
  }
  return map[action] || action
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
.appointment-approve {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 80px;
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
