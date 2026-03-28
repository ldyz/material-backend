<template>
  <div
    ref="containerRef"
    class="virtual-timeline"
    :style="{ height: containerHeight + 'px' }"
    @scroll="handleScroll"
  >
    <!-- 时间轴容器 -->
    <div
      class="virtual-timeline__content"
      :style="{
        width: totalWidth + 'px',
        height: totalHeight + 'px'
      }"
    >
      <!-- 垂直方向：可见任务行 -->
      <div
        class="virtual-timeline__rows"
        :style="{ transform: `translateY(${offsetY}px)` }"
      >
        <slot
          name="row"
          v-for="(task, index) in visibleTasks"
          :key="task.id"
          :task="task"
          :index="startIndex + index"
          :style="{ height: rowHeight + 'px' }"
        ></slot>
      </div>

      <!-- 水平方向：时间网格 -->
      <div
        class="virtual-timeline__grid"
        :style="{
          transform: `translateX(${offsetX}px)`,
          height: totalHeight + 'px'
        }"
      >
        <slot
          name="grid"
          :days="visibleDays"
          :day-width="dayWidth"
        ></slot>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, PropType } from 'vue'
import { useVirtualScroll } from '@/composables/useVirtualScroll'

/**
 * 虚拟时间轴组件
 * 同时支持垂直和水平方向的虚拟滚动
 */
interface Props {
  // 任务列表
  tasks: any[]
  // 时间轴日期
  days: any[]
  // 行高度
  rowHeight: number
  // 日宽度
  dayWidth: number
  // 容器高度
  containerHeight: number
  // 额外渲染行数
  rowOverscan?: number
  // 额外渲染天数
  dayOverscan?: number
}

const props = withDefaults(defineProps<Props>(), {
  rowOverscan: 3,
  dayOverscan: 7
})

const emit = defineEmits<{
  scroll: [event: Event]
}>()

const containerRef = ref<HTMLElement>()

// 垂直方向虚拟滚动
const {
  visibleItems: visibleTasks,
  totalHeight,
  offsetY,
  scrollTop
} = useVirtualScroll({
  items: ref(props.tasks),
  itemHeight: props.rowHeight,
  containerHeight: props.containerHeight,
  containerRef,
  overscan: props.rowOverscan
})

// 水平方向状态
const scrollLeft = ref(0)
const viewportWidth = ref(0)

// 计算总宽度
const totalWidth = computed(() => {
  return props.days.length * props.dayWidth
})

// 计算可见天数
const visibleDays = computed(() => {
  const startDay = Math.floor(scrollLeft.value / props.dayWidth) - props.dayOverscan
  const endDay = Math.min(
    props.days.length,
    Math.ceil((scrollLeft.value + viewportWidth.value) / props.dayWidth) + props.dayOverscan
  )

  return props.days.slice(
    Math.max(0, startDay),
    endDay
  )
})

// 水平方向偏移
const offsetX = computed(() => {
  const startDay = Math.floor(scrollLeft.value / props.dayWidth) - props.dayOverscan
  return Math.max(0, startDay) * props.dayWidth
})

// 起始索引
const startIndex = computed(() => {
  return Math.floor(scrollTop.value / props.rowHeight) - props.rowOverscan
})

/**
 * 处理滚动事件
 */
function handleScroll(event: Event) {
  const target = event.target as HTMLElement
  scrollTop.value = target.scrollTop
  scrollLeft.value = target.scrollLeft
  viewportWidth.value = target.clientWidth

  emit('scroll', event)
}

/**
 * 监听数据变化
 */
watch(() => props.tasks, () => {
  // 垂直滚动会自动更新
}, { deep: true })

watch(() => props.days, () => {
  // 水平滚动会自动更新
}, { deep: true })
</script>

<style scoped>
.virtual-timeline {
  position: relative;
  overflow: auto;
  -webkit-overflow-scrolling: touch;
}

.virtual-timeline__content {
  position: relative;
}

.virtual-timeline__rows {
  position: absolute;
  left: 0;
  right: 0;
  top: 0;
  will-change: transform;
}

.virtual-timeline__grid {
  position: absolute;
  left: 0;
  top: 0;
  pointer-events: none;
  will-change: transform;
}
</style>
