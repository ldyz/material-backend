<template>
  <view class="calendar-page">
    <view class="calendar-header">
      <view class="month-nav">
        <view class="nav-btn" @click="prevMonth">‹</view>
        <text class="month-text">{{ currentMonth }}</text>
        <view class="nav-btn" @click="nextMonth">›</view>
      </view>
    </view>

    <view class="calendar-body">
      <view class="week-header">
        <text v-for="day in weekDays" :key="day" class="week-day">{{ day }}</text>
      </view>

      <view class="calendar-grid">
        <view
          v-for="(item, index) in calendarDays"
          :key="index"
          class="day-cell"
          :class="{ 'is-today': item.isToday, 'is-selected': item.isSelected, 'is-other-month': item.isOtherMonth }"
          @click="selectDay(item)"
        >
          <text class="day-number">{{ item.day }}</text>
          <view class="day-dot" v-if="item.hasEvents"></view>
        </view>
      </view>
    </view>

    <view class="events-section" v-if="selectedEvents.length > 0">
      <view class="section-title">{{ selectedDate }} 的预约</view>
      <view class="event-list">
        <view
          class="event-item"
          v-for="event in selectedEvents"
          :key="event.id"
          @click="goToDetail(event.id)"
        >
          <text class="event-title">{{ event.title }}</text>
          <text class="event-time">{{ getTimeSlotLabel(event.time_slot) }}</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script>
import { getCalendarView, getTimeSlotLabel } from '@/api/appointment'

export default {
  data() {
    return {
      year: new Date().getFullYear(),
      month: new Date().getMonth() + 1,
      weekDays: ['日', '一', '二', '三', '四', '五', '六'],
      calendarDays: [],
      selectedDate: '',
      selectedEvents: [],
      eventsMap: {}
    }
  },
  computed: {
    currentMonth() {
      return `${this.year}年${this.month}月`
    }
  },
  onLoad() {
    this.initCalendar()
    this.fetchEvents()
  },
  methods: {
    getTimeSlotLabel,
    initCalendar() {
      const today = new Date()
      const firstDay = new Date(this.year, this.month - 1, 1)
      const lastDay = new Date(this.year, this.month, 0)

      const days = []
      const startDay = firstDay.getDay()
      const totalDays = lastDay.getDate()

      // 上个月的日期
      const prevMonthLastDay = new Date(this.year, this.month - 1, 0).getDate()
      for (let i = startDay - 1; i >= 0; i--) {
        days.push({
          day: prevMonthLastDay - i,
          isOtherMonth: true,
          dateStr: this.formatDateStr(this.year, this.month - 1, prevMonthLastDay - i)
        })
      }

      // 当前月的日期
      for (let i = 1; i <= totalDays; i++) {
        const dateStr = this.formatDateStr(this.year, this.month, i)
        days.push({
          day: i,
          isToday: today.getFullYear() === this.year && today.getMonth() + 1 === this.month && today.getDate() === i,
          dateStr,
          hasEvents: this.eventsMap[dateStr]?.length > 0
        })
      }

      // 下个月的日期
      const remaining = 42 - days.length
      for (let i = 1; i <= remaining; i++) {
        days.push({
          day: i,
          isOtherMonth: true,
          dateStr: this.formatDateStr(this.year, this.month + 1, i)
        })
      }

      this.calendarDays = days
    },
    formatDateStr(year, month, day) {
      const m = month < 1 ? 12 : month > 12 ? 1 : month
      const y = month < 1 ? year - 1 : month > 12 ? year + 1 : year
      return `${y}-${String(m).padStart(2, '0')}-${String(day).padStart(2, '0')}`
    },
    async fetchEvents() {
      try {
        const res = await getCalendarView({
          year: this.year,
          month: this.month
        })

        const events = res.data || []
        this.eventsMap = {}

        events.forEach(event => {
          const date = event.work_date?.split('T')[0]
          if (date) {
            if (!this.eventsMap[date]) {
              this.eventsMap[date] = []
            }
            this.eventsMap[date].push(event)
          }
        })

        this.initCalendar()
      } catch (error) {
        console.error('获取日历数据失败', error)
      }
    },
    selectDay(item) {
      if (item.isOtherMonth) return

      this.selectedDate = item.dateStr
      this.selectedEvents = this.eventsMap[item.dateStr] || []

      this.calendarDays.forEach(d => {
        d.isSelected = d.dateStr === item.dateStr
      })
    },
    prevMonth() {
      if (this.month === 1) {
        this.month = 12
        this.year--
      } else {
        this.month--
      }
      this.fetchEvents()
    },
    nextMonth() {
      if (this.month === 12) {
        this.month = 1
        this.year++
      } else {
        this.month++
      }
      this.fetchEvents()
    },
    goToDetail(id) {
      uni.navigateTo({ url: `/pages/appointment/detail?id=${id}` })
    }
  }
}
</script>

<style scoped>
.calendar-page {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.calendar-header {
  background-color: #ffffff;
  padding: 24rpx;
  border-bottom: 1rpx solid #ebedf0;
}

.month-nav {
  display: flex;
  justify-content: center;
  align-items: center;
}

.nav-btn {
  width: 64rpx;
  height: 64rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 36rpx;
  color: #1989fa;
}

.month-text {
  font-size: 32rpx;
  font-weight: 500;
  color: #323233;
  margin: 0 32rpx;
}

.calendar-body {
  background-color: #ffffff;
  padding: 24rpx;
}

.week-header {
  display: flex;
  margin-bottom: 16rpx;
}

.week-day {
  flex: 1;
  text-align: center;
  font-size: 24rpx;
  color: #969799;
}

.calendar-grid {
  display: flex;
  flex-wrap: wrap;
}

.day-cell {
  width: 14.28%;
  aspect-ratio: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  position: relative;
}

.day-number {
  font-size: 28rpx;
  color: #323233;
}

.day-cell.is-today .day-number {
  background-color: #1989fa;
  color: #ffffff;
  border-radius: 50%;
  width: 48rpx;
  height: 48rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.day-cell.is-selected .day-number {
  background-color: #e8f4ff;
  border-radius: 50%;
  width: 48rpx;
  height: 48rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.day-cell.is-other-month .day-number {
  color: #c8c9cc;
}

.day-dot {
  width: 8rpx;
  height: 8rpx;
  background-color: #1989fa;
  border-radius: 50%;
  position: absolute;
  bottom: 8rpx;
}

.events-section {
  background-color: #ffffff;
  margin-top: 24rpx;
  padding: 24rpx;
}

.section-title {
  font-size: 32rpx;
  font-weight: 500;
  color: #323233;
  margin-bottom: 24rpx;
}

.event-item {
  display: flex;
  justify-content: space-between;
  padding: 24rpx 0;
  border-bottom: 1rpx solid #ebedf0;
}

.event-item:last-child {
  border-bottom: none;
}

.event-title {
  font-size: 28rpx;
  color: #323233;
}

.event-time {
  font-size: 24rpx;
  color: #969799;
}
</style>
