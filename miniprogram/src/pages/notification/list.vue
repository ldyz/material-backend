<template>
  <view class="notification-list-page">
    <ListContainer
      :list="list"
      :loading="loading"
      :has-more="hasMore"
      @refresh="onRefresh"
      @load-more="onLoadMore"
    >
      <view
        class="notification-item"
        :class="{ unread: !item.is_read }"
        v-for="item in list"
        :key="item.id"
        @click="goToDetail(item)"
      >
        <view class="notification-header">
          <text class="notification-title">{{ item.title }}</text>
          <text class="notification-time">{{ formatTime(item.created_at) }}</text>
        </view>
        <view class="notification-content">
          <text>{{ item.content }}</text>
        </view>
      </view>
    </ListContainer>
  </view>
</template>

<script>
import { getNotifications } from '@/api/notification'
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
      pageSize: 20
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
        const res = await getNotifications({
          page: this.page,
          page_size: this.pageSize
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
    formatTime(dateStr) {
      if (!dateStr) return '-'
      const date = new Date(dateStr)
      const now = new Date()
      const diff = now - date

      if (diff < 60000) {
        return '刚刚'
      } else if (diff < 3600000) {
        return `${Math.floor(diff / 60000)}分钟前`
      } else if (diff < 86400000) {
        return `${Math.floor(diff / 3600000)}小时前`
      } else {
        return `${date.getMonth() + 1}月${date.getDate()}日`
      }
    },
    goToDetail(item) {
      if (item.related_id && item.related_type) {
        const urlMap = {
          'plan': `/pages/plan/detail?id=${item.related_id}`,
          'inbound': `/pages/inbound/detail?id=${item.related_id}`,
          'requisition': `/pages/requisition/detail?id=${item.related_id}`,
          'appointment': `/pages/appointment/detail?id=${item.related_id}`
        }
        const url = urlMap[item.related_type]
        if (url) {
          uni.navigateTo({ url })
        }
      }
    }
  }
}
</script>

<style scoped>
.notification-list-page {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.notification-item {
  background-color: #ffffff;
  padding: 24rpx;
  margin-bottom: 2rpx;
}

.notification-item.unread {
  background-color: #e8f4ff;
}

.notification-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 12rpx;
}

.notification-title {
  font-size: 30rpx;
  font-weight: 500;
  color: #323233;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.notification-time {
  font-size: 24rpx;
  color: #969799;
  margin-left: 16rpx;
}

.notification-content {
  font-size: 26rpx;
  color: #646566;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
