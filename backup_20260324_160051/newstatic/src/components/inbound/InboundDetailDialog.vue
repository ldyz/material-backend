<template>
  <Dialog
    v-model="visible"
    title="入库单详情"
    width="900px"
    :loading="loading"
    @cancel="handleClose"
  >
    <template #extra>
      <el-button
        type="primary"
        :icon="Printer"
        @click="handlePrint"
        size="small"
      >
        打印
      </el-button>
    </template>

    <el-descriptions :column="2" border>
      <el-descriptions-item label="入库单号" :span="2">
        <el-tag type="primary">{{ data.order_no || '-' }}</el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="状态">
        <el-tag :type="getStatusTagType(data.status)" size="small">
          {{ getStatusText(data.status) }}
        </el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="供应商">
        {{ data.supplier || '-' }}
      </el-descriptions-item>
      <el-descriptions-item label="联系人">
        {{ data.contact || '-' }}
      </el-descriptions-item>
      <el-descriptions-item label="验收人">
        {{ data.receiver || '-' }}
      </el-descriptions-item>
      <el-descriptions-item label="创建人">
        {{ data.creator_name || '-' }}
      </el-descriptions-item>
      <el-descriptions-item label="创建时间">
        {{ data.created_at || '-' }}
      </el-descriptions-item>
      <el-descriptions-item label="关联项目" :span="2">
        {{ data.project_name || '-' }}
      </el-descriptions-item>
      <el-descriptions-item label="备注" :span="2" v-if="data.remark">
        {{ data.remark || '-' }}
      </el-descriptions-item>
    </el-descriptions>

    <!-- 工作流状态显示 -->
    <div style="margin-top: 20px;">
      <WorkflowStatus
        v-if="data.status"
        :status="data.status"
        :status-time="data.updated_at || data.created_at"
        :status-description="getStatusDescription(data.status)"
        workflow-type="inbound"
      />
    </div>

    <!-- 物资明细 -->
    <el-divider content-position="left">物资明细</el-divider>
    <el-table
      :data="data.items"
      border
      stripe
      style="width: 100%"
      size="small"
    >
      <el-table-column prop="material_name" label="物资名称" min-width="150" show-overflow-tooltip />
      <el-table-column prop="spec" label="规格型号" width="120" show-overflow-tooltip />
      <el-table-column prop="material" label="材质" width="100" show-overflow-tooltip />
      <el-table-column prop="unit" label="单位" width="80" />
      <el-table-column prop="quantity" label="数量" width="100" align="right" />
      <el-table-column prop="unit_price" label="单价" width="100" align="right">
        <template #default="scope">
          {{ scope.row.unit_price ? formatCurrency(scope.row.unit_price) : '-' }}
        </template>
      </el-table-column>
      <el-table-column label="金额" width="120" align="right">
        <template #default="scope">
          {{ formatCurrency((scope.row.quantity || 0) * (scope.row.unit_price || 0)) }}
        </template>
      </el-table-column>
      <el-table-column prop="remark" label="备注" min-width="120" show-overflow-tooltip />
    </el-table>

    <!-- 工作流历史记录 -->
    <template v-if="workflowHistories.length > 0">
      <el-divider content-position="left">审批历史</el-divider>
      <WorkflowHistory :histories="workflowHistories" />
    </template>
  </Dialog>
</template>

<script setup>
import { ref, watch } from 'vue'
import { Printer } from '@element-plus/icons-vue'
import Dialog from '@/components/common/Dialog.vue'
import WorkflowStatus from '@/components/common/WorkflowStatus.vue'
import WorkflowHistory from '@/components/common/WorkflowHistory.vue'
import { inboundApi } from '@/api'
import { ElMessage } from 'element-plus'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  orderNo: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue'])

const visible = ref(false)
const loading = ref(false)
const data = ref({
  order_no: '',
  supplier: '',
  contact: '',
  receiver: '',
  inbound_date: '',
  remark: '',
  items: [],
  status: '',
  project_id: null,
  project_name: '',
  creator_name: '',
  created_at: '',
  updated_at: ''
})
const workflowHistories = ref([])

// 格式化货币
const formatCurrency = (value) => {
  return Number(value).toLocaleString('zh-CN', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  })
}

// 获取状态标签类型
const getStatusTagType = (status) => {
  const types = {
    pending: 'info',
    approved: 'warning',
    rejected: 'danger',
    completed: 'success'
  }
  return types[status] || 'info'
}

