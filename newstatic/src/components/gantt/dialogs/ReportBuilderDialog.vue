<template>
  <el-dialog
    v-model="dialogVisible"
    title="Report Builder"
    width="900px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <div class="report-builder">
      <!-- Report Type Selection -->
      <div class="builder-section">
        <h3>1. Select Report Type</h3>
        <el-radio-group v-model="config.type" size="large">
          <el-radio-button label="task">Task Report</el-radio-button>
          <el-radio-button label="resource">Resource Report</el-radio-button>
          <el-radio-button label="milestone">Milestone Report</el-radio-button>
          <el-radio-button label="progress">Progress Report</el-radio-button>
        </el-radio-group>
      </div>

      <!-- Date Range -->
      <div class="builder-section">
        <h3>2. Select Date Range</h3>
        <el-date-picker
          v-model="config.dateRange"
          type="daterange"
          range-separator="to"
          start-placeholder="Start date"
          end-placeholder="End date"
          style="width: 100%"
        />
      </div>

      <!-- Column Selection -->
      <div class="builder-section">
        <h3>3. Select Columns</h3>
        <el-checkbox-group v-model="config.columns">
          <el-row :gutter="12">
            <el-col
              v-for="column in availableColumns"
              :key="column.key"
              :span="6"
            >
              <el-checkbox :label="column.key">
                {{ column.label }}
              </el-checkbox>
            </el-col>
          </el-row>
        </el-checkbox-group>
      </div>

      <!-- Grouping -->
      <div class="builder-section">
        <h3>4. Group By</h3>
        <el-select
          v-model="config.groupBy"
          placeholder="Select grouping option"
          clearable
          style="width: 100%"
        >
          <el-option label="Status" value="status" />
          <el-option label="Assignee" value="assignee" />
          <el-option label="Priority" value="priority" />
          <el-option label="Phase" value="phase" />
          <el-option label="Milestone" value="milestone" />
        </el-select>
      </div>

      <!-- Filters -->
      <div class="builder-section">
        <h3>5. Filters</h3>
        <div class="filters-container">
          <div
            v-for="(filter, index) in config.filters"
            :key="index"
            class="filter-row"
          >
            <el-select
              v-model="filter.field"
              placeholder="Field"
              style="width: 150px"
              @change="handleFilterFieldChange(index)"
            >
              <el-option
                v-for="field in filterFields"
                :key="field.key"
                :label="field.label"
                :value="field.key"
              />
            </el-select>

            <el-select
              v-model="filter.operator"
              placeholder="Operator"
              style="width: 120px"
            >
              <el-option label="Equals" value="eq" />
              <el-option label="Not Equals" value="ne" />
              <el-option label="Contains" value="contains" />
              <el-option label="Greater Than" value="gt" />
              <el-option label="Less Than" value="lt" />
            </el-select>

            <el-input
              v-model="filter.value"
              placeholder="Value"
              style="flex: 1"
            />

            <el-button
              :icon="Delete"
              type="danger"
              text
              @click="removeFilter(index)"
            />
          </div>

          <el-button
            :icon="Plus"
            type="primary"
            text
            @click="addFilter"
          >
            Add Filter
          </el-button>
        </div>
      </div>

      <!-- Sorting -->
      <div class="builder-section">
        <h3>6. Sorting</h3>
        <el-select
          v-model="config.sortBy"
          placeholder="Sort by"
          style="width: 200px; margin-right: 12px"
        >
          <el-option
            v-for="column in availableColumns"
            :key="column.key"
            :label="column.label"
            :value="column.key"
          />
        </el-select>

        <el-radio-group v-model="config.sortOrder">
          <el-radio-button label="asc">Ascending</el-radio-button>
          <el-radio-button label="desc">Descending</el-radio-button>
        </el-radio-group>
      </div>

      <!-- Preview -->
      <div class="builder-section">
        <h3>7. Preview</h3>
        <el-tabs v-model="activeTab">
          <el-tab-pane label="Data Preview" name="data">
            <el-table
              :data="previewData"
              stripe
              border
              height="300"
              style="width: 100%"
            >
              <el-table-column
                v-for="column in selectedColumns"
                :key="column.key"
                :prop="column.key"
                :label="column.label"
                :width="column.width"
              />
            </el-table>
          </el-tab-pane>

          <el-tab-pane label="Summary" name="summary">
            <div class="summary-stats">
              <div class="stat-item">
                <span class="stat-label">Total Records:</span>
                <span class="stat-value">{{ previewData.length }}</span>
              </div>
              <div v-if="config.groupBy" class="stat-item">
                <span class="stat-label">Groups:</span>
                <span class="stat-value">{{ groupCount }}</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">Columns:</span>
                <span class="stat-value">{{ config.columns.length }}</span>
              </div>
            </div>
          </el-tab-pane>
        </el-tabs>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">Cancel</el-button>
        <el-button
          :icon="Download"
          @click="exportToExcel"
        >
          Export to Excel
        </el-button>
        <el-button
          :icon="Document"
          @click="exportToPDF"
        >
          Export to PDF
        </el-button>
        <el-button
          type="primary"
          :icon="Printer"
          @click="printReport"
        >
          Print
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import {
  Plus,
  Delete,
  Download,
  Document,
  Printer
} from '@element-plus/icons-vue'
import { generateReport, exportToPDF, exportToExcel, generatePrintLayout } from '@/utils/reportGenerator'
import { ElMessage } from 'element-plus'

