<template>
  <a
    :href="`#${targetId}`"
    class="skip-link"
    @click="handleClick"
  >
    <slot>{{ label }}</slot>
  </a>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface Props {
  targetId: string
  label?: string
}

const props = withDefaults(defineProps<Props>(), {
  label: '跳到主内容'
})

const targetElement = ref<HTMLElement | null>(null)

/**
 * 处理点击事件
 */
function handleClick(event: MouseEvent) {
  event.preventDefault()

  // 查找目标元素
  const target = document.getElementById(props.targetId)
  if (!target) return

  // 设置焦点
  target.focus()

  // 滚动到视图
  target.scrollIntoView({ behavior: 'smooth', block: 'start' })

  // 确保元素可获得焦点
  if (target.tabIndex === -1) {
    target.tabIndex = -1
  }
}

/**
 * 初始化
 * 确保目标元素可以获得焦点
 */
onMounted(() => {
  const target = document.getElementById(props.targetId)
  if (target && target.tabIndex < 0) {
    target.tabIndex = -1
  }
})
</script>

<style scoped>
.skip-link {
  position: fixed;
  top: -40px;
  left: 0;
  padding: 8px 16px;
  background: var(--color-primary);
  color: #fff;
  text-decoration: none;
  z-index: 9999;
  transition: top var(--transition-fast);
  font-weight: var(--font-weight-medium);
}

.skip-link:focus {
  top: 0;
}

.skip-link:hover {
  background: var(--color-primary-dark);
  text-decoration: underline;
}

/* 高对比度模式 */
@media (prefers-contrast: high) {
  .skip-link {
    border: 2px solid currentColor;
  }
}
</style>
