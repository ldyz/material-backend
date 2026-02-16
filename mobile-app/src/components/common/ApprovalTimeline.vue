<template>
  <div class="approval-timeline">
    <van-steps direction="vertical" :active="currentStep">
      <van-step v-for="(node, index) in timelineNodes" :key="index">
        <!-- 自定义节点图标 -->
        <template #icon>
          <div class="step-icon" :class="getStepClass(node)">
            <van-icon v-if="node.status === 'approved'" name="passed" size="14" />
            <van-icon v-else-if="node.status === 'rejected'" name="close" size="14" />
            <van-icon v-else-if="node.status === 'pending'" name="clock-o" size="14" />
            <van-icon v-else name="ellipsis" size="14" />
          </div>
        </template>

        <!-- 节点内容 -->
        <div class="step-content">
          <div class="step-header">
            <span class="step-title">{{ node.title }}</span>
            <van-tag
              v-if="node.status"
              :type="getStatusType(node.status)"
              size="small"
            >
              {{ getStatusText(node.status) }}
            </van-tag>
          </div>

          <!-- 审批人信息 -->
          <div v-if="node.approver" class="step-info">
            <van-icon name="user-o" size="14" />
            <span>{{ node.approver }}</span>
            <span v-if="node.approver_role" class="role-badge">{{ node.approver_role }}</span>
          </div>

          <!-- 审批时间 -->
          <div v-if="node.approved_at" class="step-time">
            <van-icon name="clock-o" size="14" />
            <span>{{ formatDateTime(node.approved_at) }}</span>
          </div>

          <!-- 审批意见 -->
          <div v-if="node.remark" class="step-remark">
            <van-icon name="comment-o" size="14" />
            <span>{{ node.remark }}</span>
          </div>

          <!-- 待审批提示 -->
          <div v-if="node.status === 'pending'" class="step-pending">
            <van-icon name="info-o" size="14" />
            <span>等待{{ node.title }}审批</span>
          </div>
        </div>
      </van-step>
    </van-steps>

    <van-empty v-if="!timelineNodes.length" description="暂无审批流程" />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { formatDateTime } from '@/composables/useDateTime'
import { formatApprovalLog, getNodeName, getActionName, getActionStatus, getStatusText as getWorkflowStatusText, getStatusType as getWorkflowStatusType } from '@/composables/useApprovalWorkflow'

const props = defineProps({
  /**
   * 审批历史记录数组
   */
  approvalLogs: {
    type: Array,
    default: () => []
  },
  /**
   * 审批流程配置（包含所有审批节点）
   * 示例: [
   *   { role: 'project_manager', title: '项目经理审批', order: 1 },
   *   { role: 'supervisor', title: '主管审批', order: 2 }
   * ]
   */
  workflowConfig: {
    type: Array,
    default: () => []
  },
  /**
   * 当前预约单状态
   */
  currentStatus: {
    type: String,
    default: ''
  }
})

/**
 * 构建时间线节点
 * 合并审批流程配置和实际审批记录
 */
