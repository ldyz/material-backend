/**
 * 工作流审批操作面板组件
 *
 * 功能说明：
 * - 显示单据的当前状态和工作流进度
 * - 提供审批操作按钮（通过、拒绝、退回、转办等）
 * - 显示可执行的操作列表
 * - 显示审批历史记录
 *
 * @module WorkflowActionPanel
 * @author Material Management System
 * @date 2025-01-27
 */
<template>
  <div class="workflow-action-panel">
    <!-- 当前状态卡片 -->
    <el-card shadow="never" class="status-card">
      <template #header>
        <div class="card-header">
          <span class="title">当前状态</span>
          <el-tag :type="getStatusTagType(status)" size="large">
            {{ getStatusText(status) }}
          </el-tag>
        </div>
      </template>

      <!-- 工作流进度条 -->
      <div v-if="workflowSteps.length > 0" class="workflow-progress">
        <el-steps :active="currentStepIndex" finish-status="success" align-center>
          <el-step
            v-for="(step, index) in workflowSteps"
            :key="index"
            :title="step.name"
            :description="step.description"
          />
        </el-steps>
      </div>

      <!-- 状态说明 -->
      <div v-if="statusDescription" class="status-description">
        <el-alert
          :title="statusDescription"
          :type="getStatusAlertType(status)"
          :closable="false"
          show-icon
        />
      </div>

      <!-- 审批人信息（多人审批时显示） -->
      <div v-if="approverInfo && approverInfo.length > 0" class="approver-info">
        <div class="approver-title">
          <el-icon><User /></el-icon>
          <span>审批人信息</span>
        </div>
        <div class="approver-list">
          <div
            v-for="(approver, index) in approverInfo"
            :key="index"
            class="approver-item"
            :class="{
              'approved': approver.status === 'approved',
              'pending': approver.status === 'pending'
            }"
          >
            <el-avatar :size="32" :icon="UserFilled" />
            <div class="approver-detail">
              <span class="approver-name">{{ approver.name }}</span>
              <span class="approver-status">{{ getApproverStatusText(approver.status) }}</span>
            </div>
            <el-icon v-if="approver.status === 'approved'" class="approver-check" color="#67c23a">
              <CircleCheck />
            </el-icon>
            <el-icon v-else class="approver-waiting" color="#909399">
              <Clock />
            </el-icon>
          </div>
        </div>
      </div>

      <!-- 可执行操作 -->
      <div v-if="availableActions.length > 0" class="action-buttons">
        <div class="action-title">
          <el-icon><Operation /></el-icon>
          <span>可执行操作</span>
        </div>
        <div class="action-list">
          <el-button
            v-for="action in availableActions"
            :key="action.key"
            :type="action.type || 'primary'"
            :icon="action.icon"
            :loading="actionLoading"
            @click="handleAction(action.key)"
          >
            {{ action.label }}
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- 审批历史 -->
    <el-card v-if="histories.length > 0" shadow="never" class="history-card">
      <template #header>
        <div class="card-header">
          <span class="title">
            <el-icon><Clock /></el-icon>
            审批历史
          </span>
        </div>
      </template>

      <el-timeline>
        <el-timeline-item
          v-for="(history, index) in histories"
          :key="index"
          :timestamp="history.timestamp"
          :type="getHistoryType(history.action)"
          placement="top"
        >
          <div class="history-item">
            <div class="history-header">
              <span class="action-name">{{ history.action_name }}</span>
              <span class="actor">{{ history.actor_name }}</span>
            </div>
            <div v-if="history.comment" class="history-comment">
              {{ history.comment }}
            </div>
            <div v-if="history.attachments && history.attachments.length > 0" class="history-attachments">
              <el-tag
                v-for="(file, idx) in history.attachments"
                :key="idx"
                size="small"
                class="attachment-tag"
              >
                {{ file.name }}
              </el-tag>
            </div>
          </div>
        </el-timeline-item>
      </el-timeline>
    </el-card>

    <!-- 审批操作对话框 -->
    <el-dialog
      v-model="actionDialogVisible"
      :title="actionDialogTitle"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="actionFormRef"
        :model="actionForm"
        :rules="actionFormRules"
        label-width="100px"
      >
        <!-- 根据操作类型显示不同的表单项 -->
        <template v-if="currentAction === 'approve'">
          <el-alert
            title="审批通过"
            type="success"
            :closable="false"
            style="margin-bottom: 20px"
          >
            确认要通过此单据吗？通过后单据将流转到下一环节
          </el-alert>
        </template>

        <template v-else-if="currentAction === 'reject'">
          <el-alert
            title="审批拒绝"
            type="error"
            :closable="false"
            style="margin-bottom: 20px"
          >
            拒绝后单据将被终止，无法继续流转
          </el-alert>

          <el-form-item label="拒绝原因" prop="reason" required>
            <el-input
              v-model="actionForm.reason"
              type="textarea"
              :rows="4"
              placeholder="请输入拒绝原因"
              maxlength="500"
              show-word-limit
            />
          </el-form-item>
        </template>

        <template v-else-if="currentAction === 'return'">
          <el-alert
            title="退回修改"
            type="warning"
            :closable="false"
            style="margin-bottom: 20px"
          >
            将单据退回到上一节点或指定节点进行修改
          </el-alert>

          <el-form-item label="退回到" prop="targetNodeId" required>
            <el-select
              v-model="actionForm.targetNodeId"
              placeholder="请选择退回节点"
              style="width: 100%"
            >
              <el-option
                v-for="node in returnableNodes"
                :key="node.id"
                :label="node.name"
                :value="node.id"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="退回原因" prop="comment" required>
            <el-input
              v-model="actionForm.comment"
              type="textarea"
              :rows="4"
              placeholder="请说明退回原因"
              maxlength="500"
              show-word-limit
            />
          </el-form-item>
        </template>

        <template v-else-if="currentAction === 'comment'">
          <el-form-item label="意见说明" prop="comment" required>
            <el-input
              v-model="actionForm.comment"
              type="textarea"
              :rows="6"
              placeholder="请输入您的意见或说明"
              maxlength="1000"
              show-word-limit
            />
          </el-form-item>

          <el-form-item label="附件">
            <el-upload
              ref="uploadRef"
              :auto-upload="false"
              :limit="3"
              :file-list="actionForm.attachments"
              :on-change="handleFileChange"
              action="#"
            >
              <el-button :icon="Upload" size="small">选择文件</el-button>
              <template #tip>
                <div class="el-upload__tip">
                  最多上传3个文件，每个文件不超过10MB
                </div>
              </template>
            </el-upload>
          </el-form-item>
        </template>

        <template v-else>
          <el-form-item label="操作说明" prop="comment">
            <el-input
              v-model="actionForm.comment"
              type="textarea"
              :rows="4"
              placeholder="请输入操作说明（可选）"
              maxlength="500"
            />
          </el-form-item>
        </template>
      </el-form>

      <template #footer>
        <el-button @click="actionDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="actionSubmitting" @click="handleSubmitAction">
          确认
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Operation, Clock, Upload, Check, Close, Back, ChatLineSquare, User, UserFilled, CircleCheck } from '@element-plus/icons-vue'

