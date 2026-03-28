<template>
  <div class="dashboard-view">
    <!-- Dashboard Controls -->
    <div class="dashboard-view__controls">
      <div class="controls-left">
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          range-separator="to"
          start-placeholder="Start date"
          end-placeholder="End date"
          @change="handleDateRangeChange"
        />

        <el-select
          v-model="selectedProject"
          placeholder="All Projects"
          clearable
          @change="handleProjectChange"
        >
          <el-option
            v-for="project in projects"
            :key="project.id"
            :label="project.name"
            :value="project.id"
          />
        </el-select>
      </div>

      <div class="controls-right">
        <el-button
          :icon="Refresh"
          @click="refreshDashboard"
          :loading="loading"
        >
          Refresh
        </el-button>

        <el-button
          :icon="Download"
          @click="exportDashboard"
        >
          Export
        </el-button>
      </div>
    </div>

    <!-- KPI Summary Cards -->
    <div class="dashboard-view__kpi">
      <StatCard
        v-for="stat in kpiStats"
        :key="stat.id"
        :icon="stat.icon"
        :title="stat.title"
        :value="stat.value"
        :trend="stat.trend"
        :trend-value="stat.trendValue"
        :color="stat.color"
        :sparkline-data="stat.sparklineData"
        @click="handleStatClick(stat)"
      />
    </div>

    <!-- Charts Grid -->
    <div class="dashboard-view__charts">
      <!-- Earned Value Chart -->
      <div class="chart-card">
        <div class="chart-card__header">
          <h3>Earned Value Analysis</h3>
          <el-button-group size="small">
            <el-button
              :type="evmPeriod === 'week' ? 'primary' : ''"
              @click="evmPeriod = 'week'"
            >
              Week
            </el-button>
            <el-button
              :type="evmPeriod === 'month' ? 'primary' : ''"
              @click="evmPeriod = 'month'"
            >
              Month
            </el-button>
            <el-button
              :type="evmPeriod === 'quarter' ? 'primary' : ''"
              @click="evmPeriod = 'quarter'"
            >
              Quarter
            </el-button>
          </el-button-group>
        </div>
        <div class="chart-card__content">
          <EarnedValueChart
            :period="evmPeriod"
            :date-range="dateRange"
            :tasks="tasks"
            :height="300"
          />
        </div>
      </div>

      <!-- Burndown Chart -->
      <div class="chart-card">
        <div class="chart-card__header">
          <h3>Burndown Chart</h3>
          <el-select
            v-model="burndownSprint"
            size="small"
            style="width: 150px"
          >
            <el-option label="Current Sprint" value="current" />
            <el-option label="Last Sprint" value="last" />
            <el-option label="All Time" value="all" />
          </el-select>
        </div>
        <div class="chart-card__content">
          <BurndownChart
            :sprint="burndownSprint"
            :date-range="dateRange"
            :tasks="tasks"
            :height="300"
          />
        </div>
      </div>

      <!-- Resource Utilization -->
      <div class="chart-card chart-card--full">
        <div class="chart-card__header">
          <h3>Resource Utilization</h3>
          <el-button-group size="small">
            <el-button
              :type="resourceView === 'chart' ? 'primary' : ''"
              @click="resourceView = 'chart'"
            >
              Chart
            </el-button>
            <el-button
              :type="resourceView === 'table' ? 'primary' : ''"
              @click="resourceView = 'table'"
            >
              Table
            </el-button>
          </el-button-group>
        </div>
        <div class="chart-card__content">
          <ResourceUtilization
            :view="resourceView"
            :date-range="dateRange"
            :tasks="tasks"
            :resources="resources"
            :height="300"
          />
        </div>
      </div>

      <!-- Milestone Tracker -->
      <div class="chart-card">
        <div class="chart-card__header">
          <h3>Milestone Tracker</h3>
          <el-tag :type="milestoneStatus.type">
            {{ milestoneStatus.text }}
          </el-tag>
        </div>
        <div class="chart-card__content">
          <MilestoneTracker
            :date-range="dateRange"
            :milestones="milestones"
            :height="300"
          />
        </div>
      </div>

      <!-- Task Distribution -->
      <div class="chart-card">
        <div class="chart-card__header">
          <h3>Task Distribution</h3>
          <el-select
            v-model="distributionBy"
            size="small"
            style="width: 150px"
          >
            <el-option label="By Status" value="status" />
            <el-option label="By Priority" value="priority" />
            <el-option label="By Assignee" value="assignee" />
          </el-select>
        </div>
        <div class="chart-card__content">
          <div ref="distributionChart" style="height: 300px"></div>
        </div>
      </div>
    </div>

    <!-- Recent Activity -->
    <div class="dashboard-view__activity">
      <div class="activity-card">
        <div class="activity-card__header">
          <h3>Recent Activity</h3>
          <el-button text @click="viewAllActivity">View All</el-button>
        </div>
        <div class="activity-card__content">
          <el-timeline>
            <el-timeline-item
              v-for="activity in recentActivity"
              :key="activity.id"
              :timestamp="activity.timestamp"
              :color="activity.color"
            >
              <div class="activity-item">
                <el-avatar :size="32" :src="activity.avatar">
                  {{ activity.user.charAt(0) }}
                </el-avatar>
                <div class="activity-details">
                  <p class="activity-text">
                    <strong>{{ activity.user }}</strong>
                    {{ activity.action }}
                    <strong>{{ activity.task }}</strong>
                  </p>
                  <p class="activity-meta">{{ activity.time }}</p>
                </div>
              </div>
            </el-timeline-item>
          </el-timeline>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { Refresh, Download } from '@element-plus/icons-vue'
