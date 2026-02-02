<template>
  <!-- 调试信息 -->
  <div style="padding: 10px; background: #fff3cd; margin: 10px 0; border-radius: 4px;">
    <strong>ResourceView 调试:</strong>
    projectId={{ projectId }} |
    tasks数量={{ tasks.length }} |
    timelineDays数量={{ timelineDays.length }} |
    dayWidth={{ dayWidth }}
  </div>

  <div class="resource-view" v-if="tasks.length > 0">
    <!-- 资源类型切换 -->
    <div class="resource-view-header">
      <div class="view-title">
        <el-icon><TrendCharts /></el-icon>
        <span>资源需求视图</span>
      </div>
      <div class="view-controls">
        <el-radio-group v-model="activeResourceType" size="small" @change="handleResourceTypeChange">
          <el-radio-button label="labor">人力</el-radio-button>
          <el-radio-button label="equipment">工机具</el-radio-button>
          <el-radio-button label="material">材料</el-radio-button>
        </el-radio-group>
      </div>
    </div>

    <!-- 日期范围 -->
    <div class="resource-timeline-container" ref="timelineContainer">
      <svg :width="timelineWidth" :height="chartHeight" class="resource-timeline">
        <!-- 背景网格 -->
        <defs>
          <pattern
            :id="`grid-${activeResourceType}`"
            :width="dayWidth"
            :height="chartHeight"
            patternUnits="userSpaceOnUse"
          >
            <line
              :x1="dayWidth"
              y1="0"
              :x2="dayWidth"
              :y2="chartHeight"
              stroke="#e4e7ed"
              stroke-width="1"
            />
          </pattern>
        </defs>
        <rect
          x="0"
          y="0"
          :width="timelineWidth"
          :height="chartHeight"
          :fill="`url(#grid-${activeResourceType})`"
        />

        <!-- Y轴线 -->
        <line
          x1="40"
          y1="10"
          x2="40"
          :y2="chartHeight - 30"
          stroke="#dcdfe6"
          stroke-width="1"
        />

        <!-- X轴线 -->
        <line
          x1="40"
          :y1="chartHeight - 30"
          :x2="timelineWidth"
          :y2="chartHeight - 30"
          stroke="#dcdfe6"
          stroke-width="1"
        />

        <!-- 每日的资源柱状图 -->
        <g v-for="(day, index) in dailyDemandData" :key="day.date">
          <!-- 柱状图 -->
          <rect
            :x="index * dayWidth + 45"
            :y="chartHeight - 30 - getBarHeight(day.demand)"
            :width="Math.max(dayWidth - 10, 20)"
            :height="getBarHeight(day.demand)"
            :fill="getBarColor(day.demand, maxDemand)"
            :class="{ 'is-weekend': day.isWeekend }"
            rx="2"
          />

          <!-- 数量标签 -->
          <text
            v-if="day.demand > 0"
            :x="index * dayWidth + dayWidth / 2 + 42"
            :y="chartHeight - 35 - getBarHeight(day.demand)"
            text-anchor="middle"
            font-size="11"
            fill="#606266"
          >
            {{ day.demand.toFixed(1) }}
          </text>

          <!-- 日期标签 -->
          <text
            :x="index * dayWidth + dayWidth / 2 + 42"
            :y="chartHeight - 10"
            text-anchor="middle"
            font-size="10"
            fill="#909399"
          >
            {{ formatDayDate(day.date) }}
          </text>
        </g>

        <!-- Y轴刻度和标签 -->
        <g v-for="tick in yAxisTicks" :key="tick.value">
          <text
            x="35"
            :y="chartHeight - 30 - tick.position + 4"
            text-anchor="end"
            font-size="10"
            fill="#909399"
          >
            {{ tick.label }}
          </text>
          <line
            x1="35"
            :y1="chartHeight - 30 - tick.position"
            x2="40"
            :y2="chartHeight - 30 - tick.position"
            stroke="#dcdfe6"
            stroke-width="1"
          />
        </g>
      </svg>
    </div>

    <!-- 图例和说明 -->
    <div class="resource-legend">
      <div class="legend-section">
        <div class="legend-title">需求等级</div>
        <div class="legend-items">
          <div class="legend-item">
            <span class="legend-color low-demand"></span>
            <span>低需求 (≤{{ (maxDemand * 0.5).toFixed(1) }})</span>
          </div>
          <div class="legend-item">
            <span class="legend-color medium-demand"></span>
            <span>中等需求 ({{ (maxDemand * 0.5).toFixed(1) }}-{{ (maxDemand * 0.8).toFixed(1) }})</span>
          </div>
          <div class="legend-item">
            <span class="legend-color high-demand"></span>
            <span>高需求 (>{{ (maxDemand * 0.8).toFixed(1) }})</span>
          </div>
        </div>
      </div>
      <div class="legend-section">
        <div class="legend-title">统计信息</div>
        <div class="legend-items">
          <div class="legend-item">
            <span class="stat-label">峰值需求:</span>
            <span class="stat-value">{{ maxDemand.toFixed(1) }} {{ unit }}</span>
          </div>
          <div class="legend-item">
            <span class="stat-label">平均需求:</span>
            <span class="stat-value">{{ averageDemand.toFixed(1) }} {{ unit }}</span>
          </div>
          <div class="legend-item">
            <span class="stat-label">总需求:</span>
            <span class="stat-value">{{ totalDemand.toFixed(1) }} {{ unit }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
  <div v-else class="resource-view-empty">
    <el-empty description="暂无任务数据" />
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { TrendCharts } from '@element-plus/icons-vue'
import { formatDate } from '@/utils/dateFormat'

const props = defineProps({
  projectId: {
    type: [Number, String],
    required: true
  },
  tasks: {
    type: Array,
    default: () => []
  },
  timelineDays: {
    type: Array,
    default: () => []
  },
  dayWidth: {
    type: Number,
    default: 40
  }
})

const activeResourceType = ref('labor')
const chartHeight = 180
const timelineContainer = ref(null)

// Debug: Log component lifecycle
onMounted(() => {
  console.log('ResourceView mounted!')
  console.log('props.projectId:', props.projectId)
  console.log('props.tasks:', props.tasks)
  console.log('props.timelineDays:', props.timelineDays)
  console.log('props.dayWidth:', props.dayWidth)
})

// 资源单位映射
const resourceUnits = {
  labor: '人/d',
  equipment: '台/d',
  material: 'kg'
}

const unit = computed(() => resourceUnits[activeResourceType.value] || '')

// 计算每日资源需求
const dailyDemandData = computed(() => {
  if (props.timelineDays.length === 0) return []

  const demandMap = new Map()
  const weekendDays = [0, 6] // 周日、周六

  // 初始化所有日期
  props.timelineDays.forEach(day => {
    const date = new Date(day.date)
    demandMap.set(day.date, {
      date: day.date,
      demand: 0,
      isWeekend: weekendDays.includes(date.getDay())
    })
  })

  // 计算每个任务的资源需求
  props.tasks.forEach(task => {
    // 检查任务是否有资源数据
    if (!task.resources || !Array.isArray(task.resources) || task.resources.length === 0) {
      return
    }

    const startDate = new Date(task.start)
    const endDate = new Date(task.end)

    for (let d = new Date(startDate); d <= endDate; d.setDate(d.getDate() + 1)) {
      const dateKey = formatDate(d)
      const dayData = demandMap.get(dateKey)

      if (dayData) {
        task.resources.forEach(resource => {
          if (resource.type === activeResourceType.value) {
            dayData.demand += resource.quantity || 0
          }
        })
      }
    }
  })

  return Array.from(demandMap.values())
})

// 最大需求量
const maxDemand = computed(() => {
  const demands = dailyDemandData.value.map(d => d.demand)
  return Math.max(...demands, 10) // 最小值为10
})

// 平均需求
const averageDemand = computed(() => {
  const demands = dailyDemandData.value.map(d => d.demand)
  const sum = demands.reduce((a, b) => a + b, 0)
  return demands.length > 0 ? sum / demands.length : 0
})

// 总需求
const totalDemand = computed(() => {
  return dailyDemandData.value.reduce((sum, day) => sum + day.demand, 0)
})

// 时间轴宽度
const timelineWidth = computed(() => {
  return (props.timelineDays.length * props.dayWidth) + 80
})

// Y轴刻度
const yAxisTicks = computed(() => {
  const ticks = []
  const tickCount = 5
  const step = maxDemand.value / tickCount

  for (let i = 0; i <= tickCount; i++) {
    const value = step * i
    const position = (value / maxDemand.value) * (chartHeight - 40)
    ticks.push({
      value,
      position,
      label: value.toFixed(0)
    })
  }

  return ticks
})

// 获取柱状图高度
const getBarHeight = (demand) => {
  const availableHeight = chartHeight - 40
  return (demand / maxDemand.value) * availableHeight
}

// 获取柱状图颜色
const getBarColor = (demand, max) => {
  const ratio = demand / max
  if (ratio <= 0.5) return '#67c23a' // 绿色
  if (ratio <= 0.8) return '#e6a23c' // 橙色
  return '#f56c6c' // 红色
}

// 格式化日期显示
const formatDayDate = (dateStr) => {
  const date = new Date(dateStr)
  return `${date.getMonth() + 1}/${date.getDate()}`
}

// 资源类型变化
const handleResourceTypeChange = () => {
  // 可以在这里添加额外的处理逻辑
}
</script>

<style scoped>
.resource-view {
  border-top: 1px solid #dcdfe6;
  background: #f5f7fa;
  padding: 16px;
}

.resource-view-empty {
  border-top: 1px solid #dcdfe6;
  background: #f5f7fa;
  padding: 32px;
  text-align: center;
}

.resource-view-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.view-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 600;
  color: #303133;
}

.view-controls {
  display: flex;
  gap: 12px;
}

.resource-timeline-container {
  margin-top: 12px;
  overflow-x: auto;
  background: #fff;
  border-radius: 4px;
  padding: 8px;
}

.resource-timeline {
  display: block;
  min-width: 100%;
}

.resource-timeline rect {
  transition: fill 0.2s, height 0.3s;
}

.resource-timeline rect:hover {
  opacity: 0.8;
}

.resource-timeline rect.is-weekend {
  fill-opacity: 0.5;
}

.resource-legend {
  display: flex;
  gap: 24px;
  margin-top: 16px;
  flex-wrap: wrap;
}

.legend-section {
  flex: 1;
  min-width: 200px;
}

.legend-title {
  font-size: 12px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 8px;
}

.legend-items {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #606266;
}

.legend-color {
  width: 16px;
  height: 12px;
  border-radius: 2px;
}

.low-demand {
  background: #67c23a;
}

.medium-demand {
  background: #e6a23c;
}

.high-demand {
  background: #f56c6c;
}

.stat-label {
  color: #909399;
}

.stat-value {
  font-weight: 600;
  color: #303133;
}
</style>
