<template>
  <view class="appointment-detail-page">
    <view class="loading-wrapper" v-if="loading">
      <view class="loading-spinner"></view>
    </view>

    <view class="detail-content" v-else-if="detail">
      <DetailInfoGroup
        title="基本信息"
        :items="[
          { label: '预约单号', value: detail.appointment_number },
          { label: '标题', value: detail.title },
          { label: '作业日期', value: formatDate(detail.work_date) },
          { label: '时间段', value: getTimeSlotLabel(detail.time_slot) },
          { label: '作业地点', value: detail.location },
          { label: '状态', value: getStatusLabel(detail.status), valueClass: 'highlight' },
          { label: '申请人', value: detail.applicant_name },
          { label: '创建时间', value: formatDateTime(detail.created_at) }
        ]"
      />

      <DetailInfoGroup
        title="作业描述"
        :items="[
          { label: '作业内容', value: detail.description || '-' }
        ]"
      />

      <ApprovalTimeline
        title="审批记录"
        :list="workflowHistory"
      />
    </view>

    <view class="bottom-bar" v-if="showActions">
      <button class="btn cancel-btn" v-if="canCancel" @click="handleCancel">取消</button>
      <button class="btn complete-btn" v-if="canComplete" @click="handleComplete">完成</button>
    </view>
  </view>
</template>

<script>
import { getAppointmentDetail, cancelAppointment, completeAppointment, getStatusLabel, getTimeSlotLabel } from '@/api/appointment'
import DetailInfoGroup from '@/components/DetailInfoGroup.vue'
import ApprovalTimeline from '@/components/ApprovalTimeline.vue'
import { useAuthStore } from '@/stores/auth'

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
    canCancel() {
      return ['pending', 'scheduled', 'draft'].includes(this.detail?.status)
    },
    canComplete() {
      return ['in_progress', 'scheduled'].includes(this.detail?.status)
    },
    showActions() {
      return this.canCancel || this.canComplete
    }
  },
  onLoad(options) {
    this.id = options.id
    this.fetchDetail()
  },
  methods: {
    getStatusLabel,
    getTimeSlotLabel,
    async fetchDetail() {
      this.loading = true
      try {
        const res = await getAppointmentDetail(this.id)
        this.detail = res.data
        this.workflowHistory = res.data?.approval_history || []
      } catch (error) {
        uni.showToast({
          title: error.message || '加载失败',
          icon: 'none'
        })
      } finally {
        this.loading = false
      }
    },
    formatDate(dateStr) {
      if (!dateStr) return '-'
      const date = new Date(dateStr)
      return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
    },
    formatDateTime(dateStr) {
      if (!dateStr) return '-'
      const date = new Date(dateStr)
      return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
    },
    async handleCancel() {
      try {
        const res = await new Promise((resolve) => {
          uni.showModal({
            title: '确认取消',
            content: '确定要取消此预约吗？',
            success: resolve
          })
        })

        if (res.confirm) {
          await cancelAppointment(this.id, {})
          uni.showToast({ title: '已取消', icon: 'success' })
          this.fetchDetail()
        }
      } catch (error) {
        uni.showToast({ title: error.message || '操作失败', icon: 'none' })
      }
    },
    async handleComplete() {
      try {
        await completeAppointment(this.id, {})
        uni.showToast({ title: '已完成', icon: 'success' })
        this.fetchDetail()
      } catch (error) {
        uni.showToast({ title: error.message || '操作失败', icon: 'none' })
      }
    }
  }
}
</script>

<style scoped>
.appointment-detail-page {
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

.complete-btn {
  background-color: #07c160;
  color: #ffffff;
}

.cancel-btn {
  background-color: #ee0a24;
  color: #ffffff;
}
</style>