import { Chart, registerables } from 'chart.js'
import StatCard from '../dashboard/StatCard.vue'
import EarnedValueChart from '../dashboard/EarnedValueChart.vue'
import BurndownChart from '../dashboard/BurndownChart.vue'
import ResourceUtilization from '../dashboard/ResourceUtilization.vue'
import MilestoneTracker from '../dashboard/MilestoneTracker.vue'

/**
 * DashboardView Component
 * Main dashboard view with KPI cards and charts
 *
 * @props {Array} tasks - Array of all tasks
 * @props {Array} projects - Array of projects
 * @props {Array} resources - Array of resources/team members
 * @props {Array} milestones - Array of milestones
 * @props {Array} activity - Array of recent activity
 *
 * @emits {Date} date-range-change - Emitted when date range changes
 * @emits {Number} project-change - Emitted when project filter changes
 * @emits {String} stat-click - Emitted when a stat card is clicked
 */

Chart.register(...registerables)

const props = defineProps({
  tasks: {
    type: Array,
    default: () => []
  },
  projects: {
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
  },
  activity: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['date-range-change', 'project-change', 'stat-click', 'refresh', 'export'])

// State
const loading = ref(false)
const dateRange = ref([
  new Date(Date.now() - 30 * 24 * 60 * 60 * 1000), // 30 days ago
  new Date()
])
const selectedProject = ref(null)
const evmPeriod = ref('month')
const burndownSprint = ref('current')
const resourceView = ref('chart')
const distributionBy = ref('status')
const distributionChart = ref(null)
let distributionChartInstance = null

// Computed
const kpiStats = computed(() => {
  const total = props.tasks.length
  const completed = props.tasks.filter(t => t.status === 'completed').length
  const inProgress = props.tasks.filter(t => t.status === 'in_progress').length
  const notStarted = props.tasks.filter(t => t.status === 'not_started').length
  const overdue = props.tasks.filter(t => t.isOverdue).length
  const criticalPath = props.tasks.filter(t => t.isOnCriticalPath).length
  const overallProgress = Math.round(
    props.tasks.reduce((sum, t) => sum + (t.progress || 0), 0) / (total || 1)
  )

  return [
    {
      id: 'total',
      icon: 'List',
      title: 'Total Tasks',
      value: total,
      trend: total > 0 ? 'up' : 'neutral',
      trendValue: '+5%',
      color: '#409EFF',
      sparklineData: generateSparklineData([10, 12, 15, 14, 18, 20, total])
    },
    {
      id: 'completed',
      icon: 'CircleCheck',
      title: 'Completed',
      value: completed,
      trend: 'up',
      trendValue: '+12%',
      color: '#67C23A',
      sparklineData: generateSparklineData([5, 8, 10, 12, 15, 18, completed])
    },
    {
      id: 'in-progress',
      icon: 'Loading',
      title: 'In Progress',
      value: inProgress,
      trend: 'up',
      trendValue: '+3%',
      color: '#E6A23C',
      sparklineData: generateSparklineData([8, 10, 9, 11, 12, 13, inProgress])
    },
    {
      id: 'overdue',
      icon: 'Warning',
      title: 'Overdue',
      value: overdue,
      trend: overdue > 0 ? 'up' : 'down',
      trendValue: overdue > 0 ? '+2' : '-5',
      color: '#F56C6C',
      sparklineData: generateSparklineData([0, 1, 2, 1, 3, 2, overdue])
    },
    {
      id: 'critical-path',
      icon: 'Guide',
      title: 'Critical Path',
      value: criticalPath,
      trend: 'neutral',
      trendValue: '0%',
      color: '#909399',
      sparklineData: generateSparklineData([5, 5, 6, 5, 5, 6, criticalPath])
    },
    {
      id: 'progress',
      icon: 'PieChart',
      title: 'Overall Progress',
      value: `${overallProgress}%`,
      trend: 'up',
      trendValue: '+8%',
      color: overallProgress >= 80 ? '#67C23A' : overallProgress >= 50 ? '#E6A23C' : '#F56C6C',
      sparklineData: generateSparklineData([30, 40, 45, 50, 55, 60, overallProgress])
    }
  ]
})

const recentActivity = computed(() => {
  return props.activity.slice(0, 5)
})

const milestoneStatus = computed(() => {
  const completed = props.milestones.filter(m => m.status === 'completed').length
  const total = props.milestones.length

  if (completed === total) {
    return { type: 'success', text: 'All Milestones Completed' }
  } else if (completed === 0) {
    return { type: 'info', text: 'Milestones Pending' }
  } else {
    return { type: 'warning', text: `${completed}/${total} Milestones Completed` }
  }
})

// Methods
const generateSparklineData = (values) => {
  return values.map(v => ({ value: v }))
}

const handleDateRangeChange = (value) => {
  emit('date-range-change', value)
}

const handleProjectChange = (value) => {
  emit('project-change', value)
}

const handleStatClick = (stat) => {
  emit('stat-click', stat)
}

const refreshDashboard = async () => {
  loading.value = true
  try {
    emit('refresh')
    await new Promise(resolve => setTimeout(resolve, 500))
  } finally {
    loading.value = false
  }
}

const exportDashboard = () => {
  emit('export', {
    dateRange: dateRange.value,
    project: selectedProject.value
  })
}

const viewAllActivity = () => {
  // Navigate to activity view
}

const updateDistributionChart = () => {
  if (!distributionChart.value) return

  if (distributionChartInstance) {
    distributionChartInstance.destroy()
  }

  const ctx = distributionChart.value.getContext('2d')

  // Prepare data based on distributionBy
  let labels = []
  let data = []
  let colors = []

  if (distributionBy.value === 'status') {
    const statusMap = {
      completed: { label: 'Completed', color: '#67C23A' },
      in_progress: { label: 'In Progress', color: '#E6A23C' },
      not_started: { label: 'Not Started', color: '#909399' }
    }

    Object.entries(statusMap).forEach(([key, { label, color }]) => {
      labels.push(label)
      colors.push(color)
      data.push(props.tasks.filter(t => t.status === key).length)
    })
  } else if (distributionBy.value === 'priority') {
    const priorityMap = {
      high: { label: 'High', color: '#F56C6C' },
      medium: { label: 'Medium', color: '#E6A23C' },
      low: { label: 'Low', color: '#409EFF' }
    }

    Object.entries(priorityMap).forEach(([key, { label, color }]) => {
      labels.push(label)
      colors.push(color)
      data.push(props.tasks.filter(t => t.priority === key).length)
    })
  } else if (distributionBy.value === 'assignee') {
    const assigneeCounts = {}
    props.tasks.forEach(task => {
      const assignee = task.assignee || 'Unassigned'
      assigneeCounts[assignee] = (assigneeCounts[assignee] || 0) + 1
    })

    Object.keys(assigneeCounts).forEach((assignee, index) => {
      labels.push(assignee)
      colors.push(`hsl(${(index * 360) / Object.keys(assigneeCounts).length}, 70%, 50%)`)
      data.push(assigneeCounts[assignee])
    })
  }

  distributionChartInstance = new Chart(ctx, {
    type: 'doughnut',
    data: {
      labels,
      datasets: [{
        data,
        backgroundColor: colors,
        borderWidth: 2,
        borderColor: '#fff'
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          position: 'right'
        },
        tooltip: {
          callbacks: {
            label: (context) => {
              const label = context.label || ''
              const value = context.parsed || 0
              const total = context.dataset.data.reduce((a, b) => a + b, 0)
              const percentage = ((value / total) * 100).toFixed(1)
              return `${label}: ${value} (${percentage}%)`
            }
          }
        }
      }
    }
  })
}

// Lifecycle
onMounted(() => {
  updateDistributionChart()
})

watch(distributionBy, () => {
  updateDistributionChart()
})

watch(() => props.tasks, () => {
  updateDistributionChart()
}, { deep: true })
</script>

<style scoped lang="scss">
.dashboard-view {
  padding: 16px;
  background: #f5f7fa;
  min-height: 100vh;

  &__controls {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
    padding: 16px;
    background: #fff;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

    .controls-left,
    .controls-right {
      display: flex;
      gap: 12px;
    }
  }

  &__kpi {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 16px;
    margin-bottom: 24px;
  }

  &__charts {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 16px;
    margin-bottom: 24px;

    .chart-card {
      background: #fff;
      border-radius: 8px;
      padding: 16px;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

      &--full {
        grid-column: 1 / -1;
      }

      &__header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 16px;

        h3 {
          margin: 0;
          font-size: 16px;
          font-weight: 600;
          color: #303133;
        }
      }

      &__content {
        position: relative;
      }
    }
  }

  &__activity {
    .activity-card {
      background: #fff;
      border-radius: 8px;
      padding: 16px;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

      &__header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 16px;

        h3 {
          margin: 0;
          font-size: 16px;
          font-weight: 600;
          color: #303133;
        }
      }

      .activity-item {
        display: flex;
        gap: 12px;

        .activity-details {
          flex: 1;

          .activity-text {
            margin: 0 0 4px 0;
            font-size: 14px;
            color: #303133;
          }

          .activity-meta {
            margin: 0;
            font-size: 12px;
            color: #909399;
          }
        }
      }
    }
  }
}

/* Responsive Design */
@media (max-width: 1200px) {
  .dashboard-view__charts {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .dashboard-view {
    padding: 12px;

    &__controls {
      flex-direction: column;
      gap: 12px;

      .controls-left,
      .controls-right {
        width: 100%;
        justify-content: center;
        flex-wrap: wrap;
      }
    }

    &__kpi {
      grid-template-columns: repeat(2, 1fr);
    }
  }
}
</style>