/**
 * ReportBuilderDialog Component
 * Dialog for building and exporting custom reports
 *
 * @props {Boolean} modelValue - Dialog visibility
 * @props {Array} tasks - Task data to report on
 * @props {Array} resources - Resource data
 * @props {Array} milestones - Milestone data
 *
 * @emits {Boolean} update:modelValue - Dialog visibility update
 * @emits {Object} report-generated - Emitted when report is generated
 */

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
  },
  milestones: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:modelValue', 'report-generated'])

// State
const dialogVisible = ref(props.modelValue)
const activeTab = ref('data')
const config = ref({
  type: 'task',
  dateRange: [],
  columns: ['name', 'status', 'assignee', 'progress'],
  groupBy: '',
  filters: [],
  sortBy: 'name',
  sortOrder: 'asc'
})

const filterFields = ref([
  { key: 'name', label: 'Name' },
  { key: 'status', label: 'Status' },
  { key: 'assignee', label: 'Assignee' },
  { key: 'priority', label: 'Priority' },
  { key: 'progress', label: 'Progress' },
  { key: 'startDate', label: 'Start Date' },
  { key: 'endDate', label: 'End Date' }
])

// Computed
const availableColumns = computed(() => {
  const baseColumns = [
    { key: 'id', label: 'ID', width: 80 },
    { key: 'name', label: 'Task Name', width: 200 },
    { key: 'status', label: 'Status', width: 120 },
    { key: 'assignee', label: 'Assignee', width: 120 },
    { key: 'priority', label: 'Priority', width: 100 },
    { key: 'progress', label: 'Progress', width: 100 },
    { key: 'startDate', label: 'Start Date', width: 120 },
    { key: 'endDate', label: 'End Date', width: 120 },
    { key: 'duration', label: 'Duration', width: 100 },
    { key: 'budget', label: 'Budget', width: 100 },
    { key: 'actualCost', label: 'Actual Cost', width: 100 }
  ]

  // Add type-specific columns
  if (config.value.type === 'resource') {
    baseColumns.push(
      { key: 'role', label: 'Role', width: 120 },
      { key: 'capacity', label: 'Capacity', width: 100 },
      { key: 'utilization', label: 'Utilization', width: 100 }
    )
  } else if (config.value.type === 'milestone') {
    baseColumns.push(
      { key: 'date', label: 'Date', width: 120 },
      { key: 'description', label: 'Description', width: 200 }
    )
  }

  return baseColumns
})

const selectedColumns = computed(() => {
  return availableColumns.value.filter(col =>
    config.value.columns.includes(col.key)
  )
})

