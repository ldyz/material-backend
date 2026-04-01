<template>
  <view class="create-page">
    <view class="form-card">
      <view class="form-item">
        <text class="form-label">领料类型</text>
        <input class="form-input" v-model="form.type" placeholder="请输入领料类型" />
      </view>
      <view class="form-item">
        <text class="form-label">所属项目</text>
        <picker :value="projectIndex" :range="projectList" range-key="name" @change="onProjectChange">
          <view class="form-picker">
            <text>{{ projectList[projectIndex]?.name || '请选择项目' }}</text>
            <text class="arrow">›</text>
          </view>
        </picker>
      </view>
      <view class="form-item">
        <text class="form-label">备注</text>
        <textarea class="form-textarea" v-model="form.remark" placeholder="请输入备注" />
      </view>
    </view>

    <view class="submit-btn-wrapper">
      <button class="submit-btn" :loading="loading" @click="handleSubmit">提交</button>
    </view>
  </view>
</template>

<script>
import { createRequisition } from '@/api/requisition'
import { getProjects } from '@/api/project'

export default {
  data() {
    return {
      loading: false,
      projectList: [],
      projectIndex: -1,
      form: {
        type: '',
        project_id: null,
        remark: ''
      }
    }
  },
  onLoad() {
    this.fetchProjects()
  },
  methods: {
    async fetchProjects() {
      try {
        const res = await getProjects({ page_size: 100 })
        this.projectList = res.data?.items || res.data || []
      } catch (error) {
        uni.showToast({ title: error.message || '加载项目失败', icon: 'none' })
      }
    },
    onProjectChange(e) {
      this.projectIndex = e.detail.value
      this.form.project_id = this.projectList[this.projectIndex]?.id
    },
    async handleSubmit() {
      if (!this.form.type) {
        uni.showToast({ title: '请输入领料类型', icon: 'none' })
        return
      }
      if (!this.form.project_id) {
        uni.showToast({ title: '请选择项目', icon: 'none' })
        return
      }

      this.loading = true
      try {
        await createRequisition(this.form)
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
