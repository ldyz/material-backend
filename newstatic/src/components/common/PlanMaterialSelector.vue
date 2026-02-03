<template>
  <div class="plan-material-selector">
    <el-button :icon="Plus" @click="handleOpen" type="primary" size="small" :disabled="!planId">
      选择物资
    </el-button>
    <span v-if="!planId" class="tip-text">请先选择物资计划</span>

    <el-table
      :data="selectedMaterials"
      border
      stripe
      size="small"
      style="width: 100%; margin-top: 10px"
    >
      <el-table-column prop="material_name" label="物资名称" min-width="150" show-overflow-tooltip />
      <el-table-column prop="spec" label="规格型号" width="120" show-overflow-tooltip />
      <el-table-column prop="material" label="材质" width="100" show-overflow-tooltip />
      <el-table-column prop="unit" label="单位" width="70" />
      <el-table-column prop="quantity" label="入库数量" width="130">
        <template #default="scope">
          <el-input-number
            v-model="scope.row.quantity"
            :min="1"
            :max="scope.row.remaining_quantity || 999999"
            size="small"
            :disabled="!editable"
            @change="handleQuantityChange(scope.row)"
            :placeholder="`最大: ${scope.row.remaining_quantity}`"
          />
          <div class="qty-hint">可入: {{ scope.row.remaining_quantity }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="unit_price" label="单价" width="100" align="right">
        <template #default="scope">
          <template v-if="editable">
            <el-input-number
              v-model="scope.row.unit_price"
              :min="0"
              :precision="2"
              size="small"
              @change="handleQuantityChange(scope.row)"
              style="width: 100%"
            />
          </template>
          <template v-else>
            {{ scope.row.unit_price ? formatCurrency(scope.row.unit_price) : '-' }}
          </template>
        </template>
      </el-table-column>
      <el-table-column label="金额" width="100" align="right">
        <template #default="scope">
          {{ formatCurrency(scope.row.quantity * scope.row.unit_price) }}
        </template>
      </el-table-column>
      <el-table-column prop="remark" label="备注" width="120" v-if="editable">
        <template #default="scope">
          <el-input
            v-model="scope.row.remark"
            placeholder="备注"
            size="small"
            @change="handleQuantityChange(scope.row)"
          />
        </template>
      </el-table-column>
      <el-table-column label="操作" width="60" fixed="right" v-if="editable">
        <template #default="scope">
          <el-button
            type="danger"
            size="small"
            :icon="Delete"
            @click="handleRemove(scope.$index)"
            link
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="total-amount" v-if="selectedMaterials.length > 0">
      <div class="summary">
        <span>共 <strong>{{ selectedMaterials.length }}</strong> 项物资</span>
        <span style="margin-left: 20px;">总金额: <strong style="color: #f56c6c; font-size: 16px;">{{ formatCurrency(totalAmount) }}</strong></span>
      </div>
    </div>

    <!-- 选择物资对话框 -->
    <el-dialog
      v-model="dialogVisible"
      title="选择计划物资"
      width="1000px"
      :close-on-click-modal="false"
    >
      <!-- 搜索栏 -->
      <div class="search-bar">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索物资名称、编码、材质"
          clearable
          style="width: 300px"
          @input="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </div>

      <!-- 物资列表 -->
      <el-table
        ref="materialTableRef"
        v-loading="loading"
        :data="planItemList"
        border
        stripe
        height="450"
        @selection-change="handleSelectionChange"
        style="margin-top: 15px"
      >
        <el-table-column type="selection" width="50" :selectable="checkSelectable" />
        <el-table-column prop="material_name" label="物资名称" min-width="150" show-overflow-tooltip />
        <el-table-column prop="spec" label="规格型号" width="120" show-overflow-tooltip />
        <el-table-column prop="material" label="材质" width="100" show-overflow-tooltip />
        <el-table-column prop="unit" label="单位" width="70" />
        <el-table-column label="到货情况" width="150" align="center">
          <template #default="scope">
            <div class="arrival-info">
              <div class="arrival-stats">
                <span class="stat-item">计划: {{ scope.row.planned_quantity }}</span>
                <span class="stat-item arrived">已到: {{ scope.row.received_quantity || 0 }}</span>
                <span class="stat-item warning">可入: {{ scope.row.remaining_quantity }}</span>
              </div>
              <el-tag v-if="scope.row.status === 'completed'" type="success" size="small">已到齐</el-tag>
              <el-tag v-else-if="scope.row.status === 'partial'" type="warning" size="small">
                {{ ((scope.row.received_quantity / scope.row.planned_quantity) * 100).toFixed(1) }}%
              </el-tag>
              <el-tag v-else type="info" size="small">未到货</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="unit_price" label="单价" width="90" align="right">
          <template #default="scope">
            {{ scope.row.unit_price ? formatCurrency(scope.row.unit_price) : '-' }}
          </template>
        </el-table-column>
      </el-table>

      <!-- 底部操作栏 -->
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleConfirm" :disabled="selectedRows.length === 0">
            确定 (已选 {{ selectedRows.length }} 项)
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Plus, Delete } from '@element-plus/icons-vue'
import { materialPlanApi, materialApi } from '@/api'