/**
 * 组件 Props
 */
const props = defineProps({
  /**
   * 单据ID
   */
  orderId: {
    type: [Number, String],
    required: true
  },
  /**
   * 单据类型（inbound/requisition）
   */
  orderType: {
    type: String,
    required: true,
    validator: (value) => ['inbound', 'requisition'].includes(value)
  },
  /**
   * 当前状态
   */
  status: {
    type: String,
    default: 'draft'
  },
  /**
   * 状态说明
   */
  statusDescription: {
    type: String,
    default: ''
  },
  /**
   * 工作流步骤
   */
  workflowSteps: {
    type: Array,
    default: () => []
  },
  /**
   * 当前步骤索引
   */
  currentStepIndex: {
    type: Number,
    default: 0
  },
  /**
   * 可执行的操作
   */
  availableActions: {
    type: Array,
    default: () => []
  },
  /**
   * 审批历史
   */
  histories: {
    type: Array,
    default: () => []
  },
  /**
   * 审批人信息（多人审批时使用）
   */
  approverInfo: {
    type: Array,
    default: () => []
  }
})

/**
 * 组件 Emits
 */
const emit = defineEmits(['action'])

// ========== 状态管理 ==========

/**
 * 操作加载状态
 */
const actionLoading = ref(false)

/**
 * 操作对话框显示状态
 */
