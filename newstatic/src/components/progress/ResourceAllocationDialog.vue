<template>
  <el-dialog
    v-model="visible"
    title="分配资源"
    width="700px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <!-- 已分配资源列表 -->
    <div class="allocated-resources">
      <h4>已分配资源</h4>
      <el-table :data="allocatedResources" border max-height="250" empty-text="暂无分配资源">
        <el-table-column prop="resource_name" label="资源名称" width="150" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getResourceTypeTag(row.type)" size="small">
              {{ getResourceTypeName(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="unit" label="单位" width="80" />
        <el-table-column prop="quantity" label="分配数量" width="140">
          <template #default="{ row }">
            <el-input-number
              v-model="row.quantity"
              :min="0"
              :step="1"
              :precision="2"
              size="small"
              controls-position="right"
              @change="handleUpdateQuantity(row)"
            />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="80" fixed="right">
          <template #default="{ $index }">
            <el-button link type="danger" @click="handleRemoveResource($index)">
              移除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 添加资源 -->
    <el-divider>添加资源</el-divider>
    <el-form :inline="true" :model="newResource">
      <el-form-item label="资源类型">
        <el-select
          v-model="newResource.type"
          placeholder="选择类型"
          style="width: 120px"
          @change="handleTypeChange"
        >
          <el-option label="人力" value="labor" />
          <el-option label="工机具" value="equipment" />
          <el-option label="材料" value="material" />
        </el-select>
      </el-form-item>
      <el-form-item label="资源">
        <el-select
          v-model="newResource.resource_id"
          filterable
          placeholder="选择资源"
          style="width: 200px"
        >
          <el-option
            v-for="res in filteredResources"
            :key="res.id"
            :label="`${res.name} (可用: ${res.quantity}${res.unit || ''})`"
            :value="res.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="数量">
        <el-input-number
          v-model="newResource.quantity"
          :min="0.01"
          :step="1"
          :precision="2"
          controls-position="right"
          style="width: 150px"
        />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="handleAddResource">添加</el-button>
      </el-form-item>
    </el-form>

    <!-- 资源使用统计 -->
    <div class="resource-summary">
      <el-descriptions :column="3" border size="small">
        <el-descriptions-item label="分配资源数">
          {{ allocatedResources.length }} 种
        </el-descriptions-item>
        <el-descriptions-item label="人力需求">
          {{ totalLabor }} 人/d
        </el-descriptions-item>
        <el-descriptions-item label="设备需求">
          {{ totalEquipment }} 台/d
        </el-descriptions-item>
      </el-descriptions>
    </div>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" @click="handleSave" :loading="saving">保存</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { progressApi } from '@/api'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  taskId: {
    type: [Number, String],
    default: null
  },
  projectId: {
    type: [Number, String],
    required: true
  }
})

const emit = defineEmits(['update:modelValue', 'saved'])

const visible = ref(false)
const saving = ref(false)
const allResources = ref([])
const allocatedResources = ref([])

const newResource = ref({
  type: 'labor',
  resource_id: null,
  quantity: 1
})

// 过滤后的资源列表
const filteredResources = computed(() => {
  if (!newResource.value.type) return []
  return allResources.value.filter(
    r => r.type === newResource.value.type && r.is_active
  )
})

// 统计信息
const totalLabor = computed(() => {
  return allocatedResources.value
    .filter(r => r.type === 'labor')
    .reduce((sum, r) => sum + (r.quantity || 0), 0)
    .toFixed(2)
})

const totalEquipment = computed(() => {
  return allocatedResources.value
    .filter(r => r.type === 'equipment')
    .reduce((sum, r) => sum + (r.quantity || 0), 0)
    .toFixed(2)
})

// 获取资源类型标签颜色
const getResourceTypeTag = (type) => {
  const colors = {
    labor: 'primary',
    equipment: 'success',
    material: 'warning'
  }
  return colors[type] || 'info'
}

