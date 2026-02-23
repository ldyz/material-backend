<template>
  <el-dialog
    :model-value="visible"
    @update:model-value="val => { if (!val) handleClose() }"
    :title="editingTask ? '编辑任务' : '新建任务'"
    width="800px"
    :close-on-click-modal="false"
  >
    <el-tabs v-model="activeTab" type="border-card">
      <!-- 基本信息 -->
      <el-tab-pane label="基本信息" name="basic">
        <el-form :model="formData" :rules="rules" ref="formRef" label-width="100px">
          <el-form-item label="任务名称" prop="name">
            <el-input v-model="formData.name" placeholder="请输入任务名称" />
          </el-form-item>
          <el-form-item label="开始日期" prop="start">
            <el-date-picker
              v-model="formData.start"
              type="date"
              placeholder="选择开始日期"
              style="width: 100%"
              value-format="YYYY-MM-DD"
            />
          </el-form-item>
          <el-form-item label="结束日期" prop="end">
            <el-date-picker
              v-model="formData.end"
              type="date"
              placeholder="选择结束日期"
              style="width: 100%"
              value-format="YYYY-MM-DD"
            />
          </el-form-item>
          <el-form-item label="进度" prop="progress">
            <el-slider v-model="formData.progress" :marks="{ 0: '0%', 50: '50%', 100: '100%' }" />
          </el-form-item>
          <el-form-item label="优先级" prop="priority">
            <el-radio-group v-model="formData.priority">
              <el-radio-button label="urgent">紧急</el-radio-button>
              <el-radio-button label="high">高</el-radio-button>
              <el-radio-button label="medium">中</el-radio-button>
              <el-radio-button label="low">低</el-radio-button>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="备注">
            <el-input v-model="formData.notes" type="textarea" :rows="3" placeholder="任务备注" />
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <!-- 资源管理 -->
      <el-tab-pane label="资源分配" name="resources">
        <div class="resources-section">
          <div class="resources-header">
            <span class="resources-title">已分配资源</span>
            <el-button type="primary" size="small" @click="showAddResourceDialog">
              <el-icon><Plus /></el-icon>
              添加资源
            </el-button>
          </div>

          <el-table :data="formData.resources" style="width: 100%" max-height="300">
            <el-table-column prop="resource_name" label="资源名称" width="180">
              <template #default="{ row }">
                {{ row.resource_name || row.name }}
              </template>
            </el-table-column>
            <el-table-column prop="type" label="类型" width="100">
              <template #default="{ row }">
                <el-tag :type="getResourceTypeTag(row.type)" size="small">
                  {{ getResourceTypeName(row.type) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="quantity" label="数量" width="100" />
            <el-table-column prop="unit" label="单位" width="80" />
            <el-table-column prop="cost" label="成本" width="100">
              <template #default="{ row }">
                {{ row.cost ? `¥${row.cost}` : '-' }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="150">
              <template #default="{ row, $index }">
                <el-button link type="primary" size="small" @click="editResource(row, $index)">
                  编辑
                </el-button>
                <el-button link type="danger" size="small" @click="removeResource($index)">
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>

          <el-empty v-if="formData.resources.length === 0" description="暂无分配资源" :image-size="80" />
        </div>
      </el-tab-pane>

      <!-- 任务依赖 -->
      <el-tab-pane label="任务依赖" name="dependencies">
        <div class="dependencies-section">
          <!-- 紧前任务 -->
          <div class="dep-group">
            <div class="dep-header">
              <span class="dep-title">
                <el-icon><Back /></el-icon>
                紧前任务（前置任务）
              </span>
              <el-button type="primary" size="small" @click="showAddPredecessorDialog">
                <el-icon><Plus /></el-icon>
                添加
              </el-button>
            </div>
            <el-table :data="formData.predecessors" style="width: 100%" max-height="200">
              <el-table-column prop="name" label="任务名称" />
              <el-table-column prop="type" label="依赖类型" width="120">
                <template #default="{ row }">
                  <el-select v-model="row.type" size="small" style="width: 100%">
                    <el-option label="完成-开始 (FS)" value="FS" />
                    <el-option label="开始-开始 (SS)" value="SS" />
                    <el-option label="完成-完成 (FF)" value="FF" />
                    <el-option label="开始-完成 (SF)" value="SF" />
                  </el-select>
                </template>
              </el-table-column>
              <el-table-column prop="lag" label="滞后(天)" width="100">
                <template #default="{ row }">
                  <el-input-number v-model="row.lag" :min="0" size="small" controls-position="right" style="width: 100%" />
                </template>
              </el-table-column>
              <el-table-column label="操作" width="80">
                <template #default="{ $index }">
                  <el-button link type="danger" size="small" @click="removePredecessor($index)">
                    删除
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
            <el-empty v-if="formData.predecessors.length === 0" description="暂无紧前任务" :image-size="60" />
          </div>

          <!-- 紧后任务 -->
          <div class="dep-group">
            <div class="dep-header">
              <span class="dep-title">
                <el-icon><Right /></el-icon>
                紧后任务（后置任务）
              </span>
              <el-button type="primary" size="small" @click="showAddSuccessorDialog">
                <el-icon><Plus /></el-icon>
                添加
              </el-button>
            </div>
            <el-table :data="formData.successors" style="width: 100%" max-height="200">
              <el-table-column prop="name" label="任务名称" />
              <el-table-column prop="type" label="依赖类型" width="120">
                <template #default="{ row }">
                  <el-select v-model="row.type" size="small" style="width: 100%">
                    <el-option label="完成-开始 (FS)" value="FS" />
                    <el-option label="开始-开始 (SS)" value="SS" />
                    <el-option label="完成-完成 (FF)" value="FF" />
                    <el-option label="开始-完成 (SF)" value="SF" />
                  </el-select>
                </template>
              </el-table-column>
              <el-table-column prop="lag" label="滞后(天)" width="100">
                <template #default="{ row }">
                  <el-input-number v-model="row.lag" :min="0" size="small" controls-position="right" style="width: 100%" />
                </template>
              </el-table-column>
              <el-table-column label="操作" width="80">
                <template #default="{ $index }">
                  <el-button link type="danger" size="small" @click="removeSuccessor($index)">
                    删除
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
            <el-empty v-if="formData.successors.length === 0" description="暂无紧后任务" :image-size="60" />
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" @click="handleSave" :loading="saving">保存</el-button>
    </template>
  </el-dialog>

  <!-- 添加/编辑资源对话框 -->
  <el-dialog
    v-model="resourceDialogVisible"
    :title="editingResourceIndex >= 0 ? '编辑资源' : '添加资源'"
    width="500px"
  >
    <el-form :model="resourceForm" :rules="resourceRules" ref="resourceFormRef" label-width="100px">
      <el-form-item label="资源类型" prop="type">
        <el-select v-model="resourceForm.type" placeholder="请选择资源类型" style="width: 100%">
          <el-option label="人力" value="labor" />
          <el-option label="设备" value="equipment" />
          <el-option label="材料" value="material" />
        </el-select>
      </el-form-item>
      <el-form-item label="资源名称" prop="resource_id">
        <el-select
          v-model="resourceForm.resource_id"
          placeholder="请选择资源"
          filterable
          style="width: 100%"
          @change="handleResourceChange"
        >
          <el-option
            v-for="resource in filteredResources"
            :key="resource.id"
            :label="resource.name"
            :value="resource.id"
          >
            <span>{{ resource.name }}</span>
            <span style="color: #8492a6; font-size: 12px; margin-left: 8px">
              (库存: {{ resource.available_quantity || 0 }}{{ resource.unit }})
            </span>
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="数量" prop="quantity">
        <el-input-number v-model="resourceForm.quantity" :min="1" controls-position="right" style="width: 100%" />
      </el-form-item>
      <el-form-item label="成本" prop="cost">
        <el-input-number v-model="resourceForm.cost" :min="0" :precision="2" controls-position="right" style="width: 100%" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="resourceDialogVisible = false">取消</el-button>
      <el-button type="primary" @click="saveResource">确定</el-button>
    </template>
  </el-dialog>

  <!-- 添加紧前任务对话框 -->
  <el-dialog
    v-model="predecessorDialogVisible"
    title="添加紧前任务"
    width="500px"
  >
    <el-form :model="dependencyForm" ref="dependencyFormRef" label-width="100px">
      <el-form-item label="选择任务" prop="taskId">
        <el-select
          v-model="dependencyForm.taskId"
          placeholder="请选择任务"
          filterable
          style="width: 100%"
        >
          <el-option
            v-for="task in availablePredecessorTasks"
            :key="task.id"
            :label="task.name"
            :value="task.id"
          >
            <span>{{ task.name }}</span>
            <span style="color: #8492a6; font-size: 12px; margin-left: 8px">
              ({{ task.start }} ~ {{ task.end }})
            </span>
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="依赖类型" prop="type">
        <el-select v-model="dependencyForm.type" style="width: 100%">
          <el-option label="完成-开始 (FS)" value="FS">
            <div>
              <div>完成-开始 (FS)</div>
              <div style="color: #8492a6; font-size: 12px">前置任务完成后，后置任务才能开始</div>
            </div>
          </el-option>
          <el-option label="开始-开始 (SS)" value="SS">
            <div>
              <div>开始-开始 (SS)</div>
              <div style="color: #8492a6; font-size: 12px">前置任务开始后，后置任务才能开始</div>
            </div>
          </el-option>
          <el-option label="完成-完成 (FF)" value="FF">
            <div>
              <div>完成-完成 (FF)</div>
              <div style="color: #8492a6; font-size: 12px">前置任务完成后，后置任务才能完成</div>
            </div>
          </el-option>
          <el-option label="开始-完成 (SF)" value="SF">
            <div>
              <div>开始-完成 (SF)</div>
              <div style="color: #8492a6; font-size: 12px">前置任务开始后，后置任务才能完成</div>
            </div>
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="滞后时间" prop="lag">
        <el-input-number v-model="dependencyForm.lag" :min="0" controls-position="right" style="width: 100%" />
        <span style="margin-left: 8px; color: #909399">天</span>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="predecessorDialogVisible = false">取消</el-button>
      <el-button type="primary" @click="addPredecessor">确定</el-button>
    </template>
  </el-dialog>

  <!-- 添加紧后任务对话框 -->
  <el-dialog
    v-model="successorDialogVisible"
    title="添加紧后任务"
    width="500px"
  >
    <el-form :model="dependencyForm" ref="dependencyFormRef" label-width="100px">
      <el-form-item label="选择任务" prop="taskId">
        <el-select
          v-model="dependencyForm.taskId"
          placeholder="请选择任务"
          filterable
          style="width: 100%"
        >
          <el-option
            v-for="task in availableSuccessorTasks"
            :key="task.id"
            :label="task.name"
            :value="task.id"
          >
            <span>{{ task.name }}</span>
            <span style="color: #8492a6; font-size: 12px; margin-left: 8px">
              ({{ task.start }} ~ {{ task.end }})
            </span>
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="依赖类型" prop="type">
        <el-select v-model="dependencyForm.type" style="width: 100%">
          <el-option label="完成-开始 (FS)" value="FS">
            <div>
              <div>完成-开始 (FS)</div>
              <div style="color: #8492a6; font-size: 12px">前置任务完成后，后置任务才能开始</div>
            </div>
          </el-option>
          <el-option label="开始-开始 (SS)" value="SS">
            <div>
              <div>开始-开始 (SS)</div>
              <div style="color: #8492a6; font-size: 12px">前置任务开始后，后置任务才能开始</div>
            </div>
          </el-option>
          <el-option label="完成-完成 (FF)" value="FF">
            <div>
              <div>完成-完成 (FF)</div>
              <div style="color: #8492a6; font-size: 12px">前置任务完成后，后置任务才能完成</div>
            </div>
          </el-option>
          <el-option label="开始-完成 (SF)" value="SF">
            <div>
              <div>开始-完成 (SF)</div>
              <div style="color: #8492a6; font-size: 12px">前置任务开始后，后置任务才能完成</div>
            </div>
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="滞后时间" prop="lag">
        <el-input-number v-model="dependencyForm.lag" :min="0" controls-position="right" style="width: 100%" />
        <span style="margin-left: 8px; color: #909399">天</span>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="successorDialogVisible = false">取消</el-button>
      <el-button type="primary" @click="addSuccessor">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, watch, computed, onMounted, onUnmounted } from 'vue'
import { Plus, Back, Right } from '@element-plus/icons-vue'
import { progressApi } from '@/api'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  editingTask: {
    type: Object,
    default: null
  },
  saving: {
    type: Boolean,
    default: false
  },
  // 所有任务列表（用于选择依赖关系）
  allTasks: {
    type: Array,
    default: () => []
  },
  // 资源库列表
  resourceLibrary: {
    type: Array,
    default: () => []
  },
  // 项目ID（用于加载资源）
  projectId: {
    type: [Number, String],
    default: null
  }
})

