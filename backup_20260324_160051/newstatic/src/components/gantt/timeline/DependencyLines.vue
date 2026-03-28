<template>
  <g class="dependency-lines">
    <path
      v-for="dep in renderableDependencies"
      :key="dep.key"
      :d="dep.path"
      :stroke="dep.stroke"
      :stroke-width="dep.strokeWidth"
      fill="none"
      :marker-end="`url(#${arrowMarkerId})`"
      class="dependency-line"
      :class="{ 'is-critical': dep.isCritical }"
    />
  </g>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  tasks: any[]
  dependencies: any[]
  arrowMarkerId: string
  arrowColor: string
  showCriticalPath: boolean
}

const props = defineProps<Props>()

/**
 * 可渲染的依赖关系
 */
const renderableDependencies = computed(() => {
  return props.dependencies.map((dep, index) => {
    const fromTask = dep.from
    const toTask = dep.to

    // 计算起点和终点
    const startX = fromTask.x + fromTask.width
    const startY = fromTask.y + fromTask.height / 2
    const endX = toTask.x
    const endY = toTask.y + toTask.height / 2

    // 计算路径（折线）
    const midX = (startX + endX) / 2
    const path = `M ${startX} ${startY} L ${midX} ${startY} L ${midX} ${endY} L ${endX} ${endY}`

    return {
      key: `${fromTask.id}-${toTask.id}-${index}`,
      path,
      stroke: dep.isCritical ? '#f56c6c' : props.arrowColor,
      strokeWidth: dep.isCritical ? 2 : 1.5,
      isCritical: dep.isCritical
    }
  })
})
</script>

<style scoped>
.dependency-line {
  transition: stroke-width 0.15s ease;
}

.dependency-line:hover {
  stroke-width: 3 !important;
}

.dependency-line.is-critical {
  stroke: #f56c6c;
  stroke-width: 2;
}
</style>
