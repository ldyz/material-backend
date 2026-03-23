<template>
  <div class="burndown-chart">
    <canvas ref="chartCanvas" :height="height"></canvas>

    <!-- Burndown Summary -->
    <div class="burndown-summary">
      <div class="burndown-metric">
        <div class="metric-label">Total Work</div>
        <div class="metric-value">{{ totalWork }} story points</div>
      </div>
      <div class="burndown-metric">
        <div class="metric-label">Remaining</div>
        <div class="metric-value metric-remaining">{{ remainingWork }} story points</div>
      </div>
      <div class="burndown-metric">
        <div class="metric-label">Completed</div>
        <div class="metric-value metric-completed">{{ completedWork }} story points</div>
      </div>
      <div class="burndown-metric">
        <div class="metric-label">Velocity</div>
        <div class="metric-value">{{ velocity }} pts/day</div>
      </div>
      <div class="burndown-metric">
        <div class="metric-label">Predicted Completion</div>
        <div class="metric-value" :class="predictionClass">
          {{ predictedCompletion }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { Line } from 'chart.js/auto'
import { format, addDays, differenceInBusinessDays } from 'date-fns'

/**
 * BurndownChart Component
 * Displays ideal vs actual burndown with completion prediction
 *
 * @props {String} sprint - Sprint identifier (current/last/all)
 * @props {Array} dateRange - Date range [start, end]
 * @props {Array} tasks - Array of tasks with story points and progress
 * @props {Number} height - Chart height in pixels
 */

const props = defineProps({
  sprint: {
    type: String,
    default: 'current',
    validator: (value) => ['current', 'last', 'all'].includes(value)
  },
  dateRange: {
    type: Array,
    default: () => [new Date(), new Date()]
  },
  tasks: {
    type: Array,
    default: () => []
  },
  height: {
    type: Number,
    default: 300
  }
})

// State
const chartCanvas = ref(null)
let chartInstance = null

// Computed
const totalWork = computed(() => {
  return props.tasks.reduce((sum, task) => sum + (task.storyPoints || 0), 0)
})

const remainingWork = computed(() => {
  return props.tasks.reduce((sum, task) => {
    const points = task.storyPoints || 0
    const remaining = 1 - (task.progress || 0) / 100
    return sum + (points * remaining)
  }, 0)
})

const completedWork = computed(() => {
  return props.tasks.reduce((sum, task) => {
    const points = task.storyPoints || 0
    const completed = (task.progress || 0) / 100
    return sum + (points * completed)
  }, 0)
})

const velocity = computed(() => {
  const startDate = props.dateRange[0]
  const endDate = props.dateRange[1]
  const daysWorked = differenceInBusinessDays(endDate, startDate) || 1

  return (completedWork.value / daysWorked).toFixed(1)
})

const predictedCompletion = computed(() => {
  if (parseFloat(velocity.value) <= 0) return 'Never'
  if (remainingWork.value <= 0) return 'Completed'

  const daysRemaining = Math.ceil(remainingWork.value / parseFloat(velocity.value))
  const completionDate = addDays(new Date(), daysRemaining)

  if (daysRemaining <= 0) return 'Completed'
  if (daysRemaining <= 3) return `${daysRemaining} days`
  if (daysRemaining <= 30) return format(completionDate, 'MMM d')
  return format(completionDate, 'MMM d, yyyy')
})

const predictionClass = computed(() => {
  const daysRemaining = Math.ceil(remainingWork.value / parseFloat(velocity.value))
  if (daysRemaining <= 0) return 'metric-success'
  if (daysRemaining <= 7) return 'metric-warning'
  return 'metric-info'
})

// Methods
const generateBurndownData = () => {
  const startDate = props.dateRange[0]
  const endDate = props.dateRange[1]
  const totalDays = differenceInBusinessDays(endDate, startDate) + 1

  const labels = []
  const idealData = []
  const actualData = []

  // Generate ideal burndown line (linear)
  for (let i = 0; i < totalDays; i++) {
    const date = addDays(startDate, i)
    labels.push(format(date, 'MMM d'))

    // Ideal: remaining work should decrease linearly
    const idealRemaining = Math.max(0, totalWork.value * (1 - i / totalDays))
    idealData.push(idealRemaining)
  }

  // Generate actual burndown based on task completion history
  // For demo purposes, we'll simulate actual data
  let cumulativeCompleted = 0
  for (let i = 0; i < totalDays; i++) {
    // Simulate work completion based on velocity
    const daysCompleted = i + 1
    const expectedCompletion = Math.min(totalWork.value, parseFloat(velocity.value) * daysCompleted)
    const actualRemaining = Math.max(0, totalWork.value - expectedCompletion)

    // Add some randomness to simulate real-world fluctuations
    const noise = (Math.random() - 0.5) * (totalWork.value * 0.05)
    actualData.push(Math.max(0, actualRemaining + noise))
  }

  // Add final actual point
  actualData[actualData.length - 1] = remainingWork.value

  return { labels, idealData, actualData }
}

const initChart = () => {
  if (!chartCanvas.value) return

  const ctx = chartCanvas.value.getContext('2d')
  const { labels, idealData, actualData } = generateBurndownData()

  if (chartInstance) {
    chartInstance.destroy()
  }

  chartInstance = new Line(ctx, {
    type: 'line',
    data: {
      labels,
      datasets: [
        {
          label: 'Ideal Burndown',
          data: idealData,
          borderColor: '#909399',
          backgroundColor: 'transparent',
          borderWidth: 2,
          borderDash: [5, 5],
          tension: 0,
          pointRadius: 0
        },
        {
          label: 'Actual Burndown',
          data: actualData,
          borderColor: '#409EFF',
          backgroundColor: 'rgba(64, 158, 255, 0.1)',
          borderWidth: 3,
          tension: 0.4,
          fill: true,
          pointBackgroundColor: '#409EFF',
          pointBorderColor: '#fff',
          pointBorderWidth: 2,
          pointRadius: 4,
          pointHoverRadius: 6
        }
      ]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      interaction: {
        intersect: false,
        mode: 'index'
      },
      plugins: {
        legend: {
          position: 'top',
          labels: {
            usePointStyle: true,
            padding: 15
          }
        },
        tooltip: {
          backgroundColor: 'rgba(0, 0, 0, 0.8)',
          padding: 12,
          titleFont: { size: 14 },
          bodyFont: { size: 13 },
          callbacks: {
            label: (context) => {
              return `${context.dataset.label}: ${Math.round(context.parsed.y)} story points`
            }
          }
        },
        annotation: {
          annotations: {
            todayLine: {
              type: 'line',
              xMin: labels.length - 1,
              xMax: labels.length - 1,
              borderColor: '#E6A23C',
              borderWidth: 2,
              borderDash: [6, 6],
              label: {
                display: true,
                content: 'Today',
                position: 'start',
                backgroundColor: '#E6A23C',
                font: { size: 11 }
              }
            }
          }
        }
      },
      scales: {
        y: {
          beginAtZero: true,
          reverse: true, // Burndown charts go down
          title: {
            display: true,
            text: 'Story Points',
            font: { size: 12, weight: '600' }
          },
          grid: {
            color: 'rgba(0, 0, 0, 0.05)'
          }
        },
        x: {
          grid: {
            display: false
          }
        }
      }
    }
  })
}

// Lifecycle
onMounted(() => {
  initChart()
})

watch(() => [props.sprint, props.dateRange, props.tasks], () => {
  initChart()
}, { deep: true })
</script>

<style scoped lang="scss">
.burndown-chart {
  width: 100%;

  canvas {
    width: 100% !important;
  }

  .burndown-summary {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
    gap: 16px;
    margin-top: 16px;
    padding-top: 16px;
    border-top: 1px solid #ebeef5;

    .burndown-metric {
      text-align: center;

      .metric-label {
        font-size: 12px;
        color: #909399;
        margin-bottom: 4px;
      }

      .metric-value {
        font-size: 16px;
        font-weight: 700;
        color: #303133;

        &.metric-remaining {
          color: #409EFF;
        }

        &.metric-completed {
          color: #67C23A;
        }

        &.metric-success {
          color: #67C23A;
        }

        &.metric-warning {
          color: #E6A23C;
        }

        &.metric-info {
          color: #909399;
        }
      }
    }
  }
}
</style>
