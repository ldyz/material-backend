<template>
  <view class="stock-detail-page">
    <view class="loading-wrapper" v-if="loading">
      <view class="loading-spinner"></view>
    </view>

    <view class="detail-content" v-else-if="detail">
      <DetailInfoGroup
        title="物料信息"
        :items="[
          { label: '物料名称', value: detail.material_name },
          { label: '规格型号', value: detail.specification },
          { label: '单位', value: detail.unit },
          { label: '库存数量', value: String(detail.quantity) },
          { label: '所属项目', value: detail.project_name },
          { label: '仓库位置', value: detail.warehouse_location }
        ]"
      />

      <view class="section">
        <view class="section-title">出入库记录</view>
        <view class="log-list" v-if="logs.length > 0">
          <view class="log-item" v-for="(item, index) in logs" :key="index">
            <view class="log-header">
              <text class="log-type" :class="item.type">{{ item.type === 'in' ? '入库' : '出库' }}</text>
              <text class="log-time">{{ formatDateTime(item.created_at) }}</text>
            </view>
            <view class="log-info">
              <text>数量: {{ item.quantity }}</text>
              <text>单号: {{ item.order_number || '-' }}</text>
            </view>
          </view>
        </view>
        <view class="empty-log" v-else>
          <text>暂无记录</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script>
import { getStockDetail, getStockLogs } from '@/api/stock'
import DetailInfoGroup from '@/components/DetailInfoGroup.vue'

export default {
  components: {
    DetailInfoGroup
  },
  data() {
    return {
      id: null,
      loading: true,
      detail: null,
      logs: []
    }
  },
  onLoad(options) {
    this.id = options.id
    this.fetchDetail()
  },
  methods: {
    async fetchDetail() {
      this.loading = true
      try {
        const res = await getStockDetail(this.id)
        this.detail = res.data

        const logsRes = await getStockLogs(this.id)
        this.logs = logsRes.data?.items || logsRes.data || []
      } catch (error) {
        uni.showToast({
          title: error.message || '加载失败',
          icon: 'none'
        })
      } finally {
        this.loading = false
      }
    },
    formatDateTime(dateStr) {
      if (!dateStr) return '-'
      const date = new Date(dateStr)
      return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
    }
  }
}
</script>

<style scoped>
.stock-detail-page {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.loading-wrapper {
  display: flex;
  justify-content: center;
  padding: 100rpx 0;
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
  to { transform: rotate(360deg); }
}

.detail-content {
  padding: 24rpx;
}

.section {
  background-color: #ffffff;
  border-radius: 16rpx;
  margin-bottom: 24rpx;
  overflow: hidden;
}

.section-title {
  font-size: 32rpx;
  font-weight: 500;
  color: #323233;
  padding: 24rpx;
  border-bottom: 1rpx solid #ebedf0;
}

.log-list {
  padding: 0 24rpx;
}

.log-item {
  padding: 24rpx 0;
  border-bottom: 1rpx solid #ebedf0;
}

.log-item:last-child {
  border-bottom: none;
}

.log-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 12rpx;
}

.log-type {
  font-size: 26rpx;
  padding: 4rpx 12rpx;
  border-radius: 4rpx;
}

.log-type.in {
  background-color: #e8f8ee;
  color: #07c160;
}

.log-type.out {
  background-color: #fff3e8;
  color: #ff976a;
}

.log-time {
  font-size: 24rpx;
  color: #969799;
}

.log-info {
  display: flex;
  font-size: 24rpx;
  color: #646566;
}

.log-info text:first-child {
  margin-right: 24rpx;
}

.empty-log {
  text-align: center;
  padding: 48rpx 0;
  color: #969799;
  font-size: 28rpx;
}
</style>
