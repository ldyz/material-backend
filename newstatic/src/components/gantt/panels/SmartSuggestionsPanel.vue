<template>
  <div class="smart-suggestions-panel">
    <!-- Panel Header -->
    <div class="panel-header">
      <div class="header-left">
        <h3>{{ $t('gantt.ai.suggestions.title') }}</h3>
        <el-badge
          v-if="suggestions.length > 0"
          :value="suggestions.length"
          class="suggestion-count"
        />
      </div>
      <div class="header-actions">
        <el-button
          :icon="Refresh"
          circle
          size="small"
          :loading="analyzing"
          @click="analyzeSchedule"
        />
        <el-button
          :icon="Close"
          circle
          size="small"
          @click="$emit('close')"
        />
      </div>
    </div>

    <!-- Score Overview -->
    <div v-if="analysis" class="score-overview">
      <div class="score-card" :class="getScoreClass(analysis.overallScore)">
        <div class="score-circle">
          <svg viewBox="0 0 36 36" class="circular-chart">
            <path
              class="circle-bg"
              d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
            />
            <path
              class="circle"
              :stroke-dasharray="`${analysis.overallScore}, 100`"
              d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
            />
          </svg>
          <div class="score-value">{{ analysis.overallScore }}</div>
        </div>
        <div class="score-label">
          <span class="label-text">{{ $t('gantt.ai.suggestions.healthScore') }}</span>
          <span class="label-status" :class="getScoreClass(analysis.overallScore)">
            {{ getScoreLabel(analysis.overallScore) }}
          </span>
        </div>
      </div>

      <div class="stats-grid">
        <div class="stat-item">
          <el-icon class="stat-icon" color="#F56C6C"><Warning /></el-icon>
          <div class="stat-content">
            <span class="stat-value">{{ analysis.risks?.length || 0 }}</span>
            <span class="stat-label">{{ $t('gantt.ai.suggestions.risks') }}</span>
          </div>
        </div>
        <div class="stat-item">
          <el-icon class="stat-icon" color="#409EFF"><Lightbulb /></el-icon>
          <div class="stat-content">
            <span class="stat-value">{{ analysis.suggestions?.length || 0 }}</span>
            <span class="stat-label">{{ $t('gantt.ai.suggestions.suggestions') }}</span>
          </div>
        </div>
        <div class="stat-item">
          <el-icon class="stat-icon" color="#67C23A"><TrendCharts /></el-icon>
          <div class="stat-content">
            <span class="stat-value">{{ analysis.optimizations?.length || 0 }}</span>
            <span class="stat-label">{{ $t('gantt.ai.suggestions.optimizations') }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Filter Tabs -->
    <div v-if="suggestions.length > 0" class="filter-tabs">
      <el-radio-group v-model="activeFilter" size="small">
        <el-radio-button label="all">
          {{ $t('gantt.ai.suggestions.filters.all') }}
          <el-badge :value="suggestions.length" />
        </el-radio-button>
        <el-radio-button label="pending">
          {{ $t('gantt.ai.suggestions.filters.pending') }}
          <el-badge :value="pendingCount" />
        </el-radio-button>
        <el-radio-button label="accepted">
          {{ $t('gantt.ai.suggestions.filters.accepted') }}
          <el-badge :value="acceptedCount" />
        </el-radio-button>
        <el-radio-button label="rejected">
          {{ $t('gantt.ai.suggestions.filters.rejected') }}
          <el-badge :value="rejectedCount" />
        </el-radio-button>
      </el-radio-group>
    </div>

    <!-- Suggestions List -->
    <div v-loading="analyzing" class="suggestions-list">
      <!-- Empty State -->
      <div v-if="!analyzing && filteredSuggestions.length === 0" class="empty-state">
        <el-empty
          :description="analysis ? $t('gantt.ai.suggestions.noSuggestions') : $t('gantt.ai.suggestions.runAnalysis')"
        >
          <el-button v-if="!analysis" type="primary" @click="analyzeSchedule">
            {{ $t('gantt.ai.suggestions.analyze') }}
          </el-button>
        </el-empty>
      </div>

      <!-- Suggestion Cards -->
      <SuggestionCard
        v-for="suggestion in filteredSuggestions"
        :key="suggestion.id"
        :suggestion="suggestion"
        @accept="handleAccept"
        @reject="handleReject"
        @dismiss="handleDismiss"
        @view-details="handleViewDetails"
      />
    </div>

    <!-- Detail Dialog -->
    <el-dialog
      v-model="detailDialogVisible"
      :title="$t('gantt.ai.suggestions.detailTitle')"
      width="600px"
    >
      <div v-if="selectedSuggestion" class="suggestion-detail">
        <div class="detail-header">
          <el-tag :type="getImpactType(selectedSuggestion.impact)" size="large">
            {{ $t(`gantt.ai.suggestions.impact.${selectedSuggestion.impact}`) }}
          </el-tag>
          <el-tag :type="getStatusType(selectedSuggestion.status)" size="large">
            {{ $t(`gantt.ai.suggestions.status.${selectedSuggestion.status}`) }}
          </el-tag>
        </div>

        <h4 class="detail-title">{{ selectedSuggestion.title }}</h4>
        <p class="detail-description">{{ selectedSuggestion.description }}</p>

        <div v-if="selectedSuggestion.actions?.length > 0" class="detail-actions">
          <h5>{{ $t('gantt.ai.suggestions.suggestedActions') }}</h5>
          <el-timeline>
            <el-timeline-item
              v-for="(action, index) in selectedSuggestion.actions"
              :key="index"
              :timestamp="action.type"
              placement="top"
            >
              <el-card>
                <h4>{{ action.description }}</h4>
                <div v-if="action.changes" class="action-changes">
                  <div v-for="(change, key) in action.changes" :key="key" class="change-item">
                    <span class="change-key">{{ key }}:</span>
                    <span class="change-value">{{ change }}</span>
                  </div>
                </div>
              </el-card>
            </el-timeline-item>
          </el-timeline>
        </div>

        <div v-if="selectedSuggestion.impactPreview" class="impact-preview">
          <h5>{{ $t('gantt.ai.suggestions.impactPreview') }}</h5>
          <el-descriptions :column="2" border>
            <el-descriptions-item
              v-for="(value, key) in selectedSuggestion.impactPreview"
              :key="key"
              :label="key"
            >
              {{ value }}
            </el-descriptions-item>
          </el-descriptions>
        </div>
      </div>

      <template #footer>
        <el-button @click="detailDialogVisible = false">
          {{ $t('common.close') }}
        </el-button>
        <el-button
          v-if="selectedSuggestion?.status === 'pending'"
          type="primary"
          @click="handleAccept(selectedSuggestion)"
        >
          {{ $t('gantt.ai.suggestions.accept') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Refresh, Close, Warning, Lightbulb, TrendCharts } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { analyzeSchedule, predictDelays, optimizeSchedule, createSuggestion } from '@/utils/aiOptimizer'
import SuggestionCard from './SuggestionCard.vue'

const props = defineProps({
  tasks: {
    type: Array,
    default: () => []
  },
  dependencies: {
    type: Array,
    default: () => []
  },
  resources: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['close', 'apply-suggestion', 'refresh'])

const { t } = useI18n()

const analyzing = ref(false)
const analysis = ref(null)
const suggestions = ref([])
const activeFilter = ref('all')
const detailDialogVisible = ref(false)
const selectedSuggestion = ref(null)

const pendingCount = computed(() =>
  suggestions.value.filter(s => s.status === 'pending').length
)

const acceptedCount = computed(() =>
  suggestions.value.filter(s => s.status === 'accepted').length
)

const rejectedCount = computed(() =>
  suggestions.value.filter(s => s.status === 'rejected').length
)

const filteredSuggestions = computed(() => {
  if (activeFilter.value === 'all') {
    return suggestions.value
  }
  return suggestions.value.filter(s => s.status === activeFilter.value)
})

/**
 * Analyze schedule
 */
async function analyzeSchedule() {
  analyzing.value = true

  try {
    // Run analysis
    const result = analyzeSchedule(props.tasks, props.dependencies, props.resources)
    analysis.value = result

    // Convert analysis results to suggestions
    const generatedSuggestions = []

    // Add risk suggestions
    result.risks?.forEach(risk => {
      generatedSuggestions.push(createSuggestion(
        'risk',
        risk.message,
        `Risk detected for task: ${risk.taskId}`,
        risk.level,
        [{
          type: 'mitigate',
          description: `Review and mitigate risk: ${risk.type}`,
          changes: { taskId: risk.taskId, riskType: risk.type }
        }]
      ))
    })

    // Add optimization suggestions
    result.suggestions?.forEach(suggestion => {
      generatedSuggestions.push({
        ...suggestion,
        status: 'pending',
        createdAt: new Date().toISOString()
      })
    })

    suggestions.value = generatedSuggestions

    ElMessage.success(t('gantt.ai.suggestions.analysisComplete'))
  } catch (error) {
    console.error('Analysis error:', error)
    ElMessage.error(t('gantt.ai.suggestions.analysisFailed'))
  } finally {
    analyzing.value = false
  }
}

/**
 * Handle accept suggestion
 */
function handleAccept(suggestion) {
  // Update status
  const index = suggestions.value.findIndex(s => s.id === suggestion.id)
  if (index !== -1) {
    suggestions.value[index].status = 'accepted'
  }

  // Emit apply event
  emit('apply-suggestion', suggestion)

  ElMessage.success(t('gantt.ai.suggestions.messages.accepted'))
  detailDialogVisible.value = false
}

/**
 * Handle reject suggestion
 */
function handleReject(suggestion) {
  const index = suggestions.value.findIndex(s => s.id === suggestion.id)
  if (index !== -1) {
    suggestions.value[index].status = 'rejected'
  }

  ElMessage.info(t('gantt.ai.suggestions.messages.rejected'))
  detailDialogVisible.value = false
}

/**
 * Handle dismiss suggestion
 */
function handleDismiss(suggestion) {
  const index = suggestions.value.findIndex(s => s.id === suggestion.id)
  if (index !== -1) {
    suggestions.value[index].status = 'dismissed'
  }

  ElMessage.info(t('gantt.ai.suggestions.messages.dismissed'))
}

/**
 * Handle view details
 */
function handleViewDetails(suggestion) {
  selectedSuggestion.value = suggestion
  detailDialogVisible.value = true
}

/**
 * Get score class
 */
function getScoreClass(score) {
  if (score >= 80) return 'excellent'
  if (score >= 60) return 'good'
  if (score >= 40) return 'fair'
  return 'poor'
}

/**
 * Get score label
 */
function getScoreLabel(score) {
  if (score >= 80) return t('gantt.ai.suggestions.score.excellent')
  if (score >= 60) return t('gantt.ai.suggestions.score.good')
  if (score >= 40) return t('gantt.ai.suggestions.score.fair')
  return t('gantt.ai.suggestions.score.poor')
}

/**
 * Get impact type
 */
function getImpactType(impact) {
  const types = {
    high: 'danger',
    medium: 'warning',
    low: 'info'
  }
  return types[impact] || 'info'
}

/**
 * Get status type
 */
function getStatusType(status) {
  const types = {
    pending: 'warning',
    accepted: 'success',
    rejected: 'danger',
    dismissed: 'info'
  }
  return types[status] || 'info'
}

// Auto-analyze on mount
onMounted(() => {
  if (props.tasks.length > 0) {
    analyzeSchedule()
  }
})
</script>

<script>
export default {
  name: 'SmartSuggestionsPanel'
}
</script>

<style scoped>
.smart-suggestions-panel {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #f5f7fa;
}

/* Header */
.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: white;
  border-bottom: 1px solid #e4e7ed;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-left h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.header-actions {
  display: flex;
  gap: 8px;
}

/* Score Overview */
.score-overview {
  padding: 20px;
  background: white;
  border-bottom: 1px solid #e4e7ed;
}

.score-card {
  display: flex;
  align-items: center;
  gap: 24px;
  margin-bottom: 20px;
  padding: 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  color: white;
}

.score-circle {
  position: relative;
  width: 80px;
  height: 80px;
}

.circular-chart {
  display: block;
  margin: 0 auto;
  max-width: 100%;
  max-height: 100%;
}

.circle-bg {
  fill: none;
  stroke: rgba(255, 255, 255, 0.2);
  stroke-width: 2;
}

.circle {
  fill: none;
  stroke: white;
  stroke-width: 2;
  stroke-linecap: round;
  transition: stroke-dasharray 0.6s ease;
}

.score-value {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 24px;
  font-weight: 700;
}

.score-label {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.label-text {
  font-size: 14px;
  opacity: 0.9;
}

.label-status {
  font-size: 18px;
  font-weight: 600;
  text-transform: capitalize;
}

/* Stats Grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #f5f7fa;
  border-radius: 8px;
}

.stat-icon {
  font-size: 24px;
}

.stat-content {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.stat-value {
  font-size: 20px;
  font-weight: 700;
  color: #303133;
}

.stat-label {
  font-size: 12px;
  color: #909399;
}

/* Filter Tabs */
.filter-tabs {
  padding: 12px 20px;
  background: white;
  border-bottom: 1px solid #e4e7ed;
}

/* Suggestions List */
.suggestions-list {
  flex: 1;
  overflow-y: auto;
  padding: 16px 20px;
}

.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 300px;
}

/* Detail Dialog */
.suggestion-detail {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.detail-header {
  display: flex;
  gap: 12px;
}

.detail-title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #303133;
}

.detail-description {
  margin: 0;
  font-size: 14px;
  color: #606266;
  line-height: 1.6;
}

.detail-actions h5,
.impact-preview h5 {
  margin: 0 0 12px 0;
  font-size: 14px;
  font-weight: 600;
  color: #303133;
}

.action-changes {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 8px;
}

.change-item {
  display: flex;
  gap: 8px;
  font-size: 13px;
}

.change-key {
  font-weight: 600;
  color: #606266;
}

.change-value {
  color: #303133;
}

/* Responsive */
@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }

  .score-card {
    flex-direction: column;
    text-align: center;
  }
}
</style>
