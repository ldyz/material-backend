<template>
  <div class="resource-histogram">
    <div class="histogram-header">
      <div class="header-left">
        <h3>{{ t('gantt.resourceHistogram.title') }}</h3>
      </div>
      <div class="header-right">
        <el-select
          v-model="selectedResourceId"
          :placeholder="t('gantt.resourceHistogram.selectResource')"
          @change="handleResourceChange"
        >
          <el-option
            v-for="resource in resources"
            :key="resource.id"
            :label="resource.name"
            :value="resource.id"
          />
        </el-select>
        <el-radio-group v-model="viewMode" @change="handleViewModeChange">
          <el-radio-button value="daily">{{ t('gantt.resourceHistogram.daily') }}</el-radio-button>
          <el-radio-button value="weekly">{{ t('gantt.resourceHistogram.weekly') }}</el-radio-button>
        </el-radio-group>
      </div>
    </div>

    <div v-if="selectedResource" class="histogram-content">
      <!-- Capacity Legend -->
      <div class="capacity-legend">
        <div class="legend-item">
          <span class="legend-color under-capacity"></span>
          <span>&lt; {{ selectedResource.capacity || 100 }}%</span>
        </div>
        <div class="legend-item">
          <span class="legend-color at-capacity"></span>
          <span>= {{ selectedResource.capacity || 100 }}%</span>
        </div>
        <div class="legend-item">
          <span class="legend-color over-capacity"></span>
          <span>&gt; {{ selectedResource.capacity || 100 }}%</span>
        </div>
      </div>

      <!-- Histogram Chart -->
      <div class="histogram-chart" ref="chartRef">
        <svg :width="chartWidth" :height="chartHeight">
          <!-- Background grid -->
          <g class="grid">
            <line
              v-for="(line, index) in gridLines"
              :key="'grid-' + index"
              :x1="line.x"
              :y1="0"
              :x2="line.x"
              :y2="chartHeight - 30"
              :stroke="index % 2 === 0 ? '#e4e7ed' : '#f5f7fa'"
              stroke-width="1"
            />
          </g>

          <!-- Capacity reference line -->
          <line
            :x1="0"
            :y1="capacityY"
            :x2="chartWidth"
            :y2="capacityY"
            stroke="#E6A23C"
            stroke-width="2"
            stroke-dasharray="5,5"
          />
          <text
            :x="chartWidth - 10"
            :y="capacityY - 5"
            text-anchor="end"
            fill="#E6A23C"
            font-size="12"
          >
            {{ selectedResource.capacity || 100 }}%
          </text>

          <!-- Allocation bars -->
          <g class="bars">
            <rect
              v-for="(bar, index) in allocationBars"
              :key="'bar-' + index"
              :x="bar.x"
              :y="bar.y"
              :width="bar.width"
              :height="bar.height"
              :fill="bar.fill"
              @mouseenter="handleBarMouseEnter($event, bar)"
              @mouseleave="handleBarMouseLeave"
              @click="handleBarClick(bar)"
              class="allocation-bar"
            >
              <title>{{ bar.tooltip }}</title>
            </rect>
          </g>

          <!-- Overallocation indicators -->
          <g v-if="hasOverallocation" class="overallocations">
            <polygon
              v-for="(indicator, index) in overallocationIndicators"
              :key="'over-' + index"
              :points="indicator.points"
              fill="#F56C6C"
              @click="handleOverallocationClick(indicator)"
              class="overallocation-indicator"
            >
              <title>{{ indicator.tooltip }}</title>
            </polygon>
          </g>

          <!-- Date labels -->
          <g class="labels">
            <text
              v-for="(label, index) in dateLabels"
              :key="'label-' + index"
              :x="label.x"
              :y="chartHeight - 10"
              text-anchor="middle"
              fill="#606266"
              font-size="11"
            >
              {{ label.text }}
            </text>
          </g>
        </svg>
      </div>

      <!-- Tooltip -->
      <div
        v-if="tooltip.visible"
        class="histogram-tooltip"
        :style="{
          left: tooltip.x + 'px',
          top: tooltip.y + 'px'
        }"
      >
        <div class="tooltip-header">{{ tooltip.date }}</div>
        <div class="tooltip-body">
          <div>{{ t('gantt.resourceHistogram.allocation') }}: {{ tooltip.allocation }}%</div>
          <div>{{ t('gantt.resourceHistogram.capacity') }}: {{ tooltip.capacity }}%</div>
          <div v-if="tooltip.overallocation" class="overallocation-text">
            {{ t('gantt.resourceHistogram.overallocation') }}: +{{ tooltip.overallocation }}%
          </div>
          <div v-if="tooltip.tasks && tooltip.tasks.length > 0" class="tooltip-tasks">
            <div class="task-list-title">{{ t('gantt.resourceHistogram.assignedTasks') }}:</div>
            <div
              v-for="task in tooltip.tasks"
              :key="task.id"
              class="task-item"
            >
              <span class="task-name">{{ task.name }}</span>
              <span class="task-allocation">{{ task.allocation }}%</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="empty-state">
      <el-empty :description="t('gantt.resourceHistogram.noResourceSelected')" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

