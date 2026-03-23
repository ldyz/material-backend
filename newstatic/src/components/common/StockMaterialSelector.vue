<template>
  <div class="stock-material-selector">
    <el-button :icon="Plus" @click="handleOpen" type="primary" size="small">
      选择库存物资
    </el-button>

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
      <el-table-column prop="quantity" label="数量" width="130">
        <template #default="scope">
          <el-input-number
            v-model="scope.row.quantity"
            :min="1"
            :max="scope.row.stock_quantity || 999999"
            size="small"
            :disabled="!editable"
            @change="handleQuantityChange(scope.row)"
          />
          <div class="qty-hint">库存: {{ scope.row.stock_quantity }}</div>
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
      </div>
    </div>

    <!-- 选择物资对话框 -->
    <el-dialog
      v-model="dialogVisible"
      title="选择库存物资"
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

      <!-- 调试信息 -->
      <div v-if="debugMode" class="debug-info">
        <el-alert type="info" :closable="false">
          <div>总数: {{ stockList.length }}</div>
          <div>加载中: {{ loading }}</div>
          <div>已选: {{ selectedRows.length }}</div>
        </el-alert>
      </div>

      <!-- 物资列表 -->
      <el-table
        ref="materialTableRef"
        v-loading="loading"
        :data="stockList"
        border
        stripe
        height="450"
        @selection-change="handleSelectionChange"
        style="margin-top: 15px"
      >
        <el-table-column type="selection" width="50" />
        <el-table-column prop="material_name" label="物资名称" min-width="150" show-overflow-tooltip />
        <el-table-column prop="spec" label="规格型号" width="120" show-overflow-tooltip />
        <el-table-column prop="material" label="材质" width="100" show-overflow-tooltip />
        <el-table-column prop="unit" label="单位" width="70" />
        <el-table-column prop="stock_quantity" label="库存数量" width="100" align="right">
          <template #default="scope">
            <el-tag :type="getStockTagType(scope.row.stock_quantity)" size="small">
              {{ scope.row.stock_quantity }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="project_name" label="项目" min-width="120" show-overflow-tooltip />
        <!-- 调试用列：显示原始数据 -->
        <el-table-column label="调试-原始数据" width="150" v-if="debugMode">
          <template #default="scope">
            <pre style="font-size: 10px; max-height: 100px; overflow: auto;">{{ JSON.stringify(scope.row, null, 2) }}</pre>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
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
import { ref, computed, reactive, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Plus, Delete } from '@element-plus/icons-vue'
import { stockApi, projectApi } from '@/api'

const props = defineProps({
  modelValue: {
    type: Array,
    default: () => []
  },
  editable: {
    type: Boolean,
    default: true
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
const stockList = ref([])
const selectedRows = ref([])
const materialTableRef = ref(null)
const projectList = ref([])
const selectedProjectId = ref(null)
const debugMode = ref(true) // 调试模式
const projectsLoaded = ref(false) // 项目是否已加载
const projectsLoading = ref(false) // 项目加载中

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 获取项目列表
const fetchProjects = async () => {
  // 防止重复加载
  if (projectsLoading.value) {
    console.log('[StockMaterialSelector] 项目数据正在加载中，跳过')
    return
  }

  // 如果数据已加载，直接返回
  if (projectsLoaded.value && projectList.value.length > 0) {
    console.log('[StockMaterialSelector] 项目数据已缓存，跳过加载')
    return
  }

  try {
    projectsLoading.value = true
    console.log('[StockMaterialSelector] 开始加载项目列表...')
    const response = await projectApi.getList({ pageSize: 1000 })
    console.log('[StockMaterialSelector] 项目API响应:', response)

    // 兼容不同的响应格式
    let projects = []
    if (response && response.success) {
      if (Array.isArray(response.data)) {
        projects = response.data
      } else if (response.data?.projects && Array.isArray(response.data.projects)) {
        projects = response.data.projects
      } else if (response.projects && Array.isArray(response.projects)) {
        projects = response.projects
      }
    }

    console.log('[StockMaterialSelector] 解析后的项目列表:', projects)
    projectList.value = projects
    projectsLoaded.value = true

    if (props.projectId) {
      selectedProjectId.value = props.projectId
    }
  } catch (error) {
    console.error('[StockMaterialSelector] 获取项目列表失败:', error)
  } finally {
    projectsLoading.value = false
  }
}

// 获取库存物资列表
const fetchStockMaterials = async () => {
  loading.value = true

  try {
    console.log('[StockMaterialSelector] ===== 开始加载库存物资 =====')
    const params = {
      page: pagination.page,
      pageSize: pagination.pageSize,
      search: searchKeyword.value || undefined,
      project_id: selectedProjectId.value || undefined
    }
    console.log('[StockMaterialSelector] 请求参数:', params)

    const response = await stockApi.getList(params)
    console.log('[StockMaterialSelector] API响应:', response)

    // 获取数据数组 - 直接使用 response.data
    let dataArray = []
    if (response && Array.isArray(response.data)) {
      dataArray = response.data
      console.log('[StockMaterialSelector] ✓ 直接使用 response.data，长度:', dataArray.length)
    } else if (response?.data?.stocks && Array.isArray(response.data.stocks)) {
      dataArray = response.data.stocks
      console.log('[StockMaterialSelector] ✓ 使用 response.data.stocks，长度:', dataArray.length)
    } else if (response?.data?.items && Array.isArray(response.data.items)) {
      dataArray = response.data.items
      console.log('[StockMaterialSelector] ✓ 使用 response.data.items，长度:', dataArray.length)
    } else {
      console.error('[StockMaterialSelector] ✗ 无法从响应中获取数据数组')
      console.log('[StockMaterialSelector] response.data:', response?.data)
      console.log('[StockMaterialSelector] response 键:', response ? Object.keys(response) : 'response 为空')
    }

    console.log('[StockMaterialSelector] 数据数组长度:', dataArray.length)

    if (dataArray.length === 0) {
      console.warn('[StockMaterialSelector] ⚠ 数据数组为空')
      stockList.value = []
      pagination.total = 0
      return
    }

    // 字段映射：确保所有必需字段都存在
    const mappedList = dataArray.map((item, index) => {
      const mapped = {
        id: item.id,
        material_id: item.material_id,
        material_name: item.material_name || item.name || '',
        spec: item.specification || item.spec || '',
        material: item.material || item.mat_material || '',
        unit: item.unit || '',
        stock_quantity: item.quantity || 0,
        project_name: item.project_name || ''
      }

      if (index === 0) {
        console.log('[StockMaterialSelector] 第一条数据映射:')
        console.log('  原始:', JSON.stringify(item, null, 2))
        console.log('  映射后:', JSON.stringify(mapped, null, 2))
      }

      return mapped
    })

    // 强制更新 stockList
    stockList.value = [...mappedList]
    console.log('[StockMaterialSelector] ✓ 已赋值给 stockList.value，长度:', stockList.value.length)

    // 更新分页总数
    pagination.total = response?.pagination?.total || dataArray.length || 0

    console.log('[StockMaterialSelector] ===== 加载完成 =====')
  } catch (error) {
    console.error('[StockMaterialSelector] ✗ 获取库存列表失败:', error)
    ElMessage.error('获取库存列表失败: ' + error.message)
    stockList.value = []
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchStockMaterials()
}

// 分页
const handlePageChange = (page) => {
  pagination.page = page
  fetchStockMaterials()
}

const handleSizeChange = (size) => {
  pagination.pageSize = size
  pagination.page = 1
  fetchStockMaterials()
}

// 打开对话框
const handleOpen = () => {
  console.log('[StockMaterialSelector] 打开选择对话框')
  dialogVisible.value = true
  fetchStockMaterials()
}

// 选择变化
const handleSelectionChange = (selection) => {
  console.log('[StockMaterialSelector] 选择变化:', selection.length, '项')
  selectedRows.value = selection
}

// 确认选择
const handleConfirm = () => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning('请至少选择一个物资')
    return
  }

  console.log('[StockMaterialSelector] 确认选择:', selectedRows.value.length, '项')

  // 合并到已选列表
  selectedRows.value.forEach(row => {
    const exists = selectedMaterials.value.find(m => m.material_id === row.material_id)
    if (exists) {
      ElMessage.info(`物资"${row.material_name}"已在列表中`)
    } else {
      selectedMaterials.value.push({
        material_id: row.material_id,
        stock_id: row.id,
        material: row.material || '',
        material_name: row.material_name,
        spec: row.specification || row.spec || '',
        unit: row.unit,
        quantity: row.stock_quantity || 0, // 默认使用最大库存数量
        stock_quantity: row.stock_quantity,
        remark: ''
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
  if (row.quantity > row.stock_quantity) {
    ElMessage.warning(`物资"${row.material_name}"数量不能超过库存数量 ${row.stock_quantity}`)
    row.quantity = row.stock_quantity
  }
  emitChange()
}

// 触发变化事件
const emitChange = () => {
  emit('change', selectedMaterials.value)
}

// 获取库存标签类型
const getStockTagType = (stock) => {
  if (stock <= 0) return 'danger'
  if (stock < 10) return 'warning'
  return 'success'
}

// 监控 stockList 的变化
watch(() => stockList.value, (newVal, oldVal) => {
  console.log('[StockMaterialSelector] ===== stockList 变化 =====')
  console.log('[StockMaterialSelector] 新值长度:', newVal?.length)
  console.log('[StockMaterialSelector] 旧值长度:', oldVal?.length)
  if (newVal && newVal.length > 0) {
    console.log('[StockMaterialSelector] 新值第一条数据:', newVal[0])
  }
  console.log('[StockMaterialSelector] ===== stockList 变化完成 =====')
}, { deep: true })

onMounted(() => {
  console.log('[StockMaterialSelector] 组件已挂载')
  fetchProjects()
})
</script>

<style scoped>
.stock-material-selector {
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

.debug-info {
  margin-bottom: 10px;
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
