<template>
  <div
    ref="timelineRef"
    class="virtual-timeline"
    :style="containerStyle"
    @scroll.passive="handleScroll"
  >
    <!-- Timeline container -->
    <div
      class="virtual-timeline__content"
      :style="contentStyle"
    >
      <!-- Timeline grid header -->
      <div
        class="virtual-timeline__header"
        :style="headerStyle"
      >
        <slot
          name="header"
          :days="visibleDays"
          :week-width="weekWidth"
          :day-width="dayWidth"
        >
          <TimelineGrid
            :days="visibleDays"
            :day-width="dayWidth"
            :offset-x="offsetX"
          />
        </slot>
      </div>

      <!-- Timeline rows (tasks) -->
      <div
        class="virtual-timeline__rows"
        :style="rowsStyle"
      >
        <RecycleScroller
          ref="scrollerRef"
          class="virtual-timeline__scroller"
          :items="tasks"
          :item-size="rowHeight"
          :buffer="buffer"
          key-field="id"
          :emit-update="true"
          @update="handleScrollerUpdate"
          @resize="handleScrollerResize"
          v-slot="{ item: task, index }"
        >
          <div
            class="virtual-timeline__row"
            :style="{
              height: `${rowHeight}px`,
              width: `${totalWidth}px`
            }"
          >
            <slot
              name="row"
              :task="task"
              :index="index"
              :days="visibleDays"
              :day-width="dayWidth"
            >
              <!-- Task bar will be rendered here -->
              <TaskBar
                :task="task"
                :day-width="dayWidth"
                :days="visibleDays"
              />
            </slot>
          </div>
        </RecycleScroller>
      </div>

      <!-- Timeline background -->
      <div
        class="virtual-timeline__background"
        :style="backgroundStyle"
      >
        <slot
          name="background"
          :days="visibleDays"
          :day-width="dayWidth"
        >
          <TimelineBackground
            :days="visibleDays"
            :day-width="dayWidth"
            :height="totalHeight"
          />
        </slot>
      </div>

      <!-- Today marker -->
      <div
        v-if="todayPosition !== null"
        class="virtual-timeline__today-marker"
        :style="todayMarkerStyle"
      />

      <!-- Dependency lines overlay -->
      <svg
        class="virtual-timeline__dependencies"
        :style="dependenciesStyle"
      >
        <slot
          name="dependencies"
          :tasks="visibleTasks"
          :day-width="dayWidth"
          :row-height="rowHeight"
        >
          <DependencyLines
            :tasks="visibleTasks"
            :day-width="dayWidth"
            :row-height="rowHeight"
          />
        </slot>
      </svg>
    </div>

    <!-- Loading indicator -->
    <div
      v-if="loading"
      class="virtual-timeline__loading"
    >
      <ElIcon class="is-loading">
        <Loading />
      </ElIcon>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { RecycleScroller } from 'vue-virtual-scroller'
import 'vue-virtual-scroller/dist/vue-virtual-scroller.css'
import { Loading } from '@element-plus/icons-vue'
import { ElIcon } from 'element-plus'
import TimelineGrid from './TimelineGrid.vue'
import TimelineBackground from './TimelineBackground.vue'
import TaskBar from './TaskBar.vue'
import DependencyLines from './DependencyLines.vue'

/**
 * Virtual Scrolling Timeline Component
 *
 * High-performance timeline with virtual scrolling support.
 * Uses vue-virtual-scroller's RecycleScroller for efficient rendering.
 *
 * Features:
 * - Virtual scrolling with dynamic row height
 * - Buffer configuration (500px default)
 * - Performance optimized rendering
 * - Sync scrolling with task list
 *
 * @props
 * @param {Array} tasks - Task list to display
 * @param {Array} days - Timeline days
 * @param {number} rowHeight - Height of each row (default: 60)
 * @param {number} dayWidth - Width of each day (default: 40)
 * @param {number} buffer - Buffer size in pixels (default: 500)
 * @param {number} containerHeight - Container height in pixels
 * @param {boolean} showToday - Show today marker (default: true)
 * @param {boolean} loading - Loading state (default: false)
 */

const props = withDefaults(defineProps({
  tasks: {
    type: Array,
    default: () => []
  },
  days: {
    type: Array,
    default: () => []
  },
  rowHeight: {
    type: Number,
    default: 60
  },
  dayWidth: {
    type: Number,
    default: 40
  },
  buffer: {
    type: Number,
    default: 500
  },
  containerHeight: {
    type: Number,
    default: 600
  },
  showToday: {
    type: Boolean,
    default: true
  },
  loading: {
    type: Boolean,
    default: false
  }
}), {
  tasks: () => [],
  days: () => [],
  rowHeight: 60,
  dayWidth: 40,
  buffer: 500,
  containerHeight: 600,
  showToday: true,
  loading: false
})

