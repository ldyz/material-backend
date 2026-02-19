<template>
  <el-dialog
    v-model="dialogVisible"
    :title="t('gantt.calendar.editTitle')"
    width="700px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <div v-loading="loading" class="calendar-dialog">
      <!-- Calendar Selection -->
      <div class="form-section">
        <label class="form-label">{{ t('gantt.calendar.calendarType') }}</label>
        <el-select
          v-model="selectedPreset"
          :placeholder="t('gantt.calendar.selectPreset')"
          @change="handlePresetChange"
        >
          <el-option
            v-for="preset in calendarStore.standardPresets"
            :key="preset.id"
            :label="preset.name"
            :value="preset.id"
          >
            <div class="preset-option">
              <div class="preset-name">{{ preset.name }}</div>
              <div class="preset-desc">{{ preset.description }}</div>
            </div>
          </el-option>
        </el-select>
      </div>

      <!-- Working Days -->
      <div class="form-section">
        <label class="form-label">{{ t('gantt.calendar.workingDays') }}</label>
        <div class="working-days-selector">
          <el-checkbox-group v-model="calendar.workingDays" @change="handleWorkingDaysChange">
            <el-checkbox
              v-for="day in weekDays"
              :key="day.value"
              :label="day.value"
              :disabled="selectedPreset !== 'custom'"
            >
              {{ day.label }}
            </el-checkbox>
          </el-checkbox-group>
        </div>
      </div>

      <!-- Working Hours -->
      <div class="form-section">
        <label class="form-label">{{ t('gantt.calendar.workingHours') }}</label>
        <div class="working-hours-input">
          <el-time-picker
            v-model="workingHoursStart"
            :placeholder="t('gantt.calendar.startTime')"
            format="HH:mm"
            value-format="HH:mm"
            :disabled="selectedPreset !== 'custom'"
            @change="handleWorkingHoursChange"
          />
          <span class="time-separator">{{ t('gantt.calendar.to') }}</span>
          <el-time-picker
            v-model="workingHoursEnd"
            :placeholder="t('gantt.calendar.endTime')"
            format="HH:mm"
            value-format="HH:mm"
            :disabled="selectedPreset !== 'custom'"
            @change="handleWorkingHoursChange"
          />
          <span class="hours-total">
            {{ t('gantt.calendar.totalHours') }}: {{ calculateWorkingHours() }}
          </span>
        </div>
      </div>

      <!-- Holidays -->
      <div class="form-section">
        <div class="section-header">
          <label class="form-label">{{ t('gantt.calendar.holidays') }}</label>
          <el-button
            type="primary"
            size="small"
            @click="handleAddHoliday"
          >
            {{ t('gantt.calendar.addHoliday') }}
          </el-button>
        </div>
        <div class="holidays-list">
          <el-table
            :data="calendar.holidays"
            stripe
            max-height="200"
            @row-click="handleHolidayClick"
          >
            <el-table-column
              prop="date"
              :label="t('gantt.calendar.date')"
              width="120"
            />
            <el-table-column
              prop="name"
              :label="t('gantt.calendar.name')"
            />
            <el-table-column
              prop="recurring"
              :label="t('gantt.calendar.recurring')"
              width="100"
              align="center"
            >
              <template #default="scope">
                <el-tag v-if="scope.row.recurring" type="info" size="small">
                  {{ t('gantt.calendar.yearly') }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column
              :label="t('common.actions')"
              width="80"
              align="center"
            >
              <template #default="scope">
                <el-button
                  type="danger"
                  size="small"
                  link
                  @click.stop="handleRemoveHoliday(scope.$index)"
                >
                  {{ t('common.delete') }}
                </el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-empty
            v-if="!calendar.holidays || calendar.holidays.length === 0"
            :description="t('gantt.calendar.noHolidays')"
            :image-size="80"
          />
        </div>
      </div>

      <!-- Exception Dates -->
      <div class="form-section">
        <div class="section-header">
          <label class="form-label">{{ t('gantt.calendar.exceptions') }}</label>
          <el-button
            type="primary"
            size="small"
            @click="handleAddException"
          >
            {{ t('gantt.calendar.addException') }}
          </el-button>
        </div>
        <div class="exceptions-list">
          <el-table
            :data="calendar.exceptions"
            stripe
            max-height="200"
          >
            <el-table-column
              prop="date"
              :label="t('gantt.calendar.date')"
              width="120"
            />
            <el-table-column
              prop="isWorkingDay"
              :label="t('gantt.calendar.type')"
              width="120"
            >
              <template #default="scope">
                <el-tag :type="scope.row.isWorkingDay ? 'success' : 'warning'" size="small">
                  {{ scope.row.isWorkingDay ? t('gantt.calendar.working') : t('gantt.calendar.nonWorking') }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column
              prop="description"
              :label="t('gantt.calendar.description')"
            />
            <el-table-column
              :label="t('common.actions')"
              width="80"
              align="center"
            >
              <template #default="scope">
                <el-button
                  type="danger"
                  size="small"
                  link
                  @click="handleRemoveException(scope.$index)"
                >
                  {{ t('common.delete') }}
                </el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-empty
            v-if="!calendar.exceptions || calendar.exceptions.length === 0"
            :description="t('gantt.calendar.noExceptions')"
            :image-size="80"
          />
        </div>
      </div>

      <!-- Preview -->
      <div class="preview-section">
        <div class="preview-header">
          <span>{{ t('gantt.calendar.preview') }}</span>
        </div>
        <CalendarPreview
          :calendar="calendar"
          :start-date="previewStartDate"
          :end-date="previewEndDate"
        />
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">{{ t('common.cancel') }}</el-button>
        <el-button
          type="primary"
          :loading="saving"
          @click="handleSave"
        >
          {{ t('common.save') }}
        </el-button>
      </div>
    </template>
  </el-dialog>

  <!-- Holiday Edit Dialog -->
  <el-dialog
    v-model="holidayDialogVisible"
    :title="editingHoliday ? t('gantt.calendar.editHoliday') : t('gantt.calendar.addHoliday')"
    width="400px"
  >
    <el-form :model="holidayForm" label-width="100px">
      <el-form-item :label="t('gantt.calendar.date')">
        <el-date-picker
          v-model="holidayForm.date"
          type="date"
          format="YYYY-MM-DD"
          value-format="YYYY-MM-DD"
        />
      </el-form-item>
      <el-form-item :label="t('gantt.calendar.name')">
        <el-input v-model="holidayForm.name" />
      </el-form-item>
      <el-form-item :label="t('gantt.calendar.recurring')">
        <el-switch v-model="holidayForm.recurring" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="holidayDialogVisible = false">{{ t('common.cancel') }}</el-button>
      <el-button type="primary" @click="saveHoliday">{{ t('common.save') }}</el-button>
    </template>
  </el-dialog>

  <!-- Exception Edit Dialog -->
  <el-dialog
    v-model="exceptionDialogVisible"
    :title="editingException ? t('gantt.calendar.editException') : t('gantt.calendar.addException')"
    width="400px"
  >
    <el-form :model="exceptionForm" label-width="120px">
      <el-form-item :label="t('gantt.calendar.date')">
        <el-date-picker
          v-model="exceptionForm.date"
          type="date"
          format="YYYY-MM-DD"
          value-format="YYYY-MM-DD"
        />
      </el-form-item>
      <el-form-item :label="t('gantt.calendar.type')">
        <el-radio-group v-model="exceptionForm.isWorkingDay">
          <el-radio :label="true">{{ t('gantt.calendar.working') }}</el-radio>
          <el-radio :label="false">{{ t('gantt.calendar.nonWorking') }}</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item :label="t('gantt.calendar.description')">
        <el-input v-model="exceptionForm.description" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="exceptionDialogVisible = false">{{ t('common.cancel') }}</el-button>
      <el-button type="primary" @click="saveException">{{ t('common.save') }}</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { useCalendarStore } from '@/stores/calendarStore'
import CalendarPreview from '../views/CalendarPreview.vue'

const { t } = useI18n()
const calendarStore = useCalendarStore()

// Props
const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  calendarId: {
    type: Number,
    default: null
  },
  projectId: {
    type: Number,
    default: null
  }
})

