<template>
  <div class="resource-utilization">
    <!-- Chart View -->
    <div v-if="view === 'chart'" class="chart-view">
      <canvas ref="chartCanvas" :height="height"></canvas>
    </div>

    <!-- Table View -->
    <div v-else class="table-view">
      <el-table :data="resourceData" stripe>
        <el-table-column prop="name" label="Resource" width="200">
          <template #default="{ row }">
            <div class="resource-cell">
              <el-avatar :size="32" :src="row.avatar">
                {{ row.name.charAt(0) }}
              </el-avatar>
              <div class="resource-info">
                <div class="resource-name">{{ row.name }}</div>
                <div class="resource-role">{{ row.role }}</div>
              </div>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="assignedTasks" label="Assigned Tasks" width="120" align="center">
          <template #default="{ row }">
            <el-tag>{{ row.assignedTasks }}</el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="totalHours" label="Total Hours" width="120" align="center">
          <template #default="{ row }">
            <span>{{ row.totalHours }}h</span>
          </template>
        </el-table-column>

        <el-table-column prop="capacity" label="Capacity" width="120" align="center">
          <template #default="{ row }">
            <span>{{ row.capacity }}h</span>
          </template>
        </el-table-column>

        <el-table-column prop="utilization" label="Utilization" width="150" align="center">
          <template #default="{ row }">
            <el-progress
              :percentage="row.utilization"
              :status="getUtilizationStatus(row.utilization)"
              :stroke-width="12"
            />
          </template>
        </el-table-column>

        <el-table-column prop="status" label="Status" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.utilization)">
              {{ getStatusLabel(row.utilization) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="Actions" width="100" align="center">
          <template #default="{ row }">
            <el-button
              :icon="View"
              size="small"
              text
              @click="viewDetails(row)"
            />
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- Summary -->
    <div class="utilization-summary">
      <div class="summary-item">
        <span class="summary-label">Total Resources:</span>
        <span class="summary-value">{{ resourceData.length }}</span>
      </div>
      <div class="summary-item">
        <span class="summary-label">Overallocated:</span>
        <span class="summary-value summary-danger">{{ overallocatedCount }}</span>
      </div>
      <div class="summary-item">
        <span class="summary-label">Underutilized:</span>
        <span class="summary-value summary-warning">{{ underutilizedCount }}</span>
      </div>
      <div class="summary-item">
        <span class="summary-label">Optimal:</span>
        <span class="summary-value summary-success">{{ optimalCount }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { View } from '@element-plus/icons-vue'
import { Doughnut, Bar } from 'chart.js/auto'
import { format, parseISO } from 'date-fns'

/**
 * ResourceUtilization Component
 * Displays resource usage via pie/donut chart or table
 *
 * @props {String} view - Display mode (chart/table)
 * @props {Array} dateRange - Date range for analysis
 * @props {Array} tasks - Array of tasks with assignments
 * @props {Array} resources - Array of team members
 * @props {Number} height - Chart height in pixels
 */

const props = defineProps({
  view: {
    type: String,
    default: 'chart',
    validator: (value) => ['chart', 'table'].includes(value)
  },
  dateRange: {
    type: Array,
    default: () => [new Date(), new Date()]
  },
  tasks: {
    type: Array,
    default: () => []
  },
  resources: {
    type: Array,
    default: () => []
  },
  height: {
    type: Number,
    default: 300
  }
})

const emit = defineEmits(['view-details'])

// State
const chartCanvas = ref(null)
let chartInstance = null

// Computed
const resourceData = computed(() => {
  return props.resources.map(resource => {
    const assignedTasks = props.tasks.filter(t => t.assignee === resource.name)
    const totalHours = assignedTasks.reduce((sum, task) => {
      return sum + (task.estimatedHours || 8)
    }, 0)
    const capacity = resource.capacity || 40 // Default 40h/week
    const utilization = Math.min(Math.round((totalHours / capacity) * 100), 100)

    return {
      ...resource,
      avatar: `https://ui-avatars.com/api/?name=${encodeURIComponent(resource.name)}&background=random`,
      assignedTasks: assignedTasks.length,
      totalHours,
      capacity,
      utilization
    }
  })
})

const overallocatedCount = computed(() => {
  return resourceData.value.filter(r => r.utilization > 100).length
})

const underutilizedCount = computed(() => {
  return resourceData.value.filter(r => r.utilization < 50).length
})

const optimalCount = computed(() => {
  return resourceData.value.filter(r => r.utilization >= 70 && r.utilization <= 100).length
})

// Methods
const getUtilizationStatus = (utilization) => {
  if (utilization > 100) return 'exception'
  if (utilization < 50) return 'warning'
  return 'success'
}

const getStatusType = (utilization) => {
  if (utilization > 100) return 'danger'
  if (utilization < 50) return 'warning'
  return 'success'
}

const getStatusLabel = (utilization) => {
  if (utilization > 100) return 'Overallocated'
  if (utilization < 50) return 'Available'
  if (utilization >= 70 && utilization <= 100) return 'Optimal'
  return 'Normal'
}

const viewDetails = (resource) => {
  emit('view-details', resource)
}

const initChart = () => {
  if (!chartCanvas.value || props.view !== 'chart') return

  const ctx = chartCanvas.value.getContext('2d')

  if (chartInstance) {
    chartInstance.destroy()
  }

  // Prepare data for donut chart
  const overallocated = resourceData.value.filter(r => r.utilization > 100).length
  const optimal = resourceData.value.filter(r => r.utilization >= 70 && r.utilization <= 100).length
  const available = resourceData.value.filter(r => r.utilization < 70).length

  chartInstance = new Doughnut(ctx, {
    type: 'doughnut',
    data: {
      labels: ['Overallocated', 'Optimal', 'Available'],
      datasets: [{
        data: [overallocated, optimal, available],
        backgroundColor: [
          '#F56C6C', // Red for overallocated
          '#67C23A', // Green for optimal
          '#409EFF'  // Blue for available
        ],
        borderWidth: 2,
        borderColor: '#fff'
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          position: 'right',
          labels: {
            usePointStyle: true,
            padding: 20,
            font: { size: 13 }
          }
        },
        tooltip: {
          backgroundColor: 'rgba(0, 0, 0, 0.8)',
          padding: 12,
          callbacks: {
            label: (context) => {
              const label = context.label || ''
              const value = context.parsed || 0
              const total = context.dataset.data.reduce((a, b) => a + b, 0)
              const percentage = ((value / total) * 100).toFixed(1)
              return `${label}: ${value} resources (${percentage}%)`
            }
          }
        }
      }
    }
  })
}

// Lifecycle
onMounted(() => {
  if (props.view === 'chart') {
    initChart()
  }
})

watch(() => [props.view, props.dateRange, props.tasks, props.resources], () => {
  if (props.view === 'chart') {
    initChart()
  }
}, { deep: true })
</script>

<style scoped lang="scss">
.resource-utilization {
  width: 100%;

  .chart-view {
    canvas {
      width: 100% !important;
      max-width: 500px;
      margin: 0 auto;
    }
  }

  .table-view {
    .resource-cell {
      display: flex;
      align-items: center;
      gap: 12px;

      .resource-info {
        .resource-name {
          font-weight: 600;
          color: #303133;
          font-size: 14px;
        }

        .resource-role {
          font-size: 12px;
          color: #909399;
          margin-top: 2px;
        }
      }
    }
  }

  .utilization-summary {
    display: flex;
    justify-content: space-around;
    margin-top: 20px;
    padding: 16px;
    background: #f5f7fa;
    border-radius: 6px;

    .summary-item {
      text-align: center;

      .summary-label {
        display: block;
        font-size: 12px;
        color: #909399;
        margin-bottom: 4px;
      }

      .summary-value {
        display: block;
        font-size: 18px;
        font-weight: 700;
        color: #303133;

        &.summary-danger {
          color: #F56C6C;
        }

        &.summary-warning {
          color: #E6A23C;
        }

        &.summary-success {
          color: #67C23A;
        }
      }
    }
  }
}

/* Responsive Design */
@media (max-width: 768px) {
  .resource-utilization {
    .utilization-summary {
      flex-direction: column;
      gap: 12px;
    }
  }
}
</style>
