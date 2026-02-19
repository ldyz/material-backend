<template>
  <el-dialog
    v-model="dialogVisible"
    title="任务模板"
    width="900px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <!-- Tabs -->
    <el-tabs v-model="activeTab" class="template-tabs">
      <!-- Pre-defined Templates -->
      <el-tab-pane label="预设模板" name="preset">
        <div class="template-list">
          <div
            v-for="template in presetTemplates"
            :key="template.id"
            class="template-card"
            :class="{ 'is-selected': selectedTemplate?.id === template.id }"
            @click="selectTemplate(template)"
          >
            <div class="template-card__header">
              <el-icon class="template-icon" :style="{ color: template.color }">
                <component :is="template.icon" />
              </el-icon>
              <div class="template-info">
                <h4 class="template-name">{{ template.name }}</h4>
                <p class="template-desc">{{ template.description }}</p>
              </div>
              <el-icon v-if="selectedTemplate?.id === template.id" class="check-icon">
                <Check />
              </el-icon>
            </div>

            <div class="template-card__body">
              <div class="template-meta">
                <span class="meta-item">
                  <el-icon><Clock /></el-icon>
                  默认工期: {{ template.defaultDuration }} 天
                </span>
                <span class="meta-item" v-if="template.defaultPriority">
                  <el-icon><Flag /></el-icon>
                  优先级: {{ priorityLabel(template.defaultPriority) }}
                </span>
              </div>

              <div class="template-preview">
                <div class="preview-label">包含字段:</div>
                <el-tag
                  v-for="field in template.fields"
                  :key="field"
                  size="small"
                  type="info"
                  class="field-tag"
                >
                  {{ fieldLabel(field) }}
                </el-tag>
              </div>
            </div>

            <div class="template-card__footer">
              <el-button
                type="primary"
                size="small"
                :disabled="!selectedTemplate || selectedTemplate.id !== template.id"
                @click.stop="useTemplate(template)"
              >
                使用此模板
              </el-button>
              <el-button
                size="small"
                @click.stop="previewTemplate(template)"
              >
                预览
              </el-button>
            </div>
          </div>
        </div>
      </el-tab-pane>

      <!-- Custom Templates -->
      <el-tab-pane label="自定义模板" name="custom">
        <div class="custom-templates">
          <div class="custom-header">
            <h4>我的模板</h4>
            <el-button type="primary" size="small" @click="showCreateTemplate = true">
              <el-icon><Plus /></el-icon>
              新建模板
            </el-button>
          </div>

          <div v-if="customTemplates.length > 0" class="template-list">
            <div
              v-for="template in customTemplates"
              :key="template.id"
              class="template-card custom"
            >
              <div class="template-card__header">
                <el-icon class="template-icon" style="color: #409eff">
                  <Document />
                </el-icon>
                <div class="template-info">
                  <h4 class="template-name">{{ template.name }}</h4>
                  <p class="template-desc">{{ template.description }}</p>
                </div>
                <div class="template-actions">
                  <el-button
                    type="primary"
                    size="small"
                    @click="useTemplate(template)"
                  >
                    使用
                  </el-button>
                  <el-dropdown @command="(cmd) => handleTemplateAction(cmd, template)">
                    <el-button size="small" :icon="MoreFilled" />
                    <template #dropdown>
                      <el-dropdown-menu>
                        <el-dropdown-item command="edit">编辑</el-dropdown-item>
                        <el-dropdown-item command="duplicate">复制</el-dropdown-item>
                        <el-dropdown-item command="delete" divided>删除</el-dropdown-item>
                      </el-dropdown-menu>
                    </template>
                  </el-dropdown>
                </div>
              </div>

              <div class="template-card__body">
                <div class="template-meta">
                  <span class="meta-item">
                    <el-icon><Clock /></el-icon>
                    工期: {{ template.defaultDuration }} 天
                  </span>
                  <span class="meta-item">
                    <el-icon><Calendar /></el-icon>
                    创建于: {{ formatDate(template.createdAt) }}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <el-empty
            v-else
            description="暂无自定义模板"
            :image-size="100"
          />
        </div>
      </el-tab-pane>

      <!-- Quick Create -->
      <el-tab-pane label="快速创建" name="quick">
        <el-form
          ref="quickFormRef"
          :model="quickForm"
          :rules="quickFormRules"
          label-width="100px"
          class="quick-form"
        >
          <el-form-item label="任务名称" prop="name">
            <el-input
              v-model="quickForm.name"
              placeholder="请输入任务名称"
              maxlength="100"
              show-word-limit
            />
          </el-form-item>

          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="任务类型" prop="type">
                <el-select v-model="quickForm.type" placeholder="选择类型">
                  <el-option label="普通任务" value="task" />
                  <el-option label="里程碑" value="milestone" />
                  <el-option label="阶段" value="phase" />
                  <el-option label="交付物" value="deliverable" />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="工期" prop="duration">
                <el-input-number
                  v-model="quickForm.duration"
                  :min="1"
                  :max="365"
                  :step="1"
                  controls-position="right"
                />
                <span class="unit-label">天</span>
              </el-form-item>
            </el-col>
          </el-row>

          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="优先级" prop="priority">
                <el-select v-model="quickForm.priority" placeholder="选择优先级">
                  <el-option label="高" value="high" />
                  <el-option label="中" value="medium" />
                  <el-option label="低" value="low" />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="进度" prop="progress">
                <el-slider
                  v-model="quickForm.progress"
                  :marks="progressMarks"
                  :step="10"
                />
              </el-form-item>
            </el-col>
          </el-row>

          <el-form-item label="开始日期" prop="startDate">
            <el-date-picker
              v-model="quickForm.startDate"
              type="date"
              placeholder="选择开始日期"
              value-format="YYYY-MM-DD"
              :disabled-date="disabledDate"
            />
          </el-form-item>

          <el-form-item label="任务描述">
            <el-input
              v-model="quickForm.description"
              type="textarea"
              :rows="3"
              placeholder="请输入任务描述（可选）"
              maxlength="500"
              show-word-limit
            />
          </el-form-item>

          <el-form-item label="保存为模板">
            <el-switch v-model="quickForm.saveAsTemplate" />
            <span v-if="quickForm.saveAsTemplate" class="template-name-input">
              <el-input
                v-model="quickForm.templateName"
                placeholder="模板名称"
                size="small"
                style="width: 200px; margin-left: 10px;"
              />
            </span>
          </el-form-item>
        </el-form>
      </el-tab-pane>
    </el-tabs>

    <!-- Preview Dialog -->
    <el-dialog
      v-model="previewVisible"
      title="模板预览"
      width="600px"
      append-to-body
    >
      <div v-if="previewTemplate" class="template-preview-detail">
        <div class="preview-header">
          <el-icon class="preview-icon" :style="{ color: previewTemplate.color }">
            <component :is="previewTemplate.icon" />
          </el-icon>
          <div>
            <h3>{{ previewTemplate.name }}</h3>
            <p>{{ previewTemplate.description }}</p>
          </div>
        </div>

        <el-divider />

        <div class="preview-sections">
          <div class="preview-section">
            <h4>基本信息</h4>
            <el-descriptions :column="2" border>
              <el-descriptions-item label="默认工期">
                {{ previewTemplate.defaultDuration }} 天
              </el-descriptions-item>
              <el-descriptions-item label="优先级">
                {{ priorityLabel(previewTemplate.defaultPriority) }}
              </el-descriptions-item>
              <el-descriptions-item label="默认进度">
                {{ previewTemplate.defaultProgress }}%
              </el-descriptions-item>
              <el-descriptions-item label="任务类型">
                {{ typeLabel(previewTemplate.type) }}
              </el-descriptions-item>
            </el-descriptions>
          </div>

          <div class="preview-section">
            <h4>包含字段</h4>
            <div class="fields-grid">
              <div
                v-for="field in previewTemplate.fields"
                :key="field"
                class="field-item"
              >
                <el-icon class="field-icon"><Check /></el-icon>
                <span>{{ fieldLabel(field) }}</span>
              </div>
            </div>
          </div>

          <div v-if="previewTemplate.defaultValues" class="preview-section">
            <h4>默认值</h4>
            <el-descriptions :column="1" border>
              <el-descriptions-item
                v-for="(value, key) in previewTemplate.defaultValues"
                :key="key"
                :label="fieldLabel(key)"
              >
                {{ value }}
              </el-descriptions-item>
            </el-descriptions>
          </div>
        </div>
      </div>
    </el-dialog>

    <!-- Create/Edit Template Dialog -->
    <el-dialog
      v-model="showCreateTemplate"
      :title="editingTemplate ? '编辑模板' : '新建模板'"
      width="600px"
      append-to-body
      @close="resetTemplateForm"
    >
      <el-form
        ref="templateFormRef"
        :model="templateForm"
        :rules="templateFormRules"
        label-width="100px"
      >
        <el-form-item label="模板名称" prop="name">
          <el-input v-model="templateForm.name" placeholder="请输入模板名称" />
        </el-form-item>

        <el-form-item label="模板描述" prop="description">
          <el-input
            v-model="templateForm.description"
            type="textarea"
            :rows="2"
            placeholder="请输入模板描述"
          />
        </el-form-item>

        <el-form-item label="默认工期" prop="defaultDuration">
          <el-input-number
            v-model="templateForm.defaultDuration"
            :min="1"
            :max="365"
          />
        </el-form-item>

        <el-form-item label="优先级" prop="defaultPriority">
          <el-select v-model="templateForm.defaultPriority">
            <el-option label="高" value="high" />
            <el-option label="中" value="medium" />
            <el-option label="低" value="low" />
          </el-select>
        </el-form-item>

        <el-form-item label="包含字段" prop="fields">
          <el-checkbox-group v-model="templateForm.fields">
            <el-checkbox label="name">任务名称</el-checkbox>
            <el-checkbox label="description">任务描述</el-checkbox>
            <el-checkbox label="startDate">开始日期</el-checkbox>
            <el-checkbox label="duration">工期</el-checkbox>
            <el-checkbox label="progress">进度</el-checkbox>
            <el-checkbox label="priority">优先级</el-checkbox>
            <el-checkbox label="assignee">负责人</el-checkbox>
            <el-checkbox label="resources">资源</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="showCreateTemplate = false">取消</el-button>
        <el-button type="primary" @click="saveCustomTemplate">保存</el-button>
      </template>
    </el-dialog>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button
        v-if="activeTab === 'quick'"
        type="primary"
        :loading="submitting"
        @click="handleQuickCreate"
      >
        创建任务
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Check,
  Clock,
  Flag,
  Plus,
  Document,
  MoreFilled,
  Calendar,
  Star,
  FolderOpened,
  Box
} from '@element-plus/icons-vue'
import { progressApi } from '@/api'
import { useUndoRedoStore } from '@/stores/undoRedoStore'
import { CreateTaskCommand } from '@/stores/undoRedoStore'
import eventBus, { GanttEvents } from '@/utils/eventBus'

