<template>
  <div class="req-approve">
    <van-nav-bar title="审批出库" left-arrow @click-left="router.back()" />

    <van-loading v-if="loading" type="spinner" vertical />

    <van-form v-else-if="requisition" @submit="handleApprove">
      <van-cell-group inset title="出库信息">
        <van-cell title="出库单号" :value="requisition.requisition_number" />
        <van-cell title="项目名称" :value="requisition.project_name || '-'" />
        <van-cell title="申请部门" :value="requisition.department_name || '-'" />
        <van-cell title="物料数量" :value="`${requisition.items?.length || 0} 项`" />
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
          @click="handleReject"
        >
          拒绝
        </van-button>
        <van-button
          round
          block
          type="primary"
          native-type="submit"
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
import { getRequisitionDetail, approveRequisition, rejectRequisition } from '@/api/requisition'

const router = useRouter()
const route = useRoute()

const loading = ref(true)
const requisition = ref(null)
const formData = ref({ remark: '' })

async function loadData() {
  loading.value = true
  try {
    const response = await getRequisitionDetail(route.params.id)
    requisition.value = response.data
  } catch (error) {
    showToast({ type: 'fail', message: '加载失败' })
  } finally {
    loading.value = false
  }
}

async function handleApprove() {
  try {
    await approveRequisition(requisition.value.id, { approval_comment: formData.value.remark })
    showToast({ type: 'success', message: '已批准' })
    setTimeout(() => router.back(), 1000)
  } catch (error) {
    showToast({ type: 'fail', message: error.message || '操作失败' })
  }
}

async function handleReject() {
  try {
    await showConfirmDialog({
      title: '拒绝确认',
      message: '确定要拒绝该出库单吗？'
    })

    if (!formData.value.remark.trim()) {
      showToast({ type: 'fail', message: '请填写拒绝原因' })
      return
    }

    await rejectRequisition(requisition.value.id, { rejection_reason: formData.value.remark })
    showToast({ type: 'success', message: '已拒绝' })
    setTimeout(() => router.back(), 1000)
  } catch (error) {
    if (error !== 'cancel') {
      showToast({ type: 'fail', message: error.message || '操作失败' })
    }
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.req-approve {
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
