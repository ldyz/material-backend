<template>
  <div class="suggestion-card" :class="[`suggestion-${suggestion.type}`, `impact-${suggestion.impact}`]">
    <!-- Card Header -->
    <div class="card-header">
      <div class="header-left">
        <el-icon :color="getIconColor()" class="type-icon">
          <component :is="getTypeIcon()" />
        </el-icon>
        <div class="header-content">
          <h4 class="suggestion-title">{{ suggestion.title }}</h4>
          <span class="suggestion-type">{{ $t(`gantt.ai.suggestions.types.${suggestion.type}`) }}</span>
        </div>
      </div>
      <div class="header-right">
        <el-tag :type="getImpactType()" size="small">
          {{ $t(`gantt.ai.suggestions.impact.${suggestion.impact}`) }}
        </el-tag>
        <el-dropdown trigger="click" @command="handleCommand">
          <el-button :icon="MoreFilled" circle text />
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="details" :icon="View">
                {{ $t('gantt.ai.suggestions.viewDetails') }}
              </el-dropdown-item>
              <el-dropdown-item command="dismiss" :icon="Close" divided>
                {{ $t('gantt.ai.suggestions.dismiss') }}
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>

    <!-- Card Body -->
    <div class="card-body">
      <p class="suggestion-description">{{ suggestion.description }}</p>

      <!-- Actions Preview -->
      <div v-if="suggestion.actions?.length > 0" class="actions-preview">
        <div class="actions-label">
          {{ $t('gantt.ai.suggestions.suggestedActions') }} ({{ suggestion.actions.length }})
        </div>
        <div class="actions-list">
          <div
            v-for="(action, index) in suggestion.actions.slice(0, 2)"
            :key="index"
            class="action-item"
          >
            <el-icon class="action-icon"><Right /></el-icon>
            <span class="action-text">{{ action.description }}</span>
          </div>
          <div v-if="suggestion.actions.length > 2" class="more-actions">
            +{{ suggestion.actions.length - 2 }} {{ $t('gantt.ai.suggestions.moreActions') }}
          </div>
        </div>
      </div>

      <!-- Impact Preview -->
      <div v-if="suggestion.impactPreview" class="impact-preview">
        <el-descriptions :column="1" size="small" border>
          <el-descriptions-item
            v-for="(value, key) in getPreviewItems()"
            :key="key"
            :label="key"
          >
            {{ value }}
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </div>

    <!-- Card Footer -->
    <div class="card-footer">
      <div class="footer-meta">
        <span class="meta-item">
          <el-icon><Clock /></el-icon>
          {{ formatTime(suggestion.createdAt) }}
        </span>
        <el-tag
          v-if="suggestion.status !== 'pending'"
          :type="getStatusType()"
          size="small"
        >
          {{ $t(`gantt.ai.suggestions.status.${suggestion.status}`) }}
        </el-tag>
      </div>

      <div v-if="suggestion.status === 'pending'" class="footer-actions">
        <el-button size="small" @click="handleReject">
          {{ $t('gantt.ai.suggestions.reject') }}
        </el-button>
        <el-button type="primary" size="small" @click="handleAccept">
          {{ $t('gantt.ai.suggestions.accept') }}
        </el-button>
      </div>

      <div v-else class="footer-status">
        <el-icon v-if="suggestion.status === 'accepted'" color="#67C23A"><CircleCheck /></el-icon>
        <el-icon v-else-if="suggestion.status === 'rejected'" color="#F56C6C"><CircleClose /></el-icon>
        <el-icon v-else color="#909399"><InfoFilled /></el-icon>
        <span>{{ $t(`gantt.ai.suggestions.status.${suggestion.status}`) }}</span>
      </div>
    </div>

    <!-- Confirmation Dialog -->
    <el-dialog
      v-model="confirmDialogVisible"
      :title="$t('gantt.ai.suggestions.confirmAction')"
      width="400px"
    >
      <div class="confirm-content">
        <el-alert
          :type="confirmAction === 'accept' ? 'success' : 'warning'"
          :closable="false"
          show-icon
        >
          <template #title>
            {{ confirmAction === 'accept'
              ? $t('gantt.ai.suggestions.confirmAccept')
              : $t('gantt.ai.suggestions.confirmReject') }}
          </template>
        </el-alert>
        <div class="confirm-details">
          <h4>{{ suggestion.title }}</h4>
          <p>{{ suggestion.description }}</p>
          <div v-if="confirmAction === 'accept' && suggestion.actions?.length > 0" class="confirm-actions">
            <h5>{{ $t('gantt.ai.suggestions.willApply') }}:</h5>
            <ul>
              <li v-for="(action, index) in suggestion.actions" :key="index">
                {{ action.description }}
              </li>
            </ul>
          </div>
        </div>
      </div>
      <template #footer>
        <el-button @click="confirmDialogVisible = false">
          {{ $t('common.cancel') }}
        </el-button>
        <el-button
          :type="confirmAction === 'accept' ? 'primary' : 'warning'"
          @click="executeAction"
        >
          {{ $t('common.confirm') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  Warning,
  Lightbulb,
  TrendCharts,
  CircleCheck,
  CircleClose,
  InfoFilled,
  MoreFilled,
  Right,
  Clock,
  View,
  Close
} from '@element-plus/icons-vue'
import { formatDistanceToNow } from 'date-fns'
import { zhCN, enUS } from 'date-fns/locale'

const props = defineProps({
  suggestion: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['accept', 'reject', 'dismiss', 'view-details'])

const { t, locale } = useI18n()
const confirmDialogVisible = ref(false)
const confirmAction = ref(null)

/**
 * Get type icon
 */
function getTypeIcon() {
  const icons = {
    risk: Warning,
    optimization: TrendCharts,
    resource: Lightbulb,
    dependency: TrendCharts
  }
  return icons[props.suggestion.type] || Lightbulb
}

/**
 * Get icon color
 */
function getIconColor() {
  const colors = {
    risk: '#F56C6C',
    optimization: '#409EFF',
    resource: '#E6A23C',
    dependency: '#67C23A'
  }
  return colors[props.suggestion.type] || '#909399'
}

/**
 * Get impact type
 */
function getImpactType() {
  const types = {
    high: 'danger',
    medium: 'warning',
    low: 'info'
  }
  return types[props.suggestion.impact] || 'info'
}

/**
 * Get status type
 */
function getStatusType() {
  const types = {
    accepted: 'success',
    rejected: 'danger',
    dismissed: 'info'
  }
  return types[props.suggestion.status] || 'info'
}

/**
 * Get preview items
 */
function getPreviewItems() {
  const preview = props.suggestion.impactPreview
  if (!preview) return {}

  // Limit to 3 items
  const items = {}
  let count = 0
  for (const [key, value] of Object.entries(preview)) {
    if (count >= 3) break
    items[key] = value
    count++
  }
  return items
}

/**
 * Format time
 */
function formatTime(dateStr) {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  const dateLocale = locale.value === 'zh' ? zhCN : enUS
  return formatDistanceToNow(date, { addSuffix: true, locale: dateLocale })
}

/**
 * Handle accept
 */
function handleAccept() {
  confirmAction.value = 'accept'
  confirmDialogVisible.value = true
}

/**
 * Handle reject
 */
function handleReject() {
  confirmAction.value = 'reject'
  confirmDialogVisible.value = true
}

/**
 * Execute action
 */
function executeAction() {
  if (confirmAction.value === 'accept') {
    emit('accept', props.suggestion)
  } else {
    emit('reject', props.suggestion)
  }
  confirmDialogVisible.value = false
}

/**
 * Handle dropdown command
 */
function handleCommand(command) {
  switch (command) {
    case 'details':
      emit('view-details', props.suggestion)
      break
    case 'dismiss':
      emit('dismiss', props.suggestion)
      break
  }
}
</script>

<script>
export default {
  name: 'SuggestionCard'
}
</script>

<style scoped>
.suggestion-card {
  background: white;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  margin-bottom: 12px;
  transition: all 0.2s ease;
  overflow: hidden;
}

.suggestion-card:hover {
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  transform: translateY(-1px);
}

/* Impact-specific borders */
.suggestion-card.impact-high {
  border-left: 4px solid #F56C6C;
}

.suggestion-card.impact-medium {
  border-left: 4px solid #E6A23C;
}

.suggestion-card.impact-low {
  border-left: 4px solid #409EFF;
}

/* Header */
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.header-left {
  display: flex;
  gap: 12px;
  flex: 1;
}

.type-icon {
  font-size: 24px;
  flex-shrink: 0;
}

.header-content {
  flex: 1;
}

.suggestion-title {
  margin: 0 0 4px 0;
  font-size: 15px;
  font-weight: 600;
  color: #303133;
  line-height: 1.4;
}

.suggestion-type {
  font-size: 12px;
  color: #909399;
  text-transform: capitalize;
}

.header-right {
  display: flex;
  gap: 8px;
  align-items: center;
}

/* Body */
.card-body {
  padding: 16px;
}

.suggestion-description {
  margin: 0 0 12px 0;
  font-size: 14px;
  color: #606266;
  line-height: 1.6;
}

/* Actions Preview */
.actions-preview {
  background: #f5f7fa;
  border-radius: 6px;
  padding: 12px;
  margin-bottom: 12px;
}

.actions-label {
  font-size: 12px;
  font-weight: 600;
  color: #909399;
  margin-bottom: 8px;
}

.actions-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.action-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #606266;
}

.action-icon {
  font-size: 14px;
  color: #409EFF;
  flex-shrink: 0;
}

.action-text {
  flex: 1;
}

.more-actions {
  font-size: 12px;
  color: #909399;
  padding-left: 22px;
}

/* Impact Preview */
.impact-preview {
  margin-top: 12px;
}

/* Footer */
.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #fafafa;
  border-top: 1px solid #f0f0f0;
}

