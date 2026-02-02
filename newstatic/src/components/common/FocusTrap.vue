<template>
  <teleport to="body">
    <transition name="fade">
      <div
        v-if="active"
        ref="trapRef"
        class="focus-trap"
        :aria-hidden="true"
        tabindex="-1"
      ></div>
    </transition>
  </teleport>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue'

interface Props {
  active: boolean
  onEscape?: () => void
}

const props = defineProps<Props>()

const trapRef = ref<HTMLDivElement>()
let cleanup: (() => void) | null = null

/**
 * 创建焦点 trap
 */
function createFocusTrap() {
  if (!trapRef.value) return

  // 获取所有可聚焦元素
  const focusableElements = document.querySelectorAll(
    'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
  ) as NodeListOf<HTMLElement>

  if (focusableElements.length === 0) return

  const firstFocusable = focusableElements[0]
  const lastFocusable = focusableElements[focusableElements.length - 1]

  /**
   * 处理 Tab 键
   */
  function handleTabKey(e: KeyboardEvent) {
    if (e.key !== 'Tab') return

    if (e.shiftKey) {
      if (document.activeElement === firstFocusable) {
        e.preventDefault()
        lastFocusable.focus()
      }
    } else {
      if (document.activeElement === lastFocusable) {
        e.preventDefault()
        firstFocusable.focus()
      }
    }
  }

  /**
   * 处理 Escape 键
   */
  function handleEscapeKey(e: KeyboardEvent) {
    if (e.key === 'Escape' && props.onEscape) {
      props.onEscape()
    }
  }

  // 添加事件监听
  document.addEventListener('keydown', handleTabKey)
  document.addEventListener('keydown', handleEscapeKey)

  // 焦点到第一个元素
  setTimeout(() => {
    firstFocusable?.focus()
  }, 100)

  // 返回清理函数
  cleanup = () => {
    document.removeEventListener('keydown', handleTabKey)
    document.removeEventListener('keydown', handleEscapeKey)
  }
}

/**
 * 监听激活状态
 */
watch(() => props.active, (isActive) => {
  if (isActive) {
    createFocusTrap()
  } else {
    if (cleanup) {
      cleanup()
      cleanup = null
    }
  }
})

onUnmounted(() => {
  if (cleanup) {
    cleanup()
  }
})
</script>

<style scoped>
.focus-trap {
  position: fixed;
  top: 0;
  left: 0;
  width: 1px;
  height: 1px;
  opacity: 0;
  pointer-events: none;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