// Props
const props = defineProps({
  tasks: {
    type: Array,
    default: () => []
  },
  resources: {
    type: Array,
    default: () => []
  },
  startDate: {
    type: String,
    required: true
  },
  endDate: {
    type: String,
    required: true
  }
})

// Emits
const emit = defineEmits(['periodClick', 'taskClick'])

// State
const chartRef = ref(null)
const chartWidth = ref(800)
const chartHeight = ref(300)
const selectedResourceId = ref(null)
const viewMode = ref('daily')

const tooltip = ref({
  visible: false,
  x: 0,
  y: 0,
  date: '',
  allocation: 0,
  capacity: 0,
  overallocation: 0,
  tasks: []
})

// Computed
const selectedResource = computed(() => {
  if (!selectedResourceId.value || !props.resources) {
    return null
  }
  return props.resources.find(r => r.id === selectedResourceId.value)
})

const chartDays = computed(() => {
  const start = new Date(props.startDate)
  const end = new Date(props.endDate)
  return Math.ceil((end - start) / (1000 * 60 * 60 * 24))
})

const barWidth = computed(() => {
  if (viewMode.value === 'daily') {
    return Math.max(10, chartWidth.value / chartDays.value - 2)
  } else {
    // Weekly view
    const weeks = Math.ceil(chartDays.value / 7)
    return Math.max(20, chartWidth.value / weeks - 4)
  }
})

const capacityY = computed(() => {
  const capacity = selectedResource.value?.capacity || 100
  return chartHeight.value - 30 - (capacity / 150) * (chartHeight.value - 30)
})

const allocationData = computed(() => {
  if (!selectedResource.value || !props.tasks) {
    return []
  }

  const data = []
  const start = new Date(props.startDate)
  const end = new Date(props.endDate)
  const days = Math.ceil((end - start) / (1000 * 60 * 60 * 24))

  for (let i = 0; i <= days; i++) {
    const currentDate = new Date(start)
    currentDate.setDate(currentDate.getDate() + i)
    const dateKey = currentDate.toISOString().split('T')[0]

    // Calculate total allocation for this date
    let totalAllocation = 0
    const assignedTasks = []

    props.tasks.forEach(task => {
      const taskStart = new Date(task.start_date)
      const taskEnd = new Date(task.end_date)

      if (currentDate >= taskStart && currentDate <= taskEnd) {
        const assignment = task.resources?.find(
          r => r.resource_id === selectedResourceId.value
        )

        if (assignment) {
          const allocation = assignment.units || 100
          totalAllocation += allocation
          assignedTasks.push({
            id: task.id,
            name: task.name,
            allocation
          })
        }
      }
    })

    data.push({
      date: dateKey,
      allocation: Math.round(totalAllocation),
      capacity: selectedResource.value.capacity || 100,
      tasks: assignedTasks,
      isOverallocated: totalAllocation > (selectedResource.value.capacity || 100)
    })
  }

  // Group by week if in weekly view
  if (viewMode.value === 'weekly') {
    const weeklyData = []
    for (let i = 0; i < data.length; i += 7) {
      const weekData = data.slice(i, i + 7)
      const maxAllocation = Math.max(...weekData.map(d => d.allocation))
      const allTasks = [...new Set(weekData.flatMap(d => d.tasks.map(t => t.id)))]
        .map(id => weekData.flatMap(d => d.tasks).find(t => t.id === id))

      weeklyData.push({
        date: weekData[0].date + ' - ' + weekData[weekData.length - 1].date,
        allocation: maxAllocation,
        capacity: weekData[0].capacity,
        tasks: allTasks,
        isOverallocated: weekData.some(d => d.isOverallocated)
      })
    }
    return weeklyData
  }

  return data
})

const gridLines = computed(() => {
  const lines = []
  const count = viewMode.value === 'daily' ? chartDays.value : Math.ceil(chartDays.value / 7)

  for (let i = 0; i <= count; i++) {
    lines.push({
      x: i * (barWidth.value + (viewMode.value === 'daily' ? 2 : 4))
    })
  }

  return lines
})

const allocationBars = computed(() => {
  return allocationData.value.map((data, index) => {
    const maxHeight = chartHeight.value - 30
    const barHeight = Math.min(maxHeight, (data.allocation / 150) * maxHeight)
    const y = maxHeight - barHeight

    let fill = '#67C23A' // Green
    if (data.allocation >= data.capacity) {
      fill = '#E6A23C' // Orange
    }
    if (data.allocation > data.capacity) {
      fill = '#F56C6C' // Red
    }

    return {
      x: index * (barWidth.value + (viewMode.value === 'daily' ? 2 : 4)),
      y,
      width: barWidth.value,
      height: barHeight,
      fill,
      data,
      tooltip: `${data.date}: ${data.allocation}%`
    }
  })
})