const actionDialogVisible = ref(false)

/**
 * 当前执行的操作
 */
const currentAction = ref('')

/**
 * 操作对话框标题
 */
const actionDialogTitle = computed(() => {
  const titles = {
    approve: '审批通过',
    reject: '审批拒绝',
    return: '退回修改',
    comment: '添加意见',
    transfer: '转办',
    delegate: '委派'
  }
  return titles[currentAction.value] || '执行操作'
})

/**
 * 操作表单
 */
const actionForm = ref({
  reason: '',
  comment: '',
  targetNodeId: null,
  attachments: []
})

/**
 * 操作表单引用
 */
const actionFormRef = ref(null)

/**
 * 操作提交中状态
 */
const actionSubmitting = ref(false)

/**
 * 可退回的节点列表
 */
const returnableNodes = ref([])

/**
 * 上传组件引用
 */
const uploadRef = ref(null)

// ========== 表单验证规则 ==========

/**
 * 操作表单验证规则
 */
const actionFormRules = computed(() => {
  const rules = {
    comment: []
  }

  if (currentAction.value === 'reject') {
    rules.reason = [
      { required: true, message: '请输入拒绝原因', trigger: 'blur' },
      { min: 5, message: '拒绝原因至少5个字符', trigger: 'blur' }
    ]
  }

  if (currentAction.value === 'return') {
    rules.targetNodeId = [
      { required: true, message: '请选择退回节点', trigger: 'change' }
    ]
    rules.comment = [
      { required: true, message: '请输入退回原因', trigger: 'blur' },
      { min: 5, message: '退回原因至少5个字符', trigger: 'blur' }
    ]
  }

  if (currentAction.value === 'comment') {
    rules.comment = [
      { required: true, message: '请输入意见说明', trigger: 'blur' },
      { min: 2, message: '意见说明至少2个字符', trigger: 'blur' }
    ]
  }

  return rules
})

// ========== 方法定义 ==========

/**
 * 获取状态标签类型
 */
const getStatusTagType = (status) => {
  const types = {
    draft: 'info',
    pending: 'warning',
    approved: 'success',
    completed: 'success',
    rejected: 'danger',
    cancelled: 'info'
  }
  return types[status] || 'info'
}

/**
 * 获取状态文本
 */
const getStatusText = (status) => {
  const texts = {
    draft: '草稿',
    pending: '待审批',
    approved: '已审批',
    completed: '已完成',
    rejected: '已拒绝',
    cancelled: '已取消'
  }
  return texts[status] || status
}

/**
 * 获取状态提示类型
 */
const getStatusAlertType = (status) => {
  const types = {
    draft: 'info',
    pending: 'warning',
    approved: 'success',
    completed: 'success',
    rejected: 'error',
    cancelled: 'info'
  }
  return types[status] || 'info'
}

/**
 * 获取历史记录类型
 */
const getHistoryType = (action) => {
  const types = {
    approve: 'success',
    reject: 'danger',
    return: 'warning',
    comment: 'primary',
    submit: 'info'
  }
  return types[action] || 'primary'
}

/**
 * 获取审批人状态文本
 */
const getApproverStatusText = (status) => {
  const texts = {
    approved: '已审批',
    pending: '待审批',
    rejected: '已拒绝',
    cancelled: '已取消'
  }
  return texts[status] || status
}

/**
 * 处理操作按钮点击
 */
const handleAction = async (actionKey) => {
  currentAction.value = actionKey

  // 重置表单
  actionForm.value = {
    reason: '',
    comment: '',
    targetNodeId: null,
    attachments: []
  }

  // 如果需要退回节点，先获取可退回的节点
  if (actionKey === 'return') {
    await fetchReturnableNodes()
  }

  // 显示操作对话框
  actionDialogVisible.value = true
}

