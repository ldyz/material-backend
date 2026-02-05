<template>
  <div class="req-issue">
    <van-nav-bar title="发放出库" left-arrow @click-left="router.back()" />

    <van-loading v-if="loading" type="spinner" vertical />

    <van-form v-else-if="requisition" @submit="handleSubmit">
      <van-cell-group inset title="出库信息">
        <van-cell title="出库单号" :value="requisition.requisition_number" />
        <van-cell title="项目名称" :value="requisition.project_name || '-'" />
      </van-cell-group>

      <van-cell-group inset title="物料发放">
        <div
          v-for="item in formData.items"
          :key="item.id"
          class="item-row"
        >
          <div class="item-header">
            <span class="item-name">{{ item.material_name }}</span>
            <van-tag type="primary">{{ item.requested_quantity }}</van-tag>
          </div>
          <van-field
            v-model.number="item.issued_quantity"
            type="number"
            label="发放数量"
            placeholder="请输入发放数量"
          />
        </div>
      </van-cell-group>

      <van-cell-group inset title="发放备注">
        <van-field
          v-model="formData.issue_remark"
          type="textarea"
          placeholder="请输入发放备注（可选）"
          rows="3"
        />
      </van-cell-group>

      <div class="footer-actions">
        <van-button round block type="primary" native-type="submit">
          确认发放
        </van-button>
      </div>
    </van-form>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showToast, showConfirmDialog } from 'vant'
import { getRequisitionDetail, issueRequisition } from '@/api/requisition'

const router = useRouter()
const route = useRoute()

const loading = ref(true)
const requisition = ref(null)
const formData = ref({
  items: [],
  issue_remark: ''
})

async function loadData() {
  loading.value = true
  try {
    const response = await getRequisitionDetail(route.params.id)
    requisition.value = response.data

    formData.value.items = (requisition.value.items || []).map(item => ({
      id: item.id,
      material_id: item.material_id,
      material_name: item.material_name,
      requested_quantity: item.requested_quantity || 0,
      issued_quantity: item.requested_quantity || 0
    }))
  } catch (error) {
    showToast({ type: 'fail', message: '加载失败' })
  } finally {
    loading.value = false
  }
}

async function handleSubmit() {
  try {
    await showConfirmDialog({
      title: '确认发放',
      message: '确定要发放这些物料吗？'
    })

    const items = formData.value.items.map(item => ({
      id: item.id,
      issued_quantity: Number(item.issued_quantity) || 0
    }))

    const totalRequested = formData.value.items.reduce((sum, item) => sum + (item.requested_quantity || 0), 0)
    const totalIssued = items.reduce((sum, item) => sum + item.issued_quantity, 0)

    if (totalIssued === 0) {
      showToast({ type: 'fail', message: '至少发放一种物料' })
      return
    }

    if (totalIssued > totalRequested) {
      showToast({ type: 'fail', message: '发放数量不能超过申请数量' })
      return
    }

    await issueRequisition(requisition.value.id, {
      items,
      issue_remark: formData.value.issue_remark
    })

    showToast({ type: 'success', message: '发放成功' })
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
.req-issue {
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

.item-row {
  padding: 12px 16px;
  background: #f7f8fa;
  margin: 8px 16px;
  border-radius: 8px;
}

.item-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.item-name {
  font-weight: bold;
  font-size: 14px;
}

.footer-actions {
  padding: 16px;
}
</style>
