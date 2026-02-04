<template>
  <div class="workflow-status">
    <!-- 工作流状态步骤条 -->
    <el-steps
      :active="currentStep"
      finish-status="success"
      align-center
      process-status="process"
    >
      <el-step
        v-for="(step, index) in workflowSteps"
        :key="index"
        :title="step.title"
        :description="step.description"
        :status="getStepStatus(step.status)"
      >
        <template #icon>
          <el-icon>
            <Document v-if="step.status === 'draft'" />
            <Loading v-else-if="step.status === 'pending'" />
            <CircleCheck v-else-if="step.status === 'approved'" />
            <CircleClose v-else-if="step.status === 'rejected'" />
            <Box v-else-if="step.status === 'issued'" />
            <Finished v-else-if="step.status === 'completed'" />
            <Warning v-else />
          </el-icon>
        </template>
      </el-step>
    </el-steps>

    <!-- 当前状态信息 -->
    <el-alert
      v-if="currentStatus"
      :type="getAlertType(currentStatus)"
      :closable="false"
      style="margin-top: 20px"
    >
      <template #title>
        <div class="status-title">
          <span>当前状态：{{ getStatusText(currentStatus) }}</span>
          <span class="status-time" v-if="statusTime">
            （{{ formatTime(statusTime) }}）
          </span>
        </div>
      </template>
      <div v-if="statusDescription" class="status-description">
        {{ statusDescription }}
      </div>
    </el-alert>

    <!-- 操作按钮组 -->
    <div class="workflow-actions" v-if="actions && actions.length > 0">
      <el-button
        v-for="(action, index) in actions"
        :key="index"
        :type="action.type || 'default'"
        :icon="action.icon"
        :disabled="action.disabled"
        @click="action.handler"
        size="small"
      >
        {{ action.label }}
      </el-button>
    </div>
  </div>
</template>

<script setup>
import { computed, reactive } from 'vue'
import {
  Document,
  Loading,
  CircleCheck,
  CircleClose,
  Box,
  Finished,
  Warning
} from '@element-plus/icons-vue'

const props = defineProps({
  // 当前状态
  status: {
    type: String,
    required: true
  },
  // 状态时间
  statusTime: {
    type: String,
    default: ''
  },
  // 状态描述
  statusDescription: {
    type: String,
    default: ''
  },
  // 工作流类型
  workflowType: {
    type: String,
    default: 'requisition' // requisition | inbound
  }
})

const emit = defineEmits(['action'])

// 当前步骤索引
const currentStep = computed(() => {
  const stepMap = {
    draft: 0,
    pending: 1,
    approved: 2,
    rejected: 1,
    issued: 3,
    completed: 2
  }
  return stepMap[props.status] || 0
})

// 工作流步骤定义
const workflowSteps = computed(() => {
  if (props.workflowType === 'requisition') {
    return [
      {
        status: 'draft',
        title: '草稿',
        description: '创建出库单草稿'
      },
      {
        status: 'pending',
        title: '待审核',
        description: '等待审核人员审核'
      },
      {
        status: 'approved',
        title: '已审核',
        description: '审核通过，准备发货'
      },
      {
        status: 'issued',
        title: '已发货',
        description: '物资已出库发货'
      }
    ]
  } else if (props.workflowType === 'inbound') {
    return [
      {
        status: 'draft',
        title: '草稿',
        description: '创建入库单草稿'
      },
      {
        status: 'pending',
        title: '待审核',
        description: '等待审核人员审核'
      },
      {
        status: 'completed',
        title: '已完成',
        description: '入库已完成，库存已更新'
      }
    ]
  }
  return []
})

// 可用操作
const actions = computed(() => {
  const actionList = []

  if (props.status === 'pending') {
    actionList.push({
      label: '审核通过',
      type: 'success',
      icon: 'CircleCheck',
      handler: () => emit('action', 'approve')
    })
    actionList.push({
      label: '审核拒绝',
      type: 'danger',
      icon: 'CircleClose',
      handler: () => emit('action', 'reject')
    })
  } else if (props.status === 'approved' && props.workflowType === 'requisition') {
    actionList.push({
      label: '发货',
      type: 'warning',
      icon: 'Box',
      handler: () => emit('action', 'issue')
    })
  }

  return actionList
})

// 获取步骤状态
const getStepStatus = (stepStatus) => {
  if (stepStatus === 'rejected') return 'error'
  if (stepStatus === props.status) return 'process'
  return 'wait'
}

// 获取告警类型
const getAlertType = (status) => {
  const types = {
    draft: 'info',
    pending: 'warning',
    approved: 'success',
    rejected: 'error',
    issued: 'info',
    completed: 'success'
  }
  return types[status] || 'info'
}

// 获取状态文本
const getStatusText = (status) => {
  const texts = {
    draft: '草稿',
    pending: '待审核',
    approved: '已审核',
    rejected: '已拒绝',
    issued: '已发货',
    completed: '已完成'
  }
  return texts[status] || status
}

// 格式化时间
const formatTime = (time) => {
  if (!time) return ''
  const date = new Date(time)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}
</script>

<style scoped>
.workflow-status {
  padding: 10px 0;
  /* 限制为2行高度：标题1行 + 描述1行 */
  max-height: 80px;
  overflow: hidden;
}

.status-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 14px;
  line-height: 1.5;
}

.status-time {
  font-size: 12px;
  color: #909399;
  font-weight: normal;
}

.status-description {
  margin-top: 4px;
  font-size: 13px;
  color: #606266;
  line-height: 1.5;
  /* 限制为1行 */
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.workflow-actions {
  margin-top: 10px;
  display: flex;
  gap: 10px;
  justify-content: center;
}

/* 隐藏步骤条的描述，只保留标题，节省空间 */
:deep(.el-step__title) {
  font-size: 13px;
  font-weight: 500;
  line-height: 1.5;
}

:deep(.el-step__description) {
  display: none;
}

:deep(.el-steps) {
  max-height: 40px;
}

:deep(.el-step) {
  flex-basis: 30% !important;
}
</style>