const previewData = computed(() => {
  // Generate preview data based on config
  let data = []

  if (config.value.type === 'task') {
    data = props.tasks
  } else if (config.value.type === 'resource') {
    data = props.resources
  } else if (config.value.type === 'milestone') {
    data = props.milestones
  } else if (config.value.type === 'progress') {
    data = props.tasks
  }

  // Apply filters
  data = data.filter(item => {
    return config.value.filters.every(filter => {
      if (!filter.field || !filter.operator || !filter.value) return true

      const itemValue = item[filter.field]
      const filterValue = filter.value

      switch (filter.operator) {
        case 'eq':
          return itemValue === filterValue
        case 'ne':
          return itemValue !== filterValue
        case 'contains':
          return String(itemValue).toLowerCase().includes(String(filterValue).toLowerCase())
        case 'gt':
          return Number(itemValue) > Number(filterValue)
        case 'lt':
          return Number(itemValue) < Number(filterValue)
        default:
          return true
      }
    })
  })

  // Apply sorting
  data = [...data].sort((a, b) => {
    const aVal = a[config.value.sortBy]
    const bVal = b[config.value.sortBy]

    if (config.value.sortOrder === 'asc') {
      return aVal > bVal ? 1 : -1
    } else {
      return aVal < bVal ? 1 : -1
    }
  })

  // Limit to 100 rows for preview
  return data.slice(0, 100)
})

const groupCount = computed(() => {
  if (!config.value.groupBy) return 0

  const groups = new Set()
  previewData.value.forEach(item => {
    groups.add(item[config.value.groupBy])
  })

  return groups.size
})

// Methods
const addFilter = () => {
  config.value.filters.push({
    field: '',
    operator: 'eq',
    value: ''
  })
}

const removeFilter = (index) => {
  config.value.filters.splice(index, 1)
}

const handleFilterFieldChange = (index) => {
  // Reset operator and value when field changes
  config.value.filters[index].operator = 'eq'
  config.value.filters[index].value = ''
}

const handleClose = () => {
  dialogVisible.value = false
  emit('update:modelValue', false)
}

const exportToExcel = async () => {
  try {
    const report = generateReport({
      data: previewData.value,
      columns: selectedColumns.value,
      config: config.value
    })

    await exportToExcel(report)
    ElMessage.success('Report exported to Excel successfully')
    emit('report-generated', report)
  } catch (error) {
    ElMessage.error('Failed to export to Excel: ' + error.message)
  }
}

const exportToPDF = async () => {
  try {
    const report = generateReport({
      data: previewData.value,
      columns: selectedColumns.value,
      config: config.value
    })

    await exportToPDF(report)
    ElMessage.success('Report exported to PDF successfully')
    emit('report-generated', report)
  } catch (error) {
    ElMessage.error('Failed to export to PDF: ' + error.message)
  }
}

const printReport = () => {
  try {
    const report = generateReport({
      data: previewData.value,
      columns: selectedColumns.value,
      config: config.value
    })

    const layout = generatePrintLayout(report)

    // Open print window
    const printWindow = window.open('', '_blank')
    printWindow.document.write(layout)
    printWindow.document.close()
    printWindow.print()

    ElMessage.success('Report sent to printer')
  } catch (error) {
    ElMessage.error('Failed to print report: ' + error.message)
  }
}

// Watchers
watch(() => props.modelValue, (newVal) => {
  dialogVisible.value = newVal
})

watch(dialogVisible, (newVal) => {
  emit('update:modelValue', newVal)
})
</script>

<style scoped lang="scss">
.report-builder {
  .builder-section {
    margin-bottom: 24px;
    padding-bottom: 24px;
    border-bottom: 1px solid #ebeef5;

    &:last-child {
      border-bottom: none;
    }

    h3 {
      margin: 0 0 12px 0;
      font-size: 14px;
      font-weight: 600;
      color: #303133;
    }

    .filters-container {
      .filter-row {
        display: flex;
        gap: 12px;
        margin-bottom: 12px;
        align-items: center;
      }
    }

    .summary-stats {
      display: flex;
      justify-content: space-around;
      padding: 16px;
      background: #f5f7fa;
      border-radius: 4px;

      .stat-item {
        text-align: center;

        .stat-label {
          display: block;
          font-size: 12px;
          color: #909399;
          margin-bottom: 4px;
        }

        .stat-value {
          display: block;
          font-size: 18px;
          font-weight: 700;
          color: #303133;
        }
      }
    }
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

/* Responsive Design */
@media (max-width: 768px) {
  .report-builder {
    .filter-row {
      flex-wrap: wrap;
    }
  }
}
</style>
