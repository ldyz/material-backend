<template>
  <view class="inbound-list-page">
    <!-- 筛选栏 -->
    <view class="filter-bar">
      <picker :value="statusIndex" :range="statusOptions" range-key="label" @change="onStatusChange">
        <view class="filter-item">
          <text>{{ statusOptions[statusIndex].label }}</text>
          <text class="arrow">▼</text>
        </view>
      </picker>
    </view>

    <!-- 列表 -->
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
        :title="item.order_number"
        :status="item.status"
        :status-text="getStatusLabel(item.status)"
        :status-color="getStatusColor(item.status)"
        :items="[
          { label: '项目', value: item.project_name },
          { label: '入库类型', value: item.type },
          { label: '创建时间', value: formatDate(item.created_at) }
        ]"
        clickable
        @click="goToDetail(item.id)"
      />
    </ListContainer>

    <!-- 新增按钮 -->
    <view class="fab-btn" @click="goToCreate">
      <text>+</text>
    </view>
  </view>
</template>

<script>
import { getInboundOrders } from '@/api/inbound'
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
        { value: 'approved', label: '已通过' },
        { value: 'rejected', label: '已拒绝' }
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
        const res = await getInboundOrders({
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
    getStatusLabel(status) {
      const labels = {
        draft: '草稿',
        pending: '待审批',
        approved: '已通过',
        rejected: '已拒绝'
      }
      return labels[status] || status
    },
    getStatusColor(status) {
      const colors = {
        draft: 'default',
        pending: 'warning',
        approved: 'success',
        rejected: 'danger'
      }
      return colors[status] || 'default'
    },
    formatDate(dateStr) {
      if (!dateStr) return '-'
      const date = new Date(dateStr)
      return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
    },
    goToDetail(id) {
      uni.navigateTo({ url: `/pages/inbound/detail?id=${id}` })
    },
    goToCreate() {
      uni.navigateTo({ url: '/pages/inbound/create' })
    }
  }
}
</script>

<style scoped>
.inbound-list-page {
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
