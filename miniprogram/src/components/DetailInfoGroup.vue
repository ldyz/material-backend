<template>
  <view class="detail-info-group">
    <view class="group-title" v-if="title">{{ title }}</view>
    <view class="info-list">
      <view
        v-for="(item, index) in items"
        :key="index"
        class="info-item"
        :class="{ 'is-link': item.clickable || item.isLink }"
        @click="handleClick(item)"
      >
        <text class="info-label">{{ item.label }}</text>
        <view class="info-value-wrapper">
          <text class="info-value" :class="item.valueClass">{{ item.value || '-' }}</text>
          <text class="info-arrow" v-if="item.clickable || item.isLink">›</text>
        </view>
      </view>
    </view>
    <slot></slot>
  </view>
</template>

<script>
export default {
  name: 'DetailInfoGroup',
  props: {
    title: {
      type: String,
      default: ''
    },
    items: {
      type: Array,
      default: () => []
    }
  },
  methods: {
    handleClick(item) {
      if (item.clickable || item.isLink) {
        this.$emit('click', item)
        if (item.onClick) {
          item.onClick()
        }
      }
    }
  }
}
</script>

<style scoped>
.detail-info-group {
  background-color: #ffffff;
  border-radius: 16rpx;
  overflow: hidden;
  margin-bottom: 24rpx;
}

.group-title {
  font-size: 32rpx;
  font-weight: 500;
  color: #323233;
  padding: 24rpx 24rpx 16rpx;
  border-bottom: 1rpx solid #ebedf0;
}

.info-list {
  padding: 0 24rpx;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20rpx 0;
  border-bottom: 1rpx solid #ebedf0;
}

.info-item:last-child {
  border-bottom: none;
}

.info-item.is-link {
  cursor: pointer;
}

.info-item.is-link:active {
  background-color: #f7f8fa;
}

.info-label {
  font-size: 28rpx;
  color: #646566;
  flex-shrink: 0;
}

.info-value-wrapper {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  flex: 1;
  min-width: 0;
  margin-left: 24rpx;
}

.info-value {
  font-size: 28rpx;
  color: #323233;
  text-align: right;
  word-break: break-all;
}

.info-value.highlight {
  color: #1989fa;
}

.info-value.success {
  color: #07c160;
}

.info-value.warning {
  color: #ff976a;
}

.info-value.danger {
  color: #ee0a24;
}

.info-arrow {
  font-size: 32rpx;
  color: #969799;
  margin-left: 8rpx;
}
</style>