// 获取状态文本
const getStatusText = (status) => {
  const texts = {
    pending: '待审核',
    approved: '已批准',
    rejected: '已拒绝',
    completed: '已完成'
  }
  return texts[status] || status
}

// 获取状态描述
const getStatusDescription = (status) => {
  const descriptions = {
    pending: '入库单等待审核',
    approved: '入库单已批准，等待入库',
    rejected: '入库单已被拒绝',
    completed: '入库单已完成入库'
  }
  return descriptions[status] || ''
}

// 获取工作流历史
const fetchWorkflowHistory = async (id) => {
  try {
    const response = await inboundApi.getWorkflowHistory(id)
    workflowHistories.value = response.data || []
  } catch (error) {
    console.error('获取工作流历史失败:', error)
  }
}

// 加载数据
const loadData = async () => {
  if (!props.orderNo) return

  loading.value = true
  try {
    // 根据单号查询入库单
    const response = await inboundApi.getList({ pageSize: 1000 })
    const order = response.data.find(item => item.order_no === props.orderNo || item.inbound_no === props.orderNo)

    if (!order) {
      ElMessage.error('未找到该入库单')
      handleClose()
      return
    }

    // 获取完整的入库单详情
    const detail = await inboundApi.getDetail(order.id)
    data.value = detail.data || detail

    // 加载工作流历史
    await fetchWorkflowHistory(order.id)
  } catch (error) {
    console.error('加载入库单详情失败:', error)
    ElMessage.error('加载入库单详情失败')
    handleClose()
  } finally {
    loading.value = false
  }
}

// 关闭对话框
const handleClose = () => {
  visible.value = false
  emit('update:modelValue', false)
}

// 打印
const handlePrint = () => {
  // 创建打印内容
  const printContent = `
    <!DOCTYPE html>
    <html>
    <head>
      <meta charset="utf-8">
      <title>入库单 - ${data.value.order_no}</title>
      <style>
        * {
          margin: 0;
          padding: 0;
          box-sizing: border-box;
        }

        body {
          font-family: "Microsoft YaHei", Arial, sans-serif;
          font-size: 12px;
          line-height: 1.5;
          color: #333;
          padding: 20px;
        }

        .header {
          text-align: center;
          margin-bottom: 20px;
          border-bottom: 2px solid #333;
          padding-bottom: 10px;
        }

        .header h1 {
          font-size: 20px;
          margin-bottom: 5px;
        }

        .info-section {
          margin-bottom: 20px;
        }

        .info-row {
          display: flex;
          margin-bottom: 8px;
          border-bottom: 1px solid #eee;
          padding-bottom: 5px;
        }

        .info-label {
          font-weight: bold;
          width: 100px;
          flex-shrink: 0;
        }

        .info-value {
          flex: 1;
        }

        .section-title {
          font-size: 14px;
          font-weight: bold;
          margin: 20px 0 10px 0;
          padding-bottom: 5px;
          border-bottom: 1px solid #333;
        }

        table {
          width: 100%;
          border-collapse: collapse;
          margin-bottom: 20px;
        }

        table th,
        table td {
          border: 1px solid #333;
          padding: 8px;
          text-align: left;
        }

        table th {
          background-color: #f5f5f5;
          font-weight: bold;
          text-align: center;
        }

        table td {
          text-align: center;
        }

        .text-left {
          text-align: left !important;
        }

        .text-right {
          text-align: right !important;
        }

        @page {
          size: A4;
          margin: 15mm;
        }

        /* 防止表格行跨页 */
        tr {
          page-break-inside: avoid;
        }

        /* 防止元素内部跨页 */
        .info-section,
        .section-title {
          page-break-inside: avoid;
        }

        /* 避免在表格后立即分页 */
        table {
          page-break-inside: avoid;
        }
      </style>
    </head>
    <body>
      <div class="header">
        <h1>入库单</h1>
      </div>

      <div class="info-section">
        <div class="info-row">
          <span class="info-label">入库单号：</span>
          <span class="info-value">${data.value.order_no || '-'}</span>
        </div>
        <div class="info-row">
          <span class="info-label">状态：</span>
          <span class="info-value">${getStatusText(data.value.status)}</span>
        </div>
        <div class="info-row">
          <span class="info-label">供应商：</span>
          <span class="info-value">${data.value.supplier || '-'}</span>
        </div>
        <div class="info-row">
          <span class="info-label">联系人：</span>
          <span class="info-value">${data.value.contact || '-'}</span>
        </div>
        <div class="info-row">
          <span class="info-label">验收人：</span>
          <span class="info-value">${data.value.receiver || '-'}</span>
        </div>
        <div class="info-row">
          <span class="info-label">创建人：</span>
          <span class="info-value">${data.value.creator_name || '-'}</span>
        </div>
        <div class="info-row">
          <span class="info-label">创建时间：</span>
          <span class="info-value">${data.value.created_at || '-'}</span>
        </div>
        <div class="info-row">
          <span class="info-label">关联项目：</span>
          <span class="info-value">${data.value.project_name || '-'}</span>
        </div>
        ${data.value.remark ? `
        <div class="info-row">
          <span class="info-label">备注：</span>
          <span class="info-value">${data.value.remark}</span>
        </div>
        ` : ''}
      </div>

      <div class="section-title">物资明细</div>
      <table>
        <thead>
          <tr>
            <th>物资名称</th>
            <th>规格型号</th>
            <th>材质</th>
            <th>单位</th>
            <th>数量</th>
            <th>单价</th>
            <th>金额</th>
          </tr>
        </thead>
        <tbody>
          ${data.value.items.map(item => `
            <tr>
              <td class="text-left">${item.material_name || '-'}</td>
              <td>${item.spec || '-'}</td>
              <td>${item.material || '-'}</td>
              <td>${item.unit || '-'}</td>
              <td>${item.quantity || 0}</td>
              <td>${item.unit_price ? Number(item.unit_price).toFixed(2) : '-'}</td>
              <td class="text-right">${((item.quantity || 0) * (item.unit_price || 0)).toFixed(2)}</td>
            </tr>
          `).join('')}
        </tbody>
      </table>

      <div style="margin-top: 30px; text-align: right; font-size: 10px; color: #999;">
        打印时间：${new Date().toLocaleString('zh-CN')}
      </div>
    </body>
    </html>
  `

  // 打开新窗口打印
  const printWindow = window.open('', '_blank')
  if (printWindow) {
    printWindow.document.write(printContent)
    printWindow.document.close()
    printWindow.onload = () => {
      printWindow.print()
      printWindow.close()
    }
  } else {
    ElMessage.error('无法打开打印窗口，请检查浏览器设置')
  }
}

