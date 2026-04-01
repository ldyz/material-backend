<template>
  <view class="create-page">
    <view class="form-card">
      <view class="form-item">
        <text class="form-label">标题</text>
        <input class="form-input" v-model="form.title" placeholder="请输入预约标题" />
      </view>
      <view class="form-item">
        <text class="form-label">作业日期</text>
        <picker mode="date" :value="form.work_date" @change="onDateChange">
          <view class="form-picker">
            <text>{{ form.work_date || '请选择日期' }}</text>
            <text class="arrow">›</text>
          </view>
        </picker>
      </view>
      <view class="form-item">
        <text class="form-label">时间段</text>
        <picker :value="timeSlotIndex" :range="timeSlotOptions" range-key="label" @change="onTimeSlotChange">
          <view class="form-picker">
            <text>{{ timeSlotOptions[timeSlotIndex]?.label || '请选择时间段' }}</text>
            <text class="arrow">›</text>
          </view>
        </picker>
      </view>
      <view class="form-item">
        <text class="form-label">作业地点</text>
        <input class="form-input" v-model="form.location" placeholder="请输入作业地点" />
      </view>
      <view class="form-item">
        <text class="form-label">作业描述</text>
        <textarea class="form-textarea" v-model="form.description" placeholder="请输入作业描述" />
      </view>
    </view>

    <view class="submit-btn-wrapper">
      <button class="submit-btn" :loading="loading" @click="handleSubmit">提交</button>
    </view>
  </view>
</template>

<script>
import { createAppointment, getTimeSlotOptions } from '@/api/appointment'

export default {
  data() {
    return {
      loading: false,
      timeSlotOptions: getTimeSlotOptions(),
      timeSlotIndex: 0,
      form: {
        title: '',
        work_date: '',
        time_slot: 'morning',
        location: '',
        description: ''
      }
    }
  },
  methods: {
    onDateChange(e) {
      this.form.work_date = e.detail.value
    },
    onTimeSlotChange(e) {
      this.timeSlotIndex = e.detail.value
      this.form.time_slot = this.timeSlotOptions[this.timeSlotIndex].value
    },
    async handleSubmit() {
      if (!this.form.title) {
        uni.showToast({ title: '请输入标题', icon: 'none' })
        return
      }
      if (!this.form.work_date) {
        uni.showToast({ title: '请选择作业日期', icon: 'none' })
        return
      }
      if (!this.form.location) {
        uni.showToast({ title: '请输入作业地点', icon: 'none' })
        return
      }

      this.loading = true
      try {
        await createAppointment(this.form)
        uni.showToast({ title: '创建成功', icon: 'success' })
        setTimeout(() => {
          uni.navigateBack()
        }, 1000)
      } catch (error) {
        uni.showToast({ title: error.message || '创建失败', icon: 'none' })
      } finally {
        this.loading = false
      }
    }
  }
}
</script>

<style scoped>
.create-page {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding: 24rpx;
}

.form-card {
  background-color: #ffffff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-bottom: 24rpx;
}

.form-item {
  margin-bottom: 24rpx;
}

.form-label {
  display: block;
  font-size: 28rpx;
  color: #323233;
  margin-bottom: 12rpx;
}

.form-input,
.form-textarea {
  width: 100%;
  background-color: #f7f8fa;
  border-radius: 8rpx;
  padding: 24rpx;
  font-size: 28rpx;
}

.form-textarea {
  min-height: 160rpx;
}

.form-picker {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #f7f8fa;
  border-radius: 8rpx;
  padding: 24rpx;
  font-size: 28rpx;
  color: #323233;
}

.arrow {
  color: #c8c9cc;
}

.submit-btn-wrapper {
  padding: 24rpx;
}

.submit-btn {
  width: 100%;
  height: 96rpx;
  background-color: #1989fa;
  border-radius: 16rpx;
  color: #ffffff;
  font-size: 32rpx;
}
</style>