.footer-meta {
  display: flex;
  gap: 12px;
  align-items: center;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #909399;
}

.footer-actions {
  display: flex;
  gap: 8px;
}

.footer-status {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #606266;
}

/* Confirm Dialog */
.confirm-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.confirm-details {
  padding: 12px;
  background: #f5f7fa;
  border-radius: 6px;
}

.confirm-details h4 {
  margin: 0 0 8px 0;
  font-size: 15px;
  font-weight: 600;
  color: #303133;
}

.confirm-details p {
  margin: 0 0 12px 0;
  font-size: 14px;
  color: #606266;
}

.confirm-actions h5 {
  margin: 0 0 8px 0;
  font-size: 13px;
  font-weight: 600;
  color: #303133;
}

.confirm-actions ul {
  margin: 0;
  padding-left: 20px;
  font-size: 13px;
  color: #606266;
}

.confirm-actions li {
  margin-bottom: 4px;
}

/* Responsive */
@media (max-width: 768px) {
  .card-header {
    flex-direction: column;
    gap: 12px;
  }

  .header-left {
    width: 100%;
  }

  .header-right {
    width: 100%;
    justify-content: flex-end;
  }

  .card-footer {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
  }

  .footer-actions {
    justify-content: stretch;
  }

  .footer-actions .el-button {
    flex: 1;
  }
}
</style>