/**
 * TaskTemplatesDialog Component
 *
 * Provides pre-defined and custom task templates for quick task creation
 * Integrates with undo/redo system
 *
 * @date 2025-02-18
 */

const props = defineProps({
  modelValue: Boolean,
  projectId: {
    type: [Number, String],
    required: true
  },
  startDate: {
    type: String,
    default: () => new Date().toISOString().split('T')[0]
  }
})

const emit = defineEmits(['update:modelValue', 'created', 'template-selected'])

// Store
const undoRedoStore = useUndoRedoStore()

// State
const dialogVisible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const activeTab = ref('preset')
const selectedTemplate = ref(null)
const previewVisible = ref(false)
const previewTemplateData = ref(null)
const showCreateTemplate = ref(false)
const editingTemplate = ref(null)
const submitting = ref(false)

// Preset templates
const presetTemplates = ref([
  {
    id: 'milestone',
    name: '里程碑',
    description: '项目关键节点，工期为0',
    icon: Star,
    color: '#f5a623',
    type: 'milestone',
    defaultDuration: 0,
    defaultPriority: 'high',
    defaultProgress: 0,
    fields: ['name', 'description', 'startDate'],
    defaultValues: {
      duration: 0,
      progress: 0
    }
  },
  {
    id: 'phase',
    name: '项目阶段',
    description: '项目的各个阶段划分',
    icon: FolderOpened,
    color: '#4a90e2',
    type: 'phase',
    defaultDuration: 30,
    defaultPriority: 'medium',
    defaultProgress: 0,
    fields: ['name', 'description', 'startDate', 'duration', 'progress'],
    defaultValues: {
      progress: 0
    }
  },
  {
    id: 'deliverable',
    name: '交付物',
    description: '项目交付成果',
    icon: Box,
    color: '#50e3c2',
    type: 'deliverable',
    defaultDuration: 7,
    defaultPriority: 'high',
    defaultProgress: 0,
    fields: ['name', 'description', 'startDate', 'duration', 'priority', 'assignee'],
    defaultValues: {
      priority: 'high'
    }
  },
  {
    id: 'task',
    name: '普通任务',
    description: '常规项目任务',
    icon: Document,
    color: '#909399',
    type: 'task',
    defaultDuration: 5,
    defaultPriority: 'medium',
    defaultProgress: 0,
    fields: ['name', 'description', 'startDate', 'duration', 'progress', 'priority', 'assignee', 'resources'],
    defaultValues: {
      priority: 'medium',
      progress: 0
    }
  },
  {
    id: 'review',
    name: '评审任务',
    description: '需要评审的任务',
    icon: Document,
    color: '#e6a23c',
    type: 'task',
    defaultDuration: 3,
    defaultPriority: 'high',
    defaultProgress: 0,
    fields: ['name', 'description', 'startDate', 'duration', 'progress', 'priority', 'assignee'],
    defaultValues: {
      priority: 'high',
      duration: 3
    }
  }
])

