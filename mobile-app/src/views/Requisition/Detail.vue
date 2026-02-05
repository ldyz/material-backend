<template>
  <div class="req-detail">
    <van-nav-bar title="出库详情" left-arrow @click-left="router.back()" />

    <van-loading v-if="loading" type="spinner" vertical />

    <div v-else-if="requisition">
      <van-cell-group inset title="基本信息">
        <van-cell title="出库单号" :value="requisition.requisition_number" />
        <van-cell title="项目名称" :value="requisition.project_name || '-'" />
        <van-cell title="申请部门" :value="requisition.department_name || '-'" />
        <van-cell title="申请日期" :value="formatDate(requisition.requisition_date)" />
        <van-cell title="状态">
          <template #value>
            <van-tag :type="getStatusType(requisition.workflow_status)">
              {{ getStatusText(requisition.workflow_status) }}
            </van-tag>
          </template>
        </van-cell>
      </van-cell-group>

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
          v-if="requisition.workflow_status === 'pending_approval'"
          block
          type="primary"
          @click="router.push(`/requisition/${requisition.id}/approve`)"
        >
          去审批
        </van-button>
        <van-button
          v-if="requisition.workflow_status === 'approved'"
          block
          type="success"
          @click="router.push(`/requisition/${requisition.id}/issue`)"
        >
          去发放
        </van-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getRequisitionDetail } from '@/api/requisition'

const router = useRouter()
const route = useRoute()

const loading = ref(true)
const requisition = ref(null)

function getStatusType(status) {
  const map = {
    draft: 'default',
    pending_approval: 'warning',
    approved: 'primary',
    rejected: 'danger',
    issued: 'success'
  }
  return map[status] || 'default'
}

function getStatusText(status) {
  const map = {
    draft: '草稿',
    pending_approval: '待审批',
    approved: '已批准',
    rejected: '已拒绝',
    issued: '已发放'
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
    const response = await getRequisitionDetail(route.params.id)
    requisition.value = response.data
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