/**
 * 获取可退回的节点列表
 */
const fetchReturnableNodes = async () => {
  try {
    // TODO: 调用API获取可退回的节点
    // 这里先模拟数据
    returnableNodes.value = [
      { id: 1, name: '提交人' },
      { id: 2, name: '部门主管' }
    ]
  } catch (error) {
    console.error('获取可退回节点失败:', error)
  }
}

/**
 * 文件选择变化
 */
const handleFileChange = (file) => {
  actionForm.value.attachments = [file]
}

/**
 * 提交操作
 */
const handleSubmitAction = async () => {
  if (!actionFormRef.value) return

  try {
    await actionFormRef.value.validate()

    actionSubmitting.value = true

    // 准备提交数据
    const data = {
      action: currentAction.value,
      comment: actionForm.value.comment,
      reason: actionForm.value.reason,
      target_node_id: actionForm.value.targetNodeId,
      attachments: actionForm.value.attachments
    }

    // 触发操作事件
    emit('action', data)

    // 关闭对话框
    actionDialogVisible.value = false
    ElMessage.success('操作成功')
  } catch (error) {
    console.error('提交操作失败:', error)
  } finally {
    actionSubmitting.value = false
  }
}

/**
 * 重置操作表单
 */
const resetActionForm = () => {
  if (actionFormRef.value) {
    actionFormRef.value.clearValidate()
  }
  actionForm.value = {
    reason: '',
    comment: '',
    targetNodeId: null,
    attachments: []
  }
}

// 监听对话框关闭
watch(actionDialogVisible, (val) => {
  if (!val) {
    resetActionForm()
  }
})
</script>

<style scoped>
.workflow-action-panel {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.status-card,
.history-card {
  border-radius: 8px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header .title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.card-header .title .el-icon {
  margin-right: 8px;
}

.workflow-progress {
  margin: 20px 0;
}

.status-description {
  margin-top: 20px;
}

.approver-info {
  margin-top: 20px;
  padding: 16px;
  background-color: #f5f7fa;
  border-radius: 8px;
}

.approver-title {
  display: flex;
  align-items: center;
  font-size: 14px;
  font-weight: 600;
  color: #606266;
  margin-bottom: 12px;
}

.approver-title .el-icon {
  margin-right: 8px;
  color: #409eff;
}

.approver-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.approver-item {
  display: flex;
  align-items: center;
  padding: 10px 12px;
  background-color: #fff;
  border-radius: 6px;
  transition: all 0.3s;
}

.approver-item.pending {
  border-left: 3px solid #e6a23c;
}

.approver-item.approved {
  border-left: 3px solid #67c23a;
  background-color: #f0f9ff;
}

.approver-detail {
  flex: 1;
  margin-left: 10px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.approver-name {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
}

.approver-status {
  font-size: 12px;
  color: #909399;
}

.approver-item.approved .approver-status {
  color: #67c23a;
}

.approver-check {
  font-size: 20px;
}

.approver-waiting {
  font-size: 18px;
}

.action-buttons {
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px solid #ebeef5;
}

.action-title {
  display: flex;
  align-items: center;
  font-size: 14px;
  font-weight: 600;
  color: #606266;
  margin-bottom: 12px;
}

.action-title .el-icon {
  margin-right: 8px;
  color: #409eff;
}

.action-list {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.action-list .el-button {
  flex-shrink: 0;
}

.history-item {
  padding-bottom: 8px;
}

.history-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.action-name {
  font-weight: 600;
  color: #303133;
}

.actor {
  font-size: 13px;
  color: #909399;
}

.history-comment {
  margin-top: 8px;
  padding: 8px 12px;
  background-color: #f5f7fa;
  border-radius: 4px;
  font-size: 14px;
  color: #606266;
  line-height: 1.6;
}

.history-attachments {
  margin-top: 8px;
}

.attachment-tag {
  margin-right: 8px;
  margin-bottom: 4px;
}
</style>