const timelineNodes = computed(() => {
  console.log('[ApprovalTimeline] approvalLogs:', props.approvalLogs)
  console.log('[ApprovalTimeline] workflowConfig:', props.workflowConfig)

  // 如果没有工作流配置，直接使用审批日志
  if (!props.workflowConfig || props.workflowConfig.length === 0) {
    const nodes = props.approvalLogs.map(log => {
      const formatted = formatApprovalLog(log)
      return {
        title: formatted.title,
        status: formatted.status,
        approver: formatted.approver,
        approver_role: formatted.approver_role,
        approved_at: formatted.approved_at,
        remark: formatted.remark
      }
    })
    console.log('[ApprovalTimeline] 无配置，节点数:', nodes.length)
    return nodes
  }

  // 有工作流配置，合并配置和实际记录
  const nodes = []

  // 先添加"提交"节点（如果有审批记录）
  if (props.approvalLogs.length > 0) {
    const firstLog = props.approvalLogs[0]
    nodes.push({
      title: '提交申请',
      status: 'approved',
      approver: firstLog.submitter_name || firstLog.actor_name || '申请人',
      approver_role: '申请人',
      approved_at: firstLog.created_at,
      remark: null,
      order: 0
    })
  }

  // 添加每个审批节点
  props.workflowConfig.forEach(config => {
    // 查找该节点的审批记录
    const log = props.approvalLogs.find(l =>
      l.node_key === config.role ||
      l.node_name === config.title ||
      l.node_type === config.role
    )

    // 确定节点状态
    let status = 'pending'
    let approver = null
    let approver_role = config.role_name || config.role
    let approved_at = null
    let remark = null

    if (log) {
      const formatted = formatApprovalLog(log)
      status = formatted.status
      approver = formatted.approver
      approver_role = formatted.approver_role || approver_role
      approved_at = formatted.approved_at
      remark = formatted.remark
    } else {
      // 没有审批记录，检查是否已经被拒绝
      if (props.currentStatus === 'rejected') {
        // 如果被拒绝了，且当前节点还没有审批记录，说明是前面的节点拒绝的
        status = 'cancelled'
      }
    }

    nodes.push({
      title: config.title,
      status,
      approver,
      approver_role,
      approved_at,
      remark,
      order: config.order
    })
  })

  // 按order排序
  const sorted = nodes.sort((a, b) => a.order - b.order)
  console.log('[ApprovalTimeline] 最终节点数:', sorted.length)
  return sorted
})

/**
 * 当前激活的步骤（最后一个已完成的步骤）
 */
const currentStep = computed(() => {
  // 找到最后一个已批准的步骤
  let lastApprovedIndex = -1
  timelineNodes.value.forEach((node, index) => {
    if (node.status === 'approved') {
      lastApprovedIndex = index
    }
  })
  return lastApprovedIndex
})

/**
 * 从日志action获取状态
 */
function getLogStatus(action) {
  const statusMap = {
    submit: 'approved',
    approve: 'approved',
    reject: 'rejected',
    cancel: 'cancelled',
    return: 'returned'
  }
  return statusMap[action] || 'pending'
}

/**
 * 获取节点步骤样式类
 */
function getStepClass(node) {
  const classMap = {
    approved: 'step-success',
    rejected: 'step-error',
    pending: 'step-warning',
    cancelled: 'step-disabled',
    returned: 'step-warning'
  }
  return classMap[node.status] || 'step-default'
}

/**
 * 获取标签类型
 */
function getStatusType(status) {
  return getWorkflowStatusType(status)
}

/**
 * 获取状态文本
 */
function getStatusText(status) {
  return getWorkflowStatusText(status)
}
</script>

<style scoped>
.approval-timeline {
  padding: 16px 0;
}

/* 自定义节点图标 */
.step-icon {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;
}

.step-success {
  background: linear-gradient(135deg, #07c160 0%, #06ad56 100%);
}

.step-error {
  background: linear-gradient(135deg, #ee0a24 0%, #d60a1f 100%);
}

.step-warning {
  background: linear-gradient(135deg, #ff976a 0%, #f3731f 100%);
}

.step-disabled {
  background: #c8c9cc;
}

.step-default {
  background: #e5e5e5;
  color: #969799;
}

/* 节点内容 */
.step-content {
  padding-left: 8px;
  flex: 1;
}

.step-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.step-title {
  font-size: 15px;
  font-weight: 500;
  color: #323233;
}

.step-info,
.step-time,
.step-remark,
.step-pending {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #646566;
  margin-bottom: 6px;
}

.step-info:last-child,
.step-time:last-child,
.step-remark:last-child,
.step-pending:last-child {
  margin-bottom: 0;
}

.role-badge {
  display: inline-block;
  padding: 2px 8px;
  background-color: #f0f0f0;
  border-radius: 4px;
  font-size: 12px;
  color: #666;
  margin-left: 4px;
}

.step-remark {
  padding: 8px;
  background-color: #f7f8fa;
  border-radius: 4px;
  color: #323233;
}

.step-pending {
  color: #ff976a;
}

/* Vant Steps 样式覆盖 */
:deep(.van-steps) {
  background-color: transparent;
}

:deep(.van-step__circle-container) {
  background-color: #f7f8fa;
}

:deep(.van-step__line) {
  background-color: #ebedf0;
}

:deep(.van-step--process .van-step__circle-container) {
  background-color: #f7f8fa;
}
</style>
