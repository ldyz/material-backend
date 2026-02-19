<template>
  <div
    class="stat-card"
    :class="{ 'is-clickable': clickable }"
    :style="{ borderTopColor: color }"
    @click="handleClick"
  >
    <!-- Card Header -->
    <div class="stat-card__header">
      <div class="stat-icon" :style="{ backgroundColor: `${color}20` }">
        <el-icon :size="24" :color="color">
          <component :is="iconComponent" />
        </el-icon>
      </div>
      <div class="stat-actions" v-if="$slots.actions">
        <slot name="actions"></slot>
      </div>
    </div>

    <!-- Card Body -->
    <div class="stat-card__body">
      <div class="stat-value">{{ value }}</div>
      <div class="stat-title">{{ title }}</div>
    </div>

    <!-- Trend Indicator -->
    <div v-if="trend" class="stat-card__trend">
      <el-icon
        :size="14"
        :color="trendColor"
      >
        <component :is="trendIcon" />
      </el-icon>
      <span :style="{ color: trendColor }">{{ trendValue }}</span>
      <span class="trend-label">vs last period</span>
    </div>

    <!-- Sparkline -->
    <div v-if="sparklineData && sparklineData.length > 0" class="stat-card__sparkline">
      <canvas ref="sparklineCanvas" :height="40"></canvas>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import * as ElementPlusIcons from '@element-plus/icons-vue'
import { Chart } from 'chart.js/auto'

/**
 * StatCard Component
 * KPI summary card with trend indicator and sparkline
 *
 * @props {String} icon - Icon name from Element Plus icons
 * @props {String} title - Card title
 * @props {String|Number} value - Main value to display
 * @props {String} trend - Trend direction (up/down/neutral)
 * @props {String} trendValue - Trend percentage or value
 * @props {String} color - Accent color
 * @props {Array} sparklineData - Array of data points for sparkline
 * @props {Boolean} clickable - Whether card is clickable
 *
 * @emits {Object} click - Emitted when card is clicked
 */

const props = defineProps({
  icon: {
    type: String,
    required: true
  },
  title: {
    type: String,
    required: true
  },
  value: {
    type: [String, Number],
    required: true
  },
  trend: {
    type: String,
    default: 'neutral',
    validator: (value) => ['up', 'down', 'neutral'].includes(value)
  },
  trendValue: {
    type: String,
    default: ''
  },
  color: {
    type: String,
    default: '#409EFF'
  },
  sparklineData: {
    type: Array,
    default: () => []
  },
  clickable: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['click'])

// State
const sparklineCanvas = ref(null)
let sparklineChart = null

// Computed
const iconComponent = computed(() => {
  return ElementPlusIcons[props.icon] || ElementPlusIcons['Document']
})

const trendIcon = computed(() => {
  if (props.trend === 'up') return 'ArrowUp'
  if (props.trend === 'down') return 'ArrowDown'
  return 'Minus'
})

const trendColor = computed(() => {
  if (props.trend === 'up') return '#67C23A'
  if (props.trend === 'down') return '#F56C6C'
  return '#909399'
})

// Methods
const handleClick = () => {
  if (props.clickable) {
    emit('click')
  }
}

const initSparkline = () => {
  if (!sparklineCanvas.value || !props.sparklineData.length) return

  const ctx = sparklineCanvas.value.getContext('2d')

  if (sparklineChart) {
    sparklineChart.destroy()
  }

  sparklineChart = new Chart(ctx, {
    type: 'line',
    data: {
      labels: props.sparklineData.map((_, i) => i),
      datasets: [{
        data: props.sparklineData.map(d => d.value),
        borderColor: props.color,
        backgroundColor: `${props.color}20`,
        borderWidth: 2,
        tension: 0.4,
        fill: true,
        pointRadius: 0,
        pointHoverRadius: 4
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          display: false
        },
        tooltip: {
          enabled: false
        }
      },
      scales: {
        x: {
          display: false
        },
        y: {
          display: false
        }
      }
    }
  })
}

// Lifecycle
onMounted(() => {
  initSparkline()
})

watch(() => props.sparklineData, () => {
  initSparkline()
}, { deep: true })
</script>

<style scoped lang="scss">
.stat-card {
  background: #fff;
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  border-top: 3px solid;
  transition: all 0.3s;

  &.is-clickable {
    cursor: pointer;

    &:hover {
      transform: translateY(-4px);
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    }
  }

  &__header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 12px;

    .stat-icon {
      width: 48px;
      height: 48px;
      border-radius: 12px;
      display: flex;
      align-items: center;
      justify-content: center;
    }

    .stat-actions {
      opacity: 0;
      transition: opacity 0.2s;
    }
  }

  &:hover &__header .stat-actions {
    opacity: 1;
  }

  &__body {
    margin-bottom: 12px;

    .stat-value {
      font-size: 28px;
      font-weight: 700;
      color: #303133;
      line-height: 1;
      margin-bottom: 4px;
    }

    .stat-title {
      font-size: 13px;
      color: #909399;
      font-weight: 500;
    }
  }

  &__trend {
    display: flex;
    align-items: center;
    gap: 4px;
    font-size: 12px;
    margin-bottom: 8px;

    .trend-label {
      color: #909399;
      margin-left: 4px;
    }
  }

  &__sparkline {
    margin-top: 8px;
    height: 40px;

    canvas {
      width: 100% !important;
      height: 100% !important;
    }
  }
}

/* Responsive Design */
@media (max-width: 768px) {
  .stat-card {
    padding: 12px;

    &__body .stat-value {
      font-size: 24px;
    }
  }
}
</style>