// Custom templates (from localStorage)
const customTemplates = ref([])

// Quick form
const quickFormRef = ref(null)
const quickForm = ref({
  name: '',
  type: 'task',
  duration: 5,
  priority: 'medium',
  progress: 0,
  startDate: props.startDate,
  description: '',
  saveAsTemplate: false,
  templateName: ''
})

const quickFormRules = {
  name: [
    { required: true, message: '请输入任务名称', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择任务类型', trigger: 'change' }
  ],
  duration: [
    { required: true, message: '请输入工期', trigger: 'blur' }
  ],
  startDate: [
    { required: true, message: '请选择开始日期', trigger: 'change' }
  ]
}

const progressMarks = {
  0: '0%',
  50: '50%',
  100: '100%'
}

// Template form
const templateFormRef = ref(null)
const templateForm = ref({
  name: '',
  description: '',
  defaultDuration: 5,
  defaultPriority: 'medium',
  defaultProgress: 0,
  fields: ['name', 'description', 'startDate', 'duration']
})

const templateFormRules = {
  name: [
    { required: true, message: '请输入模板名称', trigger: 'blur' }
  ],
  defaultDuration: [
    { required: true, message: '请输入默认工期', trigger: 'blur' }
  ]
}

// Methods
/**
 * Load custom templates from localStorage
 */
function loadCustomTemplates() {
  try {
    const saved = localStorage.getItem('gantt-task-templates')
    if (saved) {
      customTemplates.value = JSON.parse(saved)
    }
  } catch (error) {
    console.error('Failed to load custom templates:', error)
    customTemplates.value = []
  }
}

/**
 * Save custom templates to localStorage
 */
function saveCustomTemplates() {
  try {
    localStorage.setItem('gantt-task-templates', JSON.stringify(customTemplates.value))
  } catch (error) {
    console.error('Failed to save custom templates:', error)
  }
}

/**
 * Select a template
 */
function selectTemplate(template) {
  selectedTemplate.value = template
}

/**
 * Preview a template
 */
function previewTemplate(template) {
  previewTemplateData.value = template
  previewVisible.value = true
}

/**
 * Use a template to create a task
 */
async function useTemplate(template) {
  try {
    submitting.value = true

    // Build task data from template
    const taskData = {
      project_id: props.projectId,
      task_name: template.name || quickForm.value.name,
      start_date: template.defaultStartDate || props.startDate,
      duration: template.defaultDuration,
      progress: template.defaultProgress || 0,
      priority: template.defaultPriority || 'medium',
      description: template.description || ''
    }

    // Create command
    const command = new CreateTaskCommand(
      progressApi.create.bind(progressApi),
      taskData,
      null
    )

    // Execute command (integrates with undo/redo)
    const result = await undoRedoStore.executeCommand(command)

    ElMessage.success('任务已创建')
    emit('created', result)
    emit('template-selected', template)

    handleClose()
  } catch (error) {
    console.error('Failed to create task from template:', error)
    ElMessage.error(error.message || '创建任务失败')
  } finally {
    submitting.value = false
  }
}

/**
 * Quick create task
 */
async function handleQuickCreate() {
  try {
    await quickFormRef.value.validate()

    submitting.value = true

    const taskData = {
      project_id: props.projectId,
      task_name: quickForm.value.name,
      start_date: quickForm.value.startDate,
      duration: quickForm.value.type === 'milestone' ? 0 : quickForm.value.duration,
      progress: quickForm.value.progress,
      priority: quickForm.value.priority,
      description: quickForm.value.description
    }

    // Create command
    const command = new CreateTaskCommand(
      progressApi.create.bind(progressApi),
      taskData,
      null
    )

    // Execute command
    const result = await undoRedoStore.executeCommand(command)

    // Save as template if requested
    if (quickForm.value.saveAsTemplate && quickForm.value.templateName) {
      const newTemplate = {
        id: `custom-${Date.now()}`,
        name: quickForm.value.templateName,
        description: quickForm.value.description || '自定义模板',
        type: quickForm.value.type,
        defaultDuration: quickForm.value.duration,
        defaultPriority: quickForm.value.priority,
        defaultProgress: quickForm.value.progress,
        fields: ['name', 'description', 'startDate', 'duration', 'progress', 'priority'],
        defaultValues: taskData,
        createdAt: new Date().toISOString()
      }

      customTemplates.value.push(newTemplate)
      saveCustomTemplates()
    }

    ElMessage.success('任务已创建')
    emit('created', result)

    resetQuickForm()
    handleClose()
  } catch (error) {
    if (error !== false) { // Form validation error
      console.error('Failed to create task:', error)
      ElMessage.error(error.message || '创建任务失败')
    }
  } finally {
    submitting.value = false
  }
}

/**
 * Handle template action (edit, duplicate, delete)
 */
async function handleTemplateAction(command, template) {
  switch (command) {
    case 'edit':
      editingTemplate.value = template
      templateForm.value = { ...template }
      showCreateTemplate.value = true
      break

    case 'duplicate':
      const duplicate = {
        ...template,
        id: `custom-${Date.now()}`,
        name: `${template.name} (副本)`,
        createdAt: new Date().toISOString()
      }
      customTemplates.value.push(duplicate)
      saveCustomTemplates()
      ElMessage.success('模板已复制')
      break

    case 'delete':
      try {
        await ElMessageBox.confirm('确定要删除此模板吗？', '提示', {
          type: 'warning'
        })
        const index = customTemplates.value.findIndex(t => t.id === template.id)
        if (index > -1) {
          customTemplates.value.splice(index, 1)
          saveCustomTemplates()
          ElMessage.success('模板已删除')
        }
      } catch {
        // User cancelled
      }
      break
  }
}

/**
 * Save custom template
 */
async function saveCustomTemplate() {
  try {
    await templateFormRef.value.validate()

    if (editingTemplate.value) {
      // Update existing
      const index = customTemplates.value.findIndex(t => t.id === editingTemplate.value.id)
      if (index > -1) {
        customTemplates.value[index] = {
          ...customTemplates.value[index],
          ...templateForm.value
        }
      }
      ElMessage.success('模板已更新')
    } else {
      // Create new
      const newTemplate = {
        id: `custom-${Date.now()}`,
        ...templateForm.value,
        createdAt: new Date().toISOString()
      }
      customTemplates.value.push(newTemplate)
      ElMessage.success('模板已创建')
    }

    saveCustomTemplates()
    showCreateTemplate.value = false
    resetTemplateForm()
  } catch (error) {
    if (error !== false) {
      ElMessage.error('保存模板失败')
    }
  }
}

/**
 * Reset template form
 */
function resetTemplateForm() {
  editingTemplate.value = null
  templateForm.value = {
    name: '',
    description: '',
    defaultDuration: 5,
    defaultPriority: 'medium',
    defaultProgress: 0,
    fields: ['name', 'description', 'startDate', 'duration']
  }
  templateFormRef.value?.clearValidate()
}

/**
 * Reset quick form
 */
function resetQuickForm() {
  quickForm.value = {
    name: '',
    type: 'task',
    duration: 5,
    priority: 'medium',
    progress: 0,
    startDate: props.startDate,
    description: '',
    saveAsTemplate: false,
    templateName: ''
  }
  quickFormRef.value?.clearValidate()
}

/**
 * Close dialog
 */
function handleClose() {
  selectedTemplate.value = null
  resetQuickForm()
  emit('update:modelValue', false)
}

/**
 * Helper functions
 */
function priorityLabel(priority) {
  const map = {
    high: '高',
    medium: '中',
    low: '低'
  }
  return map[priority] || priority
}

function typeLabel(type) {
  const map = {
    task: '普通任务',
    milestone: '里程碑',
    phase: '阶段',
    deliverable: '交付物'
  }
  return map[type] || type
}

function fieldLabel(field) {
  const map = {
    name: '任务名称',
    description: '任务描述',
    startDate: '开始日期',
    duration: '工期',
    progress: '进度',
    priority: '优先级',
    assignee: '负责人',
    resources: '资源'
  }
  return map[field] || field
}

function formatDate(dateStr) {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

function disabledDate(time) {
  // Can't select past dates
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  return time.getTime() < today.getTime()
}

// Lifecycle
watch(() => props.modelValue, (val) => {
  if (val) {
    loadCustomTemplates()
    quickForm.value.startDate = props.startDate
  }
})

// Initialize
loadCustomTemplates()
</script>

<style scoped>
.template-tabs {
  min-height: 400px;
}

.template-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
  padding: 16px 0;
}

.template-card {
  border: 2px solid var(--el-border-color, #dcdfe6);
  border-radius: 8px;
  padding: 16px;
  cursor: pointer;
  transition: all 0.3s;
  background: #fff;
}

.template-card:hover {
  border-color: var(--el-color-primary, #409eff);
  box-shadow: 0 2px 12px rgba(64, 158, 255, 0.2);
}

.template-card.is-selected {
  border-color: var(--el-color-primary, #409eff);
  background: var(--el-color-primary-light-9, #ecf5ff);
}

.template-card.custom {
  cursor: default;
}

.template-card__header {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 12px;
}

.template-icon {
  font-size: 32px;
  flex-shrink: 0;
}

.template-info {
  flex: 1;
  min-width: 0;
}

.template-name {
  margin: 0 0 4px;
  font-size: 16px;
  font-weight: 500;
  color: var(--el-text-color-primary, #303133);
}

.template-desc {
  margin: 0;
  font-size: 12px;
  color: var(--el-text-color-secondary, #909399);
  line-height: 1.5;
}

.check-icon {
  font-size: 20px;
  color: var(--el-color-primary, #409eff);
  flex-shrink: 0;
}

.template-card__body {
  margin-bottom: 12px;
}

.template-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  margin-bottom: 12px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: var(--el-text-color-regular, #606266);
}

.template-preview {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
}

.preview-label {
  font-size: 12px;
  color: var(--el-text-color-secondary, #909399);
}

.field-tag {
  margin: 0;
}

.template-card__footer {
  display: flex;
  gap: 8px;
  padding-top: 12px;
  border-top: 1px solid var(--el-border-color-lighter, #ebeef5);
}

.template-actions {
  display: flex;
  gap: 8px;
  margin-left: auto;
}

/* Custom Templates Tab */
.custom-templates {
  padding: 16px 0;
}

.custom-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.custom-header h4 {
  margin: 0;
  font-size: 16px;
  font-weight: 500;
}

/* Quick Form */
.quick-form {
  padding: 16px 0;
}

.unit-label {
  margin-left: 8px;
  font-size: 12px;
  color: var(--el-text-color-secondary, #909399);
}

.template-name-input {
  display: inline-flex;
  align-items: center;
}

/* Preview Detail */
.template-preview-detail {
  padding: 16px 0;
}

.preview-header {
  display: flex;
  align-items: flex-start;
  gap: 16px;
}

.preview-icon {
  font-size: 48px;
}

.preview-header h3 {
  margin: 0 0 8px;
  font-size: 18px;
  font-weight: 500;
}

.preview-header p {
  margin: 0;
  font-size: 14px;
  color: var(--el-text-color-secondary, #909399);
}

.preview-sections {
  margin-top: 20px;
}

.preview-section {
  margin-bottom: 24px;
}

.preview-section h4 {
  margin: 0 0 12px;
  font-size: 14px;
  font-weight: 500;
  color: var(--el-text-color-primary, #303133);
}

.fields-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 12px;
}

.field-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: var(--el-fill-color-light, #f5f7fa);
  border-radius: 4px;
  font-size: 13px;
}

.field-icon {
  color: var(--el-color-success, #67c23a);
  font-size: 16px;
}

/* Responsive */
@media (max-width: 768px) {
  .template-list {
    grid-template-columns: 1fr;
  }

  .template-card__footer {
    flex-direction: column;
  }

  .template-card__footer .el-button {
    width: 100%;
  }
}
</style>
