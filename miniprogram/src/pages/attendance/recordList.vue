<template>
  <view class="record-list-page">
    <!-- 月份选择 -->
    <view class="month-selector">
      <view class="month-btn" @click="prevMonth">‹</view>
      <text class="month-text">{{ currentMonth }}</text>
      <view class="month-btn" @click="nextMonth">›</view>
    </view>

    <!-- 打卡记录列表 -->
    <ListContainer
      :list="list"
      :loading="loading"
      :has-more="hasMore"
      @refresh="onRefresh"
      @load-more="onLoadMore"
    >
      <view class="record-item" v-for="item in list" :key="item.id">
        <view class="record-header">
          <text class="record-type" :style="{ color: getTypeColor(item.type)">{{ getTypeLabel(item.type) }}</text>
          <text class="record-status" :class="item.status">{{ getStatusLabel(item.status) }}</text>
        </view>
        <view class="record-info">
          <text class="record-date">{{ formatDate(item.clock_in_time) }}</text>
          <text class="record-time">{{ formatTime(item.clock_in_time) }}</text>
        </view>
        <view class="record-location" v-if="item.address">
          <text>{{ item.address }}</text>
        </view>
      </view>
    </ListContainer>
  </view>
</template>

<script>
import { getMyRecords, getAttendanceTypeLabel, getAttendanceTypeColor, getStatusLabel as getRecordStatusLabel } from '@/api/attendance'
import ListContainer from '@/components/ListContainer.vue'

export default {
  components: {
    ListContainer
  },
  data() {
    return {
      list: [],
      loading: false,
      hasMore: false,
      page: 1,
      pageSize: 20,
      year: new Date().getFullYear(),
      month: new Date().getMonth() + 1
    }
  },
  computed: {
    currentMonth() {
      return `${this.year}年${this.month}月`
    }
  },
  onShow() {
    this.onRefresh()
  },
  methods: {
    getTypeLabel: getAttendanceTypeLabel,
    getTypeColor: getAttendanceTypeColor,
    getStatusLabel: getRecordStatusLabel,
    async fetchData(callback) {
      if (this.loading) return
      this.loading = true

      try {
        const res = await getMyRecords({
          page: this.page,
          page_size: this.pageSize,
          year: this.year,
          month: this.month
        })

        const items = res.data?.items || res.data || []

        if (this.page === 1) {
          this.list = items
        } else {
          this.list = [...this.list, ...items]
        }

        this.hasMore = items.length >= this.pageSize
      } catch (error) {
        uni.showToast({
          title: error.message || '加载失败',
          icon: 'none'
        })
      } finally {
        this.loading = false
        if (callback) callback()
      }
    },
    onRefresh(callback) {
      this.page = 1
      this.fetchData(callback)
    },
    onLoadMore() {
      this.page++
      this.fetchData()
    },
    prevMonth() {
      if (this.month === 1) {
        this.month = 12
        this.year--
      } else {
        this.month--
      }
      this.onRefresh()
    },
    nextMonth() {
      if (this.month === 12) {
        this.month = 1
        this.year++
      } else {
        this.month++
      }
      this.onRefresh()
    },
    formatDate(dateStr) {
      if (!dateStr) return '-'
      const date = new Date(dateStr)
      return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
    },
    formatTime(dateStr) {
      if (!dateStr) return '-'
      const date = new Date(dateStr)
      return `${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}:${String(date.getSeconds()).padStart(2, '0')}`
    }
  }
}
</script>

<style scoped>
.record-list-page {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.month-selector {
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #ffffff;
  padding: 24rpx;
  border-bottom: 1rpx solid #ebedf0;
}

.month-btn {
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

.record-item {
  background-color: #ffffff;
  border-radius: 16rpx;
  margin-bottom: 24rpx;
  padding: 24rpx;
}

.record-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 16rpx;
}

.record-type {
  font-size: 28rpx;
  font-weight: 500;
}

.record-status {
  font-size: 24rpx;
  padding: 4rpx 12rpx;
  border-radius: 4rpx;
}

.record-status.confirmed {
  background-color: #e8f8ee;
  color: #07c160;
}

.record-status.pending {
  background-color: #fff3e8;
  color: #ff976a;
}

.record-status.rejected {
  background-color: #ffebee;
  color: #ee0a24;
}

.record-info {
  display: flex;
  justify-content: space-between;
  font-size: 26rpx;
  color: #646566;
  margin-bottom: 12rpx;
}

.record-location {
  font-size: 24rpx;
  color: #969799;
}
</style>
