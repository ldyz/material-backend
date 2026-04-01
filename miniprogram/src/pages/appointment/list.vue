<template>
  <view class="appointment-list-page">
    <view class="filter-bar">
      <picker :value="statusIndex" :range="statusOptions" range-key="label" @change="onStatusChange">
        <view class="filter-item">
          <text>{{ statusOptions[statusIndex].label }}</text>
          <text class="arrow">▼</text>
        </view>
      </picker>
    </view>

    <ListContainer
      :list="list"
      :loading="loading"
      :has-more="hasMore"
      @refresh="onRefresh"
      @load-more="onLoadMore"
    >
      <ListItemCard
        v-for="item in list"
        :key="item.id"
        :title="item.title || item.appointment_number"
        :status="item.status"
        :status-text="getStatusLabel(item.status)"
        :status-color="getStatusColor(item.status)"
        :items="[
          { label: '作业日期', value: formatDate(item.work_date) },
          { label: '时间段', value: getTimeSlotLabel(item.time_slot) },
          { label: '作业地点', value: item.location }
        ]"
        clickable
        @click="goToDetail(item.id)"
      />
    </ListContainer>

    <view class="fab-btn" @click="goToCreate">
      <text>+</text>
    </view>
  </view>
</template>

<script>
import { getMyAppointments, getStatusLabel, getStatusColor, getTimeSlotLabel } from '@/api/appointment'
import ListContainer from '@/components/ListContainer.vue'
import ListItemCard from '@/components/ListItemCard.vue'

export default {
  components: {
    ListContainer,
    ListItemCard
  },
  data() {
    return {
      list: [],
      loading: false,
      hasMore: false,
      page: 1,
      pageSize: 20,
      statusIndex: 0,
      statusOptions: [
        { value: '', label: '全部状态' },
        { value: 'draft', label: '草稿' },
        { value: 'pending', label: '待审批' },
        { value: 'scheduled', label: '已排期' },
        { value: 'in_progress', label: '进行中' },
        { value: 'completed', label: '已完成' },
        { value: 'cancelled', label: '已取消' }
      ]
    }
  },
  onShow() {
    this.onRefresh()
  },
  methods: {
    async fetchData(callback) {
      if (this.loading) return
      this.loading = true

      try {
        const status = this.statusOptions[this.statusIndex].value
        const res = await getMyAppointments({
          page: this.page,
          page_size: this.pageSize,
          status: status || undefined
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
    onStatusChange(e) {
      this.statusIndex = e.detail.value
      this.onRefresh()
    },
    getStatusLabel,
    getStatusColor,
    getTimeSlotLabel,
    formatDate(dateStr) {
      if (!dateStr) return '-'
      const date = new Date(dateStr)
      return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
    },
    goToDetail(id) {
      uni.navigateTo({ url: `/pages/appointment/detail?id=${id}` })
    },
    goToCreate() {
      uni.navigateTo({ url: '/pages/appointment/create' })
    }
  }
}
</script>

<style scoped>
.appointment-list-page {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.filter-bar {
  display: flex;
  background-color: #ffffff;
  padding: 24rpx;
  border-bottom: 1rpx solid #ebedf0;
}

.filter-item {
  display: flex;
  align-items: center;
  padding: 16rpx 24rpx;
  background-color: #f7f8fa;
  border-radius: 8rpx;
  font-size: 28rpx;
  color: #323233;
}

.arrow {
  font-size: 20rpx;
  color: #969799;
  margin-left: 8rpx;
}

.fab-btn {
  position: fixed;
  right: 32rpx;
  bottom: 100rpx;
  width: 96rpx;
  height: 96rpx;
  background-color: #1989fa;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8rpx 24rpx rgba(25, 137, 250, 0.4);
}

.fab-btn text {
  font-size: 48rpx;
  color: #ffffff;
}

.fab-btn:active {
  transform: scale(0.95);
}
</style>
