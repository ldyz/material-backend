<template>
  <view class="approval-timeline">
    <view class="timeline-title" v-if="title">{{ title }}</view>
    <view class="timeline-list">
      <view
        v-for="(item, index) in list"
        :key="index"
        class="timeline-item"
        :class="{ 'is-last': index === list.length - 1 }"
      >
        <view class="timeline-line">
          <view class="timeline-dot" :class="getDotClass(item)"></view>
        </view>
        <view class="timeline-content">
          <view class="timeline-header">
            <text class="timeline-title">{{ item.title || item.action }}</text>
            <text class="timeline-status" :class="getStatusClass(item)">{{ getStatusText(item) }}</text>
          </view>
          <view class="timeline-info" v-if="item.operator_name || item.operator">
            <text class="timeline-user">{{ item.operator_name || item.operator }}</text>
            <text class="timeline-time" v-if="item.created_at || item.operated_at">{{ formatTime(item.created_at || item.operated_at) }}</text>
          </view>
          <view class="timeline-remark" v-if="item.remark || item.comment">
            <text>{{ item.remark || item.comment }}</text>
          </view>
        </view>
      </view>
    </view>
    <view class="timeline-empty" v-if="!list || list.length === 0">
      <text>暂无审批记录</text>
    </view>
  </view>
</template>

<script>
export default {
  name: 'ApprovalTimeline',
  props: {
    title: {
      type: String,
      default: ''
    },
    list: {
      type: Array,
      default: () => []
    }
  },
  methods: {
    getDotClass(item) {
      const status = item.status || item.action
      if (status === 'approved' || status === 'approve') {
        return 'dot-success'
      }
      if (status === 'rejected' || status === 'reject') {
        return 'dot-danger'
      }
      if (status === 'pending') {
        return 'dot-warning'
      }
      return ''
    },
    getStatusClass(item) {
      const status = item.status || item.action
      if (status === 'approved' || status === 'approve') {
        return 'status-success'
      }
      if (status === 'rejected' || status === 'reject') {
        return 'status-danger'
      }
      if (status === 'pending') {
        return 'status-warning'
      }
      return ''
    },
    getStatusText(item) {
      const status = item.status || item.action
      const statusMap = {
        'approved': '已通过',
        'approve': '已通过',
        'rejected': '已拒绝',
        'reject': '已拒绝',
        'pending': '待审批',
        'submitted': '已提交',
        'draft': '草稿',
        'created': '已创建',
        'resubmitted': '已重新提交'
      }
      return statusMap[status] || status
    },
    formatTime(time) {
      if (!time) return ''
      const date = new Date(time)
      const year = date.getFullYear()
      const month = String(date.getMonth() + 1).padStart(2, '0')
      const day = String(date.getDate()).padStart(2, '0')
      const hour = String(date.getHours()).padStart(2, '0')
      const minute = String(date.getMinutes()).padStart(2, '0')
      return `${year}-${month}-${day} ${hour}:${minute}`
    }
  }
}
</script>

<style scoped>
.approval-timeline {
  background-color: #ffffff;
  border-radius: 16rpx;
  padding: 24rpx;
}

.timeline-title {
  font-size: 32rpx;
  font-weight: 500;
  color: #323233;
  margin-bottom: 24rpx;
}

.timeline-list {
  position: relative;
}

.timeline-item {
  display: flex;
  padding-bottom: 24rpx;
  position: relative;
}

.timeline-item.is-last {
  padding-bottom: 0;
}

.timeline-line {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 40rpx;
  margin-right: 16rpx;
  position: relative;
}

.timeline-line::before {
  content: '';
  position: absolute;
  top: 20rpx;
  left: 50%;
  transform: translateX(-50%);
  width: 2rpx;
  height: calc(100% - 20rpx);
  background-color: #ebedf0;
}

.timeline-item.is-last .timeline-line::before {
  display: none;
}

.timeline-dot {
  width: 24rpx;
  height: 24rpx;
  border-radius: 50%;
  background-color: #969799;
  border: 4rpx solid #ffffff;
  box-shadow: 0 0 0 2rpx #ebedf0;
  z-index: 1;
}

.timeline-dot.dot-success {
  background-color: #07c160;
  box-shadow: 0 0 0 2rpx #07c160;
}

.timeline-dot.dot-danger {
  background-color: #ee0a24;
  box-shadow: 0 0 0 2rpx #ee0a24;
}

.timeline-dot.dot-warning {
  background-color: #ff976a;
  box-shadow: 0 0 0 2rpx #ff976a;
}

.timeline-content {
  flex: 1;
  min-width: 0;
}

.timeline-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8rpx;
}

.timeline-title {
  font-size: 28rpx;
  color: #323233;
  font-weight: 500;
}

.timeline-status {
  font-size: 24rpx;
}

.timeline-status.status-success {
  color: #07c160;
}

.timeline-status.status-danger {
  color: #ee0a24;
}

.timeline-status.status-warning {
  color: #ff976a;
}

.timeline-info {
  display: flex;
  align-items: center;
  font-size: 24rpx;
  color: #969799;
  margin-bottom: 8rpx;
}

.timeline-user {
  margin-right: 16rpx;
}

.timeline-remark {
  font-size: 26rpx;
  color: #646566;
  background-color: #f7f8fa;
  padding: 12rpx 16rpx;
  border-radius: 8rpx;
  margin-top: 8rpx;
}

.timeline-empty {
  text-align: center;
  padding: 32rpx 0;
  color: #969799;
  font-size: 28rpx;
}
</style>
