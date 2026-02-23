<template>
  <div class="gantt-header">
    <div class="header-tasks">
      <div class="header-column column-name">任务名称</div>
      <div class="header-column column-duration">工期</div>
      <div class="header-column column-dates">起止时间</div>
      <div class="header-column column-resources">资源</div>
    </div>
    <div class="header-timeline" :class="timelineHeightClass" ref="timelineHeaderRef" :style="{ transform: `translateX(${panOffset}px)` }">
      <!-- 根据时间轴格式动态渲染 -->
      <template v-if="timelineFormat === 'day'">
        <!-- 单层：只显示日期 -->
        <div class="timeline-days-row timeline-single-row">
          <div
            v-for="day in filteredDays"
            :key="day.date"
            class="timeline-day-cell"
            :class="{ 'is-today': day.isToday, 'is-weekend': day.isWeekend }"
            :style="{ left: day.position + 'px', width: dayWidth + 'px' }"
          >
            <div class="day-number">{{ day.day }}</div>
          </div>
        </div>
      </template>

      <template v-else-if="timelineFormat === 'month-day'">
        <!-- 双层：上层月份，下层日期 -->
        <div class="timeline-months-row" v-if="timelineHeaderMonths && timelineHeaderMonths.length > 0">
          <div
            v-for="month in timelineHeaderMonths"
            :key="month.key"
            class="timeline-month-cell"
            :style="{ left: month.position + 'px', width: month.width + 'px' }"
          >
            <div class="month-label">{{ month.label }}</div>
          </div>
        </div>
        <div class="timeline-days-row">
          <div
            v-for="day in filteredDays"
            :key="day.date"
            class="timeline-day-cell"
            :class="{ 'is-today': day.isToday, 'is-weekend': day.isWeekend }"
            :style="{ left: day.position + 'px', width: dayWidth + 'px' }"
          >
            <div class="day-number">{{ day.day }}</div>
          </div>
        </div>
      </template>

      <template v-else-if="timelineFormat === 'year-month'">
        <!-- 双层：上年、下月 -->
        <div class="timeline-years-row" v-if="timelineYears && timelineYears.length > 0">
          <div
            v-for="year in timelineYears"
            :key="year.key"
            class="timeline-year-cell"
            :style="{ left: year.position + 'px', width: year.width + 'px' }"
          >
            <div class="year-label">{{ year.label }}</div>
          </div>
        </div>
        <div class="timeline-months-row-2">
          <div
            v-for="month in timelineMonths"
            :key="month.key"
            class="timeline-month-cell-2"
            :style="{ left: month.position + 'px', width: month.width + 'px' }"
          >
            <div class="month-label-2">{{ month.month }}月</div>
          </div>
        </div>
      </template>

      <template v-else-if="timelineFormat === 'year-month-day'">
        <!-- 三层：年月日 -->
        <div class="timeline-years-row" v-if="timelineYears && timelineYears.length > 0">
          <div
            v-for="year in timelineYears"
            :key="year.key"
            class="timeline-year-cell"
            :style="{ left: year.position + 'px', width: year.width + 'px' }"
          >
            <div class="year-label">{{ year.label }}</div>
          </div>
        </div>
        <div class="timeline-months-row-2" v-if="timelineHeaderMonths && timelineHeaderMonths.length > 0">
          <div
            v-for="month in timelineHeaderMonths"
            :key="month.key"
            class="timeline-month-cell-2"
            :style="{ left: month.position + 'px', width: month.width + 'px' }"
          >
            <div class="month-label-2">{{ month.label }}</div>
          </div>
        </div>
        <div class="timeline-days-row-3">
          <div
            v-for="day in filteredDays"
            :key="day.date"
            class="timeline-day-cell"
            :class="{ 'is-today': day.isToday, 'is-weekend': day.isWeekend }"
            :style="{ left: day.position + 'px', width: dayWidth + 'px' }"
          >
            <div class="day-number">{{ day.day }}</div>
          </div>
        </div>
      </template>

      <template v-else-if="timelineFormat === 'week'">
        <!-- 单层：周 -->
        <div class="timeline-weeks-row timeline-single-row">
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
        </div>
      </template>

      <template v-else-if="timelineFormat === 'month'">
        <!-- 单层：月 -->
        <div class="timeline-months-row-3 timeline-single-row">
          <div
            v-for="month in timelineMonths"
            :key="month.key"
            class="timeline-cell timeline-cell-month"
            :style="{ left: month.position + 'px', width: month.width + 'px' }"
          >
            <div class="cell-month">{{ month.year }}-{{ month.month }}</div>
            <div class="cell-day-count">{{ month.dayCount }}天</div>
          </div>
        </div>
      </template>

      <template v-else-if="timelineFormat === 'quarter'">
        <!-- 单层：季度 -->
        <div class="timeline-quarters-row timeline-single-row">
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
import { ref, computed } from 'vue'