const emit = defineEmits(['update:visible', 'save'])

// 表单引用
const formRef = ref(null)
const resourceFormRef = ref(null)
const dependencyFormRef = ref(null)

// 当前激活的标签页
const activeTab = ref('basic')

// 默认表单数据
const defaultFormData = {
  name: '',
  start: '',
  end: '',
  progress: 0,
  priority: 'medium',
  notes: '',
  resources: [],
  predecessors: [],
  successors: []
}

const formData = ref({ ...defaultFormData })

const rules = {
  name: [{ required: true, message: '请输入任务名称', trigger: 'blur' }],
  start: [{ required: true, message: '请选择开始日期', trigger: 'change' }],
  end: [{ required: true, message: '请选择结束日期', trigger: 'change' }]
}

// ==================== 资源管理 ====================
const resourceDialogVisible = ref(false)
const editingResourceIndex = ref(-1)

const resourceForm = ref({
  type: 'labor',
  resource_id: null,
  quantity: 1,
  cost: 0
})

const resourceRules = {
  type: [{ required: true, message: '请选择资源类型', trigger: 'change' }],
  resource_id: [{ required: true, message: '请选择资源', trigger: 'change' }],
  quantity: [{ required: true, message: '请输入数量', trigger: 'blur' }]
}

// 根据类型过滤资源
const filteredResources = computed(() => {
  if (!props.resourceLibrary || props.resourceLibrary.length === 0) return []
  return props.resourceLibrary.filter(r => r.type === resourceForm.value.type)
})

