<template>
  <view class="inbound-detail-page">
    <view class="loading-wrapper" v-if="loading">
      <view class="loading-spinner"></view>
    </view>

    <view class="detail-content" v-else-if="detail">
      <DetailInfoGroup
        title="基本信息"
        :items="[
          { label: '入库单号', value: detail.order_number },
          { label: '入库类型', value: detail.type },
          { label: '项目', value: detail.project_name },
          { label: '状态', value: getStatusLabel(detail.status), valueClass: 'highlight' },
          { label: '创建时间', value: formatDateTime(detail.created_at) },
          { label: '创建人', value: detail.creator_name }
        ]"
      />

      <view class="section">
        <view class="section-title">物料清单</view>
        <view class="material-list">
          <view class="material-item" v-for="(item, index) in detail.items" :key="index">
            <view class="material-header">
              <text class="material-name">{{ item.material_name }}</text>
              <text class="material-spec">{{ item.specification }}</text>
            </view>
            <view class="material-info">
              <text class="material-unit">单位: {{ item.unit || '-' }}</text>
              <text class="material-quantity">数量: {{ item.quantity }}</text>
            </view>
          </view>
        </view>
      </view>

      <ApprovalTimeline
        title="审批记录"
        :list="workflowHistory"
      />
    </view>

    <view class="bottom-bar" v-if="showActions">
      <button class="btn reject-btn" v-if="canApprove" @click="handleReject">拒绝</button>
      <button class="btn approve-btn" v-if="canApprove" @click="handleApprove">通过</button>
    </view>
  </view>
</template>

<script>
import { getInboundDetail, approveInbound, rejectInbound } from '@/api/inbound'
import DetailInfoGroup from '@/components/DetailInfoGroup.vue'
import ApprovalTimeline from '@/components/ApprovalTimeline.vue'

export default {
  components: {
    DetailInfoGroup,
    ApprovalTimeline
  },
  data() {
    return {
      id: null,
      loading: true,
      detail: null,
      workflowHistory: []
    }
  },
  computed: {
    canApprove() {
      return this.detail?.status === 'pending'
    },
    showActions() {
      return this.canApprove
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
        const res = await getInboundDetail(this.id)
        this.detail = res.data
        this.workflowHistory = res.data?.workflow_history || []
      } catch (error) {
        uni.showToast({
          title: error.message || '加载失败',
          icon: 'none'
        })
      } finally {
        this.loading = false
      }
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
    formatDateTime(dateStr) {
      if (!dateStr) return '-'
      const date = new Date(dateStr)
      return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
    },
    async handleApprove() {
      try {
        await approveInbound(this.id, {})
        uni.showToast({ title: '审批通过', icon: 'success' })
        this.fetchDetail()
      } catch (error) {
        uni.showToast({ title: error.message || '操作失败', icon: 'none' })
      }
    },
    async handleReject() {
      try {
        const res = await new Promise((resolve) => {
          uni.showModal({
            title: '拒绝原因',
            editable: true,
            placeholderText: '请输入拒绝原因',
            success: resolve
          })
        })

        if (res.confirm) {
          await rejectInbound(this.id, { remark: res.content || '' })
          uni.showToast({ title: '已拒绝', icon: 'success' })
          this.fetchDetail()
        }
      } catch (error) {
        uni.showToast({ title: error.message || '操作失败', icon: 'none' })
      }
    }
  }
}
</script>

<style scoped>
.inbound-detail-page {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding-bottom: 120rpx;
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

.material-list {
  padding: 0 24rpx;
}

.material-item {
  padding: 24rpx 0;
  border-bottom: 1rpx solid #ebedf0;
}

.material-item:last-child {
  border-bottom: none;
}

.material-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 12rpx;
}

.material-name {
  font-size: 28rpx;
  color: #323233;
  font-weight: 500;
}

.material-spec {
  font-size: 26rpx;
  color: #969799;
}

.material-info {
  display: flex;
  font-size: 24rpx;
  color: #646566;
}

.material-unit {
  margin-right: 24rpx;
}

.bottom-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  padding: 24rpx;
  background-color: #ffffff;
  border-top: 1rpx solid #ebedf0;
  padding-bottom: constant(safe-area-inset-bottom);
  padding-bottom: env(safe-area-inset-bottom);
}

.btn {
  flex: 1;
  height: 88rpx;
  border-radius: 16rpx;
  font-size: 32rpx;
  margin: 0 12rpx;
}

.approve-btn {
  background-color: #07c160;
  color: #ffffff;
}

.reject-btn {
  background-color: #ee0a24;
  color: #ffffff;
}
</style>