const props = defineProps({
  modelValue: {
    type: Array,
    default: () => []
  },
  editable: {
    type: Boolean,
    default: true
  },
  planId: {
    type: Number,
    default: null
  },
  projectId: {
    type: Number,
    default: null
  }
})

const emit = defineEmits(['update:modelValue', 'change'])

const selectedMaterials = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const dialogVisible = ref(false)
const loading = ref(false)
const searchKeyword = ref('')
const planItemList = ref([])
const selectedRows = ref([])
const materialTableRef = ref(null)

// 计算总金额
const totalAmount = computed(() => {
  return selectedMaterials.value.reduce((sum, item) => {
    return sum + (item.quantity || 0) * (item.unit_price || 0)
  }, 0)
})

// 获取计划物资列表
const fetchPlanItems = async () => {
  if (!props.planId) {
    planItemList.value = []
    return
  }

  loading.value = true
  try {
    const { data } = await materialPlanApi.getPlanItems(props.planId)
    const allItems = data || []

    // 过滤出未到齐的物资
    const itemsWithRemaining = allItems.filter(item => {
      const remaining = (item.planned_quantity || 0) - (item.received_quantity || 0)
      return remaining > 0
    })

    // 检查是否有物资缺少material_id
    const itemsWithoutMaterial = itemsWithRemaining.filter(item => !item.material_id)

    if (itemsWithoutMaterial.length > 0) {
      console.warn('以下计划物资未关联到物资库，将自动创建：', itemsWithoutMaterial.map(i => i.material_name))

      // 自动为没有material_id的物资创建物资记录
      try {
        // 调用批量创建物资的API
        const materialsData = itemsWithoutMaterial.map(item => ({
          name: item.material_name,
          code: item.material_code || '',
          specification: item.specification || '',
          category: item.category || '',
          unit: item.unit,
          price: item.unit_price || 0,
          quantity: 0,
          project_id: props.projectId || null
        }))

        const result = await materialApi.batchCreateMaterials(materialsData)

        if (result.success && result.data?.materials) {
          // 更新item的material_id
          const materialMap = new Map()
          result.data.materials.forEach(m => {
            const key = [m.name, m.specification || ''].join('|')
            materialMap.set(key, m.id)
          })

          let updatedCount = 0
          itemsWithoutMaterial.forEach(item => {
            const key = [item.material_name, item.specification || ''].join('|')
            if (materialMap.has(key)) {
              item.material_id = materialMap.get(key)
              updatedCount++
            }
          })

          // 调用后端API更新material_plan_items表的material_id
          if (updatedCount > 0) {
            await updatePlanItemsMaterialID(props.planId, itemsWithoutMaterial)
          }

          ElMessage.success(`已自动创建并关联 ${updatedCount} 个物资记录`)
        }
      } catch (error) {
        console.error('自动创建物资失败:', error)
        ElMessage.warning(`自动创建物资失败: ${error.message}，请手动关联物资`)
      }
    }

    // 过滤出有效物资（有material_id的物资）
    const validItems = itemsWithRemaining.filter(item => item.material_id)

    planItemList.value = validItems.map(item => ({
      ...item,
      remaining_quantity: (item.planned_quantity || 0) - (item.received_quantity || 0),
      spec: item.specification || item.spec || ''
    }))

    if (planItemList.value.length === 0 && itemsWithRemaining.length > 0) {
      ElMessage.warning('该计划中的物资无法入库，请检查物资关联状态')
    }
  } catch (error) {
    console.error('获取计划物资失败:', error)
    ElMessage.error('获取计划物资失败')
  } finally {
    loading.value = false
  }
}

