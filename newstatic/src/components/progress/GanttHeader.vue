<template>
  <div class="gantt-header">
    <div class="header-tasks">
      <div class="header-column column-name">任务名称</div>
      <div class="header-column column-duration">工期</div>
      <div class="header-column column-dates">起止时间</div>
      <div class="header-column column-resources">资源</div>
    </div>
    <div class="header-timeline" ref="timelineHeaderRef">
      <!-- 多级时间轴 -->
      <template v-if="viewMode === 'day'">
        <div
          v-for="day in timelineDays"
          :key="day.date"
          class="timeline-cell"
          :class="{ 'is-today': day.isToday, 'is-weekend': day.isWeekend }"
          :style="{ left: day.position + 'px', width: dayWidth + 'px' }"
        >
          <div class="cell-date">{{ day.day }}</div>
          <div class="cell-weekday">{{ day.weekday }}</div>
        </div>
      </template>

      <template v-else-if="viewMode === 'week'">
        <div
          v-for="week in timelineWeeks"
          :key="week.key"
          class="timeline-cell timeline-cell-week"
          :class="{ 'is-current': week.isCurrent }"
          :style="{ left: week.position + 'px', width: week.width + 'px' }"
        >
          <div class="cell-week">W{{ week.weekNumber }}</div>
          <div class="cell-date-range">{{ week.start }} ~ {{ week.end }}</div>
        </div>
      </template>

      <template v-else-if="viewMode === 'month'">
        <div
          v-for="month in timelineMonths"
          :key="month.key"
          class="timeline-cell timeline-cell-month"
          :style="{ left: month.position + 'px', width: month.width + 'px' }"
        >
          <div class="cell-month">{{ month.year }}-{{ month.month }}</div>
          <div class="cell-day-count">{{ month.dayCount }}天</div>
        </div>
      </template>

      <template v-else-if="viewMode === 'quarter'">
        <div
          v-for="quarter in timelineQuarters"
          :key="quarter.key"
          class="timeline-cell timeline-cell-quarter"
          :class="{ 'is-current': quarter.isCurrent }"
          :style="{ left: quarter.position + 'px', width: quarter.width + 'px' }"
        >
          <div class="cell-quarter">Q{{ quarter.quarter }}</div>
          <div class="cell-year">{{ quarter.year }}</div>
        </div>
      </template>

      <!-- 今天标记线 -->
      <div
        v-if="todayPosition !== null"
        class="today-line"
        :style="{ left: todayPosition + 'px' }"
      ></div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

defineProps({
  viewMode: {
    type: String,
    default: 'day'
  },
  timelineDays: {
    type: Array,
    default: () => []
  },
  timelineWeeks: {
    type: Array,
    default: () => []
  },
  timelineMonths: {
    type: Array,
    default: () => []
  },
  timelineQuarters: {
    type: Array,
    default: () => []
  },
  dayWidth: {
    type: Number,
    default: 40
  },
  todayPosition: {
    type: Number,
    default: null
  }
})

const timelineHeaderRef = ref(null)

defineExpose({
  timelineHeaderRef
})
</script>

<style scoped>
/* 表头 */
.gantt-header {
  display: flex;
  background: #f5f7fa;
  border-bottom: 1px solid #dcdfe6;
  position: sticky;
  top: 0;
  z-index: 10;
  flex-shrink: 0;
}

.header-tasks {
  width: 550px;
  padding: 0;
  font-weight: bold;
  color: #303133;
  border-right: 1px solid #dcdfe6;
  flex-shrink: 0;
  display: flex;
  position: sticky;
  left: 0;
  z-index: 20;
  background: #f5f7fa;
}

.header-column {
  padding: 12px 8px;
  font-size: 12px;
  display: flex;
  align-items: center;
  border-right: 1px solid #e4e7ed;
}

.header-column:last-child {
  border-right: none;
}

.column-name {
  flex: 0 0 200px;
}

.column-duration {
  flex: 0 0 70px;
  justify-content: center;
}

.column-dates {
  flex: 0 0 150px;
  justify-content: center;
}

.column-resources {
  flex: 1;
  justify-content: center;
}

.header-timeline {
  flex: 1;
  position: relative;
  min-height: 50px;
  min-width: 800px;
}

.timeline-cell {
  position: absolute;
  top: 0;
  height: 100%;
  border-right: 1px solid #e4e7ed;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  color: #606266;
}

.timeline-cell.is-today,
.timeline-cell.is-current {
  background: #fff3e0;
}

.timeline-cell.is-weekend {
  background: #fafafa;
}

.cell-date,
.cell-week,
.cell-month,
.cell-quarter {
  font-size: 14px;
  font-weight: bold;
  color: #303133;
}

.cell-weekday,
.cell-date-range,
.cell-day-count,
.cell-year {
  font-size: 11px;
  color: #909399;
  margin-top: 2px;
}

.timeline-cell-week {
  padding: 8px 0;
}

.timeline-cell-month {
  padding: 12px 0;
}

.timeline-cell-quarter {
  padding: 16px 0;
}

.today-line {
  position: absolute;
  top: 0;
  bottom: 0;
  width: 2px;
  background: #e6a23c;
  z-index: 5;
}

.today-line::before {
  content: '今天';
  position: absolute;
  top: -20px;
  left: 50%;
  transform: translateX(-50%);
  font-size: 11px;
  color: #e6a23c;
  white-space: nowrap;
}
</style>