const props = defineProps({
  timelineFormat: {
    type: String,
    default: 'month-day'
  },
  dateDisplayFormat: {
    type: String,
    default: 'all'
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
  timelineHeaderMonths: {
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
  },
  panOffset: {
    type: Number,
    default: 0
  }
})

const timelineHeaderRef = ref(null)

// 计算年份层级数据（用于 year-month 和 year-month-day 格式）
const timelineYears = computed(() => {
  if (props.timelineFormat !== 'year-month' && props.timelineFormat !== 'year-month-day') {
    return []
  }

  const months = props.timelineMonths
  if (!months || months.length === 0) return []

  const years = []
  let currentYear = null
  let yearStartIndex = 0
  let yearStartPosition = null

  months.forEach((month, index) => {
    if (currentYear !== month.year) {
      if (currentYear !== null) {
        const width = month.position - yearStartPosition
        years.push({
          key: `${currentYear}`,
          label: `${currentYear}年`,
          position: yearStartPosition,
          width: width
        })
      }
      currentYear = month.year
      yearStartIndex = index
      yearStartPosition = month.position
    }
  })

  // 添加最后一年
  if (currentYear !== null && months.length > 0) {
    const lastMonth = months[months.length - 1]
    const width = lastMonth.position + lastMonth.width - yearStartPosition
    years.push({
      key: `${currentYear}`,
      label: `${currentYear}年`,
      position: yearStartPosition,
      width: width
    })
  }

  return years
})

// 计算时间轴高度类
const timelineHeightClass = computed(() => {
  const format = props.timelineFormat
  if (format === 'year-month-day') {
    return 'height-3layers' // 三层：年月日
  } else if (format === 'month-day' || format === 'year-month') {
    return 'height-2layers' // 双层
  }
  return 'height-1layer' // 单层
})

// 根据日期显示格式过滤日期
const filteredDays = computed(() => {
  if (!props.timelineDays || props.timelineDays.length === 0) {
    return []
  }

  const format = props.dateDisplayFormat
  const result = []
  let index = 0

  // 显示全部日期
  if (format === 'all') {
    return props.timelineDays
  }

  // 只显示奇数日期
  if (format === 'odd') {
    props.timelineDays.forEach((day) => {
      if (day.day % 2 === 1) {
        result.push({
          ...day,
          position: index * props.dayWidth
        })
        index++
      }
    })
    return result
  }

  // 间隔3天
  if (format === 'interval3') {
    props.timelineDays.forEach((day, i) => {
      if (i % 3 === 0) {
        result.push({
          ...day,
          position: index * props.dayWidth
        })
        index++
      }
    })
    return result
  }

  // 间隔5天
  if (format === 'interval5') {
    props.timelineDays.forEach((day, i) => {
      if (i % 5 === 0) {
        result.push({
          ...day,
          position: index * props.dayWidth
        })
        index++
      }
    })
    return result
  }

  // 每月1号
  if (format === 'first') {
    props.timelineDays.forEach((day) => {
      if (day.day === 1) {
        result.push({
          ...day,
          position: index * props.dayWidth
        })
        index++
      }
    })
    return result
  }

  return props.timelineDays
})

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
  z-index: 200; /* 确保在时间轴上方 */
  background: #f5f7fa;
  box-shadow: 2px 0 4px rgba(0, 0, 0, 0.1); /* 添加阴影增强视觉层次 */
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
  min-width: 800px;
  min-height: 32px; /* 单层默认高度 */
}

.header-timeline.height-1layer {
  min-height: 32px;
}

.header-timeline.height-2layers {
  min-height: 60px; /* 28px * 2 + 4px gap */
}

.header-timeline.height-3layers {
  min-height: 88px; /* 28px * 3 + 4px gap * 2 */
}

/* 单行时间轴 */
.timeline-single-row {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}

/* 双层时间轴 - 上层月份 */
.timeline-months-row {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 28px;
  border-bottom: 1px solid #dcdfe6;
}

.timeline-month-cell {
  position: absolute;
  top: 0;
  height: 100%;
  border-right: 1px solid #e4e7ed;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  font-weight: bold;
  color: #303133;
  background: #fff;
}

.month-label {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  padding: 0 4px;
}

/* 年份行 */
.timeline-years-row {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 28px;
  border-bottom: 1px solid #dcdfe6;
}

.timeline-year-cell {
  position: absolute;
  top: 0;
  height: 100%;
  border-right: 1px solid #e4e7ed;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  font-weight: bold;
  color: #303133;
  background: #fff;
}

.year-label {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  padding: 0 4px;
}

/* 双层时间轴 - 下层日期 */
.timeline-days-row {
  position: absolute;
  top: 28px;
  left: 0;
  right: 0;
  bottom: 0;
}

/* 三层时间轴 - 中层月份 */
.timeline-days-row-3 {
  position: absolute;
  top: 56px;
  left: 0;
  right: 0;
  bottom: 0;
}

/* 双层时间轴 - 下层月份 */
.timeline-months-row-2 {
  position: absolute;
  top: 28px;
  left: 0;
  right: 0;
  height: 28px;
  border-bottom: 1px solid #dcdfe6;
}

.timeline-month-cell-2 {
  position: absolute;
  top: 0;
  height: 100%;
  border-right: 1px solid #e4e7ed;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  font-weight: bold;
  color: #303133;
  background: #fff;
}

.month-label-2 {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  padding: 0 4px;
}

/* 三层时间轴 - 下层月份 */
.timeline-months-row-3 {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}

/* 周视图 */
.timeline-weeks-row {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}

/* 季度视图 */
.timeline-quarters-row {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}

/* 日期单元格 */
.timeline-day-cell {
  position: absolute;
  top: 0;
  height: 100%;
  border-right: 1px solid #e4e7ed;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  color: #606266;
}

.timeline-day-cell.is-today {
  background: #fff3e0;
}

.timeline-day-cell.is-weekend {
  background: #fafafa;
}

.day-number {
  font-size: 12px;
  font-weight: 500;
}

/* 通用单元格样式 */
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

.timeline-cell.is-current {
  background: #fff3e0;
}

.cell-week,
.cell-month,
.cell-quarter {
  font-size: 14px;
  font-weight: bold;
  color: #303133;
}

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

/* 今天标记线 */
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
