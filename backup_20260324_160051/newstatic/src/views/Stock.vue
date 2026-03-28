<template>
  <div class="stock-container">
    <el-card shadow="never">
      <!-- 工具栏 -->
      <TableToolbar>
        <template #left>
          <ProjectSelector
            v-model="searchForm.project_id"
            :projects="projectList"
            placeholder="选择项目（支持层级显示）"
            width="300px"
          />
          <el-input
            v-model="searchForm.keyword"
            placeholder="搜索物资名称、编码"
            clearable
            style="width: 250px"
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-select
            v-model="searchForm.category"
            placeholder="物资分类"
            clearable
            style="width: 150px"
          >
            <el-option label="全部" value="" />
            <el-option
              v-for="cat in categoryList"
              :key="cat.id"
              :label="cat.name"
              :value="cat.name"
            />
          </el-select>
          <el-select
            v-model="searchForm.status"
            placeholder="库存状态"
            clearable
            style="width: 150px"
          >
            <el-option label="所有状态" value="" />
            <el-option label="正常" value="normal" />
            <el-option label="库存偏低" value="low" />
            <el-option label="库存不足" value="shortage" />
          </el-select>
          <el-button :icon="Refresh" @click="handleReset">重置</el-button>
        </template>
        <template #right>
          <el-button
            type="success"
            :icon="Download"
            @click="handleExport"
            v-if="authStore.hasPermission('stock_export')"
          >
            导出
          </el-button>
        </template>
      </TableToolbar>

      <!-- 表格 -->
      <el-table
        v-loading="loading"
        :data="tableData"
        border
        stripe
        style="width: 100%"
      >
        <!-- 序号列已移除 -->
        <el-table-column prop="material_code" label="物资编码" width="130" />
        <el-table-column prop="material_name" label="物资名称" min-width="150" show-overflow-tooltip />
        <el-table-column prop="category" label="分类" width="100">
          <template #default="scope">
            <el-tag size="small">{{ scope.row.category || '-' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="specification" label="规格型号" width="120" show-overflow-tooltip />
        <el-table-column label="材质" width="100" show-overflow-tooltip>
          <template #default="scope">
            {{ scope.row.material || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="unit" label="单位" width="80" />
        <el-table-column prop="quantity" label="库存数量" width="120" align="right">
          <template #default="scope">
            <el-tag
              :type="getStockTagType(scope.row.quantity)"
              size="large"
            >
              {{ scope.row.quantity || 0 }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="safety_stock" label="安全库存" width="100" align="right">
          <template #default="scope">
            {{ scope.row.safety_stock || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="stock_status" label="库存状态" width="100">
          <template #default="scope">
            <el-tag :type="getStockStatusTagType(scope.row)" size="small">
              {{ getStockStatusText(scope.row) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="project_name" label="关联项目" min-width="150" show-overflow-tooltip>
          <template #default="scope">
            {{ scope.row.project_name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="updated_at" label="最后更新" width="160">
          <template #default="scope">
            {{ scope.row.updated_at || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="260" fixed="right">
          <template #default="scope">
            <el-button
              type="success"
              size="small"
              :icon="Plus"
              @click="handleRowIn(scope.row)"
              v-if="authStore.hasPermission('stock_in')"
            >
              入库
            </el-button>
            <el-button
              type="warning"
              size="small"
              :icon="Minus"
              @click="handleRowOut(scope.row)"
              v-if="authStore.hasPermission('stock_out')"
            >
              出库
            </el-button>
            <el-button
              type="primary"
              size="small"
              :icon="Document"
              @click="handleViewLogs(scope.row)"
            >
              日志
            </el-button>
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
        class="mt-20"
      />
    </el-card>

    <!-- 库存操作对话框 -->
    <StockOperationDialog
      v-model="operationDialogVisible"
      :material-id="operationMaterialId"
      :operation-type="operationType"
      :stock-data="operationStockData"
      @success="fetchData"
    />

    <!-- 库存日志对话框 -->
    <Dialog
      v-model="logsDialogVisible"
      title="库存日志"
      width="900px"
      :show-footer="false"
    >
      <el-table
        v-loading="logsLoading"
        :data="logsData"
        border
        stripe
        max-height="400"
      >
        <el-table-column prop="created_at" label="时间" width="160" />
        <el-table-column prop="type" label="类型" width="80">
          <template #default="scope">
            <el-tag :type="getLogTypeTag(scope.row.type)" size="small">
              {{ getLogTypeText(scope.row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="quantity" label="数量" width="80" align="right" />
        <el-table-column prop="quantity_before" label="操作前" width="80" align="right" />
        <el-table-column prop="quantity_after" label="操作后" width="80" align="right" />
        <el-table-column prop="price" label="单价" width="80" align="right">
          <template #default="scope">
            {{ scope.row.price ? formatCurrency(scope.row.price) : '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="200">
          <template #default="scope">
            <el-link
              v-if="scope.row.inbound_code || scope.row.requisition_code || isRelatedOrder(scope.row.remark)"
              type="primary"
              @click="handleViewDetail(scope.row)"
            >
              {{ formatRemark(scope.row) }}
            </el-link>
            <span v-else>{{ scope.row.remark || '-' }}</span>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="logsPagination.page"
        v-model:page-size="logsPagination.pageSize"
        :page-sizes="[10, 20, 50]"
        :total="logsPagination.total"
        layout="total, sizes, prev, pager, next"
        @size-change="fetchLogs"
        @current-change="fetchLogs"
        class="mt-20"
      />
    </Dialog>

    <!-- 入库单详情对话框 -->
    <InboundDetailDialog
      v-model="inboundDetailVisible"
      :order-no="currentOrderNo"
    />

    <!-- 出库单详情对话框 -->
    <RequisitionDetailDialog
      v-model="requisitionDetailVisible"
      :requisition-no="currentRequisitionNo"
    />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { stockApi } from '@/api'
import { ElMessage } from 'element-plus'
import {
  Search,
  Refresh,
  Plus,
  Minus,
  Download,
  Document
} from '@element-plus/icons-vue'
import Dialog from '@/components/common/Dialog.vue'
import TableToolbar from '@/components/common/TableToolbar.vue'
import InboundDetailDialog from '@/components/inbound/InboundDetailDialog.vue'
import RequisitionDetailDialog from '@/components/requisition/RequisitionDetailDialog.vue'
import ProjectSelector from '@/components/common/ProjectSelector.vue'
import StockOperationDialog from '@/components/stock/StockOperationDialog.vue'

const authStore = useAuthStore()

// 列表数据
const loading = ref(false)
const tableData = ref([])
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 搜索表单
const searchForm = reactive({
  keyword: '',
  category: '',
  project_id: '',
  status: ''
})

// 项目列表
const projectList = ref([])

// 物资分类列表
const categoryList = ref([])

// 库存操作对话框
const operationDialogVisible = ref(false)
const operationMaterialId = ref(null)
const operationType = ref('in')
const operationStockData = ref(null)

// 库存日志
const logsDialogVisible = ref(false)
const logsLoading = ref(false)
const logsData = ref([])
const logsPagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})
const currentLogStockId = ref(null)

// 详情对话框
const inboundDetailVisible = ref(false)
const requisitionDetailVisible = ref(false)
const currentOrderNo = ref('')
const currentRequisitionNo = ref('')

// 获取列表数据
// 适配统一响应格式
const fetchData = async () => {
  loading.value = true
  try {
    // 收集项目ID（包含子项目）
    let projectIds = []
    if (searchForm.project_id) {
      projectIds = collectProjectIds(searchForm.project_id, projectList.value)
    }

    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      search: searchForm.keyword || undefined,
      category: searchForm.category || undefined,
      project_ids: projectIds.length > 0 ? projectIds.join(',') : undefined,
      status: searchForm.status || undefined
    }
    const { data, pagination: pag } = await stockApi.getList(params)
    tableData.value = data || []
    pagination.total = pag?.total || 0
  } catch (error) {
    console.error('获取库存列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 收集项目及其所有子项目的ID
const collectProjectIds = (projectId, projectTree) => {
  const ids = [projectId]

  // 在树形结构中查找项目并收集子项目ID
  const findAndCollectChildren = (nodes) => {
    for (const node of nodes) {
      if (node.id === projectId) {
        // 找到目标项目，递归收集所有子项目ID
        const collectAllChildren = (project) => {
          if (project.children && project.children.length > 0) {
            for (const child of project.children) {
              ids.push(child.id)
              collectAllChildren(child)
            }
          }
        }
        collectAllChildren(node)
        return true
      }
      if (node.children && node.children.length > 0) {
        if (findAndCollectChildren(node.children)) {
          return true
        }
      }
    }
    return false
  }

  findAndCollectChildren(projectTree)
  return ids
}

// 加载项目列表
// 适配统一响应格式
const fetchProjects = async () => {
  try {
    const { projectApi } = await import('@/api')
    const { data } = await projectApi.getList({ pageSize: 1000 })
    // 构建树形结构
    projectList.value = buildProjectTree(data || [])
  } catch (error) {
    console.error('获取项目列表失败:', error)
  }
}

// 构建项目树形结构
const buildProjectTree = (projects) => {
  if (!projects || projects.length === 0) return []

  // 创建项目映射
  const projectMap = new Map()
  projects.forEach(project => {
    projectMap.set(project.id, { ...project, children: [] })
  })

  const roots = []

  // 构建树形结构
  projects.forEach(project => {
    const node = projectMap.get(project.id)
    if (!project.parent_id) {
      // 没有父项目，作为根节点
      roots.push(node)
    } else {
      // 有父项目，添加到父项目的children中
      const parent = projectMap.get(project.parent_id)
      if (parent) {
        parent.children.push(node)
      } else {
        // 父项目不在列表中，作为根节点
        roots.push(node)
      }
    }
  })

  return roots
}

// 加载物资分类列表
// 适配统一响应格式
const fetchCategories = async () => {
  try {
    const { materialApi } = await import('@/api')
    const { data } = await materialApi.getCategories()
    categoryList.value = data || []
  } catch (error) {
    console.error('获取物资分类列表失败:', error)
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchData()
}

// 重置
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.category = ''
  searchForm.project_id = ''
  searchForm.status = ''
  pagination.page = 1
  fetchData()
}

// 行入库（带物资信息）
const handleRowIn = (row) => {
  operationMaterialId.value = row.material_id
  operationType.value = 'in'
  operationStockData.value = row
  operationDialogVisible.value = true
}

// 行出库（带物资信息）
const handleRowOut = (row) => {
  operationMaterialId.value = row.material_id
  operationType.value = 'out'
  operationStockData.value = row
  operationDialogVisible.value = true
}

// 查看日志
const handleViewLogs = (row) => {
  currentLogStockId.value = row.id
  logsDialogVisible.value = true
  fetchLogs()
}

// 获取库存日志
// 适配统一响应格式
const fetchLogs = async () => {
  if (!currentLogStockId.value) return

  logsLoading.value = true
  try {
    const params = {
      page: logsPagination.page,
      page_size: logsPagination.pageSize,
      stock_id: currentLogStockId.value
    }
    const { data, pagination: pag } = await stockApi.getLogs(params)
    logsData.value = data || []
    logsPagination.total = pag?.total || 0
  } catch (error) {
    console.error('获取库存日志失败:', error)
  } finally {
    logsLoading.value = false
  }
}

// 导出
const handleExport = async () => {
  try {
    const response = await stockApi.export(searchForm)
    const blob = new Blob([response], {
      type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
    })
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `库存列表_${new Date().getTime()}.xlsx`
    a.click()
    window.URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
  } catch (error) {
    console.error('导出失败:', error)
  }
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

// 获取库存状态标签类型
const getStockStatusTagType = (row) => {
  const quantity = row.quantity || 0
  const safeStock = row.safe_stock || 0

  if (quantity <= 0) return 'danger'
  if (quantity < safeStock) return 'warning'
  return 'success'
}

// 获取库存状态文本
const getStockStatusText = (row) => {
  const quantity = row.quantity || 0
  const safeStock = row.safe_stock || 0

  if (quantity <= 0) return '库存不足'
  if (quantity < safeStock) return '库存偏低'
  return '正常'
}

// 获取日志类型标签
const getLogTypeTag = (type) => {
  const tags = {
    in: 'success',
    out: 'warning',
    adjust: 'info'
  }
  return tags[type] || 'info'
}

// 获取日志类型文本
const getLogTypeText = (type) => {
  const texts = {
    in: '入库',
    out: '出库',
    adjust: '调整'
  }
  return texts[type] || type
}

// 检查是否关联到订单
const isRelatedOrder = (remark) => {
  if (!remark) return false
  return remark.includes('入库单') || remark.includes('出库单')
}

// 格式化备注显示
const formatRemark = (row) => {
  // 优先使用直接字段
  if (row.inbound_code) {
    return `入库单-${row.inbound_code}`
  }
  if (row.requisition_code) {
    return `出库单-${row.requisition_code}`
  }

  // 从备注中提取单号（兼容旧数据）
  if (!row.remark) return '-'

  // 入库单（多种格式）：入库单、入库单审核入库、备注：入库单等
  const inboundMatch = row.remark.match(/入库单(?:审核入库)?\s*[:：]?\s*(\w+)/)
  if (inboundMatch) {
    return `入库单-${inboundMatch[1]}`
  }

  // 出库单（多种格式）：出库单、出库单发放等
  const requisitionMatch = row.remark.match(/出库单(?:发放)?\s*[:：]?\s*(\w+)/)
  if (requisitionMatch) {
    return `出库单-${requisitionMatch[1]}`
  }

  // 其他情况返回原备注
  return row.remark
}

// 查看详情
const handleViewDetail = (row) => {
  // 优先使用直接字段
  if (row.inbound_code) {
    currentOrderNo.value = row.inbound_code
    inboundDetailVisible.value = true
    return
  }

  if (row.requisition_code) {
    currentRequisitionNo.value = row.requisition_code
    requisitionDetailVisible.value = true
    return
  }

  // 从备注中提取单号（兼容旧数据）
  if (!row.remark) return

  // 入库单（多种格式）：入库单、入库单审核入库、备注：入库单等
  const inboundMatch = row.remark.match(/入库单(?:审核入库)?\s*[:：]?\s*(\w+)/)
  if (inboundMatch) {
    currentOrderNo.value = inboundMatch[1]
    inboundDetailVisible.value = true
    return
  }

  // 出库单（多种格式）：出库单、出库单发放等
  const requisitionMatch = row.remark.match(/出库单(?:发放)?\s*[:：]?\s*(\w+)/)
  if (requisitionMatch) {
    currentRequisitionNo.value = requisitionMatch[1]
    requisitionDetailVisible.value = true
    return
  }
}

// 防抖定时器
let searchTimer = null

// 即时搜索函数（带防抖）
const debouncedSearch = () => {
  // 清除之前的定时器
  if (searchTimer) {
    clearTimeout(searchTimer)
  }

  // 设置新的定时器，500ms后执行搜索
  searchTimer = setTimeout(() => {
    pagination.page = 1
    fetchData()
  }, 500)
}

// 分页处理
const handlePageChange = (page) => {
  pagination.page = page
  fetchData()
}

const handleSizeChange = (size) => {
  pagination.pageSize = size
  pagination.page = 1
  fetchData()
}

// 监听所有搜索字段变化，实现即时搜索
watch(() => searchForm.keyword, debouncedSearch)
watch(() => searchForm.category, debouncedSearch)
watch(() => searchForm.project_id, debouncedSearch)
watch(() => searchForm.status, debouncedSearch)

onMounted(() => {
  fetchProjects()
  fetchCategories()
  fetchData()
})
</script>

<style scoped>
.stock-container {
  padding: 0;
}

.mt-20 {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
