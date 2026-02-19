<template>
  <el-dialog
    v-model="dialogVisible"
    :title="t('gantt.resourceLeveling.title')"
    width="1000px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <div v-loading="loading" class="resource-leveling-dialog">
      <!-- Resource Conflicts Summary -->
      <div class="conflicts-summary">
        <el-alert
          :type="conflictCount > 0 ? 'warning' : 'success'"
          :closable="false"
        >
          <template #title>
            <span v-if="conflictCount > 0">
              {{ t('gantt.resourceLeveling.conflictsFound', { count: conflictCount }) }}
            </span>
            <span v-else>
              {{ t('gantt.resourceLeveling.noConflicts') }}
            </span>
          </template>
        </el-alert>
      </div>

      <!-- Leveling Options -->
      <div class="leveling-options">
        <div class="option-header">
          <span>{{ t('gantt.resourceLeveling.levelingOptions') }}</span>
        </div>
        <el-radio-group v-model="levelingMode" @change="handleModeChange">
          <el-radio value="manual">
            <div class="radio-content">
              <div class="radio-label">{{ t('gantt.resourceLeveling.manual') }}</div>
              <div class="radio-desc">{{ t('gantt.resourceLeveling.manualDesc') }}</div>
            </div>
          </el-radio>
          <el-radio value="auto">
            <div class="radio-content">
              <div class="radio-label">{{ t('gantt.resourceLeveling.auto') }}</div>
              <div class="radio-desc">{{ t('gantt.resourceLeveling.autoDesc') }}</div>
            </div>
          </el-radio>
        </el-radio-group>

        <!-- Auto-leveling settings -->
        <div v-if="levelingMode === 'auto'" class="auto-settings">
          <el-form label-width="200px">
            <el-form-item :label="t('gantt.resourceLeveling.priority')">
              <el-select v-model="levelingOptions.priority">
                <el-option
                  value="priority"
                  :label="t('gantt.resourceLeveling.taskPriority')"
                />
                <el-option
                  value="duration"
                  :label="t('gantt.resourceLeveling.taskDuration')"
                />
                <el-option
                  value="slack"
                  :label="t('gantt.resourceLeveling.totalSlack')"
                />
              </el-select>
            </el-form-item>
            <el-form-item :label="t('gantt.resourceLeveling.range')">
              <el-radio-group v-model="levelingOptions.range">
                <el-radio value="all">{{ t('gantt.resourceLeveling.allTasks') }}</el-radio>
                <el-radio value="selected">{{ t('gantt.resourceLeveling.selectedTasks') }}</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item :label="t('gantt.resourceLeveling.splitTasks')">
              <el-switch v-model="levelingOptions.allowSplitting" />
            </el-form-item>
            <el-form-item :label="t('gantt.resourceLeveling.adjustDependencies')">
              <el-switch v-model="levelingOptions.adjustDependencies" />
            </el-form-item>
          </el-form>
        </div>
      </div>

      <!-- Resource Conflicts List -->
      <div v-if="conflicts.length > 0" class="conflicts-list">
        <div class="list-header">
          <span>{{ t('gantt.resourceLeveling.conflictDetails') }}</span>
        </div>
        <el-table
          :data="conflicts"
          stripe
          max-height="300"
          @row-click="handleConflictClick"
        >
          <el-table-column
            prop="resourceName"
            :label="t('gantt.resourceLeveling.resource')"
            width="150"
          />
          <el-table-column
            prop="date"
            :label="t('gantt.resourceLeveling.date')"
            width="120"
          />
          <el-table-column
            prop="assigned"
            :label="t('gantt.resourceLeveling.assigned')"
            width="100"
            align="center"
          >
            <template #default="scope">
              <el-tag :type="scope.row.assigned > scope.row.capacity ? 'danger' : 'success'">
                {{ scope.row.assigned }} / {{ scope.row.capacity }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column
            prop="tasks"
            :label="t('gantt.resourceLeveling.conflictingTasks')"
          >
            <template #default="scope">
              <el-tag
                v-for="task in scope.row.tasks"
                :key="task.id"
                size="small"
                style="margin-right: 4px;"
              >
                {{ task.name }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column
            :label="t('common.actions')"
            width="100"
            align="center"
          >
            <template #default="scope">
              <el-button
                type="primary"
                size="small"
                link
                @click.stop="handleResolveConflict(scope.row)"
              >
                {{ t('gantt.resourceLeveling.resolve') }}
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- Before/After Preview -->
      <div v-if="showPreview" class="preview-comparison">
        <div class="preview-header">
          <span>{{ t('gantt.resourceLeveling.preview') }}</span>
        </div>
        <div class="comparison-container">
          <div class="comparison-panel before">
            <div class="panel-title">{{ t('gantt.resourceLeveling.before') }}</div>
            <div class="gantt-preview">
              <GanttMiniView
                :tasks="originalTasks"
                :conflicts="conflicts"
                :height="150"
              />
            </div>
          </div>
          <div class="comparison-divider">→</div>
          <div class="comparison-panel after">
            <div class="panel-title">{{ t('gantt.resourceLeveling.after') }}</div>
            <div class="gantt-preview">
              <GanttMiniView
                :tasks="leveledTasks"
                :conflicts="[]"
                :height="150"
              />
            </div>
          </div>
        </div>
        <div class="preview-stats">
          <el-row :gutter="20">
            <el-col :span="8">
              <div class="stat-item">
                <div class="stat-label">{{ t('gantt.resourceLeveling.tasksDelayed') }}</div>
                <div class="stat-value">{{ levelingStats.tasksDelayed }}</div>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="stat-item">
                <div class="stat-label">{{ t('gantt.resourceLeveling.maxDelay') }}</div>
                <div class="stat-value">{{ levelingStats.maxDelay }} {{ t('common.days') }}</div>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="stat-item">
                <div class="stat-label">{{ t('gantt.resourceLeveling.projectExtension') }}</div>
                <div class="stat-value">{{ levelingStats.projectExtension }} {{ t('common.days') }}</div>
              </div>
            </el-col>
          </el-row>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">{{ t('common.cancel') }}</el-button>
        <el-button
          type="primary"
          :loading="leveling"
          :disabled="conflictCount === 0"
          @click="handleLevelResources"
        >
          {{ t('gantt.resourceLeveling.levelResources') }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import {
  detectResourceConflicts,
  applyResourceLeveling,
  calculateLevelingStatistics
} from '@/utils/resourceLeveling'
import GanttMiniView from '../views/GanttMiniView.vue'

const { t } = useI18n()

// Props
const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  tasks: {
    type: Array,
    default: () => []
  },
  resources: {
    type: Array,
    default: () => []
  }
})

// Emits
const emit = defineEmits(['update:modelValue', 'resourcesLeveled'])

// State
const dialogVisible = ref(false)
const loading = ref(false)
const leveling = ref(false)
const showPreview = ref(false)

const levelingMode = ref('manual')
const levelingOptions = ref({
  priority: 'priority',
  range: 'all',
  allowSplitting: false,
  adjustDependencies: false
})

const conflicts = ref([])
const originalTasks = ref([])
const leveledTasks = ref([])
const levelingStats = ref({
  tasksDelayed: 0,
  maxDelay: 0,
  projectExtension: 0
})

// Computed
const conflictCount = computed(() => conflicts.value.length)

// Methods
const loadConflicts = async () => {
  loading.value = true
  try {
    conflicts.value = await detectResourceConflicts(props.tasks, props.resources)
  } catch (error) {
    console.error('Failed to load conflicts:', error)
    ElMessage.error(t('gantt.resourceLeveling.loadConflictsError'))
  } finally {
    loading.value = false
  }
}

const handleModeChange = () => {
  if (levelingMode.value === 'auto' && conflicts.value.length > 0) {
    previewLeveling()
  }
}

const previewLeveling = async () => {
  loading.value = true
  try {
    originalTasks.value = [...props.tasks]
    leveledTasks.value = await applyResourceLeveling(
      props.tasks,
      props.resources,
      levelingOptions.value
    )
    levelingStats.value = calculateLevelingStatistics(
      originalTasks.value,
      leveledTasks.value
    )
    showPreview.value = true
  } catch (error) {
    console.error('Failed to preview leveling:', error)
    ElMessage.error(t('gantt.resourceLeveling.previewError'))
  } finally {
    loading.value = false
  }
}

const handleLevelResources = async () => {
  leveling.value = true
  try {
    let resultTasks
    if (levelingMode.value === 'auto') {
      resultTasks = leveledTasks.value
    } else {
      resultTasks = await applyResourceLeveling(
        props.tasks,
        props.resources,
        levelingOptions.value
      )
    }

    ElMessage.success(t('gantt.resourceLeveling.success'))
    emit('resourcesLeveled', resultTasks)
    handleClose()
  } catch (error) {
    console.error('Failed to level resources:', error)
    ElMessage.error(t('gantt.resourceLeveling.levelingError'))
  } finally {
    leveling.value = false
  }
}

const handleConflictClick = (row) => {
  // TODO: Highlight conflicting tasks in main Gantt view
  console.log('Conflict clicked:', row)
}

const handleResolveConflict = async (conflict) => {
  // TODO: Open conflict resolution dialog
  console.log('Resolve conflict:', conflict)
}

const handleClose = () => {
  dialogVisible.value = false
  emit('update:modelValue', false)
  nextTick(() => {
    conflicts.value = []
    showPreview.value = false
    levelingMode.value = 'manual'
    leveledTasks.value = []
  })
}

// Watch modelValue changes
watch(() => props.modelValue, (newVal) => {
  dialogVisible.value = newVal
  if (newVal) {
    loadConflicts()
  }
})

watch(dialogVisible, (newVal) => {
  emit('update:modelValue', newVal)
})
</script>

<style scoped>
.resource-leveling-dialog {
  padding: 10px 0;
}

.conflicts-summary {
  margin-bottom: 24px;
}

.leveling-options {
  margin-bottom: 24px;
  padding: 16px;
  background: #f5f7fa;
  border-radius: 4px;
}

.option-header {
  font-weight: 500;
  margin-bottom: 12px;
  color: #606266;
}

.radio-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.radio-label {
  font-weight: 500;
}

.radio-desc {
  font-size: 12px;
  color: #909399;
}

.auto-settings {
  margin-top: 16px;
  padding: 16px;
  background: white;
  border-radius: 4px;
}

.conflicts-list {
  margin-bottom: 24px;
}

.list-header {
  font-weight: 500;
  margin-bottom: 12px;
  color: #606266;
}

.preview-comparison {
  margin-top: 24px;
  padding: 16px;
  background: #f5f7fa;
  border-radius: 4px;
}

.preview-header {
  font-weight: 500;
  margin-bottom: 16px;
  color: #606266;
}

.comparison-container {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
}

.comparison-panel {
  flex: 1;
  background: white;
  border-radius: 4px;
  padding: 12px;
}

.panel-title {
  font-weight: 500;
  margin-bottom: 8px;
  color: #606266;
  text-align: center;
}

.gantt-preview {
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  overflow: hidden;
}

.comparison-divider {
  font-size: 24px;
  color: #409EFF;
  font-weight: bold;
}

.preview-stats {
  padding: 16px;
  background: white;
  border-radius: 4px;
}

.stat-item {
  text-align: center;
}

.stat-label {
  font-size: 12px;
  color: #909399;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #409EFF;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
