<template>
  <el-dialog
    v-model="dialogVisible"
    title="审批预约单"
    width="700px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <div v-if="loading" class="loading-wrapper">
      <el-icon class="is-loading"><Loading /></el-icon>
      <span>加载中...</span>
    </div>

    <div v-else-if="appointment">
      <!-- 预约单信息 -->
      <el-card class="appointment-info" shadow="never">
        <template #header>
          <span>预约单信息</span>
        </template>
        <el-descriptions :column="2" border size="small">
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
          <el-descriptions-item label="作业时间">
            {{ formatDateTime(appointment.work_date, appointment.time_slot) }}
          </el-descriptions-item>
          <el-descriptions-item label="作业地点" :span="2">
            {{ appointment.work_location }}
          </el-descriptions-item>
          <el-descriptions-item label="作业内容" :span="2">
            {{ appointment.work_content }}
          </el-descriptions-item>
          <el-descriptions-item label="是否加急">
            <el-tag v-if="appointment.is_urgent" type="danger" size="small">是</el-tag>
            <el-tag v-else type="info" size="small">否</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="优先级" v-if="appointment.is_urgent">
            {{ appointment.priority }}
          </el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 当前审批信息 -->
      <el-card v-if="currentApproval" class="current-approval" shadow="never">
        <template #header>
          <span>当前审批节点</span>
        </template>
        <el-form label-width="100px">
          <el-form-item label="审批节点">
            {{ currentApproval.node_name }}
          </el-form-item>
          <el-form-item label="审批人">
            <el-tag
              v-for="approver in currentApproval.approvers"
              :key="approver.id"
              type="primary"
              size="small"
              style="margin-right: 4px"
            >
              {{ approver.name }}
            </el-tag>
          </el-form-item>
        </el-form>
      </el-card>

      <!-- 分配作业人员 -->
      <el-card class="assign-worker" shadow="never">
        <template #header>
          <span>分配作业人员</span>
        </template>
        <el-form label-width="120px">
          <el-form-item label="立即分配">
            <el-switch v-model="assignNow" />
            <span style="margin-left: 8px; color: #999; font-size: 12px">
              审批通过后自动分配
            </span>
          </el-form-item>
          <el-form-item v-if="assignNow" label="作业人员" required>
            <el-select
              v-model="selectedWorkerId"
              filterable
              placeholder="选择作业人员"
              style="width: 100%"
            >
              <el-option
                v-for="worker in availableWorkers"
                :key="worker.id"
                :label="worker.name"
                :value="worker.id"
              />
            </el-select>
          </el-form-item>
        </el-form>
      </el-card>

      <!-- 修改作业时间 -->
      <el-card v-if="appointment" class="reschedule-time" shadow="never">
        <template #header>
          <span>修改作业时间</span>
        </template>
        <el-form label-width="120px">
          <el-form-item label="调整时间">
            <el-switch v-model="reschedule" />
            <span style="margin-left: 8px; color: #999; font-size: 12px">
              审批通过后使用新时间
            </span>
          </el-form-item>

          <template v-if="reschedule">
            <el-form-item label="新作业日期" required>
              <el-date-picker
                v-model="newWorkDate"
                type="date"
                placeholder="选择日期"
                format="YYYY-MM-DD"
                value-format="YYYY-MM-DD"
                :disabled-date="disablePastDates"
                @change="handleDateChange"
              />
            </el-form-item>

            <el-form-item label="新时间段" required>
              <el-select
                v-model="newTimeSlot"
                placeholder="选择时间段"
                @change="handleTimeSlotChange"
              >
                <el-option label="上午 (8:00-11:30)" value="morning" />
                <el-option label="中午 (12:00-13:30)" value="noon" />
                <el-option label="下午 (13:30-16:30)" value="afternoon" />
                <el-option label="全天" value="full_day" />
              </el-select>
            </el-form-item>

            <el-form-item label="可用作业人员">
              <el-tag :type="availableWorkersCount > 0 ? 'success' : 'danger'">
                {{ availableWorkersCount }} 人可用
              </el-tag>
            </el-form-item>
          </template>
        </el-form>
      </el-card>

      <!-- 审批表单 -->
      <el-form ref="formRef" :model="formData" label-width="100px">
        <el-form-item label="审批意见">
          <el-input
            v-model="formData.comment"
            type="textarea"
            :rows="3"
            placeholder="请输入审批意见"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>
      </el-form>
    </div>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="danger" :loading="submitting" @click="handleReject">
        拒绝
      </el-button>
      <el-button type="success" :loading="submitting" @click="handleApprove">
        同意
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Loading } from '@element-plus/icons-vue'
import { appointmentApi } from '@/api'

const props = defineProps({
  modelValue: Boolean,
  appointmentId: [Number, String]
})

const emit = defineEmits(['update:modelValue', 'success'])

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const loading = ref(false)
const submitting = ref(false)
const appointment = ref(null)
const currentApproval = ref(null)
const availableWorkers = ref([])
const assignNow = ref(false)
const selectedWorkerId = ref(null)

