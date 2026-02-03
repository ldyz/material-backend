<template>
  <el-dialog
    :model-value="modelValue"
    :title="planDetail?.plan_name || '计划详情'"
    width="1000px"
    @update:model-value="$emit('update:modelValue', $event)"
    @open="handleOpen"
  >
    <el-skeleton v-if="loading" :rows="10" animated />
    <div v-else-if="planDetail">
      <!-- 计划基本信息 -->
      <el-descriptions :column="3" border class="mb-20">
        <el-descriptions-item label="计划编号">{{ planDetail.plan_no }}</el-descriptions-item>
        <el-descriptions-item label="所属项目">{{ planDetail.project_name }}</el-descriptions-item>
        <el-descriptions-item label="计划类型">
          <el-tag :type="getPlanTypeTagType(planDetail.plan_type)" size="small">
            {{ getPlanTypeLabel(planDetail.plan_type) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusTagType(planDetail.status)" size="small">
            {{ getStatusLabel(planDetail.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="优先级">
          <el-tag :type="getPriorityTagType(planDetail.priority)" size="small">
            {{ getPriorityLabel(planDetail.priority) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="创建人">{{ planDetail.creator_name }}</el-descriptions-item>
        <el-descriptions-item label="计划开始日期">
          {{ planDetail.planned_start_date ? formatDate(planDetail.planned_start_date) : '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="计划结束日期">
          {{ planDetail.planned_end_date ? formatDate(planDetail.planned_end_date) : '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">
          {{ formatDateTime(planDetail.created_at) }}
        </el-descriptions-item>
        <el-descriptions-item label="总预算" v-if="planDetail.total_budget">
          ¥{{ formatAmount(planDetail.total_budget) }}
        </el-descriptions-item>
        <el-descriptions-item label="实际成本" v-if="planDetail.actual_cost">
          ¥{{ formatAmount(planDetail.actual_cost) }}
        </el-descriptions-item>
        <el-descriptions-item label="项目数量">
          {{ planDetail.items_count || 0 }} 项
        </el-descriptions-item>
        <el-descriptions-item label="描述" :span="3" v-if="planDetail.description">
          {{ planDetail.description }}
        </el-descriptions-item>
        <el-descriptions-item label="备注" :span="3" v-if="planDetail.remark">
          <pre class="remark-text">{{ planDetail.remark }}</pre>
        </el-descriptions-item>
      </el-descriptions>

      <!-- 进度统计 -->
      <div v-if="planDetail.progress" class="progress-section mb-20">
        <h4 class="section-title">执行进度</h4>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-card shadow="hover">
              <div class="progress-card">
                <div class="progress-label">到货进度</div>
                <el-progress
                  :percentage="Math.round(planDetail.progress.receive_progress || 0)"
                  :stroke-width="12"
                  :color="getProgressColor(planDetail.progress.receive_progress)"
                />
              </div>
            </el-card>
          </el-col>
          <el-col :span="8">
            <el-card shadow="hover">
              <div class="progress-card">
                <div class="progress-label">发放进度</div>
                <el-progress
                  :percentage="Math.round(planDetail.progress.issue_progress || 0)"
                  :stroke-width="12"
                  :color="getProgressColor(planDetail.progress.issue_progress)"
                />
              </div>
            </el-card>
          </el-col>
          <el-col :span="8">
            <el-card shadow="hover">
              <div class="progress-card">
                <div class="progress-label">总体进度</div>
                <el-progress
                  :percentage="Math.round(planDetail.progress.overall_progress || 0)"
                  :stroke-width="12"
                  :color="getProgressColor(planDetail.progress.overall_progress)"
                  type="circle"
                  :width="80"
                />
              </div>
            </el-card>
          </el-col>
        </el-row>
      </div>

      <!-- 计划项目列表 -->
      <div class="items-section">
        <h4 class="section-title">计划项目 ({{ planDetail.items?.length || 0 }})</h4>
        <el-table
          :data="planDetail.items"
          border
          size="small"
          max-height="400"
        >
          <el-table-column prop="material_name" label="物资名称" min-width="150" show-overflow-tooltip />
          <el-table-column prop="material_code" label="物资编码" width="120" />
          <el-table-column prop="specification" label="规格型号" width="120" show-overflow-tooltip />
          <el-table-column prop="category" label="分类" width="100" />
          <el-table-column prop="unit" label="单位" width="70" />
          <el-table-column prop="planned_quantity" label="计划数量" width="90" align="right" />
          <el-table-column label="已到货" width="80" align="right">
            <template #default="scope">
              {{ scope.row.received_quantity || 0 }}
            </template>
          </el-table-column>
          <el-table-column label="已发放" width="80" align="right">
            <template #default="scope">
              {{ scope.row.issued_quantity || 0 }}
            </template>
          </el-table-column>
          <el-table-column label="剩余" width="80" align="right">
            <template #default="scope">
              {{ scope.row.remaining_quantity || 0 }}
            </template>
          </el-table-column>
          <el-table-column label="进度" width="150" align="center">
            <template #default="scope">
              <el-progress
                :percentage="Math.round(scope.row.receive_progress || 0)"
                :stroke-width="8"
                :status="scope.row.receive_progress >= 100 ? 'success' : undefined"
              />
            </template>
          </el-table-column>
          <el-table-column prop="unit_price" label="单价" width="90" align="right">
            <template #default="scope">
              {{ scope.row.unit_price ? '¥' + scope.row.unit_price : '-' }}
            </template>
          </el-table-column>
          <el-table-column prop="total_price" label="总价" width="100" align="right">
            <template #default="scope">
              {{ scope.row.total_price ? '¥' + formatAmount(scope.row.total_price) : '-' }}
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="90" align="center">
            <template #default="scope">
              <el-tag :type="getItemStatusTagType(scope.row.status)" size="small">
                {{ getItemStatusLabel(scope.row.status) }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 审批信息 -->
      <div v-if="planDetail.approver_name" class="approval-info mt-20">
        <h4 class="section-title">审批信息</h4>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="审批人">{{ planDetail.approver_name }}</el-descriptions-item>
          <el-descriptions-item label="审批时间">
            {{ planDetail.approved_at ? formatDateTime(planDetail.approved_at) : '-' }}
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </div>
    <div v-else class="empty-text">
      暂无数据
    </div>
  </el-dialog>
</template>

<script setup>
import { ref } from 'vue'
import { materialPlanApi } from '@/api'
import { ElMessage } from 'element-plus'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  planId: {
    type: Number,
    default: null
  }
})

const emit = defineEmits(['update:modelValue'])

const loading = ref(false)
const planDetail = ref(null)

const handleOpen = async () => {
  if (!props.planId) return

  loading.value = true
  try {
    const response = await materialPlanApi.getPlanDetail(props.planId)
    if (response.success) {
      planDetail.value = response.data
    } else {
      ElMessage.error('获取计划详情失败')
    }
  } catch (error) {
    console.error('获取计划详情失败:', error)
    ElMessage.error('获取计划详情失败')
  } finally {
    loading.value = false
  }
}

// 辅助函数
const getStatusLabel = (status) => {
  const labels = {
    draft: '草稿',
    pending: '待审批',
    approved: '已批准',
    active: '进行中',
    completed: '已完成',
    cancelled: '已取消',
    rejected: '已拒绝'
  }
  return labels[status] || status
}

const getStatusTagType = (status) => {
  const types = {
    draft: 'info',
    pending: 'warning',
    approved: 'success',
    active: 'primary',
    completed: 'success',
    cancelled: 'danger',
    rejected: 'danger'
  }
  return types[status] || 'info'
}

const getPlanTypeLabel = (type) => {
  const labels = {
    procurement: '采购计划',
    usage: '使用计划',
    mixed: '混合计划'
  }
  return labels[type] || type
}

const getPlanTypeTagType = (type) => {
  const types = {
    procurement: 'primary',
    usage: 'success',
    mixed: 'warning'
  }
  return types[type] || 'info'
}

const getPriorityLabel = (priority) => {
  const labels = {
    urgent: '紧急',
    high: '高',
    normal: '普通',
    low: '低'
  }
  return labels[priority] || priority
}

const getPriorityTagType = (priority) => {
  const types = {
    urgent: 'danger',
    high: 'warning',
    normal: 'info',
    low: ''
  }
  return types[priority] || 'info'
}

const getItemStatusLabel = (status) => {
  const labels = {
    pending: '待处理',
    partial: '部分完成',
    completed: '已完成',
    cancelled: '已取消'
  }
  return labels[status] || status
}

const getItemStatusTagType = (status) => {
  const types = {
    pending: 'info',
    partial: 'warning',
    completed: 'success',
    cancelled: 'danger'
  }
  return types[status] || 'info'
}

const getProgressColor = (percentage) => {
  if (percentage >= 100) return '#67c23a'
  if (percentage >= 50) return '#409eff'
  if (percentage >= 25) return '#e6a23c'
  return '#f56c6c'
}

const formatAmount = (amount) => {
  return Number(amount || 0).toLocaleString('zh-CN', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  })
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN')
}

const formatDateTime = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}
</script>

<style scoped>
.mb-20 {
  margin-bottom: 20px;
}

.mt-20 {
  margin-top: 20px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 12px;
  color: #303133;
}

.progress-section {
  padding: 15px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.progress-card {
  text-align: center;
  padding: 10px;
}

.progress-label {
  font-size: 14px;
  color: #606266;
  margin-bottom: 15px;
}

.remark-text {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
  font-family: inherit;
  font-size: 14px;
  color: #606266;
}

.empty-text {
  text-align: center;
  padding: 40px 0;
  color: #909399;
}
</style>