// 获取资源类型名称
const getResourceTypeName = (type) => {
  const names = {
    labor: '人力',
    equipment: '工机具',
    material: '材料'
  }
  return names[type] || type
}

// 加载项目资源
const loadProjectResources = async () => {
  try {
    const response = await progressApi.getProjectResources(props.projectId)
    allResources.value = response.data || []
  } catch (error) {
    console.error('加载资源列表失败:', error)
    ElMessage.error('加载资源列表失败')
  }
}

// 加载任务资源分配
const loadTaskResources = async () => {
  if (!props.taskId) return

  try {
    const response = await progressApi.getTaskResources(props.taskId)
    allocatedResources.value = response.data || []
  } catch (error) {
    console.error('加载任务资源失败:', error)
  }
}

// 资源类型变化
const handleTypeChange = () => {
  newResource.value.resource_id = null
}

// 添加资源
const handleAddResource = () => {
  if (!newResource.value.resource_id) {
    ElMessage.warning('请选择资源')
    return
  }

  if (newResource.value.quantity <= 0) {
    ElMessage.warning('请输入有效的数量')
    return
  }

  // 检查是否已分配
  const existing = allocatedResources.value.find(
    r => r.resource_id === newResource.value.resource_id
  )

  if (existing) {
    ElMessage.warning('该资源已分配，请直接修改数量')
    return
  }

  // 添加到已分配列表
  const resource = allResources.value.find(
    r => r.id === newResource.value.resource_id
  )

  if (resource) {
    allocatedResources.value.push({
      id: null,
      task_id: props.taskId,
      resource_id: resource.id,
      resource_name: resource.name,
      type: resource.type,
      unit: resource.unit,
      color: resource.color,
      quantity: newResource.value.quantity
    })

    // 重置表单
    newResource.value = {
      type: newResource.value.type,
      resource_id: null,
      quantity: 1
    }

    ElMessage.success('资源已添加')
  }
}

// 更新数量
const handleUpdateQuantity = (row) => {
  // 数量已通过 v-model 直接更新
}

// 移除资源
const handleRemoveResource = (index) => {
  allocatedResources.value.splice(index, 1)
}

// 保存
const handleSave = async () => {
  saving.value = true

  try {
    console.log('开始保存资源分配，任务ID:', props.taskId)
    console.log('分配的资源:', allocatedResources.value)

    // 先删除现有分配
    // （这里简化处理，实际应该对比差异）
    for (const resource of allocatedResources.value) {
      console.log('分配资源:', {
        taskId: props.taskId,
        resourceId: resource.resource_id,
        quantity: resource.quantity
      })

      const response = await progressApi.allocateTaskResource(props.taskId, {
        resource_id: resource.resource_id,
        quantity: resource.quantity
      })

      console.log('资源分配响应:', response)
    }

    ElMessage.success('资源分配保存成功')
    emit('saved')
    handleClose()
  } catch (error) {
    console.error('保存资源分配失败:', error)
    console.error('错误详情:', error.response?.data)
    ElMessage.error('保存失败: ' + (error.response?.data?.error || error.response?.data?.message || error.message || '未知错误'))
  } finally {
    saving.value = false
  }
}

// 关闭对话框
const handleClose = () => {
  visible.value = false
  emit('update:modelValue', false)
}

// 监听显示状态
watch(() => props.modelValue, (val) => {
  visible.value = val
  if (val) {
    loadProjectResources()
    loadTaskResources()
  }
})

watch(visible, (val) => {
  if (!val) {
    emit('update:modelValue', false)
  }
})
</script>

<style scoped>
.allocated-resources {
  margin-bottom: 20px;
}

.allocated-resources h4 {
  margin: 0 0 12px 0;
  font-size: 14px;
  font-weight: 600;
  color: #303133;
}

.resource-summary {
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px dashed #dcdfe6;
}

:deep(.el-input-number) {
  width: 100%;
}

:deep(.el-input-number .el-input__inner) {
  text-align: left;
}
</style>
