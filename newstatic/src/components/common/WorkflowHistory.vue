<template>
  <div class="workflow-history">
    <div class="timeline" v-if="histories && histories.length > 0">
      <div
        v-for="(item, index) in histories"
        :key="index"
        class="timeline-item"
      >
        <div class="timeline-dot" :class="getStatusClass(item.action)">
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
        <div class="timeline-content">
          <div class="timeline-header">
            <span class="timeline-title">{{ getActionText(item.action) }}</span>
            <span class="timeline-time">{{ formatTime(item.created_at) }}</span>
          </div>
          <div class="timeline-body" v-if="item.remark || item.description">
            <p v-if="item.remark"><strong>审核意见：</strong>{{ item.remark }}</p>
            <p v-if="item.description && !item.remark">{{ item.description }}</p>
          </div>
          <div class="timeline-footer">
            <span class="timeline-operator">
              <el-icon><User /></el-icon>
              {{ item.operator_name || item.operator || '系统' }}
            </span>
            <span class="timeline-department" v-if="item.department">
              <el-icon><OfficeBuilding /></el-icon>
              {{ item.department }}
            </span>
          </div>
        </div>
      </div>
    </div>
    <el-empty
      v-else
      description="暂无工作流记录"
      :image-size="100"
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
  User,
  OfficeBuilding
} from '@element-plus/icons-vue'

const props = defineProps({
  histories: {
    type: Array,
    default: () => []
  }
})

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
    cancelled: '已取消'
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
    cancelled: 'cancelled'
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
  padding: 20px;
}

.timeline {
  position: relative;
  padding-left: 30px;
}

.timeline::before {
  content: '';
  position: absolute;
  left: 15px;
  top: 0;
  bottom: 0;
  width: 2px;
  background: linear-gradient(180deg, #409eff 0%, #c0c4cc 100%);
}

.timeline-item {
  position: relative;
  padding-bottom: 30px;
}

.timeline-item:last-child {
  padding-bottom: 0;
}

.timeline-dot {
  position: absolute;
  left: -29px;
  top: 0;
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
}

.timeline-dot.draft {
  border-color: #909399;
  color: #909399;
}

.timeline-dot.submit {
  border-color: #409eff;
  color: #409eff;
}

.timeline-dot.pending {
  border-color: #e6a23c;
  color: #e6a23c;
}

.timeline-dot.approved {
  border-color: #67c23a;
  color: #67c23a;
}

.timeline-dot.rejected {
  border-color: #f56c6c;
  color: #f56c6c;
}

.timeline-dot.completed {
  border-color: #67c23a;
  color: #67c23a;
}

.timeline-dot.issued {
  border-color: #909399;
  color: #909399;
}

.timeline-content {
  background: #f5f7fa;
  border-radius: 8px;
  padding: 15px;
  margin-left: 20px;
}

.timeline-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.timeline-title {
  font-weight: bold;
  font-size: 14px;
  color: #303133;
}

.timeline-time {
  font-size: 12px;
  color: #909399;
}

.timeline-body {
  margin-bottom: 10px;
}

.timeline-body p {
  margin: 0;
  font-size: 13px;
  color: #606266;
  line-height: 1.6;
}

.timeline-body strong {
  color: #303133;
}

.timeline-footer {
  display: flex;
  gap: 20px;
  font-size: 12px;
  color: #909399;
}

.timeline-operator,
.timeline-department {
  display: flex;
  align-items: center;
  gap: 4px;
}
</style>