// Emits
const emit = defineEmits(['update:modelValue', 'calendarUpdated'])

// State
const dialogVisible = ref(false)
const loading = ref(false)
const saving = ref(false)
const selectedPreset = ref('standard')

const calendar = ref({
  workingDays: [1, 2, 3, 4, 5],
  workingHours: { start: '09:00', end: '17:00', hours: 8 },
  holidays: [],
  exceptions: []
})

const workingHoursStart = ref('09:00')
const workingHoursEnd = ref('17:00')

// Holiday dialog
const holidayDialogVisible = ref(false)
const editingHoliday = ref(null)
const holidayForm = ref({
  date: null,
  name: '',
  recurring: false
})

// Exception dialog
const exceptionDialogVisible = ref(false)
const editingException = ref(null)
const exceptionForm = ref({
  date: null,
  isWorkingDay: false,
  description: ''
})

// Preview date range
const previewStartDate = ref(new Date())
const previewEndDate = ref(new Date(Date.now() + 30 * 24 * 60 * 60 * 1000)) // 30 days

// Week days for selection
const weekDays = [
  { value: 0, label: t('gantt.calendar.sunday') },
  { value: 1, label: t('gantt.calendar.monday') },
  { value: 2, label: t('gantt.calendar.tuesday') },
  { value: 3, label: t('gantt.calendar.wednesday') },
  { value: 4, label: t('gantt.calendar.thursday') },
  { value: 5, label: t('gantt.calendar.friday') },
  { value: 6, label: t('gantt.calendar.saturday') }
]

