<template>
  <div class="workflow-history">
    <div class="steps-container" v-if="histories && histories.length > 0">
      <div class="steps-wrapper">
        <div
          v-for="(item, index) in histories"
          :key="index"
          class="step-item"
          :class="{ 'is-last': index === histories.length - 1 }"
        >
          <div class="step-line" v-if="!isLastItem(index)"></div>
          <div class="step-dot" :class="getStatusClass(item.action)">
            <el-icon>
              <CircleCheck v-if="item.action === 'approved'" />
              <CircleClose v-else-if="item.action === 'rejected'" />
              <Loading v-else-if="item.action === 'pending'" />
              <Edit v-else-if="item.action === 'draft'" />
              <Finished v-else-if="item.action === 'completed'" />
              <Box v-else-if="item.action === 'issued'" />
              <Document v-else />
            </el-icon>
          </div>
          <div class="step-content">
            <div class="step-title">{{ getActionText(item.action) }}</div>
            <div class="step-time">{{ formatTime(item.created_at) }}</div>
            <div class="step-remark" v-if="item.remark">{{ item.remark }}</div>
            <div class="step-description" v-else-if="item.description">{{ item.description }}</div>
            <div class="step-operator">
              <el-icon><User /></el-icon>
              {{ item.operator_name || item.operator || '系统' }}
            </div>
          </div>
        </div>
      </div>
    </div>
    <el-empty
      v-else
      description="暂无审批记录"
      :image-size="80"
    />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import {
  CircleCheck,
  CircleClose,
  Loading,
  Edit,
  Finished,
  Box,
  Document,
  User
} from '@element-plus/icons-vue'

const props = defineProps({
  histories: {
    type: Array,
    default: () => []
  }
})

// 判断是否是最后一项
const isLastItem = (index) => {
  return index === props.histories.length - 1
}

// 获取操作文本
const getActionText = (action) => {
  const texts = {
    draft: '创建草稿',
    submit: '提交审核',
    pending: '待审核',
    approved: '审核通过',
    rejected: '审核拒绝',
    completed: '已完成',
    issued: '已发货',
    cancelled: '已取消',
    approve: '审批通过',
    reject: '审批拒绝',
    return: '退回',
    comment: '评论'
  }
  return texts[action] || action
}

// 获取状态样式
const getStatusClass = (action) => {
  const classes = {
    draft: 'draft',
    submit: 'submit',
    pending: 'pending',
    approved: 'approved',
    rejected: 'rejected',
    completed: 'completed',
    issued: 'issued',
    cancelled: 'cancelled',
    approve: 'approved',
    reject: 'rejected',
    return: 'return',
    comment: 'comment'
  }
  return classes[action] || 'default'
}

// 格式化时间
const formatTime = (time) => {
  if (!time) return ''
  const date = new Date(time)
  const now = new Date()
  const diffMs = now - date
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMs / 3600000)
  const diffDays = Math.floor(diffMs / 86400000)

  if (diffMins < 1) return '刚刚'
  if (diffMins < 60) return `${diffMins}分钟前`
  if (diffHours < 24) return `${diffHours}小时前`
  if (diffDays < 7) return `${diffDays}天前`

  // 超过7天显示具体日期
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
.workflow-history {
  padding: 10px 0;
}

.steps-container {
  width: 100%;
  overflow-x: auto;
}

.steps-wrapper {
  display: flex;
  align-items: flex-start;
  gap: 0;
  min-width: max-content;
  padding: 10px 0;
}

.step-item {
  position: relative;
  display: flex;
  align-items: flex-start;
  flex-shrink: 0;
  padding: 0 20px 0 0;
}

.step-item.is-last {
  padding-right: 0;
}

.step-line {
  position: absolute;
  top: 16px;
  left: 60px;
  right: -20px;
  height: 2px;
  background: #e4e7ed;
  z-index: 0;
}

.step-dot {
  position: relative;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: white;
  border: 2px solid;
  z-index: 1;
  font-size: 16px;
  flex-shrink: 0;
}

.step-dot.draft {
  border-color: #909399;
  color: #909399;
  background: #f5f7fa;
}

.step-dot.submit {
  border-color: #409eff;
  color: #409eff;
  background: #ecf5ff;
}

.step-dot.pending {
  border-color: #e6a23c;
  color: #e6a23c;
  background: #fdf6ec;
}

.step-dot.approved {
  border-color: #67c23a;
  color: #67c23a;
  background: #f0f9ff;
}

.step-dot.rejected {
  border-color: #f56c6c;
  color: #f56c6c;
  background: #fef0f0;
}

.step-dot.completed {
  border-color: #67c23a;
  color: #67c23a;
  background: #f0f9ff;
}

.step-dot.issued {
  border-color: #909399;
  color: #909399;
  background: #f5f7fa;
}

.step-dot.return {
  border-color: #e6a23c;
  color: #e6a23c;
  background: #fdf6ec;
}

.step-dot.comment {
  border-color: #409eff;
  color: #409eff;
  background: #ecf5ff;
}

.step-content {
  margin-left: 12px;
  min-width: 150px;
  max-width: 200px;
}

.step-title {
  font-weight: bold;
  font-size: 14px;
  color: #303133;
  margin-bottom: 4px;
}

.step-time {
  font-size: 12px;
  color: #909399;
  margin-bottom: 6px;
}

.step-remark,
.step-description {
  font-size: 12px;
  color: #606266;
  line-height: 1.5;
  margin-bottom: 6px;
  word-break: break-all;
}

.step-operator {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #909399;
}

/* 确保水平滚动条样式美观 */
.steps-container::-webkit-scrollbar {
  height: 6px;
}

.steps-container::-webkit-scrollbar-track {
  background: #f5f7fa;
  border-radius: 3px;
}

.steps-container::-webkit-scrollbar-thumb {
  background: #dcdfe6;
  border-radius: 3px;
}

.steps-container::-webkit-scrollbar-thumb:hover {
  background: #c0c4cc;
}
</style>
