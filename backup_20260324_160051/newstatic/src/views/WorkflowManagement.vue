<template>
  <div class="workflow-management">
    <el-card shadow="never">
      <!-- 工具栏 -->
      <TableToolbar>
        <template #left>
          <el-select
            v-model="searchForm.module"
            placeholder="选择模块"
            clearable
            style="width: 200px"
            @change="handleSearch"
          >
            <el-option label="入库管理" value="inbound" />
            <el-option label="领用管理" value="requisition" />
          </el-select>
          <el-button type="primary" :icon="Search" @click="handleSearch">
            搜索
          </el-button>
          <el-button :icon="Refresh" @click="handleReset">重置</el-button>
        </template>
        <template #right>
          <el-button
            type="primary"
            :icon="Plus"
            @click="handleCreate"
            v-if="authStore.hasPermission('workflow_create')"
          >
            新建工作流
          </el-button>
        </template>
      </TableToolbar>

      <!-- 工作流列表 -->
      <el-table
        v-loading="loading"
        :data="tableData"
        border
        stripe
        style="width: 100%; margin-top: 16px"
      >
        <!-- 序号列已移除 -->
        <el-table-column prop="name" label="工作流名称" min-width="200" show-overflow-tooltip />
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="module" label="关联模块" width="120">
          <template #default="scope">
            <el-tag size="small">{{ getModuleText(scope.row.module) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="version" label="版本" width="80" align="center" />
        <el-table-column prop="is_active" label="状态" width="100" align="center">
          <template #default="scope">
            <el-tag :type="scope.row.is_active ? 'success' : 'info'" size="small">
              {{ scope.row.is_active ? '已激活' : '未激活' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160" />
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="scope">
            <el-button
              type="primary"
              size="small"
              :icon="Edit"
              @click="handleEdit(scope.row)"
              v-if="authStore.hasPermission('workflow_edit')"
            >
              编辑
            </el-button>
            <el-button
              :type="scope.row.is_active ? 'warning' : 'success'"
              size="small"
              @click="handleToggleActive(scope.row)"
              v-if="authStore.hasPermission('workflow_activate')"
            >
              {{ scope.row.is_active ? '停用' : '激活' }}
            </el-button>
            <el-button
              type="danger"
              size="small"
              :icon="Delete"
              @click="handleDelete(scope.row)"
              v-if="authStore.hasPermission('workflow_delete')"
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
        :page-sizes="[10, 20, 50]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
        style="margin-top: 16px; display: flex; justify-content: flex-end"
      />
    </el-card>

    <!-- 工作流编辑器对话框 -->
    <el-dialog
      v-model="editorDialogVisible"
      :title="editorTitle"
      width="90%"
      top="5vh"
      :close-on-click-modal="false"
    >
      <el-form :model="workflowForm" label-width="120px" style="margin-bottom: 16px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="工作流名称" required>
              <el-input v-model="workflowForm.name" placeholder="请输入工作流名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="关联模块" required>
              <el-select v-model="workflowForm.module" placeholder="请选择模块" :disabled="!!editingWorkflow">
                <el-option label="计划管理" value="material_plan" />
                <el-option label="入库管理" value="inbound" />
                <el-option label="领用管理" value="requisition" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="描述">
          <el-input v-model="workflowForm.description" type="textarea" :rows="2" placeholder="请输入描述" />
        </el-form-item>
      </el-form>

      <!-- 工作流编辑器 -->
      <div style="height: 500px; border: 1px solid #e0e0e0; border-radius: 4px; overflow: hidden">
        <WorkflowEditor
          ref="editorRef"
          :workflow-id="editingWorkflow?.id"
          :module="workflowForm.module"
          @save="handleSaveWorkflow"
          @change="handleWorkflowChange"
        />
      </div>

      <template #footer>
        <el-button @click="editorDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveAndClose" :loading="saving">
          保存并关闭
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, nextTick } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { workflowApi } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Plus, Edit, Delete } from '@element-plus/icons-vue'
import TableToolbar from '@/components/common/TableToolbar.vue'
import WorkflowEditor from '@/components/workflow/WorkflowEditor.vue'

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
  module: ''
})

// 编辑器对话框
const editorDialogVisible = ref(false)
const editorTitle = computed(() => editingWorkflow.value ? '编辑工作流' : '新建工作流')
const editorRef = ref(null)
const editingWorkflow = ref(null)
const saving = ref(false)
const hasChanges = ref(false)

