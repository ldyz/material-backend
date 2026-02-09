<template>
  <el-dialog
    v-model="dialogVisible"
    title="分配作业人员"
    width="600px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <div v-if="appointment">
      <el-alert
        title="预约信息"
        type="info"
        :closable="false"
        style="margin-bottom: 20px"
      >
        <p><strong>预约单号：</strong>{{ appointment.appointment_no }}</p>
        <p><strong>作业时间：</strong>{{ formatDateTime(appointment.work_date, appointment.time_slot) }}</p>
        <p><strong>作业地点：</strong>{{ appointment.work_location }}</p>
      </el-alert>

      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item label="作业人员" prop="worker_id" required>
          <el-select
            v-model="formData.worker_id"
            filterable
            placeholder="选择作业人员"
            style="width: 100%"
          >
            <el-option
              v-for="worker in availableWorkers"
              :key="worker.id"
              :label="worker.name"
              :value="worker.id"
            >
              <span>{{ worker.name }}</span>
              <span style="color: #999; font-size: 12px; margin-left: 8px">
                (ID: {{ worker.id }})
              </span>
            </el-option>
          </el-select>
          <div style="margin-top: 8px; color: #999; font-size: 12px">
            只显示该时间段可用的作业人员
          </div>
        </el-form-item>
      </el-form>
    </div>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" :loading="submitting" @click="handleSubmit">
        确认分配
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { appointmentApi } from '@/api'

const props = defineProps({
  modelValue: Boolean,
  appointment: Object
})

const emit = defineEmits(['update:modelValue', 'success'])

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const formRef = ref(null)
const submitting = ref(false)
const availableWorkers = ref([])

const formData = ref({
  worker_id: null
})

const formRules = {
  worker_id: [{ required: true, message: '请选择作业人员', trigger: 'change' }]
}

watch(() => props.modelValue, async (val) => {
  if (val && props.appointment) {
    await loadAvailableWorkers()
  }
})

async function loadAvailableWorkers() {
  if (!props.appointment) return
  try {
    const { data } = await appointmentApi.getAvailableWorkers({
      work_date: props.appointment.work_date,
      time_slot: props.appointment.time_slot
    })
    availableWorkers.value = data.data || []

    // 如果已经分配了作业人员，且在可用列表中，则选中
    if (props.appointment.assigned_worker_id) {
      const exists = availableWorkers.value.some(w => w.id === props.appointment.assigned_worker_id)
      if (exists) {
        formData.value.worker_id = props.appointment.assigned_worker_id
      }
    }
  } catch (error) {
    ElMessage.error('获取可用作业人员失败')
  }
}

async function handleSubmit() {
  try {
    await formRef.value.validate()
    submitting.value = true

    await appointmentApi.assignWorker(props.appointment.id, {
      worker_id: formData.value.worker_id
    })

    ElMessage.success('分配成功')
    emit('success')
    handleClose()
  } catch (error) {
    ElMessage.error(error.message || '分配失败')
  } finally {
    submitting.value = false
  }
}

function handleClose() {
  formData.value = { worker_id: null }
  formRef.value?.clearValidate()
  emit('update:modelValue', false)
}

function formatDateTime(dateStr, timeSlot) {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  const dateStr2 = date.toLocaleDateString('zh-CN')
  const slots = { morning: '上午', afternoon: '下午', evening: '晚上', full_day: '全天' }
  return `${dateStr2} ${slots[timeSlot] || timeSlot}`
}
</script>