// Methods
const handlePresetChange = (presetId) => {
  const preset = calendarStore.standardPresets.find(p => p.id === presetId)
  if (preset) {
    calendar.value = {
      workingDays: [...preset.workingDays],
      workingHours: { ...preset.workingHours },
      holidays: [...(preset.holidays || [])],
      exceptions: []
    }
    workingHoursStart.value = preset.workingHours.start
    workingHoursEnd.value = preset.workingHours.end
  }
}

const handleWorkingDaysChange = () => {
  selectedPreset.value = 'custom'
}

const handleWorkingHoursChange = () => {
  selectedPreset.value = 'custom'
  calendar.value.workingHours = {
    start: workingHoursStart.value,
    end: workingHoursEnd.value,
    hours: calculateWorkingHours()
  }
}

const calculateWorkingHours = () => {
  if (!workingHoursStart.value || !workingHoursEnd.value) {
    return 0
  }

  const start = new Date(`2000-01-01T${workingHoursStart.value}`)
  const end = new Date(`2000-01-01T${workingHoursEnd.value}`)
  const diff = (end - start) / (1000 * 60 * 60)

  return Math.max(0, Math.round(diff * 100) / 100)
}

const handleAddHoliday = () => {
  editingHoliday.value = null
  holidayForm.value = {
    date: new Date().toISOString().split('T')[0],
    name: '',
    recurring: false
  }
  holidayDialogVisible.value = true
}

const handleHolidayClick = (row) => {
  editingHoliday.value = row
  holidayForm.value = {
    date: row.date,
    name: row.name,
    recurring: row.recurring || false
  }
  holidayDialogVisible.value = true
}

const saveHoliday = () => {
  if (!holidayForm.value.date || !holidayForm.value.name) {
    ElMessage.warning(t('gantt.calendar.fillRequiredFields'))
    return
  }

  if (!calendar.value.holidays) {
    calendar.value.holidays = []
  }

  if (editingHoliday.value) {
    const index = calendar.value.holidays.findIndex(h => h === editingHoliday.value)
    if (index !== -1) {
      calendar.value.holidays[index] = { ...holidayForm.value }
    }
  } else {
    calendar.value.holidays.push({ ...holidayForm.value })
  }

  holidayDialogVisible.value = false
  editingHoliday.value = null
}