// 修改作业时间相关
const reschedule = ref(false)
const newWorkDate = ref('')
const newTimeSlot = ref('')
const availableWorkersCount = ref(0)

const formData = ref({
  comment: ''
})

const formRef = ref(null)

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

watch(assignNow, async (val) => {
  if (val && appointment.value) {
    await loadAvailableWorkers()
  } else {
    availableWorkers.value = []
    selectedWorkerId.value = null
  }
})

// 当开启修改时间或日期/时间改变时，加载可用作业人员数量
watch([reschedule, newWorkDate, newTimeSlot], async () => {
  if (reschedule.value && newWorkDate.value && newTimeSlot.value) {
    await loadAvailableWorkersForNewTime()
  }
})

async function loadDetail() {
  loading.value = true
  try {
    const { data } = await appointmentApi.getDetail(props.appointmentId)
    appointment.value = data.data

    // 加载当前审批节点
    const approvalRes = await appointmentApi.getCurrentApproval(props.appointmentId)
    currentApproval.value = approvalRes.data.data
  } catch (error) {
    ElMessage.error('加载详情失败')
  } finally {
    loading.value = false
  }
}

async function loadAvailableWorkers() {
  if (!appointment.value) return
  try {
    const { data } = await appointmentApi.getAvailableWorkers({
      work_date: appointment.value.work_date,
      time_slot: appointment.value.time_slot
    })
    availableWorkers.value = data.data || []
  } catch (error) {
    console.error('获取可用作业人员失败:', error)
  }
}

async function loadAvailableWorkersForNewTime() {
  try {
    const { data } = await appointmentApi.getAvailableWorkers({
      work_date: newWorkDate.value,
      time_slot: newTimeSlot.value
    })
    availableWorkersCount.value = (data.data || []).length
  } catch (error) {
    console.error('获取可用作业人员数量失败:', error)
  }
}

// 禁用过去的日期
function disablePastDates(time) {
  return time.getTime() < Date.now() - 24 * 60 * 60 * 1000
}

function handleDateChange() {
  // 日期改变时重新计算可用作业人员数量
  if (newTimeSlot.value) {
    loadAvailableWorkersForNewTime()
  }
}

function handleTimeSlotChange() {
  // 时间段改变时重新计算可用作业人员数量
  if (newWorkDate.value) {
    loadAvailableWorkersForNewTime()
  }
}

async function handleApprove() {
  try {
    // 验证：如果修改了时间但可用作业人员为0，阻止提交
    if (reschedule.value && availableWorkersCount.value === 0) {
      ElMessage.error('所选时间段没有可用作业人员，请选择其他时间')
      return
    }

    // 验证：如果开启分配但未选择作业人员
    if (assignNow.value && !selectedWorkerId.value) {
      ElMessage.warning('请选择作业人员')
      return
    }

    await ElMessageBox.confirm('确认同意此预约单？', '提示', {
      type: 'success'
    })

    submitting.value = true
    await appointmentApi.approve(props.appointmentId, {
      action: 'approve',
      comment: formData.value.comment,
      assign_now: assignNow.value,
      worker_id: selectedWorkerId.value,
      reschedule: reschedule.value,
      new_work_date: newWorkDate.value,
      new_time_slot: newTimeSlot.value
    })

    ElMessage.success('审批成功')
    emit('success')
    handleClose()
  } catch (error) {
      if (error !== 'cancel') {
        ElMessage.error(error.message || '审批失败')
      }
  } finally {
    submitting.value = false
  }
}

async function handleReject() {
  try {
    await ElMessageBox.confirm('确认拒绝此预约单？', '提示', {
      type: 'warning'
    })

    submitting.value = true
    await appointmentApi.approve(props.appointmentId, {
      action: 'reject',
      comment: formData.value.comment
    })

    ElMessage.success('已拒绝')
    emit('success')
    handleClose()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '操作失败')
    }
  } finally {
    submitting.value = false
  }
}

function handleClose() {
  emit('update:modelValue', false)
  // 重置表单
  formData.value = { comment: '' }
  assignNow.value = false
  selectedWorkerId.value = null
  availableWorkers.value = []
  reschedule.value = false
  newWorkDate.value = ''
  newTimeSlot.value = ''
  availableWorkersCount.value = 0
}

function formatDateTime(dateStr, timeSlot) {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  const dateStr2 = date.toLocaleDateString('zh-CN')
  const slots = { morning: '上午', afternoon: '下午', evening: '晚上', full_day: '全天' }
  return `${dateStr2} ${slots[timeSlot] || timeSlot}`
}

function getStatusType(status) {
  return appointmentApi.getStatusType(status)
}

function getStatusLabel(status) {
  return appointmentApi.getStatusLabel(status)
}
</script>

<style scoped>
.loading-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 0;
  color: #999;
}

.loading-wrapper .el-icon {
  font-size: 32px;
  margin-bottom: 16px;
}

.appointment-info,
.current-approval,
.assign-worker,
.reschedule-time {
  margin-bottom: 16px;
}

.assign-worker :deep(.el-form-item__content) {
  color: #666;
  font-size: 13px;
}
</style>
