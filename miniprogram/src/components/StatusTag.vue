<template>
  <view class="status-tag" :class="[`status-${color}`, { 'is-plain': plain }]" :style="customStyle">
    <text>{{ text }}</text>
  </view>
</template>

<script>
export default {
  name: 'StatusTag',
  props: {
    text: {
      type: String,
      default: ''
    },
    color: {
      type: String,
      default: 'default' // default, primary, success, warning, danger, info
    },
    plain: {
      type: Boolean,
      default: false
    },
    customStyle: {
      type: String,
      default: ''
    }
  }
}
</script>

<style scoped>
.status-tag {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 4rpx 16rpx;
  border-radius: 8rpx;
  font-size: 24rpx;
  line-height: 1.4;
}

/* 默认样式 */
.status-default {
  background-color: #f2f3f5;
  color: #646566;
}

.status-default.is-plain {
  background-color: transparent;
  border: 1px solid #ebedf0;
}

/* 主要样式 */
.status-primary {
  background-color: #1989fa;
  color: #ffffff;
}

.status-primary.is-plain {
  background-color: #e8f4ff;
  color: #1989fa;
}

/* 成功样式 */
.status-success {
  background-color: #07c160;
  color: #ffffff;
}

.status-success.is-plain {
  background-color: #e8f8ee;
  color: #07c160;
}

/* 警告样式 */
.status-warning {
  background-color: #ff976a;
  color: #ffffff;
}

.status-warning.is-plain {
  background-color: #fff3e8;
  color: #ff976a;
}

/* 危险样式 */
.status-danger {
  background-color: #ee0a24;
  color: #ffffff;
}

.status-danger.is-plain {
  background-color: #ffebee;
  color: #ee0a24;
}

/* 信息样式 */
.status-info {
  background-color: #909399;
  color: #ffffff;
}

.status-info.is-plain {
  background-color: #f4f4f5;
  color: #909399;
}
</style>
