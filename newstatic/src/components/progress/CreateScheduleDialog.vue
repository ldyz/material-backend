<template>
  <el-dialog
    v-model="visible"
    title="创建进度计划"
    width="900px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <div class="create-schedule-dialog">
      <!-- 项目信息 -->
      <el-alert
        :title="alertTitle"
        type="info"
        :closable="false"
        show-icon
        style="margin-bottom: 20px"
      />

      <!-- 任务列表 -->
      <div class="task-list-header">
        <div class="header-title">
          <span>任务列表</span>
          <el-tag type="info" size="small">{{ tasks.length }} 个任务</el-tag>
        </div>
        <el-button
          type="primary"
          size="small"
          :icon="Plus"
          @click="handleAddTask"
        >
          添加任务
        </el-button>
      </div>

      <el-table
        :data="tasks"
        border
        stripe
        style="width: 100%; margin-bottom: 20px"
        max-height="400"
      >
        <el-table-column label="序号" type="index" width="60" align="center" />
        <el-table-column label="任务名称" min-width="180">
          <template #default="scope">
            <el-input
              v-model="scope.row.name"
              placeholder="任务名称"
              size="small"
              :disabled="scope.row.loading"
            />
          </template>
        </el-table-column>
        <el-table-column label="开始日期" width="150">
          <template #default="scope">
            <el-date-picker
              v-model="scope.row.start"
              type="date"
              placeholder="开始日期"
              size="small"
              style="width: 100%"
              value-format="YYYY-MM-DD"
              :disabled="scope.row.loading"
            />
          </template>
        </el-table-column>
        <el-table-column label="结束日期" width="150">
          <template #default="scope">
            <el-date-picker
              v-model="scope.row.end"
              type="date"
              placeholder="结束日期"
              size="small"
              style="width: 100%"
              value-format="YYYY-MM-DD"
              :disabled="scope.row.loading"
            />
          </template>
        </el-table-column>
        <el-table-column label="工期(天)" width="100" align="center">
          <template #default="scope">
            <span>{{ calculateDuration(scope.row) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100" align="center">
          <template #default="scope">
            <el-button
              type="danger"
              size="small"
              :icon="Delete"
              circle
              :disabled="scope.row.loading"
              @click="handleRemoveTask(scope.$index)"
            />
          </template>
        </el-table-column>
      </el-table>

      <!-- 快速添加模板任务 -->
      <div class="quick-add-section">
        <el-divider content-position="left">
          <span style="color: #909399; font-size: 13px">快速添加常用任务</span>
        </el-divider>
        <el-space wrap>
          <el-button
            size="small"
            @click="addTemplateTask('planning')"
          >
            项目规划
          </el-button>
          <el-button
            size="small"
            @click="addTemplateTask('design')"
          >
            方案设计
          </el-button>
          <el-button
            size="small"
            @click="addTemplateTask('preparation')"
          >
            施工准备
          </el-button>
          <el-button
            size="small"
            @click="addTemplateTask('execution')"
          >
            施工执行
          </el-button>
          <el-button
            size="small"
            @click="addTemplateTask('inspection')"
          >
            质量检查
          </el-button>
          <el-button
            size="small"
            @click="addTemplateTask('completion')"
          >
            竣工验收
          </el-button>
        </el-space>
      </div>

      <!-- 说明 -->
      <el-alert
        title="提示"
        type="warning"
        :closable="false"
        show-icon
        style="margin-top: 20px"
      >
        <ul style="margin: 5px 0; padding-left: 20px">
          <li>创建进度计划后，您可以在甘特图和网络图中编辑和管理任务</li>
          <li>任务的实际进度会自动同步到项目列表中显示</li>
          <li>您可以随时调整任务的开始和结束日期，系统会自动重新计算项目进度</li>
        </ul>
      </el-alert>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose" :disabled="saving">取消</el-button>
        <el-button
          type="primary"
          @click="handleSave"
          :loading="saving"
          :disabled="tasks.length === 0"
        >
          创建计划
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, Delete } from '@element-plus/icons-vue'
import { progressApi } from '@/api'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  projectId: {
    type: [Number, String],
    default: null
  },
  projectName: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue', 'created'])

const visible = ref(false)
const saving = ref(false)
const tasks = ref([])

// 计算属性：警告标题
const alertTitle = computed(() => {
  return `为项目 "${props.projectName}" 创建进度计划`
})

// 监听 modelValue 变化
watch(() => props.modelValue, (val) => {
  visible.value = val
  if (val && props.projectId) {
    // 如果还没有任务，添加一个默认任务
    if (tasks.value.length === 0) {
      addDefaultTasks()
    }
  }
})

// 监听 visible 变化
watch(visible, (val) => {
  emit('update:modelValue', val)
  if (!val) {
    // 对话框关闭时重置
    tasks.value = []
  }
})

// 添加默认任务
const addDefaultTasks = () => {
  const today = new Date()
  const formatDate = (date) => {
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    return `${year}-${month}-${day}`
  }

  tasks.value = [
    {
      name: '项目启动',
      start: formatDate(today),
      end: formatDate(new Date(today.getTime() + 7 * 24 * 60 * 60 * 1000)),
      loading: false
    }
  ]
}

// 添加任务
const handleAddTask = () => {
  const lastTask = tasks.value[tasks.value.length - 1]
  let startDate = new Date()

  if (lastTask && lastTask.end) {
    startDate = new Date(lastTask.end)
    startDate.setDate(startDate.getDate() + 1)
  }

  const formatDate = (date) => {
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    return `${year}-${month}-${day}`
  }

  const endDate = new Date(startDate)
  endDate.setDate(endDate.getDate() + 7)

  tasks.value.push({
    name: '',
    start: formatDate(startDate),
    end: formatDate(endDate),
    loading: false
  })
}

// 删除任务
const handleRemoveTask = (index) => {
  tasks.value.splice(index, 1)
}

// 计算工期
const calculateDuration = (task) => {
  if (!task.start || !task.end) return 0
  const start = new Date(task.start)
  const end = new Date(task.end)
  const diff = Math.ceil((end - start) / (1000 * 60 * 60 * 24))
  return diff > 0 ? diff : 0
}

// 添加模板任务
const addTemplateTask = (type) => {
  const templates = {
    planning: { name: '项目规划', duration: 7 },
    design: { name: '方案设计', duration: 14 },
    preparation: { name: '施工准备', duration: 10 },
    execution: { name: '施工执行', duration: 30 },
    inspection: { name: '质量检查', duration: 7 },
    completion: { name: '竣工验收', duration: 5 }
  }

  const template = templates[type]
  if (!template) return

  const lastTask = tasks.value[tasks.value.length - 1]
  let startDate = new Date()

  if (lastTask && lastTask.end) {
    startDate = new Date(lastTask.end)
    startDate.setDate(startDate.getDate() + 1)
  }

  const formatDate = (date) => {
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    return `${year}-${month}-${day}`
  }

  const endDate = new Date(startDate)
  endDate.setDate(endDate.getDate() + template.duration)

  tasks.value.push({
    name: template.name,
    start: formatDate(startDate),
    end: formatDate(endDate),
    loading: false
  })
}

// 保存计划
const handleSave = async () => {
  console.log('CreateScheduleDialog - 开始验证，任务数量:', tasks.value.length)

  // 验证任务名称
  const emptyTasks = []
  const validTasks = []

  tasks.value.forEach((task, index) => {
    console.log(`任务 ${index}: name="${task.name}", trim="${task.name?.trim()}"`)
    if (!task.name || !task.name.trim()) {
      emptyTasks.push(index + 1)
    } else {
      validTasks.push(task)
    }
  })

  console.log('空任务索引:', emptyTasks)
  console.log('有效任务数量:', validTasks.length)

  if (validTasks.length === 0) {
    ElMessage.warning('请至少添加一个任务，并填写任务名称')
    return
  }

  if (emptyTasks.length > 0) {
    ElMessage.warning(`第 ${emptyTasks.join(', ')} 行的任务名称不能为空，请填写后再保存`)
    return
  }

  // 验证日期
  for (const task of validTasks) {
    if (!task.start || !task.end) {
      ElMessage.warning(`任务"${task.name}"的开始日期和结束日期不能为空`)
      return
    }
    if (new Date(task.start) > new Date(task.end)) {
      ElMessage.warning(`任务"${task.name}"的开始日期不能晚于结束日期`)
      return
    }
  }

  try {
    saving.value = true
    console.log('开始批量创建任务，数量:', validTasks.length)

    // 批量创建任务
    let successCount = 0
    for (const task of validTasks) {
      const taskData = {
        project_id: props.projectId,
        name: task.name,
        start_date: task.start,
        end_date: task.end,
        progress: 0,
        priority: 'medium',
        status: 'not_started'
      }

      console.log('创建任务:', taskData)
      await progressApi.create(taskData)
      successCount++
    }

    console.log(`成功创建 ${successCount} 个任务`)
    ElMessage.success(`成功创建 ${successCount} 个任务`)
    emit('created', props.projectId)
    handleClose()
  } catch (error) {
    console.error('创建计划失败:', error)
    console.error('错误详情:', error.response?.data)

    const errorMsg = error.response?.data?.error || error.response?.data?.message || error.message || '创建计划失败，请重试'
    ElMessage.error(`创建失败: ${errorMsg}`)
  } finally {
    saving.value = false
  }
}

// 关闭对话框
const handleClose = () => {
  visible.value = false
  tasks.value = []
}
</script>

<style scoped>
.create-schedule-dialog {
  padding: 10px 0;
}

.task-list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.header-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 14px;
  font-weight: 500;
  color: #303133;
}

.quick-add-section {
  margin: 20px 0;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

:deep(.el-table) {
  font-size: 13px;
}

:deep(.el-table th) {
  background-color: #f5f7fa;
}

:deep(.el-input__wrapper) {
  padding: 4px 8px;
}

:deep(.el-alert__content) {
  padding: 0 10px;
}
</style>
