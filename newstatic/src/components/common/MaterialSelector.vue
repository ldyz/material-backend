<template>
  <div class="material-selector">
    <el-button :icon="Plus" @click="handleOpen" type="primary" size="small">
      选择物资
    </el-button>

    <el-table
      :data="selectedMaterials"
      border
      stripe
      size="small"
      style="width: 100%; margin-top: 10px"
    >
      <!-- 序号列已移除 -->
      <el-table-column prop="material_name" label="物资名称" min-width="150" show-overflow-tooltip />
      <el-table-column prop="spec" label="规格型号" width="120" show-overflow-tooltip />
      <el-table-column prop="material" label="材质" width="100" show-overflow-tooltip />
      <el-table-column prop="unit" label="单位" width="70" />
      <el-table-column prop="quantity" label="数量" width="130">
        <template #default="scope">
          <el-input-number
            v-model="scope.row.quantity"
            :min="1"
            :max="scope.row.remaining_quantity || 999999"
            size="small"
            :disabled="!editable"
            @change="handleQuantityChange(scope.row)"
            :placeholder="scope.row.remaining_quantity ? `最大: ${scope.row.remaining_quantity}` : '数量'"
          />
          <div v-if="scope.row.remaining_quantity !== undefined && scope.row.planned_quantity > 0" class="qty-hint">
            剩余: {{ scope.row.remaining_quantity }}
          </div>
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
              :disabled="!editable"
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
      title="选择物资"
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
        <el-select
          v-model="selectedProjectId"
          placeholder="筛选项目"
          clearable
          filterable
          style="width: 200px; margin-left: 10px"
          @change="handleSearch"
        >
          <el-option label="全部项目" :value="null" />
          <el-option
            v-for="project in projectList"
            :key="project.id"
            :label="project.name"
            :value="project.id"
          />
        </el-select>
      </div>

      <!-- 物资列表 -->
      <el-table
        ref="materialTableRef"
        v-loading="loading"
        :data="materialList"
        border
        stripe
        height="450"
        @selection-change="handleSelectionChange"
        style="margin-top: 15px"
      >
        <el-table-column type="selection" width="50" :selectable="checkSelectable" />
        <!-- 序号列已移除 -->
        <el-table-column prop="name" label="物资名称" min-width="150" show-overflow-tooltip />
        <el-table-column prop="specification" label="规格型号" width="120" show-overflow-tooltip />
        <el-table-column prop="material" label="材质" width="100" show-overflow-tooltip />
        <el-table-column prop="unit" label="单位" width="70" />
        <el-table-column label="到货情况" width="150" align="center">
          <template #default="scope">
            <div v-if="scope.row.quantity > 0" class="arrival-info">
              <div class="arrival-stats">
                <span class="stat-item">计划: {{ scope.row.quantity }}</span>
                <span class="stat-item arrived">已到: {{ scope.row.arrived_quantity || 0 }}</span>
                <span class="stat-item" :class="{ 'warning': (scope.row.remaining_quantity || 0) < 10 }">
                  剩余: {{ scope.row.remaining_quantity || (scope.row.quantity - (scope.row.arrived_quantity || 0)) }}
                </span>
              </div>
              <el-tag v-if="scope.row.is_fully_arrived" type="success" size="small">已到齐</el-tag>
              <el-tag v-else-if="scope.row.arrival_percentage >= 80" type="warning" size="small">
                {{ (scope.row.arrival_percentage || 0).toFixed(1) }}%
              </el-tag>
              <el-tag v-else type="info" size="small">
                {{ (scope.row.arrival_percentage || 0).toFixed(1) }}%
              </el-tag>
            </div>
            <el-tag v-else type="info" size="small">未设计划</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="stock" label="库存" width="90" align="right">
          <template #default="scope">
            <el-tag :type="getStockTagType(scope.row.stock)" size="small">
              {{ scope.row.stock || 0 }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="price" label="单价" width="90" align="right">
          <template #default="scope">
            {{ scope.row.price ? formatCurrency(scope.row.price) : '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="project_name" label="项目" min-width="120" show-overflow-tooltip />
      </el-table>

      <!-- 分页 -->
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="fetchMaterials"
        @current-change="fetchMaterials"
        style="margin-top: 15px"
      />

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
import { ref, computed, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Plus, Delete } from '@element-plus/icons-vue'
import { materialApi, projectApi } from '@/api'

const props = defineProps({
  modelValue: {
    type: Array,
    default: () => []
  },
  editable: {
    type: Boolean,
    default: true
  },
  // 用于区分是入库单还是出库单，决定是否显示库存
  showStock: {
    type: Boolean,
    default: true
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
const materialList = ref([])
const selectedRows = ref([])
const materialTableRef = ref(null)
const projectList = ref([])
const selectedProjectId = ref(null)

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 计算总金额
const totalAmount = computed(() => {
  return selectedMaterials.value.reduce((sum, item) => {
    return sum + (item.quantity || 0) * (item.unit_price || 0)
  }, 0)
})

// 获取项目列表
const fetchProjects = async () => {
  try {
    const { data } = await projectApi.getList({ pageSize: 1000 })
    projectList.value = data || []
  } catch (error) {
    console.error('获取项目列表失败:', error)
  }
}

// 获取物资列表
const fetchMaterials = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      pageSize: pagination.pageSize,
      search: searchKeyword.value || undefined,
      project_id: selectedProjectId.value || undefined,
      include_children: selectedProjectId.value ? true : undefined // 包含子项目
    }
    const { data, pagination: pag } = await materialApi.getList(params)

    // 过滤掉已完全到货的物资
    const allMaterials = data || []
    const filteredMaterials = allMaterials.filter(item => {
      // 如果已完全到货，过滤掉
      if (item.is_fully_arrived) {
        return false
      }
      // 如果剩余数量为0，过滤掉
      const remaining = item.remaining_quantity || (item.planned_quantity - (item.arrived_quantity || 0))
      if (item.planned_quantity > 0 && remaining <= 0) {
        return false
      }
      return true
    })

    materialList.value = filteredMaterials
    // 注意：这里total需要调整，因为过滤后数量减少了
    // 为了简单起见，显示实际返回的数量，而不是数据库总数
    pagination.total = pag?.total || 0
  } catch (error) {
    console.error('获取物资列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchMaterials()
}

// 打开对话框
const handleOpen = () => {
  dialogVisible.value = true
  fetchMaterials()
}

// 选择变化
const handleSelectionChange = (selection) => {
  selectedRows.value = selection
}

// 确认选择
const handleConfirm = () => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning('请至少选择一个物资')
    return
  }

  // 合并到已选列表
  selectedRows.value.forEach(row => {
    const exists = selectedMaterials.value.find(m => m.material_id === row.id)
    if (exists) {
      // 如果已存在，询问是否增加数量
      ElMessage.info(`物资"${row.name}"已在列表中`)
    } else {
      // 计算剩余可入数量
      const remainingQty = row.remaining_quantity || (row.planned_quantity - (row.arrived_quantity || 0))

      selectedMaterials.value.push({
        material_id: row.id,
        material: row.material || '',
        material_name: row.name,
        spec: row.specification || row.spec || '',
        unit: row.unit,
        quantity: 1,
        stock_quantity: row.stock,
        unit_price: row.price || 0,
        remark: '',
        // 新增：到货信息用于验证
        planned_quantity: row.planned_quantity || 0,
        arrived_quantity: row.arrived_quantity || 0,
        remaining_quantity: remainingQty,
        is_fully_arrived: row.is_fully_arrived || false
      })
    }
  })

  dialogVisible.value = false
  // 清空选择
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
  // 如果有计划数量，验证不超过剩余数量
  if (row.planned_quantity > 0 && row.remaining_quantity !== undefined) {
    if (row.quantity > row.remaining_quantity) {
      ElMessage.warning(`物资"${row.material_name}"数量不能超过剩余可入数量 ${row.remaining_quantity}`)
      // 自动调整为最大值
      row.quantity = row.remaining_quantity
    }
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

// 获取库存标签类型
const getStockTagType = (stock) => {
  if (stock <= 0) return 'danger'
  if (stock < 10) return 'warning'
  return 'success'
}

// 检查行是否可选择（已到齐的物资不可选）
const checkSelectable = (row) => {
  // 如果已完全到货，禁止选择
  if (row.is_fully_arrived) {
    return false
  }
  // 如果剩余数量为0，禁止选择
  const remaining = row.remaining_quantity || (row.planned_quantity - (row.arrived_quantity || 0))
  if (row.planned_quantity > 0 && remaining <= 0) {
    return false
  }
  return true
}

onMounted(() => {
  fetchProjects()
})
</script>

<style scoped>
.material-selector {
  width: 100%;
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
