<template>
  <div class="inbound-detail">
    <van-nav-bar title="入库详情" left-arrow @click-left="router.back()" />

    <van-loading v-if="loading" type="spinner" vertical />

    <div v-else-if="order">
      <van-cell-group inset title="基本信息">
        <van-cell title="入库单号" :value="order.order_number" />
        <van-cell title="项目名称" :value="order.project_name || '-'" />
        <van-cell title="供应商" :value="order.supplier_name || '-'" />
        <van-cell title="入库日期" :value="formatDate(order.inbound_date)" />
        <van-cell title="状态">
          <template #value>
            <van-tag :type="getStatusType(order.workflow_status)">
              {{ getStatusText(order.workflow_status) }}
            </van-tag>
          </template>
        </van-cell>
      </van-cell-group>

      <van-cell-group inset title="物料明细">
        <van-empty v-if="!order.items?.length" description="暂无物料" />
        <van-cell
          v-for="item in order.items"
          :key="item.id"
          :title="item.material_name"
          :label="`数量：${item.quantity || 0}`"
        />
      </van-cell-group>

      <div v-if="order.workflow_status === 'pending_approval'" class="footer-actions">
        <van-button block type="primary" @click="router.push(`/inbound/${order.id}/approve`)">
          去审批
        </van-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getInboundDetail } from '@/api/inbound'

const router = useRouter()
const route = useRoute()

const loading = ref(true)
const order = ref(null)

function getStatusType(status) {
  const map = {
    draft: 'default',
    pending_approval: 'warning',
    approved: 'success',
    rejected: 'danger',
    completed: 'primary'
  }
  return map[status] || 'default'
}

function getStatusText(status) {
  const map = {
    draft: '草稿',
    pending_approval: '待审批',
    approved: '已批准',
    rejected: '已拒绝',
    completed: '已完成'
  }
  return map[status] || status
}

function formatDate(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('zh-CN')
}

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