const overallocationIndicators = computed(() => {
  return allocationBars.value
    .filter(bar => bar.data.isOverallocated)
    .map(bar => {
      const x = bar.x + bar.width / 2
      const y = 10
      return {
        points: `${x},${y} ${x - 6},${y - 8} ${x + 6},${y - 8}`,
        data: bar.data,
        tooltip: `${t('gantt.resourceHistogram.overallocation')}: ${bar.data.allocation}%`
      }
    })
})

const hasOverallocation = computed(() => {
  return allocationData.value.some(d => d.isOverallocated)
})

const dateLabels = computed(() => {
  const labels = []
  const data = allocationData.value
  const skipEvery = Math.ceil(data.length / 10) // Show ~10 labels max

  data.forEach((d, index) => {
    if (index % skipEvery === 0 || index === data.length - 1) {
      const x = index * (barWidth.value + (viewMode.value === 'daily' ? 2 : 4)) + barWidth.value / 2
      let text = d.date
      if (viewMode.value === 'weekly') {
        const dates = d.date.split(' - ')
        text = dates[0].substring(5) // Show MM-DD
      } else {
        text = d.date.substring(5) // Show MM-DD
      }

      labels.push({ x, text })
    }
  })

  return labels
})

// Methods
const handleResourceChange = () => {
  // Recalculate chart
}

const handleViewModeChange = () => {
  // Recalculate chart
}

const handleBarMouseEnter = (event, bar) => {
  tooltip.value = {
    visible: true,
    x: event.clientX + 10,
    y: event.clientY - 10,
    date: bar.data.date,
    allocation: bar.data.allocation,
    capacity: bar.data.capacity,
    overallocation: bar.data.allocation > bar.data.capacity
      ? bar.data.allocation - bar.data.capacity
      : 0,
    tasks: bar.data.tasks
  }
}

const handleBarMouseLeave = () => {
  tooltip.value.visible = false
}

const handleBarClick = (bar) => {
  emit('periodClick', {
    date: bar.data.date,
    allocation: bar.data.allocation,
    tasks: bar.data.tasks
  })
}

const handleOverallocationClick = (indicator) => {
  emit('periodClick', {
    date: indicator.data.date,
    allocation: indicator.data.allocation,
    tasks: indicator.data.tasks,
    isOverallocation: true
  })
}

const updateChartWidth = () => {
  if (chartRef.value) {
    chartWidth.value = chartRef.value.offsetWidth
  }
}

// Lifecycle
onMounted(() => {
  updateChartWidth()
  window.addEventListener('resize', updateChartWidth)

  // Select first resource by default
  if (props.resources && props.resources.length > 0) {
    selectedResourceId.value = props.resources[0].id
  }
})

// Watch for data changes
watch(() => [props.startDate, props.endDate], () => {
  nextTick(() => {
    updateChartWidth()
  })
})
</script>

<style scoped>
.resource-histogram {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.histogram-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid #e4e7ed;
}

.header-left h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 500;
}

.header-right {
  display: flex;
  gap: 12px;
  align-items: center;
}

.histogram-content {
  flex: 1;
  padding: 16px;
  overflow: auto;
}

.capacity-legend {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
  font-size: 12px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.legend-color {
  width: 16px;
  height: 16px;
  border-radius: 2px;
}

.legend-color.under-capacity {
  background: #67C23A;
}

.legend-color.at-capacity {
  background: #E6A23C;
}

.legend-color.over-capacity {
  background: #F56C6C;
}

.histogram-chart {
  position: relative;
  margin: 16px 0;
}

.allocation-bar {
  cursor: pointer;
  transition: opacity 0.2s;
}

.allocation-bar:hover {
  opacity: 0.8;
}

.overallocation-indicator {
  cursor: pointer;
  transition: transform 0.2s;
}

.overallocation-indicator:hover {
  transform: scale(1.2);
}

.histogram-tooltip {
  position: fixed;
  background: white;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  padding: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  z-index: 9999;
  min-width: 200px;
}

.tooltip-header {
  font-weight: 500;
  margin-bottom: 8px;
  color: #303133;
}

.tooltip-body {
  font-size: 12px;
  color: #606266;
}

.tooltip-body > div {
  margin-bottom: 4px;
}

.overallocation-text {
  color: #F56C6C;
  font-weight: 500;
}

.tooltip-tasks {
  margin-top: 8px;
  padding-top: 8px;
  border-top: 1px solid #e4e7ed;
}

.task-list-title {
  font-weight: 500;
  margin-bottom: 4px;
}

.task-item {
  display: flex;
  justify-content: space-between;
  padding: 2px 0;
}

.task-name {
  flex: 1;
}

.task-allocation {
  color: #909399;
  margin-left: 8px;
}

.empty-state {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