const emit = defineEmits([
  'scroll',
  'resize',
  'row-mouseenter',
  'row-mouseleave',
  'row-click'
])

// Refs
const timelineRef = ref(null)
const scrollerRef = ref(null)

// State
const scrollTop = ref(0)
const scrollLeft = ref(0)
const viewportWidth = ref(0)
const viewportHeight = ref(props.containerHeight)

// ==================== Computed Properties ====================

/**
 * Container style
 */
const containerStyle = computed(() => ({
  height: `${props.containerHeight}px`,
  position: 'relative',
  overflow: 'auto',
  '-webkit-overflow-scrolling': 'touch'
}))

/**
 * Content style
 */
const contentStyle = computed(() => ({
  width: `${totalWidth.value}px`,
  height: `${totalHeight.value}px`,
  position: 'relative',
  minHeight: '100%'
}))

/**
 * Total width
 */
const totalWidth = computed(() => {
  return props.days.length * props.dayWidth
})

/**
 * Total height
 */
const totalHeight = computed(() => {
  return props.tasks.length * props.rowHeight
})

/**
 * Header style
 */
const headerStyle = computed(() => ({
  position: 'sticky',
  top: 0,
  left: 0,
  right: 0,
  height: '50px',
  backgroundColor: '#fff',
  zIndex: 10,
  transform: `translateX(-${scrollLeft.value}px)`
}))

/**
 * Rows style
 */
const rowsStyle = computed(() => ({
  position: 'absolute',
  top: '50px',
  left: 0,
  right: 0,
  bottom: 0
}))

/**
 * Background style
 */
const backgroundStyle = computed(() => ({
  position: 'absolute',
  top: '50px',
  left: 0,
  width: `${totalWidth.value}px`,
  height: `${totalHeight.value}px`,
  transform: `translateX(-${scrollLeft.value}px)`,
  pointerEvents: 'none',
  zIndex: 0
}))

/**
 * Dependencies SVG style
 */
const dependenciesStyle = computed(() => ({
  position: 'absolute',
  top: '50px',
  left: 0,
  width: `${totalWidth.value}px`,
  height: `${totalHeight.value}px`,
  transform: `translate(-${scrollLeft.value}px, 0)`,
  pointerEvents: 'none',
  zIndex: 5
}))

/**
 * Calculate visible days based on scroll position
 */
const visibleDays = computed(() => {
  const startDay = Math.floor(scrollLeft.value / props.dayWidth)
  const endDay = Math.ceil((scrollLeft.value + viewportWidth.value) / props.dayWidth)

  return props.days.slice(
    Math.max(0, startDay - 7), // 7 days buffer
    Math.min(props.days.length, endDay + 7)
  )
})

/**
 * Horizontal offset for days
 */
const offsetX = computed(() => {
  const startDay = Math.floor(scrollLeft.value / props.dayWidth)
  return (Math.max(0, startDay - 7)) * props.dayWidth
})

/**
 * Week width (7 days)
 */
const weekWidth = computed(() => {
  return props.dayWidth * 7
})

/**
 * Calculate today's position
 */
const todayPosition = computed(() => {
  if (!props.showToday || props.days.length === 0) {
    return null
  }

  const today = new Date()
  const todayStr = today.toISOString().split('T')[0]

  const todayIndex = props.days.findIndex(day => {
    if (typeof day === 'string') {
      return day === todayStr
    }
    return day.date === todayStr
  })

  if (todayIndex === -1) {
    return null
  }

  return todayIndex * props.dayWidth + props.dayWidth / 2
})

/**
 * Today marker style
 */
const todayMarkerStyle = computed(() => {
  const position = todayPosition.value
  if (position === null) {
    return {}
  }

  return {
    position: 'absolute',
    top: '50px',
    left: `${position}px`,
    width: '2px',
    height: `${totalHeight.value}px`,
    backgroundColor: '#ff6b6b',
    zIndex: 6,
    pointerEvents: 'none'
  }
})

/**
 * Get visible tasks (for dependency lines)
 */
const visibleTasks = computed(() => {
  // RecycleScroller handles virtualization internally
  // Return all tasks as dependencies need full context
  return props.tasks
})

// ==================== Methods ====================

/**
 * Handle scroll event
 */