// 获取资源类型标签颜色
const getResourceTypeTag = (type) => {
  const tags = {
    labor: 'primary',
    equipment: 'success',
    material: 'warning'
  }
  return tags[type] || 'info'
}

// 获取资源类型名称
const getResourceTypeName = (type) => {
  const names = {
    labor: '人力',
    equipment: '设备',
    material: '材料'
  }
  return names[type] || type
}

// 显示添加资源对话框
const showAddResourceDialog = () => {
  editingResourceIndex.value = -1
  resourceForm.value = {
    type: 'labor',
    resource_id: null,
    quantity: 1,
    cost: 0
  }
  resourceDialogVisible.value = true
}

// 编辑资源
const editResource = (resource, index) => {
  editingResourceIndex.value = index
  resourceForm.value = {
    type: resource.type,
    resource_id: resource.resource_id || resource.id,
    quantity: resource.quantity,
    cost: resource.cost || 0
  }
  resourceDialogVisible.value = true
}

// 资源选择变化
const handleResourceChange = (resourceId) => {
  const resource = props.resourceLibrary.find(r => r.id === resourceId)
  if (resource) {
    resourceForm.value.cost = resource.unit_cost || 0
    resourceForm.value.resource_name = resource.name
    resourceForm.value.unit = resource.unit
  }
}