// 更新计划项的material_id
const updatePlanItemsMaterialID = async (planId, items) => {
  try {
    // 使用API层调用后端API更新material_plan_items
    await materialPlanApi.syncPlanMaterialIds(planId, {
      items: items.map(item => ({
        id: item.id,
        material_id: item.material_id
      }))
    })
  } catch (error) {
    console.error('更新计划项material_id失败:', error)
  }
}

// 搜索
const handleSearch = () => {
  if (!searchKeyword.value) {
    fetchPlanItems()
    return
  }
  const keyword = searchKeyword.value.toLowerCase()
  planItemList.value = planItemList.value.filter(item => {
    return (item.material_name || '').toLowerCase().includes(keyword) ||
           (item.material || '').toLowerCase().includes(keyword) ||
           (item.spec || '').toLowerCase().includes(keyword)
  })
}

// 打开对话框
const handleOpen = () => {
  if (!props.planId) {
    ElMessage.warning('请先选择物资计划')
    return
  }
  dialogVisible.value = true
  fetchPlanItems()
}

// 选择变化
const handleSelectionChange = (selection) => {
  selectedRows.value = selection
}

// 检查行是否可选择
const checkSelectable = (row) => {
  return row.remaining_quantity > 0
}

// 确认选择
const handleConfirm = () => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning('请至少选择一个物资')
    return
  }

  // 合并到已选列表
  selectedRows.value.forEach(row => {
    // 使用plan_item_id来判断唯一性，而不是material_id
    // 这样不同计划中的相同物资可以分别入库
    const exists = selectedMaterials.value.find(m => m.plan_item_id === row.id)
    if (exists) {
      ElMessage.info(`物资"${row.material_name}"已在列表中`)
    } else {
      selectedMaterials.value.push({
        material_id: row.material_id,
        material: row.material || '',
        material_name: row.material_name,
        spec: row.spec || row.specification || '',
        unit: row.unit,
        quantity: row.remaining_quantity, // 默认填入可入数量
        unit_price: row.unit_price || 0,
        remark: '',
        plan_item_id: row.id,
        planned_quantity: row.planned_quantity,
        received_quantity: row.received_quantity || 0,
        remaining_quantity: row.remaining_quantity
      })
    }
  })

  dialogVisible.value = false
  if (materialTableRef.value) {
    materialTableRef.value.clearSelection()
  }
  selectedRows.value = []
  emitChange()
}

// 移除物资
const handleRemove = (index) => {
  selectedMaterials.value.splice(index, 1)
  emitChange()
}

// 数量变化
const handleQuantityChange = (row) => {
  if (row.remaining_quantity !== undefined && row.quantity > row.remaining_quantity) {
    ElMessage.warning(`物资"${row.material_name}"数量不能超过可入数量 ${row.remaining_quantity}`)
    row.quantity = row.remaining_quantity
  }
  emitChange()
}

// 触发变化事件
const emitChange = () => {
  emit('change', selectedMaterials.value)
}

// 格式化货币
const formatCurrency = (value) => {
  return Number(value).toLocaleString('zh-CN', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  })
}

// 监听计划ID变化
watch(() => props.planId, (newVal, oldVal) => {
  if (newVal !== oldVal) {
    // 清空已选物资
    selectedMaterials.value = []
    emitChange()
  }
})
</script>

<style scoped>
.plan-material-selector {
  width: 100%;
}

.tip-text {
  margin-left: 10px;
  font-size: 12px;
  color: #909399;
}

.total-amount {
  margin-top: 15px;
  padding: 10px;
  background: #f5f7fa;
  border-radius: 4px;
}

.summary {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 14px;
}

.search-bar {
  display: flex;
  align-items: center;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.arrival-info {
  padding: 4px 0;
}

.arrival-stats {
  display: flex;
  flex-direction: column;
  gap: 2px;
  margin-bottom: 4px;
}

.stat-item {
  font-size: 11px;
  color: #606266;
}

.stat-item.arrived {
  color: #67c23a;
}

.stat-item.warning {
  color: #e6a23c;
  font-weight: bold;
}

.qty-hint {
  font-size: 10px;
  color: #909399;
  margin-top: 2px;
}

:deep(.el-input-number) {
  width: 100%;
}

:deep(.el-input-number .el-input__inner) {
  text-align: left;
}
</style>
