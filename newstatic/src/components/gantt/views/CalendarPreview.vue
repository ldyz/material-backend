<template>
  <div class="calendar-preview">
    <div class="preview-grid">
      <div
        v-for="day in daysInRange"
        :key="day.date"
        class="day-cell"
        :class="{
          'working-day': day.isWorking,
          'non-working-day': !day.isWorking,
          'holiday': day.isHoliday,
          'exception': day.isException
        }"
      >
        <div class="day-date">{{ day.dayOfMonth }}</div>
        <div class="day-name">{{ day.dayName }}</div>
        <div v-if="day.isHoliday" class="day-badge holiday-badge">
          {{ t('gantt.calendar.holiday') }}
        </div>
        <div v-if="day.isException" class="day-badge exception-badge">
          {{ day.exceptionType }}
        </div>
      </div>
    </div>
    <div class="preview-legend">
      <div class="legend-item">
        <span class="legend-box working-day"></span>
        <span>{{ t('gantt.calendar.working') }}</span>
      </div>
      <div class="legend-item">
        <span class="legend-box non-working-day"></span>
        <span>{{ t('gantt.calendar.nonWorking') }}</span>
      </div>
      <div class="legend-item">
        <span class="legend-box holiday"></span>
        <span>{{ t('gantt.calendar.holiday') }}</span>
      </div>
      <div class="legend-item">
        <span class="legend-box exception"></span>
        <span>{{ t('gantt.calendar.exception') }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useCalendarStore } from '@/stores/calendarStore'

const { t } = useI18n()
const calendarStore = useCalendarStore()

const props = defineProps({
  calendar: {
    type: Object,
    required: true
  },
  startDate: {
    type: [Date, String],
    default: () => new Date()
  },
  endDate: {
    type: [Date, String],
    default: () => {
      const date = new Date()
      date.setDate(date.getDate() + 30)
      return date
    }
  }
})

const dayNames = [
  t('gantt.calendar.sun'),
  t('gantt.calendar.mon'),
  t('gantt.calendar.tue'),
  t('gantt.calendar.wed'),
  t('gantt.calendar.thu'),
  t('gantt.calendar.fri'),
  t('gantt.calendar.sat')
]

const daysInRange = computed(() => {
  const start = new Date(props.startDate)
  const end = new Date(props.endDate)
  const days = []

  const current = new Date(start)
  while (current <= end) {
    const dateStr = current.toISOString().split('T')[0]
    const dayOfWeek = current.getDay()

    // Check if working day
    let isWorking = props.calendar.workingDays?.includes(dayOfWeek) ?? true

    // Check holidays
    let isHoliday = false
    let holidayName = ''
    if (props.calendar.holidays) {
      const holiday = props.calendar.holidays.find(h => h.date === dateStr)
      if (holiday) {
        isHoliday = true
        isWorking = false
        holidayName = holiday.name
      }
    }

    // Check exceptions
    let isException = false
    let exceptionType = ''
    if (props.calendar.exceptions) {
      const exception = props.calendar.exceptions.find(e => e.date === dateStr)
      if (exception) {
        isException = true
        isWorking = exception.isWorkingDay
        exceptionType = exception.isWorkingDay
          ? t('gantt.calendar.working')
          : t('gantt.calendar.nonWorking')
      }
    }

    days.push({
      date: dateStr,
      dayOfMonth: current.getDate(),
      dayName: dayNames[dayOfWeek],
      isWorking,
      isHoliday,
      holidayName,
      isException,
      exceptionType
    })

    current.setDate(current.getDate() + 1)
  }

  return days
})
</script>

<style scoped>
.calendar-preview {
  background: white;
  border-radius: 4px;
  padding: 16px;
}

.preview-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(50px, 1fr));
  gap: 8px;
  margin-bottom: 16px;
}

.day-cell {
  aspect-ratio: 1;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  padding: 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  position: relative;
}

.day-cell.working-day {
  background: #f0f9ff;
  border-color: #409EFF;
}

.day-cell.non-working-day {
  background: #f5f7fa;
  border-color: #e4e7ed;
}

.day-cell.holiday {
  background: #fef0f0;
  border-color: #F56C6C;
}

.day-cell.exception {
  background: #fdf6ec;
  border-color: #E6A23C;
}

.day-date {
  font-weight: bold;
  font-size: 14px;
  margin-bottom: 4px;
}

.day-name {
  color: #606266;
  font-size: 10px;
}

.day-badge {
  position: absolute;
  bottom: 2px;
  font-size: 8px;
  padding: 2px 4px;
  border-radius: 2px;
  white-space: nowrap;
}

.holiday-badge {
  background: #F56C6C;
  color: white;
}

.exception-badge {
  background: #E6A23C;
  color: white;
}

.preview-legend {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
  font-size: 12px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.legend-box {
  width: 16px;
  height: 16px;
  border-radius: 2px;
  border: 1px solid;
}

.legend-box.working-day {
  background: #f0f9ff;
  border-color: #409EFF;
}

.legend-box.non-working-day {
  background: #f5f7fa;
  border-color: #e4e7ed;
}

.legend-box.holiday {
  background: #fef0f0;
  border-color: #F56C6C;
}

.legend-box.exception {
  background: #fdf6ec;
  border-color: #E6A23C;
}
</style>
