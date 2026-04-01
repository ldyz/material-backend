<template>
  <view class="list-container">
    <!-- 下拉刷新 -->
    <scroll-view
      scroll-y
      class="scroll-view"
      :refresher-enabled="refreshable"
      :refresher-triggered="isRefreshing"
      @refresherrefresh="onRefresh"
      @scrolltolower="onLoadMore"
    >
      <!-- 列表内容 -->
      <view class="list-content">
        <slot></slot>
      </view>

      <!-- 加载状态 -->
      <view class="loading-wrapper" v-if="loading && !isRefreshing">
        <view class="loading-spinner"></view>
        <text class="loading-text">加载中...</text>
      </view>

      <!-- 加载更多 -->
      <view class="load-more" v-if="!loading && hasMore && list.length > 0" @click="onLoadMore">
        <text>加载更多</text>
      </view>

      <!-- 没有更多 -->
      <view class="no-more" v-if="!loading && !hasMore && list.length > 0">
        <text>没有更多了</text>
      </view>

      <!-- 空状态 -->
      <view class="empty-state" v-if="!loading && list.length === 0">
        <image class="empty-image" src="/static/images/empty.png" mode="aspectFit" />
        <text class="empty-text">{{ emptyText }}</text>
      </view>
    </scroll-view>
  </view>
</template>

<script>
export default {
  name: 'ListContainer',
  props: {
    list: {
      type: Array,
      default: () => []
    },
    loading: {
      type: Boolean,
      default: false
    },
    hasMore: {
      type: Boolean,
      default: false
    },
    refreshable: {
      type: Boolean,
      default: true
    },
    emptyText: {
      type: String,
      default: '暂无数据'
    }
  },
  data() {
    return {
      isRefreshing: false
    }
  },
  methods: {
    onRefresh() {
      this.isRefreshing = true
      this.$emit('refresh', () => {
        this.isRefreshing = false
      })
    },
    onLoadMore() {
      if (!this.loading && this.hasMore) {
        this.$emit('load-more')
      }
    }
  }
}
</script>

<style scoped>
.list-container {
  height: 100%;
  background-color: #f7f8fa;
}

.scroll-view {
  height: 100%;
}

.list-content {
  padding: 24rpx;
}

.loading-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48rpx 0;
}

.loading-spinner {
  width: 48rpx;
  height: 48rpx;
  border: 4rpx solid #ebedf0;
  border-top-color: #1989fa;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.loading-text {
  font-size: 28rpx;
  color: #969799;
  margin-top: 16rpx;
}

.load-more {
  text-align: center;
  padding: 24rpx 0;
  font-size: 28rpx;
  color: #969799;
}

.no-more {
  text-align: center;
  padding: 24rpx 0;
  font-size: 28rpx;
  color: #c8c9cc;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 96rpx 0;
}

.empty-image {
  width: 200rpx;
  height: 200rpx;
  margin-bottom: 24rpx;
  opacity: 0.6;
}

.empty-text {
  font-size: 28rpx;
  color: #969799;
}
</style>
