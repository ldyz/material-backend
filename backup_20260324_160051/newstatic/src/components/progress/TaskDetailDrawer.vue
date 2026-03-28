<template>
  <el-drawer
    :model-value="visible"
    @update:model-value="$emit('update:visible', $event)"
    title="任务详情"
    direction="rtl"
    size="400px"
    :close-on-click-modal="false"
  >
    <div v-if="task" class="task-detail">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="任务名称" :span="2">
          {{ task.name }}
        </el-descriptions-item>
        <el-descriptions-item label="开始日期">
          {{ task.start }}
        </el-descriptions-item>
        <el-descriptions-item label="结束日期">
          {{ task.end }}
        </el-descriptions-item>
        <el-descriptions-item label="工期">
          {{ getTaskDuration(task) }} 天
        </el-descriptions-item>
        <el-descriptions-item label="进度">
          <el-progress
            :percentage="task.progress || 0"
            :status="getProgressStatus(task.progress)"
            :stroke-width="8"
          />
        </el-descriptions-item>
        <el-descriptions-item label="状态" :span="2">
          <el-tag :type="getStatusType(task)" size="small">
            {{ getStatusText(task) }}
          </el-tag>
          <el-tag v-if="task.is_critical" type="danger" size="small" style="margin-left: 8px">
            关键路径
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="优先级">
          <el-tag :type="getPriorityType(task)" size="small">
            {{ getPriorityText(task) }}
          </el-tag>
        </el-descriptions-item>
      </el-descriptions>

      <el-divider />

      <div class="task-actions">
        <el-button type="primary" size="small" @click="$emit('edit')">
          编辑任务
        </el-button>
        <el-button size="small" @click="$emit('duplicate')">
          复制任务
        </el-button>
        <el-button type="danger" size="small" @click="$emit('delete')">
          删除任务
        </el-button>
      </div>
    </div>
  </el-drawer>
</template>

<script setup>
import { diffDays } from '@/utils/dateFormat'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  task: {
    type: Object,
    default: null
  }
})

defineEmits(['update:visible', 'edit', 'duplicate', 'delete'])

const getTaskDuration = (task) => {
  try {
    if (!task.start || !task.end) return 0
    return diffDays(task.start, task.end)
  } catch (e) {
    console.error('Error calculating task duration:', e)
    return 0
  }
}

const getProgressStatus = (progress) => {
  if (progress >= 100) return 'success'
  if (progress >= 80) return 'warning'
  return ''
}

const getStatusType = (task) => {
  const types = {
    completed: 'success',
    delayed: 'danger',
    in_progress: 'warning',
    not_started: 'info'
  }
  return types[task.status] || 'info'
}

const getStatusText = (task) => {
  const texts = {
    completed: '已完成',
    delayed: '已延期',
    in_progress: '进行中',
    not_started: '未开始'
  }
  return texts[task.status] || task.status
}

const getPriorityType = (task) => {
  const types = {
    urgent: 'danger',
    high: 'warning',
    medium: '',
    low: 'info'
  }
  return types[task.priority] || ''
}

const getPriorityText = (task) => {
  const texts = {
    urgent: '紧急',
    high: '高',
    medium: '中',
    low: '低'
  }
  return texts[task.priority] || task.priority
}
</script>

<style scoped>
.task-detail {
  padding: 0 20px;
}

.task-actions {
  display: flex;
  gap: 8px;
  margin-top: 20px;
}
</style>
