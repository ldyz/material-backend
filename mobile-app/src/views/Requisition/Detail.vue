<template>
  <div class="req-detail">
    <van-nav-bar title="出库详情" left-arrow @click-left="router.back()" />

    <van-loading v-if="loading" type="spinner" vertical />

    <div v-else-if="requisition">
      <DetailInfoGroup
        title="基本信息"
        :item="requisition"
        status-type="requisition"
        :fields="basicFields"
      />

      <van-cell-group inset title="物料明细">
        <van-empty v-if="!requisition.items?.length" description="暂无物料" />
        <van-cell
          v-for="item in requisition.items"
          :key="item.id"
          :title="item.material_name"
          :label="`申请数量：${item.requested_quantity || 0}`"
        />
      </van-cell-group>

      <div class="footer-actions">
        <van-button
          v-if="requisition.status === 'pending'"
          block
          type="primary"
          @click="router.push(`/requisition/${requisition.id}/approve`)"
        >
          去审批
        </van-button>
        <van-button
          v-if="requisition.status === 'rejected'"
          block
          type="primary"
          :loading="resubmitting"
          @click="handleResubmit"
        >
          重新提交
        </van-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showToast, showConfirmDialog } from 'vant'
import { getRequisitionDetail, resubmitRequisition } from '@/api/requisition'
import DetailInfoGroup from '@/components/common/DetailInfoGroup.vue'
import { logger } from '@/utils/logger'

const router = useRouter()
const route = useRoute()

const loading = ref(true)
const resubmitting = ref(false)
const requisition = ref(null)

const basicFields = computed(() => [
  { key: 'requisition_number', label: '出库单号' },
  { key: 'project_name', label: '项目名称' },
  { key: 'department_name', label: '申请部门' },
  { key: 'requisition_date', label: '申请日期', type: 'date' },
  { key: 'status', label: '状态', type: 'status' }
])

async function loadData() {
  loading.value = true
  try {
    const response = await getRequisitionDetail(route.params.id)
    requisition.value = response.data
  } catch (error) {
    logger.error('加载失败:', error)
  } finally {
    loading.value = false
  }
}

async function handleResubmit() {
  try {
    await showConfirmDialog({
      title: '确认重新提交',
      message: '确认重新提交该出库单进行审批？',
      teleport: '#app',
      confirmButtonColor: '#1989fa'
    })

    resubmitting.value = true
    try {
      await resubmitRequisition(requisition.value.id, {})
      showToast({ type: 'success', message: '已重新提交' })
      setTimeout(() => {
        loadData()
        router.back()
      }, 1500)
    } catch (error) {
      const errorMsg = error.error || error.message || '操作失败'
      showToast({ type: 'fail', message: errorMsg })
      throw error
    } finally {
      resubmitting.value = false
    }
  } catch (error) {
    if (error !== 'cancel') {
      logger.error('重新提交失败:', error)
    }
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.req-detail {
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
  flex-direction: column;
  gap: 12px;
}
</style>