const handleRemoveHoliday = (index) => {
  calendar.value.holidays.splice(index, 1)
}

const handleAddException = () => {
  editingException.value = null
  exceptionForm.value = {
    date: new Date().toISOString().split('T')[0],
    isWorkingDay: false,
    description: ''
  }
  exceptionDialogVisible.value = true
}

const saveException = () => {
  if (!exceptionForm.value.date) {
    ElMessage.warning(t('gantt.calendar.fillRequiredFields'))
    return
  }

  if (!calendar.value.exceptions) {
    calendar.value.exceptions = []
  }

  if (editingException.value) {
    const index = calendar.value.exceptions.findIndex(e => e === editingException.value)
    if (index !== -1) {
      calendar.value.exceptions[index] = { ...exceptionForm.value }
    }
  } else {
    calendar.value.exceptions.push({ ...exceptionForm.value })
  }

  exceptionDialogVisible.value = false
  editingException.value = null
}

const handleRemoveException = (index) => {
  calendar.value.exceptions.splice(index, 1)
}

const handleSave = async () => {
  saving.value = true

  try {
    if (props.calendarId) {
      await calendarStore.updateCalendar(props.calendarId, calendar.value)
    } else {
      await calendarStore.createCalendar({
        ...calendar.value,
        project_id: props.projectId
      })
    }

    ElMessage.success(t('gantt.calendar.saveSuccess'))
    emit('calendarUpdated', calendar.value)
    handleClose()
  } catch (error) {
    console.error('Failed to save calendar:', error)
    ElMessage.error(t('gantt.calendar.saveError'))
  } finally {
    saving.value = false
  }
}

const handleClose = () => {
  dialogVisible.value = false
  emit('update:modelValue', false)
  nextTick(() => {
    resetForm()
  })
}

const resetForm = () => {
  selectedPreset.value = 'standard'
  calendar.value = {
    workingDays: [1, 2, 3, 4, 5],
    workingHours: { start: '09:00', end: '17:00', hours: 8 },
    holidays: [],
    exceptions: []
  }
  workingHoursStart.value = '09:00'
  workingHoursEnd.value = '17:00'
}

// Watch modelValue changes
watch(() => props.modelValue, (newVal) => {
  dialogVisible.value = newVal
  if (newVal && props.calendarId) {
    // Load existing calendar
    const existingCalendar = calendarStore.getCalendarById(props.calendarId)
    if (existingCalendar) {
      calendar.value = { ...existingCalendar }
      workingHoursStart.value = existingCalendar.workingHours?.start || '09:00'
      workingHoursEnd.value = existingCalendar.workingHours?.end || '17:00'
    }
  } else if (newVal) {
    resetForm()
  }
})

watch(dialogVisible, (newVal) => {
  emit('update:modelValue', newVal)
})
</script>

<style scoped>
.calendar-dialog {
  padding: 10px 0;
}

.form-section {
  margin-bottom: 24px;
}

.form-label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #606266;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.preset-option {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.preset-name {
  font-weight: 500;
}

.preset-desc {
  font-size: 12px;
  color: #909399;
}

.working-days-selector {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.working-hours-input {
  display: flex;
  align-items: center;
  gap: 12px;
}

.time-separator {
  color: #606266;
}

.hours-total {
  margin-left: auto;
  font-size: 12px;
  color: #909399;
  font-weight: 500;
}

.holidays-list,
.exceptions-list {
  margin-top: 12px;
}

.preview-section {
  margin-top: 24px;
  padding: 16px;
  background: #f5f7fa;
  border-radius: 4px;
}

.preview-header {
  font-weight: 500;
  margin-bottom: 12px;
  color: #606266;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
