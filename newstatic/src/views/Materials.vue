<template>
  <div class="materials-container">
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
            placeholder="搜索物资名称、编码、型号"
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
          <el-button
            :type="searchForm.unstored ? 'primary' : 'default'"
            @click="toggleUnstored"
          >
            <el-icon><Filter /></el-icon>
            未入库物资
          </el-button>
          <el-button :icon="Refresh" @click="handleReset">重置</el-button>
        </template>
        <template #right>
          <el-button
            type="primary"
            :icon="Plus"
            @click="handleAdd"
            v-if="authStore.hasPermission('material_create')"
          >
            新增物资
          </el-button>
          <el-button
            type="success"
            :icon="Upload"
            @click="handleImport"
            v-if="authStore.hasPermission('material_import')"
          >
            导入
          </el-button>
          <el-button
            type="warning"
            :icon="Download"
            @click="handleExport"
            v-if="authStore.hasPermission('material_export')"
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
        <el-table-column prop="code" label="物资编码" width="130" />
        <el-table-column prop="name" label="物资名称" min-width="150" show-overflow-tooltip />
        <el-table-column prop="category" label="分类" width="100">
          <template #default="scope">
            <el-tag size="small">{{ scope.row.category || '-' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="specification" label="规格型号" width="120" show-overflow-tooltip />
        <el-table-column prop="unit" label="单位" width="80" />
        <el-table-column prop="planned_quantity" label="计划数量" width="100" align="right">
          <template #default="scope">
            {{ scope.row.planned_quantity || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="quantity" label="库存数量" width="100" align="right">
          <template #default="scope">
            <span :class="getStockClass(scope.row.quantity || 0)">
              {{ scope.row.quantity || 0 }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="到货进度" width="150" align="center">
          <template #default="scope">
            <div v-if="scope.row.planned_quantity && scope.row.planned_quantity > 0" class="arrival-progress">
              <el-progress
                :percentage="Math.min(scope.row.arrival_percentage || 0, 100)"
                :status="getProgressStatus(scope.row)"
                :stroke-width="8"
              />
            </div>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="project_name" label="项目" min-width="150" show-overflow-tooltip>
          <template #default="scope">
            {{ scope.row.project_name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="150" show-overflow-tooltip />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="scope">
            <el-button
              type="primary"
              size="small"
              :icon="Edit"
              @click="handleEdit(scope.row)"
              v-if="authStore.hasPermission('material_edit')"
            >
              编辑
            </el-button>
            <el-button
              type="danger"
              size="small"
              :icon="Delete"
              @click="handleDelete(scope.row)"
              v-if="authStore.hasPermission('material_delete')"
            >
              删除
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

    <!-- 新增/编辑对话框 -->
    <Dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="700px"
      :loading="dialogLoading"
      @confirm="handleSubmit"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="物资编码" prop="code">
              <el-input
                v-model="formData.code"
                placeholder="请输入物资编码"
                maxlength="50"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="物资名称" prop="name">
              <el-input
                v-model="formData.name"
                placeholder="请输入物资名称"
                maxlength="100"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="分类" prop="category">
              <el-select v-model="formData.category" placeholder="请选择分类" style="width: 100%">
                <el-option
                  v-for="cat in categoryList"
                  :key="cat.id"
                  :label="cat.name"
                  :value="cat.name"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="规格型号" prop="specification">
              <el-input
                v-model="formData.specification"
                placeholder="请输入规格型号"
                maxlength="100"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="单位" prop="unit">
              <el-select v-model="formData.unit" placeholder="请选择单位" style="width: 100%">
                <el-option label="吨" value="吨" />
                <el-option label="千克" value="千克" />
                <el-option label="米" value="米" />
                <el-option label="平方米" value="平方米" />
                <el-option label="立方米" value="立方米" />
                <el-option label="个" value="个" />
                <el-option label="件" value="件" />
                <el-option label="箱" value="箱" />
                <el-option label="套" value="套" />
                <el-option label="千克" value="千克" />
                <el-option label="公斤" value="公斤" />
                <el-option label="g" value="g" />
                <el-option label="kg" value="kg" />
                <el-option label="m" value="m" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="单价(元)" prop="price">
              <el-input-number
                v-model="formData.price"
                :min="0"
                :precision="2"
                :step="0.01"
                placeholder="请输入单价"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="计划数量" prop="quantity">
              <el-input-number
                v-model="formData.quantity"
                :min="0"
                :step="1"
                placeholder="请输入计划数量"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="所属项目" prop="project_id">
              <ProjectSelector
                v-model="formData.project_id"
                :projects="projectList"
                placeholder="请选择项目"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="备注" prop="remark">
          <el-input
            v-model="formData.remark"
            type="textarea"
            :rows="3"
            placeholder="请输入备注"
            maxlength="500"
          />
        </el-form-item>
      </el-form>
    </Dialog>

    <!-- 智能导入对话框 -->
    <el-dialog
      v-model="importDialogVisible"
      title="智能导入物资"
      width="900px"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
      draggable
    >
      <ExcelImportMapper
        :fields="materialFields"
        :import-api="handleBatchImport"
        :projects="projectList"
        @success="handleImportSuccess"
        @close="importDialogVisible = false"
      />
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { materialApi } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Search,
  Refresh,
  Plus,
  Upload,
  Download,
  Edit,
  Delete,
  UploadFilled,
  Filter
} from '@element-plus/icons-vue'
import Dialog from '@/components/common/Dialog.vue'
import TableToolbar from '@/components/common/TableToolbar.vue'
import ExcelImportMapper from '@/components/common/ExcelImportMapper.vue'
import ProjectSelector from '@/components/common/ProjectSelector.vue'

const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()

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
  unstored: false
})

// 项目列表
const projectList = ref([])

// 物资分类列表
const categoryList = ref([])

// 对话框
const dialogVisible = ref(false)
const dialogTitle = computed(() => formData.id ? '编辑物资' : '新增物资')
const dialogLoading = ref(false)
const formRef = ref(null)

// 表单数据
const formData = reactive({
  id: null,
  code: '',
  name: '',
  category: '',
  specification: '',
  unit: '',
  price: null,
  quantity: null,
  project_id: null,
  remark: ''
})

// 表单验证规则
const formRules = {
  name: [
    { required: true, message: '请输入物资名称', trigger: 'blur' }
  ],
  unit: [
    { required: true, message: '请选择单位', trigger: 'change' }
  ],
  project_id: [
    { required: true, message: '请选择所属项目', trigger: 'change' }
  ]
}

// 导入对话框
const importDialogVisible = ref(false)

/**
 * 物资字段定义
 * 用于 Excel 导入映射
 */
const materialFields = [
  {
    field: 'code',
    label: '物资编码',
    required: false,
    hint: '唯一编码，如：STL-001（可选）'
  },
  {
    field: 'name',
    label: '物资名称',
    required: true,
    hint: '物资名称（必填）'
  },
  {
    field: 'category',
    label: '分类',
    required: false,
    hint: '如：钢材、水泥、电气材料等（可选）'
  },
  {
    field: 'specification',
    label: '规格型号',
    required: true,
    hint: '如：Φ12mm、C30等（必填）'
  },
  {
    field: 'unit',
    label: '单位',
    required: true,
    hint: '如：吨、米、个、件等（必填）'
  },
  {
    field: 'price',
    label: '单价',
    required: false,
    hint: '单位价格，元（可选）'
  },
  {
    field: 'quantity',
    label: '数量',
    required: true,
    hint: '物资数量（必填）'
  },
  {
    field: 'quality_standard',
    label: '质量标准',
    required: true,
    hint: '质量标准或技术要求（必填）'
  },
  {
    field: 'remark',
    label: '备注',
    required: true,
    hint: '其他说明（必填）'
  }
]

/**
 * 批量导入物资
 */
const handleBatchImport = async (data) => {
  return await materialApi.batchImport(data)
}

/**
 * 导入成功处理
 */
const handleImportSuccess = () => {
  fetchData()
}

// 获取列表数据
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
      filter: searchForm.unstored ? 'unstored' : undefined
    }
    // 使用解构赋值获取 data 和 pagination（适配统一响应格式）
    const { data, pagination: pag } = await materialApi.getList(params)
    tableData.value = data || []
    pagination.total = pag?.total || 0
  } catch (error) {
    // 错误已被 request.js 自动处理
    console.error('获取物资列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 收集项目及其所有子项目的ID
const collectProjectIds = (projectId, projectTree) => {
  const ids = [projectId]

  // 在树形结构中查找项目
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
const fetchProjects = async () => {
  try {
    const { projectApi } = await import('@/api')
    // 使用解构赋值（适配统一响应格式）
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
const fetchCategories = async () => {
  try {
    // 使用解构赋值（适配统一响应格式）
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
  searchForm.unstored = false
  pagination.page = 1
  fetchData()
}

// 切换未入库筛选
const toggleUnstored = () => {
  searchForm.unstored = !searchForm.unstored
  handleSearch()
}

// 新增
const handleAdd = () => {
  resetForm()
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row) => {
  Object.assign(formData, {
    id: row.id,
    code: row.code,
    name: row.name,
    category: row.category,
    specification: row.specification || '',
    unit: row.unit,
    price: row.price,
    quantity: row.quantity || row.planned_quantity || null,
    project_id: row.project_id || null,
    remark: row.remark || ''
  })
  dialogVisible.value = true
}

// 删除
const handleDelete = (row) => {
  ElMessageBox.confirm(
    `确定要删除物资"${row.name}"吗？`,
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await materialApi.delete(row.id)
      ElMessage.success('删除成功')
      fetchData()
    } catch (error) {
      console.error('删除失败:', error)
      ElMessage.error(error?.message || '删除失败，请重试')
    }
  }).catch((error) => {
    // 用户取消操作，不需要提示
    if (error !== 'cancel') {
      console.error('操作失败:', error)
    }
  })
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    dialogLoading.value = true

    const data = {
      code: formData.code,
      name: formData.name,
      category: formData.category,
      specification: formData.specification,
      unit: formData.unit,
      price: formData.price,
      quantity: formData.quantity,
      project_id: formData.project_id ? String(formData.project_id) : '',
      remark: formData.remark
    }

    if (formData.id) {
      await materialApi.update(formData.id, data)
      ElMessage.success('更新成功')
    } else {
      await materialApi.create(data)
      ElMessage.success('创建成功')
    }

    dialogVisible.value = false
    fetchData()
  } catch (error) {
    console.error('提交失败:', error)
  } finally {
    dialogLoading.value = false
  }
}

// 导入
const handleImport = async () => {
  // 确保项目列表已加载
  if (projectList.value.length === 0) {
    await fetchProjects()
  }
  importDialogVisible.value = true
}

// 导出
const handleExport = async () => {
  try {
    const response = await materialApi.export(searchForm)
    // 处理文件下载
    const blob = new Blob([response], {
      type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
    })
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `物资列表_${new Date().getTime()}.xlsx`
    a.click()
    window.URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
  } catch (error) {
    console.error('导出失败:', error)
  }
}

// 重置表单
const resetForm = () => {
  Object.assign(formData, {
    id: null,
    code: '',
    name: '',
    category: '',
    specification: '',
    unit: '',
    price: null,
    quantity: null,
    project_id: null,
    remark: ''
  })
  if (formRef.value) {
    formRef.value.clearValidate()
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

// 获取进度条状态
const getProgressStatus = (row) => {
  if (row.is_fully_arrived) return 'success'
  const percentage = row.arrival_percentage || 0
  if (percentage >= 80) return 'warning'
  return ''
}

// 获取库存数量样式类
const getStockClass = (quantity) => {
  if (quantity <= 0) return 'stock-danger'
  if (quantity < 10) return 'stock-warning'
  return 'stock-normal'
}

// 监听路由参数，自动触发操作
watch(() => route.query, (query) => {
  if (query.action === 'create') {
    // 重置查询参数，避免重复触发
    router.replace({ query: {} })
    // 稍微延迟，等待页面加载完成
    setTimeout(() => {
      handleAdd()
    }, 100)
  }
}, { immediate: true })

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

onMounted(() => {
  fetchProjects()
  fetchCategories()
  fetchData()
})
</script>

<style scoped>
.materials-container {
  padding: 0;
}

.upload-area {
  margin: 20px 0;
}

.mt-20 {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.arrival-progress {
  padding: 4px 0;
}

.progress-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
  font-size: 12px;
}

.progress-text {
  color: #606266;
  font-weight: 500;
}

.progress-percent {
  color: #409eff;
  font-weight: bold;
}

.fully-arrived-tag {
  margin-top: 4px;
  text-align: center;
}

/* 库存数量样式 */
.stock-danger {
  color: #f56c6c;
  font-weight: bold;
}

.stock-warning {
  color: #e6a23c;
  font-weight: bold;
}

.stock-normal {
  color: #67c23a;
  font-weight: bold;
}
</style>