// 保存资源
const saveResource = async () => {
  if (!resourceFormRef.value) return

  try {
    await resourceFormRef.value.validate()

    const resource = props.resourceLibrary.find(r => r.id === resourceForm.value.resource_id)
    if (!resource) {
      return
    }

    const resourceData = {
      id: resource.id,
      resource_id: resource.id,
      resource_name: resource.name,
      type: resourceForm.value.type,
      quantity: resourceForm.value.quantity,
      cost: resourceForm.value.cost,
      unit: resource.unit,
      available_quantity: resource.available_quantity
    }

    if (editingResourceIndex.value >= 0) {
      // 编辑现有资源
      formData.value.resources[editingResourceIndex.value] = resourceData
    } else {
      // 添加新资源
      formData.value.resources.push(resourceData)
    }

    resourceDialogVisible.value = false
  } catch (error) {
    // 验证失败
  }
}

// 删除资源
const removeResource = (index) => {
  formData.value.resources.splice(index, 1)
}

// ==================== 任务依赖管理 ====================
const predecessorDialogVisible = ref(false)
const successorDialogVisible = ref(false)

const dependencyForm = ref({
  taskId: null,
  type: 'FS',
  lag: 0
})

// 可选的紧前任务（排除当前任务和已选择的）
const availablePredecessorTasks = computed(() => {
  if (!props.allTasks || props.allTasks.length === 0) return []

  const selectedIds = formData.value.predecessors.map(p => p.id)
  const currentTaskId = props.editingTask?.id

  return props.allTasks.filter(task => {
    // 排除当前任务
    if (currentTaskId && task.id === currentTaskId) return false
    // 排除已选择的任务
    if (selectedIds.includes(task.id)) return false
    // 排除紧后任务（避免循环依赖）
    if (formData.value.successors.some(s => s.id === task.id)) return false
    return true
  })
})

