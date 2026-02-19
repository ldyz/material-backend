<template>
  <div class="earned-value-chart">
    <canvas ref="chartCanvas" :height="height"></canvas>

    <!-- EVM Metrics Summary -->
    <div class="evm-summary">
      <div class="evm-metric">
        <div class="metric-label">Planned Value (PV)</div>
        <div class="metric-value metric-pv">${{ formatNumber(pv) }}</div>
      </div>
      <div class="evm-metric">
        <div class="metric-label">Earned Value (EV)</div>
        <div class="metric-value metric-ev">${{ formatNumber(ev) }}</div>
      </div>
      <div class="evm-metric">
        <div class="metric-label">Actual Cost (AC)</div>
        <div class="metric-value metric-ac">${{ formatNumber(ac) }}</div>
      </div>
      <div class="evm-metric">
        <div class="metric-label">Schedule Performance (SPI)</div>
        <div class="metric-value" :class="spiClass">{{ spi.toFixed(2) }}</div>
      </div>
      <div class="evm-metric">
        <div class="metric-label">Cost Performance (CPI)</div>
        <div class="metric-value" :class="cpiClass">{{ cpi.toFixed(2) }}</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { Line } from 'chart.js/auto'
import { format, subDays, subMonths, parseISO } from 'date-fns'

/**
 * EarnedValueChart Component
 * Displays Earned Value Management (EVM) metrics visualization
 *
 * @props {String} period - Time period (week/month/quarter)
 * @props {Array} dateRange - Date range [start, end]
 * @props {Array} tasks - Array of tasks with progress and budget data
 * @props {Number} height - Chart height in pixels
 */