const handleScroll = (event) => {
  const target = event.target
  scrollTop.value = target.scrollTop
  scrollLeft.value = target.scrollLeft

  // Update viewport dimensions
  viewportWidth.value = target.clientWidth
  viewportHeight.value = target.clientHeight

  emit('scroll', {
    scrollTop: scrollTop.value,
    scrollLeft: scrollLeft.value,
    viewportWidth: viewportWidth.value,
    viewportHeight: viewportHeight.value
  })
}

/**
 * Handle scroller update
 */
const handleScrollerUpdate = () => {
  // Force update on scroller change
  nextTick(() => {
    if (timelineRef.value) {
      viewportWidth.value = timelineRef.value.clientWidth
      viewportHeight.value = timelineRef.value.clientHeight
    }
  })
}

/**
 * Handle scroller resize
 */
const handleScrollerResize = () => {
  emit('resize', {
    width: viewportWidth.value,
    height: viewportHeight.value
  })
}

/**
 * Scroll to task
 * @param {number} taskIndex - Task index
 * @param {string} alignment - Alignment (start, center, end, auto)
 */
const scrollToTask = (taskIndex, alignment = 'auto') => {
  if (!scrollerRef.value || taskIndex < 0 || taskIndex >= props.tasks.length) {
    return
  }

  scrollerRef.value.scrollToItem(taskIndex, alignment)
}

/**
 * Scroll to position
 * @param {number} scrollTop - Scroll top position
 * @param {number} scrollLeft - Scroll left position
 */
const scrollToPosition = ({ scrollTop: top, scrollLeft: left } = {}) => {
  if (!timelineRef.value) {
    return
  }

  if (top !== undefined) {
    timelineRef.value.scrollTop = top
  }

  if (left !== undefined) {
    timelineRef.value.scrollLeft = left
  }
}

/**
 * Scroll to today
 */
const scrollToToday = () => {
  const position = todayPosition.value
  if (position !== null && timelineRef.value) {
    const centerPosition = position - viewportWidth.value / 2
    timelineRef.value.scrollLeft = Math.max(0, centerPosition)
  }
}

/**
 * Refresh layout
 */
const refresh = () => {
  if (scrollerRef.value) {
    scrollerRef.value.forceUpdate()
  }
}

/**
 * Get scroll info
 */
const getScrollInfo = () => {
  return {
    scrollTop: scrollTop.value,
    scrollLeft: scrollLeft.value,
    viewportWidth: viewportWidth.value,
    viewportHeight: viewportHeight.value,
    totalWidth: totalWidth.value,
    totalHeight: totalHeight.value
  }
}

// ==================== Lifecycle ====================

onMounted(() => {
  if (timelineRef.value) {
    viewportWidth.value = timelineRef.value.clientWidth
    viewportHeight.value = timelineRef.value.clientHeight
  }
})

// ==================== Watchers ====================

// Watch for tasks changes
watch(() => props.tasks, () => {
  refresh()
}, { deep: true })

// Watch for days changes
watch(() => props.days, () => {
  refresh()
}, { deep: true })

// Expose methods
defineExpose({
  scrollToTask,
  scrollToPosition,
  scrollToToday,
  refresh,
  getScrollInfo
})
</script>

<style scoped lang="scss">
.virtual-timeline {
  position: relative;
  overflow: auto;
  -webkit-overflow-scrolling: touch;
  user-select: none;

  &__content {
    position: relative;
    min-width: 100%;
    min-height: 100%;
  }

  &__header {
    position: sticky;
    top: 0;
    z-index: 10;
    background-color: #fff;
    will-change: transform;
  }

  &__rows {
    position: relative;
  }

  &__scroller {
    height: 100%;

    :deep(.vue-recycle-scroller__item-wrapper) {
      position: relative;
    }
  }

  &__row {
    position: relative;
    display: flex;
    align-items: center;
    border-bottom: 1px solid #f0f0f0;
    transition: background-color 0.2s;

    &:hover {
      background-color: #f5f7fa;
    }
  }

  &__background {
    position: absolute;
    pointer-events: none;
    will-change: transform;
  }

  &__today-marker {
    &::after {
      content: '今天';
      position: absolute;
      top: -20px;
      left: 50%;
      transform: translateX(-50%);
      font-size: 12px;
      color: #ff6b6b;
      white-space: nowrap;
    }
  }

  &__dependencies {
    position: absolute;
    overflow: visible;
    will-change: transform;

    :deep(svg) {
      display: block;
    }
  }

  &__loading {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    z-index: 100;
    font-size: 24px;
    color: var(--el-color-primary);
  }
}
</style>
