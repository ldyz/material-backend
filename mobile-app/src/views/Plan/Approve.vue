<template>
  <div class="plan-approve">
    <van-nav-bar title="审批计划" left-arrow @click-left="router.back()" />

    <van-loading v-if="loading" type="spinner" vertical />

    <van-form v-else-if="plan" @submit="handleApprove">
      <van-cell-group inset title="计划信息">
        <van-cell title="计划单号" :value="plan.plan_no" />
        <van-cell title="项目名称" :value="plan.project_name || '-'" />
        <van-cell title="物料数量" :value="`${plan.items?.length || 0} 项`" />
      </van-cell-group>

      <van-cell-group inset title="审批操作">
        <van-field
          v-model="formData.remark"
          type="textarea"
          label="审批意见"
          placeholder="请输入审批意见（可选）"
          rows="3"
        />
      </van-cell-group>

      <div class="footer-actions">
        <van-button
          round
          block
          type="danger"
          plain
          :loading="rejecting"
          @click="handleReject"
        >
          拒绝
        </van-button>
        <van-button
          round
          block
          type="primary"
          native-type="submit"
          :loading="approving"
        >
          批准
        </van-button>
      </div>
    </van-form>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showToast, showConfirmDialog } from 'vant'
import { getPlanDetail, approvePlan, rejectPlan } from '@/api/material_plan'
import { logger } from '@/utils/logger'

const router = useRouter()
const route = useRoute()

const loading = ref(true)
const approving = ref(false)
const rejecting = ref(false)
const plan = ref(null)
const formData = ref({ remark: '' })

async function loadData() {
  loading.value = true
  try {
    const response = await getPlanDetail(route.params.id)
    plan.value = response.data
  } catch (error) {
    const errorMsg = error.error || error.message || '加载失败'
    showToast({ type: 'fail', message: errorMsg })
  } finally {
    loading.value = false
  }
}

async function handleApprove() {
  approving.value = true
  try {
    await approvePlan(plan.value.id, { remark: formData.value.remark })
    showToast({ type: 'success', message: '已批准' })
    setTimeout(() => router.back(), 1500)
  } catch (error) {
    const errorMsg = error.error || error.message || '操作失败'
    // 使用对话框显示权限错误
    await showConfirmDialog({
      title: '审批失败',
      message: errorMsg,
      showCancelButton: false,
      confirmButtonText: '我知道了',
      teleport: '#app',
      confirmButtonColor: '#1989fa'
    }).catch(() => {})
  } finally {
    approving.value = false
  }
}

async function handleReject() {
  try {
    await showConfirmDialog({
      title: '拒绝确认',
      message: '确定要拒绝该计划吗？',
      teleport: '#app',
      confirmButtonColor: '#ee0a24'
    })

    if (!formData.value.remark.trim()) {
      showToast({ type: 'fail', message: '请填写拒绝原因' })
      return
    }

    rejecting.value = true
    try {
      await rejectPlan(plan.value.id, { remark: formData.value.remark })
      showToast({ type: 'success', message: '已拒绝' })
      setTimeout(() => router.back(), 1500)
    } catch (error) {
      const errorMsg = error.error || error.message || '操作失败'
      // 使用对话框显示权限错误
      await showConfirmDialog({
        title: '审批失败',
        message: errorMsg,
        showCancelButton: false,
        confirmButtonText: '我知道了',
        teleport: '#app',
        confirmButtonColor: '#1989fa'
      }).catch(() => {})
      throw error
    } finally {
      rejecting.value = false
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('拒绝失败:', error)
    }
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.plan-approve {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding-bottom: 16px;
}

.van-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 200px;
}

.footer-actions {
  padding: 16px;
  display: flex;
  gap: 12px;
}

.footer-actions .van-button {
  flex: 1;
}
</style>
