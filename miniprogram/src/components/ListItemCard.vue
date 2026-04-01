<template>
  <view class="list-item-card" @click="handleClick">
    <!-- 头部 -->
    <view class="card-header" v-if="title || $slots.header">
      <text class="card-title" v-if="title">{{ title }}</text>
      <slot name="header"></slot>
      <view class="card-status" v-if="status">
        <StatusTag :text="statusText" :color="statusColor" plain />
      </view>
    </view>

    <!-- 内容 -->
    <view class="card-content">
      <slot>
        <view class="info-row" v-for="(item, index) in items" :key="index">
          <text class="info-label">{{ item.label }}</text>
          <text class="info-value">{{ item.value || '-' }}</text>
        </view>
      </slot>
    </view>

    <!-- 底部 -->
    <view class="card-footer" v-if="$slots.footer || showActions">
      <slot name="footer"></slot>
      <view class="card-actions" v-if="showActions">
        <view
          v-for="(action, index) in actions"
          :key="index"
          class="action-btn"
          :class="action.type"
          @click.stop="handleAction(action)"
        >
          <text>{{ action.text }}</text>
        </view>
      </view>
    </view>

    <!-- 箭头 -->
    <view class="card-arrow" v-if="clickable">
      <text>›</text>
    </view>
  </view>
</template>

<script>
import StatusTag from './StatusTag.vue'

export default {
  name: 'ListItemCard',
  components: {
    StatusTag
  },
  props: {
    title: {
      type: String,
      default: ''
    },
    items: {
      type: Array,
      default: () => []
    },
    status: {
      type: String,
      default: ''
    },
    statusText: {
      type: String,
      default: ''
    },
    statusColor: {
      type: String,
      default: 'default'
    },
    clickable: {
      type: Boolean,
      default: false
    },
    showActions: {
      type: Boolean,
      default: false
    },
    actions: {
      type: Array,
      default: () => []
    }
  },
  methods: {
    handleClick() {
      if (this.clickable) {
        this.$emit('click')
      }
    },
    handleAction(action) {
      this.$emit('action', action)
    }
  }
}
</script>

<style scoped>
.list-item-card {
  background-color: #ffffff;
  border-radius: 16rpx;
  margin-bottom: 24rpx;
  padding: 24rpx;
  position: relative;
  overflow: hidden;
}

.list-item-card:active {
  background-color: #f7f8fa;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16rpx;
}

.card-title {
  font-size: 32rpx;
  font-weight: 500;
  color: #323233;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.card-status {
  flex-shrink: 0;
}

.card-content {
  margin-bottom: 16rpx;
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8rpx 0;
}

.info-label {
  font-size: 28rpx;
  color: #969799;
}

.info-value {
  font-size: 28rpx;
  color: #323233;
  max-width: 60%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.card-footer {
  border-top: 1rpx solid #ebedf0;
  padding-top: 16rpx;
  margin-top: 8rpx;
}

.card-actions {
  display: flex;
  justify-content: flex-end;
  gap: 16rpx;
}

.action-btn {
  padding: 12rpx 24rpx;
  border-radius: 8rpx;
  font-size: 26rpx;
}

.action-btn.primary {
  background-color: #1989fa;
  color: #ffffff;
}

.action-btn.danger {
  background-color: #ee0a24;
  color: #ffffff;
}

.action-btn.default {
  background-color: #f7f8fa;
  color: #646566;
}

.card-arrow {
  position: absolute;
  right: 24rpx;
  top: 50%;
  transform: translateY(-50%);
  font-size: 36rpx;
  color: #c8c9cc;
}
</style>
