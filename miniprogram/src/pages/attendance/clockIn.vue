<template>
  <view class="clock-in-page">
    <!-- 顶部时间显示 -->
    <view class="time-display">
      <text class="current-time">{{ currentTime }}</text>
      <text class="current-date">{{ currentDate }}</text>
    </view>

    <!-- 定位信息 -->
    <view class="location-card">
      <view class="location-header">
        <text class="location-title">当前位置</text>
        <text class="refresh-btn" @click="refreshLocation">刷新</text>
      </view>
      <view class="location-content" v-if="location">
        <text class="location-text">{{ location.address || `${location.latitude}, ${location.longitude}` }}</text>
      </view>
      <view class="location-loading" v-else-if="locationLoading">
        <text>获取位置中...</text>
      </view>
      <view class="location-error" v-else>
        <text>获取位置失败，请点击刷新</text>
      </view>
    </view>

    <!-- 打卡类型选择 -->
    <view class="type-card">
      <view class="type-title">选择打卡类型</view>
      <view class="type-list">
        <view
          v-for="item in attendanceTypes"
          :key="item.value"
          class="type-item"
          :class="{ active: selectedType === item.value }"
          :style="{ borderColor: item.color }"
          @click="selectedType = item.value"
        >
          <text class="type-label" :style="{ color: item.color }">{{ item.label }}</text>
        </view>
      </view>
    </view>

    <!-- 拍照区域 -->
    <view class="photo-card">
      <view class="photo-title">现场照片</view>
      <view class="photo-content">
        <view class="photo-item" v-if="photoPath">
          <image class="photo-preview" :src="photoPath" mode="aspectFill" @click="previewPhoto" />
          <view class="photo-delete" @click="photoPath = ''">×</view>
        </view>
        <view class="photo-add" v-else @click="takePhoto">
          <text class="photo-add-icon">+</text>
          <text class="photo-add-text">拍照打卡</text>
        </view>
      </view>
    </view>

    <!-- 打卡按钮 -->
    <view class="clock-in-btn-wrapper">
      <button
        class="clock-in-btn"
        :class="{ disabled: !canClockIn }"
        :loading="clockInLoading"
        @click="handleClockIn"
      >
        打卡
      </button>
    </view>
  </view>
</template>

<script>
import { clockIn as clockInApi, getTodayAppointments } from '@/api/attendance'
import { getLocation, takePhoto as takePhotoUtil, uploadImage, previewImage } from '@/utils/device'

export default {
  data() {
    return {
      currentTime: '',
      currentDate: '',
      location: null,
      locationLoading: false,
      selectedType: 'morning',
      attendanceTypes: [
        { value: 'morning', label: '上午打卡', color: '#1989fa' },
        { value: 'afternoon', label: '下午打卡', color: '#07c160' },
        { value: 'noon_overtime', label: '中午加班', color: '#ff976a' },
        { value: 'night_overtime', label: '晚上加班', color: '#7232dd' }
      ],
      photoPath: '',
      clockInLoading: false,
      timer: null
    }
  },
  computed: {
    canClockIn() {
      return this.location && !this.clockInLoading
    }
  },
  onLoad() {
    this.updateTime()
    this.timer = setInterval(this.updateTime, 1000)
    this.refreshLocation()
  },
  onUnload() {
    if (this.timer) {
      clearInterval(this.timer)
    }
  },
  methods: {
    updateTime() {
      const now = new Date()
      this.currentTime = `${String(now.getHours()).padStart(2, '0')}:${String(now.getMinutes()).padStart(2, '0')}:${String(now.getSeconds()).padStart(2, '0')}`

      const weekDays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
      this.currentDate = `${now.getFullYear()}年${now.getMonth() + 1}月${now.getDate()}日 ${weekDays[now.getDay()]}`
    },
    async refreshLocation() {
      this.locationLoading = true
      try {
        this.location = await getLocation()
      } catch (error) {
        uni.showToast({
          title: error.message || '获取位置失败',
          icon: 'none'
        })
      } finally {
        this.locationLoading = false
      }
    },
    async takePhoto() {
      try {
        const result = await takePhotoUtil()
        if (result.tempFilePaths && result.tempFilePaths.length > 0) {
          this.photoPath = result.tempFilePaths[0]
        }
      } catch (error) {
        uni.showToast({
          title: error.message || '拍照失败',
          icon: 'none'
        })
      }
    },
    previewPhoto() {
      if (this.photoPath) {
        previewImage([this.photoPath])
      }
    },
    async handleClockIn() {
      if (!this.canClockIn) return

      this.clockInLoading = true
      try {
        // 上传照片
        let photoUrl = ''
        if (this.photoPath) {
          const uploadRes = await uploadImage(this.photoPath)
          photoUrl = uploadRes.data?.url || ''
        }

        // 提交打卡
        await clockInApi({
          type: this.selectedType,
          latitude: this.location.latitude,
          longitude: this.location.longitude,
          address: this.location.address || '',
          photo_url: photoUrl
        })

        uni.showToast({
          title: '打卡成功',
          icon: 'success'
        })

        this.photoPath = ''
      } catch (error) {
        uni.showToast({
          title: error.message || '打卡失败',
          icon: 'none'
        })
      } finally {
        this.clockInLoading = false
      }
    }
  }
}
</script>