// 监听 modelValue 变化
watch(() => props.modelValue, (val) => {
  visible.value = val
  if (val) {
    loadData()
  }
})

// 监听 visible 变化
watch(visible, (val) => {
  if (!val) {
    emit('update:modelValue', false)
  }
})
</script>

<style scoped>
:deep(.el-descriptions__label) {
  width: 120px;
}

/* 打印样式 */
@media print {
  /* 防止表格行跨页分割 */
  :deep(.el-table__body-wrapper) {
    page-break-inside: avoid;
  }

  /* 防止描述列表跨页分割 */
  :deep(.el-descriptions) {
    page-break-inside: avoid;
  }

  /* 防止表格行跨页 */
  :deep(.el-table__row) {
    page-break-inside: avoid;
  }

  /* 防止表格单元格跨页 */
  :deep(.el-table__cell) {
    page-break-inside: avoid;
  }

  /* 确保对话框内容正确显示 */
  :deep(.el-dialog__body) {
    page-break-inside: avoid;
  }

  /* 防止分隔符后立即分页 */
  :deep(.el-divider) {
    page-break-after: avoid;
    page-break-inside: avoid;
  }

  /* 移除打印时的不需要的元素 */
  :deep(.el-button),
  :deep(.el-dialog__headerbtn) {
    display: none !important;
  }

  /* 确保内容宽度 */
  :deep(.el-dialog) {
    width: 100% !important;
    max-width: 100% !important;
  }

  /* 移除背景和阴影以节省墨水 */
  :deep(.el-card),
  :deep(.el-dialog) {
    box-shadow: none !important;
    border: none !important;
  }

  /* 分页设置 */
  @page {
    size: A4;
    margin: 15mm 15mm 15mm 15mm;
  }

  /* 打印时保持可见性 */
  body {
    visibility: visible;
  }

  /* 只显示对话框内容 */
  #app > *:not(.el-overlay) {
    display: none;
  }

  .el-overlay {
    position: static !important;
    background: white !important;
  }

  .el-dialog {
    position: static !important;
    transform: none !important;
    margin: 0 !important;
    box-shadow: none !important;
  }
}
</style>