// 可选的紧后任务（排除当前任务和已选择的）
const availableSuccessorTasks = computed(() => {
  if (!props.allTasks || props.allTasks.length === 0) return []

  const selectedIds = formData.value.successors.map(s => s.id)
  const currentTaskId = props.editingTask?.id

  return props.allTasks.filter(task => {
    // 排除当前任务
    if (currentTaskId && task.id === currentTaskId) return false
    // 排除已选择的任务
    if (selectedIds.includes(task.id)) return false
    // 排除紧前任务（避免循环依赖）
    if (formData.value.predecessors.some(p => p.id === task.id)) return false
    return true
  })
})

// 显示添加紧前任务对话框
const showAddPredecessorDialog = () => {
  dependencyForm.value = {
    taskId: null,
    type: 'FS',
    lag: 0
  }
  predecessorDialogVisible.value = true
}

// 显示添加紧后任务对话框
const showAddSuccessorDialog = () => {
  dependencyForm.value = {
    taskId: null,
    type: 'FS',
    lag: 0
  }
  successorDialogVisible.value = true
}

// 添加紧前任务
const addPredecessor = () => {
  if (!dependencyForm.value.taskId) {
    return
  }

  const task = props.allTasks.find(t => t.id === dependencyForm.value.taskId)
  if (!task) return

  formData.value.predecessors.push({
    id: task.id,
    name: task.name,
    type: dependencyForm.value.type,
    lag: dependencyForm.value.lag
  })

  predecessorDialogVisible.value = false
}

// 添加紧后任务
const addSuccessor = () => {
  if (!dependencyForm.value.taskId) {
    return
  }

  const task = props.allTasks.find(t => t.id === dependencyForm.value.taskId)
  if (!task) return

  formData.value.successors.push({
    id: task.id,
    name: task.name,
    type: dependencyForm.value.type,
    lag: dependencyForm.value.lag
  })

  successorDialogVisible.value = false
}

// 删除紧前任务
const removePredecessor = (index) => {
  formData.value.predecessors.splice(index, 1)
}

// 删除紧后任务
const removeSuccessor = (index) => {
  formData.value.successors.splice(index, 1)
}

