<template>
  <view class="about-page">
    <view class="about-header">
      <image class="app-logo" src="/static/images/logo.png" mode="aspectFit" />
      <text class="app-name">项目管理系统</text>
      <text class="app-version">版本 1.0.0</text>
    </view>

    <view class="about-content">
      <view class="info-card">
        <view class="info-item">
          <text class="info-label">应用名称</text>
          <text class="info-value">项目管理系统</text>
        </view>
        <view class="info-item">
          <text class="info-label">应用版本</text>
          <text class="info-value">1.0.0</text>
        </view>
        <view class="info-item">
          <text class="info-label">技术支持</text>
          <text class="info-value">技术部</text>
        </view>
      </view>

      <view class="desc-card">
        <text class="desc-text">
          项目管理系统是一款面向工程项目管理的综合性应用，包含物资计划、入库管理、领料管理、预约管理、考勤打卡等功能模块，助力项目高效管理。
        </text>
      </view>
    </view>
  </view>
</template>

<script>
export default {
  name: 'AboutPage'
}
</script>

<style scoped>
.about-page {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.about-header {
  background-color: #ffffff;
  padding: 64rpx 32rpx;
  text-align: center;
  margin-bottom: 24rpx;
}

.app-logo {
  width: 160rpx;
  height: 160rpx;
  margin-bottom: 24rpx;
}

.app-name {
  display: block;
  font-size: 40rpx;
  font-weight: 600;
  color: #323233;
  margin-bottom: 12rpx;
}

.app-version {
  display: block;
  font-size: 26rpx;
  color: #969799;
}

.about-content {
  padding: 0 24rpx;
}

.info-card {
  background-color: #ffffff;
  border-radius: 16rpx;
  margin-bottom: 24rpx;
}

.info-item {
  display: flex;
  justify-content: space-between;
  padding: 32rpx;
  border-bottom: 1rpx solid #ebedf0;
}

.info-item:last-child {
  border-bottom: none;
}

.info-label {
  font-size: 28rpx;
  color: #646566;
}

.info-value {
  font-size: 28rpx;
  color: #323233;
}

.desc-card {
  background-color: #ffffff;
  border-radius: 16rpx;
  padding: 32rpx;
}

.desc-text {
  font-size: 28rpx;
  color: #646566;
  line-height: 1.8;
}
</style>
