<template>
  <div class="plan-selector">
    <el-select
      :model-value="modelValue"
      :placeholder="placeholder"
      :disabled="disabled"
      :clearable="clearable"
      :filterable="filterable"
      :loading="loading"
      :remote="remote"
      :remote-method="remoteMethod"
      @update:model-value="$emit('update:modelValue', $event)"
      @change="handleChange"
      style="width: 100%"
    >
      <el-option
        v-for="plan in filteredPlans"
        :key="plan.id"
        :label="`${plan.plan_no} - ${plan.plan_name}`"
        :value="plan.id"
      >
        <div class="plan-option">
          <div class="plan-option-header">
            <span class="plan-no">{{ plan.plan_no }}</span>
            <el-tag :type="getStatusTagType(plan.status)" size="small">
              {{ getStatusLabel(plan.status) }}
            </el-tag>
          </div>
          <div class="plan-option-info">
            <span class="plan-name">{{ plan.plan_name }}</span>
            <span class="project-name">{{ plan.project_name }}</span>
          </div>
          <div class="plan-option-meta">
            <span>{{ plan.items_count || 0 }} 项</span>
            <span v-if="plan.total_budget">¥{{ formatAmount(plan.total_budget) }}</span>
          </div>
        </div>
      </el-option>
      <template #footer>
        <el-button
          v-if="showCreateButton"
          text
          :icon="Plus"
          @click="handleCreate"
          style="width: 100%"
        >
          新建计划
        </el-button>
      </template>
    </el-select>

    <!-- 计划项目列表 -->
    <div v-if="showItems && selectedPlanItems.length > 0" class="plan-items mt-10">
      <div class="items-header">
        <h4>计划项目 ({{ selectedPlanItems.length }})</h4>
        <el-button
          type="primary"
          size="small"
          :icon="Check"
          @click="handleSelectAll"
        >
          全部添加
        </el-button>
      </div>
      <el-table
        :data="selectedPlanItems"
        border
        size="small"
        max-height="300"
        @selection-change="handleSelectionChange"
      >
        <el-table-column
          v-if="selectable"
          type="selection"
          width="55"
        />
        <el-table-column prop="material_name" label="物资名称" min-width="150" show-overflow-tooltip />
        <el-table-column prop="specification" label="规格型号" width="120" show-overflow-tooltip />
        <el-table-column prop="unit" label="单位" width="70" />
        <el-table-column prop="planned_quantity" label="计划数量" width="90" align="right" />
        <el-table-column prop="remaining_quantity" label="剩余数量" width="90" align="right">
          <template #default="scope">
            {{ scope.row.remaining_quantity || scope.row.planned_quantity }}
          </template>
        </el-table-column>
        <el-table-column prop="unit_price" label="单价" width="90" align="right">
          <template #default="scope">
            {{ scope.row.unit_price ? '¥' + scope.row.unit_price : '-' }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="80" fixed="right">
          <template #default="scope">
            <el-button
              type="primary"
              size="small"
              :icon="Plus"
              @click="handleAddItem(scope.row)"
            >
              添加
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { materialPlanApi } from '@/api'
import { ElMessage } from 'element-plus'
import { Plus, Check } from '@element-plus/icons-vue'

const props = defineProps({
  modelValue: {
    type: [Number, String],
    default: null
  },
  // 项目ID筛选
  projectId: {
    type: Number,
    default: null
  },
  // 状态筛选
  status: {
    type: String,
    default: 'active' // 默认只显示进行中的计划
  },
  placeholder: {
    type: String,
    default: '请选择物资计划'
  },
  disabled: {
    type: Boolean,
    default: false
  },
  clearable: {
    type: Boolean,
    default: true
  },
  filterable: {
    type: Boolean,
    default: true
  },
  // 是否显示计划项目
  showItems: {
    type: Boolean,
    default: false
  },
  // 是否支持多选
  selectable: {
    type: Boolean,
    default: false
  },
  // 是否显示新建按钮
  showCreateButton: {
    type: Boolean,
    default: false
  },
  // 是否使用远程搜索
  remote: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'change', 'add-item', 'add-items'])

const loading = ref(false)
const plans = ref([])
const selectedPlanItems = ref([])
const selectedItems = ref([])

// 过滤后的计划列表
const filteredPlans = computed(() => {
  let result = plans.value

  // 状态筛选
  if (props.status) {
    result = result.filter(p => p.status === props.status)
  }

  // 项目筛选
  if (props.projectId) {
    result = result.filter(p => p.project_id === props.projectId)
  }

  return result
})

// 获取计划列表
const fetchPlans = async (search = '') => {
  loading.value = true
  try {
    const params = {
      page: 1,
      page_size: 100
    }

    if (search) {
      params.search = search
    }

    const response = await materialPlanApi.getPlans(params)
    if (response.success) {
      plans.value = response.data || []
    }
  } catch (error) {
    console.error('获取计划列表失败:', error)
    ElMessage.error('获取计划列表失败')
  } finally {
    loading.value = false
  }
}

// 远程搜索方法
const remoteMethod = (query) => {
  if (query) {
    fetchPlans(query)
  } else {
    fetchPlans()
  }
}

// 获取计划项目
const fetchPlanItems = async (planId) => {
  if (!planId || !props.showItems) return

  try {
    const response = await materialPlanApi.getPlanItems(planId)
    if (response.success) {
      selectedPlanItems.value = (response.data || []).filter(item => {
        // 只显示未完成的项目
        return item.status !== 'completed' && item.status !== 'cancelled'
      })
    }
  } catch (error) {
    console.error('获取计划项目失败:', error)
  }
}

// 处理选择变化
const handleChange = (value) => {
  const plan = plans.value.find(p => p.id === value)
  emit('change', plan)
  fetchPlanItems(value)
}

// 添加单个项目
const handleAddItem = (item) => {
  emit('add-item', item)
}

// 添加所有项目
const handleSelectAll = () => {
  emit('add-items', selectedPlanItems.value)
}

// 选择变化（多选）
const handleSelectionChange = (selection) => {
  selectedItems.value = selection
}

// 新建计划
const handleCreate = () => {
  // 打开新建计划对话框（由父组件处理）
  emit('create')
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

const formatAmount = (amount) => {
  return Number(amount || 0).toLocaleString('zh-CN', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  })
}

// 监听项目ID变化
watch(() => props.projectId, () => {
  fetchPlans()
})

// 初始化
onMounted(() => {
  fetchPlans()
  if (props.modelValue) {
    fetchPlanItems(props.modelValue)
  }
})

// 暴露方法
defineExpose({
  fetchPlans,
  fetchPlanItems
})
</script>

<style scoped>
.plan-selector {
  width: 100%;
}

.plan-option {
  padding: 8px 0;
}

.plan-option-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.plan-no {
  font-weight: 600;
  color: #303133;
  font-size: 14px;
}

.plan-option-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.plan-name {
  color: #606266;
  font-size: 13px;
}

.project-name {
  color: #909399;
  font-size: 12px;
}

.plan-option-meta {
  display: flex;
  gap: 15px;
  font-size: 12px;
  color: #909399;
}

.plan-items {
  margin-top: 15px;
  padding: 15px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.items-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.items-header h4 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #303133;
}

.mt-10 {
  margin-top: 10px;
}
</style>
