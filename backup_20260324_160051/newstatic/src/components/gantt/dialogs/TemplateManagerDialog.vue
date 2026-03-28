<template>
  <el-dialog
    v-model="dialogVisible"
    :title="$t('gantt.templates.title')"
    width="90%"
    :close-on-click-modal="false"
    @close="handleClose"
    class="template-manager-dialog"
  >
    <div class="template-manager">
      <!-- Sidebar -->
      <div class="template-sidebar">
        <div class="sidebar-section">
          <h3>{{ $t('gantt.templates.categories') }}</h3>
          <div class="category-list">
            <div
              v-for="category in categories"
              :key="category.id"
              class="category-item"
              :class="{ active: selectedCategory === category.id }"
              @click="selectCategory(category.id)"
            >
              <el-icon><component :is="getCategoryIcon(category.icon)" /></el-icon>
              <span>{{ category.name }}</span>
              <el-badge
                :value="getCategoryCount(category.id)"
                class="category-badge"
              />
            </div>
          </div>
        </div>

        <div class="sidebar-section" v-if="recentTemplates.length > 0">
          <h3>{{ $t('gantt.templates.recent') }}</h3>
          <div class="recent-list">
            <div
              v-for="template in recentTemplates"
              :key="template.id"
              class="recent-item"
              @click="selectTemplate(template)"
            >
              <span class="recent-name">{{ template.name }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Main Content -->
      <div class="template-content">
        <!-- Search and Actions -->
        <div class="template-header">
          <el-input
            v-model="searchQuery"
            :placeholder="$t('gantt.templates.searchPlaceholder')"
            :prefix-icon="Search"
            clearable
            class="search-input"
          />

          <div class="template-actions">
            <el-button
              :icon="Plus"
              type="primary"
              @click="showCreateDialog"
            >
              {{ $t('gantt.templates.create') }}
            </el-button>
            <el-dropdown @command="handleImport">
              <el-button :icon="Upload">
                {{ $t('gantt.templates.import') }}
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="import">
                    {{ $t('gantt.templates.importFromFile') }}
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </div>

        <!-- Template Grid -->
        <div class="template-grid" v-loading="loading">
          <div
            v-for="template in filteredTemplates"
            :key="template.id"
            class="template-card"
            @click="selectTemplate(template)"
          >
            <div class="template-card-header">
              <h4>{{ template.name }}</h4>
              <el-dropdown trigger="click" @command="(cmd) => handleTemplateAction(cmd, template)">
                <el-button :icon="MoreFilled" circle text />
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="apply" :icon="Check">
                      {{ $t('gantt.templates.apply') }}
                    </el-dropdown-item>
                    <el-dropdown-item command="edit" :icon="Edit">
                      {{ $t('common.edit') }}
                    </el-dropdown-item>
                    <el-dropdown-item command="duplicate" :icon="CopyDocument">
                      {{ $t('common.duplicate') }}
                    </el-dropdown-item>
                    <el-dropdown-item command="export" :icon="Download">
                      {{ $t('common.export') }}
                    </el-dropdown-item>
                    <el-dropdown-item command="delete" :icon="Delete" divided>
                      {{ $t('common.delete') }}
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>

            <p class="template-description">{{ template.description || $t('gantt.templates.noDescription') }}</p>

            <div class="template-meta">
              <div class="meta-item">
                <el-icon><List /></el-icon>
                <span>{{ template.tasks?.length || 0 }} {{ $t('gantt.templates.tasks') }}</span>
              </div>
              <div class="meta-item">
                <el-icon><Clock /></el-icon>
                <span>{{ getTemplateDuration(template) }}</span>
              </div>
              <div class="meta-item">
                <el-icon><Calendar /></el-icon>
                <span>{{ formatDate(template.updatedAt) }}</span>
              </div>
            </div>

            <div class="template-category">
              <el-tag :type="getCategoryTagType(template.category)" size="small">
                {{ getCategoryName(template.category) }}
              </el-tag>
            </div>
          </div>

          <!-- Empty State -->
          <div v-if="filteredTemplates.length === 0" class="empty-state">
            <el-empty :description="$t('gantt.templates.empty')">
              <el-button type="primary" @click="showCreateDialog">
                {{ $t('gantt.templates.createFirst') }}
              </el-button>
            </el-empty>
          </div>
        </div>
      </div>
    </div>

    <!-- Template Preview/Action Dialog -->
    <el-dialog
      v-model="previewDialogVisible"
      :title="selectedTemplate?.name"
      width="70%"
      append-to-body
    >
      <div v-if="selectedTemplate" class="template-preview">
        <div class="preview-info">
          <el-descriptions :column="2" border>
            <el-descriptions-item :label="$t('gantt.templates.category')">
              {{ getCategoryName(selectedTemplate.category) }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('gantt.templates.lastModified')">
              {{ formatDate(selectedTemplate.updatedAt) }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('gantt.templates.taskCount')">
              {{ selectedTemplate.tasks?.length || 0 }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('gantt.templates.dependencyCount')">
              {{ selectedTemplate.dependencies?.length || 0 }}
            </el-descriptions-item>
          </el-descriptions>
        </div>

        <div class="preview-tasks">
          <h4>{{ $t('gantt.templates.taskList') }}</h4>
          <el-table :data="selectedTemplate.tasks" max-height="300">
            <el-table-column prop="name" :label="$t('gantt.task.name')" />
            <el-table-column prop="duration" :label="$t('gantt.task.duration')" width="100" />
            <el-table-column prop="start" :label="$t('gantt.task.start')" width="100" />
          </el-table>
        </div>
      </div>

      <template #footer>
        <el-button @click="previewDialogVisible = false">
          {{ $t('common.cancel') }}
        </el-button>
        <el-button type="primary" @click="applySelectedTemplate">
          {{ $t('gantt.templates.apply') }}
        </el-button>
      </template>
    </el-dialog>

    <!-- Create/Edit Template Dialog -->
    <el-dialog
      v-model="editDialogVisible"
      :title="isEditMode ? $t('gantt.templates.edit') : $t('gantt.templates.create')"
      width="500px"
      append-to-body
    >
      <el-form
        ref="formRef"
        :model="templateForm"
        :rules="formRules"
        label-width="120px"
      >
        <el-form-item :label="$t('gantt.templates.name')" prop="name">
          <el-input v-model="templateForm.name" clearable />
        </el-form-item>

        <el-form-item :label="$t('gantt.templates.description')" prop="description">
          <el-input
            v-model="templateForm.description"
            type="textarea"
            :rows="3"
          />
        </el-form-item>

        <el-form-item :label="$t('gantt.templates.category')" prop="category">
          <el-select v-model="templateForm.category" style="width: 100%">
            <el-option
              v-for="category in categories"
              :key="category.id"
              :label="category.name"
              :value="category.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item v-if="!isEditMode" :label="$t('gantt.templates.source')">
          <el-radio-group v-model="templateSource">
            <el-radio label="current">{{ $t('gantt.templates.currentProject') }}</el-radio>
            <el-radio label="blank">{{ $t('gantt.templates.blank') }}</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="editDialogVisible = false">
          {{ $t('common.cancel') }}
        </el-button>
        <el-button type="primary" @click="handleSaveTemplate" :loading="saving">
          {{ $t('common.save') }}
        </el-button>
      </template>
    </el-dialog>

    <!-- Hidden file input for import -->
    <input
      ref="fileInputRef"
      type="file"
      accept=".json"
      style="display: none"
      @change="handleFileImport"
    />
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useTemplateStore } from '@/stores/templateStore'
import {
  Search,
  Plus,
  Upload,
  Download,
  Edit,
  Delete,
  CopyDocument,
  Check,
  MoreFilled,
  List,
  Clock,
  Calendar,
  Document,
  Code,
  Building,
  Megaphone,
  Flag
} from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { formatDate, addDays } from '@/utils/dateFormat'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  projectData: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['update:modelValue', 'apply-template', 'create-template'])

const { t } = useI18n()
const templateStore = useTemplateStore()

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const loading = ref(false)
const searchQuery = ref('')
const selectedCategory = ref('all')
const previewDialogVisible = ref(false)
const editDialogVisible = ref(false)
const isEditMode = ref(false)
const saving = ref(false)
const templateSource = ref('current')
const fileInputRef = ref(null)
const formRef = ref(null)

const selectedTemplate = ref(null)
const editingTemplate = ref(null)

const templateForm = ref({
  name: '',
  description: '',
  category: 'custom'
})

const formRules = {
  name: [
    { required: true, message: t('gantt.templates.validation.nameRequired'), trigger: 'blur' }
  ],
  category: [
    { required: true, message: t('gantt.templates.validation.categoryRequired'), trigger: 'change' }
  ]
}

const categories = computed(() => templateStore.categories)
const recentTemplates = computed(() => templateStore.recentTemplates)
const filteredTemplates = computed(() => templateStore.filteredTemplates)

/**
 * Select category
 */
function selectCategory(categoryId) {
  selectedCategory.value = categoryId
}

/**
 * Get category count
 */
function getCategoryCount(categoryId) {
  return templateStore.templates.filter(t => t.category === categoryId).length
}

/**
 * Get category icon
 */
function getCategoryIcon(iconName) {
  const icons = {
    code: Code,
    building: Building,
    megaphone: Megaphone,
    calendar: Calendar,
    document: Document
  }
  return icons[iconName] || Document
}

/**
 * Get category name
 */
function getCategoryName(categoryId) {
  const category = categories.value.find(c => c.id === categoryId)
  return category ? category.name : categoryId
}

/**
 * Get category tag type
 */
function getCategoryTagType(categoryId) {
  const types = {
    software: 'primary',
    construction: 'success',
    marketing: 'warning',
    event: 'danger',
    custom: 'info'
  }
  return types[categoryId] || 'info'
}

/**
 * Get template duration
 */
function getTemplateDuration(template) {
  if (!template.tasks || template.tasks.length === 0) {
    return t('gantt.templates.unknown')
  }

  const totalDays = template.tasks.reduce((sum, task) => sum + (task.duration || 0), 0)
  return `${totalDays} ${t('gantt.templates.days')}`
}

/**
 * Format date
 */
function formatDate(dateStr) {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return formatDate(new Date(date))
}

/**
 * Select template for preview
 */
function selectTemplate(template) {
  selectedTemplate.value = template
  previewDialogVisible.value = true
}

/**
 * Apply selected template
 */
function applySelectedTemplate() {
  if (!selectedTemplate.value) return

  const appliedData = templateStore.applyTemplateToProject(selectedTemplate.value.id, props.projectData)
  if (appliedData) {
    emit('apply-template', appliedData)
    previewDialogVisible.value = false
    dialogVisible.value = false
  }
}

/**
 * Handle template action
 */
function handleTemplateAction(command, template) {
  switch (command) {
    case 'apply':
      selectedTemplate.value = template
      applySelectedTemplate()
      break
    case 'edit':
      showEditDialog(template)
      break
    case 'duplicate':
      duplicateTemplate(template)
      break
    case 'export':
      templateStore.exportTemplate(template.id)
      break
    case 'delete':
      deleteTemplate(template)
      break
  }
}

/**
 * Show create dialog
 */
function showCreateDialog() {
  isEditMode.value = false
  editingTemplate.value = null
  templateForm.value = {
    name: '',
    description: '',
    category: 'custom'
  }
  templateSource.value = 'current'
  editDialogVisible.value = true
}

/**
 * Show edit dialog
 */
function showEditDialog(template) {
  isEditMode.value = true
  editingTemplate.value = template
  templateForm.value = {
    name: template.name,
    description: template.description || '',
    category: template.category
  }
  editDialogVisible.value = true
}

/**
 * Save template
 */
async function handleSaveTemplate() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  saving.value = true

  try {
    if (isEditMode.value) {
      // Update existing template
      templateStore.updateTemplate(editingTemplate.value.id, templateForm.value)
      ElMessage.success(t('gantt.templates.messages.updated'))
    } else {
      // Create new template
      let templateData = {}

      if (templateSource.value === 'current') {
        // Create from current project
        templateData = {
          ...templateForm.value,
          tasks: props.projectData.tasks || [],
          dependencies: props.projectData.dependencies || [],
          resources: props.projectData.resources || []
        }
      } else {
        // Create blank template
        templateData = {
          ...templateForm.value,
          tasks: [],
          dependencies: [],
          resources: []
        }
      }

      const newTemplate = templateStore.addTemplate(templateData)
      emit('create-template', newTemplate)
      ElMessage.success(t('gantt.templates.messages.created'))
    }

    editDialogVisible.value = false
  } catch (error) {
    ElMessage.error(t('gantt.templates.messages.saveFailed'))
  } finally {
    saving.value = false
  }
}

/**
 * Duplicate template
 */
function duplicateTemplate(template) {
  templateStore.duplicateTemplate(template.id)
}

/**
 * Delete template
 */
function deleteTemplate(template) {
  ElMessageBox.confirm(
    t('gantt.templates.confirmDelete', { name: template.name }),
    t('common.warning'),
    {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    }
  )
    .then(() => {
      templateStore.deleteTemplate(template.id)
      ElMessage.success(t('gantt.templates.messages.deleted'))
    })
    .catch(() => {})
}

/**
 * Handle import
 */
function handleImport(command) {
  if (command === 'import') {
    fileInputRef.value?.click()
  }
}

/**
 * Handle file import
 */
async function handleFileImport(event) {
  const file = event.target.files?.[0]
  if (!file) return

  loading.value = true

  try {
    await templateStore.importTemplate(file)
    ElMessage.success(t('gantt.templates.messages.imported'))
  } catch (error) {
    console.error('Import error:', error)
  } finally {
    loading.value = false
    // Reset file input
    event.target.value = ''
  }
}

/**
 * Handle close
 */
function handleClose() {
  // Reset state
  searchQuery.value = ''
  selectedCategory.value = 'all'
  selectedTemplate.value = null
}

// Watch for search query changes
watch(searchQuery, (value) => {
  templateStore.searchQuery = value
})

// Watch for category changes
watch(selectedCategory, (value) => {
  templateStore.selectedCategory = value
})

// Initialize on mount
onMounted(() => {
  templateStore.initialize()
})
</script>

<script>
export default {
  name: 'TemplateManagerDialog'
}
</script>

<style scoped>
.template-manager-dialog :deep(.el-dialog__body) {
  padding: 0;
}

.template-manager {
  display: flex;
  height: 600px;
}

/* Sidebar */
.template-sidebar {
  width: 250px;
  background: #f5f7fa;
  border-right: 1px solid #e4e7ed;
  padding: 20px;
  overflow-y: auto;
}

.sidebar-section {
  margin-bottom: 30px;
}

.sidebar-section h3 {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 12px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.category-list,
.recent-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.category-item,
.recent-item {
  display: flex;
  align-items: center;
  padding: 10px 12px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  color: #606266;
  font-size: 14px;
}

.category-item:hover,
.recent-item:hover {
  background: #e4e7ed;
  color: #303133;
}

.category-item.active {
  background: #409eff;
  color: white;
}

.category-item .el-icon {
  margin-right: 8px;
  font-size: 18px;
}

.category-badge {
  margin-left: auto;
}

.recent-name {
  font-size: 14px;
}

/* Main Content */
.template-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 20px;
  overflow: hidden;
}

.template-header {
  display: flex;
  gap: 16px;
  margin-bottom: 20px;
}

.search-input {
  flex: 1;
  max-width: 400px;
}

.template-actions {
  display: flex;
  gap: 8px;
}

/* Template Grid */
.template-grid {
  flex: 1;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
  overflow-y: auto;
  padding: 4px;
}

.template-card {
  background: white;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  padding: 16px;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.template-card:hover {
  border-color: #409eff;
  box-shadow: 0 2px 12px rgba(64, 158, 255, 0.2);
  transform: translateY(-2px);
}

.template-card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.template-card-header h4 {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin: 0;
  flex: 1;
  margin-right: 8px;
}

.template-description {
  font-size: 13px;
  color: #909399;
  line-height: 1.5;
  margin: 0;
  min-height: 40px;
}

.template-meta {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: #909399;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.template-category {
  margin-top: auto;
}

/* Template Preview */
.template-preview {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.preview-info {
  margin-bottom: 16px;
}

.preview-tasks {
  flex: 1;
}

.preview-tasks h4 {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 12px;
}

/* Empty State */
.empty-state {
  grid-column: 1 / -1;
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 300px;
}

/* Responsive */
@media (max-width: 768px) {
  .template-manager {
    flex-direction: column;
    height: auto;
  }

  .template-sidebar {
    width: 100%;
    border-right: none;
    border-bottom: 1px solid #e4e7ed;
    max-height: 200px;
  }

  .template-header {
    flex-direction: column;
  }

  .search-input {
    max-width: none;
  }

  .template-grid {
    grid-template-columns: 1fr;
  }
}
</style>
