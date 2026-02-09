<template>
  <el-dialog
    v-model="dialogVisible"
    :title="dialogTitle"
    width="800px"
    :close-on-click-modal="false"
    @close="handleClose"
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
            <el-input v-model="formData.work_type" placeholder="请输入作业类型" />
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
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="时间段" prop="time_slot">
            <el-select v-model="formData.time_slot" placeholder="请选择" style="width: 100%">
              <el-option label="上午 (08:00-12:00)" value="morning" />
              <el-option label="下午 (14:00-18:00)" value="afternoon" />
              <el-option label="晚上 (19:00-22:00)" value="evening" />
              <el-option label="全天" value="full_day" />
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
        <el-select
          v-model="formData.assigned_worker_id"
          placeholder="选择作业人员（可选）"
          filterable
          style="width: 100%"
        >
          <el-option
            v-for="worker in workerList"
            :key="worker.id"
            :label="worker.name"
            :value="worker.id"
          />
        </el-select>
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" :loading="submitting" @click="handleSubmit">
        {{ mode === 'create' ? '创建' : '保存' }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useProjectStore } from '@/stores/projectStore'
import { appointmentApi } from '@/api'
import ProjectSelector from '@/components/common/ProjectSelector.vue'

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
const projectList = ref([])

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
  assigned_worker_id: null
})

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

const projectStore = useProjectStore()

// TODO: 获取作业人员列表
const workerList = ref([])

watch(() => props.appointment, (val) => {
  if (val && props.mode === 'edit') {
    Object.assign(formData.value, val)
    workDate.value = val.work_date
  } else {
    resetForm()
  }
}, { immediate: true })

watch(workDate, (val) => {
  formData.value.work_date = val
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
    assigned_worker_id: null
  }
  workDate.value = ''
  formRef.value?.clearValidate()
}

async function handleSubmit() {
  try {
    await formRef.value.validate()
    submitting.value = true

    if (props.mode === 'create') {
      await appointmentApi.create(formData.value)
      ElMessage.success('创建成功')
    } else {
      await appointmentApi.update(props.appointment.id, formData.value)
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

function handleClose() {
  resetForm()
  emit('update:modelValue', false)
}

// 初始化
projectStore.fetchProjects().then(() => {
  projectList.value = projectStore.projects
})
</script>
