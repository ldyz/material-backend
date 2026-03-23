<template>
  <el-dialog
    :model-value="visible"
    @update:model-value="val => visible = val"
    title="资源库管理"
    width="900px"
    @close="handleClose"
  >
    <!-- 工具栏 -->
    <div class="resource-toolbar">
      <el-button type="primary" :icon="Plus" @click="handleAddResource">
        添加资源
      </el-button>
      <el-select
        v-model="filterType"
        placeholder="筛选类型"
        clearable
        style="width: 150px; margin-left: 12px"
        @change="loadResources"
      >
        <el-option label="全部" value="" />
        <el-option label="人力" value="labor" />
        <el-option label="工机具" value="equipment" />
        <el-option label="材料" value="material" />
      </el-select>
      <el-input
        v-model="searchKeyword"
        placeholder="搜索资源名称"
        clearable
        style="width: 200px; margin-left: 12px"
        @input="handleSearch"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
    </div>

    <!-- 资源列表 -->
    <el-table
      :data="filteredResources"
      border
      stripe
      style="margin-top: 16px"
      v-loading="loading"
    >
      <el-table-column prop="name" label="资源名称" width="180" />
      <el-table-column prop="type" label="类型" width="100">
        <template #default="{ row }">
          <el-tag :type="getResourceTypeTag(row.type)">
            {{ getResourceTypeName(row.type) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="unit" label="单位" width="80" />
      <el-table-column prop="quantity" label="可用数量" width="100">
        <template #default="{ row }">
          {{ row.quantity }} {{ row.unit }}
        </template>
      </el-table-column>
      <el-table-column prop="cost_per_unit" label="单位成本" width="100">
        <template #default="{ row }">
          ¥{{ row.cost_per_unit }}
        </template>
      </el-table-column>
      <el-table-column prop="color" label="标识颜色" width="100">
        <template #default="{ row }">
          <el-color-picker v-model="row.color" show-alpha disabled />
        </template>
      </el-table-column>
      <el-table-column prop="is_active" label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="row.is_active ? 'success' : 'info'">
            {{ row.is_active ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="handleEditResource(row)">
            编辑
          </el-button>
          <el-button link type="danger" @click="handleDeleteResource(row)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 资源编辑对话框 -->
    <el-dialog
      v-model="editDialogVisible"
      :title="isEditing ? '编辑资源' : '添加资源'"
      width="500px"
      append-to-body
    >
      <el-form
        ref="resourceFormRef"
        :model="resourceForm"
        :rules="resourceFormRules"
        label-width="100px"
      >
        <el-form-item label="资源名称" prop="name">
          <el-input v-model="resourceForm.name" placeholder="请输入资源名称" />
        </el-form-item>

        <el-form-item label="资源类型" prop="type">
          <el-select v-model="resourceForm.type" placeholder="请选择类型">
            <el-option label="人力" value="labor" />
            <el-option label="工机具" value="equipment" />
            <el-option label="材料" value="material" />
          </el-select>
        </el-form-item>

        <el-form-item label="计量单位" prop="unit">
          <el-select v-model="resourceForm.unit" placeholder="请选择单位">
            <el-option label="人/d" value="人/d" />
            <el-option label="台/d" value="台/d" />
            <el-option label="kg" value="kg" />
            <el-option label="m" value="m" />
            <el-option label="m²" value="m²" />
            <el-option label="m³" value="m³" />
            <el-option label="件" value="件" />
            <el-option label="套" value="套" />
          </el-select>
        </el-form-item>

        <el-form-item label="可用数量" prop="quantity">
          <el-input-number
            v-model="resourceForm.quantity"
            :min="0"
            :step="1"
            :precision="2"
          />
        </el-form-item>

        <el-form-item label="单位成本" prop="cost_per_unit">
          <el-input-number
            v-model="resourceForm.cost_per_unit"
            :min="0"
            :step="0.01"
            :precision="2"
            :controls="false"
          />
          <span style="margin-left: 8px; color: #909399">元</span>
        </el-form-item>

        <el-form-item label="标识颜色" prop="color">
          <el-color-picker v-model="resourceForm.color" show-alpha />
        </el-form-item>

        <el-form-item label="状态">
          <el-switch
            v-model="resourceForm.is_active"
            active-text="启用"
            inactive-text="禁用"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveResource" :loading="saving">
          保存
        </el-button>
      </template>
    </el-dialog>

    <template #footer>
      <el-button @click="visible = false">关闭</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch, toValue } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search } from '@element-plus/icons-vue'
import { progressApi } from '@/api'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  projectId: {
    type: [Number, String],
    required: true
  }
})

const emit = defineEmits(['update:modelValue', 'refresh'])

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const loading = ref(false)
const saving = ref(false)
const resources = ref([])
const filterType = ref('')
const searchKeyword = ref('')
const editDialogVisible = ref(false)
const isEditing = ref(false)
const resourceFormRef = ref(null)

// 资源表单
const resourceForm = ref({
  id: null,
  name: '',
  type: 'labor',
  unit: '人/d',
  quantity: 0,
  cost_per_unit: 0,
  color: '#409eff',
  is_active: true
})