<style scoped>
.clock-in-page {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.time-display {
  background: linear-gradient(135deg, #1989fa 0%, #07c160 100%);
  padding: 48rpx 32rpx;
  text-align: center;
}

.current-time {
  display: block;
  font-size: 80rpx;
  font-weight: 600;
  color: #ffffff;
  letter-spacing: 4rpx;
}

.current-date {
  display: block;
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.8);
  margin-top: 16rpx;
}

.location-card,
.type-card,
.photo-card {
  background-color: #ffffff;
  margin: 24rpx;
  border-radius: 16rpx;
  padding: 24rpx;
}

.location-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16rpx;
}

.location-title {
  font-size: 32rpx;
  font-weight: 500;
  color: #323233;
}

.refresh-btn {
  font-size: 26rpx;
  color: #1989fa;
}

.location-text {
  font-size: 28rpx;
  color: #646566;
}

.location-loading,
.location-error {
  font-size: 28rpx;
  color: #969799;
}

.type-title,
.photo-title {
  font-size: 32rpx;
  font-weight: 500;
  color: #323233;
  margin-bottom: 24rpx;
}

.type-list {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24rpx;
}

.type-item {
  padding: 24rpx;
  border: 2rpx solid #ebedf0;
  border-radius: 16rpx;
  text-align: center;
}

.type-item.active {
  background-color: #f7f8fa;
}

.type-label {
  font-size: 28rpx;
  font-weight: 500;
}

.photo-content {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
}

.photo-item,
.photo-add {
  width: 200rpx;
  height: 200rpx;
  border-radius: 12rpx;
  position: relative;
}

.photo-item {
  overflow: hidden;
}

.photo-preview {
  width: 100%;
  height: 100%;
}

.photo-delete {
  position: absolute;
  top: 8rpx;
  right: 8rpx;
  width: 40rpx;
  height: 40rpx;
  background-color: rgba(0, 0, 0, 0.5);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #ffffff;
  font-size: 28rpx;
}

.photo-add {
  border: 2rpx dashed #dcdee0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.photo-add-icon {
  font-size: 48rpx;
  color: #dcdee0;
}

.photo-add-text {
  font-size: 24rpx;
  color: #969799;
  margin-top: 8rpx;
}

.clock-in-btn-wrapper {
  padding: 48rpx 24rpx;
  padding-bottom: constant(safe-area-inset-bottom);
  padding-bottom: env(safe-area-inset-bottom);
}

.clock-in-btn {
  width: 100%;
  height: 96rpx;
  background-color: #1989fa;
  border-radius: 48rpx;
  color: #ffffff;
  font-size: 36rpx;
  font-weight: 500;
}

.clock-in-btn.disabled {
  background-color: #c8c9cc;
}
</style>
