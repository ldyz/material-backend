<template>
  <el-dialog
    v-model="dialogVisible"
    title="预约单详情"
    width="900px"
    @close="handleClose"
  >
    <el-descriptions :column="2" border v-if="appointment">
      <el-descriptions-item label="预约单号">
        {{ appointment.appointment_no }}
      </el-descriptions-item>
      <el-descriptions-item label="状态">
        <el-tag :type="getStatusType(appointment.status)">
          {{ getStatusLabel(appointment.status) }}
        </el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="申请人">
        {{ appointment.applicant_name }}
      </el-descriptions-item>
      <el-descriptions-item label="联系电话">
        {{ appointment.contact_phone || '-' }}
      </el-descriptions-item>
      <el-descriptions-item label="作业时间" :span="2">
        {{ formatDateTime(appointment.work_date, appointment.time_slot) }}
      </el-descriptions-item>
      <el-descriptions-item label="作业地点" :span="2">
        {{ appointment.work_location }}
      </el-descriptions-item>
      <el-descriptions-item label="作业内容" :span="2">
        {{ appointment.work_content }}
      </el-descriptions-item>
      <el-descriptions-item label="作业类型" v-if="appointment.work_type">
        {{ appointment.work_type }}
      </el-descriptions-item>
      <el-descriptions-item label="作业人员">
        {{ appointment.assigned_worker_name || '未分配' }}
      </el-descriptions-item>
      <el-descriptions-item label="是否加急">
        <el-tag v-if="appointment.is_urgent" type="danger" size="small">是</el-tag>
        <el-tag v-else type="info" size="small">否</el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="优先级" v-if="appointment.is_urgent">
        {{ appointment.priority }}
      </el-descriptions-item>
      <el-descriptions-item label="加急原因" v-if="appointment.urgent_reason" :span="2">
        {{ appointment.urgent_reason }}
      </el-descriptions-item>
      <el-descriptions-item label="创建时间">
        {{ formatFullDateTime(appointment.created_at) }}
      </el-descriptions-item>
      <el-descriptions-item label="提交时间" v-if="appointment.submitted_at">
        {{ formatFullDateTime(appointment.submitted_at) }}
      </el-descriptions-item>
      <el-descriptions-item label="审批时间" v-if="appointment.approved_at">
        {{ formatFullDateTime(appointment.approved_at) }}
      </el-descriptions-item>
      <el-descriptions-item label="完成时间" v-if="appointment.completed_at">
        {{ formatFullDateTime(appointment.completed_at) }}
      </el-descriptions-item>
    </el-descriptions>

    <!-- 审批历史 -->
    <div class="approval-history" v-if="approvalLogs.length > 0">
      <h4>审批历史</h4>
      <el-timeline>
        <el-timeline-item
          v-for="(log, index) in approvalLogs"
          :key="index"
          :timestamp="formatFullDateTime(log.created_at)"
        >
          <div class="timeline-content">
            <div class="timeline-title">{{ log.node_name }}</div>
            <div class="timeline-action">
              <el-tag :type="getActionType(log.action)" size="small">
                {{ getActionText(log.action) }}
              </el-tag>
              <span class="approver">{{ log.approver_name }}</span>
            </div>
            <div v-if="log.remark" class="timeline-remark">{{ log.remark }}</div>
          </div>
        </el-timeline-item>
      </el-timeline>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button
          v-if="canApprove()"
          type="success"
          :icon="Check"
          @click="handleApprove"
        >
          审批
        </el-button>
        <el-button @click="handleClose">关闭</el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { Check } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { appointmentApi } from '@/api'

const props = defineProps({
  modelValue: Boolean,
  appointmentId: [Number, String]
})

const emit = defineEmits(['update:modelValue', 'approve'])

const authStore = useAuthStore()

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const appointment = ref(null)
const approvalLogs = ref([])

watch(() => props.appointmentId, (id) => {
  if (id && props.modelValue) {
    loadDetail()
  }
}, { immediate: true })

watch(() => props.modelValue, (val) => {
  if (val && props.appointmentId) {
    loadDetail()
  }
})

async function loadDetail() {
  try {
    const response = await appointmentApi.getDetail(props.appointmentId)
    appointment.value = response.data

    // 加载审批历史
    const historyRes = await appointmentApi.getApprovalHistory(props.appointmentId)
    approvalLogs.value = historyRes.data || []
  } catch (error) {
    console.error('加载详情失败:', error)
  }
}

function handleClose() {
  emit('update:modelValue', false)
}

function formatDateTime(dateStr, timeSlot) {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  const dateStr2 = date.toLocaleDateString('zh-CN')
  const slots = { morning: '上午', afternoon: '下午', evening: '晚上', full_day: '全天' }
  return `${dateStr2} ${slots[timeSlot] || timeSlot}`
}

function formatFullDateTime(dateStr) {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

function getStatusType(status) {
  return appointmentApi.getStatusType(status)
}

function getStatusLabel(status) {
  return appointmentApi.getStatusLabel(status)
}

function getActionType(action) {
  const types = {
    approve: 'success',
    reject: 'danger',
    comment: 'info'
  }
  return types[action] || ''
}

function getActionText(action) {
  const texts = {
    approve: '通过',
    reject: '拒绝',
    comment: '评论'
  }
  return texts[action] || action
}

function canApprove() {
  if (!appointment.value) return false
  return appointment.value.status === 'pending' && authStore.hasPermission('appointment_approve')
}

function handleApprove() {
  emit('approve', appointment.value)
}

</script>

<style scoped>
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

.approval-history {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #eee;
}

.approval-history h4 {
  margin-bottom: 16px;
  font-size: 16px;
  font-weight: 500;
}

.timeline-content {
  padding-bottom: 8px;
}

.timeline-title {
  font-weight: 500;
  margin-bottom: 4px;
}

.timeline-action {
  display: flex;
  align-items: center;
  gap: 8px;
}

.timeline-remark {
  margin-top: 4px;
  color: #666;
  font-size: 13px;
}
</style>
