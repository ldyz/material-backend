<template>
  <div class="inbound-detail">
    <van-nav-bar title="入库详情" left-arrow @click-left="router.back()" />

    <van-loading v-if="loading" type="spinner" vertical />

    <div v-else-if="order">
      <DetailInfoGroup
        title="基本信息"
        :item="order"
        status-type="inbound"
        :fields="basicFields"
      />

      <van-cell-group inset title="物料明细">
        <van-empty v-if="!order.items?.length" description="暂无物料" />
        <van-cell
          v-for="item in order.items"
          :key="item.id"
          :title="item.material_name"
          :label="`数量：${item.quantity || 0}`"
        />
      </van-cell-group>

      <div v-if="order.status === 'pending'" class="footer-actions">
        <van-button block type="primary" @click="router.push(`/inbound/${order.id}/approve`)">
          去审批
        </van-button>
      </div>

      <div v-if="order.status === 'rejected'" class="footer-actions">
        <van-button
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
import { getInboundDetail, resubmitInbound } from '@/api/inbound'
import DetailInfoGroup from '@/components/common/DetailInfoGroup.vue'

const router = useRouter()
const route = useRoute()

const loading = ref(true)
const resubmitting = ref(false)
const order = ref(null)

const basicFields = computed(() => [
  { key: 'order_number', label: '入库单号' },
  { key: 'project_name', label: '项目名称' },
  { key: 'supplier_name', label: '供应商' },
  { key: 'inbound_date', label: '入库日期', type: 'date' },
  { key: 'status', label: '状态', type: 'status' }
])

async function loadData() {
  loading.value = true
  try {
    const response = await getInboundDetail(route.params.id)
    order.value = response.data
  } catch (error) {
    console.error('加载失败:', error)
  } finally {
    loading.value = false
  }
}

async function handleResubmit() {
  try {
    await showConfirmDialog({
      title: '确认重新提交',
      message: '确认重新提交该入库单进行审批？',
      teleport: '#app',
      confirmButtonColor: '#1989fa'
    })

    resubmitting.value = true
    try {
      await resubmitInbound(order.value.id, {})
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
      console.error('重新提交失败:', error)
    }
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.inbound-detail {
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
}
</style>