// 表单验证规则
const resourceFormRules = {
  name: [{ required: true, message: '请输入资源名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择资源类型', trigger: 'change' }],
  unit: [{ required: true, message: '请选择计量单位', trigger: 'change' }],
  quantity: [{ required: true, message: '请输入可用数量', trigger: 'blur' }]
}

// 过滤后的资源列表
const filteredResources = computed(() => {
  let result = resources.value

  // 类型筛选
  if (filterType.value) {
    result = result.filter(r => r.type === filterType.value)
  }

  // 关键词搜索
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    result = result.filter(r =>
      r.name.toLowerCase().includes(keyword)
    )
  }

  return result
})

// 获取资源类型标签样式
const getResourceTypeTag = (type) => {
  const tagMap = {
    labor: 'success',
    equipment: 'warning',
    material: 'info'
  }
  return tagMap[type] || ''
}

// 获取资源类型名称
const getResourceTypeName = (type) => {
  const nameMap = {
    labor: '人力',
    equipment: '工机具',
    material: '材料'
  }
  return nameMap[type] || type
}

// 加载资源列表
const loadResources = async () => {
  loading.value = true
  try {
    const response = await progressApi.getProjectResources(toValue(props.projectId))
    resources.value = response.data || []
    console.log('资源列表加载成功:', resources.value)
  } catch (error) {
    console.error('加载资源列表失败:', error)
    ElMessage.error('加载资源列表失败: ' + (error.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

// 搜索处理
const handleSearch = () => {
  // 搜索逻辑由 filteredResources computed 自动处理
}

// 添加资源
const handleAddResource = () => {
  isEditing.value = false
  resourceForm.value = {
    id: null,
    name: '',
    type: 'labor',
    unit: '人/d',
    quantity: 0,
    cost_per_unit: 0,
    color: '#409eff',
    is_active: true
  }
  editDialogVisible.value = true
}

// 编辑资源
const handleEditResource = (row) => {
  isEditing.value = true
  resourceForm.value = {
    id: row.id,
    name: row.name,
    type: row.type,
    unit: row.unit,
    quantity: row.quantity,
    cost_per_unit: row.cost_per_unit,
    color: row.color || '#409eff',
    is_active: row.is_active
  }
  editDialogVisible.value = true
}

// 保存资源
const handleSaveResource = async () => {
  if (!resourceFormRef.value) return

  try {
    const valid = await resourceFormRef.value.validate()
    if (!valid) return

    saving.value = true
    try {
      const projectIdValue = toValue(props.projectId)
      console.log('DEBUG - props.projectId:', props.projectId)
      console.log('DEBUG - projectIdValue:', projectIdValue)
      console.log('DEBUG - typeof projectIdValue:', typeof projectIdValue)

      const data = {
        project_id: projectIdValue,
        name: resourceForm.value.name,
        type: resourceForm.value.type,
        unit: resourceForm.value.unit,
        quantity: resourceForm.value.quantity,
        cost_per_unit: resourceForm.value.cost_per_unit,
        color: resourceForm.value.color,
        is_active: resourceForm.value.is_active
      }

      console.log('DEBUG - data.project_id:', data.project_id, 'type:', typeof data.project_id)
      console.log('DEBUG - isEditing:', isEditing.value)
      console.log('DEBUG - resourceForm.value.id:', resourceForm.value.id)
      console.log('DEBUG - Calling API with projectId:', projectIdValue, 'resourceId:', resourceForm.value.id)

      if (isEditing.value) {
        // 更新资源
        console.log('DEBUG - Calling updateResource with:', projectIdValue, resourceForm.value.id, data)
        await progressApi.updateResource(projectIdValue, resourceForm.value.id, data)
        ElMessage.success('资源更新成功')
      } else {
        // 创建资源
        console.log('DEBUG - Calling createResource with:', projectIdValue, data)
        await progressApi.createResource(projectIdValue, data)
        ElMessage.success('资源添加成功')
      }

      editDialogVisible.value = false
      await loadResources()
      emit('refresh')
    } catch (error) {
      console.error('保存资源失败:', error)
      ElMessage.error('保存资源失败: ' + (error.message || '未知错误'))
    } finally {
      saving.value = false
    }
  } catch (validationError) {
    // Validation failed, do nothing
    console.log('表单验证失败')
  }
}

// 删除资源
const handleDeleteResource = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除资源 "${row.name}" 吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    try {
      await progressApi.deleteResource(toValue(props.projectId), row.id)
      ElMessage.success('资源删除成功')
      await loadResources()
      emit('refresh')
    } catch (error) {
      console.error('删除资源失败:', error)
      ElMessage.error('删除资源失败: ' + (error.message || '未知错误'))
    }
  } catch {
    // 用户取消删除
  }
}

const handleClose = () => {
  visible.value = false
}

// 监听对话框打开
watch(() => props.modelValue, (val) => {
  if (val) {
    console.log('ResourceManagementDialog opened, props.projectId:', props.projectId)
    loadResources()
  }
})
</script>

<style scoped>
.resource-toolbar {
  display: flex;
  align-items: center;
}

.el-table {
  font-size: 13px;
}

:deep(.el-table__cell) {
  padding: 8px 0;
}
</style>
