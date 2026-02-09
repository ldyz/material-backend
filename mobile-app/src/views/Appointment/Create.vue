<template>
  <div class="appointment-create">
    <van-nav-bar
      title="创建预约单"
      left-arrow
      @click-left="router.back()"
    />

    <van-form @submit="handleSubmit">
      <!-- 基本信息 -->
      <van-cell-group title="基本信息" inset>
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
          placeholder="请输入作业类型"
        />
      </van-cell-group>

      <!-- 时间信息 -->
      <van-cell-group title="时间信息" inset>
        <van-field
          name="work_date"
          label="作业日期"
          :value="form.work_date"
          readonly
          is-link
          @click="showDatePicker = true"
          :rules="[{ required: true, message: '请选择作业日期' }]"
        />
        <van-field
          name="time_slot"
          label="时间段"
          :value="timeSlotLabel"
          readonly
          is-link
          @click="showTimeSlotPicker = true"
          :rules="[{ required: true, message: '请选择时间段' }]"
        />
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
      <van-cell-group title="作业人员（可选）" inset>
        <van-field
          name="assigned_worker_id"
          label="作业人员"
          :value="workerName"
          readonly
          is-link
          placeholder="选择作业人员"
          @click="showWorkerPicker = true"
        />
      </van-cell-group>

      <div class="submit-bar">
        <van-button round block type="primary" native-type="submit" :loading="submitting">
          创建预约单
        </van-button>
      </div>
    </van-form>

    <!-- 日期选择器 -->
    <van-popup v-model:show="showDatePicker" position="bottom">
      <van-date-picker
        v-model="currentDate"
        @confirm="onDateConfirm"
        @cancel="showDatePicker = false"
      />
    </van-popup>

    <!-- 时间段选择器 -->
    <van-popup v-model:show="showTimeSlotPicker" position="bottom">
      <van-picker
        :columns="timeSlotOptions"
        @confirm="onTimeSlotConfirm"
        @cancel="showTimeSlotPicker = false"
      />
    </van-popup>

    <!-- 作业人员选择器 -->
    <van-popup v-model:show="showWorkerPicker" position="bottom">
      <van-picker
        :columns="workerOptions"
        @confirm="onWorkerConfirm"
        @cancel="showWorkerPicker = false"
      />
    </van-popup>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { Toast } from 'vant'
import { createAppointment, getTimeSlotOptions } from '@/api/appointment'

const router = useRouter()

const form = ref({
  work_location: '',
  work_content: '',
  work_type: '',
  work_date: '',
  time_slot: '',
  contact_person: '',
  contact_phone: '',
  is_urgent: false,
  priority: 0,
  urgent_reason: '',
  assigned_worker_id: null
})

const submitting = ref(false)
const showDatePicker = ref(false)
const showTimeSlotPicker = ref(false)
const showWorkerPicker = ref(false)
const currentDate = ref(new Date())
const workerName = ref('')

const timeSlotOptions = getTimeSlotOptions().map(opt => ({
  text: opt.label,
  value: opt.value
}))

const workerOptions = ref([]) // 从API获取作业人员列表

const timeSlotLabel = computed(() => {
  if (!form.value.time_slot) return ''
  const option = timeSlotOptions.find(opt => opt.value === form.value.time_slot)
  return option ? option.text : ''
})

async function handleSubmit() {
  try {
    submitting.value = true
    await createAppointment(form.value)
    Toast.success('创建成功')
    router.back()
  } catch (error) {
    Toast.fail(error.message || '创建失败')
  } finally {
    submitting.value = false
  }
}

function onDateConfirm(value) {
  const date = new Date(value)
  form.value.work_date = date.toISOString().split('T')[0]
  showDatePicker.value = false
}

function onTimeSlotConfirm({ selectedOptions }) {
  form.value.time_slot = selectedOptions[0].value
  showTimeSlotPicker.value = false
}

function onWorkerConfirm({ selectedOptions }) {
  form.value.assigned_worker_id = selectedOptions[0].value
  workerName.value = selectedOptions[0].text
  showWorkerPicker.value = false
}
</script>

<style scoped>
.appointment-create {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 80px;
}

.van-cell-group {
  margin-bottom: 12px;
}

.submit-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 16px;
  background-color: #fff;
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.1);
}
</style>