// ==================== 表单操作 ====================
// 监听 editingTask 变化以填充表单
watch(
  () => props.editingTask,
  (newTask) => {
    if (newTask) {
      formData.value = {
        name: newTask.name || '',
        start: newTask.start || '',
        end: newTask.end || '',
        progress: newTask.progress || 0,
        priority: newTask.priority || 'medium',
        notes: newTask.notes || '',
        resources: newTask.resources || [],
        predecessors: (newTask.predecessors || []).map(p => {
          if (typeof p === 'object') return p
          // 如果只是 ID，需要查找任务信息
          const task = props.allTasks.find(t => t.id === p)
          return task ? { id: task.id, name: task.name, type: 'FS', lag: 0 } : null
        }).filter(Boolean),
        successors: (newTask.successors || []).map(s => {
          if (typeof s === 'object') return s
          // 如果只是 ID，需要查找任务信息
          const task = props.allTasks.find(t => t.id === s)
          return task ? { id: task.id, name: task.name, type: 'FS', lag: 0 } : null
        }).filter(Boolean)
      }
    } else {
      formData.value = { ...defaultFormData }
    }
  },
  { immediate: true }
)

// 对话框打开/关闭时重置表单
watch(
  () => props.visible,
  (visible) => {
    if (visible && !props.editingTask) {
      formData.value = { ...defaultFormData }
      activeTab.value = 'basic'
    }
  }
)

const handleSave = async () => {
  if (!formRef.value) return

  try {
    console.log('TaskEditDialog - 表单验证前的数据:', JSON.parse(JSON.stringify(formData.value)))
    console.log('TaskEditDialog - name字段值:', formData.value.name)
    console.log('TaskEditDialog - name字段类型:', typeof formData.value.name)

    const valid = await formRef.value.validate()
    console.log('TaskEditDialog - 验证结果:', valid)

    if (!valid) {
      console.error('TaskEditDialog - 表单验证失败')
      // 获取验证错误信息
      const fields = formRef.value.fields || {}
      for (const key in fields) {
        const field = fields[key]
        if (field && field.validateState === 'error') {
          console.error(`字段 ${key} 验证失败:`, field.validateMessage)
        }
      }
      return
    }

    console.log('TaskEditDialog - 表单验证通过')

    // 构建保存数据
    const saveData = {
      ...formData.value,
      predecessor_ids: formData.value.predecessors.map(p => ({
        task_id: p.id,
        type: p.type,
        lag: p.lag
      })),
      successor_ids: formData.value.successors.map(s => ({
        task_id: s.id,
        type: s.type,
        lag: s.lag
      }))
    }

    console.log('TaskEditDialog - 要保存的数据:', saveData)
    emit('save', saveData)
  } catch (error) {
    console.error('TaskEditDialog - 保存过程出错:', error)
    // 验证失败或其他错误
  }
}

// 关闭对话框
const handleClose = () => {
  // 重置表单状态
  resetForm()
  // 发送关闭事件
  emit('update:visible', false)
}

const resetForm = () => {
  formData.value = { ...defaultFormData }
  formRef.value?.clearValidate()
  activeTab.value = 'basic'
}

// ==================== ESC 键关闭对话框 ====================
const handleKeydown = (event) => {
  // ESC 键关闭对话框
  if (event.key === 'Escape' && props.visible) {
    // 如果子对话框打开，不关闭主对话框
    if (!resourceDialogVisible.value &&
        !predecessorDialogVisible.value &&
        !successorDialogVisible.value) {
      handleClose()
    }
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
})

defineExpose({
  resetForm
})
</script>

<style scoped>
/* 资源部分 */
.resources-section {
  padding: 10px;
}

.resources-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid #ebeef5;
}

.resources-title {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
}

/* 依赖关系部分 */
.dependencies-section {
  padding: 10px;
}

.dep-group {
  margin-bottom: 24px;
}

.dep-group:last-child {
  margin-bottom: 0;
}

.dep-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #ebeef5;
}

.dep-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  font-weight: 500;
  color: #303133;
}

.dep-title .el-icon {
  color: #409eff;
}

/* 表格样式优化 */
:deep(.el-table) {
  font-size: 13px;
}

:deep(.el-table th) {
  background-color: #fafafa;
}

:deep(.el-empty) {
  padding: 20px 0;
}
</style>