const props = defineProps({
  period: {
    type: String,
    default: 'month',
    validator: (value) => ['week', 'month', 'quarter'].includes(value)
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
const pv = computed(() => {
  // Planned Value: Budgeted cost of work scheduled
  return props.tasks.reduce((sum, task) => {
    const totalBudget = task.budget || 0
    const taskDuration = getTaskDuration(task)
    const elapsedDays = getElapsedDays(task)
    const plannedPercent = Math.min(elapsedDays / taskDuration, 1)
    return sum + (totalBudget * plannedPercent)
  }, 0)
})

const ev = computed(() => {
  // Earned Value: Budgeted cost of work performed
  return props.tasks.reduce((sum, task) => {
    const totalBudget = task.budget || 0
    const progress = (task.progress || 0) / 100
    return sum + (totalBudget * progress)
  }, 0)
})

const ac = computed(() => {
  // Actual Cost: Actual cost of work performed
  return props.tasks.reduce((sum, task) => {
    return sum + (task.actualCost || 0)
  }, 0)
})

const spi = computed(() => {
  // Schedule Performance Index: EV / PV
  return pv.value > 0 ? ev.value / pv.value : 0
})

const cpi = computed(() => {
  // Cost Performance Index: EV / AC
  return ac.value > 0 ? ev.value / ac.value : 0
})

const spiClass = computed(() => {
  if (spi.value >= 1) return 'metric-success'
  if (spi.value >= 0.9) return 'metric-warning'
  return 'metric-danger'
})

const cpiClass = computed(() => {
  if (cpi.value >= 1) return 'metric-success'
  if (cpi.value >= 0.9) return 'metric-warning'
  return 'metric-danger'
})

// Methods
const getTaskDuration = (task) => {
  if (!task.startDate || !task.endDate) return 1
  const start = parseISO(task.startDate)
  const end = parseISO(task.endDate)
  return Math.max(1, Math.ceil((end - start) / (1000 * 60 * 60 * 24)))
}

const getElapsedDays = (task) => {
  if (!task.startDate) return 0
  const start = parseISO(task.startDate)
  const now = new Date()
  return Math.max(0, Math.ceil((now - start) / (1000 * 60 * 60 * 24)))
}

const formatNumber = (num) => {
  return Math.round(num).toLocaleString()
}

const generateTimeSeriesData = () => {
  const labels = []
  const pvData = []
  const evData = []
  const acData = []

  const startDate = props.dateRange[0]
  const endDate = props.dateRange[1]
  const daysDiff = Math.ceil((endDate - startDate) / (1000 * 60 * 60 * 24))

  const interval = props.period === 'week' ? 1 :
                  props.period === 'month' ? 7 : 30

  for (let i = 0; i <= daysDiff; i += interval) {
    const currentDate = new Date(startDate)
    currentDate.setDate(startDate.getDate() + i)
    labels.push(format(currentDate, 'MMM d'))

    // Calculate cumulative values up to this date
    const pvAtDate = props.tasks.reduce((sum, task) => {
      const totalBudget = task.budget || 0
      const taskDuration = getTaskDuration(task)
      const elapsedDays = Math.max(0, Math.ceil((currentDate - parseISO(task.startDate)) / (1000 * 60 * 60 * 24)))
      const plannedPercent = Math.min(Math.max(elapsedDays / taskDuration, 0), 1)
      return sum + (totalBudget * plannedPercent)
    }, 0)

    const evAtDate = props.tasks.reduce((sum, task) => {
      const totalBudget = task.budget || 0
      const taskEndDate = parseISO(task.endDate)
      if (currentDate < taskEndDate) {
        // Assume linear progress for projection
        const taskDuration = getTaskDuration(task)
        const elapsedDays = Math.max(0, Math.ceil((currentDate - parseISO(task.startDate)) / (1000 * 60 * 60 * 24)))
        const progressPercent = Math.min(Math.max(elapsedDays / taskDuration, 0), 1)
        return sum + (totalBudget * progressPercent)
      } else {
        // Task completed by this date
        return sum + totalBudget
      }
    }, 0)

    pvData.push(pvAtDate)
    evData.push(evAtDate)
    acData.push(evAtDate * 1.1) // Simulated AC (10% over budget)
  }

  return { labels, pvData, evData, acData }
}

const initChart = () => {
  if (!chartCanvas.value) return

  const ctx = chartCanvas.value.getContext('2d')
  const { labels, pvData, evData, acData } = generateTimeSeriesData()

  if (chartInstance) {
    chartInstance.destroy()
  }

  chartInstance = new Line(ctx, {
    type: 'line',
    data: {
      labels,
      datasets: [
        {
          label: 'Planned Value (PV)',
          data: pvData,
          borderColor: '#409EFF',
          backgroundColor: 'rgba(64, 158, 255, 0.1)',
          borderWidth: 2,
          tension: 0.4,
          fill: true
        },
        {
          label: 'Earned Value (EV)',
          data: evData,
          borderColor: '#67C23A',
          backgroundColor: 'rgba(103, 194, 58, 0.1)',
          borderWidth: 2,
          tension: 0.4,
          fill: true
        },
        {
          label: 'Actual Cost (AC)',
          data: acData,
          borderColor: '#E6A23C',
          backgroundColor: 'rgba(230, 162, 60, 0.1)',
          borderWidth: 2,
          tension: 0.4,
          fill: true
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
              return `${context.dataset.label}: $${formatNumber(context.parsed.y)}`
            }
          }
        }
      },
      scales: {
        y: {
          beginAtZero: true,
          ticks: {
            callback: (value) => '$' + formatNumber(value)
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

watch(() => [props.period, props.dateRange, props.tasks], () => {
  initChart()
}, { deep: true })
</script>

<style scoped lang="scss">
.earned-value-chart {
  width: 100%;

  canvas {
    width: 100% !important;
  }

  .evm-summary {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
    gap: 16px;
    margin-top: 16px;
    padding-top: 16px;
    border-top: 1px solid #ebeef5;

    .evm-metric {
      text-align: center;

      .metric-label {
        font-size: 12px;
        color: #909399;
        margin-bottom: 4px;
      }

      .metric-value {
        font-size: 18px;
        font-weight: 700;
        color: #303133;

        &.metric-pv {
          color: #409EFF;
        }

        &.metric-ev {
          color: #67C23A;
        }

        &.metric-ac {
          color: #E6A23C;
        }

        &.metric-success {
          color: #67C23A;
        }

        &.metric-warning {
          color: #E6A23C;
        }

        &.metric-danger {
          color: #F56C6C;
        }
      }
    }
  }
}
</style>