// 工作流表单
const workflowForm = reactive({
  name: '',
  description: '',
  module: ''
})

// 获取列表数据
// 适配统一响应格式
const fetchData = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      module: searchForm.module || undefined
    }
    const { data, pagination: pag } = await workflowApi.getList(params)
    tableData.value = data || []
    pagination.total = pag?.total || 0
  } catch (error) {
    console.error('获取工作流列表失败:', error)
    ElMessage.error('获取工作流列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchData()
}

// 重置
const handleReset = () => {
  searchForm.module = ''
  pagination.page = 1
  fetchData()
}

// 新建工作流
const handleCreate = () => {
  editingWorkflow.value = null
  Object.assign(workflowForm, {
    name: '',
    description: '',
    module: searchForm.module || ''
  })
  hasChanges.value = false
  editorDialogVisible.value = true
}

// 编辑工作流
const handleEdit = async (row) => {
  editingWorkflow.value = row
  Object.assign(workflowForm, {
    name: row.name,
    description: row.description || '',
    module: row.module
  })
  hasChanges.value = false
  editorDialogVisible.value = true

  // 获取工作流详情并加载到编辑器
  try {
    const { data } = await workflowApi.getDetail(row.id)
    // 等待编辑器组件挂载后再设置数据
    await nextTick()
    if (editorRef.value && data) {
      editorRef.value.setWorkflowData(data)
    }
  } catch (error) {
    console.error('加载工作流失败:', error)
    ElMessage.error('加载工作流失败')
  }
}

// 保存工作流
const handleSaveWorkflow = async (workflowData) => {
  if (!workflowForm.name) {
    ElMessage.warning('请输入工作流名称')
    return
  }
  if (!workflowForm.module) {
    ElMessage.warning('请选择关联模块')
    return
  }

  saving.value = true
  try {
    const data = {
      name: workflowForm.name,
      description: workflowForm.description,
      module: workflowForm.module,
      ...workflowData
    }

    if (editingWorkflow.value) {
      await workflowApi.update(editingWorkflow.value.id, data)
      ElMessage.success('工作流更新成功')
    } else {
      await workflowApi.create(data)
      ElMessage.success('工作流创建成功')
    }

    hasChanges.value = false
    editorDialogVisible.value = false
    fetchData()
  } catch (error) {
    console.error('保存工作流失败:', error)
    ElMessage.error('保存工作流失败: ' + (error.message || '未知错误'))
  } finally {
    saving.value = false
  }
}

// 保存并关闭
const handleSaveAndClose = () => {
  if (editorRef.value) {
    editorRef.value.handleSave()
  }
}

// 工作流变更
const handleWorkflowChange = () => {
  hasChanges.value = true
}

// 切换激活状态
const handleToggleActive = (row) => {
  const action = row.is_active ? '停用' : '激活'
  ElMessageBox.confirm(
    `确定要${action}工作流"${row.name}"吗？`,
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      if (row.is_active) {
        await workflowApi.deactivate(row.id)
      } else {
        await workflowApi.activate(row.id)
      }
      ElMessage.success(`${action}成功`)
      fetchData()
    } catch (error) {
      console.error(`${action}失败:`, error)
      ElMessage.error(`${action}失败`)
    }
  }).catch((error) => {
    // 用户取消操作，不需要提示
    if (error !== 'cancel') {
      console.error('操作失败:', error)
    }
  })
}

// 删除工作流
const handleDelete = (row) => {
  ElMessageBox.confirm(
    `确定要删除工作流"${row.name}"吗？删除后不可恢复。`,
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await workflowApi.delete(row.id)
      ElMessage.success('删除成功')
      fetchData()
    } catch (error) {
      console.error('删除失败:', error)
      ElMessage.error('删除失败')
    }
  }).catch((error) => {
    // 用户取消操作，不需要提示
    if (error !== 'cancel') {
      console.error('操作失败:', error)
    }
  })
}

// 获取模块文本
const getModuleText = (module) => {
  const moduleMap = {
    material_plan: '计划管理',
    inbound: '入库管理',
    requisition: '领用管理'
  }
  return moduleMap[module] || module
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

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.workflow-management {
  padding: 0;
}
</style>
